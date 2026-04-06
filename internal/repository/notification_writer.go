package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationWriteModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	Type      string    `gorm:"column:type"`
	Title     string    `gorm:"column:title"`
	Message   string    `gorm:"column:message"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (notificationWriteModel) TableName() string { return "notifications" }

func createNotificationTx(tx *gorm.DB, userID uuid.UUID, notifType, title, message string) error {
	if tx == nil || userID == uuid.Nil {
		return nil
	}

	notifType = strings.TrimSpace(notifType)
	title = strings.TrimSpace(title)
	message = strings.TrimSpace(message)
	if notifType == "" || title == "" || message == "" {
		return nil
	}

	model := notificationWriteModel{
		ID:        uuid.NewString(),
		UserID:    userID.String(),
		Type:      notifType,
		Title:     truncateRunes(title, 255),
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}
	return tx.Create(&model).Error
}

func findVendorOwnerUserIDTx(tx *gorm.DB, vendorID uuid.UUID) (*uuid.UUID, error) {
	if tx == nil || vendorID == uuid.Nil {
		return nil, nil
	}

	var row struct {
		OwnerUserID string `gorm:"column:owner_user_id"`
	}
	if err := tx.Table("vendors").
		Select("owner_user_id").
		Where("id = ?", vendorID.String()).
		Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	id, err := uuid.Parse(row.OwnerUserID)
	if err != nil {
		return nil, nil
	}
	return &id, nil
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}
