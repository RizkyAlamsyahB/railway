package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

type productVariantPricingRow struct {
	ID            string   `gorm:"column:id"`
	ProductID     string   `gorm:"column:product_id"`
	SKU           string   `gorm:"column:sku"`
	VariantName   string   `gorm:"column:variant_name"`
	Price         float64  `gorm:"column:price"`
	OriginalPrice float64  `gorm:"column:original_price"`
	PromoPrice    *float64 `gorm:"column:promo_price"`
	Currency      string   `gorm:"column:currency"`
	StockOnHand   int      `gorm:"column:stock_on_hand"`
	WeightGram    *int     `gorm:"column:weight_gram"`
	IsDefault     bool     `gorm:"column:is_default"`
	IsActive      bool     `gorm:"column:is_active"`
}

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

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository backed by GORM.
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product, variants []domain.ProductVariant, images []domain.ProductImage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert product.
		pm := toProductModel(product)
		if err := tx.Create(&pm).Error; err != nil {
			return err
		}

		// 2. Insert variants.
		for i := range variants {
			vm := toProductVariantModel(&variants[i])
			if err := tx.Create(&vm).Error; err != nil {
				return err
			}
		}

		// 3. Insert image placeholders.
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
	var row productVariantPricingRow
	activeVoucherDiscountExpr := buildActiveVoucherDiscountExpr("pv.product_id")
	effectivePriceExpr := buildEffectivePriceExpr("pv.price", "pv.product_id")
	if err := r.db.WithContext(ctx).
		Table("product_variants pv").
		Select(fmt.Sprintf(`
			pv.id,
			pv.product_id,
			pv.sku,
			pv.variant_name,
			pv.price AS original_price,
			%s AS promo_price,
			%s AS price,
			pv.currency,
			pv.stock_on_hand,
			pv.weight_gram,
			pv.is_default,
			pv.is_active
		`, activeVoucherDiscountExpr, effectivePriceExpr)).
		Where("pv.id = ?", id.String()).
		Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainProductVariantFromPricingRow(&row), nil
}

