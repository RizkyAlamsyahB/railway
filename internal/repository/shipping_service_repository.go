package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type shippingServiceModel struct {
	ID       string `gorm:"column:id;primaryKey"`
	Code     string `gorm:"column:code"`
	Name     string `gorm:"column:name"`
	IsActive bool   `gorm:"column:is_active"`
}

func (shippingServiceModel) TableName() string { return "shipping_services" }

type shippingServiceRepository struct {
	db *gorm.DB
}

// NewShippingServiceRepository creates a new ShippingServiceRepository backed by GORM.
func NewShippingServiceRepository(db *gorm.DB) domain.ShippingServiceRepository {
	return &shippingServiceRepository{db: db}
}

func (r *shippingServiceRepository) ListActive(ctx context.Context) ([]domain.ShippingService, error) {
	var models []shippingServiceModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	services := make([]domain.ShippingService, len(models))
	for i, m := range models {
		services[i] = *toDomainShippingService(&m)
	}
	return services, nil
}

func (r *shippingServiceRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.ShippingService, error) {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = id.String()
	}

	var models []shippingServiceModel
	if err := r.db.WithContext(ctx).Where("id IN ? AND is_active = ?", strIDs, true).Find(&models).Error; err != nil {
		return nil, err
	}

	services := make([]domain.ShippingService, len(models))
	for i, m := range models {
		services[i] = *toDomainShippingService(&m)
	}
	return services, nil
}

func toDomainShippingService(m *shippingServiceModel) *domain.ShippingService {
	id, _ := uuid.Parse(m.ID)
	return &domain.ShippingService{
		ID:       id,
		Code:     m.Code,
		Name:     m.Name,
		IsActive: m.IsActive,
	}
}
