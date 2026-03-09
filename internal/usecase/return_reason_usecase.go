package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type returnReasonUseCase struct {
	repo domain.ReturnReasonRepository
}

// NewReturnReasonUseCase creates a new ReturnReasonUseCase.
func NewReturnReasonUseCase(repo domain.ReturnReasonRepository) domain.ReturnReasonUseCase {
	return &returnReasonUseCase{repo: repo}
}

func (uc *returnReasonUseCase) CreateReturnReason(ctx context.Context, req domain.CreateReturnReasonRequest) (*domain.ReturnReasonResponse, error) {
	now := time.Now()
	rr := &domain.ReturnReason{
		Reason:    req.Reason,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.repo.Create(ctx, rr); err != nil {
		return nil, fmt.Errorf("failed to create return reason: %w", err)
	}

	return toReturnReasonResponse(rr), nil
}

func (uc *returnReasonUseCase) ListReturnReasons(ctx context.Context) ([]domain.ReturnReasonResponse, error) {
	items, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list return reasons: %w", err)
	}
	result := make([]domain.ReturnReasonResponse, len(items))
	for i, item := range items {
		result[i] = *toReturnReasonResponse(&item)
	}
	return result, nil
}

func (uc *returnReasonUseCase) UpdateReturnReason(ctx context.Context, id int, req domain.UpdateReturnReasonRequest) (*domain.ReturnReasonResponse, error) {
	rr, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find return reason: %w", err)
	}
	if rr == nil {
		return nil, ErrReturnReasonNotFound
	}

	if req.Reason != nil {
		rr.Reason = *req.Reason
	}
	rr.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, rr); err != nil {
		return nil, fmt.Errorf("failed to update return reason: %w", err)
	}

	return toReturnReasonResponse(rr), nil
}

func (uc *returnReasonUseCase) DeleteReturnReason(ctx context.Context, id int) error {
	rr, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find return reason: %w", err)
	}
	if rr == nil {
		return ErrReturnReasonNotFound
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete return reason: %w", err)
	}
	return nil
}

func toReturnReasonResponse(rr *domain.ReturnReason) *domain.ReturnReasonResponse {
	return &domain.ReturnReasonResponse{
		ID:        rr.ID,
		Reason:    rr.Reason,
		CreatedAt: rr.CreatedAt,
		UpdatedAt: rr.UpdatedAt,
	}
}
