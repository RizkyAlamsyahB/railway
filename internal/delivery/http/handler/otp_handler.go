package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// OTPHandler handles public OTP request and verification endpoints.
type OTPHandler struct {
	useCase domain.OTPUseCase
}

// NewOTPHandler creates a new OTPHandler.
func NewOTPHandler(useCase domain.OTPUseCase) *OTPHandler {
	return &OTPHandler{useCase: useCase}
}

// RequestOTP handles POST /api/v1/otp/request.
func (h *OTPHandler) RequestOTP(c *gin.Context) {
	var req domain.RequestOTPInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	result, err := h.useCase.RequestOTP(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "OTP generated successfully", result)
}

// VerifyOTP handles POST /api/v1/otp/verify.
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req domain.VerifyOTPInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err)
		return
	}

	result, err := h.useCase.VerifyOTP(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "OTP verified successfully", result)
}
