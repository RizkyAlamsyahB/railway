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

// --- Checkout Preview DTOs ---

// CheckoutPreviewRequest is the input DTO for the checkout preview endpoint.
type CheckoutPreviewRequest struct {
	AddressID string `json:"address_id" binding:"required,uuid"`
}

// CheckoutPreviewVendorItem represents one cart item in a vendor group (for preview).
type CheckoutPreviewVendorItem struct {
	ProductVariantID uuid.UUID `json:"product_variant_id"`
	ProductName      string    `json:"product_name"`
	VariantName      string    `json:"variant_name"`
	ImageURL         *string   `json:"image_url"`
	Price            float64   `json:"price"`
	OriginalPrice    float64   `json:"original_price"`
	PromoPrice       *float64  `json:"promo_price,omitempty"`
	HasPromo         bool      `json:"has_promo"`
	Qty              int       `json:"qty"`
	Subtotal         float64   `json:"subtotal"`
	WeightGram       int       `json:"weight_gram"`
}

// CheckoutPreviewVendorGroup represents one vendor's products and shipping options.
type CheckoutPreviewVendorGroup struct {
	VendorID        uuid.UUID                   `json:"vendor_id"`
	VendorName      string                      `json:"vendor_name"`
	Items           []CheckoutPreviewVendorItem `json:"items"`
	Subtotal        float64                     `json:"subtotal"`
	TotalWeightGram int                         `json:"total_weight_gram"`
	ShippingOptions []ShippingCostOption        `json:"shipping_options"`
}

// CheckoutPreviewResponse is the output DTO for the checkout preview endpoint.
type CheckoutPreviewResponse struct {
	Address     AddressResponse              `json:"address"`
	Vendors     []CheckoutPreviewVendorGroup `json:"vendors"`
	PlatformFee float64                      `json:"platform_fee"`
}

// --- Checkout (Place Order) DTOs ---

// CheckoutRequest is the input DTO for placing an order after preview.
type CheckoutRequest struct {
	AddressID       string                   `json:"address_id" binding:"required,uuid"`
	ShippingChoices []CheckoutShippingChoice `json:"shipping_choices" binding:"required,min=1,dive"`
	Notes           *string                  `json:"notes,omitempty"`
}

// CheckoutShippingChoice represents the customer's courier selection for one vendor.
// Cost is NOT accepted from the client; it is recalculated server-side via RajaOngkir.
type CheckoutShippingChoice struct {
	VendorID    string `json:"vendor_id" binding:"required,uuid"`
	CourierCode string `json:"courier_code" binding:"required"`
	Service     string `json:"service" binding:"required"`
}

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

	// FindItemsByOrderIDs returns all items grouped by order IDs.
	FindItemsByOrderIDs(ctx context.Context, orderIDs []uuid.UUID) (map[uuid.UUID][]OrderItem, error)

	// UpdateOrderStatus sets order_status and payment_status, appending a status history row.
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, orderStatus, paymentStatus string, changedBy *uuid.UUID, notes *string) error

	// MarkOrderReceived transitions a shipped order to received and credits vendor escrow balance.
	// Returns true when the transition is applied, false when no state change is made.
	MarkOrderReceived(ctx context.Context, orderID uuid.UUID, changedBy *uuid.UUID, notes *string) (bool, error)

	// MarkOrderCompleted transitions a received order to completed and releases escrow to available balance.
	// Returns true when the transition is applied, false when no state change is made.
	MarkOrderCompleted(ctx context.Context, orderID uuid.UUID, changedBy *uuid.UUID, notes *string) (bool, error)

	// ListAutoReceiveCandidates returns shipped orders with delivered shipments older than the cutoff.
	ListAutoReceiveCandidates(ctx context.Context, deliveredBefore time.Time, limit int) ([]Order, error)

	// ListSettlementCandidates returns received orders older than the cutoff for settlement.
	ListSettlementCandidates(ctx context.Context, receivedBefore time.Time, limit int) ([]Order, error)

	// RestoreStock increments stock_on_hand for each order item's variant.
	RestoreStock(ctx context.Context, orderID uuid.UUID) error

	// ApplyExpiredWebhookUpdate atomically applies webhook side effects for an expired invoice.
	// It updates invoice status (pending -> expired), updates order status, writes status history,
	// and restores stock in a single DB transaction. Returns false when no state change is applied.
	ApplyExpiredWebhookUpdate(ctx context.Context, orderID, invoiceID uuid.UUID, notes *string) (bool, error)

	// FindByOrderNo returns the order matching the given order number, or nil.
	FindByOrderNo(ctx context.Context, orderNo string) (*Order, error)

	// ListByUser returns customer orders with optional status filter and pagination.
	ListByUser(ctx context.Context, userID uuid.UUID, params CustomerOrderListParams) ([]Order, int64, error)
}

