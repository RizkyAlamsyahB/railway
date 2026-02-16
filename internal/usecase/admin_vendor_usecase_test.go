package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// ---------- helpers ----------

func setupAdminVendorUseCase(t *testing.T) (
	*mocks.MockVendorRepository,
	*mocks.MockUserRepository,
	*mocks.MockStorageProvider,
	domain.AdminVendorUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	uc := NewAdminVendorUseCase(vendorRepo, userRepo, storage)
	return vendorRepo, userRepo, storage, uc
}

func dummyVendor(id uuid.UUID, status string) *domain.Vendor {
	ownerID := uuid.New()
	now := time.Now()
	return &domain.Vendor{
		ID:                    id,
		OwnerUserID:           ownerID,
		VendorType:            "umrah_souvenir_store",
		DisplayName:           "Toko Oleh-Oleh Haji",
		LegalName:             ptrString("PT Toko Haji"),
		ResponsiblePersonName: "Ahmad",
		Description:           ptrString("Toko oleh-oleh haji terlengkap"),
		Status:                status,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

func dummyVendorWithOwner(id uuid.UUID, ownerID uuid.UUID, status string) *domain.Vendor {
	v := dummyVendor(id, status)
	v.OwnerUserID = ownerID
	return v
}

func dummyOwner(id uuid.UUID, status string) *domain.User {
	now := time.Now()
	return &domain.User{
		ID:        id,
		Email:     "owner@example.com",
		FullName:  "Owner Name",
		Phone:     ptrString("08123456789"),
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func dummyBankAccount(vendorID uuid.UUID) *domain.VendorBankAccount {
	now := time.Now()
	return &domain.VendorBankAccount{
		ID:                 uuid.New(),
		VendorID:           vendorID,
		BankName:           "BCA",
		AccountNumber:      "1234567890",
		AccountHolderName:  "Ahmad",
		VerificationStatus: "pending",
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func dummyDocuments(vendorID uuid.UUID) []domain.VendorDocument {
	now := time.Now()
	uploaderID := uuid.New()
	return []domain.VendorDocument{
		{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            "owner_ktp",
			FileURL:            "vendors/doc1.jpg",
			MimeType:           ptrString("image/jpeg"),
			UploadedBy:         &uploaderID,
			VerificationStatus: "pending",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            "store_photo",
			FileURL:            "",
			VerificationStatus: "pending",
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}
}

// ============================================================
// List
// ============================================================

func TestAdminVendorList_Success(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID1 := uuid.New()
	vendorID2 := uuid.New()
	ownerID1 := uuid.New()
	ownerID2 := uuid.New()

	vendors := []domain.Vendor{
		*dummyVendorWithOwner(vendorID1, ownerID1, "submitted"),
		*dummyVendorWithOwner(vendorID2, ownerID2, "active"),
	}

	params := domain.VendorListParams{Page: 1, Limit: 10}

	vendorRepo.EXPECT().List(ctx, params).Return(vendors, int64(2), nil)
	userRepo.EXPECT().FindByID(ctx, ownerID1).Return(dummyOwner(ownerID1, "active"), nil)
	userRepo.EXPECT().FindByID(ctx, ownerID2).Return(dummyOwner(ownerID2, "active"), nil)

	items, meta, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if meta.TotalItems != 2 {
		t.Errorf("expected total_items=2, got %d", meta.TotalItems)
	}
	if meta.TotalPages != 1 {
		t.Errorf("expected total_pages=1, got %d", meta.TotalPages)
	}
	if items[0].OwnerName != "Owner Name" {
		t.Errorf("expected owner name 'Owner Name', got %s", items[0].OwnerName)
	}
}

func TestAdminVendorList_DefaultsInvalidPage(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	params := domain.VendorListParams{Page: 0, Limit: 0}
	expectedParams := domain.VendorListParams{Page: 1, Limit: 10}

	vendorRepo.EXPECT().List(ctx, expectedParams).Return([]domain.Vendor{}, int64(0), nil)

	items, meta, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
	if meta.TotalItems != 0 {
		t.Errorf("expected total_items=0, got %d", meta.TotalItems)
	}
}

func TestAdminVendorList_LimitExceeds100(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	params := domain.VendorListParams{Page: 1, Limit: 200}
	expectedParams := domain.VendorListParams{Page: 1, Limit: 10}

	vendorRepo.EXPECT().List(ctx, expectedParams).Return([]domain.Vendor{}, int64(0), nil)

	_, _, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAdminVendorList_RepoError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	params := domain.VendorListParams{Page: 1, Limit: 10}

	vendorRepo.EXPECT().List(ctx, params).Return(nil, int64(0), errors.New("db error"))

	_, _, err := uc.List(ctx, params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorList_FindOwnerError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendors := []domain.Vendor{*dummyVendorWithOwner(vendorID, ownerID, "submitted")}
	params := domain.VendorListParams{Page: 1, Limit: 10}

	vendorRepo.EXPECT().List(ctx, params).Return(vendors, int64(1), nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))

	_, _, err := uc.List(ctx, params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorList_OwnerNotFound(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendors := []domain.Vendor{*dummyVendorWithOwner(vendorID, ownerID, "submitted")}
	params := domain.VendorListParams{Page: 1, Limit: 10}

	vendorRepo.EXPECT().List(ctx, params).Return(vendors, int64(1), nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, nil)

	items, _, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if items[0].OwnerName != "" {
		t.Errorf("expected empty owner name, got %s", items[0].OwnerName)
	}
	if items[0].OwnerEmail != "" {
		t.Errorf("expected empty owner email, got %s", items[0].OwnerEmail)
	}
}

func TestAdminVendorList_PaginationCalculation(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	ownerID := uuid.New()
	vendors := []domain.Vendor{*dummyVendorWithOwner(uuid.New(), ownerID, "active")}
	params := domain.VendorListParams{Page: 2, Limit: 3}

	vendorRepo.EXPECT().List(ctx, params).Return(vendors, int64(7), nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)

	_, meta, err := uc.List(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// 7 items / 3 per page = ceil(2.33) = 3 pages
	if meta.TotalPages != 3 {
		t.Errorf("expected total_pages=3, got %d", meta.TotalPages)
	}
}

// ============================================================
// GetByID
// ============================================================

func TestAdminVendorGetByID_Success(t *testing.T) {
	vendorRepo, userRepo, storage, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(dummyBankAccount(vendorID), nil)

	docs := dummyDocuments(vendorID)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(docs, nil)

	// Only the first document has UploadedBy != nil && FileURL != "", so only 1 presigned URL call.
	storage.EXPECT().GeneratePresignedURL(ctx, "vendors/doc1.jpg", presignedDownloadExpiry).Return("https://signed-url.example.com/doc1.jpg", nil)

	resp, err := uc.GetByID(ctx, vendorID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != vendorID {
		t.Errorf("expected vendor id %s, got %s", vendorID, resp.ID)
	}
	if resp.Owner.FullName != "Owner Name" {
		t.Errorf("expected owner name 'Owner Name', got %s", resp.Owner.FullName)
	}
	if resp.BankAccount == nil {
		t.Fatal("expected non-nil bank account")
	}
	if len(resp.Documents) != 2 {
		t.Errorf("expected 2 documents, got %d", len(resp.Documents))
	}
	if resp.Documents[0].DownloadURL != "https://signed-url.example.com/doc1.jpg" {
		t.Errorf("expected presigned download URL, got %s", resp.Documents[0].DownloadURL)
	}
	if resp.Documents[1].DownloadURL != "" {
		t.Errorf("expected empty download URL for non-uploaded doc, got %s", resp.Documents[1].DownloadURL)
	}
}

func TestAdminVendorGetByID_NotFound(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)

	_, err := uc.GetByID(ctx, vendorID)
	if !errors.Is(err, ErrVendorNotFound) {
		t.Errorf("expected ErrVendorNotFound, got %v", err)
	}
}

func TestAdminVendorGetByID_RepoError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))

	_, err := uc.GetByID(ctx, vendorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorGetByID_OwnerError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))

	_, err := uc.GetByID(ctx, vendorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorGetByID_BankAccountError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))

	_, err := uc.GetByID(ctx, vendorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorGetByID_DocumentsError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(dummyBankAccount(vendorID), nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))

	_, err := uc.GetByID(ctx, vendorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorGetByID_PresignedURLError(t *testing.T) {
	vendorRepo, userRepo, storage, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	uploaderID := uuid.New()
	docs := []domain.VendorDocument{
		{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            "owner_ktp",
			FileURL:            "vendors/doc1.jpg",
			UploadedBy:         &uploaderID,
			VerificationStatus: "pending",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		},
	}

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(dummyBankAccount(vendorID), nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(docs, nil)
	storage.EXPECT().GeneratePresignedURL(ctx, "vendors/doc1.jpg", presignedDownloadExpiry).Return("", errors.New("s3 error"))

	_, err := uc.GetByID(ctx, vendorID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorGetByID_NoBankAccount(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "draft")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
	vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{}, nil)

	resp, err := uc.GetByID(ctx, vendorID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.BankAccount != nil {
		t.Errorf("expected nil bank account, got %+v", resp.BankAccount)
	}
}

// ============================================================
// Approve
// ============================================================

func TestAdminVendorApprove_SuccessFromSubmitted(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "pending"), nil)
	userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	resp, err := uc.Approve(ctx, vendorID, adminID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "active" {
		t.Errorf("expected status 'active', got %s", resp.Status)
	}
	if resp.VendorID != vendorID {
		t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
	}
}

func TestAdminVendorApprove_SuccessFromRejected(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "rejected")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	// Owner is already active, so Update should NOT be called.

	resp, err := uc.Approve(ctx, vendorID, adminID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "active" {
		t.Errorf("expected status 'active', got %s", resp.Status)
	}
}

func TestAdminVendorApprove_InvalidStatus(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()

	for _, status := range []string{"draft", "active", "blocked"} {
		vendor := dummyVendor(vendorID, status)
		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)

		_, err := uc.Approve(ctx, vendorID, adminID)
		if !errors.Is(err, ErrInvalidStatusTransition) {
			t.Errorf("status=%s: expected ErrInvalidStatusTransition, got %v", status, err)
		}
	}
}

func TestAdminVendorApprove_NotFound(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)

	_, err := uc.Approve(ctx, vendorID, adminID)
	if !errors.Is(err, ErrVendorNotFound) {
		t.Errorf("expected ErrVendorNotFound, got %v", err)
	}
}

func TestAdminVendorApprove_FindByIDError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))

	_, err := uc.Approve(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorApprove_UpdateStatusError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	vendor := dummyVendor(vendorID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Approve(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorApprove_FindOwnerError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))

	_, err := uc.Approve(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorApprove_ActivateOwnerError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "pending"), nil)
	userRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Approve(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Reject
// ============================================================

func TestAdminVendorReject_Success(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendor := dummyVendor(vendorID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)

	resp, err := uc.Reject(ctx, vendorID, "dokumen tidak lengkap")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "rejected" {
		t.Errorf("expected status 'rejected', got %s", resp.Status)
	}
	if resp.VendorID != vendorID {
		t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
	}
}

func TestAdminVendorReject_InvalidStatus(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()

	for _, status := range []string{"draft", "active", "rejected", "blocked"} {
		vendor := dummyVendor(vendorID, status)
		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)

		_, err := uc.Reject(ctx, vendorID, "reason")
		if !errors.Is(err, ErrInvalidStatusTransition) {
			t.Errorf("status=%s: expected ErrInvalidStatusTransition, got %v", status, err)
		}
	}
}

func TestAdminVendorReject_NotFound(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)

	_, err := uc.Reject(ctx, vendorID, "reason")
	if !errors.Is(err, ErrVendorNotFound) {
		t.Errorf("expected ErrVendorNotFound, got %v", err)
	}
}

func TestAdminVendorReject_FindByIDError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))

	_, err := uc.Reject(ctx, vendorID, "reason")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorReject_UpdateStatusError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendor := dummyVendor(vendorID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Reject(ctx, vendorID, "reason")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Block
// ============================================================

func TestAdminVendorBlock_SuccessFromActive(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendor := dummyVendor(vendorID, "active")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)

	resp, err := uc.Block(ctx, vendorID, "pelanggaran kebijakan")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "blocked" {
		t.Errorf("expected status 'blocked', got %s", resp.Status)
	}
}

func TestAdminVendorBlock_SuccessFromSubmitted(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendor := dummyVendor(vendorID, "submitted")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)

	resp, err := uc.Block(ctx, vendorID, "reason")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "blocked" {
		t.Errorf("expected status 'blocked', got %s", resp.Status)
	}
}

