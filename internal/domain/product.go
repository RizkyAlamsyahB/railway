package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Product represents the products table.
type Product struct {
	ID            uuid.UUID `json:"id"`
	VendorID      uuid.UUID `json:"vendor_id"`
	CategoryID    uuid.UUID `json:"category_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	HalalAIStatus string    `json:"halal_ai_status"`
	HalalAINotes  *string   `json:"halal_ai_notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ProductVariant represents the product_variants table.
type ProductVariant struct {
	ID            uuid.UUID `json:"id"`
	ProductID     uuid.UUID `json:"product_id"`
	SKU           string    `json:"sku"`
	VariantName   string    `json:"variant_name"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	PromoPrice    *float64  `json:"promo_price,omitempty"`
	HasPromo      bool      `json:"has_promo"`
	Currency      string    `json:"currency"`
	StockOnHand   int       `json:"stock_on_hand"`
	WeightGram    *int      `json:"weight_gram,omitempty"`
	IsDefault     bool      `json:"is_default"`
	IsActive      bool      `json:"is_active"`
}

// ProductImage represents the product_images table.
type ProductImage struct {
	ID            uuid.UUID `json:"id"`
	ProductID     uuid.UUID `json:"product_id"`
	ImageURL      string    `json:"image_url"`
	MimeType      *string   `json:"mime_type,omitempty"`
	FileSizeBytes *int      `json:"file_size_bytes,omitempty"`
	IsPrimary     bool      `json:"is_primary"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
}

// Category represents the categories table.
type Category struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	IsActive bool       `json:"is_active"`
}

// --- Request DTOs ---

// CreateProductImageInput represents a single image in the create product request.
type CreateProductImageInput struct {
	FileName    string `json:"file_name" binding:"required,max=255"`
	ContentType string `json:"content_type" binding:"required,oneof=image/jpeg image/png image/webp"`
	IsPrimary   bool   `json:"is_primary"`
}

// CreateProductVariantInput represents an additional variant in the create product request.
type CreateProductVariantInput struct {
	VariantName string  `json:"variant_name" binding:"required,max=120"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	WeightGram  *int    `json:"weight_gram,omitempty" binding:"omitempty,min=0"`
	IsActive    bool    `json:"is_active"`
}

// CreateProductRequest is the input DTO for creating a product.
type CreateProductRequest struct {
	Name        string                      `json:"name" binding:"required,max=180"`
	CategoryID  string                      `json:"category_id" binding:"required,uuid"`
	Description string                      `json:"description" binding:"required"`
	Price       float64                     `json:"price" binding:"required,gt=0"`
	Stock       int                         `json:"stock" binding:"min=0"`
	WeightGram  *int                        `json:"weight_gram,omitempty" binding:"omitempty,min=0"`
	IsActive    bool                        `json:"is_active"`
	Variants    []CreateProductVariantInput `json:"variants,omitempty" binding:"omitempty,max=20,dive"`
	Images      []CreateProductImageInput   `json:"images,omitempty" binding:"omitempty,min=1,max=10,dive"`
}

// UpdateProductVariantInput represents a single variant entry in the update product request.
// If ID is provided, the variant is updated; if absent, a new variant is created.
type UpdateProductVariantInput struct {
	ID          *string `json:"id,omitempty" binding:"omitempty,uuid"`
	VariantName string  `json:"variant_name" binding:"required,max=120"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	WeightGram  *int    `json:"weight_gram,omitempty" binding:"omitempty,min=0"`
	IsActive    bool    `json:"is_active"`
}

// UpdateProductRequest is the input DTO for updating a product (full-replace PUT semantics).
type UpdateProductRequest struct {
	Name         string                      `json:"name" binding:"required,max=180"`
	CategoryID   string                      `json:"category_id" binding:"required,uuid"`
	Description  string                      `json:"description" binding:"required"`
	Price        float64                     `json:"price" binding:"required,gt=0"`
	Stock        int                         `json:"stock" binding:"min=0"`
	WeightGram   *int                        `json:"weight_gram,omitempty" binding:"omitempty,min=0"`
	IsActive     bool                        `json:"is_active"`
	Variants     []UpdateProductVariantInput `json:"variants,omitempty" binding:"omitempty,max=20,dive"`
	KeepImageIDs []string                    `json:"keep_image_ids,omitempty" binding:"omitempty,dive,uuid"`
	NewImages    []CreateProductImageInput   `json:"new_images,omitempty" binding:"omitempty,max=10,dive"`
}

