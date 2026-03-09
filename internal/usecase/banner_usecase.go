package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// Allowed banner image content types.
var allowedBannerContentTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
}

// MaxBannerImageSize is the maximum allowed banner image size (5 MB).
const MaxBannerImageSize int64 = 5 * 1024 * 1024

type bannerUseCase struct {
	bannerRepo domain.BannerRepository
	storage    domain.StorageProvider
}

// NewBannerUseCase creates a new BannerUseCase with the given dependencies.
func NewBannerUseCase(bannerRepo domain.BannerRepository, storage domain.StorageProvider) domain.BannerUseCase {
	return &bannerUseCase{
		bannerRepo: bannerRepo,
		storage:    storage,
	}
}

func (uc *bannerUseCase) CreateBanner(ctx context.Context, req domain.CreateBannerRequest) (*domain.BannerPresignResponse, error) {
	// Validate content type.
	ct := strings.ToLower(strings.TrimSpace(req.ContentType))
	if !allowedBannerContentTypes[ct] {
		return nil, ErrInvalidBannerContentType
	}

	// Create banner entity with a placeholder image URL (will be confirmed after upload).
	bannerID := uuid.New()
	now := time.Now()

	objectKey := fmt.Sprintf("banners/%s/%d", bannerID.String(), now.UnixMilli())

	banner := &domain.Banner{
		ID:        bannerID,
		Title:     req.Title,
		ImageURL:  objectKey,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.bannerRepo.Create(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to create banner: %w", err)
	}

	// Generate presigned upload URL.
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, ct, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return &domain.BannerPresignResponse{
		BannerID:  bannerID,
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		ExpiresIn: int(PresignedUploadExpiry / time.Second),
	}, nil
}

func (uc *bannerUseCase) ConfirmBanner(ctx context.Context, bannerID uuid.UUID, req domain.ConfirmBannerRequest) (*domain.BannerResponse, error) {
	banner, err := uc.bannerRepo.FindByID(ctx, bannerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find banner: %w", err)
	}
	if banner == nil {
		return nil, ErrBannerNotFound
	}

	// Verify the object exists in storage.
	info, err := uc.storage.HeadObject(ctx, req.ObjectKey)
	if err != nil || info == nil {
		return nil, ErrBannerImageNotUploaded
	}

	// Validate content type.
	if !allowedBannerContentTypes[strings.ToLower(info.ContentType)] {
		return nil, ErrInvalidBannerContentType
	}

	// Validate file size.
	if info.ContentLength > MaxBannerImageSize {
		return nil, ErrBannerImageTooLarge
	}

	// Update banner image URL.
	banner.ImageURL = req.ObjectKey
	banner.UpdatedAt = time.Now()

	if err := uc.bannerRepo.Update(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	imageURL := uc.storage.GetURL(banner.ImageURL)
	return toBannerResponse(banner, imageURL), nil
}

func (uc *bannerUseCase) ListBanners(ctx context.Context) ([]domain.BannerResponse, error) {
	banners, err := uc.bannerRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}
	result := make([]domain.BannerResponse, len(banners))
	for i, b := range banners {
		imageURL := uc.storage.GetURL(b.ImageURL)
		result[i] = *toBannerResponse(&b, imageURL)
	}
	return result, nil
}

func (uc *bannerUseCase) GetBanner(ctx context.Context, id uuid.UUID) (*domain.BannerResponse, error) {
	banner, err := uc.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find banner: %w", err)
	}
	if banner == nil {
		return nil, ErrBannerNotFound
	}
	imageURL := uc.storage.GetURL(banner.ImageURL)
	return toBannerResponse(banner, imageURL), nil
}

func (uc *bannerUseCase) UpdateBanner(ctx context.Context, id uuid.UUID, req domain.UpdateBannerRequest) (*domain.BannerResponse, error) {
	banner, err := uc.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find banner: %w", err)
	}
	if banner == nil {
		return nil, ErrBannerNotFound
	}

	if req.Title != nil {
		banner.Title = *req.Title
	}
	banner.UpdatedAt = time.Now()

	if err := uc.bannerRepo.Update(ctx, banner); err != nil {
		return nil, fmt.Errorf("failed to update banner: %w", err)
	}

	imageURL := uc.storage.GetURL(banner.ImageURL)
	return toBannerResponse(banner, imageURL), nil
}

func (uc *bannerUseCase) DeleteBanner(ctx context.Context, id uuid.UUID) error {
	banner, err := uc.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find banner: %w", err)
	}
	if banner == nil {
		return ErrBannerNotFound
	}

	// Delete image from storage.
	if banner.ImageURL != "" {
		_ = uc.storage.Delete(ctx, banner.ImageURL) // best-effort cleanup
	}

	if err := uc.bannerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}
	return nil
}

// toBannerResponse maps a domain Banner to a BannerResponse.
func toBannerResponse(b *domain.Banner, imageURL string) *domain.BannerResponse {
	return &domain.BannerResponse{
		ID:        b.ID,
		Title:     b.Title,
		ImageURL:  imageURL,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
