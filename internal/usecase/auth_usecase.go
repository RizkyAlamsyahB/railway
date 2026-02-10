package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotAdmin           = errors.New("user does not have admin access")
	ErrAccountInactive    = errors.New("account is not active")
)

type authUseCase struct {
	userRepo  domain.UserRepository
	jwtSecret string
	jwtExpiry int
	jwtIssuer string
}

// NewAuthUseCase creates a new AuthUseCase.
func NewAuthUseCase(userRepo domain.UserRepository, jwtSecret string, jwtExpiry int, jwtIssuer string) domain.AuthUseCase {
	return &authUseCase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
		jwtIssuer: jwtIssuer,
	}
}

func (uc *authUseCase) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
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

	if user.Status != "active" {
		return nil, ErrAccountInactive
	}

	hasAdmin := false
	roleCodes := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roleCodes[i] = r.Code
		if r.Code == "admin" {
			hasAdmin = true
		}
	}
	if !hasAdmin {
		return nil, ErrNotAdmin
	}

	token, err := auth.GenerateToken(user.ID, user.Email, roleCodes, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.LoginResponse{Token: token}, nil
}
