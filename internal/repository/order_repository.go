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
	"gorm.io/gorm/clause"
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

			// Atomically decrement stock only when enough inventory remains.
			result := tx.Model(&productVariantModel{}).
				Where("id = ? AND stock_on_hand >= ?", item.ProductVariantID.String(), item.Qty).
				Update("stock_on_hand", gorm.Expr("stock_on_hand - ?", item.Qty))
			if result.Error != nil {
				return fmt.Errorf("decrement stock: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("%w: variant %s", domain.ErrStockUnavailable, item.ProductVariantID)
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

func (r *orderRepository) FindByOrderNo(ctx context.Context, orderNo string) (*domain.Order, error) {
	var model orderModel
	if err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&model).Error; err != nil {
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

func (r *orderRepository) FindItemsByOrderIDs(ctx context.Context, orderIDs []uuid.UUID) (map[uuid.UUID][]domain.OrderItem, error) {
	result := make(map[uuid.UUID][]domain.OrderItem, len(orderIDs))
	if len(orderIDs) == 0 {
		return result, nil
	}

	ids := make([]string, 0, len(orderIDs))
	for _, id := range orderIDs {
		ids = append(ids, id.String())
	}

	var models []orderItemModel
	if err := r.db.WithContext(ctx).
		Where("order_id IN ?", ids).
		Order("order_id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	for _, m := range models {
		item := toDomainOrderItem(&m)
		result[item.OrderID] = append(result[item.OrderID], *item)
	}

	return result, nil
}

func (r *orderRepository) ListByUser(ctx context.Context, userID uuid.UUID, params domain.CustomerOrderListParams) ([]domain.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&orderModel{}).Where("user_id = ?", userID.String())
	if params.Status != "" {
		q = q.Where("order_status = ?", params.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	var models []orderModel
	if err := q.
		Order("placed_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomainOrder(&m)
	}

	return orders, total, nil
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

func (r *orderRepository) MarkOrderReceived(ctx context.Context, orderID uuid.UUID, changedBy *uuid.UUID, notes *string) (bool, error) {
	applied := false

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current orderModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", orderID.String()).First(&current).Error; err != nil {
			return fmt.Errorf("find order: %w", err)
		}
		if current.OrderStatus != domain.OrderStatusShipped {
			return nil
		}

		now := time.Now()
		if err := tx.Model(&orderModel{}).
			Where("id = ?", orderID.String()).
			Updates(map[string]interface{}{
				"order_status":   domain.OrderStatusReceived,
				"payment_status": current.PaymentStatus,
				"updated_at":     now,
			}).Error; err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		var changedByStr *string
		if changedBy != nil {
			s := changedBy.String()
			changedByStr = &s
		}

		history := orderStatusHistoryModel{
			ID:        uuid.New().String(),
			OrderID:   orderID.String(),
			OldStatus: &current.OrderStatus,
			NewStatus: domain.OrderStatusReceived,
			ChangedBy: changedByStr,
			ChangedAt: now,
			Notes:     notes,
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert status history: %w", err)
		}

		res := tx.Model(&vendorBalanceModel{}).
			Where("vendor_id = ?", current.VendorID).
			Updates(map[string]interface{}{
				"escrow_balance": gorm.Expr("escrow_balance + ?", current.Subtotal),
				"total_earned":   gorm.Expr("total_earned + ?", current.Subtotal),
				"updated_at":     now,
			})
		if res.Error != nil {
			return fmt.Errorf("update escrow balance: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return errors.New("failed to apply escrow balance update")
		}

		applied = true
		return nil
	})
	if err != nil {
		return false, err
	}

	return applied, nil
}

func (r *orderRepository) MarkOrderCompleted(ctx context.Context, orderID uuid.UUID, changedBy *uuid.UUID, notes *string) (bool, error) {
	applied := false

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current orderModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", orderID.String()).First(&current).Error; err != nil {
			return fmt.Errorf("find order: %w", err)
		}
		if current.OrderStatus != domain.OrderStatusReceived {
			return nil
		}

		now := time.Now()
		if err := tx.Model(&orderModel{}).
			Where("id = ?", orderID.String()).
			Updates(map[string]interface{}{
				"order_status":   domain.OrderStatusCompleted,
				"payment_status": current.PaymentStatus,
				"updated_at":     now,
			}).Error; err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		var changedByStr *string
		if changedBy != nil {
			s := changedBy.String()
			changedByStr = &s
		}

		history := orderStatusHistoryModel{
			ID:        uuid.New().String(),
			OrderID:   orderID.String(),
			OldStatus: &current.OrderStatus,
			NewStatus: domain.OrderStatusCompleted,
			ChangedBy: changedByStr,
			ChangedAt: now,
			Notes:     notes,
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert status history: %w", err)
		}

		res := tx.Model(&vendorBalanceModel{}).
			Where("vendor_id = ? AND escrow_balance >= ?", current.VendorID, current.Subtotal).
			Updates(map[string]interface{}{
				"escrow_balance":    gorm.Expr("escrow_balance - ?", current.Subtotal),
				"available_balance": gorm.Expr("available_balance + ?", current.Subtotal),
				"updated_at":        now,
			})
		if res.Error != nil {
			return fmt.Errorf("release escrow balance: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return errors.New("failed to release escrow balance")
		}

		// Increment vendor total_sold based on order items.
		var totalQty int64
		if err := tx.Model(&orderItemModel{}).
			Select("COALESCE(SUM(qty), 0)").
			Where("order_id = ?", orderID.String()).
			Scan(&totalQty).Error; err != nil {
			return fmt.Errorf("calculate total sold: %w", err)
		}
		if totalQty > 0 {
			if err := tx.Model(&vendorModel{}).
				Where("id = ?", current.VendorID).
				Updates(map[string]interface{}{
					"total_sold": gorm.Expr("total_sold + ?", totalQty),
					"updated_at": now,
				}).Error; err != nil {
				return fmt.Errorf("increment vendor total sold: %w", err)
			}
		}

		applied = true
		return nil
	})
	if err != nil {
		return false, err
	}

	return applied, nil
}

func (r *orderRepository) ListAutoReceiveCandidates(ctx context.Context, deliveredBefore time.Time, limit int) ([]domain.Order, error) {
	var models []orderModel
	query := r.db.WithContext(ctx).
		Table("orders o").
		Select("o.*").
		Joins("JOIN shipments s ON s.order_id = o.id").
		Where("o.order_status = ?", domain.OrderStatusShipped).
		Where("s.delivered_at IS NOT NULL AND s.delivered_at <= ?", deliveredBefore).
		Order("s.delivered_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomainOrder(&m)
	}
	return orders, nil
}

func (r *orderRepository) ListSettlementCandidates(ctx context.Context, receivedBefore time.Time, limit int) ([]domain.Order, error) {
	var models []orderModel
	query := r.db.WithContext(ctx).
		Model(&orderModel{}).
		Where("order_status = ? AND updated_at <= ?", domain.OrderStatusReceived, receivedBefore).
		Order("updated_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomainOrder(&m)
	}
	return orders, nil
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

func (r *orderRepository) ApplyExpiredWebhookUpdate(ctx context.Context, orderID, invoiceID uuid.UUID, notes *string) (bool, error) {
	applied := false

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current orderModel
		if err := tx.Where("id = ?", orderID.String()).First(&current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("find order: %w", err)
		}

		// No-op for out-of-order webhooks when order already moved to another state.
		if current.OrderStatus != domain.OrderStatusPendingPayment {
			return nil
		}

		invoiceUpdate := tx.Model(&paymentInvoiceModel{}).
			Where("id = ? AND status = ?", invoiceID.String(), domain.InvoiceStatusPending).
			Updates(map[string]interface{}{
				"status":     domain.InvoiceStatusExpired,
				"updated_at": time.Now(),
			})
		if invoiceUpdate.Error != nil {
			return fmt.Errorf("update invoice: %w", invoiceUpdate.Error)
		}
		if invoiceUpdate.RowsAffected == 0 {
			return nil
		}

		if err := tx.Model(&orderModel{}).
			Where("id = ?", orderID.String()).
			Updates(map[string]interface{}{
				"order_status":   domain.OrderStatusCanceled,
				"payment_status": domain.PaymentStatusUnpaid,
				"updated_at":     time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		history := orderStatusHistoryModel{
			ID:        uuid.New().String(),
			OrderID:   orderID.String(),
			OldStatus: &current.OrderStatus,
			NewStatus: domain.OrderStatusCanceled,
			ChangedAt: time.Now(),
			Notes:     notes,
		}
		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("insert status history: %w", err)
		}

		var items []orderItemModel
		if err := tx.Where("order_id = ?", orderID.String()).Find(&items).Error; err != nil {
			return fmt.Errorf("find order items: %w", err)
		}
		for _, item := range items {
			if err := tx.Model(&productVariantModel{}).
				Where("id = ?", item.ProductVariantID).
				Update("stock_on_hand", gorm.Expr("stock_on_hand + ?", item.Qty)).Error; err != nil {
				return fmt.Errorf("restore stock for variant %s: %w", item.ProductVariantID, err)
			}
		}

		applied = true
		return nil
	})
	if err != nil {
		return false, err
	}

	return applied, nil
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
