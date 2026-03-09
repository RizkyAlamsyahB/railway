package repository

import (
	"context"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type adminContactModel struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (adminContactModel) TableName() string { return "admin_contacts" }

type adminContactRepository struct {
	db *gorm.DB
}

// NewAdminContactRepository creates an AdminContactRepository backed by GORM.
func NewAdminContactRepository(db *gorm.DB) domain.AdminContactRepository {
	return &adminContactRepository{db: db}
}

func (r *adminContactRepository) Create(ctx context.Context, c *domain.AdminContact) error {
	m := adminContactModel{
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	c.ID = m.ID
	return nil
}

func (r *adminContactRepository) FindByID(ctx context.Context, id int) (*domain.AdminContact, error) {
	var m adminContactModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toAdminContactDomain(m), nil
}

func (r *adminContactRepository) List(ctx context.Context) ([]domain.AdminContact, error) {
	q := r.db.WithContext(ctx).Model(&adminContactModel{})
	q = q.Order("id ASC")

	var rows []adminContactModel
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.AdminContact, len(rows))
	for i, row := range rows {
		items[i] = *toAdminContactDomain(row)
	}
	return items, nil
}

func (r *adminContactRepository) Update(ctx context.Context, c *domain.AdminContact) error {
	m := adminContactModel{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *adminContactRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&adminContactModel{}).Error
}

func toAdminContactDomain(m adminContactModel) *domain.AdminContact {
	return &domain.AdminContact{
		ID:        m.ID,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
