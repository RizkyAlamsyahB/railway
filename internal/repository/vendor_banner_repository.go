package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type vendorBannerModel struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	VendorID  uuid.UUID `gorm:"column:vendor_id;type:uuid"`
	Title     string    `gorm:"column:title"`
	ImageURL  string    `gorm:"column:image_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (vendorBannerModel) TableName() string { return "vendor_banners" }

func toVendorBannerDomain(m vendorBannerModel) domain.VendorBanner {
	return domain.VendorBanner{
		ID:        m.ID,
		VendorID:  m.VendorID,
		Title:     m.Title,
		ImageURL:  m.ImageURL,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toVendorBannerModel(b *domain.VendorBanner) vendorBannerModel {
	return vendorBannerModel{
		ID:        b.ID,
		VendorID:  b.VendorID,
		Title:     b.Title,
		ImageURL:  b.ImageURL,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

type vendorBannerRepository struct {
	db *gorm.DB
}

// NewVendorBannerRepository creates a VendorBannerRepository backed by GORM.
func NewVendorBannerRepository(db *gorm.DB) domain.VendorBannerRepository {
	return &vendorBannerRepository{db: db}
}

func (r *vendorBannerRepository) Create(ctx context.Context, banner *domain.VendorBanner) error {
	m := toVendorBannerModel(banner)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *vendorBannerRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.VendorBanner, error) {
	var m vendorBannerModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	b := toVendorBannerDomain(m)
	return &b, nil
}

func (r *vendorBannerRepository) ListByVendorID(ctx context.Context, vendorID uuid.UUID) ([]domain.VendorBanner, error) {
	var rows []vendorBannerModel
	if err := r.db.WithContext(ctx).Where("vendor_id = ?", vendorID).Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.VendorBanner, len(rows))
	for i, row := range rows {
		items[i] = toVendorBannerDomain(row)
	}
	return items, nil
}

func (r *vendorBannerRepository) Update(ctx context.Context, banner *domain.VendorBanner) error {
	m := toVendorBannerModel(banner)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *vendorBannerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&vendorBannerModel{}).Error
}
