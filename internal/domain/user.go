package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the users table.
type User struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FullName        string     `json:"full_name"`
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	PasswordHash    string     `json:"-"`
	Status          string     `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Role            *Role      `json:"role,omitempty"`
}

// Role represents the roles table.
type Role struct {
	ID   int16  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreateUserRequest is the input DTO for creating a user.
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email,max=255"`
	FullName  string  `json:"full_name" binding:"required,max=120"`
	Password  string  `json:"password" binding:"required,min=8"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	BirthDate *string `json:"birth_date,omitempty"`
	Role      string  `json:"role" binding:"required,oneof=admin umkm customer cs finance"`
}

// UpdateUserRequest is the input DTO for updating a user.
type UpdateUserRequest struct {
	FullName  *string `json:"full_name,omitempty" binding:"omitempty,max=120"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	BirthDate *string `json:"birth_date,omitempty"`
	Status    *string `json:"status,omitempty" binding:"omitempty,oneof=pending active blocked"`
	Role      *string `json:"role,omitempty" binding:"omitempty,oneof=admin umkm customer cs finance"`
	Password  *string `json:"password,omitempty" binding:"omitempty,min=8"`
}

// UserListParams holds query parameters for listing users.
type UserListParams struct {
	Page   int
	Limit  int
	Role   string
	Status string
	Search string
}

// UserResponse is the output DTO for a user.
type UserResponse struct {
	ID              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FullName        string     `json:"full_name"`
	BirthDate       *time.Time `json:"birth_date,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	Status          string     `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	Role            string     `json:"role"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PaginationMeta holds pagination metadata for list responses.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// AdminUserUseCase defines the interface for admin user management operations.
type AdminUserUseCase interface {
	Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	List(ctx context.Context, params UserListParams) ([]UserResponse, *PaginationMeta, error)
	GetByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserResponse, error)
	Delete(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *User, roleCode string) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, params UserListParams) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	UpdateWithRole(ctx context.Context, user *User, roleCode string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
