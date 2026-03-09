package domain

import (
	"context"
	"time"
)

// AdminContact represents a contact info entry managed by admin (free text).
type AdminContact struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request DTOs ---

type CreateAdminContactRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateAdminContactRequest struct {
	Content *string `json:"content"`
}

// --- Response DTOs ---

type AdminContactResponse struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Interfaces ---

type AdminContactRepository interface {
	Create(ctx context.Context, c *AdminContact) error
	FindByID(ctx context.Context, id int) (*AdminContact, error)
	List(ctx context.Context) ([]AdminContact, error)
	Update(ctx context.Context, c *AdminContact) error
	Delete(ctx context.Context, id int) error
}

type AdminContactUseCase interface {
	CreateAdminContact(ctx context.Context, req CreateAdminContactRequest) (*AdminContactResponse, error)
	ListAdminContacts(ctx context.Context) ([]AdminContactResponse, error)
	UpdateAdminContact(ctx context.Context, id int, req UpdateAdminContactRequest) (*AdminContactResponse, error)
	DeleteAdminContact(ctx context.Context, id int) error
}
