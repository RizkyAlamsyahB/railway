package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminCategoryHandler handles admin CRUD for product categories.
type AdminCategoryHandler struct {
	uc domain.AdminCategoryUseCase
}

// NewAdminCategoryHandler creates a new AdminCategoryHandler.
func NewAdminCategoryHandler(uc domain.AdminCategoryUseCase) *AdminCategoryHandler {
	return &AdminCategoryHandler{uc: uc}
}

// CreateCategory handles POST /admin/settings/categories
func (h *AdminCategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.CreateCategory(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.Created(c, "category created", result)
}

// ListCategories handles GET /admin/settings/categories
func (h *AdminCategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.uc.ListCategories(c.Request.Context())
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "categories loaded", categories)
}

// UpdateCategory handles PATCH /admin/settings/categories/:id
func (h *AdminCategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid category id", nil)
		return
	}

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "category updated", result)
}

// DeleteCategory handles DELETE /admin/settings/categories/:id
func (h *AdminCategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid category id", nil)
		return
	}

	if err := h.uc.DeleteCategory(c.Request.Context(), id); err != nil {
		HandleUsecaseError(c, err)
		return
	}
	response.OK(c, "category deleted", nil)
}