// UpdateProductResponse is the output DTO for a successful product update.
type UpdateProductResponse struct {
	ID            uuid.UUID                `json:"id"`
	VendorID      uuid.UUID                `json:"vendor_id"`
	CategoryID    uuid.UUID                `json:"category_id"`
	Name          string                   `json:"name"`
	Slug          string                   `json:"slug"`
	Description   string                   `json:"description"`
	Status        string                   `json:"status"`
	HalalAIStatus string                   `json:"halal_ai_status"`
	Variants      []ProductVariantResponse `json:"variants"`
	NewUploadURLs []ProductImageUploadInfo `json:"new_upload_urls,omitempty"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

// ConfirmProductImageItem represents a single image in the confirm-images request.
type ConfirmProductImageItem struct {
	ImageID   string `json:"image_id" binding:"required,uuid"`
	ObjectKey string `json:"object_key" binding:"required"`
}

// ConfirmProductImagesRequest is the input DTO for confirming product image uploads.
type ConfirmProductImagesRequest struct {
	Images []ConfirmProductImageItem `json:"images" binding:"required,min=1,dive"`
}

// --- Response DTOs ---

// ProductImageUploadInfo holds a presigned URL for image upload.
type ProductImageUploadInfo struct {
	ImageID   uuid.UUID `json:"image_id"`
	UploadURL string    `json:"upload_url"`
	ObjectKey string    `json:"object_key"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
}

// ProductVariantResponse is the output DTO for a product variant.
type ProductVariantResponse struct {
	ID            uuid.UUID `json:"id"`
	SKU           string    `json:"sku"`
	VariantName   string    `json:"variant_name"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	PromoPrice    *float64  `json:"promo_price,omitempty"`
	HasPromo      bool      `json:"has_promo"`
	Currency      string    `json:"currency"`
	StockOnHand   int       `json:"stock_on_hand"`
	WeightGram    *int      `json:"weight_gram,omitempty"`
	IsDefault     bool      `json:"is_default"`
	IsActive      bool      `json:"is_active"`
}

// CreateProductResponse is the output DTO for a successful product creation.
type CreateProductResponse struct {
	ID            uuid.UUID                `json:"id"`
	VendorID      uuid.UUID                `json:"vendor_id"`
	CategoryID    uuid.UUID                `json:"category_id"`
	Name          string                   `json:"name"`
	Slug          string                   `json:"slug"`
	Description   string                   `json:"description"`
	Status        string                   `json:"status"`
	HalalAIStatus string                   `json:"halal_ai_status"`
	Variants      []ProductVariantResponse `json:"variants"`
	UploadURLs    []ProductImageUploadInfo `json:"upload_urls,omitempty"`
	CreatedAt     time.Time                `json:"created_at"`
}

// ConfirmProductImagesResponse is the output DTO for successful image confirmation.
type ConfirmProductImagesResponse struct {
	ProductID       uuid.UUID `json:"product_id"`
	ImagesConfirmed int       `json:"images_confirmed"`
}

// ProductListParams holds query parameters for customer product listing.
type ProductListParams struct {
	Page   int
	Limit  int
	Sort   string
	Search string
}

// ProductListItem is the output DTO for customer product listing.
type ProductListItem struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	PromoPrice    *float64  `json:"promo_price,omitempty"`
	HasPromo      bool      `json:"has_promo"`
	RatingAverage float64   `json:"rating_average"`
	RatingCount   int64     `json:"rating_count"`
}

// ProductDetailCategory is category info in customer product detail.
type ProductDetailCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// ProductDetailVendor is vendor info in customer product detail.
type ProductDetailVendor struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
}

// ProductDetailImageItem is image info in customer product detail.
type ProductDetailImageItem struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
}

// ProductDetailResponse is the output DTO for customer product detail.
type ProductDetailResponse struct {
	ID            uuid.UUID                `json:"id"`
	Name          string                   `json:"name"`
	Slug          string                   `json:"slug"`
	Description   string                   `json:"description"`
	Status        string                   `json:"status"`
	HalalAIStatus string                   `json:"halal_ai_status"`
	Category      ProductDetailCategory    `json:"category"`
	Vendor        ProductDetailVendor      `json:"vendor"`
	RatingAverage float64                  `json:"rating_average"`
	RatingCount   int64                    `json:"rating_count"`
	Images        []ProductDetailImageItem `json:"images"`
	Variants      []ProductVariantResponse `json:"variants"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}

