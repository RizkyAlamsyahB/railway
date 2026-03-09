package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
)

type stubAdminUserUseCase struct {
	getMe func(ctx context.Context, adminID uuid.UUID) (*domain.AdminMeResponse, error)
}

func (s stubAdminUserUseCase) Create(context.Context, domain.CreateUserRequest) (*domain.UserResponse, error) {
	return nil, nil
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
