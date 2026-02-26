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

// --- helpers & fixtures ---

func setupCheckoutUseCase(t *testing.T) (
	*mocks.MockCartRepository,
	*mocks.MockProductRepository,
	*mocks.MockVendorRepository,
	*mocks.MockOrderRepository,
	*mocks.MockPaymentRepository,
	*mocks.MockLedgerRepository,
	*mocks.MockUserRepository,
	*mocks.MockXenditInvoiceProvider,
	domain.CheckoutUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	cartRepo := mocks.NewMockCartRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	orderRepo := mocks.NewMockOrderRepository(ctrl)
	paymentRepo := mocks.NewMockPaymentRepository(ctrl)
	ledgerRepo := mocks.NewMockLedgerRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	xenditInvoice := mocks.NewMockXenditInvoiceProvider(ctrl)

	uc := NewCheckoutUseCase(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, xenditInvoice, "https://example.com", "https://test.example.com/api/v1/webhooks/xendit/invoice")
	return cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, xenditInvoice, uc
}

func fixtureCheckoutIDs() (userID, cartID, variantID, productID, vendorID uuid.UUID) {
	userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cartID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	variantID = uuid.MustParse("00000000-0000-0000-0000-000000000003")
	productID = uuid.MustParse("00000000-0000-0000-0000-000000000004")
	vendorID = uuid.MustParse("00000000-0000-0000-0000-000000000005")
	return
}

func fixtureCart(cartID, userID uuid.UUID) *domain.Cart {
	return &domain.Cart{
		ID:     cartID,
		UserID: userID,
		Status: domain.CartStatusActive,
	}
}

func fixtureCartItems(cartID, variantID uuid.UUID) []domain.CartItem {
	return []domain.CartItem{
		{
			ID:               uuid.New(),
			CartID:           cartID,
			ProductVariantID: variantID,
			Qty:              2,
		},
	}
}

func fixtureUser(userID uuid.UUID) *domain.User {
	return &domain.User{
		ID:       userID,
		Email:    "test@example.com",
		FullName: "Test User",
	}
}

func fixtureVariant(variantID, productID uuid.UUID) *domain.ProductVariant {
	return &domain.ProductVariant{
		ID:          variantID,
		ProductID:   productID,
		SKU:         "SKU-001",
		VariantName: "Default",
		Price:       100000,
		StockOnHand: 10,
		IsActive:    true,
	}
}

func fixtureProduct(productID, vendorID uuid.UUID) *domain.Product {
	return &domain.Product{
		ID:       productID,
		VendorID: vendorID,
		Name:     "Test Product",
		Status:   domain.ProductStatusPublished,
	}
}

func fixtureVendor(vendorID uuid.UUID) *domain.Vendor {
	xenditAccID := "xa-vendor-123"
	return &domain.Vendor{
		ID:              vendorID,
		DisplayName:     "Test Vendor",
		XenditAccountID: &xenditAccID,
	}
}

func fixtureXenditInvoiceResponse() *domain.XenditInvoiceResponse {
	return &domain.XenditInvoiceResponse{
		ID:         "xinv-001",
		InvoiceURL: "https://checkout.xendit.co/invoice-001",
		ExpiryDate: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
}

// --- Checkout tests ---

func TestCheckout(t *testing.T) {
	userID, cartID, variantID, productID, vendorID := fixtureCheckoutIDs()

	tests := []struct {
		name    string
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider)
		wantErr error
	}{
		{
			name: "success",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				xi.EXPECT().CreateInvoice(gomock.Any(), "xa-vendor-123", gomock.Any()).Return(fixtureXenditInvoiceResponse(), nil)
				pmr.EXPECT().CreateInvoice(gomock.Any(), gomock.Any()).Return(nil)
				cr.EXPECT().UpdateStatus(gomock.Any(), cartID, domain.CartStatusConverted).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "empty cart - no cart found",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, nil)
			},
			wantErr: ErrCartEmpty,
		},
		{
			name: "empty cart - no items",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{}, nil)
			},
			wantErr: ErrCartEmpty,
		},
		{
			name: "variant not active",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)

				inactiveVariant := fixtureVariant(variantID, productID)
				inactiveVariant.IsActive = false
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(inactiveVariant, nil)
			},
			wantErr: ErrCartHasUnavailableItems,
		},
		{
			name: "product not published",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)

				unpublished := fixtureProduct(productID, vendorID)
				unpublished.Status = "draft"
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(unpublished, nil)
			},
			wantErr: ErrCartHasUnavailableItems,
		},
		{
			name: "insufficient stock",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)

				lowStock := fixtureVariant(variantID, productID)
				lowStock.StockOnHand = 1 // cart has qty=2
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(lowStock, nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
			},
			wantErr: ErrCheckoutStockInsufficient,
		},
		{
			name: "vendor has no xendit account",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)

				noXendit := fixtureVendor(vendorID)
				noXendit.XenditAccountID = nil
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(noXendit, nil)
			},
			wantErr: ErrVendorNoXenditAccount,
		},
		{
			name: "xendit API error",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				xi.EXPECT().CreateInvoice(gomock.Any(), "xa-vendor-123", gomock.Any()).Return(nil, errors.New("xendit 500"))
			},
			wantErr: ErrInvoiceCreationFailed,
		},
		{
			name: "find cart DB error",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, xenditInvoice, uc := setupCheckoutUseCase(t)
			tc.setup(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, xenditInvoice)

			resp, err := uc.Checkout(context.Background(), userID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp == nil || len(resp.Orders) == 0 {
				t.Fatal("expected non-empty checkout response")
			}
			order := resp.Orders[0]
			if order.VendorID != vendorID {
				t.Errorf("expected vendor %s, got %s", vendorID, order.VendorID)
			}
			// subtotal = 100_000 * 2 = 200_000
			if order.Subtotal != 200000 {
				t.Errorf("expected subtotal 200000, got %f", order.Subtotal)
			}
			// grand_total = subtotal + platformFeeTotal(6000)
			if order.GrandTotal != 206000 {
				t.Errorf("expected grand total 206000, got %f", order.GrandTotal)
			}
			if order.PlatformFee != 6000 {
				t.Errorf("expected platform fee 6000, got %f", order.PlatformFee)
			}
			if order.InvoiceURL == "" {
				t.Error("expected non-empty invoice URL")
			}
		})
	}
}

