package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminLoginHandler handles admin authentication endpoints.
type AdminLoginHandler struct {
	useCase domain.AdminAuthUseCase
}

// NewAdminLoginHandler creates a new AdminLoginHandler.
func NewAdminLoginHandler(uc domain.AdminAuthUseCase) *AdminLoginHandler {
	return &AdminLoginHandler{useCase: uc}
}

// Login handles POST /api/v1/admin/login
func (h *AdminLoginHandler) Login(c *gin.Context) {
	var req domain.AdminLoginRequest
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

func (h *AdminLoginHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials):
		response.Unauthorized(c, "invalid email or password", nil)
	case errors.Is(err, usecase.ErrAccountInactive):
		response.Unauthorized(c, "account is not active", nil)
	case errors.Is(err, usecase.ErrNotAdmin):
		response.Forbidden(c, "admin access required", nil)
	default:
		response.InternalServerError(c, "internal server error", nil)
	}
}
