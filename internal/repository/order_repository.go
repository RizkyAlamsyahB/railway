package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type orderModel struct {
	ID                      string    `gorm:"column:id;primaryKey"`
	OrderNo                 string    `gorm:"column:order_no"`
	UserID                  string    `gorm:"column:user_id"`
	VendorID                string    `gorm:"column:vendor_id"`
	ShippingAddressSnapshot string    `gorm:"column:shipping_address_snapshot;type:jsonb"`
	OrderStatus             string    `gorm:"column:order_status"`
	PaymentStatus           string    `gorm:"column:payment_status"`
	Subtotal                float64   `gorm:"column:subtotal"`
	ShippingFee             float64   `gorm:"column:shipping_fee"`
	PlatformFee             float64   `gorm:"column:platform_fee"`
	GrandTotal              float64   `gorm:"column:grand_total"`
	PlacedAt                time.Time `gorm:"column:placed_at"`
	CreatedAt               time.Time `gorm:"column:created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at"`
}

func (orderModel) TableName() string { return "orders" }

type orderItemModel struct {
	ID                  string  `gorm:"column:id;primaryKey"`
	OrderID             string  `gorm:"column:order_id"`
	ProductVariantID    string  `gorm:"column:product_variant_id"`
	ProductNameSnapshot string  `gorm:"column:product_name_snapshot"`
	SKUSnapshot         string  `gorm:"column:sku_snapshot"`
	Qty                 int     `gorm:"column:qty"`
	UnitPrice           float64 `gorm:"column:unit_price"`
	LineTotal           float64 `gorm:"column:line_total"`
}

func (orderItemModel) TableName() string { return "order_items" }

type orderStatusHistoryModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	OrderID   string    `gorm:"column:order_id"`
	OldStatus *string   `gorm:"column:old_status"`
	NewStatus string    `gorm:"column:new_status"`
	ChangedBy *string   `gorm:"column:changed_by"`
	ChangedAt time.Time `gorm:"column:changed_at"`
	Notes     *string   `gorm:"column:notes"`
}

func (orderStatusHistoryModel) TableName() string { return "order_status_history" }

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new OrderRepository backed by GORM.
func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderWithItems(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert order.
		om := toOrderModel(order)
		if err := tx.Create(&om).Error; err != nil {
			return fmt.Errorf("insert order: %w", err)
		}

		// 2. Insert order items and decrement stock.
		for _, item := range items {
			im := toOrderItemModel(&item)
			if err := tx.Create(&im).Error; err != nil {
				return fmt.Errorf("insert order item: %w", err)
			}

			// Decrement stock_on_hand with optimistic locking.
			result := tx.Model(&productVariantModel{}).
				Where("id = ? AND stock_on_hand >= ?", item.ProductVariantID.String(), item.Qty).
				Update("stock_on_hand", gorm.Expr("stock_on_hand - ?", item.Qty))
			if result.Error != nil {
				return fmt.Errorf("decrement stock: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("insufficient stock for variant %s", item.ProductVariantID)
			}
		}

		// 3. Insert initial status history.
		history := orderStatusHistoryModel{
			ID:        uuid.New().String(),
			OrderID:   order.ID.String(),
			OldStatus: nil,
			NewStatus: order.OrderStatus,
			ChangedAt: time.Now(),
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert status history: %w", err)
		}

		return nil
	})
}

func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var model orderModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainOrder(&model), nil
}

func (r *orderRepository) FindItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	var models []orderItemModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID.String()).Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.OrderItem, len(models))
	for i, m := range models {
		items[i] = *toDomainOrderItem(&m)
	}
	return items, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, orderStatus, paymentStatus string, changedBy *uuid.UUID, notes *string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current order status for history.
		var current orderModel
		if err := tx.Where("id = ?", orderID.String()).First(&current).Error; err != nil {
			return fmt.Errorf("find order: %w", err)
		}

		// Update order status.
		if err := tx.Model(&orderModel{}).
			Where("id = ?", orderID.String()).
			Updates(map[string]interface{}{
				"order_status":   orderStatus,
				"payment_status": paymentStatus,
				"updated_at":     time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		// Insert status history.
		var changedByStr *string
		if changedBy != nil {
			s := changedBy.String()
			changedByStr = &s
		}

		history := orderStatusHistoryModel{
			ID:        uuid.New().String(),
			OrderID:   orderID.String(),
			OldStatus: &current.OrderStatus,
			NewStatus: orderStatus,
			ChangedBy: changedByStr,
			ChangedAt: time.Now(),
			Notes:     notes,
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert status history: %w", err)
		}

		return nil
	})
}

