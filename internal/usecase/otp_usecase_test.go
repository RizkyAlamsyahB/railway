package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

func setupOTPUseCase(t *testing.T) (*mocks.MockOTPRepository, *mocks.MockEmailProvider, *otpUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockOTPRepository(ctrl)
	emailProvider := mocks.NewMockEmailProvider(ctrl)
	uc, err := NewOTPUseCase(repo, emailProvider, config.OTPConfig{
		CodeLength:            6,
		ExpiryMinutes:         5,
		ResendCooldownSeconds: 60,
		MaxAttempts:           5,
		ProofExpiryMinutes:    10,
		Issuer:                "test-otp-issuer",
	}, "test-otp-secret")
	if err != nil {
		t.Fatalf("NewOTPUseCase() unexpected error: %v", err)
	}
	return repo, emailProvider, uc.(*otpUseCase)
}

func TestOTPRequest(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo, emailProvider, uc := setupOTPUseCase(t)

		repo.EXPECT().FindLatestByEmailPurpose(ctx, "user@example.com", domain.OTPPurposeEmailVerification, domain.OTPChannelEmail).Return(nil, nil)
		repo.EXPECT().InvalidateActiveByEmailPurpose(ctx, "user@example.com", domain.OTPPurposeEmailVerification, domain.OTPChannelEmail).Return(nil)
		repo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
		emailProvider.EXPECT().Send(ctx, gomock.Any()).Return(nil)

		result, err := uc.RequestOTP(ctx, domain.RequestOTPInput{
			Email:   "User@Example.com",
			Purpose: domain.OTPPurposeEmailVerification,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.ExpiresAt.IsZero() || result.CooldownUntil.IsZero() {
			t.Fatal("expected expiry and cooldown timestamps")
		}
	})

	t.Run("cooldown", func(t *testing.T) {
		repo, _, uc := setupOTPUseCase(t)
		repo.EXPECT().FindLatestByEmailPurpose(ctx, "user@example.com", domain.OTPPurposeForgotPassword, domain.OTPChannelEmail).
			Return(&domain.OTPCode{CreatedAt: time.Now().Add(-10 * time.Second)}, nil)

		_, err := uc.RequestOTP(ctx, domain.RequestOTPInput{
			Email:   "user@example.com",
			Purpose: domain.OTPPurposeForgotPassword,
		})
		if !errors.Is(err, ErrOTPResendTooSoon) {
			t.Fatalf("expected ErrOTPResendTooSoon, got %v", err)
		}
	})

	t.Run("invalid purpose", func(t *testing.T) {
		repo, _, uc := setupOTPUseCase(t)

		_, err := uc.RequestOTP(ctx, domain.RequestOTPInput{
			Email:   "user@example.com",
			Purpose: "unknown",
		})
		if !errors.Is(err, ErrOTPPurposeInvalid) {
			t.Fatalf("expected ErrOTPPurposeInvalid, got %v", err)
		}

		// Ensure no repository calls were needed.
		_ = repo
	})
}

func TestOTPVerify(t *testing.T) {
	ctx := context.Background()

	t.Run("success issues proof token", func(t *testing.T) {
		repo, _, uc := setupOTPUseCase(t)
		otpID := uuid.New()
		code := "123456"
		email := "user@example.com"
		purpose := domain.OTPPurposeForgotPassword
		hash := uc.hashOTPCode(email, purpose, domain.OTPChannelEmail, code)

		repo.EXPECT().FindActiveByEmailPurpose(ctx, email, purpose, domain.OTPChannelEmail).Return(&domain.OTPCode{
			ID:           otpID,
			Email:        email,
			Purpose:      purpose,
			Channel:      domain.OTPChannelEmail,
			CodeHash:     hash,
			ExpiresAt:    time.Now().Add(5 * time.Minute),
			AttemptCount: 0,
		}, nil)
		repo.EXPECT().Consume(ctx, otpID).Return(nil)

		result, err := uc.VerifyOTP(ctx, domain.VerifyOTPInput{
			Email:   email,
			Purpose: purpose,
			Code:    code,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ProofToken == "" {
			t.Fatal("expected proof token")
		}

		claims, err := uc.VerifyProofToken(ctx, result.ProofToken, purpose)
		if err != nil {
			t.Fatalf("expected valid proof token, got %v", err)
		}
		if claims.Email != email {
			t.Fatalf("expected email %s, got %s", email, claims.Email)
		}
	})

	t.Run("invalid code increments attempts", func(t *testing.T) {
		repo, _, uc := setupOTPUseCase(t)
		otpID := uuid.New()
		email := "user@example.com"
		purpose := domain.OTPPurposeEmailVerification

		repo.EXPECT().FindActiveByEmailPurpose(ctx, email, purpose, domain.OTPChannelEmail).Return(&domain.OTPCode{
			ID:           otpID,
			Email:        email,
			Purpose:      purpose,
			Channel:      domain.OTPChannelEmail,
			CodeHash:     uc.hashOTPCode(email, purpose, domain.OTPChannelEmail, "654321"),
			ExpiresAt:    time.Now().Add(5 * time.Minute),
			AttemptCount: 0,
		}, nil)
		repo.EXPECT().IncrementAttempts(ctx, otpID).Return(nil)

		_, err := uc.VerifyOTP(ctx, domain.VerifyOTPInput{
			Email:   email,
			Purpose: purpose,
			Code:    "123456",
		})
		if !errors.Is(err, ErrOTPInvalid) {
			t.Fatalf("expected ErrOTPInvalid, got %v", err)
		}
	})

	t.Run("expired otp", func(t *testing.T) {
		repo, _, uc := setupOTPUseCase(t)
		email := "user@example.com"
		purpose := domain.OTPPurposeForgotPassword

		repo.EXPECT().FindActiveByEmailPurpose(ctx, email, purpose, domain.OTPChannelEmail).Return(&domain.OTPCode{
			ID:        uuid.New(),
			Email:     email,
			Purpose:   purpose,
			Channel:   domain.OTPChannelEmail,
			CodeHash:  uc.hashOTPCode(email, purpose, domain.OTPChannelEmail, "123456"),
			ExpiresAt: time.Now().Add(-1 * time.Minute),
		}, nil)
		repo.EXPECT().InvalidateActiveByEmailPurpose(ctx, email, purpose, domain.OTPChannelEmail).Return(nil)

		_, err := uc.VerifyOTP(ctx, domain.VerifyOTPInput{
			Email:   email,
			Purpose: purpose,
			Code:    "123456",
		})
		if !errors.Is(err, ErrOTPExpired) {
			t.Fatalf("expected ErrOTPExpired, got %v", err)
		}
	})
}

func TestOTPVerifyProofToken(t *testing.T) {
	ctx := context.Background()
	_, _, uc := setupOTPUseCase(t)

	token, _, err := uc.issueProofToken("user@example.com", domain.OTPPurposeForgotPassword)
	if err != nil {
		t.Fatalf("issueProofToken() unexpected error: %v", err)
	}

	if _, err := uc.VerifyProofToken(ctx, token, domain.OTPPurposeEmailVerification); !errors.Is(err, ErrOTPProofInvalid) {
		t.Fatalf("expected ErrOTPProofInvalid for mismatched purpose, got %v", err)
	}
}
