package payment_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func setupPayoutTestServer(t *testing.T, handler http.HandlerFunc) (domain.XenditPayoutProvider, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := payment.NewXenditPayoutClient(config.XenditConfig{
		APISecretKey: "test-api-secret-key",
		APIPublicKey: "test-api-public-key",
		BaseURL:      server.URL,
	})
	if err != nil {
		t.Fatalf("failed to create xendit payout client: %v", err)
	}

	return client, server
}

func setupInvoiceTestServer(t *testing.T, handler http.HandlerFunc) (domain.XenditInvoiceProvider, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := payment.NewXenditInvoiceClient(config.XenditConfig{
		APISecretKey: "test-api-secret-key",
		APIPublicKey: "test-api-public-key",
		BaseURL:      server.URL,
	})
	if err != nil {
		t.Fatalf("failed to create xendit invoice client: %v", err)
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

func TestXenditPayoutClient_CreatePayout(t *testing.T) {
	req := domain.XenditPayoutRequest{
		ReferenceID: "wd-123",
		ChannelCode: "ID_BCA",
		ChannelProperties: domain.XenditPayoutChannelProperties{
			AccountNumber:     "1234567890",
			AccountHolderName: "Ahmad",
		},
		Amount:      100000,
		Description: "Withdrawal test",
		Currency:    "IDR",
	}

	t.Run("success", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/v2/payouts" {
				t.Errorf("expected path /v2/payouts, got %s", r.URL.Path)
			}
			if r.Header.Get("for-user-id") != "acc_123" {
				t.Errorf("expected for-user-id acc_123, got %s", r.Header.Get("for-user-id"))
			}
			if r.Header.Get("Idempotency-key") != "idem-123" {
				t.Errorf("expected Idempotency-key idem-123, got %s", r.Header.Get("Idempotency-key"))
			}

			var body domain.XenditPayoutRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode request body: %v", err)
			}
			if body.ReferenceID != req.ReferenceID {
				t.Errorf("expected reference_id %s, got %s", req.ReferenceID, body.ReferenceID)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(domain.XenditPayoutResponse{
				ID:          "payout_123",
				ReferenceID: req.ReferenceID,
				Status:      "ACCEPTED",
				ChannelCode: req.ChannelCode,
				Amount:      req.Amount,
				Currency:    req.Currency,
			})
		})

		resp, err := client.CreatePayout(context.Background(), "acc_123", "idem-123", req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != "payout_123" {
			t.Errorf("expected payout id payout_123, got %s", resp.ID)
		}
		if resp.Status != "ACCEPTED" {
			t.Errorf("expected status ACCEPTED, got %s", resp.Status)
		}
	})

	t.Run("api error", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(domain.XenPlatformErrorResponse{
				ErrorCode: "NOT_FOUND",
				Message:   "The requested resource was not found",
			})
		})

		_, err := client.CreatePayout(context.Background(), "acc_123", "idem-123", req)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "status 404") {
			t.Errorf("expected error to contain status 404, got %v", err)
		}
	})
}

func TestXenditPayoutClient_ListPayoutChannels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/payouts_channels" {
				t.Errorf("expected path /payouts_channels, got %s", r.URL.Path)
			}
			if r.URL.Query().Get("currency") != "IDR" {
				t.Errorf("expected currency=IDR, got %s", r.URL.Query().Get("currency"))
			}
			if r.URL.Query().Get("channel_category") != "BANK" {
				t.Errorf("expected channel_category=BANK, got %s", r.URL.Query().Get("channel_category"))
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"channel_code":     "ID_BCA",
						"channel_name":     "Bank Central Asia",
						"currency":         "IDR",
						"channel_category": "BANK",
						"is_activated":     true,
					},
				},
			})
		})

		channels, err := client.ListPayoutChannels(context.Background(), domain.XenditListPayoutChannelsParams{
			Currency:        "IDR",
			ChannelCategory: "BANK",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(channels) != 1 {
			t.Fatalf("expected 1 channel, got %d", len(channels))
		}
		if channels[0].ChannelCode != "ID_BCA" {
			t.Errorf("expected channel_code ID_BCA, got %s", channels[0].ChannelCode)
		}
		if !channels[0].IsActivated {
			t.Error("expected channel to be activated")
		}
	})

	t.Run("api error", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(domain.XenPlatformErrorResponse{
				ErrorCode: "GATEWAY_ERROR",
				Message:   "xendit upstream error",
			})
		})

		_, err := client.ListPayoutChannels(context.Background(), domain.XenditListPayoutChannelsParams{
			Currency:        "IDR",
			ChannelCategory: "BANK",
		})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "status 502") {
			t.Errorf("expected error to contain status 502, got %v", err)
		}
	})
}

