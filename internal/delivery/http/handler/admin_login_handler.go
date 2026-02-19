package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
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
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "login successful", result)
}
