package usecase

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type financeUseCase struct {
	repo domain.FinanceRepository
}

// NewFinanceUseCase creates a new FinanceUseCase.
func NewFinanceUseCase(repo domain.FinanceRepository) domain.FinanceUseCase {
	return &financeUseCase{repo: repo}
}

func (uc *financeUseCase) GetDashboard(ctx context.Context, month string) (*domain.FinanceDashboardResponse, error) {
	period, monthLabel, err := uc.parseMonth(month)
	if err != nil {
		return nil, err
	}

	summary, err := uc.repo.GetDashboardSummary(ctx, period)
	if err != nil {
		return nil, err
	}

	dayStart, dayEnd := utcDayRange(time.Now().UTC())
	todayTransactions, err := uc.repo.ListTodayTransactions(ctx, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	return &domain.FinanceDashboardResponse{
		Month:             monthLabel,
		Summary:           *summary,
		TodayTransactions: todayTransactions,
	}, nil
}

func (uc *financeUseCase) ListTransactions(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceTransactionItem, *domain.PaginationMeta, error) {
	page, limit := normalizePaging(params.Page, params.Limit)
	params.Page = page
	params.Limit = limit

	if params.Status != "" && !isValidPaymentStatus(params.Status) {
		return nil, nil, ErrInvalidPaymentStatus
	}

	period, _, err := uc.parseMonth(params.Month)
	if err != nil {
		return nil, nil, err
	}

	items, total, err := uc.repo.ListTransactions(ctx, period, params)
	if err != nil {
		return nil, nil, err
	}

	meta := buildPaginationMeta(page, limit, total)
	return items, meta, nil
}

func (uc *financeUseCase) ExportTransactions(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceTransactionItem, error) {
	if params.Status != "" && !isValidPaymentStatus(params.Status) {
		return nil, ErrInvalidPaymentStatus
	}

	period, _, err := uc.parseMonth(params.Month)
	if err != nil {
		return nil, err
	}

	params.Page = 1
	params.Limit = 0

	items, _, err := uc.repo.ListTransactions(ctx, period, params)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *financeUseCase) ListPayouts(ctx context.Context, params domain.FinanceListParams) ([]domain.FinancePayoutItem, *domain.PaginationMeta, error) {
	page, limit := normalizePaging(params.Page, params.Limit)
	params.Page = page
	params.Limit = limit

	if params.Status != "" {
		normalized := normalizePayoutStatus(params.Status)
		if !isValidPayoutStatus(normalized) {
			return nil, nil, ErrInvalidPayoutStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseMonth(params.Month)
	if err != nil {
		return nil, nil, err
	}

	items, total, err := uc.repo.ListPayouts(ctx, period, params)
	if err != nil {
		return nil, nil, err
	}

	meta := buildPaginationMeta(page, limit, total)
	return items, meta, nil
}

func (uc *financeUseCase) ExportPayouts(ctx context.Context, params domain.FinanceListParams) ([]domain.FinancePayoutItem, error) {
	if params.Status != "" {
		normalized := normalizePayoutStatus(params.Status)
		if !isValidPayoutStatus(normalized) {
			return nil, ErrInvalidPayoutStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseMonth(params.Month)
	if err != nil {
		return nil, err
	}

	params.Page = 1
	params.Limit = 0

	items, _, err := uc.repo.ListPayouts(ctx, period, params)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *financeUseCase) ListRefunds(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceRefundItem, *domain.PaginationMeta, error) {
	page, limit := normalizePaging(params.Page, params.Limit)
	params.Page = page
	params.Limit = limit

	if params.Status != "" {
		normalized := normalizeRefundStatus(params.Status)
		if !isValidRefundStatus(normalized) {
			return nil, nil, ErrInvalidRefundStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseMonth(params.Month)
	if err != nil {
		return nil, nil, err
	}

	items, total, err := uc.repo.ListRefunds(ctx, period, params)
	if err != nil {
		return nil, nil, err
	}

	meta := buildPaginationMeta(page, limit, total)
	return items, meta, nil
}

func (uc *financeUseCase) ExportRefunds(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceRefundItem, error) {
	if params.Status != "" {
		normalized := normalizeRefundStatus(params.Status)
		if !isValidRefundStatus(normalized) {
			return nil, ErrInvalidRefundStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseMonth(params.Month)
	if err != nil {
		return nil, err
	}

	params.Page = 1
	params.Limit = 0

	items, _, err := uc.repo.ListRefunds(ctx, period, params)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (uc *financeUseCase) UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, status string, actorID uuid.UUID) (*domain.StatusActionResponse, error) {
	normalized := normalizeRefundStatus(status)
	if !isValidRefundStatus(normalized) {
		return nil, ErrInvalidRefundStatus
	}

	record, err := uc.repo.GetRefundRecord(ctx, refundID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrRefundNotFound
	}

	if !isRefundTransitionAllowed(record.Status, normalized) {
		return nil, ErrInvalidRefundTransition
	}

	resp, err := uc.repo.UpdateRefundStatus(ctx, refundID, normalized, actorID)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrRefundNotFound
	}
	return resp, nil
}

func (uc *financeUseCase) UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status string, actorID uuid.UUID) (*domain.StatusActionResponse, error) {
	normalized := normalizePayoutStatus(status)
	if !isValidPayoutStatus(normalized) {
		return nil, ErrInvalidPayoutStatus
	}

	record, err := uc.repo.GetPayoutRecord(ctx, payoutID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrPayoutNotFound
	}

	resp, err := uc.repo.UpdatePayoutStatus(ctx, payoutID, normalized, actorID)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrPayoutNotFound
	}
	return resp, nil
}

func (uc *financeUseCase) parseMonth(month string) (domain.FinancePeriod, string, error) {
	trimmed := strings.TrimSpace(month)
	if trimmed == "" {
		now := time.Now().UTC()
		trimmed = now.Format("2006-01")
	}

	parsed, err := time.ParseInLocation("2006-01", trimmed, time.UTC)
	if err != nil {
		return domain.FinancePeriod{}, "", ErrInvalidMonth
	}

	start := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return domain.FinancePeriod{
		Start: start,
		End:   end,
	}, trimmed, nil
}

func normalizePaging(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}

func buildPaginationMeta(page, limit int, total int64) *domain.PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if limit == 0 {
		totalPages = 1
	}
	return &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

func utcDayRange(now time.Time) (time.Time, time.Time) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return start, start.AddDate(0, 0, 1)
}

func isValidPaymentStatus(status string) bool {
	switch status {
	case domain.PaymentInvoiceStatusPending,
		domain.PaymentInvoiceStatusPaid,
		domain.PaymentInvoiceStatusFailed,
		domain.PaymentInvoiceStatusExpired,
		domain.PaymentInvoiceStatusRefunded:
		return true
	default:
		return false
	}
}

func normalizeRefundStatus(status string) string {
	if strings.EqualFold(status, "pending") {
		return domain.RefundStatusRequested
	}
	return status
}

func isValidRefundStatus(status string) bool {
	switch status {
	case domain.RefundStatusRequested,
		domain.RefundStatusApproved,
		domain.RefundStatusRejected,
		domain.RefundStatusProcessed:
		return true
	default:
		return false
	}
}

func isRefundTransitionAllowed(current, target string) bool {
	if current == target {
		return true
	}

	switch current {
	case domain.RefundStatusRequested:
		return target == domain.RefundStatusApproved || target == domain.RefundStatusRejected
	case domain.RefundStatusApproved:
		return target == domain.RefundStatusProcessed
	default:
		return false
	}
}

func normalizePayoutStatus(status string) string {
	if strings.EqualFold(status, "on_hold") {
		return domain.PayoutStatusOnHold
	}
	return status
}

func isValidPayoutStatus(status string) bool {
	switch status {
	case domain.PayoutStatusReady,
		domain.PayoutStatusSchedule,
		domain.PayoutStatusComplete,
		domain.PayoutStatusFailed,
		domain.PayoutStatusOnHold:
		return true
	default:
		return false
	}
}
