package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type addressUseCase struct {
	addressRepo domain.AddressRepository
}

// NewAddressUseCase creates a new AddressUseCase.
func NewAddressUseCase(addressRepo domain.AddressRepository) domain.AddressUseCase {
	return &addressUseCase{addressRepo: addressRepo}
}

func (uc *addressUseCase) Create(ctx context.Context, userID uuid.UUID, req domain.CreateAddressRequest) (*domain.AddressResponse, error) {
	now := time.Now()

	address := &domain.Address{
		ID:              uuid.New(),
		UserID:          userID,
		Label:           req.Label,
		RecipientName:   req.RecipientName,
		Phone:           req.Phone,
		ProvinceID:      req.ProvinceID,
		ProvinceName:    req.ProvinceName,
		CityID:          req.CityID,
		CityName:        req.CityName,
		DistrictID:      req.DistrictID,
		DistrictName:    req.DistrictName,
		SubdistrictID:   req.SubdistrictID,
		SubdistrictName: req.SubdistrictName,
		PostalCode:      req.PostalCode,
		AddressLine:     req.AddressLine,
		Notes:           req.Notes,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		IsDefault:       req.IsDefault,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// If this is the default address, clear other defaults first
	if address.IsDefault {
		if err := uc.addressRepo.ClearDefaultByUser(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to clear default addresses: %w", err)
		}
	}

	if err := uc.addressRepo.Create(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return toAddressResponse(address), nil
}

func (uc *addressUseCase) List(ctx context.Context, userID uuid.UUID) ([]domain.AddressResponse, error) {
	addresses, err := uc.addressRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}

	result := make([]domain.AddressResponse, len(addresses))
	for i := range addresses {
		result[i] = *toAddressResponse(&addresses[i])
	}
	return result, nil
}

func (uc *addressUseCase) GetByID(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*domain.AddressResponse, error) {
	address, err := uc.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	if address == nil || address.UserID != userID {
		return nil, ErrAddressNotFound
	}
	return toAddressResponse(address), nil
}

func (uc *addressUseCase) Update(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req domain.UpdateAddressRequest) (*domain.AddressResponse, error) {
	address, err := uc.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	if address == nil || address.UserID != userID {
		return nil, ErrAddressNotFound
	}

	// Apply partial updates
	if req.Label != nil {
		address.Label = req.Label
	}
	if req.RecipientName != nil {
		address.RecipientName = req.RecipientName
	}
	if req.Phone != nil {
		address.Phone = req.Phone
	}
	if req.ProvinceID != nil {
		address.ProvinceID = req.ProvinceID
	}
	if req.ProvinceName != nil {
		address.ProvinceName = req.ProvinceName
	}
	if req.CityID != nil {
		address.CityID = req.CityID
	}
	if req.CityName != nil {
		address.CityName = req.CityName
	}
	if req.DistrictID != nil {
		address.DistrictID = req.DistrictID
	}
	if req.DistrictName != nil {
		address.DistrictName = req.DistrictName
	}
	if req.SubdistrictID != nil {
		address.SubdistrictID = req.SubdistrictID
	}
	if req.SubdistrictName != nil {
		address.SubdistrictName = req.SubdistrictName
	}
	if req.PostalCode != nil {
		address.PostalCode = req.PostalCode
	}
	if req.AddressLine != nil {
		address.AddressLine = *req.AddressLine
	}
	if req.Notes != nil {
		address.Notes = req.Notes
	}
	if req.Latitude != nil {
		address.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		address.Longitude = req.Longitude
	}
	if req.IsDefault != nil {
		if *req.IsDefault {
			if err := uc.addressRepo.ClearDefaultByUser(ctx, userID); err != nil {
				return nil, fmt.Errorf("failed to clear default addresses: %w", err)
			}
		}
		address.IsDefault = *req.IsDefault
	}

	address.UpdatedAt = time.Now()

	if err := uc.addressRepo.Update(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return toAddressResponse(address), nil
}

func (uc *addressUseCase) Delete(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error {
	address, err := uc.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}
	if address == nil || address.UserID != userID {
		return ErrAddressNotFound
	}

	if err := uc.addressRepo.Delete(ctx, addressID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}

func (uc *addressUseCase) SetDefault(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*domain.AddressResponse, error) {
	address, err := uc.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	if address == nil || address.UserID != userID {
		return nil, ErrAddressNotFound
	}

	if err := uc.addressRepo.ClearDefaultByUser(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to clear default addresses: %w", err)
	}

	address.IsDefault = true
	address.UpdatedAt = time.Now()

	if err := uc.addressRepo.Update(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return toAddressResponse(address), nil
}

// toAddressResponse maps domain.Address to domain.AddressResponse.
func toAddressResponse(a *domain.Address) *domain.AddressResponse {
	return &domain.AddressResponse{
		ID:              a.ID,
		UserID:          a.UserID,
		Label:           a.Label,
		RecipientName:   a.RecipientName,
		Phone:           a.Phone,
		ProvinceID:      a.ProvinceID,
		ProvinceName:    a.ProvinceName,
		CityID:          a.CityID,
		CityName:        a.CityName,
		DistrictID:      a.DistrictID,
		DistrictName:    a.DistrictName,
		SubdistrictID:   a.SubdistrictID,
		SubdistrictName: a.SubdistrictName,
		PostalCode:      a.PostalCode,
		AddressLine:     a.AddressLine,
		Notes:           a.Notes,
		Latitude:        a.Latitude,
		Longitude:       a.Longitude,
		IsDefault:       a.IsDefault,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}
