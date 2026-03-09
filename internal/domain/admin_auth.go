package domain

import (
	"context"

	"github.com/google/uuid"
)

// AdminLoginRequest is the input DTO for admin authentication.
type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse is the output DTO for a successful admin authentication.
type AdminLoginResponse struct {
	Token    string    `json:"token"`
	ID       uuid.UUID `json:"id"`
	ImageURL string    `json:"image_url"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
}

// AdminAuthUseCase defines the interface for admin authentication operations.
type AdminAuthUseCase interface {
	Login(ctx context.Context, req AdminLoginRequest) (*AdminLoginResponse, error)
}
