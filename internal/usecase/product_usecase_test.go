package usecase

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// ---------- helpers ----------

func setupProductUseCase(t *testing.T) (
	*mocks.MockProductRepository,
	*mocks.MockVendorRepository,
	*mocks.MockCategoryRepository,
	*mocks.MockStorageProvider,
	domain.ProductUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	productRepo := mocks.NewMockProductRepository(ctrl)
	vendorRepo := mocks.NewMockVendorRepository(ctrl)
	categoryRepo := mocks.NewMockCategoryRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	uc := NewProductUseCase(productRepo, vendorRepo, categoryRepo, storage)
	return productRepo, vendorRepo, categoryRepo, storage, uc
}

func activeVendor(id uuid.UUID) *domain.Vendor {
	return &domain.Vendor{
		ID:                    id,
		OwnerUserID:           uuid.New(),
		VendorType:            domain.VendorTypeSouvenirStore,
		DisplayName:           "Toko Haji",
		ResponsiblePersonName: "Ahmad",
		Status:                "active",
	}
}

func activeCategory(id uuid.UUID) *domain.Category {
	return &domain.Category{
		ID:       id,
		Name:     "Oleh-Oleh",
		Slug:     "oleh-oleh",
		IsActive: true,
	}
}

func validCreateProductRequest(categoryID string) domain.CreateProductRequest {
	weight := 500
	return domain.CreateProductRequest{
		Name:        "Sajadah Premium",
		CategoryID:  categoryID,
		Description: "Sajadah berkualitas tinggi",
		Price:       150000,
		Stock:       100,
		WeightGram:  &weight,
		IsActive:    false,
		Images: []domain.CreateProductImageInput{
			{
				FileName:    "sajadah-front.jpg",
				ContentType: "image/jpeg",
				IsPrimary:   true,
			},
		},
	}
}

// expectSlugUnique sets up expectations so generateSlug returns the base slug (no collision).
func expectSlugUnique(productRepo *mocks.MockProductRepository, ctx context.Context) {
	productRepo.EXPECT().FindBySlug(ctx, gomock.Any()).Return(nil, nil)
}

// expectSlugCollision sets up expectations so generateSlug finds a collision and appends a suffix.
func expectSlugCollision(productRepo *mocks.MockProductRepository, ctx context.Context) {
	productRepo.EXPECT().FindBySlug(ctx, gomock.Any()).Return(&domain.Product{}, nil)
	productRepo.EXPECT().CountBySlugPrefix(ctx, gomock.Any()).Return(int64(1), nil)
	productRepo.EXPECT().FindBySlug(ctx, gomock.Any()).Return(nil, nil)
}

// ============================================================
// Create
// ============================================================

