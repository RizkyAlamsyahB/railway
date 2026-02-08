package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// HealthHandler handles health check HTTP requests.
type HealthHandler struct {
	useCase domain.HealthUseCase
}

// NewHealthHandler creates a new HealthHandler with the given usecase.
func NewHealthHandler(uc domain.HealthUseCase) *HealthHandler {
	return &HealthHandler{useCase: uc}
}

// Check handles GET /api/v1/health and returns a Hello World response.
func (h *HealthHandler) Check(c *gin.Context) {
	result := h.useCase.Check(c.Request.Context())
	response.Success(c, http.StatusOK, "OK", result)
}
