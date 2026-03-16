package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type financeRepository struct {
	db *gorm.DB
}

// NewFinanceRepository creates a new FinanceRepository backed by GORM.
func NewFinanceRepository(db *gorm.DB) domain.FinanceRepository {
	return &financeRepository{db: db}
}

type sumResult struct {
	Total float64 `gorm:"column:total"`
}

func (r *financeRepository) GetDashboardSummary(ctx context.Context, period domain.FinancePeriod) (*domain.FinanceDashboardSummary, error) {
	var totalTransactions sumResult
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(grand_total), 0) AS total").
		Where("created_at >= ? AND created_at < ?", period.Start, period.End).
		Scan(&totalTransactions).Error; err != nil {
		return nil, err
	}

	var platformCommission sumResult
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(platform_fee), 0) AS total").
		Where("created_at >= ? AND created_at < ?", period.Start, period.End).
		Scan(&platformCommission).Error; err != nil {
		return nil, err
	}

	var vendorPayout sumResult
	if err := r.db.WithContext(ctx).
		Table("payout_batches").
		Select("COALESCE(SUM(total_net), 0) AS total").
		Where("period_end >= ? AND period_end < ?", period.Start, period.End).
		Scan(&vendorPayout).Error; err != nil {
		return nil, err
	}

	var processedRefund sumResult
	if err := r.db.WithContext(ctx).
		Table("refunds").
		Select("COALESCE(SUM(amount), 0) AS total").
		Where("status = ?", domain.RefundStatusProcessed).
		Where("processed_at >= ? AND processed_at < ?", period.Start, period.End).
		Scan(&processedRefund).Error; err != nil {
		return nil, err
	}

	return &domain.FinanceDashboardSummary{
		TotalTransactions:  totalTransactions.Total,
		PlatformCommission: platformCommission.Total,
		VendorPayout:       vendorPayout.Total,
		ProcessedRefund:    processedRefund.Total,
	}, nil
}

func (r *financeRepository) GetReportSummary(ctx context.Context, period domain.FinancePeriod) (*domain.FinanceReportSummary, error) {
	var grossRevenue sumResult
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(grand_total), 0) AS total").
		Where("placed_at >= ? AND placed_at < ?", period.Start, period.End).
		Scan(&grossRevenue).Error; err != nil {
		return nil, err
	}

	var platformCommission sumResult
	if err := r.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(platform_fee), 0) AS total").
		Where("placed_at >= ? AND placed_at < ?", period.Start, period.End).
		Scan(&platformCommission).Error; err != nil {
		return nil, err
	}

	var vendorPayout sumResult
	if err := r.db.WithContext(ctx).
		Table("payout_batches").
		Select("COALESCE(SUM(total_net), 0) AS total").
		Where("period_end >= ? AND period_end < ?", period.Start, period.End).
		Scan(&vendorPayout).Error; err != nil {
		return nil, err
	}

	return &domain.FinanceReportSummary{
		GrossRevenueTotal:    grossRevenue.Total,
		PlatformCommission:   platformCommission.Total,
		VendorPayoutTotal:    vendorPayout.Total,
		TemporaryGrossProfit: grossRevenue.Total,
	}, nil
}

func (r *financeRepository) ListReportDailyIncome(ctx context.Context, period domain.FinancePeriod) ([]domain.FinanceReportDailyItem, error) {
	type row struct {
		Date   time.Time `gorm:"column:date"`
		Amount float64   `gorm:"column:amount"`
	}

	var rows []row
	if err := r.db.WithContext(ctx).
		Table("orders o").
		Select(`
			DATE(o.placed_at AT TIME ZONE 'UTC') AS date,
			COALESCE(SUM(o.grand_total), 0) AS amount
		`).
		Where("o.placed_at >= ? AND o.placed_at < ?", period.Start, period.End).
		Group("date").
		Order("date").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]domain.FinanceReportDailyItem, len(rows))
	for i, r := range rows {
		items[i] = domain.FinanceReportDailyItem{
			Date:   r.Date,
			Amount: r.Amount,
		}
	}

	return items, nil
}

