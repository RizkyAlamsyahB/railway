package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Vendor represents the vendors table.
type Vendor struct {
	ID                    uuid.UUID  `json:"id"`
	OwnerUserID           uuid.UUID  `json:"owner_user_id"`
	VendorType            string     `json:"vendor_type"`
	LegalName             *string    `json:"legal_name,omitempty"`
	DisplayName           string     `json:"display_name"`
	ResponsiblePersonName string     `json:"responsible_person_name"`
	Description           *string    `json:"description,omitempty"`
	Status                string     `json:"status"`
	ApprovedBy            *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt            *time.Time `json:"approved_at,omitempty"`
	StatusReason          *string    `json:"status_reason,omitempty"`
	XenditAccountID       *string    `json:"xendit_account_id,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// VendorDocument represents the vendor_documents table.
type VendorDocument struct {
	ID                 uuid.UUID  `json:"id"`
	VendorID           uuid.UUID  `json:"vendor_id"`
	DocType            string     `json:"doc_type"`
	FileURL            string     `json:"file_url"`
	MimeType           *string    `json:"mime_type,omitempty"`
	FileSizeBytes      *int       `json:"file_size_bytes,omitempty"`
	FileChecksum       *string    `json:"file_checksum,omitempty"`
	UploadedBy         *uuid.UUID `json:"uploaded_by,omitempty"`
	VerificationStatus string     `json:"verification_status"`
	RejectionReason    *string    `json:"rejection_reason,omitempty"`
	VerifiedBy         *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// VendorBankAccount represents the vendor_bank_accounts table.
type VendorBankAccount struct {
	ID                 uuid.UUID  `json:"id"`
	VendorID           uuid.UUID  `json:"vendor_id"`
	BankName           string     `json:"bank_name"`
	AccountNumber      string     `json:"account_number"`
	AccountHolderName  string     `json:"account_holder_name"`
	VerificationStatus string     `json:"verification_status"`
	RejectionReason    *string    `json:"rejection_reason,omitempty"`
	VerifiedBy         *uuid.UUID `json:"verified_by,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// VendorRegisterRequest is the input DTO for vendor registration.
type VendorRegisterRequest struct {
	StoreName             string  `json:"store_name" binding:"required,max=120"`
	StoreType             string  `json:"store_type" binding:"required,oneof=umrah_souvenir_store hajj_souvenir_store general_souvenir_store"`
	OwnerName             string  `json:"owner_name" binding:"required,max=120"`
	LegalName             *string `json:"legal_name,omitempty" binding:"omitempty,max=160"`
	ResponsiblePersonName string  `json:"responsible_person_name" binding:"required,max=120"`
	Phone                 string  `json:"phone" binding:"required,max=20"`
	Email                 string  `json:"email" binding:"required,email,max=255"`
	Password              string  `json:"password" binding:"required,min=8"`
	BankName              string  `json:"bank_name" binding:"required,max=120"`
	BankAccountNumber     string  `json:"bank_account_number" binding:"required,max=60"`
	BankAccountHolderName string  `json:"bank_account_holder_name" binding:"required,max=160"`
}

// PresignedUploadInfo holds a presigned upload URL for a specific document type.
type PresignedUploadInfo struct {
	DocType   string `json:"doc_type"`
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

// VendorRegisterResponse is the output DTO for a successful vendor registration.
type VendorRegisterResponse struct {
	VendorID   uuid.UUID             `json:"vendor_id"`
	UserID     uuid.UUID             `json:"user_id"`
	Token      string                `json:"token"`
	UploadURLs []PresignedUploadInfo `json:"upload_urls"`
}

// ConfirmDocumentItem represents a single document in the confirm-upload request.
type ConfirmDocumentItem struct {
	DocType   string `json:"doc_type" binding:"required"`
	ObjectKey string `json:"object_key" binding:"required"`
}

// ConfirmDocumentsRequest is the input DTO for confirming document uploads.
type ConfirmDocumentsRequest struct {
	Documents []ConfirmDocumentItem `json:"documents" binding:"required,min=1,dive"`
}

// ConfirmDocumentsResponse is the output DTO for a successful document confirmation.
type ConfirmDocumentsResponse struct {
	VendorID       uuid.UUID `json:"vendor_id"`
	VendorStatus   string    `json:"vendor_status"`
	DocumentsCount int       `json:"documents_confirmed"`
}

// VendorLoginRequest is the input DTO for vendor authentication.
type VendorLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VendorLoginResponse is the output DTO for a successful vendor authentication.
type VendorLoginResponse struct {
	Token        string    `json:"token"`
	VendorID     uuid.UUID `json:"vendor_id"`
	VendorStatus string    `json:"vendor_status"`
	DisplayName  string    `json:"display_name"`
}

// VendorListParams holds query parameters for listing vendors.
type VendorListParams struct {
	Page       int
	Limit      int
	Status     string
	VendorType string
	Search     string
}

// AdminVendorListItem is the output DTO for a vendor in the admin list.
type AdminVendorListItem struct {
	ID                    uuid.UUID `json:"id"`
	DisplayName           string    `json:"display_name"`
	LegalName             *string   `json:"legal_name,omitempty"`
	VendorType            string    `json:"vendor_type"`
	ResponsiblePersonName string    `json:"responsible_person_name"`
	Status                string    `json:"status"`
	StatusReason          *string   `json:"status_reason,omitempty"`
	XenditAccountID       *string   `json:"xendit_account_id,omitempty"`
	OwnerName             string    `json:"owner_name"`
	OwnerEmail            string    `json:"owner_email"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// AdminVendorOwnerResponse is the owner info in the admin vendor detail.
type AdminVendorOwnerResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Phone    *string   `json:"phone,omitempty"`
	Status   string    `json:"status"`
}

// AdminVendorBankAccountResponse is the bank account info in the admin vendor detail.
type AdminVendorBankAccountResponse struct {
	ID                 uuid.UUID  `json:"id"`
	BankName           string     `json:"bank_name"`
	AccountNumber      string     `json:"account_number"`
	AccountHolderName  string     `json:"account_holder_name"`
	VerificationStatus string     `json:"verification_status"`
	RejectionReason    *string    `json:"rejection_reason,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// AdminVendorDocumentResponse is a document with presigned download URL in the admin vendor detail.
type AdminVendorDocumentResponse struct {
	ID                 uuid.UUID  `json:"id"`
	DocType            string     `json:"doc_type"`
	FileURL            string     `json:"file_url"`
	DownloadURL        string     `json:"download_url"`
	MimeType           *string    `json:"mime_type,omitempty"`
	FileSizeBytes      *int       `json:"file_size_bytes,omitempty"`
	VerificationStatus string     `json:"verification_status"`
	RejectionReason    *string    `json:"rejection_reason,omitempty"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// AdminVendorDetailResponse is the full vendor detail for admin review.
type AdminVendorDetailResponse struct {
	ID                    uuid.UUID                       `json:"id"`
	VendorType            string                          `json:"vendor_type"`
	DisplayName           string                          `json:"display_name"`
	LegalName             *string                         `json:"legal_name,omitempty"`
	ResponsiblePersonName string                          `json:"responsible_person_name"`
	Description           *string                         `json:"description,omitempty"`
	Status                string                          `json:"status"`
	StatusReason          *string                         `json:"status_reason,omitempty"`
	ApprovedAt            *time.Time                      `json:"approved_at,omitempty"`
	XenditAccountID       *string                         `json:"xendit_account_id,omitempty"`
	CreatedAt             time.Time                       `json:"created_at"`
	UpdatedAt             time.Time                       `json:"updated_at"`
	Owner                 AdminVendorOwnerResponse        `json:"owner"`
	BankAccount           *AdminVendorBankAccountResponse `json:"bank_account"`
	Documents             []AdminVendorDocumentResponse   `json:"documents"`
}

// VendorRepository defines the interface for vendor data access.
type VendorRepository interface {
	// Create inserts a new vendor, its bank account, and document placeholders in a single transaction.
	Create(ctx context.Context, vendor *Vendor, bankAccount *VendorBankAccount, documents []VendorDocument) error

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

	// ConfirmDocumentsAndUpdateStatus updates the given documents and optionally sets the vendor status, all in one transaction.
	// If newStatus is empty, the vendor status is not changed.
	ConfirmDocumentsAndUpdateStatus(ctx context.Context, vendorID uuid.UUID, documents []VendorDocument, newStatus string) error

	// UpdateStatus updates the vendor's columns specified in the updates map.
	UpdateStatus(ctx context.Context, vendorID uuid.UUID, updates map[string]interface{}) error
}

// VendorUseCase defines the interface for vendor business operations.
type VendorUseCase interface {
	// Register creates a new user (with role umkm), vendor, bank account, and document placeholders,
	// then returns presigned upload URLs for required documents.
	Register(ctx context.Context, req VendorRegisterRequest) (*VendorRegisterResponse, error)

	// ConfirmDocuments verifies that documents were uploaded to S3 and marks them as confirmed.
	ConfirmDocuments(ctx context.Context, userID uuid.UUID, req ConfirmDocumentsRequest) (*ConfirmDocumentsResponse, error)

	// Login authenticates a vendor user and returns a JWT token with vendor claims.
	Login(ctx context.Context, req VendorLoginRequest) (*VendorLoginResponse, error)
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

	// GetByID returns the full vendor detail including owner, bank account, and documents with presigned download URLs.
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
