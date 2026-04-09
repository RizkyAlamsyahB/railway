package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// --- VendorOrderRepository ---

type vendorOrderRepository struct {
	db *gorm.DB
}

// NewVendorOrderRepository creates a new VendorOrderRepository backed by GORM.
func NewVendorOrderRepository(db *gorm.DB) domain.VendorOrderRepository {
	return &vendorOrderRepository{db: db}
}

func (r *vendorOrderRepository) ListByVendor(ctx context.Context, vendorID uuid.UUID, params domain.VendorOrderListParams) ([]domain.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&orderModel{}).Where("orders.vendor_id = ?", vendorID.String())

	// Status filter
	if params.Status != "" {
		q = q.Where("orders.order_status = ?", params.Status)
	}

	// Date filter
	if params.DateFrom != nil {
		q = q.Where("orders.placed_at >= ?", *params.DateFrom)
	}
	if params.DateTo != nil {
		// Include the entire day by adding 1 day
		endOfDay := params.DateTo.AddDate(0, 0, 1)
		q = q.Where("orders.placed_at < ?", endOfDay)
	}

	// Search filter (by order_no or customer name from users table)
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		q = q.Joins("LEFT JOIN users ON users.id = orders.user_id").
			Where("LOWER(orders.order_no) LIKE ? OR LOWER(users.full_name) LIKE ?", searchPattern, searchPattern)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderClause := r.buildOrderClause(params.SortBy, params.SortOrder)
	q = q.Order(orderClause)

	offset := (params.Page - 1) * params.Limit
	var models []orderModel
	if err := q.
		Offset(offset).
		Limit(params.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomainOrder(&m)
	}

	return orders, total, nil
}

// buildOrderClause builds the ORDER BY clause based on sort parameters.
func (r *vendorOrderRepository) buildOrderClause(sortBy, sortOrder string) string {
	// Default sort
	column := "orders.placed_at"
	direction := "DESC"

	// Map sortBy to actual column
	switch sortBy {
	case domain.VendorOrderSortByOrderDate:
		column = "orders.placed_at"
	case domain.VendorOrderSortByTotalPayment:
		column = "orders.grand_total"
	case domain.VendorOrderSortByCustomerName:
		// This requires join with users table which is done in search
		// For now, we sort by placed_at as fallback for customer_name sorting
		// A proper implementation would require subquery or join
		column = "orders.placed_at"
	}

	// Validate direction
	if strings.ToLower(sortOrder) == domain.VendorOrderSortOrderAsc {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func (r *vendorOrderRepository) FindByIDAndVendor(ctx context.Context, orderID, vendorID uuid.UUID) (*domain.Order, error) {
	var m orderModel
	if err := r.db.WithContext(ctx).
		Where("id = ? AND vendor_id = ?", orderID.String(), vendorID.String()).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainOrder(&m), nil
}

// --- Dashboard queries ---

func (r *vendorOrderRepository) DashboardOrderStats(ctx context.Context, vendorID uuid.UUID, periodStart, periodEnd time.Time) (int64, int64, error) {
	type statsRow struct {
		Total      int64 `gorm:"column:total"`
		Successful int64 `gorm:"column:successful"`
	}
	var row statsRow
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select(`
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE order_status IN ('completed','received')) AS successful
		`).
		Where("vendor_id = ?", vendorID.String()).
		Where("placed_at >= ? AND placed_at < ?", periodStart, periodEnd).
		Take(&row).Error; err != nil {
		return 0, 0, err
	}
	return row.Total, row.Successful, nil
}

func (r *vendorOrderRepository) DashboardTodayTransactions(ctx context.Context, vendorID uuid.UUID, today time.Time, limit int) ([]domain.VendorDashboardTransaction, error) {
	type txRow struct {
		DateTime    time.Time `gorm:"column:date_time"`
		Invoice     string    `gorm:"column:invoice"`
		ProductName string    `gorm:"column:product_name"`
		Category    string    `gorm:"column:category"`
		Price       float64   `gorm:"column:price"`
		Status      string    `gorm:"column:status"`
	}

	tomorrow := today.AddDate(0, 0, 1)
	var rows []txRow
	if err := r.db.WithContext(ctx).
		Table("orders o").
		Select(`
			o.placed_at AS date_time,
			o.order_no AS invoice,
			COALESCE((SELECT oi.product_name_snapshot FROM order_items oi WHERE oi.order_id = o.id LIMIT 1), '-') AS product_name,
			COALESCE((SELECT c.name FROM order_items oi2 JOIN product_variants pv ON pv.id = oi2.product_variant_id JOIN products p ON p.id = pv.product_id JOIN categories c ON c.id = p.category_id WHERE oi2.order_id = o.id LIMIT 1), '-') AS category,
			o.grand_total AS price,
			o.order_status AS status
		`).
		Where("o.vendor_id = ?", vendorID.String()).
		Where("o.placed_at >= ? AND o.placed_at < ?", today, tomorrow).
		Order("o.placed_at DESC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]domain.VendorDashboardTransaction, len(rows))
	for i, r := range rows {
		result[i] = domain.VendorDashboardTransaction{
			DateTime:    r.DateTime,
			Invoice:     r.Invoice,
			ProductName: r.ProductName,
			Category:    r.Category,
			Price:       r.Price,
			Status:      r.Status,
		}
	}
	return result, nil
}

func (r *vendorOrderRepository) DashboardPaymentFlow(ctx context.Context, vendorID uuid.UUID, start, end time.Time) ([]domain.VendorPaymentFlowItem, error) {
	type flowRow struct {
		Label string `gorm:"column:label"`
		Count int64  `gorm:"column:cnt"`
	}
	var rows []flowRow
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select(`
			CASE
				WHEN order_status IN ('completed','received') THEN 'Selesai'
				WHEN order_status IN ('paid','processing','packed','shipped') THEN 'Dalam Proses'
				WHEN order_status IN ('canceled','refunded') THEN 'Dibatalkan'
				ELSE 'Menunggu Pembayaran'
			END AS label,
			COUNT(*) AS cnt
		`).
		Where("vendor_id = ?", vendorID.String()).
		Where("placed_at >= ? AND placed_at < ?", start, end).
		Group("label").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	var total int64
	for _, r := range rows {
		total += r.Count
	}

	result := make([]domain.VendorPaymentFlowItem, len(rows))
	for i, r := range rows {
		pct := float64(0)
		if total > 0 {
			pct = float64(r.Count) / float64(total) * 100
		}
		result[i] = domain.VendorPaymentFlowItem{
			Label: r.Label,
			Value: math.Round(pct*100) / 100,
		}
	}
	return result, nil
}
