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

func TestCreate_Success(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	userID := uuid.New()

	req := domain.CreateUserRequest{
		Email:    "new@example.com",
		FullName: "Jane Doe",
		Password: "password123",
		Role:     "admin",
	}

	repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
	repo.EXPECT().FindByID(ctx, gomock.Any()).Return(dummyUser(userID), nil)

	resp, err := uc.Create(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Role != "admin" {
		t.Errorf("expected role admin, got %s", resp.Role)
	}
}

func TestCreate_WithBirthDate(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	userID := uuid.New()

	req := domain.CreateUserRequest{
		Email:     "new@example.com",
		FullName:  "Jane Doe",
		Password:  "password123",
		BirthDate: ptrString("1990-01-15"),
		Role:      "admin",
	}

	repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
	repo.EXPECT().FindByID(ctx, gomock.Any()).Return(dummyUser(userID), nil)

	resp, err := uc.Create(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestCreate_EmailAlreadyExists(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	req := domain.CreateUserRequest{
		Email:    "exists@example.com",
		FullName: "Existing User",
		Password: "password123",
		Role:     "admin",
	}

	existing := dummyUser(uuid.New())
	repo.EXPECT().FindByEmail(ctx, req.Email).Return(existing, nil)

	_, err := uc.Create(ctx, req)
	if !errors.Is(err, ErrEmailExists) {
		t.Errorf("expected ErrEmailExists, got %v", err)
	}
}

func TestCreate_InvalidBirthDate(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	req := domain.CreateUserRequest{
		Email:     "new@example.com",
		FullName:  "Jane Doe",
		Password:  "password123",
		BirthDate: ptrString("invalid-date"),
		Role:      "admin",
	}

	repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)

	_, err := uc.Create(ctx, req)
	if !errors.Is(err, ErrInvalidBirthDate) {
		t.Errorf("expected ErrInvalidBirthDate, got %v", err)
	}
}

func TestCreate_FindByEmailError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	req := domain.CreateUserRequest{
		Email:    "new@example.com",
		FullName: "Jane Doe",
		Password: "password123",
		Role:     "admin",
	}

	repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, errors.New("db error"))

	_, err := uc.Create(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreate_RepoCreateError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	req := domain.CreateUserRequest{
		Email:    "new@example.com",
		FullName: "Jane Doe",
		Password: "password123",
		Role:     "admin",
	}

	repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(errors.New("db error"))

	_, err := uc.Create(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreate_FindByIDAfterCreateError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	req := domain.CreateUserRequest{
		Email:    "new@example.com",
		FullName: "Jane Doe",
		Password: "password123",
		Role:     "admin",
	}

	repo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	repo.EXPECT().Create(ctx, gomock.Any(), req.Role).Return(nil)
	repo.EXPECT().FindByID(ctx, gomock.Any()).Return(nil, errors.New("db error"))

	_, err := uc.Create(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// List
// ============================================================

func TestList_Success(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	params := domain.UserListParams{Page: 1, Limit: 10}
	users := []domain.User{
		*dummyUser(uuid.New()),
		*dummyUser(uuid.New()),
	}

	repo.EXPECT().List(ctx, params).Return(users, int64(2), nil)

	responses, meta, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
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
}

func TestList_DefaultsInvalidPage(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	params := domain.UserListParams{Page: 0, Limit: 0}
	// After defaults: Page=1, Limit=10
	expectedParams := domain.UserListParams{Page: 1, Limit: 10}

	repo.EXPECT().List(ctx, expectedParams).Return([]domain.User{}, int64(0), nil)

	responses, meta, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(responses) != 0 {
		t.Errorf("expected 0 responses, got %d", len(responses))
	}
	if meta.TotalItems != 0 {
		t.Errorf("expected total_items=0, got %d", meta.TotalItems)
	}
}

func TestList_LimitExceeds100(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	params := domain.UserListParams{Page: 1, Limit: 200}
	expectedParams := domain.UserListParams{Page: 1, Limit: 10}

	repo.EXPECT().List(ctx, expectedParams).Return([]domain.User{}, int64(0), nil)

	_, _, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestList_RepoError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	params := domain.UserListParams{Page: 1, Limit: 10}

	repo.EXPECT().List(ctx, params).Return(nil, int64(0), errors.New("db error"))

	_, _, err := uc.List(ctx, params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestList_PaginationCalculation(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()

	params := domain.UserListParams{Page: 2, Limit: 3}

	repo.EXPECT().List(ctx, params).Return([]domain.User{*dummyUser(uuid.New())}, int64(7), nil)

	_, meta, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// 7 items / 3 per page = ceil(2.33) = 3 pages
	if meta.TotalPages != 3 {
		t.Errorf("expected total_pages=3, got %d", meta.TotalPages)
	}
}

// ============================================================
// GetByID
// ============================================================

func TestGetByID_Success(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	user := dummyUser(id)
	repo.EXPECT().FindByID(ctx, id).Return(user, nil)

	resp, err := uc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != id {
		t.Errorf("expected id %s, got %s", id, resp.ID)
	}
	if resp.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, resp.Email)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	repo.EXPECT().FindByID(ctx, id).Return(nil, nil)

	_, err := uc.GetByID(ctx, id)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetByID_RepoError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))

	_, err := uc.GetByID(ctx, id)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Update
// ============================================================

func TestUpdate_FullUpdateWithRole(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	existing := dummyUser(id)
	updated := dummyUser(id)
	updated.FullName = "Updated Name"

	req := domain.UpdateUserRequest{
		FullName:  ptrString("Updated Name"),
		Phone:     ptrString("08999999999"),
		BirthDate: ptrString("1995-06-15"),
		Status:    ptrString("blocked"),
		Role:      ptrString("cs"),
		Password:  ptrString("newpassword123"),
	}

	repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
	repo.EXPECT().UpdateWithRole(ctx, gomock.Any(), "cs").Return(nil)
	repo.EXPECT().FindByID(ctx, id).Return(updated, nil)

	resp, err := uc.Update(ctx, id, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestUpdate_PartialUpdateWithoutRole(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	existing := dummyUser(id)
	updated := dummyUser(id)

	req := domain.UpdateUserRequest{
		FullName: ptrString("New Name"),
	}

	repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
	repo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
	repo.EXPECT().FindByID(ctx, id).Return(updated, nil)

	resp, err := uc.Update(ctx, id, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}

func TestUpdate_UserNotFound(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	req := domain.UpdateUserRequest{FullName: ptrString("Name")}

	repo.EXPECT().FindByID(ctx, id).Return(nil, nil)

	_, err := uc.Update(ctx, id, req)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdate_InvalidBirthDate(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	existing := dummyUser(id)
	req := domain.UpdateUserRequest{
		BirthDate: ptrString("not-a-date"),
	}

	repo.EXPECT().FindByID(ctx, id).Return(existing, nil)

	_, err := uc.Update(ctx, id, req)
	if !errors.Is(err, ErrInvalidBirthDate) {
		t.Errorf("expected ErrInvalidBirthDate, got %v", err)
	}
}

func TestUpdate_FindByIDError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	req := domain.UpdateUserRequest{FullName: ptrString("Name")}

	repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))

	_, err := uc.Update(ctx, id, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdate_RepoUpdateError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	existing := dummyUser(id)
	req := domain.UpdateUserRequest{FullName: ptrString("Name")}

	repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
	repo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Update(ctx, id, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdate_RepoUpdateWithRoleError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	existing := dummyUser(id)
	req := domain.UpdateUserRequest{
		FullName: ptrString("Name"),
		Role:     ptrString("admin"),
	}

	repo.EXPECT().FindByID(ctx, id).Return(existing, nil)
	repo.EXPECT().UpdateWithRole(ctx, gomock.Any(), "admin").Return(errors.New("db error"))

	_, err := uc.Update(ctx, id, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Delete
// ============================================================

func TestDelete_Success(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()
	actorID := uuid.New()

	user := dummyUser(id)
	repo.EXPECT().FindByID(ctx, id).Return(user, nil)
	repo.EXPECT().Delete(ctx, id).Return(nil)

	err := uc.Delete(ctx, id, actorID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete_CannotDeleteSelf(t *testing.T) {
	_, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()

	err := uc.Delete(ctx, id, id) // actorID == id
	if !errors.Is(err, ErrCannotDeleteSelf) {
		t.Errorf("expected ErrCannotDeleteSelf, got %v", err)
	}
}

func TestDelete_UserNotFound(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()
	actorID := uuid.New()

	repo.EXPECT().FindByID(ctx, id).Return(nil, nil)

	err := uc.Delete(ctx, id, actorID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDelete_FindByIDError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()
	actorID := uuid.New()

	repo.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))

	err := uc.Delete(ctx, id, actorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_RepoDeleteError(t *testing.T) {
	repo, uc := setupUseCase(t)
	ctx := context.Background()
	id := uuid.New()
	actorID := uuid.New()

	user := dummyUser(id)
	repo.EXPECT().FindByID(ctx, id).Return(user, nil)
	repo.EXPECT().Delete(ctx, id).Return(errors.New("db error"))

	err := uc.Delete(ctx, id, actorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// toUserResponse (helper)
// ============================================================

func TestToUserResponse_WithRole(t *testing.T) {
	user := dummyUser(uuid.New())
	resp := toUserResponse(user)

	if resp.Role != "admin" {
		t.Errorf("expected role admin, got %s", resp.Role)
	}
	if resp.ID != user.ID {
		t.Errorf("expected id %s, got %s", user.ID, resp.ID)
	}
	if resp.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, resp.Email)
	}
}

func TestToUserResponse_WithoutRole(t *testing.T) {
	user := dummyUser(uuid.New())
	user.Role = nil

	resp := toUserResponse(user)
	if resp.Role != "" {
		t.Errorf("expected empty role, got %s", resp.Role)
	}
}
