package usecase

import (
	"context"
	"fmt"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type shippingUseCase struct {
	rajaOngkir domain.RajaOngkirProvider
}

// NewShippingUseCase creates a new ShippingUseCase for location lookups.
func NewShippingUseCase(rajaOngkir domain.RajaOngkirProvider) domain.ShippingUseCase {
	return &shippingUseCase{rajaOngkir: rajaOngkir}
}

func (uc *shippingUseCase) GetProvinces(ctx context.Context) ([]domain.ROProvince, error) {
	provinces, err := uc.rajaOngkir.GetProvinces(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRajaOngkirFailed, err)
	}
	return provinces, nil
}

func (uc *shippingUseCase) GetCities(ctx context.Context, provinceID string) ([]domain.ROCity, error) {
	cities, err := uc.rajaOngkir.GetCitiesByProvince(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRajaOngkirFailed, err)
	}
	return cities, nil
}

func (uc *shippingUseCase) GetDistricts(ctx context.Context, cityID string) ([]domain.RODistrict, error) {
	districts, err := uc.rajaOngkir.GetDistrictsByCity(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRajaOngkirFailed, err)
	}
	return districts, nil
}

func (uc *shippingUseCase) GetSubdistricts(ctx context.Context, districtID string) ([]domain.ROSubdistrict, error) {
	subdistricts, err := uc.rajaOngkir.GetSubdistrictsByDistrict(ctx, districtID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRajaOngkirFailed, err)
	}
	return subdistricts, nil
}
