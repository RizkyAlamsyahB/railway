package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type ledgerAccountModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	Code        string `gorm:"column:code"`
	Name        string `gorm:"column:name"`
	AccountType string `gorm:"column:account_type"`
	NormalSide  string `gorm:"column:normal_side"`
	IsActive    bool   `gorm:"column:is_active"`
}

func (ledgerAccountModel) TableName() string { return "ledger_accounts" }

type ledgerJournalModel struct {
	ID          string    `gorm:"column:id;primaryKey"`
	JournalNo   string   `gorm:"column:journal_no"`
	SourceType  string   `gorm:"column:source_type"`
	SourceID    string   `gorm:"column:source_id"`
	EventTime   time.Time `gorm:"column:event_time"`
	Description *string  `gorm:"column:description"`
	Status      string   `gorm:"column:status"`
	CreatedBy   *string  `gorm:"column:created_by"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (ledgerJournalModel) TableName() string { return "ledger_journals" }

type ledgerLineModel struct {
	ID        string  `gorm:"column:id;primaryKey"`
	JournalID string  `gorm:"column:journal_id"`
	AccountID string  `gorm:"column:account_id"`
	Debit     float64 `gorm:"column:debit"`
	Credit    float64 `gorm:"column:credit"`
	Currency  string  `gorm:"column:currency"`
	Reference *string `gorm:"column:reference"`
}

func (ledgerLineModel) TableName() string { return "ledger_lines" }

type ledgerRepository struct {
	db *gorm.DB
}

// NewLedgerRepository creates a new LedgerRepository backed by GORM.
func NewLedgerRepository(db *gorm.DB) domain.LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) FindAccountByCode(ctx context.Context, code string) (*domain.LedgerAccount, error) {
	var model ledgerAccountModel
	if err := r.db.WithContext(ctx).Where("code = ? AND is_active = true", code).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainLedgerAccount(&model), nil
}

func (r *ledgerRepository) CreateJournalWithLines(ctx context.Context, journal *domain.LedgerJournal, lines []domain.LedgerLine) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Insert journal.
		jm := toLedgerJournalModel(journal)
		if err := tx.Create(&jm).Error; err != nil {
			return fmt.Errorf("insert ledger journal: %w", err)
		}

		// Insert lines.
		for _, line := range lines {
			lm := toLedgerLineModel(&line)
			if err := tx.Create(&lm).Error; err != nil {
				return fmt.Errorf("insert ledger line: %w", err)
			}
		}

		return nil
	})
}

// Mapper helpers.

func toDomainLedgerAccount(m *ledgerAccountModel) *domain.LedgerAccount {
	id, _ := uuid.Parse(m.ID)
	return &domain.LedgerAccount{
		ID:          id,
		Code:        m.Code,
		Name:        m.Name,
		AccountType: m.AccountType,
		NormalSide:  m.NormalSide,
		IsActive:    m.IsActive,
	}
}

func toLedgerJournalModel(j *domain.LedgerJournal) ledgerJournalModel {
	var createdBy *string
	if j.CreatedBy != nil {
		s := j.CreatedBy.String()
		createdBy = &s
	}

	return ledgerJournalModel{
		ID:          j.ID.String(),
		JournalNo:   j.JournalNo,
		SourceType:  j.SourceType,
		SourceID:    j.SourceID.String(),
		EventTime:   j.EventTime,
		Description: j.Description,
		Status:      j.Status,
		CreatedBy:   createdBy,
		CreatedAt:   j.CreatedAt,
	}
}

func toLedgerLineModel(l *domain.LedgerLine) ledgerLineModel {
	return ledgerLineModel{
		ID:        l.ID.String(),
		JournalID: l.JournalID.String(),
		AccountID: l.AccountID.String(),
		Debit:     l.Debit,
		Credit:    l.Credit,
		Currency:  l.Currency,
		Reference: l.Reference,
	}
}
