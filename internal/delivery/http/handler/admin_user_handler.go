package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
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
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	response.OK(c, "user retrieved successfully", result)
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
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	response.OK(c, "user deleted successfully", nil)
}

func (h *AdminUserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		response.NotFound(c, "user not found", nil)
	case errors.Is(err, usecase.ErrEmailExists):
		response.BadRequest(c, "email already exists", nil)
	case errors.Is(err, usecase.ErrInvalidBirthDate):
		response.BadRequest(c, "invalid birth_date format, expected YYYY-MM-DD", nil)
	case errors.Is(err, usecase.ErrCannotDeleteSelf):
		response.BadRequest(c, "cannot delete your own account", nil)
	default:
		response.InternalServerError(c, "internal server error", err.Error())
	}
}
