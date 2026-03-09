package repository

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type productReviewModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	ProductID   string    `gorm:"column:product_id"`
	OrderID     string    `gorm:"column:order_id"`
	OrderItemID string    `gorm:"column:order_item_id"`
	UserID      string    `gorm:"column:user_id"`
	Rating      int       `gorm:"column:rating"`
	ReviewText  string    `gorm:"column:review_text"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (productReviewModel) TableName() string { return "product_reviews" }

type productReviewImageModel struct {
	ID            string    `gorm:"column:id;primaryKey"`
	ReviewID      string    `gorm:"column:review_id"`
	ObjectKey     string    `gorm:"column:object_key"`
	MimeType      *string   `gorm:"column:mime_type"`
	FileSizeBytes *int      `gorm:"column:file_size_bytes"`
	SortOrder     int       `gorm:"column:sort_order"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (productReviewImageModel) TableName() string { return "product_review_images" }

type productReviewStatsModel struct {
	ProductID     string    `gorm:"column:product_id;primaryKey"`
	TotalReviews  int64     `gorm:"column:total_reviews"`
	TotalStars    int64     `gorm:"column:total_stars"`
	AverageRating float64   `gorm:"column:average_rating"`
	Star0Count    int64     `gorm:"column:star_0_count"`
	Star1Count    int64     `gorm:"column:star_1_count"`
	Star2Count    int64     `gorm:"column:star_2_count"`
	Star3Count    int64     `gorm:"column:star_3_count"`
	Star4Count    int64     `gorm:"column:star_4_count"`
	Star5Count    int64     `gorm:"column:star_5_count"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (productReviewStatsModel) TableName() string { return "product_review_stats" }

type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new ReviewRepository backed by GORM.
func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) ExistsByUserAndOrderItem(ctx context.Context, userID, orderItemID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&productReviewModel{}).
		Where("user_id = ? AND order_item_id = ?", userID.String(), orderItemID.String()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *reviewRepository) CreateWithImages(ctx context.Context, review *domain.ProductReview, images []domain.ProductReviewImage) error {
	star1, star2, star3, star4, star5 := 0, 0, 0, 0, 0
	switch review.Rating {
	case 1:
		star1 = 1
	case 2:
		star2 = 1
	case 3:
		star3 = 1
	case 4:
		star4 = 1
	case 5:
		star5 = 1
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rm := toProductReviewModel(review)
		if err := tx.Create(&rm).Error; err != nil {
			return err
		}

		if len(images) > 0 {
			models := make([]productReviewImageModel, len(images))
			for i := range images {
				models[i] = toProductReviewImageModel(&images[i])
			}
			if err := tx.Create(&models).Error; err != nil {
				return err
			}
		}

		if err := tx.Exec(`
INSERT INTO product_review_stats (
	product_id, total_reviews, total_stars, average_rating,
	star_0_count, star_1_count, star_2_count, star_3_count, star_4_count, star_5_count, updated_at
) VALUES (?, 1, ?, ?, 0, ?, ?, ?, ?, ?, now())
ON CONFLICT (product_id) DO UPDATE SET
	total_reviews = product_review_stats.total_reviews + 1,
	total_stars = product_review_stats.total_stars + EXCLUDED.total_stars,
	star_1_count = product_review_stats.star_1_count + EXCLUDED.star_1_count,
	star_2_count = product_review_stats.star_2_count + EXCLUDED.star_2_count,
	star_3_count = product_review_stats.star_3_count + EXCLUDED.star_3_count,
	star_4_count = product_review_stats.star_4_count + EXCLUDED.star_4_count,
	star_5_count = product_review_stats.star_5_count + EXCLUDED.star_5_count,
	average_rating = ROUND((product_review_stats.total_stars + EXCLUDED.total_stars)::numeric / NULLIF(product_review_stats.total_reviews + 1, 0), 2),
	updated_at = now()
`, review.ProductID.String(), review.Rating, float64(review.Rating), star1, star2, star3, star4, star5).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *reviewRepository) ListByProductID(ctx context.Context, productID uuid.UUID, params domain.ProductReviewListParams) ([]domain.ProductReviewListItem, int64, error) {
	base := r.db.WithContext(ctx).
		Table("product_reviews pr").
		Where("pr.product_id = ?", productID.String()).
		Where("pr.status = ?", domain.ReviewStatusPublished)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type reviewRow struct {
		ReviewID     string    `gorm:"column:review_id"`
		UserID       string    `gorm:"column:user_id"`
		ReviewerName string    `gorm:"column:reviewer_name"`
		Rating       int       `gorm:"column:rating"`
		ReviewText   string    `gorm:"column:review_text"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}

	offset := (params.Page - 1) * params.Limit
	var rows []reviewRow
	if err := base.
		Select("pr.id AS review_id, pr.user_id, u.full_name AS reviewer_name, pr.rating, pr.review_text, pr.created_at").
		Joins("JOIN users u ON u.id = pr.user_id").
		Order("pr.created_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return []domain.ProductReviewListItem{}, total, nil
	}

	reviewIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		reviewIDs = append(reviewIDs, row.ReviewID)
	}

	type imageRow struct {
		ReviewID  string `gorm:"column:review_id"`
		ObjectKey string `gorm:"column:object_key"`
	}
	var imgRows []imageRow
	if err := r.db.WithContext(ctx).
		Table("product_review_images").
		Select("review_id, object_key").
		Where("review_id IN ?", reviewIDs).
		Order("sort_order ASC").
		Scan(&imgRows).Error; err != nil {
		return nil, 0, err
	}

	imageMap := make(map[string][]string, len(rows))
	for _, img := range imgRows {
		imageMap[img.ReviewID] = append(imageMap[img.ReviewID], img.ObjectKey)
	}

	items := make([]domain.ProductReviewListItem, len(rows))
	for i, row := range rows {
		reviewID, _ := uuid.Parse(row.ReviewID)
		userID, _ := uuid.Parse(row.UserID)
		items[i] = domain.ProductReviewListItem{
			ReviewID:     reviewID,
			UserID:       userID,
			ReviewerName: row.ReviewerName,
			Rating:       row.Rating,
			ReviewText:   row.ReviewText,
			ImageKeys:    imageMap[row.ReviewID],
			CreatedAt:    row.CreatedAt,
		}
	}

	return items, total, nil
}

