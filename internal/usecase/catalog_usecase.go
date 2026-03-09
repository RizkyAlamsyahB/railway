package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type catalogUseCase struct {
	categoryRepo domain.CategoryRepository
	productRepo  domain.ProductRepository
	storage      domain.StorageProvider
}

// NewCatalogUseCase creates a new CatalogUseCase.
func NewCatalogUseCase(
	categoryRepo domain.CategoryRepository,
	productRepo domain.ProductRepository,
	storage domain.StorageProvider,
) domain.CatalogUseCase {
	return &catalogUseCase{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
		storage:      storage,
	}
}

func (uc *catalogUseCase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.ListActive(ctx)
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

func (uc *catalogUseCase) GetProductDetail(ctx context.Context, id uuid.UUID) (*domain.ProductDetailResponse, error) {
	detail, err := uc.productRepo.GetPublishedDetailForCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product detail: %w", err)
	}
	if detail == nil {
		return nil, ErrProductNotFound
	}

	for i := range detail.Images {
		if !isAbsoluteURL(detail.Images[i].URL) {
			detail.Images[i].URL = uc.storage.GetURL(detail.Images[i].URL)
		}
	}

	return detail, nil
}

func isAbsoluteURL(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}
