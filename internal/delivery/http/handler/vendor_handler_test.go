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

type stubVendorUseCase struct {
	getMe func(ctx context.Context, vendorID uuid.UUID) (*domain.VendorProfileResponse, error)
}

func (s stubVendorUseCase) Register(context.Context, domain.VendorRegisterRequest) (*domain.VendorRegisterResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) ConfirmDocuments(context.Context, uuid.UUID, domain.ConfirmDocumentsRequest) (*domain.ConfirmDocumentsResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) Login(context.Context, domain.VendorLoginRequest) (*domain.VendorLoginResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) GetMe(ctx context.Context, vendorID uuid.UUID) (*domain.VendorProfileResponse, error) {
	if s.getMe == nil {
		return nil, nil
	}
	return s.getMe(ctx, vendorID)
}

func (s stubVendorUseCase) GetBalance(context.Context, uuid.UUID) (*domain.VendorBalanceResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) ListPayoutChannels(context.Context, uuid.UUID) (*domain.VendorPayoutChannelsResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) RequestWithdrawal(context.Context, uuid.UUID, domain.VendorWithdrawRequest) (*domain.VendorWithdrawResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) HandlePayoutWebhook(context.Context, domain.XenditPayoutWebhookPayload) error {
	return nil
}

type handlerEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestVendorHandlerGetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		vendorID := uuid.New()
		h := NewVendorHandler(stubVendorUseCase{
			getMe: func(ctx context.Context, gotVendorID uuid.UUID) (*domain.VendorProfileResponse, error) {
				if gotVendorID != vendorID {
					t.Fatalf("expected vendor ID %s, got %s", vendorID, gotVendorID)
				}
				return &domain.VendorProfileResponse{
					VendorID:     vendorID,
					ImageURL:     "https://cdn.example.com/vendors/profile.jpg",
					Email:        "toko.mabrur@example.com",
					Name:         "Abu Bakar Shidiq Basalamah",
					VendorType:   domain.VendorTypeSouvenirStore,
					VendorStatus: domain.VendorStatusActive,
					DisplayName:  "Toko Oleh-Oleh Haji Mabrur",
				}, nil
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/vendors/me", nil)
		c.Request = req
		c.Set(middleware.ContextKeyVendorID, vendorID)

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
		if env.Message != "vendor profile retrieved successfully" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var data domain.VendorProfileResponse
		if err := json.Unmarshal(env.Data, &data); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if data.Name != "Abu Bakar Shidiq Basalamah" {
			t.Fatalf("expected owner full name, got %s", data.Name)
		}
		if data.DisplayName != "Toko Oleh-Oleh Haji Mabrur" {
			t.Fatalf("expected display name, got %s", data.DisplayName)
		}
	})

	t.Run("invalid vendor id in context", func(t *testing.T) {
		h := NewVendorHandler(stubVendorUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/vendors/me", nil)
		c.Request = req
		c.Set(middleware.ContextKeyVendorID, 123)

		h.GetMe(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "invalid vendor ID in token" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})

	t.Run("usecase error is mapped", func(t *testing.T) {
		vendorID := uuid.New()
		h := NewVendorHandler(stubVendorUseCase{
			getMe: func(context.Context, uuid.UUID) (*domain.VendorProfileResponse, error) {
				return nil, usecase.ErrVendorNotFound
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/vendors/me", nil)
		c.Request = req
		c.Set(middleware.ContextKeyVendorID, vendorID)

		h.GetMe(c)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "vendor not found" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})
}