func (r *financeRepository) ListTodayTransactions(ctx context.Context, dayStart, dayEnd time.Time) ([]domain.FinanceTodayTransactionItem, error) {
	type row struct {
		DateTime time.Time      `gorm:"column:date_time"`
		Invoice  sql.NullString `gorm:"column:invoice"`
		Product  string         `gorm:"column:product"`
		Category sql.NullString `gorm:"column:category"`
		Price    float64        `gorm:"column:price"`
		Status   sql.NullString `gorm:"column:status"`
	}

	var rows []row
	query := r.db.WithContext(ctx).
		Table("orders o").
		Joins("JOIN order_items oi ON oi.order_id = o.id").
		Joins("LEFT JOIN product_variants pv ON pv.id = oi.product_variant_id").
		Joins("LEFT JOIN products p ON p.id = pv.product_id").
		Joins("LEFT JOIN categories c ON c.id = p.category_id").
		Joins("LEFT JOIN payment_invoices pi ON pi.order_id = o.id").
		Select(`
			o.placed_at AS date_time,
			pi.external_invoice_id AS invoice,
			oi.product_name_snapshot AS product,
			c.name AS category,
			oi.unit_price AS price,
			CASE
				WHEN o.payment_status = 'refunded' THEN 'refunded'
				ELSE pi.status
			END AS status
		`).
		Where("o.placed_at >= ? AND o.placed_at < ?", dayStart, dayEnd).
		Order("o.placed_at DESC").
		Limit(10)

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]domain.FinanceTodayTransactionItem, len(rows))
	for i, r := range rows {
		invoice := ""
		if r.Invoice.Valid {
			invoice = r.Invoice.String
		}
		category := ""
		if r.Category.Valid {
			category = r.Category.String
		}
		status := ""
		if r.Status.Valid {
			status = r.Status.String
		}
		items[i] = domain.FinanceTodayTransactionItem{
			DateTime: r.DateTime,
			Invoice:  invoice,
			Product:  r.Product,
			Category: category,
			Price:    r.Price,
			Status:   status,
		}
	}

	return items, nil
}

func (r *financeRepository) ListTransactions(ctx context.Context, period domain.FinancePeriod, params domain.FinanceListParams) ([]domain.FinanceTransactionItem, int64, error) {
	base := r.transactionBaseQuery(ctx, period, params)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		Date     time.Time      `gorm:"column:date"`
		Invoice  string         `gorm:"column:invoice"`
		Customer string         `gorm:"column:customer"`
		Vendor   string         `gorm:"column:vendor"`
		Method   sql.NullString `gorm:"column:method"`
		Amount   float64        `gorm:"column:amount"`
		Status   string         `gorm:"column:status"`
	}

	query := base.Select(`
		o.created_at AS date,
		pi.external_invoice_id AS invoice,
		u.full_name AS customer,
		v.display_name AS vendor,
		pi.payment_method AS method,
		pi.amount AS amount,
		CASE
			WHEN o.payment_status = 'refunded' THEN 'refunded'
			ELSE pi.status
		END AS status
	`).Order("o.created_at DESC")

	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	var rows []row
	if err := query.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.FinanceTransactionItem, len(rows))
	for i, r := range rows {
		var method *string
		if r.Method.Valid {
			m := r.Method.String
			method = &m
		}
		items[i] = domain.FinanceTransactionItem{
			Date:     r.Date,
			Invoice:  r.Invoice,
			Customer: r.Customer,
			Vendor:   r.Vendor,
			Method:   method,
			Amount:   r.Amount,
			Status:   r.Status,
		}
	}

	return items, total, nil
}

func (r *financeRepository) GetTransactionSummary(ctx context.Context, period domain.FinancePeriod) (*domain.FinanceTransactionSummary, error) {
	type summaryRow struct {
		Total         float64 `gorm:"column:total"`
		Paid          float64 `gorm:"column:paid"`
		Pending       float64 `gorm:"column:pending"`
		FailedExpired float64 `gorm:"column:failed_expired"`
	}

	var row summaryRow
	if err := r.db.WithContext(ctx).
		Table("payment_invoices pi").
		Joins("JOIN orders o ON o.id = pi.order_id").
		Select(`
			COALESCE(SUM(pi.amount), 0) AS total,
			COALESCE(SUM(CASE WHEN pi.status = 'paid' AND o.payment_status <> 'refunded' THEN pi.amount ELSE 0 END), 0) AS paid,
			COALESCE(SUM(CASE WHEN pi.status = 'pending' THEN pi.amount ELSE 0 END), 0) AS pending,
			COALESCE(SUM(CASE WHEN pi.status IN ('failed', 'expired') THEN pi.amount ELSE 0 END), 0) AS failed_expired
		`).
		Where("o.created_at >= ? AND o.created_at < ?", period.Start, period.End).
		Scan(&row).Error; err != nil {
		return nil, err
	}

	return &domain.FinanceTransactionSummary{
		TotalTransactions:  row.Total,
		PaidTotal:          row.Paid,
		PendingTotal:       row.Pending,
		FailedExpiredTotal: row.FailedExpired,
	}, nil
}

