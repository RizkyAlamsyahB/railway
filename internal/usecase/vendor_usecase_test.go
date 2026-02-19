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
	uc := NewVendorUseCase(userRepo, vendorRepo, storage, "test-secret", 3600, "test-issuer")
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
