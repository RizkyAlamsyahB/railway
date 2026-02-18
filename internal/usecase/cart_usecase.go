package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// Cart sentinel errors.
var (
	ErrCartItemNotFound    = errors.New("cart item not found")
	ErrCartItemNotOwned    = errors.New("cart item does not belong to your cart")
	ErrVariantNotFound     = errors.New("product variant not found")
	ErrVariantNotActive    = errors.New("product variant is not active")
	ErrProductNotAvailable = errors.New("product is not available")
	ErrInsufficientStock   = errors.New("insufficient stock")
)

type cartUseCase struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
}

// NewCartUseCase creates a new CartUseCase.
func NewCartUseCase(
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
) domain.CartUseCase {
	return &cartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// getOrCreateCart finds the user's active cart or creates one lazily.
func (uc *cartUseCase) getOrCreateCart(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find cart: %w", err)
	}
	if cart != nil {
		return cart, nil
	}

	now := time.Now()
	cart = &domain.Cart{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.cartRepo.Create(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}
	return cart, nil
}

func (uc *cartUseCase) GetCart(ctx context.Context, userID uuid.UUID) (*domain.GetCartResponse, error) {
	cart, err := uc.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	items, err := uc.cartRepo.FindItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find cart items: %w", err)
	}

	cartItems := make([]domain.CartItemResponse, 0, len(items))
	var total float64

	for _, item := range items {
		variant, err := uc.productRepo.FindVariantByID(ctx, item.ProductVariantID)
		if err != nil {
			return nil, fmt.Errorf("failed to find variant: %w", err)
		}
		// Skip items whose variant has been deleted.
		if variant == nil {
			continue
		}

		product, err := uc.productRepo.FindByID(ctx, variant.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to find product: %w", err)
		}
		if product == nil {
			continue
		}

		// Find primary image.
		var imageURL *string
		images, err := uc.productRepo.FindImagesByProductID(ctx, variant.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to find product images: %w", err)
		}
		for _, img := range images {
			if img.IsPrimary {
				imageURL = &img.ImageURL
				break
			}
		}

		isAvailable := variant.IsActive && product.Status == "published"
		subtotal := variant.Price * float64(item.Qty)

		if isAvailable {
			total += subtotal
		}

		cartItems = append(cartItems, domain.CartItemResponse{
			ID:               item.ID,
			ProductVariantID: item.ProductVariantID,
			VariantName:      variant.VariantName,
			SKU:              variant.SKU,
			ProductID:        product.ID,
			ProductName:      product.Name,
			ProductSlug:      product.Slug,
			ImageURL:         imageURL,
			Price:            variant.Price,
			Currency:         variant.Currency,
			Qty:              item.Qty,
			Subtotal:         subtotal,
			StockOnHand:      variant.StockOnHand,
			IsAvailable:      isAvailable,
			CreatedAt:        item.CreatedAt,
		})
	}

	return &domain.GetCartResponse{
		ID:        cart.ID,
		Status:    cart.Status,
		Items:     cartItems,
		ItemCount: len(cartItems),
		Total:     total,
		Currency:  "IDR",
		UpdatedAt: cart.UpdatedAt,
	}, nil
}

func (uc *cartUseCase) AddItem(ctx context.Context, userID uuid.UUID, req domain.AddCartItemRequest) (*domain.CartItemActionResponse, error) {
	variantID, err := uuid.Parse(req.ProductVariantID)
	if err != nil {
		return nil, ErrVariantNotFound
	}

	// Validate variant.
	variant, err := uc.productRepo.FindVariantByID(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find variant: %w", err)
	}
	if variant == nil {
		return nil, ErrVariantNotFound
	}
	if !variant.IsActive {
		return nil, ErrVariantNotActive
	}

	// Validate product.
	product, err := uc.productRepo.FindByID(ctx, variant.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	if product.Status != "published" {
		return nil, ErrProductNotAvailable
	}

	// Get or create cart.
	cart, err := uc.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if variant already in cart and validate stock.
	existingItem, err := uc.cartRepo.FindItemByCartAndVariant(ctx, cart.ID, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing cart item: %w", err)
	}

	totalQty := req.Qty
	if existingItem != nil {
		totalQty += existingItem.Qty
	}
	if totalQty > variant.StockOnHand {
		return nil, ErrInsufficientStock
	}

	// Upsert item.
	item := &domain.CartItem{
		ID:               uuid.New(),
		CartID:           cart.ID,
		ProductVariantID: variantID,
		Qty:              req.Qty,
		CreatedAt:        time.Now(),
	}
	if err := uc.cartRepo.UpsertItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to upsert cart item: %w", err)
	}

	// If upserted (existing item), return existing item ID.
	itemID := item.ID
	if existingItem != nil {
		itemID = existingItem.ID
	}

	return &domain.CartItemActionResponse{
		CartID: cart.ID,
		ItemID: itemID,
	}, nil
}

func (uc *cartUseCase) UpdateItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req domain.UpdateCartItemRequest) (*domain.CartItemActionResponse, error) {
	cart, err := uc.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Find the item.
	item, err := uc.cartRepo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find cart item: %w", err)
	}
	if item == nil {
		return nil, ErrCartItemNotFound
	}
	if item.CartID != cart.ID {
		return nil, ErrCartItemNotOwned
	}

	// Validate stock.
	variant, err := uc.productRepo.FindVariantByID(ctx, item.ProductVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find variant: %w", err)
	}
	if variant == nil {
		return nil, ErrVariantNotFound
	}
	if req.Qty > variant.StockOnHand {
		return nil, ErrInsufficientStock
	}

	if err := uc.cartRepo.UpdateItemQty(ctx, itemID, req.Qty); err != nil {
		return nil, fmt.Errorf("failed to update cart item qty: %w", err)
	}

	return &domain.CartItemActionResponse{
		CartID: cart.ID,
		ItemID: itemID,
	}, nil
}

func (uc *cartUseCase) RemoveItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	cart, err := uc.getOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	item, err := uc.cartRepo.FindItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to find cart item: %w", err)
	}
	if item == nil {
		return ErrCartItemNotFound
	}
	if item.CartID != cart.ID {
		return ErrCartItemNotOwned
	}

	if err := uc.cartRepo.DeleteItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	return nil
}

func (uc *cartUseCase) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find cart: %w", err)
	}
	if cart == nil {
		return nil
	}

	items, err := uc.cartRepo.FindItemsByCartID(ctx, cart.ID)
	if err != nil {
		return fmt.Errorf("failed to find cart items: %w", err)
	}
	if len(items) == 0 {
		return nil
	}

	if err := uc.cartRepo.ClearItems(ctx, cart.ID); err != nil {
		return fmt.Errorf("failed to clear cart items: %w", err)
	}

	return nil
}
