package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Entities ---

// Order represents the orders table.
type Order struct {
	ID                      uuid.UUID              `json:"id"`
	OrderNo                 string                 `json:"order_no"`
	UserID                  uuid.UUID              `json:"user_id"`
	VendorID                uuid.UUID              `json:"vendor_id"`
	ShippingAddressSnapshot map[string]interface{} `json:"shipping_address_snapshot"`
	OrderStatus             string                 `json:"order_status"`
	PaymentStatus           string                 `json:"payment_status"`
	Subtotal                float64                `json:"subtotal"`
	ShippingFee             float64                `json:"shipping_fee"`
	PlatformFee             float64                `json:"platform_fee"`
	GrandTotal              float64                `json:"grand_total"`
	PlacedAt                time.Time              `json:"placed_at"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
}

// OrderItem represents the order_items table.
type OrderItem struct {
	ID                  uuid.UUID `json:"id"`
	OrderID             uuid.UUID `json:"order_id"`
	ProductVariantID    uuid.UUID `json:"product_variant_id"`
	ProductNameSnapshot string    `json:"product_name_snapshot"`
	SKUSnapshot         string    `json:"sku_snapshot"`
	Qty                 int       `json:"qty"`
	UnitPrice           float64   `json:"unit_price"`
	LineTotal           float64   `json:"line_total"`
}

// OrderStatusHistory represents the order_status_history table.
type OrderStatusHistory struct {
	ID        uuid.UUID  `json:"id"`
	OrderID   uuid.UUID  `json:"order_id"`
	OldStatus *string    `json:"old_status"`
	NewStatus string     `json:"new_status"`
	ChangedBy *uuid.UUID `json:"changed_by"`
	ChangedAt time.Time  `json:"changed_at"`
	Notes     *string    `json:"notes"`
}

// PaymentInvoice represents the payment_invoices table.
type PaymentInvoice struct {
	ID                uuid.UUID              `json:"id"`
	OrderID           uuid.UUID              `json:"order_id"`
	Gateway           string                 `json:"gateway"`
	XenditInvoiceID   *string                `json:"xendit_invoice_id"`
	ExternalInvoiceID string                 `json:"external_invoice_id"`
	InvoiceURL        *string                `json:"invoice_url"`
	PaymentMethod     *string                `json:"payment_method"`
	PaymentChannel    *string                `json:"payment_channel"`
	Amount            float64                `json:"amount"`
	Currency          string                 `json:"currency"`
	Status            string                 `json:"status"`
	ExpiresAt         *time.Time             `json:"expires_at"`
	PaidAt            *time.Time             `json:"paid_at"`
	RawPayload        map[string]interface{} `json:"raw_payload"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// PaymentEvent represents the payment_events table.
type PaymentEvent struct {
	ID               uuid.UUID              `json:"id"`
	PaymentInvoiceID uuid.UUID              `json:"payment_invoice_id"`
	EventType        string                 `json:"event_type"`
	ExternalEventID  string                 `json:"external_event_id"`
	Payload          map[string]interface{} `json:"payload"`
	ReceivedAt       time.Time              `json:"received_at"`
}

// LedgerAccount represents the ledger_accounts table.
type LedgerAccount struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	AccountType string    `json:"account_type"`
	NormalSide  string    `json:"normal_side"`
	IsActive    bool      `json:"is_active"`
}

