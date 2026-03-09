package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// --- GORM models ---

type courierModel struct {
	ID       int    `gorm:"primaryKey;column:id"`
	Code     string `gorm:"column:code"`
	Name     string `gorm:"column:name"`
	LogoURL  string `gorm:"column:logo_url"`
	IsActive bool   `gorm:"column:is_active"`
}

func (courierModel) TableName() string { return "couriers" }

type vendorCourierModel struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;column:id;default:gen_random_uuid()"`
	VendorID  uuid.UUID `gorm:"type:uuid;column:vendor_id"`
	CourierID int       `gorm:"column:courier_id"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (vendorCourierModel) TableName() string { return "vendor_couriers" }

// --- CourierRepository ---

type courierRepository struct {
	db *gorm.DB
}

// NewCourierRepository creates a new CourierRepository backed by GORM.
func NewCourierRepository(db *gorm.DB) domain.CourierRepository {
	return &courierRepository{db: db}
}

func (r *courierRepository) FindAll(ctx context.Context, onlyActive bool) ([]domain.Courier, error) {
	var models []courierModel
	q := r.db.WithContext(ctx)
	if onlyActive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Courier, len(models))
	for i, m := range models {
		result[i] = domain.Courier{
			ID:       m.ID,
			Code:     m.Code,
			Name:     m.Name,
			LogoURL:  m.LogoURL,
			IsActive: m.IsActive,
		}
	}
	return result, nil
}

func (r *courierRepository) FindByID(ctx context.Context, id int) (*domain.Courier, error) {
	var m courierModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Courier{
		ID:       m.ID,
		Code:     m.Code,
		Name:     m.Name,
		LogoURL:  m.LogoURL,
		IsActive: m.IsActive,
	}, nil
}

// --- VendorCourierRepository ---

type vendorCourierRepository struct {
	db *gorm.DB
}

// NewVendorCourierRepository creates a new VendorCourierRepository backed by GORM.
func NewVendorCourierRepository(db *gorm.DB) domain.VendorCourierRepository {
	return &vendorCourierRepository{db: db}
}

func (r *vendorCourierRepository) FindByVendorID(ctx context.Context, vendorID uuid.UUID) ([]domain.VendorCourierResponse, error) {
	var results []domain.VendorCourierResponse
	err := r.db.WithContext(ctx).
		Table("vendor_couriers vc").
		Select("vc.id, vc.courier_id, c.code AS courier_code, c.name AS courier_name, c.logo_url AS logo_url, vc.is_active, vc.created_at").
		Joins("JOIN couriers c ON c.id = vc.courier_id").
		Where("vc.vendor_id = ?", vendorID).
		Order("c.name ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *vendorCourierRepository) SetCouriers(ctx context.Context, vendorID uuid.UUID, courierIDs []int) ([]domain.VendorCourierResponse, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing
		if err := tx.Where("vendor_id = ?", vendorID).Delete(&vendorCourierModel{}).Error; err != nil {
			return err
		}
		// Insert new
		for _, cid := range courierIDs {
			m := vendorCourierModel{
				ID:        uuid.New(),
				VendorID:  vendorID,
				CourierID: cid,
				IsActive:  true,
			}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Return updated list
	return r.FindByVendorID(ctx, vendorID)
}

func (r *vendorCourierRepository) DeleteByVendorAndCourier(ctx context.Context, vendorID uuid.UUID, courierID int) error {
	result := r.db.WithContext(ctx).
		Where("vendor_id = ? AND courier_id = ?", vendorID, courierID).
		Delete(&vendorCourierModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
