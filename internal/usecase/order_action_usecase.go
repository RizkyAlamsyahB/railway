package usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type orderActionUseCase struct {
	orderRepo    domain.OrderRepository
	shipmentRepo domain.ShipmentRepository
	paymentRepo  domain.PaymentRepository
	vendorRepo   domain.VendorRepository
	productRepo  domain.ProductRepository
	rajaOngkir   domain.RajaOngkirProvider
}

// NewOrderActionUseCase creates a new OrderActionUseCase.
func NewOrderActionUseCase(
	orderRepo domain.OrderRepository,
	shipmentRepo domain.ShipmentRepository,
	paymentRepo domain.PaymentRepository,
	vendorRepo domain.VendorRepository,
	productRepo domain.ProductRepository,
	rajaOngkir domain.RajaOngkirProvider,
) domain.OrderActionUseCase {
	return &orderActionUseCase{
		orderRepo:    orderRepo,
		shipmentRepo: shipmentRepo,
		paymentRepo:  paymentRepo,
		vendorRepo:   vendorRepo,
		productRepo:  productRepo,
		rajaOngkir:   rajaOngkir,
	}
}

func (uc *orderActionUseCase) CompleteByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*domain.ReceiveOrderResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}

	if order.OrderStatus != domain.OrderStatusReceived {
		return nil, ErrInvalidOrderCompletionTransition
	}

	notes := "Order completed by customer"
	applied, err := uc.orderRepo.MarkOrderCompleted(ctx, orderID, &userID, &notes)
	if err != nil {
		return nil, fmt.Errorf("failed to mark order completed: %w", err)
	}
	if !applied {
		return nil, ErrInvalidOrderCompletionTransition
	}

	return &domain.ReceiveOrderResponse{
		OrderID:       orderID,
		OrderStatus:   domain.OrderStatusCompleted,
		PaymentStatus: order.PaymentStatus,
		CompletedAt:   time.Now(),
	}, nil
}

func (uc *orderActionUseCase) CancelByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*domain.CancelOrderResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}

	switch order.OrderStatus {
	case domain.OrderStatusPendingPayment, domain.OrderStatusPaid, domain.OrderStatusProcessing:
		// allowed
	default:
		return nil, ErrInvalidOrderCancelTransition
	}

	notes := "Order canceled by customer"
	if err := uc.orderRepo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusCanceled, order.PaymentStatus, &userID, &notes); err != nil {
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}

	// Restore stock for canceled order (only if payment was made / stock was decremented at checkout).
	if err := uc.orderRepo.RestoreStock(ctx, orderID); err != nil {
		fmt.Printf("WARN: failed to restore stock for order %s: %v\n", orderID, err)
	}

	return &domain.CancelOrderResponse{
		OrderID:     orderID,
		OrderStatus: domain.OrderStatusCanceled,
		CanceledAt:  time.Now(),
	}, nil
}

