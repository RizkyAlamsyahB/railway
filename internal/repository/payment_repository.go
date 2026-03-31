package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type paymentInvoiceModel struct {
	ID                string     `gorm:"column:id;primaryKey"`
	OrderID           string     `gorm:"column:order_id"`
	Gateway           string     `gorm:"column:gateway"`
	XenditInvoiceID   *string    `gorm:"column:xendit_invoice_id"`
	ExternalInvoiceID string     `gorm:"column:external_invoice_id"`
	InvoiceURL        *string    `gorm:"column:invoice_url"`
	PaymentMethod     *string    `gorm:"column:payment_method"`
	PaymentChannel    *string    `gorm:"column:payment_channel"`
	Amount            float64    `gorm:"column:amount"`
	Currency          string     `gorm:"column:currency"`
	Status            string     `gorm:"column:status"`
	ExpiresAt         *time.Time `gorm:"column:expires_at"`
	PaidAt            *time.Time `gorm:"column:paid_at"`
	RawPayload        *string    `gorm:"column:raw_payload;type:jsonb"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
}

func (paymentInvoiceModel) TableName() string { return "payment_invoices" }

type paymentEventModel struct {
	ID               string    `gorm:"column:id;primaryKey"`
	PaymentInvoiceID string    `gorm:"column:payment_invoice_id"`
	EventType        string    `gorm:"column:event_type"`
	ExternalEventID  string    `gorm:"column:external_event_id"`
	Payload          string    `gorm:"column:payload;type:jsonb"`
	ReceivedAt       time.Time `gorm:"column:received_at"`
}

func (paymentEventModel) TableName() string { return "payment_events" }

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PaymentRepository backed by GORM.
func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreateInvoice(ctx context.Context, invoice *domain.PaymentInvoice) error {
	model := toPaymentInvoiceModel(invoice)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *paymentRepository) FindInvoiceByExternalID(ctx context.Context, externalID string) (*domain.PaymentInvoice, error) {
	var model paymentInvoiceModel
	if err := r.db.WithContext(ctx).Where("external_invoice_id = ?", externalID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainPaymentInvoice(&model), nil
}

func (r *paymentRepository) FindInvoiceByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.PaymentInvoice, error) {
	var model paymentInvoiceModel
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID.String()).Order("created_at DESC").First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainPaymentInvoice(&model), nil
}

func (r *paymentRepository) UpdateInvoiceStatus(ctx context.Context, invoiceID uuid.UUID, status string, paidAt *time.Time, paymentMethod, paymentChannel *string, rawPayload map[string]interface{}) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	if paymentMethod != nil {
		updates["payment_method"] = *paymentMethod
	}
	if paymentChannel != nil {
		updates["payment_channel"] = *paymentChannel
	}
	if rawPayload != nil {
		data, _ := json.Marshal(rawPayload)
		s := string(data)
		updates["raw_payload"] = s
	}

	return r.db.WithContext(ctx).Model(&paymentInvoiceModel{}).
		Where("id = ?", invoiceID.String()).
		Updates(updates).Error
}

func (r *paymentRepository) UpdateInvoiceStatusIfCurrent(
	ctx context.Context,
	invoiceID uuid.UUID,
	expectedCurrentStatus, nextStatus string,
	paidAt *time.Time,
	paymentMethod, paymentChannel *string,
	rawPayload map[string]interface{},
) (bool, error) {
	updates := map[string]interface{}{
		"status":     nextStatus,
		"updated_at": time.Now(),
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	if paymentMethod != nil {
		updates["payment_method"] = *paymentMethod
	}
	if paymentChannel != nil {
		updates["payment_channel"] = *paymentChannel
	}
	if rawPayload != nil {
		data, _ := json.Marshal(rawPayload)
		s := string(data)
		updates["raw_payload"] = s
	}

	result := r.db.WithContext(ctx).Model(&paymentInvoiceModel{}).
		Where("id = ? AND status = ?", invoiceID.String(), expectedCurrentStatus).
		Updates(updates)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *paymentRepository) CreateEvent(ctx context.Context, event *domain.PaymentEvent) error {
	model := toPaymentEventModel(event)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *paymentRepository) EventExistsByExternalID(ctx context.Context, externalEventID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&paymentEventModel{}).
		Where("external_event_id = ?", externalEventID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *paymentRepository) ListForAdmin(ctx context.Context, params domain.AdminPaymentListParams) ([]domain.AdminPaymentListItem, int64, error) {
	base := r.adminPaymentBaseQuery(ctx, params)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		InvoiceID     string         `gorm:"column:invoice_id"`
		OrderID       string         `gorm:"column:order_id"`
		OrderNo       string         `gorm:"column:order_no"`
		CustomerName  string         `gorm:"column:customer_name"`
		VendorName    string         `gorm:"column:vendor_name"`
		PaymentMethod sql.NullString `gorm:"column:payment_method"`
		Amount        float64        `gorm:"column:amount"`
		Status        string         `gorm:"column:status"`
		CreatedAt     time.Time      `gorm:"column:created_at"`
	}

	query := base.Select(`
		pi.id AS invoice_id,
		pi.order_id AS order_id,
		o.order_no AS order_no,
		u.full_name AS customer_name,
		v.display_name AS vendor_name,
		pi.payment_method AS payment_method,
		pi.amount AS amount,
		pi.status AS status,
		pi.created_at AS created_at
	`)

	orderClause := r.buildAdminPaymentOrderClause(params.SortBy, params.SortOrder)
	query = query.Order(orderClause)

	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	var rows []row
	if err := query.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.AdminPaymentListItem, len(rows))
	for i, r := range rows {
		invoiceID, _ := uuid.Parse(r.InvoiceID)
		orderID, _ := uuid.Parse(r.OrderID)
		var method *string
		if r.PaymentMethod.Valid {
			method = &r.PaymentMethod.String
		}
		items[i] = domain.AdminPaymentListItem{
			InvoiceID:     invoiceID,
			OrderID:       orderID,
			OrderNo:       r.OrderNo,
			CustomerName:  r.CustomerName,
			VendorName:    r.VendorName,
			PaymentMethod: method,
			Amount:        r.Amount,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt,
		}
	}

	return items, total, nil
}

func (r *paymentRepository) adminPaymentBaseQuery(ctx context.Context, params domain.AdminPaymentListParams) *gorm.DB {
	query := r.db.WithContext(ctx).
		Table("payment_invoices pi").
		Joins("JOIN orders o ON o.id = pi.order_id").
		Joins("JOIN users u ON u.id = o.user_id").
		Joins("JOIN vendors v ON v.id = o.vendor_id")

	if params.Status != "" {
		query = query.Where("pi.status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"(LOWER(u.full_name) LIKE ? OR LOWER(o.order_no) LIKE ?)",
			search, search,
		)
	}

	return query
}

func (r *paymentRepository) buildAdminPaymentOrderClause(sortBy, sortOrder string) string {
	column := "pi.created_at"
	direction := "DESC"

	switch sortBy {
	case "amount":
		column = "pi.amount"
	case "customer_name":
		column = "u.full_name"
	case "status":
		column = "pi.status"
	}

	if strings.ToLower(sortOrder) == "asc" {
		direction = "ASC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

// Mapper helpers.

func toPaymentInvoiceModel(i *domain.PaymentInvoice) paymentInvoiceModel {
	var rawPayload *string
	if i.RawPayload != nil {
		data, _ := json.Marshal(i.RawPayload)
		s := string(data)
		rawPayload = &s
	}

	return paymentInvoiceModel{
		ID:                i.ID.String(),
		OrderID:           i.OrderID.String(),
		Gateway:           i.Gateway,
		XenditInvoiceID:   i.XenditInvoiceID,
		ExternalInvoiceID: i.ExternalInvoiceID,
		InvoiceURL:        i.InvoiceURL,
		PaymentMethod:     i.PaymentMethod,
		PaymentChannel:    i.PaymentChannel,
		Amount:            i.Amount,
		Currency:          i.Currency,
		Status:            i.Status,
		ExpiresAt:         i.ExpiresAt,
		PaidAt:            i.PaidAt,
		RawPayload:        rawPayload,
		CreatedAt:         i.CreatedAt,
		UpdatedAt:         i.UpdatedAt,
	}
}

func toDomainPaymentInvoice(m *paymentInvoiceModel) *domain.PaymentInvoice {
	id, _ := uuid.Parse(m.ID)
	orderID, _ := uuid.Parse(m.OrderID)

	var rawPayload map[string]interface{}
	if m.RawPayload != nil {
		_ = json.Unmarshal([]byte(*m.RawPayload), &rawPayload)
	}

	return &domain.PaymentInvoice{
		ID:                id,
		OrderID:           orderID,
		Gateway:           m.Gateway,
		XenditInvoiceID:   m.XenditInvoiceID,
		ExternalInvoiceID: m.ExternalInvoiceID,
		InvoiceURL:        m.InvoiceURL,
		PaymentMethod:     m.PaymentMethod,
		PaymentChannel:    m.PaymentChannel,
		Amount:            m.Amount,
		Currency:          m.Currency,
		Status:            m.Status,
		ExpiresAt:         m.ExpiresAt,
		PaidAt:            m.PaidAt,
		RawPayload:        rawPayload,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func toPaymentEventModel(e *domain.PaymentEvent) paymentEventModel {
	payload, _ := json.Marshal(e.Payload)
	return paymentEventModel{
		ID:               e.ID.String(),
		PaymentInvoiceID: e.PaymentInvoiceID.String(),
		EventType:        e.EventType,
		ExternalEventID:  e.ExternalEventID,
		Payload:          string(payload),
		ReceivedAt:       e.ReceivedAt,
	}
}
