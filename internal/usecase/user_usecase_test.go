package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
	"go.uber.org/mock/gomock"
)

// ---------- helpers ----------

func setupUserUseCase(t *testing.T) (
	*mocks.MockUserRepository,
	*mocks.MockEmailVerificationTokenRepository,
	*mocks.MockEmailProvider,
	domain.UserUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenRepo := mocks.NewMockEmailVerificationTokenRepository(ctrl)
	emailProvider := mocks.NewMockEmailProvider(ctrl)
	uc := NewUserUseCase(userRepo, tokenRepo, emailProvider, "http://localhost:8080", "test-secret", 24, "test-issuer")
	return userRepo, tokenRepo, emailProvider, uc
}

// ============================================================
// Register
// ============================================================

func TestUserRegister_Success(t *testing.T) {
	userRepo, tokenRepo, emailProvider, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.RegisterCustomerRequest{
		FullName: "John Doe",
		Phone:    "08123456789",
		Email:    "john@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	userRepo.EXPECT().FindByPhone(ctx, req.Phone).Return(nil, nil)
	userRepo.EXPECT().Create(ctx, gomock.Any(), "customer").Return(nil)
	tokenRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	emailProvider.EXPECT().Send(ctx, gomock.Any()).Return(nil)

	resp, err := uc.Register(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, resp.Email)
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestUserRegister_EmailAlreadyExists(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.RegisterCustomerRequest{
		FullName: "John Doe",
		Phone:    "08123456789",
		Email:    "existing@example.com",
		Password: "password123",
	}

	existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)

	_, err := uc.Register(ctx, req)
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Errorf("expected ErrEmailAlreadyRegistered, got %v", err)
	}
}

