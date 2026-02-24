package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORM model structs (internal to repository layer).

type cartModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (cartModel) TableName() string { return "carts" }

type cartItemModel struct {
	ID               string    `gorm:"column:id;primaryKey"`
	CartID           string    `gorm:"column:cart_id"`
	ProductVariantID string    `gorm:"column:product_variant_id"`
	Qty              int       `gorm:"column:qty"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (cartItemModel) TableName() string { return "cart_items" }

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new CartRepository backed by GORM.
func NewCartRepository(db *gorm.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	var model cartModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND status = 'active'", userID.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainCart(&model), nil
}

func (r *cartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	model := toCartModel(cart)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *cartRepository) FindItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]domain.CartItem, error) {
	var models []cartItemModel
	if err := r.db.WithContext(ctx).Where("cart_id = ?", cartID.String()).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.CartItem, len(models))
	for i, m := range models {
		items[i] = *toDomainCartItem(&m)
	}
	return items, nil
}

func (r *cartRepository) FindItemByID(ctx context.Context, itemID uuid.UUID) (*domain.CartItem, error) {
	var model cartItemModel
	if err := r.db.WithContext(ctx).Where("id = ?", itemID.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainCartItem(&model), nil
}

func (r *cartRepository) FindItemByCartAndVariant(ctx context.Context, cartID uuid.UUID, variantID uuid.UUID) (*domain.CartItem, error) {
	var model cartItemModel
	if err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_variant_id = ?", cartID.String(), variantID.String()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainCartItem(&model), nil
}

func (r *cartRepository) UpsertItem(ctx context.Context, item *domain.CartItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := toCartItemModel(item)

		// Upsert: insert or increment qty on conflict.
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "cart_id"},
				{Name: "product_variant_id"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"qty": gorm.Expr("cart_items.qty + ?", item.Qty),
			}),
		}).Create(&model).Error; err != nil {
			return err
		}

		// Update cart's updated_at.
		if err := tx.Model(&cartModel{}).
			Where("id = ?", item.CartID.String()).
			Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *cartRepository) UpdateItemQty(ctx context.Context, itemID uuid.UUID, qty int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the item to find its cart_id.
		var model cartItemModel
		if err := tx.Where("id = ?", itemID.String()).First(&model).Error; err != nil {
			return err
		}

		// Update qty.
		if err := tx.Model(&cartItemModel{}).
			Where("id = ?", itemID.String()).
			Update("qty", qty).Error; err != nil {
			return err
		}

		// Update cart's updated_at.
		if err := tx.Model(&cartModel{}).
			Where("id = ?", model.CartID).
			Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *cartRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the item to find its cart_id.
		var model cartItemModel
		if err := tx.Where("id = ?", itemID.String()).First(&model).Error; err != nil {
			return err
		}

		// Delete the item.
		if err := tx.Where("id = ?", itemID.String()).Delete(&cartItemModel{}).Error; err != nil {
			return err
		}

		// Update cart's updated_at.
		if err := tx.Model(&cartModel{}).
			Where("id = ?", model.CartID).
			Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *cartRepository) ClearItems(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all items for this cart.
		if err := tx.Where("cart_id = ?", cartID.String()).Delete(&cartItemModel{}).Error; err != nil {
			return err
		}

		// Update cart's updated_at.
		if err := tx.Model(&cartModel{}).
			Where("id = ?", cartID.String()).
			Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *cartRepository) UpdateStatus(ctx context.Context, cartID uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&cartModel{}).
		Where("id = ?", cartID.String()).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// Mapper helpers.

func toCartModel(c *domain.Cart) cartModel {
	return cartModel{
		ID:        c.ID.String(),
		UserID:    c.UserID.String(),
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toDomainCart(m *cartModel) *domain.Cart {
	id, _ := uuid.Parse(m.ID)
	userID, _ := uuid.Parse(m.UserID)

	return &domain.Cart{
		ID:        id,
		UserID:    userID,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toCartItemModel(i *domain.CartItem) cartItemModel {
	return cartItemModel{
		ID:               i.ID.String(),
		CartID:           i.CartID.String(),
		ProductVariantID: i.ProductVariantID.String(),
		Qty:              i.Qty,
		CreatedAt:        i.CreatedAt,
	}
}

func toDomainCartItem(m *cartItemModel) *domain.CartItem {
	id, _ := uuid.Parse(m.ID)
	cartID, _ := uuid.Parse(m.CartID)
	variantID, _ := uuid.Parse(m.ProductVariantID)

	return &domain.CartItem{
		ID:               id,
		CartID:           cartID,
		ProductVariantID: variantID,
		Qty:              m.Qty,
		CreatedAt:        m.CreatedAt,
	}
}
