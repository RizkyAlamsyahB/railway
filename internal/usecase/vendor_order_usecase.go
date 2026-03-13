package usecase

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type vendorOrderUseCase struct {
	vendorOrderRepo domain.VendorOrderRepository
	orderRepo       domain.OrderRepository
	shipmentRepo    domain.ShipmentRepository
	paymentRepo     domain.PaymentRepository
	userRepo        domain.UserRepository
	rajaOngkir      domain.RajaOngkirProvider
}

// NewVendorOrderUseCase creates a new VendorOrderUseCase.
func NewVendorOrderUseCase(
	vendorOrderRepo domain.VendorOrderRepository,
	orderRepo domain.OrderRepository,
	shipmentRepo domain.ShipmentRepository,
	paymentRepo domain.PaymentRepository,
	userRepo domain.UserRepository,
	rajaOngkir domain.RajaOngkirProvider,
) domain.VendorOrderUseCase {
	return &vendorOrderUseCase{
		vendorOrderRepo: vendorOrderRepo,
		orderRepo:       orderRepo,
		shipmentRepo:    shipmentRepo,
		paymentRepo:     paymentRepo,
		userRepo:        userRepo,
		rajaOngkir:      rajaOngkir,
	}
}

// ListOrders returns paginated orders for a vendor.
func (uc *vendorOrderUseCase) ListOrders(ctx context.Context, vendorID uuid.UUID, params domain.VendorOrderListParams) ([]domain.VendorOrderListItem, *domain.PaginationMeta, error) {
	params = normalizeVendorOrderListParams(params)
	if params.Status != "" && !isValidOrderStatus(params.Status) {
		return nil, nil, ErrInvalidOrderStatus
	}

	orders, total, err := uc.vendorOrderRepo.ListByVendor(ctx, vendorID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list vendor orders: %w", err)
	}

	meta := buildVendorOrderPaginationMeta(params.Page, params.Limit, total)
	if len(orders) == 0 {
		return []domain.VendorOrderListItem{}, meta, nil
	}

	// Fetch order items and customer names in bulk.
	orderIDs := make([]uuid.UUID, 0, len(orders))
	userIDs := make(map[uuid.UUID]bool)
	for _, o := range orders {
		orderIDs = append(orderIDs, o.ID)
		userIDs[o.UserID] = true
	}

	itemsByOrderID, err := uc.orderRepo.FindItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	// Fetch customer names.
	customerNames := make(map[uuid.UUID]string)
	for uid := range userIDs {
		user, err := uc.userRepo.FindByID(ctx, uid)
		if err == nil && user != nil {
			customerNames[uid] = user.FullName
		}
	}

	result := make([]domain.VendorOrderListItem, 0, len(orders))
	for _, order := range orders {
		rawItems := itemsByOrderID[order.ID]
		items := make([]domain.VendorOrderProductItem, 0, len(rawItems))
		for _, it := range rawItems {
			items = append(items, domain.VendorOrderProductItem{
				ProductName:     it.ProductNameSnapshot,
				SelectedVariant: it.SKUSnapshot,
				Qty:             it.Qty,
				UnitPrice:       it.UnitPrice,
				LineTotal:       it.LineTotal,
			})
		}

		result = append(result, domain.VendorOrderListItem{
			OrderID:      order.ID,
			OrderNo:      order.OrderNo,
			OrderDate:    order.PlacedAt,
			CustomerName: customerNames[order.UserID],
			Items:        items,
			TotalPayment: order.GrandTotal,
			Status:       order.OrderStatus,
		})
	}

	return result, meta, nil
}

