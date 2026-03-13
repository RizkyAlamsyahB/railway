package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type vendorOnboardingModel struct {
	ID                  string     `gorm:"column:id;primaryKey"`
	Email               string     `gorm:"column:email"`
	Status              string     `gorm:"column:status"`
	PasswordHash        *string    `gorm:"column:password_hash"`
	StoreName           *string    `gorm:"column:store_name"`
	VendorType          *string    `gorm:"column:vendor_type"`
	BusinessLegalType   *string    `gorm:"column:business_legal_type"`
	DocumentIDType      *string    `gorm:"column:document_id_type"`
	NIK                 *string    `gorm:"column:nik"`
	OwnerName           *string    `gorm:"column:owner_name"`
	BirthDate           *time.Time `gorm:"column:birth_date"`
	DocumentIDObjectKey *string    `gorm:"column:document_id_object_key"`
	OTPVerifiedAt       time.Time  `gorm:"column:otp_verified_at"`
	CompletedAt         *time.Time `gorm:"column:completed_at"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

func (vendorOnboardingModel) TableName() string { return "vendor_onboardings" }

type vendorOnboardingRepository struct {
	db *gorm.DB
}

// NewVendorOnboardingRepository creates a new VendorOnboardingRepository backed by GORM.
func NewVendorOnboardingRepository(db *gorm.DB) domain.VendorOnboardingRepository {
	return &vendorOnboardingRepository{db: db}
}

func (r *vendorOnboardingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
	var model vendorOnboardingModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendorOnboarding(&model), nil
}

func (r *vendorOnboardingRepository) FindByEmail(ctx context.Context, email string) (*domain.VendorOnboarding, error) {
	var model vendorOnboardingModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendorOnboarding(&model), nil
}

func (r *vendorOnboardingRepository) Upsert(ctx context.Context, onboarding *domain.VendorOnboarding) error {
	model := toVendorOnboardingModel(onboarding)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"id",
			"status",
			"password_hash",
			"store_name",
			"vendor_type",
			"business_legal_type",
			"document_id_type",
			"nik",
			"owner_name",
			"birth_date",
			"document_id_object_key",
			"otp_verified_at",
			"completed_at",
			"updated_at",
		}),
	}).Create(&model).Error
}

func (r *vendorOnboardingRepository) Update(ctx context.Context, onboarding *domain.VendorOnboarding, fields ...string) error {
	model := toVendorOnboardingModel(onboarding)
	return r.db.WithContext(ctx).
		Model(&vendorOnboardingModel{}).
		Where("id = ?", onboarding.ID.String()).
		Select(fields).
		Updates(&model).Error
}

func (r *vendorOnboardingRepository) MarkCompleted(ctx context.Context, id uuid.UUID, completedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&vendorOnboardingModel{}).
		Where("id = ?", id.String()).
		Updates(map[string]any{
			"status":       domain.VendorOnboardingStatusCompleted,
			"completed_at": completedAt,
			"updated_at":   completedAt,
		}).Error
}

func toVendorOnboardingModel(onboarding *domain.VendorOnboarding) vendorOnboardingModel {
	return vendorOnboardingModel{
		ID:                  onboarding.ID.String(),
		Email:               onboarding.Email,
		Status:              onboarding.Status,
		PasswordHash:        onboarding.PasswordHash,
		StoreName:           onboarding.StoreName,
		VendorType:          onboarding.VendorType,
		BusinessLegalType:   onboarding.BusinessLegalType,
		DocumentIDType:      onboarding.DocumentIDType,
		NIK:                 onboarding.NIK,
		OwnerName:           onboarding.OwnerName,
		BirthDate:           onboarding.BirthDate,
		DocumentIDObjectKey: onboarding.DocumentIDObjectKey,
		OTPVerifiedAt:       onboarding.OTPVerifiedAt,
		CompletedAt:         onboarding.CompletedAt,
		CreatedAt:           onboarding.CreatedAt,
		UpdatedAt:           onboarding.UpdatedAt,
	}
}

func toDomainVendorOnboarding(model *vendorOnboardingModel) *domain.VendorOnboarding {
	id, _ := uuid.Parse(model.ID)
	return &domain.VendorOnboarding{
		ID:                  id,
		Email:               model.Email,
		Status:              model.Status,
		PasswordHash:        model.PasswordHash,
		StoreName:           model.StoreName,
		VendorType:          model.VendorType,
		BusinessLegalType:   model.BusinessLegalType,
		DocumentIDType:      model.DocumentIDType,
		NIK:                 model.NIK,
		OwnerName:           model.OwnerName,
		BirthDate:           model.BirthDate,
		DocumentIDObjectKey: model.DocumentIDObjectKey,
		OTPVerifiedAt:       model.OTPVerifiedAt,
		CompletedAt:         model.CompletedAt,
		CreatedAt:           model.CreatedAt,
		UpdatedAt:           model.UpdatedAt,
	}
}
