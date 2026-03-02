package usecase

import (
	"context"
	"fmt"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type faqUseCase struct {
	faqRepo domain.FAQRepository
}

// NewFAQUseCase creates a new FAQUseCase.
func NewFAQUseCase(faqRepo domain.FAQRepository) domain.FAQUseCase {
	return &faqUseCase{faqRepo: faqRepo}
}

func (uc *faqUseCase) ListFAQs(ctx context.Context, params domain.FAQListParams) ([]domain.FAQCategoryGroup, error) {
	faqs, err := uc.faqRepo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list FAQs: %w", err)
	}

	// Group by category, preserving order from DB (category ASC, sort_order ASC)
	var groups []domain.FAQCategoryGroup
	catIndex := map[string]int{}

	for _, f := range faqs {
		idx, exists := catIndex[f.Category]
		if !exists {
			idx = len(groups)
			catIndex[f.Category] = idx
			groups = append(groups, domain.FAQCategoryGroup{
				Category: f.Category,
			})
		}
		groups[idx].Items = append(groups[idx].Items, domain.FAQResponse{
			ID:       f.ID,
			Category: f.Category,
			Question: f.Question,
			Answer:   f.Answer,
		})
	}
	return groups, nil
}

func (uc *faqUseCase) ListCategories(ctx context.Context) ([]string, error) {
	return uc.faqRepo.ListCategories(ctx)
}
