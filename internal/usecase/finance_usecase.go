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

type financeUseCase struct {
	repo         domain.FinanceRepository
	xenditPayout domain.XenditPayoutProvider
}

var wibLocation = time.FixedZone("WIB", 7*60*60)

// NewFinanceUseCase creates a new FinanceUseCase.
func NewFinanceUseCase(repo domain.FinanceRepository, xenditPayout domain.XenditPayoutProvider) domain.FinanceUseCase {
	return &financeUseCase{
		repo:         repo,
		xenditPayout: xenditPayout,
	}
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

func (uc *financeUseCase) GetFinancialReport(ctx context.Context, month string) (*domain.FinanceReportResponse, error) {
	period, monthLabel, err := uc.parseMonth(month)
	if err != nil {
		return nil, err
	}

	summary, err := uc.repo.GetReportSummary(ctx, period)
	if err != nil {
		return nil, err
	}

	dailyIncome, err := uc.repo.ListReportDailyIncome(ctx, period)
	if err != nil {
		return nil, err
	}

	return &domain.FinanceReportResponse{
		Month:       monthLabel,
		Summary:     *summary,
		DailyIncome: dailyIncome,
	}, nil
}

func (uc *financeUseCase) ListTransactions(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceTransactionItem, *domain.PaginationMeta, error) {
	page, limit := normalizePaging(params.Page, params.Limit)
	params.Page = page
	params.Limit = limit

	if params.Status != "" {
		normalized := strings.ToLower(strings.TrimSpace(params.Status))
		if !isValidPaymentStatus(normalized) {
			return nil, nil, ErrInvalidPaymentStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseListPeriod(params)
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

func (uc *financeUseCase) GetTransactionSummary(ctx context.Context, month string) (*domain.FinanceTransactionSummary, error) {
	period, _, err := uc.parseMonth(month)
	if err != nil {
		return nil, err
	}

	return uc.repo.GetTransactionSummary(ctx, period)
}

func (uc *financeUseCase) ExportTransactions(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceTransactionItem, error) {
	if params.Status != "" {
		normalized := strings.ToLower(strings.TrimSpace(params.Status))
		if !isValidPaymentStatus(normalized) {
			return nil, ErrInvalidPaymentStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseListPeriod(params)
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

	period, _, err := uc.parseListPeriod(params)
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

func (uc *financeUseCase) GetPayoutSummary(ctx context.Context, month string) (*domain.FinancePayoutSummary, error) {
	period, _, err := uc.parseMonth(month)
	if err != nil {
		return nil, err
	}

	return uc.repo.GetPayoutSummary(ctx, period)
}

func (uc *financeUseCase) ExportPayouts(ctx context.Context, params domain.FinanceListParams) ([]domain.FinancePayoutItem, error) {
	if params.Status != "" {
		normalized := normalizePayoutStatus(params.Status)
		if !isValidPayoutStatus(normalized) {
			return nil, ErrInvalidPayoutStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseListPeriod(params)
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

	period, _, err := uc.parseListPeriod(params)
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

func (uc *financeUseCase) GetRefundSummary(ctx context.Context, month string) (*domain.FinanceRefundSummary, error) {
	period, _, err := uc.parseMonth(month)
	if err != nil {
		return nil, err
	}

	return uc.repo.GetRefundSummary(ctx, period)
}

func (uc *financeUseCase) ExportRefunds(ctx context.Context, params domain.FinanceListParams) ([]domain.FinanceRefundItem, error) {
	if params.Status != "" {
		normalized := normalizeRefundStatus(params.Status)
		if !isValidRefundStatus(normalized) {
			return nil, ErrInvalidRefundStatus
		}
		params.Status = normalized
	}

	period, _, err := uc.parseListPeriod(params)
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

func (uc *financeUseCase) UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, req domain.UpdateRefundStatusRequest, actorID uuid.UUID) (*domain.StatusActionResponse, error) {
	normalized := normalizeRefundStatus(req.Status)
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

	// Option 1 policy:
	// once an order has been settled into a completed payout batch, refund cannot proceed.
	if normalized == domain.RefundStatusApproved || normalized == domain.RefundStatusProcessed {
		isSettled, err := uc.repo.IsOrderSettlementCompleted(ctx, record.OrderID)
		if err != nil {
			return nil, err
		}
		if isSettled {
			return nil, ErrRefundBlockedByCompletedSettlement
		}
	}

	if normalized == domain.RefundStatusProcessed {
		if record.Status == domain.RefundStatusProcessed {
			return &domain.StatusActionResponse{
				ID:     refundID,
				Status: domain.RefundStatusProcessed,
			}, nil
		}
		return uc.processRefundDisbursement(ctx, refundID, req, actorID)
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

func (uc *financeUseCase) HandleRefundPayoutWebhook(ctx context.Context, payload domain.XenditPayoutWebhookPayload) error {
	referenceID := strings.TrimSpace(payload.Data.ReferenceID)
	if referenceID == "" {
		return nil
	}

	payoutStatus := normalizeXenditPayoutWebhookStatus(payload.Event, payload.Data.Status)
	var failedReason *string
	if fr := buildPayoutFailureReason(payload); fr != nil {
		failedReason = fr
	} else if strings.EqualFold(payload.Event, domain.XenditPayoutWebhookEventFailed) || strings.EqualFold(payload.Event, domain.XenditPayoutWebhookEventReversed) {
		msg := fmt.Sprintf("xendit payout event: %s", strings.TrimSpace(payload.Event))
		failedReason = &msg
	}

	completed := false
	if strings.EqualFold(payload.Event, domain.XenditPayoutWebhookEventSucceeded) || strings.EqualFold(payload.Data.Status, domain.XenditPayoutStatusSucceeded) {
		completed = true
	}

	var completedAt *time.Time
	if completed {
		now := time.Now().UTC()
		completedAt = &now
	}

	_, err := uc.repo.ApplyRefundPayoutWebhookUpdate(
		ctx,
		referenceID,
		strings.TrimSpace(payload.Data.ID),
		payoutStatus,
		failedReason,
		completedAt,
	)
	return err
}

func (uc *financeUseCase) HandleRefundGatewayWebhook(ctx context.Context, payload domain.XenditRefundWebhookPayload) error {
	referenceID := strings.TrimSpace(payload.Data.ReferenceID)
	if referenceID == "" {
		return nil
	}

	xenditRefundID := strings.TrimSpace(payload.Data.ID)
	status := strings.TrimSpace(payload.Data.Status)
	var failureCode *string
	if fc := strings.TrimSpace(payload.Data.FailureCode); fc != "" {
		failureCode = &fc
	}

	updated, err := uc.repo.ApplyRefundGatewayWebhookUpdate(ctx, referenceID, xenditRefundID, status, failureCode)
	if err != nil {
		return fmt.Errorf("failed to apply refund gateway webhook: %w", err)
	}

	// If refund succeeded, update order payment_status to refunded
	if updated && status == domain.XenditRefundStatusSucceeded {
		// Extract order ID from reference ID (format: "refund-<uuid>")
		refundIDStr := strings.TrimPrefix(referenceID, "refund-")
		refundID, err := uuid.Parse(refundIDStr)
		if err != nil {
			return nil // non-critical
		}

		record, err := uc.repo.GetRefundRecord(ctx, refundID)
		if err != nil || record == nil {
			return nil // non-critical
		}

		if err := uc.repo.UpdateOrderPaymentStatus(ctx, record.OrderID, domain.PaymentStatusRefunded); err != nil {
			fmt.Printf("WARN: failed to update order payment status for refund %s: %v\n", refundID, err)
		}
	}

	return nil
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

func (uc *financeUseCase) processRefundDisbursement(ctx context.Context, refundID uuid.UUID, req domain.UpdateRefundStatusRequest, actorID uuid.UUID) (*domain.StatusActionResponse, error) {
	refundCtx, err := uc.repo.GetRefundDisbursementContext(ctx, refundID)
	if err != nil {
		return nil, err
	}
	if refundCtx == nil {
		return nil, ErrRefundNotFound
	}

	if isPayoutStatusProcessing(refundCtx.PayoutStatus) {
		return nil, ErrRefundPayoutInProgress
	}

	strategy := resolveRefundStrategy(refundCtx.PaymentMethod, refundCtx.PaymentChannel)
	if strategy == refundStrategyUnsupported {
		return nil, ErrRefundStrategyNotSupported
	}

	useQRFallback := true
	if req.UseDisbursementFallbackForQR != nil {
		useQRFallback = *req.UseDisbursementFallbackForQR
	}
	useEWalletFallback := true
	if req.UseDisbursementFallbackForEWallet != nil {
		useEWalletFallback = *req.UseDisbursementFallbackForEWallet
	}
	effectiveStrategy := strategy
	if strategy == refundStrategyQRGateway && !useQRFallback {
		return nil, ErrRefundStrategyNotSupported
	}
	if strategy == refundStrategyEWalletGateway && !useEWalletFallback {
		return nil, ErrRefundStrategyNotSupported
	}
	if strategy == refundStrategyQRGateway && useQRFallback {
		effectiveStrategy = refundStrategyQRDisbursementFallback
	}
	if strategy == refundStrategyEWalletGateway && useEWalletFallback {
		effectiveStrategy = refundStrategyEWalletDisbursementFallback
	}

	// Manual ON_HOLD policy:
	// if vendor balance is insufficient, hold refund and stop before calling gateway.
	balanceSnapshot, err := uc.repo.GetRefundVendorBalanceSnapshot(ctx, refundCtx.OrderID)
	if err != nil {
		return nil, err
	}
	if balanceSnapshot != nil && !balanceSnapshot.Sufficient {
		reason := fmt.Sprintf(
			"ON_HOLD: insufficient vendor balance. required %.2f from %s, available %.2f",
			balanceSnapshot.RequiredAmount,
			balanceSnapshot.BalanceSource,
			balanceSnapshot.AvailableAmount,
		)
		if err := uc.repo.SetRefundPayoutFailed(ctx, refundID, "ON_HOLD", reason); err != nil {
			return nil, err
		}
		return &domain.StatusActionResponse{
			ID:     refundID,
			Status: domain.RefundStatusApproved,
		}, nil
	}

	channelCode := strings.TrimSpace(req.DestinationChannelCode)
	if channelCode == "" {
		channelCode = strings.TrimSpace(derefString(refundCtx.DestinationChannelCode))
	}
	if channelCode == "" {
		channelCode = domain.PayoutChannelIDBCA
	}
	bankName := strings.TrimSpace(req.DestinationBankName)
	if bankName == "" {
		bankName = strings.TrimSpace(derefString(refundCtx.DestinationBankName))
	}
	accountNumber := strings.TrimSpace(req.DestinationAccountNumber)
	if accountNumber == "" {
		accountNumber = strings.TrimSpace(derefString(refundCtx.DestinationAccountNumber))
	}
	accountHolderName := strings.TrimSpace(req.DestinationAccountHolderName)
	if accountHolderName == "" {
		accountHolderName = strings.TrimSpace(derefString(refundCtx.DestinationAccountHolderName))
	}
	if accountNumber == "" || accountHolderName == "" {
		return nil, ErrRefundDestinationRequired
	}

	accountLast4 := maskAccountLast4(accountNumber)
	referenceID := fmt.Sprintf("refund-%s", refundID.String())
	idempotencyKey := fmt.Sprintf("%s-%d", referenceID, time.Now().UTC().UnixNano())
	description := fmt.Sprintf("Refund for order %s", refundCtx.OrderID.String())

	payoutReq := domain.XenditPayoutRequest{
		ReferenceID: referenceID,
		ChannelCode: channelCode,
		ChannelProperties: domain.XenditPayoutChannelProperties{
			AccountNumber:     accountNumber,
			AccountHolderName: accountHolderName,
		},
		Amount:      refundCtx.Amount,
		Description: description,
		Currency:    refundCtx.Currency,
	}

	payoutResp, payoutErr := uc.xenditPayout.CreatePayout(ctx, derefString(refundCtx.VendorXenditAccountID), idempotencyKey, payoutReq)
	if payoutErr != nil {
		_ = uc.repo.SetRefundPayoutFailed(ctx, refundID, "FAILED", payoutErr.Error())
		return nil, fmt.Errorf("%w: %v", ErrRefundDisbursementFailed, payoutErr)
	}

	payoutID := strings.TrimSpace(payoutResp.ID)
	payoutStatus := strings.TrimSpace(payoutResp.Status)
	if payoutStatus == "" {
		payoutStatus = "PROCESSING"
	}

	if err := uc.repo.SetRefundPayoutInitiated(
		ctx,
		refundID,
		actorID,
		effectiveStrategy,
		referenceID,
		channelCode,
		bankName,
		accountHolderName,
		accountLast4,
		&payoutID,
		payoutStatus,
	); err != nil {
		return nil, err
	}

	if strings.EqualFold(payoutStatus, domain.XenditPayoutStatusSucceeded) {
		now := time.Now().UTC()
		if _, err := uc.repo.ApplyRefundPayoutWebhookUpdate(ctx, referenceID, payoutID, payoutStatus, nil, &now); err != nil {
			return nil, err
		}
		return &domain.StatusActionResponse{
			ID:     refundID,
			Status: domain.RefundStatusProcessed,
		}, nil
	}

	// Refund remains approved until payout webhook marks it succeeded.
	return &domain.StatusActionResponse{
		ID:     refundID,
		Status: domain.RefundStatusApproved,
	}, nil
}

func (uc *financeUseCase) parseMonth(month string) (domain.FinancePeriod, string, error) {
	trimmed := strings.TrimSpace(month)
	if trimmed == "" {
		now := time.Now().In(wibLocation)
		trimmed = now.Format("2006-01")
	}

	parsed, err := time.ParseInLocation("2006-01", trimmed, wibLocation)
	if err != nil {
		return domain.FinancePeriod{}, "", ErrInvalidMonth
	}

	startLocal := time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, wibLocation)
	endLocal := startLocal.AddDate(0, 1, 0)

	return domain.FinancePeriod{
		Start: startLocal.UTC(),
		End:   endLocal.UTC(),
	}, parsed.Format("2006-01"), nil
}

func (uc *financeUseCase) parseListPeriod(params domain.FinanceListParams) (domain.FinancePeriod, string, error) {
	monthPeriod, monthLabel, err := uc.parseMonth(params.Month)
	if err != nil {
		return domain.FinancePeriod{}, "", err
	}

	startLocal := monthPeriod.Start.In(wibLocation)
	endLocal := monthPeriod.End.In(wibLocation)

	var fromStart *time.Time
	if strings.TrimSpace(params.DateFrom) != "" {
		v, err := parseDateStartWIB(params.DateFrom)
		if err != nil {
			return domain.FinancePeriod{}, "", err
		}
		fromStart = &v
	}

	var toStart *time.Time
	if strings.TrimSpace(params.DateTo) != "" {
		v, err := parseDateStartWIB(params.DateTo)
		if err != nil {
			return domain.FinancePeriod{}, "", err
		}
		toStart = &v
	}

	if fromStart != nil && toStart != nil && fromStart.After(*toStart) {
		return domain.FinancePeriod{}, "", ErrInvalidDateRange
	}

	if fromStart != nil && fromStart.After(startLocal) {
		startLocal = *fromStart
	}
	if toStart != nil {
		toEndLocal := toStart.AddDate(0, 0, 1)
		if toEndLocal.Before(endLocal) {
			endLocal = toEndLocal
		}
	}

	if !startLocal.Before(endLocal) {
		endLocal = startLocal
	}

	return domain.FinancePeriod{
		Start: startLocal.UTC(),
		End:   endLocal.UTC(),
	}, monthLabel, nil
}

func parseDateStartWIB(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	parsed, err := time.ParseInLocation("2006-01-02", trimmed, wibLocation)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, wibLocation), nil
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
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "pending" {
		return domain.RefundStatusRequested
	}
	return normalized
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
	normalized := strings.ToLower(strings.TrimSpace(status))
	switch {
	case normalized == "on hold", normalized == "on_hold":
		return domain.PayoutStatusOnHold
	case normalized == "complete", normalized == "completed":
		return domain.PayoutStatusComplete
	default:
		return normalized
	}
}

func isValidPayoutStatus(status string) bool {
	switch status {
	case domain.PayoutStatusSchedule,
		domain.PayoutStatusComplete,
		domain.PayoutStatusFailed,
		domain.PayoutStatusOnHold:
		return true
	default:
		return false
	}
}

const (
	refundStrategyVADisbursement              = "va_disbursement"
	refundStrategyQRGateway                   = "qr_gateway_refund"
	refundStrategyEWalletGateway              = "ewallet_gateway_refund"
	refundStrategyQRDisbursementFallback      = "qr_disbursement_fallback"
	refundStrategyEWalletDisbursementFallback = "ewallet_disbursement_fallback"
	refundStrategyDisbursementDefault         = "disbursement_default"
	refundStrategyUnsupported                 = "unsupported"
)

func resolveRefundStrategy(paymentMethod, paymentChannel *string) string {
	normMethod := normalizePaymentToken(derefString(paymentMethod))
	normChannel := normalizePaymentToken(derefString(paymentChannel))
	combined := normMethod + " " + normChannel

	if strings.Contains(combined, "QRIS") || strings.Contains(combined, "QR CODE") || strings.Contains(combined, "QR_CODE") || strings.Contains(combined, "QR") {
		return refundStrategyQRGateway
	}

	if strings.Contains(combined, "EWALLET") ||
		strings.Contains(combined, "E WALLET") ||
		strings.Contains(combined, "DANA") ||
		strings.Contains(combined, "OVO") ||
		strings.Contains(combined, "GOPAY") ||
		strings.Contains(combined, "SHOPEEPAY") ||
		strings.Contains(combined, "LINKAJA") ||
		strings.Contains(combined, "ASTRAPAY") {
		return refundStrategyEWalletGateway
	}

	if strings.Contains(combined, "VIRTUAL ACCOUNT") ||
		strings.Contains(combined, "VIRTUAL_ACCOUNT") ||
		strings.Contains(combined, "BANK TRANSFER") ||
		strings.Contains(combined, "BANK_TRANSFER") ||
		strings.Contains(combined, "VA") {
		return refundStrategyVADisbursement
	}

	if combined == "" {
		return refundStrategyUnsupported
	}
	return refundStrategyDisbursementDefault
}

func normalizePaymentToken(raw string) string {
	raw = strings.ToUpper(strings.TrimSpace(raw))
	raw = strings.ReplaceAll(raw, "-", " ")
	raw = strings.ReplaceAll(raw, "/", " ")
	raw = strings.Join(strings.Fields(raw), " ")
	return raw
}

func maskAccountLast4(accountNumber string) string {
	digits := strings.TrimSpace(accountNumber)
	if len(digits) <= 4 {
		return digits
	}
	return digits[len(digits)-4:]
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func isPayoutStatusProcessing(status *string) bool {
	s := strings.ToUpper(strings.TrimSpace(derefString(status)))
	switch s {
	case "PROCESSING", "PENDING", "SCHEDULED", "ACCEPTED":
		return true
	default:
		return false
	}
}
