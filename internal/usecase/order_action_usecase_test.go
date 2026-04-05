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

type orderActionMocks struct {
	orderRepo    *mocks.MockOrderRepository
	shipmentRepo *mocks.MockShipmentRepository
	paymentRepo  *mocks.MockPaymentRepository
	vendorRepo   *mocks.MockVendorRepository
	productRepo  *mocks.MockProductRepository
	rajaOngkir   *mocks.MockRajaOngkirProvider
	storage      *mocks.MockStorageProvider
}

type stubCustomerRefundRepo struct {
	hasOpen    bool
	hasOpenErr error
	createResp *domain.CreateOrderRefundResponse
	createErr  error
}

type stubReturnReasonRepo struct {
	reason *domain.ReturnReason
	err    error
}

func (s *stubReturnReasonRepo) Create(context.Context, *domain.ReturnReason) error { return nil }
func (s *stubReturnReasonRepo) FindByID(context.Context, int) (*domain.ReturnReason, error) {
	return s.reason, s.err
}
func (s *stubReturnReasonRepo) List(context.Context) ([]domain.ReturnReason, error) { return nil, nil }
func (s *stubReturnReasonRepo) Update(context.Context, *domain.ReturnReason) error  { return nil }
func (s *stubReturnReasonRepo) Delete(context.Context, int) error                   { return nil }

func (s *stubCustomerRefundRepo) HasOpenRefundByOrderID(context.Context, uuid.UUID) (bool, error) {
	return s.hasOpen, s.hasOpenErr
}

func (s *stubCustomerRefundRepo) CreateRefundRequest(context.Context, domain.CreateCustomerRefundInput) (*domain.CreateOrderRefundResponse, error) {
	return s.createResp, s.createErr
}

func (s *stubCustomerRefundRepo) SubmitRefundDestination(context.Context, uuid.UUID, string, string, string, string, string, *uuid.UUID) error {
	return nil
}

func (s *stubCustomerRefundRepo) FindByID(context.Context, uuid.UUID) (*domain.RefundDetail, error) {
	return nil, nil
}

func (s *stubCustomerRefundRepo) FindAwaitingDestinationByOrderAndUser(context.Context, uuid.UUID, uuid.UUID) (*domain.RefundDetail, error) {
	return nil, nil
}

func (s *stubCustomerRefundRepo) FindLatestByOrderID(context.Context, uuid.UUID) (*domain.RefundDetail, error) {
	return nil, nil
}

func setupOrderActionUseCase(t *testing.T) (*orderActionMocks, domain.OrderActionUseCase) {
	return setupOrderActionUseCaseWithDeps(t, nil, &stubReturnReasonRepo{}, nil)
}

func setupOrderActionUseCaseWithRefundRepo(t *testing.T, refundRepo domain.CustomerRefundRepository) (*orderActionMocks, domain.OrderActionUseCase) {
	return setupOrderActionUseCaseWithDeps(t, refundRepo, &stubReturnReasonRepo{}, nil)
}

func setupOrderActionUseCaseWithDeps(
	t *testing.T,
	refundRepo domain.CustomerRefundRepository,
	returnReasonRepo domain.ReturnReasonRepository,
	storage domain.StorageProvider,
) (*orderActionMocks, domain.OrderActionUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	m := &orderActionMocks{
		orderRepo:    mocks.NewMockOrderRepository(ctrl),
		shipmentRepo: mocks.NewMockShipmentRepository(ctrl),
		paymentRepo:  mocks.NewMockPaymentRepository(ctrl),
		vendorRepo:   mocks.NewMockVendorRepository(ctrl),
		productRepo:  mocks.NewMockProductRepository(ctrl),
		rajaOngkir:   mocks.NewMockRajaOngkirProvider(ctrl),
		storage:      mocks.NewMockStorageProvider(ctrl),
	}
	if storage == nil {
		storage = m.storage
	}
	uc := NewOrderActionUseCase(m.orderRepo, m.shipmentRepo, m.paymentRepo, refundRepo, returnReasonRepo, m.vendorRepo, m.productRepo, m.rajaOngkir, storage, nil, nil, nil)
	return m, uc
}

func TestOrderActionUseCase_CompleteByCustomer(t *testing.T) {
	m, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	m.orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusReceived,
		PaymentStatus: domain.PaymentStatusPaid,
	}, nil)
	m.orderRepo.EXPECT().MarkOrderCompleted(ctx, orderID, &userID, gomock.Any()).Return(true, nil)

	res, err := uc.CompleteByCustomer(ctx, userID, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusCompleted {
		t.Fatalf("expected completed status, got %s", res.OrderStatus)
	}
}

