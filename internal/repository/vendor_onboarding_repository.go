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
	ID                   string     `gorm:"column:id;primaryKey"`
	Email                string     `gorm:"column:email"`
	Status               string     `gorm:"column:status"`
	PasswordHash         *string    `gorm:"column:password_hash"`
	StoreName            *string    `gorm:"column:store_name"`
	VendorType           *string    `gorm:"column:vendor_type"`
	BusinessLegalType    *string    `gorm:"column:business_legal_type"`
	DocumentIDType       *string    `gorm:"column:document_id_type"`
	NIK                  *string    `gorm:"column:nik"`
	OwnerName            *string    `gorm:"column:owner_name"`
	BirthDate            *time.Time `gorm:"column:birth_date"`
	DocumentIDObjectKey  *string    `gorm:"column:document_id_object_key"`
	NIB                  *string    `gorm:"column:nib"`
	CompanyName          *string    `gorm:"column:company_name"`
	EstablishedDate      *time.Time `gorm:"column:established_date"`
	RegisteredAddress    *string    `gorm:"column:registered_address"`
	NIBDocumentObjectKey *string    `gorm:"column:nib_document_object_key"`
	OTPVerifiedAt        time.Time  `gorm:"column:otp_verified_at"`
	CompletedAt          *time.Time `gorm:"column:completed_at"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`
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
	return r.toDomainVendorOnboarding(&model)
}

func (r *vendorOnboardingRepository) FindByEmail(ctx context.Context, email string) (*domain.VendorOnboarding, error) {
	var model vendorOnboardingModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainVendorOnboarding(&model)
}

func (r *vendorOnboardingRepository) Upsert(ctx context.Context, onboarding *domain.VendorOnboarding) error {
	model, err := r.toVendorOnboardingModel(onboarding)
	if err != nil {
		return err
	}
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
			"nib",
			"company_name",
			"established_date",
			"registered_address",
			"nib_document_object_key",
			"otp_verified_at",
			"completed_at",
			"updated_at",
		}),
	}).Create(&model).Error
}

func (r *vendorOnboardingRepository) Update(ctx context.Context, onboarding *domain.VendorOnboarding, fields ...string) error {
	model, err := r.toVendorOnboardingModel(onboarding)
	if err != nil {
		return err
	}
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
			"status":                  domain.VendorOnboardingStatusCompleted,
			"password_hash":           nil,
			"store_name":              nil,
			"vendor_type":             nil,
			"business_legal_type":     nil,
			"document_id_type":        nil,
			"nik":                     nil,
			"owner_name":              nil,
			"birth_date":              nil,
			"document_id_object_key":  nil,
			"nib":                     nil,
			"company_name":            nil,
			"established_date":        nil,
			"registered_address":      nil,
			"nib_document_object_key": nil,
			"completed_at":            completedAt,
			"updated_at":              completedAt,
		}).Error
}

