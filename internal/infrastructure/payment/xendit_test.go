package payment_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/infrastructure/payment"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (domain.XenPlatformProvider, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := payment.NewXenditClient(config.XenditConfig{
		APISecretKey: "test-api-secret-key",
		APIPublicKey: "test-api-public-key",
		BaseURL:      server.URL,
	})
	if err != nil {
		t.Fatalf("failed to create xendit client: %v", err)
	}

	return client, server
}

func TestNewXenditClient(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.XenditConfig
		wantErr bool
	}{
		{
			name:    "empty API secret key returns error",
			cfg:     config.XenditConfig{APISecretKey: "", APIPublicKey: "test-api-public-key", BaseURL: "https://api.xendit.co"},
			wantErr: true,
		},
		{
			name:    "empty API public key returns error",
			cfg:     config.XenditConfig{APISecretKey: "test-api-secret-key", APIPublicKey: "", BaseURL: "https://api.xendit.co"},
			wantErr: true,
		},
		{
			name:    "valid config returns client",
			cfg:     config.XenditConfig{APISecretKey: "test-api-secret-key", APIPublicKey: "test-api-public-key", BaseURL: "https://api.xendit.co"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := payment.NewXenditClient(tt.cfg)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("expected non-nil client")
			}
		})
	}
}

func TestXenditClient_CreateAccount(t *testing.T) {
	tests := []struct {
		name       string
		req        domain.XenPlatformCreateAccountRequest
		serverFunc http.HandlerFunc
		wantErr    bool
		wantID     string
		wantStatus string
	}{
		{
			name: "success",
			req: domain.XenPlatformCreateAccountRequest{
				Email: "seller@example.com",
				Type:  "OWNED",
				PublicProfile: &domain.XenPlatformPublicProfile{
					BusinessName: "Test Store",
				},
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/v2/accounts" {
					t.Errorf("expected path /v2/accounts, got %s", r.URL.Path)
				}
				if r.Header.Get("Authorization") == "" {
					t.Error("expected Authorization header")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
				}

				var body domain.XenPlatformCreateAccountRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}
				if body.Email != "seller@example.com" {
					t.Errorf("expected email seller@example.com, got %s", body.Email)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(domain.XenPlatformAccount{
					ID:     "acc_123",
					Type:   "OWNED",
					Email:  "seller@example.com",
					Status: "LIVE",
					PublicProfile: &domain.XenPlatformPublicProfile{
						BusinessName: "Test Store",
					},
				})
			},
			wantErr:    false,
			wantID:     "acc_123",
			wantStatus: "LIVE",
		},
		{
			name: "API error",
			req: domain.XenPlatformCreateAccountRequest{
				Email: "seller@example.com",
				Type:  "OWNED",
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(domain.XenPlatformErrorResponse{
					ErrorCode: "API_VALIDATION_ERROR",
					Message:   "email is required",
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setupTestServer(t, tt.serverFunc)

			account, err := client.CreateAccount(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if account.ID != tt.wantID {
				t.Errorf("expected ID %s, got %s", tt.wantID, account.ID)
			}
			if account.Status != tt.wantStatus {
				t.Errorf("expected status %s, got %s", tt.wantStatus, account.Status)
			}
		})
	}
}

func TestXenditClient_GetAccount(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		serverFunc http.HandlerFunc
		wantErr    bool
		wantID     string
	}{
		{
			name: "success",
			id:   "acc_123",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/v2/accounts/acc_123" {
					t.Errorf("expected path /v2/accounts/acc_123, got %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(domain.XenPlatformAccount{
					ID:     "acc_123",
					Type:   "OWNED",
					Email:  "seller@example.com",
					Status: "LIVE",
				})
			},
			wantErr: false,
			wantID:  "acc_123",
		},
		{
			name: "not found",
			id:   "acc_nonexistent",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(domain.XenPlatformErrorResponse{
					ErrorCode: "DATA_NOT_FOUND",
					Message:   "Account not found",
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setupTestServer(t, tt.serverFunc)

			account, err := client.GetAccount(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if account.ID != tt.wantID {
				t.Errorf("expected ID %s, got %s", tt.wantID, account.ID)
			}
		})
	}
}

func TestXenditClient_ListAccounts(t *testing.T) {
	tests := []struct {
		name        string
		params      domain.XenPlatformListAccountsParams
		serverFunc  http.HandlerFunc
		wantErr     bool
		wantCount   int
		wantHasMore bool
	}{
		{
			name: "with filters",
			params: domain.XenPlatformListAccountsParams{
				Type:         "OWNED",
				Limit:        5,
				BusinessName: "Test",
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}

				q := r.URL.Query()
				if q.Get("type") != "OWNED" {
					t.Errorf("expected type=OWNED, got %s", q.Get("type"))
				}
				if q.Get("limit") != "5" {
					t.Errorf("expected limit=5, got %s", q.Get("limit"))
				}
				if q.Get("public_profile.business_name") != "Test" {
					t.Errorf("expected public_profile.business_name=Test, got %s", q.Get("public_profile.business_name"))
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(domain.XenPlatformListAccountsResponse{
					Data: []domain.XenPlatformAccount{
						{ID: "acc_1", Status: "LIVE"},
						{ID: "acc_2", Status: "LIVE"},
					},
					HasMore: true,
					Links: []domain.XenPlatformListLink{
						{Href: "/v2/accounts?after_id=acc_2", Rel: "next", Method: "GET"},
					},
				})
			},
			wantErr:     false,
			wantCount:   2,
			wantHasMore: true,
		},
		{
			name:   "empty params",
			params: domain.XenPlatformListAccountsParams{},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != "" {
					t.Errorf("expected no query params, got %s", r.URL.RawQuery)
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(domain.XenPlatformListAccountsResponse{
					Data:    []domain.XenPlatformAccount{},
					HasMore: false,
				})
			},
			wantErr:     false,
			wantCount:   0,
			wantHasMore: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setupTestServer(t, tt.serverFunc)

			result, err := client.ListAccounts(context.Background(), tt.params)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result.Data) != tt.wantCount {
				t.Errorf("expected %d accounts, got %d", tt.wantCount, len(result.Data))
			}
			if result.HasMore != tt.wantHasMore {
				t.Errorf("expected has_more=%v, got %v", tt.wantHasMore, result.HasMore)
			}
		})
	}
}

