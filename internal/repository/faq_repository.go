package repository

import (
	"context"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// faqModel is the GORM model for the faqs table.
type faqModel struct {
	ID        int    `gorm:"column:id;primaryKey"`
	Category  string `gorm:"column:category"`
	Question  string `gorm:"column:question"`
	Answer    string `gorm:"column:answer"`
	SortOrder int    `gorm:"column:sort_order"`
	IsActive  bool   `gorm:"column:is_active"`
}

func (faqModel) TableName() string { return "faqs" }

type faqRepository struct {
	db *gorm.DB
}

// NewFAQRepository creates a new FAQRepository backed by GORM.
func NewFAQRepository(db *gorm.DB) domain.FAQRepository {
	return &faqRepository{db: db}
}

func (r *faqRepository) List(ctx context.Context, params domain.FAQListParams) ([]domain.FAQ, error) {
	q := r.db.WithContext(ctx).Model(&faqModel{}).Where("is_active = TRUE")

	if params.Category != "" {
		q = q.Where("category = ?", params.Category)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		q = q.Where("(LOWER(question) LIKE LOWER(?) OR LOWER(answer) LIKE LOWER(?))", like, like)
	}

	q = q.Order("category ASC, sort_order ASC, id ASC")

	var rows []faqModel
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]domain.FAQ, len(rows))
	for i, row := range rows {
		items[i] = domain.FAQ{
			ID:        row.ID,
			Category:  row.Category,
			Question:  row.Question,
			Answer:    row.Answer,
			SortOrder: row.SortOrder,
			IsActive:  row.IsActive,
		}
	}
	return items, nil
}

func (r *faqRepository) ListCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).
		Model(&faqModel{}).
		Where("is_active = TRUE").
		Distinct("category").
		Order("category ASC").
		Pluck("category", &categories).Error
	return categories, err
}
