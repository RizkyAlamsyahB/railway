package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Vendor represents the vendors table.
type Vendor struct {
	ID              uuid.UUID  `json:"id"`
	OwnerUserID     uuid.UUID  `json:"owner_user_id"`
	VendorType      *string    `json:"vendor_type"`
	LegalName       *string    `json:"legal_name,omitempty"`
	DisplayName     *string    `json:"display_name"`
	Description     *string    `json:"description,omitempty"`
	ProvinceID      *string    `json:"province_id,omitempty"`
	ProvinceName    *string    `json:"province_name,omitempty"`
	CityID          *string    `json:"city_id,omitempty"`
	CityName        *string    `json:"city_name,omitempty"`
	DistrictID      *string    `json:"district_id,omitempty"`
	DistrictName    *string    `json:"district_name,omitempty"`
	SubdistrictID   *string    `json:"subdistrict_id,omitempty"`
	SubdistrictName *string    `json:"subdistrict_name,omitempty"`
	PostalCode      *string    `json:"postal_code,omitempty"`
	AddressLine     *string    `json:"address_line,omitempty"`
	Status          string     `json:"status"`
	ApprovedBy      *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	StatusReason    *string    `json:"status_reason,omitempty"`
	XenditAccountID *string    `json:"xendit_account_id,omitempty"`
	TotalSold       int        `json:"total_sold"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// VendorDocument represents the vendor_documents table.
type VendorDocument struct {
	ID            uuid.UUID  `json:"id"`
	VendorID      uuid.UUID  `json:"vendor_id"`
	DocType       string     `json:"doc_type"`
	FileURL       string     `json:"file_url"`
	MimeType      *string    `json:"mime_type,omitempty"`
	FileSizeBytes *int       `json:"file_size_bytes,omitempty"`
	FileChecksum  *string    `json:"file_checksum,omitempty"`
	UploadedBy    *uuid.UUID `json:"uploaded_by,omitempty"`
	VerifiedBy    *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// VendorBankAccount represents the vendor_bank_accounts table.
type VendorBankAccount struct {
	ID                uuid.UUID `json:"id"`
	VendorID          uuid.UUID `json:"vendor_id"`
	ChannelCode       string    `json:"channel_code"`
	BankName          string    `json:"bank_name"`
	AccountNumber     string    `json:"account_number"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountLast4      string    `json:"account_last4"`
	IsDefault         bool      `json:"is_default"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// VendorBankAccountListItem is the DTO returned when listing vendor bank accounts.
type VendorBankAccountListItem struct {
	ID                uuid.UUID `json:"id"`
	ChannelCode       string    `json:"channel_code"`
	BankName          string    `json:"bank_name"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountLast4      string    `json:"account_last4"`
	IsDefault         bool      `json:"is_default"`
}

// CreateVendorBankAccountRequest is the input DTO for adding a vendor bank account.
type CreateVendorBankAccountRequest struct {
	ChannelCode       string `json:"channel_code" binding:"required,max=40"`
	BankName          string `json:"bank_name" binding:"required,max=80"`
	AccountNumber     string `json:"account_number" binding:"required,max=64"`
	AccountHolderName string `json:"account_holder_name" binding:"required,max=120"`
	IsDefault         bool   `json:"is_default"`
}

// CreateVendorBankAccountResponse is the output DTO after creating a vendor bank account.
type CreateVendorBankAccountResponse struct {
	ID                uuid.UUID `json:"id"`
	ChannelCode       string    `json:"channel_code"`
	BankName          string    `json:"bank_name"`
	AccountHolderName string    `json:"account_holder_name"`
	AccountLast4      string    `json:"account_last4"`
	IsDefault         bool      `json:"is_default"`
}

