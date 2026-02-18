package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// ---------- helpers ----------

func setupCartUseCase(t *testing.T) (
	*mocks.MockCartRepository,
	*mocks.MockProductRepository,
	domain.CartUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	cartRepo := mocks.NewMockCartRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	uc := NewCartUseCase(cartRepo, productRepo)
	return cartRepo, productRepo, uc
}

// fixtureIDs returns a reusable set of deterministic UUIDs for test fixtures.
func fixtureIDs() (userID, cartID, variantID, productID, itemID uuid.UUID) {
	userID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	cartID = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	variantID = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	productID = uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	itemID = uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	return
}

func activeCart(cartID, userID uuid.UUID) *domain.Cart {
	now := time.Now()
	return &domain.Cart{
		ID:        cartID,
		UserID:    userID,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func activeVariant(variantID, productID uuid.UUID) *domain.ProductVariant {
	return &domain.ProductVariant{
		ID:          variantID,
		ProductID:   productID,
		SKU:         "SKU-001",
		VariantName: "Default",
		Price:       50000,
		Currency:    "IDR",
		StockOnHand: 10,
		IsDefault:   true,
		IsActive:    true,
	}
}

func publishedProduct(productID uuid.UUID) *domain.Product {
	return &domain.Product{
		ID:     productID,
		Name:   "Kurma Ajwa",
		Slug:   "kurma-ajwa",
		Status: "published",
	}
}

var errDB = errors.New("db error")

// ============================================================
// GetCart
// ============================================================

func TestGetCart(t *testing.T) {
	userID, cartID, variantID, productID, itemID := fixtureIDs()
	now := time.Now()

	imgURL := "https://cdn.example.com/img.jpg"

	tests := []struct {
		name      string
		setup     func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository)
		wantErr   bool
		wantItems int
		wantTotal float64
	}{
		{
			name: "success - cart with one available item",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2, CreatedAt: now},
				}, nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				pr.EXPECT().FindImagesByProductID(gomock.Any(), productID).Return([]domain.ProductImage{
					{ID: uuid.New(), ProductID: productID, ImageURL: imgURL, IsPrimary: true},
				}, nil)
			},
			wantErr:   false,
			wantItems: 1,
			wantTotal: 100000, // 50000 * 2
		},
		{
			name: "success - empty cart (no items)",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{}, nil)
			},
			wantErr:   false,
			wantItems: 0,
			wantTotal: 0,
		},
		{
			name: "success - cart auto-created when not found",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, nil)
				cr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), gomock.Any()).Return([]domain.CartItem{}, nil)
			},
			wantErr:   false,
			wantItems: 0,
			wantTotal: 0,
		},
		{
			name: "success - skips item with deleted variant",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1, CreatedAt: now},
				}, nil)
				// Variant not found (deleted).
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(nil, nil)
			},
			wantErr:   false,
			wantItems: 0,
			wantTotal: 0,
		},
		{
			name: "success - skips item with deleted product",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1, CreatedAt: now},
				}, nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, nil) // product deleted
			},
			wantErr:   false,
			wantItems: 0,
			wantTotal: 0,
		},
		{
			name: "success - unavailable item not counted in total",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				inactiveVariant := activeVariant(variantID, productID)
				inactiveVariant.IsActive = false

				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2, CreatedAt: now},
				}, nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(inactiveVariant, nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				pr.EXPECT().FindImagesByProductID(gomock.Any(), productID).Return([]domain.ProductImage{}, nil)
			},
			wantErr:   false,
			wantItems: 1,
			wantTotal: 0, // inactive variant → not counted
		},
		{
			name: "error - FindByUserID fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: true,
		},
		{
			name: "error - Create cart fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, nil)
				cr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errDB)
			},
			wantErr: true,
		},
		{
			name: "error - FindItemsByCartID fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(nil, errDB)
			},
			wantErr: true,
		},
		{
			name: "error - FindVariantByID fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1, CreatedAt: now},
				}, nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(nil, errDB)
			},
			wantErr: true,
		},
		{
			name: "error - FindByID (product) fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1, CreatedAt: now},
				}, nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, errDB)
			},
			wantErr: true,
		},
		{
			name: "error - FindImagesByProductID fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1, CreatedAt: now},
				}, nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				pr.EXPECT().FindImagesByProductID(gomock.Any(), productID).Return(nil, errDB)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, uc := setupCartUseCase(t)
			tc.setup(cartRepo, productRepo)

			resp, err := uc.GetCart(context.Background(), userID)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(resp.Items) != tc.wantItems {
				t.Errorf("expected %d items, got %d", tc.wantItems, len(resp.Items))
			}
			if resp.Total != tc.wantTotal {
				t.Errorf("expected total %.2f, got %.2f", tc.wantTotal, resp.Total)
			}
			if resp.Currency != "IDR" {
				t.Errorf("expected currency IDR, got %s", resp.Currency)
			}
		})
	}
}

