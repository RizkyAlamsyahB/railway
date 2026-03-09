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

// VendorCourierHandler handles HTTP requests for vendor courier management.
type VendorCourierHandler struct {
	uc domain.VendorCourierUseCase
}

// NewVendorCourierHandler creates a new VendorCourierHandler.
func NewVendorCourierHandler(uc domain.VendorCourierUseCase) *VendorCourierHandler {
	return &VendorCourierHandler{uc: uc}
}

// ListCouriers returns all available couriers.
// GET /api/v1/shipping/couriers
func (h *VendorCourierHandler) ListCouriers(c *gin.Context) {
	couriers, err := h.uc.ListCouriers(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "couriers loaded", couriers)
}

// GetVendorCouriers returns the authenticated vendor's selected couriers.
// GET /api/v1/vendors/couriers
func (h *VendorCourierHandler) GetVendorCouriers(c *gin.Context) {
	vendorID := c.MustGet(middleware.ContextKeyVendorID).(uuid.UUID)

	couriers, err := h.uc.GetVendorCouriers(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "vendor couriers loaded", couriers)
}

// SetVendorCouriers replaces the authenticated vendor's courier list.
// PUT /api/v1/vendors/couriers
func (h *VendorCourierHandler) SetVendorCouriers(c *gin.Context) {
	vendorID := c.MustGet(middleware.ContextKeyVendorID).(uuid.UUID)

	var req domain.SetVendorCouriersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	couriers, err := h.uc.SetVendorCouriers(c.Request.Context(), vendorID, req.CourierIDs)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "vendor couriers updated", couriers)
}

// RemoveVendorCourier removes a specific courier from the vendor's list.
// DELETE /api/v1/vendors/couriers/:courierId
func (h *VendorCourierHandler) RemoveVendorCourier(c *gin.Context) {
	vendorID := c.MustGet(middleware.ContextKeyVendorID).(uuid.UUID)

	courierID, err := strconv.Atoi(c.Param("courierId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid courier id", nil)
		return
	}

	if err := h.uc.RemoveVendorCourier(c.Request.Context(), vendorID, courierID); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "courier removed from vendor", nil)
}
