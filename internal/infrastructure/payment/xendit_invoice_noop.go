package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// noopXenditInvoiceClient is a dev-only bypass that returns dummy invoice responses
// without calling the real Xendit API. Enable via XENDIT_BYPASS=true.
type noopXenditInvoiceClient struct{}

// NewNoopXenditInvoiceClient creates a no-op Xendit invoice provider for development.
func NewNoopXenditInvoiceClient() domain.XenditInvoiceProvider {
	log.Println("WARNING: Xendit invoice bypass is ENABLED — no real invoices will be created")
	return &noopXenditInvoiceClient{}
}

func (c *noopXenditInvoiceClient) CreateInvoice(_ context.Context, forUserID string, req domain.XenditInvoiceRequest) (*domain.XenditInvoiceResponse, error) {
	fakeID := fmt.Sprintf("noop-inv-%s", uuid.New().String()[:8])
	expiresAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	log.Printf("[XENDIT-BYPASS] CreateInvoice: forUserID=%s externalID=%s amount=%.0f → fakeID=%s",
		forUserID, req.ExternalID, req.Amount, fakeID)

	return &domain.XenditInvoiceResponse{
		ID:         fakeID,
		ExternalID: req.ExternalID,
		Status:     "PENDING",
		InvoiceURL: fmt.Sprintf("https://checkout-bypass.example.com/%s", fakeID),
		ExpiryDate: expiresAt,
	}, nil
}

func (c *noopXenditInvoiceClient) ExpireInvoice(_ context.Context, forUserID string, invoiceID string) error {
	log.Printf("[XENDIT-BYPASS] ExpireInvoice: forUserID=%s invoiceID=%s → skipped", forUserID, invoiceID)
	return nil
}
