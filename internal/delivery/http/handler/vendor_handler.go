package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// VendorHandler handles vendor-related endpoints.
type VendorHandler struct {
	useCase domain.VendorUseCase
}

// NewVendorHandler creates a new VendorHandler.
func NewVendorHandler(uc domain.VendorUseCase) *VendorHandler {
	return &VendorHandler{useCase: uc}
}

// Register handles POST /api/v1/vendors/register
func (h *VendorHandler) Register(c *gin.Context) {
	var req domain.VendorRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.Register(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "vendor registration successful", result)
}

// Login handles POST /api/v1/vendors/login
func (h *VendorHandler) Login(c *gin.Context) {
	var req domain.VendorLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.Login(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "login successful", result)
}

// ConfirmDocuments handles POST /api/v1/vendors/documents/confirm
func (h *VendorHandler) ConfirmDocuments(c *gin.Context) {
	var req domain.ConfirmDocumentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	result, err := h.useCase.ConfirmDocuments(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "documents confirmed successfully", result)
}

// GetBalance handles GET /api/v1/vendors/balance
func (h *VendorHandler) GetBalance(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.GetBalance(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor balance retrieved", result)
}

// RequestWithdrawal handles POST /api/v1/vendors/withdrawals
func (h *VendorHandler) RequestWithdrawal(c *gin.Context) {
	var req domain.VendorWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.RequestWithdrawal(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "withdrawal request submitted successfully", result)
}
