package usecase

import (
	"context"
	"fmt"
	"math"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminReportUseCase struct {
	repo domain.AdminReportRepository
}

// NewAdminReportUseCase creates a new admin report use case.
func NewAdminReportUseCase(repo domain.AdminReportRepository) domain.AdminReportUseCase {
	return &adminReportUseCase{repo: repo}
}

func (uc *adminReportUseCase) GetReport(ctx context.Context, params domain.AdminReportParams) (*domain.AdminReportResponse, *domain.PaginationMeta, error) {
	params = normalizeAdminReportParams(params)

	summary, err := uc.repo.GetSummary(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load admin report summary: %w", err)
	}

	items, total, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load admin report items: %w", err)
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return &domain.AdminReportResponse{
		Year:    params.Year,
		Summary: *summary,
		Items:   items,
	}, meta, nil
}

func (uc *adminReportUseCase) ExportReport(ctx context.Context, params domain.AdminReportParams) ([]domain.AdminReportItem, error) {
	params = normalizeAdminReportParams(params)
	params.Page = 1
	params.Limit = 0

	items, _, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to export admin report items: %w", err)
	}

	return items, nil
}

func normalizeAdminReportParams(params domain.AdminReportParams) domain.AdminReportParams {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}
	return params
}
