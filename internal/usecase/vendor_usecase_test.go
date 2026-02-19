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

// ---------- helpers ----------

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

// ============================================================
// Register
// ============================================================

func TestRegister_Success(t *testing.T) {
	userRepo, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	// Email not taken.
	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)

	// Presigned URL generated for each doc type (7 total: 6 required + 1 optional).
	storage.EXPECT().
		GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
		Return("https://s3.example.com/upload", nil).
		Times(len(allDocTypes))

	// User created with umkm role.
	userRepo.EXPECT().Create(ctx, gomock.Any(), "umkm").Return(nil)

	// Vendor + bank account + documents created.
	vendorRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	resp, err := uc.Register(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
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
}

func TestRegister_EmailAlreadyRegistered(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)
	// User exists but has no vendor.
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(nil, nil)

	_, err := uc.Register(ctx, req)
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Errorf("expected ErrEmailAlreadyRegistered, got %v", err)
	}
}

func TestRegister_VendorAlreadyExists(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
	existingVendor := &domain.Vendor{ID: uuid.New(), OwnerUserID: existingUser.ID}

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(existingVendor, nil)

	_, err := uc.Register(ctx, req)
	if !errors.Is(err, ErrVendorAlreadyExists) {
		t.Errorf("expected ErrVendorAlreadyExists, got %v", err)
	}
}

func TestRegister_FindByEmailError(t *testing.T) {
	userRepo, _, _, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, errors.New("db error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegister_FindVendorByOwnerError(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	existingUser := &domain.User{ID: uuid.New(), Email: req.Email}
	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(existingUser, nil)
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, existingUser.ID).Return(nil, errors.New("db error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegister_PresignedURLError(t *testing.T) {
	userRepo, _, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	// First call to GeneratePresignedUploadURL fails.
	storage.EXPECT().
		GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
		Return("", errors.New("s3 error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegister_CreateUserError(t *testing.T) {
	userRepo, _, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	storage.EXPECT().
		GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
		Return("https://s3.example.com/upload", nil).
		Times(len(allDocTypes))
	userRepo.EXPECT().Create(ctx, gomock.Any(), "umkm").Return(errors.New("db error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRegister_CreateVendorError(t *testing.T) {
	userRepo, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()
	req := validRegisterRequest()

	userRepo.EXPECT().FindByEmail(ctx, req.Email).Return(nil, nil)
	storage.EXPECT().
		GeneratePresignedUploadURL(ctx, gomock.Any(), "", PresignedUploadExpiry).
		Return("https://s3.example.com/upload", nil).
		Times(len(allDocTypes))
	userRepo.EXPECT().Create(ctx, gomock.Any(), "umkm").Return(nil)
	vendorRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Register(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// ConfirmDocuments
// ============================================================

func TestConfirmDocuments_Success_AllRequiredUploaded(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)

	// Build existing docs for ALL required types (uploaded_by is nil initially).
	existingDocs := make([]domain.VendorDocument, 0, len(requiredDocTypes))
	confirmItems := make([]domain.ConfirmDocumentItem, 0, len(requiredDocTypes))
	for _, docType := range requiredDocTypes {
		objKey := "vendors/" + vendorID.String() + "/documents/" + docType + "/file"
		existingDocs = append(existingDocs, domain.VendorDocument{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            docType,
			FileURL:            objKey,
			VerificationStatus: "pending",
		})
		confirmItems = append(confirmItems, domain.ConfirmDocumentItem{
			DocType:   docType,
			ObjectKey: objKey,
		})
	}

	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(existingDocs, nil)

	// HeadObject succeeds for each document.
	storage.EXPECT().
		HeadObject(ctx, gomock.Any()).
		Return(&domain.ObjectInfo{ContentType: "image/jpeg", ContentLength: 1024}, nil).
		Times(len(requiredDocTypes))

	// All required docs are uploaded → status transitions to "submitted".
	vendorRepo.EXPECT().
		ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), "submitted").
		Return(nil)

	req := domain.ConfirmDocumentsRequest{Documents: confirmItems}
	resp, err := uc.ConfirmDocuments(ctx, userID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.VendorStatus != "submitted" {
		t.Errorf("expected status 'submitted', got %s", resp.VendorStatus)
	}
	if resp.DocumentsCount != len(requiredDocTypes) {
		t.Errorf("expected %d documents confirmed, got %d", len(requiredDocTypes), resp.DocumentsCount)
	}
}

func TestConfirmDocuments_Success_PartialUpload(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)

	// Build existing docs for ALL required types, but only confirm one.
	existingDocs := make([]domain.VendorDocument, 0, len(requiredDocTypes))
	for _, docType := range requiredDocTypes {
		existingDocs = append(existingDocs, domain.VendorDocument{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            docType,
			FileURL:            "vendors/" + vendorID.String() + "/documents/" + docType + "/file",
			VerificationStatus: "pending",
		})
	}
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(existingDocs, nil)

	objectKey := "vendors/" + vendorID.String() + "/documents/owner_ktp/file"
	storage.EXPECT().
		HeadObject(ctx, objectKey).
		Return(&domain.ObjectInfo{ContentType: "image/png", ContentLength: 2048}, nil)

	// Only 1 of 6 required docs uploaded → status stays "draft" (empty newStatus).
	vendorRepo.EXPECT().
		ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), "").
		Return(nil)

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{
			{DocType: "owner_ktp", ObjectKey: objectKey},
		},
	}
	resp, err := uc.ConfirmDocuments(ctx, userID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.VendorStatus != "draft" {
		t.Errorf("expected status 'draft', got %s", resp.VendorStatus)
	}
	if resp.DocumentsCount != 1 {
		t.Errorf("expected 1 document confirmed, got %d", resp.DocumentsCount)
	}
}

func TestConfirmDocuments_VendorNotFound(t *testing.T) {
	_, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()
	userID := uuid.New()

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, nil)

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "key"}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if !errors.Is(err, ErrVendorNotFound) {
		t.Errorf("expected ErrVendorNotFound, got %v", err)
	}
}

func TestConfirmDocuments_FindVendorError(t *testing.T) {
	_, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()
	userID := uuid.New()

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, errors.New("db error"))

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "key"}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmDocuments_FindDocumentsError(t *testing.T) {
	_, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: "key"}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmDocuments_DocumentNotFound(t *testing.T) {
	_, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	// Return empty list of existing docs.
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{}, nil)

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "nonexistent_type", ObjectKey: "key"}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if !errors.Is(err, ErrDocumentNotFound) {
		t.Errorf("expected ErrDocumentNotFound, got %v", err)
	}
}

