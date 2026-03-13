package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

const (
	platformAdminFee = 5000                              // Rp 5,000
	platformAppFee   = 1000                              // Rp 1,000
	platformFeeTotal = platformAdminFee + platformAppFee // Rp 6,000
	invoiceDuration  = 86400                             // 24 hours in seconds
)

type checkoutUseCase struct {
	cartRepo          domain.CartRepository
	productRepo       domain.ProductRepository
	vendorRepo        domain.VendorRepository
	orderRepo         domain.OrderRepository
	paymentRepo       domain.PaymentRepository
	ledgerRepo        domain.LedgerRepository
	userRepo          domain.UserRepository
	addressRepo       domain.AddressRepository
	vendorCourierRepo domain.VendorCourierRepository
	shipmentRepo      domain.ShipmentRepository
	rajaOngkir        domain.RajaOngkirProvider
	storage           domain.StorageProvider
	xenditInvoice     domain.XenditInvoiceProvider
	frontendURL       string
	webhookURL        string
}

// NewCheckoutUseCase creates a new CheckoutUseCase.
func NewCheckoutUseCase(
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	vendorRepo domain.VendorRepository,
	orderRepo domain.OrderRepository,
	paymentRepo domain.PaymentRepository,
	ledgerRepo domain.LedgerRepository,
	userRepo domain.UserRepository,
	addressRepo domain.AddressRepository,
	vendorCourierRepo domain.VendorCourierRepository,
	shipmentRepo domain.ShipmentRepository,
	rajaOngkir domain.RajaOngkirProvider,
	storage domain.StorageProvider,
	xenditInvoice domain.XenditInvoiceProvider,
	frontendURL string,
	webhookURL string,
) domain.CheckoutUseCase {
	return &checkoutUseCase{
		cartRepo:          cartRepo,
		productRepo:       productRepo,
		vendorRepo:        vendorRepo,
		orderRepo:         orderRepo,
		paymentRepo:       paymentRepo,
		ledgerRepo:        ledgerRepo,
		userRepo:          userRepo,
		addressRepo:       addressRepo,
		vendorCourierRepo: vendorCourierRepo,
		shipmentRepo:      shipmentRepo,
		rajaOngkir:        rajaOngkir,
		storage:           storage,
		xenditInvoice:     xenditInvoice,
		frontendURL:       frontendURL,
		webhookURL:        webhookURL,
	}
}

// enrichedItem holds cart item data enriched with variant and product info.
type enrichedItem struct {
	cartItem domain.CartItem
	variant  domain.ProductVariant
	product  domain.Product
}

type createdCheckoutUnit struct {
	orderID             uuid.UUID
	vendorXenditAccount string
	xenditInvoiceID     *string
	paymentInvoiceID    *uuid.UUID
}

const defaultWeightGram = 500

