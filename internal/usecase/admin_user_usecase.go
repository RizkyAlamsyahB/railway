package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

type adminUserUseCase struct {
	userRepo domain.UserRepository
}

// NewAdminUserUseCase creates a new AdminUserUseCase.
func NewAdminUserUseCase(userRepo domain.UserRepository) domain.AdminUserUseCase {
	return &adminUserUseCase{userRepo: userRepo}
}

func (uc *adminUserUseCase) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.UserResponse, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	if req.Phone != nil {
		existingPhone, err := uc.userRepo.FindByPhone(ctx, *req.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to check phone: %w", err)
		}
		if existingPhone != nil {
			return nil, ErrPhoneAlreadyRegistered
		}
	}

	var birthDate *time.Time
	if req.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return nil, ErrInvalidBirthDate
		}
		birthDate = &t
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		ID:              uuid.New(),
		Email:           req.Email,
		FullName:        req.FullName,
		BirthDate:       birthDate,
		Phone:           req.Phone,
		PasswordHash:    hash,
		Status:          domain.UserStatusActive,
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.userRepo.Create(ctx, user, req.Role); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	created, err := uc.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}

	return toUserResponse(created), nil
}

func (uc *adminUserUseCase) List(ctx context.Context, params domain.UserListParams) ([]domain.UserResponse, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	users, total, err := uc.userRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list users: %w", err)
	}

	responses := make([]domain.UserResponse, len(users))
	for i, u := range users {
		responses[i] = *toUserResponse(&u)
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return responses, meta, nil
}

func (uc *adminUserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return toUserResponse(user), nil
}

func (uc *adminUserUseCase) GetMe(ctx context.Context, adminID uuid.UUID) (*domain.AdminMeResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.Role == nil || user.Role.Code != domain.RoleAdmin {
		return nil, ErrNotAdmin
	}

	return toAdminMeResponse(user), nil
}

func (uc *adminUserUseCase) Update(ctx context.Context, id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.BirthDate != nil {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return nil, ErrInvalidBirthDate
		}
		user.BirthDate = &t
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Password != nil {
		hash, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hash
	}
	user.UpdatedAt = time.Now()

	if req.Role != nil {
		if err := uc.userRepo.UpdateWithRole(ctx, user, *req.Role); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		if err := uc.userRepo.Update(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	updated, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return toUserResponse(updated), nil
}

func (uc *adminUserUseCase) Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error {
	if id == actorID {
		return ErrCannotDeleteSelf
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