func TestOrderActionUseCase_CompleteByCustomer_InvalidTransition(t *testing.T) {
	m, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	m.orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusPendingPayment,
		PaymentStatus: domain.PaymentStatusUnpaid,
	}, nil)

	_, err := uc.CompleteByCustomer(ctx, userID, orderID)
	if !errors.Is(err, ErrInvalidOrderCompletionTransition) {
		t.Fatalf("expected ErrInvalidOrderCompletionTransition, got %v", err)
	}
}

func TestOrderActionUseCase_CancelByCustomer(t *testing.T) {
	m, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	m.orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusProcessing,
		PaymentStatus: domain.PaymentStatusPaid,
	}, nil)
	m.orderRepo.EXPECT().UpdateOrderStatus(ctx, orderID, domain.OrderStatusCanceled, domain.PaymentStatusPaid, &userID, gomock.Any()).Return(nil)
	m.orderRepo.EXPECT().RestoreStock(ctx, orderID).Return(nil)
	m.paymentRepo.EXPECT().FindInvoiceByOrderID(ctx, orderID).Return(nil, nil)

	res, err := uc.CancelByCustomer(ctx, userID, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusCanceled {
		t.Fatalf("expected canceled status, got %s", res.OrderStatus)
	}
}

func TestOrderActionUseCase_CancelByCustomer_InvalidTransition(t *testing.T) {
	m, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	m.orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusShipped,
		PaymentStatus: domain.PaymentStatusPaid,
	}, nil)

	_, err := uc.CancelByCustomer(ctx, userID, orderID)
	if !errors.Is(err, ErrInvalidOrderCancelTransition) {
		t.Fatalf("expected ErrInvalidOrderCancelTransition, got %v", err)
	}
}

