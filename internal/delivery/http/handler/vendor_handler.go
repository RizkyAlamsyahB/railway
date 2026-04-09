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

// UpdateProfile handles PUT /api/v1/vendors/me/profile
func (h *VendorHandler) UpdateProfile(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	var req domain.VendorUpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.UpdateProfile(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor profile updated successfully", result)
}

// SubmitSouvenirStoreProposal handles POST /api/v1/vendors/register/propose/souvenir_store
func (h *VendorHandler) SubmitSouvenirStoreProposal(c *gin.Context) {
	var req domain.VendorSouvenirStoreProposalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	result, err := h.useCase.SubmitSouvenirStoreProposal(c.Request.Context(), vendorID, userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "vendor proposal submitted successfully", result)
}

// PresignSouvenirStoreProposalDocuments handles POST /api/v1/vendors/register/propose/souvenir_store/presign
func (h *VendorHandler) PresignSouvenirStoreProposalDocuments(c *gin.Context) {
	var req domain.VendorSouvenirStoreProposalPresignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	userID, ok := extractUserID(c)
	if !ok {
		response.BadRequest(c, "invalid user ID in token", nil)
		return
	}

	result, err := h.useCase.PresignSouvenirStoreProposalDocuments(c.Request.Context(), vendorID, userID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "presigned upload URL generated", result)
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

// CreateBankAccount handles POST /api/v1/vendors/bank-accounts
func (h *VendorHandler) CreateBankAccount(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	var req domain.CreateVendorBankAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "validation failed", err.Error())
		return
	}

	result, err := h.useCase.CreateBankAccount(c.Request.Context(), vendorID, req)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.Created(c, "bank account created", result)
}

// ListBankAccounts handles GET /api/v1/vendors/bank-accounts
func (h *VendorHandler) ListBankAccounts(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	items, err := h.useCase.ListBankAccounts(c.Request.Context(), vendorID)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "bank accounts loaded", items)
}

// DeleteBankAccount handles DELETE /api/v1/vendors/bank-accounts/:id
func (h *VendorHandler) DeleteBankAccount(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid bank account id", nil)
		return
	}

	if err := h.useCase.DeleteBankAccount(c.Request.Context(), vendorID, id); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "bank account deleted", nil)
}

// SetDefaultBankAccount handles PATCH /api/v1/vendors/bank-accounts/:id/default
func (h *VendorHandler) SetDefaultBankAccount(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid bank account id", nil)
		return
	}

	if err := h.useCase.SetDefaultBankAccount(c.Request.Context(), vendorID, id); err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.OK(c, "default bank account updated", nil)
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

// ListWithdrawals handles GET /api/v1/vendors/withdrawals
func (h *VendorHandler) ListWithdrawals(c *gin.Context) {
	vendorID, ok := extractVendorID(c)
	if !ok {
		response.BadRequest(c, "invalid vendor ID in token", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	params := domain.VendorWithdrawalListParams{
		Page:   page,
		Limit:  limit,
		Status: c.Query("status"),
	}

	items, meta, err := h.useCase.ListWithdrawals(c.Request.Context(), vendorID, params)
	if err != nil {
		HandleUsecaseError(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "vendor withdrawals retrieved successfully", items, meta)
}
