package domain

import (
	"context"
	"time"
)

// AdminDashboardParams holds the selected dashboard period.
type AdminDashboardParams struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// AdminDashboardSummary holds aggregated metrics for the admin dashboard.
type AdminDashboardSummary struct {
	PPIUAndPIHKCount    int64 `json:"ppiu_and_pihk_count"`
	DormitoryCount      int64 `json:"dormitory_count"`
	StoreCount          int64 `json:"store_count"`
	PendingPaymentCount int64 `json:"pending_payment_count"`
	NeedsReviewCount    int64 `json:"needs_review_count"`
	OrderCount          int64 `json:"order_count"`
}

// AdminDashboardTransactionItem represents one row in the admin dashboard transaction table.
type AdminDashboardTransactionItem struct {
	DateTime     time.Time `json:"date_time"`
	Invoice      string    `json:"invoice"`
	CustomerName string    `json:"customer_name"`
	ProductName  string    `json:"product_name"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"`
}

// AdminDashboardResponse is the output DTO for the admin dashboard.
type AdminDashboardResponse struct {
	Month        int                             `json:"month"`
	Year         int                             `json:"year"`
	Summary      AdminDashboardSummary           `json:"summary"`
	Transactions []AdminDashboardTransactionItem `json:"transactions"`
}

// AdminDashboardRepository defines the data access contract for the admin dashboard.
type AdminDashboardRepository interface {
	GetSummary(ctx context.Context, params AdminDashboardParams) (*AdminDashboardSummary, error)
	ListTransactions(ctx context.Context, params AdminDashboardParams, limit int) ([]AdminDashboardTransactionItem, error)
}

// AdminDashboardUseCase defines the business contract for the admin dashboard.
type AdminDashboardUseCase interface {
	GetDashboard(ctx context.Context, params AdminDashboardParams) (*AdminDashboardResponse, error)
}
