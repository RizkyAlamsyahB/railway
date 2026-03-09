package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// ShippingHandler handles HTTP requests for shipping location lookups (RajaOngkir proxy).
type ShippingHandler struct {
	uc domain.ShippingUseCase
}

// NewShippingHandler creates a new ShippingHandler.
func NewShippingHandler(uc domain.ShippingUseCase) *ShippingHandler {
	return &ShippingHandler{uc: uc}
}

// GetProvinces godoc
// GET /api/v1/shipping/provinces
func (h *ShippingHandler) GetProvinces(c *gin.Context) {
	provinces, err := h.uc.GetProvinces(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "provinces loaded", provinces)
}

// GetCities godoc
// GET /api/v1/shipping/cities?province_id=
func (h *ShippingHandler) GetCities(c *gin.Context) {
	provinceID := c.Query("province_id")
	if provinceID == "" {
		response.Error(c, http.StatusBadRequest, "province_id is required", nil)
		return
	}

	cities, err := h.uc.GetCities(c.Request.Context(), provinceID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "cities loaded", cities)
}

// GetDistricts godoc
// GET /api/v1/shipping/districts?city_id=
func (h *ShippingHandler) GetDistricts(c *gin.Context) {
	cityID := c.Query("city_id")
	if cityID == "" {
		response.Error(c, http.StatusBadRequest, "city_id is required", nil)
		return
	}

	districts, err := h.uc.GetDistricts(c.Request.Context(), cityID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "districts loaded", districts)
}

// GetSubdistricts godoc
// GET /api/v1/shipping/subdistricts?district_id=
func (h *ShippingHandler) GetSubdistricts(c *gin.Context) {
	districtID := c.Query("district_id")
	if districtID == "" {
		response.Error(c, http.StatusBadRequest, "district_id is required", nil)
		return
	}

	subdistricts, err := h.uc.GetSubdistricts(c.Request.Context(), districtID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "subdistricts loaded", subdistricts)
}
