package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type productModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	VendorID      string    `gorm:"column:vendor_id"`
	CategoryID    string    `gorm:"column:category_id"`
	Name          string    `gorm:"column:name"`
	Slug          string    `gorm:"column:slug"`
	Description   string    `gorm:"column:description"`
	Status        string    `gorm:"column:status"`
	HalalAIStatus string    `gorm:"column:halal_ai_status"`
	HalalAINotes  *string   `gorm:"column:halal_ai_notes"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (productModel) TableName() string { return "products" }

type productVariantModel struct {
	ID          string  `gorm:"column:id;primaryKey"`
	ProductID   string  `gorm:"column:product_id"`
	SKU         string  `gorm:"column:sku"`
	VariantName string  `gorm:"column:variant_name"`
	Price       float64 `gorm:"column:price"`
	Currency    string  `gorm:"column:currency"`
	StockOnHand int     `gorm:"column:stock_on_hand"`
	WeightGram  *int    `gorm:"column:weight_gram"`
	IsDefault   bool    `gorm:"column:is_default"`
	IsActive    bool    `gorm:"column:is_active"`
}

func (productVariantModel) TableName() string { return "product_variants" }

type productImageModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	ProductID     string    `gorm:"column:product_id"`
	ImageURL      string    `gorm:"column:image_url"`
	MimeType      *string   `gorm:"column:mime_type"`
	FileSizeBytes *int      `gorm:"column:file_size_bytes"`
	IsPrimary     bool      `gorm:"column:is_primary"`
	SortOrder     int       `gorm:"column:sort_order"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (productImageModel) TableName() string { return "product_images" }

type productShippingServiceModel struct {
	ProductID         string    `gorm:"column:product_id;primaryKey"`
	ShippingServiceID string    `gorm:"column:shipping_service_id;primaryKey"`
	CreatedAt         time.Time `gorm:"column:created_at"`
}

func (productShippingServiceModel) TableName() string { return "product_shipping_services" }

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository backed by GORM.
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product, variants []domain.ProductVariant, images []domain.ProductImage, shippingServiceIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert product.
		pm := toProductModel(product)
		if err := tx.Create(&pm).Error; err != nil {
			return err
		}

		// 2. Insert shipping service associations (before commit for deferred trigger).
		now := time.Now()
		for _, ssID := range shippingServiceIDs {
			pss := productShippingServiceModel{
				ProductID:         product.ID.String(),
				ShippingServiceID: ssID.String(),
				CreatedAt:         now,
			}
			if err := tx.Create(&pss).Error; err != nil {
				return err
			}
		}

		// 3. Insert variants.
		for i := range variants {
			vm := toProductVariantModel(&variants[i])
			if err := tx.Create(&vm).Error; err != nil {
				return err
			}
		}

		// 4. Insert image placeholders.
		for i := range images {
			im := toProductImageModel(&images[i])
			if err := tx.Create(&im).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var model productModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainProduct(&model), nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var model productModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainProduct(&model), nil
}

func (r *productRepository) CountBySlugPrefix(ctx context.Context, prefix string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&productModel{}).
		Where("slug LIKE ?", prefix+"%").
		Count(&count).Error
	return count, err
}

func (r *productRepository) FindVariantsByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	var models []productVariantModel
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID.String()).Find(&models).Error; err != nil {
		return nil, err
	}

	variants := make([]domain.ProductVariant, len(models))
	for i, m := range models {
		variants[i] = *toDomainProductVariant(&m)
	}
	return variants, nil
}

func (r *productRepository) FindImagesByProductID(ctx context.Context, productID uuid.UUID) ([]domain.ProductImage, error) {
	var models []productImageModel
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID.String()).Order("sort_order ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	images := make([]domain.ProductImage, len(models))
	for i, m := range models {
		images[i] = *toDomainProductImage(&m)
	}
	return images, nil
}

func (r *productRepository) UpdateImages(ctx context.Context, images []domain.ProductImage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range images {
			im := toProductImageModel(&images[i])
			if err := tx.Model(&productImageModel{ID: im.ID}).
				Select("image_url", "mime_type", "file_size_bytes").
				Updates(&im).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *productRepository) CountByVendorID(ctx context.Context, vendorID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&productModel{}).
		Where("vendor_id = ?", vendorID.String()).
		Count(&count).Error
	return count, err
}

func (r *productRepository) FindVariantByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	var model productVariantModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainProductVariant(&model), nil
}

// Mapper helpers.

func toProductModel(p *domain.Product) productModel {
	return productModel{
		ID:            p.ID.String(),
		VendorID:      p.VendorID.String(),
		CategoryID:    p.CategoryID.String(),
		Name:          p.Name,
		Slug:          p.Slug,
		Description:   p.Description,
		Status:        p.Status,
		HalalAIStatus: p.HalalAIStatus,
		HalalAINotes:  p.HalalAINotes,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func toDomainProduct(m *productModel) *domain.Product {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)
	categoryID, _ := uuid.Parse(m.CategoryID)

	return &domain.Product{
		ID:            id,
		VendorID:      vendorID,
		CategoryID:    categoryID,
		Name:          m.Name,
		Slug:          m.Slug,
		Description:   m.Description,
		Status:        m.Status,
		HalalAIStatus: m.HalalAIStatus,
		HalalAINotes:  m.HalalAINotes,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func toProductVariantModel(v *domain.ProductVariant) productVariantModel {
	return productVariantModel{
		ID:          v.ID.String(),
		ProductID:   v.ProductID.String(),
		SKU:         v.SKU,
		VariantName: v.VariantName,
		Price:       v.Price,
		Currency:    v.Currency,
		StockOnHand: v.StockOnHand,
		WeightGram:  v.WeightGram,
		IsDefault:   v.IsDefault,
		IsActive:    v.IsActive,
	}
}

func toDomainProductVariant(m *productVariantModel) *domain.ProductVariant {
	id, _ := uuid.Parse(m.ID)
	productID, _ := uuid.Parse(m.ProductID)

	return &domain.ProductVariant{
		ID:          id,
		ProductID:   productID,
		SKU:         m.SKU,
		VariantName: m.VariantName,
		Price:       m.Price,
		Currency:    m.Currency,
		StockOnHand: m.StockOnHand,
		WeightGram:  m.WeightGram,
		IsDefault:   m.IsDefault,
		IsActive:    m.IsActive,
	}
}

func toProductImageModel(img *domain.ProductImage) productImageModel {
	return productImageModel{
		ID:            img.ID.String(),
		ProductID:     img.ProductID.String(),
		ImageURL:      img.ImageURL,
		MimeType:      img.MimeType,
		FileSizeBytes: img.FileSizeBytes,
		IsPrimary:     img.IsPrimary,
		SortOrder:     img.SortOrder,
		CreatedAt:     img.CreatedAt,
	}
}

func toDomainProductImage(m *productImageModel) *domain.ProductImage {
	id, _ := uuid.Parse(m.ID)
	productID, _ := uuid.Parse(m.ProductID)

	return &domain.ProductImage{
		ID:            id,
		ProductID:     productID,
		ImageURL:      m.ImageURL,
		MimeType:      m.MimeType,
		FileSizeBytes: m.FileSizeBytes,
		IsPrimary:     m.IsPrimary,
		SortOrder:     m.SortOrder,
		CreatedAt:     m.CreatedAt,
	}
}