func (uc *orderActionUseCase) GetOrderDetail(ctx context.Context, userID, orderID uuid.UUID) (*domain.CustomerOrderDetailResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}

	// Fetch order items.
	rawItems, err := uc.orderRepo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	// Resolve product images.
	imageMap := uc.resolveProductImages(ctx, rawItems)

	items := make([]domain.CustomerOrderDetailItem, 0, len(rawItems))
	for _, it := range rawItems {
		items = append(items, domain.CustomerOrderDetailItem{
			Name:     it.ProductNameSnapshot,
			Variant:  it.SKUSnapshot,
			ImageURL: imageMap[it.ProductVariantID],
			Price:    it.UnitPrice,
			Qty:      it.Qty,
			Subtotal: it.LineTotal,
		})
	}

	// Fetch vendor/store info.
	store := domain.CustomerOrderDetailStore{VendorID: order.VendorID}
	vendor, err := uc.vendorRepo.FindByID(ctx, order.VendorID)
	if err == nil && vendor != nil {
		store.Name = vendor.DisplayName
	}

	// Parse recipient address.
	address := uc.buildCustomerAddress(order.ShippingAddressSnapshot)

	// Fetch shipment (may be nil).
	shipment, _ := uc.shipmentRepo.FindByOrderID(ctx, orderID)

	// Build shipping info.
	var shipping *domain.CustomerOrderDetailShipping
	var shippedAt *time.Time
	var deliveredAt *time.Time
	if shipment != nil && shipment.TrackingNo != "" {
		shippingInfo := uc.buildCustomerShippingInfo(ctx, shipment, order.OrderStatus)
		shipping = &shippingInfo.CustomerOrderDetailShipping
		shippedAt = shipment.ShippedAt
		deliveredAt = shippingInfo.deliveredAt

		// Auto-update to received if courier confirms delivery.
		if shippingInfo.deliveryDetected && order.OrderStatus == domain.OrderStatusShipped {
			notes := "Auto-updated: courier confirmed delivery"
			uc.orderRepo.MarkOrderReceived(ctx, orderID, nil, &notes)
			if deliveredAt != nil {
				uc.shipmentRepo.UpdateDeliveredAt(ctx, orderID, *deliveredAt)
			}
		}
	}

	// Fetch payment info.
	var paymentMethod string
	var invoiceID *string
	var paidAt *time.Time
	invoice, _ := uc.paymentRepo.FindInvoiceByOrderID(ctx, orderID)
	if invoice != nil {
		paymentMethod = buildPaymentMethodDisplay(invoice.PaymentMethod)
		if invoice.ExternalInvoiceID != "" {
			invoiceID = &invoice.ExternalInvoiceID
		}
		paidAt = invoice.PaidAt
	}

	// Build actions.
	actions := domain.CustomerOrderActions{
		CanCancel:   canCancelOrder(order.OrderStatus),
		CanComplete: order.OrderStatus == domain.OrderStatusReceived,
		CanReview:   order.OrderStatus == domain.OrderStatusCompleted,
		CanContact:  order.OrderStatus != domain.OrderStatusPendingPayment && order.OrderStatus != domain.OrderStatusCanceled,
	}

	return &domain.CustomerOrderDetailResponse{
		Shipping: shipping,
		Address:  address,
		Store:    store,
		Items:    items,
		Payment: domain.CustomerOrderDetailPayment{
			Method:    paymentMethod,
			InvoiceID: invoiceID,
		},
		Order: domain.CustomerOrderDetailMeta{
			OrderNo:       order.OrderNo,
			OrderTime:     order.PlacedAt,
			PaymentTime:   paidAt,
			ShippingTime:  shippedAt,
			DeliveredTime: deliveredAt,
		},
		Summary: domain.CustomerOrderDetailSummary{
			Subtotal:    order.Subtotal,
			ShippingFee: order.ShippingFee,
			ServiceFee:  order.PlatformFee,
			GrandTotal:  order.GrandTotal,
		},
		Actions: actions,
	}, nil
}

func (uc *orderActionUseCase) GetOrderShippingInfo(ctx context.Context, userID, orderID uuid.UUID) (*domain.VendorOrderShippingInfoResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}

	shipment, _ := uc.shipmentRepo.FindByOrderID(ctx, orderID)

	resp := &domain.VendorOrderShippingInfoResponse{
		EstimatedArrival: nil,
		Stages:           buildCustomerShippingStages(order.OrderStatus),
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

	resp.Courier = domain.ShippingInfoCourier{
		Name:       courierCodeToName(shipment.CourierCode),
		Code:       shipment.CourierCode,
		TrackingNo: shipment.TrackingNo,
	}

	if shipment.ShippedAt != nil && shipment.ETD != "" {
		resp.EstimatedArrival = calculateEstimatedArrival(*shipment.ShippedAt, shipment.ETD)
	}

	if shipment.TrackingNo != "" {
		trackResult, err := uc.rajaOngkir.TrackWaybill(ctx, shipment.TrackingNo, shipment.CourierCode)
		if err == nil && trackResult != nil {
			events := make([]domain.ShippingInfoEvent, 0, len(trackResult.Manifest))
			for _, m := range trackResult.Manifest {
				events = append(events, domain.ShippingInfoEvent{
					DateTime:    fmt.Sprintf("%s %s", m.ManifestDate, m.ManifestTime),
					Description: m.ManifestDescription,
					Location:    m.CityName,
				})
			}
			sort.Slice(events, func(i, j int) bool {
				return events[i].DateTime > events[j].DateTime
			})
			resp.Events = events

			if trackResult.Delivered {
				resp.Stages = buildCustomerShippingStages(domain.OrderStatusReceived)

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
				resp.Stages = buildCustomerShippingStagesFromEvent(events[0].Description)
			}
		}
	}

	return resp, nil
}

