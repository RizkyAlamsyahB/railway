package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// VendorDashboardHandler handles vendor seller dashboard endpoints.
type VendorDashboardHandler struct {
	useCase domain.VendorDashboardUseCase
}

// NewVendorDashboardHandler creates a new VendorDashboardHandler.
func NewVendorDashboardHandler(useCase domain.VendorDashboardUseCase) *VendorDashboardHandler {
	return &VendorDashboardHandler{useCase: useCase}
}

// GetDashboard handles GET /api/v1/vendors/dashboard
func (h *VendorDashboardHandler) GetDashboard(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

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

	result, err := h.useCase.GetDashboard(c.Request.Context(), vendorID, month, year)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor dashboard retrieved successfully", result)
}
