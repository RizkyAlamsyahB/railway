package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type adminVendorUseCase struct {
	vendorRepo  domain.VendorRepository
	userRepo    domain.UserRepository
	storage     domain.StorageProvider
	xenPlatform domain.XenPlatformProvider
}

// NewAdminVendorUseCase creates a new AdminVendorUseCase.
func NewAdminVendorUseCase(
	vendorRepo domain.VendorRepository,
	userRepo domain.UserRepository,
	storage domain.StorageProvider,
	xenPlatform domain.XenPlatformProvider,
) domain.AdminVendorUseCase {
	return &adminVendorUseCase{
		vendorRepo:  vendorRepo,
		userRepo:    userRepo,
		storage:     storage,
		xenPlatform: xenPlatform,
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
			StatusReason:          v.StatusReason,
			XenditAccountID:       v.XenditAccountID,
			RegisteredAddress:     v.RegisteredAddress,
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
			downloadURL, err := uc.storage.GeneratePresignedURL(ctx, doc.FileURL, PresignedDownloadExpiry)
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
		RegisteredAddress:     vendor.RegisteredAddress,
		Status:                vendor.Status,
		StatusReason:          vendor.StatusReason,
		ApprovedAt:            vendor.ApprovedAt,
		XenditAccountID:       vendor.XenditAccountID,
		CreatedAt:             vendor.CreatedAt,
		UpdatedAt:             vendor.UpdatedAt,
		Owner:                 ownerResp,
		BankAccount:           bankAccountResp,
		Documents:             docResponses,
	}, nil
}

func (uc *adminVendorUseCase) Approve(ctx context.Context, vendorID uuid.UUID, adminID uuid.UUID) (*domain.AdminVendorActionResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	if vendor.Status != domain.VendorStatusSubmitted && vendor.Status != domain.VendorStatusRejected {
		return nil, fmt.Errorf("%w: cannot approve vendor with status %q", ErrInvalidStatusTransition, vendor.Status)
	}

	// Fetch owner early — we need the email for Xendit account creation.
	owner, err := uc.userRepo.FindByID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor owner: %w", err)
	}
	if owner == nil {
		return nil, fmt.Errorf("failed to find vendor owner: owner not found")
	}

	// Create OWNED sub-account on Xendit XenPlatform.
	xenAccount, err := uc.xenPlatform.CreateAccount(ctx, domain.XenPlatformCreateAccountRequest{
		Email: owner.Email,
		Type:  "OWNED",
		PublicProfile: &domain.XenPlatformPublicProfile{
			BusinessName: vendor.DisplayName,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrXenditAccountCreation, err)
	}

	now := time.Now()
	adminIDStr := adminID.String()
	updates := map[string]interface{}{
		"status":            domain.VendorStatusActive,
		"approved_by":       adminIDStr,
		"approved_at":       now,
		"status_reason":     nil,
		"xendit_account_id": xenAccount.ID,
		"updated_at":        now,
	}

	if err := uc.vendorRepo.UpdateStatus(ctx, vendorID, updates); err != nil {
		return nil, fmt.Errorf("failed to approve vendor: %w", err)
	}

	// Update owner user status to active.
	if owner.Status != domain.UserStatusActive {
		owner.Status = domain.UserStatusActive
		owner.UpdatedAt = now
		if err := uc.userRepo.Update(ctx, owner); err != nil {
			return nil, fmt.Errorf("failed to activate vendor owner: %w", err)
		}
	}

	return &domain.AdminVendorActionResponse{
		VendorID: vendorID,
		Status:   domain.VendorStatusActive,
		Message:  "vendor approved successfully",
	}, nil
}

func (uc *adminVendorUseCase) Reject(ctx context.Context, vendorID uuid.UUID, reason string) (*domain.AdminVendorActionResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	if vendor.Status != domain.VendorStatusSubmitted {
		return nil, fmt.Errorf("%w: cannot reject vendor with status %q", ErrInvalidStatusTransition, vendor.Status)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        domain.VendorStatusRejected,
		"status_reason": reason,
		"updated_at":    now,
	}

	if err := uc.vendorRepo.UpdateStatus(ctx, vendorID, updates); err != nil {
		return nil, fmt.Errorf("failed to reject vendor: %w", err)
	}

	return &domain.AdminVendorActionResponse{
		VendorID: vendorID,
		Status:   domain.VendorStatusRejected,
		Message:  "vendor rejected successfully",
	}, nil
}

func (uc *adminVendorUseCase) Block(ctx context.Context, vendorID uuid.UUID, reason string) (*domain.AdminVendorActionResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	if vendor.Status != domain.VendorStatusActive && vendor.Status != domain.VendorStatusSubmitted {
		return nil, fmt.Errorf("%w: cannot block vendor with status %q", ErrInvalidStatusTransition, vendor.Status)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        domain.VendorStatusBlocked,
		"status_reason": reason,
		"updated_at":    now,
	}

	if err := uc.vendorRepo.UpdateStatus(ctx, vendorID, updates); err != nil {
		return nil, fmt.Errorf("failed to block vendor: %w", err)
	}

	return &domain.AdminVendorActionResponse{
		VendorID: vendorID,
		Status:   domain.VendorStatusBlocked,
		Message:  "vendor blocked successfully",
	}, nil
}

func (uc *adminVendorUseCase) Unblock(ctx context.Context, vendorID uuid.UUID, adminID uuid.UUID) (*domain.AdminVendorActionResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	if vendor.Status != domain.VendorStatusBlocked {
		return nil, fmt.Errorf("%w: cannot unblock vendor with status %q", ErrInvalidStatusTransition, vendor.Status)
	}

	now := time.Now()
	adminIDStr := adminID.String()
	updates := map[string]interface{}{
		"status":        domain.VendorStatusActive,
		"approved_by":   adminIDStr,
		"approved_at":   now,
		"status_reason": nil,
		"updated_at":    now,
	}

	if err := uc.vendorRepo.UpdateStatus(ctx, vendorID, updates); err != nil {
		return nil, fmt.Errorf("failed to unblock vendor: %w", err)
	}

	// Update owner user status to active if not already.
	owner, err := uc.userRepo.FindByID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor owner: %w", err)
	}
	if owner != nil && owner.Status != domain.UserStatusActive {
		owner.Status = domain.UserStatusActive
		owner.UpdatedAt = now
		if err := uc.userRepo.Update(ctx, owner); err != nil {
			return nil, fmt.Errorf("failed to activate vendor owner: %w", err)
		}
	}

	return &domain.AdminVendorActionResponse{
		VendorID: vendorID,
		Status:   domain.VendorStatusActive,
		Message:  "vendor unblocked successfully",
	}, nil
}
