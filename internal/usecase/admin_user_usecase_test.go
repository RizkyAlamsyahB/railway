package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// ---------- helpers ----------

func setupUseCase(t *testing.T) (*mocks.MockUserRepository, domain.AdminUserUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockUserRepository(ctrl)
	uc := NewAdminUserUseCase(repo)
	return repo, uc
}

func ptrString(s string) *string { return &s }

func dummyUser(id uuid.UUID) *domain.User {
	now := time.Now()
	return &domain.User{
		ID:              id,
		Email:           "user@example.com",
		FullName:        "John Doe",
		Phone:           ptrString("08123456789"),
		PasswordHash:    "$2a$10$hashedpassword",
		Status:          "active",
		EmailVerifiedAt: &now,
		Role:            &domain.Role{ID: 1, Code: "admin", Name: "Admin"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// ============================================================
// Create
// ============================================================

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		req       domain.CreateUserRequest
		setupMock func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest)
		wantErr   error
		wantNil   bool
		checkResp func(t *testing.T, resp *domain.UserResponse)
	}{
		{
			name: "success",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				userID := uuid.New()
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, gomock.Any()).Times(0)
				repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
				repo.EXPECT().FindByID(ctx, gomock.Any()).Return(dummyUser(userID), nil)
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.Role != "admin" {
					t.Errorf("expected role admin, got %s", resp.Role)
				}
			},
		},
		{
			name: "with birth date",
			req: domain.CreateUserRequest{
				Email:     "new@example.com",
				FullName:  "Jane Doe",
				Password:  "password123",
				BirthDate: ptrString("1990-01-15"),
				Role:      "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				userID := uuid.New()
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, gomock.Any()).Times(0)
				repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
				repo.EXPECT().FindByID(ctx, gomock.Any()).Return(dummyUser(userID), nil)
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
			},
		},
		{
			name: "email already exists",
			req: domain.CreateUserRequest{
				Email:    "exists@example.com",
				FullName: "Existing User",
				Password: "password123",
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				existing := dummyUser(uuid.New())
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(existing, nil)
			},
			wantErr: ErrEmailExists,
		},
		{
			name: "phone already exists",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Phone:    ptrString("08123456789"),
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, *req.Phone).Return(dummyUser(uuid.New()), nil)
			},
			wantErr: ErrPhoneAlreadyRegistered,
		},
		{
			name: "invalid birth date",
			req: domain.CreateUserRequest{
				Email:     "new@example.com",
				FullName:  "Jane Doe",
				Password:  "password123",
				BirthDate: ptrString("invalid-date"),
				Role:      "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, gomock.Any()).Times(0)
			},
			wantErr: ErrInvalidBirthDate,
		},
		{
			name: "success with phone",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Phone:    ptrString("08123456789"),
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				userID := uuid.New()
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, *req.Phone).Return(nil, nil)
				repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
				repo.EXPECT().FindByID(ctx, gomock.Any()).Return(dummyUser(userID), nil)
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
			},
		},
		{
			name: "FindByEmail error",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, errors.New("db error"))
			},
			wantNil: true,
		},
		{
			name: "FindByPhone error",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Phone:    ptrString("08123456789"),
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, *req.Phone).Return(nil, errors.New("db error"))
			},
			wantNil: true,
		},
		{
			name: "repo Create error",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, gomock.Any()).Times(0)
				repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(errors.New("db error"))
			},
			wantNil: true,
		},
		{
			name: "FindByID after Create error",
			req: domain.CreateUserRequest{
				Email:    "new@example.com",
				FullName: "Jane Doe",
				Password: "password123",
				Role:     "admin",
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, req domain.CreateUserRequest) {
				repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				repo.EXPECT().FindByPhone(ctx, gomock.Any()).Times(0)
				repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
				repo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupUseCase(t)
			ctx := context.Background()

			tc.setupMock(repo, ctx, tc.req)

			resp, err := uc.Create(ctx, tc.req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantNil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp)
			}
		})
	}
}

