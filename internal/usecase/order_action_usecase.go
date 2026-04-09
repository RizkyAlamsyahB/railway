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
	orderRepo        domain.OrderRepository
	shipmentRepo     domain.ShipmentRepository
	paymentRepo      domain.PaymentRepository
	refundRepo       domain.CustomerRefundRepository
	returnReasonRepo domain.ReturnReasonRepository
	vendorRepo       domain.VendorRepository
	productRepo      domain.ProductRepository
	rajaOngkir       domain.RajaOngkirProvider
	storage          domain.StorageProvider
	xenditRefund     domain.XenditRefundProvider
	bankAccountRepo  domain.UserBankAccountRepository
	financeRepo      domain.FinanceRepository
}

// NewOrderActionUseCase creates a new OrderActionUseCase.
func NewOrderActionUseCase(
	orderRepo domain.OrderRepository,
	shipmentRepo domain.ShipmentRepository,
	paymentRepo domain.PaymentRepository,
	refundRepo domain.CustomerRefundRepository,
	returnReasonRepo domain.ReturnReasonRepository,
	vendorRepo domain.VendorRepository,
	productRepo domain.ProductRepository,
	rajaOngkir domain.RajaOngkirProvider,
	storage domain.StorageProvider,
	xenditRefund domain.XenditRefundProvider,
	bankAccountRepo domain.UserBankAccountRepository,
	financeRepo domain.FinanceRepository,
) domain.OrderActionUseCase {
	return &orderActionUseCase{
		orderRepo:        orderRepo,
		shipmentRepo:     shipmentRepo,
		paymentRepo:      paymentRepo,
		refundRepo:       refundRepo,
		returnReasonRepo: returnReasonRepo,
		vendorRepo:       vendorRepo,
		productRepo:      productRepo,
		rajaOngkir:       rajaOngkir,
		storage:          storage,
		xenditRefund:     xenditRefund,
		bankAccountRepo:  bankAccountRepo,
		financeRepo:      financeRepo,
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

	resp := &domain.CancelOrderResponse{
		OrderID:     orderID,
		OrderStatus: domain.OrderStatusCanceled,
		CanceledAt:  time.Now(),
	}

	// Auto-refund for paid orders
	if order.PaymentStatus == domain.PaymentStatusPaid {
		refundInfo := uc.initiateAutoRefund(ctx, order, "Order canceled by customer")
		if refundInfo != nil {
			resp.RefundInfo = refundInfo
		}
	}

	return resp, nil
}

func (uc *orderActionUseCase) RequestRefundByCustomer(ctx context.Context, userID, orderID uuid.UUID, req domain.CreateOrderRefundRequest) (*domain.CreateOrderRefundResponse, error) {
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

	if !isOrderEligibleForRefundRequest(order) {
		return nil, ErrOrderRefundNotEligible
	}

	if uc.refundRepo == nil || uc.returnReasonRepo == nil || uc.storage == nil {
		return nil, fmt.Errorf("refund dependencies are not configured")
	}

	hasOpenRefund, err := uc.refundRepo.HasOpenRefundByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing refund request: %w", err)
	}
	if hasOpenRefund {
		return nil, ErrRefundAlreadyRequested
	}

	invoice, err := uc.paymentRepo.FindInvoiceByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find payment invoice: %w", err)
	}
	if invoice == nil {
		return nil, ErrInvoiceNotFound
	}

	returnReason, err := uc.returnReasonRepo.FindByID(ctx, req.ReturnReasonID)
	if err != nil {
		return nil, fmt.Errorf("failed to find return reason: %w", err)
	}
	if returnReason == nil {
		return nil, ErrRefundReasonNotFound
	}

	channelCode := strings.TrimSpace(req.DestinationChannelCode)
	if channelCode == "" {
		channelCode = domain.PayoutChannelIDBCA
	}

	reason := strings.TrimSpace(returnReason.Reason)
	description := strings.TrimSpace(req.Description)
	var descriptionPtr *string
	if description != "" {
		descriptionPtr = &description
	}

	accountNumber := strings.TrimSpace(req.DestinationAccountNumber)
	accountHolderName := strings.TrimSpace(req.DestinationAccountHolderName)
	bankName := strings.TrimSpace(req.DestinationBankName)
	if accountNumber == "" || accountHolderName == "" || bankName == "" {
		return nil, ErrRefundDestinationRequired
	}

	evidences := make([]domain.CreateCustomerRefundEvidenceInput, 0, len(req.ImageObjectKeys)+1)
	imageURLs := make([]string, 0, len(req.ImageObjectKeys))
	for i, objectKey := range req.ImageObjectKeys {
		key := strings.TrimSpace(objectKey)
		if key == "" {
			continue
		}

		info, err := uc.storage.HeadObject(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to verify refund image object: %w", err)
		}
		if info == nil {
			return nil, ErrRefundEvidenceNotUploaded
		}

		contentType := normalizeContentType(info.ContentType)
		if !isRefundEvidenceImageContentType(contentType) {
			return nil, ErrInvalidRefundEvidenceContentType
		}
		if info.ContentLength > RefundEvidenceMaxBytes {
			return nil, ErrRefundEvidenceTooLarge
		}

		fileSize := int(info.ContentLength)
		evidences = append(evidences, domain.CreateCustomerRefundEvidenceInput{
			ObjectKey:     key,
			MimeType:      contentType,
			FileSizeBytes: fileSize,
			MediaType:     "image",
			SortOrder:     i,
		})
		imageURLs = append(imageURLs, uc.storage.GetURL(key))
	}
	if len(evidences) > 5 {
		return nil, ErrTooManyRefundEvidenceImages
	}

	var videoURL *string
	if req.VideoObjectKey != nil && strings.TrimSpace(*req.VideoObjectKey) != "" {
		key := strings.TrimSpace(*req.VideoObjectKey)
		info, err := uc.storage.HeadObject(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("failed to verify refund video object: %w", err)
		}
		if info == nil {
			return nil, ErrRefundEvidenceNotUploaded
		}

		contentType := normalizeContentType(info.ContentType)
		if !isRefundEvidenceVideoContentType(contentType) {
			return nil, ErrInvalidRefundEvidenceContentType
		}
		if info.ContentLength > RefundEvidenceMaxBytes {
			return nil, ErrRefundEvidenceTooLarge
		}

		fileSize := int(info.ContentLength)
		evidences = append(evidences, domain.CreateCustomerRefundEvidenceInput{
			ObjectKey:     key,
			MimeType:      contentType,
			FileSizeBytes: fileSize,
			MediaType:     "video",
			SortOrder:     len(evidences),
		})
		url := uc.storage.GetURL(key)
		videoURL = &url
	}

	now := time.Now().UTC()
	resp, err := uc.refundRepo.CreateRefundRequest(ctx, domain.CreateCustomerRefundInput{
		OrderID:                       orderID,
		PaymentInvoiceID:              invoice.ID,
		Amount:                        order.GrandTotal,
		ReturnReasonID:                req.ReturnReasonID,
		Reason:                        reason,
		Description:                   descriptionPtr,
		RequestedBy:                   userID,
		RequestedAt:                   now,
		DestinationChannelCode:        channelCode,
		DestinationBankName:           bankName,
		DestinationAccountNumber:      accountNumber,
		DestinationAccountHolderName:  accountHolderName,
		DestinationAccountNumberLast4: maskAccountLast4(accountNumber),
		Evidences:                     evidences,
	})
	if err != nil {
		return nil, err
	}

	if resp != nil {
		resp.ImageURLs = imageURLs
		resp.VideoURL = videoURL
	}

	return resp, nil
}

