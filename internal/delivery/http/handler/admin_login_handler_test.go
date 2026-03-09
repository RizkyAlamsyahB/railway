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
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
)

type stubAdminAuthUseCase struct {
	login func(context.Context, domain.AdminLoginRequest) (*domain.AdminLoginResponse, error)
}

func (s stubAdminAuthUseCase) Login(ctx context.Context, req domain.AdminLoginRequest) (*domain.AdminLoginResponse, error) {
	if s.login == nil {
		return nil, nil
	}
	return s.login(ctx, req)
}

func TestAdminLoginHandlerLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		adminID := uuid.New()
		h := NewAdminLoginHandler(stubAdminAuthUseCase{
			login: func(_ context.Context, req domain.AdminLoginRequest) (*domain.AdminLoginResponse, error) {
				if req.Email != "admin@example.com" {
					t.Fatalf("expected email admin@example.com, got %s", req.Email)
				}
				if req.Password != "password123" {
					t.Fatalf("expected password password123, got %s", req.Password)
				}
				return &domain.AdminLoginResponse{
					Token:    "jwt-token",
					ID:       adminID,
					ImageURL: "https://cdn.example.com/admins/profile.jpg",
					Email:    "admin@example.com",
					Name:     "Admin Ops",
				}, nil
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/login",
			strings.NewReader(`{"email":"admin@example.com","password":"password123"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

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
		if env.Message != "login successful" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var data domain.AdminLoginResponse
		if err := json.Unmarshal(env.Data, &data); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if data.Token != "jwt-token" {
			t.Fatalf("expected token jwt-token, got %s", data.Token)
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
		if data.Name != "Admin Ops" {
			t.Fatalf("unexpected name: %s", data.Name)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		h := NewAdminLoginHandler(stubAdminAuthUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/login",
			strings.NewReader(`{"email":"not-an-email","password":"password123"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env responseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "validation failed" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
		if env.Data != nil {
			t.Fatal("expected nil data on validation error")
		}
	})

	t.Run("usecase error is mapped", func(t *testing.T) {
		h := NewAdminLoginHandler(stubAdminAuthUseCase{
			login: func(context.Context, domain.AdminLoginRequest) (*domain.AdminLoginResponse, error) {
				return nil, usecase.ErrInvalidCredentials
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(
			http.MethodPost,
			"/api/v1/admin/login",
			strings.NewReader(`{"email":"admin@example.com","password":"wrongpassword"}`),
		)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}

		var env responseEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "invalid email or password" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
		if env.Data != nil {
			t.Fatal("expected nil data on error response")
		}
	})
}

type responseEnvelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
	Meta    interface{} `json:"meta"`
}
