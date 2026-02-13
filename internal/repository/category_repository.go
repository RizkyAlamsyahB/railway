package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type categoryModel struct {
	ID       string  `gorm:"column:id;primaryKey"`
	ParentID *string `gorm:"column:parent_id"`
	Name     string  `gorm:"column:name"`
	Slug     string  `gorm:"column:slug"`
	IsActive bool    `gorm:"column:is_active"`
}

func (categoryModel) TableName() string { return "categories" }

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository backed by GORM.
func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) ListActive(ctx context.Context) ([]domain.Category, error) {
	var models []categoryModel
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]domain.Category, len(models))
	for i, m := range models {
		categories[i] = *toDomainCategory(&m)
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var model categoryModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainCategory(&model), nil
}

func toDomainCategory(m *categoryModel) *domain.Category {
	id, _ := uuid.Parse(m.ID)

	c := &domain.Category{
		ID:       id,
		Name:     m.Name,
		Slug:     m.Slug,
		IsActive: m.IsActive,
	}
	if m.ParentID != nil {
		parentID, _ := uuid.Parse(*m.ParentID)
		c.ParentID = &parentID
	}
	return c
}