func (uc *orderActionUseCase) GetOrderInvoice(ctx context.Context, userID, orderID uuid.UUID) (*domain.VendorOrderInvoiceResponse, error) {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotOwned
	}

	rawItems, err := uc.orderRepo.FindItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	imageMap := uc.resolveProductImages(ctx, rawItems)

	items := make([]domain.InvoiceItem, 0, len(rawItems))
	for _, it := range rawItems {
		items = append(items, domain.InvoiceItem{
			Name:     it.ProductNameSnapshot,
			Variant:  it.SKUSnapshot,
			Price:    it.UnitPrice,
			Quantity: it.Qty,
			Subtotal: it.LineTotal,
			ImageURL: imageMap[it.ProductVariantID],
		})
	}

	recipient := uc.buildCustomerAddress(order.ShippingAddressSnapshot)

	var paymentMethod string
	var paidAtStr string
	invoice, _ := uc.paymentRepo.FindInvoiceByOrderID(ctx, orderID)
	if invoice != nil {
		paymentMethod = buildPaymentMethodDisplay(invoice.PaymentMethod)
		if invoice.PaidAt != nil {
			paidAtStr = invoice.PaidAt.Format("2006-01-02 15:04")
		}
	}

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
				Street:     recipient.Street,
				District:   recipient.District,
				City:       recipient.City,
				Province:   recipient.Province,
				PostalCode: recipient.PostalCode,
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

func (uc *orderActionUseCase) ListByCustomer(ctx context.Context, userID uuid.UUID, params domain.CustomerOrderListParams) ([]domain.CustomerOrderListItem, *domain.PaginationMeta, error) {
	params = normalizeCustomerOrderListParams(params)
	if params.Status != "" && !isValidOrderStatus(params.Status) {
		return nil, nil, ErrInvalidOrderStatus
	}

	orders, total, err := uc.orderRepo.ListByUser(ctx, userID, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list orders: %w", err)
	}

	meta := buildOrderPaginationMeta(params.Page, params.Limit, total)
	if len(orders) == 0 {
		return []domain.CustomerOrderListItem{}, meta, nil
	}

	orderIDs := make([]uuid.UUID, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	itemsByOrderID, err := uc.orderRepo.FindItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list order items: %w", err)
	}

	// Collect all order items to resolve images.
	var allItems []domain.OrderItem
	for _, items := range itemsByOrderID {
		allItems = append(allItems, items...)
	}
	imageMap := uc.resolveProductImages(ctx, allItems)

	result := make([]domain.CustomerOrderListItem, 0, len(orders))
	for _, order := range orders {
		rawItems := itemsByOrderID[order.ID]
		items := make([]domain.CustomerOrderProductItem, 0, len(rawItems))
		for _, it := range rawItems {
			items = append(items, domain.CustomerOrderProductItem{
				ProductName:     it.ProductNameSnapshot,
				SelectedVariant: it.SKUSnapshot,
				ImageURL:        imageMap[it.ProductVariantID],
				Qty:             it.Qty,
				UnitPrice:       it.UnitPrice,
			})
		}

		result = append(result, domain.CustomerOrderListItem{
			OrderID:      order.ID,
			Items:        items,
			TotalPayment: order.GrandTotal,
			OrderDate:    order.PlacedAt,
			Status:       order.OrderStatus,
		})
	}

	return result, meta, nil
}

func (uc *orderActionUseCase) ListOrderStatuses(_ context.Context) []string {
	result := make([]string, len(domain.OrderStatuses))
	copy(result, domain.OrderStatuses)
	return result
}

// --- helpers ---

// resolveProductImages maps product_variant_id → primary image URL for order items.
func (uc *orderActionUseCase) resolveProductImages(ctx context.Context, items []domain.OrderItem) map[uuid.UUID]*string {
	result := make(map[uuid.UUID]*string, len(items))

	// Collect unique variant IDs.
	variantIDs := make(map[uuid.UUID]bool)
	for _, it := range items {
		variantIDs[it.ProductVariantID] = true
	}

	// Look up variant → product mapping.
	variantToProduct := make(map[uuid.UUID]uuid.UUID)
	productIDs := make(map[uuid.UUID]bool)
	for vid := range variantIDs {
		variant, err := uc.productRepo.FindVariantByID(ctx, vid)
		if err == nil && variant != nil {
			variantToProduct[vid] = variant.ProductID
			productIDs[variant.ProductID] = true
		}
	}

	// Look up primary image for each product.
	productImage := make(map[uuid.UUID]string)
	for pid := range productIDs {
		images, err := uc.productRepo.FindImagesByProductID(ctx, pid)
		if err == nil {
			for _, img := range images {
				if img.IsPrimary {
					productImage[pid] = img.ImageURL
					break
				}
			}
			// Fallback to first image if no primary.
			if _, ok := productImage[pid]; !ok && len(images) > 0 {
				productImage[pid] = images[0].ImageURL
			}
		}
	}

	// Map variant ID → image URL.
	for _, it := range items {
		if pid, ok := variantToProduct[it.ProductVariantID]; ok {
			if url, ok := productImage[pid]; ok {
				urlCopy := url
				result[it.ProductVariantID] = &urlCopy
			}
		}
	}

	return result
}