func TestConfirmDocuments_ObjectNotUploaded(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}
	objectKey := "vendors/doc/file"

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
		{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: objectKey, VerificationStatus: "pending"},
	}, nil)
	// HeadObject returns nil (object does not exist).
	storage.EXPECT().HeadObject(ctx, objectKey).Return(nil, nil)

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: objectKey}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if !errors.Is(err, ErrObjectNotUploaded) {
		t.Errorf("expected ErrObjectNotUploaded, got %v", err)
	}
}

func TestConfirmDocuments_HeadObjectError(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}
	objectKey := "vendors/doc/file"

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
		{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: objectKey, VerificationStatus: "pending"},
	}, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(nil, errors.New("s3 error"))

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: objectKey}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmDocuments_InvalidContentType(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}
	objectKey := "vendors/doc/file"

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
		{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: objectKey, VerificationStatus: "pending"},
	}, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType:   "text/html",
		ContentLength: 1024,
	}, nil)

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: objectKey}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if !errors.Is(err, ErrInvalidDocumentContent) {
		t.Errorf("expected ErrInvalidDocumentContent, got %v", err)
	}
}

func TestConfirmDocuments_DocumentSizeOverflow(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}
	objectKey := "vendors/doc/file"

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
		{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: objectKey, VerificationStatus: "pending"},
	}, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType:   "image/jpeg",
		ContentLength: int64(math.MaxInt32) + 1,
	}, nil)

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: objectKey}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if !errors.Is(err, ErrDocumentSizeOverflow) {
		t.Errorf("expected ErrDocumentSizeOverflow, got %v", err)
	}
}

