package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"gorm.io/gorm"
)

type vendorVoucherUseCase struct {
	voucherRepo domain.VendorVoucherRepository
	productRepo domain.ProductRepository
}

// NewVendorVoucherUseCase creates a new VendorVoucherUseCase.
func NewVendorVoucherUseCase(
	voucherRepo domain.VendorVoucherRepository,
	productRepo domain.ProductRepository,
) domain.VendorVoucherUseCase {
	return &vendorVoucherUseCase{
		voucherRepo: voucherRepo,
		productRepo: productRepo,
	}
}

func (uc *vendorVoucherUseCase) CreateVoucher(ctx context.Context, vendorID uuid.UUID, req domain.CreateVendorVoucherRequest) (*domain.VendorVoucherResponse, error) {
	code, err := normalizeVoucherCode(req.Code)
	if err != nil {
		return nil, err
	}

	if req.EndsAt.Before(req.StartsAt) {
		return nil, ErrInvalidVoucherPeriod
	}

	if req.QuotaTotal < 1 {
		return nil, ErrInvalidVoucherQuota
	}

	existingByCode, err := uc.voucherRepo.FindByCode(ctx, vendorID, code)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing voucher code: %w", err)
	}
	if existingByCode != nil {
		return nil, ErrVoucherCodeExists
	}

	productIDs, err := parseAndValidateVoucherProductIDs(req.ProductIDs)
	if err != nil {
		return nil, err
	}
	if err := uc.validateProductOwnership(ctx, vendorID, productIDs); err != nil {
		return nil, err
	}

	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	voucher := &domain.VendorVoucher{
		ID:             uuid.New(),
		VendorID:       vendorID,
		Code:           code,
		Name:           strings.TrimSpace(req.Name),
		Description:    req.Description,
		DiscountAmount: req.DiscountAmount,
		QuotaTotal:     req.QuotaTotal,
		QuotaUsed:      0,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
		IsActive:       isActive,
		ProductIDs:     productIDs,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := uc.voucherRepo.Create(ctx, voucher); err != nil {
		if isVoucherCodeConflict(err) {
			return nil, ErrVoucherCodeExists
		}
		return nil, fmt.Errorf("failed to create voucher: %w", err)
	}

	return toVendorVoucherResponse(voucher, now), nil
}

func (uc *vendorVoucherUseCase) ListVouchers(ctx context.Context, vendorID uuid.UUID, params domain.VendorVoucherListParams) ([]domain.VendorVoucherResponse, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}
	params.Search = strings.TrimSpace(params.Search)

	items, total, err := uc.voucherRepo.ListByVendor(ctx, vendorID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list vouchers: %w", err)
	}

	now := time.Now()
	resp := make([]domain.VendorVoucherResponse, 0, len(items))
	for i := range items {
		resp = append(resp, *toVendorVoucherResponse(&items[i], now))
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return resp, meta, nil
}

func (uc *vendorVoucherUseCase) GetVoucher(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) (*domain.VendorVoucherResponse, error) {
	voucher, err := uc.voucherRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load voucher: %w", err)
	}
	if voucher == nil || voucher.VendorID != vendorID {
		return nil, ErrVoucherNotFound
	}

	return toVendorVoucherResponse(voucher, time.Now()), nil
}

