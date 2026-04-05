package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// UserBankAccountHandler handles HTTP requests for user bank account management.
type UserBankAccountHandler struct {
	uc domain.UserBankAccountUseCase
}

// NewUserBankAccountHandler creates a new UserBankAccountHandler.
func NewUserBankAccountHandler(uc domain.UserBankAccountUseCase) *UserBankAccountHandler {
	return &UserBankAccountHandler{uc: uc}
}

// Create handles POST /users/bank-accounts
func (h *UserBankAccountHandler) Create(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	var req domain.CreateUserBankAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	result, err := h.uc.Create(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "bank account created", result)
}

// List handles GET /users/bank-accounts
func (h *UserBankAccountHandler) List(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	items, err := h.uc.List(c.Request.Context(), userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "bank accounts loaded", items)
}

// Delete handles DELETE /users/bank-accounts/:id
func (h *UserBankAccountHandler) Delete(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid bank account id", nil)
		return
	}

	if err := h.uc.Delete(c.Request.Context(), userID, id); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "bank account deleted", nil)
}

// SetDefault handles PATCH /users/bank-accounts/:id/default
func (h *UserBankAccountHandler) SetDefault(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid bank account id", nil)
		return
	}

	if err := h.uc.SetDefault(c.Request.Context(), userID, id); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "default bank account updated", nil)
}
