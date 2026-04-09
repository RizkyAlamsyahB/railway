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
	addressRepo  domain.AddressRepository
	storage      domain.StorageProvider
}

// NewProductUseCase creates a new ProductUseCase.
func NewProductUseCase(
	productRepo domain.ProductRepository,
	vendorRepo domain.VendorRepository,
	categoryRepo domain.CategoryRepository,
	addressRepo domain.AddressRepository,
	storage domain.StorageProvider,
) domain.ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		vendorRepo:   vendorRepo,
		categoryRepo: categoryRepo,
		addressRepo:  addressRepo,
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

	// 2. Validate vendor has a warehouse address (default address with district).
	warehouseAddr, err := uc.addressRepo.FindDefaultByUserID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check vendor warehouse address: %w", err)
	}
	if warehouseAddr == nil || warehouseAddr.DistrictID == nil || *warehouseAddr.DistrictID == "" {
		return nil, ErrVendorNoWarehouseAddress
	}

	// 3. Parse and validate category.
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
	defaultSKU := generateSKU(vendorDisplayNameOrEmpty(vendor), req.Name, int(productCount), 0)
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
		sku := generateSKU(vendorDisplayNameOrEmpty(vendor), req.Name, int(productCount), i+1)
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
		ProductType:   productType(req.ProductType),
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
			ID:            v.ID,
			SKU:           v.SKU,
			VariantName:   v.VariantName,
			Price:         v.Price,
			OriginalPrice: v.Price,
			PromoPrice:    nil,
			HasPromo:      false,
			Currency:      v.Currency,
			StockOnHand:   v.StockOnHand,
			WeightGram:    v.WeightGram,
			IsDefault:     v.IsDefault,
			IsActive:      v.IsActive,
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