func (uc *vendorVoucherUseCase) UpdateVoucher(ctx context.Context, vendorID uuid.UUID, id uuid.UUID, req domain.UpdateVendorVoucherRequest) (*domain.VendorVoucherResponse, error) {
	voucher, err := uc.voucherRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load voucher: %w", err)
	}
	if voucher == nil || voucher.VendorID != vendorID {
		return nil, ErrVoucherNotFound
	}

	if req.Code != nil {
		code, err := normalizeVoucherCode(*req.Code)
		if err != nil {
			return nil, err
		}
		if code != voucher.Code {
			existingByCode, err := uc.voucherRepo.FindByCode(ctx, vendorID, code)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing voucher code: %w", err)
			}
			if existingByCode != nil && existingByCode.ID != voucher.ID {
				return nil, ErrVoucherCodeExists
			}
		}
		voucher.Code = code
	}

	if req.Name != nil {
		voucher.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		voucher.Description = req.Description
	}
	if req.DiscountAmount != nil {
		voucher.DiscountAmount = *req.DiscountAmount
	}
	if req.QuotaTotal != nil {
		if *req.QuotaTotal < voucher.QuotaUsed {
			return nil, ErrInvalidVoucherQuota
		}
		voucher.QuotaTotal = *req.QuotaTotal
	}
	if req.StartsAt != nil {
		voucher.StartsAt = *req.StartsAt
	}
	if req.EndsAt != nil {
		voucher.EndsAt = *req.EndsAt
	}
	if voucher.EndsAt.Before(voucher.StartsAt) {
		return nil, ErrInvalidVoucherPeriod
	}
	if req.IsActive != nil {
		voucher.IsActive = *req.IsActive
	}

	replaceProducts := false
	if req.ProductIDs != nil {
		productIDs, err := parseAndValidateVoucherProductIDs(req.ProductIDs)
		if err != nil {
			return nil, err
		}
		if err := uc.validateProductOwnership(ctx, vendorID, productIDs); err != nil {
			return nil, err
		}
		voucher.ProductIDs = productIDs
		replaceProducts = true
	}

	voucher.UpdatedAt = time.Now()

	if err := uc.voucherRepo.Update(ctx, voucher, replaceProducts); err != nil {
		if isVoucherCodeConflict(err) {
			return nil, ErrVoucherCodeExists
		}
		return nil, fmt.Errorf("failed to update voucher: %w", err)
	}

	return toVendorVoucherResponse(voucher, time.Now()), nil
}

func (uc *vendorVoucherUseCase) DeleteVoucher(ctx context.Context, vendorID uuid.UUID, id uuid.UUID) error {
	if err := uc.voucherRepo.Delete(ctx, id, vendorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVoucherNotFound
		}
		return fmt.Errorf("failed to delete voucher: %w", err)
	}
	return nil
}

func (uc *vendorVoucherUseCase) validateProductOwnership(ctx context.Context, vendorID uuid.UUID, productIDs []uuid.UUID) error {
	for _, productID := range productIDs {
		product, err := uc.productRepo.FindByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("failed to check product %s: %w", productID.String(), err)
		}
		if product == nil {
			return ErrProductNotFound
		}
		if product.VendorID != vendorID {
			return ErrProductNotOwned
		}
	}
	return nil
}

func parseAndValidateVoucherProductIDs(raw []string) ([]uuid.UUID, error) {
	if len(raw) == 0 {
		return nil, ErrVoucherProductsRequired
	}

	seen := make(map[uuid.UUID]struct{}, len(raw))
	result := make([]uuid.UUID, 0, len(raw))
	for _, item := range raw {
		id, err := uuid.Parse(item)
		if err != nil {
			return nil, ErrVoucherProductsRequired
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}

	if len(result) == 0 {
		return nil, ErrVoucherProductsRequired
	}

	return result, nil
}

func normalizeVoucherCode(code string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if normalized == "" {
		return "", ErrInvalidVoucherCode
	}
	return normalized, nil
}

func toVendorVoucherResponse(v *domain.VendorVoucher, now time.Time) *domain.VendorVoucherResponse {
	remaining := v.QuotaTotal - v.QuotaUsed
	if remaining < 0 {
		remaining = 0
	}
	currentlyActive := v.IsActive &&
		!now.Before(v.StartsAt) &&
		!now.After(v.EndsAt) &&
		v.QuotaUsed < v.QuotaTotal

	return &domain.VendorVoucherResponse{
		ID:                v.ID,
		VendorID:          v.VendorID,
		Code:              v.Code,
		Name:              v.Name,
		Description:       v.Description,
		DiscountAmount:    v.DiscountAmount,
		QuotaTotal:        v.QuotaTotal,
		QuotaUsed:         v.QuotaUsed,
		QuotaRemaining:    remaining,
		StartsAt:          v.StartsAt,
		EndsAt:            v.EndsAt,
		IsActive:          v.IsActive,
		IsCurrentlyActive: currentlyActive,
		ProductIDs:        v.ProductIDs,
		CreatedAt:         v.CreatedAt,
		UpdatedAt:         v.UpdatedAt,
	}
}

func isVoucherCodeConflict(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "uq_vendor_vouchers_vendor_code")
}

func (uc *vendorVoucherUseCase) GetSummary(ctx context.Context, vendorID uuid.UUID) (*domain.VendorVoucherSummary, error) {
	summary, err := uc.voucherRepo.GetSummary(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher summary: %w", err)
	}
	return summary, nil
}
