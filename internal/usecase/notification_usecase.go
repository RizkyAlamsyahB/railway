package usecase

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type notificationUseCase struct {
	repo domain.NotificationRepository
}

// NewNotificationUseCase creates a new NotificationUseCase.
func NewNotificationUseCase(repo domain.NotificationRepository) domain.NotificationUseCase {
	return &notificationUseCase{repo: repo}
}

func (uc *notificationUseCase) List(ctx context.Context, userID uuid.UUID, params domain.NotificationListParams) ([]domain.NotificationItem, *domain.PaginationMeta, error) {
	page, limit := normalizeNotificationPaging(params.Page, params.Limit)
	params.Page = page
	params.Limit = limit

	items, total, err := uc.repo.List(ctx, userID, params)
	if err != nil {
		return nil, nil, err
	}

	meta := buildNotificationMeta(page, limit, total)
	return items, meta, nil
}

func (uc *notificationUseCase) UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return uc.repo.UnreadCount(ctx, userID)
}

func (uc *notificationUseCase) MarkRead(ctx context.Context, userID, notificationID uuid.UUID) (*domain.NotificationItem, error) {
	item, err := uc.repo.MarkRead(ctx, userID, notificationID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNotificationNotFound
	}
	return item, nil
}

func (uc *notificationUseCase) MarkAllRead(ctx context.Context, userID uuid.UUID) (int64, error) {
	return uc.repo.MarkAllRead(ctx, userID)
}

func normalizeNotificationPaging(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}

func buildNotificationMeta(page, limit int, total int64) *domain.PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if limit == 0 {
		totalPages = 1
	}
	return &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}