func (uc *productUseCase) UpdateProduct(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID, req domain.UpdateProductRequest) (*domain.UpdateProductResponse, error) {
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

	// 2. Load and authorize product.
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

	// 3. Validate category.
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

	// 4. Determine new status.
	status := domain.ProductStatusDraft
	if req.IsActive {
		status = domain.ProductStatusPublished
	}

	// 5. Conditional slug regeneration.
	slug := product.Slug
	if req.Name != product.Name {
		newSlug, err := uc.generateSlug(ctx, req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to generate slug: %w", err)
		}
		slug = newSlug
	}

	// 6. Load existing variants.
	existingVariants, err := uc.productRepo.FindVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to load variants: %w", err)
	}
	existingVariantMap := make(map[uuid.UUID]*domain.ProductVariant, len(existingVariants))
	var defaultVariant domain.ProductVariant
	for i := range existingVariants {
		existingVariantMap[existingVariants[i].ID] = &existingVariants[i]
		if existingVariants[i].IsDefault {
			defaultVariant = existingVariants[i]
		}
	}

	// 7. Process variants.
	productCount, err := uc.productRepo.CountByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to count vendor products: %w", err)
	}

	// Always update default variant from top-level fields.
	defaultVariant.Price = req.Price
	defaultVariant.StockOnHand = req.Stock
	defaultVariant.WeightGram = req.WeightGram
	defaultVariant.IsActive = true

	variantsToUpsert := []domain.ProductVariant{defaultVariant}
	requestedVariantIDs := make(map[uuid.UUID]bool)

	for i, v := range req.Variants {
		if v.ID == nil {
			// New variant.
			sku := generateSKU(vendorDisplayNameOrEmpty(vendor), req.Name, int(productCount), len(existingVariants)+i)
			variantsToUpsert = append(variantsToUpsert, domain.ProductVariant{
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
			})
		} else {
			// Existing variant.
			variantID, err := uuid.Parse(*v.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid variant id: %w", err)
			}
			existing, found := existingVariantMap[variantID]
			if !found {
				return nil, ErrVariantNotFound
			}
			requestedVariantIDs[variantID] = true
			existing.VariantName = v.VariantName
			existing.Price = v.Price
			existing.StockOnHand = v.Stock
			existing.WeightGram = v.WeightGram
			existing.IsActive = v.IsActive
			variantsToUpsert = append(variantsToUpsert, *existing)
		}
	}

	// Deactivate non-default variants not present in the request.
	var variantIDsToDeactivate []uuid.UUID
	for id, v := range existingVariantMap {
		if !v.IsDefault && !requestedVariantIDs[id] {
			variantIDsToDeactivate = append(variantIDsToDeactivate, id)
		}
	}

	// 8. Process images.
	existingImages, err := uc.productRepo.FindImagesByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to load images: %w", err)
	}
	existingImageMap := make(map[uuid.UUID]*domain.ProductImage, len(existingImages))
	for i := range existingImages {
		existingImageMap[existingImages[i].ID] = &existingImages[i]
	}

	keepSet := make(map[uuid.UUID]bool, len(req.KeepImageIDs))
	for _, rawID := range req.KeepImageIDs {
		id, err := uuid.Parse(rawID)
		if err != nil {
			return nil, fmt.Errorf("invalid keep_image_id: %w", err)
		}
		if _, exists := existingImageMap[id]; !exists {
			return nil, ErrImageNotFound
		}
		keepSet[id] = true
	}

	if len(keepSet)+len(req.NewImages) > 10 {
		return nil, ErrTooManyImages
	}

	// Validate primary image rules on new images.
	if len(req.NewImages) > 0 {
		newPrimaryCount := 0
		for _, img := range req.NewImages {
			if img.IsPrimary {
				newPrimaryCount++
			}
		}
		if newPrimaryCount > 1 {
			return nil, ErrDuplicatePrimaryImage
		}
		keptPrimaryCount := 0
		for id := range keepSet {
			if existingImageMap[id].IsPrimary {
				keptPrimaryCount++
			}
		}
		if keptPrimaryCount+newPrimaryCount == 0 {
			return nil, ErrNoPrimaryImage
		}
	}

	// Build list of image IDs to delete.
	var imageIDsToDelete []uuid.UUID
	for id := range existingImageMap {
		if !keepSet[id] {
			imageIDsToDelete = append(imageIDsToDelete, id)
		}
	}

	// Generate presigned URLs for new images.
	now := time.Now()
	newImageEntities := make([]domain.ProductImage, 0, len(req.NewImages))
	newUploadInfos := make([]domain.ProductImageUploadInfo, 0, len(req.NewImages))

	for i, img := range req.NewImages {
		imageID := uuid.New()
		objectKey := fmt.Sprintf("products/%s/images/%s/%s",
			productID.String(), imageID.String(), sanitizeFileName(img.FileName))

		uploadURL, err := uc.storage.GeneratePresignedUploadURL(
			ctx, objectKey, img.ContentType, PresignedUploadExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for image %d: %w", i, err)
		}

		sortOrder := len(keepSet) + i
		newImageEntities = append(newImageEntities, domain.ProductImage{
			ID:        imageID,
			ProductID: productID,
			ImageURL:  objectKey,
			IsPrimary: img.IsPrimary,
			SortOrder: sortOrder,
			CreatedAt: now,
		})
		newUploadInfos = append(newUploadInfos, domain.ProductImageUploadInfo{
			ImageID:   imageID,
			UploadURL: uploadURL,
			ObjectKey: objectKey,
			SortOrder: sortOrder,
			IsPrimary: img.IsPrimary,
		})
	}

	// 9. Mutate product entity.
	product.Name = req.Name
	product.Slug = slug
	product.Description = req.Description
	product.CategoryID = categoryID
	product.Status = status
	product.ProductType = productType(req.ProductType)
	product.UpdatedAt = now

	// 10. Persist in one transaction.
	if err := uc.productRepo.UpdateProduct(
		ctx, product, variantsToUpsert, variantIDsToDeactivate, imageIDsToDelete, newImageEntities,
	); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// 11. Build response.
	variantResponses := make([]domain.ProductVariantResponse, len(variantsToUpsert))
	for i, v := range variantsToUpsert {
		variantResponses[i] = domain.ProductVariantResponse{
			ID:            v.ID,
			SKU:           v.SKU,
			VariantName:   v.VariantName,
			Price:         v.Price,
			OriginalPrice: v.Price,
			PromoPrice:    nil,
			HasPromo:      false,
			Currency:      v.Currency,
			StockOnHand:   v.StockOnHand,
			WeightGram:    v.WeightGram,
			IsDefault:     v.IsDefault,
			IsActive:      v.IsActive,
		}
	}

	return &domain.UpdateProductResponse{
		ID:            productID,
		VendorID:      vendorID,
		CategoryID:    categoryID,
		Name:          req.Name,
		Slug:          slug,
		Description:   req.Description,
		Status:        status,
		HalalAIStatus: product.HalalAIStatus,
		Variants:      variantResponses,
		NewUploadURLs: newUploadInfos,
		UpdatedAt:     now,
	}, nil
}

