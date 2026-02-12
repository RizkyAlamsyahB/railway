package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// presignedDownloadExpiry is the validity duration for presigned download URLs.
const presignedDownloadExpiry = 30 * time.Minute

type adminVendorUseCase struct {
	vendorRepo domain.VendorRepository
	userRepo   domain.UserRepository
	storage    domain.StorageProvider
}

// NewAdminVendorUseCase creates a new AdminVendorUseCase.
func NewAdminVendorUseCase(
	vendorRepo domain.VendorRepository,
	userRepo domain.UserRepository,
	storage domain.StorageProvider,
) domain.AdminVendorUseCase {
	return &adminVendorUseCase{
		vendorRepo: vendorRepo,
		userRepo:   userRepo,
		storage:    storage,
	}
}

func (uc *adminVendorUseCase) List(ctx context.Context, params domain.VendorListParams) ([]domain.AdminVendorListItem, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	vendors, total, err := uc.vendorRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list vendors: %w", err)
	}

	items := make([]domain.AdminVendorListItem, len(vendors))
	for i, v := range vendors {
		item := domain.AdminVendorListItem{
			ID:                    v.ID,
			DisplayName:           v.DisplayName,
			LegalName:             v.LegalName,
			VendorType:            v.VendorType,
			ResponsiblePersonName: v.ResponsiblePersonName,
			Status:                v.Status,
			CreatedAt:             v.CreatedAt,
			UpdatedAt:             v.UpdatedAt,
		}

		owner, err := uc.userRepo.FindByID(ctx, v.OwnerUserID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to find owner for vendor %s: %w", v.ID, err)
		}
		if owner != nil {
			item.OwnerName = owner.FullName
			item.OwnerEmail = owner.Email
		}

		items[i] = item
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}

func (uc *adminVendorUseCase) GetByID(ctx context.Context, id uuid.UUID) (*domain.AdminVendorDetailResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	// Fetch owner user info.
	var ownerResp domain.AdminVendorOwnerResponse
	owner, err := uc.userRepo.FindByID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor owner: %w", err)
	}
	if owner != nil {
		ownerResp = domain.AdminVendorOwnerResponse{
			ID:       owner.ID,
			Email:    owner.Email,
			FullName: owner.FullName,
			Phone:    owner.Phone,
			Status:   owner.Status,
		}
	}

	// Fetch bank account.
	var bankAccountResp *domain.AdminVendorBankAccountResponse
	bankAccount, err := uc.vendorRepo.FindBankAccountByVendorID(ctx, vendor.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bank account: %w", err)
	}
	if bankAccount != nil {
		bankAccountResp = &domain.AdminVendorBankAccountResponse{
			ID:                 bankAccount.ID,
			BankName:           bankAccount.BankName,
			AccountNumber:      bankAccount.AccountNumber,
			AccountHolderName:  bankAccount.AccountHolderName,
			VerificationStatus: bankAccount.VerificationStatus,
			RejectionReason:    bankAccount.RejectionReason,
			VerifiedAt:         bankAccount.VerifiedAt,
			CreatedAt:          bankAccount.CreatedAt,
			UpdatedAt:          bankAccount.UpdatedAt,
		}
	}

	// Fetch documents and generate presigned download URLs.
	documents, err := uc.vendorRepo.FindDocumentsByVendorID(ctx, vendor.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	docResponses := make([]domain.AdminVendorDocumentResponse, len(documents))
	for i, doc := range documents {
		docResp := domain.AdminVendorDocumentResponse{
			ID:                 doc.ID,
			DocType:            doc.DocType,
			FileURL:            doc.FileURL,
			MimeType:           doc.MimeType,
			FileSizeBytes:      doc.FileSizeBytes,
			VerificationStatus: doc.VerificationStatus,
			RejectionReason:    doc.RejectionReason,
			VerifiedAt:         doc.VerifiedAt,
			CreatedAt:          doc.CreatedAt,
			UpdatedAt:          doc.UpdatedAt,
		}

		// Only generate presigned URL if the document has been uploaded.
		if doc.UploadedBy != nil && doc.FileURL != "" {
			downloadURL, err := uc.storage.GeneratePresignedURL(ctx, doc.FileURL, presignedDownloadExpiry)
			if err != nil {
				return nil, fmt.Errorf("failed to generate presigned URL for document %s: %w", doc.DocType, err)
			}
			docResp.DownloadURL = downloadURL
		}

		docResponses[i] = docResp
	}

	return &domain.AdminVendorDetailResponse{
		ID:                    vendor.ID,
		VendorType:            vendor.VendorType,
		DisplayName:           vendor.DisplayName,
		LegalName:             vendor.LegalName,
		ResponsiblePersonName: vendor.ResponsiblePersonName,
		Description:           vendor.Description,
		Status:                vendor.Status,
		ApprovedAt:            vendor.ApprovedAt,
		CreatedAt:             vendor.CreatedAt,
		UpdatedAt:             vendor.UpdatedAt,
		Owner:                 ownerResp,
		BankAccount:           bankAccountResp,
		Documents:             docResponses,
	}, nil
}
