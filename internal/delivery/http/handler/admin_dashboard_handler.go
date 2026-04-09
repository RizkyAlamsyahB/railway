package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminDashboardHandler handles admin dashboard endpoints.
type AdminDashboardHandler struct {
	useCase domain.AdminDashboardUseCase
}

// NewAdminDashboardHandler creates a new AdminDashboardHandler.
func NewAdminDashboardHandler(useCase domain.AdminDashboardUseCase) *AdminDashboardHandler {
	return &AdminDashboardHandler{useCase: useCase}
}

// GetDashboard handles GET /api/v1/admin/dashboard
func (h *AdminDashboardHandler) GetDashboard(c *gin.Context) {
	now := time.Now()
	month := queryInt(c, "month", int(now.Month()))
	year := queryInt(c, "year", now.Year())

	if month < 1 || month > 12 {
		response.BadRequest(c, "month must be between 1 and 12", nil)
		return
	}
	if year < 2000 || year > 2100 {
		response.BadRequest(c, "year must be between 2000 and 2100", nil)
		return
	}

	result, err := h.useCase.GetDashboard(c.Request.Context(), domain.AdminDashboardParams{
		Month: month,
		Year:  year,
	})
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "admin dashboard retrieved successfully", result)
}
