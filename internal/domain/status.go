package domain

// User statuses.
const (
	UserStatusPending = "pending"
	UserStatusActive  = "active"
	UserStatusBlocked = "blocked"
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

// Cart statuses.
const (
	CartStatusActive = "active"
)

// Halal AI review statuses.
const (
	HalalAIStatusPending = "pending"
)

// Verification statuses (documents, bank accounts).
const (
	VerificationStatusPending = "pending"
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
	TicketStatusResolved   = "resolved"
	TicketStatusClosed     = "closed"
)

// Ticket sources.
const (
	TicketSourceApp = "app"
	TicketSourceWeb = "web"
)

// Chat conversation statuses.
const (
	ChatConvStatusOpen   = "open"
	ChatConvStatusClosed = "closed"
)