// ============================================================
// List
// ============================================================

func TestList(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.UserListParams
		setupMock func(repo *mocks.MockUserRepository, ctx context.Context)
		wantErr   bool
		checkResp func(t *testing.T, responses []domain.UserResponse, meta *domain.PaginationMeta)
	}{
		{
			name:   "success",
			params: domain.UserListParams{Page: 1, Limit: 10},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				users := []domain.User{
					*dummyUser(uuid.New()),
					*dummyUser(uuid.New()),
				}
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 1, Limit: 10}).Return(users, int64(2), nil)
			},
			checkResp: func(t *testing.T, responses []domain.UserResponse, meta *domain.PaginationMeta) {
				t.Helper()
				if len(responses) != 2 {
					t.Errorf("expected 2 responses, got %d", len(responses))
				}
				if meta.TotalItems != 2 {
					t.Errorf("expected total_items=2, got %d", meta.TotalItems)
				}
				if meta.TotalPages != 1 {
					t.Errorf("expected total_pages=1, got %d", meta.TotalPages)
				}
				if meta.Page != 1 {
					t.Errorf("expected page=1, got %d", meta.Page)
				}
			},
		},
		{
			name:   "defaults invalid page and limit",
			params: domain.UserListParams{Page: 0, Limit: 0},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				// After defaults: Page=1, Limit=10
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 1, Limit: 10}).Return([]domain.User{}, int64(0), nil)
			},
			checkResp: func(t *testing.T, responses []domain.UserResponse, meta *domain.PaginationMeta) {
				t.Helper()
				if len(responses) != 0 {
					t.Errorf("expected 0 responses, got %d", len(responses))
				}
				if meta.TotalItems != 0 {
					t.Errorf("expected total_items=0, got %d", meta.TotalItems)
				}
			},
		},
		{
			name:   "limit exceeds 100",
			params: domain.UserListParams{Page: 1, Limit: 200},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 1, Limit: 10}).Return([]domain.User{}, int64(0), nil)
			},
		},
		{
			name:   "repo error",
			params: domain.UserListParams{Page: 1, Limit: 10},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 1, Limit: 10}).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "pagination calculation",
			params: domain.UserListParams{Page: 2, Limit: 3},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 2, Limit: 3}).Return([]domain.User{*dummyUser(uuid.New())}, int64(7), nil)
			},
			checkResp: func(t *testing.T, responses []domain.UserResponse, meta *domain.PaginationMeta) {
				t.Helper()
				// 7 items / 3 per page = ceil(2.33) = 3 pages
				if meta.TotalPages != 3 {
					t.Errorf("expected total_pages=3, got %d", meta.TotalPages)
				}
			},
		},
		{
			name:   "sort by full_name asc",
			params: domain.UserListParams{Page: 1, Limit: 10, SortBy: "full_name", SortOrder: "asc"},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 1, Limit: 10, SortBy: "full_name", SortOrder: "asc"}).Return([]domain.User{*dummyUser(uuid.New())}, int64(1), nil)
			},
			checkResp: func(t *testing.T, responses []domain.UserResponse, meta *domain.PaginationMeta) {
				t.Helper()
				if len(responses) != 1 {
					t.Errorf("expected 1 response, got %d", len(responses))
				}
			},
		},
		{
			name:   "sort by email desc",
			params: domain.UserListParams{Page: 1, Limit: 10, SortBy: "email", SortOrder: "desc"},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context) {
				repo.EXPECT().List(ctx, domain.UserListParams{Page: 1, Limit: 10, SortBy: "email", SortOrder: "desc"}).Return([]domain.User{*dummyUser(uuid.New()), *dummyUser(uuid.New())}, int64(2), nil)
			},
			checkResp: func(t *testing.T, responses []domain.UserResponse, meta *domain.PaginationMeta) {
				t.Helper()
				if len(responses) != 2 {
					t.Errorf("expected 2 responses, got %d", len(responses))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupUseCase(t)
			ctx := context.Background()

			tc.setupMock(repo, ctx)

			responses, meta, err := uc.List(ctx, tc.params)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, responses, meta)
			}
		})
	}
}

