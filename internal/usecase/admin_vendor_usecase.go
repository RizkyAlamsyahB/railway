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
			ID:        v.ID,
			Status:    v.Status,
			StoreType: v.VendorType,
			StoreName: v.DisplayName,
			Address: buildAdminVendorListAddress(
				v.ProvinceName,
				v.CityName,
				v.DistrictName,
				v.SubdistrictName,
				v.PostalCode,
			),
			CreatedAt: v.CreatedAt,
		}

		owner, err := uc.userRepo.FindByID(ctx, v.OwnerUserID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to find owner for vendor %s: %w", v.ID, err)
		}
		if owner != nil {
			item.Email = owner.Email
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

func buildAdminVendorListAddress(
	province *string,
	city *string,
	district *string,
	subdistrict *string,
	postalCode *string,
) *domain.AdminVendorListAddress {
	if province == nil && city == nil && district == nil && subdistrict == nil && postalCode == nil {
		return nil
	}

	return &domain.AdminVendorListAddress{
		Province:    province,
		City:        city,
		District:    district,
		Subdistrict: subdistrict,
		PostalCode:  postalCode,
	}
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
	ownerEmail := ""
	var ownerName *string
	var ownerPhone *string
	var ownerEmailPtr *string
	owner, err := uc.userRepo.FindByID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor owner: %w", err)
	}
	if owner != nil {
		ownerEmail = owner.Email
		ownerEmailPtr = &owner.Email
		if owner.FullName != "" {
			ownerName = &owner.FullName
		}
		ownerPhone = owner.Phone
	}

	// Fetch responsible person data.
	responsiblePerson, err := uc.vendorRepo.FindResponsiblePersonByVendorID(ctx, vendor.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find responsible person: %w", err)
	}

	var responsibleNIK *string
	if responsiblePerson != nil && responsiblePerson.NIK != "" {
		responsibleNIK = &responsiblePerson.NIK
	}

	// Fetch documents and map them into KTP + required documents.
	documents, err := uc.vendorRepo.FindDocumentsByVendorID(ctx, vendor.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	requirementsDocuments := make([]domain.AdminVendorRequirementDocumentResponse, 0, len(documents))
	var ktp *domain.AdminVendorKTPResponse
	for _, doc := range documents {
		fileURL := doc.FileURL
		if doc.FileURL != "" {
			signedURL, err := uc.storage.GeneratePresignedURL(ctx, doc.FileURL, PresignedDownloadExpiry)
			if err != nil {
				return nil, fmt.Errorf("failed to generate presigned URL for document %s: %w", doc.DocType, err)
			}
			fileURL = signedURL
		}

		switch doc.DocType {
		case domain.VendorDocumentTypeOwnerDocumentID:
			ktp = &domain.AdminVendorKTPResponse{
				ID:            doc.ID,
				DocType:       doc.DocType,
				FileURL:       fileURL,
				MimeType:      doc.MimeType,
				FileSizeBytes: doc.FileSizeBytes,
			}
		case domain.VendorDocumentTypeBusinessNIB, domain.VendorDocumentTypeHalalCertificate:
			requirementsDocuments = append(requirementsDocuments, domain.AdminVendorRequirementDocumentResponse{
				ID:            doc.ID,
				DocType:       doc.DocType,
				FileURL:       fileURL,
				MimeType:      doc.MimeType,
				FileSizeBytes: doc.FileSizeBytes,
			})
		}
	}

	return &domain.AdminVendorDetailResponse{
		ID:               vendor.ID,
		StoreName:        vendor.DisplayName,
		StoreDescription: vendor.Description,
		StoreType:        vendor.VendorType,
		Status:           vendor.Status,
		StatusReason:     vendor.StatusReason,
		Email:            ownerEmail,
		Address: domain.AdminVendorDetailAddress{
			Province:    vendor.ProvinceName,
			City:        vendor.CityName,
			District:    vendor.DistrictName,
			Subdistrict: vendor.SubdistrictName,
			PostalCode:  vendor.PostalCode,
			AddressLine: vendor.AddressLine,
		},
		ResponsiblePerson: domain.AdminVendorResponsiblePersonResponse{
			Name:  ownerName,
			Phone: ownerPhone,
			Email: ownerEmailPtr,
			NIK:   responsibleNIK,
			KTP:   ktp,
		},
		RequirementsDocuments: requirementsDocuments,
		CreatedAt:             vendor.CreatedAt,
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
			BusinessName: vendorDisplayNameOrFallback(vendor, owner.Email),
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
