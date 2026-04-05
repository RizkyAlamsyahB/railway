package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type userBankAccountUseCase struct {
	repo domain.UserBankAccountRepository
}

// NewUserBankAccountUseCase creates a new UserBankAccountUseCase.
func NewUserBankAccountUseCase(repo domain.UserBankAccountRepository) domain.UserBankAccountUseCase {
	return &userBankAccountUseCase{repo: repo}
}

func (uc *userBankAccountUseCase) Create(ctx context.Context, userID uuid.UUID, req domain.CreateUserBankAccountRequest) (*domain.CreateUserBankAccountResponse, error) {
	accountNumber := req.AccountNumber
	last4 := accountNumber
	if len(last4) > 4 {
		last4 = last4[len(last4)-4:]
	}

	account := &domain.UserBankAccount{
		ID:                uuid.New(),
		UserID:            userID,
		ChannelCode:       req.ChannelCode,
		BankName:          req.BankName,
		AccountNumber:     accountNumber,
		AccountHolderName: req.AccountHolderName,
		AccountLast4:      last4,
		IsDefault:         req.IsDefault,
	}

	if err := uc.repo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	return &domain.CreateUserBankAccountResponse{
		ID:                account.ID,
		ChannelCode:       account.ChannelCode,
		BankName:          account.BankName,
		AccountHolderName: account.AccountHolderName,
		AccountLast4:      account.AccountLast4,
		IsDefault:         account.IsDefault,
	}, nil
}

func (uc *userBankAccountUseCase) List(ctx context.Context, userID uuid.UUID) ([]domain.UserBankAccountListItem, error) {
	items, err := uc.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list bank accounts: %w", err)
	}
	return items, nil
}

func (uc *userBankAccountUseCase) Delete(ctx context.Context, userID, id uuid.UUID) error {
	account, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find bank account: %w", err)
	}
	if account == nil {
		return ErrBankAccountNotFound
	}
	if account.UserID != userID {
		return ErrBankAccountNotOwned
	}

	if err := uc.repo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("failed to delete bank account: %w", err)
	}
	return nil
}

func (uc *userBankAccountUseCase) SetDefault(ctx context.Context, userID, id uuid.UUID) error {
	account, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find bank account: %w", err)
	}
	if account == nil {
		return ErrBankAccountNotFound
	}
	if account.UserID != userID {
		return ErrBankAccountNotOwned
	}

	if err := uc.repo.SetDefault(ctx, userID, id); err != nil {
		return fmt.Errorf("failed to set default bank account: %w", err)
	}
	return nil
}
