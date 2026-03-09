package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
	"go.uber.org/mock/gomock"
)

func setupAdminAuthUseCase(t *testing.T) (*mocks.MockUserRepository, domain.AdminAuthUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := NewAdminAuthUseCase(userRepo, "test-secret", 24, "test-issuer")
	return userRepo, uc
}

func TestAdminAuthUseCaseLogin(t *testing.T) {
	ctx := context.Background()
	imageURL := "https://cdn.example.com/admins/profile.jpg"

	testCases := []struct {
		name       string
		req        domain.AdminLoginRequest
		setupMocks func(*mocks.MockUserRepository)
		wantErr    error
		assertResp func(*testing.T, *domain.AdminLoginResponse)
	}{
		{
			name: "success with image url",
			req:  domain.AdminLoginRequest{Email: "admin@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "admin@example.com",
					FullName:     "Admin Ops",
					ImageURL:     &imageURL,
					PasswordHash: mustHashPasswordAdmin("password123"),
					Status:       domain.UserStatusActive,
					Role:         &domain.Role{ID: 1, Code: domain.RoleAdmin, Name: "Admin"},
				}
				userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(user, nil)
			},
			assertResp: func(t *testing.T, resp *domain.AdminLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Token == "" {
					t.Fatal("expected non-empty token")
				}
				if resp.ID == uuid.Nil {
					t.Fatal("expected non-empty user ID")
				}
				if resp.ImageURL != imageURL {
					t.Fatalf("expected image URL %q, got %q", imageURL, resp.ImageURL)
				}
				if resp.Email != "admin@example.com" {
					t.Fatalf("expected email admin@example.com, got %s", resp.Email)
				}
				if resp.Name != "Admin Ops" {
					t.Fatalf("expected name Admin Ops, got %s", resp.Name)
				}
			},
		},
		{
			name: "success without image url defaults empty string",
			req:  domain.AdminLoginRequest{Email: "admin@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "admin@example.com",
					FullName:     "Admin Ops",
					PasswordHash: mustHashPasswordAdmin("password123"),
					Status:       domain.UserStatusActive,
					Role:         &domain.Role{ID: 1, Code: domain.RoleAdmin, Name: "Admin"},
				}
				userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(user, nil)
			},
			assertResp: func(t *testing.T, resp *domain.AdminLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.ImageURL != "" {
					t.Fatalf("expected empty image URL, got %q", resp.ImageURL)
				}
			},
		},
		{
			name: "invalid credentials user not found",
			req:  domain.AdminLoginRequest{Email: "unknown@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "invalid credentials wrong password",
			req:  domain.AdminLoginRequest{Email: "admin@example.com", Password: "wrongpassword"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "admin@example.com",
					PasswordHash: mustHashPasswordAdmin("correctpassword"),
					Status:       domain.UserStatusActive,
					Role:         &domain.Role{ID: 1, Code: domain.RoleAdmin, Name: "Admin"},
				}
				userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(user, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "account inactive",
			req:  domain.AdminLoginRequest{Email: "admin@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "admin@example.com",
					PasswordHash: mustHashPasswordAdmin("password123"),
					Status:       domain.UserStatusBlocked,
					Role:         &domain.Role{ID: 1, Code: domain.RoleAdmin, Name: "Admin"},
				}
				userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(user, nil)
			},
			wantErr: ErrAccountInactive,
		},
		{
			name: "not admin",
			req:  domain.AdminLoginRequest{Email: "user@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				user := &domain.User{
					ID:           uuid.New(),
					Email:        "user@example.com",
					PasswordHash: mustHashPasswordAdmin("password123"),
					Status:       domain.UserStatusActive,
					Role:         &domain.Role{ID: 2, Code: domain.RoleCustomer, Name: "Customer"},
				}
				userRepo.EXPECT().FindByEmail(ctx, "user@example.com").Return(user, nil)
			},
			wantErr: ErrNotAdmin,
		},
		{
			name: "find by email error",
			req:  domain.AdminLoginRequest{Email: "admin@example.com", Password: "password123"},
			setupMocks: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, uc := setupAdminAuthUseCase(t)
			if tc.setupMocks != nil {
				tc.setupMocks(userRepo)
			}

			resp, err := uc.Login(ctx, tc.req)
			if tc.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tc.wantErr, ErrInvalidCredentials) ||
					errors.Is(tc.wantErr, ErrAccountInactive) ||
					errors.Is(tc.wantErr, ErrNotAdmin) {
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

func mustHashPasswordAdmin(password string) string {
	hash, err := auth.HashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}