// ============================================================
// AddItem
// ============================================================

func TestAddItem(t *testing.T) {
	userID, cartID, variantID, productID, itemID := fixtureIDs()

	tests := []struct {
		name    string
		req     domain.AddCartItemRequest
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository)
		wantErr error
	}{
		{
			name: "success - new item added",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 2},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), cartID, variantID).Return(nil, nil)
				cr.EXPECT().UpsertItem(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "success - existing item incremented",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 3},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), cartID, variantID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 5}, nil,
				)
				cr.EXPECT().UpsertItem(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "success - cart auto-created",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, nil)
				cr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), gomock.Any(), variantID).Return(nil, nil)
				cr.EXPECT().UpsertItem(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "error - invalid variant UUID",
			req:     domain.AddCartItemRequest{ProductVariantID: "not-a-uuid", Qty: 1},
			setup:   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {},
			wantErr: ErrVariantNotFound,
		},
		{
			name: "error - variant not found",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(nil, nil)
			},
			wantErr: ErrVariantNotFound,
		},
		{
			name: "error - variant not active",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				v := activeVariant(variantID, productID)
				v.IsActive = false
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(v, nil)
			},
			wantErr: ErrVariantNotActive,
		},
		{
			name: "error - product not found",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, nil)
			},
			wantErr: ErrProductNotFound,
		},
		{
			name: "error - product not published",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				p := publishedProduct(productID)
				p.Status = "draft"
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(p, nil)
			},
			wantErr: ErrProductNotAvailable,
		},
		{
			name: "error - insufficient stock (new item)",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 11},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil) // stock = 10
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), cartID, variantID).Return(nil, nil)
			},
			wantErr: ErrInsufficientStock,
		},
		{
			name: "error - insufficient stock (existing item qty + new qty exceeds stock)",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 3},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil) // stock = 10
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), cartID, variantID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 8}, nil,
				) // 8 + 3 = 11 > 10
			},
			wantErr: ErrInsufficientStock,
		},
		{
			name: "error - FindVariantByID DB error",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - FindByID (product) DB error",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - getOrCreateCart fails",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - FindItemByCartAndVariant DB error",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), cartID, variantID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - UpsertItem DB error",
			req:  domain.AddCartItemRequest{ProductVariantID: variantID.String(), Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(publishedProduct(productID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByCartAndVariant(gomock.Any(), cartID, variantID).Return(nil, nil)
				cr.EXPECT().UpsertItem(gomock.Any(), gomock.Any()).Return(errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, uc := setupCartUseCase(t)
			tc.setup(cartRepo, productRepo)

			resp, err := uc.AddItem(context.Background(), userID, tc.req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.CartID == uuid.Nil {
				t.Error("expected non-nil cart ID")
			}
			if resp.ItemID == uuid.Nil {
				t.Error("expected non-nil item ID")
			}
		})
	}
}

// ============================================================
// UpdateItem
// ============================================================

func TestUpdateItem(t *testing.T) {
	userID, cartID, variantID, productID, itemID := fixtureIDs()
	otherCartID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")

	tests := []struct {
		name    string
		itemID  uuid.UUID
		req     domain.UpdateCartItemRequest
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository)
		wantErr error
	}{
		{
			name:   "success",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 5},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2}, nil,
				)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil) // stock = 10
				cr.EXPECT().UpdateItemQty(gomock.Any(), itemID, 5).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "error - item not found",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(nil, nil)
			},
			wantErr: ErrCartItemNotFound,
		},
		{
			name:   "error - item not owned",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: otherCartID, ProductVariantID: variantID, Qty: 2}, nil,
				)
			},
			wantErr: ErrCartItemNotOwned,
		},
		{
			name:   "error - variant not found",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2}, nil,
				)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(nil, nil)
			},
			wantErr: ErrVariantNotFound,
		},
		{
			name:   "error - insufficient stock",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 15},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2}, nil,
				)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil) // stock = 10
			},
			wantErr: ErrInsufficientStock,
		},
		{
			name:   "error - getOrCreateCart fails",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "error - FindItemByID DB error",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "error - FindVariantByID DB error",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 1},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2}, nil,
				)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "error - UpdateItemQty fails",
			itemID: itemID,
			req:    domain.UpdateCartItemRequest{Qty: 3},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 2}, nil,
				)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(activeVariant(variantID, productID), nil)
				cr.EXPECT().UpdateItemQty(gomock.Any(), itemID, 3).Return(errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, uc := setupCartUseCase(t)
			tc.setup(cartRepo, productRepo)

			resp, err := uc.UpdateItem(context.Background(), userID, tc.itemID, tc.req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.CartID != cartID {
				t.Errorf("expected cart ID %s, got %s", cartID, resp.CartID)
			}
			if resp.ItemID != tc.itemID {
				t.Errorf("expected item ID %s, got %s", tc.itemID, resp.ItemID)
			}
		})
	}
}

