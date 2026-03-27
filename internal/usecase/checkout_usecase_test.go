package usecase

import (
	"context"
	"errors"
	"fmt"
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
	*mocks.MockAddressRepository,
	*mocks.MockVendorCourierRepository,
	*mocks.MockRajaOngkirProvider,
	*mocks.MockStorageProvider,
	*mocks.MockXenditInvoiceProvider,
	*mocks.MockShipmentRepository,
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
	addressRepo := mocks.NewMockAddressRepository(ctrl)
	vendorCourierRepo := mocks.NewMockVendorCourierRepository(ctrl)
	shipmentRepo := mocks.NewMockShipmentRepository(ctrl)
	rajaOngkir := mocks.NewMockRajaOngkirProvider(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	xenditInvoice := mocks.NewMockXenditInvoiceProvider(ctrl)

	uc := NewCheckoutUseCase(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, shipmentRepo, rajaOngkir, storage, xenditInvoice, "https://example.com", "https://test.example.com/api/v1/webhooks/xendit/invoice")
	return cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, rajaOngkir, storage, xenditInvoice, shipmentRepo, uc
}

func fixtureCheckoutIDs() (userID, cartID, variantID, productID, vendorID, addressID uuid.UUID) {
	userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	cartID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	variantID = uuid.MustParse("00000000-0000-0000-0000-000000000003")
	productID = uuid.MustParse("00000000-0000-0000-0000-000000000004")
	vendorID = uuid.MustParse("00000000-0000-0000-0000-000000000005")
	addressID = uuid.MustParse("00000000-0000-0000-0000-000000000006")
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
		OwnerUserID:     uuid.MustParse("00000000-0000-0000-0000-000000000010"),
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

func fixtureAddress(addressID, userID uuid.UUID) *domain.Address {
	districtID := "1376"
	return &domain.Address{
		ID:          addressID,
		UserID:      userID,
		DistrictID:  &districtID,
		AddressLine: "Jl. Test No. 1",
		IsDefault:   true,
	}
}

func fixtureCheckoutRequest(addressID, vendorID uuid.UUID) domain.CheckoutRequest {
	return domain.CheckoutRequest{
		AddressID: addressID.String(),
		ShippingChoices: []domain.CheckoutShippingChoice{
			{
				VendorID:    vendorID.String(),
				CourierCode: "jne",
				Service:     "REG",
			},
		},
	}
}

func fixtureVendorWarehouseAddress(vendorOwnerUserID uuid.UUID) *domain.Address {
	districtID := "5678"
	return &domain.Address{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000077"),
		UserID:      vendorOwnerUserID,
		DistrictID:  &districtID,
		AddressLine: "Jl. Vendor Warehouse No. 1",
		IsDefault:   true,
	}
}

func fixtureVendorCouriers() []domain.VendorCourierResponse {
	return []domain.VendorCourierResponse{
		{
			ID:          uuid.MustParse("00000000-0000-0000-0000-000000000088"),
			CourierCode: "jne",
			CourierName: "JNE",
			IsActive:    true,
		},
	}
}

func fixtureShippingOptions() []domain.ShippingCostOption {
	return []domain.ShippingCostOption{
		{
			Name:        "Jalur Nugraha Ekakurir (JNE)",
			Code:        "jne",
			Service:     "REG",
			Description: "Layanan Reguler",
			Cost:        11000,
			ETD:         "2-3",
		},
	}
}

// --- Checkout tests ---

func TestCheckout(t *testing.T) {
	userID, cartID, variantID, productID, vendorID, addressID := fixtureCheckoutIDs()
	req := fixtureCheckoutRequest(addressID, vendorID)

	tests := []struct {
		name    string
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository)
		wantErr error
	}{
		{
			name: "success",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				ar.EXPECT().FindDefaultByUserID(gomock.Any(), fixtureVendor(vendorID).OwnerUserID).Return(fixtureVendorWarehouseAddress(fixtureVendor(vendorID).OwnerUserID), nil)
				vcr.EXPECT().FindByVendorID(gomock.Any(), vendorID).Return(fixtureVendorCouriers(), nil)
				ro.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 1000, "jne").Return(fixtureShippingOptions(), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				xi.EXPECT().CreateInvoice(gomock.Any(), "xa-vendor-123", gomock.Any()).Return(fixtureXenditInvoiceResponse(), nil)
				pmr.EXPECT().CreateInvoice(gomock.Any(), gomock.Any()).Return(nil)
				cr.EXPECT().UpdateStatus(gomock.Any(), cartID, domain.CartStatusConverted).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "address not found",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(nil, nil)
			},
			wantErr: ErrAddressNotFound,
		},
		{
			name: "address not owned",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				otherUserID := uuid.MustParse("00000000-0000-0000-0000-000000000099")
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, otherUserID), nil)
			},
			wantErr: ErrAddressNotOwned,
		},
		{
			name: "empty cart - no cart found",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, nil)
			},
			wantErr: ErrCartEmpty,
		},
		{
			name: "empty cart - no items",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return([]domain.CartItem{}, nil)
			},
			wantErr: ErrCartEmpty,
		},
		{
			name: "variant not active",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
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
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
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
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
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
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
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
			name: "final stock conflict maps to checkout insufficient stock",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				ar.EXPECT().FindDefaultByUserID(gomock.Any(), fixtureVendor(vendorID).OwnerUserID).Return(fixtureVendorWarehouseAddress(fixtureVendor(vendorID).OwnerUserID), nil)
				vcr.EXPECT().FindByVendorID(gomock.Any(), vendorID).Return(fixtureVendorCouriers(), nil)
				ro.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 1000, "jne").Return(fixtureShippingOptions(), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("%w: variant %s", domain.ErrStockUnavailable, variantID))
			},
			wantErr: ErrCheckoutStockInsufficient,
		},
		{
			name: "xendit API error",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				ar.EXPECT().FindDefaultByUserID(gomock.Any(), fixtureVendor(vendorID).OwnerUserID).Return(fixtureVendorWarehouseAddress(fixtureVendor(vendorID).OwnerUserID), nil)
				vcr.EXPECT().FindByVendorID(gomock.Any(), vendorID).Return(fixtureVendorCouriers(), nil)
				ro.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 1000, "jne").Return(fixtureShippingOptions(), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				xi.EXPECT().CreateInvoice(gomock.Any(), "xa-vendor-123", gomock.Any()).Return(nil, errors.New("xendit 500"))
				or.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().RestoreStock(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: ErrInvoiceCreationFailed,
		},
		{
			name: "payment invoice save error triggers compensation",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				ar.EXPECT().FindDefaultByUserID(gomock.Any(), fixtureVendor(vendorID).OwnerUserID).Return(fixtureVendorWarehouseAddress(fixtureVendor(vendorID).OwnerUserID), nil)
				vcr.EXPECT().FindByVendorID(gomock.Any(), vendorID).Return(fixtureVendorCouriers(), nil)
				ro.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 1000, "jne").Return(fixtureShippingOptions(), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				xi.EXPECT().CreateInvoice(gomock.Any(), "xa-vendor-123", gomock.Any()).Return(fixtureXenditInvoiceResponse(), nil)
				pmr.EXPECT().CreateInvoice(gomock.Any(), gomock.Any()).Return(errDB)
				xi.EXPECT().ExpireInvoice(gomock.Any(), "xa-vendor-123", "xinv-001").Return(nil)
				or.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().RestoreStock(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: errDB,
		},
		{
			name: "compensation failure returns checkout compensation failed",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
				cr.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(fixtureCartItems(cartID, variantID), nil)
				ur.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)
				pr.EXPECT().FindVariantByID(gomock.Any(), variantID).Return(fixtureVariant(variantID, productID), nil)
				pr.EXPECT().FindByID(gomock.Any(), productID).Return(fixtureProduct(productID, vendorID), nil)
				vr.EXPECT().FindByID(gomock.Any(), vendorID).Return(fixtureVendor(vendorID), nil)
				ar.EXPECT().FindDefaultByUserID(gomock.Any(), fixtureVendor(vendorID).OwnerUserID).Return(fixtureVendorWarehouseAddress(fixtureVendor(vendorID).OwnerUserID), nil)
				vcr.EXPECT().FindByVendorID(gomock.Any(), vendorID).Return(fixtureVendorCouriers(), nil)
				ro.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 1000, "jne").Return(fixtureShippingOptions(), nil)
				or.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				sr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				xi.EXPECT().CreateInvoice(gomock.Any(), "xa-vendor-123", gomock.Any()).Return(nil, errors.New("xendit 500"))
				or.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, gomock.Any(), gomock.Any()).Return(errDB)
				or.EXPECT().RestoreStock(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: ErrCheckoutCompensationFailed,
		},
		{
			name: "find cart DB error",
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				ar.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
				cr.EXPECT().FindByUserID(gomock.Any(), userID).Return(nil, errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, rajaOngkir, storage, xenditInvoice, shipmentRepo, uc := setupCheckoutUseCase(t)
			tc.setup(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, rajaOngkir, storage, xenditInvoice, shipmentRepo)

			resp, err := uc.Checkout(context.Background(), userID, req)

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
			// grand_total = subtotal + shippingFee(11000) + platformFeeTotal(6000) = 217000
			if order.GrandTotal != 217000 {
				t.Errorf("expected grand total 217000, got %f", order.GrandTotal)
			}
			if order.ShippingFee != 11000 {
				t.Errorf("expected shipping fee 11000, got %f", order.ShippingFee)
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

func TestCheckoutCompensatesPreviousVendorWhenLaterVendorRunsOutOfStock(t *testing.T) {
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000101")
	cartID := uuid.MustParse("00000000-0000-0000-0000-000000000102")
	addressID := uuid.MustParse("00000000-0000-0000-0000-000000000103")
	variantID1 := uuid.MustParse("00000000-0000-0000-0000-000000000104")
	productID1 := uuid.MustParse("00000000-0000-0000-0000-000000000105")
	vendorID1 := uuid.MustParse("00000000-0000-0000-0000-000000000106")
	variantID2 := uuid.MustParse("00000000-0000-0000-0000-000000000107")
	productID2 := uuid.MustParse("00000000-0000-0000-0000-000000000108")
	vendorID2 := uuid.MustParse("00000000-0000-0000-0000-000000000109")

	cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, rajaOngkir, storage, xenditInvoice, shipmentRepo, uc := setupCheckoutUseCase(t)

	cartItems := []domain.CartItem{
		{
			ID:               uuid.New(),
			CartID:           cartID,
			ProductVariantID: variantID1,
			Qty:              1,
		},
		{
			ID:               uuid.New(),
			CartID:           cartID,
			ProductVariantID: variantID2,
			Qty:              1,
		},
	}

	vendor1 := fixtureVendor(vendorID1)
	vendor1Account := "xa-vendor-1"
	vendor1.XenditAccountID = &vendor1Account
	vendor1.DisplayName = "Vendor 1"

	vendor2 := fixtureVendor(vendorID2)
	vendor2Account := "xa-vendor-2"
	vendor2.XenditAccountID = &vendor2Account
	vendor2.DisplayName = "Vendor 2"

	req := domain.CheckoutRequest{
		AddressID: addressID.String(),
		ShippingChoices: []domain.CheckoutShippingChoice{
			{VendorID: vendorID1.String(), CourierCode: "jne", Service: "REG"},
			{VendorID: vendorID2.String(), CourierCode: "jne", Service: "REG"},
		},
	}

	addressRepo.EXPECT().FindByID(gomock.Any(), addressID).Return(fixtureAddress(addressID, userID), nil)
	cartRepo.EXPECT().FindByUserID(gomock.Any(), userID).Return(fixtureCart(cartID, userID), nil)
	cartRepo.EXPECT().FindItemsByCartID(gomock.Any(), cartID).Return(cartItems, nil)
	userRepo.EXPECT().FindByID(gomock.Any(), userID).Return(fixtureUser(userID), nil)

	productRepo.EXPECT().FindVariantByID(gomock.Any(), variantID1).Return(fixtureVariant(variantID1, productID1), nil)
	productRepo.EXPECT().FindByID(gomock.Any(), productID1).Return(fixtureProduct(productID1, vendorID1), nil)
	productRepo.EXPECT().FindVariantByID(gomock.Any(), variantID2).Return(fixtureVariant(variantID2, productID2), nil)
	productRepo.EXPECT().FindByID(gomock.Any(), productID2).Return(fixtureProduct(productID2, vendorID2), nil)

	vendorRepo.EXPECT().FindByID(gomock.Any(), vendorID1).Return(vendor1, nil)
	addressRepo.EXPECT().FindDefaultByUserID(gomock.Any(), vendor1.OwnerUserID).Return(fixtureVendorWarehouseAddress(vendor1.OwnerUserID), nil)
	vendorCourierRepo.EXPECT().FindByVendorID(gomock.Any(), vendorID1).Return(fixtureVendorCouriers(), nil)
	rajaOngkir.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 500, "jne").Return(fixtureShippingOptions(), nil)
	orderRepo.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	shipmentRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	xenditInvoice.EXPECT().CreateInvoice(gomock.Any(), vendor1Account, gomock.Any()).Return(fixtureXenditInvoiceResponse(), nil)
	paymentRepo.EXPECT().CreateInvoice(gomock.Any(), gomock.Any()).Return(nil)

	vendorRepo.EXPECT().FindByID(gomock.Any(), vendorID2).Return(vendor2, nil)
	addressRepo.EXPECT().FindDefaultByUserID(gomock.Any(), vendor2.OwnerUserID).Return(fixtureVendorWarehouseAddress(vendor2.OwnerUserID), nil)
	vendorCourierRepo.EXPECT().FindByVendorID(gomock.Any(), vendorID2).Return(fixtureVendorCouriers(), nil)
	rajaOngkir.EXPECT().CalculateDomesticCost(gomock.Any(), "5678", "1376", 500, "jne").Return(fixtureShippingOptions(), nil)
	orderRepo.EXPECT().CreateOrderWithItems(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fmt.Errorf("%w: variant %s", domain.ErrStockUnavailable, variantID2))

	xenditInvoice.EXPECT().ExpireInvoice(gomock.Any(), vendor1Account, "xinv-001").Return(nil)
	paymentRepo.EXPECT().UpdateInvoiceStatus(gomock.Any(), gomock.Any(), domain.InvoiceStatusFailed, nil, nil, nil, gomock.Any()).Return(nil)
	orderRepo.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, nil, gomock.Any()).Return(nil)
	orderRepo.EXPECT().RestoreStock(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := uc.Checkout(context.Background(), userID, req)

	if resp != nil {
		t.Fatalf("expected nil response, got %+v", resp)
	}
	if !errors.Is(err, ErrCheckoutStockInsufficient) {
		t.Fatalf("expected error %v, got %v", ErrCheckoutStockInsufficient, err)
	}

	_ = storage
	_ = ledgerRepo
}

// --- HandleWebhook tests ---

func TestHandleWebhook(t *testing.T) {
	invoiceID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	orderID := uuid.MustParse("00000000-0000-0000-0000-000000000011")
	vendorID := uuid.MustParse("00000000-0000-0000-0000-000000000012")

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
		ID:          orderID,
		VendorID:    vendorID,
		OrderNo:     "ORD-20250101-ABCD1234",
		Subtotal:    200000,
		OrderStatus: domain.OrderStatusPendingPayment,
	}
	canceledOrder := &domain.Order{
		ID:          orderID,
		VendorID:    vendorID,
		OrderNo:     "ORD-20250101-ABCD1234",
		Subtotal:    200000,
		OrderStatus: domain.OrderStatusCanceled,
	}
	paidOrder := &domain.Order{
		ID:          orderID,
		VendorID:    vendorID,
		OrderNo:     "ORD-20250101-ABCD1234",
		Subtotal:    200000,
		OrderStatus: domain.OrderStatusPaid,
	}

	acctReceivable := &domain.LedgerAccount{ID: uuid.New(), Code: "1100"}
	acctVendorPayable := &domain.LedgerAccount{ID: uuid.New(), Code: "2100"}
	acctPlatformFee := &domain.LedgerAccount{ID: uuid.New(), Code: "4100"}
	acctAdminFee := &domain.LedgerAccount{ID: uuid.New(), Code: "4200"}

	tests := []struct {
		name    string
		payload domain.XenditWebhookPayload
		setup   func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository)
		wantErr error
	}{
		{
			name:    "PAID - success",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), paidPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				// handlePaid
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(order, nil)
				pmr.EXPECT().UpdateInvoiceStatusIfCurrent(
					gomock.Any(),
					invoiceID,
					domain.InvoiceStatusPending,
					domain.InvoiceStatusPaid,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(true, nil)
				or.EXPECT().UpdateOrderStatus(gomock.Any(), orderID, domain.OrderStatusPaid, domain.PaymentStatusPaid, gomock.Any(), gomock.Any()).Return(nil)
				// recordPaymentLedger
				lr.EXPECT().FindAccountByCode(gomock.Any(), "1100").Return(acctReceivable, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "2100").Return(acctVendorPayable, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "4100").Return(acctPlatformFee, nil)
				lr.EXPECT().FindAccountByCode(gomock.Any(), "4200").Return(acctAdminFee, nil)
				lr.EXPECT().CreateJournalWithLines(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "PAID - invoice transition already applied",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), paidPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(order, nil)
				pmr.EXPECT().UpdateInvoiceStatusIfCurrent(
					gomock.Any(),
					invoiceID,
					domain.InvoiceStatusPending,
					domain.InvoiceStatusPaid,
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			name:    "PAID - ignored for canceled order",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), paidPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(canceledOrder, nil)
			},
			wantErr: nil,
		},
		{
			name:    "EXPIRED - success",
			payload: expiredPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:EXPIRED").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), expiredPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				// handleExpired
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(order, nil)
				or.EXPECT().ApplyExpiredWebhookUpdate(gomock.Any(), orderID, invoiceID, gomock.Any()).Return(true, nil)
			},
			wantErr: nil,
		},
		{
			name:    "EXPIRED - ignored for paid order",
			payload: expiredPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:EXPIRED").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), expiredPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(paidOrder, nil)
			},
			wantErr: nil,
		},
		{
			name:    "EXPIRED - transition already applied",
			payload: expiredPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:EXPIRED").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), expiredPayload.ExternalID).Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
				or.EXPECT().FindByID(gomock.Any(), orderID).Return(order, nil)
				or.EXPECT().ApplyExpiredWebhookUpdate(gomock.Any(), orderID, invoiceID, gomock.Any()).Return(false, nil)
			},
			wantErr: nil,
		},
		{
			name:    "idempotent - already processed",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(true, nil)
			},
			wantErr: nil,
		},
		{
			name:    "invoice not found",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PAID").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), paidPayload.ExternalID).Return(nil, nil)
			},
			wantErr: ErrInvoiceNotFound,
		},
		{
			name:    "event existence check DB error",
			payload: paidPayload,
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
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
			setup: func(cr *mocks.MockCartRepository, pr *mocks.MockProductRepository, vr *mocks.MockVendorRepository, or *mocks.MockOrderRepository, pmr *mocks.MockPaymentRepository, lr *mocks.MockLedgerRepository, ur *mocks.MockUserRepository, ar *mocks.MockAddressRepository, vcr *mocks.MockVendorCourierRepository, ro *mocks.MockRajaOngkirProvider, sp *mocks.MockStorageProvider, xi *mocks.MockXenditInvoiceProvider, sr *mocks.MockShipmentRepository) {
				pmr.EXPECT().EventExistsByExternalID(gomock.Any(), "xinv-001:PENDING").Return(false, nil)
				pmr.EXPECT().FindInvoiceByExternalID(gomock.Any(), "INV-ORD-20250101-ABCD1234-abcd1234").Return(invoice, nil)
				pmr.EXPECT().CreateEvent(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, rajaOngkir, storage, xenditInvoice, shipmentRepo, uc := setupCheckoutUseCase(t)
			tc.setup(cartRepo, productRepo, vendorRepo, orderRepo, paymentRepo, ledgerRepo, userRepo, addressRepo, vendorCourierRepo, rajaOngkir, storage, xenditInvoice, shipmentRepo)

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
