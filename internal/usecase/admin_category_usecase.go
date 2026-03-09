package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminCategoryUseCase struct {
	categoryRepo domain.AdminCategoryRepository
}

// NewAdminCategoryUseCase creates a new AdminCategoryUseCase.
func NewAdminCategoryUseCase(categoryRepo domain.AdminCategoryRepository) domain.AdminCategoryUseCase {
	return &adminCategoryUseCase{categoryRepo: categoryRepo}
}

func (uc *adminCategoryUseCase) CreateCategory(ctx context.Context, req domain.CreateCategoryRequest) (*domain.CategoryResponse, error) {
	slug := generateCategorySlug(req.Name)

	// Check slug uniqueness.
	existing, err := uc.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug: %w", err)
	}
	if existing != nil {
		return nil, ErrCategorySlugExists
	}

	category := &domain.Category{
		ID:       uuid.New(),
		Name:     req.Name,
		Slug:     slug,
		IsActive: true,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return toCategoryResponse(category), nil
}

func (uc *adminCategoryUseCase) ListCategories(ctx context.Context) ([]domain.CategoryResponse, error) {
	categories, err := uc.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	result := make([]domain.CategoryResponse, len(categories))
	for i, c := range categories {
		result[i] = *toCategoryResponse(&c)
	}
	return result, nil
}

func (uc *adminCategoryUseCase) UpdateCategory(ctx context.Context, id uuid.UUID, req domain.UpdateCategoryRequest) (*domain.CategoryResponse, error) {
	category, err := uc.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil {
		return nil, ErrAdminCategoryNotFound
	}

	if req.Name != nil {
		category.Name = *req.Name
		category.Slug = generateCategorySlug(*req.Name)

		// Check slug uniqueness (excluding self).
		existing, err := uc.categoryRepo.FindBySlug(ctx, category.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, ErrCategorySlugExists
		}
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return toCategoryResponse(category), nil
}

func (uc *adminCategoryUseCase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	category, err := uc.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil {
		return ErrAdminCategoryNotFound
	}

	if err := uc.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func toCategoryResponse(c *domain.Category) *domain.CategoryResponse {
	return &domain.CategoryResponse{
		ID:       c.ID,
		ParentID: c.ParentID,
		Name:     c.Name,
		Slug:     c.Slug,
		IsActive: c.IsActive,
	}
}

func generateCategorySlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove non-alphanumeric except hyphens.
	result := make([]byte, 0, len(slug))
	for i := 0; i < len(slug); i++ {
		c := slug[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			result = append(result, c)
		}
	}
	return string(result)
}
