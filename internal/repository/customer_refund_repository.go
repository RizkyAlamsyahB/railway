package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
	"gorm.io/gorm"
)

type refundEvidenceModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	RefundID      string    `gorm:"column:refund_id"`
	ObjectKey     string    `gorm:"column:object_key"`
	MimeType      *string   `gorm:"column:mime_type"`
	FileSizeBytes *int      `gorm:"column:file_size_bytes"`
	MediaType     string    `gorm:"column:media_type"`
	SortOrder     int       `gorm:"column:sort_order"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (refundEvidenceModel) TableName() string { return "refund_evidences" }

type customerRefundRepository struct {
	db          *gorm.DB
	fieldCipher *sensitivedata.FieldCipher
}

// NewCustomerRefundRepository creates a CustomerRefundRepository backed by GORM.
func NewCustomerRefundRepository(db *gorm.DB, fieldCipher *sensitivedata.FieldCipher) domain.CustomerRefundRepository {
	return &customerRefundRepository{
		db:          db,
		fieldCipher: fieldCipher,
	}
}

func (r *customerRefundRepository) HasOpenRefundByOrderID(ctx context.Context, orderID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("refunds").
		Where("order_id = ?", orderID.String()).
		Where("status IN ?", []string{
			domain.RefundStatusRequested,
			domain.RefundStatusApproved,
			domain.RefundStatusProcessed,
		}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *customerRefundRepository) CreateRefundRequest(ctx context.Context, input domain.CreateCustomerRefundInput) (*domain.CreateOrderRefundResponse, error) {
	accountNumberEncrypted := input.DestinationAccountNumber
	if strings.TrimSpace(accountNumberEncrypted) != "" {
		if r.fieldCipher != nil {
			encrypted, err := r.fieldCipher.EncryptString(accountNumberEncrypted)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt refund destination account number: %w", err)
			}
			accountNumberEncrypted = encrypted
		}
	}

	status := input.Status
	if status == "" {
		status = domain.RefundStatusRequested
	}

	refundID := uuid.New()
	updates := map[string]interface{}{
		"id":                              refundID.String(),
		"order_id":                        input.OrderID.String(),
		"payment_invoice_id":              input.PaymentInvoiceID.String(),
		"amount":                          input.Amount,
		"reason":                          input.Reason,
		"reason_detail":                   input.Description,
		"status":                          status,
		"refund_method":                   input.RefundMethod,
		"requested_by":                    input.RequestedBy.String(),
		"requested_at":                    input.RequestedAt,
		"destination_channel_code":        input.DestinationChannelCode,
		"destination_bank_name":           input.DestinationBankName,
		"destination_account_number":      accountNumberEncrypted,
		"destination_account_holder_name": input.DestinationAccountHolderName,
		"destination_account_last4":       input.DestinationAccountNumberLast4,
	}

	if input.ReturnReasonID != 0 {
		updates["return_reason_id"] = input.ReturnReasonID
	}

	if input.UserBankAccountID != nil {
		updates["user_bank_account_id"] = input.UserBankAccountID.String()
	}

	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("refunds").Create(updates).Error; err != nil {
			return err
		}

		if len(input.Evidences) == 0 {
			// Best-effort notifications: do not fail refund request if notification insert fails.
			_ = r.createRefundNotificationsTx(tx, input, status)
			return nil
		}

		models := make([]refundEvidenceModel, 0, len(input.Evidences))
		now := time.Now().UTC()
		for _, evidence := range input.Evidences {
			var mimeType *string
			if strings.TrimSpace(evidence.MimeType) != "" {
				v := evidence.MimeType
				mimeType = &v
			}
			var fileSizeBytes *int
			if evidence.FileSizeBytes > 0 {
				v := evidence.FileSizeBytes
				fileSizeBytes = &v
			}
			models = append(models, refundEvidenceModel{
				ID:            uuid.New().String(),
				RefundID:      refundID.String(),
				ObjectKey:     evidence.ObjectKey,
				MimeType:      mimeType,
				FileSizeBytes: fileSizeBytes,
				MediaType:     evidence.MediaType,
				SortOrder:     evidence.SortOrder,
				CreatedAt:     now,
			})
		}

		if err := tx.Create(&models).Error; err != nil {
			return err
		}

		// Best-effort notifications: do not fail refund request if notification insert fails.
		_ = r.createRefundNotificationsTx(tx, input, status)
		return nil
	}); err != nil {
		return nil, err
	}

	return &domain.CreateOrderRefundResponse{
		RefundID:                      refundID,
		OrderID:                       input.OrderID,
		Status:                        domain.RefundStatusRequested,
		Amount:                        input.Amount,
		ReturnReasonID:                input.ReturnReasonID,
		Reason:                        input.Reason,
		Description:                   input.Description,
		DestinationChannelCode:        input.DestinationChannelCode,
		DestinationBankName:           input.DestinationBankName,
		DestinationAccountHolderName:  input.DestinationAccountHolderName,
		DestinationAccountNumberLast4: input.DestinationAccountNumberLast4,
		RequestedAt:                   input.RequestedAt,
	}, nil
}

func (r *customerRefundRepository) createRefundNotificationsTx(tx *gorm.DB, input domain.CreateCustomerRefundInput, status string) error {
	orderReference := input.OrderID.String()

	var orderRow struct {
		OrderNo  string `gorm:"column:order_no"`
		VendorID string `gorm:"column:vendor_id"`
	}
	if err := tx.Table("orders").
		Select("order_no, vendor_id").
		Where("id = ?", input.OrderID.String()).
		Take(&orderRow).Error; err == nil {
		if strings.TrimSpace(orderRow.OrderNo) != "" {
			orderReference = orderRow.OrderNo
		}
	}

	customerTitle, customerMessage := buildCustomerRefundNotification(status, orderReference)
	if err := createNotificationTx(tx, input.RequestedBy, domain.NotificationTypeRefund, customerTitle, customerMessage); err != nil {
		return err
	}

	if orderRow.VendorID == "" {
		return nil
	}
	vendorID, err := uuid.Parse(orderRow.VendorID)
	if err != nil {
		return nil
	}

	ownerID, err := findVendorOwnerUserIDTx(tx, vendorID)
	if err != nil || ownerID == nil || *ownerID == input.RequestedBy {
		return err
	}

	vendorTitle := "Pengajuan Refund Baru"
	vendorMessage := fmt.Sprintf("Customer mengajukan refund untuk pesanan %s.", orderReference)
	return createNotificationTx(tx, *ownerID, domain.NotificationTypeRefund, vendorTitle, vendorMessage)
}

func buildCustomerRefundNotification(status, orderReference string) (title, message string) {
	switch status {
	case domain.RefundStatusProcessing:
		return "Refund Sedang Diproses", fmt.Sprintf("Refund untuk pesanan %s sedang diproses.", orderReference)
	case domain.RefundStatusAwaitingDestination:
		return "Refund Menunggu Rekening Tujuan", fmt.Sprintf("Refund untuk pesanan %s memerlukan rekening tujuan.", orderReference)
	default:
		return "Pengajuan Refund Diterima", fmt.Sprintf("Pengajuan refund untuk pesanan %s sudah diterima.", orderReference)
	}
}

func (r *customerRefundRepository) SubmitRefundDestination(ctx context.Context, refundID uuid.UUID, channelCode, bankName, accountNumber, accountHolderName, accountLast4 string, userBankAccountID *uuid.UUID) error {
	encrypted := accountNumber
	if strings.TrimSpace(encrypted) != "" && r.fieldCipher != nil {
		enc, err := r.fieldCipher.EncryptString(encrypted)
		if err != nil {
			return fmt.Errorf("failed to encrypt refund destination account number: %w", err)
		}
		encrypted = enc
	}

	updates := map[string]interface{}{
		"status":                          domain.RefundStatusRequested,
		"destination_channel_code":        channelCode,
		"destination_bank_name":           bankName,
		"destination_account_number":      encrypted,
		"destination_account_holder_name": accountHolderName,
		"destination_account_last4":       accountLast4,
	}
	if userBankAccountID != nil {
		updates["user_bank_account_id"] = userBankAccountID.String()
	}

	result := r.db.WithContext(ctx).
		Table("refunds").
		Where("id = ? AND status = ?", refundID.String(), domain.RefundStatusAwaitingDestination).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("refund not found or not in awaiting_destination status")
	}
	return nil
}

func (r *customerRefundRepository) FindByID(ctx context.Context, refundID uuid.UUID) (*domain.RefundDetail, error) {
	var result struct {
		ID                     string  `gorm:"column:id"`
		OrderID                string  `gorm:"column:order_id"`
		PaymentInvoiceID       string  `gorm:"column:payment_invoice_id"`
		Amount                 float64 `gorm:"column:amount"`
		Status                 string  `gorm:"column:status"`
		RefundMethod           *string `gorm:"column:refund_method"`
		DestinationChannelCode *string `gorm:"column:destination_channel_code"`
		DestinationAccountNum  *string `gorm:"column:destination_account_number"`
	}
	if err := r.db.WithContext(ctx).
		Table("refunds").
		Where("id = ?", refundID.String()).
		First(&result).Error; err != nil {
		return nil, err
	}

	id, _ := uuid.Parse(result.ID)
	orderID, _ := uuid.Parse(result.OrderID)
	paymentInvoiceID, _ := uuid.Parse(result.PaymentInvoiceID)
	return &domain.RefundDetail{
		ID:                       id,
		OrderID:                  orderID,
		PaymentInvoiceID:         paymentInvoiceID,
		Amount:                   result.Amount,
		Status:                   result.Status,
		RefundMethod:             result.RefundMethod,
		DestinationChannelCode:   result.DestinationChannelCode,
		DestinationAccountNumber: result.DestinationAccountNum,
	}, nil
}

func (r *customerRefundRepository) FindLatestByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.RefundDetail, error) {
	var result struct {
		ID                      string  `gorm:"column:id"`
		OrderID                 string  `gorm:"column:order_id"`
		PaymentInvoiceID        string  `gorm:"column:payment_invoice_id"`
		Amount                  float64 `gorm:"column:amount"`
		Status                  string  `gorm:"column:status"`
		RefundMethod            *string `gorm:"column:refund_method"`
		DestinationChannelCode  *string `gorm:"column:destination_channel_code"`
		DestinationAccountNum   *string `gorm:"column:destination_account_number"`
		DestinationBankName     *string `gorm:"column:destination_bank_name"`
		DestinationAccountLast4 *string `gorm:"column:destination_account_last4"`
	}
	if err := r.db.WithContext(ctx).
		Table("refunds").
		Where("order_id = ?", orderID.String()).
		Order("requested_at DESC").
		First(&result).Error; err != nil {
		return nil, err
	}

	id, _ := uuid.Parse(result.ID)
	oid, _ := uuid.Parse(result.OrderID)
	pid, _ := uuid.Parse(result.PaymentInvoiceID)
	return &domain.RefundDetail{
		ID:                            id,
		OrderID:                       oid,
		PaymentInvoiceID:              pid,
		Amount:                        result.Amount,
		Status:                        result.Status,
		RefundMethod:                  result.RefundMethod,
		DestinationChannelCode:        result.DestinationChannelCode,
		DestinationAccountNumber:      result.DestinationAccountNum,
		DestinationBankName:           result.DestinationBankName,
		DestinationAccountNumberLast4: result.DestinationAccountLast4,
	}, nil
}

func (r *customerRefundRepository) FindAwaitingDestinationByOrderAndUser(ctx context.Context, orderID, userID uuid.UUID) (*domain.RefundDetail, error) {
	var result struct {
		ID                     string  `gorm:"column:id"`
		OrderID                string  `gorm:"column:order_id"`
		PaymentInvoiceID       string  `gorm:"column:payment_invoice_id"`
		Amount                 float64 `gorm:"column:amount"`
		Status                 string  `gorm:"column:status"`
		RefundMethod           *string `gorm:"column:refund_method"`
		DestinationChannelCode *string `gorm:"column:destination_channel_code"`
		DestinationAccountNum  *string `gorm:"column:destination_account_number"`
	}
	if err := r.db.WithContext(ctx).
		Table("refunds").
		Where("order_id = ? AND requested_by = ? AND status = ?", orderID.String(), userID.String(), domain.RefundStatusAwaitingDestination).
		First(&result).Error; err != nil {
		return nil, err
	}

	id, _ := uuid.Parse(result.ID)
	oid, _ := uuid.Parse(result.OrderID)
	pid, _ := uuid.Parse(result.PaymentInvoiceID)
	return &domain.RefundDetail{
		ID:                       id,
		OrderID:                  oid,
		PaymentInvoiceID:         pid,
		Amount:                   result.Amount,
		Status:                   result.Status,
		RefundMethod:             result.RefundMethod,
		DestinationChannelCode:   result.DestinationChannelCode,
		DestinationAccountNumber: result.DestinationAccountNum,
	}, nil
}
