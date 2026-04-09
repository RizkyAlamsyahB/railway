package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type adminReportRepository struct {
	db *gorm.DB
}

// NewAdminReportRepository creates a new admin report repository backed by GORM.
func NewAdminReportRepository(db *gorm.DB) domain.AdminReportRepository {
	return &adminReportRepository{db: db}
}

func (r *adminReportRepository) GetSummary(ctx context.Context, params domain.AdminReportParams) (*domain.AdminReportSummary, error) {
	type row struct {
		SuccessfulTransactions int64   `gorm:"column:successful_transactions"`
		TotalTransactions      int64   `gorm:"column:total_transactions"`
		CustomerCount          int64   `gorm:"column:customer_count"`
		PlatformCommission     float64 `gorm:"column:platform_commission"`
	}

	yearStart, yearEnd := adminReportYearRange(params.Year)

	var result row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE((
				SELECT COUNT(*)
				FROM payment_invoices pi
				JOIN orders o ON o.id = pi.order_id
				WHERE o.placed_at >= ? AND o.placed_at < ?
				  AND pi.status = 'paid'
				  AND o.payment_status <> 'refunded'
			), 0) AS successful_transactions,
			COALESCE((
				SELECT COUNT(*)
				FROM payment_invoices pi
				JOIN orders o ON o.id = pi.order_id
				WHERE o.placed_at >= ? AND o.placed_at < ?
			), 0) AS total_transactions,
			COALESCE((
				SELECT COUNT(*)
				FROM users u
				JOIN roles r ON r.id = u.role_id
				WHERE r.code = ?
			), 0) AS customer_count,
			COALESCE((
				SELECT SUM(o.platform_fee)
				FROM orders o
				WHERE o.placed_at >= ? AND o.placed_at < ?
			), 0) AS platform_commission
	`, yearStart, yearEnd, yearStart, yearEnd, domain.RoleCustomer, yearStart, yearEnd).Scan(&result).Error; err != nil {
		return nil, err
	}

	successRate := 0.0
	if result.TotalTransactions > 0 {
		successRate = (float64(result.SuccessfulTransactions) / float64(result.TotalTransactions)) * 100
	}

	return &domain.AdminReportSummary{
		SuccessTransactionRate: successRate,
		TotalTransactions:      result.TotalTransactions,
		CustomerCount:          result.CustomerCount,
		PlatformCommission:     result.PlatformCommission,
	}, nil
}

func (r *adminReportRepository) List(ctx context.Context, params domain.AdminReportParams) ([]domain.AdminReportItem, int64, error) {
	type row struct {
		Period           time.Time `gorm:"column:period_date"`
		PeriodLabel      string    `gorm:"column:period_label"`
		TotalTransaction float64   `gorm:"column:total_transaction"`
		DownPayment      float64   `gorm:"column:dp"`
		Settlement       float64   `gorm:"column:settlement"`
		Refund           float64   `gorm:"column:refund"`
	}

	yearStart, yearEnd := adminReportYearRange(params.Year)
	periodFilterClause, periodFilterArgs, err := adminReportPeriodFilter(params)
	if err != nil {
		return nil, 0, err
	}
	baseSQL := `
		WITH months AS (
			SELECT generate_series(?::timestamptz, (?::timestamptz - INTERVAL '1 month'), INTERVAL '1 month') AS period_date
		),
		transaction_totals AS (
			SELECT
				date_trunc('month', o.placed_at AT TIME ZONE 'UTC')::timestamptz AS period_date,
				COALESCE(SUM(o.grand_total), 0) AS total_transaction,
				COALESCE(SUM(CASE WHEN pi.status = 'pending' THEN pi.amount ELSE 0 END), 0) AS dp,
				COALESCE(SUM(CASE WHEN pi.status = 'paid' AND o.payment_status <> 'refunded' THEN pi.amount ELSE 0 END), 0) AS settlement
			FROM orders o
			LEFT JOIN payment_invoices pi ON pi.order_id = o.id
			WHERE o.placed_at >= ? AND o.placed_at < ?
			GROUP BY 1
		),
		refund_totals AS (
			SELECT
				date_trunc('month', r.requested_at AT TIME ZONE 'UTC')::timestamptz AS period_date,
				COALESCE(SUM(CASE WHEN r.status = 'processed' THEN r.amount ELSE 0 END), 0) AS refund
			FROM refunds r
			WHERE r.requested_at >= ? AND r.requested_at < ?
			GROUP BY 1
		)
		SELECT
			m.period_date AS period_date,
			TO_CHAR(m.period_date AT TIME ZONE 'UTC', 'Mon YYYY') AS period_label,
			COALESCE(tt.total_transaction, 0) AS total_transaction,
			COALESCE(tt.dp, 0) AS dp,
			COALESCE(tt.settlement, 0) AS settlement,
			COALESCE(rt.refund, 0) AS refund
		FROM months m
		LEFT JOIN transaction_totals tt ON tt.period_date = m.period_date
		LEFT JOIN refund_totals rt ON rt.period_date = m.period_date
	`

	args := []any{yearStart, yearEnd, yearStart, yearEnd, yearStart, yearEnd}
	conditions := make([]string, 0, 2)
	if periodFilterClause != "" {
		conditions = append(conditions, periodFilterClause)
		args = append(args, periodFilterArgs...)
	}
	if params.Search != "" {
		conditions = append(conditions, "TO_CHAR(m.period_date AT TIME ZONE 'UTC', 'Mon YYYY') ILIKE ?")
		args = append(args, "%"+params.Search+"%")
	}
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) AS total FROM (%s %s) report_rows", baseSQL, whereClause)
	type countRow struct {
		Total int64 `gorm:"column:total"`
	}
	var totalRow countRow
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&totalRow).Error; err != nil {
		return nil, 0, err
	}

	orderClause := r.adminReportOrderClause(params.SortBy, params.SortOrder)
	querySQL := fmt.Sprintf("%s %s ORDER BY %s", baseSQL, whereClause, orderClause)
	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		querySQL += " LIMIT ? OFFSET ?"
		args = append(args, params.Limit, offset)
	}

	var rows []row
	if err := r.db.WithContext(ctx).Raw(querySQL, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.AdminReportItem, len(rows))
	for i, row := range rows {
		items[i] = domain.AdminReportItem{
			Period:           strings.TrimSpace(row.PeriodLabel),
			TotalTransaction: row.TotalTransaction,
			DownPayment:      row.DownPayment,
			Settlement:       row.Settlement,
			Refund:           row.Refund,
		}
	}

	return items, totalRow.Total, nil
}

func (r *adminReportRepository) adminReportOrderClause(sortBy, sortOrder string) string {
	direction := "DESC"
	if strings.EqualFold(sortOrder, "asc") {
		direction = "ASC"
	}

	column := "period_date"
	switch sortBy {
	case "total_transaction":
		column = "total_transaction"
	case "dp":
		column = "dp"
	case "settlement":
		column = "settlement"
	case "refund":
		column = "refund"
	case "period":
		column = "period_date"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func adminReportYearRange(year int) (time.Time, time.Time) {
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	return start, end
}

func adminReportPeriodFilter(params domain.AdminReportParams) (string, []any, error) {
	var conditions []string
	var args []any

	if params.FromDate != "" {
		from, err := time.Parse("2006-01-02", params.FromDate)
		if err != nil {
			return "", nil, err
		}
		fromMonthStart := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
		conditions = append(conditions, "m.period_date >= ?")
		args = append(args, fromMonthStart)
	}

	if params.ToDate != "" {
		to, err := time.Parse("2006-01-02", params.ToDate)
		if err != nil {
			return "", nil, err
		}
		toMonthEndExclusive := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
		conditions = append(conditions, "m.period_date < ?")
		args = append(args, toMonthEndExclusive)
	}

	if len(conditions) == 0 {
		return "", nil, nil
	}

	return strings.Join(conditions, " AND "), args, nil
}
