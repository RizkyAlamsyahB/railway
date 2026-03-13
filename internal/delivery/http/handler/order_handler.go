package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// OrderActionHandler handles customer order action endpoints.
type OrderActionHandler struct {
	useCase domain.OrderActionUseCase
}

// NewOrderActionHandler creates a new OrderActionHandler.
func NewOrderActionHandler(useCase domain.OrderActionUseCase) *OrderActionHandler {
	return &OrderActionHandler{useCase: useCase}
}

// Complete handles POST /api/v1/users/orders/:orderId/complete.
func (h *OrderActionHandler) Complete(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.CompleteByCustomer(c.Request.Context(), userID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order marked as received", res)
}

// ListOrders handles GET /api/v1/users/orders.
func (h *OrderActionHandler) ListOrders(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, meta, err := h.useCase.ListByCustomer(c.Request.Context(), userID, domain.CustomerOrderListParams{
		Page:   page,
		Limit:  limit,
		Status: c.Query("status"),
	})
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "orders retrieved successfully", items, meta)
}

// ListStatuses handles GET /api/v1/order-statuses.
func (h *OrderActionHandler) ListStatuses(c *gin.Context) {
	statuses := h.useCase.ListOrderStatuses(c.Request.Context())
	response.OK(c, "order statuses loaded", statuses)
}