// Preview returns cart items grouped by vendor with shipping options from RajaOngkir.
func (uc *checkoutUseCase) Preview(ctx context.Context, userID uuid.UUID, req domain.CheckoutPreviewRequest) (*domain.CheckoutPreviewResponse, error) {
	// 1. Validate the selected delivery address.
	addressID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}
	address, err := uc.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("find address: %w", err)
	}
	if address == nil {
		return nil, ErrAddressNotFound
	}
	if address.UserID != userID {
		return nil, ErrAddressNotOwned
	}
	if address.DistrictID == nil || *address.DistrictID == "" {
		return nil, ErrAddressNoDistrict
	}

	// 2. Get active cart and items.
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find cart: %w", err)
	}
	if cart == nil {
		return nil, ErrCartEmpty
	}
	cartItems, err := uc.cartRepo.FindItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("find cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, ErrCartEmpty
	}

	// 3. Enrich cart items and group by vendor.
	vendorGroups := make(map[uuid.UUID][]enrichedItem)
	for _, ci := range cartItems {
		variant, err := uc.productRepo.FindVariantByID(ctx, ci.ProductVariantID)
		if err != nil {
			return nil, fmt.Errorf("find variant: %w", err)
		}
		if variant == nil || !variant.IsActive {
			return nil, ErrCartHasUnavailableItems
		}
		product, err := uc.productRepo.FindByID(ctx, variant.ProductID)
		if err != nil {
			return nil, fmt.Errorf("find product: %w", err)
		}
		if product == nil || product.Status != domain.ProductStatusPublished {
			return nil, ErrCartHasUnavailableItems
		}
		vendorGroups[product.VendorID] = append(vendorGroups[product.VendorID], enrichedItem{
			cartItem: ci,
			variant:  *variant,
			product:  *product,
		})
	}

	// 4. Build preview vendor groups with shipping options.
	var vendorPreviews []domain.CheckoutPreviewVendorGroup
	for vendorID, items := range vendorGroups {
		vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
		if err != nil {
			return nil, fmt.Errorf("find vendor: %w", err)
		}
		if vendor == nil {
			return nil, ErrVendorNotFound
		}

		// 4a. Resolve vendor warehouse origin (default address of vendor owner).
		warehouseAddr, err := uc.addressRepo.FindDefaultByUserID(ctx, vendor.OwnerUserID)
		if err != nil {
			return nil, fmt.Errorf("find vendor warehouse address: %w", err)
		}
		if warehouseAddr == nil || warehouseAddr.DistrictID == nil || *warehouseAddr.DistrictID == "" {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorWarehouseNotFound, vendor.DisplayName)
		}

		// 4b. Get vendor courier selections.
		vendorCouriers, err := uc.vendorCourierRepo.FindByVendorID(ctx, vendorID)
		if err != nil {
			return nil, fmt.Errorf("find vendor couriers: %w", err)
		}
		if len(vendorCouriers) == 0 {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorNoCouriersConfigured, vendor.DisplayName)
		}

		// 4c. Build items and calculate subtotal + total weight.
		var vendorSubtotal float64
		var totalWeightGram int
		previewItems := make([]domain.CheckoutPreviewVendorItem, 0, len(items))
		for _, ei := range items {
			lineTotal := ei.variant.Price * float64(ei.cartItem.Qty)
			vendorSubtotal += lineTotal

			weight := defaultWeightGram
			if ei.variant.WeightGram != nil && *ei.variant.WeightGram > 0 {
				weight = *ei.variant.WeightGram
			}
			totalWeightGram += weight * ei.cartItem.Qty

			imageURL := uc.resolvePrimaryImageURL(ctx, ei.product.ID)

			previewItems = append(previewItems, domain.CheckoutPreviewVendorItem{
				ProductVariantID: ei.variant.ID,
				ProductName:      ei.product.Name,
				VariantName:      ei.variant.VariantName,
				ImageURL:         imageURL,
				Price:            ei.variant.Price,
				Qty:              ei.cartItem.Qty,
				Subtotal:         lineTotal,
				WeightGram:       weight,
			})
		}

		// 4d. Build colon-separated courier codes and call RajaOngkir.
		courierCodes := make([]string, 0, len(vendorCouriers))
		for _, vc := range vendorCouriers {
			if vc.IsActive {
				courierCodes = append(courierCodes, vc.CourierCode)
			}
		}
		var shippingOptions []domain.ShippingCostOption
		if len(courierCodes) > 0 && totalWeightGram > 0 {
			opts, err := uc.rajaOngkir.CalculateDomesticCost(
				ctx,
				*warehouseAddr.DistrictID,
				*address.DistrictID,
				totalWeightGram,
				strings.Join(courierCodes, ":"),
			)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrShippingCostFailed, err)
			}
			shippingOptions = opts
		}

		vendorPreviews = append(vendorPreviews, domain.CheckoutPreviewVendorGroup{
			VendorID:        vendorID,
			VendorName:      vendor.DisplayName,
			Items:           previewItems,
			Subtotal:        vendorSubtotal,
			TotalWeightGram: totalWeightGram,
			ShippingOptions: shippingOptions,
		})
	}

	return &domain.CheckoutPreviewResponse{
		Address:     *toAddressResponse(address),
		Vendors:     vendorPreviews,
		PlatformFee: platformFeeTotal,
	}, nil
}