func TestCreate_SuccessCases(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "draft",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()

				req := validCreateProductRequest(categoryID.String())
				req.IsActive = false

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/upload", nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				resp, err := uc.Create(ctx, vendorID, req)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.Status != "draft" {
					t.Errorf("expected status 'draft', got %s", resp.Status)
				}
				if resp.VendorID != vendorID {
					t.Errorf("expected vendor_id %s, got %s", vendorID, resp.VendorID)
				}
				if resp.CategoryID != categoryID {
					t.Errorf("expected category_id %s, got %s", categoryID, resp.CategoryID)
				}
				if resp.Name != "Sajadah Premium" {
					t.Errorf("expected name 'Sajadah Premium', got %s", resp.Name)
				}
				if resp.HalalAIStatus != "pending" {
					t.Errorf("expected halal_ai_status 'pending', got %s", resp.HalalAIStatus)
				}
				if len(resp.Variants) != 1 {
					t.Fatalf("expected 1 variant (default), got %d", len(resp.Variants))
				}
				if !resp.Variants[0].IsDefault {
					t.Error("expected default variant to have IsDefault=true")
				}
				if resp.Variants[0].Price != 150000 {
					t.Errorf("expected price 150000, got %f", resp.Variants[0].Price)
				}
				if len(resp.UploadURLs) != 1 {
					t.Fatalf("expected 1 upload URL, got %d", len(resp.UploadURLs))
				}
				if resp.UploadURLs[0].UploadURL != "https://presigned.example.com/upload" {
					t.Errorf("expected presigned upload URL, got %s", resp.UploadURLs[0].UploadURL)
				}
			},
		},
		{
			name: "published",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()

				req := validCreateProductRequest(categoryID.String())
				req.IsActive = true

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/upload", nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				resp, err := uc.Create(ctx, vendorID, req)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.Status != "published" {
					t.Errorf("expected status 'published', got %s", resp.Status)
				}
			},
		},
		{
			name: "with additional variants",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				weight1 := 600

				req := validCreateProductRequest(categoryID.String())
				req.Variants = []domain.CreateProductVariantInput{
					{VariantName: "Large", Price: 200000, Stock: 50, WeightGram: &weight1, IsActive: true},
					{VariantName: "XL", Price: 250000, Stock: 30, WeightGram: nil, IsActive: true},
				}

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(5), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/upload", nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				resp, err := uc.Create(ctx, vendorID, req)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(resp.Variants) != 3 {
					t.Fatalf("expected 3 variants, got %d", len(resp.Variants))
				}
				if !resp.Variants[0].IsDefault {
					t.Error("expected first variant to be default")
				}
				if resp.Variants[1].VariantName != "Large" {
					t.Errorf("expected variant name 'Large', got %s", resp.Variants[1].VariantName)
				}
				if resp.Variants[2].VariantName != "XL" {
					t.Errorf("expected variant name 'XL', got %s", resp.Variants[2].VariantName)
				}
				if resp.Variants[1].IsDefault || resp.Variants[2].IsDefault {
					t.Error("additional variants should not be default")
				}
			},
		},
		{
			name: "no images",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()

				req := validCreateProductRequest(categoryID.String())
				req.Images = nil

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				resp, err := uc.Create(ctx, vendorID, req)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if len(resp.UploadURLs) != 0 {
					t.Errorf("expected 0 upload URLs, got %d", len(resp.UploadURLs))
				}
			},
		},
		{
			name: "slug collision",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()

				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugCollision(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/upload", nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				resp, err := uc.Create(ctx, vendorID, req)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if resp.Slug == "" {
					t.Error("expected non-empty slug")
				}
			},
		},
		{
			name: "variant sku generation",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())
				req.Variants = []domain.CreateProductVariantInput{
					{VariantName: "Red", Price: 160000, Stock: 20, IsActive: true},
				}

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(3), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/upload", nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				resp, err := uc.Create(ctx, vendorID, req)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				defaultSKU := resp.Variants[0].SKU
				variantSKU := resp.Variants[1].SKU
				if defaultSKU == "" || variantSKU == "" {
					t.Error("expected non-empty SKUs")
				}
				if defaultSKU == variantSKU {
					t.Errorf("default and variant SKU should differ, both are %s", defaultSKU)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestCreate_ErrorCases(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "vendor not found",
			run: func(t *testing.T) {
				_, vendorRepo, _, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				req := validCreateProductRequest(uuid.New().String())
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, nil)

				_, err := uc.Create(ctx, vendorID, req)
				if !errors.Is(err, ErrVendorNotFound) {
					t.Errorf("expected ErrVendorNotFound, got %v", err)
				}
			},
		},
		{
			name: "vendor repo error",
			run: func(t *testing.T) {
				_, vendorRepo, _, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				req := validCreateProductRequest(uuid.New().String())
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(nil, errors.New("db error"))

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
		{
			name: "vendor not active",
			run: func(t *testing.T) {
				_, vendorRepo, _, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				req := validCreateProductRequest(uuid.New().String())
				for _, status := range []string{"draft", "submitted", "blocked", "rejected"} {
					vendor := activeVendor(vendorID)
					vendor.Status = status
					vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(vendor, nil)

					_, err := uc.Create(ctx, vendorID, req)
					if !errors.Is(err, ErrVendorNotActive) {
						t.Errorf("status=%s: expected ErrVendorNotActive, got %v", status, err)
					}
				}
			},
		},
		{
			name: "invalid category id",
			run: func(t *testing.T) {
				_, vendorRepo, _, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				req := validCreateProductRequest("not-a-uuid")
				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
		{
			name: "category not found",
			run: func(t *testing.T) {
				_, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(nil, nil)

				_, err := uc.Create(ctx, vendorID, req)
				if !errors.Is(err, ErrCategoryNotFound) {
					t.Errorf("expected ErrCategoryNotFound, got %v", err)
				}
			},
		},
		{
			name: "category not active",
			run: func(t *testing.T) {
				_, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				inactiveCat := activeCategory(categoryID)
				inactiveCat.IsActive = false

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(inactiveCat, nil)

				_, err := uc.Create(ctx, vendorID, req)
				if !errors.Is(err, ErrCategoryNotFound) {
					t.Errorf("expected ErrCategoryNotFound, got %v", err)
				}
			},
		},
		{
			name: "category repo error",
			run: func(t *testing.T) {
				_, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(nil, errors.New("db error"))

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
		{
			name: "too many images",
			run: func(t *testing.T) {
				_, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())
				images := make([]domain.CreateProductImageInput, 11)
				for i := range images {
					images[i] = domain.CreateProductImageInput{
						FileName:    "image.jpg",
						ContentType: "image/jpeg",
						IsPrimary:   i == 0,
					}
				}
				req.Images = images

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)

				_, err := uc.Create(ctx, vendorID, req)
				if !errors.Is(err, ErrTooManyImages) {
					t.Errorf("expected ErrTooManyImages, got %v", err)
				}
			},
		},
		{
			name: "no primary image",
			run: func(t *testing.T) {
				_, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())
				req.Images = []domain.CreateProductImageInput{
					{FileName: "img.jpg", ContentType: "image/jpeg", IsPrimary: false},
				}

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)

				_, err := uc.Create(ctx, vendorID, req)
				if !errors.Is(err, ErrNoPrimaryImage) {
					t.Errorf("expected ErrNoPrimaryImage, got %v", err)
				}
			},
		},
		{
			name: "duplicate primary image",
			run: func(t *testing.T) {
				_, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())
				req.Images = []domain.CreateProductImageInput{
					{FileName: "img1.jpg", ContentType: "image/jpeg", IsPrimary: true},
					{FileName: "img2.jpg", ContentType: "image/jpeg", IsPrimary: true},
				}

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)

				_, err := uc.Create(ctx, vendorID, req)
				if !errors.Is(err, ErrDuplicatePrimaryImage) {
					t.Errorf("expected ErrDuplicatePrimaryImage, got %v", err)
				}
			},
		},
		{
			name: "slug generation error",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				productRepo.EXPECT().FindBySlug(ctx, gomock.Any()).Return(nil, errors.New("db error"))

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
		{
			name: "count by vendor id error",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, _, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), errors.New("db error"))

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
		{
			name: "presigned upload url error",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("", errors.New("s3 error"))

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
		{
			name: "product repo create error",
			run: func(t *testing.T) {
				productRepo, vendorRepo, categoryRepo, storage, uc := setupProductUseCase(t)
				ctx := context.Background()

				vendorID := uuid.New()
				categoryID := uuid.New()
				req := validCreateProductRequest(categoryID.String())

				vendorRepo.EXPECT().FindByID(ctx, vendorID).Return(activeVendor(vendorID), nil)
				categoryRepo.EXPECT().FindByID(ctx, categoryID).Return(activeCategory(categoryID), nil)
				expectSlugUnique(productRepo, ctx)
				productRepo.EXPECT().CountByVendorID(ctx, vendorID).Return(int64(0), nil)
				storage.EXPECT().GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).Return("https://presigned.example.com/upload", nil)
				productRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))

				_, err := uc.Create(ctx, vendorID, req)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

