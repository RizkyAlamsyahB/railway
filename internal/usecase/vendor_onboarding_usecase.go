package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

const vendorOnboardingDocumentDocType = "owner_document_id"

func (uc *vendorUseCase) RequestRegistrationOTP(ctx context.Context, req domain.VendorRegisterOTPRequest) (*domain.VendorRegisterOTPResponse, error) {
	email := normalizeVendorEmail(req.Email)
	if err := uc.ensureVendorRegistrationEmailAvailable(ctx, email); err != nil {
		return nil, err
	}

	result, err := uc.otpUseCase.RequestOTP(ctx, domain.RequestOTPInput{
		Email:   email,
		Purpose: domain.OTPPurposeEmailVerification,
	})
	if err != nil {
		return nil, err
	}

	return &domain.VendorRegisterOTPResponse{
		ExpiresAt:     result.ExpiresAt,
		CooldownUntil: result.CooldownUntil,
	}, nil
}

func (uc *vendorUseCase) VerifyRegistrationOTP(ctx context.Context, req domain.VendorVerifyRegistrationOTPRequest) (*domain.VendorVerifyRegistrationOTPResponse, error) {
	email := normalizeVendorEmail(req.Email)
	if err := uc.ensureVendorRegistrationEmailAvailable(ctx, email); err != nil {
		return nil, err
	}

	if _, err := uc.otpUseCase.VerifyOTP(ctx, domain.VerifyOTPInput{
		Email:   email,
		Purpose: domain.OTPPurposeEmailVerification,
		Code:    req.Code,
	}); err != nil {
		return nil, err
	}

	now := time.Now()
	onboarding, err := uc.onboardingRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor onboarding: %w", err)
	}
	if onboarding == nil {
		onboarding = &domain.VendorOnboarding{
			ID:        uuid.New(),
			Email:     email,
			CreatedAt: now,
		}
	}

	onboarding.Status = domain.VendorOnboardingStatusOTPVerified
	onboarding.PasswordHash = nil
	onboarding.StoreName = nil
	onboarding.VendorType = nil
	onboarding.BusinessLegalType = nil
	onboarding.DocumentIDType = nil
	onboarding.NIK = nil
	onboarding.OwnerName = nil
	onboarding.BirthDate = nil
	onboarding.DocumentIDObjectKey = nil
	onboarding.OTPVerifiedAt = now
	onboarding.CompletedAt = nil
	onboarding.UpdatedAt = now

	if err := uc.onboardingRepo.Upsert(ctx, onboarding); err != nil {
		return nil, fmt.Errorf("failed to store vendor onboarding: %w", err)
	}

	token, expiresAt, err := auth.GenerateVendorOnboardingToken(
		onboarding.ID,
		email,
		uc.jwtSecret,
		time.Duration(uc.jwtExpiry)*time.Hour,
		uc.jwtIssuer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate vendor onboarding token: %w", err)
	}

	return &domain.VendorVerifyRegistrationOTPResponse{
		OnboardingToken:     token,
		OnboardingExpiresAt: expiresAt,
		Status:              onboarding.Status,
	}, nil
}