// GetOrderDetail returns full order detail for a vendor.
func (uc *vendorOrderUseCase) GetOrderDetail(ctx context.Context, vendorID, orderID uuid.UUID) (*domain.VendorOrderDetailResponse, error) {
	order, err := uc.vendorOrderRepo.FindByIDAndVendor(ctx, orderID, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotBelongToVendor
	}

	// Fetch order items.
	rawItems, err := uc.orderRepo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	items := make([]domain.VendorOrderProductItem, 0, len(rawItems))
	for _, it := range rawItems {
		items = append(items, domain.VendorOrderProductItem{
			ProductName:     it.ProductNameSnapshot,
			SelectedVariant: it.SKUSnapshot,
			Qty:             it.Qty,
			UnitPrice:       it.UnitPrice,
			LineTotal:       it.LineTotal,
		})
	}

	// Fetch customer name.
	customerName := ""
	user, err := uc.userRepo.FindByID(ctx, order.UserID)
	if err == nil && user != nil {
		customerName = user.FullName
	}

	// Fetch shipment (may be nil if not shipped yet).
	shipment, _ := uc.shipmentRepo.FindByOrderID(ctx, orderID)
	var shipmentResp *domain.VendorOrderShipmentResponse
	if shipment != nil {
		resp := domain.VendorOrderShipmentResponse{
			CourierCode:    shipment.CourierCode,
			ServiceType:    shipment.ServiceType,
			TrackingNo:     shipment.TrackingNo,
			ETD:            shipment.ETD,
			ShipmentStatus: shipment.ShipmentStatus,
			ShippedAt:      shipment.ShippedAt,
			DeliveredAt:    shipment.DeliveredAt,
		}
		if shipment.ShippedAt != nil {
			resp.EstimatedArrival = calculateEstimatedArrival(*shipment.ShippedAt, shipment.ETD)
		}
		shipmentResp = &resp
	}

	// Fetch payment info.
	var paymentInfo *domain.VendorOrderPaymentInfo
	invoice, _ := uc.paymentRepo.FindInvoiceByOrderID(ctx, orderID)
	if invoice != nil {
		paymentInfo = &domain.VendorOrderPaymentInfo{
			Status:         invoice.Status,
			PaymentMethod:  invoice.PaymentMethod,
			PaymentChannel: invoice.PaymentChannel,
			PaidAt:         invoice.PaidAt,
		}
	}

	return &domain.VendorOrderDetailResponse{
		OrderID:      order.ID,
		OrderNo:      order.OrderNo,
		OrderDate:    order.PlacedAt,
		Status:       order.OrderStatus,
		CustomerName: customerName,
		Items:        items,
		Subtotal:     order.Subtotal,
		ShippingFee:  order.ShippingFee,
		PlatformFee:  order.PlatformFee,
		GrandTotal:   order.GrandTotal,
		Address:      order.ShippingAddressSnapshot,
		Shipment:     shipmentResp,
		PaymentInfo:  paymentInfo,
	}, nil
}

// AcceptOrder transitions an order from paid → processing.
func (uc *vendorOrderUseCase) AcceptOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID) (*domain.AcceptRejectOrderResponse, error) {
	order, err := uc.vendorOrderRepo.FindByIDAndVendor(ctx, orderID, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotBelongToVendor
	}

	if order.OrderStatus != domain.OrderStatusPaid {
		return nil, ErrInvalidOrderAcceptTransition
	}

	notes := "Order accepted by vendor"
	if err := uc.orderRepo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusProcessing, order.PaymentStatus, &userID, &notes); err != nil {
		return nil, fmt.Errorf("failed to accept order: %w", err)
	}

	return &domain.AcceptRejectOrderResponse{
		OrderID:     orderID,
		OrderStatus: domain.OrderStatusProcessing,
		UpdatedAt:   time.Now(),
	}, nil
}

// RejectOrder transitions an order from paid/processing → canceled and restores stock.
func (uc *vendorOrderUseCase) RejectOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID) (*domain.AcceptRejectOrderResponse, error) {
	order, err := uc.vendorOrderRepo.FindByIDAndVendor(ctx, orderID, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotBelongToVendor
	}

	switch order.OrderStatus {
	case domain.OrderStatusPaid, domain.OrderStatusProcessing:
		// allowed
	default:
		return nil, ErrInvalidOrderRejectTransition
	}

	notes := "Order rejected by vendor"
	if err := uc.orderRepo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusCanceled, order.PaymentStatus, &userID, &notes); err != nil {
		return nil, fmt.Errorf("failed to reject order: %w", err)
	}

	// Restore stock for the rejected order.
	if err := uc.orderRepo.RestoreStock(ctx, orderID); err != nil {
		// Log but do not fail — the status change already committed.
		fmt.Printf("WARN: failed to restore stock for order %s: %v\n", orderID, err)
	}

	return &domain.AcceptRejectOrderResponse{
		OrderID:     orderID,
		OrderStatus: domain.OrderStatusCanceled,
		UpdatedAt:   time.Now(),
	}, nil
}

