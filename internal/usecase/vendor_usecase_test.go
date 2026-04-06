package usecase

import (
	"context"
	"errors"
	"strings"
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

type stubShippingUseCase struct {
	getProvinces    func(context.Context) ([]domain.ROProvince, error)
	getCities       func(context.Context, string) ([]domain.ROCity, error)
	getDistricts    func(context.Context, string) ([]domain.RODistrict, error)
	getSubdistricts func(context.Context, string) ([]domain.ROSubdistrict, error)
}

func (s *stubShippingUseCase) GetProvinces(ctx context.Context) ([]domain.ROProvince, error) {
	if s.getProvinces == nil {
		return nil, nil
	}
	return s.getProvinces(ctx)
}

func (s *stubShippingUseCase) GetCities(ctx context.Context, provinceID string) ([]domain.ROCity, error) {
	if s.getCities == nil {
		return nil, nil
	}
	return s.getCities(ctx, provinceID)
}

func (s *stubShippingUseCase) GetDistricts(ctx context.Context, cityID string) ([]domain.RODistrict, error) {
	if s.getDistricts == nil {
		return nil, nil
	}
	return s.getDistricts(ctx, cityID)
}

func (s *stubShippingUseCase) GetSubdistricts(ctx context.Context, districtID string) ([]domain.ROSubdistrict, error) {
	if s.getSubdistricts == nil {
		return nil, nil
	}
	return s.getSubdistricts(ctx, districtID)
}

func setupVendorUseCase(t *testing.T) (
	*mocks.MockUserRepository,
	*mocks.MockVendorRepository,
	*mocks.MockStorageProvider,
	*stubShippingUseCase,
	*stubOTPUseCase,
	*stubVendorOnboardingRepo,
	domain.VendorUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	shipping := &stubShippingUseCase{}
	otpUC := &stubOTPUseCase{}
	onboardingRepo := &stubVendorOnboardingRepo{}
	uc := NewVendorUseCase(otpUC, userRepo, vendorRepo, onboardingRepo, shipping, storage, nil, "test-secret", 3600, "test-issuer")
	return userRepo, vendorRepo, storage, shipping, otpUC, onboardingRepo, uc
}

func TestRequestRegistrationOTP(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo, _, _, _, otpUC, _, uc := setupVendorUseCase(t)
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
		userRepo, vendorRepo, _, _, _, _, uc := setupVendorUseCase(t)
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
	userRepo, _, _, _, otpUC, onboardingRepo, uc := setupVendorUseCase(t)
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
	if resp == nil || resp.EmailToken == "" {
		t.Fatal("expected email token")
	}
}

func TestSetRegistrationPassword_InvalidStep(t *testing.T) {
	_, _, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
	ctx := context.Background()
	onboardingID := uuid.New()
	onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
		if id != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", id)
		}
		return &domain.VendorOnboarding{
			ID:     onboardingID,
			Email:  "vendor@example.com",
			Status: "invalid_step",
		}, nil
	}

	_, err := uc.SetRegistrationPassword(ctx, onboardingID, domain.VendorRegistrationPasswordRequest{Password: "password123"})
	if !errors.Is(err, ErrVendorOnboardingStep) {
		t.Fatalf("expected ErrVendorOnboardingStep, got %v", err)
	}
}