// VendorResponsiblePerson stores vendor-specific KYC data for the owner contact.
type VendorResponsiblePerson struct {
	ID        uuid.UUID `json:"id"`
	VendorID  uuid.UUID `json:"vendor_id"`
	UserID    uuid.UUID `json:"user_id"`
	NIK       string    `json:"nik"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PresignedUploadInfo holds a presigned upload URL for a specific document type.
type PresignedUploadInfo struct {
	DocType   string `json:"doc_type"`
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

const (
	VendorDocumentTypeOwnerDocumentID  = "owner_document_id"
	VendorDocumentTypeBusinessNIB      = "business_nib"
	VendorDocumentTypeBusinessNPWP     = "business_npwp"
	VendorDocumentTypeHalalCertificate = "halal_certificate"
	VendorDocumentTypeStorePhoto       = "store_photo"
	VendorDocumentTypeBankAccountProof = "bank_account_proof"
	VendorDocumentTypeBusinessLogo     = "business_logo"
	VendorDocumentTypeBusinessBanner   = "business_banner"
)

// VendorLoginRequest is the input DTO for vendor authentication.
type VendorLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VendorLoginResponse is the output DTO for a successful vendor authentication.
type VendorLoginResponse struct {
	AccessToken  string    `json:"access_token"`
	VendorID     uuid.UUID `json:"vendor_id"`
	Email        string    `json:"email"`
	ImageURL     *string   `json:"image_url"`
	StoreName    *string   `json:"store_name"`
	VendorType   *string   `json:"vendor_type"`
	VendorStatus string    `json:"vendor_status"`
}

// VendorMeResponse is the output DTO for the authenticated vendor profile endpoint.
// It mirrors the vendor login response shape without an access token.
type VendorMeResponse struct {
	VendorID     uuid.UUID `json:"vendor_id"`
	ImageURL     *string   `json:"image_url"`
	Email        string    `json:"email"`
	VendorType   *string   `json:"vendor_type"`
	VendorStatus string    `json:"vendor_status"`
	StoreName    *string   `json:"store_name"`
	StatusReason *string   `json:"status_reason"`
}

// VendorUpdateProfileRequest is the input DTO for updating the vendor's profile (name & description).
type VendorUpdateProfileRequest struct {
	StoreName        string  `json:"store_name" binding:"required,max=120"`
	StoreDescription *string `json:"store_description" binding:"omitempty"`
}

type VendorSouvenirStoreProposalAddress struct {
	ProvinceID    string `json:"province_id" binding:"required,max=20"`
	CityID        string `json:"city_id" binding:"required,max=20"`
	DistrictID    string `json:"district_id" binding:"required,max=20"`
	SubdistrictID string `json:"subdistrict_id" binding:"required,max=20"`
	PostalCode    string `json:"postal_code" binding:"required,max=10"`
	AddressLine   string `json:"address_line" binding:"required"`
}

type VendorSouvenirStoreProposalResponsiblePerson struct {
	Name        string `json:"name" binding:"required,max=120"`
	Phone       string `json:"phone" binding:"required,max=20"`
	Email       string `json:"email" binding:"required,email,max=255"`
	NIK         string `json:"nik" binding:"required,max=32"`
	KTPObjectID string `json:"ktp_object_id" binding:"required"`
}

type VendorSouvenirStoreProposalOtherDocuments struct {
	NIBObjectID              string `json:"nib_object_id" binding:"required"`
	HalalCertificateObjectID string `json:"halal_certificate_object_id" binding:"required"`
}

// VendorSouvenirStoreProposalRequest is the input DTO for vendor onboarding proposal submission.
type VendorSouvenirStoreProposalRequest struct {
	StoreName         string                                       `json:"store_name" binding:"required,max=120"`
	StoreDescription  string                                       `json:"store_description" binding:"required"`
	Address           VendorSouvenirStoreProposalAddress           `json:"address" binding:"required"`
	ResponsiblePerson VendorSouvenirStoreProposalResponsiblePerson `json:"responsible_person" binding:"required"`
	OtherDocuments    VendorSouvenirStoreProposalOtherDocuments    `json:"other_documents" binding:"required"`
}

// VendorSouvenirStoreProposalResponse is the output DTO after vendor proposal submission.
type VendorSouvenirStoreProposalResponse struct {
	VendorID uuid.UUID `json:"vendor_id"`
	Status   string    `json:"status"`
}

type VendorSouvenirStoreProposalPresignRequest struct {
	ResponsiblePersonKTP string `json:"responsible_person_ktp" binding:"required"`
	NIB                  string `json:"nib" binding:"required"`
	HalalCertificate     string `json:"halal_certificate" binding:"required"`
}

type VendorSouvenirStoreProposalPresignItem struct {
	ObjectID     string `json:"object_id"`
	PresignedURL string `json:"presigned_url"`
}

type VendorSouvenirStoreProposalPresignResponse struct {
	ResponsiblePersonKTP VendorSouvenirStoreProposalPresignItem `json:"responsible_person_ktp"`
	NIB                  VendorSouvenirStoreProposalPresignItem `json:"nib"`
	HalalCertificate     VendorSouvenirStoreProposalPresignItem `json:"halal_certificate"`
}

// SubmitSouvenirStoreProposalInput bundles all data needed to atomically submit a vendor proposal.
type SubmitSouvenirStoreProposalInput struct {
	Vendor            *Vendor
	OwnerUser         *User
	ResponsiblePerson *VendorResponsiblePerson
	Documents         []VendorDocument
}

// VendorListParams holds query parameters for listing vendors.
type VendorListParams struct {
	Page       int
	Limit      int
	Status     string
	VendorType string
	Search     string
}

// AdminVendorListAddress is the nested address object in the admin vendor list response.
type AdminVendorListAddress struct {
	Province    *string `json:"province,omitempty"`
	City        *string `json:"city,omitempty"`
	District    *string `json:"district,omitempty"`
	Subdistrict *string `json:"subdistrict,omitempty"`
	PostalCode  *string `json:"postal_code,omitempty"`
}

// AdminVendorListItem is the output DTO for a vendor in the admin list.
type AdminVendorListItem struct {
	ID        uuid.UUID               `json:"id"`
	Status    string                  `json:"status"`
	StoreType *string                 `json:"store_type"`
	StoreName *string                 `json:"store_name"`
	Email     string                  `json:"email"`
	Address   *AdminVendorListAddress `json:"address"`
	CreatedAt time.Time               `json:"created_at"`
}

// AdminVendorDetailAddress is the nested address object in admin vendor detail.
type AdminVendorDetailAddress struct {
	Province    *string `json:"province"`
	City        *string `json:"city"`
	District    *string `json:"district"`
	Subdistrict *string `json:"subdistrict"`
	PostalCode  *string `json:"postal_code"`
	AddressLine *string `json:"address_line"`
}

// AdminVendorKTPResponse is the KTP document object in responsible_person.
type AdminVendorKTPResponse struct {
	ID            uuid.UUID `json:"id"`
	DocType       string    `json:"doc_type"`
	FileURL       string    `json:"file_url"`
	MimeType      *string   `json:"mime_type,omitempty"`
	FileSizeBytes *int      `json:"file_size_bytes,omitempty"`
}

// AdminVendorResponsiblePersonResponse is the responsible person section in admin vendor detail.
type AdminVendorResponsiblePersonResponse struct {
	Name  *string                 `json:"name"`
	Phone *string                 `json:"phone"`
	Email *string                 `json:"email"`
	NIK   *string                 `json:"nik"`
	KTP   *AdminVendorKTPResponse `json:"ktp"`
}

// AdminVendorRequirementDocumentResponse is one required document in admin vendor detail.
type AdminVendorRequirementDocumentResponse struct {
	ID            uuid.UUID `json:"id"`
	DocType       string    `json:"doc_type"`
	FileURL       string    `json:"file_url"`
	MimeType      *string   `json:"mime_type,omitempty"`
	FileSizeBytes *int      `json:"file_size_bytes,omitempty"`
}

// AdminVendorDetailResponse is the full vendor detail for admin review.
type AdminVendorDetailResponse struct {
	ID                    uuid.UUID                                `json:"id"`
	StoreName             *string                                  `json:"store_name"`
	StoreDescription      *string                                  `json:"store_description"`
	StoreType             *string                                  `json:"store_type"`
	Status                string                                   `json:"status"`
	StatusReason          *string                                  `json:"status_reason,omitempty"`
	Email                 string                                   `json:"email"`
	Address               AdminVendorDetailAddress                 `json:"address"`
	ResponsiblePerson     AdminVendorResponsiblePersonResponse     `json:"responsible_person"`
	RequirementsDocuments []AdminVendorRequirementDocumentResponse `json:"requirements_documents"`
	CreatedAt             time.Time                                `json:"created_at"`
}

// VendorRepository defines the interface for vendor data access.
type VendorRepository interface {
	// Create inserts a new vendor, its bank account, and document placeholders in a single transaction.
	Create(ctx context.Context, vendor *Vendor, bankAccount *VendorBankAccount, documents []VendorDocument) error

	// CreateMinimal inserts a new vendor and initial documents without creating a bank account.
	CreateMinimal(ctx context.Context, vendor *Vendor, documents []VendorDocument) error

	// FindByID returns the vendor with the given ID, or nil if not found.
	FindByID(ctx context.Context, id uuid.UUID) (*Vendor, error)

	// FindByOwnerUserID returns the vendor owned by the given user, or nil if none.
	FindByOwnerUserID(ctx context.Context, userID uuid.UUID) (*Vendor, error)

	// List returns vendors matching the given params, plus total count for pagination.
	List(ctx context.Context, params VendorListParams) ([]Vendor, int64, error)

	// FindDocumentsByVendorID returns all documents for the given vendor.
	FindDocumentsByVendorID(ctx context.Context, vendorID uuid.UUID) ([]VendorDocument, error)

	// FindBankAccountByVendorID returns the bank account for the given vendor, or nil if not found.
	FindBankAccountByVendorID(ctx context.Context, vendorID uuid.UUID) (*VendorBankAccount, error)

	// FindBankAccountByID returns a bank account by its primary key, or nil if not found.
	FindBankAccountByID(ctx context.Context, id uuid.UUID) (*VendorBankAccount, error)

	// ListBankAccountsByVendorID returns all bank accounts owned by the given vendor.
	ListBankAccountsByVendorID(ctx context.Context, vendorID uuid.UUID) ([]VendorBankAccountListItem, error)

	// FindResponsiblePersonByVendorID returns the vendor responsible person row, or nil if not found.
	FindResponsiblePersonByVendorID(ctx context.Context, vendorID uuid.UUID) (*VendorResponsiblePerson, error)

	// UpsertBankAccount creates or updates the vendor bank account for the given vendor.
	UpsertBankAccount(ctx context.Context, bankAccount *VendorBankAccount) error

	// CreateBankAccount inserts a new bank account for the given vendor.
	CreateBankAccount(ctx context.Context, bankAccount *VendorBankAccount) error

	// DeleteBankAccount removes a vendor bank account by id and vendor ownership.
	DeleteBankAccount(ctx context.Context, vendorID, bankAccountID uuid.UUID) error

	// SetDefaultBankAccount sets one vendor bank account as default and unsets the previous default.
	SetDefaultBankAccount(ctx context.Context, vendorID, bankAccountID uuid.UUID) error

	// ConfirmDocumentsAndUpdateStatus updates the given documents and optionally sets the vendor status, all in one transaction.
	// If newStatus is empty, the vendor status is not changed.
	ConfirmDocumentsAndUpdateStatus(ctx context.Context, vendorID uuid.UUID, documents []VendorDocument, newStatus string) error

	// UpdateStatus updates the vendor's columns specified in the updates map.
	UpdateStatus(ctx context.Context, vendorID uuid.UUID, updates map[string]any) error

	// UpdateProfile updates the display name (store name) and description for the vendor.
	UpdateProfile(ctx context.Context, vendorID uuid.UUID, displayName string, description *string) error

	// SubmitSouvenirStoreProposal atomically updates owner user profile, vendor profile, KYC data, and required documents.
	SubmitSouvenirStoreProposal(ctx context.Context, input SubmitSouvenirStoreProposalInput) error

	// IncrementTotalSold atomically adds qty to the vendor's total_sold counter.
	IncrementTotalSold(ctx context.Context, vendorID uuid.UUID, qty int) error

	// FindDocumentByVendorIDAndDocType returns a single document by vendor ID and doc_type, or nil.
	FindDocumentByVendorIDAndDocType(ctx context.Context, vendorID uuid.UUID, docType string) (*VendorDocument, error)

	// GetVendorAggregatedRating returns the aggregated average rating and total reviews
	// across all products belonging to the vendor.
	GetVendorAggregatedRating(ctx context.Context, vendorID uuid.UUID) (avgRating float64, totalReviews int64, err error)

	// FindPayoutBatchByID returns the payout batch with the given ID, or nil if not found.
	FindPayoutBatchByID(ctx context.Context, id uuid.UUID) (*PayoutBatch, error)

	// UpdatePayoutBatch updates the payout batch record with Xendit response data.
	UpdatePayoutBatch(ctx context.Context, batch *PayoutBatch) error

	// --- Balance & Withdrawal ---

	// InitBalance creates a zero-balance row for a newly registered vendor.
	InitBalance(ctx context.Context, vendorID uuid.UUID) error

	// GetBalance returns the vendor's current balance, or nil if not found.
	GetBalance(ctx context.Context, vendorID uuid.UUID) (*VendorBalance, error)

	// CreditBalance adds amount to available_balance and total_earned atomically.
	CreditBalance(ctx context.Context, vendorID uuid.UUID, amount float64) error

	// DebitBalance moves amount from available_balance to pending_balance atomically.
	// Returns error if available_balance < amount.
	DebitBalance(ctx context.Context, vendorID uuid.UUID, amount float64) error

	// CompleteWithdrawal moves amount from pending_balance to total_withdrawn atomically.
	CompleteWithdrawal(ctx context.Context, vendorID uuid.UUID, amount float64) error

	// FailWithdrawal moves amount from pending_balance back to available_balance atomically.
	FailWithdrawal(ctx context.Context, vendorID uuid.UUID, amount float64) error

	// CreateWithdrawal inserts a new vendor_withdrawals record.
	CreateWithdrawal(ctx context.Context, withdrawal *VendorWithdrawal) error

	// FindWithdrawalByID returns withdrawal by its primary key, or nil if not found.
	FindWithdrawalByID(ctx context.Context, id uuid.UUID) (*VendorWithdrawal, error)

	// UpdateWithdrawal updates an existing vendor_withdrawals record.
	UpdateWithdrawal(ctx context.Context, withdrawal *VendorWithdrawal) error

	// ApplyWithdrawalWebhookUpdate applies status/xendit updates and optional balance movement atomically.
	ApplyWithdrawalWebhookUpdate(ctx context.Context, withdrawalID uuid.UUID, newStatus string, xenditStatus string, xenditPayoutID *string, failedReason *string, balanceAction string, completion *WithdrawalCompletionData) error

	// ListWithdrawals returns paginated withdrawals for a vendor.
	ListWithdrawals(ctx context.Context, vendorID uuid.UUID, params VendorWithdrawalListParams) ([]VendorWithdrawal, int64, error)
}

// Balance actions applied during payout webhook processing.
const (
	WithdrawalBalanceActionNone             = "none"
	WithdrawalBalanceActionComplete         = "complete"
	WithdrawalBalanceActionFail             = "fail"
	WithdrawalBalanceActionReverseCompleted = "reverse_completed"
)

// Fee resolution statuses for vendor withdrawals.
const (
	WithdrawalFeeStatusPending  = "pending"
	WithdrawalFeeStatusResolved = "resolved"
	WithdrawalFeeStatusFailed   = "failed"
)

// WithdrawalCompletionData contains the fee-aware completion values from Xendit transaction data.
type WithdrawalCompletionData struct {
	FeeActual           float64 `json:"fee_actual"`
	NetAmount           float64 `json:"net_amount"`
	TotalDeducted       float64 `json:"total_deducted"`
	XenditTransactionID *string `json:"xendit_transaction_id,omitempty"`
}

// PayoutBatch represents the payout_batches table (used by the finance admin module).
type PayoutBatch struct {
	ID             uuid.UUID  `json:"id"`
	VendorID       uuid.UUID  `json:"vendor_id"`
	PeriodStart    time.Time  `json:"period_start"`
	PeriodEnd      time.Time  `json:"period_end"`
	Status         string     `json:"status"`
	TotalGross     float64    `json:"total_gross"`
	TotalFee       float64    `json:"total_fee"`
	TotalNet       float64    `json:"total_net"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	XenditPayoutID *string    `json:"xendit_payout_id,omitempty"`
	ChannelCode    *string    `json:"channel_code,omitempty"`
	Description    *string    `json:"description,omitempty"`
	XenditStatus   *string    `json:"xendit_status,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// VendorBalance represents the vendor_balances table.
type VendorBalance struct {
	VendorID         uuid.UUID `json:"vendor_id"`
	AvailableBalance float64   `json:"available_balance"`
	PendingBalance   float64   `json:"pending_balance"`
	EscrowBalance    float64   `json:"escrow_balance"`
	TotalEarned      float64   `json:"total_earned"`
	TotalWithdrawn   float64   `json:"total_withdrawn"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// VendorWithdrawal represents the vendor_withdrawals table.
type VendorWithdrawal struct {
	ID                  uuid.UUID `json:"id"`
	VendorID            uuid.UUID `json:"vendor_id"`
	Amount              float64   `json:"amount"`
	ChannelCode         string    `json:"channel_code"`
	Status              string    `json:"status"`
	FeeEstimated        float64   `json:"fee_estimated"`
	FeeActual           *float64  `json:"fee_actual,omitempty"`
	AmountNet           float64   `json:"amount_net"`
	TotalDeducted       float64   `json:"total_deducted"`
	FeeStatus           string    `json:"fee_status"`
	XenditPayoutID      *string   `json:"xendit_payout_id,omitempty"`
	XenditStatus        *string   `json:"xendit_status,omitempty"`
	XenditTransactionID *string   `json:"xendit_transaction_id,omitempty"`
	Description         *string   `json:"description,omitempty"`
	FailedReason        *string   `json:"failed_reason,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

const (
	VendorTypeSouvenirStore = "souvenir_store"
	VendorTypePPIU          = "ppiu"
	VendorTypeHajjDormitory = "hajj_dormitory"
)

// VendorWithdrawRequest is the input DTO for vendor self-service withdrawal.
type VendorWithdrawRequest struct {
	Amount        float64   `json:"amount" binding:"required,gte=10000"`
	ChannelCode   string    `json:"channel_code" binding:"required"`
	BankAccountID uuid.UUID `json:"bank_account_id" binding:"required"`
}

// VendorWithdrawResponse is the output DTO for a successful withdrawal request.
type VendorWithdrawResponse struct {
	WithdrawalID       uuid.UUID `json:"withdrawal_id"`
	XenditPayoutID     string    `json:"xendit_payout_id"`
	Status             string    `json:"status"`
	XenditStatus       string    `json:"xendit_status"`
	Amount             float64   `json:"amount"`
	RequestedAmount    float64   `json:"requested_amount"`
	EstimatedFee       float64   `json:"estimated_fee"`
	EstimatedNetAmount float64   `json:"estimated_net_amount"`
	ChannelCode        string    `json:"channel_code"`
}

// VendorWithdrawalListParams is the query params for vendor withdrawal history.
type VendorWithdrawalListParams struct {
	Page   int
	Limit  int
	Status string
}

// VendorWithdrawalListItem is one item in vendor withdrawal history.
type VendorWithdrawalListItem struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	PayoutDate time.Time `json:"payout_date"`
}

// VendorBalanceResponse is the output DTO for the vendor balance endpoint.
type VendorBalanceResponse struct {
	AvailableBalance float64 `json:"available_balance"`
	PendingBalance   float64 `json:"pending_balance"`
	EscrowBalance    float64 `json:"escrow_balance"`
	TotalEarned      float64 `json:"total_earned"`
	TotalWithdrawn   float64 `json:"total_withdrawn"`
}

// VendorPayoutChannelItem represents one payout channel item for vendor frontend.
type VendorPayoutChannelItem struct {
	ChannelCode     string `json:"channel_code"`
	ChannelName     string `json:"channel_name"`
	Currency        string `json:"currency"`
	ChannelCategory string `json:"channel_category"`
}

// VendorPayoutChannelsResponse is the output DTO for vendor payout channels endpoint.
type VendorPayoutChannelsResponse struct {
	Channels []VendorPayoutChannelItem `json:"channels"`
}

// VendorUseCase defines the interface for vendor business operations.
type VendorUseCase interface {
	RequestRegistrationOTP(ctx context.Context, req VendorRegisterOTPRequest) (*VendorRegisterOTPResponse, error)
	VerifyRegistrationOTP(ctx context.Context, req VendorVerifyRegistrationOTPRequest) (*VendorVerifyRegistrationOTPResponse, error)
	SetRegistrationPassword(ctx context.Context, onboardingID uuid.UUID, req VendorRegistrationPasswordRequest) (*VendorLoginResponse, error)

	// Login authenticates a vendor user and returns a JWT token with vendor claims.
	Login(ctx context.Context, req VendorLoginRequest) (*VendorLoginResponse, error)

	// GetMe returns the authenticated vendor profile by vendor ID from auth claims.
	GetMe(ctx context.Context, vendorID uuid.UUID) (*VendorMeResponse, error)

	// UpdateProfile updates the vendor's profile properties like store name and description.
	UpdateProfile(ctx context.Context, vendorID uuid.UUID, req VendorUpdateProfileRequest) (*VendorMeResponse, error)

	// SubmitSouvenirStoreProposal submits the required vendor profile and documents for review.
	SubmitSouvenirStoreProposal(ctx context.Context, vendorID uuid.UUID, userID uuid.UUID, req VendorSouvenirStoreProposalRequest) (*VendorSouvenirStoreProposalResponse, error)

	// PresignSouvenirStoreProposalDocuments creates presigned upload URLs for required proposal documents.
	PresignSouvenirStoreProposalDocuments(ctx context.Context, vendorID uuid.UUID, userID uuid.UUID, req VendorSouvenirStoreProposalPresignRequest) (*VendorSouvenirStoreProposalPresignResponse, error)

	// GetBalance returns the vendor's current balance.
	GetBalance(ctx context.Context, vendorID uuid.UUID) (*VendorBalanceResponse, error)

	// ListPayoutChannels returns payout channels for vendor withdrawals.
	ListPayoutChannels(ctx context.Context, vendorID uuid.UUID) (*VendorPayoutChannelsResponse, error)

	// CreateBankAccount adds a vendor bank account.
	CreateBankAccount(ctx context.Context, vendorID uuid.UUID, req CreateVendorBankAccountRequest) (*CreateVendorBankAccountResponse, error)

	// ListBankAccounts returns all vendor bank accounts.
	ListBankAccounts(ctx context.Context, vendorID uuid.UUID) ([]VendorBankAccountListItem, error)

	// DeleteBankAccount removes a vendor bank account by id.
	DeleteBankAccount(ctx context.Context, vendorID, bankAccountID uuid.UUID) error

	// SetDefaultBankAccount sets one bank account as default.
	SetDefaultBankAccount(ctx context.Context, vendorID, bankAccountID uuid.UUID) error

	// RequestWithdrawal initiates a self-service withdrawal to the vendor's bank account via Xendit.
	RequestWithdrawal(ctx context.Context, vendorID uuid.UUID, req VendorWithdrawRequest) (*VendorWithdrawResponse, error)

	// ListWithdrawals returns paginated withdrawal history for the authenticated vendor.
	ListWithdrawals(ctx context.Context, vendorID uuid.UUID, params VendorWithdrawalListParams) ([]VendorWithdrawalListItem, *PaginationMeta, error)

	// HandlePayoutWebhook processes Xendit payout webhook callbacks for vendor withdrawals.
	HandlePayoutWebhook(ctx context.Context, payload XenditPayoutWebhookPayload) error
}

// --- Vendor Dashboard ---

// VendorDashboardResponse is the response DTO for the vendor seller dashboard.
type VendorDashboardResponse struct {
	Month               string                       `json:"month"`
	Year                int                          `json:"year"`
	Summary             VendorDashboardSummary       `json:"summary"`
	TodayTransactions   []VendorDashboardTransaction `json:"today_transactions"`
	RecentNotifications []NotificationItem           `json:"recent_notifications"`
	PaymentFlow         []VendorPaymentFlowItem      `json:"payment_flow"`
}

// VendorDashboardSummary holds the KPI summary cards for the vendor dashboard.
type VendorDashboardSummary struct {
	TotalTransactions            int64   `json:"total_transactions"`
	SuccessfulTransactionPercent float64 `json:"successful_transaction_percent"`
	Revenue                      float64 `json:"revenue"`
	PendingSettlement            float64 `json:"pending_settlement"`
}

// VendorDashboardTransaction represents a single row in the "today's transactions" table.
type VendorDashboardTransaction struct {
	DateTime    time.Time `json:"date_time"`
	Invoice     string    `json:"invoice"`
	ProductName string    `json:"product_name"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
}

// VendorPaymentFlowItem represents one slice in the payment flow donut chart.
type VendorPaymentFlowItem struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// VendorDashboardUseCase defines the interface for the vendor seller dashboard.
type VendorDashboardUseCase interface {
	GetDashboard(ctx context.Context, vendorID uuid.UUID, month, year int) (*VendorDashboardResponse, error)
}

// --- Store Detail (Public) ---

// StoreInfoResponse represents the store header info in the public store detail.
type StoreInfoResponse struct {
	ID            uuid.UUID `json:"id"`
	DisplayName   string    `json:"display_name"`
	Description   *string   `json:"description,omitempty"`
	Location      string    `json:"location"`
	ImageURL      string    `json:"image_url"`
	Rating        float64   `json:"rating"`
	TotalReviews  int64     `json:"total_reviews"`
	TotalSold     int       `json:"total_sold"`
	TotalProducts int64     `json:"total_products"`
}

// StoreProductItem represents a product card in the public store detail page.
type StoreProductItem struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	PromoPrice    *float64  `json:"promo_price,omitempty"`
	HasPromo      bool      `json:"has_promo"`
	Currency      string    `json:"currency"`
	ImageURL      string    `json:"image_url"`
	RatingAverage float64   `json:"rating_average"`
	RatingCount   int64     `json:"rating_count"`
	TotalSold     int64     `json:"total_sold"`
}