func TestAdminVendorBlock_InvalidStatus(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()

	for _, status := range []string{"draft", "rejected", "blocked"} {
		vendor := dummyVendor(vendorID, status)
		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)

		_, err := uc.Block(ctx, vendorID, "reason")
		if !errors.Is(err, ErrInvalidStatusTransition) {
			t.Errorf("status=%s: expected ErrInvalidStatusTransition, got %v", status, err)
		}
	}
}

func TestAdminVendorBlock_NotFound(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)

	_, err := uc.Block(ctx, vendorID, "reason")
	if !errors.Is(err, ErrVendorNotFound) {
		t.Errorf("expected ErrVendorNotFound, got %v", err)
	}
}

func TestAdminVendorBlock_UpdateStatusError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	vendor := dummyVendor(vendorID, "active")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Block(ctx, vendorID, "reason")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Unblock
// ============================================================

func TestAdminVendorUnblock_Success(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "blocked"), nil)
	userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	resp, err := uc.Unblock(ctx, vendorID, adminID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "active" {
		t.Errorf("expected status 'active', got %s", resp.Status)
	}
	if resp.VendorID != vendorID {
		t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
	}
}

func TestAdminVendorUnblock_OwnerAlreadyActive(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
	// Owner is already active, so Update should NOT be called.

	resp, err := uc.Unblock(ctx, vendorID, adminID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != "active" {
		t.Errorf("expected status 'active', got %s", resp.Status)
	}
}

func TestAdminVendorUnblock_InvalidStatus(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()

	for _, status := range []string{"draft", "submitted", "active", "rejected"} {
		vendor := dummyVendor(vendorID, status)
		vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)

		_, err := uc.Unblock(ctx, vendorID, adminID)
		if !errors.Is(err, ErrInvalidStatusTransition) {
			t.Errorf("status=%s: expected ErrInvalidStatusTransition, got %v", status, err)
		}
	}
}

func TestAdminVendorUnblock_NotFound(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)

	_, err := uc.Unblock(ctx, vendorID, adminID)
	if !errors.Is(err, ErrVendorNotFound) {
		t.Errorf("expected ErrVendorNotFound, got %v", err)
	}
}

func TestAdminVendorUnblock_UpdateStatusError(t *testing.T) {
	vendorRepo, _, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	vendor := dummyVendor(vendorID, "blocked")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Unblock(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorUnblock_FindOwnerError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))

	_, err := uc.Unblock(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAdminVendorUnblock_ActivateOwnerError(t *testing.T) {
	vendorRepo, userRepo, _, uc := setupAdminVendorUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	adminID := uuid.New()
	ownerID := uuid.New()
	vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

	vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
	vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
	userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "blocked"), nil)
	userRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))

	_, err := uc.Unblock(ctx, vendorID, adminID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
