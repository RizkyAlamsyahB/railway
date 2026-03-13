package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// --- VendorOrderRepository ---

type vendorOrderRepository struct {
	db *gorm.DB
}

// NewVendorOrderRepository creates a new VendorOrderRepository backed by GORM.
func NewVendorOrderRepository(db *gorm.DB) domain.VendorOrderRepository {
	return &vendorOrderRepository{db: db}
}

func (r *vendorOrderRepository) ListByVendor(ctx context.Context, vendorID uuid.UUID, params domain.VendorOrderListParams) ([]domain.Order, int64, error) {
	q := r.db.WithContext(ctx).Model(&orderModel{}).Where("vendor_id = ?", vendorID.String())
	if params.Status != "" {
		q = q.Where("order_status = ?", params.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	var models []orderModel
	if err := q.
		Order("placed_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]domain.Order, len(models))
	for i, m := range models {
		orders[i] = *toDomainOrder(&m)
	}

	return orders, total, nil
}

func (r *vendorOrderRepository) FindByIDAndVendor(ctx context.Context, orderID, vendorID uuid.UUID) (*domain.Order, error) {
	var m orderModel
	if err := r.db.WithContext(ctx).
		Where("id = ? AND vendor_id = ?", orderID.String(), vendorID.String()).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainOrder(&m), nil
}
