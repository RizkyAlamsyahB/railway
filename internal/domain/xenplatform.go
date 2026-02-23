package domain

import "context"

// XenPlatformPublicProfile holds the public-facing profile of a Xendit sub-account.
type XenPlatformPublicProfile struct {
	BusinessName string `json:"business_name"`
}

// XenPlatformCreateAccountRequest is the request body for creating a Xendit sub-account.
type XenPlatformCreateAccountRequest struct {
	Email         string                    `json:"email"`
	Type          string                    `json:"type"`
	PublicProfile *XenPlatformPublicProfile `json:"public_profile,omitempty"`
}

// XenPlatformAccount represents a Xendit sub-account returned from the API.
type XenPlatformAccount struct {
	ID              string                    `json:"id"`
	Type            string                    `json:"type"`
	Email           string                    `json:"email"`
	PublicProfile   *XenPlatformPublicProfile `json:"public_profile,omitempty"`
	Status          string                    `json:"status"`
	AccountHolderID *string                   `json:"account_holder_id,omitempty"`
}

// XenPlatformListAccountsParams holds query parameters for listing Xendit sub-accounts.
type XenPlatformListAccountsParams struct {
	Email        string
	Status       string
	BusinessName string
	Type         string
	CreatedGTE   string
	CreatedLTE   string
	Limit        int
	AfterID      string
}

// XenPlatformListLink represents a pagination link in the list response.
type XenPlatformListLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// XenPlatformListAccountsResponse is the response envelope for listing Xendit sub-accounts.
type XenPlatformListAccountsResponse struct {
	Data    []XenPlatformAccount  `json:"data"`
	HasMore bool                  `json:"has_more"`
	Links   []XenPlatformListLink `json:"links"`
}

// XenPlatformUpdateAccountRequest is the request body for updating a Xendit sub-account.
type XenPlatformUpdateAccountRequest struct {
	Email           *string                   `json:"email,omitempty"`
	PublicProfile   *XenPlatformPublicProfile `json:"public_profile,omitempty"`
	AccountHolderID *string                   `json:"account_holder_id,omitempty"`
}

// XenPlatformErrorResponse represents an error returned by the Xendit API.
type XenPlatformErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// XenPlatformProvider defines the interface for Xendit XenPlatform account management operations.
// Implementations target the Xendit REST API at https://api.xendit.co.
type XenPlatformProvider interface {
	// CreateAccount creates a new sub-account (OWNED or MANAGED) on Xendit.
	CreateAccount(ctx context.Context, req XenPlatformCreateAccountRequest) (*XenPlatformAccount, error)

	// ListAccounts returns a paginated, filtered list of sub-accounts.
	ListAccounts(ctx context.Context, params XenPlatformListAccountsParams) (*XenPlatformListAccountsResponse, error)

	// GetAccount retrieves a single sub-account by its Xendit ID.
	GetAccount(ctx context.Context, id string) (*XenPlatformAccount, error)

	// UpdateAccount partially updates a sub-account (email, public_profile, account_holder_id).
	UpdateAccount(ctx context.Context, id string, req XenPlatformUpdateAccountRequest) (*XenPlatformAccount, error)
}
