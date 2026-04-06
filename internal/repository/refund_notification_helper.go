package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type refundNotificationTarget struct {
	OrderNo     string
	CustomerID  uuid.UUID
	VendorOwner *uuid.UUID
}

func (r *financeRepository) loadRefundNotificationTargetTx(tx *gorm.DB, refundID string) (*refundNotificationTarget, error) {
	type row struct {
		OrderNo  string `gorm:"column:order_no"`
		UserID   string `gorm:"column:user_id"`
		VendorID string `gorm:"column:vendor_id"`
	}

	var data row
	if err := tx.Table("refunds r").
		Joins("JOIN orders o ON o.id = r.order_id").
		Select("o.order_no, o.user_id, o.vendor_id").
		Where("r.id = ?", refundID).
		Take(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	customerID, err := uuid.Parse(strings.TrimSpace(data.UserID))
	if err != nil {
		return nil, nil
	}

	var ownerID *uuid.UUID
	vendorID, err := uuid.Parse(strings.TrimSpace(data.VendorID))
	if err == nil {
		ownerID, err = findVendorOwnerUserIDTx(tx, vendorID)
		if err != nil {
			return nil, err
		}
	}

	return &refundNotificationTarget{
		OrderNo:     strings.TrimSpace(data.OrderNo),
		CustomerID:  customerID,
		VendorOwner: ownerID,
	}, nil
}

func (r *financeRepository) notifyRefundStatusChangedTx(tx *gorm.DB, refundID, status string) {
	target, err := r.loadRefundNotificationTargetTx(tx, refundID)
	if err != nil || target == nil {
		return
	}

	orderRef := target.OrderNo
	if orderRef == "" {
		orderRef = refundID
	}

	customerTitle, customerMessage := buildCustomerRefundStatusNotification(status, orderRef)
	_ = createNotificationTx(tx, target.CustomerID, domain.NotificationTypeRefund, customerTitle, customerMessage)

	if target.VendorOwner != nil && *target.VendorOwner != target.CustomerID {
		vendorTitle, vendorMessage := buildVendorRefundStatusNotification(status, orderRef)
		_ = createNotificationTx(tx, *target.VendorOwner, domain.NotificationTypeRefund, vendorTitle, vendorMessage)
	}
}

func (r *financeRepository) notifyRefundProcessingInitiatedTx(tx *gorm.DB, refundID string) {
	target, err := r.loadRefundNotificationTargetTx(tx, refundID)
	if err != nil || target == nil {
		return
	}

	orderRef := target.OrderNo
	if orderRef == "" {
		orderRef = refundID
	}

	_ = createNotificationTx(
		tx,
		target.CustomerID,
		domain.NotificationTypeRefund,
		"Refund Sedang Diproses",
		fmt.Sprintf("Refund untuk pesanan %s sedang diproses oleh tim finance.", orderRef),
	)

	if target.VendorOwner != nil && *target.VendorOwner != target.CustomerID {
		_ = createNotificationTx(
			tx,
			*target.VendorOwner,
			domain.NotificationTypeRefund,
			"Refund Diproses Finance",
			fmt.Sprintf("Refund untuk pesanan %s sedang diproses oleh tim finance.", orderRef),
		)
	}
}

func (r *financeRepository) notifyRefundPayoutFailedTx(tx *gorm.DB, refundID, payoutStatus, failedReason string) {
	target, err := r.loadRefundNotificationTargetTx(tx, refundID)
	if err != nil || target == nil {
		return
	}

	orderRef := target.OrderNo
	if orderRef == "" {
		orderRef = refundID
	}

	normalized := strings.ToUpper(strings.TrimSpace(payoutStatus))
	customerTitle := "Proses Refund Tertunda"
	customerMessage := fmt.Sprintf("Proses refund untuk pesanan %s tertunda. Tim finance akan melakukan tindak lanjut.", orderRef)
	vendorTitle := "Refund Tertunda"
	vendorMessage := fmt.Sprintf("Proses refund untuk pesanan %s tertunda dan perlu tindak lanjut.", orderRef)

	if normalized == "FAILED" {
		customerTitle = "Proses Refund Gagal"
		customerMessage = fmt.Sprintf("Proses refund untuk pesanan %s gagal. Tim finance akan meninjau ulang.", orderRef)
		vendorTitle = "Proses Refund Gagal"
		vendorMessage = fmt.Sprintf("Proses refund untuk pesanan %s gagal dan perlu ditinjau ulang.", orderRef)
	}
	if strings.TrimSpace(failedReason) != "" {
		customerMessage = fmt.Sprintf("%s (%s)", customerMessage, failedReason)
		vendorMessage = fmt.Sprintf("%s (%s)", vendorMessage, failedReason)
	}

	_ = createNotificationTx(tx, target.CustomerID, domain.NotificationTypeRefund, customerTitle, customerMessage)
	if target.VendorOwner != nil && *target.VendorOwner != target.CustomerID {
		_ = createNotificationTx(tx, *target.VendorOwner, domain.NotificationTypeRefund, vendorTitle, vendorMessage)
	}
}

func buildCustomerRefundStatusNotification(status, orderRef string) (title, message string) {
	switch status {
	case domain.RefundStatusApproved:
		return "Refund Disetujui", fmt.Sprintf("Pengajuan refund untuk pesanan %s telah disetujui finance.", orderRef)
	case domain.RefundStatusRejected:
		return "Refund Ditolak", fmt.Sprintf("Pengajuan refund untuk pesanan %s ditolak oleh finance.", orderRef)
	case domain.RefundStatusProcessed:
		return "Refund Berhasil Diproses", fmt.Sprintf("Refund untuk pesanan %s telah berhasil diproses.", orderRef)
	case domain.RefundStatusRequested:
		return "Status Refund Diperbarui", fmt.Sprintf("Status refund untuk pesanan %s telah diperbarui.", orderRef)
	default:
		return "Status Refund Diperbarui", fmt.Sprintf("Status refund untuk pesanan %s berubah menjadi %s.", orderRef, status)
	}
}

func buildVendorRefundStatusNotification(status, orderRef string) (title, message string) {
	switch status {
	case domain.RefundStatusApproved:
		return "Refund Disetujui Finance", fmt.Sprintf("Pengajuan refund pesanan %s telah disetujui finance.", orderRef)
	case domain.RefundStatusRejected:
		return "Refund Ditolak Finance", fmt.Sprintf("Pengajuan refund pesanan %s ditolak oleh finance.", orderRef)
	case domain.RefundStatusProcessed:
		return "Refund Selesai Diproses", fmt.Sprintf("Refund untuk pesanan %s sudah selesai diproses.", orderRef)
	case domain.RefundStatusRequested:
		return "Status Refund Diperbarui", fmt.Sprintf("Status refund untuk pesanan %s telah diperbarui.", orderRef)
	default:
		return "Status Refund Diperbarui", fmt.Sprintf("Status refund untuk pesanan %s berubah menjadi %s.", orderRef, status)
	}
}
