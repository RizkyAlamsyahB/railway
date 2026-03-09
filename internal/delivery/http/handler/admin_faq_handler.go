package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminFAQHandler handles admin CRUD for FAQ items.
type AdminFAQHandler struct {
	uc domain.AdminFAQUseCase
}

// NewAdminFAQHandler creates a new AdminFAQHandler.
func NewAdminFAQHandler(uc domain.AdminFAQUseCase) *AdminFAQHandler {
	return &AdminFAQHandler{uc: uc}
}

// CreateFAQ handles POST /admin/settings/faq
func (h *AdminFAQHandler) CreateFAQ(c *gin.Context) {
	var req domain.CreateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.CreateFAQ(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "FAQ created", result)
}

// ListFAQs handles GET /admin/settings/faq
func (h *AdminFAQHandler) ListFAQs(c *gin.Context) {
	faqs, err := h.uc.ListFAQs(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "FAQs loaded", faqs)
}

// UpdateFAQ handles PATCH /admin/settings/faq/:id
func (h *AdminFAQHandler) UpdateFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid FAQ id", nil)
		return
	}

	var req domain.UpdateFAQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.UpdateFAQ(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "FAQ updated", result)
}

// DeleteFAQ handles DELETE /admin/settings/faq/:id
func (h *AdminFAQHandler) DeleteFAQ(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid FAQ id", nil)
		return
	}

	if err := h.uc.DeleteFAQ(c.Request.Context(), id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "FAQ deleted", nil)
}
