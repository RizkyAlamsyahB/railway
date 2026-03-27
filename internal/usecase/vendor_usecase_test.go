package usecase

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
	"go.uber.org/mock/gomock"
)

type stubOTPUseCase struct {
	requestOTP  func(context.Context, domain.RequestOTPInput) (*domain.RequestOTPResult, error)
	verifyOTP   func(context.Context, domain.VerifyOTPInput) (*domain.VerifyOTPResult, error)
	verifyProof func(context.Context, string, string) (*domain.OTPProofClaims, error)
}

func (s *stubOTPUseCase) RequestOTP(ctx context.Context, input domain.RequestOTPInput) (*domain.RequestOTPResult, error) {
	if s.requestOTP == nil {
		return nil, nil
	}
	return s.requestOTP(ctx, input)
}

func (s *stubOTPUseCase) VerifyOTP(ctx context.Context, input domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) {
	if s.verifyOTP == nil {
		return nil, nil
	}
	return s.verifyOTP(ctx, input)
}

func (s *stubOTPUseCase) VerifyProofToken(ctx context.Context, token string, expectedPurpose string) (*domain.OTPProofClaims, error) {
	if s.verifyProof == nil {
		return nil, nil
	}
	return s.verifyProof(ctx, token, expectedPurpose)
}

type stubVendorOnboardingRepo struct {
	findByID             func(context.Context, uuid.UUID) (*domain.VendorOnboarding, error)
	findByEmail          func(context.Context, string) (*domain.VendorOnboarding, error)
	upsert               func(context.Context, *domain.VendorOnboarding) error
	update               func(context.Context, *domain.VendorOnboarding, ...string) error
	markCompleted        func(context.Context, uuid.UUID, time.Time) error
	finalizeRegistration func(context.Context, domain.VendorRegistrationFinalizeInput) error
}

func (s *stubVendorOnboardingRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
	if s.findByID == nil {
		return nil, nil
	}
	return s.findByID(ctx, id)
}

func (s *stubVendorOnboardingRepo) FindByEmail(ctx context.Context, email string) (*domain.VendorOnboarding, error) {
	if s.findByEmail == nil {
		return nil, nil
	}
	return s.findByEmail(ctx, email)
}

func (s *stubVendorOnboardingRepo) Upsert(ctx context.Context, onboarding *domain.VendorOnboarding) error {
	if s.upsert == nil {
		return nil
	}
	return s.upsert(ctx, onboarding)
}

func (s *stubVendorOnboardingRepo) Update(ctx context.Context, onboarding *domain.VendorOnboarding, fields ...string) error {
	if s.update == nil {
		return nil
	}
	return s.update(ctx, onboarding, fields...)
}

func (s *stubVendorOnboardingRepo) MarkCompleted(ctx context.Context, id uuid.UUID, completedAt time.Time) error {
	if s.markCompleted == nil {
		return nil
	}
	return s.markCompleted(ctx, id, completedAt)
}

func (s *stubVendorOnboardingRepo) FinalizeRegistration(ctx context.Context, input domain.VendorRegistrationFinalizeInput) error {
	if s.finalizeRegistration == nil {
		return nil
	}
	return s.finalizeRegistration(ctx, input)
}

func setupVendorUseCase(t *testing.T) (
	*mocks.MockUserRepository,
	*mocks.MockVendorRepository,
	*mocks.MockStorageProvider,
	*stubOTPUseCase,
	*stubVendorOnboardingRepo,
	domain.VendorUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	otpUC := &stubOTPUseCase{}
	onboardingRepo := &stubVendorOnboardingRepo{}
	uc := NewVendorUseCase(otpUC, userRepo, vendorRepo, onboardingRepo, storage, nil, "test-secret", 3600, "test-issuer")
	return userRepo, vendorRepo, storage, otpUC, onboardingRepo, uc
}

func TestRequestRegistrationOTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo, _, _, otpUC, _, uc := setupVendorUseCase(t)
		ctx := context.Background()

		userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(nil, nil)
		otpUC.requestOTP = func(_ context.Context, input domain.RequestOTPInput) (*domain.RequestOTPResult, error) {
			if input.Email != "vendor@example.com" {
				t.Fatalf("unexpected email: %s", input.Email)
			}
			if input.Purpose != domain.OTPPurposeEmailVerification {
				t.Fatalf("unexpected purpose: %s", input.Purpose)
			}
			return &domain.RequestOTPResult{
				ExpiresAt:     time.Now().Add(5 * time.Minute),
				CooldownUntil: time.Now().Add(1 * time.Minute),
			}, nil
		}

		resp, err := uc.RequestRegistrationOTP(ctx, domain.VendorRegisterOTPRequest{Email: "Vendor@Example.com"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.ExpiresAt.IsZero() || resp.CooldownUntil.IsZero() {
			t.Fatal("expected otp response with timestamps")
		}
	})

	t.Run("vendor already exists", func(t *testing.T) {
		userRepo, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		existingUser := &domain.User{ID: uuid.New(), Email: "vendor@example.com"}

		userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(existingUser, nil)
		vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(&domain.Vendor{ID: uuid.New(), OwnerUserID: existingUser.ID}, nil)

		_, err := uc.RequestRegistrationOTP(ctx, domain.VendorRegisterOTPRequest{Email: "vendor@example.com"})
		if !errors.Is(err, ErrVendorAlreadyExists) {
			t.Fatalf("expected ErrVendorAlreadyExists, got %v", err)
		}
	})
}

func TestVerifyRegistrationOTP(t *testing.T) {
	userRepo, _, _, otpUC, onboardingRepo, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(nil, nil)
	otpUC.verifyOTP = func(_ context.Context, input domain.VerifyOTPInput) (*domain.VerifyOTPResult, error) {
		if input.Email != "vendor@example.com" {
			t.Fatalf("unexpected email: %s", input.Email)
		}
		return &domain.VerifyOTPResult{}, nil
	}
	onboardingRepo.findByEmail = func(_ context.Context, email string) (*domain.VendorOnboarding, error) {
		if email != "vendor@example.com" {
			t.Fatalf("unexpected email lookup: %s", email)
		}
		return nil, nil
	}
	onboardingRepo.upsert = func(_ context.Context, onboarding *domain.VendorOnboarding) error {
		if onboarding.Status != domain.VendorOnboardingStatusOTPVerified {
			t.Fatalf("unexpected onboarding status: %s", onboarding.Status)
		}
		if onboarding.ID == uuid.Nil {
			t.Fatal("expected onboarding ID")
		}
		return nil
	}

	resp, err := uc.VerifyRegistrationOTP(ctx, domain.VendorVerifyRegistrationOTPRequest{
		Email: "vendor@example.com",
		Code:  "123456",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.OnboardingToken == "" {
		t.Fatal("expected onboarding token")
	}
	if resp.Status != domain.VendorOnboardingStatusOTPVerified {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
}

func TestGetRegistrationStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
		ctx := context.Background()
		onboardingID := uuid.New()

		onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
			if id != onboardingID {
				t.Fatalf("unexpected onboarding id: %s", id)
			}
			return &domain.VendorOnboarding{
				ID:     onboardingID,
				Email:  "vendor@example.com",
				Status: domain.VendorOnboardingStatusPasswordSet,
			}, nil
		}

		resp, err := uc.GetRegistrationStatus(ctx, onboardingID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.Email != "vendor@example.com" {
			t.Fatalf("unexpected email: %s", resp.Email)
		}
		if resp.Status != domain.VendorOnboardingStatusPasswordSet {
			t.Fatalf("unexpected status: %s", resp.Status)
		}
	})

	t.Run("completed onboarding remains readable", func(t *testing.T) {
		_, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
		ctx := context.Background()
		onboardingID := uuid.New()
		completedAt := time.Now()

		onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
			if id != onboardingID {
				t.Fatalf("unexpected onboarding id: %s", id)
			}
			return &domain.VendorOnboarding{
				ID:          onboardingID,
				Email:       "vendor@example.com",
				Status:      domain.VendorOnboardingStatusCompleted,
				CompletedAt: &completedAt,
			}, nil
		}

		resp, err := uc.GetRegistrationStatus(ctx, onboardingID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.Status != domain.VendorOnboardingStatusCompleted {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})

	t.Run("onboarding not found", func(t *testing.T) {
		_, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
		ctx := context.Background()
		onboardingID := uuid.New()

		onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
			if id != onboardingID {
				t.Fatalf("unexpected onboarding id: %s", id)
			}
			return nil, nil
		}

		_, err := uc.GetRegistrationStatus(ctx, onboardingID)
		if !errors.Is(err, ErrVendorOnboardingNotFound) {
			t.Fatalf("expected ErrVendorOnboardingNotFound, got %v", err)
		}
	})
}

