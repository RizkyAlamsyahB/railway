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
	*mocks.MockXenPlatformProvider,
	domain.AdminVendorUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	xenPlatform := mocks.NewMockXenPlatformProvider(ctrl)
	uc := NewAdminVendorUseCase(vendorRepo, userRepo, storage, xenPlatform)
	return vendorRepo, userRepo, storage, xenPlatform, uc
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

func TestAdminVendorList(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.VendorListParams
		setupMock func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context)
		wantErr   bool
		checkResp func(t *testing.T, items []domain.AdminVendorListItem, meta *domain.PaginationMeta)
	}{
		{
			name:   "success",
			params: domain.VendorListParams{Page: 1, Limit: 10},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorID1 := uuid.New()
				vendorID2 := uuid.New()
				ownerID1 := uuid.New()
				ownerID2 := uuid.New()

				v2 := dummyVendorWithOwner(vendorID2, ownerID2, "active")
				v2.XenditAccountID = ptrString("xen_account_123")

				vendors := []domain.Vendor{
					*dummyVendorWithOwner(vendorID1, ownerID1, "submitted"),
					*v2,
				}

				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return(vendors, int64(2), nil)
				userRepo.EXPECT().FindByID(ctx, ownerID1).Return(dummyOwner(ownerID1, "active"), nil)
				userRepo.EXPECT().FindByID(ctx, ownerID2).Return(dummyOwner(ownerID2, "active"), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminVendorListItem, meta *domain.PaginationMeta) {
				t.Helper()
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
				if items[0].XenditAccountID != nil {
					t.Errorf("expected nil xendit_account_id for submitted vendor, got %v", *items[0].XenditAccountID)
				}
				if items[1].XenditAccountID == nil || *items[1].XenditAccountID != "xen_account_123" {
					t.Errorf("expected xendit_account_id 'xen_account_123' for active vendor, got %v", items[1].XenditAccountID)
				}
			},
		},
		{
			name:   "defaults invalid page and limit",
			params: domain.VendorListParams{Page: 0, Limit: 0},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return([]domain.Vendor{}, int64(0), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminVendorListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 0 {
					t.Errorf("expected 0 items, got %d", len(items))
				}
				if meta.TotalItems != 0 {
					t.Errorf("expected total_items=0, got %d", meta.TotalItems)
				}
			},
		},
		{
			name:   "limit exceeds 100",
			params: domain.VendorListParams{Page: 1, Limit: 200},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return([]domain.Vendor{}, int64(0), nil)
			},
		},
		{
			name:   "repo error",
			params: domain.VendorListParams{Page: 1, Limit: 10},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "find owner error",
			params: domain.VendorListParams{Page: 1, Limit: 10},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorID := uuid.New()
				ownerID := uuid.New()
				vendors := []domain.Vendor{*dummyVendorWithOwner(vendorID, ownerID, "submitted")}

				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return(vendors, int64(1), nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "owner not found",
			params: domain.VendorListParams{Page: 1, Limit: 10},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorID := uuid.New()
				ownerID := uuid.New()
				vendors := []domain.Vendor{*dummyVendorWithOwner(vendorID, ownerID, "submitted")}

				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return(vendors, int64(1), nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminVendorListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if items[0].OwnerName != "" {
					t.Errorf("expected empty owner name, got %s", items[0].OwnerName)
				}
				if items[0].OwnerEmail != "" {
					t.Errorf("expected empty owner email, got %s", items[0].OwnerEmail)
				}
			},
		},
		{
			name:   "pagination calculation",
			params: domain.VendorListParams{Page: 2, Limit: 3},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				ownerID := uuid.New()
				vendors := []domain.Vendor{*dummyVendorWithOwner(uuid.New(), ownerID, "active")}

				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 2, Limit: 3}).Return(vendors, int64(7), nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminVendorListItem, meta *domain.PaginationMeta) {
				t.Helper()
				// 7 items / 3 per page = ceil(2.33) = 3 pages
				if meta.TotalPages != 3 {
					t.Errorf("expected total_pages=3, got %d", meta.TotalPages)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vendorRepo, userRepo, _, _, uc := setupAdminVendorUseCase(t)
			ctx := context.Background()

			tc.setupMock(vendorRepo, userRepo, ctx)

			items, meta, err := uc.List(ctx, tc.params)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, items, meta)
			}
		})
	}
}

// ============================================================
// GetByID
// ============================================================

