package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ReceiveOrderResponse is the output DTO for customer order received confirmation.
type ReceiveOrderResponse struct {
	OrderID       uuid.UUID `json:"order_id"`
	OrderStatus   string    `json:"order_status"`
	PaymentStatus string    `json:"payment_status"`
	CompletedAt   time.Time `json:"completed_at"`
}

// CancelOrderResponse is the output DTO for customer order cancellation.
type CancelOrderResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderStatus string    `json:"order_status"`
	CanceledAt  time.Time `json:"canceled_at"`
}

// CustomerOrderListParams holds query parameters for customer order listing.
type CustomerOrderListParams struct {
	Page   int
	Limit  int
	Status string
}

// CustomerOrderProductItem is product and variant info in a customer order row.
type CustomerOrderProductItem struct {
	ProductName     string  `json:"product_name"`
	SelectedVariant string  `json:"selected_variant"`
	ImageURL        *string `json:"image_url"`
	Qty             int     `json:"qty"`
	UnitPrice       float64 `json:"unit_price"`
}

// CustomerOrderListItem is the output DTO for customer order listing.
type CustomerOrderListItem struct {
	OrderID      uuid.UUID                  `json:"order_id"`
	Items        []CustomerOrderProductItem `json:"items"`
	TotalPayment float64                    `json:"total_payment"`
	OrderDate    time.Time                  `json:"order_date"`
	Status       string                     `json:"status"`
}

// --- Customer Order Detail DTOs ---

// CustomerOrderDetailResponse is the full detail DTO for a customer order.
type CustomerOrderDetailResponse struct {
	Shipping *CustomerOrderDetailShipping `json:"shipping,omitempty"`
	Address  CustomerOrderDetailAddress   `json:"address"`
	Store    CustomerOrderDetailStore     `json:"store"`
	Items    []CustomerOrderDetailItem    `json:"items"`
	Payment  CustomerOrderDetailPayment   `json:"payment"`
	Order    CustomerOrderDetailMeta      `json:"order"`
	Summary  CustomerOrderDetailSummary   `json:"summary"`
	Actions  CustomerOrderActions         `json:"actions"`
}

// CustomerOrderDetailShipping is shipping info on customer order detail.
type CustomerOrderDetailShipping struct {
	Courier        string  `json:"courier"`
	TrackingNumber string  `json:"tracking_number"`
	Status         string  `json:"status"`
	StatusDate     *string `json:"status_date,omitempty"`
}

// CustomerOrderDetailAddress is recipient address in customer order detail.
type CustomerOrderDetailAddress struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Street     string `json:"street"`
	District   string `json:"district"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}

// CustomerOrderDetailStore is vendor/store info.
type CustomerOrderDetailStore struct {
	VendorID uuid.UUID `json:"vendor_id"`
	Name     string    `json:"name"`
}

// CustomerOrderDetailItem is a product item in the order detail.
type CustomerOrderDetailItem struct {
	Name     string  `json:"name"`
	Variant  string  `json:"variant"`
	ImageURL *string `json:"image_url"`
	Price    float64 `json:"price"`
	Qty      int     `json:"qty"`
	Subtotal float64 `json:"subtotal"`
}

// CustomerOrderDetailPayment is payment info in order detail.
type CustomerOrderDetailPayment struct {
	Method    string  `json:"method"`
	InvoiceID *string `json:"invoice_id,omitempty"`
}

// CustomerOrderDetailMeta is order metadata/timestamps.
type CustomerOrderDetailMeta struct {
	OrderNo       string     `json:"order_no"`
	OrderTime     time.Time  `json:"order_time"`
	PaymentTime   *time.Time `json:"payment_time,omitempty"`
	ShippingTime  *time.Time `json:"shipping_time,omitempty"`
	DeliveredTime *time.Time `json:"delivered_time,omitempty"`
}

// CustomerOrderDetailSummary is the order totals summary.
type CustomerOrderDetailSummary struct {
	Subtotal    float64 `json:"subtotal"`
	ShippingFee float64 `json:"shipping_fee"`
	ServiceFee  float64 `json:"service_fee"`
	GrandTotal  float64 `json:"grand_total"`
}

// CustomerOrderActions indicates which actions the customer can perform.
type CustomerOrderActions struct {
	CanCancel   bool `json:"can_cancel"`
	CanComplete bool `json:"can_complete"`
	CanReview   bool `json:"can_review"`
	CanContact  bool `json:"can_contact"`
}

// OrderActionUseCase defines customer order lifecycle actions.
type OrderActionUseCase interface {
	// CompleteByCustomer marks a received order as completed.
	CompleteByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*ReceiveOrderResponse, error)
	// CancelByCustomer cancels an order that is still cancellable.
	CancelByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*CancelOrderResponse, error)
	// GetOrderDetail returns full order detail for a customer.
	GetOrderDetail(ctx context.Context, userID, orderID uuid.UUID) (*CustomerOrderDetailResponse, error)
	// GetOrderShippingInfo returns shipping/tracking info for a customer order.
	GetOrderShippingInfo(ctx context.Context, userID, orderID uuid.UUID) (*VendorOrderShippingInfoResponse, error)
	// GetOrderInvoice returns invoice data for a customer order.
	GetOrderInvoice(ctx context.Context, userID, orderID uuid.UUID) (*VendorOrderInvoiceResponse, error)
	// ListByCustomer returns customer order rows with pagination metadata.
	ListByCustomer(ctx context.Context, userID uuid.UUID, params CustomerOrderListParams) ([]CustomerOrderListItem, *PaginationMeta, error)
	// ListOrderStatuses returns available order statuses.
	ListOrderStatuses(ctx context.Context) []string
}
