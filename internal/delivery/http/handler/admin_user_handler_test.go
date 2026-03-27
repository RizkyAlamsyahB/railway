package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
)

type stubAdminUserUseCase struct {
	create func(context.Context, domain.CreateUserRequest) (*domain.UserResponse, error)
	getMe  func(ctx context.Context, adminID uuid.UUID) (*domain.AdminMeResponse, error)
}

func (s stubAdminUserUseCase) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	if s.create == nil {
		return nil, nil
	}
	return s.create(ctx, req)
}

func (s stubAdminUserUseCase) List(context.Context, domain.UserListParams) ([]domain.UserResponse, *domain.PaginationMeta, error) {
	return nil, nil, nil
}

func (s stubAdminUserUseCase) GetByID(context.Context, uuid.UUID) (*domain.UserResponse, error) {
	return nil, nil
}

func (s stubAdminUserUseCase) GetMe(ctx context.Context, adminID uuid.UUID) (*domain.AdminMeResponse, error) {
	if s.getMe == nil {
		return nil, nil
	}
	return s.getMe(ctx, adminID)
}

func (s stubAdminUserUseCase) Update(context.Context, uuid.UUID, domain.UpdateUserRequest) (*domain.UserResponse, error) {
	return nil, nil
}

func (s stubAdminUserUseCase) Delete(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

func TestAdminUserHandlerGetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		adminID := uuid.New()
		h := NewAdminUserHandler(stubAdminUserUseCase{
			getMe: func(_ context.Context, gotAdminID uuid.UUID) (*domain.AdminMeResponse, error) {
				if gotAdminID != adminID {
					t.Fatalf("expected admin ID %s, got %s", adminID, gotAdminID)
				}
				return &domain.AdminMeResponse{
					ID:       adminID,
					ImageURL: "https://cdn.example.com/admins/profile.jpg",
					Email:    "admin@example.com",
					Name:     "System Administrator",
				}, nil
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/me", nil)
		c.Request = req
		c.Set(middleware.ContextKeyUserID, adminID)

		h.GetMe(c)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !env.Success {
			t.Fatal("expected success=true")
		}
		if env.Message != "admin profile retrieved successfully" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var data domain.AdminMeResponse
		if err := json.Unmarshal(env.Data, &data); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if data.ID != adminID {
			t.Fatalf("expected ID %s, got %s", adminID, data.ID)
		}
		if data.ImageURL != "https://cdn.example.com/admins/profile.jpg" {
			t.Fatalf("unexpected image URL: %s", data.ImageURL)
		}
		if data.Email != "admin@example.com" {
			t.Fatalf("unexpected email: %s", data.Email)
		}
		if data.Name != "System Administrator" {
			t.Fatalf("unexpected name: %s", data.Name)
		}
	})

	t.Run("invalid user id in context", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/me", nil)
		c.Request = req
		c.Set(middleware.ContextKeyUserID, 123)

		h.GetMe(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "invalid user ID in token" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})

	t.Run("usecase error is mapped", func(t *testing.T) {
		adminID := uuid.New()
		h := NewAdminUserHandler(stubAdminUserUseCase{
			getMe: func(context.Context, uuid.UUID) (*domain.AdminMeResponse, error) {
				return nil, usecase.ErrUserNotFound
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/me", nil)
		c.Request = req
		c.Set(middleware.ContextKeyUserID, adminID)

		h.GetMe(c)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "user not found" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})
}

type adminUserResponseEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Errors  json.RawMessage `json:"errors"`
	Meta    json.RawMessage `json:"meta"`
}

type adminUserRequestErrorEnvelope struct {
	Fields  map[string][]string `json:"fields"`
	Request []string            `json:"request"`
}

func TestAdminUserHandlerCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("validation error returns structured field errors", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/users",
			strings.NewReader(`{"email":"not-an-email","password":"short","role":"owner"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env adminUserResponseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "validation failed" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var errs adminUserRequestErrorEnvelope
		if err := json.Unmarshal(env.Errors, &errs); err != nil {
			t.Fatalf("failed to decode errors: %v", err)
		}
		if len(errs.Request) != 0 {
			t.Fatalf("expected empty request errors, got %v", errs.Request)
		}
		if got := errs.Fields["email"]; len(got) != 1 || got[0] != "must be a valid email address" {
			t.Fatalf("unexpected email errors: %v", got)
		}
		if got := errs.Fields["full_name"]; len(got) != 1 || got[0] != "is required" {
			t.Fatalf("unexpected full_name errors: %v", got)
		}
		if got := errs.Fields["password"]; len(got) != 1 || got[0] != "must be at least 8 characters" {
			t.Fatalf("unexpected password errors: %v", got)
		}
		if got := errs.Fields["role"]; len(got) != 1 || got[0] != "must be one of: admin, umkm, customer, cs, finance" {
			t.Fatalf("unexpected role errors: %v", got)
		}
	})

	t.Run("malformed json returns request error", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/users",
			strings.NewReader(`{"email":"admin@example.com"`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env adminUserResponseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		var errs adminUserRequestErrorEnvelope
		if err := json.Unmarshal(env.Errors, &errs); err != nil {
			t.Fatalf("failed to decode errors: %v", err)
		}
		if len(errs.Fields) != 0 {
			t.Fatalf("expected empty field errors, got %v", errs.Fields)
		}
		if len(errs.Request) != 1 || errs.Request[0] != "request body must be valid JSON" {
			t.Fatalf("unexpected request errors: %v", errs.Request)
		}
	})

	t.Run("email conflict returns field error", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{
			create: func(_ context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
				if req.Email != "admin@example.com" {
					t.Fatalf("unexpected email: %s", req.Email)
				}
				return nil, usecase.ErrEmailExists
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/users",
			strings.NewReader(`{"email":"admin@example.com","full_name":"Admin","password":"password123","role":"admin"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected status 409, got %d", w.Code)
		}

		var env adminUserResponseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "email already exists" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var errs adminUserRequestErrorEnvelope
		if err := json.Unmarshal(env.Errors, &errs); err != nil {
			t.Fatalf("failed to decode errors: %v", err)
		}
		if got := errs.Fields["email"]; len(got) != 1 || got[0] != "email is already registered" {
			t.Fatalf("unexpected email errors: %v", got)
		}
		if len(errs.Request) != 0 {
			t.Fatalf("expected empty request errors, got %v", errs.Request)
		}
	})

	t.Run("phone conflict returns field error", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{
			create: func(_ context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
				if req.Phone == nil || *req.Phone != "08123456789" {
					t.Fatalf("unexpected phone: %v", req.Phone)
				}
				return nil, usecase.ErrPhoneAlreadyRegistered
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/users",
			strings.NewReader(`{"email":"admin@example.com","full_name":"Admin","password":"password123","phone":"08123456789","role":"admin"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		if w.Code != http.StatusConflict {
			t.Fatalf("expected status 409, got %d", w.Code)
		}

		var env adminUserResponseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "phone number already registered" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var errs adminUserRequestErrorEnvelope
		if err := json.Unmarshal(env.Errors, &errs); err != nil {
			t.Fatalf("failed to decode errors: %v", err)
		}
		if got := errs.Fields["phone"]; len(got) != 1 || got[0] != "phone number is already registered" {
			t.Fatalf("unexpected phone errors: %v", got)
		}
	})

	t.Run("invalid birth date returns field error", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{
			create: func(context.Context, domain.CreateUserRequest) (*domain.UserResponse, error) {
				return nil, usecase.ErrInvalidBirthDate
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/users",
			strings.NewReader(`{"email":"admin@example.com","full_name":"Admin","password":"password123","birth_date":"17-03-2026","role":"admin"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env adminUserResponseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "validation failed" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var errs adminUserRequestErrorEnvelope
		if err := json.Unmarshal(env.Errors, &errs); err != nil {
			t.Fatalf("failed to decode errors: %v", err)
		}
		if got := errs.Fields["birth_date"]; len(got) != 1 || got[0] != "must use YYYY-MM-DD format" {
			t.Fatalf("unexpected birth_date errors: %v", got)
		}
	})

	t.Run("unexpected usecase error is sanitized", func(t *testing.T) {
		h := NewAdminUserHandler(stubAdminUserUseCase{
			create: func(context.Context, domain.CreateUserRequest) (*domain.UserResponse, error) {
				return nil, context.DeadlineExceeded
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/users",
			strings.NewReader(`{"email":"admin@example.com","full_name":"Admin","password":"password123","role":"admin"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", w.Code)
		}

		var env adminUserResponseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "internal server error" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var errs adminUserRequestErrorEnvelope
		if err := json.Unmarshal(env.Errors, &errs); err != nil {
			t.Fatalf("failed to decode errors: %v", err)
		}
		if len(errs.Fields) != 0 || len(errs.Request) != 0 {
			t.Fatalf("expected empty structured errors, got fields=%v request=%v", errs.Fields, errs.Request)
		}
	})
}
