package usecase

import (
	"context"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminDashboardUseCase struct {
	repo domain.AdminDashboardRepository
}

// NewAdminDashboardUseCase creates a new admin dashboard use case.
func NewAdminDashboardUseCase(repo domain.AdminDashboardRepository) domain.AdminDashboardUseCase {
	return &adminDashboardUseCase{repo: repo}
}

func (uc *adminDashboardUseCase) GetDashboard(ctx context.Context, params domain.AdminDashboardParams) (*domain.AdminDashboardResponse, error) {
	summary, err := uc.repo.GetSummary(ctx, params)
	if err != nil {
		return nil, err
	}

	transactions, err := uc.repo.ListTransactions(ctx, params, 10)
	if err != nil {
		return nil, err
	}

	return &domain.AdminDashboardResponse{
		Month:        params.Month,
		Year:         params.Year,
		Summary:      *summary,
		Transactions: transactions,
	}, nil
}