func TestConfirmDocuments_ConfirmRepoError(t *testing.T) {
	_, vendorRepo, storage, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	vendor := &domain.Vendor{ID: vendorID, OwnerUserID: userID, Status: "draft"}
	objectKey := "vendors/doc/file"

	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{
		{ID: uuid.New(), VendorID: vendorID, DocType: "owner_ktp", FileURL: objectKey, VerificationStatus: "pending"},
	}, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType: "image/jpeg", ContentLength: 1024,
	}, nil)
	vendorRepo.EXPECT().
		ConfirmDocumentsAndUpdateStatus(ctx, vendorID, gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	req := domain.ConfirmDocumentsRequest{
		Documents: []domain.ConfirmDocumentItem{{DocType: "owner_ktp", ObjectKey: objectKey}},
	}
	_, err := uc.ConfirmDocuments(ctx, userID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Login
// ============================================================

func TestLogin_Success(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	// Hash a known password for test.
	password := "password123"
	hash, err := hashPasswordForTest(password)
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
	}

	user := &domain.User{
		ID:           userID,
		Email:        "vendor@example.com",
		PasswordHash: hash,
		Role:         &domain.Role{Code: "umkm", Name: "UMKM"},
	}
	vendor := &domain.Vendor{
		ID:          vendorID,
		OwnerUserID: userID,
		DisplayName: "Toko Haji",
		Status:      "active",
	}

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(user, nil)
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)

	req := domain.VendorLoginRequest{Email: "vendor@example.com", Password: password}
	resp, loginErr := uc.Login(ctx, req)
	if loginErr != nil {
		t.Fatalf("expected no error, got %v", loginErr)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.VendorID != vendorID {
		t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
	}
	if resp.VendorStatus != "active" {
		t.Errorf("expected status 'active', got %s", resp.VendorStatus)
	}
	if resp.DisplayName != "Toko Haji" {
		t.Errorf("expected display name 'Toko Haji', got %s", resp.DisplayName)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo, _, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userRepo.EXPECT().FindByEmail(ctx, "unknown@example.com").Return(nil, nil)

	req := domain.VendorLoginRequest{Email: "unknown@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrVendorInvalidCredentials) {
		t.Errorf("expected ErrVendorInvalidCredentials, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo, _, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	hash, _ := hashPasswordForTest("correct_password")
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "vendor@example.com",
		PasswordHash: hash,
		Role:         &domain.Role{Code: "umkm"},
	}

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(user, nil)

	req := domain.VendorLoginRequest{Email: "vendor@example.com", Password: "wrong_password"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrVendorInvalidCredentials) {
		t.Errorf("expected ErrVendorInvalidCredentials, got %v", err)
	}
}

func TestLogin_NotVendorRole(t *testing.T) {
	userRepo, _, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	hash, _ := hashPasswordForTest("password123")
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "admin@example.com",
		PasswordHash: hash,
		Role:         &domain.Role{Code: "admin"},
	}

	userRepo.EXPECT().FindByEmail(ctx, "admin@example.com").Return(user, nil)

	req := domain.VendorLoginRequest{Email: "admin@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrNotVendor) {
		t.Errorf("expected ErrNotVendor, got %v", err)
	}
}

func TestLogin_NilRole(t *testing.T) {
	userRepo, _, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	hash, _ := hashPasswordForTest("password123")
	user := &domain.User{
		ID:           uuid.New(),
		Email:        "norole@example.com",
		PasswordHash: hash,
		Role:         nil,
	}

	userRepo.EXPECT().FindByEmail(ctx, "norole@example.com").Return(user, nil)

	req := domain.VendorLoginRequest{Email: "norole@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrNotVendor) {
		t.Errorf("expected ErrNotVendor, got %v", err)
	}
}

func TestLogin_NoVendorProfile(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	hash, _ := hashPasswordForTest("password123")
	user := &domain.User{
		ID:           userID,
		Email:        "vendor@example.com",
		PasswordHash: hash,
		Role:         &domain.Role{Code: "umkm"},
	}

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(user, nil)
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, nil)

	req := domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrNoVendorProfile) {
		t.Errorf("expected ErrNoVendorProfile, got %v", err)
	}
}

func TestLogin_VendorBlocked(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	vendorID := uuid.New()
	hash, _ := hashPasswordForTest("password123")
	user := &domain.User{
		ID:           userID,
		Email:        "vendor@example.com",
		PasswordHash: hash,
		Role:         &domain.Role{Code: "umkm"},
	}
	vendor := &domain.Vendor{
		ID:          vendorID,
		OwnerUserID: userID,
		Status:      "blocked",
	}

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(user, nil)
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(vendor, nil)

	req := domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if !errors.Is(err, ErrVendorAccountBlocked) {
		t.Errorf("expected ErrVendorAccountBlocked, got %v", err)
	}
}

func TestLogin_FindUserError(t *testing.T) {
	userRepo, _, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(nil, errors.New("db error"))

	req := domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_FindVendorError(t *testing.T) {
	userRepo, vendorRepo, _, uc := setupVendorUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	hash, _ := hashPasswordForTest("password123")
	user := &domain.User{
		ID:           userID,
		Email:        "vendor@example.com",
		PasswordHash: hash,
		Role:         &domain.Role{Code: "umkm"},
	}

	userRepo.EXPECT().FindByEmail(ctx, "vendor@example.com").Return(user, nil)
	vendorRepo.EXPECT().FindByOwnerUserID(ctx, userID).Return(nil, errors.New("db error"))

	req := domain.VendorLoginRequest{Email: "vendor@example.com", Password: "password123"}
	_, err := uc.Login(ctx, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---------- test helper ----------

// hashPasswordForTest creates a bcrypt hash for use in login test fixtures.
func hashPasswordForTest(pw string) (string, error) {
	return auth.HashPassword(pw)
}
