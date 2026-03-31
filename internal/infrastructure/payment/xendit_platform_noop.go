package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

// noopXenPlatformClient is a dev-only bypass that returns dummy sub-account responses
// without calling the real Xendit XenPlatform API. Enable via XENDIT_BYPASS=true.
type noopXenPlatformClient struct{}

// NewNoopXenPlatformClient creates a no-op XenPlatformProvider for development.
func NewNoopXenPlatformClient() domain.XenPlatformProvider {
	log.Println("WARNING: Xendit XenPlatform bypass is ENABLED — no real sub-accounts will be created")
	return &noopXenPlatformClient{}
}

func (c *noopXenPlatformClient) CreateAccount(_ context.Context, req domain.XenPlatformCreateAccountRequest) (*domain.XenPlatformAccount, error) {
	fakeID := "noop-xen-" + uuid.New().String()[:8]
	log.Printf("[XENDIT-BYPASS] XenPlatform.CreateAccount: email=%s type=%s → fakeID=%s", req.Email, req.Type, fakeID)
	return &domain.XenPlatformAccount{
		ID:     fakeID,
		Type:   req.Type,
		Email:  req.Email,
		Status: "ACTIVE",
	}, nil
}

func (c *noopXenPlatformClient) ListAccounts(_ context.Context, _ domain.XenPlatformListAccountsParams) (*domain.XenPlatformListAccountsResponse, error) {
	log.Println("[XENDIT-BYPASS] XenPlatform.ListAccounts → empty list")
	return &domain.XenPlatformListAccountsResponse{Data: []domain.XenPlatformAccount{}, HasMore: false}, nil
}

func (c *noopXenPlatformClient) GetAccount(_ context.Context, id string) (*domain.XenPlatformAccount, error) {
	log.Printf("[XENDIT-BYPASS] XenPlatform.GetAccount: id=%s → dummy", id)
	return &domain.XenPlatformAccount{ID: id, Status: "ACTIVE"}, nil
}

func (c *noopXenPlatformClient) UpdateAccount(_ context.Context, id string, _ domain.XenPlatformUpdateAccountRequest) (*domain.XenPlatformAccount, error) {
	log.Printf("[XENDIT-BYPASS] XenPlatform.UpdateAccount: id=%s → skipped", id)
	return &domain.XenPlatformAccount{ID: id, Status: "ACTIVE"}, nil
}
