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
	Page   int
	Limit  int
	Status string
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

// VendorOrderProductItem is product info in a vendor order row.
type VendorOrderProductItem struct {
	ProductName     string  `json:"product_name"`
	SelectedVariant string  `json:"selected_variant"`
	Qty             int     `json:"qty"`
	UnitPrice       float64 `json:"unit_price"`
	LineTotal       float64 `json:"line_total"`
}

// VendorOrderDetailResponse is the full detail DTO for a vendor order.
type VendorOrderDetailResponse struct {
	OrderID      uuid.UUID                    `json:"order_id"`
	OrderNo      string                       `json:"order_no"`
	OrderDate    time.Time                    `json:"order_date"`
	Status       string                       `json:"status"`
	CustomerName string                       `json:"customer_name"`
	Items        []VendorOrderProductItem     `json:"items"`
	Subtotal     float64                      `json:"subtotal"`
	ShippingFee  float64                      `json:"shipping_fee"`
	PlatformFee  float64                      `json:"platform_fee"`
	GrandTotal   float64                      `json:"grand_total"`
	Address      map[string]interface{}       `json:"shipping_address"`
	Shipment     *VendorOrderShipmentResponse `json:"shipment"`
	PaymentInfo  *VendorOrderPaymentInfo      `json:"payment_info"`
}

// VendorOrderShipmentResponse is shipping info within a vendor order.
type VendorOrderShipmentResponse struct {
	CourierCode      string     `json:"courier_code"`
	ServiceType      string     `json:"service_type"`
	TrackingNo       string     `json:"tracking_no"`
	ETD              string     `json:"etd"`
	ShipmentStatus   string     `json:"shipment_status"`
	ShippedAt        *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
	EstimatedArrival *string    `json:"estimated_arrival,omitempty"`
}

// VendorOrderPaymentInfo is payment info within a vendor order.
type VendorOrderPaymentInfo struct {
	Status         string     `json:"status"`
	PaymentMethod  *string    `json:"payment_method,omitempty"`
	PaymentChannel *string    `json:"payment_channel,omitempty"`
	PaidAt         *time.Time `json:"paid_at,omitempty"`
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

// --- Usecase interface ---

// VendorOrderUseCase defines the interface for vendor order management.
type VendorOrderUseCase interface {
	// ListOrders returns paginated orders for a vendor.
	ListOrders(ctx context.Context, vendorID uuid.UUID, params VendorOrderListParams) ([]VendorOrderListItem, *PaginationMeta, error)
	// GetOrderDetail returns full order detail for a vendor.
	GetOrderDetail(ctx context.Context, vendorID, orderID uuid.UUID) (*VendorOrderDetailResponse, error)
	// AcceptOrder transitions an order from paid → processing.
	AcceptOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID) (*AcceptRejectOrderResponse, error)
	// RejectOrder transitions an order from paid/processing → canceled.
	RejectOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID) (*AcceptRejectOrderResponse, error)
	// ShipOrder inputs tracking number and transitions from processing → shipped.
	ShipOrder(ctx context.Context, vendorID, orderID, userID uuid.UUID, req ShipOrderRequest) (*ShipOrderResponse, error)
	// TrackWaybill returns tracking info for an order's shipment.
	TrackWaybill(ctx context.Context, vendorID, orderID uuid.UUID) (*TrackWaybillResponse, error)
}
