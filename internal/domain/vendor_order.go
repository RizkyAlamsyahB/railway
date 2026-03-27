package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Shipment entity ---

// Shipment represents the shipments table (1:1 with orders).
type Shipment struct {
	ID             uuid.UUID  `json:"id"`
	OrderID        uuid.UUID  `json:"order_id"`
	CourierCode    string     `json:"courier_code"`
	ServiceType    string     `json:"service_type"`
	TrackingNo     string     `json:"tracking_no"`
	ETD            string     `json:"etd"`
	ShipmentStatus string     `json:"shipment_status"`
	ShippedAt      *time.Time `json:"shipped_at"`
	DeliveredAt    *time.Time `json:"delivered_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// --- Vendor Order DTOs ---

// VendorOrderListParams holds query parameters for vendor order listing.
type VendorOrderListParams struct {
	Page      int
	Limit     int
	Status    string
	Search    string     // Search by order_no or customer_name
	DateFrom  *time.Time // Filter orders from this date (inclusive)
	DateTo    *time.Time // Filter orders to this date (inclusive)
	SortBy    string     // Sort field: order_date, customer_name, total_payment
	SortOrder string     // Sort order: asc, desc
}

// VendorOrderSortBy constants.
const (
	VendorOrderSortByOrderDate    = "order_date"
	VendorOrderSortByCustomerName = "customer_name"
	VendorOrderSortByTotalPayment = "total_payment"
)

// VendorOrderSortOrder constants.
const (
	VendorOrderSortOrderAsc  = "asc"
	VendorOrderSortOrderDesc = "desc"
)

// VendorOrderExportItem is the output DTO for vendor order CSV export.
type VendorOrderExportItem struct {
	OrderNo      string    `json:"order_no"`
	OrderDate    time.Time `json:"order_date"`
	CustomerName string    `json:"customer_name"`
	ProductName  string    `json:"product_name"`
	Variant      string    `json:"variant"`
	Qty          int       `json:"qty"`
	UnitPrice    float64   `json:"unit_price"`
	LineTotal    float64   `json:"line_total"`
	Subtotal     float64   `json:"subtotal"`
	ShippingFee  float64   `json:"shipping_fee"`
	PlatformFee  float64   `json:"platform_fee"`
	GrandTotal   float64   `json:"grand_total"`
	Status       string    `json:"status"`
}

// VendorOrderListItem is the output DTO for vendor order listing.
type VendorOrderListItem struct {
	OrderID      uuid.UUID                `json:"order_id"`
	OrderNo      string                   `json:"order_no"`
	OrderDate    time.Time                `json:"order_date"`
	CustomerName string                   `json:"customer_name"`
	Items        []VendorOrderProductItem `json:"items"`
	TotalPayment float64                  `json:"total_payment"`
	Status       string                   `json:"status"`
}

// VendorOrderProductItem is product info in a vendor order row (used in list).
type VendorOrderProductItem struct {
	ProductName     string  `json:"product_name"`
	SelectedVariant string  `json:"selected_variant"`
	Qty             int     `json:"qty"`
	UnitPrice       float64 `json:"unit_price"`
	LineTotal       float64 `json:"line_total"`
}

// --- Order Detail Response (redesigned for frontend) ---

// VendorOrderDetailResponse is the full detail DTO for a vendor order.
type VendorOrderDetailResponse struct {
	Shipping  VendorOrderDetailShipping  `json:"shipping"`
	Recipient VendorOrderDetailRecipient `json:"recipient"`
	Store     VendorOrderDetailStore     `json:"store"`
	Items     []VendorOrderDetailItem    `json:"items"`
	Payment   VendorOrderDetailPayment   `json:"payment"`
	Order     VendorOrderDetailOrder     `json:"order"`
	Summary   VendorOrderDetailSummary   `json:"summary"`
}

// VendorOrderDetailShipping is shipping/tracking info.
type VendorOrderDetailShipping struct {
	Courier        string     `json:"courier"`
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}

// VendorOrderDetailRecipientAddress is the nested address structure.
type VendorOrderDetailRecipientAddress struct {
	Street     string `json:"street"`
	District   string `json:"district"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}

// VendorOrderDetailRecipient is recipient info.
type VendorOrderDetailRecipient struct {
	Name    string                            `json:"name"`
	Phone   string                            `json:"phone"`
	Address VendorOrderDetailRecipientAddress `json:"address"`
}

// VendorOrderDetailStore is store/vendor info.
type VendorOrderDetailStore struct {
	Name string `json:"name"`
}

