package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type vendorBannerUseCase struct {
	bannerRepo domain.VendorBannerRepository
	storage    domain.StorageProvider
}

// NewVendorBannerUseCase creates a new VendorBannerUseCase.
func NewVendorBannerUseCase(bannerRepo domain.VendorBannerRepository, storage domain.StorageProvider) domain.VendorBannerUseCase {
	return &vendorBannerUseCase{
		bannerRepo: bannerRepo,
		storage:    storage,
	}
}

func (uc *vendorBannerUseCase) CreateBanner(ctx context.Context, vendorID uuid.UUID, req domain.CreateVendorBannerRequest) (*domain.VendorBannerPresignResponse, error) {
	ct := strings.ToLower(strings.TrimSpace(req.ContentType))
	if !allowedBannerContentTypes[ct] {
		return nil, ErrInvalidVendorBannerContentType
	}

	bannerID := uuid.New()
	now := time.Now()
	objectKey := fmt.Sprintf("vendor-banners/%s/%s/%d", vendorID.String(), bannerID.String(), now.UnixMilli())

	banner := &domain.VendorBanner{
		ID:        bannerID,
		VendorID:  vendorID,
		Title:     req.Title,
		ImageURL:  objectKey,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.bannerRepo.Create(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to create vendor banner: %w", err)
	}

	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, ct, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &domain.VendorBannerPresignResponse{
		BannerID:  bannerID,
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		ExpiresIn: int(PresignedUploadExpiry / time.Second),
	}, nil
}

func (uc *vendorBannerUseCase) ConfirmBanner(ctx context.Context, vendorID uuid.UUID, bannerID uuid.UUID, req domain.ConfirmVendorBannerRequest) (*domain.VendorBannerResponse, error) {
	banner, err := uc.bannerRepo.FindByID(ctx, bannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor banner: %w", err)
	}
	if banner == nil || banner.VendorID != vendorID {
		return nil, ErrVendorBannerNotFound
	}

	info, err := uc.storage.HeadObject(ctx, req.ObjectKey)
	if err != nil || info == nil {
		return nil, ErrVendorBannerImageNotUploaded
	}

	if !allowedBannerContentTypes[strings.ToLower(info.ContentType)] {
		return nil, ErrInvalidVendorBannerContentType
	}

	if info.ContentLength > MaxBannerImageSize {
		return nil, ErrVendorBannerImageTooLarge
	}

	banner.ImageURL = req.ObjectKey
	banner.UpdatedAt = time.Now()

	if err := uc.bannerRepo.Update(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to update vendor banner: %w", err)
	}

	imageURL := uc.storage.GetURL(banner.ImageURL)
	return toVendorBannerResponse(banner, imageURL), nil
}

func (uc *vendorBannerUseCase) ListBanners(ctx context.Context, vendorID uuid.UUID) ([]domain.VendorBannerResponse, error) {
	banners, err := uc.bannerRepo.ListByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to list vendor banners: %w", err)
	}
	result := make([]domain.VendorBannerResponse, len(banners))
	for i, b := range banners {
		imageURL := uc.storage.GetURL(b.ImageURL)
		result[i] = *toVendorBannerResponse(&b, imageURL)
	}
	return result, nil
}

func (uc *vendorBannerUseCase) GetBanner(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) (*domain.VendorBannerResponse, error) {
	banner, err := uc.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor banner: %w", err)
	}
	if banner == nil || banner.VendorID != vendorID {
		return nil, ErrVendorBannerNotFound
	}
	imageURL := uc.storage.GetURL(banner.ImageURL)
	return toVendorBannerResponse(banner, imageURL), nil
}

func (uc *vendorBannerUseCase) UpdateBanner(ctx context.Context, vendorID uuid.UUID, id uuid.UUID, req domain.UpdateVendorBannerRequest) (*domain.VendorBannerResponse, error) {
	banner, err := uc.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor banner: %w", err)
	}
	if banner == nil || banner.VendorID != vendorID {
		return nil, ErrVendorBannerNotFound
	}

	if req.Title != nil {
		banner.Title = *req.Title
	}
	banner.UpdatedAt = time.Now()

	if err := uc.bannerRepo.Update(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to update vendor banner: %w", err)
	}

	imageURL := uc.storage.GetURL(banner.ImageURL)
	return toVendorBannerResponse(banner, imageURL), nil
}

func (uc *vendorBannerUseCase) DeleteBanner(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) error {
	banner, err := uc.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find vendor banner: %w", err)
	}
	if banner == nil || banner.VendorID != vendorID {
		return ErrVendorBannerNotFound
	}

	if err := uc.bannerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete vendor banner: %w", err)
	}

	// Best-effort cleanup of S3 object.
	if banner.ImageURL != "" {
		_ = uc.storage.Delete(ctx, banner.ImageURL)
	}

	return nil
}

func (uc *vendorBannerUseCase) ListByVendorIDPublic(ctx context.Context, vendorID uuid.UUID) ([]domain.VendorBannerResponse, error) {
	return uc.ListBanners(ctx, vendorID)
}

func toVendorBannerResponse(b *domain.VendorBanner, imageURL string) *domain.VendorBannerResponse {
	return &domain.VendorBannerResponse{
		ID:        b.ID,
		VendorID:  b.VendorID,
		Title:     b.Title,
		ImageURL:  imageURL,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
