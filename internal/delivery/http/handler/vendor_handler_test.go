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
)

type stubVendorUseCase struct {
	setRegistrationPass func(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPasswordRequest) (*domain.VendorLoginResponse, error)
	getMe               func(ctx context.Context, vendorID uuid.UUID) (*domain.VendorMeResponse, error)
	submitProposal      func(ctx context.Context, vendorID uuid.UUID, userID uuid.UUID, req domain.VendorSouvenirStoreProposalRequest) (*domain.VendorSouvenirStoreProposalResponse, error)
	presignProposalDocs func(ctx context.Context, vendorID uuid.UUID, userID uuid.UUID, req domain.VendorSouvenirStoreProposalPresignRequest) (*domain.VendorSouvenirStoreProposalPresignResponse, error)
}

func (s stubVendorUseCase) RequestRegistrationOTP(context.Context, domain.VendorRegisterOTPRequest) (*domain.VendorRegisterOTPResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) VerifyRegistrationOTP(context.Context, domain.VendorVerifyRegistrationOTPRequest) (*domain.VendorVerifyRegistrationOTPResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) SetRegistrationPassword(ctx context.Context, onboardingID uuid.UUID, req domain.VendorRegistrationPasswordRequest) (*domain.VendorLoginResponse, error) {
	if s.setRegistrationPass == nil {
		return nil, nil
	}
	return s.setRegistrationPass(ctx, onboardingID, req)
}

func (s stubVendorUseCase) Login(context.Context, domain.VendorLoginRequest) (*domain.VendorLoginResponse, error) {
	return nil, nil
}

func (s stubVendorUseCase) GetMe(ctx context.Context, vendorID uuid.UUID) (*domain.VendorMeResponse, error) {
	if s.getMe == nil {
		return nil, nil
	}
	return s.getMe(ctx, vendorID)
}

func (s stubVendorUseCase) SubmitSouvenirStoreProposal(ctx context.Context, vendorID uuid.UUID, userID uuid.UUID, req domain.VendorSouvenirStoreProposalRequest) (*domain.VendorSouvenirStoreProposalResponse, error) {
	if s.submitProposal == nil {
		return nil, nil
	}
	return s.submitProposal(ctx, vendorID, userID, req)
}

