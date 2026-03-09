package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// BannerHandler handles HTTP requests for banner management.
type BannerHandler struct {
	uc domain.BannerUseCase
}

// NewBannerHandler creates a new BannerHandler.
func NewBannerHandler(uc domain.BannerUseCase) *BannerHandler {
	return &BannerHandler{uc: uc}
}

// CreateBanner godoc
// POST /admin/settings/banners
// Creates a banner record and returns a presigned upload URL for the image.
func (h *BannerHandler) CreateBanner(c *gin.Context) {
	var req domain.CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.CreateBanner(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "banner created, upload image using the presigned URL", result)
}

// ConfirmBanner godoc
// POST /admin/settings/banners/:id/confirm
// Confirms that the banner image has been uploaded to S3.
func (h *BannerHandler) ConfirmBanner(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	var req domain.ConfirmBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.ConfirmBanner(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "banner image confirmed", result)
}

// ListBanners godoc
// GET /admin/settings/banners
// Lists all banners.
func (h *BannerHandler) ListBanners(c *gin.Context) {
	banners, err := h.uc.ListBanners(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "banners loaded", banners)
}

// ListActiveBanners godoc
// GET /banners
// Lists banners for public/customer display.
func (h *BannerHandler) ListActiveBanners(c *gin.Context) {
	banners, err := h.uc.ListBanners(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "failed to load banners", nil)
		return
	}
	response.OK(c, "banners loaded", banners)
}

// GetBanner godoc
// GET /admin/settings/banners/:id
func (h *BannerHandler) GetBanner(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	result, err := h.uc.GetBanner(c.Request.Context(), id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "banner loaded", result)
}

// UpdateBanner godoc
// PATCH /admin/settings/banners/:id
func (h *BannerHandler) UpdateBanner(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	var req domain.UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.UpdateBanner(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "banner updated", result)
}

// DeleteBanner godoc
// DELETE /admin/settings/banners/:id
func (h *BannerHandler) DeleteBanner(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid banner id", nil)
		return
	}

	if err := h.uc.DeleteBanner(c.Request.Context(), id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "banner deleted", nil)
}
