package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/delivery/http/middleware"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
)

type stubVendorUseCase struct {
	getRegistrationStatus                 func(ctx context.Context, onboardingID uuid.UUID) (*domain.VendorOnboardingStatusResponse, error)
	getMe                                 func(ctx context.Context, vendorID uuid.UUID) (*domain.VendorProfileResponse, error)
	presignRegistrationIndividualDocument func(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPresignDocumentRequest) (*domain.VendorRegistrationPresignDocumentResponse, error)
}

func (s stubVendorUseCase) RequestRegistrationOTP(context.Context, domain.VendorRegisterOTPRequest) (*domain.VendorRegisterOTPResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) VerifyRegistrationOTP(context.Context, domain.VendorVerifyRegistrationOTPRequest) (*domain.VendorVerifyRegistrationOTPResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) GetRegistrationStatus(ctx context.Context, onboardingID uuid.UUID) (*domain.VendorOnboardingStatusResponse, error) {
	if s.getRegistrationStatus == nil {
		return nil, nil
	}
	return s.getRegistrationStatus(ctx, onboardingID)
}

func (s stubVendorUseCase) SetRegistrationPassword(context.Context, uuid.UUID, domain.VendorRegistrationPasswordRequest) (*domain.VendorOnboardingProgressResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) PresignRegistrationIndividualDocument(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPresignDocumentRequest) (*domain.VendorRegistrationPresignDocumentResponse, error) {
	if s.presignRegistrationIndividualDocument == nil {
		return nil, nil
	}
	return s.presignRegistrationIndividualDocument(ctx, onboardingID, req)
}

func (s stubVendorUseCase) PresignRegistrationCorporateDocument(context.Context, uuid.UUID) (*domain.VendorRegistrationPresignDocumentResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) SubmitRegistrationIndividual(context.Context, uuid.UUID, domain.VendorRegistrationIndividualLegalRequest) (*domain.VendorLoginResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) SubmitRegistrationCorporate(context.Context, uuid.UUID, domain.VendorRegistrationCorporateLegalRequest) (*domain.VendorLoginResponse, error) {
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

func TestVendorHandlerGetRegistrationStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		onboardingID := uuid.New()
		h := NewVendorHandler(stubVendorUseCase{
			getRegistrationStatus: func(ctx context.Context, gotOnboardingID uuid.UUID) (*domain.VendorOnboardingStatusResponse, error) {
				if gotOnboardingID != onboardingID {
					t.Fatalf("expected onboarding ID %s, got %s", onboardingID, gotOnboardingID)
				}
				return &domain.VendorOnboardingStatusResponse{
					Status: domain.VendorOnboardingStatusPasswordSet,
					Email:  "toko.mabrur@example.com",
				}, nil
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/vendors/register/status", nil)
		c.Request = req
		c.Set(middleware.ContextKeyVendorOnboardingID, onboardingID)

		h.GetRegistrationStatus(c)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "vendor registration status retrieved" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var data domain.VendorOnboardingStatusResponse
		if err := json.Unmarshal(env.Data, &data); err != nil {
			t.Fatalf("failed to decode data: %v", err)
		}
		if data.Status != domain.VendorOnboardingStatusPasswordSet {
			t.Fatalf("unexpected status: %s", data.Status)
		}
		if data.Email != "toko.mabrur@example.com" {
			t.Fatalf("unexpected email: %s", data.Email)
		}
	})

	t.Run("invalid onboarding id in context", func(t *testing.T) {
		h := NewVendorHandler(stubVendorUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/vendors/register/status", nil)
		c.Request = req
		c.Set(middleware.ContextKeyVendorOnboardingID, 123)

		h.GetRegistrationStatus(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "invalid vendor onboarding ID in token" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})

	t.Run("usecase error is mapped", func(t *testing.T) {
		onboardingID := uuid.New()
		h := NewVendorHandler(stubVendorUseCase{
			getRegistrationStatus: func(context.Context, uuid.UUID) (*domain.VendorOnboardingStatusResponse, error) {
				return nil, usecase.ErrVendorOnboardingNotFound
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/vendors/register/status", nil)
		c.Request = req
		c.Set(middleware.ContextKeyVendorOnboardingID, onboardingID)

		h.GetRegistrationStatus(c)

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}

		var env handlerEnvelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if env.Message != "vendor onboarding not found" {
			t.Fatalf("unexpected message: %s", env.Message)
		}
	})
}

func TestVendorHandlerPresignRegistrationIndividualDocument(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		onboardingID := uuid.New()
		h := NewVendorHandler(stubVendorUseCase{
			presignRegistrationIndividualDocument: func(ctx context.Context, gotOnboardingID uuid.UUID, req domain.VendorRegistrationPresignDocumentRequest) (*domain.VendorRegistrationPresignDocumentResponse, error) {
				if gotOnboardingID != onboardingID {
					t.Fatalf("expected onboarding ID %s, got %s", onboardingID, gotOnboardingID)
				}
				if req.DocumentIDType != domain.VendorDocumentIDTypeKTP {
					t.Fatalf("unexpected document id type: %s", req.DocumentIDType)
				}
				if req.ContentType != "image/jpeg" {
					t.Fatalf("unexpected content type: %s", req.ContentType)
				}
				return &domain.VendorRegistrationPresignDocumentResponse{
					UploadURL:       "https://upload.example.com",
					ObjectKey:       "vendor-onboardings/x/documents/owner_document_id/file",
					ExpiresAt:       time.Now().Add(time.Minute),
					DocumentIDType:  req.DocumentIDType,
					DocumentDocType: domain.VendorDocumentTypeOwnerDocumentID,
				}, nil
			},
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/vendors/register/souvenir-store/individual/presign", strings.NewReader(`{"document_id_type":"ktp","content_type":"image/jpeg"}`))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Set(middleware.ContextKeyVendorOnboardingID, onboardingID)

		h.PresignRegistrationIndividualDocument(c)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("missing content type", func(t *testing.T) {
		h := NewVendorHandler(stubVendorUseCase{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/vendors/register/souvenir-store/individual/presign", strings.NewReader(`{"document_id_type":"ktp"}`))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req
		c.Set(middleware.ContextKeyVendorOnboardingID, uuid.New())

		h.PresignRegistrationIndividualDocument(c)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}
	})
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
