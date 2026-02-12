package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
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
		h.handleError(c, err)
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
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	response.OK(c, "documents confirmed successfully", result)
}

func (h *VendorHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrEmailAlreadyRegistered):
		response.Error(c, 409, "email already registered", nil)
	case errors.Is(err, usecase.ErrVendorAlreadyExists):
		response.Error(c, 409, "user already has a vendor", nil)
	case errors.Is(err, usecase.ErrVendorNotFound):
		response.NotFound(c, "vendor not found", nil)
	case errors.Is(err, usecase.ErrDocumentNotFound):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrObjectNotUploaded):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrInvalidDocumentContent):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrDocumentSizeOverflow):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrVendorInvalidCredentials):
		response.Unauthorized(c, "invalid email or password", nil)
	case errors.Is(err, usecase.ErrVendorAccountBlocked):
		response.Forbidden(c, "vendor account is blocked", nil)
	case errors.Is(err, usecase.ErrNotVendor):
		response.Forbidden(c, "vendor access required", nil)
	case errors.Is(err, usecase.ErrNoVendorProfile):
		response.Forbidden(c, "no vendor profile found", nil)
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
}
