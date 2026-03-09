package domain

import (
	"context"
	"time"
)

// ReturnReason represents a return reason entry managed by admin.
type ReturnReason struct {
	ID        int       `json:"id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request DTOs ---

type CreateReturnReasonRequest struct {
	Reason string `json:"reason" binding:"required,max=255"`
}

type UpdateReturnReasonRequest struct {
	Reason *string `json:"reason" binding:"omitempty,max=255"`
}

// --- Response DTOs ---

type ReturnReasonResponse struct {
	ID        int       `json:"id"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Interfaces ---

type ReturnReasonRepository interface {
	Create(ctx context.Context, r *ReturnReason) error
	FindByID(ctx context.Context, id int) (*ReturnReason, error)
	List(ctx context.Context) ([]ReturnReason, error)
	Update(ctx context.Context, r *ReturnReason) error
	Delete(ctx context.Context, id int) error
}

type ReturnReasonUseCase interface {
	CreateReturnReason(ctx context.Context, req CreateReturnReasonRequest) (*ReturnReasonResponse, error)
	ListReturnReasons(ctx context.Context) ([]ReturnReasonResponse, error)
	UpdateReturnReason(ctx context.Context, id int, req UpdateReturnReasonRequest) (*ReturnReasonResponse, error)
	DeleteReturnReason(ctx context.Context, id int) error
}
