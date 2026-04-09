package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// ---------- helpers ----------

func setupAdminPaymentUseCase(t *testing.T) (*mocks.MockPaymentRepository, domain.AdminPaymentUseCase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockPaymentRepository(ctrl)
	uc := NewAdminPaymentUseCase(repo)
	return repo, uc
}

func dummyPaymentItems(n int) []domain.AdminPaymentListItem {
	items := make([]domain.AdminPaymentListItem, n)
	method := "BANK_TRANSFER"
	for i := range items {
		items[i] = domain.AdminPaymentListItem{
			InvoiceID:     uuid.New(),
			OrderID:       uuid.New(),
			OrderNo:       "ORD-20260327-ABC" + string(rune('0'+i)),
			CustomerName:  "Customer " + string(rune('A'+i)),
			VendorName:    "Vendor " + string(rune('A'+i)),
			PaymentMethod: &method,
			Amount:        float64((i + 1) * 50000),
			Status:        "paid",
			CreatedAt:     time.Now(),
		}
	}
	return items
}

// ============================================================
// List
// ============================================================

func TestAdminPaymentList(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.AdminPaymentListParams
		setupMock func(repo *mocks.MockPaymentRepository, ctx context.Context)
		wantErr   bool
		checkResp func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta)
	}{
		{
			name:   "success with results",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 10},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				items := dummyPaymentItems(2)
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10}).Return(items, int64(2), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 2 {
					t.Errorf("expected 2 items, got %d", len(items))
				}
				if meta.TotalItems != 2 {
					t.Errorf("expected total_items=2, got %d", meta.TotalItems)
				}
				if meta.TotalPages != 1 {
					t.Errorf("expected total_pages=1, got %d", meta.TotalPages)
				}
				if meta.Page != 1 {
					t.Errorf("expected page=1, got %d", meta.Page)
				}
				if meta.Limit != 10 {
					t.Errorf("expected limit=10, got %d", meta.Limit)
				}
			},
		},
		{
			name:   "empty results",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 10},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10}).Return([]domain.AdminPaymentListItem{}, int64(0), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 0 {
					t.Errorf("expected 0 items, got %d", len(items))
				}
				if meta.TotalItems != 0 {
					t.Errorf("expected total_items=0, got %d", meta.TotalItems)
				}
				if meta.TotalPages != 0 {
					t.Errorf("expected total_pages=0, got %d", meta.TotalPages)
				}
			},
		},
		{
			name:   "defaults invalid page and limit",
			params: domain.AdminPaymentListParams{Page: 0, Limit: 0},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10}).Return([]domain.AdminPaymentListItem{}, int64(0), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if meta.Page != 1 {
					t.Errorf("expected page=1 after default, got %d", meta.Page)
				}
				if meta.Limit != 10 {
					t.Errorf("expected limit=10 after default, got %d", meta.Limit)
				}
			},
		},
		{
			name:   "negative page defaults to 1",
			params: domain.AdminPaymentListParams{Page: -5, Limit: 10},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10}).Return([]domain.AdminPaymentListItem{}, int64(0), nil)
			},
		},
		{
			name:   "limit exceeds 100 defaults to 10",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 200},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10}).Return([]domain.AdminPaymentListItem{}, int64(0), nil)
			},
		},
		{
			name:   "pagination calculation - multiple pages",
			params: domain.AdminPaymentListParams{Page: 2, Limit: 3},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				items := dummyPaymentItems(3)
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 2, Limit: 3}).Return(items, int64(7), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				// 7 items / 3 per page = ceil(2.33) = 3 pages
				if meta.TotalPages != 3 {
					t.Errorf("expected total_pages=3, got %d", meta.TotalPages)
				}
				if meta.Page != 2 {
					t.Errorf("expected page=2, got %d", meta.Page)
				}
			},
		},
		{
			name:   "with status filter",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 10, Status: "paid"},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				items := dummyPaymentItems(1)
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10, Status: "paid"}).Return(items, int64(1), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 1 {
					t.Errorf("expected 1 item, got %d", len(items))
				}
			},
		},
		{
			name:   "with search filter",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 10, Search: "ahmad"},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				items := dummyPaymentItems(1)
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10, Search: "ahmad"}).Return(items, int64(1), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 1 {
					t.Errorf("expected 1 item, got %d", len(items))
				}
			},
		},
		{
			name:   "with sort params",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 10, SortBy: "amount", SortOrder: "desc"},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				items := dummyPaymentItems(2)
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10, SortBy: "amount", SortOrder: "desc"}).Return(items, int64(2), nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem, meta *domain.PaginationMeta) {
				t.Helper()
				if len(items) != 2 {
					t.Errorf("expected 2 items, got %d", len(items))
				}
			},
		},
		{
			name:   "repo error",
			params: domain.AdminPaymentListParams{Page: 1, Limit: 10},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ListForAdmin(ctx, domain.AdminPaymentListParams{Page: 1, Limit: 10}).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupAdminPaymentUseCase(t)
			ctx := context.Background()

			tc.setupMock(repo, ctx)

			items, meta, err := uc.List(ctx, tc.params)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, items, meta)
			}
		})
	}
}

