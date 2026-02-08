package usecase

import (
	"context"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type healthUseCase struct{}

// NewHealthUseCase creates a new instance of HealthUseCase.
func NewHealthUseCase() domain.HealthUseCase {
	return &healthUseCase{}
}

func (uc *healthUseCase) Check(_ context.Context) domain.HealthResponse {
	return domain.HealthResponse{
		Message: "Hello World",
	}
}
