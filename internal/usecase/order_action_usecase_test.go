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

func setupOrderActionUseCase(t *testing.T) (*mocks.MockOrderRepository, domain.OrderActionUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	orderRepo := mocks.NewMockOrderRepository(ctrl)
	uc := NewOrderActionUseCase(orderRepo)
	return orderRepo, uc
}

func TestOrderActionUseCase_CompleteByCustomer(t *testing.T) {
	orderRepo, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusShipped,
		PaymentStatus: domain.PaymentStatusPaid,
	}, nil)
	orderRepo.EXPECT().MarkOrderReceived(ctx, orderID, &userID, gomock.Any()).Return(true, nil)

	res, err := uc.CompleteByCustomer(ctx, userID, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.OrderStatus != domain.OrderStatusReceived {
		t.Fatalf("expected received status, got %s", res.OrderStatus)
	}
}

func TestOrderActionUseCase_CompleteByCustomer_InvalidTransition(t *testing.T) {
	orderRepo, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
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

func TestOrderActionUseCase_ListByCustomer(t *testing.T) {
	orderRepo, uc := setupOrderActionUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	orderID := uuid.New()
	orderDate := time.Date(2026, 3, 1, 8, 30, 0, 0, time.UTC)

	orderRepo.EXPECT().
		ListByUser(ctx, userID, domain.CustomerOrderListParams{Page: 1, Limit: 10, Status: domain.OrderStatusPaid}).
		Return([]domain.Order{
			{
				ID:          orderID,
				GrandTotal:  206000,
				PlacedAt:    orderDate,
				OrderStatus: domain.OrderStatusPaid,
			},
		}, int64(1), nil)

	orderRepo.EXPECT().
		FindItemsByOrderIDs(ctx, []uuid.UUID{orderID}).
		Return(map[uuid.UUID][]domain.OrderItem{
			orderID: {
				{
					ProductNameSnapshot: "Kurma Ajwa",
					SKUSnapshot:         "500gr",
					Qty:                 2,
				},
			},
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
	orderRepo, uc := setupOrderActionUseCase(t)
	ctx := context.Background()
	userID := uuid.New()

	orderRepo.EXPECT().
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
