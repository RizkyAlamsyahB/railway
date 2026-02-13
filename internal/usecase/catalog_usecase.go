package usecase

import (
	"context"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type catalogUseCase struct {
	categoryRepo domain.CategoryRepository
	shippingRepo domain.ShippingServiceRepository
}

// NewCatalogUseCase creates a new CatalogUseCase.
func NewCatalogUseCase(
	categoryRepo domain.CategoryRepository,
	shippingRepo domain.ShippingServiceRepository,
) domain.CatalogUseCase {
	return &catalogUseCase{
		categoryRepo: categoryRepo,
		shippingRepo: shippingRepo,
	}
}

func (uc *catalogUseCase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.ListActive(ctx)
}

func (uc *catalogUseCase) ListShippingServices(ctx context.Context) ([]domain.ShippingService, error) {
	return uc.shippingRepo.ListActive(ctx)
}
