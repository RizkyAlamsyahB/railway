package usecase

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
	"go.uber.org/mock/gomock"
)

func setupVendorUseCase(t *testing.T) (
	*mocks.MockUserRepository,
	*mocks.MockVendorRepository,
	*mocks.MockStorageProvider,
	domain.VendorUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	uc := NewVendorUseCase(userRepo, vendorRepo, storage, nil, "test-secret", 3600, "test-issuer")
	return userRepo, vendorRepo, storage, uc
}

func validRegisterRequest() domain.VendorRegisterRequest {
	return domain.VendorRegisterRequest{
		StoreName:             "Toko Oleh-Oleh Haji",
		StoreType:             "umrah_souvenir_store",
		OwnerName:             "Ahmad",
		LegalName:             ptrString("PT Toko Haji"),
		ResponsiblePersonName: "Ahmad",
		Phone:                 "08123456789",
		Email:                 "vendor@example.com",
		Password:              "password123",
		BankName:              "BCA",
		BankAccountNumber:     "1234567890",
		BankAccountHolderName: "Ahmad",
	}
}

func TestRegister_TableDriven(t *testing.T) {
	type testCase struct {
		name       string
		setupMocks func(
			ctx context.Context,
			req domain.VendorRegisterRequest,
			userRepo *mocks.MockUserRepository,
			vendorRepo *mocks.MockVendorRepository,
			storage *mocks.MockStorageProvider,
		)
		wantErr    error
		wantAnyErr bool
		assertResp func(t *testing.T, resp *domain.VendorRegisterResponse)
	}

	tests := []testCase{
		{
			name: "success",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				storage.EXPECT().
					GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
					Return("https://s3.example.com/upload", nil).
					Times(len(allDocTypes))
				userRepo.EXPECT().Create(ctx, gomock.Any(), domain.RoleUMKM).Return(nil)
				vendorRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				vendorRepo.EXPECT().InitBalance(ctx, gomock.Any()).Return(nil)
			},
			assertResp: func(t *testing.T, resp *domain.VendorRegisterResponse) {
				t.Helper()
				if resp == nil {
					t.Fatal("expected response, got nil")
				}
				if resp.Token == "" {
					t.Error("expected non-empty token")
				}
				if len(resp.UploadURLs) != len(allDocTypes) {
					t.Errorf("expected %d upload URLs, got %d", len(allDocTypes), len(resp.UploadURLs))
				}
				if resp.VendorID == uuid.Nil {
					t.Error("expected non-nil vendor ID")
				}
				if resp.UserID == uuid.Nil {
					t.Error("expected non-nil user ID")
				}
			},
		},
		{
			name: "email already registered without vendor",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(nil, nil)
			},
			wantErr: ErrEmailAlreadyRegistered,
		},
		{
			name: "vendor already exists",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
				existingVendor := &domain.Vendor{ID: uuid.New(), OwnerUserID: existingUser.ID}
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(existingVendor, nil)
			},
			wantErr: ErrVendorAlreadyExists,
		},
		{
			name: "find by email error",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name: "find vendor by owner error",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				_ *mocks.MockStorageProvider,
			) {
				existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(nil, errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name: "presigned URL error",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				storage.EXPECT().
					GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
					Return("", errors.New("s3 error"))
			},
			wantAnyErr: true,
		},
		{
			name: "create user error",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				_ *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				storage.EXPECT().
					GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
					Return("https://s3.example.com/upload", nil).
					Times(len(allDocTypes))
				userRepo.EXPECT().Create(ctx, gomock.Any(), domain.RoleUMKM).Return(errors.New("db error"))
			},
			wantAnyErr: true,
		},
		{
			name: "create vendor error",
			setupMocks: func(
				ctx context.Context,
				req domain.VendorRegisterRequest,
				userRepo *mocks.MockUserRepository,
				vendorRepo *mocks.MockVendorRepository,
				storage *mocks.MockStorageProvider,
			) {
				userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
				storage.EXPECT().
					GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
					Return("https://s3.example.com/upload", nil).
					Times(len(allDocTypes))
				userRepo.EXPECT().Create(ctx, gomock.Any(), domain.RoleUMKM).Return(nil)
				vendorRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantAnyErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo, vendorRepo, storage, uc := setupVendorUseCase(t)
			ctx := context.Background()
			req := validRegisterRequest()

			tc.setupMocks(ctx, req, userRepo, vendorRepo, storage)

			resp, err := uc.Register(ctx, req)

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

	buildAllRequiredDocs := func(vendorID uuid.UUID) ([]domain.VendorDocument, []domain.ConfirmDocumentItem) {
		existingDocs := make([]domain.VendorDocument, 0, len(requiredDocTypes))
		confirmItems := make([]domain.ConfirmDocumentItem, 0, len(requiredDocTypes))
		for _, docType := range requiredDocTypes {
			objKey := "vendors/" + vendorID.String() + "/documents/" + docType + "/file"
			existingDocs = append(existingDocs, domain.VendorDocument{
				ID:                 uuid.New(),
				VendorID:           vendorID,
				DocType:            docType,
				FileURL:            objKey,
				VerificationStatus: domain.VerificationStatusPending,
			})
			confirmItems = append(confirmItems, domain.ConfirmDocumentItem{
				DocType:   docType,
				ObjectKey: objKey,
			})
		}
		return existingDocs, confirmItems
	}

	tests := []testCase{
		{
			name: "success all required uploaded",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				_, items := buildAllRequiredDocs(vendorID)
				return domain.ConfirmDocumentsRequest{Documents: items}
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
				existingDocs, _ := buildAllRequiredDocs(vendorID)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(existingDocs, nil)
				storage.EXPECT().
					HeadObject(ctx, gomock.Any()).
					Return(&domain.ObjectInfo{ContentType: "image/jpeg", ContentLength: 1024}, nil).
					Times(len(requiredDocTypes))
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), domain.VendorStatusSubmitted).
					Return(nil)
			},
			wantStatus: domain.VendorStatusSubmitted,
			wantDocCnt: len(requiredDocTypes),
		},
		{
			name: "success partial upload",
			reqBuilder: func(vendorID uuid.UUID) domain.ConfirmDocumentsRequest {
				objectKey := "vendors/" + vendorID.String() + "/documents/owner_ktp/file"
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{
						{DocType: "owner_ktp", ObjectKey: objectKey},
					},
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
				existingDocs, _ := buildAllRequiredDocs(vendorID)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(existingDocs, nil)
				storage.EXPECT().
					HeadObject(ctx, req.Documents[0].ObjectKey).
					Return(&domain.ObjectInfo{ContentType: "image/png", ContentLength: 2048}, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), "").
					Return(nil)
			},
			wantStatus: domain.VendorStatusDraft,
			wantDocCnt: 1,
		},
		{
			name: "vendor not found",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "key"}},
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
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "key"}},
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
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "key"}},
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
			name: "document not found",
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
			wantErr: ErrDocumentNotFound,
		},
		{
			name: "object not uploaded",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "vendors/doc/file"}},
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
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
					{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: req.Documents[0].ObjectKey, VerificationStatus: domain.VerificationStatusPending},
				}, nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(nil, nil)
			},
			wantErr: ErrObjectNotUploaded,
		},
		{
			name: "head object error",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "vendors/doc/file"}},
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
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
					{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: req.Documents[0].ObjectKey, VerificationStatus: domain.VerificationStatusPending},
				}, nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(nil, errors.New("s3 error"))
			},
			wantAnyErr: true,
		},
		{
			name: "invalid content type",
			reqBuilder: func(_ uuid.UUID) domain.ConfirmDocumentsRequest {
				return domain.ConfirmDocumentsRequest{
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "vendors/doc/file"}},
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
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
					{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: req.Documents[0].ObjectKey, VerificationStatus: domain.VerificationStatusPending},
				}, nil)
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
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "vendors/doc/file"}},
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
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
					{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: req.Documents[0].ObjectKey, VerificationStatus: domain.VerificationStatusPending},
				}, nil)
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
					Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "vendors/doc/file"}},
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
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
					{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: req.Documents[0].ObjectKey, VerificationStatus: domain.VerificationStatusPending},
				}, nil)
				storage.EXPECT().HeadObject(ctx, req.Documents[0].ObjectKey).Return(&domain.ObjectInfo{
					ContentType:   "image/jpeg",
					ContentLength: 1024,
				}, nil)
				vendorRepo.EXPECT().
					ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantAnyErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, vendorRepo, storage, uc := setupVendorUseCase(t)
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
				hash, err := hashPasswordForTest("password123")
				if err != nil {
					t.Fatalf("failed to hash password for test: %v", err)
				}
				userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(&domain.User{
					ID:           userID,
					Email:        "vendor@example.com",
					PasswordHash: hash,
					Role:         &domain.Role{Code: domain.RoleUMKM, Name: "UMKM"},
				}, nil)
				vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(&domain.Vendor{
					ID:          vendorID,
					OwnerUserID: userID,
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
				if resp.VendorStatus != domain.VendorStatusActive {
					t.Errorf("expected status %q, got %q", domain.VendorStatusActive, resp.VendorStatus)
				}
				if resp.DisplayName != "Toko Haji" {
					t.Errorf("expected display name %q, got %q", "Toko Haji", resp.DisplayName)
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
			userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
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
	uc := NewVendorUseCase(userRepo, vendorRepo, storage, xenditPayout, "test-secret", 3600, "test-issuer")
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
		userRepo,
		vendorRepo,
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