// VendorOrderDetailItem is product item info.
type VendorOrderDetailItem struct {
	Name     string  `json:"name"`
	Variant  string  `json:"variant"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

// VendorOrderDetailPayment is payment info.
type VendorOrderDetailPayment struct {
	Method      string  `json:"method"`
	TotalAmount float64 `json:"total_amount"`
}

// VendorOrderDetailOrder is order metadata/timestamps.
type VendorOrderDetailOrder struct {
	OrderNumber   string     `json:"order_number"`
	OrderTime     time.Time  `json:"order_time"`
	PaymentTime   *time.Time `json:"payment_time,omitempty"`
	ShippingTime  *time.Time `json:"shipping_time,omitempty"`
	DeliveredTime *time.Time `json:"delivered_time,omitempty"`
}

// VendorOrderDetailSummary is the order summary with totals.
type VendorOrderDetailSummary struct {
	ItemsTotal  float64 `json:"items_total"`
	ShippingFee float64 `json:"shipping_fee"`
	ServiceFee  float64 `json:"service_fee"`
	GrandTotal  float64 `json:"grand_total"`
}

// ShipOrderRequest is the input DTO for vendor shipping an order (input resi).
// Courier and service type are already set during checkout; only tracking number is needed.
type ShipOrderRequest struct {
	TrackingNo string `json:"tracking_no" binding:"required"`
}

// AcceptRejectOrderResponse is the output DTO for accept/reject actions.
type AcceptRejectOrderResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderStatus string    `json:"order_status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShipOrderResponse is the output DTO for ship action (input resi).
type ShipOrderResponse struct {
	OrderID          uuid.UUID `json:"order_id"`
	OrderStatus      string    `json:"order_status"`
	TrackingNo       string    `json:"tracking_no"`
	CourierCode      string    `json:"courier_code"`
	ShippedAt        time.Time `json:"shipped_at"`
	ETD              string    `json:"etd"`
	EstimatedArrival *string   `json:"estimated_arrival,omitempty"`
}

// --- Tracking AWB DTOs ---

// TrackWaybillResponse is the output DTO for waybill tracking.
type TrackWaybillResponse struct {
	Delivered      bool                `json:"delivered"`
	Summary        TrackWaybillSummary `json:"summary"`
	Details        TrackWaybillDetails `json:"details"`
	DeliveryStatus TrackDeliveryStatus `json:"delivery_status"`
	Manifest       []TrackManifestItem `json:"manifest"`
}

// TrackWaybillSummary is a summary of the tracked waybill.
type TrackWaybillSummary struct {
	CourierCode   string `json:"courier_code"`
	CourierName   string `json:"courier_name"`
	WaybillNumber string `json:"waybill_number"`
	ServiceCode   string `json:"service_code"`
	WaybillDate   string `json:"waybill_date"`
	ShipperName   string `json:"shipper_name"`
	ReceiverName  string `json:"receiver_name"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Status        string `json:"status"`
}

// TrackWaybillDetails holds detailed waybill info.
type TrackWaybillDetails struct {
	WaybillNumber    string `json:"waybill_number"`
	WaybillDate      string `json:"waybill_date"`
	WaybillTime      string `json:"waybill_time"`
	Weight           string `json:"weight"`
	Origin           string `json:"origin"`
	Destination      string `json:"destination"`
	ShipperName      string `json:"shipper_name"`
	ShipperAddress1  string `json:"shipper_address1"`
	ShipperAddress2  string `json:"shipper_address2"`
	ShipperAddress3  string `json:"shipper_address3"`
	ShipperCity      string `json:"shipper_city"`
	ReceiverName     string `json:"receiver_name"`
	ReceiverAddress1 string `json:"receiver_address1"`
	ReceiverAddress2 string `json:"receiver_address2"`
	ReceiverAddress3 string `json:"receiver_address3"`
	ReceiverCity     string `json:"receiver_city"`
}

// TrackDeliveryStatus holds final delivery status.
type TrackDeliveryStatus struct {
	Status      string `json:"status"`
	PodReceiver string `json:"pod_receiver"`
	PodDate     string `json:"pod_date"`
	PodTime     string `json:"pod_time"`
}

// TrackManifestItem is one event in the tracking manifest.
type TrackManifestItem struct {
	ManifestCode        string `json:"manifest_code"`
	ManifestDescription string `json:"manifest_description"`
	ManifestDate        string `json:"manifest_date"`
	ManifestTime        string `json:"manifest_time"`
	CityName            string `json:"city_name"`
}

// --- Repository interfaces ---

// ShipmentRepository handles shipment data access.
type ShipmentRepository interface {
	// Create inserts a new shipment record.
	Create(ctx context.Context, shipment *Shipment) error
	// FindByOrderID returns the shipment for the given order, or nil.
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*Shipment, error)
	// UpdateTrackingAndShip sets tracking number, status to "shipped", and shipped_at on the existing shipment.
	UpdateTrackingAndShip(ctx context.Context, orderID uuid.UUID, trackingNo string, shippedAt time.Time) error
}

// OrderRepository extensions for vendor order management.
// These are added to the existing OrderRepository interface via VendorOrderRepository.

// VendorOrderRepository extends order data access for vendor operations.
type VendorOrderRepository interface {
	// ListByVendor returns vendor orders with optional status filter and pagination.
	ListByVendor(ctx context.Context, vendorID uuid.UUID, params VendorOrderListParams) ([]Order, int64, error)
	// FindByIDAndVendor returns an order only if it belongs to the given vendor.
	FindByIDAndVendor(ctx context.Context, orderID, vendorID uuid.UUID) (*Order, error)
}

// --- Invoice DTOs ---

// InvoiceActions indicates what actions are available for the invoice.
type InvoiceActions struct {
	CanCopy     bool `json:"can_copy"`
	CanPrint    bool `json:"can_print"`
	CanDownload bool `json:"can_download"`
}

// InvoiceInfo holds invoice metadata.
type InvoiceInfo struct {
	InvoiceNumber string         `json:"invoice_number"`
	IssuedAt      string         `json:"issued_at"`
	Actions       InvoiceActions `json:"actions"`
}

// InvoiceCustomerAddress holds customer address info for invoice.
type InvoiceCustomerAddress struct {
	Street     string `json:"street"`
	District   string `json:"district"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}

// InvoiceCustomer holds customer info for invoice.
type InvoiceCustomer struct {
	Name    string                 `json:"name"`
	Phone   string                 `json:"phone"`
	Address InvoiceCustomerAddress `json:"address"`
}

// InvoicePayment holds payment info for invoice.
type InvoicePayment struct {
	Method string `json:"method"`
	PaidAt string `json:"paid_at"`
}

// InvoiceItem holds a single item in the invoice.
type InvoiceItem struct {
	Name     string  `json:"name"`
	Variant  string  `json:"variant"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
	ImageURL *string `json:"image_url"`
}

// InvoiceSummary holds invoice totals.
type InvoiceSummary struct {
	Subtotal    float64 `json:"subtotal"`
	ShippingFee float64 `json:"shipping_fee"`
	ServiceFee  float64 `json:"service_fee"`
	Total       float64 `json:"total"`
	Currency    string  `json:"currency"`
}

// VendorOrderInvoiceResponse is the output DTO for vendor order invoice.
type VendorOrderInvoiceResponse struct {
	Invoice  InvoiceInfo     `json:"invoice"`
	Customer InvoiceCustomer `json:"customer"`
	Payment  InvoicePayment  `json:"payment"`
	Items    []InvoiceItem   `json:"items"`
	Summary  InvoiceSummary  `json:"summary"`
}

// --- Shipping Info DTOs ---

// ShippingInfoStage represents a stage in the shipping timeline.
type ShippingInfoStage struct {
	Status    string `json:"status"`
	Label     string `json:"label"`
	Completed bool   `json:"completed"`
	Active    bool   `json:"active"`
}

// ShippingInfoCourier holds courier info for shipping.
type ShippingInfoCourier struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	TrackingNo string `json:"tracking_no"`
}

// ShippingInfoEvent represents a single tracking event.
type ShippingInfoEvent struct {
	DateTime    string `json:"date_time"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// VendorOrderShippingInfoResponse is the output DTO for shipping info.
type VendorOrderShippingInfoResponse struct {
	EstimatedArrival *string             `json:"estimated_arrival,omitempty"`
	Stages           []ShippingInfoStage `json:"stages"`
	Courier          ShippingInfoCourier `json:"courier"`
	Events           []ShippingInfoEvent `json:"events"`
}

// --- Usecase interface ---

// VendorOrderUseCase defines the interface for vendor order management.
type VendorOrderUseCase interface {
	// ListOrders returns paginated orders for a vendor.
	ListOrders(ctx context.Context, vendorID uuid.UUID, params VendorOrderListParams) ([]VendorOrderListItem, *PaginationMeta, error)
	// ExportOrders returns all orders for CSV export (no pagination).
	ExportOrders(ctx context.Context, vendorID uuid.UUID, params VendorOrderListParams) ([]VendorOrderExportItem, error)
	// GetOrderDetail returns full order detail for a vendor.
	GetOrderDetail(ctx context.Context, vendorID, orderID uuid.UUID) (*VendorOrderDetailResponse, error)
	// GetOrderInvoice returns invoice data for an order.
	GetOrderInvoice(ctx context.Context, vendorID, orderID uuid.UUID) (*VendorOrderInvoiceResponse, error)
	// GetOrderShippingInfo returns shipping/tracking info for an order.
	GetOrderShippingInfo(ctx context.Context, vendorID, orderID uuid.UUID) (*VendorOrderShippingInfoResponse, error)
	// AcceptOrder transitions an order from paid → processing.
	AcceptOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID) (*AcceptRejectOrderResponse, error)
	// RejectOrder transitions an order from paid/processing → canceled.
	RejectOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID) (*AcceptRejectOrderResponse, error)
	// ShipOrder inputs tracking number and transitions from processing → shipped.
	ShipOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID, req ShipOrderRequest) (*ShipOrderResponse, error)
	// TrackWaybill returns tracking info for an order's shipment.
	TrackWaybill(ctx context.Context, vendorID, orderID uuid.UUID) (*TrackWaybillResponse, error)
}
