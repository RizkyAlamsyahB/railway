package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type productUseCase struct {
	productRepo  domain.ProductRepository
	vendorRepo   domain.VendorRepository
	categoryRepo domain.CategoryRepository
	storage      domain.StorageProvider
}

// NewProductUseCase creates a new ProductUseCase.
func NewProductUseCase(
	productRepo domain.ProductRepository,
	vendorRepo domain.VendorRepository,
	categoryRepo domain.CategoryRepository,
	storage domain.StorageProvider,
) domain.ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		vendorRepo:   vendorRepo,
		categoryRepo: categoryRepo,
		storage:      storage,
	}
}

func (uc *productUseCase) Create(ctx context.Context, vendorID uuid.UUID, req domain.CreateProductRequest) (*domain.CreateProductResponse, error) {
	// 1. Validate vendor is active.
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}
	if vendor.Status != domain.VendorStatusActive {
		return nil, ErrVendorNotActive
	}

	// 2. Parse and validate category.
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category_id: %w", err)
	}
	category, err := uc.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil || !category.IsActive {
		return nil, ErrCategoryNotFound
	}

	// 3. Determine product status.
	status := domain.ProductStatusDraft
	if req.IsActive {
		status = domain.ProductStatusPublished
	}

	// 4. Validate images.
	if len(req.Images) > 10 {
		return nil, ErrTooManyImages
	}
	if len(req.Images) > 0 {
		primaryCount := 0
		for _, img := range req.Images {
			if img.IsPrimary {
				primaryCount++
			}
		}
		if primaryCount == 0 {
			return nil, ErrNoPrimaryImage
		}
		if primaryCount > 1 {
			return nil, ErrDuplicatePrimaryImage
		}
	}

	// 5. Generate slug.
	slug, err := uc.generateSlug(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate slug: %w", err)
	}

	// 6. Get product count for SKU generation.
	productCount, err := uc.productRepo.CountByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to count vendor products: %w", err)
	}

	now := time.Now()
	productID := uuid.New()

	// 7. Build default variant.
	defaultSKU := generateSKU(vendor.DisplayName, req.Name, int(productCount), 0)
	defaultVariant := domain.ProductVariant{
		ID:          uuid.New(),
		ProductID:   productID,
		SKU:         defaultSKU,
		VariantName: "Default",
		Price:       req.Price,
		Currency:    "IDR",
		StockOnHand: req.Stock,
		WeightGram:  req.WeightGram,
		IsDefault:   true,
		IsActive:    true,
	}

	allVariants := []domain.ProductVariant{defaultVariant}

	// 8. Build additional variants.
	for i, v := range req.Variants {
		sku := generateSKU(vendor.DisplayName, req.Name, int(productCount), i+1)
		variant := domain.ProductVariant{
			ID:          uuid.New(),
			ProductID:   productID,
			SKU:         sku,
			VariantName: v.VariantName,
			Price:       v.Price,
			Currency:    "IDR",
			StockOnHand: v.Stock,
			WeightGram:  v.WeightGram,
			IsDefault:   false,
			IsActive:    v.IsActive,
		}
		allVariants = append(allVariants, variant)
	}

	// 9. Generate presigned upload URLs for images.
	uploadInfos := make([]domain.ProductImageUploadInfo, 0, len(req.Images))
	imageEntities := make([]domain.ProductImage, 0, len(req.Images))

	for i, img := range req.Images {
		imageID := uuid.New()
		objectKey := fmt.Sprintf("products/%s/images/%s/%s",
			productID.String(), imageID.String(), sanitizeFileName(img.FileName))

		uploadURL, err := uc.storage.GeneratePresignedUploadURL(
			ctx, objectKey, img.ContentType, PresignedUploadExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for image %d: %w", i, err)
		}

		imageEntity := domain.ProductImage{
			ID:        imageID,
			ProductID: productID,
			ImageURL:  objectKey,
			IsPrimary: img.IsPrimary,
			SortOrder: i,
			CreatedAt: now,
		}
		imageEntities = append(imageEntities, imageEntity)

		uploadInfos = append(uploadInfos, domain.ProductImageUploadInfo{
			ImageID:   imageID,
			UploadURL: uploadURL,
			ObjectKey: objectKey,
			SortOrder: i,
			IsPrimary: img.IsPrimary,
		})
	}

	// 10. Build product entity.
	product := &domain.Product{
		ID:            productID,
		VendorID:      vendorID,
		CategoryID:    categoryID,
		Name:          req.Name,
		Slug:          slug,
		Description:   req.Description,
		Status:        status,
		HalalAIStatus: domain.HalalAIStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 11. Persist everything in one transaction.
	if err := uc.productRepo.Create(ctx, product, allVariants, imageEntities); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// 12. Build response.
	variantResponses := make([]domain.ProductVariantResponse, len(allVariants))
	for i, v := range allVariants {
		variantResponses[i] = domain.ProductVariantResponse{
			ID:          v.ID,
			SKU:         v.SKU,
			VariantName: v.VariantName,
			Price:       v.Price,
			Currency:    v.Currency,
			StockOnHand: v.StockOnHand,
			WeightGram:  v.WeightGram,
			IsDefault:   v.IsDefault,
			IsActive:    v.IsActive,
		}
	}

	return &domain.CreateProductResponse{
		ID:            productID,
		VendorID:      vendorID,
		CategoryID:    categoryID,
		Name:          req.Name,
		Slug:          slug,
		Description:   req.Description,
		Status:        status,
		HalalAIStatus: domain.HalalAIStatusPending,
		Variants:      variantResponses,
		UploadURLs:    uploadInfos,
		CreatedAt:     now,
	}, nil
}

