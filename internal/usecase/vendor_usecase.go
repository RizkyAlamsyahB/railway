package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrVendorAlreadyExists    = errors.New("user already has a vendor")
	ErrVendorNotFound         = errors.New("vendor not found")
	ErrDocumentNotFound       = errors.New("document not found for this vendor")
	ErrObjectNotUploaded      = errors.New("object not found in storage")
	ErrInvalidDocumentContent = errors.New("invalid document content type")
	ErrDocumentSizeOverflow   = errors.New("document file size exceeds supported limit")
)

// Required document types for vendor registration.
var vendorDocTypes = []string{
	"owner_ktp",
	"owner_passport",
	"business_npwp",
	"store_photo",
	"bank_account_proof",
	"business_logo",
	"business_banner",
}

var allowedDocumentContentTypes = map[string]struct{}{
	"image/jpeg":      {},
	"image/png":       {},
	"image/webp":      {},
	"application/pdf": {},
}

// presignedUploadExpiry is the validity duration for presigned upload URLs.
const presignedUploadExpiry = 30 * time.Minute

type vendorUseCase struct {
	userRepo   domain.UserRepository
	vendorRepo domain.VendorRepository
	storage    domain.StorageProvider
	jwtSecret  string
	jwtExpiry  int
	jwtIssuer  string
}

// NewVendorUseCase creates a new VendorUseCase.
func NewVendorUseCase(
	userRepo domain.UserRepository,
	vendorRepo domain.VendorRepository,
	storage domain.StorageProvider,
	jwtSecret string,
	jwtExpiry int,
	jwtIssuer string,
) domain.VendorUseCase {
	return &vendorUseCase{
		userRepo:   userRepo,
		vendorRepo: vendorRepo,
		storage:    storage,
		jwtSecret:  jwtSecret,
		jwtExpiry:  jwtExpiry,
		jwtIssuer:  jwtIssuer,
	}
}

