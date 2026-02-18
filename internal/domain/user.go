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
	FindByPhone(ctx context.Context, phone string) (*User, error)
	List(ctx context.Context, params UserListParams) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	UpdateWithRole(ctx context.Context, user *User, roleCode string) error
	ActivateUser(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// --- Customer Registration & Email Verification ---

// EmailVerificationToken represents the email_verification_tokens table.
type EmailVerificationToken struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Email         string
	TokenHash     string
	ExpiresAt     time.Time
	ConsumedAt    *time.Time
	InvalidatedAt *time.Time
	CreatedAt     time.Time
}

// RegisterCustomerRequest is the input DTO for customer self-registration.
type RegisterCustomerRequest struct {
	FullName string `json:"full_name" binding:"required,max=120"`
	Phone    string `json:"phone" binding:"required,max=20"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterCustomerResponse is the output DTO for a successful customer registration.
type RegisterCustomerResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	Message string    `json:"message"`
}

// ResendVerificationRequest is the input DTO for resending verification email.
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// LoginRequest is the input DTO for user authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the output DTO for a successful user authentication.
type LoginResponse struct {
	Token string `json:"token"`
	User  UserResponse `json:"user"`
}

// EmailVerificationTokenRepository defines the interface for email verification token data access.
type EmailVerificationTokenRepository interface {
	Create(ctx context.Context, token *EmailVerificationToken) error
	FindActiveByTokenHash(ctx context.Context, tokenHash string) (*EmailVerificationToken, error)
	Consume(ctx context.Context, id uuid.UUID) error
	InvalidateByUserID(ctx context.Context, userID uuid.UUID) error
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*EmailVerificationToken, error)
}

// UserUseCase defines the interface for customer-facing user operations.
type UserUseCase interface {
	Register(ctx context.Context, req RegisterCustomerRequest) (*RegisterCustomerResponse, error)
	VerifyEmail(ctx context.Context, rawToken string) error
	ResendVerification(ctx context.Context, req ResendVerificationRequest) error
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}
