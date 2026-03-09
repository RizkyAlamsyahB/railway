package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminFAQUseCase struct {
	faqRepo domain.FAQRepository
}

// NewAdminFAQUseCase creates a new AdminFAQUseCase.
func NewAdminFAQUseCase(faqRepo domain.FAQRepository) domain.AdminFAQUseCase {
	return &adminFAQUseCase{faqRepo: faqRepo}
}

func (uc *adminFAQUseCase) CreateFAQ(ctx context.Context, req domain.CreateFAQRequest) (*domain.AdminFAQResponse, error) {
	// Validate category.
	if !isValidFAQCategory(req.Category) {
		return nil, ErrInvalidFAQCategory
	}

	now := time.Now()
	faq := &domain.FAQ{
		Category:  req.Category,
		Question:  req.Question,
		Answer:    req.Answer,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.faqRepo.Create(ctx, faq); err != nil {
		return nil, fmt.Errorf("failed to create FAQ: %w", err)
	}

	return toAdminFAQResponse(faq), nil
}

func (uc *adminFAQUseCase) ListFAQs(ctx context.Context) ([]domain.AdminFAQResponse, error) {
	faqs, err := uc.faqRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list FAQs: %w", err)
	}
	result := make([]domain.AdminFAQResponse, len(faqs))
	for i, f := range faqs {
		result[i] = *toAdminFAQResponse(&f)
	}
	return result, nil
}

func (uc *adminFAQUseCase) UpdateFAQ(ctx context.Context, id int, req domain.UpdateFAQRequest) (*domain.AdminFAQResponse, error) {
	faq, err := uc.faqRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find FAQ: %w", err)
	}
	if faq == nil {
		return nil, ErrFAQNotFound
	}

	if req.Category != nil {
		if !isValidFAQCategory(*req.Category) {
			return nil, ErrInvalidFAQCategory
		}
		faq.Category = *req.Category
	}
	if req.Question != nil {
		faq.Question = *req.Question
	}
	if req.Answer != nil {
		faq.Answer = *req.Answer
	}
	faq.UpdatedAt = time.Now()

	if err := uc.faqRepo.Update(ctx, faq); err != nil {
		return nil, fmt.Errorf("failed to update FAQ: %w", err)
	}

	return toAdminFAQResponse(faq), nil
}

func (uc *adminFAQUseCase) DeleteFAQ(ctx context.Context, id int) error {
	faq, err := uc.faqRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find FAQ: %w", err)
	}
	if faq == nil {
		return ErrFAQNotFound
	}

	if err := uc.faqRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete FAQ: %w", err)
	}
	return nil
}

func toAdminFAQResponse(f *domain.FAQ) *domain.AdminFAQResponse {
	return &domain.AdminFAQResponse{
		ID:        f.ID,
		Category:  f.Category,
		Question:  f.Question,
		Answer:    f.Answer,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func isValidFAQCategory(cat string) bool {
	for _, v := range domain.ValidFAQCategories {
		if v == cat {
			return true
		}
	}
	return false
}
