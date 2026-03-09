package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// AdminUserHandler handles admin user management endpoints.
type AdminUserHandler struct {
	useCase domain.AdminUserUseCase
}

// NewAdminUserHandler creates a new AdminUserHandler.
func NewAdminUserHandler(uc domain.AdminUserUseCase) *AdminUserHandler {
	return &AdminUserHandler{useCase: uc}
}

// Create handles POST /api/v1/admin/users
func (h *AdminUserHandler) Create(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.Create(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "user created successfully", result)
}

// List handles GET /api/v1/admin/users
func (h *AdminUserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	params := domain.UserListParams{
		Page:   page,
		Limit:  limit,
		Role:   c.Query("role"),
		Status: c.Query("status"),
		Search: c.Query("search"),
	}

	users, meta, err := h.useCase.List(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, "failed to list users", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "users retrieved successfully", users, meta)
}

// GetByID handles GET /api/v1/admin/users/:id
func (h *AdminUserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID", "id must be a valid UUID")
		return
	}

	result, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "user retrieved successfully", result)
}

// GetMe handles GET /api/v1/admin/me
func (h *AdminUserHandler) GetMe(c *gin.Context) {
	adminID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	result, err := h.useCase.GetMe(c.Request.Context(), adminID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "admin profile retrieved successfully", result)
}

// Update handles PUT /api/v1/admin/users/:id
func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID", "id must be a valid UUID")
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.Update(c.Request.Context(), id, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "user updated successfully", result)
}

// Delete handles DELETE /api/v1/admin/users/:id
func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid user ID", "id must be a valid UUID")
		return
	}

	actorID, _ := c.Get(middleware.ContextKeyUserID)
	actorUUID, _ := actorID.(uuid.UUID)

	if err := h.useCase.Delete(c.Request.Context(), id, actorUUID); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "user deleted successfully", nil)
}
