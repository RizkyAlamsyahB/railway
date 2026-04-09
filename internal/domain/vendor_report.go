package domain

import (
	"context"

	"github.com/google/uuid"
)

// VendorReportParams holds query parameters for the vendor revenue report page.
type VendorReportParams struct {
	Month    string
	DateFrom string
	DateTo   string
	Page     int
	Limit    int
	Search   string
}

// VendorReportSummary holds the KPI cards for the vendor revenue report page.
type VendorReportSummary struct {
	SuccessTransactionRate float64 `json:"success_transaction_rate"`
	TotalTransactions      int64   `json:"total_transactions"`
	GrossRevenue           float64 `json:"gross_revenue"`
	NetRevenue             float64 `json:"net_revenue"`
}

// VendorReportDailyItem represents one daily row in the vendor revenue report table.
type VendorReportDailyItem struct {
	Date               string  `json:"date"`
	TransactionCount   int64   `json:"transaction_count"`
	GrossRevenue       float64 `json:"gross_revenue"`
	PlatformCommission float64 `json:"platform_commission"`
	NetRevenue         float64 `json:"net_revenue"`
}

// VendorReportResponse is the output DTO for the vendor revenue report endpoint.
type VendorReportResponse struct {
	Summary VendorReportSummary     `json:"summary"`
	Items   []VendorReportDailyItem `json:"items"`
}

// VendorReportRepository defines the data access contract for vendor revenue reports.
type VendorReportRepository interface {
	GetSummary(ctx context.Context, vendorID uuid.UUID, period FinancePeriod) (*VendorReportSummary, error)
	ListDailyRevenue(ctx context.Context, vendorID uuid.UUID, period FinancePeriod, params VendorReportParams) ([]VendorReportDailyItem, int64, error)
}

// VendorReportUseCase defines the business contract for vendor revenue reports.
type VendorReportUseCase interface {
	GetReport(ctx context.Context, vendorID uuid.UUID, params VendorReportParams) (*VendorReportResponse, *PaginationMeta, error)
	ExportReport(ctx context.Context, vendorID uuid.UUID, params VendorReportParams) ([]VendorReportDailyItem, error)
}
