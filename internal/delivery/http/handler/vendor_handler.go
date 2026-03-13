package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

func extractVendorID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(middleware.ContextKeyVendorID)
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

func extractVendorOnboardingID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(middleware.ContextKeyVendorOnboardingID)
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

// VendorHandler handles vendor-related endpoints.
type VendorHandler struct {
	useCase domain.VendorUseCase
}

// NewVendorHandler creates a new VendorHandler.
func NewVendorHandler(uc domain.VendorUseCase) *VendorHandler {
	return &VendorHandler{useCase: uc}
}

// RequestRegistrationOTP handles POST /api/v1/vendors/register/request-otp
func (h *VendorHandler) RequestRegistrationOTP(c *gin.Context) {
	var req domain.VendorRegisterOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.RequestRegistrationOTP(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "vendor registration otp sent", result)
}

// VerifyRegistrationOTP handles POST /api/v1/vendors/register/verify-otp
func (h *VendorHandler) VerifyRegistrationOTP(c *gin.Context) {
	var req domain.VendorVerifyRegistrationOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.VerifyRegistrationOTP(c.Request.Context(), req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor registration otp verified", result)
}

// SetRegistrationPassword handles POST /api/v1/vendors/register/password
func (h *VendorHandler) SetRegistrationPassword(c *gin.Context) {
	var req domain.VendorRegistrationPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	onboardingID, ok := extractVendorOnboardingID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor onboarding ID in token", nil)
		return
	}

	result, err := h.useCase.SetRegistrationPassword(c.Request.Context(), onboardingID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor registration password saved", result)
}

// SaveRegistrationStore handles POST /api/v1/vendors/register/store
func (h *VendorHandler) SaveRegistrationStore(c *gin.Context) {
	var req domain.VendorRegistrationStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	onboardingID, ok := extractVendorOnboardingID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor onboarding ID in token", nil)
		return
	}

	result, err := h.useCase.SaveRegistrationStore(c.Request.Context(), onboardingID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor registration store saved", result)
}

// PresignRegistrationLegalDocument handles POST /api/v1/vendors/register/legal-document/presign
func (h *VendorHandler) PresignRegistrationLegalDocument(c *gin.Context) {
	var req domain.VendorRegistrationPresignDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	onboardingID, ok := extractVendorOnboardingID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor onboarding ID in token", nil)
		return
	}

	result, err := h.useCase.PresignRegistrationLegalDocument(c.Request.Context(), onboardingID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor registration legal document presigned", result)
}

// SubmitRegistrationLegal handles POST /api/v1/vendors/register/legal-document/submit
func (h *VendorHandler) SubmitRegistrationLegal(c *gin.Context) {
	var req domain.VendorRegistrationLegalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	onboardingID, ok := extractVendorOnboardingID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor onboarding ID in token", nil)
		return
	}

	result, err := h.useCase.SubmitRegistrationLegal(c.Request.Context(), onboardingID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor registration completed", result)
}

// Login handles POST /api/v1/vendors/login
func (h *VendorHandler) Login(c *gin.Context) {
	var req domain.VendorLoginRequest
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

// GetMe handles GET /api/v1/vendors/me
func (h *VendorHandler) GetMe(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	result, err := h.useCase.GetMe(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor profile retrieved successfully", result)
}

// ConfirmDocuments handles POST /api/v1/vendors/documents/confirm
func (h *VendorHandler) ConfirmDocuments(c *gin.Context) {
	var req domain.ConfirmDocumentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, _ := userIDVal.(uuid.UUID)

	result, err := h.useCase.ConfirmDocuments(c.Request.Context(), userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "documents confirmed successfully", result)
}

// GetBalance handles GET /api/v1/vendors/balance
func (h *VendorHandler) GetBalance(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.GetBalance(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor balance retrieved", result)
}

// ListPayoutChannels handles GET /api/v1/vendors/payout-channels
func (h *VendorHandler) ListPayoutChannels(c *gin.Context) {
	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.ListPayoutChannels(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "payout channels retrieved successfully", result)
}

// RequestWithdrawal handles POST /api/v1/vendors/withdrawals
func (h *VendorHandler) RequestWithdrawal(c *gin.Context) {
	var req domain.VendorWithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorIDVal, _ := c.Get(middleware.ContextKeyVendorID)
	vendorID, _ := vendorIDVal.(uuid.UUID)

	result, err := h.useCase.RequestWithdrawal(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "withdrawal request submitted successfully", result)
}
