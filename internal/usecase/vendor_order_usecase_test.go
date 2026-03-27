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

func setupVendorOrderUseCase(t *testing.T) (
	*mocks.MockVendorOrderRepository,
	*mocks.MockOrderRepository,
	*mocks.MockShipmentRepository,
	*mocks.MockPaymentRepository,
	*mocks.MockUserRepository,
	*mocks.MockVendorRepository,
	*mocks.MockRajaOngkirProvider,
	domain.VendorOrderUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	vendorOrderRepo := mocks.NewMockVendorOrderRepository(ctrl)
	orderRepo := mocks.NewMockOrderRepository(ctrl)
	shipmentRepo := mocks.NewMockShipmentRepository(ctrl)
	paymentRepo := mocks.NewMockPaymentRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	rajaOngkir := mocks.NewMockRajaOngkirProvider(ctrl)
	uc := NewVendorOrderUseCase(vendorOrderRepo, orderRepo, shipmentRepo, paymentRepo, userRepo, vendorRepo, rajaOngkir)
	return vendorOrderRepo, orderRepo, shipmentRepo, paymentRepo, userRepo, vendorRepo, rajaOngkir, uc
}

// --- AcceptOrder tests ---

func TestVendorOrderUseCase_AcceptOrder_Success(t *testing.T) {
	vendorOrderRepo, orderRepo, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	order := &domain.Order{
		ID:            orderID,
		VendorID:      vendorID,
		OrderStatus:   domain.OrderStatusPaid,
		PaymentStatus: "paid",
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)
	orderRepo.EXPECT().UpdateOrderStatus(ctx, orderID, domain.OrderStatusProcessing, "paid", gomock.Any(), gomock.Any()).Return(nil)

	res, err := uc.AcceptOrder(ctx, vendorID, orderID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusProcessing {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusProcessing, res.OrderStatus)
	}
}

func TestVendorOrderUseCase_AcceptOrder_NotPaid(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	order := &domain.Order{
		ID:          orderID,
		VendorID:    vendorID,
		OrderStatus: domain.OrderStatusProcessing,
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)

	_, err := uc.AcceptOrder(ctx, vendorID, orderID, userID)
	if !errors.Is(err, ErrInvalidOrderAcceptTransition) {
		t.Fatalf("expected ErrInvalidOrderAcceptTransition, got %v", err)
	}
}

func TestVendorOrderUseCase_AcceptOrder_NotBelongToVendor(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(nil, nil)

	_, err := uc.AcceptOrder(ctx, vendorID, orderID, userID)
	if !errors.Is(err, ErrOrderNotBelongToVendor) {
		t.Fatalf("expected ErrOrderNotBelongToVendor, got %v", err)
	}
}

// --- RejectOrder tests ---

func TestVendorOrderUseCase_RejectOrder_Success_FromPaid(t *testing.T) {
	vendorOrderRepo, orderRepo, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	order := &domain.Order{
		ID:            orderID,
		VendorID:      vendorID,
		OrderStatus:   domain.OrderStatusPaid,
		PaymentStatus: "paid",
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)
	orderRepo.EXPECT().UpdateOrderStatus(ctx, orderID, domain.OrderStatusCanceled, "paid", gomock.Any(), gomock.Any()).Return(nil)
	orderRepo.EXPECT().RestoreStock(ctx, orderID).Return(nil)

	res, err := uc.RejectOrder(ctx, vendorID, orderID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusCanceled {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusCanceled, res.OrderStatus)
	}
}

func TestVendorOrderUseCase_RejectOrder_Success_FromProcessing(t *testing.T) {
	vendorOrderRepo, orderRepo, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	order := &domain.Order{
		ID:            orderID,
		VendorID:      vendorID,
		OrderStatus:   domain.OrderStatusProcessing,
		PaymentStatus: "paid",
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)
	orderRepo.EXPECT().UpdateOrderStatus(ctx, orderID, domain.OrderStatusCanceled, "paid", gomock.Any(), gomock.Any()).Return(nil)
	orderRepo.EXPECT().RestoreStock(ctx, orderID).Return(nil)

	res, err := uc.RejectOrder(ctx, vendorID, orderID, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusCanceled {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusCanceled, res.OrderStatus)
	}
}

