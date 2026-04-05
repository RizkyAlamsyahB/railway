package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Cart represents the carts table.
type Cart struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartItem represents the cart_items table.
type CartItem struct {
	ID               uuid.UUID `json:"id"`
	CartID           uuid.UUID `json:"cart_id"`
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	Qty              int       `json:"qty"`
	IsSelected       bool      `json:"is_selected"`
	CreatedAt        time.Time `json:"created_at"`
}

// --- Request DTOs ---

// AddCartItemRequest is the input DTO for adding an item to the cart.
type AddCartItemRequest struct {
	ProductVariantID string `json:"product_variant_id" binding:"required,uuid"`
	Qty              int    `json:"qty" binding:"required,min=1"`
}

// UpdateCartItemRequest is the input DTO for updating a cart item's quantity.
type UpdateCartItemRequest struct {
	Qty int `json:"qty" binding:"required,min=1"`
}

// UpdateCartItemSelectionRequest is the input DTO for updating a cart item's selected state.
type UpdateCartItemSelectionRequest struct {
	IsSelected bool `json:"is_selected"`
}

// --- Response DTOs ---

// CartItemResponse is a single item in the cart response with enriched product data.
type CartItemResponse struct {
	ID               uuid.UUID `json:"id"`
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	VariantName      string    `json:"variant_name"`
	SKU              string    `json:"sku"`
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`
	ProductSlug      string    `json:"product_slug"`
	ImageURL         *string   `json:"image_url"`
	Price            float64   `json:"price"`
	OriginalPrice    float64   `json:"original_price"`
	PromoPrice       *float64  `json:"promo_price,omitempty"`
	HasPromo         bool      `json:"has_promo"`
	Currency         string    `json:"currency"`
	Qty              int       `json:"qty"`
	Subtotal         float64   `json:"subtotal"`
	StockOnHand      int       `json:"stock_on_hand"`
	IsAvailable      bool      `json:"is_available"`
	IsSelected       bool      `json:"is_selected"`
	CreatedAt        time.Time `json:"created_at"`
}

// GetCartResponse is the output DTO for the full cart view.
type GetCartResponse struct {
	ID        uuid.UUID          `json:"id"`
	Status    string             `json:"status"`
	Items     []CartItemResponse `json:"items"`
	ItemCount int                `json:"item_count"`
	Total     float64            `json:"total"`
	Currency  string             `json:"currency"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// CartItemActionResponse is the output DTO for add/update/delete item operations.
type CartItemActionResponse struct {
	CartID uuid.UUID `json:"cart_id"`
	ItemID uuid.UUID `json:"item_id"`
}

// --- Repository Interface ---

// CartRepository defines the interface for cart data access.
type CartRepository interface {
	// FindByUserID returns the active cart for a user, or nil if none exists.
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Cart, error)

	// Create creates a new cart for a user.
	Create(ctx context.Context, cart *Cart) error

	// FindItemsByCartID returns all items in a cart.
	FindItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]CartItem, error)

	// FindSelectedItemsByCartID returns selected items in a cart.
	FindSelectedItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]CartItem, error)

	// FindItemByID returns a single cart item, or nil if not found.
	FindItemByID(ctx context.Context, itemID uuid.UUID) (*CartItem, error)

	// FindItemByCartAndVariant returns the cart item for a given cart and variant, or nil.
	FindItemByCartAndVariant(ctx context.Context, cartID uuid.UUID, variantID uuid.UUID) (*CartItem, error)

	// UpsertItem inserts a new cart item or increments qty if the variant already exists.
	UpsertItem(ctx context.Context, item *CartItem) error

	// UpdateItemQty updates the qty of a cart item and the cart's updated_at.
	UpdateItemQty(ctx context.Context, itemID uuid.UUID, qty int) error

	// UpdateItemSelection updates the selected state of a cart item and the cart's updated_at.
	UpdateItemSelection(ctx context.Context, itemID uuid.UUID, isSelected bool) error

	// DeleteItem removes a cart item and updates the cart's updated_at.
	DeleteItem(ctx context.Context, itemID uuid.UUID) error

	// ClearItems removes all items from a cart and updates the cart's updated_at.
	ClearItems(ctx context.Context, cartID uuid.UUID) error

	// DeleteItemsByIDs removes the specified items from a cart and updates the cart's updated_at.
	DeleteItemsByIDs(ctx context.Context, cartID uuid.UUID, itemIDs []uuid.UUID) error
}

// --- Usecase Interface ---

// CartUseCase defines the interface for cart business operations.
type CartUseCase interface {
	// GetCart returns the user's cart with enriched item details and live prices.
	GetCart(ctx context.Context, userID uuid.UUID) (*GetCartResponse, error)

	// AddItem adds a product variant to the user's cart, creating the cart if needed.
	AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*CartItemActionResponse, error)

	// UpdateItem changes the quantity of a cart item.
	UpdateItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req UpdateCartItemRequest) (*CartItemActionResponse, error)

	// UpdateItemSelection changes whether a cart item is selected for checkout.
	UpdateItemSelection(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req UpdateCartItemSelectionRequest) (*CartItemActionResponse, error)

	// RemoveItem removes a single item from the cart.
	RemoveItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error

	// ClearCart removes all items from the user's cart.
	ClearCart(ctx context.Context, userID uuid.UUID) error
}
