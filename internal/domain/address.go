package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Address Entity ---

// Address represents the addresses table (shared by customer & vendor warehouse).
type Address struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Label           *string   `json:"label,omitempty"`
	RecipientName   *string   `json:"recipient_name,omitempty"`
	Phone           *string   `json:"phone,omitempty"`
	ProvinceID      *string   `json:"province_id,omitempty"`
	ProvinceName    *string   `json:"province_name,omitempty"`
	CityID          *string   `json:"city_id,omitempty"`
	CityName        *string   `json:"city_name,omitempty"`
	DistrictID      *string   `json:"district_id,omitempty"`
	DistrictName    *string   `json:"district_name,omitempty"`
	SubdistrictID   *string   `json:"subdistrict_id,omitempty"`
	SubdistrictName *string   `json:"subdistrict_name,omitempty"`
	PostalCode      *string   `json:"postal_code,omitempty"`
	AddressLine     string    `json:"address_line"`
	Notes           *string   `json:"notes,omitempty"`
	Latitude        *float64  `json:"latitude,omitempty"`
	Longitude       *float64  `json:"longitude,omitempty"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// --- RajaOngkir Location DTOs ---

// ROProvince represents a province from RajaOngkir.
type ROProvince struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ROCity represents a city from RajaOngkir.
type ROCity struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	PostalCode string `json:"postal_code"`
}

// RODistrict represents a district (kecamatan) from RajaOngkir.
type RODistrict struct {
	ID     string `json:"id"`
	CityID string `json:"city_id"`
	Name   string `json:"name"`
}

// ROSubdistrict represents a subdistrict (kelurahan) from RajaOngkir.
type ROSubdistrict struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
	ZipCode    string `json:"zip_code,omitempty"`
}

// --- Request / Response DTOs ---

// CreateAddressRequest is the input DTO for creating a new address.
type CreateAddressRequest struct {
	Label           *string  `json:"label,omitempty" binding:"omitempty,max=40"`
	RecipientName   *string  `json:"recipient_name,omitempty" binding:"omitempty,max=120"`
	Phone           *string  `json:"phone,omitempty" binding:"omitempty,max=20"`
	ProvinceID      *string  `json:"province_id,omitempty" binding:"omitempty,max=20"`
	ProvinceName    *string  `json:"province_name,omitempty" binding:"omitempty,max=100"`
	CityID          *string  `json:"city_id,omitempty" binding:"omitempty,max=20"`
	CityName        *string  `json:"city_name,omitempty" binding:"omitempty,max=100"`
	DistrictID      *string  `json:"district_id,omitempty" binding:"omitempty,max=20"`
	DistrictName    *string  `json:"district_name,omitempty" binding:"omitempty,max=100"`
	SubdistrictID   *string  `json:"subdistrict_id,omitempty" binding:"omitempty,max=20"`
	SubdistrictName *string  `json:"subdistrict_name,omitempty" binding:"omitempty,max=100"`
	PostalCode      *string  `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	AddressLine     string   `json:"address_line" binding:"required"`
	Notes           *string  `json:"notes,omitempty"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	IsDefault       bool     `json:"is_default"`
}

// UpdateAddressRequest is the input DTO for updating an existing address.
type UpdateAddressRequest struct {
	Label           *string  `json:"label,omitempty" binding:"omitempty,max=40"`
	RecipientName   *string  `json:"recipient_name,omitempty" binding:"omitempty,max=120"`
	Phone           *string  `json:"phone,omitempty" binding:"omitempty,max=20"`
	ProvinceID      *string  `json:"province_id,omitempty" binding:"omitempty,max=20"`
	ProvinceName    *string  `json:"province_name,omitempty" binding:"omitempty,max=100"`
	CityID          *string  `json:"city_id,omitempty" binding:"omitempty,max=20"`
	CityName        *string  `json:"city_name,omitempty" binding:"omitempty,max=100"`
	DistrictID      *string  `json:"district_id,omitempty" binding:"omitempty,max=20"`
	DistrictName    *string  `json:"district_name,omitempty" binding:"omitempty,max=100"`
	SubdistrictID   *string  `json:"subdistrict_id,omitempty" binding:"omitempty,max=20"`
	SubdistrictName *string  `json:"subdistrict_name,omitempty" binding:"omitempty,max=100"`
	PostalCode      *string  `json:"postal_code,omitempty" binding:"omitempty,max=10"`
	AddressLine     *string  `json:"address_line,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
	Latitude        *float64 `json:"latitude,omitempty"`
	Longitude       *float64 `json:"longitude,omitempty"`
	IsDefault       *bool    `json:"is_default,omitempty"`
}

