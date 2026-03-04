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

func (c *xenditClient) ExpireInvoice(ctx context.Context, forUserID string, invoiceID string) error {
	headers := map[string]string{
		"for-user-id": forUserID,
	}

	path := fmt.Sprintf("/v2/invoices/%s/expire!", url.PathEscape(invoiceID))
	resp, err := c.doRequestWithHeaders(ctx, http.MethodPost, path, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to call Xendit expire invoice API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// --- XenditPayoutProvider implementation ---

// NewXenditPayoutClient creates a new XenditPayoutProvider backed by the Xendit REST API.
func NewXenditPayoutClient(cfg config.XenditConfig) (domain.XenditPayoutProvider, error) {
	return newXenditClientInternal(cfg)
}

func (c *xenditClient) ListPayoutChannels(ctx context.Context, params domain.XenditListPayoutChannelsParams) ([]domain.XenditPayoutChannel, error) {
	query := url.Values{}
	if params.Currency != "" {
		query.Set("currency", params.Currency)
	}
	if params.ChannelCategory != "" {
		query.Set("channel_category", params.ChannelCategory)
	}

	path := "/payouts_channels"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit list payout channels API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payout channels response: %w", err)
	}

	channels, err := parseXenditPayoutChannels(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payout channels response: %w", err)
	}

	return channels, nil
}

func (c *xenditClient) CreatePayout(ctx context.Context, forUserID string, idempotencyKey string, req domain.XenditPayoutRequest) (*domain.XenditPayoutResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create payout request: %w", err)
	}

	headers := map[string]string{
		"for-user-id":     forUserID,
		"Idempotency-key": idempotencyKey,
	}

	resp, err := c.doRequestWithHeaders(ctx, http.MethodPost, "/v2/payouts", bytes.NewReader(body), headers)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit create payout API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var payout domain.XenditPayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&payout); err != nil {
		return nil, fmt.Errorf("failed to decode create payout response: %w", err)
	}

	return &payout, nil
}

func (c *xenditClient) GetTransactionByReference(ctx context.Context, forUserID string, referenceID string) (*domain.XenditTransaction, error) {
	query := url.Values{}
	query.Set("reference_id", referenceID)
	query.Set("limit", "1")

	path := "/transactions?" + query.Encode()
	headers := map[string]string{
		"for-user-id": forUserID,
	}

	resp, err := c.doRequestWithHeaders(ctx, http.MethodGet, path, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to call Xendit transactions API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read transactions response: %w", err)
	}

	items, err := parseXenditTransactions(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transactions response: %w", err)
	}
	if len(items) == 0 {
		return nil, domain.ErrXenditTransactionNotFound
	}

	fee, err := extractTransactionFee(items[0])
	if err != nil {
		return nil, err
	}

	return &domain.XenditTransaction{
		ID:          items[0].ID,
		ReferenceID: items[0].ReferenceID,
		Amount:      items[0].Amount,
		Fee:         fee,
		Status:      items[0].Status,
		Currency:    items[0].Currency,
		ChannelCode: items[0].ChannelCode,
	}, nil
}

type xenditTransactionListResponse struct {
	Data []xenditTransactionItem `json:"data"`
}

type xenditPayoutChannelsListResponse struct {
	Data []xenditPayoutChannelItem `json:"data"`
}

type xenditPayoutChannelItem struct {
	ChannelCode     string `json:"channel_code"`
	ChannelName     string `json:"channel_name"`
	Currency        string `json:"currency"`
	ChannelCategory string `json:"channel_category"`
	IsActivated     *bool  `json:"is_activated"`
	IsAvailable     *bool  `json:"is_available"`
	IsEnabled       *bool  `json:"is_enabled"`
}

type xenditTransactionItem struct {
	ID          string                 `json:"id"`
	ReferenceID string                 `json:"reference_id"`
	Amount      float64                `json:"amount"`
	Fee         *float64               `json:"fee"`
	FeeAmount   *float64               `json:"fee_amount"`
	Fees        []xenditTransactionFee `json:"fees"`
	Status      string                 `json:"status"`
	Currency    string                 `json:"currency"`
	ChannelCode string                 `json:"channel_code"`
}

type xenditTransactionFee struct {
	Amount float64 `json:"amount"`
}

func parseXenditPayoutChannels(bodyBytes []byte) ([]domain.XenditPayoutChannel, error) {
	var wrapped xenditPayoutChannelsListResponse
	if err := json.Unmarshal(bodyBytes, &wrapped); err == nil && wrapped.Data != nil {
		return toDomainPayoutChannels(wrapped.Data), nil
	}

	var list []xenditPayoutChannelItem
	if err := json.Unmarshal(bodyBytes, &list); err == nil {
		return toDomainPayoutChannels(list), nil
	}

	var single xenditPayoutChannelItem
	if err := json.Unmarshal(bodyBytes, &single); err == nil && single.ChannelCode != "" {
		return toDomainPayoutChannels([]xenditPayoutChannelItem{single}), nil
	}

	return nil, fmt.Errorf("unexpected payout channels response payload")
}

func toDomainPayoutChannels(items []xenditPayoutChannelItem) []domain.XenditPayoutChannel {
	result := make([]domain.XenditPayoutChannel, 0, len(items))
	for _, item := range items {
		isActivated := true
		if item.IsActivated != nil {
			isActivated = *item.IsActivated
		} else if item.IsAvailable != nil {
			isActivated = *item.IsAvailable
		} else if item.IsEnabled != nil {
			isActivated = *item.IsEnabled
		}

		result = append(result, domain.XenditPayoutChannel{
			ChannelCode:     item.ChannelCode,
			ChannelName:     item.ChannelName,
			Currency:        item.Currency,
			ChannelCategory: item.ChannelCategory,
			IsActivated:     isActivated,
		})
	}
	return result
}

func parseXenditTransactions(bodyBytes []byte) ([]xenditTransactionItem, error) {
	var wrapped xenditTransactionListResponse
	if err := json.Unmarshal(bodyBytes, &wrapped); err == nil && wrapped.Data != nil {
		return wrapped.Data, nil
	}

	var list []xenditTransactionItem
	if err := json.Unmarshal(bodyBytes, &list); err == nil {
		return list, nil
	}

	var single xenditTransactionItem
	if err := json.Unmarshal(bodyBytes, &single); err == nil && single.ID != "" {
		return []xenditTransactionItem{single}, nil
	}

	return nil, fmt.Errorf("unexpected transactions response payload")
}

func extractTransactionFee(item xenditTransactionItem) (float64, error) {
	if item.Fee != nil {
		return *item.Fee, nil
	}
	if item.FeeAmount != nil {
		return *item.FeeAmount, nil
	}
	if len(item.Fees) > 0 {
		total := 0.0
		for _, f := range item.Fees {
			total += f.Amount
		}
		return total, nil
	}
	return 0, domain.ErrXenditTransactionFeeUnavailable
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
