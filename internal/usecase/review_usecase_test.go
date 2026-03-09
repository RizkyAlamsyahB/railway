package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

func setupReviewUseCase(t *testing.T) (
	*mocks.MockReviewRepository,
	*mocks.MockOrderRepository,
	*mocks.MockProductRepository,
	*mocks.MockStorageProvider,
	domain.ReviewUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	reviewRepo := mocks.NewMockReviewRepository(ctrl)
	orderRepo := mocks.NewMockOrderRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	storage := mocks.NewMockStorageProvider(ctrl)
	uc := NewReviewUseCase(reviewRepo, orderRepo, productRepo, storage)
	return reviewRepo, orderRepo, productRepo, storage, uc
}

func TestReviewUseCase_PresignImage(t *testing.T) {
	_, _, _, storage, uc := setupReviewUseCase(t)
	ctx := context.Background()

	storage.EXPECT().
		GeneratePresignedUploadURL(ctx, gomock.Any(), "image/jpeg", PresignedUploadExpiry).
		Return("https://upload.example.com", nil)

	res, err := uc.PresignImage(ctx, domain.PresignReviewImageRequest{ContentType: "image/jpeg"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.UploadURL == "" || res.ObjectKey == "" {
		t.Fatalf("expected upload URL and object key")
	}

	_, err = uc.PresignImage(ctx, domain.PresignReviewImageRequest{ContentType: "video/mp4"})
	if !errors.Is(err, ErrInvalidReviewImageContentType) {
		t.Fatalf("expected ErrInvalidReviewImageContentType, got %v", err)
	}
}

func TestReviewUseCase_Create_Success(t *testing.T) {
	reviewRepo, orderRepo, productRepo, storage, uc := setupReviewUseCase(t)
	ctx := context.Background()

	userID := uuid.New()
	orderID := uuid.New()
	orderItemID := uuid.New()
	variantID := uuid.New()
	productID := uuid.New()

	orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:            orderID,
		UserID:        userID,
		OrderStatus:   domain.OrderStatusCompleted,
		PaymentStatus: domain.PaymentStatusPaid,
	}, nil)
	reviewRepo.EXPECT().ExistsByUserAndOrderItem(ctx, userID, orderItemID).Return(false, nil)
	orderRepo.EXPECT().FindItemsByOrderID(ctx, orderID).Return([]domain.OrderItem{{
		ID:               orderItemID,
		ProductVariantID: variantID,
	}}, nil)
	productRepo.EXPECT().FindVariantByID(ctx, variantID).Return(&domain.ProductVariant{
		ID:        variantID,
		ProductID: productID,
	}, nil)
	storage.EXPECT().HeadObject(ctx, "reviews/obj-1").Return(&domain.ObjectInfo{
		ContentType:   "image/jpeg",
		ContentLength: 1024,
	}, nil)
	reviewRepo.EXPECT().CreateWithImages(ctx, gomock.Any(), gomock.Any()).Return(nil)
	storage.EXPECT().GetURL("reviews/obj-1").Return("https://cdn.example.com/reviews/obj-1")

	res, err := uc.Create(ctx, userID, domain.CreateProductReviewRequest{
		OrderID:     orderID.String(),
		OrderItemID: orderItemID.String(),
		Rating:      5,
		ReviewText:  "Bagus sekali",
		Images: []domain.ReviewImageInput{{
			ObjectKey: "reviews/obj-1",
		}},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ProductID != productID {
		t.Fatalf("expected product ID %s, got %s", productID, res.ProductID)
	}
	if len(res.ImageURLs) != 1 {
		t.Fatalf("expected 1 image URL, got %d", len(res.ImageURLs))
	}
}

func TestReviewUseCase_Create_NotEligible(t *testing.T) {
	reviewRepo, orderRepo, productRepo, storage, uc := setupReviewUseCase(t)
	_ = reviewRepo
	_ = productRepo
	_ = storage

	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()
	orderItemID := uuid.New()

	orderRepo.EXPECT().FindByID(ctx, orderID).Return(&domain.Order{
		ID:          orderID,
		UserID:      userID,
		OrderStatus: domain.OrderStatusPaid,
	}, nil)

	_, err := uc.Create(ctx, userID, domain.CreateProductReviewRequest{
		OrderID:     orderID.String(),
		OrderItemID: orderItemID.String(),
		Rating:      5,
		ReviewText:  "Mantap",
	})
	if !errors.Is(err, ErrReviewNotAllowed) {
		t.Fatalf("expected ErrReviewNotAllowed, got %v", err)
	}
}