func TestXenditInvoiceClient_ExpireInvoice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, _ := setupInvoiceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/v2/invoices/inv_123/expire!" {
				t.Errorf("expected path /v2/invoices/inv_123/expire!, got %s", r.URL.Path)
			}
			if r.Header.Get("for-user-id") != "acc_123" {
				t.Errorf("expected for-user-id acc_123, got %s", r.Header.Get("for-user-id"))
			}
			w.WriteHeader(http.StatusOK)
		})

		if err := client.ExpireInvoice(context.Background(), "acc_123", "inv_123"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("api error", func(t *testing.T) {
		client, _ := setupInvoiceTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(domain.XenPlatformErrorResponse{
				ErrorCode: "INVOICE_ALREADY_PAID",
				Message:   "Invoice has already been paid",
			})
		})

		err := client.ExpireInvoice(context.Background(), "acc_123", "inv_123")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "status 409") {
			t.Errorf("expected error to contain status 409, got %v", err)
		}
	})
}

func TestXenditPayoutClient_GetTransactionByReference(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/transactions" {
				t.Errorf("expected path /transactions, got %s", r.URL.Path)
			}
			if r.Header.Get("for-user-id") != "acc_123" {
				t.Errorf("expected for-user-id acc_123, got %s", r.Header.Get("for-user-id"))
			}
			if r.URL.Query().Get("reference_id") != "wd-123" {
				t.Errorf("expected reference_id wd-123, got %s", r.URL.Query().Get("reference_id"))
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":           "txn_123",
						"reference_id": "wd-123",
						"amount":       96000,
						"fee":          4000,
						"status":       "SUCCEEDED",
						"currency":     "IDR",
						"channel_code": "ID_BCA",
					},
				},
			})
		})

		tx, err := client.GetTransactionByReference(context.Background(), "acc_123", "wd-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tx.ID != "txn_123" {
			t.Errorf("expected tx id txn_123, got %s", tx.ID)
		}
		if tx.Fee != 4000 {
			t.Errorf("expected fee 4000, got %f", tx.Fee)
		}
		if tx.Amount != 96000 {
			t.Errorf("expected amount 96000, got %f", tx.Amount)
		}
	})

	t.Run("not found", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{},
			})
		})

		_, err := client.GetTransactionByReference(context.Background(), "acc_123", "wd-404")
		if !errors.Is(err, domain.ErrXenditTransactionNotFound) {
			t.Fatalf("expected ErrXenditTransactionNotFound, got %v", err)
		}
	})

	t.Run("fee unavailable", func(t *testing.T) {
		client, _ := setupPayoutTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{
						"id":           "txn_123",
						"reference_id": "wd-123",
						"amount":       100000,
						"status":       "SUCCEEDED",
						"currency":     "IDR",
						"channel_code": "ID_BCA",
					},
				},
			})
		})

		_, err := client.GetTransactionByReference(context.Background(), "acc_123", "wd-123")
		if !errors.Is(err, domain.ErrXenditTransactionFeeUnavailable) {
			t.Fatalf("expected ErrXenditTransactionFeeUnavailable, got %v", err)
		}
	})
}
