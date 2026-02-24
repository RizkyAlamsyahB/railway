package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

// GORM model structs (internal to repository layer).

type vendorModel struct {
	ID                    string     `gorm:"column:id;primaryKey"`
	OwnerUserID           string     `gorm:"column:owner_user_id"`
	VendorType            string     `gorm:"column:vendor_type"`
	LegalName             *string    `gorm:"column:legal_name"`
	DisplayName           string     `gorm:"column:display_name"`
	ResponsiblePersonName string     `gorm:"column:responsible_person_name"`
	Description           *string    `gorm:"column:description"`
	Status                string     `gorm:"column:status"`
	ApprovedBy            *string    `gorm:"column:approved_by"`
	ApprovedAt            *time.Time `gorm:"column:approved_at"`
	StatusReason          *string    `gorm:"column:status_reason"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	UpdatedAt             time.Time  `gorm:"column:updated_at"`
}

func (vendorModel) TableName() string { return "vendors" }

type vendorBankAccountModel struct {
	ID                 string     `gorm:"column:id;primaryKey"`
	VendorID           string     `gorm:"column:vendor_id"`
	BankName           string     `gorm:"column:bank_name"`
	AccountNumber      string     `gorm:"column:account_number"`
	AccountHolderName  string     `gorm:"column:account_holder_name"`
	VerificationStatus string     `gorm:"column:verification_status"`
	RejectionReason    *string    `gorm:"column:rejection_reason"`
	VerifiedBy         *string    `gorm:"column:verified_by"`
	VerifiedAt         *time.Time `gorm:"column:verified_at"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
}

func (vendorBankAccountModel) TableName() string { return "vendor_bank_accounts" }

