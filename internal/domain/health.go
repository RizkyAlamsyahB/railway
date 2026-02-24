package domain

import "context"

// HealthResponse represents the response from a health check.
type HealthResponse struct {
	Message string `json:"message"`
}

// HealthUseCase defines the interface for health check operations.
type HealthUseCase interface {
	Check(ctx context.Context) HealthResponse
}