func (uc *vendorUseCase) SetRegistrationPassword(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPasswordRequest) (*domain.VendorOnboardingProgressResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusOTPVerified {
		return nil, ErrVendorOnboardingStep
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	onboarding.PasswordHash = &passwordHash
	onboarding.Status = domain.VendorOnboardingStatusPasswordSet
	onboarding.UpdatedAt = now

	if err := uc.onboardingRepo.Update(ctx, onboarding, "password_hash", "status", "updated_at"); err != nil {
		return nil, fmt.Errorf("failed to update vendor onboarding password: %w", err)
	}

	return &domain.VendorOnboardingProgressResponse{Status: onboarding.Status}, nil
}

func (uc *vendorUseCase) SaveRegistrationStore(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationStoreRequest) (*domain.VendorOnboardingProgressResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusPasswordSet {
		return nil, ErrVendorOnboardingStep
	}

	now := time.Now()
	storeName := strings.TrimSpace(req.StoreName)
	vendorType := strings.TrimSpace(req.VendorType)
	onboarding.StoreName = &storeName
	onboarding.VendorType = &vendorType
	onboarding.Status = domain.VendorOnboardingStatusStoreInfoComplete
	onboarding.UpdatedAt = now

	if err := uc.onboardingRepo.Update(ctx, onboarding, "store_name", "vendor_type", "status", "updated_at"); err != nil {
		return nil, fmt.Errorf("failed to update vendor onboarding store info: %w", err)
	}

	return &domain.VendorOnboardingProgressResponse{Status: onboarding.Status}, nil
}

func (uc *vendorUseCase) PresignRegistrationLegalDocument(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPresignDocumentRequest) (*domain.VendorRegistrationPresignDocumentResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusStoreInfoComplete {
		return nil, ErrVendorOnboardingStep
	}

	objectKey := buildVendorOnboardingDocumentObjectKey(onboardingID)
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, "", PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL for vendor onboarding document: %w", err)
	}

	return &domain.VendorRegistrationPresignDocumentResponse{
		UploadURL:       uploadURL,
		ObjectKey:       objectKey,
		ExpiresAt:       time.Now().Add(PresignedUploadExpiry),
		DocumentIDType:  req.DocumentIDType,
		DocumentDocType: vendorOnboardingDocumentDocType,
	}, nil
}

