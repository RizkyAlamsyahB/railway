package usecase

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type vendorReportUseCase struct {
	repo domain.VendorReportRepository
}

// NewVendorReportUseCase creates a new VendorReportUseCase.
func NewVendorReportUseCase(repo domain.VendorReportRepository) domain.VendorReportUseCase {
	return &vendorReportUseCase{repo: repo}
}

func (uc *vendorReportUseCase) GetReport(ctx context.Context, vendorID uuid.UUID, params domain.VendorReportParams) (*domain.VendorReportResponse, *domain.PaginationMeta, error) {
	page, limit := normalizePaging(params.Page, params.Limit)
	params.Page = page
	params.Limit = limit

	period, err := uc.parseReportPeriod(params)
	if err != nil {
		return nil, nil, err
	}

	summary, err := uc.repo.GetSummary(ctx, vendorID, period)
	if err != nil {
		return nil, nil, err
	}

	items, total, err := uc.repo.ListDailyRevenue(ctx, vendorID, period, params)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if limit == 0 {
		totalPages = 1
	}
	meta := &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return &domain.VendorReportResponse{
		Summary: *summary,
		Items:   items,
	}, meta, nil
}

func (uc *vendorReportUseCase) ExportReport(ctx context.Context, vendorID uuid.UUID, params domain.VendorReportParams) ([]domain.VendorReportDailyItem, error) {
	period, err := uc.parseReportPeriod(params)
	if err != nil {
		return nil, err
	}

	params.Page = 1
	params.Limit = 0

	items, _, err := uc.repo.ListDailyRevenue(ctx, vendorID, period, params)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *vendorReportUseCase) parseReportPeriod(params domain.VendorReportParams) (domain.FinancePeriod, error) {
	trimmed := strings.TrimSpace(params.Month)
	if trimmed == "" {
		now := time.Now().In(wibLocation)
		trimmed = now.Format("2006-01")
	}

	parsed, err := time.ParseInLocation("2006-01", trimmed, wibLocation)
	if err != nil {
		return domain.FinancePeriod{}, ErrInvalidMonth
	}

	startLocal := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, wibLocation)
	endLocal := startLocal.AddDate(0, 1, 0)

	// Apply optional date_from / date_to narrowing.
	if df := strings.TrimSpace(params.DateFrom); df != "" {
		v, err := time.ParseInLocation("2006-01-02", df, wibLocation)
		if err != nil {
			return domain.FinancePeriod{}, ErrInvalidDate
		}
		fromStart := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, wibLocation)
		if fromStart.After(startLocal) {
			startLocal = fromStart
		}
	}

	if dt := strings.TrimSpace(params.DateTo); dt != "" {
		v, err := time.ParseInLocation("2006-01-02", dt, wibLocation)
		if err != nil {
			return domain.FinancePeriod{}, ErrInvalidDate
		}
		toEnd := time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, wibLocation).AddDate(0, 0, 1)
		if toEnd.Before(endLocal) {
			endLocal = toEnd
		}
	}

	if !startLocal.Before(endLocal) {
		return domain.FinancePeriod{}, ErrInvalidDateRange
	}

	return domain.FinancePeriod{
		Start: startLocal.UTC(),
		End:   endLocal.UTC(),
	}, nil
}
