package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

// Sentinel errors for user use case.
var (
	ErrPhoneAlreadyRegistered   = errors.New("phone number already registered")
	ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
	ErrUserAlreadyActive        = errors.New("user is already active")
	ErrResendTooSoon            = errors.New("please wait before requesting another verification email")
	ErrUserNotPending           = errors.New("user is not in pending status")
	ErrUserInvalidCredentials   = errors.New("invalid email or password")
	ErrUserAccountNotActive     = errors.New("account is not active")
	ErrUserEmailNotVerified     = errors.New("email is not verified")
	ErrUserAccountBlocked       = errors.New("account is blocked")
)

const (
	verificationTokenExpiry  = 24 * time.Hour
	resendCooldown           = 1 * time.Minute
	verificationTokenByteLen = 32
)

type userUseCase struct {
	userRepo      domain.UserRepository
	tokenRepo     domain.EmailVerificationTokenRepository
	emailProvider domain.EmailProvider
	backendURL    string
	jwtSecret     string
	jwtExpiry     int
	jwtIssuer     string
}

// NewUserUseCase creates a new UserUseCase.
func NewUserUseCase(
	userRepo domain.UserRepository,
	tokenRepo domain.EmailVerificationTokenRepository,
	emailProvider domain.EmailProvider,
	backendURL string,
	jwtSecret string,
	jwtExpiry int,
	jwtIssuer string,
) domain.UserUseCase {
	return &userUseCase{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		emailProvider: emailProvider,
		backendURL:    backendURL,
		jwtSecret:     jwtSecret,
		jwtExpiry:     jwtExpiry,
		jwtIssuer:     jwtIssuer,
	}
}

func (uc *userUseCase) Register(ctx context.Context, req domain.RegisterCustomerRequest) (*domain.RegisterCustomerResponse, error) {
	// Check email uniqueness
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	// Check phone uniqueness
	existingPhone, err := uc.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check phone: %w", err)
	}
	if existingPhone != nil {
		return nil, ErrPhoneAlreadyRegistered
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with role customer, status pending
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		FullName:     req.FullName,
		Phone:        &req.Phone,
		PasswordHash: passwordHash,
		Status:       "pending",
	}

	if err := uc.userRepo.Create(ctx, user, "customer"); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate verification token and send email
	if err := uc.createAndSendVerificationToken(ctx, user.ID, req.Email); err != nil {
		return nil, fmt.Errorf("failed to send verification email: %w", err)
	}

	return &domain.RegisterCustomerResponse{
		UserID:  user.ID,
		Email:   req.Email,
		Message: "Registration successful. Please check your email to verify your account.",
	}, nil
}

func (uc *userUseCase) VerifyEmail(ctx context.Context, rawToken string) error {
	// Decode and hash the token
	tokenBytes, err := base64.URLEncoding.DecodeString(rawToken)
	if err != nil {
		return ErrInvalidVerificationToken
	}

	hash := sha256.Sum256(tokenBytes)
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])

	// Find active token
	token, err := uc.tokenRepo.FindActiveByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to find verification token: %w", err)
	}
	if token == nil {
		return ErrInvalidVerificationToken
	}

	// Check expiry
	if time.Now().After(token.ExpiresAt) {
		return ErrInvalidVerificationToken
	}

	// Activate user
	if err := uc.userRepo.ActivateUser(ctx, token.UserID); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	// Consume the token
	if err := uc.tokenRepo.Consume(ctx, token.ID); err != nil {
		return fmt.Errorf("failed to consume token: %w", err)
	}

	return nil
}