func TestAdminVendorGetByID(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID)
		wantErr   error
		wantAny   bool // expect any non-nil error (not a specific sentinel)
		checkResp func(t *testing.T, resp *domain.AdminVendorDetailResponse, vendorID uuid.UUID)
	}{
		{
			name: "success",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(dummyBankAccount(vendorID), nil)

				docs := dummyDocuments(vendorID)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(docs, nil)

				// Only the first document has UploadedBy != nil && FileURL != "", so only 1 presigned URL call.
				storage.EXPECT().GeneratePresignedURL(ctx, "vendors/doc1.jpg", PresignedDownloadExpiry).Return("https://signed-url.example.com/doc1.jpg", nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorDetailResponse, vendorID uuid.UUID) {
				t.Helper()
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
			},
		},
		{
			name: "not found",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "repo error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "owner error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "bank account error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "documents error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(dummyBankAccount(vendorID), nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "presigned URL error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
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
				storage.EXPECT().GeneratePresignedURL(ctx, "vendors/doc1.jpg", PresignedDownloadExpiry).Return("", errors.New("s3 error"))
			},
			wantAny: true,
		},
		{
			name: "no bank account",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "draft")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindBankAccountByVendorID(ctx, vendorID).Return(nil, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{}, nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorDetailResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.BankAccount != nil {
					t.Errorf("expected nil bank account, got %+v", resp.BankAccount)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vendorRepo, userRepo, storage, _, uc := setupAdminVendorUseCase(t)
			ctx := context.Background()
			vendorID := uuid.New()

			tc.setupMock(vendorRepo, userRepo, storage, ctx, vendorID)

			resp, err := uc.GetByID(ctx, vendorID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp, vendorID)
			}
		})
	}
}

// ============================================================
// Approve
// ============================================================

func TestAdminVendorApprove(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID)
		wantErr   error
		wantAny   bool
		checkResp func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID)
	}{
		{
			name: "success from submitted",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "pending"), nil)
				xenPlatform.EXPECT().CreateAccount(ctx, domain.XenPlatformCreateAccountRequest{
					Email: "owner@example.com",
					Type:  "OWNED",
					PublicProfile: &domain.XenPlatformPublicProfile{
						BusinessName: "Toko Oleh-Oleh Haji",
					},
				}).Return(&domain.XenPlatformAccount{
					ID:     "xnd_acc_123",
					Type:   "OWNED",
					Email:  "owner@example.com",
					Status: "LIVE",
				}, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.Status != "active" {
					t.Errorf("expected status 'active', got %s", resp.Status)
				}
				if resp.VendorID != vendorID {
					t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
				}
			},
		},
		{
			name: "success from rejected",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "rejected")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				xenPlatform.EXPECT().CreateAccount(ctx, gomock.Any()).Return(&domain.XenPlatformAccount{
					ID:     "xnd_acc_456",
					Type:   "OWNED",
					Email:  "owner@example.com",
					Status: "LIVE",
				}, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				// Owner is already active, so Update should NOT be called.
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.Status != "active" {
					t.Errorf("expected status 'active', got %s", resp.Status)
				}
			},
		},
		{
			name: "xendit create account error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "pending"), nil)
				xenPlatform.EXPECT().CreateAccount(ctx, gomock.Any()).Return(nil, errors.New("xendit API error"))
			},
			wantErr: ErrXenditAccountCreation,
		},
		{
			name: "invalid status draft",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "draft"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status active",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "active"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status blocked",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "blocked"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "not found",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "FindByID error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "find owner error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "UpdateStatus error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "pending"), nil)
				xenPlatform.EXPECT().CreateAccount(ctx, gomock.Any()).Return(&domain.XenPlatformAccount{
					ID:     "xnd_acc_789",
					Type:   "OWNED",
					Email:  "owner@example.com",
					Status: "LIVE",
				}, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "activate owner error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, xenPlatform *mocks.MockXenPlatformProvider, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "pending"), nil)
				xenPlatform.EXPECT().CreateAccount(ctx, gomock.Any()).Return(&domain.XenPlatformAccount{
					ID:     "xnd_acc_789",
					Type:   "OWNED",
					Email:  "owner@example.com",
					Status: "LIVE",
				}, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				userRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vendorRepo, userRepo, _, xenPlatform, uc := setupAdminVendorUseCase(t)
			ctx := context.Background()
			vendorID := uuid.New()
			adminID := uuid.New()

			tc.setupMock(vendorRepo, userRepo, xenPlatform, ctx, vendorID, adminID)

			resp, err := uc.Approve(ctx, vendorID, adminID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp, vendorID)
			}
		})
	}
}

