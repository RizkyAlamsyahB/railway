package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminVendorHandler handles admin vendor review endpoints.
type AdminVendorHandler struct {
	useCase domain.AdminVendorUseCase
}

// NewAdminVendorHandler creates a new AdminVendorHandler.
func NewAdminVendorHandler(uc domain.AdminVendorUseCase) *AdminVendorHandler {
	return &AdminVendorHandler{useCase: uc}
}

// handleError maps known usecase errors to proper HTTP responses.
func (h *AdminVendorHandler) handleError(c *gin.Context, err error, action string) {
	if errors.Is(err, usecase.ErrVendorNotFound) {
		response.NotFound(c, "vendor not found", nil)
		return
	}
	if errors.Is(err, usecase.ErrInvalidStatusTransition) {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.InternalServerError(c, "failed to "+action, err.Error())
}

// List handles GET /api/v1/admin/vendors
func (h *AdminVendorHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.VendorListParams{
		Page:       page,
		Limit:      limit,
		Status:     c.Query("status"),
		VendorType: c.Query("vendor_type"),
		Search:     c.Query("search"),
	}

	vendors, meta, err := h.useCase.List(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to list vendors", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendors retrieved successfully", vendors, meta)
}

// GetByID handles GET /api/v1/admin/vendors/:id
func (h *AdminVendorHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor ID", "id must be a valid UUID")
		return
	}

	vendor, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err, "get vendor")
		return
	}

	response.OK(c, "vendor retrieved successfully", vendor)
}

// Approve handles PATCH /api/v1/admin/vendors/:id/approve
func (h *AdminVendorHandler) Approve(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor ID", "id must be a valid UUID")
		return
	}

	adminIDVal, _ := c.Get(middleware.ContextKeyUserID)
	adminID, _ := adminIDVal.(uuid.UUID)

	result, err := h.useCase.Approve(c.Request.Context(), vendorID, adminID)
	if err != nil {
		h.handleError(c, err, "approve vendor")
		return
	}

	response.OK(c, result.Message, result)
}

// Reject handles PATCH /api/v1/admin/vendors/:id/reject
func (h *AdminVendorHandler) Reject(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor ID", "id must be a valid UUID")
		return
	}

	var req domain.AdminVendorReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required", err.Error())
		return
	}

	result, err := h.useCase.Reject(c.Request.Context(), vendorID, req.Reason)
	if err != nil {
		h.handleError(c, err, "reject vendor")
		return
	}

	response.OK(c, result.Message, result)
}

// Block handles PATCH /api/v1/admin/vendors/:id/block
func (h *AdminVendorHandler) Block(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor ID", "id must be a valid UUID")
		return
	}

	var req domain.AdminVendorReasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "reason is required", err.Error())
		return
	}

	result, err := h.useCase.Block(c.Request.Context(), vendorID, req.Reason)
	if err != nil {
		h.handleError(c, err, "block vendor")
		return
	}

	response.OK(c, result.Message, result)
}

// Unblock handles PATCH /api/v1/admin/vendors/:id/unblock
func (h *AdminVendorHandler) Unblock(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor ID", "id must be a valid UUID")
		return
	}

	adminIDVal, _ := c.Get(middleware.ContextKeyUserID)
	adminID, _ := adminIDVal.(uuid.UUID)

	result, err := h.useCase.Unblock(c.Request.Context(), vendorID, adminID)
	if err != nil {
		h.handleError(c, err, "unblock vendor")
		return
	}

	response.OK(c, result.Message, result)
}