func (uc *userUseCase) ResendVerification(ctx context.Context, req domain.ResendVerificationRequest) error {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		// Return nil to not leak whether email exists
		return nil
	}

	// Check user is still pending
	if user.Status == "active" {
		return ErrUserAlreadyActive
	}
	if user.Status != "pending" {
		return ErrUserNotPending
	}

	// Rate limit: check last token created_at
	lastToken, err := uc.tokenRepo.FindLatestByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to check resend cooldown: %w", err)
	}
	if lastToken != nil && time.Since(lastToken.CreatedAt) < resendCooldown {
		return ErrResendTooSoon
	}

	// Invalidate existing tokens
	if err := uc.tokenRepo.InvalidateByUserID(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to invalidate old tokens: %w", err)
	}

	// Create and send new token
	if err := uc.createAndSendVerificationToken(ctx, user.ID, user.Email); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// createAndSendVerificationToken generates a token, stores its hash, and sends the verification email.
func (uc *userUseCase) createAndSendVerificationToken(ctx context.Context, userID uuid.UUID, email string) error {
	// Generate random token
	rawToken := make([]byte, verificationTokenByteLen)
	if _, err := rand.Read(rawToken); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Hash for storage
	hash := sha256.Sum256(rawToken)
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])

	// Store token
	token := &domain.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    userID,
		Email:     email,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(verificationTokenExpiry),
		CreatedAt: time.Now(),
	}

	if err := uc.tokenRepo.Create(ctx, token); err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// Build verification URL (points to backend GET endpoint)
	rawTokenB64 := base64.URLEncoding.EncodeToString(rawToken)
	verifyURL := fmt.Sprintf("%s/api/v1/users/verify-email?token=%s", uc.backendURL, rawTokenB64)

	// Send email
	emailMsg := domain.EmailMessage{
		To:      []string{email},
		Subject: "Verifikasi Email Anda",
		Body:    buildVerificationEmailBody(verifyURL),
		IsHTML:  true,
	}

	if err := uc.emailProvider.Send(ctx, emailMsg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildVerificationEmailBody creates the HTML email body for email verification.
func buildVerificationEmailBody(verifyURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0;padding:0;font-family:Arial,Helvetica,sans-serif;background-color:#f4f4f4;">
    <div style="max-width:600px;margin:0 auto;padding:20px;">
        <div style="background-color:#ffffff;border-radius:8px;padding:40px;text-align:center;">
            <h1 style="color:#333333;margin-bottom:20px;">Verifikasi Email Anda</h1>
            <p style="color:#666666;font-size:16px;line-height:1.5;margin-bottom:30px;">
                Terima kasih telah mendaftar. Silakan klik tombol di bawah untuk memverifikasi alamat email Anda.
            </p>
            <a href="%s" style="display:inline-block;background-color:#4CAF50;color:#ffffff;text-decoration:none;padding:14px 32px;border-radius:6px;font-size:16px;font-weight:bold;">
                Verifikasi Email
            </a>
            <p style="color:#999999;font-size:13px;margin-top:30px;">
                Link ini hanya berlaku selama 24 jam. Jika Anda tidak mendaftar, abaikan email ini.
            </p>
        </div>
    </div>
</body>
</html>`, verifyURL)
}

func (uc *userUseCase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// 1. Find user by email.
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserInvalidCredentials
	}

	// 2. Verify password.
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrUserInvalidCredentials
	}

	// 3. Check account is not blocked.
	if user.Status == "blocked" {
		return nil, ErrUserAccountBlocked
	}

	// 4. Check account is active.
	if user.Status != "active" {
		return nil, ErrUserAccountNotActive
	}

	// 5. Check email is verified.
	if user.EmailVerifiedAt == nil {
		return nil, ErrUserEmailNotVerified
	}

	// 6. Determine role code.
	roleCode := ""
	if user.Role != nil {
		roleCode = user.Role.Code
	}

	// 7. Generate JWT token.
	token, err := auth.GenerateToken(user.ID, user.Email, roleCode, nil, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token: token,
		User: domain.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			FullName:        user.FullName,
			BirthDate:       user.BirthDate,
			Phone:           user.Phone,
			Status:          user.Status,
			EmailVerifiedAt: user.EmailVerifiedAt,
			Role:            roleCode,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		},
	}, nil
}
