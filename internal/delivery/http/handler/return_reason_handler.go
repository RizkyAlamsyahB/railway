package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// ReturnReasonHandler handles admin CRUD for return reasons.
type ReturnReasonHandler struct {
	uc domain.ReturnReasonUseCase
}

// NewReturnReasonHandler creates a new ReturnReasonHandler.
func NewReturnReasonHandler(uc domain.ReturnReasonUseCase) *ReturnReasonHandler {
	return &ReturnReasonHandler{uc: uc}
}

// CreateReturnReason handles POST /admin/settings/return-reasons
func (h *ReturnReasonHandler) CreateReturnReason(c *gin.Context) {
	var req domain.CreateReturnReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.CreateReturnReason(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "return reason created", result)
}

// ListReturnReasons handles GET /admin/settings/return-reasons
func (h *ReturnReasonHandler) ListReturnReasons(c *gin.Context) {
	reasons, err := h.uc.ListReturnReasons(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "return reasons loaded", reasons)
}

// ListActiveReturnReasons handles GET /return-reasons (public)
func (h *ReturnReasonHandler) ListActiveReturnReasons(c *gin.Context) {
	reasons, err := h.uc.ListReturnReasons(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "failed to load return reasons", nil)
		return
	}
	response.OK(c, "return reasons loaded", reasons)
}

// UpdateReturnReason handles PATCH /admin/settings/return-reasons/:id
func (h *ReturnReasonHandler) UpdateReturnReason(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid return reason id", nil)
		return
	}

	var req domain.UpdateReturnReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.UpdateReturnReason(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "return reason updated", result)
}

// DeleteReturnReason handles DELETE /admin/settings/return-reasons/:id
func (h *ReturnReasonHandler) DeleteReturnReason(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid return reason id", nil)
		return
	}

	if err := h.uc.DeleteReturnReason(c.Request.Context(), id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "return reason deleted", nil)
}
