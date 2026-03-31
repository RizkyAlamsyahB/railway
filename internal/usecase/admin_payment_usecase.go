package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminPaymentUseCase struct {
	paymentRepo domain.PaymentRepository
}

// NewAdminPaymentUseCase creates a new AdminPaymentUseCase.
func NewAdminPaymentUseCase(paymentRepo domain.PaymentRepository) domain.AdminPaymentUseCase {
	return &adminPaymentUseCase{paymentRepo: paymentRepo}
}

func (uc *adminPaymentUseCase) List(ctx context.Context, params domain.AdminPaymentListParams) ([]domain.AdminPaymentListItem, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	items, total, err := uc.paymentRepo.ListForAdmin(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list payment invoices: %w", err)
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}