func TestSetRegistrationPassword_InvalidStep(t *testing.T) {
	_, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
	ctx := context.Background()
	onboardingID := uuid.New()
	onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
		if id != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", id)
		}
		return &domain.VendorOnboarding{
			ID:     onboardingID,
			Email:  "vendor@example.com",
			Status: domain.VendorOnboardingStatusPasswordSet,
		}, nil
	}

	_, err := uc.SetRegistrationPassword(ctx, onboardingID, domain.VendorRegistrationPasswordRequest{Password: "password123"})
	if !errors.Is(err, ErrVendorOnboardingStep) {
		t.Fatalf("expected ErrVendorOnboardingStep, got %v", err)
	}
}

func TestPresignRegistrationIndividualDocument_TableDriven(t *testing.T) {
	t.Run("success jpeg", func(t *testing.T) {
		_, _, storage, _, onboardingRepo, uc := setupVendorUseCase(t)
		ctx := context.Background()
		onboardingID := uuid.New()

		onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
			if id != onboardingID {
				t.Fatalf("unexpected onboarding id: %s", id)
			}
			return &domain.VendorOnboarding{
				ID:     onboardingID,
				Email:  "vendor@example.com",
				Status: domain.VendorOnboardingStatusPasswordSet,
			}, nil
		}
		storage.EXPECT().
			GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).
			Return("https://upload.example.com", nil)

		resp, err := uc.PresignRegistrationIndividualDocument(ctx, onboardingID, domain.VendorRegistrationPresignDocumentRequest{
			DocumentIDType: domain.VendorDocumentIDTypeKTP,
			ContentType:    "image/jpeg",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.UploadURL == "" {
			t.Fatal("expected upload url")
		}
	})

	t.Run("success normalized pdf", func(t *testing.T) {
		_, _, storage, _, onboardingRepo, uc := setupVendorUseCase(t)
		ctx := context.Background()
		onboardingID := uuid.New()

		onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
			if id != onboardingID {
				t.Fatalf("unexpected onboarding id: %s", id)
			}
			return &domain.VendorOnboarding{
				ID:     onboardingID,
				Email:  "vendor@example.com",
				Status: domain.VendorOnboardingStatusPasswordSet,
			}, nil
		}
		storage.EXPECT().
			GeneratePresignedUploadURL(ctx, gomock.Any(), "application/pdf", PresignedUploadExpiry).
			Return("https://upload.example.com", nil)

		resp, err := uc.PresignRegistrationIndividualDocument(ctx, onboardingID, domain.VendorRegistrationPresignDocumentRequest{
			DocumentIDType: domain.VendorDocumentIDTypePassport,
			ContentType:    " APPLICATION/PDF ; charset=utf-8 ",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.UploadURL == "" {
			t.Fatal("expected upload url")
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		_, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
		ctx := context.Background()
		onboardingID := uuid.New()

		onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
			if id != onboardingID {
				t.Fatalf("unexpected onboarding id: %s", id)
			}
			return &domain.VendorOnboarding{
				ID:     onboardingID,
				Email:  "vendor@example.com",
				Status: domain.VendorOnboardingStatusPasswordSet,
			}, nil
		}

		_, err := uc.PresignRegistrationIndividualDocument(ctx, onboardingID, domain.VendorRegistrationPresignDocumentRequest{
			DocumentIDType: domain.VendorDocumentIDTypeKTP,
			ContentType:    "image/png",
		})
		if !errors.Is(err, ErrInvalidDocumentContent) {
			t.Fatalf("expected ErrInvalidDocumentContent, got %v", err)
		}
	})
}

func TestSubmitRegistrationIndividual_Success(t *testing.T) {
	userRepo, _, storage, _, onboardingRepo, uc := setupVendorUseCase(t)
	ctx := context.Background()
	onboardingID := uuid.New()
	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
		if id != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", id)
		}
		return &domain.VendorOnboarding{
			ID:            onboardingID,
			Email:         "vendor@example.com",
			Status:        domain.VendorOnboardingStatusPasswordSet,
			PasswordHash:  &passwordHash,
			OTPVerifiedAt: time.Now(),
		}, nil
	}
	onboardingRepo.update = func(_ context.Context, onboarding *domain.VendorOnboarding, fields ...string) error {
		if onboarding.OwnerName == nil || *onboarding.OwnerName != "Ahmad" {
			t.Fatalf("unexpected owner name: %+v", onboarding.OwnerName)
		}
		if len(fields) == 0 {
			t.Fatal("expected updated fields")
		}
		return nil
	}
	onboardingRepo.finalizeRegistration = func(_ context.Context, input domain.VendorRegistrationFinalizeInput) error {
		if input.OnboardingID != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", input.OnboardingID)
		}
		if input.UserRole != domain.RoleUMKM {
			t.Fatalf("unexpected user role: %s", input.UserRole)
		}
		return nil
	}

	objectKey := buildVendorOnboardingDocumentObjectKey(onboardingID)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		Key:           objectKey,
		ContentType:   "image/jpeg",
		ContentLength: 1024,
	}, nil)
	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(nil, nil)

	resp, err := uc.SubmitRegistrationIndividual(ctx, onboardingID, domain.VendorRegistrationIndividualLegalRequest{
		StoreName:           "Toko Haji",
		DocumentIDType:      domain.VendorDocumentIDTypeKTP,
		NIK:                 "3173000000000001",
		OwnerName:           "Ahmad",
		BirthDate:           "1990-01-02",
		DocumentIDObjectKey: objectKey,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.Token == "" {
		t.Fatal("expected auth token in response")
	}
	if resp.VendorStatus != domain.VendorStatusDraft {
		t.Fatalf("unexpected vendor status: %s", resp.VendorStatus)
	}
}

func TestSubmitRegistrationIndividual_InvalidContentType(t *testing.T) {
	_, _, storage, _, onboardingRepo, uc := setupVendorUseCase(t)
	ctx := context.Background()
	onboardingID := uuid.New()
	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
		if id != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", id)
		}
		return &domain.VendorOnboarding{
			ID:            onboardingID,
			Email:         "vendor@example.com",
			Status:        domain.VendorOnboardingStatusPasswordSet,
			PasswordHash:  &passwordHash,
			OTPVerifiedAt: time.Now(),
		}, nil
	}

	objectKey := buildVendorOnboardingDocumentObjectKey(onboardingID)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		Key:           objectKey,
		ContentType:   "image/png",
		ContentLength: 1024,
	}, nil)

	_, err = uc.SubmitRegistrationIndividual(ctx, onboardingID, domain.VendorRegistrationIndividualLegalRequest{
		StoreName:           "Toko Haji",
		DocumentIDType:      domain.VendorDocumentIDTypeKTP,
		NIK:                 "3173000000000001",
		OwnerName:           "Ahmad",
		BirthDate:           "1990-01-02",
		DocumentIDObjectKey: objectKey,
	})
	if !errors.Is(err, ErrInvalidDocumentContent) {
		t.Fatalf("expected ErrInvalidDocumentContent, got %v", err)
	}
}