func (r *financeRepository) transactionBaseQuery(ctx context.Context, period domain.FinancePeriod, params domain.FinanceListParams) *gorm.DB {
	query := r.db.WithContext(ctx).
		Table("payment_invoices pi").
		Joins("JOIN orders o ON o.id = pi.order_id").
		Joins("JOIN users u ON u.id = o.user_id").
		Joins("JOIN vendors v ON v.id = o.vendor_id").
		Where("o.created_at >= ? AND o.created_at < ?", period.Start, period.End)

	if params.Status != "" {
		if params.Status == domain.PaymentInvoiceStatusRefunded {
			query = query.Where("o.payment_status = ?", domain.PaymentInvoiceStatusRefunded)
		} else {
			query = query.Where("pi.status = ?", params.Status)
		}
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"(LOWER(pi.external_invoice_id) LIKE ? OR LOWER(u.full_name) LIKE ? OR LOWER(v.display_name) LIKE ?)",
			search, search, search,
		)
	}

	return query
}

func (r *financeRepository) ListPayouts(ctx context.Context, period domain.FinancePeriod, params domain.FinanceListParams) ([]domain.FinancePayoutItem, int64, error) {
	base := r.payoutBaseQuery(ctx, period, params)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID             string  `gorm:"column:id"`
		Vendor         string  `gorm:"column:vendor"`
		VendorType     string  `gorm:"column:vendor_type"`
		OrderCompleted int64   `gorm:"column:order_completed"`
		Nominal        float64 `gorm:"column:nominal"`
		Commission     float64 `gorm:"column:commission"`
		NetPayout      float64 `gorm:"column:net_payout"`
		Status         string  `gorm:"column:status"`
	}

	query := base.Select(`
		pb.id AS id,
		v.display_name AS vendor,
		v.vendor_type AS vendor_type,
		(
			SELECT COUNT(1)
			FROM orders o
			WHERE o.vendor_id = pb.vendor_id
				AND o.order_status = 'completed'
				AND o.placed_at >= pb.period_start
				AND o.placed_at < (pb.period_end + INTERVAL '1 day')
		) AS order_completed,
		pb.total_gross AS nominal,
		pb.total_fee AS commission,
		pb.total_net AS net_payout,
		pb.status AS status
	`).Order("pb.paid_at DESC")

	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	var rows []row
	if err := query.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.FinancePayoutItem, len(rows))
	for i, r := range rows {
		id, _ := uuid.Parse(r.ID)
		items[i] = domain.FinancePayoutItem{
			ID:             id,
			Vendor:         r.Vendor,
			VendorType:     r.VendorType,
			OrderCompleted: int(r.OrderCompleted),
			Nominal:        r.Nominal,
			Commission:     r.Commission,
			NetPayout:      r.NetPayout,
			Status:         r.Status,
		}
	}

	return items, total, nil
}

func (r *financeRepository) GetPayoutSummary(ctx context.Context, period domain.FinancePeriod) (*domain.FinancePayoutSummary, error) {
	type summaryRow struct {
		OnHold    float64 `gorm:"column:on_hold"`
		Schedule  float64 `gorm:"column:schedule"`
		Completed float64 `gorm:"column:completed"`
		Failed    float64 `gorm:"column:failed"`
	}

	var row summaryRow
	if err := r.db.WithContext(ctx).
		Table("payout_batches pb").
		Select(`
			COALESCE(SUM(CASE WHEN pb.status = 'on_hold' THEN pb.total_net ELSE 0 END), 0) AS on_hold,
			COALESCE(SUM(CASE WHEN pb.status = 'schedule' THEN pb.total_net ELSE 0 END), 0) AS schedule,
			COALESCE(SUM(CASE WHEN pb.status = 'completed' THEN pb.total_net ELSE 0 END), 0) AS completed,
			COALESCE(SUM(CASE WHEN pb.status = 'failed' THEN pb.total_net ELSE 0 END), 0) AS failed
		`).
		Where("pb.paid_at >= ? AND pb.paid_at < ?", period.Start, period.End).
		Scan(&row).Error; err != nil {
		return nil, err
	}

	return &domain.FinancePayoutSummary{
		OnHoldTotal:    row.OnHold,
		ScheduleTotal:  row.Schedule,
		CompletedTotal: row.Completed,
		FailedTotal:    row.Failed,
	}, nil
}

