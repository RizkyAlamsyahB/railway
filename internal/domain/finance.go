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
	Month    string
	DateFrom string
	DateTo   string
	Page     int
	Limit    int
	Status   string
	Search   string
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
	GrossRevenueTotal    float64 `json:"gross_revenue_total"`
	PlatformCommission   float64 `json:"platform_commission"`
	VendorPayoutTotal    float64 `json:"vendor_payout_total"`
	TemporaryGrossProfit float64 `json:"temporary_gross_profit"`
}

// FinanceReportDailyItem represents daily aggregated income.
type FinanceReportDailyItem struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

// FinanceReportResponse is the output DTO for the financial report.
type FinanceReportResponse struct {
	Month       string                   `json:"month"`
	Summary     FinanceReportSummary     `json:"summary"`
	DailyIncome []FinanceReportDailyItem `json:"daily_income"`
}

// --- Transactions & Payments ---

// FinanceTransactionItem represents a row in the transactions & payments table.
type FinanceTransactionItem struct {
	Date     time.Time `json:"date"`
	Invoice  string    `json:"invoice"`
	Customer string    `json:"customer"`
	Vendor   *string   `json:"vendor"`
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
	Vendor         *string   `json:"vendor"`
	VendorType     *string   `json:"vendor_type"`
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
	Vendor   *string   `json:"vendor"`
	OrderID  uuid.UUID `json:"order_id"`
	Reason   *string   `json:"reason,omitempty"`
	Amount   float64   `json:"amount"`
	Status   string    `json:"status"`
}

// FinanceRefundSummary holds aggregated counts for the refund & dispute page.
type FinanceRefundSummary struct {
	SubmittedCount     int64   `json:"submitted_count"`
	DisputeActiveCount int64   `json:"dispute_active_count"`
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
	Status                            string `json:"status" binding:"required"`
	DestinationChannelCode            string `json:"destination_channel_code,omitempty"`
	DestinationBankName               string `json:"destination_bank_name,omitempty"`
	DestinationAccountNumber          string `json:"destination_account_number,omitempty"`
	DestinationAccountHolderName      string `json:"destination_account_holder_name,omitempty"`
	UseDisbursementFallbackForQR      *bool  `json:"use_disbursement_fallback_for_qr,omitempty"`
	UseDisbursementFallbackForEWallet *bool  `json:"use_disbursement_fallback_for_ewallet,omitempty"`
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

// RefundDisbursementContext holds data required to process refund payout.
type RefundDisbursementContext struct {
	RefundID                     uuid.UUID
	OrderID                      uuid.UUID
	Amount                       float64
	Currency                     string
	Status                       string
	PaymentMethod                *string
	PaymentChannel               *string
	VendorXenditAccountID        *string
	PayoutReferenceID            *string
	PayoutStatus                 *string
	DestinationChannelCode       *string
	DestinationBankName          *string
	DestinationAccountNumber     *string
	DestinationAccountHolderName *string
	DestinationAccountLast4      *string
}

// RefundVendorBalanceSnapshot is a pre-check snapshot for refund processing.
type RefundVendorBalanceSnapshot struct {
	Sufficient      bool
	BalanceSource   string
	RequiredAmount  float64
	AvailableAmount float64
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
	GetRefundDisbursementContext(ctx context.Context, refundID uuid.UUID) (*RefundDisbursementContext, error)
	GetRefundVendorBalanceSnapshot(ctx context.Context, orderID uuid.UUID) (*RefundVendorBalanceSnapshot, error)
	GetPayoutRecord(ctx context.Context, payoutID uuid.UUID) (*PayoutRecord, error)
	IsOrderSettlementCompleted(ctx context.Context, orderID uuid.UUID) (bool, error)
	SetRefundPayoutInitiated(ctx context.Context, refundID uuid.UUID, actorID uuid.UUID, strategy string, referenceID string, channelCode string, bankName string, accountHolderName string, accountLast4 string, payoutID *string, payoutStatus string) error
	SetRefundPayoutFailed(ctx context.Context, refundID uuid.UUID, payoutStatus string, failedReason string) error
	ApplyRefundPayoutWebhookUpdate(ctx context.Context, referenceID string, payoutID string, payoutStatus string, failedReason *string, completedAt *time.Time) (bool, error)
	UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
	UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
	SetRefundGatewayInitiated(ctx context.Context, refundID uuid.UUID, xenditRefundID string, refundMethod string, payoutReferenceID string) error
	ApplyRefundGatewayWebhookUpdate(ctx context.Context, referenceID string, xenditRefundID string, status string, failureCode *string) (bool, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderID uuid.UUID, paymentStatus string) error
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

	UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, req UpdateRefundStatusRequest, actorID uuid.UUID) (*StatusActionResponse, error)
	UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status string, actorID uuid.UUID) (*StatusActionResponse, error)
	HandleRefundPayoutWebhook(ctx context.Context, payload XenditPayoutWebhookPayload) error
	HandleRefundGatewayWebhook(ctx context.Context, payload XenditRefundWebhookPayload) error
}