func (uc *vendorUseCase) SubmitRegistrationLegal(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationLegalRequest) (*domain.VendorLoginResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusStoreInfoComplete || onboarding.PasswordHash == nil || onboarding.StoreName == nil || onboarding.VendorType == nil {
		return nil, ErrVendorOnboardingStep
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, ErrInvalidBirthDate
	}

	expectedPrefix := vendorOnboardingDocumentPrefix(onboardingID)
	if !strings.HasPrefix(req.DocumentIDObjectKey, expectedPrefix) {
		return nil, ErrInvalidDocumentObjectKey
	}

	info, err := uc.storage.HeadObject(ctx, req.DocumentIDObjectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify onboarding document object: %w", err)
	}
	if info == nil {
		return nil, ErrObjectNotUploaded
	}
	contentType := normalizeContentType(info.ContentType)
	if !isAllowedDocumentContentType(contentType) {
		return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidDocumentContent, vendorOnboardingDocumentDocType, info.ContentType)
	}
	if info.ContentLength > int64(math.MaxInt32) {
		return nil, fmt.Errorf("%w: %s (%d bytes)", ErrDocumentSizeOverflow, vendorOnboardingDocumentDocType, info.ContentLength)
	}

	if err := uc.ensureVendorRegistrationEmailAvailable(ctx, onboarding.Email); err != nil {
		return nil, err
	}

	now := time.Now()
	businessLegalType := req.BusinessLegalType
	documentIDType := req.DocumentIDType
	nik := strings.TrimSpace(req.NIK)
	ownerName := strings.TrimSpace(req.OwnerName)
	objectKey := req.DocumentIDObjectKey

	onboarding.BusinessLegalType = &businessLegalType
	onboarding.DocumentIDType = &documentIDType
	onboarding.NIK = &nik
	onboarding.OwnerName = &ownerName
	onboarding.BirthDate = &birthDate
	onboarding.DocumentIDObjectKey = &objectKey
	onboarding.UpdatedAt = now
	if err := uc.onboardingRepo.Update(
		ctx,
		onboarding,
		"business_legal_type",
		"document_id_type",
		"nik",
		"owner_name",
		"birth_date",
		"document_id_object_key",
		"updated_at",
	); err != nil {
		return nil, fmt.Errorf("failed to update vendor onboarding legal info: %w", err)
	}

	userID := uuid.New()
	emailVerifiedAt := onboarding.OTPVerifiedAt
	user := &domain.User{
		ID:              userID,
		Email:           onboarding.Email,
		FullName:        ownerName,
		BirthDate:       &birthDate,
		PasswordHash:    *onboarding.PasswordHash,
		Status:          domain.UserStatusActive,
		EmailVerifiedAt: &emailVerifiedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := uc.userRepo.Create(ctx, user, domain.RoleUMKM); err != nil {
		return nil, fmt.Errorf("failed to create vendor user: %w", err)
	}

	vendorID := uuid.New()
	vendor := &domain.Vendor{
		ID:                    vendorID,
		OwnerUserID:           userID,
		VendorType:            *onboarding.VendorType,
		DisplayName:           *onboarding.StoreName,
		ResponsiblePersonName: ownerName,
		Status:                domain.VendorStatusDraft,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	fileSize := int(info.ContentLength)
	document := domain.VendorDocument{
		ID:                 uuid.New(),
		VendorID:           vendorID,
		DocType:            vendorOnboardingDocumentDocType,
		FileURL:            objectKey,
		MimeType:           &contentType,
		FileSizeBytes:      &fileSize,
		UploadedBy:         &userID,
		VerificationStatus: domain.VerificationStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := uc.vendorRepo.CreateMinimal(ctx, vendor, []domain.VendorDocument{document}); err != nil {
		return nil, fmt.Errorf("failed to create vendor profile: %w", err)
	}
	if err := uc.vendorRepo.InitBalance(ctx, vendorID); err != nil {
		return nil, fmt.Errorf("failed to initialize vendor balance: %w", err)
	}
	if err := uc.onboardingRepo.MarkCompleted(ctx, onboarding.ID, now); err != nil {
		return nil, fmt.Errorf("failed to finalize vendor onboarding: %w", err)
	}

	token, err := auth.GenerateToken(userID, onboarding.Email, domain.RoleUMKM, &vendorID, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate vendor auth token: %w", err)
	}

	return &domain.VendorLoginResponse{
		Token:        token,
		VendorID:     vendorID,
		ImageURL:     "",
		Email:        onboarding.Email,
		Name:         ownerName,
		VendorType:   *onboarding.VendorType,
		VendorStatus: domain.VendorStatusDraft,
		DisplayName:  *onboarding.StoreName,
	}, nil
}

func (uc *vendorUseCase) ensureVendorRegistrationEmailAvailable(ctx context.Context, email string) error {
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existing == nil {
		return nil
	}

	existingVendor, err := uc.vendorRepo.FindByOwnerUserID(ctx, existing.ID)
	if err != nil {
		return fmt.Errorf("failed to check existing vendor: %w", err)
	}
	if existingVendor != nil {
		return ErrVendorAlreadyExists
	}

	return ErrEmailAlreadyRegistered
}

func (uc *vendorUseCase) getActiveOnboarding(ctx context.Context, onboardingID uuid.UUID) (*domain.VendorOnboarding, error) {
	onboarding, err := uc.onboardingRepo.FindByID(ctx, onboardingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor onboarding: %w", err)
	}
	if onboarding == nil {
		return nil, ErrVendorOnboardingNotFound
	}
	if onboarding.Status == domain.VendorOnboardingStatusCompleted || onboarding.CompletedAt != nil {
		return nil, ErrVendorOnboardingDone
	}
	return onboarding, nil
}

func normalizeVendorEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func buildVendorOnboardingDocumentObjectKey(onboardingID uuid.UUID) string {
	return fmt.Sprintf("%s%s", vendorOnboardingDocumentPrefix(onboardingID), uuid.New().String())
}

func vendorOnboardingDocumentPrefix(onboardingID uuid.UUID) string {
	return fmt.Sprintf("vendor-onboardings/%s/documents/%s/", onboardingID.String(), vendorOnboardingDocumentDocType)
}