func (r *orderRepository) RestoreStock(ctx context.Context, orderID uuid.UUID) error {
	var items []orderItemModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID.String()).Find(&items).Error; err != nil {
		return fmt.Errorf("find order items: %w", err)
	}
	for _, item := range items {
		if err := r.db.WithContext(ctx).Model(&productVariantModel{}).
			Where("id = ?", item.ProductVariantID).
			Update("stock_on_hand", gorm.Expr("stock_on_hand + ?", item.Qty)).Error; err != nil {
			return fmt.Errorf("restore stock for variant %s: %w", item.ProductVariantID, err)
		}
	}
	return nil
}

// Mapper helpers.

func toOrderModel(o *domain.Order) orderModel {
	snapshot, _ := json.Marshal(o.ShippingAddressSnapshot)
	return orderModel{
		ID:                      o.ID.String(),
		OrderNo:                 o.OrderNo,
		UserID:                  o.UserID.String(),
		VendorID:                o.VendorID.String(),
		ShippingAddressSnapshot: string(snapshot),
		OrderStatus:             o.OrderStatus,
		PaymentStatus:           o.PaymentStatus,
		Subtotal:                o.Subtotal,
		ShippingFee:             o.ShippingFee,
		PlatformFee:             o.PlatformFee,
		GrandTotal:              o.GrandTotal,
		PlacedAt:                o.PlacedAt,
		CreatedAt:               o.CreatedAt,
		UpdatedAt:               o.UpdatedAt,
	}
}

func toDomainOrder(m *orderModel) *domain.Order {
	id, _ := uuid.Parse(m.ID)
	userID, _ := uuid.Parse(m.UserID)
	vendorID, _ := uuid.Parse(m.VendorID)

	var snapshot map[string]interface{}
	_ = json.Unmarshal([]byte(m.ShippingAddressSnapshot), &snapshot)

	return &domain.Order{
		ID:                      id,
		OrderNo:                 m.OrderNo,
		UserID:                  userID,
		VendorID:                vendorID,
		ShippingAddressSnapshot: snapshot,
		OrderStatus:             m.OrderStatus,
		PaymentStatus:           m.PaymentStatus,
		Subtotal:                m.Subtotal,
		ShippingFee:             m.ShippingFee,
		PlatformFee:             m.PlatformFee,
		GrandTotal:              m.GrandTotal,
		PlacedAt:                m.PlacedAt,
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

func toOrderItemModel(i *domain.OrderItem) orderItemModel {
	return orderItemModel{
		ID:                  i.ID.String(),
		OrderID:             i.OrderID.String(),
		ProductVariantID:    i.ProductVariantID.String(),
		ProductNameSnapshot: i.ProductNameSnapshot,
		SKUSnapshot:         i.SKUSnapshot,
		Qty:                 i.Qty,
		UnitPrice:           i.UnitPrice,
		LineTotal:           i.LineTotal,
	}
}

func toDomainOrderItem(m *orderItemModel) *domain.OrderItem {
	id, _ := uuid.Parse(m.ID)
	orderID, _ := uuid.Parse(m.OrderID)
	variantID, _ := uuid.Parse(m.ProductVariantID)

	return &domain.OrderItem{
		ID:                  id,
		OrderID:             orderID,
		ProductVariantID:    variantID,
		ProductNameSnapshot: m.ProductNameSnapshot,
		SKUSnapshot:         m.SKUSnapshot,
		Qty:                 m.Qty,
		UnitPrice:           m.UnitPrice,
		LineTotal:           m.LineTotal,
	}
}