// StoreDetailResponse is the aggregated response for the public store detail endpoint.
type StoreDetailResponse struct {
	Store               StoreInfoResponse      `json:"store"`
	Banners             []VendorBannerResponse `json:"banners"`
	RecommendedProducts []StoreProductItem     `json:"recommended_products"`
	NewProducts         []StoreProductItem     `json:"new_products"`
	Products            []StoreProductItem     `json:"products"`
	ProductsPagination  *PaginationMeta        `json:"products_pagination"`
}

// VendorReviewItem represents a single review in the public store endpoint.
type VendorReviewItem struct {
	ID        uuid.UUID           `json:"id"`
	User      VendorReviewUser    `json:"user"`
	Rating    int                 `json:"rating"`
	Comment   string              `json:"comment"`
	CreatedAt time.Time           `json:"created_at"`
	Variant   VendorReviewVariant `json:"variant"`
	ImageKeys []string            `json:"-"`
	Images    []string            `json:"images"`
}

type VendorReviewUser struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
}

type VendorReviewVariant struct {
	Name string `json:"name"`
}

type VendorReviewSummary struct {
	AverageRating   float64          `json:"average_rating"`
	TotalReviews    int64            `json:"total_reviews"`
	RatingBreakdown map[string]int64 `json:"rating_breakdown"`
}

