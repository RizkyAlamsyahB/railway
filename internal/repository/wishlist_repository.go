package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type wishlistItemModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	ProductID string    `gorm:"column:product_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (wishlistItemModel) TableName() string { return "wishlist_items" }

type wishlistRepository struct {
	db *gorm.DB
}

// NewWishlistRepository creates a new WishlistRepository backed by GORM.
func NewWishlistRepository(db *gorm.DB) domain.WishlistRepository {
	return &wishlistRepository{db: db}
}

func (r *wishlistRepository) Add(ctx context.Context, item *domain.WishlistItem) (bool, error) {
	model := toWishlistItemModel(item)
	tx := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "product_id"},
		},
		DoNothing: true,
	}).Create(&model)
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, nil
}

func (r *wishlistRepository) DeleteByUserAndProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID.String(), productID.String()).
		Delete(&wishlistItemModel{}).Error
}

func (r *wishlistRepository) Exists(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&wishlistItemModel{}).
		Where("user_id = ? AND product_id = ?", userID.String(), productID.String()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *wishlistRepository) ListByUser(ctx context.Context, userID uuid.UUID, params domain.WishlistListParams) ([]domain.WishlistItemResponse, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Model(&wishlistItemModel{}).
		Where("user_id = ?", userID.String()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type wishlistRow struct {
		ProductID   string    `gorm:"column:product_id"`
		ProductName string    `gorm:"column:product_name"`
		ProductSlug string    `gorm:"column:product_slug"`
		VendorID    string    `gorm:"column:vendor_id"`
		VendorName  string    `gorm:"column:vendor_name"`
		ImageURL    *string   `gorm:"column:image_url"`
		Price       *float64  `gorm:"column:price"`
		Currency    *string   `gorm:"column:currency"`
		IsAvailable bool      `gorm:"column:is_available"`
		AddedAt     time.Time `gorm:"column:added_at"`
	}

	offset := (params.Page - 1) * params.Limit

	var rows []wishlistRow
	err := r.db.WithContext(ctx).
		Table("wishlist_items wi").
		Select(`
			wi.product_id AS product_id,
			p.name AS product_name,
			p.slug AS product_slug,
			p.vendor_id AS vendor_id,
			v.display_name AS vendor_name,
			pi.image_url AS image_url,
			pv.min_price AS price,
			pv.currency AS currency,
			(p.status = ?) AS is_available,
			wi.created_at AS added_at
		`, domain.ProductStatusPublished).
		Joins("JOIN products p ON p.id = wi.product_id").
		Joins("JOIN vendors v ON v.id = p.vendor_id").
		Joins(`
			LEFT JOIN (
				SELECT product_id, MIN(price) AS min_price, MIN(currency) AS currency
				FROM product_variants
				WHERE is_active = true
				GROUP BY product_id
			) pv ON pv.product_id = p.id
		`).
		Joins("LEFT JOIN product_images pi ON pi.product_id = p.id AND pi.is_primary = true").
		Where("wi.user_id = ?", userID.String()).
		Order("wi.created_at DESC").
		Order("wi.id DESC").
		Offset(offset).
		Limit(params.Limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	items := make([]domain.WishlistItemResponse, len(rows))
	for i, row := range rows {
		productID, _ := uuid.Parse(row.ProductID)
		vendorID, _ := uuid.Parse(row.VendorID)

		items[i] = domain.WishlistItemResponse{
			ProductID:   productID,
			ProductName: row.ProductName,
			ProductSlug: row.ProductSlug,
			VendorID:    vendorID,
			VendorName:  row.VendorName,
			ImageURL:    row.ImageURL,
			Price:       row.Price,
			Currency:    row.Currency,
			IsAvailable: row.IsAvailable,
			AddedAt:     row.AddedAt,
		}
	}

	return items, total, nil
}

func toWishlistItemModel(item *domain.WishlistItem) wishlistItemModel {
	return wishlistItemModel{
		ID:        item.ID.String(),
		UserID:    item.UserID.String(),
		ProductID: item.ProductID.String(),
		CreatedAt: item.CreatedAt,
	}
}
