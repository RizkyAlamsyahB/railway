package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// VendorVoucher represents voucher discount data owned by vendor.
// Voucher applies per item (restricted by product list).
type VendorVoucher struct {
	ID             uuid.UUID   `json:"id"`
	VendorID       uuid.UUID   `json:"vendor_id"`
	Code           string      `json:"code"`
	Name           string      `json:"name"`
	Description    *string     `json:"description,omitempty"`
	DiscountAmount float64     `json:"discount_amount"`
	QuotaTotal     int         `json:"quota_total"`
	QuotaUsed      int         `json:"quota_used"`
	StartsAt       time.Time   `json:"starts_at"`
	EndsAt         time.Time   `json:"ends_at"`
	IsActive       bool        `json:"is_active"`
	ProductIDs     []uuid.UUID `json:"product_ids,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// VendorVoucherListParams holds query params for vendor voucher listing.
type VendorVoucherListParams struct {
	Page      int
	Limit     int
	Search    string
	IsActive  *bool
	ProductID *uuid.UUID
}

// CreateVendorVoucherRequest is payload to create a voucher.
type CreateVendorVoucherRequest struct {
	Code           string    `json:"code" binding:"required,max=50"`
	Name           string    `json:"name" binding:"required,max=120"`
	Description    *string   `json:"description" binding:"omitempty,max=500"`
	DiscountAmount float64   `json:"discount_amount" binding:"required,gt=0"`
	QuotaTotal     int       `json:"quota_total" binding:"required,min=1"`
	StartsAt       time.Time `json:"starts_at" binding:"required"`
	EndsAt         time.Time `json:"ends_at" binding:"required"`
	IsActive       *bool     `json:"is_active,omitempty"`
	ProductIDs     []string  `json:"product_ids" binding:"required,min=1,dive,uuid"`
}

// UpdateVendorVoucherRequest is payload to update a voucher.
// Field is optional and applied only when present.
type UpdateVendorVoucherRequest struct {
	Code           *string    `json:"code,omitempty" binding:"omitempty,max=50"`
	Name           *string    `json:"name,omitempty" binding:"omitempty,max=120"`
	Description    *string    `json:"description,omitempty" binding:"omitempty,max=500"`
	DiscountAmount *float64   `json:"discount_amount,omitempty" binding:"omitempty,gt=0"`
	QuotaTotal     *int       `json:"quota_total,omitempty" binding:"omitempty,min=1"`
	StartsAt       *time.Time `json:"starts_at,omitempty" binding:"omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty" binding:"omitempty"`
	IsActive       *bool      `json:"is_active,omitempty"`
	ProductIDs     []string   `json:"product_ids,omitempty" binding:"omitempty,dive,uuid"`
}

// VendorVoucherResponse is output representation for vendor voucher CRUD.
type VendorVoucherResponse struct {
	ID                uuid.UUID   `json:"id"`
	VendorID          uuid.UUID   `json:"vendor_id"`
	Code              string      `json:"code"`
	Name              string      `json:"name"`
	Description       *string     `json:"description,omitempty"`
	DiscountAmount    float64     `json:"discount_amount"`
	QuotaTotal        int         `json:"quota_total"`
	QuotaUsed         int         `json:"quota_used"`
	QuotaRemaining    int         `json:"quota_remaining"`
	StartsAt          time.Time   `json:"starts_at"`
	EndsAt            time.Time   `json:"ends_at"`
	IsActive          bool        `json:"is_active"`
	IsCurrentlyActive bool        `json:"is_currently_active"`
	ProductIDs        []uuid.UUID `json:"product_ids"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

// VendorVoucherSummary holds aggregated voucher statistics for a vendor.
type VendorVoucherSummary struct {
	ActiveCount   int `json:"active_count"`
	UpcomingCount int `json:"upcoming_count"`
	ExpiredCount  int `json:"expired_count"`
	TotalUsage    int `json:"total_usage"`
}

// VendorVoucherRepository defines data access operations for vouchers.
type VendorVoucherRepository interface {
	Create(ctx context.Context, voucher *VendorVoucher) error
	FindByID(ctx context.Context, id uuid.UUID) (*VendorVoucher, error)
	FindByCode(ctx context.Context, vendorID uuid.UUID, code string) (*VendorVoucher, error)
	ListByVendor(ctx context.Context, vendorID uuid.UUID, params VendorVoucherListParams) ([]VendorVoucher, int64, error)
	Update(ctx context.Context, voucher *VendorVoucher, replaceProducts bool) error
	Delete(ctx context.Context, id uuid.UUID, vendorID uuid.UUID) error
	GetSummary(ctx context.Context, vendorID uuid.UUID) (*VendorVoucherSummary, error)
}

// VendorVoucherUseCase defines business operations for vendor vouchers.
type VendorVoucherUseCase interface {
	CreateVoucher(ctx context.Context, vendorID uuid.UUID, req CreateVendorVoucherRequest) (*VendorVoucherResponse, error)
	ListVouchers(ctx context.Context, vendorID uuid.UUID, params VendorVoucherListParams) ([]VendorVoucherResponse, *PaginationMeta, error)
	GetVoucher(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) (*VendorVoucherResponse, error)
	UpdateVoucher(ctx context.Context, vendorID uuid.UUID, id uuid.UUID, req UpdateVendorVoucherRequest) (*VendorVoucherResponse, error)
	DeleteVoucher(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) error
	GetSummary(ctx context.Context, vendorID uuid.UUID) (*VendorVoucherSummary, error)
}
