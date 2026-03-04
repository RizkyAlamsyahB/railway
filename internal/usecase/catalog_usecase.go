package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type catalogUseCase struct {
	categoryRepo domain.CategoryRepository
	shippingRepo domain.ShippingServiceRepository
	productRepo  domain.ProductRepository
}

// NewCatalogUseCase creates a new CatalogUseCase.
func NewCatalogUseCase(
	categoryRepo domain.CategoryRepository,
	shippingRepo domain.ShippingServiceRepository,
	productRepo domain.ProductRepository,
) domain.CatalogUseCase {
	return &catalogUseCase{
		categoryRepo: categoryRepo,
		shippingRepo: shippingRepo,
		productRepo:  productRepo,
	}
}

func (uc *catalogUseCase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.ListActive(ctx)
}

func (uc *catalogUseCase) ListShippingServices(ctx context.Context) ([]domain.ShippingService, error) {
	return uc.shippingRepo.ListActive(ctx)
}

func (uc *catalogUseCase) ListProducts(ctx context.Context, params domain.ProductListParams) ([]domain.ProductListItem, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	params.Sort = strings.ToLower(strings.TrimSpace(params.Sort))
	if params.Sort != "cheapest" {
		params.Sort = "newest"
	}
	params.Search = strings.TrimSpace(params.Search)

	items, total, err := uc.productRepo.ListPublishedForCustomer(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list products: %w", err)
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}
