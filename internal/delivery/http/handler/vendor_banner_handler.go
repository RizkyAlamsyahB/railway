package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// VendorBannerHandler handles HTTP requests for vendor banner management.
type VendorBannerHandler struct {
	uc domain.VendorBannerUseCase
}

// NewVendorBannerHandler creates a new VendorBannerHandler.
func NewVendorBannerHandler(uc domain.VendorBannerUseCase) *VendorBannerHandler {
	return &VendorBannerHandler{uc: uc}
}

// vendorIDFromContext extracts vendor_id from auth context.
func (h *VendorBannerHandler) vendorIDFromContext(c *gin.Context) (uuid.UUID, bool) {
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

// CreateBanner godoc
// POST /vendors/banners
// Creates a vendor banner and returns a presigned upload URL for the image.
func (h *VendorBannerHandler) CreateBanner(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	var req domain.CreateVendorBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.CreateBanner(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "vendor banner created, upload image using the presigned URL", result)
}

// ConfirmBanner godoc
// POST /vendors/banners/:id/confirm
// Confirms that the banner image has been uploaded to S3.
func (h *VendorBannerHandler) ConfirmBanner(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	var req domain.ConfirmVendorBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.ConfirmBanner(c.Request.Context(), vendorID, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "vendor banner image confirmed", result)
}

// ListBanners godoc
// GET /vendors/banners
// Lists all banners for the authenticated vendor.
func (h *VendorBannerHandler) ListBanners(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	banners, err := h.uc.ListBanners(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "vendor banners loaded", banners)
}

// GetBanner godoc
// GET /vendors/banners/:id
func (h *VendorBannerHandler) GetBanner(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	result, err := h.uc.GetBanner(c.Request.Context(), vendorID, id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "vendor banner loaded", result)
}

// UpdateBanner godoc
// PATCH /vendors/banners/:id
func (h *VendorBannerHandler) UpdateBanner(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	var req domain.UpdateVendorBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.UpdateBanner(c.Request.Context(), vendorID, id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "vendor banner updated", result)
}

// DeleteBanner godoc
// DELETE /vendors/banners/:id
func (h *VendorBannerHandler) DeleteBanner(c *gin.Context) {
	vendorID, ok := h.vendorIDFromContext(c)
	if !ok {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	if err := h.uc.DeleteBanner(c.Request.Context(), vendorID, id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "vendor banner deleted", nil)
}

// ListVendorBannersPublic godoc
// GET /vendors/:id/banners (public)
// Lists banners for a specific vendor store page.
func (h *VendorBannerHandler) ListVendorBannersPublic(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor id", nil)
		return
	}

	banners, err := h.uc.ListByVendorIDPublic(c.Request.Context(), vendorID)
	if err != nil {
		response.InternalServerError(c, "failed to load vendor banners", nil)
		return
	}
	response.OK(c, "vendor banners loaded", banners)
}
