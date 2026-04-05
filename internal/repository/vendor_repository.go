package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/sensitivedata"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GORM model structs (internal to repository layer).

type vendorModel struct {
	ID              string     `gorm:"column:id;primaryKey"`
	OwnerUserID     string     `gorm:"column:owner_user_id"`
	VendorType      *string    `gorm:"column:vendor_type"`
	LegalName       *string    `gorm:"column:legal_name"`
	DisplayName     *string    `gorm:"column:display_name"`
	Description     *string    `gorm:"column:description"`
	ProvinceID      *string    `gorm:"column:province_id"`
	ProvinceName    *string    `gorm:"column:province_name"`
	CityID          *string    `gorm:"column:city_id"`
	CityName        *string    `gorm:"column:city_name"`
	DistrictID      *string    `gorm:"column:district_id"`
	DistrictName    *string    `gorm:"column:district_name"`
	SubdistrictID   *string    `gorm:"column:subdistrict_id"`
	SubdistrictName *string    `gorm:"column:subdistrict_name"`
	PostalCode      *string    `gorm:"column:postal_code"`
	AddressLine     *string    `gorm:"column:address_line"`
	Status          string     `gorm:"column:status"`
	ApprovedBy      *string    `gorm:"column:approved_by"`
	ApprovedAt      *time.Time `gorm:"column:approved_at"`
	StatusReason    *string    `gorm:"column:status_reason"`
	XenditAccountID *string    `gorm:"column:xendit_account_id"`
	TotalSold       int        `gorm:"column:total_sold"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
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
	ID            string     `gorm:"column:id;primaryKey"`
	VendorID      string     `gorm:"column:vendor_id"`
	DocType       string     `gorm:"column:doc_type"`
	FileURL       string     `gorm:"column:file_url"`
	MimeType      *string    `gorm:"column:mime_type"`
	FileSizeBytes *int       `gorm:"column:file_size_bytes"`
	FileChecksum  *string    `gorm:"column:file_checksum"`
	UploadedBy    *string    `gorm:"column:uploaded_by"`
	VerifiedBy    *string    `gorm:"column:verified_by"`
	VerifiedAt    *time.Time `gorm:"column:verified_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (vendorDocumentModel) TableName() string { return "vendor_documents" }

type vendorResponsiblePersonModel struct {
	ID        string    `gorm:"column:id;primaryKey"`
	VendorID  string    `gorm:"column:vendor_id"`
	UserID    string    `gorm:"column:user_id"`
	NIK       string    `gorm:"column:nik"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (vendorResponsiblePersonModel) TableName() string { return "vendor_responsible_persons" }

type vendorRepository struct {
	db          *gorm.DB
	fieldCipher *sensitivedata.FieldCipher
}

// NewVendorRepository creates a new VendorRepository backed by GORM.
func NewVendorRepository(db *gorm.DB, fieldCipher *sensitivedata.FieldCipher) domain.VendorRepository {
	return &vendorRepository{db: db, fieldCipher: fieldCipher}
}

func (r *vendorRepository) Create(ctx context.Context, vendor *domain.Vendor, bankAccount *domain.VendorBankAccount, documents []domain.VendorDocument) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		vm := toVendorModel(vendor)
		if err := tx.Create(&vm).Error; err != nil {
			return err
		}

		bam, err := r.toVendorBankAccountModel(bankAccount)
		if err != nil {
			return err
		}
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

func (r *vendorRepository) CreateMinimal(ctx context.Context, vendor *domain.Vendor, documents []domain.VendorDocument) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		vm := toVendorModel(vendor)
		if err := tx.Create(&vm).Error; err != nil {
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
	return r.toDomainVendorBankAccount(&model)
}

func (r *vendorRepository) FindResponsiblePersonByVendorID(ctx context.Context, vendorID uuid.UUID) (*domain.VendorResponsiblePerson, error) {
	var model vendorResponsiblePersonModel
	if err := r.db.WithContext(ctx).Where("vendor_id = ?", vendorID.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainVendorResponsiblePerson(&model)
}

func (r *vendorRepository) UpsertBankAccount(ctx context.Context, bankAccount *domain.VendorBankAccount) error {
	model, err := r.toVendorBankAccountModel(bankAccount)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "vendor_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"bank_name":           model.BankName,
			"account_number":      model.AccountNumber,
			"account_holder_name": model.AccountHolderName,
			"verification_status": model.VerificationStatus,
			"rejection_reason":    model.RejectionReason,
			"verified_by":         model.VerifiedBy,
			"verified_at":         model.VerifiedAt,
			"updated_at":          model.UpdatedAt,
		}),
	}).Create(&model).Error
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
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "vendor_id"}, {Name: "doc_type"}},
				DoUpdates: clause.Assignments(map[string]any{
					"file_url":        dm.FileURL,
					"mime_type":       dm.MimeType,
					"file_size_bytes": dm.FileSizeBytes,
					"uploaded_by":     dm.UploadedBy,
					"updated_at":      dm.UpdatedAt,
				}),
			}).Create(&dm).Error; err != nil {
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

func (r *vendorRepository) SubmitSouvenirStoreProposal(ctx context.Context, input domain.SubmitSouvenirStoreProposalInput) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.OwnerUser == nil {
			return errors.New("missing owner user")
		}
		if input.Vendor == nil {
			return errors.New("missing vendor")
		}
		if input.ResponsiblePerson == nil {
			return errors.New("missing responsible person")
		}

		userUpdates := map[string]any{
			"full_name":  input.OwnerUser.FullName,
			"phone":      input.OwnerUser.Phone,
			"updated_at": input.OwnerUser.UpdatedAt,
		}
		if err := tx.Model(&userModel{}).
			Where("id = ?", input.OwnerUser.ID.String()).
			Updates(userUpdates).Error; err != nil {
			return err
		}

		vendorUpdates := map[string]any{
			"vendor_type":    input.Vendor.VendorType,
			"display_name":   input.Vendor.DisplayName,
			"description":    input.Vendor.Description,
			"province_id":    input.Vendor.ProvinceID,
			"city_id":        input.Vendor.CityID,
			"district_id":    input.Vendor.DistrictID,
			"subdistrict_id": input.Vendor.SubdistrictID,
			"postal_code":    input.Vendor.PostalCode,
			"address_line":   input.Vendor.AddressLine,
			"status":         input.Vendor.Status,
			"approved_by":    nil,
			"approved_at":    nil,
			"status_reason":  nil,
			"updated_at":     input.Vendor.UpdatedAt,
		}
		if err := tx.Model(&vendorModel{}).
			Where("id = ?", input.Vendor.ID.String()).
			Updates(vendorUpdates).Error; err != nil {
			return err
		}

		responsiblePersonModel, err := r.toVendorResponsiblePersonModel(input.ResponsiblePerson)
		if err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "vendor_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"user_id":    responsiblePersonModel.UserID,
				"nik":        responsiblePersonModel.NIK,
				"updated_at": responsiblePersonModel.UpdatedAt,
			}),
		}).Create(&responsiblePersonModel).Error; err != nil {
			return err
		}

		for i := range input.Documents {
			documentModel := toVendorDocumentModel(&input.Documents[i])
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "vendor_id"}, {Name: "doc_type"}},
				DoUpdates: clause.Assignments(map[string]any{
					"file_url":        documentModel.FileURL,
					"mime_type":       documentModel.MimeType,
					"file_size_bytes": documentModel.FileSizeBytes,
					"file_checksum":   documentModel.FileChecksum,
					"uploaded_by":     documentModel.UploadedBy,
					"verified_by":     nil,
					"verified_at":     nil,
					"updated_at":      documentModel.UpdatedAt,
				}),
			}).Create(&documentModel).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

type payoutBatchModel struct {
	ID             string     `gorm:"column:id;primaryKey"`
	VendorID       string     `gorm:"column:vendor_id"`
	PeriodStart    time.Time  `gorm:"column:period_start"`
	PeriodEnd      time.Time  `gorm:"column:period_end"`
	Status         string     `gorm:"column:status"`
	TotalGross     float64    `gorm:"column:total_gross"`
	TotalFee       float64    `gorm:"column:total_fee"`
	TotalNet       float64    `gorm:"column:total_net"`
	PaidAt         *time.Time `gorm:"column:paid_at"`
	CreatedBy      string     `gorm:"column:created_by"`
	XenditPayoutID *string    `gorm:"column:xendit_payout_id"`
	ChannelCode    *string    `gorm:"column:channel_code"`
	Description    *string    `gorm:"column:description"`
	XenditStatus   *string    `gorm:"column:xendit_status"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

func (payoutBatchModel) TableName() string { return "payout_batches" }

func (r *vendorRepository) FindPayoutBatchByID(ctx context.Context, id uuid.UUID) (*domain.PayoutBatch, error) {
	var model payoutBatchModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainPayoutBatch(&model), nil
}

func (r *vendorRepository) UpdatePayoutBatch(ctx context.Context, batch *domain.PayoutBatch) error {
	return r.db.WithContext(ctx).Model(&payoutBatchModel{}).
		Where("id = ?", batch.ID.String()).
		Updates(map[string]interface{}{
			"status":           batch.Status,
			"xendit_payout_id": batch.XenditPayoutID,
			"channel_code":     batch.ChannelCode,
			"description":      batch.Description,
			"xendit_status":    batch.XenditStatus,
			"updated_at":       time.Now(),
		}).Error
}

func toDomainPayoutBatch(m *payoutBatchModel) *domain.PayoutBatch {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)
	createdBy, _ := uuid.Parse(m.CreatedBy)

	return &domain.PayoutBatch{
		ID:             id,
		VendorID:       vendorID,
		PeriodStart:    m.PeriodStart,
		PeriodEnd:      m.PeriodEnd,
		Status:         m.Status,
		TotalGross:     m.TotalGross,
		TotalFee:       m.TotalFee,
		TotalNet:       m.TotalNet,
		PaidAt:         m.PaidAt,
		CreatedBy:      createdBy,
		XenditPayoutID: m.XenditPayoutID,
		ChannelCode:    m.ChannelCode,
		Description:    m.Description,
		XenditStatus:   m.XenditStatus,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// --- Balance & Withdrawal GORM models ---

type vendorBalanceModel struct {
	VendorID         string    `gorm:"column:vendor_id;primaryKey"`
	AvailableBalance float64   `gorm:"column:available_balance"`
	PendingBalance   float64   `gorm:"column:pending_balance"`
	EscrowBalance    float64   `gorm:"column:escrow_balance"`
	TotalEarned      float64   `gorm:"column:total_earned"`
	TotalWithdrawn   float64   `gorm:"column:total_withdrawn"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (vendorBalanceModel) TableName() string { return "vendor_balances" }

type vendorWithdrawalModel struct {
	ID                  string    `gorm:"column:id;primaryKey"`
	VendorID            string    `gorm:"column:vendor_id"`
	Amount              float64   `gorm:"column:amount"`
	ChannelCode         string    `gorm:"column:channel_code"`
	Status              string    `gorm:"column:status"`
	FeeEstimated        float64   `gorm:"column:fee_estimated"`
	FeeActual           *float64  `gorm:"column:fee_actual"`
	AmountNet           float64   `gorm:"column:amount_net"`
	TotalDeducted       float64   `gorm:"column:total_deducted"`
	FeeStatus           string    `gorm:"column:fee_status"`
	XenditPayoutID      *string   `gorm:"column:xendit_payout_id"`
	XenditStatus        *string   `gorm:"column:xendit_status"`
	XenditTransactionID *string   `gorm:"column:xendit_transaction_id"`
	Description         *string   `gorm:"column:description"`
	FailedReason        *string   `gorm:"column:failed_reason"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
}

func (vendorWithdrawalModel) TableName() string { return "vendor_withdrawals" }

func (r *vendorRepository) InitBalance(ctx context.Context, vendorID uuid.UUID) error {
	m := vendorBalanceModel{
		VendorID:         vendorID.String(),
		AvailableBalance: 0,
		PendingBalance:   0,
		EscrowBalance:    0,
		TotalEarned:      0,
		TotalWithdrawn:   0,
		UpdatedAt:        time.Now(),
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *vendorRepository) GetBalance(ctx context.Context, vendorID uuid.UUID) (*domain.VendorBalance, error) {
	var m vendorBalanceModel
	if err := r.db.WithContext(ctx).Where("vendor_id = ?", vendorID.String()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendorBalance(&m), nil
}

func (r *vendorRepository) CreditBalance(ctx context.Context, vendorID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&vendorBalanceModel{}).
		Where("vendor_id = ?", vendorID.String()).
		Updates(map[string]interface{}{
			"available_balance": gorm.Expr("available_balance + ?", amount),
			"total_earned":      gorm.Expr("total_earned + ?", amount),
			"updated_at":        time.Now(),
		}).Error
}

func (r *vendorRepository) DebitBalance(ctx context.Context, vendorID uuid.UUID, amount float64) error {
	result := r.db.WithContext(ctx).Model(&vendorBalanceModel{}).
		Where("vendor_id = ? AND available_balance >= ?", vendorID.String(), amount).
		Updates(map[string]interface{}{
			"available_balance": gorm.Expr("available_balance - ?", amount),
			"pending_balance":   gorm.Expr("pending_balance + ?", amount),
			"updated_at":        time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient balance or vendor not found")
	}
	return nil
}

func (r *vendorRepository) CompleteWithdrawal(ctx context.Context, vendorID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&vendorBalanceModel{}).
		Where("vendor_id = ?", vendorID.String()).
		Updates(map[string]interface{}{
			"pending_balance": gorm.Expr("pending_balance - ?", amount),
			"total_withdrawn": gorm.Expr("total_withdrawn + ?", amount),
			"updated_at":      time.Now(),
		}).Error
}

func (r *vendorRepository) FailWithdrawal(ctx context.Context, vendorID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).Model(&vendorBalanceModel{}).
		Where("vendor_id = ?", vendorID.String()).
		Updates(map[string]interface{}{
			"pending_balance":   gorm.Expr("pending_balance - ?", amount),
			"available_balance": gorm.Expr("available_balance + ?", amount),
			"updated_at":        time.Now(),
		}).Error
}

func (r *vendorRepository) CreateWithdrawal(ctx context.Context, withdrawal *domain.VendorWithdrawal) error {
	m := toVendorWithdrawalModel(withdrawal)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *vendorRepository) FindWithdrawalByID(ctx context.Context, id uuid.UUID) (*domain.VendorWithdrawal, error) {
	var m vendorWithdrawalModel
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendorWithdrawal(&m), nil
}

func (r *vendorRepository) UpdateWithdrawal(ctx context.Context, withdrawal *domain.VendorWithdrawal) error {
	return r.db.WithContext(ctx).Model(&vendorWithdrawalModel{}).
		Where("id = ?", withdrawal.ID.String()).
		Updates(map[string]interface{}{
			"status":                withdrawal.Status,
			"fee_estimated":         withdrawal.FeeEstimated,
			"fee_actual":            withdrawal.FeeActual,
			"amount_net":            withdrawal.AmountNet,
			"total_deducted":        withdrawal.TotalDeducted,
			"fee_status":            withdrawal.FeeStatus,
			"xendit_payout_id":      withdrawal.XenditPayoutID,
			"xendit_status":         withdrawal.XenditStatus,
			"xendit_transaction_id": withdrawal.XenditTransactionID,
			"description":           withdrawal.Description,
			"failed_reason":         withdrawal.FailedReason,
			"updated_at":            time.Now(),
		}).Error
}

func (r *vendorRepository) ApplyWithdrawalWebhookUpdate(
	ctx context.Context,
	withdrawalID uuid.UUID,
	newStatus string,
	xenditStatus string,
	xenditPayoutID *string,
	failedReason *string,
	balanceAction string,
	completion *domain.WithdrawalCompletionData,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wm vendorWithdrawalModel
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", withdrawalID.String()).
			First(&wm).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"updated_at": time.Now(),
		}
		if xenditStatus != "" {
			updates["xendit_status"] = xenditStatus
		}
		if newStatus != "" {
			updates["status"] = newStatus
		}
		if xenditPayoutID != nil && *xenditPayoutID != "" {
			updates["xendit_payout_id"] = *xenditPayoutID
		}
		if failedReason != nil && *failedReason != "" {
			updates["failed_reason"] = *failedReason
		}
		if completion != nil {
			updates["fee_actual"] = completion.FeeActual
			updates["amount_net"] = completion.NetAmount
			updates["total_deducted"] = completion.TotalDeducted
			updates["fee_status"] = domain.WithdrawalFeeStatusResolved
			if completion.XenditTransactionID != nil && *completion.XenditTransactionID != "" {
				updates["xendit_transaction_id"] = *completion.XenditTransactionID
			}
		} else if newStatus == domain.WithdrawalStatusFailed {
			updates["fee_status"] = domain.WithdrawalFeeStatusResolved
		}

		if err := tx.Model(&vendorWithdrawalModel{}).
			Where("id = ?", withdrawalID.String()).
			Updates(updates).Error; err != nil {
			return err
		}

		if balanceAction == domain.WithdrawalBalanceActionNone {
			return nil
		}

		var balUpdates map[string]interface{}
		balanceQuery := tx.Model(&vendorBalanceModel{}).Where("vendor_id = ?", wm.VendorID)

		switch balanceAction {
		case domain.WithdrawalBalanceActionComplete:
			totalDeducted := wm.TotalDeducted
			if completion != nil {
				totalDeducted = completion.TotalDeducted
			}
			if totalDeducted <= 0 {
				totalDeducted = wm.Amount
			}

			diff := totalDeducted - wm.Amount
			balanceQuery = balanceQuery.Where("pending_balance >= ?", wm.Amount)
			if diff > 0 {
				balanceQuery = balanceQuery.Where("available_balance >= ?", diff)
			}

			balUpdates = map[string]interface{}{
				"pending_balance": gorm.Expr("pending_balance - ?", wm.Amount),
				"total_withdrawn": gorm.Expr("total_withdrawn + ?", totalDeducted),
				"updated_at":      time.Now(),
			}
			if diff > 0 {
				balUpdates["available_balance"] = gorm.Expr("available_balance - ?", diff)
			} else if diff < 0 {
				balUpdates["available_balance"] = gorm.Expr("available_balance + ?", -diff)
			}
		case domain.WithdrawalBalanceActionFail:
			balanceQuery = balanceQuery.Where("pending_balance >= ?", wm.Amount)
			balUpdates = map[string]interface{}{
				"pending_balance":   gorm.Expr("pending_balance - ?", wm.Amount),
				"available_balance": gorm.Expr("available_balance + ?", wm.Amount),
				"updated_at":        time.Now(),
			}
		case domain.WithdrawalBalanceActionReverseCompleted:
			totalDeducted := wm.TotalDeducted
			if totalDeducted <= 0 {
				totalDeducted = wm.Amount
			}
			balanceQuery = balanceQuery.Where("total_withdrawn >= ?", totalDeducted)
			balUpdates = map[string]interface{}{
				"total_withdrawn":   gorm.Expr("total_withdrawn - ?", totalDeducted),
				"available_balance": gorm.Expr("available_balance + ?", totalDeducted),
				"updated_at":        time.Now(),
			}
		default:
			return errors.New("invalid withdrawal balance action")
		}

		res := balanceQuery.Updates(balUpdates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("failed to apply withdrawal balance update")
		}
		return nil
	})
}

