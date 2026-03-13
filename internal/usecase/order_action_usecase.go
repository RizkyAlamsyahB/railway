package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type orderActionUseCase struct {
	orderRepo domain.OrderRepository
}

// NewOrderActionUseCase creates a new OrderActionUseCase.
func NewOrderActionUseCase(orderRepo domain.OrderRepository) domain.OrderActionUseCase {
	return &orderActionUseCase{orderRepo: orderRepo}
}

func (uc *orderActionUseCase) CompleteByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*domain.ReceiveOrderResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}

	if order.OrderStatus != domain.OrderStatusShipped {
		return nil, ErrInvalidOrderCompletionTransition
	}

	notes := "Order marked as received by customer"
	applied, err := uc.orderRepo.MarkOrderReceived(ctx, orderID, &userID, &notes)
	if err != nil {
		return nil, fmt.Errorf("failed to mark order received: %w", err)
	}
	if !applied {
		return nil, ErrInvalidOrderCompletionTransition
	}

	return &domain.ReceiveOrderResponse{
		OrderID:       orderID,
		OrderStatus:   domain.OrderStatusReceived,
		PaymentStatus: order.PaymentStatus,
		ReceivedAt:    time.Now(),
	}, nil
}

func (uc *orderActionUseCase) ListByCustomer(ctx context.Context, userID uuid.UUID, params domain.CustomerOrderListParams) ([]domain.CustomerOrderListItem, *domain.PaginationMeta, error) {
	params = normalizeCustomerOrderListParams(params)
	if params.Status != "" && !isValidOrderStatus(params.Status) {
		return nil, nil, ErrInvalidOrderStatus
	}

	orders, total, err := uc.orderRepo.ListByUser(ctx, userID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list orders: %w", err)
	}

	meta := buildOrderPaginationMeta(params.Page, params.Limit, total)
	if len(orders) == 0 {
		return []domain.CustomerOrderListItem{}, meta, nil
	}

	orderIDs := make([]uuid.UUID, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	itemsByOrderID, err := uc.orderRepo.FindItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list order items: %w", err)
	}

	result := make([]domain.CustomerOrderListItem, 0, len(orders))
	for _, order := range orders {
		rawItems := itemsByOrderID[order.ID]
		items := make([]domain.CustomerOrderProductItem, 0, len(rawItems))
		for _, it := range rawItems {
			items = append(items, domain.CustomerOrderProductItem{
				ProductName:     it.ProductNameSnapshot,
				SelectedVariant: it.SKUSnapshot,
				Qty:             it.Qty,
			})
		}

		result = append(result, domain.CustomerOrderListItem{
			OrderID:      order.ID,
			Items:        items,
			TotalPayment: order.GrandTotal,
			OrderDate:    order.PlacedAt,
			Status:       order.OrderStatus,
		})
	}

	return result, meta, nil
}

func (uc *orderActionUseCase) ListOrderStatuses(_ context.Context) []string {
	result := make([]string, len(domain.OrderStatuses))
	copy(result, domain.OrderStatuses)
	return result
}

func normalizeCustomerOrderListParams(params domain.CustomerOrderListParams) domain.CustomerOrderListParams {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}
	params.Status = strings.ToLower(strings.TrimSpace(params.Status))
	return params
}

func buildOrderPaginationMeta(page, limit int, total int64) *domain.PaginationMeta {
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

func isValidOrderStatus(status string) bool {
	for _, allowed := range domain.OrderStatuses {
		if status == allowed {
			return true
		}
	}
	return false
}
