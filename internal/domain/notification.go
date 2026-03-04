package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Notification types.
const (
	NotificationTypeChat          = "chat"
	NotificationTypeTicket        = "ticket"
	NotificationTypeVendorProduct = "vendor_product"
	NotificationTypeOrder         = "order"
	NotificationTypeRefund        = "refund"
	NotificationTypePayout        = "payout"
	NotificationTypeSystem        = "system"
)

// NotificationItem represents a single notification.
type NotificationItem struct {
	ID        uuid.UUID  `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

// NotificationListParams holds pagination parameters.
type NotificationListParams struct {
	Page  int
	Limit int
}

// NotificationUnreadCount is the response payload for unread counts.
type NotificationUnreadCount struct {
	UnreadCount int64 `json:"unread_count"`
}

// NotificationBulkReadResponse is returned by mark-all-read.
type NotificationBulkReadResponse struct {
	Updated int64 `json:"updated"`
}

// NotificationRepository defines data access for notifications.
type NotificationRepository interface {
	List(ctx context.Context, userID uuid.UUID, params NotificationListParams) ([]NotificationItem, int64, error)
	UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkRead(ctx context.Context, userID, notificationID uuid.UUID) (*NotificationItem, error)
	MarkAllRead(ctx context.Context, userID uuid.UUID) (int64, error)
}

// NotificationUseCase defines notification use cases.
type NotificationUseCase interface {
	List(ctx context.Context, userID uuid.UUID, params NotificationListParams) ([]NotificationItem, *PaginationMeta, error)
	UnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkRead(ctx context.Context, userID, notificationID uuid.UUID) (*NotificationItem, error)
	MarkAllRead(ctx context.Context, userID uuid.UUID) (int64, error)
}