func (r *financeRepository) payoutBaseQuery(ctx context.Context, period domain.FinancePeriod, params domain.FinanceListParams) *gorm.DB {
	query := r.db.WithContext(ctx).
		Table("payout_batches pb").
		Joins("JOIN vendors v ON v.id = pb.vendor_id").
		Where("pb.paid_at >= ? AND pb.paid_at < ?", period.Start, period.End)

	if params.Status != "" {
		query = query.Where("pb.status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(v.display_name) LIKE ?", search)
	}

	return query
}

func (r *financeRepository) ListRefunds(ctx context.Context, period domain.FinancePeriod, params domain.FinanceListParams) ([]domain.FinanceRefundItem, int64, error) {
	base := r.refundBaseQuery(ctx, period, params)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID       string         `gorm:"column:id"`
		Customer string         `gorm:"column:customer"`
		Vendor   string         `gorm:"column:vendor"`
		OrderID  string         `gorm:"column:order_id"`
		Reason   sql.NullString `gorm:"column:reason"`
		Amount   float64        `gorm:"column:amount"`
		Status   string         `gorm:"column:status"`
	}

	query := base.Select(`
		r.id AS id,
		u.full_name AS customer,
		v.display_name AS vendor,
		r.order_id AS order_id,
		r.reason AS reason,
		r.amount AS amount,
		r.status AS status
	`).Order("r.requested_at DESC")

	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	var rows []row
	if err := query.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.FinanceRefundItem, len(rows))
	for i, r := range rows {
		id, _ := uuid.Parse(r.ID)
		orderID, _ := uuid.Parse(r.OrderID)
		var reason *string
		if r.Reason.Valid {
			reason = &r.Reason.String
		}
		items[i] = domain.FinanceRefundItem{
			ID:       id,
			Customer: r.Customer,
			Vendor:   r.Vendor,
			OrderID:  orderID,
			Reason:   reason,
			Amount:   r.Amount,
			Status:   r.Status,
		}
	}

	return items, total, nil
}

func (r *financeRepository) GetRefundSummary(ctx context.Context, period domain.FinancePeriod) (*domain.FinanceRefundSummary, error) {
	type summaryRow struct {
		Submitted      int64   `gorm:"column:submitted"`
		DisputeActive  int64   `gorm:"column:dispute_active"`
		ApprovedAmount float64 `gorm:"column:approved_amount"`
		Processed      int64   `gorm:"column:processed"`
	}

	var row summaryRow
	if err := r.db.WithContext(ctx).
		Table("refunds r").
		Joins("JOIN orders o ON o.id = r.order_id").
		Select(`
			COALESCE(SUM(CASE WHEN r.status IN ('requested', 'approved', 'rejected', 'processed') THEN 1 ELSE 0 END), 0) AS submitted,
			COALESCE(SUM(CASE WHEN r.status = 'requested' THEN 1 ELSE 0 END), 0) AS dispute_active,
			COALESCE(SUM(CASE WHEN r.status = 'approved' THEN r.amount ELSE 0 END), 0) AS approved_amount,
			COALESCE(SUM(CASE WHEN r.status = 'processed' THEN 1 ELSE 0 END), 0) AS processed
		`).
		Where("r.requested_at >= ? AND r.requested_at < ?", period.Start, period.End).
		Scan(&row).Error; err != nil {
		return nil, err
	}

	return &domain.FinanceRefundSummary{
		SubmittedCount:     row.Submitted,
		DisputeActiveCount: row.DisputeActive,
		ApprovedAmount:     row.ApprovedAmount,
		ProcessedCount:     row.Processed,
	}, nil
}

