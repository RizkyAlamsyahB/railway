package payment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type xenditClient struct {
	httpClient *http.Client
	baseURL    string
	authHeader string
}

// newXenditClientInternal creates the shared xenditClient from config.
func newXenditClientInternal(cfg config.XenditConfig) (*xenditClient, error) {
	if cfg.APISecretKey == "" {
		return nil, fmt.Errorf("XENDIT_API_SECRET_KEY is required")
	}
	if cfg.APIPublicKey == "" {
		return nil, fmt.Errorf("XENDIT_API_PUBLIC_KEY is required")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("XENDIT_BASE_URL is required")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(cfg.APISecretKey + ":"))

	return &xenditClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    cfg.BaseURL,
		authHeader: "Basic " + encoded,
	}, nil
}

// NewXenditClient creates a new XenPlatformProvider backed by the Xendit REST API.
func NewXenditClient(cfg config.XenditConfig) (domain.XenPlatformProvider, error) {
	return newXenditClientInternal(cfg)
}

// NewXenditInvoiceClient creates a new XenditInvoiceProvider backed by the Xendit REST API.
func NewXenditInvoiceClient(cfg config.XenditConfig) (domain.XenditInvoiceProvider, error) {
	return newXenditClientInternal(cfg)
}

// --- XenPlatformProvider implementation ---

func (c *xenditClient) CreateAccount(ctx context.Context, req domain.XenPlatformCreateAccountRequest) (*domain.XenPlatformAccount, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create account request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/accounts", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit create account API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var account domain.XenPlatformAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("failed to decode create account response: %w", err)
	}

	return &account, nil
}

func (c *xenditClient) ListAccounts(ctx context.Context, params domain.XenPlatformListAccountsParams) (*domain.XenPlatformListAccountsResponse, error) {
	query := url.Values{}
	if params.Email != "" {
		query.Set("email", params.Email)
	}
	if params.Status != "" {
		query.Set("status", params.Status)
	}
	if params.BusinessName != "" {
		query.Set("public_profile.business_name", params.BusinessName)
	}
	if params.Type != "" {
		query.Set("type", params.Type)
	}
	if params.CreatedGTE != "" {
		query.Set("created[gte]", params.CreatedGTE)
	}
	if params.CreatedLTE != "" {
		query.Set("created[lte]", params.CreatedLTE)
	}
	if params.Limit > 0 {
		query.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.AfterID != "" {
		query.Set("after_id", params.AfterID)
	}

	path := "/v2/accounts"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit list accounts API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result domain.XenPlatformListAccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode list accounts response: %w", err)
	}

	return &result, nil
}

func (c *xenditClient) GetAccount(ctx context.Context, id string) (*domain.XenPlatformAccount, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/v2/accounts/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit get account API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var account domain.XenPlatformAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("failed to decode get account response: %w", err)
	}

	return &account, nil
}

func (c *xenditClient) UpdateAccount(ctx context.Context, id string, req domain.XenPlatformUpdateAccountRequest) (*domain.XenPlatformAccount, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update account request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, "/v2/accounts/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit update account API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var account domain.XenPlatformAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, fmt.Errorf("failed to decode update account response: %w", err)
	}

	return &account, nil
}

// --- XenditInvoiceProvider implementation ---

func (c *xenditClient) CreateInvoice(ctx context.Context, forUserID string, req domain.XenditInvoiceRequest) (*domain.XenditInvoiceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create invoice request: %w", err)
	}

	headers := map[string]string{
		"for-user-id": forUserID,
	}

	resp, err := c.doRequestWithHeaders(ctx, http.MethodPost, "/v2/invoices", bytes.NewReader(body), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit create invoice API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var invoice domain.XenditInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, fmt.Errorf("failed to decode create invoice response: %w", err)
	}

	return &invoice, nil
}

// --- HTTP helpers ---

func (c *xenditClient) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	return c.doRequestWithHeaders(ctx, method, path, body, nil)
}

func (c *xenditClient) doRequestWithHeaders(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}

func (c *xenditClient) handleErrorResponse(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("xendit API error (status %d): failed to read response body: %w", resp.StatusCode, err)
	}

	var xenErr domain.XenPlatformErrorResponse
	if err := json.Unmarshal(bodyBytes, &xenErr); err != nil {
		return fmt.Errorf("xendit API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return fmt.Errorf("xendit API error (status %d, code %s): %s", resp.StatusCode, xenErr.ErrorCode, xenErr.Message)
}