// buildCustomerAddress parses shipping address snapshot into customer address DTO.
func (uc *orderActionUseCase) buildCustomerAddress(snapshot map[string]interface{}) domain.CustomerOrderDetailAddress {
	getString := func(key string) string {
		if v, ok := snapshot[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}

	return domain.CustomerOrderDetailAddress{
		Name:       getString("recipient_name"),
		Phone:      getString("phone"),
		Street:     getString("address_line"),
		District:   getString("district_name"),
		City:       getString("city_name"),
		Province:   getString("province_name"),
		PostalCode: getString("postal_code"),
	}
}

// customerShippingInfoInternal holds parsed shipping info with internal fields.
type customerShippingInfoInternal struct {
	domain.CustomerOrderDetailShipping
	deliveredAt      *time.Time
	deliveryDetected bool
}

// buildCustomerShippingInfo builds shipping info by checking tracking API.
func (uc *orderActionUseCase) buildCustomerShippingInfo(ctx context.Context, shipment *domain.Shipment, orderStatus string) customerShippingInfoInternal {
	result := customerShippingInfoInternal{
		CustomerOrderDetailShipping: domain.CustomerOrderDetailShipping{
			Courier:        courierCodeToName(shipment.CourierCode),
			TrackingNumber: shipment.TrackingNo,
			Status:         orderStatusToShippingStatus(orderStatus),
		},
		deliveredAt: shipment.DeliveredAt,
	}

	if shipment.TrackingNo != "" {
		trackResult, err := uc.rajaOngkir.TrackWaybill(ctx, shipment.TrackingNo, shipment.CourierCode)
		if err == nil && trackResult != nil && len(trackResult.Manifest) > 0 {
			sort.Slice(trackResult.Manifest, func(i, j int) bool {
				di := trackResult.Manifest[i].ManifestDate + " " + trackResult.Manifest[i].ManifestTime
				dj := trackResult.Manifest[j].ManifestDate + " " + trackResult.Manifest[j].ManifestTime
				return di > dj
			})

			result.Status = trackResult.Manifest[0].ManifestDescription
			statusDate := fmt.Sprintf("%s %s", trackResult.Manifest[0].ManifestDate, trackResult.Manifest[0].ManifestTime)
			result.StatusDate = &statusDate

			if trackResult.Delivered && trackResult.DeliveryStatus.PodDate != "" {
				result.deliveryDetected = true
				deliveredAt := parseTrackingDateTime(trackResult.DeliveryStatus.PodDate, trackResult.DeliveryStatus.PodTime)
				if deliveredAt != nil {
					result.deliveredAt = deliveredAt
				}
			}
		}
	}

	return result
}

func canCancelOrder(status string) bool {
	switch status {
	case domain.OrderStatusPendingPayment, domain.OrderStatusPaid, domain.OrderStatusProcessing:
		return true
	default:
		return false
	}
}

// buildCustomerShippingStages builds shipping stages for customer shipping info.
func buildCustomerShippingStages(orderStatus string) []domain.ShippingInfoStage {
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

// buildCustomerShippingStagesFromEvent infers stages from tracking event description.
func buildCustomerShippingStagesFromEvent(description string) []domain.ShippingInfoStage {
	desc := strings.ToLower(description)

	switch {
	case strings.Contains(desc, "diterima") || strings.Contains(desc, "delivered") || strings.Contains(desc, "received"):
		return buildCustomerShippingStages(domain.OrderStatusReceived)
	case strings.Contains(desc, "dikirim ke alamat") || strings.Contains(desc, "out for delivery") || strings.Contains(desc, "akan dikirim ke alamat penerima"):
		stages := buildCustomerShippingStages(domain.OrderStatusShipped)
		stages[2].Completed = true
		stages[3].Active = true
		return stages
	case strings.Contains(desc, "dikirimkan ke") || strings.Contains(desc, "transit") || strings.Contains(desc, "sampai di") || strings.Contains(desc, "gateway"):
		return buildCustomerShippingStages(domain.OrderStatusShipped)
	default:
		stages := buildCustomerShippingStages(domain.OrderStatusShipped)
		stages[1].Active = true
		stages[2].Active = false
		stages[2].Completed = false
		return stages
	}
}

func normalizeCustomerOrderListParams(params domain.CustomerOrderListParams) domain.CustomerOrderListParams {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}
	params.Status = strings.ToLower(strings.TrimSpace(params.Status))
	return params
}

func buildOrderPaginationMeta(page, limit int, total int64) *domain.PaginationMeta {
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

func isValidOrderStatus(status string) bool {
	for _, allowed := range domain.OrderStatuses {
		if status == allowed {
			return true
		}
	}
	return false
}
