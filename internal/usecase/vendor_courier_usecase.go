package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type vendorCourierUseCase struct {
	courierRepo       domain.CourierRepository
	vendorCourierRepo domain.VendorCourierRepository
}

// NewVendorCourierUseCase creates a new VendorCourierUseCase.
func NewVendorCourierUseCase(
	courierRepo domain.CourierRepository,
	vendorCourierRepo domain.VendorCourierRepository,
) domain.VendorCourierUseCase {
	return &vendorCourierUseCase{
		courierRepo:       courierRepo,
		vendorCourierRepo: vendorCourierRepo,
	}
}

// ListCouriers returns all active couriers.
func (uc *vendorCourierUseCase) ListCouriers(ctx context.Context) ([]domain.CourierResponse, error) {
	couriers, err := uc.courierRepo.FindAll(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to list couriers: %w", err)
	}
	resp := make([]domain.CourierResponse, len(couriers))
	for i, c := range couriers {
		resp[i] = domain.CourierResponse{
			ID:       c.ID,
			Code:     c.Code,
			Name:     c.Name,
			LogoURL:  c.LogoURL,
			IsActive: c.IsActive,
		}
	}
	return resp, nil
}

// GetVendorCouriers returns a vendor's selected couriers.
func (uc *vendorCourierUseCase) GetVendorCouriers(ctx context.Context, vendorID uuid.UUID) ([]domain.VendorCourierResponse, error) {
	results, err := uc.vendorCourierRepo.FindByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor couriers: %w", err)
	}
	return results, nil
}

// SetVendorCouriers replaces a vendor's courier list. Validates all courier IDs exist and are active.
func (uc *vendorCourierUseCase) SetVendorCouriers(ctx context.Context, vendorID uuid.UUID, courierIDs []int) ([]domain.VendorCourierResponse, error) {
	// Validate each courier ID
	for _, id := range courierIDs {
		c, err := uc.courierRepo.FindByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to verify courier %d: %w", id, err)
		}
		if c == nil {
			return nil, ErrCourierNotFound
		}
		if !c.IsActive {
			return nil, ErrCourierNotActive
		}
	}

	results, err := uc.vendorCourierRepo.SetCouriers(ctx, vendorID, courierIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to set vendor couriers: %w", err)
	}
	return results, nil
}

// RemoveVendorCourier removes one courier from a vendor.
func (uc *vendorCourierUseCase) RemoveVendorCourier(ctx context.Context, vendorID uuid.UUID, courierID int) error {
	err := uc.vendorCourierRepo.DeleteByVendorAndCourier(ctx, vendorID, courierID)
	if err != nil {
		return ErrVendorCourierNotFound
	}
	return nil
}