func TestConfirmDocuments_TableDriven(t *testing.T) {
	type testCase struct {
		name       string
		setupMocks func(
			ctx context.Context,
			userID uuid.UUID,
			vendorID uuid.UUID,
			vendor *domain.Vendor,
			req domain.ConfirmDocumentsRequest,
			vendorRepo *mocks.MockVendorRepository,
			storage *mocks.MockStorageProvider,
		)
		reqBuilder func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest
		wantErr    error
		wantAnyErr bool
		wantStatus string
		wantDocCnt int
	}

	buildOnboardingDocument := func(vendorID uuid.UUID, docType string) []domain.VendorDocument {
		uploaderID := uuid.New()
		return []domain.VendorDocument{
			{
				ID:                 uuid.New(),
				VendorID:           vendorID,
				DocType:            docType,
				FileURL:            "vendor-onboardings/onb/documents/" + docType + "/file",
				UploadedBy:         &uploaderID,
				VerificationStatus: domain.VerificationStatusPending,
			},
		}
	}

	buildCompletionItems := func(vendorID uuid.UUID, docTypes ...string) []domain.ConfirmDocumentItem {
		items := make([]domain.ConfirmDocumentItem, 0, len(docTypes))
		for _, docType := range docTypes {
			items = append(items, domain.ConfirmDocumentItem{
				DocType:   docType,
				ObjectKey: "vendors/" + vendorID.String() + "/documents/" + docType + "/file",
			})
		}
		return items
	}

	tests := []testCase{
		{
			name: "success all required uploaded",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: buildCompletionItems(
						vendorID,
						domain.VendorDocumentTypeStorePhoto,
						domain.VendorDocumentTypeBankAccountProof,
						domain.VendorDocumentTypeBusinessLogo,
						domain.VendorDocumentTypeBusinessBanner,
					),
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().
					HeadObject(ctx, gomock.Any()).
					Return(&domain.ObjectInfo{ContentType: "image/jpeg", ContentLength: 1024}, nil).
					Times(4)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(&domain.VendorBankAccount{
					ID:       uuid.New(),
					VendorID: vendorID,
				}, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), domain.VendorStatusSubmitted).
					Return(nil)
			},
			wantStatus: domain.VendorStatusSubmitted,
			wantDocCnt: 4,
		},
		{
			name: "success partial upload",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: buildCompletionItems(vendorID, domain.VendorDocumentTypeStorePhoto),
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				req domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().
					HeadObject(ctx, req.Documents[0].ObjectKey).
					Return(&domain.ObjectInfo{ContentType: "image/png", ContentLength: 2048}, nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), "").
					Return(nil)
			},
			wantStatus: domain.VendorStatusDraft,
			wantDocCnt: 1,
		},
		{
			name: "corporate requires npwp",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: buildCompletionItems(
						vendorID,
						domain.VendorDocumentTypeStorePhoto,
						domain.VendorDocumentTypeBankAccountProof,
						domain.VendorDocumentTypeBusinessLogo,
						domain.VendorDocumentTypeBusinessBanner,
					),
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				legalType := domain.VendorBusinessLegalTypeCorporate
				vendor.BusinessLegalType = &legalType
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeBusinessNIB), nil)
				storage.EXPECT().
					HeadObject(ctx, gomock.Any()).
					Return(&domain.ObjectInfo{ContentType: "image/jpeg", ContentLength: 1024}, nil).
					Times(4)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(&domain.VendorBankAccount{
					ID:       uuid.New(),
					VendorID: vendorID,
				}, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), "").
					Return(nil)
			},
			wantStatus: domain.VendorStatusDraft,
			wantDocCnt: 4,
		},
		{
			name: "corporate submitted when nib and npwp complete",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: buildCompletionItems(
						vendorID,
						domain.VendorDocumentTypeStorePhoto,
						domain.VendorDocumentTypeBankAccountProof,
						domain.VendorDocumentTypeBusinessLogo,
						domain.VendorDocumentTypeBusinessBanner,
						domain.VendorDocumentTypeBusinessNPWP,
					),
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				legalType := domain.VendorBusinessLegalTypeCorporate
				vendor.BusinessLegalType = &legalType
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeBusinessNIB), nil)
				storage.EXPECT().
					HeadObject(ctx, gomock.Any()).
					Return(&domain.ObjectInfo{ContentType: "image/jpeg", ContentLength: 1024}, nil).
					Times(5)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(&domain.VendorBankAccount{
					ID:       uuid.New(),
					VendorID: vendorID,
				}, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), domain.VendorStatusSubmitted).
					Return(nil)
			},
			wantStatus: domain.VendorStatusSubmitted,
			wantDocCnt: 5,
		},
		{
			name: "corporate without business nib stays draft",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: buildCompletionItems(
						vendorID,
						domain.VendorDocumentTypeStorePhoto,
						domain.VendorDocumentTypeBankAccountProof,
						domain.VendorDocumentTypeBusinessLogo,
						domain.VendorDocumentTypeBusinessBanner,
						domain.VendorDocumentTypeBusinessNPWP,
					),
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				legalType := domain.VendorBusinessLegalTypeCorporate
				vendor.BusinessLegalType = &legalType
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().
					HeadObject(ctx, gomock.Any()).
					Return(&domain.ObjectInfo{ContentType: "image/jpeg", ContentLength: 1024}, nil).
					Times(5)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(&domain.VendorBankAccount{
					ID:       uuid.New(),
					VendorID: vendorID,
				}, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), "").
					Return(nil)
			},
			wantStatus: domain.VendorStatusDraft,
			wantDocCnt: 5,
		},
		{
			name: "vendor not found",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "key"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				_ uuid.UUID,
				_ *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "find vendor error",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "key"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				_ uuid.UUID,
				_ *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name: "find documents error",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "key"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name: "invalid document type",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: "nonexistent_type", ObjectKey: "key"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				_ domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{}, nil)
			},
			wantErr: ErrInvalidVendorDocumentType,
		},
		{
			name: "object not uploaded",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "vendors/doc/file"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				req domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(nil, nil)
			},
			wantErr: ErrObjectNotUploaded,
		},
		{
			name: "head object error",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "vendors/doc/file"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				req domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(nil, errors.New("s3 error"))
			},
			wantAnyErr: true,
		},
		{
			name: "invalid content type",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "vendors/doc/file"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				req domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(&domain.ObjectInfo{
					ContentType:   "text/html",
					ContentLength: 1024,
				}, nil)
			},
			wantErr: ErrInvalidDocumentContent,
		},
		{
			name: "document size overflow",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "vendors/doc/file"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				req domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(&domain.ObjectInfo{
					ContentType:   "image/jpeg",
					ContentLength: int64(math.MaxInt32) + 1,
				}, nil)
			},
			wantErr: ErrDocumentSizeOverflow,
		},
		{
			name: "confirm repo error",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: domain.VendorDocumentTypeStorePhoto, ObjectKey: "vendors/doc/file"}},
				}
			},
			setupMocks: func(
				ctx context.Context,
				userID uuid.UUID,
				vendorID uuid.UUID,
				vendor *domain.Vendor,
				req domain.ConfirmDocumentsRequest,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(buildOnboardingDocument(vendorID, domain.VendorDocumentTypeOwnerDocumentID), nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(&domain.ObjectInfo{
					ContentType:   "image/jpeg",
					ContentLength: 1024,
				}, nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantAnyErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, vendorRepo, storage, _, _, uc := setupVendorUseCase(t)
			ctx := context.Background()
			userID := uuid.New()
			vendorID := uuid.New()
			vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: domain.VendorStatusDraft}
			req := tc.reqBuilder(vendorID)

			tc.setupMocks(ctx, userID, vendorID, vendor, req, vendorRepo, storage)

			resp, err := uc.ConfirmDocuments(ctx, userID, req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAnyErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp == nil {
				t.Fatal("expected response, got nil")
			}
			if resp.VendorStatus != tc.wantStatus {
				t.Errorf("expected status %q, got %q", tc.wantStatus, resp.VendorStatus)
			}
			if resp.DocumentsCount != tc.wantDocCnt {
				t.Errorf("expected %d documents confirmed, got %d", tc.wantDocCnt, resp.DocumentsCount)
			}
		})
	}
}

