package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type stubAdminDashboardUseCase struct {
	getDashboard func(ctx context.Context, params domain.AdminDashboardParams) (*domain.AdminDashboardResponse, error)
}

func (s stubAdminDashboardUseCase) GetDashboard(ctx context.Context, params domain.AdminDashboardParams) (*domain.AdminDashboardResponse, error) {
	return s.getDashboard(ctx, params)
}

type adminDashboardResponseEnvelope struct {
	Success bool                          `json:"success"`
	Message string                        `json:"message"`
	Data    domain.AdminDashboardResponse `json:"data"`
}

func TestAdminDashboardHandlerGetDashboard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedTime := time.Date(2026, time.March, 1, 20, 15, 0, 0, time.UTC)
	h := NewAdminDashboardHandler(stubAdminDashboardUseCase{
		getDashboard: func(ctx context.Context, params domain.AdminDashboardParams) (*domain.AdminDashboardResponse, error) {
			if params.Month != 3 || params.Year != 2026 {
				t.Fatalf("unexpected params: %+v", params)
			}

			return &domain.AdminDashboardResponse{
				Month: 3,
				Year:  2026,
				Summary: domain.AdminDashboardSummary{
					PPIUAndPIHKCount:    128,
					DormitoryCount:      56,
					StoreCount:          342,
					PendingPaymentCount: 17,
					NeedsReviewCount:    24,
					OrderCount:          1245,
				},
				Transactions: []domain.AdminDashboardTransactionItem{
					{
						DateTime:     expectedTime,
						Invoice:      "INV-20260301-007",
						CustomerName: "Hendra Wijaya",
						ProductName:  "Paket Umroh VIP",
						Price:        35000000,
						Status:       "paid",
					},
				},
			}, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard?month=3&year=2026", nil)

	h.GetDashboard(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var env adminDashboardResponseEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !env.Success {
		t.Fatalf("expected success=true")
	}
	if env.Message != "admin dashboard retrieved successfully" {
		t.Fatalf("unexpected message: %s", env.Message)
	}
	if env.Data.Month != 3 || env.Data.Year != 2026 {
		t.Fatalf("unexpected period: %+v", env.Data)
	}
	if len(env.Data.Transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(env.Data.Transactions))
	}
	if env.Data.Transactions[0].DateTime != expectedTime {
		t.Fatalf("unexpected transaction time: %v", env.Data.Transactions[0].DateTime)
	}
}

func TestAdminDashboardHandlerGetDashboardInvalidMonth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAdminDashboardHandler(stubAdminDashboardUseCase{
		getDashboard: func(ctx context.Context, params domain.AdminDashboardParams) (*domain.AdminDashboardResponse, error) {
			t.Fatal("usecase should not be called for invalid params")
			return nil, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard?month=13&year=2026", nil)

	h.GetDashboard(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
