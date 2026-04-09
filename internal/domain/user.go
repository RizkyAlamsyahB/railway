package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the users table.
type User struct {
	ID                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	FullName            string     `json:"full_name"`
	ImageURL            *string    `json:"image_url,omitempty"`
	BirthDate           *time.Time `json:"birth_date,omitempty"`
	Phone               *string    `json:"phone,omitempty"`
	PasswordHash        string     `json:"-"`
	Status              string     `json:"status"`
	EmailVerifiedAt     *time.Time `json:"email_verified_at,omitempty"`
	PasswordChangedAt   *time.Time `json:"-"`
	DeletionReason      *string    `json:"-"`
	DeletionRequestedAt *time.Time `json:"-"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	Role                *Role      `json:"role,omitempty"`
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
	Page      int
	Limit     int
	Role      string
	Status    string
	Search    string
	SortBy    string // allowed: "full_name", "email", "role", "status", "created_at"
	SortOrder string // allowed: "asc", "desc"
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

// AdminMeResponse is the output DTO for the authenticated admin profile endpoint.
type AdminMeResponse struct {
	ID       uuid.UUID `json:"id"`
	ImageURL string    `json:"image_url"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
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
	GetMe(ctx context.Context, adminID uuid.UUID) (*AdminMeResponse, error)
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
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, passwordHash string) error
	ActivateUser(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID, reason string) error
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

// DeleteAccountRequest is the input DTO for customer account deletion.
type DeleteAccountRequest struct {
	Reason   string `json:"reason" binding:"required,max=50"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest is the input DTO for authenticated password change.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ResetPasswordRequest is the input DTO for resetting password via OTP proof token.
type ResetPasswordRequest struct {
	ProofToken  string `json:"proof_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// LoginRequest is the input DTO for user authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the output DTO for a successful user authentication.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// GoogleLoginRequest is the input DTO for Google OAuth login.
type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
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
	LoginWithGoogle(ctx context.Context, req GoogleLoginRequest) (*LoginResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req ChangePasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
	DeleteAccount(ctx context.Context, userID uuid.UUID, req DeleteAccountRequest) error
}
