package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type otpCodeModel struct {
	ID            string     `gorm:"column:id;primaryKey"`
	UserID        *string    `gorm:"column:user_id"`
	Email         string     `gorm:"column:email"`
	Purpose       string     `gorm:"column:purpose"`
	Channel       string     `gorm:"column:channel"`
	CodeHash      string     `gorm:"column:code_hash"`
	ExpiresAt     time.Time  `gorm:"column:expires_at"`
	AttemptCount  int        `gorm:"column:attempt_count"`
	ConsumedAt    *time.Time `gorm:"column:consumed_at"`
	InvalidatedAt *time.Time `gorm:"column:invalidated_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}

func (otpCodeModel) TableName() string { return "otp_codes" }

type otpRepository struct {
	db *gorm.DB
}

// NewOTPRepository creates a new OTPRepository backed by GORM.
func NewOTPRepository(db *gorm.DB) domain.OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *domain.OTPCode) error {
	model := toOTPCodeModel(otp)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *otpRepository) FindActiveByEmailPurpose(ctx context.Context, email, purpose, channel string) (*domain.OTPCode, error) {
	var model otpCodeModel
	err := r.db.WithContext(ctx).
		Where("lower(email) = lower(?) AND purpose = ? AND channel = ? AND consumed_at IS NULL AND invalidated_at IS NULL", email, purpose, channel).
		Order("created_at DESC").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainOTPCode(&model), nil
}

func (r *otpRepository) FindLatestByEmailPurpose(ctx context.Context, email, purpose, channel string) (*domain.OTPCode, error) {
	var model otpCodeModel
	err := r.db.WithContext(ctx).
		Where("lower(email) = lower(?) AND purpose = ? AND channel = ?", email, purpose, channel).
		Order("created_at DESC").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainOTPCode(&model), nil
}

func (r *otpRepository) InvalidateActiveByEmailPurpose(ctx context.Context, email, purpose, channel string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&otpCodeModel{}).
		Where("lower(email) = lower(?) AND purpose = ? AND channel = ? AND consumed_at IS NULL AND invalidated_at IS NULL", email, purpose, channel).
		Update("invalidated_at", now).Error
}

func (r *otpRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&otpCodeModel{}).
		Where("id = ?", id.String()).
		UpdateColumn("attempt_count", gorm.Expr("attempt_count + 1")).Error
}

func (r *otpRepository) Consume(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&otpCodeModel{}).
		Where("id = ?", id.String()).
		Update("consumed_at", now).Error
}

func toOTPCodeModel(otp *domain.OTPCode) otpCodeModel {
	var userID *string
	if otp.UserID != nil {
		s := otp.UserID.String()
		userID = &s
	}
	return otpCodeModel{
		ID:            otp.ID.String(),
		UserID:        userID,
		Email:         otp.Email,
		Purpose:       otp.Purpose,
		Channel:       otp.Channel,
		CodeHash:      otp.CodeHash,
		ExpiresAt:     otp.ExpiresAt,
		AttemptCount:  otp.AttemptCount,
		ConsumedAt:    otp.ConsumedAt,
		InvalidatedAt: otp.InvalidatedAt,
		CreatedAt:     otp.CreatedAt,
	}
}

func toDomainOTPCode(model *otpCodeModel) *domain.OTPCode {
	id, _ := uuid.Parse(model.ID)
	var userID *uuid.UUID
	if model.UserID != nil {
		parsed, err := uuid.Parse(*model.UserID)
		if err == nil {
			userID = &parsed
		}
	}
	return &domain.OTPCode{
		ID:            id,
		UserID:        userID,
		Email:         model.Email,
		Purpose:       model.Purpose,
		Channel:       model.Channel,
		CodeHash:      model.CodeHash,
		ExpiresAt:     model.ExpiresAt,
		AttemptCount:  model.AttemptCount,
		ConsumedAt:    model.ConsumedAt,
		InvalidatedAt: model.InvalidatedAt,
		CreatedAt:     model.CreatedAt,
	}
}
