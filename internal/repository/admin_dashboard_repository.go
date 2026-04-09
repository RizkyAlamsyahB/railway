package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type adminDashboardRepository struct {
	db *gorm.DB
}

// NewAdminDashboardRepository creates a new admin dashboard repository backed by GORM.
func NewAdminDashboardRepository(db *gorm.DB) domain.AdminDashboardRepository {
	return &adminDashboardRepository{db: db}
}

func (r *adminDashboardRepository) GetSummary(ctx context.Context, params domain.AdminDashboardParams) (*domain.AdminDashboardSummary, error) {
	type row struct {
		PPIUAndPIHKCount    int64 `gorm:"column:ppiu_and_pihk_count"`
		DormitoryCount      int64 `gorm:"column:dormitory_count"`
		StoreCount          int64 `gorm:"column:store_count"`
		PendingPaymentCount int64 `gorm:"column:pending_payment_count"`
		NeedsReviewCount    int64 `gorm:"column:needs_review_count"`
		OrderCount          int64 `gorm:"column:order_count"`
	}

	periodStart := time.Date(params.Year, time.Month(params.Month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0)

	var result row
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			(SELECT COUNT(*) FROM vendors WHERE status = 'active' AND vendor_type = 'ppiu') AS ppiu_and_pihk_count,
			(SELECT COUNT(*) FROM vendors WHERE status = 'active' AND vendor_type = 'hajj_dormitory') AS dormitory_count,
			(SELECT COUNT(*) FROM vendors WHERE status = 'active' AND vendor_type = 'souvenir_store') AS store_count,
			(SELECT COUNT(*) FROM payment_invoices WHERE status = 'pending' AND created_at >= ? AND created_at < ?) AS pending_payment_count,
			(SELECT COUNT(*) FROM vendors WHERE status = 'submitted') AS needs_review_count,
			(SELECT COUNT(*) FROM orders WHERE placed_at >= ? AND placed_at < ?) AS order_count
	`, periodStart, periodEnd, periodStart, periodEnd).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &domain.AdminDashboardSummary{
		PPIUAndPIHKCount:    result.PPIUAndPIHKCount,
		DormitoryCount:      result.DormitoryCount,
		StoreCount:          result.StoreCount,
		PendingPaymentCount: result.PendingPaymentCount,
		NeedsReviewCount:    result.NeedsReviewCount,
		OrderCount:          result.OrderCount,
	}, nil
}

func (r *adminDashboardRepository) ListTransactions(ctx context.Context, params domain.AdminDashboardParams, limit int) ([]domain.AdminDashboardTransactionItem, error) {
	type row struct {
		DateTime     time.Time      `gorm:"column:date_time"`
		Invoice      sql.NullString `gorm:"column:invoice"`
		CustomerName string         `gorm:"column:customer_name"`
		ProductName  string         `gorm:"column:product_name"`
		Price        float64        `gorm:"column:price"`
		Status       sql.NullString `gorm:"column:status"`
	}

	periodStart := time.Date(params.Year, time.Month(params.Month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0)

	var rows []row
	if err := r.db.WithContext(ctx).
		Table("orders o").
		Joins("JOIN users u ON u.id = o.user_id").
		Joins("JOIN order_items oi ON oi.order_id = o.id").
		Joins("LEFT JOIN payment_invoices pi ON pi.order_id = o.id").
		Select(`
			o.placed_at AS date_time,
			COALESCE(NULLIF(pi.external_invoice_id, ''), o.order_no) AS invoice,
			u.full_name AS customer_name,
			oi.product_name_snapshot AS product_name,
			oi.unit_price AS price,
			CASE
				WHEN o.payment_status = 'refunded' THEN 'refunded'
				WHEN pi.status IS NOT NULL THEN pi.status
				ELSE o.order_status
			END AS status
		`).
		Where("o.placed_at >= ? AND o.placed_at < ?", periodStart, periodEnd).
		Order("o.placed_at DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]domain.AdminDashboardTransactionItem, len(rows))
	for i, row := range rows {
		invoice := ""
		if row.Invoice.Valid {
			invoice = row.Invoice.String
		}

		status := ""
		if row.Status.Valid {
			status = row.Status.String
		}

		items[i] = domain.AdminDashboardTransactionItem{
			DateTime:     row.DateTime,
			Invoice:      invoice,
			CustomerName: row.CustomerName,
			ProductName:  row.ProductName,
			Price:        row.Price,
			Status:       status,
		}
	}

	return items, nil
}