// ShipOrder inputs tracking number and transitions from processing → shipped.
func (uc *vendorOrderUseCase) ShipOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID, req domain.ShipOrderRequest) (*domain.ShipOrderResponse, error) {
	order, err := uc.vendorOrderRepo.FindByIDAndVendor(ctx, orderID, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotBelongToVendor
	}

	if order.OrderStatus != domain.OrderStatusProcessing {
		return nil, ErrInvalidOrderShipTransition
	}

	if strings.TrimSpace(req.TrackingNo) == "" {
		return nil, ErrTrackingNumberRequired
	}

	// Find existing shipment (created during checkout with courier info).
	shipment, err := uc.shipmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shipment: %w", err)
	}
	if shipment == nil {
		return nil, ErrShipmentNotFound
	}

	now := time.Now()
	trackingNo := strings.TrimSpace(req.TrackingNo)

	// Update existing shipment with tracking number and ship status.
	if err := uc.shipmentRepo.UpdateTrackingAndShip(ctx, orderID, trackingNo, now); err != nil {
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}

	// Update order status to shipped.
	notes := fmt.Sprintf("Order shipped by vendor, tracking: %s (%s %s)", trackingNo, shipment.CourierCode, shipment.ServiceType)
	if err := uc.orderRepo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusShipped, order.PaymentStatus, &userID, &notes); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return &domain.ShipOrderResponse{
		OrderID:          orderID,
		OrderStatus:      domain.OrderStatusShipped,
		TrackingNo:       trackingNo,
		CourierCode:      shipment.CourierCode,
		ShippedAt:        now,
		ETD:              shipment.ETD,
		EstimatedArrival: calculateEstimatedArrival(now, shipment.ETD),
	}, nil
}

// TrackWaybill returns tracking info for an order's shipment.
func (uc *vendorOrderUseCase) TrackWaybill(ctx context.Context, vendorID, orderID uuid.UUID) (*domain.TrackWaybillResponse, error) {
	order, err := uc.vendorOrderRepo.FindByIDAndVendor(ctx, orderID, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotBelongToVendor
	}

	shipment, err := uc.shipmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shipment: %w", err)
	}
	if shipment == nil {
		return nil, ErrShipmentNotFound
	}

	result, err := uc.rajaOngkir.TrackWaybill(ctx, shipment.TrackingNo, shipment.CourierCode)
	if err != nil {
		return nil, ErrTrackingFailed
	}

	return result, nil
}

// --- helpers ---

func normalizeVendorOrderListParams(params domain.VendorOrderListParams) domain.VendorOrderListParams {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}
	params.Status = strings.ToLower(strings.TrimSpace(params.Status))
	return params
}

func buildVendorOrderPaginationMeta(page, limit int, total int64) *domain.PaginationMeta {
	totalPages := 0
	if limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
	return &domain.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}

// etdNumbersRe matches one or more digits in an ETD string.
var etdNumbersRe = regexp.MustCompile(`\d+`)

// parseMaxETDDays extracts the maximum number of days from a RajaOngkir ETD string.
// Examples: "1-2 day" → 2, "3 day" → 3, "1-2" → 2, "" → 0.
func parseMaxETDDays(etd string) int {
	matches := etdNumbersRe.FindAllString(etd, -1)
	maxDays := 0
	for _, m := range matches {
		n, err := strconv.Atoi(m)
		if err == nil && n > maxDays {
			maxDays = n
		}
	}
	return maxDays
}

// calculateEstimatedArrival computes estimated arrival date from shipped time + ETD string.
// Returns nil if etd is empty or unparseable.
func calculateEstimatedArrival(shippedAt time.Time, etd string) *string {
	days := parseMaxETDDays(etd)
	if days <= 0 {
		return nil
	}
	arrival := shippedAt.Add(time.Duration(days) * 24 * time.Hour)
	formatted := arrival.Format("2006-01-02")
	return &formatted
}