// resolvePrimaryImageURL finds the primary product image and returns a presigned download URL.
func (uc *checkoutUseCase) resolvePrimaryImageURL(ctx context.Context, productID uuid.UUID) *string {
	images, err := uc.productRepo.FindImagesByProductID(ctx, productID)
	if err != nil || len(images) == 0 {
		return nil
	}
	for _, img := range images {
		if img.IsPrimary {
			presigned, err := uc.storage.GeneratePresignedURL(ctx, img.ImageURL, PresignedDownloadExpiry)
			if err != nil {
				return nil
			}
			return &presigned
		}
	}
	return nil
}

func (uc *checkoutUseCase) Checkout(ctx context.Context, userID uuid.UUID, req domain.CheckoutRequest) (*domain.CheckoutResponse, error) {
	// 1. Validate the selected delivery address.
	addressID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}
	address, err := uc.addressRepo.FindByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("find address: %w", err)
	}
	if address == nil {
		return nil, ErrAddressNotFound
	}
	if address.UserID != userID {
		return nil, ErrAddressNotOwned
	}
	if address.DistrictID == nil || *address.DistrictID == "" {
		return nil, ErrAddressNoDistrict
	}

	// Build shipping choices map keyed by vendor ID.
	shippingMap := make(map[uuid.UUID]domain.CheckoutShippingChoice, len(req.ShippingChoices))
	for _, sc := range req.ShippingChoices {
		vid, err := uuid.Parse(sc.VendorID)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid vendor_id %s", ErrInvalidShippingChoice, sc.VendorID)
		}
		shippingMap[vid] = sc
	}

	// 2. Get active cart.
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find cart: %w", err)
	}
	if cart == nil {
		return nil, ErrCartEmpty
	}

	// 3. Get cart items.
	cartItems, err := uc.cartRepo.FindItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("find cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, ErrCartEmpty
	}

	// 4. Get user info for Xendit customer data.
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	// 5. Enrich and validate cart items, group by vendor.
	vendorGroups := make(map[uuid.UUID][]enrichedItem)

	for _, ci := range cartItems {
		variant, err := uc.productRepo.FindVariantByID(ctx, ci.ProductVariantID)
		if err != nil {
			return nil, fmt.Errorf("find variant: %w", err)
		}
		if variant == nil || !variant.IsActive {
			return nil, ErrCartHasUnavailableItems
		}

		product, err := uc.productRepo.FindByID(ctx, variant.ProductID)
		if err != nil {
			return nil, fmt.Errorf("find product: %w", err)
		}
		if product == nil || product.Status != domain.ProductStatusPublished {
			return nil, ErrCartHasUnavailableItems
		}

		if ci.Qty > variant.StockOnHand {
			return nil, fmt.Errorf("%w: %s", ErrCheckoutStockInsufficient, variant.VariantName)
		}

		vendorGroups[product.VendorID] = append(vendorGroups[product.VendorID], enrichedItem{
			cartItem: ci,
			variant:  *variant,
			product:  *product,
		})
	}

	// 6. For each vendor group: create order -> call Xendit -> save invoice.
	var results []domain.CheckoutOrderResult
	var createdUnits []createdCheckoutUnit
	now := time.Now()

	// Build shipping address snapshot from the selected address.
	shippingSnapshot := buildAddressSnapshot(address)

	for vendorID, items := range vendorGroups {
		// 6a. Look up shipping choice for this vendor.
		sc, ok := shippingMap[vendorID]
		if !ok {
			return nil, fmt.Errorf("%w: missing shipping choice for vendor %s", ErrInvalidShippingChoice, vendorID)
		}

		// 6b. Verify vendor has Xendit account.
		vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
		if err != nil {
			return nil, fmt.Errorf("find vendor: %w", err)
		}
		if vendor == nil || vendor.XenditAccountID == nil || *vendor.XenditAccountID == "" {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorNoXenditAccount, vendorID)
		}

		// 6c. Resolve vendor warehouse origin (default address of vendor owner).
		warehouseAddr, err := uc.addressRepo.FindDefaultByUserID(ctx, vendor.OwnerUserID)
		if err != nil {
			return nil, fmt.Errorf("find vendor warehouse address: %w", err)
		}
		if warehouseAddr == nil || warehouseAddr.DistrictID == nil || *warehouseAddr.DistrictID == "" {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorWarehouseNotFound, vendor.DisplayName)
		}

		// 6d. Get vendor courier selections.
		vendorCouriers, err := uc.vendorCourierRepo.FindByVendorID(ctx, vendorID)
		if err != nil {
			return nil, fmt.Errorf("find vendor couriers: %w", err)
		}
		if len(vendorCouriers) == 0 {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorNoCouriersConfigured, vendor.DisplayName)
		}

		// 6e. Build order items, calculate subtotal and total weight.
		orderID := uuid.New()
		orderNo := generateOrderNo(now)
		createdUnit := createdCheckoutUnit{
			orderID:             orderID,
			vendorXenditAccount: *vendor.XenditAccountID,
		}

		var subtotal float64
		var totalWeightGram int
		orderItems := make([]domain.OrderItem, 0, len(items))
		xenditItems := make([]domain.XenditInvoiceItem, 0, len(items))

		for _, ei := range items {
			lineTotal := ei.variant.Price * float64(ei.cartItem.Qty)
			subtotal += lineTotal

			weight := defaultWeightGram
			if ei.variant.WeightGram != nil && *ei.variant.WeightGram > 0 {
				weight = *ei.variant.WeightGram
			}
			totalWeightGram += weight * ei.cartItem.Qty

			orderItems = append(orderItems, domain.OrderItem{
				ID:                  uuid.New(),
				OrderID:             orderID,
				ProductVariantID:    ei.variant.ID,
				ProductNameSnapshot: ei.product.Name,
				SKUSnapshot:         ei.variant.SKU,
				Qty:                 ei.cartItem.Qty,
				UnitPrice:           ei.variant.Price,
				LineTotal:           lineTotal,
			})

			xenditItems = append(xenditItems, domain.XenditInvoiceItem{
				Name:     fmt.Sprintf("%s - %s", ei.product.Name, ei.variant.VariantName),
				Quantity: ei.cartItem.Qty,
				Price:    ei.variant.Price,
			})
		}

		// 6f. Call RajaOngkir to recalculate shipping cost server-side.
		courierCodes := make([]string, 0, len(vendorCouriers))
		for _, vc := range vendorCouriers {
			if vc.IsActive {
				courierCodes = append(courierCodes, vc.CourierCode)
			}
		}
		if len(courierCodes) == 0 {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorNoCouriersConfigured, vendor.DisplayName)
		}

		shippingOptions, err := uc.rajaOngkir.CalculateDomesticCost(
			ctx,
			*warehouseAddr.DistrictID,
			*address.DistrictID,
			totalWeightGram,
			strings.Join(courierCodes, ":"),
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrShippingCostFailed, err)
		}

		// 6g. Find the matching courier+service from RajaOngkir response.
		var shippingFee float64
		var matchedETD string
		serviceFound := false
		for _, opt := range shippingOptions {
			if strings.EqualFold(opt.Code, sc.CourierCode) && strings.EqualFold(opt.Service, sc.Service) {
				shippingFee = float64(opt.Cost)
				matchedETD = opt.ETD
				serviceFound = true
				break
			}
		}
		if !serviceFound {
			return nil, fmt.Errorf("%w: courier=%s service=%s for vendor %s",
				ErrShippingServiceNotFound, sc.CourierCode, sc.Service, vendor.DisplayName)
		}

		grandTotal := subtotal + shippingFee + platformFeeTotal

		order := &domain.Order{
			ID:                      orderID,
			OrderNo:                 orderNo,
			UserID:                  userID,
			VendorID:                vendorID,
			ShippingAddressSnapshot: shippingSnapshot,
			OrderStatus:             domain.OrderStatusPendingPayment,
			PaymentStatus:           domain.PaymentStatusUnpaid,
			Subtotal:                subtotal,
			ShippingFee:             shippingFee,
			PlatformFee:             platformFeeTotal,
			GrandTotal:              grandTotal,
			PlacedAt:                now,
			CreatedAt:               now,
			UpdatedAt:               now,
		}

		// 6h. Create order + decrement stock in DB transaction.
		if err := uc.orderRepo.CreateOrderWithItems(ctx, order, orderItems); err != nil {
			if errors.Is(err, domain.ErrStockUnavailable) {
				return nil, uc.failCheckoutWithCompensation(ctx, createdUnits, ErrCheckoutStockInsufficient)
			}
			checkoutErr := fmt.Errorf("create order: %w", err)
			return nil, uc.failCheckoutWithCompensation(ctx, createdUnits, checkoutErr)
		}

		// 6h-2. Create shipment record with courier info and ETD from RajaOngkir.
		shipment := &domain.Shipment{
			ID:             uuid.New(),
			OrderID:        orderID,
			CourierCode:    sc.CourierCode,
			ServiceType:    sc.Service,
			ETD:            matchedETD,
			ShipmentStatus: "waiting_pickup",
		}
		if err := uc.shipmentRepo.Create(ctx, shipment); err != nil {
			checkoutErr := fmt.Errorf("create shipment: %w", err)
			return nil, uc.failCheckoutWithCompensation(ctx, createdUnits, checkoutErr)
		}

		createdUnits = append(createdUnits, createdUnit)
		unitIdx := len(createdUnits) - 1

		// 6i. Call Xendit to create invoice.
		externalInvoiceID := fmt.Sprintf("INV-%s-%s", orderNo, uuid.New().String()[:8])

		xenditReq := domain.XenditInvoiceRequest{
			ExternalID:      externalInvoiceID,
			Amount:          grandTotal,
			Description:     fmt.Sprintf("Pembayaran order %s", orderNo),
			InvoiceDuration: invoiceDuration,
			Customer: &domain.XenditInvoiceCustomer{
				GivenNames: user.FullName,
				Email:      user.Email,
			},
			SuccessRedirectURL: fmt.Sprintf("%s/payment/success?order_id=%s", uc.frontendURL, orderID),
			FailureRedirectURL: fmt.Sprintf("%s/payment/failed?order_id=%s", uc.frontendURL, orderID),
			CallbackURL:        uc.webhookURL,
			Currency:           "IDR",
			Items:              xenditItems,
			Metadata: map[string]interface{}{
				"order_id":  orderID.String(),
				"order_no":  orderNo,
				"vendor_id": vendorID.String(),
			},
		}

		xenditResp, err := uc.xenditInvoice.CreateInvoice(ctx, *vendor.XenditAccountID, xenditReq)
		if err != nil {
			checkoutErr := fmt.Errorf("%w: %v", ErrInvoiceCreationFailed, err)
			return nil, uc.failCheckoutWithCompensation(ctx, createdUnits, checkoutErr)
		}
		createdUnits[unitIdx].xenditInvoiceID = &xenditResp.ID

		// 6j. Parse expiry date.
		expiresAt, _ := time.Parse(time.RFC3339, xenditResp.ExpiryDate)

		// 6k. Save payment invoice record.
		paymentInvoice := &domain.PaymentInvoice{
			ID:                uuid.New(),
			OrderID:           orderID,
			Gateway:           "xendit",
			XenditInvoiceID:   &xenditResp.ID,
			ExternalInvoiceID: externalInvoiceID,
			InvoiceURL:        &xenditResp.InvoiceURL,
			Amount:            grandTotal,
			Currency:          "IDR",
			Status:            domain.InvoiceStatusPending,
			ExpiresAt:         &expiresAt,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		if err := uc.paymentRepo.CreateInvoice(ctx, paymentInvoice); err != nil {
			checkoutErr := fmt.Errorf("save payment invoice: %w", err)
			return nil, uc.failCheckoutWithCompensation(ctx, createdUnits, checkoutErr)
		}
		paymentInvoiceID := paymentInvoice.ID
		createdUnits[unitIdx].paymentInvoiceID = &paymentInvoiceID

		results = append(results, domain.CheckoutOrderResult{
			OrderID:     orderID,
			OrderNo:     orderNo,
			VendorID:    vendorID,
			VendorName:  vendor.DisplayName,
			Subtotal:    subtotal,
			ShippingFee: shippingFee,
			PlatformFee: platformFeeTotal,
			GrandTotal:  grandTotal,
			InvoiceURL:  xenditResp.InvoiceURL,
			ExpiresAt:   expiresAt,
		})
	}

	// 7. Mark cart as converted.
	if err := uc.cartRepo.UpdateStatus(ctx, cart.ID, domain.CartStatusConverted); err != nil {
		return nil, fmt.Errorf("update cart status: %w", err)
	}

	return &domain.CheckoutResponse{Orders: results}, nil
}

