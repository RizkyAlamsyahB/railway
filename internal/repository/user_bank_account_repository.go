package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
	"gorm.io/gorm"
)

type userBankAccountModel struct {
	ID                string    `gorm:"column:id;primaryKey"`
	UserID            string    `gorm:"column:user_id"`
	ChannelCode       string    `gorm:"column:channel_code"`
	BankName          string    `gorm:"column:bank_name"`
	AccountNumber     string    `gorm:"column:account_number"`
	AccountHolderName string    `gorm:"column:account_holder_name"`
	AccountLast4      string    `gorm:"column:account_last4"`
	IsDefault         bool      `gorm:"column:is_default"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (userBankAccountModel) TableName() string { return "user_bank_accounts" }

type userBankAccountRepository struct {
	db          *gorm.DB
	fieldCipher *sensitivedata.FieldCipher
}

// NewUserBankAccountRepository creates a UserBankAccountRepository backed by GORM.
func NewUserBankAccountRepository(db *gorm.DB, fieldCipher *sensitivedata.FieldCipher) domain.UserBankAccountRepository {
	return &userBankAccountRepository{db: db, fieldCipher: fieldCipher}
}

func (r *userBankAccountRepository) Create(ctx context.Context, account *domain.UserBankAccount) error {
	encrypted := account.AccountNumber
	if strings.TrimSpace(encrypted) != "" && r.fieldCipher != nil {
		enc, err := r.fieldCipher.EncryptString(encrypted)
		if err != nil {
			return fmt.Errorf("failed to encrypt account number: %w", err)
		}
		encrypted = enc
	}

	model := userBankAccountModel{
		ID:                account.ID.String(),
		UserID:            account.UserID.String(),
		ChannelCode:       account.ChannelCode,
		BankName:          account.BankName,
		AccountNumber:     encrypted,
		AccountHolderName: account.AccountHolderName,
		AccountLast4:      account.AccountLast4,
		IsDefault:         account.IsDefault,
		CreatedAt:         account.CreatedAt,
		UpdatedAt:         account.UpdatedAt,
	}

	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *userBankAccountRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserBankAccountListItem, error) {
	var models []userBankAccountModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID.String()).
		Order("is_default DESC, created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.UserBankAccountListItem, 0, len(models))
	for _, m := range models {
		id, _ := uuid.Parse(m.ID)
		items = append(items, domain.UserBankAccountListItem{
			ID:                id,
			ChannelCode:       m.ChannelCode,
			BankName:          m.BankName,
			AccountHolderName: m.AccountHolderName,
			AccountLast4:      m.AccountLast4,
			IsDefault:         m.IsDefault,
		})
	}
	return items, nil
}

func (r *userBankAccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.UserBankAccount, error) {
	var m userBankAccountModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&m)
}

func (r *userBankAccountRepository) FindDefaultByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserBankAccount, error) {
	var m userBankAccountModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = true", userID.String()).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&m)
}

func (r *userBankAccountRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id.String(), userID.String()).
		Delete(&userBankAccountModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userBankAccountRepository) SetDefault(ctx context.Context, userID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset current default
		if err := tx.Table("user_bank_accounts").
			Where("user_id = ? AND is_default = true", userID.String()).
			Update("is_default", false).Error; err != nil {
			return err
		}
		// Set new default
		result := tx.Table("user_bank_accounts").
			Where("id = ? AND user_id = ?", id.String(), userID.String()).
			Update("is_default", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *userBankAccountRepository) toDomain(m *userBankAccountModel) (*domain.UserBankAccount, error) {
	accountNumber := m.AccountNumber
	if strings.TrimSpace(accountNumber) != "" && r.fieldCipher != nil {
		dec, err := r.fieldCipher.DecryptString(accountNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt account number: %w", err)
		}
		accountNumber = dec
	}

	id, _ := uuid.Parse(m.ID)
	userID, _ := uuid.Parse(m.UserID)
	return &domain.UserBankAccount{
		ID:                id,
		UserID:            userID,
		ChannelCode:       m.ChannelCode,
		BankName:          m.BankName,
		AccountNumber:     accountNumber,
		AccountHolderName: m.AccountHolderName,
		AccountLast4:      m.AccountLast4,
		IsDefault:         m.IsDefault,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}, nil
}