// ============================================================
// Reject
// ============================================================

func TestAdminVendorReject(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID)
		wantErr   error
		wantAny   bool
		checkResp func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID)
	}{
		{
			name: "success",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendor := dummyVendor(vendorID, "submitted")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.Status != "rejected" {
					t.Errorf("expected status 'rejected', got %s", resp.Status)
				}
				if resp.VendorID != vendorID {
					t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
				}
			},
		},
		{
			name: "invalid status draft",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "draft"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status active",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "active"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status rejected",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "rejected"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status blocked",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "blocked"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "not found",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "FindByID error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "UpdateStatus error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendor := dummyVendor(vendorID, "submitted")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vendorRepo, _, _, _, uc := setupAdminVendorUseCase(t)
			ctx := context.Background()
			vendorID := uuid.New()

			tc.setupMock(vendorRepo, ctx, vendorID)

			resp, err := uc.Reject(ctx, vendorID, "dokumen tidak lengkap")

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp, vendorID)
			}
		})
	}
}

// ============================================================
// Block
// ============================================================

func TestAdminVendorBlock(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID)
		wantErr   error
		wantAny   bool
		checkResp func(t *testing.T, resp *domain.AdminVendorActionResponse)
	}{
		{
			name: "success from active",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendor := dummyVendor(vendorID, "active")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse) {
				t.Helper()
				if resp.Status != "blocked" {
					t.Errorf("expected status 'blocked', got %s", resp.Status)
				}
			},
		},
		{
			name: "success from submitted",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendor := dummyVendor(vendorID, "submitted")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse) {
				t.Helper()
				if resp.Status != "blocked" {
					t.Errorf("expected status 'blocked', got %s", resp.Status)
				}
			},
		},
		{
			name: "invalid status draft",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "draft"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status rejected",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "rejected"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status blocked",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "blocked"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "not found",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "UpdateStatus error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, ctx context.Context, vendorID uuid.UUID) {
				vendor := dummyVendor(vendorID, "active")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vendorRepo, _, _, _, uc := setupAdminVendorUseCase(t)
			ctx := context.Background()
			vendorID := uuid.New()

			tc.setupMock(vendorRepo, ctx, vendorID)

			resp, err := uc.Block(ctx, vendorID, "pelanggaran kebijakan")

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp)
			}
		})
	}
}

// ============================================================
// Unblock
// ============================================================

func TestAdminVendorUnblock(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID)
		wantErr   error
		wantAny   bool
		checkResp func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID)
	}{
		{
			name: "success",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "blocked"), nil)
				userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.Status != "active" {
					t.Errorf("expected status 'active', got %s", resp.Status)
				}
				if resp.VendorID != vendorID {
					t.Errorf("expected vendor id %s, got %s", vendorID, resp.VendorID)
				}
			},
		},
		{
			name: "owner already active",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				// Owner is already active, so Update should NOT be called.
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorActionResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.Status != "active" {
					t.Errorf("expected status 'active', got %s", resp.Status)
				}
			},
		},
		{
			name: "invalid status draft",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "draft"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status submitted",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "submitted"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status active",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "active"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "invalid status rejected",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(dummyVendor(vendorID, "rejected"), nil)
			},
			wantErr: ErrInvalidStatusTransition,
		},
		{
			name: "not found",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)
			},
			wantErr: ErrVendorNotFound,
		},
		{
			name: "UpdateStatus error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				vendor := dummyVendor(vendorID, "blocked")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "find owner error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "activate owner error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context, vendorID, adminID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "blocked")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				vendorRepo.EXPECT().UpdateStatus(ctx, vendorID, gomock.Any()).Return(nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "blocked"), nil)
				userRepo.EXPECT().Update(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			wantAny: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vendorRepo, userRepo, _, _, uc := setupAdminVendorUseCase(t)
			ctx := context.Background()
			vendorID := uuid.New()
			adminID := uuid.New()

			tc.setupMock(vendorRepo, userRepo, ctx, vendorID, adminID)

			resp, err := uc.Unblock(ctx, vendorID, adminID)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.wantAny {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, resp, vendorID)
			}
		})
	}
}
