package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	VendorOnboardingStatusOTPVerified       = "otp_verified"
	VendorOnboardingStatusPasswordSet       = "password_set"
	VendorOnboardingStatusStoreInfoComplete = "store_info_completed"
	VendorOnboardingStatusCompleted         = "completed"
)

const (
	VendorBusinessLegalTypeIndividual = "perorangan"
	VendorBusinessLegalTypeCorporate  = "korporasi"
)

const (
	VendorDocumentIDTypeKTP      = "ktp"
	VendorDocumentIDTypePassport = "passport"
)

// VendorOnboarding stores draft vendor registration data between onboarding steps.
type VendorOnboarding struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	Status              string     `json:"status"`
	PasswordHash        *string    `json:"-"`
	StoreName           *string    `json:"store_name,omitempty"`
	VendorType          *string    `json:"vendor_type,omitempty"`
	BusinessLegalType   *string    `json:"business_legal_type,omitempty"`
	DocumentIDType      *string    `json:"document_id_type,omitempty"`
	NIK                 *string    `json:"nik,omitempty"`
	OwnerName           *string    `json:"owner_name,omitempty"`
	BirthDate           *time.Time `json:"birth_date,omitempty"`
	DocumentIDObjectKey *string    `json:"document_id_object_key,omitempty"`
	OTPVerifiedAt       time.Time  `json:"otp_verified_at"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
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
	OnboardingToken     string    `json:"onboarding_token"`
	OnboardingExpiresAt time.Time `json:"onboarding_expires_at"`
	Status              string    `json:"status"`
}

type VendorRegistrationPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type VendorRegistrationStoreRequest struct {
	StoreName  string `json:"store_name" binding:"required,max=120"`
	VendorType string `json:"vendor_type" binding:"required,oneof=souvenir_store ppiu hajj_dormitory"`
}

type VendorRegistrationPresignDocumentRequest struct {
	DocumentIDType string `json:"document_id_type" binding:"required,oneof=ktp passport"`
}

type VendorRegistrationPresignDocumentResponse struct {
	UploadURL       string    `json:"upload_url"`
	ObjectKey       string    `json:"object_key"`
	ExpiresAt       time.Time `json:"expires_at"`
	DocumentIDType  string    `json:"document_id_type"`
	DocumentDocType string    `json:"doc_type"`
}

type VendorRegistrationLegalRequest struct {
	BusinessLegalType   string `json:"business_legal_type" binding:"required,oneof=perorangan korporasi"`
	DocumentIDType      string `json:"document_id_type" binding:"required,oneof=ktp passport"`
	NIK                 string `json:"nik" binding:"required,max=32"`
	OwnerName           string `json:"owner_name" binding:"required,max=120"`
	BirthDate           string `json:"birth_date" binding:"required"`
	DocumentIDObjectKey string `json:"document_id_object_key" binding:"required"`
}

type VendorOnboardingProgressResponse struct {
	Status string `json:"status"`
}

// VendorOnboardingRepository defines persistence for vendor onboarding drafts.
type VendorOnboardingRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*VendorOnboarding, error)
	FindByEmail(ctx context.Context, email string) (*VendorOnboarding, error)
	Upsert(ctx context.Context, onboarding *VendorOnboarding) error
	Update(ctx context.Context, onboarding *VendorOnboarding, fields ...string) error
	MarkCompleted(ctx context.Context, id uuid.UUID, completedAt time.Time) error
}