func (uc *productUseCase) ConfirmImages(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID, req domain.ConfirmProductImagesRequest) (*domain.ConfirmProductImagesResponse, error) {
	// 1. Verify product exists and belongs to vendor.
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	if product.VendorID != vendorID {
		return nil, ErrProductNotOwned
	}

	// 2. Load existing images.
	existingImages, err := uc.productRepo.FindImagesByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find images: %w", err)
	}

	imageMap := make(map[uuid.UUID]*domain.ProductImage, len(existingImages))
	for i := range existingImages {
		imageMap[existingImages[i].ID] = &existingImages[i]
	}

	// 3. Validate each image and verify it exists in S3.
	updatedImages := make([]domain.ProductImage, 0, len(req.Images))

	for _, item := range req.Images {
		imageID, err := uuid.Parse(item.ImageID)
		if err != nil {
			return nil, fmt.Errorf("invalid image_id: %w", err)
		}

		img, exists := imageMap[imageID]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrImageNotFound, item.ImageID)
		}

		info, err := uc.storage.HeadObject(ctx, item.ObjectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to verify image %s: %w", item.ObjectKey, err)
		}
		if info == nil {
			return nil, fmt.Errorf("%w: %s", ErrImageNotUploaded, item.ObjectKey)
		}

		contentType := normalizeContentType(info.ContentType)
		if !isAllowedImageContentType(contentType) {
			return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidImageContentType, item.ImageID, info.ContentType)
		}
		if info.ContentLength > int64(math.MaxInt32) {
			return nil, fmt.Errorf("%w: %s (%d bytes)", ErrImageSizeOverflow, item.ImageID, info.ContentLength)
		}

		fileSize := int(info.ContentLength)
		img.ImageURL = item.ObjectKey
		img.MimeType = &contentType
		img.FileSizeBytes = &fileSize

		updatedImages = append(updatedImages, *img)
	}

	// 4. Persist updates.
	if err := uc.productRepo.UpdateImages(ctx, updatedImages); err != nil {
		return nil, fmt.Errorf("failed to confirm images: %w", err)
	}

	return &domain.ConfirmProductImagesResponse{
		ProductID:       productID,
		ImagesConfirmed: len(updatedImages),
	}, nil
}

// generateSlug creates a URL-friendly slug from a product name, ensuring uniqueness.
func (uc *productUseCase) generateSlug(ctx context.Context, name string) (string, error) {
	base := slugify(name)
	if len(base) > 200 {
		base = base[:200]
	}

	existing, err := uc.productRepo.FindBySlug(ctx, base)
	if err != nil {
		return "", err
	}
	if existing == nil {
		return base, nil
	}

	count, err := uc.productRepo.CountBySlugPrefix(ctx, base)
	if err != nil {
		return "", err
	}

	for i := count; i < count+100; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i+1)
		existing, err := uc.productRepo.FindBySlug(ctx, candidate)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return candidate, nil
		}
	}

	return fmt.Sprintf("%s-%s", base, uuid.New().String()[:8]), nil
}

// slugify converts a name into a URL-friendly slug.
func slugify(name string) string {
	slug := strings.ToLower(name)

	var b strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	slug = b.String()

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	slug = strings.Trim(slug, "-")
	return slug
}

// generateSKU produces a unique SKU for a product variant.
// Format: {VENDOR_5CHARS}-{PRODUCT_5CHARS}{SEQ_3DIGITS}-{VARIANT_2DIGITS}
func generateSKU(vendorName, productName string, productSeq int, variantIdx int) string {
	vendorCode := extractCode(vendorName, 5)
	productCode := extractCode(productName, 5)
	return fmt.Sprintf("%s-%s%03d-%02d", vendorCode, productCode, productSeq, variantIdx)
}

// extractCode extracts an uppercase code of the given length from a name.
func extractCode(name string, length int) string {
	upper := strings.ToUpper(name)

	var alphanumeric []rune
	for _, r := range upper {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			alphanumeric = append(alphanumeric, r)
		}
	}

	if len(alphanumeric) == 0 {
		alphanumeric = []rune("PRODX")
	}

	if len(alphanumeric) >= length {
		return string(alphanumeric[:length])
	}

	result := string(alphanumeric)
	for len(result) < length {
		result += "X"
	}
	return result
}

// sanitizeFileName strips unsafe characters from file names.
func sanitizeFileName(name string) string {
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			b.WriteRune(r)
		}
	}
	result := b.String()
	if result == "" {
		return "image"
	}
	return result
}
