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
	OrderID     uuid.UUID              `json:"order_id"`
	OrderStatus string                 `json:"order_status"`
	CanceledAt  time.Time              `json:"canceled_at"`
	RefundInfo  *CancelOrderRefundInfo `json:"refund_info,omitempty"`
}

// CancelOrderRefundInfo describes the refund created upon paid-order cancellation.
type CancelOrderRefundInfo struct {
	RefundID     uuid.UUID `json:"refund_id"`
	RefundMethod string    `json:"refund_method"` // "gateway" or "disbursement"
	Status       string    `json:"status"`
	Message      string    `json:"message"`
}

// CreateOrderRefundRequest is the input DTO when customer submits a refund request.
type CreateOrderRefundRequest struct {
	ReturnReasonID               int      `json:"return_reason_id" binding:"required,min=1"`
	Description                  string   `json:"description" binding:"omitempty,max=1000"`
	DestinationChannelCode       string   `json:"destination_channel_code" binding:"omitempty,max=40"`
	DestinationBankName          string   `json:"destination_bank_name" binding:"required,max=80"`
	DestinationAccountNumber     string   `json:"destination_account_number" binding:"required,max=64"`
	DestinationAccountHolderName string   `json:"destination_account_holder_name" binding:"required,max=120"`
	ImageObjectKeys              []string `json:"image_object_keys" binding:"omitempty,max=5,dive,required"`
	VideoObjectKey               *string  `json:"video_object_key" binding:"omitempty"`
}

// CreateOrderRefundResponse is returned after customer submits a refund request.
type CreateOrderRefundResponse struct {
	RefundID                      uuid.UUID `json:"refund_id"`
	OrderID                       uuid.UUID `json:"order_id"`
	Status                        string    `json:"status"`
	Amount                        float64   `json:"amount"`
	ReturnReasonID                int       `json:"return_reason_id"`
	Reason                        string    `json:"reason"`
	Description                   *string   `json:"description,omitempty"`
	DestinationChannelCode        string    `json:"destination_channel_code"`
	DestinationBankName           string    `json:"destination_bank_name"`
	DestinationAccountHolderName  string    `json:"destination_account_holder_name"`
	DestinationAccountNumberLast4 string    `json:"destination_account_number_last4"`
	ImageURLs                     []string  `json:"image_urls,omitempty"`
	VideoURL                      *string   `json:"video_url,omitempty"`
	RequestedAt                   time.Time `json:"requested_at"`
}

// PresignRefundEvidenceRequest is used to generate a direct upload URL for refund evidence.
type PresignRefundEvidenceRequest struct {
	ContentType string `json:"content_type" binding:"required"`
}

