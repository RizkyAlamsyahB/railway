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
	*mocks.MockProductRepository,
	*mocks.MockStorageProvider,
	domain.CatalogUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	categoryRepo := mocks.NewMockCategoryRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	uc := NewCatalogUseCase(categoryRepo, productRepo, storage)
	return categoryRepo, productRepo, storage, uc
}

func TestCatalogUseCase_ListProducts_DefaultsAndMeta(t *testing.T) {
	_, productRepo, _, uc := setupCatalogUseCase(t)
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
	_, productRepo, _, uc := setupCatalogUseCase(t)
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
	_, productRepo, _, uc := setupCatalogUseCase(t)
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

func TestCatalogUseCase_GetProductDetail_SuccessWithImageURLNormalization(t *testing.T) {
	_, productRepo, storage, uc := setupCatalogUseCase(t)
	ctx := context.Background()

	productID := uuid.New()
	imageID1 := uuid.New()
	imageID2 := uuid.New()

	productRepo.EXPECT().
		GetPublishedDetailForCustomer(ctx, productID).
		Return(&domain.ProductDetailResponse{
			ID: productID,
			Images: []domain.ProductDetailImageItem{
				{ID: imageID1, URL: "products/abc/image-1.jpg"},
				{ID: imageID2, URL: "https://cdn.example.com/products/image-2.jpg"},
			},
		}, nil)

	storage.EXPECT().
		GetURL("products/abc/image-1.jpg").
		Return("https://storage.example.com/products/abc/image-1.jpg")

	res, err := uc.GetProductDetail(ctx, productID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Images[0].URL != "https://storage.example.com/products/abc/image-1.jpg" {
		t.Fatalf("expected normalized URL for first image, got %s", res.Images[0].URL)
	}
	if res.Images[1].URL != "https://cdn.example.com/products/image-2.jpg" {
		t.Fatalf("expected second image URL to stay unchanged, got %s", res.Images[1].URL)
	}
}

func TestCatalogUseCase_GetProductDetail_NotFound(t *testing.T) {
	_, productRepo, _, uc := setupCatalogUseCase(t)
	ctx := context.Background()
	productID := uuid.New()

	productRepo.EXPECT().
		GetPublishedDetailForCustomer(ctx, productID).
		Return(nil, nil)

	_, err := uc.GetProductDetail(ctx, productID)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestCatalogUseCase_GetProductDetail_RepoError(t *testing.T) {
	_, productRepo, _, uc := setupCatalogUseCase(t)
	ctx := context.Background()
	productID := uuid.New()
	dbErr := errors.New("db error")

	productRepo.EXPECT().
		GetPublishedDetailForCustomer(ctx, productID).
		Return(nil, dbErr)

	_, err := uc.GetProductDetail(ctx, productID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got %v", err)
	}
}
