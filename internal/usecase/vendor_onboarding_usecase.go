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

const vendorOnboardingDocumentDocType = domain.VendorDocumentTypeOwnerDocumentID

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
	onboarding.NIB = nil
	onboarding.CompanyName = nil
	onboarding.EstablishedDate = nil
	onboarding.RegisteredAddress = nil
	onboarding.NIBDocumentObjectKey = nil
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

func (uc *vendorUseCase) GetRegistrationStatus(ctx context.Context, onboardingID uuid.UUID) (*domain.VendorOnboardingStatusResponse, error) {
	onboarding, err := uc.onboardingRepo.FindByID(ctx, onboardingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor onboarding: %w", err)
	}
	if onboarding == nil {
		return nil, ErrVendorOnboardingNotFound
	}

	return &domain.VendorOnboardingStatusResponse{
		Status: onboarding.Status,
		Email:  onboarding.Email,
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

func (uc *vendorUseCase) PresignRegistrationIndividualDocument(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPresignDocumentRequest) (*domain.VendorRegistrationPresignDocumentResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusPasswordSet {
		return nil, ErrVendorOnboardingStep
	}

	contentType := normalizeContentType(req.ContentType)
	if !isAllowedIndividualOnboardingContentType(contentType) {
		return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidDocumentContent, vendorOnboardingDocumentDocType, req.ContentType)
	}

	objectKey := buildVendorOnboardingDocumentObjectKey(onboardingID)
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, contentType, PresignedUploadExpiry)
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

func (uc *vendorUseCase) PresignRegistrationCorporateDocument(ctx context.Context, onboardingID uuid.UUID) (*domain.VendorRegistrationPresignDocumentResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusPasswordSet {
		return nil, ErrVendorOnboardingStep
	}

	objectKey := buildVendorOnboardingNIBDocumentObjectKey(onboardingID)
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, "application/pdf", PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL for vendor onboarding nib document: %w", err)
	}

	return &domain.VendorRegistrationPresignDocumentResponse{
		UploadURL:       uploadURL,
		ObjectKey:       objectKey,
		ExpiresAt:       time.Now().Add(PresignedUploadExpiry),
		DocumentIDType:  "",
		DocumentDocType: domain.VendorDocumentTypeBusinessNIB,
	}, nil
}