func (r *financeRepository) refundBaseQuery(ctx context.Context, period domain.FinancePeriod, params domain.FinanceListParams) *gorm.DB {
	query := r.db.WithContext(ctx).
		Table("refunds r").
		Joins("JOIN orders o ON o.id = r.order_id").
		Joins("JOIN users u ON u.id = o.user_id").
		Joins("JOIN vendors v ON v.id = o.vendor_id").
		Where("r.requested_at >= ? AND r.requested_at < ?", period.Start, period.End)

	if params.Status != "" {
		query = query.Where("r.status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"(CAST(r.id AS TEXT) ILIKE ? OR CAST(r.order_id AS TEXT) ILIKE ? OR LOWER(u.full_name) LIKE ? OR LOWER(v.display_name) LIKE ?)",
			search, search, search, search,
		)
	}

	return query
}

func (r *financeRepository) GetRefundRecord(ctx context.Context, refundID uuid.UUID) (*domain.RefundRecord, error) {
	type row struct {
		ID      string `gorm:"column:id"`
		OrderID string `gorm:"column:order_id"`
		Status  string `gorm:"column:status"`
	}

	var result row
	err := r.db.WithContext(ctx).
		Table("refunds").
		Select("id, order_id, status").
		Where("id = ?", refundID.String()).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	id, _ := uuid.Parse(result.ID)
	orderID, _ := uuid.Parse(result.OrderID)
	return &domain.RefundRecord{
		ID:      id,
		OrderID: orderID,
		Status:  result.Status,
	}, nil
}

func (r *financeRepository) GetPayoutRecord(ctx context.Context, payoutID uuid.UUID) (*domain.PayoutRecord, error) {
	type row struct {
		ID     string `gorm:"column:id"`
		Status string `gorm:"column:status"`
	}

	var result row
	err := r.db.WithContext(ctx).
		Table("payout_batches").
		Select("id, status").
		Where("id = ?", payoutID.String()).
		First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	id, _ := uuid.Parse(result.ID)
	return &domain.PayoutRecord{
		ID:     id,
		Status: result.Status,
	}, nil
}

func (r *financeRepository) UpdateRefundStatus(ctx context.Context, refundID uuid.UUID, status string, actorID uuid.UUID) (*domain.StatusActionResponse, error) {
	now := time.Now().UTC()

	updates := map[string]interface{}{
		"status": status,
	}
	if status != domain.RefundStatusRequested {
		updates["processed_by"] = actorID.String()
		updates["processed_at"] = now
	}

	return r.updateRefundTransaction(ctx, refundID, status, updates)
}

func (r *financeRepository) updateRefundTransaction(ctx context.Context, refundID uuid.UUID, status string, updates map[string]interface{}) (*domain.StatusActionResponse, error) {
	var resp *domain.StatusActionResponse
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID      string `gorm:"column:id"`
			OrderID string `gorm:"column:order_id"`
		}

		if err := tx.Table("refunds").
			Select("id, order_id").
			Where("id = ?", refundID.String()).
			First(&refund).Error; err != nil {
			return err
		}

		if err := tx.Table("refunds").
			Where("id = ?", refundID.String()).
			Updates(updates).Error; err != nil {
			return err
		}

		if status == domain.RefundStatusProcessed {
			if err := tx.Table("orders").
				Where("id = ?", refund.OrderID).
				Update("payment_status", domain.PaymentInvoiceStatusRefunded).Error; err != nil {
				return err
			}
		}

		id, _ := uuid.Parse(refund.ID)
		resp = &domain.StatusActionResponse{
			ID:     id,
			Status: status,
		}
		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return resp, nil
}

func (r *financeRepository) UpdatePayoutStatus(ctx context.Context, payoutID uuid.UUID, status string, actorID uuid.UUID) (*domain.StatusActionResponse, error) {
	now := time.Now().UTC()

	updates := map[string]interface{}{
		"status": status,
	}
	if status == domain.PayoutStatusComplete {
		updates["paid_at"] = now
	}

	var resp *domain.StatusActionResponse
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var payout struct {
			ID string `gorm:"column:id"`
		}

		if err := tx.Table("payout_batches").
			Select("id").
			Where("id = ?", payoutID.String()).
			First(&payout).Error; err != nil {
			return err
		}

		if err := tx.Table("payout_batches").
			Where("id = ?", payoutID.String()).
			Updates(updates).Error; err != nil {
			return err
		}

		id, _ := uuid.Parse(payout.ID)
		resp = &domain.StatusActionResponse{
			ID:     id,
			Status: status,
		}
		return nil
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return resp, nil
}
