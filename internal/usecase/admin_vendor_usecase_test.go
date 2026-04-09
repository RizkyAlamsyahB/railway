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
	vendorType := domain.VendorTypeSouvenirStore
	displayName := "Toko Oleh-Oleh Haji"
	provinceID := "31"
	provinceName := "DKI Jakarta"
	cityID := "3171"
	cityName := "Jakarta Pusat"
	districtID := "317101"
	districtName := "Menteng"
	subdistrictID := "3171011001"
	subdistrictName := "Pegangsaan"
	postalCode := "10320"
	addressLine := "Jl. Pegangsaan Barat No. 12"
	return &domain.Vendor{
		ID:              id,
		OwnerUserID:     ownerID,
		VendorType:      &vendorType,
		DisplayName:     &displayName,
		LegalName:       ptrString("PT Toko Haji"),
		Description:     ptrString("Toko oleh-oleh haji terlengkap"),
		ProvinceID:      &provinceID,
		ProvinceName:    &provinceName,
		CityID:          &cityID,
		CityName:        &cityName,
		DistrictID:      &districtID,
		DistrictName:    &districtName,
		SubdistrictID:   &subdistrictID,
		SubdistrictName: &subdistrictName,
		PostalCode:      &postalCode,
		AddressLine:     &addressLine,
		Status:          status,
		CreatedAt:       now,
		UpdatedAt:       now,
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
		ID:                uuid.New(),
		VendorID:          vendorID,
		BankName:          "BCA",
		AccountNumber:     "1234567890",
		AccountHolderName: "Ahmad",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func dummyDocuments(vendorID uuid.UUID) []domain.VendorDocument {
	now := time.Now()
	uploaderID := uuid.New()
	return []domain.VendorDocument{
		{
			ID:            uuid.New(),
			VendorID:      vendorID,
			DocType:       domain.VendorDocumentTypeOwnerDocumentID,
			FileURL:       "vendors/doc1.jpg",
			MimeType:      ptrString("image/jpeg"),
			FileSizeBytes: ptrInt(1024),
			UploadedBy:    &uploaderID,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			VendorID:      vendorID,
			DocType:       domain.VendorDocumentTypeBusinessNIB,
			FileURL:       "vendors/nib.pdf",
			MimeType:      ptrString("application/pdf"),
			FileSizeBytes: ptrInt(2048),
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			VendorID:      vendorID,
			DocType:       domain.VendorDocumentTypeHalalCertificate,
			FileURL:       "vendors/halal.pdf",
			MimeType:      ptrString("application/pdf"),
			FileSizeBytes: ptrInt(4096),
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
}

func dummyResponsiblePerson(vendorID uuid.UUID, userID uuid.UUID) *domain.VendorResponsiblePerson {
	now := time.Now()
	return &domain.VendorResponsiblePerson{
		ID:        uuid.New(),
		VendorID:  vendorID,
		UserID:    userID,
		NIK:       "3173010101010001",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func ptrInt(v int) *int { return &v }

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

				vendors := []domain.Vendor{
					*dummyVendorWithOwner(vendorID1, ownerID1, "submitted"),
					*dummyVendorWithOwner(vendorID2, ownerID2, "active"),
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
				if items[0].StoreName == nil || *items[0].StoreName != "Toko Oleh-Oleh Haji" {
					t.Errorf("expected store_name 'Toko Oleh-Oleh Haji', got %v", items[0].StoreName)
				}
				if items[0].StoreType == nil || *items[0].StoreType != domain.VendorTypeSouvenirStore {
					t.Errorf("expected store_type '%s', got %v", domain.VendorTypeSouvenirStore, items[0].StoreType)
				}
				if items[0].Email != "owner@example.com" {
					t.Errorf("expected owner email 'owner@example.com', got %s", items[0].Email)
				}
				if items[0].Address == nil || items[0].Address.Province == nil || *items[0].Address.Province != "DKI Jakarta" {
					t.Errorf("expected address.province DKI Jakarta, got %+v", items[0].Address)
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
				if items[0].Email != "" {
					t.Errorf("expected empty email, got %s", items[0].Email)
				}
			},
		},
		{
			name:   "nullable fields for draft-like vendor",
			params: domain.VendorListParams{Page: 1, Limit: 10},
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, ctx context.Context) {
				vendorID := uuid.New()
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "draft")
				vendor.VendorType = nil
				vendor.DisplayName = nil
				vendor.ProvinceName = nil
				vendor.CityName = nil
				vendor.DistrictName = nil
				vendor.SubdistrictName = nil
				vendor.PostalCode = nil

				vendorRepo.EXPECT().List(ctx, domain.VendorListParams{Page: 1, Limit: 10}).Return([]domain.Vendor{*vendor}, int64(1), nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminVendorListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 1 {
					t.Fatalf("expected 1 item, got %d", len(items))
				}
				if items[0].StoreType != nil {
					t.Errorf("expected store_type nil, got %v", *items[0].StoreType)
				}
				if items[0].StoreName != nil {
					t.Errorf("expected store_name nil, got %v", *items[0].StoreName)
				}
				if items[0].Address != nil {
					t.Errorf("expected address nil, got %+v", items[0].Address)
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
				vendorRepo.EXPECT().FindResponsiblePersonByVendorID(ctx, vendorID).Return(dummyResponsiblePerson(vendorID, ownerID), nil)

				docs := dummyDocuments(vendorID)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(docs, nil)
				storage.EXPECT().GeneratePresignedURL(ctx, "vendors/doc1.jpg", PresignedDownloadExpiry).Return("https://signed-url.example.com/doc1.jpg", nil)
				storage.EXPECT().GeneratePresignedURL(ctx, "vendors/nib.pdf", PresignedDownloadExpiry).Return("https://signed-url.example.com/nib.pdf", nil)
				storage.EXPECT().GeneratePresignedURL(ctx, "vendors/halal.pdf", PresignedDownloadExpiry).Return("https://signed-url.example.com/halal.pdf", nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorDetailResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.ID != vendorID {
					t.Errorf("expected vendor id %s, got %s", vendorID, resp.ID)
				}
				if resp.StoreName == nil || *resp.StoreName != "Toko Oleh-Oleh Haji" {
					t.Errorf("expected store_name 'Toko Oleh-Oleh Haji', got %v", resp.StoreName)
				}
				if resp.Email != "owner@example.com" {
					t.Errorf("expected owner email 'owner@example.com', got %s", resp.Email)
				}
				if resp.Address.AddressLine == nil || *resp.Address.AddressLine != "Jl. Pegangsaan Barat No. 12" {
					t.Errorf("expected structured address line, got %v", resp.Address.AddressLine)
				}
				if resp.ResponsiblePerson.Name == nil || *resp.ResponsiblePerson.Name != "Owner Name" {
					t.Errorf("expected responsible person name 'Owner Name', got %v", resp.ResponsiblePerson.Name)
				}
				if resp.ResponsiblePerson.NIK == nil || *resp.ResponsiblePerson.NIK != "3173010101010001" {
					t.Errorf("expected responsible person NIK, got %v", resp.ResponsiblePerson.NIK)
				}
				if resp.ResponsiblePerson.KTP == nil || resp.ResponsiblePerson.KTP.DocType != domain.VendorDocumentTypeOwnerDocumentID {
					t.Errorf("expected KTP doc_type %s, got %+v", domain.VendorDocumentTypeOwnerDocumentID, resp.ResponsiblePerson.KTP)
				}
				if resp.ResponsiblePerson.KTP != nil && resp.ResponsiblePerson.KTP.FileURL != "https://signed-url.example.com/doc1.jpg" {
					t.Errorf("expected signed KTP file_url, got %s", resp.ResponsiblePerson.KTP.FileURL)
				}
				if len(resp.RequirementsDocuments) != 2 {
					t.Errorf("expected 2 requirement documents, got %d", len(resp.RequirementsDocuments))
				}
				if len(resp.RequirementsDocuments) == 2 {
					requiredURLs := map[string]string{}
					for _, doc := range resp.RequirementsDocuments {
						requiredURLs[doc.DocType] = doc.FileURL
					}
					if requiredURLs[domain.VendorDocumentTypeBusinessNIB] != "https://signed-url.example.com/nib.pdf" {
						t.Errorf("expected signed NIB file_url, got %s", requiredURLs[domain.VendorDocumentTypeBusinessNIB])
					}
					if requiredURLs[domain.VendorDocumentTypeHalalCertificate] != "https://signed-url.example.com/halal.pdf" {
						t.Errorf("expected signed halal file_url, got %s", requiredURLs[domain.VendorDocumentTypeHalalCertificate])
					}
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
			name: "responsible person error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindResponsiblePersonByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))
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
				vendorRepo.EXPECT().FindResponsiblePersonByVendorID(ctx, vendorID).Return(dummyResponsiblePerson(vendorID, ownerID), nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(nil, errors.New("db error"))
			},
			wantAny: true,
		},
		{
			name: "presigned URL error",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "submitted")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindResponsiblePersonByVendorID(ctx, vendorID).Return(dummyResponsiblePerson(vendorID, ownerID), nil)
				docs := dummyDocuments(vendorID)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return(docs, nil)
				storage.EXPECT().GeneratePresignedURL(ctx, "vendors/doc1.jpg", PresignedDownloadExpiry).Return("", errors.New("s3 error"))
			},
			wantAny: true,
		},
		{
			name: "no responsible person and no documents",
			setupMock: func(vendorRepo *mocks.MockVendorRepository, userRepo *mocks.MockUserRepository, storage *mocks.MockStorageProvider, ctx context.Context, vendorID uuid.UUID) {
				ownerID := uuid.New()
				vendor := dummyVendorWithOwner(vendorID, ownerID, "draft")

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)
				userRepo.EXPECT().FindByID(ctx, ownerID).Return(dummyOwner(ownerID, "active"), nil)
				vendorRepo.EXPECT().FindResponsiblePersonByVendorID(ctx, vendorID).Return(nil, nil)
				vendorRepo.EXPECT().FindDocumentsByVendorID(ctx, vendorID).Return([]domain.VendorDocument{}, nil)
			},
			checkResp: func(t *testing.T, resp *domain.AdminVendorDetailResponse, vendorID uuid.UUID) {
				t.Helper()
				if resp.ResponsiblePerson.KTP != nil {
					t.Errorf("expected nil KTP, got %+v", resp.ResponsiblePerson.KTP)
				}
				if len(resp.RequirementsDocuments) != 0 {
					t.Errorf("expected no requirement documents, got %d", len(resp.RequirementsDocuments))
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
