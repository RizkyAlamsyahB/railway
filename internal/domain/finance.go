package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// FinancePeriod represents a time range filter (UTC).
type FinancePeriod struct {
	Start time.Time
	End   time.Time
}

// FinanceListParams holds common query parameters for finance list endpoints.
type FinanceListParams struct {
	Month  string
	Page   int
	Limit  int
	Status string
	Search string
}

// --- Dashboard ---

// FinanceDashboardSummary holds aggregated finance metrics.
type FinanceDashboardSummary struct {
	TotalTransactions  float64 `json:"total_transactions"`
	PlatformCommission float64 `json:"platform_commission"`
	VendorPayout       float64 `json:"vendor_payout"`
	ProcessedRefund    float64 `json:"processed_refund"`
}

// FinanceTodayTransactionItem represents a row in the "Transaksi Hari Ini" table.
type FinanceTodayTransactionItem struct {
	DateTime time.Time `json:"date_time"`
	Invoice  string    `json:"invoice"`
	Product  string    `json:"product"`
	Category string    `json:"category"`
	Price    float64   `json:"price"`
	Status   string    `json:"status"`
}

// FinanceDashboardResponse is the output DTO for the finance dashboard.
type FinanceDashboardResponse struct {
	Month             string                        `json:"month"`
	Summary           FinanceDashboardSummary       `json:"summary"`
	TodayTransactions []FinanceTodayTransactionItem `json:"today_transactions"`
}

// --- Financial Reports ---

// FinanceReportSummary holds aggregated totals for the financial report.
type FinanceReportSummary struct {
	GrossRevenueTotal     float64 `json:"gross_revenue_total"`
	PlatformCommission    float64 `json:"platform_commission"`
	VendorPayoutTotal     float64 `json:"vendor_payout_total"`
	TemporaryGrossProfit  float64 `json:"temporary_gross_profit"`
}

// FinanceReportDailyItem represents daily aggregated income.
type FinanceReportDailyItem struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

// FinanceReportResponse is the output DTO for the financial report.
type FinanceReportResponse struct {
	Month        string                   `json:"month"`
	Summary      FinanceReportSummary     `json:"summary"`
	DailyIncome  []FinanceReportDailyItem `json:"daily_income"`
}

// --- Transactions & Payments ---

// FinanceTransactionItem represents a row in the transactions & payments table.
type FinanceTransactionItem struct {
	Date     time.Time `json:"date"`
	Invoice  string    `json:"invoice"`
	Customer string    `json:"customer"`
	Vendor   string    `json:"vendor"`
	Method   *string   `json:"method,omitempty"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"`
}

// FinanceTransactionSummary holds aggregated totals for the transactions page.
type FinanceTransactionSummary struct {
	TotalTransactions  float64 `json:"total_transactions"`
	PaidTotal          float64 `json:"paid_total"`
	PendingTotal       float64 `json:"pending_total"`
	FailedExpiredTotal float64 `json:"failed_expired_total"`
}

// --- Settlement & Payout ---

// FinancePayoutItem represents a row in the settlement & payout table.
type FinancePayoutItem struct {
	ID             uuid.UUID `json:"id"`
	Vendor         string    `json:"vendor"`
	VendorType     string    `json:"vendor_type"`
	OrderCompleted int       `json:"order_completed"`
	Nominal        float64   `json:"nominal"`
	Commission     float64   `json:"commission"`
	NetPayout      float64   `json:"net_payout"`
	Status         string    `json:"status"`
}

// FinancePayoutSummary holds aggregated totals for the settlement & payout page.
type FinancePayoutSummary struct {
	OnHoldTotal    float64 `json:"on_hold_total"`
	ScheduleTotal  float64 `json:"schedule_total"`
	CompletedTotal float64 `json:"completed_total"`
	FailedTotal    float64 `json:"failed_total"`
}

// --- Refund & Dispute ---

