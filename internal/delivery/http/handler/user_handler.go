package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// extractUserID retrieves the authenticated user's UUID from the Gin context.
func extractUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		return uuid.UUID{}, false
	}
	switch v := val.(type) {
	case uuid.UUID:
		return v, true
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return uuid.UUID{}, false
		}
		return id, true
	default:
		return uuid.UUID{}, false
	}
}

// UserHandler handles customer-facing user HTTP requests.
type UserHandler struct {
	useCase     domain.UserUseCase
	frontendURL string
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(useCase domain.UserUseCase, frontendURL string) *UserHandler {
	return &UserHandler{
		useCase:     useCase,
		frontendURL: frontendURL,
	}
}

// Register handles POST /api/v1/users/register.
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	resp, err := h.useCase.Register(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "Registration successful. Please check your email to verify your account.", resp)
}

// VerifyEmail handles GET /api/v1/users/verify-email?token=xxx.
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/verify-email?status=failed&reason=missing_token", h.frontendURL))
		return
	}

	err := h.useCase.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		reason := "invalid_token"
		if errors.Is(err, usecase.ErrInvalidVerificationToken) {
			reason = "invalid_or_expired"
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/verify-email?status=failed&reason=%s", h.frontendURL, reason))
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/verify-email?status=success", h.frontendURL))
}

// ResendVerification handles POST /api/v1/users/resend-verification.
func (h *UserHandler) ResendVerification(c *gin.Context) {
	var req domain.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	err := h.useCase.ResendVerification(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "If your email is registered and pending verification, a new verification email has been sent.", nil)
}

// Login handles POST /api/v1/users/login.
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
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

// GetMe handles GET /api/v1/users/me.
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	result, err := h.useCase.GetMe(c.Request.Context(), userID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "user profile retrieved successfully", result)
}