// ============================================================
// GetByID
// ============================================================

func TestGetByID(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID)
		wantErr   error
		wantAny   bool // expect any non-nil error (not a specific sentinel)
		checkResp func(t *testing.T, resp *domain.UserResponse, id uuid.UUID)
	}{
		{
			name: "success",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				user := dummyUser(id)
				repo.EXPECT().FindByID(ctx, id).Return(user, nil)
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse, id uuid.UUID) {
				t.Helper()
				if resp.ID != id {
					t.Errorf("expected id %s, got %s", id, resp.ID)
				}
				if resp.Email != "user@example.com" {
					t.Errorf("expected email user@example.com, got %s", resp.Email)
				}
			},
		},
		{
			name: "not found",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "repo error",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupUseCase(t)
			ctx := context.Background()
			id := uuid.New()

			tc.setupMock(repo, ctx, id)

			resp, err := uc.GetByID(ctx, id)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp, id)
			}
		})
	}
}

// ============================================================
// GetMe
// ============================================================

func TestGetMe(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID)
		wantErr   error
		wantAny   bool
		checkResp func(t *testing.T, resp *domain.AdminMeResponse, id uuid.UUID)
	}{
		{
			name: "success with image url",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				user := dummyUser(id)
				imageURL := "https://cdn.example.com/admins/profile.jpg"
				user.ImageURL = &imageURL
				repo.EXPECT().FindByID(ctx, id).Return(user, nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminMeResponse, id uuid.UUID) {
				t.Helper()
				if resp.ID != id {
					t.Errorf("expected id %s, got %s", id, resp.ID)
				}
				if resp.ImageURL != "https://cdn.example.com/admins/profile.jpg" {
					t.Errorf("expected image URL to be mapped, got %q", resp.ImageURL)
				}
				if resp.Name != "John Doe" {
					t.Errorf("expected name John Doe, got %s", resp.Name)
				}
			},
		},
		{
			name: "success with nil image url becomes empty string",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				user := dummyUser(id)
				user.ImageURL = nil
				repo.EXPECT().FindByID(ctx, id).Return(user, nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminMeResponse, _ uuid.UUID) {
				t.Helper()
				if resp.ImageURL != "" {
					t.Errorf("expected empty image URL, got %q", resp.ImageURL)
				}
			},
		},
		{
			name: "not found",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "not admin",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				user := dummyUser(id)
				user.Role = &domain.Role{ID: 2, Code: domain.RoleCustomer, Name: "Customer"}
				repo.EXPECT().FindByID(ctx, id).Return(user, nil)
			},
			wantErr: ErrNotAdmin,
		},
		{
			name: "repo error",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupUseCase(t)
			ctx := context.Background()
			id := uuid.New()

			tc.setupMock(repo, ctx, id)

			resp, err := uc.GetMe(ctx, id)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp, id)
			}
		})
	}
}

// ============================================================
// Update
// ============================================================