type vendorDocumentModel struct {
	ID                 string     `gorm:"column:id;primaryKey"`
	VendorID           string     `gorm:"column:vendor_id"`
	DocType            string     `gorm:"column:doc_type"`
	FileURL            string     `gorm:"column:file_url"`
	MimeType           *string    `gorm:"column:mime_type"`
	FileSizeBytes      *int       `gorm:"column:file_size_bytes"`
	FileChecksum       *string    `gorm:"column:file_checksum"`
	UploadedBy         *string    `gorm:"column:uploaded_by"`
	VerificationStatus string     `gorm:"column:verification_status"`
	RejectionReason    *string    `gorm:"column:rejection_reason"`
	VerifiedBy         *string    `gorm:"column:verified_by"`
	VerifiedAt         *time.Time `gorm:"column:verified_at"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
}

func (vendorDocumentModel) TableName() string { return "vendor_documents" }

type vendorRepository struct {
	db *gorm.DB
}

// NewVendorRepository creates a new VendorRepository backed by GORM.
func NewVendorRepository(db *gorm.DB) domain.VendorRepository {
	return &vendorRepository{db: db}
}

func (r *vendorRepository) Create(ctx context.Context, vendor *domain.Vendor, bankAccount *domain.VendorBankAccount, documents []domain.VendorDocument) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		vm := toVendorModel(vendor)
		if err := tx.Create(&vm).Error; err != nil {
			return err
		}

		bam := toVendorBankAccountModel(bankAccount)
		if err := tx.Create(&bam).Error; err != nil {
			return err
		}

		for i := range documents {
			dm := toVendorDocumentModel(&documents[i])
			if err := tx.Create(&dm).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *vendorRepository) FindByOwnerUserID(ctx context.Context, userID uuid.UUID) (*domain.Vendor, error) {
	var model vendorModel
	if err := r.db.WithContext(ctx).Where("owner_user_id = ?", userID.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendor(&model), nil
}

func (r *vendorRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Vendor, error) {
	var model vendorModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendor(&model), nil
}

func (r *vendorRepository) List(ctx context.Context, params domain.VendorListParams) ([]domain.Vendor, int64, error) {
	query := r.db.WithContext(ctx).Model(&vendorModel{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.VendorType != "" {
		query = query.Where("vendor_type = ?", params.VendorType)
	}
	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("(LOWER(display_name) LIKE ? OR LOWER(legal_name) LIKE ?)", search, search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []vendorModel
	offset := (params.Page - 1) * params.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.Limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	vendors := make([]domain.Vendor, len(models))
	for i, m := range models {
		vendors[i] = *toDomainVendor(&m)
	}
	return vendors, total, nil
}

func (r *vendorRepository) FindBankAccountByVendorID(ctx context.Context, vendorID uuid.UUID) (*domain.VendorBankAccount, error) {
	var model vendorBankAccountModel
	if err := r.db.WithContext(ctx).Where("vendor_id = ?", vendorID.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendorBankAccount(&model), nil
}

func (r *vendorRepository) FindDocumentsByVendorID(ctx context.Context, vendorID uuid.UUID) ([]domain.VendorDocument, error) {
	var models []vendorDocumentModel
	if err := r.db.WithContext(ctx).Where("vendor_id = ?", vendorID.String()).Find(&models).Error; err != nil {
		return nil, err
	}

	docs := make([]domain.VendorDocument, len(models))
	for i, m := range models {
		docs[i] = *toDomainVendorDocument(&m)
	}
	return docs, nil
}

func (r *vendorRepository) ConfirmDocumentsAndUpdateStatus(ctx context.Context, vendorID uuid.UUID, documents []domain.VendorDocument, newStatus string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range documents {
			dm := toVendorDocumentModel(&documents[i])
			if err := tx.Model(&vendorDocumentModel{ID: dm.ID}).
				Select("file_url", "mime_type", "file_size_bytes", "uploaded_by", "updated_at").
				Updates(&dm).Error; err != nil {
				return err
			}
		}

		if newStatus != "" {
			if err := tx.Model(&vendorModel{}).
				Where("id = ?", vendorID.String()).
				Updates(map[string]interface{}{
					"status":     newStatus,
					"updated_at": time.Now(),
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *vendorRepository) UpdateStatus(ctx context.Context, vendorID uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&vendorModel{}).
		Where("id = ?", vendorID.String()).
		Updates(updates).Error
}

// Mapper helpers.

func toVendorModel(v *domain.Vendor) vendorModel {
	m := vendorModel{
		ID:                    v.ID.String(),
		OwnerUserID:           v.OwnerUserID.String(),
		VendorType:            v.VendorType,
		LegalName:             v.LegalName,
		DisplayName:           v.DisplayName,
		ResponsiblePersonName: v.ResponsiblePersonName,
		Description:           v.Description,
		Status:                v.Status,
		StatusReason:          v.StatusReason,
		CreatedAt:             v.CreatedAt,
		UpdatedAt:             v.UpdatedAt,
	}
	if v.ApprovedBy != nil {
		s := v.ApprovedBy.String()
		m.ApprovedBy = &s
	}
	m.ApprovedAt = v.ApprovedAt
	return m
}

func toDomainVendor(m *vendorModel) *domain.Vendor {
	id, _ := uuid.Parse(m.ID)
	ownerID, _ := uuid.Parse(m.OwnerUserID)

	v := &domain.Vendor{
		ID:                    id,
		OwnerUserID:           ownerID,
		VendorType:            m.VendorType,
		LegalName:             m.LegalName,
		DisplayName:           m.DisplayName,
		ResponsiblePersonName: m.ResponsiblePersonName,
		Description:           m.Description,
		Status:                m.Status,
		StatusReason:          m.StatusReason,
		ApprovedAt:            m.ApprovedAt,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
	if m.ApprovedBy != nil {
		approvedBy, _ := uuid.Parse(*m.ApprovedBy)
		v.ApprovedBy = &approvedBy
	}
	return v
}

func toVendorBankAccountModel(ba *domain.VendorBankAccount) vendorBankAccountModel {
	m := vendorBankAccountModel{
		ID:                 ba.ID.String(),
		VendorID:           ba.VendorID.String(),
		BankName:           ba.BankName,
		AccountNumber:      ba.AccountNumber,
		AccountHolderName:  ba.AccountHolderName,
		VerificationStatus: ba.VerificationStatus,
		RejectionReason:    ba.RejectionReason,
		CreatedAt:          ba.CreatedAt,
		UpdatedAt:          ba.UpdatedAt,
	}
	if ba.VerifiedBy != nil {
		s := ba.VerifiedBy.String()
		m.VerifiedBy = &s
	}
	m.VerifiedAt = ba.VerifiedAt
	return m
}

func toDomainVendorBankAccount(m *vendorBankAccountModel) *domain.VendorBankAccount {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)

	ba := &domain.VendorBankAccount{
		ID:                 id,
		VendorID:           vendorID,
		BankName:           m.BankName,
		AccountNumber:      m.AccountNumber,
		AccountHolderName:  m.AccountHolderName,
		VerificationStatus: m.VerificationStatus,
		RejectionReason:    m.RejectionReason,
		VerifiedAt:         m.VerifiedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
	if m.VerifiedBy != nil {
		verifiedBy, _ := uuid.Parse(*m.VerifiedBy)
		ba.VerifiedBy = &verifiedBy
	}
	return ba
}

func toVendorDocumentModel(d *domain.VendorDocument) vendorDocumentModel {
	m := vendorDocumentModel{
		ID:                 d.ID.String(),
		VendorID:           d.VendorID.String(),
		DocType:            d.DocType,
		FileURL:            d.FileURL,
		MimeType:           d.MimeType,
		FileSizeBytes:      d.FileSizeBytes,
		FileChecksum:       d.FileChecksum,
		VerificationStatus: d.VerificationStatus,
		RejectionReason:    d.RejectionReason,
		VerifiedAt:         d.VerifiedAt,
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
	if d.UploadedBy != nil {
		s := d.UploadedBy.String()
		m.UploadedBy = &s
	}
	if d.VerifiedBy != nil {
		s := d.VerifiedBy.String()
		m.VerifiedBy = &s
	}
	return m
}

func toDomainVendorDocument(m *vendorDocumentModel) *domain.VendorDocument {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)

	d := &domain.VendorDocument{
		ID:                 id,
		VendorID:           vendorID,
		DocType:            m.DocType,
		FileURL:            m.FileURL,
		MimeType:           m.MimeType,
		FileSizeBytes:      m.FileSizeBytes,
		FileChecksum:       m.FileChecksum,
		VerificationStatus: m.VerificationStatus,
		RejectionReason:    m.RejectionReason,
		VerifiedAt:         m.VerifiedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
	if m.UploadedBy != nil {
		uploadedBy, _ := uuid.Parse(*m.UploadedBy)
		d.UploadedBy = &uploadedBy
	}
	if m.VerifiedBy != nil {
		verifiedBy, _ := uuid.Parse(*m.VerifiedBy)
		d.VerifiedBy = &verifiedBy
	}
	return d
}