// --- HandleWebhook tests ---

func TestHandleWebhook(t *testing.T) {
	invoiceID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000011")

	paidAt := time.Now().Format(time.RFC3339)
	paidPayload := domain.XenditWebhookPayload{
		ID:             "xinv-001",
		ExternalID:     "INV-ORD-20250101-ABCD1234-abcd1234",
		Status:         "PAID",
		Amount:         206000,
		PaidAmount:     206000,
		PaidAt:         &paidAt,
		PaymentMethod:  "BANK_TRANSFER",
		PaymentChannel: "BCA",
	}

	expiredPayload := domain.XenditWebhookPayload{
		ID:         "xinv-001",
		ExternalID: "INV-ORD-20250101-ABCD1234-abcd1234",
		Status:     "EXPIRED",
		Amount:     206000,
	}

	invoice := &domain.PaymentInvoice{
		ID:                invoiceID,
		OrderID:           orderID,
		ExternalInvoiceID: "INV-ORD-20250101-ABCD1234-abcd1234",
		Amount:            206000,
		Status:            domain.InvoiceStatusPending,
	}

	order := &domain.Order{
		ID:       orderID,
		OrderNo:  "ORD-20250101-ABCD1234",
		Subtotal: 200000,
	}

	acctReceivable := &domain.LedgerAccount{ID: uuid.New(), Code: "1100"}
	acctVendorPayable := &domain.LedgerAccount{ID: uuid.New(), Code: "2100"}
	acctPlatformFee := &domain.LedgerAccount{ID: uuid.New(), Code: "4100"}
	acctAdminFee := &domain.LedgerAccount{ID: uuid.New(), Code: "4200"}

	tests := []struct {
		name    string
		payload domain.XenditWebhookPayload
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider)
		wantErr error
	}{
		{
			name:    "PAID - success",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), paidPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				// handlePaid
				pmr.EXPECT().UpdateInvoiceStatus(gomock.Any(), invoiceID, domain.InvoiceStatusPaid, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().UpdateOrderStatus(gomock.Any(), orderID, domain.OrderStatusPaid, domain.PaymentStatusPaid, gomock.Any(), gomock.Any()).Return(nil)
				// recordPaymentLedger
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(order, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "1100").Return(acctReceivable, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "2100").Return(acctVendorPayable, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "4100").Return(acctPlatformFee, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "4200").Return(acctAdminFee, nil)
				lr.EXPECT().CreateJournalWithLines(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "EXPIRED - success",
			payload: expiredPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:EXPIRED").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), expiredPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				// handleExpired
				pmr.EXPECT().UpdateInvoiceStatus(gomock.Any(), invoiceID, domain.InvoiceStatusExpired, nil, nil, nil, nil).Return(nil)
				or.EXPECT().UpdateOrderStatus(gomock.Any(), orderID, domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().RestoreStock(gomock.Any(), orderID).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "idempotent - already processed",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(true, nil)
			},
			wantErr: nil,
		},
		{
			name:    "invoice not found",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), paidPayload.ExternalID).Return(nil, nil)
			},
			wantErr: ErrInvoiceNotFound,
		},
		{
			name:    "event existence check DB error",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "unknown status - ignored",
			payload: domain.XenditWebhookPayload{
				ID:         "xinv-001",
				ExternalID: "INV-ORD-20250101-ABCD1234-abcd1234",
				Status:     "PENDING",
				Amount:     206000,
			},
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, xi *mocks.MockXenditInvoiceProvider) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PENDING").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), "INV-ORD-20250101-ABCD1234-abcd1234").Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, xenditInvoice, uc := setupCheckoutUseCase(t)
			tc.setup(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, xenditInvoice)

			err := uc.HandleWebhook(context.Background(), tc.payload)

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
