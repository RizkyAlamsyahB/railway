package domain

import "context"

// AdminReportParams holds query parameters for the admin reports page.
type AdminReportParams struct {
	Year      int
	Page      int
	Limit     int
	Search    string
	SortBy    string
	SortOrder string
	FromDate  string
	ToDate    string
}

// AdminReportSummary holds the KPI cards for the admin report page.
type AdminReportSummary struct {
	SuccessTransactionRate float64 `json:"success_transaction_rate"`
	TotalTransactions      int64   `json:"total_transactions"`
	CustomerCount          int64   `json:"customer_count"`
	PlatformCommission     float64 `json:"platform_commission"`
}

// AdminReportItem represents one monthly row in the admin report table.
type AdminReportItem struct {
	Period           string  `json:"period"`
	TotalTransaction float64 `json:"total_transaction"`
	DownPayment      float64 `json:"dp"`
	Settlement       float64 `json:"settlement"`
	Refund           float64 `json:"refund"`
}

// AdminReportResponse is the output DTO for the admin report endpoint.
type AdminReportResponse struct {
	Year    int                `json:"year"`
	Summary AdminReportSummary `json:"summary"`
	Items   []AdminReportItem  `json:"items"`
}

// AdminReportRepository defines the data access contract for admin reports.
type AdminReportRepository interface {
	GetSummary(ctx context.Context, params AdminReportParams) (*AdminReportSummary, error)
	List(ctx context.Context, params AdminReportParams) ([]AdminReportItem, int64, error)
}

// AdminReportUseCase defines the business contract for admin reports.
type AdminReportUseCase interface {
	GetReport(ctx context.Context, params AdminReportParams) (*AdminReportResponse, *PaginationMeta, error)
	ExportReport(ctx context.Context, params AdminReportParams) ([]AdminReportItem, error)
}