// FinanceRefundItem represents a row in the refund & dispute table.
type FinanceRefundItem struct {
	ID       uuid.UUID `json:"id"`
	Customer string    `json:"customer"`
	Vendor   string    `json:"vendor"`
	OrderID  uuid.UUID `json:"order_id"`
	Reason   *string   `json:"reason,omitempty"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"`
}

// FinanceRefundSummary holds aggregated counts for the refund & dispute page.
type FinanceRefundSummary struct {
	SubmittedCount    int64   `json:"submitted_count"`
	DisputeActiveCount int64  `json:"dispute_active_count"`
	ApprovedAmount     float64 `json:"approved_amount"`
	ProcessedCount     int64   `json:"processed_count"`
}

// RefundRecord is a minimal refund record used for status transitions.
type RefundRecord struct {
	ID      uuid.UUID
	OrderID uuid.UUID
	Status  string
}

// PayoutRecord is a minimal payout record used for status transitions.
type PayoutRecord struct {
	ID     uuid.UUID
	Status string
}

// --- Status update requests ---

// UpdateRefundStatusRequest is the input DTO for updating a refund status.
type UpdateRefundStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdatePayoutStatusRequest is the input DTO for updating a payout status.
type UpdatePayoutStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// StatusActionResponse is the output DTO for status update actions.
type StatusActionResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}

// --- Repository Interface ---

// FinanceRepository defines the interface for finance data access.
type FinanceRepository interface {
	GetDashboardSummary(ctx context.Context, period FinancePeriod) (*FinanceDashboardSummary, error)
	ListTodayTransactions(ctx context.Context, dayStart, dayEnd time.Time) ([]FinanceTodayTransactionItem, error)
	GetReportSummary(ctx context.Context, period FinancePeriod) (*FinanceReportSummary, error)
	ListReportDailyIncome(ctx context.Context, period FinancePeriod) ([]FinanceReportDailyItem, error)

	ListTransactions(ctx context.Context, period FinancePeriod, params FinanceListParams) ([]FinanceTransactionItem, int64, error)
	GetTransactionSummary(ctx context.Context, period FinancePeriod) (*FinanceTransactionSummary, error)
	ListPayouts(ctx context.Context, period FinancePeriod, params FinanceListParams) ([]FinancePayoutItem, int64, error)
	GetPayoutSummary(ctx context.Context, period FinancePeriod) (*FinancePayoutSummary, error)
	ListRefunds(ctx context.Context, period FinancePeriod, params FinanceListParams) ([]FinanceRefundItem, int64, error)
	GetRefundSummary(ctx context.Context, period FinancePeriod) (*FinanceRefundSummary, error)

	GetRefundRecord(ctx context.Context, refundID uuid.UUID) (*RefundRecord, error)
	GetPayoutRecord(ctx context.Context, payoutID uuid.UUID) (*PayoutRecord, error)
	UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
	UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
}

// --- Usecase Interface ---

// FinanceUseCase defines the interface for finance business operations.
type FinanceUseCase interface {
	GetDashboard(ctx context.Context, month string) (*FinanceDashboardResponse, error)
	GetFinancialReport(ctx context.Context, month string) (*FinanceReportResponse, error)
	ListTransactions(ctx context.Context, params FinanceListParams) ([]FinanceTransactionItem, *PaginationMeta, error)
	GetTransactionSummary(ctx context.Context, month string) (*FinanceTransactionSummary, error)
	ListPayouts(ctx context.Context, params FinanceListParams) ([]FinancePayoutItem, *PaginationMeta, error)
	GetPayoutSummary(ctx context.Context, month string) (*FinancePayoutSummary, error)
	ListRefunds(ctx context.Context, params FinanceListParams) ([]FinanceRefundItem, *PaginationMeta, error)
	GetRefundSummary(ctx context.Context, month string) (*FinanceRefundSummary, error)

	ExportTransactions(ctx context.Context, params FinanceListParams) ([]FinanceTransactionItem, error)
	ExportPayouts(ctx context.Context, params FinanceListParams) ([]FinancePayoutItem, error)
	ExportRefunds(ctx context.Context, params FinanceListParams) ([]FinanceRefundItem, error)

	UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
	UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
}
