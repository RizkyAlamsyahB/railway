package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// VendorOrderHandler handles vendor order management endpoints.
type VendorOrderHandler struct {
	useCase domain.VendorOrderUseCase
}

// NewVendorOrderHandler creates a new VendorOrderHandler.
func NewVendorOrderHandler(useCase domain.VendorOrderUseCase) *VendorOrderHandler {
	return &VendorOrderHandler{useCase: useCase}
}

// ListOrders handles GET /api/v1/vendors/orders.
func (h *VendorOrderHandler) ListOrders(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, meta, err := h.useCase.ListOrders(c.Request.Context(), vendorID, domain.VendorOrderListParams{
		Page:   page,
		Limit:  limit,
		Status: c.Query("status"),
	})
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendor orders retrieved successfully", items, meta)
}

// GetOrderDetail handles GET /api/v1/vendors/orders/:orderId.
func (h *VendorOrderHandler) GetOrderDetail(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.GetOrderDetail(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order detail retrieved successfully", res)
}

// AcceptOrder handles POST /api/v1/vendors/orders/:orderId/accept.
func (h *VendorOrderHandler) AcceptOrder(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "user ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.AcceptOrder(c.Request.Context(), vendorID, orderID, userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order accepted successfully", res)
}

// RejectOrder handles POST /api/v1/vendors/orders/:orderId/reject.
func (h *VendorOrderHandler) RejectOrder(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "user ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.RejectOrder(c.Request.Context(), vendorID, orderID, userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order rejected successfully", res)
}

// ShipOrder handles POST /api/v1/vendors/orders/:orderId/ship.
func (h *VendorOrderHandler) ShipOrder(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "user ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	var req domain.ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	res, err := h.useCase.ShipOrder(c.Request.Context(), vendorID, orderID, userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "order shipped successfully", res)
}

// TrackWaybill handles GET /api/v1/vendors/orders/:orderId/tracking.
func (h *VendorOrderHandler) TrackWaybill(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "vendor ID not found in token", nil)
		return
	}

	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		response.BadRequest(c, "invalid order ID", "orderId must be a valid UUID")
		return
	}

	res, err := h.useCase.TrackWaybill(c.Request.Context(), vendorID, orderID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "tracking info retrieved successfully", res)
}
