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

func setupWishlistUseCase(t *testing.T) (
	*mocks.MockWishlistRepository,
	*mocks.MockProductRepository,
	*mocks.MockStorageProvider,
	domain.WishlistUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	wishlistRepo := mocks.NewMockWishlistRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	storageProvider := mocks.NewMockStorageProvider(ctrl)
	uc := NewWishlistUseCase(wishlistRepo, productRepo, storageProvider)
	return wishlistRepo, productRepo, storageProvider, uc
}

func TestWishlistUseCaseAddItem(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	productID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	errDB := errors.New("db error")

	tests := []struct {
		name    string
		req     domain.AddWishlistItemRequest
		setup   func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository)
		wantErr error
		wantDup bool
	}{
		{
			name: "success - newly added",
			req:  domain.AddWishlistItemRequest{ProductID: productID.String()},
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: domain.ProductStatusPublished,
				}, nil)
				wr.EXPECT().Add(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			wantDup: false,
		},
		{
			name: "success - idempotent existing item",
			req:  domain.AddWishlistItemRequest{ProductID: productID.String()},
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: domain.ProductStatusPublished,
				}, nil)
				wr.EXPECT().Add(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			wantDup: true,
		},
		{
			name: "error - invalid UUID",
			req:  domain.AddWishlistItemRequest{ProductID: "not-a-uuid"},
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
			},
			wantErr: ErrProductNotFound,
		},
		{
			name: "error - product not found",
			req:  domain.AddWishlistItemRequest{ProductID: productID.String()},
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, nil)
			},
			wantErr: ErrProductNotFound,
		},
		{
			name: "error - product not available",
			req:  domain.AddWishlistItemRequest{ProductID: productID.String()},
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: domain.ProductStatusDraft,
				}, nil)
			},
			wantErr: ErrProductNotAvailable,
		},
		{
			name: "error - repository failure",
			req:  domain.AddWishlistItemRequest{ProductID: productID.String()},
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: domain.ProductStatusPublished,
				}, nil)
				wr.EXPECT().Add(gomock.Any(), gomock.Any()).Return(false, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wishlistRepo, productRepo, _, uc := setupWishlistUseCase(t)
			tc.setup(wishlistRepo, productRepo)

			resp, err := uc.AddItem(context.Background(), userID, tc.req)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if errors.Is(tc.wantErr, ErrProductNotFound) || errors.Is(tc.wantErr, ErrProductNotAvailable) {
					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected %v, got %v", tc.wantErr, err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.AlreadyExists != tc.wantDup {
				t.Fatalf("expected already_exists=%v, got %v", tc.wantDup, resp.AlreadyExists)
			}
			if !resp.IsInWishlist {
				t.Fatal("expected is_in_wishlist=true")
			}
		})
	}
}