func (uc *checkoutUseCase) failCheckoutWithCompensation(ctx context.Context, createdUnits []createdCheckoutUnit, checkoutErr error) error {
	if len(createdUnits) == 0 {
		return checkoutErr
	}

	if err := uc.compensateCheckoutFailure(ctx, createdUnits, checkoutErr); err != nil {
		return fmt.Errorf("%w: checkout error: %v; compensation error: %v", ErrCheckoutCompensationFailed, checkoutErr, err)
	}

	return checkoutErr
}

func (uc *checkoutUseCase) compensateCheckoutFailure(ctx context.Context, createdUnits []createdCheckoutUnit, checkoutErr error) error {
	var compensationErrs []string
	rollbackNotes := fmt.Sprintf("Checkout rolled back: %v", checkoutErr)
	rawPayload := map[string]interface{}{
		"rollback_reason": checkoutErr.Error(),
	}

	for i := len(createdUnits) - 1; i >= 0; i-- {
		unit := createdUnits[i]

		if unit.xenditInvoiceID != nil {
			if err := uc.xenditInvoice.ExpireInvoice(ctx, unit.vendorXenditAccount, *unit.xenditInvoiceID); err != nil {
				compensationErrs = append(compensationErrs, fmt.Sprintf("expire invoice for order %s: %v", unit.orderID, err))
			}
		}

		if unit.paymentInvoiceID != nil {
			if err := uc.paymentRepo.UpdateInvoiceStatus(ctx, *unit.paymentInvoiceID, domain.InvoiceStatusFailed, nil, nil, nil, rawPayload); err != nil {
				compensationErrs = append(compensationErrs, fmt.Sprintf("mark payment invoice failed for order %s: %v", unit.orderID, err))
			}
		}

		if err := uc.orderRepo.UpdateOrderStatus(ctx, unit.orderID, domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, nil, &rollbackNotes); err != nil {
			compensationErrs = append(compensationErrs, fmt.Sprintf("cancel order %s: %v", unit.orderID, err))
		}

		if err := uc.orderRepo.RestoreStock(ctx, unit.orderID); err != nil {
			compensationErrs = append(compensationErrs, fmt.Sprintf("restore stock for order %s: %v", unit.orderID, err))
		}
	}

	if len(compensationErrs) > 0 {
		return errors.New(strings.Join(compensationErrs, "; "))
	}
	return nil
}

