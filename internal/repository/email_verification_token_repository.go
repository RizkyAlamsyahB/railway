package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type emailVerificationTokenModel struct {
	ID            string     `gorm:"column:id;primaryKey"`
	UserID        string     `gorm:"column:user_id"`
	Email         string     `gorm:"column:email"`
	TokenHash     string     `gorm:"column:token_hash"`
	ExpiresAt     time.Time  `gorm:"column:expires_at"`
	ConsumedAt    *time.Time `gorm:"column:consumed_at"`
	InvalidatedAt *time.Time `gorm:"column:invalidated_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}

func (emailVerificationTokenModel) TableName() string { return "email_verification_tokens" }

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

// NewEmailVerificationTokenRepository creates a new EmailVerificationTokenRepository backed by GORM.
func NewEmailVerificationTokenRepository(db *gorm.DB) domain.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(ctx context.Context, token *domain.EmailVerificationToken) error {
	model := toEmailVerificationTokenModel(token)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *emailVerificationTokenRepository) FindActiveByTokenHash(ctx context.Context, tokenHash string) (*domain.EmailVerificationToken, error) {
	var model emailVerificationTokenModel
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND consumed_at IS NULL AND invalidated_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainEmailVerificationToken(&model), nil
}

func (r *emailVerificationTokenRepository) Consume(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&emailVerificationTokenModel{}).
		Where("id = ?", id.String()).
		Update("consumed_at", now).Error
}

func (r *emailVerificationTokenRepository) InvalidateByUserID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&emailVerificationTokenModel{}).
		Where("user_id = ? AND consumed_at IS NULL AND invalidated_at IS NULL", userID.String()).
		Update("invalidated_at", now).Error
}

func (r *emailVerificationTokenRepository) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*domain.EmailVerificationToken, error) {
	var model emailVerificationTokenModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID.String()).
		Order("created_at DESC").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainEmailVerificationToken(&model), nil
}

func toEmailVerificationTokenModel(t *domain.EmailVerificationToken) emailVerificationTokenModel {
	return emailVerificationTokenModel{
		ID:            t.ID.String(),
		UserID:        t.UserID.String(),
		Email:         t.Email,
		TokenHash:     t.TokenHash,
		ExpiresAt:     t.ExpiresAt,
		ConsumedAt:    t.ConsumedAt,
		InvalidatedAt: t.InvalidatedAt,
		CreatedAt:     t.CreatedAt,
	}
}

func toDomainEmailVerificationToken(m *emailVerificationTokenModel) *domain.EmailVerificationToken {
	id, _ := uuid.Parse(m.ID)
	userID, _ := uuid.Parse(m.UserID)
	return &domain.EmailVerificationToken{
		ID:            id,
		UserID:        userID,
		Email:         m.Email,
		TokenHash:     m.TokenHash,
		ExpiresAt:     m.ExpiresAt,
		ConsumedAt:    m.ConsumedAt,
		InvalidatedAt: m.InvalidatedAt,
		CreatedAt:     m.CreatedAt,
	}
}