func TestUpdate(t *testing.T) {
	tests := []struct {
		name      string
		req       domain.UpdateUserRequest
		setupMock func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest)
		wantErr   error
		wantAny   bool
		checkResp func(t *testing.T, resp *domain.UserResponse)
	}{
		{
			name: "full update with role",
			req: domain.UpdateUserRequest{
				FullName:  ptrString("Updated Name"),
				Phone:     ptrString("08999999999"),
				BirthDate: ptrString("1995-06-15"),
				Status:    ptrString("blocked"),
				Role:      ptrString("cs"),
				Password:  ptrString("newpassword123"),
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				existing := dummyUser(id)
				updated := dummyUser(id)
				updated.FullName = "Updated Name"

				repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
				repo.EXPECT().UpdateWithRole(ctx, gomock.Any(), "cs").Return(nil)
				repo.EXPECT().FindByID(ctx, id).Return(updated, nil)
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
			},
		},
		{
			name: "partial update without role",
			req: domain.UpdateUserRequest{
				FullName: ptrString("New Name"),
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				existing := dummyUser(id)
				updated := dummyUser(id)

				repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
				repo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
				repo.EXPECT().FindByID(ctx, id).Return(updated, nil)
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
			},
		},
		{
			name: "user not found",
			req:  domain.UpdateUserRequest{FullName: ptrString("Name")},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "invalid birth date",
			req: domain.UpdateUserRequest{
				BirthDate: ptrString("not-a-date"),
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				existing := dummyUser(id)
				repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
			},
			wantErr: ErrInvalidBirthDate,
		},
		{
			name: "FindByID error",
			req:  domain.UpdateUserRequest{FullName: ptrString("Name")},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "repo Update error",
			req:  domain.UpdateUserRequest{FullName: ptrString("Name")},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				existing := dummyUser(id)
				repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
				repo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "repo UpdateWithRole error",
			req: domain.UpdateUserRequest{
				FullName: ptrString("Name"),
				Role:     ptrString("admin"),
			},
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) {
				existing := dummyUser(id)
				repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
				repo.EXPECT().UpdateWithRole(ctx, gomock.Any(), "admin").Return(errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupUseCase(t)
			ctx := context.Background()
			id := uuid.New()

			tc.setupMock(repo, ctx, id, tc.req)

			resp, err := uc.Update(ctx, id, tc.req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp)
			}
		})
	}
}

// ============================================================
// Delete
// ============================================================

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		sameActor bool // if true, actorID == id (delete self)
		setupMock func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID)
		wantErr   error
		wantAny   bool
	}{
		{
			name: "success",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				user := dummyUser(id)
				repo.EXPECT().FindByID(ctx, id).Return(user, nil)
				repo.EXPECT().Delete(ctx, id).Return(nil)
			},
		},
		{
			name:      "cannot delete self",
			sameActor: true,
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				// no repo calls expected
			},
			wantErr: ErrCannotDeleteSelf,
		},
		{
			name: "user not found",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "FindByID error",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "repo Delete error",
			setupMock: func(repo *mocks.MockUserRepository, ctx context.Context, id uuid.UUID) {
				user := dummyUser(id)
				repo.EXPECT().FindByID(ctx, id).Return(user, nil)
				repo.EXPECT().Delete(ctx, id).Return(errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupUseCase(t)
			ctx := context.Background()
			id := uuid.New()

			actorID := uuid.New()
			if tc.sameActor {
				actorID = id
			}

			tc.setupMock(repo, ctx, id)

			err := uc.Delete(ctx, id, actorID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

// ============================================================
// toUserResponse (helper)
// ============================================================

func TestToUserResponse(t *testing.T) {
	tests := []struct {
		name      string
		modifyFn  func(u *domain.User)
		checkResp func(t *testing.T, resp *domain.UserResponse, user *domain.User)
	}{
		{
			name:     "with role",
			modifyFn: nil, // use dummyUser as-is
			checkResp: func(t *testing.T, resp *domain.UserResponse, user *domain.User) {
				t.Helper()
				if resp.Role != "admin" {
					t.Errorf("expected role admin, got %s", resp.Role)
				}
				if resp.ID != user.ID {
					t.Errorf("expected id %s, got %s", user.ID, resp.ID)
				}
				if resp.Email != user.Email {
					t.Errorf("expected email %s, got %s", user.Email, resp.Email)
				}
			},
		},
		{
			name: "without role",
			modifyFn: func(u *domain.User) {
				u.Role = nil
			},
			checkResp: func(t *testing.T, resp *domain.UserResponse, user *domain.User) {
				t.Helper()
				if resp.Role != "" {
					t.Errorf("expected empty role, got %s", resp.Role)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user := dummyUser(uuid.New())
			if tc.modifyFn != nil {
				tc.modifyFn(user)
			}

			resp := toUserResponse(user)
			tc.checkResp(t, resp, user)
		})
	}
}
