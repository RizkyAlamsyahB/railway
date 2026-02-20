package usecase

import (
	"context"
	"fmt"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

type csAuthUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
	jwtExpiry int
	jwtIssuer string
}

// NewCSAuthUseCase creates a new CSAuthUseCase.
func NewCSAuthUseCase(userRepo domain.UserRepository, jwtSecret string, jwtExpiry int, jwtIssuer string) domain.CSAuthUseCase {
	return &csAuthUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
		jwtIssuer: jwtIssuer,
	}
}

func (uc *csAuthUseCase) Login(ctx context.Context, req domain.CSLoginRequest) (*domain.CSLoginResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status != domain.UserStatusActive {
		return nil, ErrAccountInactive
	}

	if user.Role == nil || user.Role.Code != domain.RoleCS {
		return nil, ErrNotCS
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role.Code, nil, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.CSLoginResponse{
		Token:    token,
		UserID:   user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}