// PresignRefundEvidenceResponse is returned after presign generation for refund evidence.
type PresignRefundEvidenceResponse struct {
	UploadURL   string `json:"upload_url"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	ExpiresIn   int    `json:"expires_in"`
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
	OrderID       uuid.UUID                  `json:"order_id"`
	Items         []CustomerOrderProductItem `json:"items"`
	TotalPayment  float64                    `json:"total_payment"`
	OrderDate     time.Time                  `json:"order_date"`
	Status        string                     `json:"status"`
	PaymentStatus string                     `json:"payment_status"`
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
	Refund   *CustomerOrderDetailRefund   `json:"refund,omitempty"`
}

// CustomerOrderDetailRefund is refund info shown in customer order detail.
type CustomerOrderDetailRefund struct {
	RefundID                      uuid.UUID `json:"refund_id"`
	Status                        string    `json:"status"`
	RefundMethod                  string    `json:"refund_method"`
	Amount                        float64   `json:"amount"`
	Message                       string    `json:"message"`
	DestinationBankName           string    `json:"destination_bank_name,omitempty"`
	DestinationAccountNumberLast4 string    `json:"destination_account_number_last4,omitempty"`
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
	VendorID    uuid.UUID `json:"vendor_id"`
	OwnerUserID uuid.UUID `json:"owner_user_id"`
	Name        string    `json:"name"`
}

// CustomerOrderDetailItem is a product item in the order detail.
type CustomerOrderDetailItem struct {
	OrderItemID uuid.UUID `json:"order_item_id"`
	Name        string    `json:"name"`
	Variant     string    `json:"variant"`
	ImageURL    *string   `json:"image_url"`
	Price       float64   `json:"price"`
	Qty         int       `json:"qty"`
	Subtotal    float64   `json:"subtotal"`
}

// CustomerOrderDetailPayment is payment info in order detail.
type CustomerOrderDetailPayment struct {
	Method    string  `json:"method"`
	InvoiceID *string `json:"invoice_id,omitempty"`
}

// CustomerOrderDetailMeta is order metadata/timestamps.
type CustomerOrderDetailMeta struct {
	OrderNo       string     `json:"order_no"`
	Status        string     `json:"status"`
	PaymentStatus string     `json:"payment_status"`
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
	CanCancel              bool `json:"can_cancel"`
	CanComplete            bool `json:"can_complete"`
	CanReview              bool `json:"can_review"`
	CanContact             bool `json:"can_contact"`
	CanRefund              bool `json:"can_refund"`
	NeedsRefundDestination bool `json:"needs_refund_destination"`
}

// CreateCustomerRefundInput is the persistence input for creating a refund request.
type CreateCustomerRefundInput struct {
	OrderID                       uuid.UUID
	PaymentInvoiceID              uuid.UUID
	Amount                        float64
	ReturnReasonID                int
	Reason                        string
	Description                   *string
	RequestedBy                   uuid.UUID
	RequestedAt                   time.Time
	Status                        string // defaults to RefundStatusRequested
	RefundMethod                  string // "gateway" or "disbursement"
	DestinationChannelCode        string
	DestinationBankName           string
	DestinationAccountNumber      string
	DestinationAccountHolderName  string
	DestinationAccountNumberLast4 string
	UserBankAccountID             *uuid.UUID
	Evidences                     []CreateCustomerRefundEvidenceInput
}

// CreateCustomerRefundEvidenceInput is one evidence object attached to a refund request.
type CreateCustomerRefundEvidenceInput struct {
	ObjectKey     string
	MimeType      string
	FileSizeBytes int
	MediaType     string
	SortOrder     int
}

// CustomerRefundRepository handles customer refund request persistence.
type CustomerRefundRepository interface {
	HasOpenRefundByOrderID(ctx context.Context, orderID uuid.UUID) (bool, error)
	CreateRefundRequest(ctx context.Context, input CreateCustomerRefundInput) (*CreateOrderRefundResponse, error)
	SubmitRefundDestination(ctx context.Context, refundID uuid.UUID, channelCode, bankName, accountNumber, accountHolderName, accountLast4 string, userBankAccountID *uuid.UUID) error
	FindByID(ctx context.Context, refundID uuid.UUID) (*RefundDetail, error)
	FindAwaitingDestinationByOrderAndUser(ctx context.Context, orderID, userID uuid.UUID) (*RefundDetail, error)
	FindLatestByOrderID(ctx context.Context, orderID uuid.UUID) (*RefundDetail, error)
}

// RefundDetail is a read model for a specific refund record.
type RefundDetail struct {
	ID                            uuid.UUID
	OrderID                       uuid.UUID
	PaymentInvoiceID              uuid.UUID
	Amount                        float64
	Status                        string
	RefundMethod                  *string
	DestinationChannelCode        *string
	DestinationAccountNumber      *string
	DestinationBankName           *string
	DestinationAccountNumberLast4 *string
}

// OrderActionUseCase defines customer order lifecycle actions.
type OrderActionUseCase interface {
	// CompleteByCustomer marks a received order as completed.
	CompleteByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*ReceiveOrderResponse, error)
	// CancelByCustomer cancels an order that is still cancellable.
	CancelByCustomer(ctx context.Context, userID, orderID uuid.UUID) (*CancelOrderResponse, error)
	// RequestRefundByCustomer submits a new refund request for an eligible order.
	RequestRefundByCustomer(ctx context.Context, userID, orderID uuid.UUID, req CreateOrderRefundRequest) (*CreateOrderRefundResponse, error)
	// SubmitRefundDestination fills in the bank account for a VA refund awaiting destination.
	SubmitRefundDestination(ctx context.Context, userID, orderID uuid.UUID, req SubmitRefundDestinationRequest) (*SubmitRefundDestinationResponse, error)
	// PresignRefundEvidenceByCustomer generates presigned URL for refund evidence upload.
	PresignRefundEvidenceByCustomer(ctx context.Context, userID, orderID uuid.UUID, req PresignRefundEvidenceRequest) (*PresignRefundEvidenceResponse, error)
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

// SubmitRefundDestinationRequest is the input when customer submits bank for VA refund.
type SubmitRefundDestinationRequest struct {
	UserBankAccountID            *uuid.UUID `json:"user_bank_account_id"`
	DestinationChannelCode       string     `json:"destination_channel_code" binding:"omitempty,max=40"`
	DestinationBankName          string     `json:"destination_bank_name" binding:"omitempty,max=80"`
	DestinationAccountNumber     string     `json:"destination_account_number" binding:"omitempty,max=64"`
	DestinationAccountHolderName string     `json:"destination_account_holder_name" binding:"omitempty,max=120"`
	SaveBankAccount              bool       `json:"save_bank_account"`
}

// SubmitRefundDestinationResponse is the output after submitting refund destination.
type SubmitRefundDestinationResponse struct {
	RefundID uuid.UUID `json:"refund_id"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
}