func TestSetRegistrationPassword_Success(t *testing.T) {
	userRepo, _, _, _, _, onboardingRepo, uc := setupVendorUseCase(t)
	ctx := context.Background()
	onboardingID := uuid.New()

	onboardingRepo.findByID = func(_ context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
		if id != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", id)
		}
		return &domain.VendorOnboarding{
			ID:            onboardingID,
			Email:         "vendor@example.com",
			Status:        domain.VendorOnboardingStatusOTPVerified,
			OTPVerifiedAt: time.Now(),
		}, nil
	}
	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(nil, nil)
	onboardingRepo.finalizeRegistration = func(_ context.Context, input domain.VendorRegistrationFinalizeInput) error {
		if input.OnboardingID != onboardingID {
			t.Fatalf("unexpected onboarding id: %s", input.OnboardingID)
		}
		if input.UserRole != domain.RoleUMKM {
			t.Fatalf("unexpected user role: %s", input.UserRole)
		}
		if input.User == nil || input.User.PasswordHash == "" {
			t.Fatal("expected user with password hash")
		}
		if input.Vendor == nil || input.Vendor.Status != domain.VendorStatusDraft {
			t.Fatalf("unexpected vendor: %+v", input.Vendor)
		}
		if input.Vendor != nil && input.Vendor.VendorType != nil {
			t.Fatalf("expected nil vendor type for new vendor, got %+v", input.Vendor.VendorType)
		}
		return nil
	}

	resp, err := uc.SetRegistrationPassword(ctx, onboardingID, domain.VendorRegistrationPasswordRequest{
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.AccessToken == "" {
		t.Fatal("expected auth token in response")
	}
	if resp.ImageURL != nil {
		t.Fatalf("expected nil image URL, got %v", *resp.ImageURL)
	}
	if resp.StoreName != nil {
		t.Fatalf("expected nil store name, got %v", resp.StoreName)
	}
	if resp.VendorType != nil {
		t.Fatalf("expected nil vendor type, got %v", resp.VendorType)
	}
	if resp.VendorStatus != domain.VendorStatusDraft {
		t.Fatalf("unexpected vendor status: %s", resp.VendorStatus)
	}
}

func TestSubmitSouvenirStoreProposal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo, vendorRepo, _, shipping, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		userID := uuid.New()
		vendor := &domain.Vendor{
			ID:          vendorID,
			OwnerUserID: userID,
			Status:      domain.VendorStatusDraft,
		}
		owner := &domain.User{
			ID:       userID,
			Email:    "owner@example.com",
			FullName: "Old Name",
		}

		shipping.getProvinces = func(context.Context) ([]domain.ROProvince, error) {
			return []domain.ROProvince{{ID: "31", Name: "DKI Jakarta"}}, nil
		}
		shipping.getCities = func(_ context.Context, provinceID string) ([]domain.ROCity, error) {
			if provinceID != "31" {
				t.Fatalf("unexpected provinceID: %s", provinceID)
			}
			return []domain.ROCity{{ID: "3171", ProvinceID: "31", Name: "Jakarta Pusat"}}, nil
		}
		shipping.getDistricts = func(_ context.Context, cityID string) ([]domain.RODistrict, error) {
			if cityID != "3171" {
				t.Fatalf("unexpected cityID: %s", cityID)
			}
			return []domain.RODistrict{{ID: "317101", CityID: "3171", Name: "Menteng"}}, nil
		}
		shipping.getSubdistricts = func(_ context.Context, districtID string) ([]domain.ROSubdistrict, error) {
			if districtID != "317101" {
				t.Fatalf("unexpected districtID: %s", districtID)
			}
			return []domain.ROSubdistrict{{ID: "3171011001", DistrictID: "317101", Name: "Pegangsaan"}}, nil
		}

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
		userRepo.EXPECT().FindByID(ctx, userID).Return(owner, nil)
		userRepo.EXPECT().FindByPhone(ctx, "08123456789").Return(owner, nil)
		vendorRepo.EXPECT().
			SubmitSouvenirStoreProposal(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, input domain.SubmitSouvenirStoreProposalInput) error {
				if input.OwnerUser == nil || input.OwnerUser.FullName != "Ahmad" {
					t.Fatalf("unexpected owner user: %+v", input.OwnerUser)
				}
				if input.OwnerUser.Phone == nil || *input.OwnerUser.Phone != "08123456789" {
					t.Fatalf("unexpected owner phone: %+v", input.OwnerUser.Phone)
				}
				if input.Vendor == nil || input.Vendor.DisplayName == nil || *input.Vendor.DisplayName != "Toko Haji" {
					t.Fatalf("unexpected vendor payload: %+v", input.Vendor)
				}
				if input.Vendor.VendorType == nil || *input.Vendor.VendorType != domain.VendorTypeSouvenirStore {
					t.Fatalf("unexpected vendor type: %+v", input.Vendor.VendorType)
				}
				if input.Vendor.ProvinceName == nil || *input.Vendor.ProvinceName != "DKI Jakarta" {
					t.Fatalf("unexpected province name: %+v", input.Vendor.ProvinceName)
				}
				if input.Vendor.CityName == nil || *input.Vendor.CityName != "Jakarta Pusat" {
					t.Fatalf("unexpected city name: %+v", input.Vendor.CityName)
				}
				if input.Vendor.DistrictName == nil || *input.Vendor.DistrictName != "Menteng" {
					t.Fatalf("unexpected district name: %+v", input.Vendor.DistrictName)
				}
				if input.Vendor.SubdistrictName == nil || *input.Vendor.SubdistrictName != "Pegangsaan" {
					t.Fatalf("unexpected subdistrict name: %+v", input.Vendor.SubdistrictName)
				}
				if input.ResponsiblePerson == nil || input.ResponsiblePerson.NIK != "3173010101010001" {
					t.Fatalf("unexpected responsible person: %+v", input.ResponsiblePerson)
				}
				if len(input.Documents) != 3 {
					t.Fatalf("expected 3 documents, got %d", len(input.Documents))
				}
				docTypes := map[string]string{}
				for _, doc := range input.Documents {
					docTypes[doc.DocType] = doc.FileURL
				}
				if docTypes[domain.VendorDocumentTypeOwnerDocumentID] != "vendors/docs/ktp" {
					t.Fatalf("unexpected KTP doc: %+v", docTypes)
				}
				if docTypes[domain.VendorDocumentTypeBusinessNIB] != "vendors/docs/nib" {
					t.Fatalf("unexpected NIB doc: %+v", docTypes)
				}
				if docTypes[domain.VendorDocumentTypeHalalCertificate] != "vendors/docs/halal" {
					t.Fatalf("unexpected halal doc: %+v", docTypes)
				}
				return nil
			})

		resp, err := uc.SubmitSouvenirStoreProposal(ctx, vendorID, userID, domain.VendorSouvenirStoreProposalRequest{
			StoreName:        "Toko Haji",
			StoreDescription: "Pusat oleh-oleh",
			Address: domain.VendorSouvenirStoreProposalAddress{
				ProvinceID:    "31",
				CityID:        "3171",
				DistrictID:    "317101",
				SubdistrictID: "3171011001",
				PostalCode:    "10110",
				AddressLine:   "Jl. Wahid Hasyim No. 10",
			},
			ResponsiblePerson: domain.VendorSouvenirStoreProposalResponsiblePerson{
				Name:        "Ahmad",
				Phone:       "08123456789",
				Email:       "responsible@example.com",
				NIK:         "3173010101010001",
				KTPObjectID: "vendors/docs/ktp",
			},
			OtherDocuments: domain.VendorSouvenirStoreProposalOtherDocuments{
				NIBObjectID:              "vendors/docs/nib",
				HalalCertificateObjectID: "vendors/docs/halal",
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.Status != domain.VendorStatusSubmitted {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})

	t.Run("responsible person email can differ from owner email", func(t *testing.T) {
		userRepo, vendorRepo, _, shipping, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		userID := uuid.New()
		vendor := &domain.Vendor{
			ID:          vendorID,
			OwnerUserID: userID,
			Status:      domain.VendorStatusDraft,
		}
		owner := &domain.User{
			ID:       userID,
			Email:    "owner@example.com",
			FullName: "Owner Lama",
		}

		shipping.getProvinces = func(context.Context) ([]domain.ROProvince, error) {
			return []domain.ROProvince{{ID: "31", Name: "DKI Jakarta"}}, nil
		}
		shipping.getCities = func(context.Context, string) ([]domain.ROCity, error) {
			return []domain.ROCity{{ID: "3171", ProvinceID: "31", Name: "Jakarta Pusat"}}, nil
		}
		shipping.getDistricts = func(context.Context, string) ([]domain.RODistrict, error) {
			return []domain.RODistrict{{ID: "317101", CityID: "3171", Name: "Menteng"}}, nil
		}
		shipping.getSubdistricts = func(context.Context, string) ([]domain.ROSubdistrict, error) {
			return []domain.ROSubdistrict{{ID: "3171011001", DistrictID: "317101", Name: "Pegangsaan"}}, nil
		}

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
		userRepo.EXPECT().FindByID(ctx, userID).Return(owner, nil)
		userRepo.EXPECT().FindByPhone(ctx, "08123456789").Return(owner, nil)
		vendorRepo.EXPECT().
			SubmitSouvenirStoreProposal(ctx, gomock.Any()).
			DoAndReturn(func(_ context.Context, input domain.SubmitSouvenirStoreProposalInput) error {
				if input.OwnerUser == nil || input.OwnerUser.FullName != "Ahmad" {
					t.Fatalf("unexpected owner user: %+v", input.OwnerUser)
				}
				if input.ResponsiblePerson == nil || input.ResponsiblePerson.UserID != userID {
					t.Fatalf("unexpected responsible person: %+v", input.ResponsiblePerson)
				}
				return nil
			})

		resp, err := uc.SubmitSouvenirStoreProposal(ctx, vendorID, userID, domain.VendorSouvenirStoreProposalRequest{
			StoreName:        "Toko Haji",
			StoreDescription: "Pusat oleh-oleh",
			Address: domain.VendorSouvenirStoreProposalAddress{
				ProvinceID:    "31",
				CityID:        "3171",
				DistrictID:    "317101",
				SubdistrictID: "3171011001",
				PostalCode:    "10110",
				AddressLine:   "Jl. Wahid Hasyim No. 10",
			},
			ResponsiblePerson: domain.VendorSouvenirStoreProposalResponsiblePerson{
				Name:        "Ahmad",
				Phone:       "08123456789",
				Email:       "different@example.com",
				NIK:         "3173010101010001",
				KTPObjectID: "vendors/docs/ktp",
			},
			OtherDocuments: domain.VendorSouvenirStoreProposalOtherDocuments{
				NIBObjectID:              "vendors/docs/nib",
				HalalCertificateObjectID: "vendors/docs/halal",
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil || resp.Status != domain.VendorStatusSubmitted {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})

	t.Run("invalid location hierarchy", func(t *testing.T) {
		tests := []struct {
			name            string
			configureLookup func(*stubShippingUseCase)
		}{
			{
				name: "province not found",
				configureLookup: func(shipping *stubShippingUseCase) {
					shipping.getProvinces = func(context.Context) ([]domain.ROProvince, error) {
						return []domain.ROProvince{{ID: "32", Name: "Jawa Barat"}}, nil
					}
				},
			},
			{
				name: "city not found for province",
				configureLookup: func(shipping *stubShippingUseCase) {
					shipping.getProvinces = func(context.Context) ([]domain.ROProvince, error) {
						return []domain.ROProvince{{ID: "31", Name: "DKI Jakarta"}}, nil
					}
					shipping.getCities = func(context.Context, string) ([]domain.ROCity, error) {
						return []domain.ROCity{{ID: "9999", ProvinceID: "31", Name: "Kota Lain"}}, nil
					}
				},
			},
			{
				name: "district not found for city",
				configureLookup: func(shipping *stubShippingUseCase) {
					shipping.getProvinces = func(context.Context) ([]domain.ROProvince, error) {
						return []domain.ROProvince{{ID: "31", Name: "DKI Jakarta"}}, nil
					}
					shipping.getCities = func(context.Context, string) ([]domain.ROCity, error) {
						return []domain.ROCity{{ID: "3171", ProvinceID: "31", Name: "Jakarta Pusat"}}, nil
					}
					shipping.getDistricts = func(context.Context, string) ([]domain.RODistrict, error) {
						return []domain.RODistrict{{ID: "8888", CityID: "3171", Name: "District Lain"}}, nil
					}
				},
			},
			{
				name: "subdistrict not found for district",
				configureLookup: func(shipping *stubShippingUseCase) {
					shipping.getProvinces = func(context.Context) ([]domain.ROProvince, error) {
						return []domain.ROProvince{{ID: "31", Name: "DKI Jakarta"}}, nil
					}
					shipping.getCities = func(context.Context, string) ([]domain.ROCity, error) {
						return []domain.ROCity{{ID: "3171", ProvinceID: "31", Name: "Jakarta Pusat"}}, nil
					}
					shipping.getDistricts = func(context.Context, string) ([]domain.RODistrict, error) {
						return []domain.RODistrict{{ID: "317101", CityID: "3171", Name: "Menteng"}}, nil
					}
					shipping.getSubdistricts = func(context.Context, string) ([]domain.ROSubdistrict, error) {
						return []domain.ROSubdistrict{{ID: "7777", DistrictID: "317101", Name: "Subdistrict Lain"}}, nil
					}
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				userRepo, vendorRepo, _, shipping, _, _, uc := setupVendorUseCase(t)
				ctx := context.Background()
				vendorID := uuid.New()
				userID := uuid.New()

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					Status:      domain.VendorStatusDraft,
				}, nil)
				userRepo.EXPECT().FindByID(ctx, userID).Return(&domain.User{
					ID:       userID,
					Email:    "owner@example.com",
					FullName: "Owner",
				}, nil)
				userRepo.EXPECT().FindByPhone(ctx, "08123456789").Return(&domain.User{
					ID: userID,
				}, nil)

				tc.configureLookup(shipping)

				_, err := uc.SubmitSouvenirStoreProposal(ctx, vendorID, userID, domain.VendorSouvenirStoreProposalRequest{
					StoreName:        "Toko Haji",
					StoreDescription: "Pusat oleh-oleh",
					Address: domain.VendorSouvenirStoreProposalAddress{
						ProvinceID:    "31",
						CityID:        "3171",
						DistrictID:    "317101",
						SubdistrictID: "3171011001",
						PostalCode:    "10110",
						AddressLine:   "Jl. Wahid Hasyim No. 10",
					},
					ResponsiblePerson: domain.VendorSouvenirStoreProposalResponsiblePerson{
						Name:        "Ahmad",
						Phone:       "08123456789",
						Email:       "responsible@example.com",
						NIK:         "3173010101010001",
						KTPObjectID: "vendors/docs/ktp",
					},
					OtherDocuments: domain.VendorSouvenirStoreProposalOtherDocuments{
						NIBObjectID:              "vendors/docs/nib",
						HalalCertificateObjectID: "vendors/docs/halal",
					},
				})
				if !errors.Is(err, ErrInvalidLocationSelection) {
					t.Fatalf("expected ErrInvalidLocationSelection, got %v", err)
				}
			})
		}
	})

	t.Run("vendor not owned", func(t *testing.T) {
		_, vendorRepo, _, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		userID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:          vendorID,
			OwnerUserID: uuid.New(),
			Status:      domain.VendorStatusDraft,
		}, nil)

		_, err := uc.SubmitSouvenirStoreProposal(ctx, vendorID, userID, domain.VendorSouvenirStoreProposalRequest{})
		if !errors.Is(err, ErrVendorNotOwned) {
			t.Fatalf("expected ErrVendorNotOwned, got %v", err)
		}
	})
}

func TestPresignSouvenirStoreProposalDocuments(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		_, vendorRepo, storage, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		userID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:          vendorID,
			OwnerUserID: userID,
			Status:      domain.VendorStatusDraft,
		}, nil)
		storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/png", PresignedUploadExpiry).Return("https://presigned.example.com/ktp", nil)
		storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/nib", nil)
		storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "application/pdf", PresignedUploadExpiry).Return("https://presigned.example.com/halal", nil)

		resp, err := uc.PresignSouvenirStoreProposalDocuments(ctx, vendorID, userID, domain.VendorSouvenirStoreProposalPresignRequest{
			ResponsiblePersonKTP: "image/png",
			NIB:                  "image/jpeg",
			HalalCertificate:     "application/pdf",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp == nil {
			t.Fatal("expected response, got nil")
		}
		if !strings.Contains(resp.ResponsiblePersonKTP.ObjectID, domain.VendorDocumentTypeOwnerDocumentID) {
			t.Fatalf("unexpected ktp object id: %s", resp.ResponsiblePersonKTP.ObjectID)
		}
		if !strings.Contains(resp.NIB.ObjectID, domain.VendorDocumentTypeBusinessNIB) {
			t.Fatalf("unexpected nib object id: %s", resp.NIB.ObjectID)
		}
		if !strings.Contains(resp.HalalCertificate.ObjectID, domain.VendorDocumentTypeHalalCertificate) {
			t.Fatalf("unexpected halal object id: %s", resp.HalalCertificate.ObjectID)
		}
	})

	t.Run("invalid content type", func(t *testing.T) {
		_, vendorRepo, _, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()
		userID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:          vendorID,
			OwnerUserID: userID,
			Status:      domain.VendorStatusDraft,
		}, nil)

		_, err := uc.PresignSouvenirStoreProposalDocuments(ctx, vendorID, userID, domain.VendorSouvenirStoreProposalPresignRequest{
			ResponsiblePersonKTP: "text/plain",
			NIB:                  "image/jpeg",
			HalalCertificate:     "application/pdf",
		})
		if !errors.Is(err, ErrInvalidDocumentContent) {
			t.Fatalf("expected ErrInvalidDocumentContent, got %v", err)
		}
	})

	t.Run("vendor not owned", func(t *testing.T) {
		_, vendorRepo, _, _, _, _, uc := setupVendorUseCase(t)
		ctx := context.Background()
		vendorID := uuid.New()

		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
			ID:          vendorID,
			OwnerUserID: uuid.New(),
			Status:      domain.VendorStatusDraft,
		}, nil)

		_, err := uc.PresignSouvenirStoreProposalDocuments(ctx, vendorID, uuid.New(), domain.VendorSouvenirStoreProposalPresignRequest{
			ResponsiblePersonKTP: "image/png",
			NIB:                  "image/jpeg",
			HalalCertificate:     "application/pdf",
		})
		if !errors.Is(err, ErrVendorNotOwned) {
			t.Fatalf("expected ErrVendorNotOwned, got %v", err)
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
				vendorType := domain.VendorTypeSouvenirStore
				displayName := "Toko Haji"
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					VendorType:  &vendorType,
					DisplayName: &displayName,
					Status:      domain.VendorStatusActive,
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.AccessToken == "" {
					t.Error("expected non-empty access token")
				}
				if resp.VendorID == uuid.Nil {
					t.Error("expected non-empty vendor ID")
				}
				if resp.ImageURL == nil || *resp.ImageURL != "https://cdn.example.com/vendors/profile.jpg" {
					t.Errorf("expected image URL %q, got %v", "https://cdn.example.com/vendors/profile.jpg", resp.ImageURL)
				}
				if resp.Email != "vendor@example.com" {
					t.Errorf("expected email %q, got %q", "vendor@example.com", resp.Email)
				}
				if resp.StoreName == nil || *resp.StoreName != "Toko Haji" {
					t.Errorf("expected store name %q, got %v", "Toko Haji", resp.StoreName)
				}
				if resp.VendorType == nil || *resp.VendorType != "souvenir_store" {
					t.Errorf("expected vendor type %q, got %v", "souvenir_store", resp.VendorType)
				}
				if resp.VendorStatus != domain.VendorStatusActive {
					t.Errorf("expected status %q, got %q", domain.VendorStatusActive, resp.VendorStatus)
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
				vendorType := domain.VendorTypePPIU
				displayName := "Toko Haji"
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					VendorType:  &vendorType,
					DisplayName: &displayName,
					Status:      domain.VendorStatusSubmitted,
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.AccessToken == "" {
					t.Error("expected non-empty access token")
				}
				if resp.ImageURL != nil {
					t.Errorf("expected nil image URL, got %v", *resp.ImageURL)
				}
				if resp.StoreName == nil || *resp.StoreName != "Toko Haji" {
					t.Errorf("expected store name %q, got %v", "Toko Haji", resp.StoreName)
				}
				if resp.VendorType == nil || *resp.VendorType != "ppiu" {
					t.Errorf("expected vendor type %q, got %v", "ppiu", resp.VendorType)
				}
				if resp.VendorStatus != domain.VendorStatusSubmitted {
					t.Errorf("expected status %q, got %q", domain.VendorStatusSubmitted, resp.VendorStatus)
				}
			},
		},
		{
			name: "success with nil vendor type",
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
				displayName := "Toko Haji"
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
					VendorType:  nil,
					DisplayName: &displayName,
					Status:      domain.VendorStatusDraft,
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorLoginResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.VendorType != nil {
					t.Errorf("expected nil vendor type, got %v", resp.VendorType)
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
			userRepo, vendorRepo, _, _, _, _, uc := setupVendorUseCase(t)
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
		assertResp func(*testing.T, *domain.VendorMeResponse, uuid.UUID)
	}{
		{
			name:     "success",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				imageURL := "https://cdn.example.com/vendors/profile.jpg"
				vendorType := domain.VendorTypeSouvenirStore
				displayName := "Toko Mabrur"
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					VendorType:  &vendorType,
					Status:      domain.VendorStatusActive,
					DisplayName: &displayName,
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
			assertResp: func(t *testing.T, resp *domain.VendorMeResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.VendorID != vendorID {
					t.Errorf("expected vendor ID %s, got %s", vendorID, resp.VendorID)
				}
				if resp.ImageURL == nil || *resp.ImageURL != "https://cdn.example.com/vendors/profile.jpg" {
					t.Errorf("expected image URL to match, got %v", resp.ImageURL)
				}
				if resp.Email != "toko.mabrur@example.com" {
					t.Errorf("expected email toko.mabrur@example.com, got %s", resp.Email)
				}
				if resp.VendorType == nil || *resp.VendorType != domain.VendorTypeSouvenirStore {
					t.Errorf("expected vendor type %s, got %v", domain.VendorTypeSouvenirStore, resp.VendorType)
				}
				if resp.VendorStatus != domain.VendorStatusActive {
					t.Errorf("expected vendor status %s, got %s", domain.VendorStatusActive, resp.VendorStatus)
				}
				if resp.StoreName == nil || *resp.StoreName != "Toko Mabrur" {
					t.Errorf("expected store name Toko Mabrur, got %v", resp.StoreName)
				}
			},
		},
		{
			name:     "success without image",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendorType := domain.VendorTypePPIU
				displayName := "Toko Tanpa Foto"
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					VendorType:  &vendorType,
					Status:      domain.VendorStatusSubmitted,
					DisplayName: &displayName,
				}, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(&domain.User{
					ID:       ownerID,
					Email:    "vendor@example.com",
					FullName: "Owner Vendor",
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorMeResponse, _ uuid.UUID) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.ImageURL != nil {
					t.Errorf("expected nil image URL, got %v", resp.ImageURL)
				}
				if resp.VendorType == nil || *resp.VendorType != domain.VendorTypePPIU {
					t.Errorf("expected vendor type %s, got %v", domain.VendorTypePPIU, resp.VendorType)
				}
				if resp.StoreName == nil || *resp.StoreName != "Toko Tanpa Foto" {
					t.Errorf("expected store name Toko Tanpa Foto, got %v", resp.StoreName)
				}
			},
		},
		{
			name:     "success with nil vendor type",
			vendorID: uuid.New(),
			setupMocks: func(ctx context.Context, userRepo *mocks.MockUserRepository, vendorRepo *mocks.MockVendorRepository, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					VendorType:  nil,
					Status:      domain.VendorStatusDraft,
					DisplayName: nil,
				}, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(&domain.User{
					ID:       ownerID,
					Email:    "vendor@example.com",
					FullName: "Owner Vendor",
				}, nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorMeResponse, _ uuid.UUID) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected non-nil response")
				}
				if resp.VendorType != nil {
					t.Errorf("expected nil vendor type, got %v", resp.VendorType)
				}
				if resp.StoreName != nil {
					t.Errorf("expected nil store name, got %v", resp.StoreName)
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
				displayName := "Toko Mabrur"
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: ownerID,
					DisplayName: &displayName,
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
			userRepo, vendorRepo, _, _, _, _, uc := setupVendorUseCase(t)
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
	uc := NewVendorUseCase(&stubOTPUseCase{}, userRepo, vendorRepo, &stubVendorOnboardingRepo{}, &stubShippingUseCase{}, storage, xenditPayout, "test-secret", 3600, "test-issuer")
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
		DisplayName:     ptrString("Toko Haji"),
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
		&stubShippingUseCase{},
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
