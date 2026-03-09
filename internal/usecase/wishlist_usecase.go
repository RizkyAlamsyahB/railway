package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type wishlistUseCase struct {
	wishlistRepo domain.WishlistRepository
	productRepo  domain.ProductRepository
	storage      domain.StorageProvider
}

// NewWishlistUseCase creates a new WishlistUseCase.
func NewWishlistUseCase(
	wishlistRepo domain.WishlistRepository,
	productRepo domain.ProductRepository,
	storage domain.StorageProvider,
) domain.WishlistUseCase {
	return &wishlistUseCase{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
		storage:      storage,
	}
}

func (uc *wishlistUseCase) AddItem(ctx context.Context, userID uuid.UUID, req domain.AddWishlistItemRequest) (*domain.WishlistAddItemResponse, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	if product.Status != domain.ProductStatusPublished {
		return nil, ErrProductNotAvailable
	}

	created, err := uc.wishlistRepo.Add(ctx, &domain.WishlistItem{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add wishlist item: %w", err)
	}

	return &domain.WishlistAddItemResponse{
		ProductID:     productID,
		IsInWishlist:  true,
		AlreadyExists: !created,
	}, nil
}

func (uc *wishlistUseCase) ListItems(ctx context.Context, userID uuid.UUID, params domain.WishlistListParams) ([]domain.WishlistItemResponse, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	items, total, err := uc.wishlistRepo.ListByUser(ctx, userID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list wishlist items: %w", err)
	}

	for i := range items {
		if items[i].ImageURL == nil || *items[i].ImageURL == "" {
			continue
		}
		if isAbsoluteWishlistURL(*items[i].ImageURL) {
			continue
		}
		absolute := uc.storage.GetURL(*items[i].ImageURL)
		items[i].ImageURL = &absolute
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}

func (uc *wishlistUseCase) RemoveItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	if err := uc.wishlistRepo.DeleteByUserAndProduct(ctx, userID, productID); err != nil {
		return fmt.Errorf("failed to remove wishlist item: %w", err)
	}
	return nil
}

func (uc *wishlistUseCase) GetProductStatus(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*domain.WishlistProductStatusResponse, error) {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	inWishlist, err := uc.wishlistRepo.Exists(ctx, userID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to check wishlist status: %w", err)
	}

	return &domain.WishlistProductStatusResponse{
		ProductID:    productID,
		IsInWishlist: inWishlist,
		IsAvailable:  product.Status == domain.ProductStatusPublished,
	}, nil
}

func isAbsoluteWishlistURL(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}
