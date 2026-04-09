package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
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
	otpUseCase    domain.OTPUseCase
	googleOAuth   domain.GoogleOAuthVerifier
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
	otpUseCase domain.OTPUseCase,
	googleOAuth domain.GoogleOAuthVerifier,
	backendURL string,
	jwtSecret string,
	jwtExpiry int,
	jwtIssuer string,
) domain.UserUseCase {
	return &userUseCase{
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		emailProvider: emailProvider,
		otpUseCase:    otpUseCase,
		googleOAuth:   googleOAuth,
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
		Status:       domain.UserStatusPending,
	}

	if err := uc.userRepo.Create(ctx, user, domain.RoleCustomer); err != nil {
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
	if user.Status == domain.UserStatusActive {
		return ErrUserAlreadyActive
	}
	if user.Status != domain.UserStatusPending {
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
	if user.Status == domain.UserStatusBlocked {
		return nil, ErrUserAccountBlocked
	}

	// 4. Check account is not deactivated.
	if user.Status == domain.UserStatusDeactivated {
		return nil, ErrAccountAlreadyDeactivated
	}

	// 5. Check account is active.
	if user.Status != domain.UserStatusActive {
		return nil, ErrUserAccountNotActive
	}

	// 5. Check email is verified.
	if user.EmailVerifiedAt == nil {
		return nil, ErrUserEmailNotVerified
	}

	// 6. Determine role code for JWT.
	roleCode := ""
	if user.Role != nil {
		roleCode = user.Role.Code
	}

	// 7. Generate JWT token.
	token, err := auth.GenerateToken(user.ID, user.Email, roleCode, nil, user.PasswordChangedAt, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token: token,
		User:  *toUserResponse(user),
	}, nil
}

func (uc *userUseCase) LoginWithGoogle(ctx context.Context, req domain.GoogleLoginRequest) (*domain.LoginResponse, error) {
	if uc.googleOAuth == nil {
		return nil, ErrGoogleOAuthDisabled
	}

	profile, err := uc.googleOAuth.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		return nil, ErrGoogleIDTokenInvalid
	}
	if profile == nil || strings.TrimSpace(profile.Email) == "" {
		return nil, ErrGoogleIDTokenInvalid
	}
	if !profile.EmailVerified {
		return nil, ErrGoogleEmailNotVerified
	}

	user, err := uc.userRepo.FindByEmail(ctx, profile.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		passwordHash, err := auth.HashPassword(uuid.NewString())
		if err != nil {
			return nil, fmt.Errorf("failed to generate user password hash: %w", err)
		}

		now := time.Now()
		var imageURL *string
		if profile.PictureURL != "" {
			imageURL = &profile.PictureURL
		}

		fullName := profile.Name
		if fullName == "" {
			fullName = profile.Email
		}

		newUser := &domain.User{
			ID:              uuid.New(),
			Email:           profile.Email,
			FullName:        fullName,
			ImageURL:        imageURL,
			PasswordHash:    passwordHash,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &now,
			Role:            &domain.Role{Code: domain.RoleCustomer},
		}

		if err := uc.userRepo.Create(ctx, newUser, domain.RoleCustomer); err != nil {
			return nil, fmt.Errorf("failed to create oauth user: %w", err)
		}

		user, err = uc.userRepo.FindByID(ctx, newUser.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to reload oauth user: %w", err)
		}
		if user == nil {
			return nil, ErrUserNotFound
		}
	}

	if user.Role == nil || user.Role.Code != domain.RoleCustomer {
		return nil, ErrGoogleAccountNotCustomer
	}
	if user.Status == domain.UserStatusBlocked {
		return nil, ErrUserAccountBlocked
	}

	if user.Status != domain.UserStatusActive || user.EmailVerifiedAt == nil {
		if err := uc.userRepo.ActivateUser(ctx, user.ID); err != nil {
			return nil, fmt.Errorf("failed to activate oauth user: %w", err)
		}
		user, err = uc.userRepo.FindByID(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to reload oauth user: %w", err)
		}
		if user == nil {
			return nil, ErrUserNotFound
		}
	}

	token, err := auth.GenerateToken(user.ID, user.Email, domain.RoleCustomer, nil, user.PasswordChangedAt, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token: token,
		User:  *toUserResponse(user),
	}, nil
}

func (uc *userUseCase) GetMe(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return toUserResponse(user), nil
}

func (uc *userUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req domain.ChangePasswordRequest) error {
	// 1. Validate password strength.
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		return fmt.Errorf("%w: %s", ErrPasswordTooWeak, err.Error())
	}

	// 2. Find user.
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 3. Verify old password.
	if err := auth.CheckPassword(req.OldPassword, user.PasswordHash); err != nil {
		return ErrOldPasswordIncorrect
	}

	// 4. Ensure new password differs from old.
	if auth.CheckPassword(req.NewPassword, user.PasswordHash) == nil {
		return ErrNewPasswordSameAsOld
	}

	// 5. Hash and update.
	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := uc.userRepo.UpdatePasswordHash(ctx, userID, newHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (uc *userUseCase) ResetPassword(ctx context.Context, req domain.ResetPasswordRequest) error {
	// 1. Verify proof token from OTP flow.
	claims, err := uc.otpUseCase.VerifyProofToken(ctx, req.ProofToken, domain.OTPPurposeForgotPassword)
	if err != nil {
		return err
	}

	// 2. Validate password strength.
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		return fmt.Errorf("%w: %s", ErrPasswordTooWeak, err.Error())
	}

	// 3. Find user by email from proof claims.
	user, err := uc.userRepo.FindByEmail(ctx, claims.Email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 4. Hash and update.
	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := uc.userRepo.UpdatePasswordHash(ctx, user.ID, newHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (uc *userUseCase) DeleteAccount(ctx context.Context, userID uuid.UUID, req domain.DeleteAccountRequest) error {
	// 1. Find user.
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 2. Check not already deactivated.
	if user.Status == domain.UserStatusDeactivated {
		return ErrAccountAlreadyDeactivated
	}

	// 3. Verify password.
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return ErrPasswordIncorrect
	}

	// 4. Deactivate account.
	if err := uc.userRepo.Deactivate(ctx, userID, req.Reason); err != nil {
		return fmt.Errorf("failed to deactivate account: %w", err)
	}

	return nil
}