// PaymentRepository defines the interface for payment invoice and event data access.
type PaymentRepository interface {
	// CreateInvoice inserts a new payment_invoices row.
	CreateInvoice(ctx context.Context, invoice *PaymentInvoice) error

	// FindInvoiceByExternalID returns the invoice matching the given external_invoice_id, or nil.
	FindInvoiceByExternalID(ctx context.Context, externalID string) (*PaymentInvoice, error)

	// FindInvoiceByOrderID returns the latest invoice for the given order ID, or nil.
	FindInvoiceByOrderID(ctx context.Context, orderID uuid.UUID) (*PaymentInvoice, error)

	// UpdateInvoiceStatus updates a payment invoice's status and related payment fields.
	UpdateInvoiceStatus(ctx context.Context, invoiceID uuid.UUID, status string, paidAt *time.Time, paymentMethod, paymentChannel *string, rawPayload map[string]interface{}) error

	// UpdateInvoiceStatusIfCurrent updates invoice status only when current status matches expectedCurrentStatus.
	// Returns true when update is applied, false when precondition is not met.
	UpdateInvoiceStatusIfCurrent(
		ctx context.Context,
		invoiceID uuid.UUID,
		expectedCurrentStatus, nextStatus string,
		paidAt *time.Time,
		paymentMethod, paymentChannel *string,
		rawPayload map[string]interface{},
	) (bool, error)

	// CreateEvent inserts a payment_events row.
	CreateEvent(ctx context.Context, event *PaymentEvent) error

	// EventExistsByExternalID checks if a payment event with the given external_event_id already exists.
	EventExistsByExternalID(ctx context.Context, externalEventID string) (bool, error)

	// ListForAdmin returns payment invoices joined with orders, users, and vendors
	// matching the given params, plus total count for pagination.
	ListForAdmin(ctx context.Context, params AdminPaymentListParams) ([]AdminPaymentListItem, int64, error)

	// SummaryForAdmin returns aggregated payment amounts grouped by status.
	SummaryForAdmin(ctx context.Context, params AdminPaymentListParams) (*AdminPaymentSummary, error)

	// ExportForAdmin returns all payment invoices matching the params without pagination.
	ExportForAdmin(ctx context.Context, params AdminPaymentListParams) ([]AdminPaymentListItem, error)
}

// LedgerRepository defines the interface for ledger journal data access.
type LedgerRepository interface {
	// FindAccountByCode returns the ledger account matching the given code, or nil.
	FindAccountByCode(ctx context.Context, code string) (*LedgerAccount, error)

	// CreateJournalWithLines inserts a ledger journal and its lines in a single transaction.
	CreateJournalWithLines(ctx context.Context, journal *LedgerJournal, lines []LedgerLine) error
}

// --- Admin Payment DTOs ---

// AdminPaymentListParams holds query parameters for admin payment invoice listing.
type AdminPaymentListParams struct {
	Page      int
	Limit     int
	Status    string
	Search    string
	SortBy    string
	SortOrder string
	StartDate *time.Time
	EndDate   *time.Time
}

// AdminPaymentListItem is the output DTO for a payment invoice in the admin list.
type AdminPaymentListItem struct {
	InvoiceID     uuid.UUID `json:"invoice_id"`
	OrderID       uuid.UUID `json:"order_id"`
	OrderNo       string    `json:"order_no"`
	CustomerName  string    `json:"customer_name"`
	VendorName    string    `json:"vendor_name"`
	PaymentMethod *string   `json:"payment_method"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// AdminPaymentSummary holds aggregated payment amounts by status.
type AdminPaymentSummary struct {
	TotalAmount   float64 `json:"total_amount"`
	PaidAmount    float64 `json:"paid_amount"`
	PendingAmount float64 `json:"pending_amount"`
	FailedAmount  float64 `json:"failed_amount"`
}

// AdminPaymentUseCase defines the interface for admin payment listing operations.
type AdminPaymentUseCase interface {
	// List returns a paginated list of payment invoices with related order/customer/vendor info.
	List(ctx context.Context, params AdminPaymentListParams) ([]AdminPaymentListItem, *PaginationMeta, error)
	// Summary returns aggregated payment amounts grouped by status.
	Summary(ctx context.Context, params AdminPaymentListParams) (*AdminPaymentSummary, error)
	// Export returns all payment invoices matching the params (without pagination) for CSV export.
	Export(ctx context.Context, params AdminPaymentListParams) ([]AdminPaymentListItem, error)
}

// --- Usecase Interface ---

// CheckoutUseCase defines the interface for checkout business operations.
type CheckoutUseCase interface {
	// Preview returns cart items grouped by vendor with shipping options from RajaOngkir.
	Preview(ctx context.Context, userID uuid.UUID, req CheckoutPreviewRequest) (*CheckoutPreviewResponse, error)

	// Checkout validates the cart, creates orders per vendor, creates Xendit invoices,
	// and returns invoice URLs. Accepts shipping choices from the preview step.
	Checkout(ctx context.Context, userID uuid.UUID, req CheckoutRequest) (*CheckoutResponse, error)

	// HandleWebhook processes an incoming Xendit webhook callback.
	HandleWebhook(ctx context.Context, payload XenditWebhookPayload) error

	// SimulatePayment marks an order as paid without calling Xendit.
	// Used for local development when XENDIT_BYPASS=true.
	SimulatePayment(ctx context.Context, orderID uuid.UUID) error
}
