package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type financeRepository struct {
	db          *gorm.DB
	fieldCipher *sensitivedata.FieldCipher
}

// NewFinanceRepository creates a new FinanceRepository backed by GORM.
func NewFinanceRepository(db *gorm.DB, fieldCipher *sensitivedata.FieldCipher) domain.FinanceRepository {
	return &financeRepository{
		db:          db,
		fieldCipher: fieldCipher,
	}
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
		Vendor   *string        `gorm:"column:vendor"`
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
		Vendor         *string `gorm:"column:vendor"`
		VendorType     *string `gorm:"column:vendor_type"`
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
		Vendor   *string        `gorm:"column:vendor"`
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

func (r *financeRepository) GetRefundDisbursementContext(ctx context.Context, refundID uuid.UUID) (*domain.RefundDisbursementContext, error) {
	type row struct {
		RefundID                     string         `gorm:"column:refund_id"`
		OrderID                      string         `gorm:"column:order_id"`
		Amount                       float64        `gorm:"column:amount"`
		Status                       string         `gorm:"column:status"`
		Currency                     sql.NullString `gorm:"column:currency"`
		PaymentMethod                sql.NullString `gorm:"column:payment_method"`
		PaymentChannel               sql.NullString `gorm:"column:payment_channel"`
		VendorXenditAccountID        sql.NullString `gorm:"column:vendor_xendit_account_id"`
		PayoutReferenceID            sql.NullString `gorm:"column:payout_reference_id"`
		PayoutStatus                 sql.NullString `gorm:"column:payout_status"`
		DestinationChannelCode       sql.NullString `gorm:"column:destination_channel_code"`
		DestinationBankName          sql.NullString `gorm:"column:destination_bank_name"`
		DestinationAccountNumber     sql.NullString `gorm:"column:destination_account_number"`
		DestinationAccountHolderName sql.NullString `gorm:"column:destination_account_holder_name"`
		DestinationAccountLast4      sql.NullString `gorm:"column:destination_account_last4"`
	}

	var result row
	err := r.db.WithContext(ctx).
		Table("refunds r").
		Joins("JOIN payment_invoices pi ON pi.id = r.payment_invoice_id").
		Joins("JOIN orders o ON o.id = r.order_id").
		Joins("JOIN vendors v ON v.id = o.vendor_id").
		Select(`
			r.id AS refund_id,
			r.order_id AS order_id,
			r.amount AS amount,
			r.status AS status,
			pi.currency AS currency,
			pi.payment_method AS payment_method,
			pi.payment_channel AS payment_channel,
			v.xendit_account_id AS vendor_xendit_account_id,
			r.payout_reference_id AS payout_reference_id,
			r.payout_status AS payout_status,
			r.destination_channel_code AS destination_channel_code,
			r.destination_bank_name AS destination_bank_name,
			r.destination_account_number AS destination_account_number,
			r.destination_account_holder_name AS destination_account_holder_name,
			r.destination_account_last4 AS destination_account_last4
		`).
		Where("r.id = ?", refundID.String()).
		Take(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	refundUUID, _ := uuid.Parse(result.RefundID)
	orderUUID, _ := uuid.Parse(result.OrderID)
	currency := "IDR"
	if result.Currency.Valid && strings.TrimSpace(result.Currency.String) != "" {
		currency = result.Currency.String
	}

	ctxResult := &domain.RefundDisbursementContext{
		RefundID: refundUUID,
		OrderID:  orderUUID,
		Amount:   result.Amount,
		Currency: currency,
		Status:   result.Status,
	}

	if result.PaymentMethod.Valid {
		v := result.PaymentMethod.String
		ctxResult.PaymentMethod = &v
	}
	if result.PaymentChannel.Valid {
		v := result.PaymentChannel.String
		ctxResult.PaymentChannel = &v
	}
	if result.VendorXenditAccountID.Valid {
		v := result.VendorXenditAccountID.String
		ctxResult.VendorXenditAccountID = &v
	}
	if result.PayoutReferenceID.Valid {
		v := result.PayoutReferenceID.String
		ctxResult.PayoutReferenceID = &v
	}
	if result.PayoutStatus.Valid {
		v := result.PayoutStatus.String
		ctxResult.PayoutStatus = &v
	}
	if result.DestinationChannelCode.Valid {
		v := strings.TrimSpace(result.DestinationChannelCode.String)
		if v != "" {
			ctxResult.DestinationChannelCode = &v
		}
	}
	if result.DestinationBankName.Valid {
		v := strings.TrimSpace(result.DestinationBankName.String)
		if v != "" {
			ctxResult.DestinationBankName = &v
		}
	}
	if result.DestinationAccountHolderName.Valid {
		v := strings.TrimSpace(result.DestinationAccountHolderName.String)
		if v != "" {
			ctxResult.DestinationAccountHolderName = &v
		}
	}
	if result.DestinationAccountLast4.Valid {
		v := strings.TrimSpace(result.DestinationAccountLast4.String)
		if v != "" {
			ctxResult.DestinationAccountLast4 = &v
		}
	}
	if result.DestinationAccountNumber.Valid && strings.TrimSpace(result.DestinationAccountNumber.String) != "" {
		decrypted := result.DestinationAccountNumber.String
		if r.fieldCipher != nil {
			var err error
			decrypted, err = r.fieldCipher.DecryptString(result.DestinationAccountNumber.String)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt refund destination account number: %w", err)
			}
		}
		decrypted = strings.TrimSpace(decrypted)
		if decrypted != "" {
			ctxResult.DestinationAccountNumber = &decrypted
		}
	}

	return ctxResult, nil
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

func (r *financeRepository) SetRefundPayoutInitiated(ctx context.Context, refundID uuid.UUID, actorID uuid.UUID, strategy string, referenceID string, channelCode string, bankName string, accountHolderName string, accountLast4 string, payoutID *string, payoutStatus string) error {
	updates := map[string]interface{}{
		"payout_strategy":                        strategy,
		"payout_reference_id":                    referenceID,
		"payout_channel_code":                    channelCode,
		"payout_destination_bank_name":           bankName,
		"payout_destination_account_holder_name": accountHolderName,
		"payout_destination_account_last4":       accountLast4,
		"payout_status":                          payoutStatus,
		"payout_failed_reason":                   nil,
		"payout_requested_by":                    actorID.String(),
		"payout_requested_at":                    time.Now().UTC(),
	}
	if payoutID != nil && strings.TrimSpace(*payoutID) != "" {
		updates["payout_id"] = strings.TrimSpace(*payoutID)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("refunds").
			Where("id = ?", refundID.String()).
			Updates(updates).Error; err != nil {
			return err
		}

		r.notifyRefundProcessingInitiatedTx(tx, refundID.String())
		return nil
	})
}

func (r *financeRepository) SetRefundPayoutFailed(ctx context.Context, refundID uuid.UUID, payoutStatus string, failedReason string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("refunds").
			Where("id = ?", refundID.String()).
			Updates(map[string]interface{}{
				"payout_status":        payoutStatus,
				"payout_failed_reason": failedReason,
			}).Error; err != nil {
			return err
		}

		r.notifyRefundPayoutFailedTx(tx, refundID.String(), payoutStatus, failedReason)
		return nil
	})
}

