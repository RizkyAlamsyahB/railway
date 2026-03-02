package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// FAQHandler handles HTTP requests for the FAQ / Pusat Bantuan feature.
type FAQHandler struct {
	uc domain.FAQUseCase
}

// NewFAQHandler creates a new FAQHandler.
func NewFAQHandler(uc domain.FAQUseCase) *FAQHandler {
	return &FAQHandler{uc: uc}
}

// ListFAQs godoc
// GET /api/v1/faq?category=&search=
func (h *FAQHandler) ListFAQs(c *gin.Context) {
	params := domain.FAQListParams{
		Category: c.Query("category"),
		Search:   c.Query("search"),
	}

	groups, err := h.uc.ListFAQs(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load FAQ", nil)
		return
	}

	response.Success(c, http.StatusOK, "FAQ loaded", groups)
}

// ListCategories godoc
// GET /api/v1/faq/categories
func (h *FAQHandler) ListCategories(c *gin.Context) {
	cats, err := h.uc.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to load categories", nil)
		return
	}
	response.Success(c, http.StatusOK, "FAQ categories loaded", cats)
}