func (s stubVendorUseCase) PresignSouvenirStoreProposalDocuments(ctx context.Context, vendorID uuid.UUID, userID uuid.UUID, req domain.VendorSouvenirStoreProposalPresignRequest) (*domain.VendorSouvenirStoreProposalPresignResponse, error) {
	if s.presignProposalDocs == nil {
		return nil, nil
	}
	return s.presignProposalDocs(ctx, vendorID, userID, req)
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

type vendorLoginResponsePayload struct {
	AccessToken  string  `json:"access_token"`
	VendorID     string  `json:"vendor_id"`
	Email        string  `json:"email"`
	ImageURL     *string `json:"image_url"`
	StoreName    *string `json:"store_name"`
	VendorType   *string `json:"vendor_type"`
	VendorStatus string  `json:"vendor_status"`
}

func TestVendorHandlerSetRegistrationPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	onboardingID := uuid.New()
	storeName := "Toko Oleh-Oleh Haji Mabrur"
	vendorID := uuid.New()
	h := NewVendorHandler(stubVendorUseCase{
		setRegistrationPass: func(ctx context.Context, gotOnboardingID uuid.UUID, req domain.VendorRegistrationPasswordRequest) (*domain.VendorLoginResponse, error) {
			if gotOnboardingID != onboardingID {
				t.Fatalf("expected onboarding ID %s, got %s", onboardingID, gotOnboardingID)
			}
			if req.Password != "password123" {
				t.Fatalf("unexpected password: %s", req.Password)
			}
			return &domain.VendorLoginResponse{
				AccessToken:  "token",
				VendorID:     vendorID,
				Email:        "vendor@example.com",
				ImageURL:     nil,
				StoreName:    &storeName,
				VendorType:   nil,
				VendorStatus: domain.VendorStatusDraft,
			}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vendors/register/password", strings.NewReader(`{"password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Set(middleware.ContextKeyVendorOnboardingID, onboardingID)

	h.SetRegistrationPassword(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var env handlerEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if env.Message != "vendor registration completed" {
		t.Fatalf("unexpected message: %s", env.Message)
	}

	var data vendorLoginResponsePayload
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("failed to decode response data: %v", err)
	}
	if data.AccessToken != "token" {
		t.Fatalf("unexpected access token: %s", data.AccessToken)
	}
	if data.VendorID != vendorID.String() {
		t.Fatalf("unexpected vendor ID: %s", data.VendorID)
	}
	if data.Email != "vendor@example.com" {
		t.Fatalf("unexpected email: %s", data.Email)
	}
	if data.ImageURL != nil {
		t.Fatalf("expected nil image URL, got %v", *data.ImageURL)
	}
	if data.StoreName == nil || *data.StoreName != storeName {
		t.Fatalf("unexpected store name: %v", data.StoreName)
	}
	if data.VendorType != nil {
		t.Fatalf("expected nil vendor type, got %v", data.VendorType)
	}
	if data.VendorStatus != domain.VendorStatusDraft {
		t.Fatalf("unexpected vendor status: %s", data.VendorStatus)
	}
}

func TestVendorHandlerGetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		vendorID := uuid.New()
		imageURL := "https://cdn.example.com/vendors/profile.jpg"
		storeName := "Toko Oleh-Oleh Haji Mabrur"
		vendorType := domain.VendorTypeSouvenirStore
		h := NewVendorHandler(stubVendorUseCase{
			getMe: func(ctx context.Context, gotVendorID uuid.UUID) (*domain.VendorMeResponse, error) {
				if gotVendorID != vendorID {
					t.Fatalf("expected vendor ID %s, got %s", vendorID, gotVendorID)
				}
				return &domain.VendorMeResponse{
					VendorID:     vendorID,
					ImageURL:     &imageURL,
					Email:        "toko.mabrur@example.com",
					VendorType:   &vendorType,
					VendorStatus: domain.VendorStatusActive,
					StoreName:    &storeName,
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
		if env.Message != "vendor profile retrieved successfully" {
			t.Fatalf("unexpected message: %s", env.Message)
		}

		var data vendorLoginResponsePayload
		if err := json.Unmarshal(env.Data, &data); err != nil {
			t.Fatalf("failed to decode response data: %v", err)
		}
		if data.AccessToken != "" {
			t.Fatalf("expected empty access token, got %s", data.AccessToken)
		}
		if data.VendorID != vendorID.String() {
			t.Fatalf("unexpected vendor ID: %s", data.VendorID)
		}
		if data.Email != "toko.mabrur@example.com" {
			t.Fatalf("unexpected email: %s", data.Email)
		}
		if data.ImageURL == nil || *data.ImageURL != imageURL {
			t.Fatalf("unexpected image URL: %v", data.ImageURL)
		}
		if data.StoreName == nil || *data.StoreName != storeName {
			t.Fatalf("unexpected store name: %v", data.StoreName)
		}
		if data.VendorType == nil || *data.VendorType != vendorType {
			t.Fatalf("unexpected vendor type: %v", data.VendorType)
		}
		if data.VendorStatus != domain.VendorStatusActive {
			t.Fatalf("unexpected vendor status: %s", data.VendorStatus)
		}
	})
}

func TestVendorHandlerSubmitSouvenirStoreProposal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	vendorID := uuid.New()
	userID := uuid.New()
	h := NewVendorHandler(stubVendorUseCase{
		submitProposal: func(ctx context.Context, gotVendorID uuid.UUID, gotUserID uuid.UUID, req domain.VendorSouvenirStoreProposalRequest) (*domain.VendorSouvenirStoreProposalResponse, error) {
			if gotVendorID != vendorID {
				t.Fatalf("expected vendor ID %s, got %s", vendorID, gotVendorID)
			}
			if gotUserID != userID {
				t.Fatalf("expected user ID %s, got %s", userID, gotUserID)
			}
			if req.Address.SubdistrictID != "3171011001" {
				t.Fatalf("unexpected subdistrict ID: %s", req.Address.SubdistrictID)
			}
			if req.ResponsiblePerson.Email != "responsible@example.com" {
				t.Fatalf("unexpected responsible person email: %s", req.ResponsiblePerson.Email)
			}
			return &domain.VendorSouvenirStoreProposalResponse{
				VendorID: vendorID,
				Status:   domain.VendorStatusSubmitted,
			}, nil
		},
	})

	body := `{
		"store_name":"Toko Haji",
		"store_description":"Pusat oleh-oleh haji dan umroh",
		"address":{
			"province_id":"31",
			"city_id":"3171",
			"district_id":"317101",
			"subdistrict_id":"3171011001",
			"postal_code":"10110",
			"address_line":"Jl. KH. Wahid Hasyim No. 10"
		},
		"responsible_person":{
			"name":"Ahmad",
			"phone":"08123456789",
			"email":"responsible@example.com",
			"nik":"3173010101010001",
			"ktp_object_id":"vendors/documents/ktp"
		},
		"other_documents":{
			"nib_object_id":"vendors/documents/nib",
			"halal_certificate_object_id":"vendors/documents/halal"
		}
	}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vendors/register/propose/souvenir_store", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Set(middleware.ContextKeyVendorID, vendorID)
	c.Set(middleware.ContextKeyUserID, userID)

	h.SubmitSouvenirStoreProposal(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var env handlerEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if env.Message != "vendor proposal submitted successfully" {
		t.Fatalf("unexpected message: %s", env.Message)
	}
}

func TestVendorHandlerPresignSouvenirStoreProposalDocuments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	vendorID := uuid.New()
	userID := uuid.New()
	h := NewVendorHandler(stubVendorUseCase{
		presignProposalDocs: func(ctx context.Context, gotVendorID uuid.UUID, gotUserID uuid.UUID, req domain.VendorSouvenirStoreProposalPresignRequest) (*domain.VendorSouvenirStoreProposalPresignResponse, error) {
			if gotVendorID != vendorID {
				t.Fatalf("expected vendor ID %s, got %s", vendorID, gotVendorID)
			}
			if gotUserID != userID {
				t.Fatalf("expected user ID %s, got %s", userID, gotUserID)
			}
			if req.ResponsiblePersonKTP != "image/png" {
				t.Fatalf("unexpected responsible person ktp content type: %s", req.ResponsiblePersonKTP)
			}
			return &domain.VendorSouvenirStoreProposalPresignResponse{
				ResponsiblePersonKTP: domain.VendorSouvenirStoreProposalPresignItem{
					ObjectID:     "vendors/a/documents/owner_document_id/1",
					PresignedURL: "https://presigned.example.com/ktp",
				},
				NIB: domain.VendorSouvenirStoreProposalPresignItem{
					ObjectID:     "vendors/a/documents/business_nib/2",
					PresignedURL: "https://presigned.example.com/nib",
				},
				HalalCertificate: domain.VendorSouvenirStoreProposalPresignItem{
					ObjectID:     "vendors/a/documents/halal_certificate/3",
					PresignedURL: "https://presigned.example.com/halal",
				},
			}, nil
		},
	})

	body := `{
		"responsible_person_ktp":"image/png",
		"nib":"image/jpeg",
		"halal_certificate":"application/pdf"
	}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vendors/register/propose/souvenir_store/presign", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Set(middleware.ContextKeyVendorID, vendorID)
	c.Set(middleware.ContextKeyUserID, userID)

	h.PresignSouvenirStoreProposalDocuments(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var env handlerEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if env.Message != "presigned upload URL generated" {
		t.Fatalf("unexpected message: %s", env.Message)
	}
}
