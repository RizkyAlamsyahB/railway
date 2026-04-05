package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	VendorOnboardingStatusOTPVerified = "otp_verified"
	VendorOnboardingStatusCompleted   = "completed"
)

// VendorOnboarding stores draft vendor registration data between onboarding steps.
type VendorOnboarding struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	Status        string     `json:"status"`
	PasswordHash  *string    `json:"-"`
	OTPVerifiedAt time.Time  `json:"otp_verified_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type VendorRegisterOTPRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

type VendorRegisterOTPResponse struct {
	ExpiresAt     time.Time `json:"expires_at"`
	CooldownUntil time.Time `json:"cooldown_until"`
}

type VendorVerifyRegistrationOTPRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
	Code  string `json:"code" binding:"required"`
}

type VendorVerifyRegistrationOTPResponse struct {
	EmailToken string `json:"email_token"`
}

type VendorRegistrationPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type VendorOnboardingStatusResponse struct {
	Status string `json:"status"`
	Email  string `json:"email"`
}

// VendorRegistrationFinalizeInput bundles all data needed to atomically
// complete vendor onboarding: create user, vendor, balance, and mark the
// onboarding record as completed.
type VendorRegistrationFinalizeInput struct {
	User         *User
	UserRole     string
	Vendor       *Vendor
	OnboardingID uuid.UUID
	CompletedAt  time.Time
}

// VendorOnboardingRepository defines persistence for vendor onboarding drafts.
type VendorOnboardingRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*VendorOnboarding, error)
	FindByEmail(ctx context.Context, email string) (*VendorOnboarding, error)
	Upsert(ctx context.Context, onboarding *VendorOnboarding) error
	Update(ctx context.Context, onboarding *VendorOnboarding, fields ...string) error
	MarkCompleted(ctx context.Context, id uuid.UUID, completedAt time.Time) error
	// FinalizeRegistration atomically creates the user, vendor profile, documents,
	// vendor balance, and marks the onboarding as completed — all in one DB transaction.
	FinalizeRegistration(ctx context.Context, input VendorRegistrationFinalizeInput) error
}