func (uc *orderActionUseCase) PresignRefundEvidenceByCustomer(ctx context.Context, userID, orderID uuid.UUID, req domain.PresignRefundEvidenceRequest) (*domain.PresignRefundEvidenceResponse, error) {
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
	if uc.storage == nil {
		return nil, fmt.Errorf("storage dependency is not configured")
	}

	ct := normalizeContentType(req.ContentType)
	if !isAllowedRefundEvidenceContentType(ct) {
		return nil, ErrInvalidRefundEvidenceContentType
	}

	objectKey := fmt.Sprintf("refund-evidences/%s/%s/%s", orderID.String(), time.Now().Format("2006/01/02"), uuid.New().String())
	uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, ct, PresignedUploadExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refund evidence presigned upload URL: %w", err)
	}

	return &domain.PresignRefundEvidenceResponse{
		UploadURL:   uploadURL,
		ObjectKey:   objectKey,
		ContentType: ct,
		ExpiresIn:   int(PresignedUploadExpiry.Seconds()),
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
			OrderItemID: it.ID,
			Name:        it.ProductNameSnapshot,
			Variant:     it.SKUSnapshot,
			ImageURL:    imageMap[it.ProductVariantID],
			Price:       it.UnitPrice,
			Qty:         it.Qty,
			Subtotal:    it.LineTotal,
		})
	}

	// Fetch vendor/store info.
	store := domain.CustomerOrderDetailStore{VendorID: order.VendorID}
	vendor, err := uc.vendorRepo.FindByID(ctx, order.VendorID)
	if err == nil && vendor != nil {
		store.Name = vendorDisplayNameOrEmpty(vendor)
		store.OwnerUserID = vendor.OwnerUserID
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

	// Build refund info (if any) — applies to canceled, received, and completed orders.
	var refundInfo *domain.CustomerOrderDetailRefund
	if order.OrderStatus == domain.OrderStatusCanceled ||
		order.OrderStatus == domain.OrderStatusReceived ||
		order.OrderStatus == domain.OrderStatusCompleted {
		refundDetail, _ := uc.refundRepo.FindLatestByOrderID(ctx, orderID)
		if refundDetail != nil {
			method := ""
			if refundDetail.RefundMethod != nil {
				method = *refundDetail.RefundMethod
			}
			info := &domain.CustomerOrderDetailRefund{
				RefundID:     refundDetail.ID,
				Status:       refundDetail.Status,
				RefundMethod: method,
				Amount:       refundDetail.Amount,
			}
			switch refundDetail.Status {
			case domain.RefundStatusAwaitingDestination:
				info.Message = "Silakan pilih rekening tujuan untuk refund"
				actions.NeedsRefundDestination = true
			case domain.RefundStatusProcessing:
				info.Message = "Refund sedang diproses"
			case domain.RefundStatusRequested:
				info.Message = "Refund sedang ditinjau"
			case domain.RefundStatusProcessed:
				info.Message = "Refund telah berhasil diproses"
			case domain.RefundStatusRejected:
				info.Message = "Refund ditolak"
			default:
				info.Message = "Refund sedang diproses"
			}
			if refundDetail.DestinationBankName != nil {
				info.DestinationBankName = *refundDetail.DestinationBankName
			}
			if refundDetail.DestinationAccountNumberLast4 != nil {
				info.DestinationAccountNumberLast4 = *refundDetail.DestinationAccountNumberLast4
			}
			refundInfo = info
		}

		// Set CanRefund: order is eligible and no open refund exists yet.
		if isOrderEligibleForRefundRequest(order) && refundDetail == nil {
			actions.CanRefund = true
		}
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
			Status:        order.OrderStatus,
			PaymentStatus: order.PaymentStatus,
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
		Refund:  refundInfo,
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
			OrderID:       order.ID,
			Items:         items,
			TotalPayment:  order.GrandTotal,
			OrderDate:     order.PlacedAt,
			Status:        order.OrderStatus,
			PaymentStatus: order.PaymentStatus,
		})
	}

	return result, meta, nil
}

