package repository

import (
	"context"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// faqModel is the GORM model for the faqs table.
type faqModel struct {
	ID       int    `gorm:"column:id;primaryKey"`
	Category string `gorm:"column:category"`
	Question string `gorm:"column:question"`
	Answer   string `gorm:"column:answer"`
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
	q := r.db.WithContext(ctx).Model(&faqModel{})

	if params.Category != "" {
		q = q.Where("category = ?", params.Category)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		q = q.Where("(LOWER(question) LIKE LOWER(?) OR LOWER(answer) LIKE LOWER(?))", like, like)
	}

	q = q.Order("category ASC, id ASC")

	var rows []faqModel
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]domain.FAQ, len(rows))
	for i, row := range rows {
		items[i] = domain.FAQ{
			ID:       row.ID,
			Category: row.Category,
			Question: row.Question,
			Answer:   row.Answer,
		}
	}
	return items, nil
}

func (r *faqRepository) ListCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).
		Model(&faqModel{}).
		Distinct("category").
		Order("category ASC").
		Pluck("category", &categories).Error
	return categories, err
}

func (r *faqRepository) FindByID(ctx context.Context, id int) (*domain.FAQ, error) {
	var m faqModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toFAQDomain(m), nil
}

func (r *faqRepository) Create(ctx context.Context, faq *domain.FAQ) error {
	m := faqModel{
		Category: faq.Category,
		Question: faq.Question,
		Answer:   faq.Answer,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	faq.ID = m.ID
	return nil
}

func (r *faqRepository) Update(ctx context.Context, faq *domain.FAQ) error {
	m := faqModel{
		ID:       faq.ID,
		Category: faq.Category,
		Question: faq.Question,
		Answer:   faq.Answer,
	}
	return r.db.WithContext(ctx).Save(&m).Error
}

func (r *faqRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&faqModel{}).Error
}

func (r *faqRepository) ListAll(ctx context.Context) ([]domain.FAQ, error) {
	var rows []faqModel
	if err := r.db.WithContext(ctx).
		Order("category ASC, id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]domain.FAQ, len(rows))
	for i, row := range rows {
		items[i] = *toFAQDomain(row)
	}
	return items, nil
}

func toFAQDomain(m faqModel) *domain.FAQ {
	return &domain.FAQ{
		ID:       m.ID,
		Category: m.Category,
		Question: m.Question,
		Answer:   m.Answer,
	}
}
