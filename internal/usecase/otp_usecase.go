package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type otpUseCase struct {
	repo          domain.OTPRepository
	emailProvider domain.EmailProvider
	cfg           config.OTPConfig
	signingSecret string
}

type otpProofClaims struct {
	Purpose string `json:"purpose"`
	Channel string `json:"channel"`
	Type    string `json:"typ"`
	jwt.RegisteredClaims
}

// NewOTPUseCase creates a reusable OTP service.
func NewOTPUseCase(repo domain.OTPRepository, emailProvider domain.EmailProvider, cfg config.OTPConfig, signingSecret string) (domain.OTPUseCase, error) {
	if cfg.CodeLength <= 0 {
		return nil, fmt.Errorf("OTP_CODE_LENGTH must be greater than zero")
	}
	if cfg.ExpiryMinutes <= 0 {
		return nil, fmt.Errorf("OTP_EXPIRY_MINUTES must be greater than zero")
	}
	if cfg.ResendCooldownSeconds < 0 {
		return nil, fmt.Errorf("OTP_RESEND_COOLDOWN_SECONDS must be zero or greater")
	}
	if cfg.MaxAttempts <= 0 {
		return nil, fmt.Errorf("OTP_MAX_ATTEMPTS must be greater than zero")
	}
	if cfg.ProofExpiryMinutes <= 0 {
		return nil, fmt.Errorf("OTP_PROOF_EXPIRY_MINUTES must be greater than zero")
	}
	if strings.TrimSpace(signingSecret) == "" {
		return nil, ErrOTPSecretNotConfigured
	}
	return &otpUseCase{
		repo:          repo,
		emailProvider: emailProvider,
		cfg:           cfg,
		signingSecret: signingSecret,
	}, nil
}

func (uc *otpUseCase) RequestOTP(ctx context.Context, input domain.RequestOTPInput) (*domain.RequestOTPResult, error) {
	email := normalizeEmail(input.Email)
	if err := validateOTPPurpose(input.Purpose); err != nil {
		return nil, err
	}

	now := time.Now()
	channel := domain.OTPChannelEmail
	lastOTP, err := uc.repo.FindLatestByEmailPurpose(ctx, email, input.Purpose, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to check otp cooldown: %w", err)
	}

	cooldownUntil := now.Add(time.Duration(uc.cfg.ResendCooldownSeconds) * time.Second)
	if lastOTP != nil {
		lastCooldownUntil := lastOTP.CreatedAt.Add(time.Duration(uc.cfg.ResendCooldownSeconds) * time.Second)
		if now.Before(lastCooldownUntil) {
			return nil, ErrOTPResendTooSoon
		}
	}

	if err := uc.repo.InvalidateActiveByEmailPurpose(ctx, email, input.Purpose, channel); err != nil {
		return nil, fmt.Errorf("failed to invalidate active otp: %w", err)
	}

	code, err := uc.generateNumericCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate otp code: %w", err)
	}
	codeHash := uc.hashOTPCode(email, input.Purpose, channel, code)
	expiresAt := now.Add(time.Duration(uc.cfg.ExpiryMinutes) * time.Minute)
	otp := &domain.OTPCode{
		ID:           uuid.New(),
		UserID:       input.UserID,
		Email:        email,
		Purpose:      input.Purpose,
		Channel:      channel,
		CodeHash:     codeHash,
		ExpiresAt:    expiresAt,
		AttemptCount: 0,
		CreatedAt:    now,
	}

	if err := uc.repo.Create(ctx, otp); err != nil {
		return nil, fmt.Errorf("failed to store otp: %w", err)
	}

	msg := domain.EmailMessage{
		To:      []string{email},
		Subject: buildOTPEmailSubject(input.Purpose),
		Body:    buildOTPEmailBody(code, input.Purpose, uc.cfg.ExpiryMinutes),
		IsHTML:  true,
	}
	if err := uc.emailProvider.Send(ctx, msg); err != nil {
		_ = uc.repo.InvalidateActiveByEmailPurpose(ctx, email, input.Purpose, channel)
		return nil, fmt.Errorf("failed to send otp email: %w", err)
	}

	return &domain.RequestOTPResult{
		ExpiresAt:     expiresAt,
		CooldownUntil: cooldownUntil,
	}, nil
}

