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

func TestUserRegister(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		req        domain.RegisterCustomerRequest
		setupMocks func(*mocks.MockUserRepository, *mocks.MockEmailVerificationTokenRepository, *mocks.MockEmailProvider)
		wantErr    error
		assertResp func(*testing.T, *domain.RegisterCustomerResponse)
	}{
		{
			name: "success",
			req: domain.RegisterCustomerRequest{
				FullName: "John Doe",
				Phone:    "08123456789",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockEmailVerificationTokenRepository, emailProvider *mocks.MockEmailProvider) {
				userRepo.EXPECT().FindByEmail(ctx, "john@example.com").Return(nil, nil)
				userRepo.EXPECT().FindByPhone(ctx, "08123456789").Return(nil, nil)
				userRepo.EXPECT().Create(ctx, gomock.Any(), "customer").Return(nil)
				tokenRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				emailProvider.EXPECT().Send(ctx, gomock.Any()).Return(nil)
			},
			assertResp: func(t *testing.T, resp *domain.RegisterCustomerResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Email != "john@example.com" {
					t.Errorf("expected email john@example.com, got %s", resp.Email)
				}
				if resp.Message == "" {
					t.Error("expected non-empty message")
				}
			},
		},
		{
			name: "email already exists",
			req: domain.RegisterCustomerRequest{
				FullName: "John Doe",
				Phone:    "08123456789",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				existingUser := &domain.User{ID: uuid.New(), Email: "existing@example.com"}
				userRepo.EXPECT().FindByEmail(ctx, "existing@example.com").Return(existingUser, nil)
			},
			wantErr: ErrEmailAlreadyRegistered,
		},
		{
			name: "phone already exists",
			req: domain.RegisterCustomerRequest{
				FullName: "John Doe",
				Phone:    "08123456789",
				Email:    "new@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				phone := "08123456789"
				existingUser := &domain.User{ID: uuid.New(), Phone: &phone}
				userRepo.EXPECT().FindByEmail(ctx, "new@example.com").Return(nil, nil)
				userRepo.EXPECT().FindByPhone(ctx, "08123456789").Return(existingUser, nil)
			},
			wantErr: ErrPhoneAlreadyRegistered,
		},
		{
			name: "find by email error",
			req: domain.RegisterCustomerRequest{
				FullName: "John Doe",
				Phone:    "08123456789",
				Email:    "new@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				userRepo.EXPECT().FindByEmail(ctx, "new@example.com").Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "create user error",
			req: domain.RegisterCustomerRequest{
				FullName: "John Doe",
				Phone:    "08123456789",
				Email:    "new@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				userRepo.EXPECT().FindByEmail(ctx, "new@example.com").Return(nil, nil)
				userRepo.EXPECT().FindByPhone(ctx, "08123456789").Return(nil, nil)
				userRepo.EXPECT().Create(ctx, gomock.Any(), "customer").Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, tokenRepo, emailProvider, uc := setupUserUseCase(t)
			if tc.setupMocks != nil {
				tc.setupMocks(userRepo, tokenRepo, emailProvider)
			}

			resp, err := uc.Register(ctx, tc.req)

			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tc.wantErr, ErrEmailAlreadyRegistered) || errors.Is(tc.wantErr, ErrPhoneAlreadyRegistered) {
					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected %v, got %v", tc.wantErr, err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.assertResp != nil {
				tc.assertResp(t, resp)
			}
		})
	}
}

