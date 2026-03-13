package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	OTPChannelEmail = "email"
)

const (
	OTPPurposeEmailVerification = "email_verification"
	OTPPurposeForgotPassword    = "forgot_password"
)

var OTPPurposes = []string{
	OTPPurposeEmailVerification,
	OTPPurposeForgotPassword,
}

// OTPCode represents a single OTP challenge.
type OTPCode struct {
	ID            uuid.UUID
	UserID        *uuid.UUID
	Email         string
	Purpose       string
	Channel       string
	CodeHash      string
	ExpiresAt     time.Time
	AttemptCount  int
	ConsumedAt    *time.Time
	InvalidatedAt *time.Time
	CreatedAt     time.Time
}

// RequestOTPInput is the input DTO for requesting a new OTP.
type RequestOTPInput struct {
	Email   string     `json:"email" binding:"required,email,max=255"`
	Purpose string     `json:"purpose" binding:"required"`
	UserID  *uuid.UUID `json:"user_id,omitempty"`
}

// RequestOTPResult is the output DTO for a successful OTP request.
type RequestOTPResult struct {
	ExpiresAt     time.Time `json:"expires_at"`
	CooldownUntil time.Time `json:"cooldown_until"`
}

// VerifyOTPInput is the input DTO for verifying an OTP code.
type VerifyOTPInput struct {
	Email   string `json:"email" binding:"required,email,max=255"`
	Purpose string `json:"purpose" binding:"required"`
	Code    string `json:"code" binding:"required"`
}

// VerifyOTPResult is the output DTO for a successful OTP verification.
type VerifyOTPResult struct {
	ProofToken     string    `json:"proof_token"`
	ProofExpiresAt time.Time `json:"proof_expires_at"`
	Email          string    `json:"email"`
	Purpose        string    `json:"purpose"`
}

// OTPProofClaims contains the verified information encoded in an OTP proof token.
type OTPProofClaims struct {
	Email      string    `json:"email"`
	Purpose    string    `json:"purpose"`
	Channel    string    `json:"channel"`
	VerifiedAt time.Time `json:"verified_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// OTPRepository defines the interface for OTP data access.
type OTPRepository interface {
	Create(ctx context.Context, otp *OTPCode) error
	FindActiveByEmailPurpose(ctx context.Context, email, purpose, channel string) (*OTPCode, error)
	FindLatestByEmailPurpose(ctx context.Context, email, purpose, channel string) (*OTPCode, error)
	InvalidateActiveByEmailPurpose(ctx context.Context, email, purpose, channel string) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
	Consume(ctx context.Context, id uuid.UUID) error
}

// OTPUseCase defines the reusable OTP business operations.
type OTPUseCase interface {
	RequestOTP(ctx context.Context, input RequestOTPInput) (*RequestOTPResult, error)
	VerifyOTP(ctx context.Context, input VerifyOTPInput) (*VerifyOTPResult, error)
	VerifyProofToken(ctx context.Context, token string, expectedPurpose string) (*OTPProofClaims, error)
}