func TestSaveBankAccount_TableDriven(t *testing.T) {
	buildDocuments := func(vendorID uuid.UUID, initialDocType string, docTypes ...string) []domain.VendorDocument {
		userID := uuid.New()
		documents := []domain.VendorDocument{
			{
				ID:                 uuid.New(),
				VendorID:           vendorID,
				DocType:            initialDocType,
				FileURL:            "vendor-onboardings/onb/documents/" + initialDocType + "/file",
				UploadedBy:         &userID,
				VerificationStatus: domain.VerificationStatusPending,
			},
		}

		for _, docType := range docTypes {
			documents = append(documents, domain.VendorDocument{
				ID:                 uuid.New(),
				VendorID:           vendorID,
				DocType:            docType,
				FileURL:            "vendors/" + vendorID.String() + "/documents/" + docType + "/file",
				UploadedBy:         &userID,
				VerificationStatus: domain.VerificationStatusPending,
			})
		}

		return documents
	}

	t.Run("create new bank account", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:     vendorID,
			Status: domain.VendorStatusDraft,
		}, nil)
		vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
		vendorRepo.EXPECT().UpsertBankAccount(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, bankAccount *domain.VendorBankAccount) error {
			if bankAccount.VendorID != vendorID {
				t.Fatalf("unexpected vendor id: %s", bankAccount.VendorID)
			}
			if bankAccount.VerificationStatus != domain.VerificationStatusPending {
				t.Fatalf("unexpected verification status: %s", bankAccount.VerificationStatus)
			}
			return nil
		})
		vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(nil, nil)

		resp, err := uc.SaveBankAccount(ctx, vendorID, domain.VendorSaveBankAccountRequest{
			BankName:          "Bank Syariah Indonesia",
			AccountNumber:     "1234567890",
			AccountHolderName: "Ahmad",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.BankAccount.AccountNumber != "1234567890" {
			t.Fatal("expected bank account response")
		}
		if resp.VendorStatus != domain.VendorStatusDraft {
			t.Fatalf("expected draft status, got %s", resp.VendorStatus)
		}
	})

	t.Run("individual submitted when completion already complete", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:     vendorID,
			Status: domain.VendorStatusDraft,
		}, nil)
		vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
		vendorRepo.EXPECT().UpsertBankAccount(ctx, gomock.Any()).Return(nil)
		vendorRepo.EXPECT().
			FindDocumentsByVendorID(ctx, vendorID).
			Return(buildDocuments(
				vendorID,
				domain.VendorDocumentTypeOwnerDocumentID,
				domain.VendorDocumentTypeStorePhoto,
				domain.VendorDocumentTypeBankAccountProof,
				domain.VendorDocumentTypeBusinessLogo,
				domain.VendorDocumentTypeBusinessBanner,
			), nil)
		vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).DoAndReturn(func(_ context.Context, _ uuid.UUID, updates map[string]any) error {
			if updates["status"] != domain.VendorStatusSubmitted {
				t.Fatalf("unexpected status update: %+v", updates)
			}
			if _, ok := updates["updated_at"]; !ok {
				t.Fatalf("missing updated_at in updates: %+v", updates)
			}
			return nil
		})

		resp, err := uc.SaveBankAccount(ctx, vendorID, domain.VendorSaveBankAccountRequest{
			BankName:          "Bank Syariah Indonesia",
			AccountNumber:     "1234567890",
			AccountHolderName: "Ahmad",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.VendorStatus != domain.VendorStatusSubmitted {
			t.Fatalf("expected submitted status, got %s", resp.VendorStatus)
		}
	})

	t.Run("corporate stays draft when npwp missing", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		legalType := domain.VendorBusinessLegalTypeCorporate

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:                vendorID,
			Status:            domain.VendorStatusDraft,
			BusinessLegalType: &legalType,
		}, nil)
		vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
		vendorRepo.EXPECT().UpsertBankAccount(ctx, gomock.Any()).Return(nil)
		vendorRepo.EXPECT().
			FindDocumentsByVendorID(ctx, vendorID).
			Return(buildDocuments(
				vendorID,
				domain.VendorDocumentTypeBusinessNIB,
				domain.VendorDocumentTypeStorePhoto,
				domain.VendorDocumentTypeBankAccountProof,
				domain.VendorDocumentTypeBusinessLogo,
				domain.VendorDocumentTypeBusinessBanner,
			), nil)

		resp, err := uc.SaveBankAccount(ctx, vendorID, domain.VendorSaveBankAccountRequest{
			BankName:          "Bank Syariah Indonesia",
			AccountNumber:     "1234567890",
			AccountHolderName: "Ahmad",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.VendorStatus != domain.VendorStatusDraft {
			t.Fatalf("expected draft status, got %s", resp.VendorStatus)
		}
	})

	t.Run("corporate submitted when npwp complete", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		legalType := domain.VendorBusinessLegalTypeCorporate

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:                vendorID,
			Status:            domain.VendorStatusDraft,
			BusinessLegalType: &legalType,
		}, nil)
		vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
		vendorRepo.EXPECT().UpsertBankAccount(ctx, gomock.Any()).Return(nil)
		vendorRepo.EXPECT().
			FindDocumentsByVendorID(ctx, vendorID).
			Return(buildDocuments(
				vendorID,
				domain.VendorDocumentTypeBusinessNIB,
				domain.VendorDocumentTypeStorePhoto,
				domain.VendorDocumentTypeBankAccountProof,
				domain.VendorDocumentTypeBusinessLogo,
				domain.VendorDocumentTypeBusinessBanner,
				domain.VendorDocumentTypeBusinessNPWP,
			), nil)
		vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)

		resp, err := uc.SaveBankAccount(ctx, vendorID, domain.VendorSaveBankAccountRequest{
			BankName:          "Bank Syariah Indonesia",
			AccountNumber:     "1234567890",
			AccountHolderName: "Ahmad",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.VendorStatus != domain.VendorStatusSubmitted {
			t.Fatalf("expected submitted status, got %s", resp.VendorStatus)
		}
	})

	t.Run("corporate stays draft when only owner document exists", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		legalType := domain.VendorBusinessLegalTypeCorporate

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:                vendorID,
			Status:            domain.VendorStatusDraft,
			BusinessLegalType: &legalType,
		}, nil)
		vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
		vendorRepo.EXPECT().UpsertBankAccount(ctx, gomock.Any()).Return(nil)
		vendorRepo.EXPECT().
			FindDocumentsByVendorID(ctx, vendorID).
			Return(buildDocuments(
				vendorID,
				domain.VendorDocumentTypeOwnerDocumentID,
				domain.VendorDocumentTypeStorePhoto,
				domain.VendorDocumentTypeBankAccountProof,
				domain.VendorDocumentTypeBusinessLogo,
				domain.VendorDocumentTypeBusinessBanner,
				domain.VendorDocumentTypeBusinessNPWP,
			), nil)

		resp, err := uc.SaveBankAccount(ctx, vendorID, domain.VendorSaveBankAccountRequest{
			BankName:          "Bank Syariah Indonesia",
			AccountNumber:     "1234567890",
			AccountHolderName: "Ahmad",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response")
		}
		if resp.VendorStatus != domain.VendorStatusDraft {
			t.Fatalf("expected draft status, got %s", resp.VendorStatus)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:     vendorID,
			Status: domain.VendorStatusSubmitted,
		}, nil)

		_, err := uc.SaveBankAccount(ctx, vendorID, domain.VendorSaveBankAccountRequest{
			BankName:          "BSI",
			AccountNumber:     "123",
			AccountHolderName: "Ahmad",
		})
		if !errors.Is(err, ErrInvalidStatusTransition) {
			t.Fatalf("expected ErrInvalidStatusTransition, got %v", err)
		}
	})
}