// VendorProductListParams holds query parameters for vendor product listing.
type VendorProductListParams struct {
	Page      int
	Limit     int
	Search    string
	Status    string // filter: "draft", "published", or "" (all)
	SortBy    string // allowed: "name", "price", "stock", "created_at"
	SortOrder string // allowed: "asc", "desc"
}

// VendorProductListItem is the output DTO for vendor product listing.
type VendorProductListItem struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CategoryName string    `json:"category_name"`
	Price        float64   `json:"price"`
	Stock        int       `json:"stock"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// --- Repository Interfaces ---

// ProductRepository defines the interface for product data access.
type ProductRepository interface {
	// Create inserts a product, its variants, and image placeholders in a single transaction.
	Create(ctx context.Context, product *Product, variants []ProductVariant,
		images []ProductImage) error

	// FindByID returns the product with the given ID, or nil if not found.
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)

	// FindBySlug returns the product with the given slug, or nil if not found.
	FindBySlug(ctx context.Context, slug string) (*Product, error)

	// CountBySlugPrefix counts products whose slug starts with the given prefix.
	CountBySlugPrefix(ctx context.Context, prefix string) (int64, error)

	// FindVariantsByProductID returns all variants for a product.
	FindVariantsByProductID(ctx context.Context, productID uuid.UUID) ([]ProductVariant, error)

	// FindImagesByProductID returns all images for a product.
	FindImagesByProductID(ctx context.Context, productID uuid.UUID) ([]ProductImage, error)

	// UpdateImages updates image metadata for confirmed images.
	UpdateImages(ctx context.Context, images []ProductImage) error

	// CountByVendorID returns the total number of products for a vendor.
	CountByVendorID(ctx context.Context, vendorID uuid.UUID) (int64, error)

	// FindVariantByID returns a single product variant by its ID, or nil if not found.
	FindVariantByID(ctx context.Context, id uuid.UUID) (*ProductVariant, error)

	// ListPublishedForCustomer returns published products that have at least one
	// active variant with available stock.
	ListPublishedForCustomer(ctx context.Context, params ProductListParams) ([]ProductListItem, int64, error)

	// GetPublishedDetailForCustomer returns a single published product detail for customers.
	GetPublishedDetailForCustomer(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error)

	// ListByVendor returns paginated products owned by a specific vendor.
	ListByVendor(ctx context.Context, vendorID uuid.UUID, params VendorProductListParams) ([]VendorProductListItem, int64, error)

	// UpdateProduct persists all product edits atomically in a single transaction.
	// It updates the product row, upserts variants, deactivates removed variants,
	// deletes image rows, and inserts new image placeholders.
	UpdateProduct(ctx context.Context, product *Product,
		variantsToUpsert []ProductVariant,
		variantIDsToDeactivate []uuid.UUID,
		imageIDsToDelete []uuid.UUID,
		newImages []ProductImage) error
}

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	// ListActive returns all active categories.
	ListActive(ctx context.Context) ([]Category, error)

	// FindByID returns the category with the given ID, or nil if not found.
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
}

// --- Usecase Interfaces ---

// ProductUseCase defines the interface for product business operations.
type ProductUseCase interface {
	// Create creates a new product with variants and image placeholders.
	Create(ctx context.Context, vendorID uuid.UUID, req CreateProductRequest) (*CreateProductResponse, error)

	// ConfirmImages verifies product images were uploaded to S3 and updates metadata.
	ConfirmImages(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID,
		req ConfirmProductImagesRequest) (*ConfirmProductImagesResponse, error)

	// ListProducts returns paginated products belonging to the authenticated vendor.
	ListProducts(ctx context.Context, vendorID uuid.UUID, params VendorProductListParams) ([]VendorProductListItem, *PaginationMeta, error)

	// UpdateProduct updates an existing product owned by the authenticated vendor.
	UpdateProduct(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID,
		req UpdateProductRequest) (*UpdateProductResponse, error)
}

// CatalogUseCase defines the interface for public catalog queries.
type CatalogUseCase interface {
	// ListCategories returns all active categories.
	ListCategories(ctx context.Context) ([]Category, error)

	// ListProducts returns products available for customers.
	ListProducts(ctx context.Context, params ProductListParams) ([]ProductListItem, *PaginationMeta, error)

	// GetProductDetail returns customer-facing product detail by ID.
	GetProductDetail(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error)
}
