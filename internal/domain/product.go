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
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	SKU         string    `json:"sku"`
	VariantName string    `json:"variant_name"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	StockOnHand int       `json:"stock_on_hand"`
	WeightGram  *int      `json:"weight_gram,omitempty"`
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
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

// ShippingService represents the shipping_services table.
type ShippingService struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
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
	Name               string                      `json:"name" binding:"required,max=180"`
	CategoryID         string                      `json:"category_id" binding:"required,uuid"`
	Description        string                      `json:"description" binding:"required"`
	Price              float64                     `json:"price" binding:"required,gt=0"`
	Stock              int                         `json:"stock" binding:"min=0"`
	WeightGram         *int                        `json:"weight_gram,omitempty" binding:"omitempty,min=0"`
	IsActive           bool                        `json:"is_active"`
	ShippingServiceIDs []string                    `json:"shipping_service_ids" binding:"required,min=1,dive,uuid"`
	Variants           []CreateProductVariantInput  `json:"variants,omitempty" binding:"omitempty,max=20,dive"`
	Images             []CreateProductImageInput    `json:"images,omitempty" binding:"omitempty,min=1,max=10,dive"`
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
	ID          uuid.UUID `json:"id"`
	SKU         string    `json:"sku"`
	VariantName string    `json:"variant_name"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	StockOnHand int       `json:"stock_on_hand"`
	WeightGram  *int      `json:"weight_gram,omitempty"`
	IsDefault   bool      `json:"is_default"`
	IsActive    bool      `json:"is_active"`
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

// --- Repository Interfaces ---

// ProductRepository defines the interface for product data access.
type ProductRepository interface {
	// Create inserts a product, its variants, image placeholders, and
	// shipping service associations in a single transaction.
	Create(ctx context.Context, product *Product, variants []ProductVariant,
		images []ProductImage, shippingServiceIDs []uuid.UUID) error

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
}

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	// ListActive returns all active categories.
	ListActive(ctx context.Context) ([]Category, error)

	// FindByID returns the category with the given ID, or nil if not found.
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
}

// ShippingServiceRepository defines the interface for shipping service data access.
type ShippingServiceRepository interface {
	// ListActive returns all active shipping services.
	ListActive(ctx context.Context) ([]ShippingService, error)

	// FindByIDs returns active shipping services matching the given IDs.
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]ShippingService, error)
}

// --- Usecase Interfaces ---

// ProductUseCase defines the interface for product business operations.
type ProductUseCase interface {
	// Create creates a new product with variants, image placeholders, and shipping service associations.
	Create(ctx context.Context, vendorID uuid.UUID, req CreateProductRequest) (*CreateProductResponse, error)

	// ConfirmImages verifies product images were uploaded to S3 and updates metadata.
	ConfirmImages(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID,
		req ConfirmProductImagesRequest) (*ConfirmProductImagesResponse, error)
}

// CatalogUseCase defines the interface for public catalog queries.
type CatalogUseCase interface {
	// ListCategories returns all active categories.
	ListCategories(ctx context.Context) ([]Category, error)

	// ListShippingServices returns all active shipping services.
	ListShippingServices(ctx context.Context) ([]ShippingService, error)
}