// ============================================================
// ConfirmImages
// ============================================================

func TestConfirmImages_Success(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/img1.jpg"

	product := &domain.Product{
		ID:       productID,
		VendorID: vendorID,
		Status:   "draft",
	}

	existingImages := []domain.ProductImage{
		{
			ID:        imageID,
			ProductID: productID,
			ImageURL:  objectKey,
			IsPrimary: true,
		},
	}

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		Key:           objectKey,
		ContentType:   "image/jpeg",
		ContentLength: 1024,
	}, nil)
	productRepo.EXPECT().UpdateImages(ctx, gomock.Any()).Return(nil)

	resp, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ProductID != productID {
		t.Errorf("expected product_id %s, got %s", productID, resp.ProductID)
	}
	if resp.ImagesConfirmed != 1 {
		t.Errorf("expected 1 image confirmed, got %d", resp.ImagesConfirmed)
	}
}

func TestConfirmImages_MultipleImages(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID1 := uuid.New()
	imageID2 := uuid.New()
	objectKey1 := "products/img1.jpg"
	objectKey2 := "products/img2.png"

	product := &domain.Product{
		ID:       productID,
		VendorID: vendorID,
	}

	existingImages := []domain.ProductImage{
		{ID: imageID1, ProductID: productID, ImageURL: objectKey1, IsPrimary: true},
		{ID: imageID2, ProductID: productID, ImageURL: objectKey2, IsPrimary: false},
	}

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID1.String(), ObjectKey: objectKey1},
			{ImageID: imageID2.String(), ObjectKey: objectKey2},
		},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey1).Return(&domain.ObjectInfo{
		ContentType:   "image/jpeg",
		ContentLength: 2048,
	}, nil)
	storage.EXPECT().HeadObject(ctx, objectKey2).Return(&domain.ObjectInfo{
		ContentType:   "image/png",
		ContentLength: 4096,
	}, nil)
	productRepo.EXPECT().UpdateImages(ctx, gomock.Any()).Return(nil)

	resp, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ImagesConfirmed != 2 {
		t.Errorf("expected 2 images confirmed, got %d", resp.ImagesConfirmed)
	}
}