func TestXenditClient_UpdateAccount(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		req        domain.XenPlatformUpdateAccountRequest
		serverFunc http.HandlerFunc
		wantErr    bool
		wantName   string
	}{
		{
			name: "update business name",
			id:   "acc_123",
			req: domain.XenPlatformUpdateAccountRequest{
				PublicProfile: &domain.XenPlatformPublicProfile{
					BusinessName: "New Name",
				},
			},
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPatch {
					t.Errorf("expected PATCH, got %s", r.Method)
				}
				if r.URL.Path != "/v2/accounts/acc_123" {
					t.Errorf("expected path /v2/accounts/acc_123, got %s", r.URL.Path)
				}

				bodyBytes, _ := io.ReadAll(r.Body)
				var body map[string]interface{}
				json.Unmarshal(bodyBytes, &body)

				// Verify omitempty works: email and account_holder_id should not be present
				if _, exists := body["email"]; exists {
					t.Error("expected email to be omitted")
				}
				if _, exists := body["account_holder_id"]; exists {
					t.Error("expected account_holder_id to be omitted")
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(domain.XenPlatformAccount{
					ID:     "acc_123",
					Type:   "OWNED",
					Email:  "seller@example.com",
					Status: "LIVE",
					PublicProfile: &domain.XenPlatformPublicProfile{
						BusinessName: "New Name",
					},
				})
			},
			wantErr:  false,
			wantName: "New Name",
		},
		{
			name: "link account holder",
			id:   "acc_123",
			req: func() domain.XenPlatformUpdateAccountRequest {
				holderID := "holder_456"
				return domain.XenPlatformUpdateAccountRequest{
					AccountHolderID: &holderID,
				}
			}(),
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				bodyBytes, _ := io.ReadAll(r.Body)
				var body map[string]interface{}
				json.Unmarshal(bodyBytes, &body)

				if body["account_holder_id"] != "holder_456" {
					t.Errorf("expected account_holder_id=holder_456, got %v", body["account_holder_id"])
				}

				holderID := "holder_456"
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(domain.XenPlatformAccount{
					ID:              "acc_123",
					Type:            "OWNED",
					Email:           "seller@example.com",
					Status:          "LIVE",
					AccountHolderID: &holderID,
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := setupTestServer(t, tt.serverFunc)

			account, err := client.UpdateAccount(context.Background(), tt.id, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantName != "" && account.PublicProfile != nil && account.PublicProfile.BusinessName != tt.wantName {
				t.Errorf("expected business name %s, got %s", tt.wantName, account.PublicProfile.BusinessName)
			}
		})
	}
}