func (r *financeRepository) ApplyRefundPayoutWebhookUpdate(ctx context.Context, referenceID string, payoutID string, payoutStatus string, failedReason *string, completedAt *time.Time) (bool, error) {
	applied := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID      string         `gorm:"column:id"`
			OrderID string         `gorm:"column:order_id"`
			Status  string         `gorm:"column:status"`
			Reason  sql.NullString `gorm:"column:payout_failed_reason"`
		}

		if err := tx.Table("refunds").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, order_id, status, payout_failed_reason").
			Where("payout_reference_id = ?", referenceID).
			First(&refund).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}

		var adjustment refundBalanceAdjustmentResult
		if completedAt != nil {
			var err error
			adjustment, err = r.applyRefundSettlementBalanceAdjustment(tx, refund.OrderID)
			if err != nil {
				return err
			}
		}

		updates := map[string]interface{}{
			"payout_status": payoutStatus,
		}
		if strings.TrimSpace(payoutID) != "" {
			updates["payout_id"] = strings.TrimSpace(payoutID)
		}
		if failedReason != nil {
			updates["payout_failed_reason"] = *failedReason
		}
		if completedAt != nil {
			updates["status"] = domain.RefundStatusProcessed
			updates["processed_at"] = *completedAt
			updates["payout_completed_at"] = *completedAt
			if adjustment.Shortfall > 0 {
				updates["payout_failed_reason"] = buildRefundBalanceShortfallReason(adjustment)
			} else {
				updates["payout_failed_reason"] = nil
			}
		}

		if err := tx.Table("refunds").
			Where("id = ?", refund.ID).
			Updates(updates).Error; err != nil {
			return err
		}

		if completedAt != nil && refund.Status != domain.RefundStatusProcessed {
			r.notifyRefundStatusChangedTx(tx, refund.ID, domain.RefundStatusProcessed)
		}

		applied = true
		return nil
	})
	if err != nil {
		return false, err
	}

	return applied, nil
}

