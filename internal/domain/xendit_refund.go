package domain

import "context"

// Refund method constants.
const (
	RefundMethodGateway      = "gateway"      // QR/E-Wallet — auto refund via Xendit Refund API
	RefundMethodDisbursement = "disbursement" // VA — payout to customer bank account
)

// XenditRefundRequest is the request body for creating a Xendit refund.
type XenditRefundRequest struct {
	PaymentRequestID string  `json:"payment_request_id,omitempty"`
	InvoiceID        string  `json:"invoice_id,omitempty"`
	ReferenceID      string  `json:"reference_id"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Reason           string  `json:"reason"`
}

// XenditRefundResponse is the response from the Xendit Create Refund API.
type XenditRefundResponse struct {
	ID          string  `json:"id"`
	PaymentID   string  `json:"payment_id"`
	InvoiceID   string  `json:"invoice_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	Reason      string  `json:"reason"`
	ReferenceID string  `json:"reference_id"`
	FailureCode string  `json:"failure_code,omitempty"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
}

// Xendit refund status constants.
const (
	XenditRefundStatusSucceeded = "SUCCEEDED"
	XenditRefundStatusFailed    = "FAILED"
	XenditRefundStatusPending   = "PENDING"
)

// Xendit refund webhook event constants.
const (
	XenditRefundWebhookEventSucceeded = "refund.succeeded"
	XenditRefundWebhookEventFailed    = "refund.failed"
)

// XenditRefundWebhookData is the "data" object from the refund webhook payload.
type XenditRefundWebhookData struct {
	ID                string  `json:"id"`
	PaymentID         string  `json:"payment_id"`
	InvoiceID         string  `json:"invoice_id"`
	Amount            float64 `json:"amount"`
	PaymentMethodType string  `json:"payment_method_type"`
	ChannelCode       string  `json:"channel_code"`
	Currency          string  `json:"currency"`
	Status            string  `json:"status"`
	Reason            string  `json:"reason"`
	ReferenceID       string  `json:"reference_id"`
	FailureCode       string  `json:"failure_code"`
	Created           string  `json:"created"`
	Updated           string  `json:"updated"`
}

// XenditRefundWebhookPayload is the refund webhook body sent by Xendit.
type XenditRefundWebhookPayload struct {
	Event      string                  `json:"event"`
	BusinessID string                  `json:"business_id"`
	Created    string                  `json:"created"`
	Data       XenditRefundWebhookData `json:"data"`
}

// XenditRefundProvider defines the interface for Xendit Refund operations.
type XenditRefundProvider interface {
	// CreateRefund creates a refund for a QR/E-Wallet payment.
	// forUserID is the vendor's Xendit account ID (sent as the for-user-id header).
	CreateRefund(ctx context.Context, forUserID string, idempotencyKey string, req XenditRefundRequest) (*XenditRefundResponse, error)
}