func TestUserVerifyEmail(t *testing.T) {
	ctx := context.Background()

	buildTokenAndHash := func(rawToken []byte) (string, string) {
		rawTokenB64 := base64.URLEncoding.EncodeToString(rawToken)
		hash := sha256.Sum256(rawToken)
		tokenHash := base64.URLEncoding.EncodeToString(hash[:])
		return rawTokenB64, tokenHash
	}

	testCases := []struct {
		name       string
		rawToken   []byte
		malformed  string
		setupMocks func(*mocks.MockUserRepository, *mocks.MockEmailVerificationTokenRepository, string)
		wantErr    error
	}{
		{
			name: "success",
			rawToken: func() []byte {
				raw := make([]byte, 32)
				for i := range raw {
					raw[i] = byte(i)
				}
				return raw
			}(),
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockEmailVerificationTokenRepository, tokenHash string) {
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
			},
		},
		{
			name:     "invalid token",
			rawToken: make([]byte, 32),
			setupMocks: func(_ *mocks.MockUserRepository, tokenRepo *mocks.MockEmailVerificationTokenRepository, tokenHash string) {
				tokenRepo.EXPECT().FindActiveByTokenHash(ctx, tokenHash).Return(nil, nil)
			},
			wantErr: ErrInvalidVerificationToken,
		},
		{
			name: "expired token",
			rawToken: func() []byte {
				raw := make([]byte, 32)
				for i := range raw {
					raw[i] = byte(i + 1)
				}
				return raw
			}(),
			setupMocks: func(_ *mocks.MockUserRepository, tokenRepo *mocks.MockEmailVerificationTokenRepository, tokenHash string) {
				token := &domain.EmailVerificationToken{
					ID:        uuid.New(),
					UserID:    uuid.New(),
					Email:     "john@example.com",
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(-1 * time.Hour),
					CreatedAt: time.Now().Add(-25 * time.Hour),
				}
				tokenRepo.EXPECT().FindActiveByTokenHash(ctx, tokenHash).Return(token, nil)
			},
			wantErr: ErrInvalidVerificationToken,
		},
		{
			name:      "malformed base64",
			malformed: "not-valid-base64!!!",
			wantErr:   ErrInvalidVerificationToken,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, tokenRepo, _, uc := setupUserUseCase(t)

			rawToken := tc.malformed
			if tc.rawToken != nil {
				encoded, tokenHash := buildTokenAndHash(tc.rawToken)
				rawToken = encoded
				if tc.setupMocks != nil {
					tc.setupMocks(userRepo, tokenRepo, tokenHash)
				}
			} else if tc.setupMocks != nil {
				tc.setupMocks(userRepo, tokenRepo, "")
			}

			err := uc.VerifyEmail(ctx, rawToken)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestUserResendVerification(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		req        domain.ResendVerificationRequest
		setupMocks func(*mocks.MockUserRepository, *mocks.MockEmailVerificationTokenRepository, *mocks.MockEmailProvider)
		wantErr    error
	}{
		{
			name: "success",
			req:  domain.ResendVerificationRequest{Email: "pending@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockEmailVerificationTokenRepository, emailProvider *mocks.MockEmailProvider) {
				user := &domain.User{ID: uuid.New(), Email: "pending@example.com", Status: "pending"}
				oldToken := &domain.EmailVerificationToken{ID: uuid.New(), UserID: user.ID, CreatedAt: time.Now().Add(-5 * time.Minute)}

				userRepo.EXPECT().FindByEmail(ctx, "pending@example.com").Return(user, nil)
				tokenRepo.EXPECT().FindLatestByUserID(ctx, user.ID).Return(oldToken, nil)
				tokenRepo.EXPECT().InvalidateByUserID(ctx, user.ID).Return(nil)
				tokenRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				emailProvider.EXPECT().Send(ctx, gomock.Any()).Return(nil)
			},
		},
		{
			name: "user already active",
			req:  domain.ResendVerificationRequest{Email: "active@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				user := &domain.User{ID: uuid.New(), Email: "active@example.com", Status: "active"}
				userRepo.EXPECT().FindByEmail(ctx, "active@example.com").Return(user, nil)
			},
			wantErr: ErrUserAlreadyActive,
		},
		{
			name: "rate limit",
			req:  domain.ResendVerificationRequest{Email: "pending@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, tokenRepo *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				user := &domain.User{ID: uuid.New(), Email: "pending@example.com", Status: "pending"}
				recentToken := &domain.EmailVerificationToken{ID: uuid.New(), UserID: user.ID, CreatedAt: time.Now().Add(-10 * time.Second)}

				userRepo.EXPECT().FindByEmail(ctx, "pending@example.com").Return(user, nil)
				tokenRepo.EXPECT().FindLatestByUserID(ctx, user.ID).Return(recentToken, nil)
			},
			wantErr: ErrResendTooSoon,
		},
		{
			name: "user not found",
			req:  domain.ResendVerificationRequest{Email: "unknown@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				userRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)
			},
		},
		{
			name: "user blocked",
			req:  domain.ResendVerificationRequest{Email: "blocked@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, _ *mocks.MockEmailVerificationTokenRepository, _ *mocks.MockEmailProvider) {
				user := &domain.User{ID: uuid.New(), Email: "blocked@example.com", Status: "blocked"}
				userRepo.EXPECT().FindByEmail(ctx, "blocked@example.com").Return(user, nil)
			},
			wantErr: ErrUserNotPending,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, tokenRepo, emailProvider, uc := setupUserUseCase(t)
			if tc.setupMocks != nil {
				tc.setupMocks(userRepo, tokenRepo, emailProvider)
			}

			err := uc.ResendVerification(ctx, tc.req)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	verifiedAt := now.Add(-1 * time.Hour)
	phone := "08123456789"

	testCases := []struct {
		name       string
		req        domain.LoginRequest
		setupMocks func(*mocks.MockUserRepository)
		wantErr    error
		assertResp func(*testing.T, *domain.LoginResponse)
	}{
		{
			name: "success",
			req:  domain.LoginRequest{Email: "john@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
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
				userRepo.EXPECT().FindByEmail(ctx, "john@example.com").Return(user, nil)
			},
			assertResp: func(t *testing.T, resp *domain.LoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Token == "" {
					t.Error("expected non-empty token")
				}
				if resp.User.Email != "john@example.com" {
					t.Errorf("expected email john@example.com, got %s", resp.User.Email)
				}
				if resp.User.Role != "customer" {
					t.Errorf("expected role customer, got %s", resp.User.Role)
				}
			},
		},
		{
			name: "invalid credentials user not found",
			req:  domain.LoginRequest{Email: "unknown@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)
			},
			wantErr: ErrUserInvalidCredentials,
		},
		{
			name: "invalid credentials wrong password",
			req:  domain.LoginRequest{Email: "john@example.com", Password: "wrongpassword"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: uuid.New(), Email: "john@example.com", PasswordHash: mustHashPassword("correctpassword"), Status: "active"}
				userRepo.EXPECT().FindByEmail(ctx, "john@example.com").Return(user, nil)
			},
			wantErr: ErrUserInvalidCredentials,
		},
		{
			name: "account blocked",
			req:  domain.LoginRequest{Email: "blocked@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: uuid.New(), Email: "blocked@example.com", PasswordHash: mustHashPassword("password123"), Status: "blocked"}
				userRepo.EXPECT().FindByEmail(ctx, "blocked@example.com").Return(user, nil)
			},
			wantErr: ErrUserAccountBlocked,
		},
		{
			name: "account not active",
			req:  domain.LoginRequest{Email: "pending@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: uuid.New(), Email: "pending@example.com", PasswordHash: mustHashPassword("password123"), Status: "pending"}
				userRepo.EXPECT().FindByEmail(ctx, "pending@example.com").Return(user, nil)
			},
			wantErr: ErrUserAccountNotActive,
		},
		{
			name: "email not verified",
			req:  domain.LoginRequest{Email: "unverified@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{ID: uuid.New(), Email: "unverified@example.com", PasswordHash: mustHashPassword("password123"), Status: "active", EmailVerifiedAt: nil}
				userRepo.EXPECT().FindByEmail(ctx, "unverified@example.com").Return(user, nil)
			},
			wantErr: ErrUserEmailNotVerified,
		},
		{
			name: "find by email error",
			req:  domain.LoginRequest{Email: "any@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByEmail(ctx, "any@example.com").Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, _, _, uc := setupUserUseCase(t)
			if tc.setupMocks != nil {
				tc.setupMocks(userRepo)
			}

			resp, err := uc.Login(ctx, tc.req)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tc.wantErr, ErrUserInvalidCredentials) ||
					errors.Is(tc.wantErr, ErrUserAccountBlocked) ||
					errors.Is(tc.wantErr, ErrUserAccountNotActive) ||
					errors.Is(tc.wantErr, ErrUserEmailNotVerified) {
					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected %v, got %v", tc.wantErr, err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.assertResp != nil {
				tc.assertResp(t, resp)
			}
		})
	}
}

func TestUserGetMe(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		userID     uuid.UUID
		setupMocks func(*mocks.MockUserRepository, uuid.UUID)
		wantErr    error
		assertResp func(*testing.T, *domain.UserResponse, uuid.UUID)
	}{
		{
			name:   "success",
			userID: uuid.New(),
			setupMocks: func(userRepo *mocks.MockUserRepository, userID uuid.UUID) {
				now := time.Now()
				verifiedAt := now.Add(-1 * time.Hour)
				phone := "08123456789"
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
			},
			assertResp: func(t *testing.T, resp *domain.UserResponse, userID uuid.UUID) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.ID != userID {
					t.Errorf("expected user ID %s, got %s", userID, resp.ID)
				}
				if resp.Email != "john@example.com" {
					t.Errorf("expected email john@example.com, got %s", resp.Email)
				}
				if resp.FullName != "John Doe" {
					t.Errorf("expected full name John Doe, got %s", resp.FullName)
				}
				if resp.Role != "customer" {
					t.Errorf("expected role customer, got %s", resp.Role)
				}
			},
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			setupMocks: func(userRepo *mocks.MockUserRepository, userID uuid.UUID) {
				userRepo.EXPECT().FindByID(ctx, userID).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:   "repo error",
			userID: uuid.New(),
			setupMocks: func(userRepo *mocks.MockUserRepository, userID uuid.UUID) {
				userRepo.EXPECT().FindByID(ctx, userID).Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, _, _, uc := setupUserUseCase(t)
			if tc.setupMocks != nil {
				tc.setupMocks(userRepo, tc.userID)
			}

			resp, err := uc.GetMe(ctx, tc.userID)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tc.wantErr, ErrUserNotFound) {
					if !errors.Is(err, tc.wantErr) {
						t.Fatalf("expected %v, got %v", tc.wantErr, err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.assertResp != nil {
				tc.assertResp(t, resp, tc.userID)
			}
		})
	}
}

func mustHashPassword(password string) string {
	hash, err := auth.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}
