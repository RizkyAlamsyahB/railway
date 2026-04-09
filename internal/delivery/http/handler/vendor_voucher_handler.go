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

// VendorVoucherHandler handles HTTP requests for vendor voucher management.
type VendorVoucherHandler struct {
	uc domain.VendorVoucherUseCase
}

// NewVendorVoucherHandler creates a new VendorVoucherHandler.
func NewVendorVoucherHandler(uc domain.VendorVoucherUseCase) *VendorVoucherHandler {
	return &VendorVoucherHandler{uc: uc}
}

func (h *VendorVoucherHandler) vendorIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	vendorIDVal, exists := c.Get(middleware.ContextKeyVendorID)
	if !exists {
		response.Unauthorized(c, "vendor id not found in token", nil)
		return uuid.Nil, false
	}
	vendorID, ok := vendorIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "invalid vendor id in token", nil)
		return uuid.Nil, false
	}
	return vendorID, true
}

// CreateVoucher handles POST /api/v1/vendors/vouchers
func (h *VendorVoucherHandler) CreateVoucher(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	var req domain.CreateVendorVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.uc.CreateVoucher(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "voucher created successfully", result)
}

// ListVouchers handles GET /api/v1/vendors/vouchers
func (h *VendorVoucherHandler) ListVouchers(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.VendorVoucherListParams{
		Page:   page,
		Limit:  limit,
		Search: c.Query("search"),
	}

	if raw := c.Query("is_active"); raw != "" {
		val, err := strconv.ParseBool(raw)
		if err != nil {
			response.BadRequest(c, "invalid is_active value", nil)
			return
		}
		params.IsActive = &val
	}

	if raw := c.Query("product_id"); raw != "" {
		productID, err := uuid.Parse(raw)
		if err != nil {
			response.BadRequest(c, "invalid product_id", nil)
			return
		}
		params.ProductID = &productID
	}

	items, meta, err := h.uc.ListVouchers(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendor vouchers retrieved successfully", items, meta)
}

// GetVoucher handles GET /api/v1/vendors/vouchers/:id
func (h *VendorVoucherHandler) GetVoucher(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid voucher id", nil)
		return
	}

	result, err := h.uc.GetVoucher(c.Request.Context(), vendorID, id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "voucher retrieved successfully", result)
}

// UpdateVoucher handles PATCH /api/v1/vendors/vouchers/:id
func (h *VendorVoucherHandler) UpdateVoucher(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid voucher id", nil)
		return
	}

	var req domain.UpdateVendorVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.uc.UpdateVoucher(c.Request.Context(), vendorID, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "voucher updated successfully", result)
}

// DeleteVoucher handles DELETE /api/v1/vendors/vouchers/:id
func (h *VendorVoucherHandler) DeleteVoucher(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid voucher id", nil)
		return
	}

	if err := h.uc.DeleteVoucher(c.Request.Context(), vendorID, id); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "voucher deleted successfully", nil)
}

// GetSummary handles GET /api/v1/vendors/vouchers/summary
func (h *VendorVoucherHandler) GetSummary(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	summary, err := h.uc.GetSummary(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "voucher summary retrieved successfully", summary)
}
