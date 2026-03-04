package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// NotificationHandler handles notification endpoints.
type NotificationHandler struct {
	uc domain.NotificationUseCase
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(uc domain.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

// List handles GET /api/v1/notifications
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, meta, err := h.uc.List(c.Request.Context(), userID, domain.NotificationListParams{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "notifications retrieved successfully", items, meta)
}

// UnreadCount handles GET /api/v1/notifications/unread-count
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	count, err := h.uc.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "unread notifications retrieved successfully", domain.NotificationUnreadCount{
		UnreadCount: count,
	})
}

// MarkRead handles PATCH /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid notification ID", "id must be a valid UUID")
		return
	}

	item, err := h.uc.MarkRead(c.Request.Context(), userID, id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "notification marked as read", item)
}

// MarkAllRead handles PATCH /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	updated, err := h.uc.MarkAllRead(c.Request.Context(), userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "notifications marked as read", domain.NotificationBulkReadResponse{
		Updated: updated,
	})
}
