package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

func (uc *vendorUseCase) PresignSouvenirStoreProposalDocuments(
	ctx context.Context,
	vendorID uuid.UUID,
	userID uuid.UUID,
	req domain.VendorSouvenirStoreProposalPresignRequest,
) (*domain.VendorSouvenirStoreProposalPresignResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}
	if vendor.OwnerUserID != userID {
		return nil, ErrVendorNotOwned
	}
	if vendor.Status != domain.VendorStatusDraft && vendor.Status != domain.VendorStatusRejected {
		return nil, fmt.Errorf("%w: cannot presign proposal documents for vendor with status %q", ErrInvalidStatusTransition, vendor.Status)
	}

	responsiblePersonKTPContentType := normalizeContentType(req.ResponsiblePersonKTP)
	nibContentType := normalizeContentType(req.NIB)
	halalCertificateContentType := normalizeContentType(req.HalalCertificate)

	if !isAllowedDocumentContentType(responsiblePersonKTPContentType) ||
		!isAllowedDocumentContentType(nibContentType) ||
		!isAllowedDocumentContentType(halalCertificateContentType) {
		return nil, ErrInvalidDocumentContent
	}

	responsiblePersonKTPObjectID := fmt.Sprintf("vendors/%s/documents/%s/%s", vendorID.String(), domain.VendorDocumentTypeOwnerDocumentID, uuid.NewString())
	responsiblePersonKTPPresignedURL, err := uc.storage.GeneratePresignedUploadURL(ctx, responsiblePersonKTPObjectID, responsiblePersonKTPContentType, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL for responsible person ktp: %w", err)
	}

	nibObjectID := fmt.Sprintf("vendors/%s/documents/%s/%s", vendorID.String(), domain.VendorDocumentTypeBusinessNIB, uuid.NewString())
	nibPresignedURL, err := uc.storage.GeneratePresignedUploadURL(ctx, nibObjectID, nibContentType, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL for nib: %w", err)
	}

	halalCertificateObjectID := fmt.Sprintf("vendors/%s/documents/%s/%s", vendorID.String(), domain.VendorDocumentTypeHalalCertificate, uuid.NewString())
	halalCertificatePresignedURL, err := uc.storage.GeneratePresignedUploadURL(ctx, halalCertificateObjectID, halalCertificateContentType, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned upload URL for halal certificate: %w", err)
	}

	return &domain.VendorSouvenirStoreProposalPresignResponse{
		ResponsiblePersonKTP: domain.VendorSouvenirStoreProposalPresignItem{
			ObjectID:     responsiblePersonKTPObjectID,
			PresignedURL: responsiblePersonKTPPresignedURL,
		},
		NIB: domain.VendorSouvenirStoreProposalPresignItem{
			ObjectID:     nibObjectID,
			PresignedURL: nibPresignedURL,
		},
		HalalCertificate: domain.VendorSouvenirStoreProposalPresignItem{
			ObjectID:     halalCertificateObjectID,
			PresignedURL: halalCertificatePresignedURL,
		},
	}, nil
}

func (uc *vendorUseCase) SubmitSouvenirStoreProposal(
	ctx context.Context,
	vendorID uuid.UUID,
	userID uuid.UUID,
	req domain.VendorSouvenirStoreProposalRequest,
) (*domain.VendorSouvenirStoreProposalResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}
	if vendor.OwnerUserID != userID {
		return nil, ErrVendorNotOwned
	}
	if vendor.Status != domain.VendorStatusDraft && vendor.Status != domain.VendorStatusRejected {
		return nil, fmt.Errorf("%w: cannot submit vendor with status %q", ErrInvalidStatusTransition, vendor.Status)
	}

	owner, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor owner: %w", err)
	}
	if owner == nil {
		return nil, ErrUserNotFound
	}

	responsibleEmail := normalizeVendorEmail(req.ResponsiblePerson.Email)
	if !strings.EqualFold(owner.Email, responsibleEmail) {
		return nil, ErrVendorResponsibleEmail
	}

	responsiblePhone := strings.TrimSpace(req.ResponsiblePerson.Phone)
	if responsiblePhone != "" {
		existingUser, err := uc.userRepo.FindByPhone(ctx, responsiblePhone)
		if err != nil {
			return nil, fmt.Errorf("failed to check phone uniqueness: %w", err)
		}
		if existingUser != nil && existingUser.ID != owner.ID {
			return nil, ErrPhoneAlreadyRegistered
		}
	}

	now := time.Now()
	storeName := strings.TrimSpace(req.StoreName)
	storeDescription := strings.TrimSpace(req.StoreDescription)
	vendorType := domain.VendorTypeSouvenirStore
	provinceID := strings.TrimSpace(req.Address.ProvinceID)
	cityID := strings.TrimSpace(req.Address.CityID)
	districtID := strings.TrimSpace(req.Address.DistrictID)
	subdistrictID := strings.TrimSpace(req.Address.SubdistrictID)
	postalCode := strings.TrimSpace(req.Address.PostalCode)
	addressLine := strings.TrimSpace(req.Address.AddressLine)
	responsibleName := strings.TrimSpace(req.ResponsiblePerson.Name)
	nik := strings.TrimSpace(req.ResponsiblePerson.NIK)

	vendor.VendorType = &vendorType
	vendor.DisplayName = &storeName
	vendor.Description = &storeDescription
	vendor.ProvinceID = &provinceID
	vendor.CityID = &cityID
	vendor.DistrictID = &districtID
	vendor.SubdistrictID = &subdistrictID
	vendor.PostalCode = &postalCode
	vendor.AddressLine = &addressLine
	vendor.Status = domain.VendorStatusSubmitted
	vendor.ApprovedBy = nil
	vendor.ApprovedAt = nil
	vendor.StatusReason = nil
	vendor.UpdatedAt = now

	owner.FullName = responsibleName
	owner.Phone = &responsiblePhone
	owner.UpdatedAt = now

	responsiblePerson := &domain.VendorResponsiblePerson{
		ID:        uuid.New(),
		VendorID:  vendor.ID,
		UserID:    owner.ID,
		NIK:       nik,
		CreatedAt: now,
		UpdatedAt: now,
	}

	documents := []domain.VendorDocument{
		{
			ID:         uuid.New(),
			VendorID:   vendor.ID,
			DocType:    domain.VendorDocumentTypeOwnerDocumentID,
			FileURL:    strings.TrimSpace(req.ResponsiblePerson.KTPObjectID),
			UploadedBy: &owner.ID,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         uuid.New(),
			VendorID:   vendor.ID,
			DocType:    domain.VendorDocumentTypeBusinessNIB,
			FileURL:    strings.TrimSpace(req.OtherDocuments.NIBObjectID),
			UploadedBy: &owner.ID,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         uuid.New(),
			VendorID:   vendor.ID,
			DocType:    domain.VendorDocumentTypeHalalCertificate,
			FileURL:    strings.TrimSpace(req.OtherDocuments.HalalCertificateObjectID),
			UploadedBy: &owner.ID,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := uc.vendorRepo.SubmitSouvenirStoreProposal(ctx, domain.SubmitSouvenirStoreProposalInput{
		Vendor:            vendor,
		OwnerUser:         owner,
		ResponsiblePerson: responsiblePerson,
		Documents:         documents,
	}); err != nil {
		return nil, fmt.Errorf("failed to submit vendor proposal: %w", err)
	}

	return &domain.VendorSouvenirStoreProposalResponse{
		VendorID: vendor.ID,
		Status:   domain.VendorStatusSubmitted,
	}, nil
}
