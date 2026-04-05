package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// StoreHandler handles public store detail endpoints.
type StoreHandler struct {
	uc domain.StoreUseCase
}

// NewStoreHandler creates a new StoreHandler.
func NewStoreHandler(uc domain.StoreUseCase) *StoreHandler {
	return &StoreHandler{uc: uc}
}

// GetStoreDetail handles GET /api/v1/vendors/:id
// Returns the public store detail page with store info, banners, and product sections.
func (h *StoreHandler) GetStoreDetail(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor id", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.uc.GetStoreDetail(c.Request.Context(), vendorID, page, limit)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "store detail loaded", result)
}

// GetStoreReviews handles GET /api/v1/vendors/:id/reviews
func (h *StoreHandler) GetStoreReviews(c *gin.Context) {
	vendorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid vendor id", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.uc.GetStoreReviews(c.Request.Context(), vendorID, page, limit)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "store reviews loaded", result)
}