func (uc *orderActionUseCase) ListOrderStatuses(_ context.Context) []string {
	result := make([]string, len(domain.OrderStatuses))
	copy(result, domain.OrderStatuses)
	return result
}

func isOrderEligibleForRefundRequest(order *domain.Order) bool {
	if order == nil {
		return false
	}
	if order.PaymentStatus != domain.PaymentStatusPaid {
		return false
	}
	return order.OrderStatus == domain.OrderStatusReceived || order.OrderStatus == domain.OrderStatusCompleted
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
				if isAbsoluteURL(url) {
					urlCopy := url
					result[it.ProductVariantID] = &urlCopy
					continue
				}

				if uc.storage == nil {
					continue
				}

				presignedURL, err := uc.storage.GeneratePresignedURL(ctx, url, PresignedDownloadExpiry)
				if err != nil {
					continue
				}

				urlCopy := presignedURL
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

// initiateAutoRefund creates a refund record and, for gateway-eligible payment methods,
// calls the Xendit Refund API. For VA payments, if the customer has a default bank account
// it stays awaiting_destination so they can confirm destination.
// This is a best-effort helper — errors are logged but do not fail the parent operation.
func (uc *orderActionUseCase) initiateAutoRefund(ctx context.Context, order *domain.Order, reason string) *domain.CancelOrderRefundInfo {
	invoice, err := uc.paymentRepo.FindInvoiceByOrderID(ctx, order.ID)
	if err != nil || invoice == nil {
		fmt.Printf("WARN: failed to find invoice for order %s: %v\n", order.ID, err)
		return nil
	}

	strategy := resolveRefundStrategy(invoice.PaymentMethod, invoice.PaymentChannel)

	refundInput := domain.CreateCustomerRefundInput{
		OrderID:          order.ID,
		PaymentInvoiceID: invoice.ID,
		Amount:           order.GrandTotal,
		Reason:           reason,
		RequestedBy:      order.UserID,
		RequestedAt:      time.Now().UTC(),
	}

	switch {
	case strategy == refundStrategyQRGateway || strategy == refundStrategyEWalletGateway:
		// Gateway refund via Xendit Refund API
		refundInput.Status = domain.RefundStatusProcessing
		refundInput.RefundMethod = domain.RefundMethodGateway

		resp, err := uc.refundRepo.CreateRefundRequest(ctx, refundInput)
		if err != nil {
			fmt.Printf("WARN: failed to create refund record for order %s: %v\n", order.ID, err)
			return nil
		}

		// Call Xendit Refund API
		gatewayInfo := uc.callXenditGatewayRefund(ctx, order, invoice, resp.RefundID)
		if gatewayInfo != nil {
			return gatewayInfo
		}

		return &domain.CancelOrderRefundInfo{
			RefundID:     resp.RefundID,
			RefundMethod: domain.RefundMethodGateway,
			Status:       domain.RefundStatusProcessing,
			Message:      "Refund sedang diproses melalui metode pembayaran asal",
		}

	default:
		// VA / bank transfer / other: disbursement path
		refundInput.Status = domain.RefundStatusAwaitingDestination
		refundInput.RefundMethod = domain.RefundMethodDisbursement

		resp, err := uc.refundRepo.CreateRefundRequest(ctx, refundInput)
		if err != nil {
			fmt.Printf("WARN: failed to create refund record for order %s: %v\n", order.ID, err)
			return nil
		}

		return &domain.CancelOrderRefundInfo{
			RefundID:     resp.RefundID,
			RefundMethod: domain.RefundMethodDisbursement,
			Status:       domain.RefundStatusAwaitingDestination,
			Message:      "Silakan pilih rekening tujuan untuk refund",
		}
	}
}

// callXenditGatewayRefund calls the Xendit Refund API for gateway-eligible payments.
func (uc *orderActionUseCase) callXenditGatewayRefund(ctx context.Context, order *domain.Order, invoice *domain.PaymentInvoice, refundID uuid.UUID) *domain.CancelOrderRefundInfo {
	vendor, err := uc.vendorRepo.FindByID(ctx, order.VendorID)
	if err != nil || vendor == nil || vendor.XenditAccountID == nil {
		fmt.Printf("WARN: cannot get vendor xendit account for order %s: %v\n", order.ID, err)
		return nil
	}

	referenceID := fmt.Sprintf("refund-%s", refundID.String())
	idempotencyKey := fmt.Sprintf("%s-%d", referenceID, time.Now().UTC().UnixNano())

	xenditReq := domain.XenditRefundRequest{
		InvoiceID:   strings.TrimSpace(derefString(invoice.XenditInvoiceID)),
		ReferenceID: referenceID,
		Amount:      order.GrandTotal,
		Currency:    invoice.Currency,
		Reason:      "CANCELLATION",
	}

	resp, err := uc.xenditRefund.CreateRefund(ctx, *vendor.XenditAccountID, idempotencyKey, xenditReq)
	if err != nil {
		fmt.Printf("WARN: xendit refund API failed for order %s: %v\n", order.ID, err)
		return nil
	}

	// Update refund record with xendit refund ID
	if err := uc.financeRepo.SetRefundGatewayInitiated(ctx, refundID, resp.ID, domain.RefundMethodGateway, referenceID); err != nil {
		fmt.Printf("WARN: failed to update refund gateway initiated for order %s: %v\n", order.ID, err)
	}

	// Update order payment_status to refunded if already succeeded
	if resp.Status == domain.XenditRefundStatusSucceeded {
		if err := uc.financeRepo.UpdateOrderPaymentStatus(ctx, order.ID, domain.PaymentStatusRefunded); err != nil {
			fmt.Printf("WARN: failed to update payment status for order %s: %v\n", order.ID, err)
		}
	}

	return nil
}

func (uc *orderActionUseCase) SubmitRefundDestination(ctx context.Context, userID, orderID uuid.UUID, req domain.SubmitRefundDestinationRequest) (*domain.SubmitRefundDestinationResponse, error) {
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

	refund, err := uc.refundRepo.FindAwaitingDestinationByOrderAndUser(ctx, orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find refund: %w", err)
	}
	if refund == nil {
		return nil, ErrRefundNotFound
	}

	var channelCode, bankName, accountNumber, accountHolderName, accountLast4 string
	var bankAccountID *uuid.UUID

	if req.UserBankAccountID != nil {
		// Use existing saved bank account
		ba, err := uc.bankAccountRepo.FindByID(ctx, *req.UserBankAccountID)
		if err != nil {
			return nil, fmt.Errorf("failed to find bank account: %w", err)
		}
		if ba == nil {
			return nil, ErrBankAccountNotFound
		}
		if ba.UserID != userID {
			return nil, ErrBankAccountNotOwned
		}
		channelCode = ba.ChannelCode
		bankName = ba.BankName
		accountNumber = ba.AccountNumber
		accountHolderName = ba.AccountHolderName
		accountLast4 = ba.AccountLast4
		bankAccountID = req.UserBankAccountID
	} else {
		// Manual entry
		channelCode = strings.TrimSpace(req.DestinationChannelCode)
		bankName = strings.TrimSpace(req.DestinationBankName)
		accountNumber = strings.TrimSpace(req.DestinationAccountNumber)
		accountHolderName = strings.TrimSpace(req.DestinationAccountHolderName)
		if accountNumber == "" || accountHolderName == "" {
			return nil, ErrRefundDestinationRequired
		}
		accountLast4 = accountNumber
		if len(accountLast4) > 4 {
			accountLast4 = accountLast4[len(accountLast4)-4:]
		}
	}

	if err := uc.refundRepo.SubmitRefundDestination(ctx, refund.ID, channelCode, bankName, accountNumber, accountHolderName, accountLast4, bankAccountID); err != nil {
		return nil, fmt.Errorf("failed to submit refund destination: %w", err)
	}

	return &domain.SubmitRefundDestinationResponse{
		RefundID: refund.ID,
		Status:   domain.RefundStatusRequested,
		Message:  "Refund destination submitted, pending processing",
	}, nil
}