// VendorReviewResponse is the aggregated response for the store reviews endpoint.
type VendorReviewResponse struct {
	Reviews    []VendorReviewItem  `json:"reviews"`
	Summary    VendorReviewSummary `json:"summary"`
	Pagination *PaginationMeta     `json:"pagination"`
}

// StoreUseCase defines the interface for public store detail operations.
type StoreUseCase interface {
	// GetStoreDetail returns the full public store detail page data.
	GetStoreDetail(ctx context.Context, vendorID uuid.UUID, productsPage int, productsLimit int) (*StoreDetailResponse, error)

	// GetStoreReviews returns the public paginated reviews for a store.
	GetStoreReviews(ctx context.Context, vendorID uuid.UUID, page, limit int) (*VendorReviewResponse, error)
}

// AdminVendorReasonRequest is the input DTO for reject/block actions that require a reason.
type AdminVendorReasonRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// AdminVendorActionResponse is the output DTO for admin vendor status actions.
type AdminVendorActionResponse struct {
	VendorID uuid.UUID `json:"vendor_id"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
}

// AdminVendorUseCase defines the interface for admin vendor review operations.
type AdminVendorUseCase interface {
	// List returns a paginated list of vendors with owner info.
	List(ctx context.Context, params VendorListParams) ([]AdminVendorListItem, *PaginationMeta, error)

	// GetByID returns the full vendor detail for admin review.
	GetByID(ctx context.Context, id uuid.UUID) (*AdminVendorDetailResponse, error)

	// Approve transitions a vendor from submitted/rejected to active.
	Approve(ctx context.Context, vendorID uuid.UUID, adminID uuid.UUID) (*AdminVendorActionResponse, error)

	// Reject transitions a vendor from submitted to rejected.
	Reject(ctx context.Context, vendorID uuid.UUID, reason string) (*AdminVendorActionResponse, error)

	// Block transitions a vendor from active/submitted to blocked.
	Block(ctx context.Context, vendorID uuid.UUID, reason string) (*AdminVendorActionResponse, error)

	// Unblock transitions a vendor from blocked to active.
	Unblock(ctx context.Context, vendorID uuid.UUID, adminID uuid.UUID) (*AdminVendorActionResponse, error)
}
