package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type vendorOnboardingModel struct {
	ID            string     `gorm:"column:id;primaryKey"`
	Email         string     `gorm:"column:email"`
	Status        string     `gorm:"column:status"`
	PasswordHash  *string    `gorm:"column:password_hash"`
	OTPVerifiedAt time.Time  `gorm:"column:otp_verified_at"`
	CompletedAt   *time.Time `gorm:"column:completed_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (vendorOnboardingModel) TableName() string { return "vendor_onboardings" }

type vendorOnboardingRepository struct {
	db          *gorm.DB
	fieldCipher *sensitivedata.FieldCipher
}

// NewVendorOnboardingRepository creates a new VendorOnboardingRepository backed by GORM.
func NewVendorOnboardingRepository(db *gorm.DB, fieldCipher *sensitivedata.FieldCipher) domain.VendorOnboardingRepository {
	return &vendorOnboardingRepository{db: db, fieldCipher: fieldCipher}
}

func (r *vendorOnboardingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.VendorOnboarding, error) {
	var model vendorOnboardingModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainVendorOnboarding(&model), nil
}

func (r *vendorOnboardingRepository) FindByEmail(ctx context.Context, email string) (*domain.VendorOnboarding, error) {
	var model vendorOnboardingModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainVendorOnboarding(&model), nil
}

func (r *vendorOnboardingRepository) Upsert(ctx context.Context, onboarding *domain.VendorOnboarding) error {
	model := r.toVendorOnboardingModel(onboarding)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"id",
			"status",
			"password_hash",
			"otp_verified_at",
			"completed_at",
			"updated_at",
		}),
	}).Create(&model).Error
}

func (r *vendorOnboardingRepository) Update(ctx context.Context, onboarding *domain.VendorOnboarding, fields ...string) error {
	model := r.toVendorOnboardingModel(onboarding)
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
			"status":        domain.VendorOnboardingStatusCompleted,
			"password_hash": nil,
			"completed_at":  completedAt,
			"updated_at":    completedAt,
		}).Error
}

func (r *vendorOnboardingRepository) FinalizeRegistration(ctx context.Context, input domain.VendorRegistrationFinalizeInput) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role roleModel
		if err := tx.Where("code = ?", input.UserRole).First(&role).Error; err != nil {
			return fmt.Errorf("role '%s' not found: %w", input.UserRole, err)
		}

		um := toUserModel(input.User)
		um.RoleID = role.ID
		if err := tx.Create(&um).Error; err != nil {
			return err
		}

		vm := toVendorModel(input.Vendor)
		if err := tx.Create(&vm).Error; err != nil {
			return err
		}

		balance := vendorBalanceModel{
			VendorID:  input.Vendor.ID.String(),
			UpdatedAt: input.CompletedAt,
		}
		if err := tx.Create(&balance).Error; err != nil {
			return err
		}

		if err := tx.Model(&vendorOnboardingModel{}).
			Where("id = ?", input.OnboardingID.String()).
			Updates(map[string]any{
				"status":        domain.VendorOnboardingStatusCompleted,
				"password_hash": nil,
				"completed_at":  input.CompletedAt,
				"updated_at":    input.CompletedAt,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *vendorOnboardingRepository) toVendorOnboardingModel(onboarding *domain.VendorOnboarding) vendorOnboardingModel {
	return vendorOnboardingModel{
		ID:            onboarding.ID.String(),
		Email:         onboarding.Email,
		Status:        onboarding.Status,
		PasswordHash:  onboarding.PasswordHash,
		OTPVerifiedAt: onboarding.OTPVerifiedAt,
		CompletedAt:   onboarding.CompletedAt,
		CreatedAt:     onboarding.CreatedAt,
		UpdatedAt:     onboarding.UpdatedAt,
	}
}

func (r *vendorOnboardingRepository) toDomainVendorOnboarding(model *vendorOnboardingModel) *domain.VendorOnboarding {
	id, _ := uuid.Parse(model.ID)
	return &domain.VendorOnboarding{
		ID:            id,
		Email:         model.Email,
		Status:        model.Status,
		PasswordHash:  model.PasswordHash,
		OTPVerifiedAt: model.OTPVerifiedAt,
		CompletedAt:   model.CompletedAt,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
