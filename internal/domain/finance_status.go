package domain

// Payment invoice statuses.
const (
	PaymentInvoiceStatusPending  = "pending"
	PaymentInvoiceStatusPaid     = "paid"
	PaymentInvoiceStatusFailed   = "failed"
	PaymentInvoiceStatusExpired  = "expired"
	PaymentInvoiceStatusRefunded = "refunded"
)

// Refund statuses.
const (
	RefundStatusRequested = "requested"
	RefundStatusApproved  = "approved"
	RefundStatusRejected  = "rejected"
	RefundStatusProcessed = "processed"
)

// Payout batch statuses.
const (
	PayoutStatusReady    = "ready"
	PayoutStatusSchedule = "schedule"
	PayoutStatusComplete = "complete"
	PayoutStatusFailed   = "failed"
	PayoutStatusOnHold   = "on hold"
)