func TestOrderActionUseCase_ListByCustomer(t *testing.T) {
	m, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	orderID := uuid.New()
	variantID := uuid.New()
	orderDate := time.Date(2026, 3, 1, 8, 30, 0, 0, time.UTC)

	m.orderRepo.EXPECT().
		ListByUser(ctx, userID, domain.CustomerOrderListParams{Page: 1, Limit: 10, Status: domain.OrderStatusPaid}).
		Return([]domain.Order{
			{
				ID:          orderID,
				GrandTotal:  206000,
				PlacedAt:    orderDate,
				OrderStatus: domain.OrderStatusPaid,
			},
		}, int64(1), nil)

	m.orderRepo.EXPECT().
		FindItemsByOrderIDs(ctx, []uuid.UUID{orderID}).
		Return(map[uuid.UUID][]domain.OrderItem{
			orderID: {
				{
					ProductVariantID:    variantID,
					ProductNameSnapshot: "Kurma Ajwa",
					SKUSnapshot:         "500gr",
					Qty:                 2,
					UnitPrice:           103000,
				},
			},
		}, nil)

	// Product image resolution
	productID := uuid.New()
	m.productRepo.EXPECT().FindVariantByID(ctx, variantID).Return(&domain.ProductVariant{
		ID:        variantID,
		ProductID: productID,
	}, nil)
	m.productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return([]domain.ProductImage{
		{IsPrimary: true, ImageURL: "https://example.com/img.jpg"},
	}, nil)

	items, meta, err := uc.ListByCustomer(ctx, userID, domain.CustomerOrderListParams{
		Page:   0, // should normalize to 1
		Limit:  0, // should normalize to 10
		Status: "PAID",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 order, got %d", len(items))
	}
	if items[0].OrderID != orderID {
		t.Fatalf("expected order id %s, got %s", orderID, items[0].OrderID)
	}
	if len(items[0].Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(items[0].Items))
	}
	if items[0].Items[0].SelectedVariant != "500gr" {
		t.Fatalf("expected selected_variant 500gr, got %s", items[0].Items[0].SelectedVariant)
	}
	if items[0].Items[0].ImageURL == nil || *items[0].Items[0].ImageURL != "https://example.com/img.jpg" {
		t.Fatalf("expected image_url, got %v", items[0].Items[0].ImageURL)
	}
	if meta.Page != 1 || meta.Limit != 10 || meta.TotalItems != 1 || meta.TotalPages != 1 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

func TestOrderActionUseCase_ListByCustomer_InvalidStatus(t *testing.T) {
	_, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	_, _, err := uc.ListByCustomer(ctx, uuid.New(), domain.CustomerOrderListParams{
		Page:   1,
		Limit:  10,
		Status: "unknown_status",
	})
	if !errors.Is(err, ErrInvalidOrderStatus) {
		t.Fatalf("expected ErrInvalidOrderStatus, got %v", err)
	}
}

func TestOrderActionUseCase_ListByCustomer_Empty(t *testing.T) {
	m, uc := setupOrderActionUseCase(t)
	ctx := context.Background()
	userID := uuid.New()

	m.orderRepo.EXPECT().
		ListByUser(ctx, userID, domain.CustomerOrderListParams{Page: 2, Limit: 10, Status: ""}).
		Return([]domain.Order{}, int64(0), nil)

	items, meta, err := uc.ListByCustomer(ctx, userID, domain.CustomerOrderListParams{
		Page:  2,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty items, got %d", len(items))
	}
	if meta.Page != 2 || meta.Limit != 10 || meta.TotalItems != 0 || meta.TotalPages != 0 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

func TestOrderActionUseCase_RequestRefundByCustomer_Success(t *testing.T) {
	refundID := uuid.New()
	orderID := uuid.New()
	userID := uuid.New()
	invoiceID := uuid.New()
	returnReasonID := 1

	refundRepo := &stubCustomerRefundRepo{
		hasOpen: false,
		createResp: &domain.CreateOrderRefundResponse{
			RefundID:                      refundID,
			OrderID:                       orderID,
			Status:                        domain.RefundStatusRequested,
			Amount:                        452500,
			ReturnReasonID:                returnReasonID,
			Reason:                        "Barang rusak",
			DestinationChannelCode:        domain.PayoutChannelIDBCA,
			DestinationBankName:           "Bank Central Asia",
			DestinationAccountHolderName:  "Budi Santoso",
			DestinationAccountNumberLast4: "7890",
			RequestedAt:                   time.Now().UTC(),
		},
	}
	m, uc := setupOrderActionUseCaseWithDeps(t, refundRepo, &stubReturnReasonRepo{
		reason: &domain.ReturnReason{
			ID:     returnReasonID,
			Reason: "Barang rusak",
		},
	}, nil)
	ctx := context.Background()

	m.orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusReceived,
		PaymentStatus: domain.PaymentStatusPaid,
		GrandTotal:    452500,
	}, nil)
	m.paymentRepo.EXPECT().FindInvoiceByOrderID(ctx, orderID).Return(&domain.PaymentInvoice{
		ID: invoiceID,
	}, nil)

	res, err := uc.RequestRefundByCustomer(ctx, userID, orderID, domain.CreateOrderRefundRequest{
		ReturnReasonID:               returnReasonID,
		Description:                  "Pecah saat diterima",
		DestinationBankName:          "Bank Central Asia",
		DestinationAccountNumber:     "1234567890",
		DestinationAccountHolderName: "Budi Santoso",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil || res.RefundID != refundID {
		t.Fatalf("expected refund id %s, got %+v", refundID, res)
	}
}

func TestOrderActionUseCase_RequestRefundByCustomer_AlreadyExists(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	refundRepo := &stubCustomerRefundRepo{hasOpen: true}
	m, uc := setupOrderActionUseCaseWithRefundRepo(t, refundRepo)
	ctx := context.Background()

	m.orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusCompleted,
		PaymentStatus: domain.PaymentStatusPaid,
		GrandTotal:    300000,
	}, nil)

	_, err := uc.RequestRefundByCustomer(ctx, userID, orderID, domain.CreateOrderRefundRequest{
		ReturnReasonID:               1,
		DestinationBankName:          "Bank Central Asia",
		DestinationAccountNumber:     "1234567890",
		DestinationAccountHolderName: "Budi Santoso",
	})
	if !errors.Is(err, ErrRefundAlreadyRequested) {
		t.Fatalf("expected ErrRefundAlreadyRequested, got %v", err)
	}
}

func TestOrderActionUseCase_ListOrderStatuses(t *testing.T) {
	_, uc := setupOrderActionUseCase(t)

	statuses := uc.ListOrderStatuses(context.Background())
	if len(statuses) != len(domain.OrderStatuses) {
		t.Fatalf("expected %d statuses, got %d", len(domain.OrderStatuses), len(statuses))
	}
	for i := range domain.OrderStatuses {
		if statuses[i] != domain.OrderStatuses[i] {
			t.Fatalf("expected status %s at index %d, got %s", domain.OrderStatuses[i], i, statuses[i])
		}
	}
}