func (uc *checkoutUseCase) HandleWebhook(ctx context.Context, payload domain.XenditWebhookPayload) error {
	// 1. Idempotency check.
	externalEventID := fmt.Sprintf("%s:%s", payload.ID, payload.Status)
	exists, err := uc.paymentRepo.EventExistsByExternalID(ctx, externalEventID)
	if err != nil {
		return fmt.Errorf("check event existence: %w", err)
	}
	if exists {
		return nil // already processed
	}

	// 2. Find payment invoice by external_id.
	invoice, err := uc.paymentRepo.FindInvoiceByExternalID(ctx, payload.ExternalID)
	if err != nil {
		return fmt.Errorf("find payment invoice: %w", err)
	}
	if invoice == nil {
		return ErrInvoiceNotFound
	}

	// 3. Record the event for idempotency.
	payloadMap := map[string]interface{}{
		"id":              payload.ID,
		"external_id":     payload.ExternalID,
		"status":          payload.Status,
		"amount":          payload.Amount,
		"paid_amount":     payload.PaidAmount,
		"paid_at":         payload.PaidAt,
		"payment_method":  payload.PaymentMethod,
		"payment_channel": payload.PaymentChannel,
	}
	event := &domain.PaymentEvent{
		ID:               uuid.New(),
		PaymentInvoiceID: invoice.ID,
		EventType:        fmt.Sprintf("invoice.%s", strings.ToLower(payload.Status)),
		ExternalEventID:  externalEventID,
		Payload:          payloadMap,
		ReceivedAt:       time.Now(),
	}
	if err := uc.paymentRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("record payment event: %w", err)
	}

	// 4. Handle based on status.
	switch payload.Status {
	case "PAID":
		return uc.handlePaid(ctx, invoice, payload)
	case "EXPIRED":
		return uc.handleExpired(ctx, invoice)
	default:
		return nil
	}
}