// ============================================================
// Summary
// ============================================================

func TestAdminPaymentSummary(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.AdminPaymentListParams
		setupMock func(repo *mocks.MockPaymentRepository, ctx context.Context)
		wantErr   bool
		checkResp func(t *testing.T, summary *domain.AdminPaymentSummary)
	}{
		{
			name:   "success",
			params: domain.AdminPaymentListParams{},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				summary := &domain.AdminPaymentSummary{
					TotalAmount:   250000,
					PaidAmount:    200000,
					PendingAmount: 30000,
					FailedAmount:  20000,
				}
				repo.EXPECT().SummaryForAdmin(ctx, gomock.Any()).Return(summary, nil)
			},
			checkResp: func(t *testing.T, summary *domain.AdminPaymentSummary) {
				if summary.TotalAmount != 250000 {
					t.Errorf("expected TotalAmount 250000, got %v", summary.TotalAmount)
				}
				if summary.PaidAmount != 200000 {
					t.Errorf("expected PaidAmount 200000, got %v", summary.PaidAmount)
				}
			},
		},
		{
			name: "with date range",
			params: func() domain.AdminPaymentListParams {
				start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
				end := time.Date(2026, 4, 7, 0, 0, 0, 0, time.UTC)
				return domain.AdminPaymentListParams{StartDate: &start, EndDate: &end}
			}(),
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				summary := &domain.AdminPaymentSummary{TotalAmount: 100000, PaidAmount: 100000}
				repo.EXPECT().SummaryForAdmin(ctx, gomock.Any()).Return(summary, nil)
			},
			checkResp: func(t *testing.T, summary *domain.AdminPaymentSummary) {
				if summary.TotalAmount != 100000 {
					t.Errorf("expected TotalAmount 100000, got %v", summary.TotalAmount)
				}
			},
		},
		{
			name:   "repo error",
			params: domain.AdminPaymentListParams{},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().SummaryForAdmin(ctx, gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupAdminPaymentUseCase(t)
			ctx := context.Background()

			tc.setupMock(repo, ctx)

			summary, err := uc.Summary(ctx, tc.params)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, summary)
			}
		})
	}
}

// ============================================================
// Export
// ============================================================

func TestAdminPaymentExport(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.AdminPaymentListParams
		setupMock func(repo *mocks.MockPaymentRepository, ctx context.Context)
		wantErr   bool
		checkResp func(t *testing.T, items []domain.AdminPaymentListItem)
	}{
		{
			name:   "success",
			params: domain.AdminPaymentListParams{},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				items := dummyPaymentItems(3)
				repo.EXPECT().ExportForAdmin(ctx, gomock.Any()).Return(items, nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem) {
				if len(items) != 3 {
					t.Errorf("expected 3 items, got %d", len(items))
				}
			},
		},
		{
			name:   "empty results",
			params: domain.AdminPaymentListParams{},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ExportForAdmin(ctx, gomock.Any()).Return([]domain.AdminPaymentListItem{}, nil)
			},
			checkResp: func(t *testing.T, items []domain.AdminPaymentListItem) {
				if len(items) != 0 {
					t.Errorf("expected 0 items, got %d", len(items))
				}
			},
		},
		{
			name:   "repo error",
			params: domain.AdminPaymentListParams{},
			setupMock: func(repo *mocks.MockPaymentRepository, ctx context.Context) {
				repo.EXPECT().ExportForAdmin(ctx, gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo, uc := setupAdminPaymentUseCase(t)
			ctx := context.Background()

			tc.setupMock(repo, ctx)

			items, err := uc.Export(ctx, tc.params)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.checkResp != nil {
				tc.checkResp(t, items)
			}
		})
	}
}