func (r *vendorOnboardingRepository) FinalizeRegistration(ctx context.Context, input domain.VendorRegistrationFinalizeInput) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Look up role (mirrors userRepository.Create logic)
		var role roleModel
		if err := tx.Where("code = ?", input.UserRole).First(&role).Error; err != nil {
			return fmt.Errorf("role '%s' not found: %w", input.UserRole, err)
		}

		// 2. Create user
		um := toUserModel(input.User)
		um.RoleID = role.ID
		if err := tx.Create(&um).Error; err != nil {
			return err
		}

		// 3. Create vendor
		vm := toVendorModel(input.Vendor)
		if err := tx.Create(&vm).Error; err != nil {
			return err
		}

		// 4. Create documents
		for i := range input.Documents {
			dm := toVendorDocumentModel(&input.Documents[i])
			if err := tx.Create(&dm).Error; err != nil {
				return err
			}
		}

		// 5. Create vendor balance
		balance := vendorBalanceModel{
			VendorID:  input.Vendor.ID.String(),
			UpdatedAt: input.CompletedAt,
		}
		if err := tx.Create(&balance).Error; err != nil {
			return err
		}

		// 6. Mark onboarding completed
		if err := tx.Model(&vendorOnboardingModel{}).
			Where("id = ?", input.OnboardingID.String()).
			Updates(map[string]any{
				"status":                  domain.VendorOnboardingStatusCompleted,
				"password_hash":           nil,
				"store_name":              nil,
				"vendor_type":             nil,
				"business_legal_type":     nil,
				"document_id_type":        nil,
				"nik":                     nil,
				"owner_name":              nil,
				"birth_date":              nil,
				"document_id_object_key":  nil,
				"nib":                     nil,
				"company_name":            nil,
				"established_date":        nil,
				"registered_address":      nil,
				"nib_document_object_key": nil,
				"completed_at":            input.CompletedAt,
				"updated_at":              input.CompletedAt,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *vendorOnboardingRepository) toVendorOnboardingModel(onboarding *domain.VendorOnboarding) (vendorOnboardingModel, error) {
	model := vendorOnboardingModel{
		ID:                   onboarding.ID.String(),
		Email:                onboarding.Email,
		Status:               onboarding.Status,
		PasswordHash:         onboarding.PasswordHash,
		StoreName:            onboarding.StoreName,
		VendorType:           onboarding.VendorType,
		BusinessLegalType:    onboarding.BusinessLegalType,
		DocumentIDType:       onboarding.DocumentIDType,
		NIK:                  onboarding.NIK,
		OwnerName:            onboarding.OwnerName,
		BirthDate:            onboarding.BirthDate,
		DocumentIDObjectKey:  onboarding.DocumentIDObjectKey,
		NIB:                  onboarding.NIB,
		CompanyName:          onboarding.CompanyName,
		EstablishedDate:      onboarding.EstablishedDate,
		RegisteredAddress:    onboarding.RegisteredAddress,
		NIBDocumentObjectKey: onboarding.NIBDocumentObjectKey,
		OTPVerifiedAt:        onboarding.OTPVerifiedAt,
		CompletedAt:          onboarding.CompletedAt,
		CreatedAt:            onboarding.CreatedAt,
		UpdatedAt:            onboarding.UpdatedAt,
	}

	if model.NIK != nil && *model.NIK != "" {
		encrypted, err := r.fieldCipher.EncryptString(*model.NIK)
		if err != nil {
			return vendorOnboardingModel{}, fmt.Errorf("failed to encrypt vendor onboarding nik: %w", err)
		}
		model.NIK = &encrypted
	}

	if model.NIB != nil && *model.NIB != "" {
		encrypted, err := r.fieldCipher.EncryptString(*model.NIB)
		if err != nil {
			return vendorOnboardingModel{}, fmt.Errorf("failed to encrypt vendor onboarding nib: %w", err)
		}
		model.NIB = &encrypted
	}

	return model, nil
}

func (r *vendorOnboardingRepository) toDomainVendorOnboarding(model *vendorOnboardingModel) (*domain.VendorOnboarding, error) {
	id, _ := uuid.Parse(model.ID)
	onboarding := &domain.VendorOnboarding{
		ID:                   id,
		Email:                model.Email,
		Status:               model.Status,
		PasswordHash:         model.PasswordHash,
		StoreName:            model.StoreName,
		VendorType:           model.VendorType,
		BusinessLegalType:    model.BusinessLegalType,
		DocumentIDType:       model.DocumentIDType,
		NIK:                  model.NIK,
		OwnerName:            model.OwnerName,
		BirthDate:            model.BirthDate,
		DocumentIDObjectKey:  model.DocumentIDObjectKey,
		NIB:                  model.NIB,
		CompanyName:          model.CompanyName,
		EstablishedDate:      model.EstablishedDate,
		RegisteredAddress:    model.RegisteredAddress,
		NIBDocumentObjectKey: model.NIBDocumentObjectKey,
		OTPVerifiedAt:        model.OTPVerifiedAt,
		CompletedAt:          model.CompletedAt,
		CreatedAt:            model.CreatedAt,
		UpdatedAt:            model.UpdatedAt,
	}

	if onboarding.NIK != nil && *onboarding.NIK != "" {
		decrypted, err := r.fieldCipher.DecryptString(*onboarding.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt vendor onboarding nik: %w", err)
		}
		onboarding.NIK = &decrypted
	}

	if onboarding.NIB != nil && *onboarding.NIB != "" {
		decrypted, err := r.fieldCipher.DecryptString(*onboarding.NIB)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt vendor onboarding nib: %w", err)
		}
		onboarding.NIB = &decrypted
	}

	return onboarding, nil
}
