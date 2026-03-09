package domain

import (
	"context"

	"github.com/google/uuid"
)

// --- Admin Category Request DTOs ---

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,max=80"`
}

type UpdateCategoryRequest struct {
	Name     *string `json:"name" binding:"omitempty,max=80"`
	IsActive *bool   `json:"is_active"`
}

// --- Admin Category Response DTOs ---

type CategoryResponse struct {
	ID       uuid.UUID  `json:"id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	IsActive bool       `json:"is_active"`
}

// --- Admin Category Interfaces ---

// AdminCategoryRepository extends category persistence with write operations.
type AdminCategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	List(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)
	FindBySlug(ctx context.Context, slug string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// AdminCategoryUseCase defines admin business operations for categories.
type AdminCategoryUseCase interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)
	ListCategories(ctx context.Context) ([]CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, req UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}