// AddressResponse is the output DTO for address data.
type AddressResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Label           *string   `json:"label,omitempty"`
	RecipientName   *string   `json:"recipient_name,omitempty"`
	Phone           *string   `json:"phone,omitempty"`
	ProvinceID      *string   `json:"province_id,omitempty"`
	ProvinceName    *string   `json:"province_name,omitempty"`
	CityID          *string   `json:"city_id,omitempty"`
	CityName        *string   `json:"city_name,omitempty"`
	DistrictID      *string   `json:"district_id,omitempty"`
	DistrictName    *string   `json:"district_name,omitempty"`
	SubdistrictID   *string   `json:"subdistrict_id,omitempty"`
	SubdistrictName *string   `json:"subdistrict_name,omitempty"`
	PostalCode      *string   `json:"postal_code,omitempty"`
	AddressLine     string    `json:"address_line"`
	Notes           *string   `json:"notes,omitempty"`
	Latitude        *float64  `json:"latitude,omitempty"`
	Longitude       *float64  `json:"longitude,omitempty"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// --- Shipping Cost DTOs ---

// ShippingCostOption represents one courier service option returned by RajaOngkir.
type ShippingCostOption struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int    `json:"cost"`
	ETD         string `json:"etd"`
}

// --- Interfaces ---

// RajaOngkirProvider defines the interface for RajaOngkir API integration.
type RajaOngkirProvider interface {
	GetProvinces(ctx context.Context) ([]ROProvince, error)
	GetCitiesByProvince(ctx context.Context, provinceID string) ([]ROCity, error)
	GetDistrictsByCity(ctx context.Context, cityID string) ([]RODistrict, error)
	GetSubdistrictsByDistrict(ctx context.Context, districtID string) ([]ROSubdistrict, error)

	// CalculateDomesticCost calls RajaOngkir domestic cost API.
	// origin/destination are district IDs, weight in grams, couriers colon-separated codes.
	CalculateDomesticCost(ctx context.Context, originDistrictID, destDistrictID string, weightGram int, courierCodes string) ([]ShippingCostOption, error)

	// TrackWaybill tracks a shipment by AWB number and courier code via RajaOngkir.
	TrackWaybill(ctx context.Context, awbNumber, courierCode string) (*TrackWaybillResponse, error)
}

// AddressRepository defines the interface for address data access.
type AddressRepository interface {
	Create(ctx context.Context, address *Address) error
	FindByID(ctx context.Context, id uuid.UUID) (*Address, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Address, error)
	FindDefaultByUserID(ctx context.Context, userID uuid.UUID) (*Address, error)
	Update(ctx context.Context, address *Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	ClearDefaultByUser(ctx context.Context, userID uuid.UUID) error
}

// AddressUseCase defines the interface for address business logic.
type AddressUseCase interface {
	Create(ctx context.Context, userID uuid.UUID, req CreateAddressRequest) (*AddressResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]AddressResponse, error)
	GetByID(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*AddressResponse, error)
	Update(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req UpdateAddressRequest) (*AddressResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error
	SetDefault(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) (*AddressResponse, error)
}

// ShippingUseCase defines the interface for shipping location lookup (RajaOngkir proxy).
type ShippingUseCase interface {
	GetProvinces(ctx context.Context) ([]ROProvince, error)
	GetCities(ctx context.Context, provinceID string) ([]ROCity, error)
	GetDistricts(ctx context.Context, cityID string) ([]RODistrict, error)
	GetSubdistricts(ctx context.Context, districtID string) ([]ROSubdistrict, error)
}
