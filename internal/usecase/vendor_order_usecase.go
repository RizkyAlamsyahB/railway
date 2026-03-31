package usecase

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
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
	vendorRepo      domain.VendorRepository
	rajaOngkir      domain.RajaOngkirProvider
}

// NewVendorOrderUseCase creates a new VendorOrderUseCase.
func NewVendorOrderUseCase(
	vendorOrderRepo domain.VendorOrderRepository,
	orderRepo domain.OrderRepository,
	shipmentRepo domain.ShipmentRepository,
	paymentRepo domain.PaymentRepository,
	userRepo domain.UserRepository,
	vendorRepo domain.VendorRepository,
	rajaOngkir domain.RajaOngkirProvider,
) domain.VendorOrderUseCase {
	return &vendorOrderUseCase{
		vendorOrderRepo: vendorOrderRepo,
		orderRepo:       orderRepo,
		shipmentRepo:    shipmentRepo,
		paymentRepo:     paymentRepo,
		userRepo:        userRepo,
		vendorRepo:      vendorRepo,
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

// ExportOrders returns all orders for CSV export (no pagination).
func (uc *vendorOrderUseCase) ExportOrders(ctx context.Context, vendorID uuid.UUID, params domain.VendorOrderListParams) ([]domain.VendorOrderExportItem, error) {
	// Remove pagination for export - get all matching records
	params.Page = 1
	params.Limit = 10000 // Max export limit

	params = normalizeVendorOrderListParams(params)
	if params.Status != "" && !isValidOrderStatus(params.Status) {
		return nil, ErrInvalidOrderStatus
	}

	orders, _, err := uc.vendorOrderRepo.ListByVendor(ctx, vendorID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list vendor orders for export: %w", err)
	}

	if len(orders) == 0 {
		return []domain.VendorOrderExportItem{}, nil
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
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	// Fetch customer names.
	customerNames := make(map[uuid.UUID]string)
	for uid := range userIDs {
		user, err := uc.userRepo.FindByID(ctx, uid)
		if err == nil && user != nil {
			customerNames[uid] = user.FullName
		}
	}

	// Build export items - one row per order item
	result := make([]domain.VendorOrderExportItem, 0)
	for _, order := range orders {
		rawItems := itemsByOrderID[order.ID]
		for _, it := range rawItems {
			result = append(result, domain.VendorOrderExportItem{
				OrderNo:      order.OrderNo,
				OrderDate:    order.PlacedAt,
				CustomerName: customerNames[order.UserID],
				ProductName:  it.ProductNameSnapshot,
				Variant:      it.SKUSnapshot,
				Qty:          it.Qty,
				UnitPrice:    it.UnitPrice,
				LineTotal:    it.LineTotal,
				Subtotal:     order.Subtotal,
				ShippingFee:  order.ShippingFee,
				PlatformFee:  order.PlatformFee,
				GrandTotal:   order.GrandTotal,
				Status:       order.OrderStatus,
			})
		}
	}

	return result, nil
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

	items := make([]domain.VendorOrderDetailItem, 0, len(rawItems))
	for _, it := range rawItems {
		items = append(items, domain.VendorOrderDetailItem{
			Name:     it.ProductNameSnapshot,
			Variant:  it.SKUSnapshot,
			Price:    it.UnitPrice,
			Quantity: it.Qty,
			Subtotal: it.LineTotal,
		})
	}

	// Fetch vendor/store name.
	storeName := ""
	vendor, err := uc.vendorRepo.FindByID(ctx, order.VendorID)
	if err == nil && vendor != nil {
		storeName = vendor.DisplayName
	}

	// Fetch shipment (may be nil if not shipped yet).
	shipment, _ := uc.shipmentRepo.FindByOrderID(ctx, orderID)

	// Build shipping info with tracking status.
	shippingInfo := uc.buildShippingInfo(ctx, shipment, order.OrderStatus)

	// Auto-update order status to "received" if courier confirms delivery.
	if shippingInfo.deliveryDetected && order.OrderStatus == domain.OrderStatusShipped {
		notes := "Auto-updated: courier confirmed delivery"
		uc.orderRepo.MarkOrderReceived(ctx, orderID, nil, &notes)
		if shippingInfo.DeliveredAt != nil {
			uc.shipmentRepo.UpdateDeliveredAt(ctx, orderID, *shippingInfo.DeliveredAt)
		}
	}

	// Parse recipient from shipping address snapshot.
	recipient := uc.buildRecipientInfo(order.ShippingAddressSnapshot)

	// Fetch payment info.
	var paymentMethod string
	var paidAt *time.Time
	invoice, _ := uc.paymentRepo.FindInvoiceByOrderID(ctx, orderID)
	if invoice != nil {
		paymentMethod = buildPaymentMethodDisplay(invoice.PaymentMethod)
		paidAt = invoice.PaidAt
	}

	return &domain.VendorOrderDetailResponse{
		Shipping:  shippingInfo.VendorOrderDetailShipping,
		Recipient: recipient,
		Store: domain.VendorOrderDetailStore{
			Name: storeName,
		},
		Items: items,
		Payment: domain.VendorOrderDetailPayment{
			Method:      paymentMethod,
			TotalAmount: order.GrandTotal,
		},
		Order: domain.VendorOrderDetailOrder{
			OrderNumber:   order.OrderNo,
			OrderTime:     order.PlacedAt,
			PaymentTime:   paidAt,
			ShippingTime:  shippingInfo.shippedAt,
			DeliveredTime: shippingInfo.DeliveredAt,
		},
		Summary: domain.VendorOrderDetailSummary{
			ItemsTotal:  order.Subtotal,
			ShippingFee: order.ShippingFee,
			ServiceFee:  order.PlatformFee,
			GrandTotal:  order.GrandTotal,
		},
	}, nil
}

// shippingInfoInternal holds shipping info with internal shippedAt for order timestamps.
type shippingInfoInternal struct {
	domain.VendorOrderDetailShipping
	shippedAt        *time.Time
	deliveryDetected bool
}

// buildShippingInfo builds shipping info by hitting tracking API if available.
func (uc *vendorOrderUseCase) buildShippingInfo(ctx context.Context, shipment *domain.Shipment, orderStatus string) shippingInfoInternal {
	result := shippingInfoInternal{
		VendorOrderDetailShipping: domain.VendorOrderDetailShipping{
			Courier:        "",
			TrackingNumber: "",
			Status:         orderStatusToShippingStatus(orderStatus),
			DeliveredAt:    nil,
		},
		shippedAt: nil,
	}

	if shipment == nil {
		return result
	}

	result.Courier = shipment.CourierCode
	result.TrackingNumber = shipment.TrackingNo
	result.DeliveredAt = shipment.DeliveredAt
	result.shippedAt = shipment.ShippedAt

	// Try to get real-time tracking status if tracking number exists.
	if shipment.TrackingNo != "" {
		trackResult, err := uc.rajaOngkir.TrackWaybill(ctx, shipment.TrackingNo, shipment.CourierCode)
		if err == nil && trackResult != nil && len(trackResult.Manifest) > 0 {
			// Sort manifest newest-first by date+time so index 0 is the latest event.
			sort.Slice(trackResult.Manifest, func(i, j int) bool {
				di := trackResult.Manifest[i].ManifestDate + " " + trackResult.Manifest[i].ManifestTime
				dj := trackResult.Manifest[j].ManifestDate + " " + trackResult.Manifest[j].ManifestTime
				return di > dj
			})

			// Get the latest manifest description as status.
			result.Status = trackResult.Manifest[0].ManifestDescription

			// Update delivered_at from tracking if delivered.
			if trackResult.Delivered && trackResult.DeliveryStatus.PodDate != "" {
				result.deliveryDetected = true
				deliveredAt := parseTrackingDateTime(trackResult.DeliveryStatus.PodDate, trackResult.DeliveryStatus.PodTime)
				if deliveredAt != nil {
					result.DeliveredAt = deliveredAt
				}
			}
		}
	}

	return result
}

// buildRecipientInfo parses shipping address snapshot into recipient info.
func (uc *vendorOrderUseCase) buildRecipientInfo(snapshot map[string]interface{}) domain.VendorOrderDetailRecipient {
	getString := func(key string) string {
		if v, ok := snapshot[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}

	return domain.VendorOrderDetailRecipient{
		Name:  getString("recipient_name"),
		Phone: getString("phone"),
		Address: domain.VendorOrderDetailRecipientAddress{
			Street:     getString("address_line"),
			District:   getString("district_name"),
			City:       getString("city_name"),
			Province:   getString("province_name"),
			PostalCode: getString("postal_code"),
		},
	}
}

// buildPaymentMethodDisplay returns the payment method from payment_invoices.
func buildPaymentMethodDisplay(method *string) string {
	if method == nil || *method == "" {
		return ""
	}
	return *method
}

// orderStatusToShippingStatus converts order status to a default shipping status description.
func orderStatusToShippingStatus(status string) string {
	switch status {
	case domain.OrderStatusPendingPayment:
		return "Menunggu pembayaran"
	case domain.OrderStatusPaid:
		return "Pembayaran diterima, menunggu konfirmasi penjual"
	case domain.OrderStatusProcessing:
		return "Pesanan sedang diproses"
	case domain.OrderStatusPacked:
		return "Pesanan sudah dikemas"
	case domain.OrderStatusShipped:
		return "Pesanan dalam pengiriman"
	case domain.OrderStatusReceived:
		return "Pesanan diterima"
	case domain.OrderStatusCompleted:
		return "Pesanan selesai"
	case domain.OrderStatusCanceled:
		return "Pesanan dibatalkan"
	case domain.OrderStatusRefunded:
		return "Pesanan direfund"
	default:
		return status
	}
}

// parseTrackingDateTime parses date and time strings from tracking API into time.Time.
func parseTrackingDateTime(dateStr, timeStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	dateTimeStr := dateStr
	if timeStr != "" {
		dateTimeStr = dateStr + " " + timeStr
	}

	// Try common formats.
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateTimeStr); err == nil {
			return &t
		}
	}

	return nil
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

// GetOrderInvoice returns invoice data for an order.
func (uc *vendorOrderUseCase) GetOrderInvoice(ctx context.Context, vendorID, orderID uuid.UUID) (*domain.VendorOrderInvoiceResponse, error) {
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

	items := make([]domain.InvoiceItem, 0, len(rawItems))
	for _, it := range rawItems {
		items = append(items, domain.InvoiceItem{
			Name:     it.ProductNameSnapshot,
			Variant:  it.SKUSnapshot,
			Price:    it.UnitPrice,
			Quantity: it.Qty,
			Subtotal: it.LineTotal,
			ImageURL: nil, // Can be populated if image URL is stored
		})
	}

	// Parse recipient from shipping address snapshot.
	recipient := uc.buildRecipientInfo(order.ShippingAddressSnapshot)

	// Fetch payment info.
	var paymentMethod string
	var paidAtStr string
	invoice, _ := uc.paymentRepo.FindInvoiceByOrderID(ctx, orderID)
	if invoice != nil {
		paymentMethod = buildPaymentMethodDisplay(invoice.PaymentMethod)
		if invoice.PaidAt != nil {
			paidAtStr = invoice.PaidAt.Format("2006-01-02 15:04")
		}
	}

	// Use external_invoice_id from payment_invoices when available.
	invoiceNumber := order.OrderNo
	if invoice != nil && invoice.ExternalInvoiceID != "" {
		invoiceNumber = invoice.ExternalInvoiceID
	}

	return &domain.VendorOrderInvoiceResponse{
		Invoice: domain.InvoiceInfo{
			InvoiceNumber: invoiceNumber,
			IssuedAt:      order.PlacedAt.Format("2006-01-02 15:04"),
			Actions: domain.InvoiceActions{
				CanCopy:     true,
				CanPrint:    true,
				CanDownload: true,
			},
		},
		Customer: domain.InvoiceCustomer{
			Name:  recipient.Name,
			Phone: recipient.Phone,
			Address: domain.InvoiceCustomerAddress{
				Street:     recipient.Address.Street,
				District:   recipient.Address.District,
				City:       recipient.Address.City,
				Province:   recipient.Address.Province,
				PostalCode: recipient.Address.PostalCode,
			},
		},
		Payment: domain.InvoicePayment{
			Method: paymentMethod,
			PaidAt: paidAtStr,
		},
		Items: items,
		Summary: domain.InvoiceSummary{
			Subtotal:    order.Subtotal,
			ShippingFee: order.ShippingFee,
			ServiceFee:  order.PlatformFee,
			Total:       order.GrandTotal,
			Currency:    "IDR",
		},
	}, nil
}

// GetOrderShippingInfo returns shipping/tracking info for an order.
func (uc *vendorOrderUseCase) GetOrderShippingInfo(ctx context.Context, vendorID, orderID uuid.UUID) (*domain.VendorOrderShippingInfoResponse, error) {
	order, err := uc.vendorOrderRepo.FindByIDAndVendor(ctx, orderID, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotBelongToVendor
	}

	// Fetch shipment.
	shipment, _ := uc.shipmentRepo.FindByOrderID(ctx, orderID)

	// Build response with default values.
	resp := &domain.VendorOrderShippingInfoResponse{
		EstimatedArrival: nil,
		Stages:           uc.buildShippingStages(order.OrderStatus),
		Courier: domain.ShippingInfoCourier{
			Name:       "",
			Code:       "",
			TrackingNo: "",
		},
		Events: []domain.ShippingInfoEvent{},
	}

	if shipment == nil {
		return resp, nil
	}

	// Set courier info.
	resp.Courier = domain.ShippingInfoCourier{
		Name:       courierCodeToName(shipment.CourierCode),
		Code:       shipment.CourierCode,
		TrackingNo: shipment.TrackingNo,
	}

	// Calculate estimated arrival if shipped.
	if shipment.ShippedAt != nil && shipment.ETD != "" {
		resp.EstimatedArrival = calculateEstimatedArrival(*shipment.ShippedAt, shipment.ETD)
	}

	// Try to get real-time tracking events if tracking number exists.
	if shipment.TrackingNo != "" {
		trackResult, err := uc.rajaOngkir.TrackWaybill(ctx, shipment.TrackingNo, shipment.CourierCode)
		if err == nil && trackResult != nil {
			// Build events from manifest.
			events := make([]domain.ShippingInfoEvent, 0, len(trackResult.Manifest))
			for _, m := range trackResult.Manifest {
				events = append(events, domain.ShippingInfoEvent{
					DateTime:    fmt.Sprintf("%s %s", m.ManifestDate, m.ManifestTime),
					Description: m.ManifestDescription,
					Location:    m.CityName,
				})
			}

			// Sort events newest-first by datetime.
			sort.Slice(events, func(i, j int) bool {
				return events[i].DateTime > events[j].DateTime
			})
			resp.Events = events

			// Determine stages from real-time tracking status.
			if trackResult.Delivered {
				resp.Stages = uc.buildShippingStages(domain.OrderStatusReceived)

				// Auto-update order status to "received" if courier confirms delivery.
				if order.OrderStatus == domain.OrderStatusShipped {
					notes := "Auto-updated: courier confirmed delivery"
					uc.orderRepo.MarkOrderReceived(ctx, orderID, nil, &notes)
					if trackResult.DeliveryStatus.PodDate != "" {
						deliveredAt := parseTrackingDateTime(trackResult.DeliveryStatus.PodDate, trackResult.DeliveryStatus.PodTime)
						if deliveredAt != nil {
							uc.shipmentRepo.UpdateDeliveredAt(ctx, orderID, *deliveredAt)
						}
					}
				}
			} else if len(events) > 0 {
				// Infer stage from the newest manifest event description.
				resp.Stages = uc.buildShippingStagesFromEvent(events[0].Description)
			}
		}
	}

	return resp, nil
}

// buildShippingStages builds shipping stage indicators based on order status.
func (uc *vendorOrderUseCase) buildShippingStages(orderStatus string) []domain.ShippingInfoStage {
	stages := []domain.ShippingInfoStage{
		{Status: "processing", Label: "Diproses", Completed: false, Active: false},
		{Status: "shipped", Label: "Dikirim", Completed: false, Active: false},
		{Status: "in_transit", Label: "Dalam Perjalanan", Completed: false, Active: false},
		{Status: "delivered", Label: "Tiba di Tujuan", Completed: false, Active: false},
	}

	switch orderStatus {
	case domain.OrderStatusProcessing, domain.OrderStatusPacked:
		stages[0].Active = true
	case domain.OrderStatusShipped:
		stages[0].Completed = true
		stages[1].Completed = true
		stages[2].Active = true
	case domain.OrderStatusReceived, domain.OrderStatusCompleted:
		stages[0].Completed = true
		stages[1].Completed = true
		stages[2].Completed = true
		stages[3].Completed = true
	}

	return stages
}

// buildShippingStagesFromEvent infers shipping stages from the latest tracking event description.
func (uc *vendorOrderUseCase) buildShippingStagesFromEvent(description string) []domain.ShippingInfoStage {
	desc := strings.ToLower(description)

	switch {
	case strings.Contains(desc, "diterima") || strings.Contains(desc, "delivered") || strings.Contains(desc, "received"):
		return uc.buildShippingStages(domain.OrderStatusReceived)
	case strings.Contains(desc, "dikirim ke alamat") || strings.Contains(desc, "out for delivery") || strings.Contains(desc, "akan dikirim ke alamat penerima"):
		// Last-mile delivery — in transit, almost delivered.
		stages := uc.buildShippingStages(domain.OrderStatusShipped)
		stages[2].Completed = true
		stages[3].Active = true
		return stages
	case strings.Contains(desc, "dikirimkan ke") || strings.Contains(desc, "transit") || strings.Contains(desc, "sampai di") || strings.Contains(desc, "gateway"):
		return uc.buildShippingStages(domain.OrderStatusShipped)
	default:
		// Manifes / picked up — shipped stage.
		stages := uc.buildShippingStages(domain.OrderStatusShipped)
		stages[1].Active = true
		stages[2].Active = false
		stages[2].Completed = false
		return stages
	}
}

// courierCodeToName converts courier code to display name.
func courierCodeToName(code string) string {
	names := map[string]string{
		"jne":      "JNE",
		"jnt":      "J&T Express",
		"sicepat":  "SiCepat",
		"pos":      "Pos Indonesia",
		"tiki":     "TIKI",
		"anteraja": "AnterAja",
		"ninja":    "Ninja Xpress",
		"lion":     "Lion Parcel",
		"idx":      "ID Express",
	}
	if name, ok := names[strings.ToLower(code)]; ok {
		return name
	}
	return strings.ToUpper(code)
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
