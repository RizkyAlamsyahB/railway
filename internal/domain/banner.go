package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Banner represents a homepage banner managed by admin.
type Banner struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request DTOs ---

// CreateBannerRequest is sent by admin to create a new banner.
// The flow is: presign → upload to S3 → create banner with the object key.
type CreateBannerRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	ContentType string `json:"content_type" binding:"required"`
}

// ConfirmBannerRequest confirms that the image has been uploaded to S3.
type ConfirmBannerRequest struct {
	ObjectKey string `json:"object_key" binding:"required"`
}

// UpdateBannerRequest allows updating the banner title.
type UpdateBannerRequest struct {
	Title *string `json:"title" binding:"omitempty,max=255"`
}

// --- Response DTOs ---

// BannerResponse is the public representation of a banner.
type BannerResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BannerPresignResponse is returned after requesting a presigned upload URL.
type BannerPresignResponse struct {
	BannerID  uuid.UUID `json:"banner_id"`
	UploadURL string    `json:"upload_url"`
	ObjectKey string    `json:"object_key"`
	ExpiresIn int       `json:"expires_in"`
}

// --- Interfaces ---

// BannerRepository defines persistence operations for banners.
type BannerRepository interface {
	Create(ctx context.Context, banner *Banner) error
	FindByID(ctx context.Context, id uuid.UUID) (*Banner, error)
	List(ctx context.Context) ([]Banner, error)
	Update(ctx context.Context, banner *Banner) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// BannerUseCase defines business operations for banners.
type BannerUseCase interface {
	CreateBanner(ctx context.Context, req CreateBannerRequest) (*BannerPresignResponse, error)
	ConfirmBanner(ctx context.Context, bannerID uuid.UUID, req ConfirmBannerRequest) (*BannerResponse, error)
	ListBanners(ctx context.Context) ([]BannerResponse, error)
	GetBanner(ctx context.Context, id uuid.UUID) (*BannerResponse, error)
	UpdateBanner(ctx context.Context, id uuid.UUID, req UpdateBannerRequest) (*BannerResponse, error)
	DeleteBanner(ctx context.Context, id uuid.UUID) error
}