func TestConfirmImages_ProductNotFound(t *testing.T) {
	productRepo, _, _, _, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()

	productRepo.EXPECT().FindByID(ctx, productID).Return(nil, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: uuid.New().String(), ObjectKey: "products/img.jpg"},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if !errors.Is(err, ErrProductNotFound) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestConfirmImages_FindByIDError(t *testing.T) {
	productRepo, _, _, _, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()

	productRepo.EXPECT().FindByID(ctx, productID).Return(nil, errors.New("db error"))

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: uuid.New().String(), ObjectKey: "products/img.jpg"},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmImages_ProductNotOwned(t *testing.T) {
	productRepo, _, _, _, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	otherVendorID := uuid.New()
	productID := uuid.New()

	product := &domain.Product{
		ID:       productID,
		VendorID: otherVendorID,
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: uuid.New().String(), ObjectKey: "products/img.jpg"},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if !errors.Is(err, ErrProductNotOwned) {
		t.Errorf("expected ErrProductNotOwned, got %v", err)
	}
}

func TestConfirmImages_FindImagesError(t *testing.T) {
	productRepo, _, _, _, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()

	product := &domain.Product{ID: productID, VendorID: vendorID}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(nil, errors.New("db error"))

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: uuid.New().String(), ObjectKey: "products/img.jpg"},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmImages_InvalidImageID(t *testing.T) {
	productRepo, _, _, _, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()

	product := &domain.Product{ID: productID, VendorID: vendorID}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return([]domain.ProductImage{}, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: "not-a-uuid", ObjectKey: "products/img.jpg"},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err == nil {
		t.Fatal("expected error for invalid image_id, got nil")
	}
}

func TestConfirmImages_ImageNotFound(t *testing.T) {
	productRepo, _, _, _, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	unknownImageID := uuid.New()

	product := &domain.Product{ID: productID, VendorID: vendorID}

	// Empty image list — requested image doesn't exist
	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return([]domain.ProductImage{}, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: unknownImageID.String(), ObjectKey: "products/img.jpg"},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if !errors.Is(err, ErrImageNotFound) {
		t.Errorf("expected ErrImageNotFound, got %v", err)
	}
}

func TestConfirmImages_HeadObjectError(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/img.jpg"

	product := &domain.Product{ID: productID, VendorID: vendorID}
	existingImages := []domain.ProductImage{
		{ID: imageID, ProductID: productID, ImageURL: objectKey},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(nil, errors.New("s3 error"))

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmImages_ImageNotUploaded(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/img.jpg"

	product := &domain.Product{ID: productID, VendorID: vendorID}
	existingImages := []domain.ProductImage{
		{ID: imageID, ProductID: productID, ImageURL: objectKey},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(nil, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if !errors.Is(err, ErrImageNotUploaded) {
		t.Errorf("expected ErrImageNotUploaded, got %v", err)
	}
}

func TestConfirmImages_InvalidContentType(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/img.txt"

	product := &domain.Product{ID: productID, VendorID: vendorID}
	existingImages := []domain.ProductImage{
		{ID: imageID, ProductID: productID, ImageURL: objectKey},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType:   "text/plain",
		ContentLength: 1024,
	}, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if !errors.Is(err, ErrInvalidImageContentType) {
		t.Errorf("expected ErrInvalidImageContentType, got %v", err)
	}
}

func TestConfirmImages_ContentTypeNormalization(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/img.jpg"

	product := &domain.Product{ID: productID, VendorID: vendorID}
	existingImages := []domain.ProductImage{
		{ID: imageID, ProductID: productID, ImageURL: objectKey},
	}

	// Content type with charset suffix should be normalized
	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType:   "  IMAGE/JPEG; charset=utf-8  ",
		ContentLength: 1024,
	}, nil)
	productRepo.EXPECT().UpdateImages(ctx, gomock.Any()).Return(nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	resp, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ImagesConfirmed != 1 {
		t.Errorf("expected 1 image confirmed, got %d", resp.ImagesConfirmed)
	}
}

func TestConfirmImages_ImageSizeOverflow(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/huge.jpg"

	product := &domain.Product{ID: productID, VendorID: vendorID}
	existingImages := []domain.ProductImage{
		{ID: imageID, ProductID: productID, ImageURL: objectKey},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType:   "image/jpeg",
		ContentLength: int64(math.MaxInt32) + 1,
	}, nil)

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if !errors.Is(err, ErrImageSizeOverflow) {
		t.Errorf("expected ErrImageSizeOverflow, got %v", err)
	}
}

func TestConfirmImages_UpdateImagesError(t *testing.T) {
	productRepo, _, _, storage, uc := setupProductUseCase(t)
	ctx := context.Background()

	vendorID := uuid.New()
	productID := uuid.New()
	imageID := uuid.New()
	objectKey := "products/img.jpg"

	product := &domain.Product{ID: productID, VendorID: vendorID}
	existingImages := []domain.ProductImage{
		{ID: imageID, ProductID: productID, ImageURL: objectKey},
	}

	productRepo.EXPECT().FindByID(ctx, productID).Return(product, nil)
	productRepo.EXPECT().FindImagesByProductID(ctx, productID).Return(existingImages, nil)
	storage.EXPECT().HeadObject(ctx, objectKey).Return(&domain.ObjectInfo{
		ContentType:   "image/jpeg",
		ContentLength: 1024,
	}, nil)
	productRepo.EXPECT().UpdateImages(ctx, gomock.Any()).Return(errors.New("db error"))

	req := domain.ConfirmProductImagesRequest{
		Images: []domain.ConfirmProductImageItem{
			{ImageID: imageID.String(), ObjectKey: objectKey},
		},
	}

	_, err := uc.ConfirmImages(ctx, vendorID, productID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ============================================================
// Helper functions (unit tests)
// ============================================================

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Sajadah Premium", "sajadah-premium"},
		{"Hello  World", "hello-world"},
		{"Product (Special)", "product-special"},
		{"ABC", "abc"},
		{"---test---", "test"},
		{"123 numbers", "123-numbers"},
	}

	for _, tc := range tests {
		got := slugify(tc.input)
		if got != tc.expected {
			t.Errorf("slugify(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestGenerateSKU(t *testing.T) {
	tests := []struct {
		vendorName  string
		productName string
		productSeq  int
		variantIdx  int
		expected    string
	}{
		{"Toko Haji", "Sajadah Premium", 0, 0, "TOKOH-SAJAD000-00"},
		{"Toko Haji", "Sajadah Premium", 5, 1, "TOKOH-SAJAD005-01"},
		{"AB", "CD", 10, 2, "ABXXX-CDXXX010-02"},
		{"", "", 0, 0, "PRODX-PRODX000-00"},
	}

	for _, tc := range tests {
		got := generateSKU(tc.vendorName, tc.productName, tc.productSeq, tc.variantIdx)
		if got != tc.expected {
			t.Errorf("generateSKU(%q, %q, %d, %d) = %q, want %q",
				tc.vendorName, tc.productName, tc.productSeq, tc.variantIdx, got, tc.expected)
		}
	}
}

func TestExtractCode(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		expected string
	}{
		{"Toko Haji", 5, "TOKOH"},
		{"AB", 5, "ABXXX"},
		{"", 5, "PRODX"},
		{"A-B_C!D@E#F", 5, "ABCDE"},
		{"ABCDEFGH", 3, "ABC"},
	}

	for _, tc := range tests {
		got := extractCode(tc.name, tc.length)
		if got != tc.expected {
			t.Errorf("extractCode(%q, %d) = %q, want %q", tc.name, tc.length, got, tc.expected)
		}
	}
}

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"photo.jpg", "photo.jpg"},
		{"my image (1).png", "myimage1.png"},
		{"hello world!.webp", "helloworld.webp"},
		{"file-name_v2.jpg", "file-name_v2.jpg"},
		{"", "image"},
		{"!!!###", "image"},
	}

	for _, tc := range tests {
		got := sanitizeFileName(tc.input)
		if got != tc.expected {
			t.Errorf("sanitizeFileName(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestNormalizeImageContentType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"image/jpeg", "image/jpeg"},
		{"  IMAGE/JPEG  ", "image/jpeg"},
		{"image/png; charset=utf-8", "image/png"},
		{"IMAGE/WEBP; boundary=something", "image/webp"},
	}

	for _, tc := range tests {
		got := normalizeContentType(tc.input)
		if got != tc.expected {
			t.Errorf("normalizeContentType(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestIsAllowedImageContentType(t *testing.T) {
	allowed := []string{"image/jpeg", "image/png", "image/webp"}
	for _, ct := range allowed {
		if !isAllowedImageContentType(ct) {
			t.Errorf("expected %q to be allowed", ct)
		}
	}

	disallowed := []string{"text/plain", "image/gif", "image/bmp", "application/octet-stream"}
	for _, ct := range disallowed {
		if isAllowedImageContentType(ct) {
			t.Errorf("expected %q to be disallowed", ct)
		}
	}
}