func TestVendorOrderUseCase_RejectOrder_InvalidTransition(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()

	order := &domain.Order{
		ID:          orderID,
		VendorID:    vendorID,
		OrderStatus: domain.OrderStatusShipped,
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)

	_, err := uc.RejectOrder(ctx, vendorID, orderID, uuid.New())
	if !errors.Is(err, ErrInvalidOrderRejectTransition) {
		t.Fatalf("expected ErrInvalidOrderRejectTransition, got %v", err)
	}
}

// --- ShipOrder tests ---

func TestVendorOrderUseCase_ShipOrder_Success(t *testing.T) {
	vendorOrderRepo, orderRepo, shipmentRepo, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	order := &domain.Order{
		ID:            orderID,
		VendorID:      vendorID,
		OrderStatus:   domain.OrderStatusProcessing,
		PaymentStatus: "paid",
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)
	shipmentRepo.EXPECT().FindByOrderID(ctx, orderID).Return(&domain.Shipment{
		ID: uuid.New(), OrderID: orderID, CourierCode: "jne", ServiceType: "REG",
		ETD: "1-2 day", ShipmentStatus: "waiting_pickup",
	}, nil)
	shipmentRepo.EXPECT().UpdateTrackingAndShip(ctx, orderID, "JNE12345678", gomock.Any()).Return(nil)
	orderRepo.EXPECT().UpdateOrderStatus(ctx, orderID, domain.OrderStatusShipped, "paid", gomock.Any(), gomock.Any()).Return(nil)

	req := domain.ShipOrderRequest{
		TrackingNo: "JNE12345678",
	}

	res, err := uc.ShipOrder(ctx, vendorID, orderID, userID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusShipped {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusShipped, res.OrderStatus)
	}
	if res.TrackingNo != "JNE12345678" {
		t.Fatalf("expected tracking JNE12345678, got %s", res.TrackingNo)
	}
	if res.ETD != "1-2 day" {
		t.Fatalf("expected ETD '1-2 day', got %s", res.ETD)
	}
	if res.EstimatedArrival == nil {
		t.Fatalf("expected estimated arrival to be set")
	}
}

func TestVendorOrderUseCase_ShipOrder_NotProcessing(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()

	order := &domain.Order{
		ID:          orderID,
		VendorID:    vendorID,
		OrderStatus: domain.OrderStatusPaid,
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)

	req := domain.ShipOrderRequest{TrackingNo: "JNE12345678"}

	_, err := uc.ShipOrder(ctx, vendorID, orderID, uuid.New(), req)
	if !errors.Is(err, ErrInvalidOrderShipTransition) {
		t.Fatalf("expected ErrInvalidOrderShipTransition, got %v", err)
	}
}

func TestVendorOrderUseCase_ShipOrder_EmptyTrackingNo(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()

	order := &domain.Order{
		ID:          orderID,
		VendorID:    vendorID,
		OrderStatus: domain.OrderStatusProcessing,
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)

	req := domain.ShipOrderRequest{TrackingNo: "  "}

	_, err := uc.ShipOrder(ctx, vendorID, orderID, uuid.New(), req)
	if !errors.Is(err, ErrTrackingNumberRequired) {
		t.Fatalf("expected ErrTrackingNumberRequired, got %v", err)
	}
}

// --- TrackWaybill tests ---

func TestVendorOrderUseCase_TrackWaybill_Success(t *testing.T) {
	vendorOrderRepo, _, shipmentRepo, _, _, _, rajaOngkir, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	now := time.Now()

	order := &domain.Order{ID: orderID, VendorID: vendorID, OrderStatus: domain.OrderStatusShipped}
	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)

	shipment := &domain.Shipment{
		ID:          uuid.New(),
		OrderID:     orderID,
		CourierCode: "jne",
		TrackingNo:  "JNE12345",
		ShippedAt:   &now,
	}
	shipmentRepo.EXPECT().FindByOrderID(ctx, orderID).Return(shipment, nil)

	trackResp := &domain.TrackWaybillResponse{
		Delivered: false,
		Summary: domain.TrackWaybillSummary{
			CourierCode:   "jne",
			WaybillNumber: "JNE12345",
			Status:        "ON PROCESS",
		},
	}
	rajaOngkir.EXPECT().TrackWaybill(ctx, "JNE12345", "jne").Return(trackResp, nil)

	res, err := uc.TrackWaybill(ctx, vendorID, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Summary.WaybillNumber != "JNE12345" {
		t.Fatalf("expected waybill JNE12345, got %s", res.Summary.WaybillNumber)
	}
}

