package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new NotificationRepository backed by GORM.
func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) List(ctx context.Context, userID uuid.UUID, params domain.NotificationListParams) ([]domain.NotificationItem, int64, error) {
	base := r.db.WithContext(ctx).
		Table("notifications").
		Where("user_id = ?", userID.String())

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID        string       `gorm:"column:id"`
		Type      string       `gorm:"column:type"`
		Title     string       `gorm:"column:title"`
		Message   string       `gorm:"column:message"`
		CreatedAt time.Time    `gorm:"column:created_at"`
		ReadAt    sql.NullTime `gorm:"column:read_at"`
	}

	query := base.
		Select("id, type, title, message, created_at, read_at").
		Order("created_at DESC")

	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		query = query.Offset(offset).Limit(params.Limit)
	}

	var rows []row
	if err := query.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	items := make([]domain.NotificationItem, len(rows))
	for i, r := range rows {
		id, _ := uuid.Parse(r.ID)
		var readAt *time.Time
		if r.ReadAt.Valid {
			readAt = &r.ReadAt.Time
		}
		items[i] = domain.NotificationItem{
			ID:        id,
			Type:      r.Type,
			Title:     r.Title,
			Message:   r.Message,
			CreatedAt: r.CreatedAt,
			ReadAt:    readAt,
		}
	}

	return items, total, nil
}

func (r *notificationRepository) UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).
		Table("notifications").
		Where("user_id = ? AND read_at IS NULL", userID.String()).
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *notificationRepository) MarkRead(ctx context.Context, userID, notificationID uuid.UUID) (*domain.NotificationItem, error) {
	type row struct {
		ID        string       `gorm:"column:id"`
		Type      string       `gorm:"column:type"`
		Title     string       `gorm:"column:title"`
		Message   string       `gorm:"column:message"`
		CreatedAt time.Time    `gorm:"column:created_at"`
		ReadAt    sql.NullTime `gorm:"column:read_at"`
	}

	var record row
	if err := r.db.WithContext(ctx).
		Table("notifications").
		Select("id, type, title, message, created_at, read_at").
		Where("id = ? AND user_id = ?", notificationID.String(), userID.String()).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if !record.ReadAt.Valid {
		now := time.Now().UTC()
		if err := r.db.WithContext(ctx).
			Table("notifications").
			Where("id = ? AND user_id = ?", notificationID.String(), userID.String()).
			Update("read_at", now).Error; err != nil {
			return nil, err
		}
		record.ReadAt = sql.NullTime{Time: now, Valid: true}
	}

	id, _ := uuid.Parse(record.ID)
	var readAt *time.Time
	if record.ReadAt.Valid {
		readAt = &record.ReadAt.Time
	}

	return &domain.NotificationItem{
		ID:        id,
		Type:      record.Type,
		Title:     record.Title,
		Message:   record.Message,
		CreatedAt: record.CreatedAt,
		ReadAt:    readAt,
	}, nil
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) (int64, error) {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).
		Table("notifications").
		Where("user_id = ? AND read_at IS NULL", userID.String()).
		Update("read_at", now)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
