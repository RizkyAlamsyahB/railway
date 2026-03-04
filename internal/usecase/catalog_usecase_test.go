package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

func setupCatalogUseCase(t *testing.T) (
	*mocks.MockCategoryRepository,
	*mocks.MockShippingServiceRepository,
	*mocks.MockProductRepository,
	domain.CatalogUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	categoryRepo := mocks.NewMockCategoryRepository(ctrl)
	shippingRepo := mocks.NewMockShippingServiceRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	uc := NewCatalogUseCase(categoryRepo, shippingRepo, productRepo)
	return categoryRepo, shippingRepo, productRepo, uc
}

func TestCatalogUseCase_ListProducts_DefaultsAndMeta(t *testing.T) {
	_, _, productRepo, uc := setupCatalogUseCase(t)
	ctx := context.Background()

	expected := domain.ProductListParams{
		Page:   1,
		Limit:  10,
		Sort:   "newest",
		Search: "sajadah",
	}

	items := []domain.ProductListItem{
		{
			ID:    uuid.New(),
			Name:  "Sajadah Premium",
			Price: 150000,
		},
	}

	productRepo.EXPECT().
		ListPublishedForCustomer(ctx, expected).
		Return(items, int64(25), nil)

	res, meta, err := uc.ListProducts(ctx, domain.ProductListParams{
		Page:   0,
		Limit:  0,
		Sort:   "unknown",
		Search: "  sajadah  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 item, got %d", len(res))
	}
	if res[0].Name != "Sajadah Premium" {
		t.Errorf("expected item name Sajadah Premium, got %s", res[0].Name)
	}
	if meta.Page != 1 || meta.Limit != 10 {
		t.Errorf("unexpected meta page/limit: %+v", meta)
	}
	if meta.TotalItems != 25 || meta.TotalPages != 3 {
		t.Errorf("unexpected meta totals: %+v", meta)
	}
}

func TestCatalogUseCase_ListProducts_CheapestSort(t *testing.T) {
	_, _, productRepo, uc := setupCatalogUseCase(t)
	ctx := context.Background()

	expected := domain.ProductListParams{
		Page:   2,
		Limit:  20,
		Sort:   "cheapest",
		Search: "",
	}

	productRepo.EXPECT().
		ListPublishedForCustomer(ctx, expected).
		Return([]domain.ProductListItem{}, int64(0), nil)

	_, meta, err := uc.ListProducts(ctx, domain.ProductListParams{
		Page:   2,
		Limit:  20,
		Sort:   " CHEAPEST ",
		Search: "   ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if meta.Page != 2 || meta.Limit != 20 {
		t.Errorf("unexpected meta page/limit: %+v", meta)
	}
}

func TestCatalogUseCase_ListProducts_RepoError(t *testing.T) {
	_, _, productRepo, uc := setupCatalogUseCase(t)
	ctx := context.Background()

	dbErr := errors.New("db error")
	productRepo.EXPECT().
		ListPublishedForCustomer(ctx, domain.ProductListParams{
			Page:   1,
			Limit:  10,
			Sort:   "newest",
			Search: "",
		}).
		Return(nil, int64(0), dbErr)

	_, _, err := uc.ListProducts(ctx, domain.ProductListParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got %v", err)
	}
}
