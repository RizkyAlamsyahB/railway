package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type bannerModel struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Title     string    `gorm:"column:title"`
	ImageURL  string    `gorm:"column:image_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (bannerModel) TableName() string { return "banners" }

func toBannerDomain(m bannerModel) domain.Banner {
	return domain.Banner{
		ID:        m.ID,
		Title:     m.Title,
		ImageURL:  m.ImageURL,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toBannerModel(b *domain.Banner) bannerModel {
	return bannerModel{
		ID:        b.ID,
		Title:     b.Title,
		ImageURL:  b.ImageURL,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

type bannerRepository struct {
	db *gorm.DB
}

// NewBannerRepository creates a BannerRepository backed by GORM.
func NewBannerRepository(db *gorm.DB) domain.BannerRepository {
	return &bannerRepository{db: db}
}

func (r *bannerRepository) Create(ctx context.Context, banner *domain.Banner) error {
	m := toBannerModel(banner)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *bannerRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Banner, error) {
	var m bannerModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	b := toBannerDomain(m)
	return &b, nil
}

func (r *bannerRepository) List(ctx context.Context) ([]domain.Banner, error) {
	q := r.db.WithContext(ctx).Model(&bannerModel{})
	q = q.Order("created_at DESC")

	var rows []bannerModel
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.Banner, len(rows))
	for i, row := range rows {
		items[i] = toBannerDomain(row)
	}
	return items, nil
}

func (r *bannerRepository) Update(ctx context.Context, banner *domain.Banner) error {
	m := toBannerModel(banner)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *bannerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&bannerModel{}).Error
}