func (uc *vendorUseCase) Register(ctx context.Context, req domain.VendorRegisterRequest) (*domain.VendorRegisterResponse, error) {
	// 1. Check if email is already taken. If user exists, check for existing vendor.
	existing, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		existingVendor, err := uc.vendorRepo.FindByOwnerUserID(ctx, existing.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing vendor: %w", err)
		}
		if existingVendor != nil {
			return nil, ErrVendorAlreadyExists
		}
		return nil, ErrEmailAlreadyRegistered
	}

	// 2. Hash password.
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Prepare all entities in memory.
	now := time.Now()
	phone := req.Phone
	userID := uuid.New()
	vendorID := uuid.New()

	user := &domain.User{
		ID:           userID,
		Email:        req.Email,
		FullName:     req.OwnerName,
		Phone:        &phone,
		PasswordHash: hash,
		Status:       "pending",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	vendor := &domain.Vendor{
		ID:                    vendorID,
		OwnerUserID:           userID,
		VendorType:            req.StoreType,
		LegalName:             req.LegalName,
		DisplayName:           req.StoreName,
		ResponsiblePersonName: req.ResponsiblePersonName,
		Status:                "draft",
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	bankAccount := &domain.VendorBankAccount{
		ID:                 uuid.New(),
		VendorID:           vendorID,
		BankName:           req.BankName,
		AccountNumber:      req.BankAccountNumber,
		AccountHolderName:  req.BankAccountHolderName,
		VerificationStatus: "pending",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// 4. Generate presigned upload URLs BEFORE any DB writes.
	//    Also prepare VendorDocument records for each document type.
	uploadURLs := make([]domain.PresignedUploadInfo, 0, len(vendorDocTypes))
	documents := make([]domain.VendorDocument, 0, len(vendorDocTypes))

	for _, docType := range vendorDocTypes {
		objectKey := fmt.Sprintf("vendors/%s/documents/%s/%s", vendorID.String(), docType, uuid.New().String())

		// Keep content type unsigned so client can upload with the file's actual MIME type.
		uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, "", presignedUploadExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for %s: %w", docType, err)
		}

		uploadURLs = append(uploadURLs, domain.PresignedUploadInfo{
			DocType:   docType,
			UploadURL: uploadURL,
			ObjectKey: objectKey,
		})

		documents = append(documents, domain.VendorDocument{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            docType,
			FileURL:            objectKey,
			VerificationStatus: "pending",
			CreatedAt:          now,
			UpdatedAt:          now,
		})
	}

	// 5. Persist user to DB.
	if err := uc.userRepo.Create(ctx, user, "umkm"); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. Persist vendor + bank account + documents in one transaction.
	if err := uc.vendorRepo.Create(ctx, vendor, bankAccount, documents); err != nil {
		return nil, fmt.Errorf("failed to create vendor: %w", err)
	}

	// 7. Generate JWT token.
	token, err := auth.GenerateToken(userID, user.Email, []string{"umkm"}, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.VendorRegisterResponse{
		VendorID:   vendorID,
		UserID:     userID,
		Token:      token,
		UploadURLs: uploadURLs,
	}, nil
}

func (uc *vendorUseCase) ConfirmDocuments(ctx context.Context, userID uuid.UUID, req domain.ConfirmDocumentsRequest) (*domain.ConfirmDocumentsResponse, error) {
	// 1. Find vendor by owner_user_id.
	vendor, err := uc.vendorRepo.FindByOwnerUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	// 2. Load existing documents for this vendor.
	existingDocs, err := uc.vendorRepo.FindDocumentsByVendorID(ctx, vendor.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Build lookup map by doc_type.
	docMap := make(map[string]*domain.VendorDocument, len(existingDocs))
	for i := range existingDocs {
		docMap[existingDocs[i].DocType] = &existingDocs[i]
	}

	// 3. Validate each requested document and verify it exists in S3.
	now := time.Now()
	updatedDocs := make([]domain.VendorDocument, 0, len(req.Documents))

	for _, item := range req.Documents {
		doc, exists := docMap[item.DocType]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrDocumentNotFound, item.DocType)
		}

		// Verify object actually exists in S3 via HeadObject.
		info, err := uc.storage.HeadObject(ctx, item.ObjectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to verify object %s: %w", item.ObjectKey, err)
		}
		if info == nil {
			return nil, fmt.Errorf("%w: %s", ErrObjectNotUploaded, item.ObjectKey)
		}
		contentType := normalizeContentType(info.ContentType)
		if !isAllowedDocumentContentType(contentType) {
			return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidDocumentContent, item.DocType, info.ContentType)
		}
		if info.ContentLength > int64(math.MaxInt32) {
			return nil, fmt.Errorf("%w: %s (%d bytes)", ErrDocumentSizeOverflow, item.DocType, info.ContentLength)
		}

		// Prepare update.
		fileSize := int(info.ContentLength)
		doc.FileURL = item.ObjectKey
		doc.MimeType = &contentType
		doc.FileSizeBytes = &fileSize
		doc.UploadedBy = &userID
		doc.UpdatedAt = now

		updatedDocs = append(updatedDocs, *doc)
	}

	// 4. Determine whether all 7 required documents now have uploads.
	for i := range updatedDocs {
		docMap[updatedDocs[i].DocType] = &updatedDocs[i]
	}

	allUploaded := true
	for _, docType := range vendorDocTypes {
		doc, exists := docMap[docType]
		if !exists || doc.UploadedBy == nil {
			allUploaded = false
			break
		}
	}

	// 5. Determine new status.
	newStatus := ""
	if allUploaded && vendor.Status == "draft" {
		newStatus = "submitted"
	}

	// 6. Persist all document updates + optional status change in one transaction.
	if err := uc.vendorRepo.ConfirmDocumentsAndUpdateStatus(ctx, vendor.ID, updatedDocs, newStatus); err != nil {
		return nil, fmt.Errorf("failed to confirm documents: %w", err)
	}

	vendorStatus := vendor.Status
	if newStatus != "" {
		vendorStatus = newStatus
	}

	return &domain.ConfirmDocumentsResponse{
		VendorID:       vendor.ID,
		VendorStatus:   vendorStatus,
		DocumentsCount: len(req.Documents),
	}, nil
}

func normalizeContentType(contentType string) string {
	normalized := strings.TrimSpace(strings.ToLower(contentType))
	if idx := strings.Index(normalized, ";"); idx >= 0 {
		normalized = strings.TrimSpace(normalized[:idx])
	}
	return normalized
}

func isAllowedDocumentContentType(contentType string) bool {
	_, ok := allowedDocumentContentTypes[contentType]
	return ok
}
