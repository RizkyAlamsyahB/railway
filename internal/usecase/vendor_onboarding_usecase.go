package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

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
	onboarding.OTPVerifiedAt = now
	onboarding.CompletedAt = nil
	onboarding.UpdatedAt = now

	if err := uc.onboardingRepo.Upsert(ctx, onboarding); err != nil {
		return nil, fmt.Errorf("failed to store vendor onboarding: %w", err)
	}

	token, err := auth.GenerateVendorOnboardingToken(
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
		EmailToken: token,
	}, nil
}

func (uc *vendorUseCase) SetRegistrationPassword(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPasswordRequest) (*domain.VendorLoginResponse, error) {
	onboarding, err := uc.getActiveOnboarding(ctx, onboardingID)
	if err != nil {
		return nil, err
	}
	if onboarding.Status != domain.VendorOnboardingStatusOTPVerified {
		return nil, ErrVendorOnboardingStep
	}
	if err := uc.ensureVendorRegistrationEmailAvailable(ctx, onboarding.Email); err != nil {
		return nil, err
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	userID := uuid.New()
	vendorID := uuid.New()

	user := &domain.User{
		ID:              userID,
		Email:           onboarding.Email,
		FullName:        onboarding.Email,
		PasswordHash:    passwordHash,
		Status:          domain.UserStatusActive,
		EmailVerifiedAt: &onboarding.OTPVerifiedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	vendor := &domain.Vendor{
		ID:          vendorID,
		OwnerUserID: userID,
		VendorType:  nil,
		DisplayName: nil,
		Status:      domain.VendorStatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.onboardingRepo.FinalizeRegistration(ctx, domain.VendorRegistrationFinalizeInput{
		User:         user,
		UserRole:     domain.RoleUMKM,
		Vendor:       vendor,
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
		AccessToken:  token,
		VendorID:     vendorID,
		Email:        onboarding.Email,
		ImageURL:     nil,
		StoreName:    vendor.DisplayName,
		VendorType:   vendor.VendorType,
		VendorStatus: vendor.Status,
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