func TestUserRegister_PhoneAlreadyExists(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.RegisterCustomerRequest{
		FullName: "John Doe",
		Phone:    "08123456789",
		Email:    "new@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)

	existingUser := &domain.User{ID: uuid.New(), Phone: &req.Phone}
	userRepo.EXPECT().FindByPhone(ctx, req.Phone).Return(existingUser, nil)

	_, err := uc.Register(ctx, req)
	if !errors.Is(err, ErrPhoneAlreadyRegistered) {
		t.Errorf("expected ErrPhoneAlreadyRegistered, got %v", err)
	}
}

func TestUserRegister_FindByEmailError(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.RegisterCustomerRequest{
		FullName: "John Doe",
		Phone:    "08123456789",
		Email:    "new@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, errors.New("db error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserRegister_CreateUserError(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.RegisterCustomerRequest{
		FullName: "John Doe",
		Phone:    "08123456789",
		Email:    "new@example.com",
		Password: "password123",
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	userRepo.EXPECT().FindByPhone(ctx, req.Phone).Return(nil, nil)
	userRepo.EXPECT().Create(ctx, gomock.Any(), "customer").Return(errors.New("db error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// VerifyEmail
// ============================================================

func TestUserVerifyEmail_Success(t *testing.T) {
	userRepo, tokenRepo, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	// Create a valid raw token
	rawToken := make([]byte, 32)
	for i := range rawToken {
		rawToken[i] = byte(i)
	}
	rawTokenB64 := base64.URLEncoding.EncodeToString(rawToken)
	hash := sha256.Sum256(rawToken)
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])

	tokenID := uuid.New()
	userID := uuid.New()

	token := &domain.EmailVerificationToken{
		ID:        tokenID,
		UserID:    userID,
		Email:     "john@example.com",
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	tokenRepo.EXPECT().FindActiveByTokenHash(ctx, tokenHash).Return(token, nil)
	userRepo.EXPECT().ActivateUser(ctx, userID).Return(nil)
	tokenRepo.EXPECT().Consume(ctx, tokenID).Return(nil)

	err := uc.VerifyEmail(ctx, rawTokenB64)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserVerifyEmail_InvalidToken(t *testing.T) {
	_, tokenRepo, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	// Use a valid base64 but no matching token in DB
	rawToken := make([]byte, 32)
	rawTokenB64 := base64.URLEncoding.EncodeToString(rawToken)
	hash := sha256.Sum256(rawToken)
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])

	tokenRepo.EXPECT().FindActiveByTokenHash(ctx, tokenHash).Return(nil, nil)

	err := uc.VerifyEmail(ctx, rawTokenB64)
	if !errors.Is(err, ErrInvalidVerificationToken) {
		t.Errorf("expected ErrInvalidVerificationToken, got %v", err)
	}
}

func TestUserVerifyEmail_ExpiredToken(t *testing.T) {
	_, tokenRepo, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	rawToken := make([]byte, 32)
	for i := range rawToken {
		rawToken[i] = byte(i + 1)
	}
	rawTokenB64 := base64.URLEncoding.EncodeToString(rawToken)
	hash := sha256.Sum256(rawToken)
	tokenHash := base64.URLEncoding.EncodeToString(hash[:])

	token := &domain.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Email:     "john@example.com",
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // expired
		CreatedAt: time.Now().Add(-25 * time.Hour),
	}

	tokenRepo.EXPECT().FindActiveByTokenHash(ctx, tokenHash).Return(token, nil)

	err := uc.VerifyEmail(ctx, rawTokenB64)
	if !errors.Is(err, ErrInvalidVerificationToken) {
		t.Errorf("expected ErrInvalidVerificationToken, got %v", err)
	}
}

func TestUserVerifyEmail_MalformedBase64(t *testing.T) {
	_, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	err := uc.VerifyEmail(ctx, "not-valid-base64!!!")
	if !errors.Is(err, ErrInvalidVerificationToken) {
		t.Errorf("expected ErrInvalidVerificationToken, got %v", err)
	}
}

// ============================================================
// ResendVerification
// ============================================================

func TestUserResendVerification_Success(t *testing.T) {
	userRepo, tokenRepo, emailProvider, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.ResendVerificationRequest{Email: "pending@example.com"}

	user := &domain.User{
		ID:     uuid.New(),
		Email:  req.Email,
		Status: "pending",
	}

	oldToken := &domain.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now().Add(-5 * time.Minute), // created 5 min ago, past cooldown
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)
	tokenRepo.EXPECT().FindLatestByUserID(ctx, user.ID).Return(oldToken, nil)
	tokenRepo.EXPECT().InvalidateByUserID(ctx, user.ID).Return(nil)
	tokenRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	emailProvider.EXPECT().Send(ctx, gomock.Any()).Return(nil)

	err := uc.ResendVerification(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserResendVerification_UserAlreadyActive(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.ResendVerificationRequest{Email: "active@example.com"}

	user := &domain.User{
		ID:     uuid.New(),
		Email:  req.Email,
		Status: "active",
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)

	err := uc.ResendVerification(ctx, req)
	if !errors.Is(err, ErrUserAlreadyActive) {
		t.Errorf("expected ErrUserAlreadyActive, got %v", err)
	}
}

func TestUserResendVerification_RateLimit(t *testing.T) {
	userRepo, tokenRepo, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.ResendVerificationRequest{Email: "pending@example.com"}

	user := &domain.User{
		ID:     uuid.New(),
		Email:  req.Email,
		Status: "pending",
	}

	recentToken := &domain.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now().Add(-10 * time.Second), // created 10 sec ago, within cooldown
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)
	tokenRepo.EXPECT().FindLatestByUserID(ctx, user.ID).Return(recentToken, nil)

	err := uc.ResendVerification(ctx, req)
	if !errors.Is(err, ErrResendTooSoon) {
		t.Errorf("expected ErrResendTooSoon, got %v", err)
	}
}

func TestUserResendVerification_UserNotFound(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.ResendVerificationRequest{Email: "unknown@example.com"}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)

	// Should return nil (not leak whether email exists)
	err := uc.ResendVerification(ctx, req)
	if err != nil {
		t.Fatalf("expected nil error for unknown email, got %v", err)
	}
}

func TestUserResendVerification_UserBlocked(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	req := domain.ResendVerificationRequest{Email: "blocked@example.com"}

	user := &domain.User{
		ID:     uuid.New(),
		Email:  req.Email,
		Status: "blocked",
	}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(user, nil)

	err := uc.ResendVerification(ctx, req)
	if !errors.Is(err, ErrUserNotPending) {
		t.Errorf("expected ErrUserNotPending, got %v", err)
	}
}

// ============================================================
// Login
// ============================================================

func TestUserLogin_Success(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	now := time.Now()
	verifiedAt := now.Add(-1 * time.Hour)
	phone := "08123456789"

	user := &domain.User{
		ID:              uuid.New(),
		Email:           "john@example.com",
		FullName:        "John Doe",
		Phone:           &phone,
		PasswordHash:    mustHashPassword("password123"),
		Status:          "active",
		EmailVerifiedAt: &verifiedAt,
		Role:            &domain.Role{ID: 3, Code: "customer", Name: "Customer"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	userRepo.EXPECT().FindByEmail(ctx, user.Email).Return(user, nil)

	req := domain.LoginRequest{Email: user.Email, Password: "password123"}
	resp, err := uc.Login(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.User.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, resp.User.Email)
	}
	if resp.User.Role != "customer" {
		t.Errorf("expected role customer, got %s", resp.User.Role)
	}
}

func TestUserLogin_InvalidCredentials_UserNotFound(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	userRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)

	req := domain.LoginRequest{Email: "unknown@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrUserInvalidCredentials) {
		t.Errorf("expected ErrUserInvalidCredentials, got %v", err)
	}
}

func TestUserLogin_InvalidCredentials_WrongPassword(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "john@example.com",
		PasswordHash: mustHashPassword("correctpassword"),
		Status:       "active",
	}

	userRepo.EXPECT().FindByEmail(ctx, user.Email).Return(user, nil)

	req := domain.LoginRequest{Email: user.Email, Password: "wrongpassword"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrUserInvalidCredentials) {
		t.Errorf("expected ErrUserInvalidCredentials, got %v", err)
	}
}

func TestUserLogin_AccountBlocked(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "blocked@example.com",
		PasswordHash: mustHashPassword("password123"),
		Status:       "blocked",
	}

	userRepo.EXPECT().FindByEmail(ctx, user.Email).Return(user, nil)

	req := domain.LoginRequest{Email: user.Email, Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrUserAccountBlocked) {
		t.Errorf("expected ErrUserAccountBlocked, got %v", err)
	}
}

func TestUserLogin_AccountNotActive(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	user := &domain.User{
		ID:           uuid.New(),
		Email:        "pending@example.com",
		PasswordHash: mustHashPassword("password123"),
		Status:       "pending",
	}

	userRepo.EXPECT().FindByEmail(ctx, user.Email).Return(user, nil)

	req := domain.LoginRequest{Email: user.Email, Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrUserAccountNotActive) {
		t.Errorf("expected ErrUserAccountNotActive, got %v", err)
	}
}

func TestUserLogin_EmailNotVerified(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	user := &domain.User{
		ID:              uuid.New(),
		Email:           "unverified@example.com",
		PasswordHash:    mustHashPassword("password123"),
		Status:          "active",
		EmailVerifiedAt: nil,
	}

	userRepo.EXPECT().FindByEmail(ctx, user.Email).Return(user, nil)

	req := domain.LoginRequest{Email: user.Email, Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrUserEmailNotVerified) {
		t.Errorf("expected ErrUserEmailNotVerified, got %v", err)
	}
}

func TestUserLogin_FindByEmailError(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	userRepo.EXPECT().FindByEmail(ctx, "any@example.com").Return(nil, errors.New("db error"))

	req := domain.LoginRequest{Email: "any@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// GetMe
// ============================================================

func TestUserGetMe_Success(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	now := time.Now()
	verifiedAt := now.Add(-1 * time.Hour)
	phone := "08123456789"

	userID := uuid.New()
	user := &domain.User{
		ID:              userID,
		Email:           "john@example.com",
		FullName:        "John Doe",
		Phone:           &phone,
		Status:          "active",
		EmailVerifiedAt: &verifiedAt,
		Role:            &domain.Role{ID: 3, Code: "customer", Name: "Customer"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	userRepo.EXPECT().FindByID(ctx, userID).Return(user, nil)

	resp, err := uc.GetMe(ctx, userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, resp.ID)
	}
	if resp.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, resp.Email)
	}
	if resp.FullName != user.FullName {
		t.Errorf("expected full name %s, got %s", user.FullName, resp.FullName)
	}
	if resp.Role != "customer" {
		t.Errorf("expected role customer, got %s", resp.Role)
	}
}

func TestUserGetMe_UserNotFound(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	userRepo.EXPECT().FindByID(ctx, userID).Return(nil, nil)

	_, err := uc.GetMe(ctx, userID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserGetMe_RepoError(t *testing.T) {
	userRepo, _, _, uc := setupUserUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	userRepo.EXPECT().FindByID(ctx, userID).Return(nil, errors.New("db error"))

	_, err := uc.GetMe(ctx, userID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// mustHashPassword is a test helper that hashes a password or panics.
func mustHashPassword(password string) string {
	hash, err := auth.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}
