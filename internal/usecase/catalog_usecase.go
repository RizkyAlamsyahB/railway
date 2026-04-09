package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type catalogUseCase struct {
	categoryRepo      domain.CategoryRepository
	productRepo       domain.ProductRepository
	vendorRepo        domain.VendorRepository
	addressRepo       domain.AddressRepository
	vendorCourierRepo domain.VendorCourierRepository
	rajaOngkir        domain.RajaOngkirProvider
	storage           domain.StorageProvider
}

// NewCatalogUseCase creates a new CatalogUseCase.
func NewCatalogUseCase(
	categoryRepo domain.CategoryRepository,
	productRepo domain.ProductRepository,
	vendorRepo domain.VendorRepository,
	addressRepo domain.AddressRepository,
	vendorCourierRepo domain.VendorCourierRepository,
	rajaOngkir domain.RajaOngkirProvider,
	storage domain.StorageProvider,
) domain.CatalogUseCase {
	return &catalogUseCase{
		categoryRepo:      categoryRepo,
		productRepo:       productRepo,
		vendorRepo:        vendorRepo,
		addressRepo:       addressRepo,
		vendorCourierRepo: vendorCourierRepo,
		rajaOngkir:        rajaOngkir,
		storage:           storage,
	}
}

func (uc *catalogUseCase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return uc.categoryRepo.ListActive(ctx)
}

func (uc *catalogUseCase) ListProducts(ctx context.Context, params domain.ProductListParams) ([]domain.ProductListItem, *domain.PaginationMeta, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}

	params.Sort = strings.ToLower(strings.TrimSpace(params.Sort))
	if params.Sort != "cheapest" && params.Sort != "bestseller" && params.Sort != "expensive" {
		params.Sort = "newest"
	}
	params.Search = strings.TrimSpace(params.Search)

	items, total, err := uc.productRepo.ListPublishedForCustomer(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list products: %w", err)
	}

	for i := range items {
		if items[i].ImageURL != "" && !isAbsoluteURL(items[i].ImageURL) {
			items[i].ImageURL = uc.storage.GetURL(items[i].ImageURL)
		}
	}

	meta := &domain.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalItems: total,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}

	return items, meta, nil
}

func (uc *catalogUseCase) GetProductDetail(ctx context.Context, id uuid.UUID) (*domain.ProductDetailResponse, error) {
	detail, err := uc.productRepo.GetPublishedDetailForCustomer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product detail: %w", err)
	}
	if detail == nil {
		return nil, ErrProductNotFound
	}

	for i := range detail.Images {
		if !isAbsoluteURL(detail.Images[i].URL) {
			detail.Images[i].URL = uc.storage.GetURL(detail.Images[i].URL)
		}
	}

	return detail, nil
}

func isAbsoluteURL(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://")
}

const catalogDefaultWeightGram = 500

func (uc *catalogUseCase) GetShippingEstimate(ctx context.Context, productID uuid.UUID, userID *uuid.UUID) (*domain.ProductShippingEstimateResponse, error) {
	// No logged-in user → no estimate.
	if userID == nil {
		return nil, nil
	}

	// Get customer default address.
	custAddr, err := uc.addressRepo.FindDefaultByUserID(ctx, *userID)
	if err != nil {
		return nil, fmt.Errorf("find customer address: %w", err)
	}
	if custAddr == nil || custAddr.DistrictID == nil || *custAddr.DistrictID == "" {
		return nil, nil
	}

	// Get product to find vendor.
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}
	if product == nil || product.Status != domain.ProductStatusPublished {
		return nil, nil
	}

	// Get vendor and warehouse address.
	vendor, err := uc.vendorRepo.FindByID(ctx, product.VendorID)
	if err != nil {
		return nil, fmt.Errorf("find vendor: %w", err)
	}
	if vendor == nil || vendor.Status != "active" {
		return nil, nil
	}

	warehouseAddr, err := uc.addressRepo.FindDefaultByUserID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("find vendor warehouse: %w", err)
	}
	if warehouseAddr == nil || warehouseAddr.DistrictID == nil || *warehouseAddr.DistrictID == "" {
		return nil, nil
	}

	// Get vendor couriers.
	vendorCouriers, err := uc.vendorCourierRepo.FindByVendorID(ctx, product.VendorID)
	if err != nil {
		return nil, fmt.Errorf("find vendor couriers: %w", err)
	}
	if len(vendorCouriers) == 0 {
		return nil, nil
	}

	// Get default variant weight.
	variants, err := uc.productRepo.FindVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("find variants: %w", err)
	}
	weightGram := catalogDefaultWeightGram
	for _, v := range variants {
		if v.IsDefault && v.IsActive {
			if v.WeightGram != nil && *v.WeightGram > 0 {
				weightGram = *v.WeightGram
			}
			break
		}
	}

	// Build courier codes.
	var codes []string
	for _, vc := range vendorCouriers {
		if vc.IsActive {
			codes = append(codes, vc.CourierCode)
		}
	}
	if len(codes) == 0 {
		return nil, nil
	}

	// Call RajaOngkir.
	options, err := uc.rajaOngkir.CalculateDomesticCost(ctx, *warehouseAddr.DistrictID, *custAddr.DistrictID, weightGram, strings.Join(codes, ":"))
	if err != nil {
		return nil, nil // silently return no estimate on API error
	}
	if len(options) == 0 {
		return nil, nil
	}

	// Find cheapest option.
	cheapest := options[0]
	for _, opt := range options[1:] {
		if opt.Cost < cheapest.Cost {
			cheapest = opt
		}
	}

	return &domain.ProductShippingEstimateResponse{
		CourierName: cheapest.Name,
		CourierCode: cheapest.Code,
		Service:     cheapest.Service,
		Cost:        cheapest.Cost,
		ETD:         cheapest.ETD,
	}, nil
}