func TestWishlistUseCaseListItems(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	productID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	vendorID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	errDB := errors.New("db error")

	tests := []struct {
		name      string
		params    domain.WishlistListParams
		setup     func(wr *mocks.MockWishlistRepository, sp *mocks.MockStorageProvider)
		wantErr   bool
		wantPage  int
		wantLimit int
	}{
		{
			name:   "success - image key converted to absolute URL",
			params: domain.WishlistListParams{Page: 1, Limit: 10},
			setup: func(wr *mocks.MockWishlistRepository, sp *mocks.MockStorageProvider) {
				key := "products/p1/main.jpg"
				abs := "https://cdn.example.com/products/p1/main.jpg"
				wr.EXPECT().ListByUser(gomock.Any(), userID, domain.WishlistListParams{Page: 1, Limit: 10}).Return([]domain.WishlistItemResponse{
					{
						ProductID:   productID,
						ProductName: "Kurma Ajwa",
						ProductSlug: "kurma-ajwa",
						VendorID:    vendorID,
						VendorName:  "Toko A",
						ImageURL:    &key,
						IsAvailable: true,
					},
				}, int64(1), nil)
				sp.EXPECT().GetURL(key).Return(abs)
			},
			wantPage:  1,
			wantLimit: 10,
		},
		{
			name:   "success - default page/limit applied",
			params: domain.WishlistListParams{Page: 0, Limit: 101},
			setup: func(wr *mocks.MockWishlistRepository, sp *mocks.MockStorageProvider) {
				wr.EXPECT().ListByUser(gomock.Any(), userID, domain.WishlistListParams{Page: 1, Limit: 10}).Return([]domain.WishlistItemResponse{}, int64(0), nil)
			},
			wantPage:  1,
			wantLimit: 10,
		},
		{
			name:   "error - repository failure",
			params: domain.WishlistListParams{Page: 1, Limit: 10},
			setup: func(wr *mocks.MockWishlistRepository, sp *mocks.MockStorageProvider) {
				wr.EXPECT().ListByUser(gomock.Any(), userID, domain.WishlistListParams{Page: 1, Limit: 10}).Return(nil, int64(0), errDB)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wishlistRepo, _, storageProvider, uc := setupWishlistUseCase(t)
			tc.setup(wishlistRepo, storageProvider)

			items, meta, err := uc.ListItems(context.Background(), userID, tc.params)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if meta.Page != tc.wantPage {
				t.Fatalf("expected page=%d, got %d", tc.wantPage, meta.Page)
			}
			if meta.Limit != tc.wantLimit {
				t.Fatalf("expected limit=%d, got %d", tc.wantLimit, meta.Limit)
			}
			if len(items) > 0 && items[0].ImageURL != nil && *items[0].ImageURL == "products/p1/main.jpg" {
				t.Fatal("expected absolute image URL conversion")
			}
		})
	}
}

func TestWishlistUseCaseRemoveItem(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	productID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	errDB := errors.New("db error")

	tests := []struct {
		name    string
		setup   func(wr *mocks.MockWishlistRepository)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(wr *mocks.MockWishlistRepository) {
				wr.EXPECT().DeleteByUserAndProduct(gomock.Any(), userID, productID).Return(nil)
			},
		},
		{
			name: "error - repository failure",
			setup: func(wr *mocks.MockWishlistRepository) {
				wr.EXPECT().DeleteByUserAndProduct(gomock.Any(), userID, productID).Return(errDB)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wishlistRepo, _, _, uc := setupWishlistUseCase(t)
			tc.setup(wishlistRepo)

			err := uc.RemoveItem(context.Background(), userID, productID)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestWishlistUseCaseGetProductStatus(t *testing.T) {
	userID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	productID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	errDB := errors.New("db error")

	tests := []struct {
		name          string
		setup         func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository)
		wantErr       error
		wantInWish    bool
		wantAvailable bool
	}{
		{
			name: "success - published and in wishlist",
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: domain.ProductStatusPublished,
				}, nil)
				wr.EXPECT().Exists(gomock.Any(), userID, productID).Return(true, nil)
			},
			wantInWish:    true,
			wantAvailable: true,
		},
		{
			name: "success - unavailable and not in wishlist",
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: "blocked",
				}, nil)
				wr.EXPECT().Exists(gomock.Any(), userID, productID).Return(false, nil)
			},
			wantInWish:    false,
			wantAvailable: false,
		},
		{
			name: "error - product not found",
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, nil)
			},
			wantErr: ErrProductNotFound,
		},
		{
			name: "error - product repo failure",
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - wishlist repo failure",
			setup: func(wr *mocks.MockWishlistRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(&domain.Product{
					ID:     productID,
					Status: domain.ProductStatusPublished,
				}, nil)
				wr.EXPECT().Exists(gomock.Any(), userID, productID).Return(false, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wishlistRepo, productRepo, _, uc := setupWishlistUseCase(t)
			tc.setup(wishlistRepo, productRepo)

			resp, err := uc.GetProductStatus(context.Background(), userID, productID)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if errors.Is(tc.wantErr, ErrProductNotFound) && !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.IsInWishlist != tc.wantInWish {
				t.Fatalf("expected is_in_wishlist=%v, got %v", tc.wantInWish, resp.IsInWishlist)
			}
			if resp.IsAvailable != tc.wantAvailable {
				t.Fatalf("expected is_available=%v, got %v", tc.wantAvailable, resp.IsAvailable)
			}
		})
	}
}
