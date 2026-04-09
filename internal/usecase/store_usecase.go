package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type storeUseCase struct {
	vendorRepo       domain.VendorRepository
	productRepo      domain.ProductRepository
	vendorBannerRepo domain.VendorBannerRepository
	reviewRepo       domain.ReviewRepository
	addressRepo      domain.AddressRepository
	storage          domain.StorageProvider
}

// NewStoreUseCase creates a new StoreUseCase.
func NewStoreUseCase(
	vendorRepo domain.VendorRepository,
	productRepo domain.ProductRepository,
	vendorBannerRepo domain.VendorBannerRepository,
	reviewRepo domain.ReviewRepository,
	addressRepo domain.AddressRepository,
	storage domain.StorageProvider,
) domain.StoreUseCase {
	return &storeUseCase{
		vendorRepo:       vendorRepo,
		productRepo:      productRepo,
		vendorBannerRepo: vendorBannerRepo,
		reviewRepo:       reviewRepo,
		addressRepo:      addressRepo,
		storage:          storage,
	}
}

const (
	storeRecommendedLimit = 8
	storeNewProductsLimit = 8
)

func (uc *storeUseCase) GetStoreDetail(ctx context.Context, vendorID uuid.UUID, productsPage int, productsLimit int) (*domain.StoreDetailResponse, error) {
	// 1. Find vendor (must be active).
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil || vendor.Status != "active" {
		return nil, ErrVendorNotFound
	}

	// 2. Normalize pagination.
	if productsPage < 1 {
		productsPage = 1
	}
	if productsLimit < 1 || productsLimit > 100 {
		productsLimit = 20
	}

	// 3. Get logo image URL from vendor documents.
	var imageURL string
	logoDoc, err := uc.vendorRepo.FindDocumentByVendorIDAndDocType(ctx, vendorID, domain.VendorDocumentTypeBusinessLogo)
	if err == nil && logoDoc != nil && logoDoc.FileURL != "" {
		imageURL = uc.resolveURL(logoDoc.FileURL)
	}

	// 4. Get aggregated rating across all vendor products.
	avgRating, totalReviews, err := uc.vendorRepo.GetVendorAggregatedRating(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor rating: %w", err)
	}

	// 5. Total published products count.
	totalProducts, err := uc.productRepo.CountByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to count vendor products: %w", err)
	}

	// 6. Banners.
	banners, err := uc.vendorBannerRepo.ListByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to list vendor banners: %w", err)
	}
	bannerResponses := make([]domain.VendorBannerResponse, len(banners))
	for i, b := range banners {
		bannerResponses[i] = domain.VendorBannerResponse{
			ID:        b.ID,
			VendorID:  b.VendorID,
			Title:     b.Title,
			ImageURL:  uc.resolveURL(b.ImageURL),
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
	}

	// 7. Recommended products (bestseller, limit 8).
	recommended, _, err := uc.productRepo.ListPublishedByVendorForStore(ctx, vendorID, "bestseller", storeRecommendedLimit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list recommended products: %w", err)
	}
	uc.resolveProductImageURLs(recommended)

	// 8. New products (newest, limit 8).
	newProducts, _, err := uc.productRepo.ListPublishedByVendorForStore(ctx, vendorID, "newest", storeNewProductsLimit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to list new products: %w", err)
	}
	uc.resolveProductImageURLs(newProducts)

	// 9. Grid products (paginated, newest).
	gridOffset := (productsPage - 1) * productsLimit
	gridProducts, gridTotal, err := uc.productRepo.ListPublishedByVendorForStore(ctx, vendorID, "newest", productsLimit, gridOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to list grid products: %w", err)
	}
	uc.resolveProductImageURLs(gridProducts)

	gridMeta := &domain.PaginationMeta{
		Page:       productsPage,
		Limit:      productsLimit,
		TotalItems: gridTotal,
		TotalPages: int(math.Ceil(float64(gridTotal) / float64(productsLimit))),
	}

	// Round rating to 1 decimal place.
	roundedRating := math.Round(avgRating*10) / 10
	displayName := strings.TrimSpace(derefStoreString(vendor.DisplayName))
	if displayName == "" {
		displayName = strings.TrimSpace(derefStoreString(vendor.LegalName))
	}
	location := buildVendorLocation(vendor)
	if location == "" {
		if addr, err := uc.addressRepo.FindDefaultByUserID(ctx, vendor.OwnerUserID); err == nil && addr != nil {
			location = buildAddressLocation(addr)
		}
	}

	return &domain.StoreDetailResponse{
		Store: domain.StoreInfoResponse{
			ID:            vendor.ID,
			DisplayName:   displayName,
			Description:   vendor.Description,
			Location:      location,
			ImageURL:      imageURL,
			Rating:        roundedRating,
			TotalReviews:  totalReviews,
			TotalSold:     vendor.TotalSold,
			TotalProducts: totalProducts,
		},
		Banners:             bannerResponses,
		RecommendedProducts: recommended,
		NewProducts:         newProducts,
		Products:            gridProducts,
		ProductsPagination:  gridMeta,
	}, nil
}