func TestVendorOrderUseCase_TrackWaybill_NoShipment(t *testing.T) {
	vendorOrderRepo, _, shipmentRepo, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()

	order := &domain.Order{ID: orderID, VendorID: vendorID, OrderStatus: domain.OrderStatusPaid}
	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)
	shipmentRepo.EXPECT().FindByOrderID(ctx, orderID).Return(nil, nil)

	_, err := uc.TrackWaybill(ctx, vendorID, orderID)
	if !errors.Is(err, ErrShipmentNotFound) {
		t.Fatalf("expected ErrShipmentNotFound, got %v", err)
	}
}

// --- ListOrders tests ---

func TestVendorOrderUseCase_ListOrders_Success(t *testing.T) {
	vendorOrderRepo, orderRepo, _, _, userRepo, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	userID := uuid.New()
	orderID := uuid.New()

	orders := []domain.Order{
		{
			ID:          orderID,
			UserID:      userID,
			VendorID:    vendorID,
			OrderNo:     "ORD-001",
			OrderStatus: domain.OrderStatusPaid,
			GrandTotal:  100000,
			PlacedAt:    time.Now(),
		},
	}

	vendorOrderRepo.EXPECT().ListByVendor(ctx, vendorID, gomock.Any()).Return(orders, int64(1), nil)
	orderRepo.EXPECT().FindItemsByOrderIDs(ctx, []uuid.UUID{orderID}).Return(map[uuid.UUID][]domain.OrderItem{
		orderID: {
			{ID: uuid.New(), OrderID: orderID, ProductNameSnapshot: "Sajadah Premium", SKUSnapshot: "L", Qty: 1, UnitPrice: 100000, LineTotal: 100000},
		},
	}, nil)
	userRepo.EXPECT().FindByID(ctx, userID).Return(&domain.User{ID: userID, FullName: "John Doe"}, nil)

	items, meta, err := uc.ListOrders(ctx, vendorID, domain.VendorOrderListParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].CustomerName != "John Doe" {
		t.Fatalf("expected customer name John Doe, got %s", items[0].CustomerName)
	}
	if meta.TotalItems != 1 {
		t.Fatalf("expected total 1, got %d", meta.TotalItems)
	}
}

func TestVendorOrderUseCase_ListOrders_InvalidStatus(t *testing.T) {
	_, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()

	_, _, err := uc.ListOrders(ctx, vendorID, domain.VendorOrderListParams{Page: 1, Limit: 10, Status: "nonexistent"})
	if !errors.Is(err, ErrInvalidOrderStatus) {
		t.Fatalf("expected ErrInvalidOrderStatus, got %v", err)
	}
}

