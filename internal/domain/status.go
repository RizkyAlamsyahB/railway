package domain

// User statuses.
const (
	UserStatusPending     = "pending"
	UserStatusActive      = "active"
	UserStatusBlocked     = "blocked"
	UserStatusDeactivated = "deactivated"
)

// Vendor statuses.
const (
	VendorStatusDraft     = "draft"
	VendorStatusSubmitted = "submitted"
	VendorStatusActive    = "active"
	VendorStatusRejected  = "rejected"
	VendorStatusBlocked   = "blocked"
)

// Product statuses.
const (
	ProductStatusDraft     = "draft"
	ProductStatusPublished = "published"
)

// Product types.
const (
	ProductTypeSingle  = "single"
	ProductTypeVariant = "variant"
	ProductTypePackage = "package"
)

// Cart statuses.
const (
	CartStatusActive    = "active"
	CartStatusConverted = "converted"
)

// Order statuses.
const (
	OrderStatusPendingPayment = "pending_payment"
	OrderStatusPaid           = "paid"
	OrderStatusProcessing     = "processing"
	OrderStatusPacked         = "packed"
	OrderStatusShipped        = "shipped"
	OrderStatusReceived       = "received"
	OrderStatusCompleted      = "completed"
	OrderStatusCanceled       = "canceled"
	OrderStatusRefunded       = "refunded"
)

// OrderStatuses lists all allowed order statuses.
var OrderStatuses = []string{
	OrderStatusPendingPayment,
	OrderStatusPaid,
	OrderStatusProcessing,
	OrderStatusPacked,
	OrderStatusShipped,
	OrderStatusReceived,
	OrderStatusCompleted,
	OrderStatusCanceled,
	OrderStatusRefunded,
}

// Payment statuses (on orders).
const (
	PaymentStatusUnpaid   = "unpaid"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)

// Payment invoice statuses.
const (
	InvoiceStatusPending = "pending"
	InvoiceStatusPaid    = "paid"
	InvoiceStatusExpired = "expired"
	InvoiceStatusFailed  = "failed"
)

// Halal AI review statuses.
const (
	HalalAIStatusPending = "pending"
)

// Role codes.
const (
	RoleAdmin    = "admin"
	RoleUMKM     = "umkm"
	RoleCustomer = "customer"
	RoleCS       = "cs"
	RoleFinance  = "finance"
)

// Ticket statuses.
const (
	TicketStatusOpen       = "open"
	TicketStatusOnProgress = "on_progress"
	TicketStatusClosed     = "closed"
)

// Ticket sources.
const (
	TicketSourceApp = "app"
	TicketSourceWeb = "web"
)

// Review statuses.
const (
	ReviewStatusPublished = "published"
	ReviewStatusHidden    = "hidden"
)