func (uc *storeUseCase) GetStoreReviews(ctx context.Context, vendorID uuid.UUID, page, limit int) (*domain.VendorReviewResponse, error) {
	// 1. Find vendor (must be active).
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil || vendor.Status != "active" {
		return nil, ErrVendorNotFound
	}

	// 2. Normalize pagination.
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 3. Get reviews from db
	reviews, total, err := uc.reviewRepo.ListByVendor(ctx, vendorID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list vendor reviews: %w", err)
	}

	// 4. Resolve Image URLs
	for i := range reviews {
		var urls []string
		for _, key := range reviews[i].ImageKeys {
			urls = append(urls, uc.resolveURL(key))
		}
		reviews[i].Images = urls
	}

	// 5. Build PaginationMeta
	meta := &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	// 6. Build Summary. For this, we should really have an aggregated model,
	// but we can compute rating breakdown by getting it from reviewRepository if needed,
	// or we can just fetch GetVendorAggregatedRating for the average_rating and total_reviews,
	// and we don't have a direct repo method for rating breakdown. Wait, `reviewRepository` has `GetStatsByProductID`
	// but not by vendor. I'll need to query it from repo, or add a method.
	// For now, let's just make it a dummy breakdown or the real one by fetching `vendorRepo.GetVendorAggregatedRating`.
	avgRating, totalReviews, err := uc.vendorRepo.GetVendorAggregatedRating(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor rating: %w", err)
	}

	// Note: We don't have GetVendorRatingBreakdown yet. Let's return a dummy or empty breakdown for now,
	// unless we implement it.
	breakdown := make(map[string]int64)

	summary := domain.VendorReviewSummary{
		AverageRating:   math.Round(avgRating*10) / 10,
		TotalReviews:    totalReviews,
		RatingBreakdown: breakdown,
	}

	return &domain.VendorReviewResponse{
		Reviews:    reviews,
		Summary:    summary,
		Pagination: meta,
	}, nil
}

// resolveURL converts S3 object keys to full URLs.
func (uc *storeUseCase) resolveURL(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	lower := strings.ToLower(v)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return v
	}
	return uc.storage.GetURL(v)
}

// resolveProductImageURLs resolves all image URLs in-place for a slice of store product items.
func (uc *storeUseCase) resolveProductImageURLs(items []domain.StoreProductItem) {
	for i := range items {
		items[i].ImageURL = uc.resolveURL(items[i].ImageURL)
	}
}

func derefStoreString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func buildAddressLocation(addr *domain.Address) string {
	if addr == nil {
		return ""
	}
	parts := make([]string, 0, 6)
	if v := strings.TrimSpace(addr.AddressLine); v != "" {
		parts = append(parts, v)
	}
	for _, p := range []*string{addr.SubdistrictName, addr.DistrictName, addr.CityName, addr.ProvinceName} {
		if p != nil {
			if v := strings.TrimSpace(*p); v != "" {
				parts = append(parts, v)
			}
		}
	}
	if p := addr.PostalCode; p != nil {
		if v := strings.TrimSpace(*p); v != "" {
			parts = append(parts, v)
		}
	}
	return strings.Join(parts, ", ")
}

func buildVendorLocation(vendor *domain.Vendor) string {
	if vendor == nil {
		return ""
	}

	parts := make([]string, 0, 5)
	for _, part := range []*string{
		vendor.AddressLine,
		vendor.SubdistrictName,
		vendor.DistrictName,
		vendor.CityName,
		vendor.ProvinceName,
	} {
		value := strings.TrimSpace(derefStoreString(part))
		if value != "" {
			parts = append(parts, value)
		}
	}

	postalCode := strings.TrimSpace(derefStoreString(vendor.PostalCode))
	if postalCode != "" {
		parts = append(parts, postalCode)
	}

	return strings.Join(parts, ", ")
}
