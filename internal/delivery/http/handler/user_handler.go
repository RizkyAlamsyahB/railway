package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

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
		h.handleError(c, err)
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
		h.handleError(c, err)
		return
	}

	response.OK(c, "If your email is registered and pending verification, a new verification email has been sent.", nil)
}

// handleError maps usecase sentinel errors to HTTP responses.
func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrEmailAlreadyRegistered):
		response.Error(c, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, usecase.ErrPhoneAlreadyRegistered):
		response.Error(c, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, usecase.ErrInvalidVerificationToken):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrUserAlreadyActive):
		response.Error(c, http.StatusConflict, err.Error(), nil)
	case errors.Is(err, usecase.ErrResendTooSoon):
		response.Error(c, http.StatusTooManyRequests, err.Error(), nil)
	case errors.Is(err, usecase.ErrUserNotPending):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.InternalServerError(c, "An unexpected error occurred", nil)
	}
}