// ============================================================
// RemoveItem
// ============================================================

func TestRemoveItem(t *testing.T) {
	userID, cartID, _, _, itemID := fixtureIDs()
	otherCartID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	variantID := uuid.New()

	tests := []struct {
		name    string
		itemID  uuid.UUID
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository)
		wantErr error
	}{
		{
			name:   "success",
			itemID: itemID,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1}, nil,
				)
				cr.EXPECT().DeleteItem(gomock.Any(), itemID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "error - item not found",
			itemID: itemID,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(nil, nil)
			},
			wantErr: ErrCartItemNotFound,
		},
		{
			name:   "error - item not owned",
			itemID: itemID,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: otherCartID, ProductVariantID: variantID, Qty: 1}, nil,
				)
			},
			wantErr: ErrCartItemNotOwned,
		},
		{
			name:   "error - getOrCreateCart fails",
			itemID: itemID,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "error - FindItemByID DB error",
			itemID: itemID,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name:   "error - DeleteItem fails",
			itemID: itemID,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemByID(gomock.Any(), itemID).Return(
					&domain.CartItem{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1}, nil,
				)
				cr.EXPECT().DeleteItem(gomock.Any(), itemID).Return(errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, uc := setupCartUseCase(t)
			tc.setup(cartRepo, productRepo)

			err := uc.RemoveItem(context.Background(), userID, tc.itemID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

// ============================================================
// ClearCart
// ============================================================

func TestClearCart(t *testing.T) {
	userID, cartID, variantID, _, itemID := fixtureIDs()

	tests := []struct {
		name    string
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository)
		wantErr error
	}{
		{
			name: "success - clears items",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 3},
				}, nil)
				cr.EXPECT().ClearItems(gomock.Any(), cartID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "success - no cart exists (noop)",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, nil)
			},
			wantErr: nil,
		},
		{
			name: "success - cart exists but empty (noop)",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{}, nil)
			},
			wantErr: nil,
		},
		{
			name: "error - FindByUserID fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - FindItemsByCartID fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "error - ClearItems fails",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(activeCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{
					{ID: itemID, CartID: cartID, ProductVariantID: variantID, Qty: 1},
				}, nil)
				cr.EXPECT().ClearItems(gomock.Any(), cartID).Return(errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, uc := setupCartUseCase(t)
			tc.setup(cartRepo, productRepo)

			err := uc.ClearCart(context.Background(), userID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
