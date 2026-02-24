package domain

import (
	"context"

	"github.com/google/uuid"
)

// CSLoginRequest is the input DTO for customer service authentication.
type CSLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CSLoginResponse is the output DTO for a successful customer service authentication.
type CSLoginResponse struct {
	Token    string    `json:"token"`
	UserID   uuid.UUID `json:"user_id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

// CSAuthUseCase defines the interface for customer service authentication operations.
type CSAuthUseCase interface {
	Login(ctx context.Context, req CSLoginRequest) (*CSLoginResponse, error)
}
