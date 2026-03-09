package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminContactUseCase struct {
	repo domain.AdminContactRepository
}

// NewAdminContactUseCase creates a new AdminContactUseCase.
func NewAdminContactUseCase(repo domain.AdminContactRepository) domain.AdminContactUseCase {
	return &adminContactUseCase{repo: repo}
}

func (uc *adminContactUseCase) CreateAdminContact(ctx context.Context, req domain.CreateAdminContactRequest) (*domain.AdminContactResponse, error) {
	now := time.Now()
	c := &domain.AdminContact{
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("failed to create admin contact: %w", err)
	}

	return toAdminContactResponse(c), nil
}

func (uc *adminContactUseCase) ListAdminContacts(ctx context.Context) ([]domain.AdminContactResponse, error) {
	items, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list admin contacts: %w", err)
	}
	result := make([]domain.AdminContactResponse, len(items))
	for i, item := range items {
		result[i] = *toAdminContactResponse(&item)
	}
	return result, nil
}

func (uc *adminContactUseCase) UpdateAdminContact(ctx context.Context, id int, req domain.UpdateAdminContactRequest) (*domain.AdminContactResponse, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin contact: %w", err)
	}
	if c == nil {
		return nil, ErrAdminContactNotFound
	}

	if req.Content != nil {
		c.Content = *req.Content
	}
	c.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("failed to update admin contact: %w", err)
	}

	return toAdminContactResponse(c), nil
}

func (uc *adminContactUseCase) DeleteAdminContact(ctx context.Context, id int) error {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find admin contact: %w", err)
	}
	if c == nil {
		return ErrAdminContactNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete admin contact: %w", err)
	}
	return nil
}

func toAdminContactResponse(c *domain.AdminContact) *domain.AdminContactResponse {
	return &domain.AdminContactResponse{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
