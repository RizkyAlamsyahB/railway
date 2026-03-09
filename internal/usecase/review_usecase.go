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

type reviewUseCase struct {
	reviewRepo  domain.ReviewRepository
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
	storage     domain.StorageProvider
}

// NewReviewUseCase creates a new ReviewUseCase.
func NewReviewUseCase(
	reviewRepo domain.ReviewRepository,
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
	storage domain.StorageProvider,
) domain.ReviewUseCase {
	return &reviewUseCase{
		reviewRepo:  reviewRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		storage:     storage,
	}
}

func (uc *reviewUseCase) PresignImage(ctx context.Context, req domain.PresignReviewImageRequest) (*domain.PresignReviewImageResponse, error) {
	ct := normalizeContentType(req.ContentType)
	if !isAllowedImageContentType(ct) {
		return nil, ErrInvalidReviewImageContentType
	}

	objectKey := fmt.Sprintf("reviews/%s/%s", time.Now().Format("2006/01/02"), uuid.New().String())
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, ct, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return &domain.PresignReviewImageResponse{
		UploadURL:   uploadURL,
		ObjectKey:   objectKey,
		ContentType: ct,
		ExpiresIn:   int(PresignedUploadExpiry.Seconds()),
	}, nil
}

func (uc *reviewUseCase) Create(ctx context.Context, userID uuid.UUID, req domain.CreateProductReviewRequest) (*domain.CreateProductReviewResponse, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, ErrInvalidReviewRating
	}
	trimmedReviewText := strings.TrimSpace(req.ReviewText)
	if trimmedReviewText == "" {
		return nil, ErrReviewTextRequired
	}
	if len(req.Images) > 5 {
		return nil, ErrTooManyReviewImages
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}
	orderItemID, err := uuid.Parse(req.OrderItemID)
	if err != nil {
		return nil, fmt.Errorf("invalid order_item_id: %w", err)
	}

	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}
	if order.OrderStatus != domain.OrderStatusCompleted {
		return nil, ErrReviewNotAllowed
	}

	exists, err := uc.reviewRepo.ExistsByUserAndOrderItem(ctx, userID, orderItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if exists {
		return nil, ErrReviewAlreadyExists
	}

	items, err := uc.orderRepo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order items: %w", err)
	}

	var targetItem *domain.OrderItem
	for i := range items {
		if items[i].ID == orderItemID {
			targetItem = &items[i]
			break
		}
	}
	if targetItem == nil {
		return nil, ErrOrderItemNotFound
	}

	variant, err := uc.productRepo.FindVariantByID(ctx, targetItem.ProductVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product variant: %w", err)
	}
	if variant == nil {
		return nil, ErrVariantNotFound
	}

	now := time.Now()
	reviewID := uuid.New()
	review := &domain.ProductReview{
		ID:          reviewID,
		ProductID:   variant.ProductID,
		OrderID:     orderID,
		OrderItemID: orderItemID,
		UserID:      userID,
		Rating:      req.Rating,
		ReviewText:  trimmedReviewText,
		Status:      domain.ReviewStatusPublished,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	images := make([]domain.ProductReviewImage, 0, len(req.Images))
	imageURLs := make([]string, 0, len(req.Images))
	for i, img := range req.Images {
		info, err := uc.storage.HeadObject(ctx, img.ObjectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to verify image object: %w", err)
		}
		if info == nil {
			return nil, ErrReviewImageNotUploaded
		}

		contentType := normalizeContentType(info.ContentType)
		if !isAllowedImageContentType(contentType) {
			return nil, ErrInvalidReviewImageContentType
		}
		if info.ContentLength > ReviewImageMaxBytes {
			return nil, ErrReviewImageTooLarge
		}
		if info.ContentLength > int64(math.MaxInt32) {
			return nil, ErrReviewImageTooLarge
		}

		fileSize := int(info.ContentLength)
		images = append(images, domain.ProductReviewImage{
			ID:            uuid.New(),
			ReviewID:      reviewID,
			ObjectKey:     img.ObjectKey,
			MimeType:      &contentType,
			FileSizeBytes: &fileSize,
			SortOrder:     i,
			CreatedAt:     now,
		})
		imageURLs = append(imageURLs, uc.storage.GetURL(img.ObjectKey))
	}

	if err := uc.reviewRepo.CreateWithImages(ctx, review, images); err != nil {
		if isUniqueConstraintViolation(err) {
			return nil, ErrReviewAlreadyExists
		}
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return &domain.CreateProductReviewResponse{
		ReviewID:    review.ID,
		ProductID:   review.ProductID,
		OrderID:     review.OrderID,
		OrderItemID: review.OrderItemID,
		Rating:      review.Rating,
		ReviewText:  review.ReviewText,
		ImageURLs:   imageURLs,
		CreatedAt:   review.CreatedAt,
	}, nil
}

func (uc *reviewUseCase) ListByProduct(ctx context.Context, productID uuid.UUID, params domain.ProductReviewListParams) ([]domain.ProductReviewListItem, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	items, total, err := uc.reviewRepo.ListByProductID(ctx, productID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list product reviews: %w", err)
	}

	for i := range items {
		items[i].ImageURLs = make([]string, 0, len(items[i].ImageKeys))
		for _, key := range items[i].ImageKeys {
			items[i].ImageURLs = append(items[i].ImageURLs, uc.storage.GetURL(key))
		}
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}

func (uc *reviewUseCase) GetSummary(ctx context.Context, productID uuid.UUID) (*domain.ProductReviewSummaryResponse, error) {
	stats, err := uc.reviewRepo.GetStatsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product review summary: %w", err)
	}

	distribution := make([]domain.RatingDistributionItem, 0, 6)
	counts := []int64{
		stats.Star0Count,
		stats.Star1Count,
		stats.Star2Count,
		stats.Star3Count,
		stats.Star4Count,
		stats.Star5Count,
	}

	for star := 0; star <= 5; star++ {
		count := counts[star]
		pct := 0.0
		if stats.TotalReviews > 0 {
			pct = (float64(count) / float64(stats.TotalReviews)) * 100
		}
		distribution = append(distribution, domain.RatingDistributionItem{
			Star:       star,
			Count:      count,
			Percentage: math.Round(pct*100) / 100,
		})
	}

	return &domain.ProductReviewSummaryResponse{
		ProductID:     productID,
		AverageRating: math.Round(stats.AverageRating*100) / 100,
		TotalReviews:  stats.TotalReviews,
		Distribution:  distribution,
	}, nil
}
