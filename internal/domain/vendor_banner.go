package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// VendorBanner represents a storefront banner owned by a vendor.
type VendorBanner struct {
	ID        uuid.UUID `json:"id"`
	VendorID  uuid.UUID `json:"vendor_id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request DTOs ---

// CreateVendorBannerRequest is sent by vendor to create a new storefront banner.
type CreateVendorBannerRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	ContentType string `json:"content_type" binding:"required"`
}

// ConfirmVendorBannerRequest confirms that the image has been uploaded to S3.
type ConfirmVendorBannerRequest struct {
	ObjectKey string `json:"object_key" binding:"required"`
}

// UpdateVendorBannerRequest allows updating the vendor banner title.
type UpdateVendorBannerRequest struct {
	Title *string `json:"title" binding:"omitempty,max=255"`
}

// --- Response DTOs ---

// VendorBannerResponse is the public representation of a vendor banner.
type VendorBannerResponse struct {
	ID        uuid.UUID `json:"id"`
	VendorID  uuid.UUID `json:"vendor_id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// VendorBannerPresignResponse is returned after requesting a presigned upload URL.
type VendorBannerPresignResponse struct {
	BannerID  uuid.UUID `json:"banner_id"`
	UploadURL string    `json:"upload_url"`
	ObjectKey string    `json:"object_key"`
	ExpiresIn int       `json:"expires_in"`
}

// --- Interfaces ---

// VendorBannerRepository defines persistence operations for vendor banners.
type VendorBannerRepository interface {
	Create(ctx context.Context, banner *VendorBanner) error
	FindByID(ctx context.Context, id uuid.UUID) (*VendorBanner, error)
	ListByVendorID(ctx context.Context, vendorID uuid.UUID) ([]VendorBanner, error)
	Update(ctx context.Context, banner *VendorBanner) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// VendorBannerUseCase defines business operations for vendor banners.
type VendorBannerUseCase interface {
	CreateBanner(ctx context.Context, vendorID uuid.UUID, req CreateVendorBannerRequest) (*VendorBannerPresignResponse, error)
	ConfirmBanner(ctx context.Context, vendorID uuid.UUID, bannerID uuid.UUID, req ConfirmVendorBannerRequest) (*VendorBannerResponse, error)
	ListBanners(ctx context.Context, vendorID uuid.UUID) ([]VendorBannerResponse, error)
	GetBanner(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) (*VendorBannerResponse, error)
	UpdateBanner(ctx context.Context, vendorID uuid.UUID, id uuid.UUID, req UpdateVendorBannerRequest) (*VendorBannerResponse, error)
	DeleteBanner(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) error
	// ListByVendorIDPublic is the public endpoint (no auth) for vendor store detail page.
	ListByVendorIDPublic(ctx context.Context, vendorID uuid.UUID) ([]VendorBannerResponse, error)
}
