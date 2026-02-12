package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		if errors.Is(err, usecase.ErrVendorNotFound) {
			response.NotFound(c, "vendor not found", nil)
			return
		}
		response.InternalServerError(c, "failed to get vendor", err.Error())
		return
	}

	response.OK(c, "vendor retrieved successfully", vendor)
}