func (uc *vendorUseCase) SubmitRegistrationIndividual(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationIndividualLegalRequest) (*domain.VendorLoginResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusPasswordSet || onboarding.PasswordHash == nil {
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
	if !isAllowedIndividualOnboardingContentType(contentType) {
		return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidDocumentContent, vendorOnboardingDocumentDocType, info.ContentType)
	}
	if info.ContentLength > int64(math.MaxInt32) {
		return nil, fmt.Errorf("%w: %s (%d bytes)", ErrDocumentSizeOverflow, vendorOnboardingDocumentDocType, info.ContentLength)
	}

	if err := uc.ensureVendorRegistrationEmailAvailable(ctx, onboarding.Email); err != nil {
		return nil, err
	}

	now := time.Now()
	businessLegalType := domain.VendorBusinessLegalTypeIndividual
	vendorType := domain.VendorTypeSouvenirStore
	storeName := strings.TrimSpace(req.StoreName)
	documentIDType := req.DocumentIDType
	nik := strings.TrimSpace(req.NIK)
	ownerName := strings.TrimSpace(req.OwnerName)
	objectKey := req.DocumentIDObjectKey

	onboarding.StoreName = &storeName
	onboarding.VendorType = &vendorType
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
		"store_name",
		"vendor_type",
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
	vendorID := uuid.New()
	vendor := &domain.Vendor{
		ID:                    vendorID,
		OwnerUserID:           userID,
		VendorType:            vendorType,
		BusinessLegalType:     &businessLegalType,
		DisplayName:           storeName,
		ResponsiblePersonName: ownerName,
		Status:                domain.VendorStatusDraft,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	fileSize := int(info.ContentLength)
	document := domain.VendorDocument{
		ID:                 uuid.New(),
		VendorID:           vendorID,
		DocType:            domain.VendorDocumentTypeOwnerDocumentID,
		FileURL:            objectKey,
		MimeType:           &contentType,
		FileSizeBytes:      &fileSize,
		UploadedBy:         &userID,
		VerificationStatus: domain.VerificationStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := uc.onboardingRepo.FinalizeRegistration(ctx, domain.VendorRegistrationFinalizeInput{
		User:         user,
		UserRole:     domain.RoleUMKM,
		Vendor:       vendor,
		Documents:    []domain.VendorDocument{document},
		OnboardingID: onboarding.ID,
		CompletedAt:  now,
	}); err != nil {
		return nil, fmt.Errorf("failed to finalize vendor registration: %w", err)
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
		VendorType:   vendorType,
		VendorStatus: domain.VendorStatusDraft,
		DisplayName:  storeName,
	}, nil
}

func (uc *vendorUseCase) SubmitRegistrationCorporate(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationCorporateLegalRequest) (*domain.VendorLoginResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusPasswordSet || onboarding.PasswordHash == nil {
		return nil, ErrVendorOnboardingStep
	}

	establishedDate, err := time.Parse("2006-01-02", req.EstablishedDate)
	if err != nil {
		return nil, ErrInvalidBirthDate
	}

	expectedPrefix := vendorOnboardingNIBDocumentPrefix(onboardingID)
	if !strings.HasPrefix(req.NIBDocumentObjectKey, expectedPrefix) {
		return nil, ErrInvalidDocumentObjectKey
	}

	info, err := uc.storage.HeadObject(ctx, req.NIBDocumentObjectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify onboarding nib document object: %w", err)
	}
	if info == nil {
		return nil, ErrObjectNotUploaded
	}
	contentType := normalizeContentType(info.ContentType)
	if !isAllowedNIBDocumentContentType(contentType) {
		return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidDocumentContent, domain.VendorDocumentTypeBusinessNIB, info.ContentType)
	}
	if info.ContentLength > int64(math.MaxInt32) {
		return nil, fmt.Errorf("%w: %s (%d bytes)", ErrDocumentSizeOverflow, domain.VendorDocumentTypeBusinessNIB, info.ContentLength)
	}

	if err := uc.ensureVendorRegistrationEmailAvailable(ctx, onboarding.Email); err != nil {
		return nil, err
	}

	now := time.Now()
	businessLegalType := domain.VendorBusinessLegalTypeCorporate
	vendorType := domain.VendorTypeSouvenirStore
	storeName := strings.TrimSpace(req.StoreName)
	nib := strings.TrimSpace(req.NIB)
	companyName := strings.TrimSpace(req.CompanyName)
	registeredAddress := strings.TrimSpace(req.RegisteredAddress)
	nibObjectKey := req.NIBDocumentObjectKey

	onboarding.StoreName = &storeName
	onboarding.VendorType = &vendorType
	onboarding.BusinessLegalType = &businessLegalType
	onboarding.NIB = &nib
	onboarding.CompanyName = &companyName
	onboarding.EstablishedDate = &establishedDate
	onboarding.RegisteredAddress = &registeredAddress
	onboarding.NIBDocumentObjectKey = &nibObjectKey
	onboarding.UpdatedAt = now
	if err := uc.onboardingRepo.Update(
		ctx,
		onboarding,
		"store_name",
		"vendor_type",
		"business_legal_type",
		"nib",
		"company_name",
		"established_date",
		"registered_address",
		"nib_document_object_key",
		"updated_at",
	); err != nil {
		return nil, fmt.Errorf("failed to update vendor onboarding corporate info: %w", err)
	}

	userID := uuid.New()
	emailVerifiedAt := onboarding.OTPVerifiedAt
	user := &domain.User{
		ID:              userID,
		Email:           onboarding.Email,
		FullName:        companyName,
		BirthDate:       nil,
		PasswordHash:    *onboarding.PasswordHash,
		Status:          domain.UserStatusActive,
		EmailVerifiedAt: &emailVerifiedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	vendorID := uuid.New()
	vendor := &domain.Vendor{
		ID:                    vendorID,
		OwnerUserID:           userID,
		VendorType:            vendorType,
		BusinessLegalType:     &businessLegalType,
		DisplayName:           storeName,
		ResponsiblePersonName: companyName,
		RegisteredAddress:     &registeredAddress,
		Status:                domain.VendorStatusDraft,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	fileSize := int(info.ContentLength)
	document := domain.VendorDocument{
		ID:                 uuid.New(),
		VendorID:           vendorID,
		DocType:            domain.VendorDocumentTypeBusinessNIB,
		FileURL:            nibObjectKey,
		MimeType:           &contentType,
		FileSizeBytes:      &fileSize,
		UploadedBy:         &userID,
		VerificationStatus: domain.VerificationStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := uc.onboardingRepo.FinalizeRegistration(ctx, domain.VendorRegistrationFinalizeInput{
		User:         user,
		UserRole:     domain.RoleUMKM,
		Vendor:       vendor,
		Documents:    []domain.VendorDocument{document},
		OnboardingID: onboarding.ID,
		CompletedAt:  now,
	}); err != nil {
		return nil, fmt.Errorf("failed to finalize vendor registration: %w", err)
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
		Name:         companyName,
		VendorType:   vendorType,
		VendorStatus: domain.VendorStatusDraft,
		DisplayName:  storeName,
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

func buildVendorOnboardingNIBDocumentObjectKey(onboardingID uuid.UUID) string {
	return fmt.Sprintf("%s%s", vendorOnboardingNIBDocumentPrefix(onboardingID), uuid.New().String())
}

func vendorOnboardingNIBDocumentPrefix(onboardingID uuid.UUID) string {
	return fmt.Sprintf("vendor-onboardings/%s/documents/%s/", onboardingID.String(), domain.VendorDocumentTypeBusinessNIB)
}
