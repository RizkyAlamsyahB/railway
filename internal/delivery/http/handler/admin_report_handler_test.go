package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type stubAdminReportUseCase struct {
	getReport    func(ctx context.Context, params domain.AdminReportParams) (*domain.AdminReportResponse, *domain.PaginationMeta, error)
	exportReport func(ctx context.Context, params domain.AdminReportParams) ([]domain.AdminReportItem, error)
}

func (s stubAdminReportUseCase) GetReport(ctx context.Context, params domain.AdminReportParams) (*domain.AdminReportResponse, *domain.PaginationMeta, error) {
	return s.getReport(ctx, params)
}

func (s stubAdminReportUseCase) ExportReport(ctx context.Context, params domain.AdminReportParams) ([]domain.AdminReportItem, error) {
	return s.exportReport(ctx, params)
}

type adminReportEnvelope struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    domain.AdminReportResponse `json:"data"`
	Meta    domain.PaginationMeta      `json:"meta"`
}

func TestAdminReportHandlerGetReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAdminReportHandler(stubAdminReportUseCase{
		getReport: func(ctx context.Context, params domain.AdminReportParams) (*domain.AdminReportResponse, *domain.PaginationMeta, error) {
			if params.Year != 2026 || params.Page != 2 || params.Limit != 5 || params.Search != "Jan" {
				t.Fatalf("unexpected params: %+v", params)
			}

			return &domain.AdminReportResponse{
					Year: 2026,
					Summary: domain.AdminReportSummary{
						SuccessTransactionRate: 82,
						TotalTransactions:      8,
						CustomerCount:          400,
						PlatformCommission:     430000000,
					},
					Items: []domain.AdminReportItem{
						{
							Period:           "Jan 2026",
							TotalTransaction: 4000000,
							DownPayment:      3000000,
							Settlement:       4000000,
							Refund:           50000,
						},
					},
				}, &domain.PaginationMeta{
					Page:       2,
					Limit:      5,
					TotalItems: 12,
					TotalPages: 3,
				}, nil
		},
		exportReport: func(ctx context.Context, params domain.AdminReportParams) ([]domain.AdminReportItem, error) {
			return nil, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/reports?year=2026&page=2&limit=5&search=Jan", nil)

	h.GetReport(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var env adminReportEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !env.Success {
		t.Fatalf("expected success=true")
	}
	if env.Message != "admin report retrieved successfully" {
		t.Fatalf("unexpected message: %s", env.Message)
	}
	if env.Data.Year != 2026 {
		t.Fatalf("unexpected year: %d", env.Data.Year)
	}
	if env.Meta.Page != 2 || env.Meta.TotalPages != 3 {
		t.Fatalf("unexpected meta: %+v", env.Meta)
	}
}

func TestAdminReportHandlerGetReportInvalidYear(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewAdminReportHandler(stubAdminReportUseCase{
		getReport: func(ctx context.Context, params domain.AdminReportParams) (*domain.AdminReportResponse, *domain.PaginationMeta, error) {
			t.Fatal("usecase should not be called")
			return nil, nil, nil
		},
		exportReport: func(ctx context.Context, params domain.AdminReportParams) ([]domain.AdminReportItem, error) {
			t.Fatal("usecase should not be called")
			return nil, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/admin/reports?year=1800", nil)

	h.GetReport(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
