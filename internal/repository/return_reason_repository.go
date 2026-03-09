package repository

import (
	"context"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type returnReasonModel struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Reason    string    `gorm:"column:reason"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (returnReasonModel) TableName() string { return "return_reasons" }

type returnReasonRepository struct {
	db *gorm.DB
}

// NewReturnReasonRepository creates a ReturnReasonRepository backed by GORM.
func NewReturnReasonRepository(db *gorm.DB) domain.ReturnReasonRepository {
	return &returnReasonRepository{db: db}
}

func (r *returnReasonRepository) Create(ctx context.Context, rr *domain.ReturnReason) error {
	m := returnReasonModel{
		Reason:    rr.Reason,
		CreatedAt: rr.CreatedAt,
		UpdatedAt: rr.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	rr.ID = m.ID
	return nil
}

func (r *returnReasonRepository) FindByID(ctx context.Context, id int) (*domain.ReturnReason, error) {
	var m returnReasonModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toReturnReasonDomain(m), nil
}

func (r *returnReasonRepository) List(ctx context.Context) ([]domain.ReturnReason, error) {
	q := r.db.WithContext(ctx).Model(&returnReasonModel{})
	q = q.Order("id ASC")

	var rows []returnReasonModel
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.ReturnReason, len(rows))
	for i, row := range rows {
		items[i] = *toReturnReasonDomain(row)
	}
	return items, nil
}

func (r *returnReasonRepository) Update(ctx context.Context, rr *domain.ReturnReason) error {
	m := returnReasonModel{
		ID:        rr.ID,
		Reason:    rr.Reason,
		CreatedAt: rr.CreatedAt,
		UpdatedAt: rr.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *returnReasonRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&returnReasonModel{}).Error
}

func toReturnReasonDomain(m returnReasonModel) *domain.ReturnReason {
	return &domain.ReturnReason{
		ID:        m.ID,
		Reason:    m.Reason,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
