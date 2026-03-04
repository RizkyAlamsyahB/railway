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
	PayoutStatusSchedule = "schedule"
	PayoutStatusComplete = "completed"
	PayoutStatusFailed   = "failed"
	PayoutStatusOnHold   = "on_hold"
)

// Vendor withdrawal statuses.
const (
	WithdrawalStatusPending    = "pending"
	WithdrawalStatusProcessing = "processing"
	WithdrawalStatusCompleted  = "completed"
	WithdrawalStatusFailed     = "failed"
)
