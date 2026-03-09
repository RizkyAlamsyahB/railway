package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// adminCategoryRepository reuses the existing categoryModel from category_repository.go.

type adminCategoryRepository struct {
	db *gorm.DB
}

// NewAdminCategoryRepository creates an AdminCategoryRepository backed by GORM.
func NewAdminCategoryRepository(db *gorm.DB) domain.AdminCategoryRepository {
	return &adminCategoryRepository{db: db}
}

func (r *adminCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	m := fromDomainCategory(category)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *adminCategoryRepository) List(ctx context.Context) ([]domain.Category, error) {
	var models []categoryModel
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	categories := make([]domain.Category, len(models))
	for i, m := range models {
		categories[i] = *toDomainCategory(&m)
	}
	return categories, nil
}

func (r *adminCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var model categoryModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainCategory(&model), nil
}

func (r *adminCategoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var model categoryModel
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainCategory(&model), nil
}

func (r *adminCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	m := fromDomainCategory(category)
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *adminCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&categoryModel{}).Error
}

func fromDomainCategory(c *domain.Category) categoryModel {
	m := categoryModel{
		ID:       c.ID.String(),
		Name:     c.Name,
		Slug:     c.Slug,
		IsActive: c.IsActive,
	}
	if c.ParentID != nil {
		s := c.ParentID.String()
		m.ParentID = &s
	}
	return m
}