func TestVendorOrderUseCase_ListOrders_Empty(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()

	vendorOrderRepo.EXPECT().ListByVendor(ctx, vendorID, gomock.Any()).Return([]domain.Order{}, int64(0), nil)

	items, meta, err := uc.ListOrders(ctx, vendorID, domain.VendorOrderListParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
	if meta.TotalItems != 0 {
		t.Fatalf("expected total 0, got %d", meta.TotalItems)
	}
}

// --- GetOrderDetail tests ---

func TestVendorOrderUseCase_GetOrderDetail_Success(t *testing.T) {
	vendorOrderRepo, orderRepo, shipmentRepo, paymentRepo, _, vendorRepo, rajaOngkir, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()

	order := &domain.Order{
		ID:          orderID,
		UserID:      userID,
		VendorID:    vendorID,
		OrderNo:     "ORD-001",
		OrderStatus: domain.OrderStatusShipped,
		Subtotal:    90000,
		ShippingFee: 10000,
		PlatformFee: 5000,
		GrandTotal:  105000,
		PlacedAt:    time.Now(),
		ShippingAddressSnapshot: map[string]interface{}{
			"recipient_name": "Jane Doe",
			"phone":          "08123456789",
			"address_line":   "Jl. Merdeka No. 1",
			"district_name":  "Menteng",
			"city_name":      "Jakarta",
			"province_name":  "DKI Jakarta",
			"postal_code":    "10110",
		},
	}

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(order, nil)
	orderRepo.EXPECT().FindItemsByOrderID(ctx, orderID).Return([]domain.OrderItem{
		{ID: uuid.New(), OrderID: orderID, ProductNameSnapshot: "Sajadah", SKUSnapshot: "L", Qty: 1, UnitPrice: 90000, LineTotal: 90000},
	}, nil)
	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{ID: vendorID, DisplayName: "Toko Sajadah"}, nil)

	now := time.Now()
	shipmentRepo.EXPECT().FindByOrderID(ctx, orderID).Return(&domain.Shipment{
		ID: uuid.New(), OrderID: orderID, CourierCode: "jne", TrackingNo: "JNE123",
		ETD: "2-3 day", ShipmentStatus: "shipped", ShippedAt: &now,
	}, nil)

	// Tracking API call
	rajaOngkir.EXPECT().TrackWaybill(ctx, "JNE123", "jne").Return(&domain.TrackWaybillResponse{
		Delivered: false,
		Manifest: []domain.TrackManifestItem{
			{ManifestDescription: "Paket dalam perjalanan ke hub tujuan", ManifestDate: "2025-01-02", ManifestTime: "10:00:00"},
		},
	}, nil)

	method := "VIRTUAL_ACCOUNT"
	channel := "BCA"
	paymentRepo.EXPECT().FindInvoiceByOrderID(ctx, orderID).Return(&domain.PaymentInvoice{
		ID: uuid.New(), OrderID: orderID, Status: "paid", PaidAt: &now,
		PaymentMethod: &method, PaymentChannel: &channel,
	}, nil)

	res, err := uc.GetOrderDetail(ctx, vendorID, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check store name
	if res.Store.Name != "Toko Sajadah" {
		t.Fatalf("expected store name 'Toko Sajadah', got %s", res.Store.Name)
	}

	// Check recipient
	if res.Recipient.Name != "Jane Doe" {
		t.Fatalf("expected recipient name 'Jane Doe', got %s", res.Recipient.Name)
	}
	if res.Recipient.Address.City != "Jakarta" {
		t.Fatalf("expected city 'Jakarta', got %s", res.Recipient.Address.City)
	}

	// Check shipping status from tracking API
	if res.Shipping.Status != "Paket dalam perjalanan ke hub tujuan" {
		t.Fatalf("expected shipping status from tracking API, got %s", res.Shipping.Status)
	}

	// Check payment method
	if res.Payment.Method != "VIRTUAL_ACCOUNT" {
		t.Fatalf("expected payment method 'VIRTUAL_ACCOUNT', got %s", res.Payment.Method)
	}

	// Check summary
	if res.Summary.ItemsTotal != 90000 {
		t.Fatalf("expected items total 90000, got %f", res.Summary.ItemsTotal)
	}
	if res.Summary.ShippingFee != 10000 {
		t.Fatalf("expected shipping fee 10000, got %f", res.Summary.ShippingFee)
	}
	if res.Summary.ServiceFee != 5000 {
		t.Fatalf("expected service fee 5000, got %f", res.Summary.ServiceFee)
	}
	if res.Summary.GrandTotal != 105000 {
		t.Fatalf("expected grand total 105000, got %f", res.Summary.GrandTotal)
	}
}

func TestVendorOrderUseCase_GetOrderDetail_NotFound(t *testing.T) {
	vendorOrderRepo, _, _, _, _, _, _, uc := setupVendorOrderUseCase(t)
	ctx := context.Background()
	vendorID := uuid.New()
	orderID := uuid.New()

	vendorOrderRepo.EXPECT().FindByIDAndVendor(ctx, orderID, vendorID).Return(nil, nil)

	_, err := uc.GetOrderDetail(ctx, vendorID, orderID)
	if !errors.Is(err, ErrOrderNotBelongToVendor) {
		t.Fatalf("expected ErrOrderNotBelongToVendor, got %v", err)
	}
}

// --- parseMaxETDDays tests ---

func TestParseMaxETDDays(t *testing.T) {
	tests := []struct {
		etd      string
		expected int
	}{
		{"1-2 day", 2},
		{"3 day", 3},
		{"1-2", 2},
		{"5-7 hari", 7},
		{"", 0},
		{"unknown", 0},
		{"10-14 day", 14},
	}
	for _, tc := range tests {
		got := parseMaxETDDays(tc.etd)
		if got != tc.expected {
			t.Errorf("parseMaxETDDays(%q) = %d, want %d", tc.etd, got, tc.expected)
		}
	}
}
