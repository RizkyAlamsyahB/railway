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
	cartRepo      domain.CartRepository
	productRepo   domain.ProductRepository
	vendorRepo    domain.VendorRepository
	orderRepo     domain.OrderRepository
	paymentRepo   domain.PaymentRepository
	ledgerRepo    domain.LedgerRepository
	userRepo      domain.UserRepository
	xenditInvoice domain.XenditInvoiceProvider
	frontendURL   string
	webhookURL    string
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
	xenditInvoice domain.XenditInvoiceProvider,
	frontendURL string,
	webhookURL string,
) domain.CheckoutUseCase {
	return &checkoutUseCase{
		cartRepo:      cartRepo,
		productRepo:   productRepo,
		vendorRepo:    vendorRepo,
		orderRepo:     orderRepo,
		paymentRepo:   paymentRepo,
		ledgerRepo:    ledgerRepo,
		userRepo:      userRepo,
		xenditInvoice: xenditInvoice,
		frontendURL:   frontendURL,
		webhookURL:    webhookURL,
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

func (uc *checkoutUseCase) Checkout(ctx context.Context, userID uuid.UUID) (*domain.CheckoutResponse, error) {
	// 1. Get active cart.
	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find cart: %w", err)
	}
	if cart == nil {
		return nil, ErrCartEmpty
	}

	// 2. Get cart items.
	cartItems, err := uc.cartRepo.FindItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("find cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return nil, ErrCartEmpty
	}

	// 3. Get user info for Xendit customer data.
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	// 4. Enrich and validate cart items, group by vendor.
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

	// 5. For each vendor group: create order -> call Xendit -> save invoice.
	var results []domain.CheckoutOrderResult
	var createdUnits []createdCheckoutUnit
	now := time.Now()

	for vendorID, items := range vendorGroups {
		// 5a. Verify vendor has Xendit account.
		vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
		if err != nil {
			return nil, fmt.Errorf("find vendor: %w", err)
		}
		if vendor == nil || vendor.XenditAccountID == nil || *vendor.XenditAccountID == "" {
			return nil, fmt.Errorf("%w: vendor %s", ErrVendorNoXenditAccount, vendorID)
		}

		// 5b. Build order.
		orderID := uuid.New()
		orderNo := generateOrderNo(now)
		createdUnit := createdCheckoutUnit{
			orderID:             orderID,
			vendorXenditAccount: *vendor.XenditAccountID,
		}

		var subtotal float64
		orderItems := make([]domain.OrderItem, 0, len(items))
		xenditItems := make([]domain.XenditInvoiceItem, 0, len(items))

		for _, ei := range items {
			lineTotal := ei.variant.Price * float64(ei.cartItem.Qty)
			subtotal += lineTotal

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

		grandTotal := subtotal + platformFeeTotal

		shippingSnapshot := map[string]interface{}{
			"name":    user.FullName,
			"address": "Placeholder - shipping not yet implemented",
		}

		order := &domain.Order{
			ID:                      orderID,
			OrderNo:                 orderNo,
			UserID:                  userID,
			VendorID:                vendorID,
			ShippingAddressSnapshot: shippingSnapshot,
			OrderStatus:             domain.OrderStatusPendingPayment,
			PaymentStatus:           domain.PaymentStatusUnpaid,
			Subtotal:                subtotal,
			ShippingFee:             0,
			PlatformFee:             platformFeeTotal,
			GrandTotal:              grandTotal,
			PlacedAt:                now,
			CreatedAt:               now,
			UpdatedAt:               now,
		}

		// 5c. Create order + decrement stock in DB transaction.
		if err := uc.orderRepo.CreateOrderWithItems(ctx, order, orderItems); err != nil {
			checkoutErr := fmt.Errorf("create order: %w", err)
			return nil, uc.failCheckoutWithCompensation(ctx, createdUnits, checkoutErr)
		}
		createdUnits = append(createdUnits, createdUnit)
		unitIdx := len(createdUnits) - 1

		// 5d. Call Xendit to create invoice.
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

		// 5e. Parse expiry date.
		expiresAt, _ := time.Parse(time.RFC3339, xenditResp.ExpiryDate)

		// 5f. Save payment invoice record.
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
			ShippingFee: 0,
			PlatformFee: platformFeeTotal,
			GrandTotal:  grandTotal,
			InvoiceURL:  xenditResp.InvoiceURL,
			ExpiresAt:   expiresAt,
		})
	}

	// 6. Mark cart as converted.
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

	if err := uc.paymentRepo.UpdateInvoiceStatus(ctx, invoice.ID, domain.InvoiceStatusPaid,
		paidAt, &payload.PaymentMethod, &payload.PaymentChannel, rawPayload); err != nil {
		return fmt.Errorf("update invoice status: %w", err)
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

	// 4. Credit vendor balance (net = subtotal, i.e. grand_total - platform_fee).
	//    This makes the revenue immediately available for vendor self-service withdrawal.
	netAmount := order.Subtotal // vendor revenue = subtotal (platform fee is not theirs)
	if err := uc.vendorRepo.CreditBalance(ctx, order.VendorID, netAmount); err != nil {
		return fmt.Errorf("credit vendor balance: %w", err)
	}

	return nil
}

func (uc *checkoutUseCase) handleExpired(ctx context.Context, invoice *domain.PaymentInvoice) error {
	// 1. Update payment invoice to expired.
	if err := uc.paymentRepo.UpdateInvoiceStatus(ctx, invoice.ID, domain.InvoiceStatusExpired,
		nil, nil, nil, nil); err != nil {
		return fmt.Errorf("update invoice status: %w", err)
	}

	// 2. Update order to canceled.
	notes := "Payment link expired"
	if err := uc.orderRepo.UpdateOrderStatus(ctx, invoice.OrderID,
		domain.OrderStatusCanceled, domain.PaymentStatusUnpaid, nil, &notes); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	// 3. Restore stock.
	if err := uc.orderRepo.RestoreStock(ctx, invoice.OrderID); err != nil {
		return fmt.Errorf("restore stock: %w", err)
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

func generateOrderNo(t time.Time) string {
	return fmt.Sprintf("ORD-%s-%s",
		t.Format("20060102"),
		strings.ToUpper(uuid.New().String()[:8]),
	)
}