func (r *reviewRepository) GetStatsByProductID(ctx context.Context, productID uuid.UUID) (*domain.ProductReviewStats, error) {
	var model productReviewStatsModel
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID.String()).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.ProductReviewStats{ProductID: productID}, nil
		}
		return nil, err
	}

	return toDomainProductReviewStats(&model), nil
}

func toProductReviewModel(r *domain.ProductReview) productReviewModel {
	return productReviewModel{
		ID:          r.ID.String(),
		ProductID:   r.ProductID.String(),
		OrderID:     r.OrderID.String(),
		OrderItemID: r.OrderItemID.String(),
		UserID:      r.UserID.String(),
		Rating:      r.Rating,
		ReviewText:  r.ReviewText,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toProductReviewImageModel(img *domain.ProductReviewImage) productReviewImageModel {
	return productReviewImageModel{
		ID:            img.ID.String(),
		ReviewID:      img.ReviewID.String(),
		ObjectKey:     img.ObjectKey,
		MimeType:      img.MimeType,
		FileSizeBytes: img.FileSizeBytes,
		SortOrder:     img.SortOrder,
		CreatedAt:     img.CreatedAt,
	}
}

func toDomainProductReviewStats(m *productReviewStatsModel) *domain.ProductReviewStats {
	productID, _ := uuid.Parse(m.ProductID)
	return &domain.ProductReviewStats{
		ProductID:     productID,
		TotalReviews:  m.TotalReviews,
		TotalStars:    m.TotalStars,
		AverageRating: math.Round(m.AverageRating*100) / 100,
		Star0Count:    m.Star0Count,
		Star1Count:    m.Star1Count,
		Star2Count:    m.Star2Count,
		Star3Count:    m.Star3Count,
		Star4Count:    m.Star4Count,
		Star5Count:    m.Star5Count,
	}
}

// Ensure compiler verifies interface implementation.
var _ domain.ReviewRepository = (*reviewRepository)(nil)
