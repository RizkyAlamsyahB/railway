package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type vendorReportRepository struct {
	db *gorm.DB
}

// NewVendorReportRepository creates a new vendor report repository backed by GORM.
func NewVendorReportRepository(db *gorm.DB) domain.VendorReportRepository {
	return &vendorReportRepository{db: db}
}

func (r *vendorReportRepository) GetSummary(ctx context.Context, vendorID uuid.UUID, period domain.FinancePeriod) (*domain.VendorReportSummary, error) {
	type row struct {
		SuccessfulTransactions int64   `gorm:"column:successful_transactions"`
		TotalTransactions      int64   `gorm:"column:total_transactions"`
		GrossRevenue           float64 `gorm:"column:gross_revenue"`
		NetRevenue             float64 `gorm:"column:net_revenue"`
	}

	var result row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN pi.status = 'paid' AND o.payment_status <> 'refunded' THEN 1 ELSE 0 END), 0) AS successful_transactions,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(CASE WHEN pi.status = 'paid' AND o.payment_status <> 'refunded' THEN o.grand_total ELSE 0 END), 0) AS gross_revenue,
			COALESCE(SUM(CASE WHEN pi.status = 'paid' AND o.payment_status <> 'refunded' THEN o.grand_total - o.platform_fee ELSE 0 END), 0) AS net_revenue
		FROM orders o
		JOIN payment_invoices pi ON pi.order_id = o.id
		WHERE o.vendor_id = ?
		  AND o.placed_at >= ? AND o.placed_at < ?
	`, vendorID.String(), period.Start, period.End).Scan(&result).Error; err != nil {
		return nil, err
	}

	successRate := 0.0
	if result.TotalTransactions > 0 {
		successRate = (float64(result.SuccessfulTransactions) / float64(result.TotalTransactions)) * 100
	}

	return &domain.VendorReportSummary{
		SuccessTransactionRate: successRate,
		TotalTransactions:      result.TotalTransactions,
		GrossRevenue:           result.GrossRevenue,
		NetRevenue:             result.NetRevenue,
	}, nil
}

func (r *vendorReportRepository) ListDailyRevenue(ctx context.Context, vendorID uuid.UUID, period domain.FinancePeriod, params domain.VendorReportParams) ([]domain.VendorReportDailyItem, int64, error) {
	type row struct {
		Date               time.Time `gorm:"column:date"`
		TransactionCount   int64     `gorm:"column:transaction_count"`
		GrossRevenue       float64   `gorm:"column:gross_revenue"`
		PlatformCommission float64   `gorm:"column:platform_commission"`
		NetRevenue         float64   `gorm:"column:net_revenue"`
	}

	baseSQL := `
		SELECT
			DATE(o.placed_at AT TIME ZONE 'UTC') AS date,
			COUNT(*) AS transaction_count,
			COALESCE(SUM(o.grand_total), 0) AS gross_revenue,
			COALESCE(SUM(o.platform_fee), 0) AS platform_commission,
			COALESCE(SUM(o.grand_total - o.platform_fee), 0) AS net_revenue
		FROM orders o
		JOIN payment_invoices pi ON pi.order_id = o.id
		WHERE o.vendor_id = ?
		  AND o.placed_at >= ? AND o.placed_at < ?
		  AND pi.status = 'paid'
		  AND o.payment_status <> 'refunded'
	`

	args := []any{vendorID.String(), period.Start, period.End}

	searchClause := ""
	if strings.TrimSpace(params.Search) != "" {
		searchClause = " AND CAST(DATE(o.placed_at AT TIME ZONE 'UTC') AS TEXT) LIKE ?"
		args = append(args, "%"+strings.TrimSpace(params.Search)+"%")
	}

	groupOrderSQL := " GROUP BY DATE(o.placed_at AT TIME ZONE 'UTC') ORDER BY date DESC"

	// Count total distinct dates
	countSQL := fmt.Sprintf("SELECT COUNT(*) AS total FROM (%s%s%s) daily_rows", baseSQL, searchClause, groupOrderSQL)
	type countRow struct {
		Total int64 `gorm:"column:total"`
	}
	var totalRow countRow
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&totalRow).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated rows
	querySQL := baseSQL + searchClause + groupOrderSQL
	queryArgs := make([]any, len(args))
	copy(queryArgs, args)

	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		querySQL += " LIMIT ? OFFSET ?"
		queryArgs = append(queryArgs, params.Limit, offset)
	}

	var rows []row
	if err := r.db.WithContext(ctx).Raw(querySQL, queryArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.VendorReportDailyItem, len(rows))
	for i, r := range rows {
		items[i] = domain.VendorReportDailyItem{
			Date:               r.Date.Format("2006-01-02"),
			TransactionCount:   r.TransactionCount,
			GrossRevenue:       r.GrossRevenue,
			PlatformCommission: r.PlatformCommission,
			NetRevenue:         r.NetRevenue,
		}
	}

	return items, totalRow.Total, nil
}
