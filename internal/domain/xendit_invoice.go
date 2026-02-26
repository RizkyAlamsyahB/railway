package domain

import "context"

// XenditInvoiceRequest is the request body for creating a Xendit invoice (Payment Link).
type XenditInvoiceRequest struct {
	ExternalID         string                 `json:"external_id"`
	Amount             float64                `json:"amount"`
	Description        string                 `json:"description"`
	InvoiceDuration    int                    `json:"invoice_duration"`
	Customer           *XenditInvoiceCustomer `json:"customer,omitempty"`
	SuccessRedirectURL string                 `json:"success_redirect_url,omitempty"`
	FailureRedirectURL string                 `json:"failure_redirect_url,omitempty"`
	CallbackURL        string                 `json:"callback_url,omitempty"`
	Currency           string                 `json:"currency"`
	Items              []XenditInvoiceItem    `json:"items,omitempty"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// XenditInvoiceCustomer holds the customer data sent to Xendit.
type XenditInvoiceCustomer struct {
	GivenNames   string `json:"given_names,omitempty"`
	Email        string `json:"email,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
}

// XenditInvoiceItem is a line item displayed on the Xendit checkout page.
type XenditInvoiceItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// XenditInvoiceResponse is the response from the Xendit Create Invoice API.
type XenditInvoiceResponse struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	InvoiceURL string `json:"invoice_url"`
	ExpiryDate string `json:"expiry_date"`
}

// XenditInvoiceProvider defines the interface for Xendit Invoice (Payment Link) operations.
type XenditInvoiceProvider interface {
	// CreateInvoice creates a Xendit invoice on behalf of a sub-account.
	// forUserID is the vendor's Xendit account ID (sent as the for-user-id header).
	CreateInvoice(ctx context.Context, forUserID string, req XenditInvoiceRequest) (*XenditInvoiceResponse, error)
}