func (uc *checkoutUseCase) handlePaid(ctx context.Context, invoice *domain.PaymentInvoice, payload domain.XenditWebhookPayload) error {
	// Guard against late paid webhooks for orders already canceled/settled.
	order, err := uc.orderRepo.FindByID(ctx, invoice.OrderID)
	if err != nil {
		return fmt.Errorf("find order for paid webhook: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}
	if order.OrderStatus != domain.OrderStatusPendingPayment {
		return nil
	}

	// 1. Update payment invoice to paid.
	var paidAt *time.Time
	if payload.PaidAt != nil {
		t, err := time.Parse(time.RFC3339, *payload.PaidAt)
		if err == nil {
			paidAt = &t
		}
	}
	rawPayload := map[string]interface{}{"webhook_payload": payload}

	updated, err := uc.paymentRepo.UpdateInvoiceStatusIfCurrent(
		ctx,
		invoice.ID,
		domain.InvoiceStatusPending,
		domain.InvoiceStatusPaid,
		paidAt,
		&payload.PaymentMethod,
		&payload.PaymentChannel,
		rawPayload,
	)
	if err != nil {
		return fmt.Errorf("update invoice status: %w", err)
	}
	if !updated {
		return nil
	}

	// 2. Update order status to paid.
	notes := fmt.Sprintf("Payment received via %s/%s", payload.PaymentMethod, payload.PaymentChannel)
	if err := uc.orderRepo.UpdateOrderStatus(ctx, invoice.OrderID,
		domain.OrderStatusPaid, domain.PaymentStatusPaid, nil, &notes); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	// 3. Record ledger journal (double-entry bookkeeping).
	if err := uc.recordPaymentLedger(ctx, invoice, order, payload.PaidAmount); err != nil {
		return fmt.Errorf("record ledger: %w", err)
	}

	return nil
}

func (uc *checkoutUseCase) handleExpired(ctx context.Context, invoice *domain.PaymentInvoice) error {
	order, err := uc.orderRepo.FindByID(ctx, invoice.OrderID)
	if err != nil {
		return fmt.Errorf("find order for expired webhook: %w", err)
	}
	if order == nil {
		return ErrOrderNotFound
	}
	if order.OrderStatus != domain.OrderStatusPendingPayment {
		return nil
	}

	notes := "Payment link expired"
	applied, err := uc.orderRepo.ApplyExpiredWebhookUpdate(ctx, invoice.OrderID, invoice.ID, &notes)
	if err != nil {
		return fmt.Errorf("apply expired webhook update: %w", err)
	}
	if !applied {
		return nil
	}

	return nil
}

func (uc *checkoutUseCase) recordPaymentLedger(ctx context.Context, invoice *domain.PaymentInvoice, order *domain.Order, paidAmount float64) error {
	// Find ledger accounts.
	gatewayReceivable, err := uc.ledgerRepo.FindAccountByCode(ctx, "1100")
	if err != nil || gatewayReceivable == nil {
		return fmt.Errorf("ledger account 1100 not found")
	}
	vendorPayable, err := uc.ledgerRepo.FindAccountByCode(ctx, "2100")
	if err != nil || vendorPayable == nil {
		return fmt.Errorf("ledger account 2100 not found")
	}
	platformFeeRevenue, err := uc.ledgerRepo.FindAccountByCode(ctx, "4100")
	if err != nil || platformFeeRevenue == nil {
		return fmt.Errorf("ledger account 4100 not found")
	}
	adminFeeRevenue, err := uc.ledgerRepo.FindAccountByCode(ctx, "4200")
	if err != nil || adminFeeRevenue == nil {
		return fmt.Errorf("ledger account 4200 not found")
	}

	vendorAmount := order.Subtotal
	journalNo := fmt.Sprintf("JRN-PAY-%s-%s", order.OrderNo, uuid.New().String()[:8])
	desc := fmt.Sprintf("Payment received for order %s", order.OrderNo)
	orderNoRef := order.OrderNo

	journal := &domain.LedgerJournal{
		ID:          uuid.New(),
		JournalNo:   journalNo,
		SourceType:  "payment_invoice",
		SourceID:    invoice.ID,
		EventTime:   time.Now(),
		Description: &desc,
		Status:      "posted",
		CreatedAt:   time.Now(),
	}

	lines := []domain.LedgerLine{
		// DR Payment Gateway Receivable = full paid amount
		{ID: uuid.New(), JournalID: journal.ID, AccountID: gatewayReceivable.ID, Debit: paidAmount, Credit: 0, Currency: "IDR", Reference: &orderNoRef},
		// CR Vendor Payable = subtotal (amount owed to vendor)
		{ID: uuid.New(), JournalID: journal.ID, AccountID: vendorPayable.ID, Debit: 0, Credit: vendorAmount, Currency: "IDR", Reference: &orderNoRef},
		// CR Platform Fee Revenue = app fee (Rp 1,000)
		{ID: uuid.New(), JournalID: journal.ID, AccountID: platformFeeRevenue.ID, Debit: 0, Credit: float64(platformAppFee), Currency: "IDR", Reference: &orderNoRef},
		// CR Admin Fee Revenue = admin fee (Rp 5,000)
		{ID: uuid.New(), JournalID: journal.ID, AccountID: adminFeeRevenue.ID, Debit: 0, Credit: float64(platformAdminFee), Currency: "IDR", Reference: &orderNoRef},
	}

	return uc.ledgerRepo.CreateJournalWithLines(ctx, journal, lines)
}

// buildAddressSnapshot creates a JSON-compatible map from a domain.Address for order storage.
func buildAddressSnapshot(a *domain.Address) map[string]interface{} {
	snap := map[string]interface{}{
		"address_id":   a.ID.String(),
		"address_line": a.AddressLine,
		"is_default":   a.IsDefault,
	}
	if a.Label != nil {
		snap["label"] = *a.Label
	}
	if a.RecipientName != nil {
		snap["recipient_name"] = *a.RecipientName
	}
	if a.Phone != nil {
		snap["phone"] = *a.Phone
	}
	if a.ProvinceName != nil {
		snap["province_name"] = *a.ProvinceName
	}
	if a.CityName != nil {
		snap["city_name"] = *a.CityName
	}
	if a.DistrictName != nil {
		snap["district_name"] = *a.DistrictName
	}
	if a.SubdistrictName != nil {
		snap["subdistrict_name"] = *a.SubdistrictName
	}
	if a.PostalCode != nil {
		snap["postal_code"] = *a.PostalCode
	}
	return snap
}

func generateOrderNo(t time.Time) string {
	return fmt.Sprintf("ORD-%s-%s",
		t.Format("20060102"),
		strings.ToUpper(uuid.New().String()[:8]),
	)
}