func (r *productRepository) GetPublishedDetailForCustomer(ctx context.Context, id uuid.UUID) (*domain.ProductDetailResponse, error) {
	type productDetailRow struct {
		ID                string    `gorm:"column:id"`
		Name              string    `gorm:"column:name"`
		Slug              string    `gorm:"column:slug"`
		Description       string    `gorm:"column:description"`
		Status            string    `gorm:"column:status"`
		HalalAIStatus     string    `gorm:"column:halal_ai_status"`
		CategoryID        string    `gorm:"column:category_id"`
		CategoryName      string    `gorm:"column:category_name"`
		CategorySlug      string    `gorm:"column:category_slug"`
		VendorID          string    `gorm:"column:vendor_id"`
		VendorDisplayName string    `gorm:"column:vendor_display_name"`
		RatingAverage     float64   `gorm:"column:rating_average"`
		RatingCount       int64     `gorm:"column:rating_count"`
		CreatedAt         time.Time `gorm:"column:created_at"`
		UpdatedAt         time.Time `gorm:"column:updated_at"`
	}

	var row productDetailRow
	if err := r.db.WithContext(ctx).
		Table("products p").
		Select(`
			p.id,
			p.name,
			p.slug,
			p.description,
			p.status,
			p.halal_ai_status,
			p.category_id,
			c.name AS category_name,
			c.slug AS category_slug,
			p.vendor_id,
			v.display_name AS vendor_display_name,
			COALESCE(prs.average_rating, 0) AS rating_average,
			COALESCE(prs.total_reviews, 0) AS rating_count,
			p.created_at,
			p.updated_at
		`).
		Joins("JOIN categories c ON c.id = p.category_id").
		Joins("JOIN vendors v ON v.id = p.vendor_id").
		Joins("LEFT JOIN product_review_stats prs ON prs.product_id = p.id").
		Where("p.id = ?", id.String()).
		Where("p.status = ?", domain.ProductStatusPublished).
		Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	activeVoucherDiscountExpr := buildActiveVoucherDiscountExpr("pv.product_id")
	effectivePriceExpr := buildEffectivePriceExpr("pv.price", "pv.product_id")
	var variantRows []productVariantPricingRow
	if err := r.db.WithContext(ctx).
		Table("product_variants pv").
		Select(fmt.Sprintf(`
			pv.id,
			pv.product_id,
			pv.sku,
			pv.variant_name,
			pv.price AS original_price,
			%s AS promo_price,
			%s AS price,
			pv.currency,
			pv.stock_on_hand,
			pv.weight_gram,
			pv.is_default,
			pv.is_active
		`, activeVoucherDiscountExpr, effectivePriceExpr)).
		Where("pv.product_id = ?", id.String()).
		Where("pv.is_active = ?", true).
		Order("pv.is_default DESC").
		Order("price ASC").
		Find(&variantRows).Error; err != nil {
		return nil, err
	}

	variantItems := make([]domain.ProductVariantResponse, len(variantRows))
	for i := range variantRows {
		v := toDomainProductVariantFromPricingRow(&variantRows[i])
		variantItems[i] = domain.ProductVariantResponse{
			ID:            v.ID,
			SKU:           v.SKU,
			VariantName:   v.VariantName,
			Price:         v.Price,
			OriginalPrice: v.OriginalPrice,
			PromoPrice:    v.PromoPrice,
			HasPromo:      v.HasPromo,
			Currency:      v.Currency,
			StockOnHand:   v.StockOnHand,
			WeightGram:    v.WeightGram,
			IsDefault:     v.IsDefault,
			IsActive:      v.IsActive,
		}
	}

	var imageModels []productImageModel
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", id.String()).
		Order("sort_order ASC").
		Find(&imageModels).Error; err != nil {
		return nil, err
	}

	imageItems := make([]domain.ProductDetailImageItem, len(imageModels))
	for i := range imageModels {
		imageID, _ := uuid.Parse(imageModels[i].ID)
		imageItems[i] = domain.ProductDetailImageItem{
			ID:        imageID,
			URL:       imageModels[i].ImageURL,
			IsPrimary: imageModels[i].IsPrimary,
			SortOrder: imageModels[i].SortOrder,
		}
	}

	productID, _ := uuid.Parse(row.ID)
	categoryID, _ := uuid.Parse(row.CategoryID)
	vendorID, _ := uuid.Parse(row.VendorID)

	return &domain.ProductDetailResponse{
		ID:            productID,
		Name:          row.Name,
		Slug:          row.Slug,
		Description:   row.Description,
		Status:        row.Status,
		HalalAIStatus: row.HalalAIStatus,
		Category: domain.ProductDetailCategory{
			ID:   categoryID,
			Name: row.CategoryName,
			Slug: row.CategorySlug,
		},
		Vendor: domain.ProductDetailVendor{
			ID:          vendorID,
			DisplayName: row.VendorDisplayName,
		},
		RatingAverage: row.RatingAverage,
		RatingCount:   row.RatingCount,
		Images:        imageItems,
		Variants:      variantItems,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

func (r *productRepository) ListPublishedForCustomer(ctx context.Context, params domain.ProductListParams) ([]domain.ProductListItem, int64, error) {
	activeVoucherDiscountExpr := buildActiveVoucherDiscountExpr("products.id")
	effectivePriceExpr := buildEffectivePriceExpr("pv.price", "products.id")

	query := r.db.WithContext(ctx).Model(&productModel{}).
		Joins("JOIN product_variants pv ON pv.product_id = products.id").
		Joins("LEFT JOIN product_review_stats prs ON prs.product_id = products.id").
		Where("products.status = ?", domain.ProductStatusPublished).
		Where("pv.is_active = ?", true).
		Where("pv.stock_on_hand > 0")

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(products.name) LIKE ?", search)
	}

	var total int64
	if err := query.Distinct("products.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type productListRow struct {
		ID            string   `gorm:"column:id"`
		Name          string   `gorm:"column:name"`
		Price         float64  `gorm:"column:price"`
		OriginalPrice float64  `gorm:"column:original_price"`
		PromoPrice    *float64 `gorm:"column:promo_price"`
		RatingAverage float64  `gorm:"column:rating_average"`
		RatingCount   int64    `gorm:"column:rating_count"`
	}

	offset := (params.Page - 1) * params.Limit
	listQuery := query.Select(fmt.Sprintf(`
			products.id,
			products.name,
			MIN(pv.price) AS original_price,
			MIN(%s) AS promo_price,
			MIN(%s) AS price,
			products.created_at,
			COALESCE(prs.average_rating, 0) AS rating_average,
			COALESCE(prs.total_reviews, 0) AS rating_count
		`, activeVoucherDiscountExpr, effectivePriceExpr)).
		Group("products.id, products.name, products.created_at, prs.average_rating, prs.total_reviews")

	if params.Sort == "cheapest" {
		listQuery = listQuery.Order("price ASC").Order("products.created_at DESC")
	} else {
		listQuery = listQuery.Order("products.created_at DESC")
	}

	var rows []productListRow
	if err := listQuery.Offset(offset).Limit(params.Limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.ProductListItem, len(rows))
	for i, row := range rows {
		id, _ := uuid.Parse(row.ID)
		promoPrice, hasPromo := normalizePromoInfo(row.OriginalPrice, row.Price)
		items[i] = domain.ProductListItem{
			ID:            id,
			Name:          row.Name,
			Price:         row.Price,
			OriginalPrice: row.OriginalPrice,
			PromoPrice:    promoPrice,
			HasPromo:      hasPromo,
			RatingAverage: row.RatingAverage,
			RatingCount:   row.RatingCount,
		}
	}

	return items, total, nil
}

func (r *productRepository) ListByVendor(ctx context.Context, vendorID uuid.UUID, params domain.VendorProductListParams) ([]domain.VendorProductListItem, int64, error) {
	type vendorProductRow struct {
		ID           string    `gorm:"column:id"`
		Name         string    `gorm:"column:name"`
		CategoryName string    `gorm:"column:category_name"`
		Price        float64   `gorm:"column:price"`
		Stock        int       `gorm:"column:stock"`
		Status       string    `gorm:"column:status"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}

	baseQuery := r.db.WithContext(ctx).
		Table("products p").
		Joins("JOIN categories c ON c.id = p.category_id").
		Joins("LEFT JOIN product_variants pv ON pv.product_id = p.id AND pv.is_active = true").
		Where("p.vendor_id = ?", vendorID.String())

	if params.Status != "" {
		baseQuery = baseQuery.Where("p.status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		baseQuery = baseQuery.Where("LOWER(p.name) LIKE ?", search)
	}

	// Count distinct products.
	var total int64
	countQuery := r.db.WithContext(ctx).
		Table("products p").
		Joins("JOIN categories c ON c.id = p.category_id").
		Where("p.vendor_id = ?", vendorID.String())

	if params.Status != "" {
		countQuery = countQuery.Where("p.status = ?", params.Status)
	}
	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		countQuery = countQuery.Where("LOWER(p.name) LIKE ?", search)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build select with aggregates.
	listQuery := baseQuery.
		Select(`
			p.id,
			p.name,
			c.name AS category_name,
			COALESCE(MIN(pv.price), 0) AS price,
			COALESCE(SUM(pv.stock_on_hand), 0) AS stock,
			p.status,
			p.created_at
		`).
		Group("p.id, p.name, c.name, p.status, p.created_at")

	// Apply sorting.
	switch params.SortBy {
	case "name":
		listQuery = listQuery.Order("p.name " + params.SortOrder)
	case "price":
		listQuery = listQuery.Order("MIN(pv.price) " + params.SortOrder)
	case "stock":
		listQuery = listQuery.Order("SUM(pv.stock_on_hand) " + params.SortOrder)
	default: // "created_at"
		listQuery = listQuery.Order("p.created_at " + params.SortOrder)
	}

	offset := (params.Page - 1) * params.Limit
	var rows []vendorProductRow
	if err := listQuery.Offset(offset).Limit(params.Limit).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.VendorProductListItem, len(rows))
	for i, row := range rows {
		id, _ := uuid.Parse(row.ID)
		items[i] = domain.VendorProductListItem{
			ID:           id,
			Name:         row.Name,
			CategoryName: row.CategoryName,
			Price:        row.Price,
			Stock:        row.Stock,
			Status:       row.Status,
			CreatedAt:    row.CreatedAt,
		}
	}

	return items, total, nil
}

func (r *productRepository) UpdateProduct(
	ctx context.Context,
	product *domain.Product,
	variantsToUpsert []domain.ProductVariant,
	variantIDsToDeactivate []uuid.UUID,
	imageIDsToDelete []uuid.UUID,
	newImages []domain.ProductImage,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update product row (excludes halal_ai_status).
		pm := toProductModel(product)
		if err := tx.Model(&productModel{ID: pm.ID}).
			Select("name", "slug", "description", "category_id", "status", "updated_at").
			Updates(&pm).Error; err != nil {
			return err
		}

		// 2. Upsert variants.
		for i := range variantsToUpsert {
			vm := toProductVariantModel(&variantsToUpsert[i])
			if variantsToUpsert[i].ID == (uuid.UUID{}) {
				if err := tx.Create(&vm).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&productVariantModel{ID: vm.ID}).
					Select("variant_name", "price", "stock_on_hand", "weight_gram", "is_active").
					Updates(&vm).Error; err != nil {
					return err
				}
			}
		}

		// 3. Deactivate removed non-default variants.
		for _, id := range variantIDsToDeactivate {
			if err := tx.Model(&productVariantModel{}).
				Where("id = ? AND is_default = false", id.String()).
				Update("is_active", false).Error; err != nil {
				return err
			}
		}

		// 4. Delete removed image rows.
		if len(imageIDsToDelete) > 0 {
			ids := make([]string, len(imageIDsToDelete))
			for i, id := range imageIDsToDelete {
				ids[i] = id.String()
			}
			if err := tx.Where("id IN ? AND product_id = ?", ids, product.ID.String()).
				Delete(&productImageModel{}).Error; err != nil {
				return err
			}
		}

		// 5. Insert new image placeholders.
		for i := range newImages {
			im := toProductImageModel(&newImages[i])
			if err := tx.Create(&im).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func buildActiveVoucherDiscountExpr(productIDCol string) string {
	return fmt.Sprintf(`
		(
			SELECT MAX(vv.discount_amount)
			FROM vendor_voucher_products vvp
			JOIN vendor_vouchers vv ON vv.id = vvp.voucher_id
			WHERE vvp.product_id = %s
				AND vv.is_active = TRUE
				AND vv.starts_at <= NOW()
				AND vv.ends_at >= NOW()
				AND vv.quota_used < vv.quota_total
		)
	`, productIDCol)
}

func buildEffectivePriceExpr(basePriceCol, productIDCol string) string {
	activeVoucherDiscountExpr := buildActiveVoucherDiscountExpr(productIDCol)
	return fmt.Sprintf(`
		GREATEST(
			%s - COALESCE(%s, 0),
			0
		)
	`, basePriceCol, activeVoucherDiscountExpr)
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
		ID:            id,
		ProductID:     productID,
		SKU:           m.SKU,
		VariantName:   m.VariantName,
		Price:         m.Price,
		OriginalPrice: m.Price,
		PromoPrice:    nil,
		HasPromo:      false,
		Currency:      m.Currency,
		StockOnHand:   m.StockOnHand,
		WeightGram:    m.WeightGram,
		IsDefault:     m.IsDefault,
		IsActive:      m.IsActive,
	}
}

func toDomainProductVariantFromPricingRow(row *productVariantPricingRow) *domain.ProductVariant {
	id, _ := uuid.Parse(row.ID)
	productID, _ := uuid.Parse(row.ProductID)
	promoPrice, hasPromo := normalizePromoInfo(row.OriginalPrice, row.Price)

	return &domain.ProductVariant{
		ID:            id,
		ProductID:     productID,
		SKU:           row.SKU,
		VariantName:   row.VariantName,
		Price:         row.Price,
		OriginalPrice: row.OriginalPrice,
		PromoPrice:    promoPrice,
		HasPromo:      hasPromo,
		Currency:      row.Currency,
		StockOnHand:   row.StockOnHand,
		WeightGram:    row.WeightGram,
		IsDefault:     row.IsDefault,
		IsActive:      row.IsActive,
	}
}

func normalizePromoInfo(originalPrice, effectivePrice float64) (*float64, bool) {
	if effectivePrice < originalPrice {
		p := effectivePrice
		return &p, true
	}
	return nil, false
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