func (uc *productUseCase) ListProducts(ctx context.Context, vendorID uuid.UUID, params domain.VendorProductListParams) ([]domain.VendorProductListItem, *domain.PaginationMeta, error) {
	// Normalize params.
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	params.Search = strings.TrimSpace(params.Search)

	switch params.Status {
	case domain.ProductStatusDraft, domain.ProductStatusPublished:
		// valid
	default:
		params.Status = "" // all
	}

	switch params.SortBy {
	case "name", "price", "stock":
		// valid
	default:
		params.SortBy = "created_at"
	}

	switch strings.ToLower(params.SortOrder) {
	case "asc":
		params.SortOrder = "ASC"
	default:
		params.SortOrder = "DESC"
	}

	items, total, err := uc.productRepo.ListByVendor(ctx, vendorID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list vendor products: %w", err)
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}

func (uc *productUseCase) GetProductByID(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID) (*domain.VendorProductDetailResponse, error) {
	// 1. Load and authorize product.
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

	// 2. Load variants.
	variants, err := uc.productRepo.FindVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to load variants: %w", err)
	}

	// 3. Load images.
	images, err := uc.productRepo.FindImagesByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to load images: %w", err)
	}

	// 4. Extract default variant fields for top-level price/stock/weight.
	var price float64
	var stock int
	var weightGram *int
	for _, v := range variants {
		if v.IsDefault {
			price = v.Price
			stock = v.StockOnHand
			weightGram = v.WeightGram
			break
		}
	}

	// 5. Build variant responses.
	variantResponses := make([]domain.ProductVariantResponse, len(variants))
	for i, v := range variants {
		variantResponses[i] = domain.ProductVariantResponse{
			ID:            v.ID,
			SKU:           v.SKU,
			VariantName:   v.VariantName,
			Price:         v.Price,
			OriginalPrice: v.Price,
			Currency:      v.Currency,
			StockOnHand:   v.StockOnHand,
			WeightGram:    v.WeightGram,
			IsDefault:     v.IsDefault,
			IsActive:      v.IsActive,
		}
	}

	// 6. Build image responses.
	imageItems := make([]domain.VendorProductDetailImageItem, len(images))
	for i, img := range images {
		imageURL := img.ImageURL
		if imageURL != "" && !isAbsoluteURL(imageURL) {
			presignedURL, err := uc.storage.GeneratePresignedURL(ctx, imageURL, PresignedDownloadExpiry)
			if err != nil {
				return nil, fmt.Errorf("failed to generate presigned URL for image %s: %w", img.ID.String(), err)
			}
			imageURL = presignedURL
		}

		imageItems[i] = domain.VendorProductDetailImageItem{
			ID:        img.ID,
			URL:       imageURL,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
		}
	}

	return &domain.VendorProductDetailResponse{
		ID:            product.ID,
		VendorID:      product.VendorID,
		CategoryID:    product.CategoryID,
		Name:          product.Name,
		Slug:          product.Slug,
		Description:   product.Description,
		Status:        product.Status,
		HalalAIStatus: product.HalalAIStatus,
		IsActive:      product.Status == domain.ProductStatusPublished,
		Price:         price,
		Stock:         stock,
		WeightGram:    weightGram,
		Variants:      variantResponses,
		Images:        imageItems,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}, nil
}

func (uc *productUseCase) DeleteProduct(ctx context.Context, vendorID uuid.UUID, productID uuid.UUID) error {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to find product: %w", err)
	}
	if product == nil {
		return ErrProductNotFound
	}
	if product.VendorID != vendorID {
		return ErrProductNotOwned
	}
	if err := uc.productRepo.DeleteProduct(ctx, productID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func productType(pt string) string {
	switch pt {
	case domain.ProductTypeSingle, domain.ProductTypeVariant, domain.ProductTypePackage:
		return pt
	default:
		return domain.ProductTypeSingle
	}
}
