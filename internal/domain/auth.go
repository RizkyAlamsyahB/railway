package domain

import "context"

// LoginRequest is the input DTO for user authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the output DTO for a successful authentication.
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthUseCase defines the interface for authentication operations.
type AuthUseCase interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}