func (r *financeRepository) IsOrderSettlementCompleted(ctx context.Context, orderID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("payout_items pi").
		Joins("JOIN payout_batches pb ON pb.id = pi.payout_batch_id").
		Where("pi.order_id = ?", orderID.String()).
		Where("pb.status = ?", domain.PayoutStatusComplete).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *financeRepository) GetRefundVendorBalanceSnapshot(ctx context.Context, orderID uuid.UUID) (*domain.RefundVendorBalanceSnapshot, error) {
	type row struct {
		OrderStatus      string  `gorm:"column:order_status"`
		Subtotal         float64 `gorm:"column:subtotal"`
		AvailableBalance float64 `gorm:"column:available_balance"`
		EscrowBalance    float64 `gorm:"column:escrow_balance"`
	}

	var result row
	if err := r.db.WithContext(ctx).
		Table("orders o").
		Joins("LEFT JOIN vendor_balances vb ON vb.vendor_id = o.vendor_id").
		Select(`
			o.order_status AS order_status,
			o.subtotal AS subtotal,
			COALESCE(vb.available_balance, 0) AS available_balance,
			COALESCE(vb.escrow_balance, 0) AS escrow_balance
		`).
		Where("o.id = ?", orderID.String()).
		Take(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	snapshot := &domain.RefundVendorBalanceSnapshot{
		Sufficient:      true,
		BalanceSource:   "none",
		RequiredAmount:  result.Subtotal,
		AvailableAmount: 0,
	}

	switch result.OrderStatus {
	case domain.OrderStatusCompleted:
		snapshot.BalanceSource = "available_balance"
		snapshot.AvailableAmount = result.AvailableBalance
		snapshot.Sufficient = result.AvailableBalance >= result.Subtotal
	case domain.OrderStatusReceived:
		snapshot.BalanceSource = "escrow_balance"
		snapshot.AvailableAmount = result.EscrowBalance
		snapshot.Sufficient = result.EscrowBalance >= result.Subtotal
	}

	return snapshot, nil
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
			Status  string `gorm:"column:status"`
		}

		if err := tx.Table("refunds").
			Select("id, order_id, status").
			Where("id = ?", refundID.String()).
			First(&refund).Error; err != nil {
			return err
		}

		if status == domain.RefundStatusProcessed {
			adjustment, err := r.applyRefundSettlementBalanceAdjustment(tx, refund.OrderID)
			if err != nil {
				return err
			}
			if adjustment.Shortfall > 0 {
				updates["payout_failed_reason"] = buildRefundBalanceShortfallReason(adjustment)
			}
		}

		if err := tx.Table("refunds").
			Where("id = ?", refundID.String()).
			Updates(updates).Error; err != nil {
			return err
		}

		if refund.Status != status {
			r.notifyRefundStatusChangedTx(tx, refund.ID, status)
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

func (r *financeRepository) SetRefundGatewayInitiated(ctx context.Context, refundID uuid.UUID, xenditRefundID string, refundMethod string, payoutReferenceID string) error {
	return r.db.WithContext(ctx).Table("refunds").
		Where("id = ?", refundID.String()).
		Updates(map[string]interface{}{
			"xendit_refund_id":    xenditRefundID,
			"refund_method":       refundMethod,
			"payout_reference_id": payoutReferenceID,
			"status":              domain.RefundStatusProcessing,
		}).Error
}

func (r *financeRepository) ApplyRefundGatewayWebhookUpdate(ctx context.Context, referenceID string, xenditRefundID string, status string, failureCode *string) (bool, error) {
	applied := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var refund struct {
			ID      string `gorm:"column:id"`
			OrderID string `gorm:"column:order_id"`
			Status  string `gorm:"column:status"`
		}
		if err := tx.Table("refunds").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, order_id, status").
			Where("payout_reference_id = ?", referenceID).
			First(&refund).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}

		updates := map[string]interface{}{
			"xendit_refund_id": xenditRefundID,
		}

		nextStatus := ""
		switch status {
		case "SUCCEEDED":
			now := time.Now().UTC()
			adjustment, err := r.applyRefundSettlementBalanceAdjustment(tx, refund.OrderID)
			if err != nil {
				return err
			}
			nextStatus = domain.RefundStatusProcessed
			updates["status"] = nextStatus
			updates["processed_at"] = now
			if adjustment.Shortfall > 0 {
				updates["payout_failed_reason"] = buildRefundBalanceShortfallReason(adjustment)
			}
		case "FAILED":
			nextStatus = domain.RefundStatusRejected
			updates["status"] = nextStatus
			if failureCode != nil {
				updates["payout_failed_reason"] = "Gateway refund failed: " + *failureCode
			}
		default:
			return nil
		}

		if err := tx.Table("refunds").
			Where("id = ?", refund.ID).
			Updates(updates).Error; err != nil {
			return err
		}

		if nextStatus != "" && nextStatus != refund.Status {
			r.notifyRefundStatusChangedTx(tx, refund.ID, nextStatus)
		}

		applied = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return applied, nil
}

type refundBalanceAdjustmentResult struct {
	Expected  float64
	Deducted  float64
	Shortfall float64
	Source    string
}

func (r *financeRepository) applyRefundSettlementBalanceAdjustment(tx *gorm.DB, orderID string) (refundBalanceAdjustmentResult, error) {
	result := refundBalanceAdjustmentResult{}

	var order struct {
		ID            string  `gorm:"column:id"`
		VendorID      string  `gorm:"column:vendor_id"`
		OrderStatus   string  `gorm:"column:order_status"`
		PaymentStatus string  `gorm:"column:payment_status"`
		Subtotal      float64 `gorm:"column:subtotal"`
	}
	if err := tx.Table("orders").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id, vendor_id, order_status, payment_status, subtotal").
		Where("id = ?", orderID).
		First(&order).Error; err != nil {
		return result, err
	}

	if order.PaymentStatus == domain.PaymentStatusRefunded || order.PaymentStatus == domain.PaymentInvoiceStatusRefunded {
		return result, nil
	}

	var balance struct {
		AvailableBalance float64 `gorm:"column:available_balance"`
		EscrowBalance    float64 `gorm:"column:escrow_balance"`
	}
	if err := tx.Table("vendor_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("available_balance, escrow_balance").
		Where("vendor_id = ?", order.VendorID).
		First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, fmt.Errorf("vendor balance not found for vendor %s", order.VendorID)
		}
		return result, err
	}

	now := time.Now().UTC()
	result.Expected = order.Subtotal

	switch order.OrderStatus {
	case domain.OrderStatusCompleted:
		result.Source = "available_balance"
		result.Deducted = math.Min(balance.AvailableBalance, order.Subtotal)
		result.Shortfall = math.Max(order.Subtotal-result.Deducted, 0)
		if result.Deducted > 0 {
			res := tx.Table("vendor_balances").
				Where("vendor_id = ?", order.VendorID).
				Updates(map[string]interface{}{
					"available_balance": gorm.Expr("available_balance - ?", result.Deducted),
					"total_earned":      gorm.Expr("GREATEST(total_earned - ?, 0)", result.Deducted),
					"updated_at":        now,
				})
			if res.Error != nil {
				return result, fmt.Errorf("failed to reverse vendor available balance for refund: %w", res.Error)
			}
		}
	case domain.OrderStatusReceived:
		result.Source = "escrow_balance"
		result.Deducted = math.Min(balance.EscrowBalance, order.Subtotal)
		result.Shortfall = math.Max(order.Subtotal-result.Deducted, 0)
		if result.Deducted > 0 {
			res := tx.Table("vendor_balances").
				Where("vendor_id = ?", order.VendorID).
				Updates(map[string]interface{}{
					"escrow_balance": gorm.Expr("escrow_balance - ?", result.Deducted),
					"total_earned":   gorm.Expr("GREATEST(total_earned - ?, 0)", result.Deducted),
					"updated_at":     now,
				})
			if res.Error != nil {
				return result, fmt.Errorf("failed to reverse vendor escrow balance for refund: %w", res.Error)
			}
		}
	}

	if err := tx.Table("orders").
		Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"payment_status": domain.PaymentInvoiceStatusRefunded,
			"updated_at":     now,
		}).Error; err != nil {
		return result, err
	}

	return result, nil
}

func buildRefundBalanceShortfallReason(adjustment refundBalanceAdjustmentResult) string {
	source := strings.TrimSpace(adjustment.Source)
	if source == "" {
		source = "vendor_balance"
	}
	return fmt.Sprintf(
		"gateway/webhook succeeded so refund cannot be put on hold: deducted %.2f from %s, unrecovered %.2f",
		adjustment.Deducted,
		source,
		adjustment.Shortfall,
	)
}

func (r *financeRepository) UpdateOrderPaymentStatus(ctx context.Context, orderID uuid.UUID, paymentStatus string) error {
	return r.db.WithContext(ctx).Table("orders").
		Where("id = ?", orderID.String()).
		Update("payment_status", paymentStatus).Error
}
