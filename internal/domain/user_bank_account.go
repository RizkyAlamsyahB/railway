package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserBankAccount represents a customer's saved bank account.
type UserBankAccount struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	ChannelCode       string    `json:"channel_code"`
	BankName          string    `json:"bank_name"`
	AccountNumber     string    `json:"-"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountLast4      string    `json:"account_last4"`
	IsDefault         bool      `json:"is_default"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserBankAccountListItem is the DTO returned when listing bank accounts (no full account number).
type UserBankAccountListItem struct {
	ID                uuid.UUID `json:"id"`
	ChannelCode       string    `json:"channel_code"`
	BankName          string    `json:"bank_name"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountLast4      string    `json:"account_last4"`
	IsDefault         bool      `json:"is_default"`
}

// CreateUserBankAccountRequest is the input DTO for adding a bank account.
type CreateUserBankAccountRequest struct {
	ChannelCode       string `json:"channel_code" binding:"required,max=40"`
	BankName          string `json:"bank_name" binding:"required,max=80"`
	AccountNumber     string `json:"account_number" binding:"required,max=64"`
	AccountHolderName string `json:"account_holder_name" binding:"required,max=120"`
	IsDefault         bool   `json:"is_default"`
}

// CreateUserBankAccountResponse is the output DTO after creating a bank account.
type CreateUserBankAccountResponse struct {
	ID                uuid.UUID `json:"id"`
	ChannelCode       string    `json:"channel_code"`
	BankName          string    `json:"bank_name"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountLast4      string    `json:"account_last4"`
	IsDefault         bool      `json:"is_default"`
}

// UserBankAccountRepository defines data access for user bank accounts.
type UserBankAccountRepository interface {
	Create(ctx context.Context, account *UserBankAccount) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]UserBankAccountListItem, error)
	FindByID(ctx context.Context, id uuid.UUID) (*UserBankAccount, error)
	Delete(ctx context.Context, userID, id uuid.UUID) error
	SetDefault(ctx context.Context, userID, id uuid.UUID) error
	FindDefaultByUserID(ctx context.Context, userID uuid.UUID) (*UserBankAccount, error)
}

// UserBankAccountUseCase defines use cases for user bank accounts.
type UserBankAccountUseCase interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateUserBankAccountRequest) (*CreateUserBankAccountResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]UserBankAccountListItem, error)
	Delete(ctx context.Context, userID, id uuid.UUID) error
	SetDefault(ctx context.Context, userID, id uuid.UUID) error
}