func TestPresignDocument_TableDriven(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, vendorRepo, storage, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:     vendorID,
			Status: domain.VendorStatusDraft,
		}, nil)
		storage.EXPECT().
			GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
			Return("https://upload.example.com", nil)

		resp, err := uc.PresignDocument(ctx, vendorID, domain.VendorDocumentPresignRequest{
			DocType: domain.VendorDocumentTypeStorePhoto,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.UploadURL == "" {
			t.Fatal("expected upload url")
		}
	})

	t.Run("invalid doc type", func(t *testing.T) {
		_, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:     vendorID,
			Status: domain.VendorStatusDraft,
		}, nil)

		_, err := uc.PresignDocument(ctx, vendorID, domain.VendorDocumentPresignRequest{
			DocType: domain.VendorDocumentTypeOwnerDocumentID,
		})
		if !errors.Is(err, ErrInvalidVendorDocumentType) {
			t.Fatalf("expected ErrInvalidVendorDocumentType, got %v", err)
		}
	})
}

func TestLogin_TableDriven(t *testing.T) {
	type testCase struct {
		name       string
		req        domain.VendorLoginRequest
		setupMocks func(
			ctx context.Context,
			userRepo *mocks.MockUserRepository,
			vendorRepo *mocks.MockVendorRepository,
		)
		wantErr    error
		wantAnyErr bool
		assertResp func(t *testing.T, resp *domain.VendorLoginResponse)
	}

	tests := []testCase{
		{
			name: "success",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
			) {
				userID := uuid.New()
				vendorID := uuid.New()
				imageURL := "https://cdn.example.com/vendors/profile.jpg"
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           userID,
					Email:        "vendor@example.com",
					FullName:     "Abu Bakar Shidiq Basalamah",
					ImageURL:     &imageURL,
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM, Name: "UMKM"},
				}, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					VendorType:  "souvenir_store",
					DisplayName: "Toko Haji",
					Status:      domain.VendorStatusActive,
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.Token == "" {
					t.Error("expected non-empty token")
				}
				if resp.VendorID == uuid.Nil {
					t.Error("expected non-empty vendor ID")
				}
				if resp.ImageURL != "https://cdn.example.com/vendors/profile.jpg" {
					t.Errorf("expected image URL %q, got %q", "https://cdn.example.com/vendors/profile.jpg", resp.ImageURL)
				}
				if resp.Email != "vendor@example.com" {
					t.Errorf("expected email %q, got %q", "vendor@example.com", resp.Email)
				}
				if resp.Name != "Abu Bakar Shidiq Basalamah" {
					t.Errorf("expected name %q, got %q", "Abu Bakar Shidiq Basalamah", resp.Name)
				}
				if resp.VendorType != "souvenir_store" {
					t.Errorf("expected vendor type %q, got %q", "souvenir_store", resp.VendorType)
				}
				if resp.VendorStatus != domain.VendorStatusActive {
					t.Errorf("expected status %q, got %q", domain.VendorStatusActive, resp.VendorStatus)
				}
				if resp.DisplayName != "Toko Haji" {
					t.Errorf("expected display name %q, got %q", "Toko Haji", resp.DisplayName)
				}
			},
		},
		{
			name: "success without vendor image",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
			) {
				userID := uuid.New()
				vendorID := uuid.New()
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           userID,
					Email:        "vendor@example.com",
					FullName:     "Abu Bakar Shidiq Basalamah",
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM, Name: "UMKM"},
				}, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					VendorType:  "ppiu",
					DisplayName: "Toko Haji",
					Status:      domain.VendorStatusSubmitted,
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.ImageURL != "" {
					t.Errorf("expected empty image URL, got %q", resp.ImageURL)
				}
				if resp.Name != "Abu Bakar Shidiq Basalamah" {
					t.Errorf("expected name %q, got %q", "Abu Bakar Shidiq Basalamah", resp.Name)
				}
				if resp.VendorType != "ppiu" {
					t.Errorf("expected vendor type %q, got %q", "ppiu", resp.VendorType)
				}
				if resp.DisplayName != "Toko Haji" {
					t.Errorf("expected display name %q, got %q", "Toko Haji", resp.DisplayName)
				}
				if resp.VendorStatus != domain.VendorStatusSubmitted {
					t.Errorf("expected status %q, got %q", domain.VendorStatusSubmitted, resp.VendorStatus)
				}
			},
		},
		{
			name: "user not found",
			req:  domain.VendorLoginRequest{Email: "unknown@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
			) {
				userRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)
			},
			wantErr: ErrVendorInvalidCredentials,
		},
		{
			name: "wrong password",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "wrong_password"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
			) {
				hash, err := hashPasswordForTest("correct_password")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           uuid.New(),
					Email:        "vendor@example.com",
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM},
				}, nil)
			},
			wantErr: ErrVendorInvalidCredentials,
		},
		{
			name: "not vendor role",
			req:  domain.VendorLoginRequest{Email: "admin@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
			) {
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(&domain.User{
					ID:           uuid.New(),
					Email:        "admin@example.com",
					PasswordHash: hash,
					Role:         &domain.Role{Code: "admin"},
				}, nil)
			},
			wantErr: ErrNotVendor,
		},
		{
			name: "nil role",
			req:  domain.VendorLoginRequest{Email: "norole@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
			) {
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "norole@example.com").Return(&domain.User{
					ID:           uuid.New(),
					Email:        "norole@example.com",
					PasswordHash: hash,
					Role:         nil,
				}, nil)
			},
			wantErr: ErrNotVendor,
		},
		{
			name: "no vendor profile",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
			) {
				userID := uuid.New()
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           userID,
					Email:        "vendor@example.com",
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM},
				}, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, nil)
			},
			wantErr: ErrNoVendorProfile,
		},
		{
			name: "vendor blocked",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
			) {
				userID := uuid.New()
				vendorID := uuid.New()
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           userID,
					Email:        "vendor@example.com",
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM},
				}, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					Status:      domain.VendorStatusBlocked,
				}, nil)
			},
			wantErr: ErrVendorAccountBlocked,
		},
		{
			name: "find user error",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
			) {
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name: "find vendor error",
			req:  domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
			) {
				userID := uuid.New()
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           userID,
					Email:        "vendor@example.com",
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM},
				}, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
			ctx := context.Background()

			tc.setupMocks(ctx, userRepo, vendorRepo)

			resp, err := uc.Login(ctx, tc.req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAnyErr {
				if err == nil {
					t.Fatal("expected error, got nil")
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

func TestGetMe_TableDriven(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		vendorID   uuid.UUID
		setupMocks func(context.Context, *mocks.MockUserRepository, *mocks.MockVendorRepository, uuid.UUID)
		wantErr    error
		wantAnyErr bool
		assertResp func(*testing.T, *domain.VendorProfileResponse, uuid.UUID)
	}{
		{
			name:     "success",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				imageURL := "https://cdn.example.com/vendors/profile.jpg"
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					VendorType:  domain.VendorTypeSouvenirStore,
					Status:      domain.VendorStatusActive,
					DisplayName: "Toko Mabrur",
				}, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(&domain.User{
					ID:        ownerID,
					Email:     "toko.mabrur@example.com",
					FullName:  "Abu Bakar Shidiq Basalamah",
					ImageURL:  &imageURL,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorProfileResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.VendorID != vendorID {
					t.Errorf("expected vendor ID %s, got %s", vendorID, resp.VendorID)
				}
				if resp.ImageURL != "https://cdn.example.com/vendors/profile.jpg" {
					t.Errorf("expected image URL to match, got %q", resp.ImageURL)
				}
				if resp.Email != "toko.mabrur@example.com" {
					t.Errorf("expected email toko.mabrur@example.com, got %s", resp.Email)
				}
				if resp.Name != "Abu Bakar Shidiq Basalamah" {
					t.Errorf("expected owner full name, got %s", resp.Name)
				}
				if resp.VendorType != domain.VendorTypeSouvenirStore {
					t.Errorf("expected vendor type %s, got %s", domain.VendorTypeSouvenirStore, resp.VendorType)
				}
				if resp.VendorStatus != domain.VendorStatusActive {
					t.Errorf("expected vendor status %s, got %s", domain.VendorStatusActive, resp.VendorStatus)
				}
				if resp.DisplayName != "Toko Mabrur" {
					t.Errorf("expected display name Toko Mabrur, got %s", resp.DisplayName)
				}
			},
		},
		{
			name:     "success without image",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					VendorType:  domain.VendorTypePPIU,
					Status:      domain.VendorStatusSubmitted,
					DisplayName: "Toko Tanpa Foto",
				}, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(&domain.User{
					ID:       ownerID,
					Email:    "vendor@example.com",
					FullName: "Owner Vendor",
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorProfileResponse, _ uuid.UUID) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.ImageURL != "" {
					t.Errorf("expected empty image URL, got %q", resp.ImageURL)
				}
				if resp.Name != "Owner Vendor" {
					t.Errorf("expected owner name, got %s", resp.Name)
				}
			},
		},
		{
			name:     "vendor not found",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, _ *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name:     "owner user not found",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					DisplayName: "Toko Mabrur",
				}, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, nil)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:     "vendor repo error",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, _ *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name:     "user repo error",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
				}, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			userRepo, vendorRepo, _, _, _, uc := setupVendorUseCase(t)
			if tc.setupMocks != nil {
				tc.setupMocks(ctx, userRepo, vendorRepo, tc.vendorID)
			}

			resp, err := uc.GetMe(ctx, tc.vendorID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAnyErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.assertResp != nil {
				tc.assertResp(t, resp, tc.vendorID)
			}
		})
	}
}

