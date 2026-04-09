package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

const vendorProposalDocumentMaxBytes int64 = 10 * 1024 * 1024 // 10 MB

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
	locationNames, err := uc.resolveVendorProposalLocationNames(ctx, provinceID, cityID, districtID, subdistrictID)
	if err != nil {
		return nil, err
	}

	vendor.VendorType = &vendorType
	vendor.DisplayName = &storeName
	vendor.Description = &storeDescription
	vendor.ProvinceID = &provinceID
	vendor.ProvinceName = &locationNames.ProvinceName
	vendor.CityID = &cityID
	vendor.CityName = &locationNames.CityName
	vendor.DistrictID = &districtID
	vendor.DistrictName = &locationNames.DistrictName
	vendor.SubdistrictID = &subdistrictID
	vendor.SubdistrictName = &locationNames.SubdistrictName
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

	for i := range documents {
		expectedPrefix := fmt.Sprintf("vendors/%s/documents/%s/", vendorID.String(), documents[i].DocType)
		if !strings.HasPrefix(documents[i].FileURL, expectedPrefix) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidDocumentObjectKey, documents[i].DocType)
		}

		info, err := uc.storage.HeadObject(ctx, documents[i].FileURL)
		if err != nil {
			return nil, fmt.Errorf("failed to verify document object %s: %w", documents[i].DocType, err)
		}
		if info == nil {
			return nil, fmt.Errorf("%w: %s", ErrObjectNotUploaded, documents[i].DocType)
		}

		contentType := normalizeContentType(info.ContentType)
		if !isAllowedDocumentContentType(contentType) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidDocumentContent, documents[i].DocType)
		}
		if info.ContentLength <= 0 || info.ContentLength > vendorProposalDocumentMaxBytes {
			return nil, fmt.Errorf("%w: %s", ErrDocumentSizeOverflow, documents[i].DocType)
		}

		fileSize := int(info.ContentLength)
		documents[i].MimeType = &contentType
		documents[i].FileSizeBytes = &fileSize
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

type vendorProposalLocationNames struct {
	ProvinceName    string
	CityName        string
	DistrictName    string
	SubdistrictName string
}

func (uc *vendorUseCase) resolveVendorProposalLocationNames(
	ctx context.Context,
	provinceID string,
	cityID string,
	districtID string,
	subdistrictID string,
) (*vendorProposalLocationNames, error) {
	if uc.shipping == nil {
		return nil, fmt.Errorf("shipping use case is not configured")
	}

	provinces, err := uc.shipping.GetProvinces(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load provinces: %w", err)
	}
	province, ok := findLocationByID(provinces, provinceID, func(item domain.ROProvince) string { return item.ID })
	if !ok {
		return nil, fmt.Errorf("%w: province_id %q not found", ErrInvalidLocationSelection, provinceID)
	}

	cities, err := uc.shipping.GetCities(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cities: %w", err)
	}
	city, ok := findLocationByID(cities, cityID, func(item domain.ROCity) string { return item.ID })
	if !ok {
		return nil, fmt.Errorf("%w: city_id %q is invalid for province_id %q", ErrInvalidLocationSelection, cityID, provinceID)
	}

	districts, err := uc.shipping.GetDistricts(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to load districts: %w", err)
	}
	district, ok := findLocationByID(districts, districtID, func(item domain.RODistrict) string { return item.ID })
	if !ok {
		return nil, fmt.Errorf("%w: district_id %q is invalid for city_id %q", ErrInvalidLocationSelection, districtID, cityID)
	}

	subdistricts, err := uc.shipping.GetSubdistricts(ctx, districtID)
	if err != nil {
		return nil, fmt.Errorf("failed to load subdistricts: %w", err)
	}
	subdistrict, ok := findLocationByID(subdistricts, subdistrictID, func(item domain.ROSubdistrict) string { return item.ID })
	if !ok {
		return nil, fmt.Errorf("%w: subdistrict_id %q is invalid for district_id %q", ErrInvalidLocationSelection, subdistrictID, districtID)
	}

	return &vendorProposalLocationNames{
		ProvinceName:    strings.TrimSpace(province.Name),
		CityName:        strings.TrimSpace(city.Name),
		DistrictName:    strings.TrimSpace(district.Name),
		SubdistrictName: strings.TrimSpace(subdistrict.Name),
	}, nil
}

func findLocationByID[T any](items []T, targetID string, getID func(T) string) (T, bool) {
	var zero T
	for _, item := range items {
		if strings.TrimSpace(getID(item)) == targetID {
			return item, true
		}
	}
	return zero, false
}