// --- Balance & Withdrawal mappers ---

func toDomainVendorBalance(m *vendorBalanceModel) *domain.VendorBalance {
	vendorID, _ := uuid.Parse(m.VendorID)
	return &domain.VendorBalance{
		VendorID:         vendorID,
		AvailableBalance: m.AvailableBalance,
		PendingBalance:   m.PendingBalance,
		EscrowBalance:    m.EscrowBalance,
		TotalEarned:      m.TotalEarned,
		TotalWithdrawn:   m.TotalWithdrawn,
		UpdatedAt:        m.UpdatedAt,
	}
}

func toDomainVendorWithdrawal(m *vendorWithdrawalModel) *domain.VendorWithdrawal {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)
	return &domain.VendorWithdrawal{
		ID:                  id,
		VendorID:            vendorID,
		Amount:              m.Amount,
		ChannelCode:         m.ChannelCode,
		Status:              m.Status,
		FeeEstimated:        m.FeeEstimated,
		FeeActual:           m.FeeActual,
		AmountNet:           m.AmountNet,
		TotalDeducted:       m.TotalDeducted,
		FeeStatus:           m.FeeStatus,
		XenditPayoutID:      m.XenditPayoutID,
		XenditStatus:        m.XenditStatus,
		XenditTransactionID: m.XenditTransactionID,
		Description:         m.Description,
		FailedReason:        m.FailedReason,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func toVendorWithdrawalModel(w *domain.VendorWithdrawal) vendorWithdrawalModel {
	return vendorWithdrawalModel{
		ID:                  w.ID.String(),
		VendorID:            w.VendorID.String(),
		Amount:              w.Amount,
		ChannelCode:         w.ChannelCode,
		Status:              w.Status,
		FeeEstimated:        w.FeeEstimated,
		FeeActual:           w.FeeActual,
		AmountNet:           w.AmountNet,
		TotalDeducted:       w.TotalDeducted,
		FeeStatus:           w.FeeStatus,
		XenditPayoutID:      w.XenditPayoutID,
		XenditStatus:        w.XenditStatus,
		XenditTransactionID: w.XenditTransactionID,
		Description:         w.Description,
		FailedReason:        w.FailedReason,
		CreatedAt:           w.CreatedAt,
		UpdatedAt:           w.UpdatedAt,
	}
}

// Mapper helpers.

func toVendorModel(v *domain.Vendor) vendorModel {
	m := vendorModel{
		ID:              v.ID.String(),
		OwnerUserID:     v.OwnerUserID.String(),
		VendorType:      v.VendorType,
		LegalName:       v.LegalName,
		DisplayName:     v.DisplayName,
		Description:     v.Description,
		ProvinceID:      v.ProvinceID,
		ProvinceName:    v.ProvinceName,
		CityID:          v.CityID,
		CityName:        v.CityName,
		DistrictID:      v.DistrictID,
		DistrictName:    v.DistrictName,
		SubdistrictID:   v.SubdistrictID,
		SubdistrictName: v.SubdistrictName,
		PostalCode:      v.PostalCode,
		AddressLine:     v.AddressLine,
		Status:          v.Status,
		StatusReason:    v.StatusReason,
		XenditAccountID: v.XenditAccountID,
		TotalSold:       v.TotalSold,
		CreatedAt:       v.CreatedAt,
		UpdatedAt:       v.UpdatedAt,
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
		ID:              id,
		OwnerUserID:     ownerID,
		VendorType:      m.VendorType,
		LegalName:       m.LegalName,
		DisplayName:     m.DisplayName,
		Description:     m.Description,
		ProvinceID:      m.ProvinceID,
		ProvinceName:    m.ProvinceName,
		CityID:          m.CityID,
		CityName:        m.CityName,
		DistrictID:      m.DistrictID,
		DistrictName:    m.DistrictName,
		SubdistrictID:   m.SubdistrictID,
		SubdistrictName: m.SubdistrictName,
		PostalCode:      m.PostalCode,
		AddressLine:     m.AddressLine,
		Status:          m.Status,
		StatusReason:    m.StatusReason,
		XenditAccountID: m.XenditAccountID,
		TotalSold:       m.TotalSold,
		ApprovedAt:      m.ApprovedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
	if m.ApprovedBy != nil {
		approvedBy, _ := uuid.Parse(*m.ApprovedBy)
		v.ApprovedBy = &approvedBy
	}
	return v
}

func (r *vendorRepository) toVendorBankAccountModel(ba *domain.VendorBankAccount) (vendorBankAccountModel, error) {
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

	if m.AccountNumber != "" {
		encrypted, err := r.fieldCipher.EncryptString(m.AccountNumber)
		if err != nil {
			return vendorBankAccountModel{}, fmt.Errorf("failed to encrypt vendor bank account number: %w", err)
		}
		m.AccountNumber = encrypted
	}

	return m, nil
}

func (r *vendorRepository) toVendorResponsiblePersonModel(rp *domain.VendorResponsiblePerson) (vendorResponsiblePersonModel, error) {
	m := vendorResponsiblePersonModel{
		ID:        rp.ID.String(),
		VendorID:  rp.VendorID.String(),
		UserID:    rp.UserID.String(),
		NIK:       rp.NIK,
		CreatedAt: rp.CreatedAt,
		UpdatedAt: rp.UpdatedAt,
	}

	if m.NIK != "" {
		encrypted, err := r.fieldCipher.EncryptString(m.NIK)
		if err != nil {
			return vendorResponsiblePersonModel{}, fmt.Errorf("failed to encrypt responsible person nik: %w", err)
		}
		m.NIK = encrypted
	}

	return m, nil
}

func (r *vendorRepository) toDomainVendorBankAccount(m *vendorBankAccountModel) (*domain.VendorBankAccount, error) {
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

	if ba.AccountNumber != "" {
		decrypted, err := r.fieldCipher.DecryptString(ba.AccountNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt vendor bank account number: %w", err)
		}
		ba.AccountNumber = decrypted
	}

	return ba, nil
}

func (r *vendorRepository) toDomainVendorResponsiblePerson(m *vendorResponsiblePersonModel) (*domain.VendorResponsiblePerson, error) {
	id, _ := uuid.Parse(m.ID)
	vendorID, _ := uuid.Parse(m.VendorID)
	userID, _ := uuid.Parse(m.UserID)

	rp := &domain.VendorResponsiblePerson{
		ID:        id,
		VendorID:  vendorID,
		UserID:    userID,
		NIK:       m.NIK,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if rp.NIK != "" {
		decrypted, err := r.fieldCipher.DecryptString(rp.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt responsible person nik: %w", err)
		}
		rp.NIK = decrypted
	}

	return rp, nil
}

func toVendorDocumentModel(d *domain.VendorDocument) vendorDocumentModel {
	m := vendorDocumentModel{
		ID:            d.ID.String(),
		VendorID:      d.VendorID.String(),
		DocType:       d.DocType,
		FileURL:       d.FileURL,
		MimeType:      d.MimeType,
		FileSizeBytes: d.FileSizeBytes,
		FileChecksum:  d.FileChecksum,
		VerifiedAt:    d.VerifiedAt,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
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
		ID:            id,
		VendorID:      vendorID,
		DocType:       m.DocType,
		FileURL:       m.FileURL,
		MimeType:      m.MimeType,
		FileSizeBytes: m.FileSizeBytes,
		FileChecksum:  m.FileChecksum,
		VerifiedAt:    m.VerifiedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
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

func (r *vendorRepository) IncrementTotalSold(ctx context.Context, vendorID uuid.UUID, qty int) error {
	return r.db.WithContext(ctx).Model(&vendorModel{}).
		Where("id = ?", vendorID.String()).
		Updates(map[string]interface{}{
			"total_sold": gorm.Expr("total_sold + ?", qty),
			"updated_at": time.Now(),
		}).Error
}

func (r *vendorRepository) FindDocumentByVendorIDAndDocType(ctx context.Context, vendorID uuid.UUID, docType string) (*domain.VendorDocument, error) {
	var model vendorDocumentModel
	if err := r.db.WithContext(ctx).
		Where("vendor_id = ? AND doc_type = ?", vendorID.String(), docType).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainVendorDocument(&model), nil
}

func (r *vendorRepository) GetVendorAggregatedRating(ctx context.Context, vendorID uuid.UUID) (float64, int64, error) {
	type ratingRow struct {
		AvgRating    float64 `gorm:"column:avg_rating"`
		TotalReviews int64   `gorm:"column:total_reviews"`
	}
	var row ratingRow
	if err := r.db.WithContext(ctx).
		Table("product_review_stats prs").
		Select(`
			COALESCE(
				CASE WHEN SUM(prs.total_reviews) > 0
					THEN SUM(prs.total_stars)::float / SUM(prs.total_reviews)::float
					ELSE 0
				END, 0
			) AS avg_rating,
			COALESCE(SUM(prs.total_reviews), 0) AS total_reviews
		`).
		Joins("JOIN products p ON p.id = prs.product_id").
		Where("p.vendor_id = ?", vendorID.String()).
		Scan(&row).Error; err != nil {
		return 0, 0, err
	}
	return row.AvgRating, row.TotalReviews, nil
}
