package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ProductReview represents the product_reviews table.
type ProductReview struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	OrderID     uuid.UUID `json:"order_id"`
	OrderItemID uuid.UUID `json:"order_item_id"`
	UserID      uuid.UUID `json:"user_id"`
	Rating      int       `json:"rating"`
	ReviewText  string    `json:"review_text"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductReviewImage represents the product_review_images table.
type ProductReviewImage struct {
	ID            uuid.UUID `json:"id"`
	ReviewID      uuid.UUID `json:"review_id"`
	ObjectKey     string    `json:"object_key"`
	MimeType      *string   `json:"mime_type,omitempty"`
	FileSizeBytes *int      `json:"file_size_bytes,omitempty"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
}

// ProductReviewStats represents pre-aggregated product review metrics.
type ProductReviewStats struct {
	ProductID     uuid.UUID `json:"product_id"`
	TotalReviews  int64     `json:"total_reviews"`
	TotalStars    int64     `json:"total_stars"`
	AverageRating float64   `json:"average_rating"`
	Star0Count    int64     `json:"star_0_count"`
	Star1Count    int64     `json:"star_1_count"`
	Star2Count    int64     `json:"star_2_count"`
	Star3Count    int64     `json:"star_3_count"`
	Star4Count    int64     `json:"star_4_count"`
	Star5Count    int64     `json:"star_5_count"`
}

// ReviewImageInput references a previously uploaded object key.
type ReviewImageInput struct {
	ObjectKey string `json:"object_key" binding:"required"`
}

// CreateProductReviewRequest is the input DTO for creating a product review.
type CreateProductReviewRequest struct {
	OrderID     string             `json:"order_id" binding:"required,uuid"`
	OrderItemID string             `json:"order_item_id" binding:"required,uuid"`
	Rating      int                `json:"rating" binding:"required,min=1,max=5"`
	ReviewText  string             `json:"review_text" binding:"required,max=2000"`
	Images      []ReviewImageInput `json:"images,omitempty" binding:"omitempty,max=5,dive"`
}

// CreateProductReviewResponse is the output DTO for successful review creation.
type CreateProductReviewResponse struct {
	ReviewID    uuid.UUID `json:"review_id"`
	ProductID   uuid.UUID `json:"product_id"`
	OrderID     uuid.UUID `json:"order_id"`
	OrderItemID uuid.UUID `json:"order_item_id"`
	Rating      int       `json:"rating"`
	ReviewText  string    `json:"review_text"`
	ImageURLs   []string  `json:"image_urls"`
	CreatedAt   time.Time `json:"created_at"`
}

// PresignReviewImageRequest is the input DTO for review image upload URL generation.
type PresignReviewImageRequest struct {
	ContentType string `json:"content_type" binding:"required"`
}

// PresignReviewImageResponse is the output DTO for review image upload URL generation.
type PresignReviewImageResponse struct {
	UploadURL   string `json:"upload_url"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ProductReviewListParams holds query parameters for listing product reviews.
type ProductReviewListParams struct {
	Page  int
	Limit int
}

// ProductReviewListItem is one review item in a product review list response.
type ProductReviewListItem struct {
	ReviewID     uuid.UUID `json:"review_id"`
	UserID       uuid.UUID `json:"user_id"`
	ReviewerName string    `json:"reviewer_name"`
	Rating       int       `json:"rating"`
	ReviewText   string    `json:"review_text"`
	ImageKeys    []string  `json:"-"`
	ImageURLs    []string  `json:"image_urls"`
	CreatedAt    time.Time `json:"created_at"`
}

// RatingDistributionItem represents a star bucket in the summary response.
type RatingDistributionItem struct {
	Star       int     `json:"star"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// ProductReviewSummaryResponse is the aggregated summary of reviews for a product.
type ProductReviewSummaryResponse struct {
	ProductID     uuid.UUID                `json:"product_id"`
	AverageRating float64                  `json:"average_rating"`
	TotalReviews  int64                    `json:"total_reviews"`
	Distribution  []RatingDistributionItem `json:"distribution"`
}

// ReviewRepository defines the interface for review data access.
type ReviewRepository interface {
	// ExistsByUserAndOrderItem checks if a review already exists for a user/order item pair.
	ExistsByUserAndOrderItem(ctx context.Context, userID, orderItemID uuid.UUID) (bool, error)

	// CreateWithImages creates a review and optional images while updating product stats atomically.
	CreateWithImages(ctx context.Context, review *ProductReview, images []ProductReviewImage) error

	// ListByProductID returns product reviews and total count for pagination.
	ListByProductID(ctx context.Context, productID uuid.UUID, params ProductReviewListParams) ([]ProductReviewListItem, int64, error)

	// GetStatsByProductID returns pre-aggregated rating stats for a product.
	GetStatsByProductID(ctx context.Context, productID uuid.UUID) (*ProductReviewStats, error)
}

// ReviewUseCase defines the interface for review business operations.
type ReviewUseCase interface {
	// PresignImage creates a presigned URL for direct review image uploads.
	PresignImage(ctx context.Context, req PresignReviewImageRequest) (*PresignReviewImageResponse, error)

	// Create creates a new product review for an eligible purchased item.
	Create(ctx context.Context, userID uuid.UUID, req CreateProductReviewRequest) (*CreateProductReviewResponse, error)

	// ListByProduct returns paginated published reviews for a product.
	ListByProduct(ctx context.Context, productID uuid.UUID, params ProductReviewListParams) ([]ProductReviewListItem, *PaginationMeta, error)

	// GetSummary returns aggregated review metrics for a product.
	GetSummary(ctx context.Context, productID uuid.UUID) (*ProductReviewSummaryResponse, error)
}
