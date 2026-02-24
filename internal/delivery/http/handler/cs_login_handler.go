package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// CSLoginHandler handles customer service authentication endpoints.
type CSLoginHandler struct {
	useCase domain.CSAuthUseCase
}

// NewCSLoginHandler creates a new CSLoginHandler.
func NewCSLoginHandler(uc domain.CSAuthUseCase) *CSLoginHandler {
	return &CSLoginHandler{useCase: uc}
}

// Login handles POST /api/v1/cs/login
func (h *CSLoginHandler) Login(c *gin.Context) {
	var req domain.CSLoginRequest
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