// LedgerJournal represents the ledger_journals table.
type LedgerJournal struct {
	ID          uuid.UUID  `json:"id"`
	JournalNo   string     `json:"journal_no"`
	SourceType  string     `json:"source_type"`
	SourceID    uuid.UUID  `json:"source_id"`
	EventTime   time.Time  `json:"event_time"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	CreatedBy   *uuid.UUID `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

// LedgerLine represents the ledger_lines table.
type LedgerLine struct {
	ID        uuid.UUID `json:"id"`
	JournalID uuid.UUID `json:"journal_id"`
	AccountID uuid.UUID `json:"account_id"`
	Debit     float64   `json:"debit"`
	Credit    float64   `json:"credit"`
	Currency  string    `json:"currency"`
	Reference *string   `json:"reference"`
}

// --- Request/Response DTOs ---

// CheckoutOrderResult represents a single order created during checkout.
type CheckoutOrderResult struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderNo     string    `json:"order_no"`
	VendorID    uuid.UUID `json:"vendor_id"`
	VendorName  string    `json:"vendor_name"`
	Subtotal    float64   `json:"subtotal"`
	ShippingFee float64   `json:"shipping_fee"`
	PlatformFee float64   `json:"platform_fee"`
	GrandTotal  float64   `json:"grand_total"`
	InvoiceURL  string    `json:"invoice_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// CheckoutResponse is the output DTO for the checkout endpoint.
type CheckoutResponse struct {
	Orders []CheckoutOrderResult `json:"orders"`
}

// XenditWebhookPayload represents the incoming webhook body from Xendit.
type XenditWebhookPayload struct {
	ID             string                 `json:"id"`
	ExternalID     string                 `json:"external_id"`
	UserID         string                 `json:"user_id"`
	Status         string                 `json:"status"`
	Amount         float64                `json:"amount"`
	PaidAmount     float64                `json:"paid_amount"`
	PaidAt         *string                `json:"paid_at"`
	PaymentMethod  string                 `json:"payment_method"`
	PaymentChannel string                 `json:"payment_channel"`
	Currency       string                 `json:"currency"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// --- Repository Interfaces ---

// OrderRepository defines the interface for order data access.
type OrderRepository interface {
	// CreateOrderWithItems creates an order, its items, and initial status history
	// in one transaction. It also decrements stock_on_hand for each variant.
	CreateOrderWithItems(ctx context.Context, order *Order, items []OrderItem) error

	// FindByID returns an order by its primary key, or nil if not found.
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)

	// FindItemsByOrderID returns all items for an order.
	FindItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderItem, error)

	// UpdateOrderStatus sets order_status and payment_status, appending a status history row.
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, orderStatus, paymentStatus string, changedBy *uuid.UUID, notes *string) error

	// RestoreStock increments stock_on_hand for each order item's variant.
	RestoreStock(ctx context.Context, orderID uuid.UUID) error
}

// PaymentRepository defines the interface for payment invoice and event data access.
type PaymentRepository interface {
	// CreateInvoice inserts a new payment_invoices row.
	CreateInvoice(ctx context.Context, invoice *PaymentInvoice) error

	// FindInvoiceByExternalID returns the invoice matching the given external_invoice_id, or nil.
	FindInvoiceByExternalID(ctx context.Context, externalID string) (*PaymentInvoice, error)

	// UpdateInvoiceStatus updates a payment invoice's status and related payment fields.
	UpdateInvoiceStatus(ctx context.Context, invoiceID uuid.UUID, status string, paidAt *time.Time, paymentMethod, paymentChannel *string, rawPayload map[string]interface{}) error

	// CreateEvent inserts a payment_events row.
	CreateEvent(ctx context.Context, event *PaymentEvent) error

	// EventExistsByExternalID checks if a payment event with the given external_event_id already exists.
	EventExistsByExternalID(ctx context.Context, externalEventID string) (bool, error)
}

// LedgerRepository defines the interface for ledger journal data access.
type LedgerRepository interface {
	// FindAccountByCode returns the ledger account matching the given code, or nil.
	FindAccountByCode(ctx context.Context, code string) (*LedgerAccount, error)

	// CreateJournalWithLines inserts a ledger journal and its lines in a single transaction.
	CreateJournalWithLines(ctx context.Context, journal *LedgerJournal, lines []LedgerLine) error
}

// --- Usecase Interface ---

// CheckoutUseCase defines the interface for checkout business operations.
type CheckoutUseCase interface {
	// Checkout validates the cart, creates orders per vendor, creates Xendit invoices,
	// and returns invoice URLs.
	Checkout(ctx context.Context, userID uuid.UUID) (*CheckoutResponse, error)

	// HandleWebhook processes an incoming Xendit webhook callback.
	HandleWebhook(ctx context.Context, payload XenditWebhookPayload) error
}
