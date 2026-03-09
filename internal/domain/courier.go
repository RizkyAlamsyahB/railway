package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Courier Entity ---

// Courier represents a shipping courier in the couriers master table.
type Courier struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	LogoURL  string `json:"logo_url,omitempty"`
	IsActive bool   `json:"is_active"`
}

// VendorCourier represents the vendor_couriers join table row.
type VendorCourier struct {
	ID        uuid.UUID `json:"id"`
	VendorID  uuid.UUID `json:"vendor_id"`
	CourierID int       `json:"courier_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// --- DTOs ---

// CourierResponse is the public representation of a courier.
type CourierResponse struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	LogoURL  string `json:"logo_url,omitempty"`
	IsActive bool   `json:"is_active"`
}

// VendorCourierResponse is the response for a vendor's selected courier.
type VendorCourierResponse struct {
	ID          uuid.UUID `json:"id"`
	CourierID   int       `json:"courier_id"`
	CourierCode string    `json:"courier_code"`
	CourierName string    `json:"courier_name"`
	LogoURL     string    `json:"logo_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// SetVendorCouriersRequest is the input to set (replace) a vendor's courier list.
type SetVendorCouriersRequest struct {
	CourierIDs []int `json:"courier_ids" binding:"required,min=1"`
}

// --- Interfaces ---

// CourierRepository handles courier data access.
type CourierRepository interface {
	// FindAll returns all couriers (optionally only active ones).
	FindAll(ctx context.Context, onlyActive bool) ([]Courier, error)
	// FindByID returns a courier by ID.
	FindByID(ctx context.Context, id int) (*Courier, error)
}

// VendorCourierRepository handles vendor_couriers data access.
type VendorCourierRepository interface {
	// FindByVendorID returns all courier selections for a vendor (joined with courier).
	FindByVendorID(ctx context.Context, vendorID uuid.UUID) ([]VendorCourierResponse, error)
	// SetCouriers replaces a vendor's courier list (delete old + insert new) in a transaction.
	SetCouriers(ctx context.Context, vendorID uuid.UUID, courierIDs []int) ([]VendorCourierResponse, error)
	// DeleteByVendorAndCourier removes a specific courier from a vendor.
	DeleteByVendorAndCourier(ctx context.Context, vendorID uuid.UUID, courierID int) error
}

// VendorCourierUseCase handles business logic for vendor courier management.
type VendorCourierUseCase interface {
	// ListCouriers returns all available couriers.
	ListCouriers(ctx context.Context) ([]CourierResponse, error)
	// GetVendorCouriers returns a vendor's selected couriers.
	GetVendorCouriers(ctx context.Context, vendorID uuid.UUID) ([]VendorCourierResponse, error)
	// SetVendorCouriers replaces a vendor's courier list.
	SetVendorCouriers(ctx context.Context, vendorID uuid.UUID, courierIDs []int) ([]VendorCourierResponse, error)
	// RemoveVendorCourier removes one courier from a vendor.
	RemoveVendorCourier(ctx context.Context, vendorID uuid.UUID, courierID int) error
}