func hashPasswordForTest(pw string) (string, error) {
	return auth.HashPassword(pw)
}

// --- Setup helper for withdrawal/balance tests (includes XenditPayoutProvider mock) ---

func setupVendorWithdrawalUseCase(t *testing.T) (
	*mocks.MockUserRepository,
	*mocks.MockVendorRepository,
	*mocks.MockStorageProvider,
	*mocks.MockXenditPayoutProvider,
	domain.VendorUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	xenditPayout := mocks.NewMockXenditPayoutProvider(ctrl)
	uc := NewVendorUseCase(&stubOTPUseCase{}, userRepo, vendorRepo, &stubVendorOnboardingRepo{}, storage, xenditPayout, "test-secret", 3600, "test-issuer")
	return userRepo, vendorRepo, storage, xenditPayout, uc
}

func TestGetBalance_TableDriven(t *testing.T) {
	type testCase struct {
		name       string
		setupMocks func(ctx context.Context, vendorID uuid.UUID, vendorRepo *mocks.MockVendorRepository)
		wantErr    error
		wantAnyErr bool
		assertResp func(t *testing.T, resp *domain.VendorBalanceResponse)
	}

	tests := []testCase{
		{
			name: "success",
			setupMocks: func(ctx context.Context, vendorID uuid.UUID, vendorRepo *mocks.MockVendorRepository) {
				vendorRepo.EXPECT().GetBalance(ctx, vendorID).Return(&domain.VendorBalance{
					VendorID:         vendorID,
					AvailableBalance: 500000,
					PendingBalance:   100000,
					EscrowBalance:    250000,
					TotalEarned:      1500000,
					TotalWithdrawn:   900000,
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorBalanceResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.AvailableBalance != 500000 {
					t.Errorf("expected available_balance 500000, got %f", resp.AvailableBalance)
				}
				if resp.PendingBalance != 100000 {
					t.Errorf("expected pending_balance 100000, got %f", resp.PendingBalance)
				}
				if resp.EscrowBalance != 250000 {
					t.Errorf("expected escrow_balance 250000, got %f", resp.EscrowBalance)
				}
			},
		},
		{
			name: "balance not found",
			setupMocks: func(ctx context.Context, vendorID uuid.UUID, vendorRepo *mocks.MockVendorRepository) {
				vendorRepo.EXPECT().GetBalance(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrBalanceNotFound,
		},
		{
			name: "repo error",
			setupMocks: func(ctx context.Context, vendorID uuid.UUID, vendorRepo *mocks.MockVendorRepository) {
				vendorRepo.EXPECT().GetBalance(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, vendorRepo, _, _, uc := setupVendorWithdrawalUseCase(t)
			ctx := context.Background()
			vendorID := uuid.New()

			tc.setupMocks(ctx, vendorID, vendorRepo)

			resp, err := uc.GetBalance(ctx, vendorID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAnyErr {
				if err == nil {
					t.Fatal("expected error, got nil")
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

func TestListPayoutChannels_TableDriven(t *testing.T) {
	vendorID := uuid.New()
	activeVendor := &domain.Vendor{
		ID:     vendorID,
		Status: domain.VendorStatusActive,
	}

	type testCase struct {
		name       string
		setupMocks func(ctx context.Context, vendorRepo *mocks.MockVendorRepository, xenditPayout *mocks.MockXenditPayoutProvider)
		wantErr    error
		assertResp func(t *testing.T, resp *domain.VendorPayoutChannelsResponse)
	}

	tests := []testCase{
		{
			name: "success with sorting and filtering inactive channel",
			setupMocks: func(ctx context.Context, vendorRepo *mocks.MockVendorRepository, xenditPayout *mocks.MockXenditPayoutProvider) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return([]domain.XenditPayoutChannel{
						{ChannelCode: "ID_BRI", ChannelName: "BRI", Currency: "IDR", ChannelCategory: "BANK", IsActivated: true},
						{ChannelCode: "ID_BCA", ChannelName: "BCA", Currency: "IDR", ChannelCategory: "BANK", IsActivated: true},
						{ChannelCode: "ID_BNI", ChannelName: "BNI", Currency: "IDR", ChannelCategory: "BANK", IsActivated: false},
					}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorPayoutChannelsResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if len(resp.Channels) != 2 {
					t.Fatalf("expected 2 active channels, got %d", len(resp.Channels))
				}
				if resp.Channels[0].ChannelCode != "ID_BCA" || resp.Channels[1].ChannelCode != "ID_BRI" {
					t.Fatalf("expected sorted channels [ID_BCA ID_BRI], got [%s %s]", resp.Channels[0].ChannelCode, resp.Channels[1].ChannelCode)
				}
			},
		},
		{
			name: "vendor not found",
			setupMocks: func(ctx context.Context, vendorRepo *mocks.MockVendorRepository, _ *mocks.MockXenditPayoutProvider) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "vendor not active",
			setupMocks: func(ctx context.Context, vendorRepo *mocks.MockVendorRepository, _ *mocks.MockXenditPayoutProvider) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:     vendorID,
					Status: domain.VendorStatusDraft,
				}, nil)
			},
			wantErr: ErrVendorNotActive,
		},
		{
			name: "xendit unavailable",
			setupMocks: func(ctx context.Context, vendorRepo *mocks.MockVendorRepository, xenditPayout *mocks.MockXenditPayoutProvider) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return(nil, errors.New("xendit down"))
			},
			wantErr: ErrPayoutChannelsUnavailable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, vendorRepo, _, xenditPayout, uc := setupVendorWithdrawalUseCase(t)
			ctx := context.Background()

			tc.setupMocks(ctx, vendorRepo, xenditPayout)
			resp, err := uc.ListPayoutChannels(ctx, vendorID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
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

func TestRequestWithdrawal_TableDriven(t *testing.T) {
	xenditAccountID := "xnd_test_123"
	vendorID := uuid.New()
	userID := uuid.New()

	activeVendor := &domain.Vendor{
		ID:              vendorID,
		OwnerUserID:     userID,
		DisplayName:     "Toko Haji",
		Status:          domain.VendorStatusActive,
		XenditAccountID: &xenditAccountID,
	}

	goodBalance := &domain.VendorBalance{
		VendorID:         vendorID,
		AvailableBalance: 500000,
		PendingBalance:   0,
		TotalEarned:      500000,
		TotalWithdrawn:   0,
	}

	goodBankAccount := &domain.VendorBankAccount{
		ID:                uuid.New(),
		VendorID:          vendorID,
		BankName:          "BCA",
		AccountNumber:     "1234567890",
		AccountHolderName: "Ahmad",
	}

	ownerUser := &domain.User{
		ID:    userID,
		Email: "vendor@example.com",
	}

	availableChannels := []domain.XenditPayoutChannel{
		{
			ChannelCode:     "ID_BCA",
			ChannelName:     "Bank Central Asia",
			Currency:        "IDR",
			ChannelCategory: "BANK",
			IsActivated:     true,
		},
		{
			ChannelCode:     "ID_BRI",
			ChannelName:     "Bank Rakyat Indonesia",
			Currency:        "IDR",
			ChannelCategory: "BANK",
			IsActivated:     true,
		},
	}

	type testCase struct {
		name       string
		req        domain.VendorWithdrawRequest
		setupMocks func(
			ctx context.Context,
			userRepo *mocks.MockUserRepository,
			vendorRepo *mocks.MockVendorRepository,
			xenditPayout *mocks.MockXenditPayoutProvider,
		)
		wantErr    error
		wantAnyErr bool
		assertResp func(t *testing.T, resp *domain.VendorWithdrawResponse)
	}

	tests := []testCase{
		{
			name: "success",
			req:  domain.VendorWithdrawRequest{Amount: 100000, ChannelCode: "ID_BCA"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return(availableChannels, nil)
				vendorRepo.EXPECT().GetBalance(ctx, vendorID).Return(goodBalance, nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(goodBankAccount, nil)
				userRepo.EXPECT().FindByID(ctx, userID).Return(ownerUser, nil)
				vendorRepo.EXPECT().DebitBalance(ctx, vendorID, float64(100000)).Return(nil)
				vendorRepo.EXPECT().CreateWithdrawal(ctx, gomock.Any()).Return(nil)
				xenditPayout.EXPECT().
					CreatePayout(ctx, xenditAccountID, gomock.Any(), gomock.Any()).
					Return(&domain.XenditPayoutResponse{
						ID:          "xnd_payout_123",
						Status:      "ACCEPTED",
						ChannelCode: "ID_BCA",
					}, nil)
				vendorRepo.EXPECT().UpdateWithdrawal(ctx, gomock.Any()).Return(nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorWithdrawResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.WithdrawalID == uuid.Nil {
					t.Error("expected non-nil withdrawal ID")
				}
				if resp.XenditPayoutID != "xnd_payout_123" {
					t.Errorf("expected xendit_payout_id xnd_payout_123, got %s", resp.XenditPayoutID)
				}
				if resp.Status != domain.WithdrawalStatusProcessing {
					t.Errorf("expected status processing, got %s", resp.Status)
				}
				if resp.Amount != 100000 {
					t.Errorf("expected amount 100000, got %f", resp.Amount)
				}
				if resp.RequestedAmount != 100000 {
					t.Errorf("expected requested_amount 100000, got %f", resp.RequestedAmount)
				}
				if resp.EstimatedFee != 0 {
					t.Errorf("expected estimated_fee 0, got %f", resp.EstimatedFee)
				}
				if resp.EstimatedNetAmount != 100000 {
					t.Errorf("expected estimated_net_amount 100000, got %f", resp.EstimatedNetAmount)
				}
			},
		},
		{
			name: "invalid channel code",
			req:  domain.VendorWithdrawRequest{Amount: 100000, ChannelCode: "INVALID"},
			setupMocks: func(
				ctx context.Context,
				_ *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return(availableChannels, nil)
			},
			wantErr: ErrInvalidChannelCode,
		},
		{
			name: "below minimum withdrawal",
			req:  domain.VendorWithdrawRequest{Amount: 5000, ChannelCode: "ID_BCA"},
			setupMocks: func(
				_ context.Context,
				_ *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				// no mocks needed — fails at minimum amount check
			},
			wantErr: ErrBelowMinWithdrawal,
		},
		{
			name: "vendor not active",
			req:  domain.VendorWithdrawRequest{Amount: 100000, ChannelCode: "ID_BCA"},
			setupMocks: func(
				ctx context.Context,
				_ *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:     vendorID,
					Status: domain.VendorStatusDraft,
				}, nil)
			},
			wantErr: ErrVendorNotActive,
		},
		{
			name: "no xendit account",
			req:  domain.VendorWithdrawRequest{Amount: 100000, ChannelCode: "ID_BCA"},
			setupMocks: func(
				ctx context.Context,
				_ *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:              vendorID,
					Status:          domain.VendorStatusActive,
					XenditAccountID: nil,
				}, nil)
			},
			wantErr: ErrVendorNoXenditAccount,
		},
		{
			name: "insufficient balance",
			req:  domain.VendorWithdrawRequest{Amount: 999999, ChannelCode: "ID_BCA"},
			setupMocks: func(
				ctx context.Context,
				_ *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return(availableChannels, nil)
				vendorRepo.EXPECT().GetBalance(ctx, vendorID).Return(&domain.VendorBalance{
					VendorID:         vendorID,
					AvailableBalance: 100000,
				}, nil)
			},
			wantErr: ErrInsufficientBalance,
		},
		{
			name: "xendit payout failed",
			req:  domain.VendorWithdrawRequest{Amount: 100000, ChannelCode: "ID_BCA"},
			setupMocks: func(
				ctx context.Context,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return(availableChannels, nil)
				vendorRepo.EXPECT().GetBalance(ctx, vendorID).Return(goodBalance, nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(goodBankAccount, nil)
				userRepo.EXPECT().FindByID(ctx, userID).Return(ownerUser, nil)
				vendorRepo.EXPECT().DebitBalance(ctx, vendorID, float64(100000)).Return(nil)
				vendorRepo.EXPECT().CreateWithdrawal(ctx, gomock.Any()).Return(nil)
				xenditPayout.EXPECT().
					CreatePayout(ctx, xenditAccountID, gomock.Any(), gomock.Any()).
					Return(nil, errors.New("xendit error"))
				vendorRepo.EXPECT().FailWithdrawal(ctx, vendorID, float64(100000)).Return(nil)
				vendorRepo.EXPECT().UpdateWithdrawal(ctx, gomock.Any()).Return(nil)
			},
			wantErr: ErrXenditPayoutFailed,
		},
		{
			name: "payout channels unavailable",
			req:  domain.VendorWithdrawRequest{Amount: 100000, ChannelCode: "ID_BCA"},
			setupMocks: func(
				ctx context.Context,
				_ *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				xenditPayout *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor, nil)
				xenditPayout.EXPECT().
					ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
						Currency:        "IDR",
						ChannelCategory: "BANK",
					}).
					Return(nil, errors.New("xendit unavailable"))
			},
			wantErr: ErrPayoutChannelsUnavailable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, vendorRepo, _, xenditPayout, uc := setupVendorWithdrawalUseCase(t)
			ctx := context.Background()

			tc.setupMocks(ctx, userRepo, vendorRepo, xenditPayout)

			resp, err := uc.RequestWithdrawal(ctx, vendorID, tc.req)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAnyErr {
				if err == nil {
					t.Fatal("expected error, got nil")
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

func TestRequestWithdrawal_NetAmountBelowMin(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	xenditPayout := mocks.NewMockXenditPayoutProvider(ctrl)

	uc := NewVendorUseCase(
		&stubOTPUseCase{},
		userRepo,
		vendorRepo,
		&stubVendorOnboardingRepo{},
		storage,
		xenditPayout,
		"test-secret",
		3600,
		"test-issuer",
		VendorWithdrawalPolicy{
			FeeEstimateFixed: 6000,
			MinNetAmount:     10000,
		},
	)

	_, err := uc.RequestWithdrawal(context.Background(), uuid.New(), domain.VendorWithdrawRequest{
		Amount:      10000,
		ChannelCode: domain.PayoutChannelIDBCA,
	})
	if !errors.Is(err, ErrWithdrawalNetAmountTooSmall) {
		t.Fatalf("expected ErrWithdrawalNetAmountTooSmall, got %v", err)
	}
}

func TestHandlePayoutWebhook_TableDriven(t *testing.T) {
	withdrawalID := uuid.New()
	vendorID := uuid.New()

	type testCase struct {
		name       string
		payload    domain.XenditPayoutWebhookPayload
		setupMocks func(
			ctx context.Context,
			vendorRepo *mocks.MockVendorRepository,
			xenditPayout *mocks.MockXenditPayoutProvider,
		)
		wantErr bool
	}

	tests := []testCase{
		{
			name: "succeeded moves processing to completed",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventSucceeded,
				Data: domain.XenditPayoutWebhookData{
					ID:          "xnd_payout_123",
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusSucceeded,
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(&domain.VendorWithdrawal{
					ID:       withdrawalID,
					VendorID: vendorID,
					Amount:   100000,
					Status:   domain.WithdrawalStatusProcessing,
				}, nil)
				vendorRepo.EXPECT().
					ApplyWithdrawalWebhookUpdate(ctx, withdrawalID, domain.WithdrawalStatusCompleted, domain.XenditPayoutStatusSucceeded, gomock.Any(), gomock.Nil(), domain.WithdrawalBalanceActionComplete, gomock.Any()).
					DoAndReturn(func(_ context.Context, _ uuid.UUID, _ string, _ string, xenditPayoutID *string, _ *string, _ string, completion *domain.WithdrawalCompletionData) error {
						if xenditPayoutID == nil || *xenditPayoutID != "xnd_payout_123" {
							t.Fatalf("unexpected xendit_payout_id: %+v", xenditPayoutID)
						}
						if completion == nil || completion.FeeActual != 3000 || completion.NetAmount != 97000 || completion.TotalDeducted != 100000 {
							t.Fatalf("unexpected completion data: %+v", completion)
						}
						return nil
					})
			},
		},
		{
			name: "failed moves processing to failed and stores reason",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventFailed,
				Data: domain.XenditPayoutWebhookData{
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusFailed,
					FailureCode: "PAYOUT_REJECTED",
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(&domain.VendorWithdrawal{
					ID:       withdrawalID,
					VendorID: vendorID,
					Amount:   100000,
					Status:   domain.WithdrawalStatusProcessing,
				}, nil)
				vendorRepo.EXPECT().
					ApplyWithdrawalWebhookUpdate(ctx, withdrawalID, domain.WithdrawalStatusFailed, domain.XenditPayoutStatusFailed, gomock.Nil(), gomock.Any(), domain.WithdrawalBalanceActionFail, gomock.Nil()).
					DoAndReturn(func(_ context.Context, _ uuid.UUID, _ string, _ string, _ *string, failedReason *string, _ string, _ *domain.WithdrawalCompletionData) error {
						if failedReason == nil || *failedReason == "" {
							t.Fatal("expected failed reason to be set")
						}
						return nil
					})
			},
		},
		{
			name: "reversed rolls back completed withdrawal",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventReversed,
				Data: domain.XenditPayoutWebhookData{
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusReversed,
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(&domain.VendorWithdrawal{
					ID:            withdrawalID,
					VendorID:      vendorID,
					Amount:        100000,
					TotalDeducted: 100000,
					Status:        domain.WithdrawalStatusCompleted,
				}, nil)
				vendorRepo.EXPECT().
					ApplyWithdrawalWebhookUpdate(ctx, withdrawalID, domain.WithdrawalStatusFailed, domain.XenditPayoutStatusReversed, gomock.Nil(), gomock.Nil(), domain.WithdrawalBalanceActionReverseCompleted, gomock.Nil()).
					Return(nil)
			},
		},
		{
			name: "unknown withdrawal is ignored",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventSucceeded,
				Data: domain.XenditPayoutWebhookData{
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusSucceeded,
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(nil, nil)
			},
		},
		{
			name: "duplicate succeeded on completed only updates xendit fields",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventSucceeded,
				Data: domain.XenditPayoutWebhookData{
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusSucceeded,
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(&domain.VendorWithdrawal{
					ID:       withdrawalID,
					VendorID: vendorID,
					Amount:   100000,
					Status:   domain.WithdrawalStatusCompleted,
				}, nil)
				vendorRepo.EXPECT().
					ApplyWithdrawalWebhookUpdate(ctx, withdrawalID, "", domain.XenditPayoutStatusSucceeded, gomock.Nil(), gomock.Nil(), domain.WithdrawalBalanceActionNone, gomock.Nil()).
					Return(nil)
			},
		},
		{
			name: "conflicting failed callback on completed is ignored",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventFailed,
				Data: domain.XenditPayoutWebhookData{
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusFailed,
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(&domain.VendorWithdrawal{
					ID:       withdrawalID,
					VendorID: vendorID,
					Amount:   100000,
					Status:   domain.WithdrawalStatusCompleted,
				}, nil)
			},
		},
		{
			name: "succeeded with amount below fixed fee returns error",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventSucceeded,
				Data: domain.XenditPayoutWebhookData{
					ID:          "xnd_payout_123",
					ReferenceID: withdrawalID.String(),
					Status:      domain.XenditPayoutStatusSucceeded,
				},
			},
			setupMocks: func(
				ctx context.Context,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
				vendorRepo.EXPECT().FindWithdrawalByID(ctx, withdrawalID).Return(&domain.VendorWithdrawal{
					ID:       withdrawalID,
					VendorID: vendorID,
					Amount:   2500,
					Status:   domain.WithdrawalStatusProcessing,
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "invalid reference id returns error",
			payload: domain.XenditPayoutWebhookPayload{
				Event: domain.XenditPayoutWebhookEventSucceeded,
				Data: domain.XenditPayoutWebhookData{
					ReferenceID: "not-a-uuid",
					Status:      domain.XenditPayoutStatusSucceeded,
				},
			},
			setupMocks: func(
				_ context.Context,
				_ *mocks.MockVendorRepository,
				_ *mocks.MockXenditPayoutProvider,
			) {
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, vendorRepo, _, xenditPayout, uc := setupVendorWithdrawalUseCase(t)
			ctx := context.Background()

			tc.setupMocks(ctx, vendorRepo, xenditPayout)
			err := uc.HandlePayoutWebhook(ctx, tc.payload)

			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
