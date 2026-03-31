package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminPaymentHandler handles admin payment listing endpoints.
type AdminPaymentHandler struct {
	useCase domain.AdminPaymentUseCase
}

// NewAdminPaymentHandler creates a new AdminPaymentHandler.
func NewAdminPaymentHandler(uc domain.AdminPaymentUseCase) *AdminPaymentHandler {
	return &AdminPaymentHandler{useCase: uc}
}

// List handles GET /api/v1/admin/payments
func (h *AdminPaymentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.AdminPaymentListParams{
		Page:      page,
		Limit:     limit,
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	payments, meta, err := h.useCase.List(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to list payments", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "payments retrieved successfully", payments, meta)
}