func (uc *otpUseCase) VerifyOTP(ctx context.Context, input domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) {
	email := normalizeEmail(input.Email)
	if err := validateOTPPurpose(input.Purpose); err != nil {
		return nil, err
	}

	otp, err := uc.repo.FindActiveByEmailPurpose(ctx, email, input.Purpose, domain.OTPChannelEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to find otp: %w", err)
	}
	if otp == nil {
		return nil, ErrOTPInvalid
	}

	now := time.Now()
	if now.After(otp.ExpiresAt) {
		if err := uc.repo.InvalidateActiveByEmailPurpose(ctx, email, input.Purpose, domain.OTPChannelEmail); err != nil {
			return nil, fmt.Errorf("failed to invalidate expired otp: %w", err)
		}
		return nil, ErrOTPExpired
	}

	if otp.AttemptCount >= uc.cfg.MaxAttempts {
		if err := uc.repo.InvalidateActiveByEmailPurpose(ctx, email, input.Purpose, domain.OTPChannelEmail); err != nil {
			return nil, fmt.Errorf("failed to invalidate exhausted otp: %w", err)
		}
		return nil, ErrOTPTooManyAttempts
	}

	expectedHash := uc.hashOTPCode(email, input.Purpose, domain.OTPChannelEmail, strings.TrimSpace(input.Code))
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(otp.CodeHash)) != 1 {
		if err := uc.repo.IncrementAttempts(ctx, otp.ID); err != nil {
			return nil, fmt.Errorf("failed to increment otp attempts: %w", err)
		}
		if otp.AttemptCount+1 >= uc.cfg.MaxAttempts {
			if err := uc.repo.InvalidateActiveByEmailPurpose(ctx, email, input.Purpose, domain.OTPChannelEmail); err != nil {
				return nil, fmt.Errorf("failed to invalidate otp after max attempts: %w", err)
			}
			return nil, ErrOTPTooManyAttempts
		}
		return nil, ErrOTPInvalid
	}

	if err := uc.repo.Consume(ctx, otp.ID); err != nil {
		return nil, fmt.Errorf("failed to consume otp: %w", err)
	}

	proofToken, proofExpiresAt, err := uc.issueProofToken(email, input.Purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to issue otp proof token: %w", err)
	}

	return &domain.VerifyOTPResult{
		ProofToken:     proofToken,
		ProofExpiresAt: proofExpiresAt,
		Email:          email,
		Purpose:        input.Purpose,
	}, nil
}

func (uc *otpUseCase) VerifyProofToken(_ context.Context, tokenString string, expectedPurpose string) (*domain.OTPProofClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &otpProofClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.signingSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrOTPProofExpired
		}
		return nil, ErrOTPProofInvalid
	}

	claims, ok := token.Claims.(*otpProofClaims)
	if !ok || !token.Valid {
		return nil, ErrOTPProofInvalid
	}
	if claims.Type != "otp_proof" {
		return nil, ErrOTPProofInvalid
	}
	if claims.Purpose != expectedPurpose {
		return nil, ErrOTPProofInvalid
	}
	if claims.Channel != domain.OTPChannelEmail {
		return nil, ErrOTPProofInvalid
	}
	if claims.Subject == "" || claims.IssuedAt == nil || claims.ExpiresAt == nil {
		return nil, ErrOTPProofInvalid
	}

	return &domain.OTPProofClaims{
		Email:      claims.Subject,
		Purpose:    claims.Purpose,
		Channel:    claims.Channel,
		VerifiedAt: claims.IssuedAt.Time,
		ExpiresAt:  claims.ExpiresAt.Time,
	}, nil
}

func (uc *otpUseCase) generateNumericCode() (string, error) {
	limit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(uc.cfg.CodeLength)), nil)
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", uc.cfg.CodeLength, n.Int64()), nil
}

func (uc *otpUseCase) hashOTPCode(email, purpose, channel, code string) string {
	mac := hmac.New(sha256.New, []byte(uc.signingSecret))
	mac.Write([]byte(email))
	mac.Write([]byte{':'})
	mac.Write([]byte(purpose))
	mac.Write([]byte{':'})
	mac.Write([]byte(channel))
	mac.Write([]byte{':'})
	mac.Write([]byte(code))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (uc *otpUseCase) issueProofToken(email, purpose string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(uc.cfg.ProofExpiryMinutes) * time.Minute)
	claims := otpProofClaims{
		Purpose: purpose,
		Channel: domain.OTPChannelEmail,
		Type:    "otp_proof",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			Issuer:    uc.cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(uc.signingSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateOTPPurpose(purpose string) error {
	switch purpose {
	case domain.OTPPurposeEmailVerification, domain.OTPPurposeForgotPassword:
		return nil
	default:
		return ErrOTPPurposeInvalid
	}
}

func buildOTPEmailSubject(purpose string) string {
	switch purpose {
	case domain.OTPPurposeForgotPassword:
		return "Kode OTP Reset Password"
	default:
		return "Kode OTP Verifikasi"
	}
}

func buildOTPEmailBody(code string, purpose string, expiryMinutes int) string {
	title := "Verifikasi OTP"
	description := "Gunakan kode OTP berikut untuk melanjutkan proses Anda."
	if purpose == domain.OTPPurposeForgotPassword {
		title = "Reset Password"
		description = "Gunakan kode OTP berikut untuk melanjutkan proses reset password Anda."
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0;padding:0;font-family:Arial,Helvetica,sans-serif;background-color:#f4f4f4;">
    <div style="max-width:600px;margin:0 auto;padding:20px;">
        <div style="background-color:#ffffff;border-radius:8px;padding:40px;text-align:center;">
            <h1 style="color:#333333;margin-bottom:20px;">%s</h1>
            <p style="color:#666666;font-size:16px;line-height:1.5;margin-bottom:24px;">%s</p>
            <div style="display:inline-block;padding:16px 28px;border-radius:8px;background:#0f766e;color:#ffffff;font-size:32px;font-weight:bold;letter-spacing:8px;">%s</div>
            <p style="color:#999999;font-size:13px;margin-top:30px;">Kode ini berlaku selama %d menit. Jangan bagikan kode ini kepada siapa pun.</p>
        </div>
    </div>
</body>
</html>`, title, description, code, expiryMinutes)
}
