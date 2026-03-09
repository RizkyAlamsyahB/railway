package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// WishlistItem represents the wishlist_items table.
type WishlistItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

// AddWishlistItemRequest is the input DTO for adding product to wishlist.
type AddWishlistItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
}

// WishlistListParams holds query parameters for wishlist listing.
type WishlistListParams struct {
	Page  int
	Limit int
}

// WishlistAddItemResponse is the output DTO for add-to-wishlist operation.
type WishlistAddItemResponse struct {
	ProductID     uuid.UUID `json:"product_id"`
	IsInWishlist  bool      `json:"is_in_wishlist"`
	AlreadyExists bool      `json:"already_exists"`
}

// WishlistItemResponse is the output DTO for a single wishlist item.
type WishlistItemResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductSlug string    `json:"product_slug"`
	VendorID    uuid.UUID `json:"vendor_id"`
	VendorName  string    `json:"vendor_name"`
	ImageURL    *string   `json:"image_url"`
	Price       *float64  `json:"price"`
	Currency    *string   `json:"currency"`
	IsAvailable bool      `json:"is_available"`
	AddedAt     time.Time `json:"added_at"`
}

// WishlistProductStatusResponse is the output DTO for product wishlist status.
type WishlistProductStatusResponse struct {
	ProductID    uuid.UUID `json:"product_id"`
	IsInWishlist bool      `json:"is_in_wishlist"`
	IsAvailable  bool      `json:"is_available"`
}

// WishlistRepository defines the interface for wishlist data access.
type WishlistRepository interface {
	// Add inserts wishlist item and returns true when newly created.
	Add(ctx context.Context, item *WishlistItem) (bool, error)

	// DeleteByUserAndProduct deletes wishlist item by user and product.
	DeleteByUserAndProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error

	// Exists checks whether a wishlist item exists for given user and product.
	Exists(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (bool, error)

	// ListByUser returns paginated wishlist items with enriched product info.
	ListByUser(ctx context.Context, userID uuid.UUID, params WishlistListParams) ([]WishlistItemResponse, int64, error)
}

// WishlistUseCase defines the interface for wishlist business operations.
type WishlistUseCase interface {
	// AddItem adds a product to user's wishlist.
	AddItem(ctx context.Context, userID uuid.UUID, req AddWishlistItemRequest) (*WishlistAddItemResponse, error)

	// ListItems returns paginated wishlist items.
	ListItems(ctx context.Context, userID uuid.UUID, params WishlistListParams) ([]WishlistItemResponse, *PaginationMeta, error)

	// RemoveItem removes a product from user's wishlist.
	RemoveItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error

	// GetProductStatus returns wishlist status for specific product.
	GetProductStatus(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*WishlistProductStatusResponse, error)
}
