package domain

import (
	"context"
	"errors"
)

// Payout channel code constants.
const (
	PayoutChannelIDBCA = "ID_BCA"
)

// ValidPayoutChannelCodes is the set of supported payout channel codes.
var ValidPayoutChannelCodes = map[string]bool{
	PayoutChannelIDBCA: true,
}

// XenditPayoutChannelProperties holds the destination account details.
type XenditPayoutChannelProperties struct {
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
}

// XenditPayoutReceiptNotification holds email recipients for payout receipts.
type XenditPayoutReceiptNotification struct {
	EmailTo []string `json:"email_to,omitempty"`
}

// XenditPayoutRequest is the request body for creating a Xendit payout.
type XenditPayoutRequest struct {
	ReferenceID         string                           `json:"reference_id"`
	ChannelCode         string                           `json:"channel_code"`
	ChannelProperties   XenditPayoutChannelProperties    `json:"channel_properties"`
	Amount              float64                          `json:"amount"`
	Description         string                           `json:"description"`
	Currency            string                           `json:"currency"`
	ReceiptNotification *XenditPayoutReceiptNotification `json:"receipt_notification,omitempty"`
}

// XenditPayoutResponse is the response from the Xendit Create Payout API.
type XenditPayoutResponse struct {
	ID                   string                         `json:"id"`
	Amount               float64                        `json:"amount"`
	ChannelCode          string                         `json:"channel_code"`
	Currency             string                         `json:"currency"`
	ReferenceID          string                         `json:"reference_id"`
	Status               string                         `json:"status"`
	Created              string                         `json:"created"`
	Updated              string                         `json:"updated"`
	EstimatedArrivalTime string                         `json:"estimated_arrival_time"`
	BusinessID           string                         `json:"business_id"`
	FailureCode          string                         `json:"failure_code,omitempty"`
	ChannelProperties    *XenditPayoutChannelProperties `json:"channel_properties,omitempty"`
}

// XenditPayoutProvider defines the interface for Xendit Payout operations.
type XenditPayoutProvider interface {
	// CreatePayout creates a payout on behalf of a sub-account.
	// forUserID is the vendor's Xendit account ID (sent as the for-user-id header).
	// idempotencyKey is used to prevent duplicate payouts (sent as the Idempotency-key header).
	CreatePayout(ctx context.Context, forUserID string, idempotencyKey string, req XenditPayoutRequest) (*XenditPayoutResponse, error)

	// GetTransactionByReference fetches the latest payout transaction by reference ID.
	// forUserID is the vendor's Xendit account ID (sent as the for-user-id header).
	GetTransactionByReference(ctx context.Context, forUserID string, referenceID string) (*XenditTransaction, error)
}

var (
	ErrXenditTransactionNotFound       = errors.New("xendit transaction not found")
	ErrXenditTransactionFeeUnavailable = errors.New("xendit transaction fee is unavailable")
)

// XenditTransaction represents a payout transaction entry from the Transactions API.
type XenditTransaction struct {
	ID          string  `json:"id"`
	ReferenceID string  `json:"reference_id"`
	Amount      float64 `json:"amount"`
	Fee         float64 `json:"fee"`
	Status      string  `json:"status"`
	Currency    string  `json:"currency"`
	ChannelCode string  `json:"channel_code"`
}

// Xendit payout webhook event constants.
const (
	XenditPayoutWebhookEventSucceeded = "payout.succeeded"
	XenditPayoutWebhookEventFailed    = "payout.failed"
	XenditPayoutWebhookEventReversed  = "payout.reversed"
)

// Xendit payout status constants seen in webhook payload.
const (
	XenditPayoutStatusSucceeded = "SUCCEEDED"
	XenditPayoutStatusFailed    = "FAILED"
	XenditPayoutStatusCancelled = "CANCELLED"
	XenditPayoutStatusReversed  = "REVERSED"
)

// XenditPayoutWebhookData represents the "data" object from payout webhook payload.
type XenditPayoutWebhookData struct {
	ID            string  `json:"id"`
	ReferenceID   string  `json:"reference_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	FailureCode   string  `json:"failure_code"`
	ChannelCode   string  `json:"channel_code"`
	Created       string  `json:"created"`
	Updated       string  `json:"updated"`
	CompletedDate string  `json:"completed_date"`
}

// XenditPayoutWebhookPayload represents the payout webhook body sent by Xendit.
type XenditPayoutWebhookPayload struct {
	BusinessID string                  `json:"business_id"`
	Event      string                  `json:"event"`
	Data       XenditPayoutWebhookData `json:"data"`
}
