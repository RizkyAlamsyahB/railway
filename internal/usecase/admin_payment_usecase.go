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

func (uc *adminPaymentUseCase) Summary(ctx context.Context, params domain.AdminPaymentListParams) (*domain.AdminPaymentSummary, error) {
	summary, err := uc.paymentRepo.SummaryForAdmin(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}
	return summary, nil
}

func (uc *adminPaymentUseCase) Export(ctx context.Context, params domain.AdminPaymentListParams) ([]domain.AdminPaymentListItem, error) {
	items, err := uc.paymentRepo.ExportForAdmin(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to export payment invoices: %w", err)
	}
	return items, nil
}
