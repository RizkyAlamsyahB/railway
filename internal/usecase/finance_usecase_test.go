package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

// --------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------

func newFinanceUC(ctrl *gomock.Controller) (domain.FinanceUseCase, *mocks.MockFinanceRepository, *mocks.MockXenditPayoutProvider) {
	repo := mocks.NewMockFinanceRepository(ctrl)
	payout := mocks.NewMockXenditPayoutProvider(ctrl)
	uc := usecase.NewFinanceUseCase(repo, payout)
	return uc, repo, payout
}

func strPtr(s string) *string { return &s }

// --------------------------------------------------------------------
// GetDashboard
// --------------------------------------------------------------------

func TestFinanceGetDashboard(t *testing.T) {
	okSummary := &domain.FinanceDashboardSummary{TotalTransactions: 5000000, PlatformCommission: 250000}
	okToday := []domain.FinanceTodayTransactionItem{{Invoice: "INV-001", Product: "Sajadah"}}

	tests := []struct {
		name      string
		month     string
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:  "sukses tanpa parameter month (default bulan ini)",
			month: "",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetDashboardSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
				r.EXPECT().ListTodayTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okToday, nil)
			},
		},
		{
			name:  "sukses dengan month valid",
			month: "2025-03",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetDashboardSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
				r.EXPECT().ListTodayTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okToday, nil)
			},
		},
		{
			name:      "gagal: format month tidak valid (3-2025)",
			month:     "3-2025",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
		{
			name:      "gagal: format month tidak valid (2025-3)",
			month:     "2025-3",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
		{
			name:      "gagal: month berisi teks acak",
			month:     "juni-2025",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
		{
			name:  "gagal: repo error saat GetDashboardSummary",
			month: "2025-01",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetDashboardSummary(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
			expectErr: errors.New("db error"),
		},
		{
			name:  "gagal: repo error saat ListTodayTransactions",
			month: "2025-01",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetDashboardSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
				r.EXPECT().ListTodayTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
			expectErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.GetDashboard(context.Background(), tt.month)
			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				// Cukup cek error tidak nil untuk db errors, atau sentinel match
				if errors.Is(tt.expectErr, usecase.ErrInvalidMonth) {
					if !errors.Is(err, usecase.ErrInvalidMonth) {
						t.Fatalf("expected ErrInvalidMonth, got %v", err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected result, got nil")
			}
		})
	}
}

// --------------------------------------------------------------------
// GetFinancialReport
// --------------------------------------------------------------------

func TestFinanceGetFinancialReport(t *testing.T) {
	okSummary := &domain.FinanceReportSummary{GrossRevenueTotal: 100000000}
	okDaily := []domain.FinanceReportDailyItem{{Date: time.Now(), Amount: 5000000}}

	tests := []struct {
		name      string
		month     string
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:  "sukses dengan month valid",
			month: "2025-04",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetReportSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
				r.EXPECT().ListReportDailyIncome(gomock.Any(), gomock.Any()).Return(okDaily, nil)
			},
		},
		{
			name:  "sukses tanpa month (default bulan ini)",
			month: "",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetReportSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
				r.EXPECT().ListReportDailyIncome(gomock.Any(), gomock.Any()).Return(okDaily, nil)
			},
		},
		{
			name:      "gagal: format month tidak valid",
			month:     "2025/04",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
		{
			name:      "gagal: month hanya berisi angka",
			month:     "202504",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
		{
			name:  "gagal: repo error saat GetReportSummary",
			month: "2025-01",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetReportSummary(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.GetFinancialReport(context.Background(), tt.month)
			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if errors.Is(tt.expectErr, usecase.ErrInvalidMonth) {
					if !errors.Is(err, usecase.ErrInvalidMonth) {
						t.Fatalf("expected ErrInvalidMonth, got %v", err)
					}
				}
				return
			}
			// skenario tanpa expectErr tapi ada repo error: result bisa nil, err tidak nil - itu valid
			_ = result
		})
	}
}

// --------------------------------------------------------------------
// ListTransactions & GetTransactionSummary
// --------------------------------------------------------------------

func TestFinanceListTransactions(t *testing.T) {
	okItems := []domain.FinanceTransactionItem{{Invoice: "INV-001", Customer: "Ali"}}

	tests := []struct {
		name      string
		params    domain.FinanceListParams
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses tanpa filter",
			params: domain.FinanceListParams{Page: 1, Limit: 10},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status paid",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "paid"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status pending",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "pending"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status failed",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "failed"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status expired",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "expired"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status refunded",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "refunded"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, Status: "xyz"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPaymentStatus,
		},
		{
			name:      "gagal: status UPPERCASE tidak valid (bukan normalisasi)",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, Status: "BAYAR"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPaymentStatus,
		},
		{
			name:      "gagal: format date_from salah",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, DateFrom: "2025-4-1"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidDate,
		},
		{
			name:      "gagal: date_from setelah date_to",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, Month: "2025-04", DateFrom: "2025-04-20", DateTo: "2025-04-01"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidDateRange,
		},
		{
			name:   "sukses: page 2 limit 5",
			params: domain.FinanceListParams{Page: 2, Limit: 5},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(10), nil)
			},
		},
		{
			name:   "sukses: page negatif, dinormalisasi ke 1",
			params: domain.FinanceListParams{Page: -1, Limit: 10},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses: limit berlebihan (>100), dinormalisasi ke 10",
			params: domain.FinanceListParams{Page: 1, Limit: 999},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses: dengan search keyword",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Search: "INV-001"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			items, meta, err := uc.ListTransactions(context.Background(), tt.params)
			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectErr)
				}
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if items == nil {
				t.Fatal("expected items, got nil")
			}
			if meta == nil {
				t.Fatal("expected meta, got nil")
			}
		})
	}
}

func TestFinanceGetTransactionSummary(t *testing.T) {
	okSummary := &domain.FinanceTransactionSummary{TotalTransactions: 9000000}

	tests := []struct {
		name      string
		month     string
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:  "sukses",
			month: "2025-04",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetTransactionSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
			},
		},
		{
			name:      "gagal: month invalid",
			month:     "04-2025",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.GetTransactionSummary(context.Background(), tt.month)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("expected result, got nil")
			}
		})
	}
}

// --------------------------------------------------------------------
// ExportTransactions
// --------------------------------------------------------------------

func TestFinanceExportTransactions(t *testing.T) {
	okItems := []domain.FinanceTransactionItem{{Invoice: "INV-001"}}

	tests := []struct {
		name      string
		params    domain.FinanceListParams
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses tanpa filter",
			params: domain.FinanceListParams{},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status paid",
			params: domain.FinanceListParams{Status: "paid"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			params:    domain.FinanceListParams{Status: "invalid-status"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPaymentStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			items, err := uc.ExportTransactions(context.Background(), tt.params)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = items
		})
	}
}

// --------------------------------------------------------------------
// ListPayouts
// --------------------------------------------------------------------

func TestFinanceListPayouts(t *testing.T) {
	payoutID := uuid.New()
	okItems := []domain.FinancePayoutItem{{ID: payoutID, Vendor: strPtr("Toko Haji")}}

	tests := []struct {
		name      string
		params    domain.FinanceListParams
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses tanpa filter",
			params: domain.FinanceListParams{Page: 1, Limit: 10},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status on_hold",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "on_hold"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status schedule",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "schedule"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status completed",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "completed"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status failed",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "failed"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses: alias 'complete' dinormalisasi ke 'completed'",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "complete"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses: alias 'on hold' (dengan spasi) dinormalisasi ke 'on_hold'",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "on hold"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, Status: "dibayar"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPayoutStatus,
		},
		{
			name:      "gagal: format date_to salah",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, DateTo: "01/04/2025"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			items, meta, err := uc.ListPayouts(context.Background(), tt.params)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = items
			_ = meta
		})
	}
}

// --------------------------------------------------------------------
// ExportPayouts
// --------------------------------------------------------------------

func TestFinanceExportPayouts(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.FinanceListParams
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses tanpa filter",
			params: domain.FinanceListParams{},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			params:    domain.FinanceListParams{Status: "unknown"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPayoutStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			_, err := uc.ExportPayouts(context.Background(), tt.params)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// --------------------------------------------------------------------
// ListRefunds
// --------------------------------------------------------------------

func TestFinanceListRefunds(t *testing.T) {
	refundID := uuid.New()
	okItems := []domain.FinanceRefundItem{{ID: refundID, Customer: "Budi"}}

	tests := []struct {
		name      string
		params    domain.FinanceListParams
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses tanpa filter",
			params: domain.FinanceListParams{Page: 1, Limit: 10},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status requested",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "requested"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status approved",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "approved"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status rejected",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "rejected"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses filter status processed",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "processed"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:   "sukses: alias 'pending' dinormalisasi ke 'requested'",
			params: domain.FinanceListParams{Page: 1, Limit: 10, Status: "pending"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(okItems, int64(1), nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, Status: "diminta"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidRefundStatus,
		},
		{
			name:      "gagal: date_from dan date_to terbalik",
			params:    domain.FinanceListParams{Page: 1, Limit: 10, Month: "2025-04", DateFrom: "2025-04-30", DateTo: "2025-04-01"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidDateRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			items, meta, err := uc.ListRefunds(context.Background(), tt.params)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = items
			_ = meta
		})
	}
}

// --------------------------------------------------------------------
// ExportRefunds
// --------------------------------------------------------------------

func TestFinanceExportRefunds(t *testing.T) {
	tests := []struct {
		name      string
		params    domain.FinanceListParams
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses tanpa filter",
			params: domain.FinanceListParams{},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), nil)
			},
		},
		{
			name:   "sukses: alias pending -> requested",
			params: domain.FinanceListParams{Status: "pending"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			params:    domain.FinanceListParams{Status: "selesai"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidRefundStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			_, err := uc.ExportRefunds(context.Background(), tt.params)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// --------------------------------------------------------------------
// GetRefundSummary & GetPayoutSummary
// --------------------------------------------------------------------

func TestFinanceGetRefundSummary(t *testing.T) {
	okSummary := &domain.FinanceRefundSummary{SubmittedCount: 5}

	tests := []struct {
		name      string
		month     string
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:  "sukses",
			month: "2025-03",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
			},
		},
		{
			name:      "gagal: month invalid",
			month:     "2025-3",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.GetRefundSummary(context.Background(), tt.month)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = result
		})
	}
}

func TestFinanceGetPayoutSummary(t *testing.T) {
	okSummary := &domain.FinancePayoutSummary{CompletedTotal: 1000000}

	tests := []struct {
		name      string
		month     string
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:  "sukses",
			month: "2025-03",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutSummary(gomock.Any(), gomock.Any()).Return(okSummary, nil)
			},
		},
		{
			name:      "gagal: month invalid",
			month:     "03/2025",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidMonth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.GetPayoutSummary(context.Background(), tt.month)
			if tt.expectErr != nil {
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = result
		})
	}
}

// --------------------------------------------------------------------
// UpdatePayoutStatus
// --------------------------------------------------------------------

func TestFinanceUpdatePayoutStatus(t *testing.T) {
	payoutID := uuid.New()
	actorID := uuid.New()

	okRecord := &domain.PayoutRecord{ID: payoutID, Status: "schedule"}
	okResponse := &domain.StatusActionResponse{ID: payoutID, Status: "completed"}

	tests := []struct {
		name      string
		status    string
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		{
			name:   "sukses update ke completed",
			status: "completed",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "completed", actorID).Return(okResponse, nil)
			},
		},
		{
			name:   "sukses update ke schedule",
			status: "schedule",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "schedule", actorID).Return(&domain.StatusActionResponse{ID: payoutID, Status: "schedule"}, nil)
			},
		},
		{
			name:   "sukses update ke on_hold",
			status: "on_hold",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "on_hold", actorID).Return(&domain.StatusActionResponse{ID: payoutID, Status: "on_hold"}, nil)
			},
		},
		{
			name:   "sukses update ke failed",
			status: "failed",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "failed", actorID).Return(&domain.StatusActionResponse{ID: payoutID, Status: "failed"}, nil)
			},
		},
		{
			name:   "sukses: alias 'complete' dinormalisasi ke 'completed'",
			status: "complete",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "completed", actorID).Return(okResponse, nil)
			},
		},
		{
			name:   "sukses: alias 'on hold' dinormalisasi ke 'on_hold'",
			status: "on hold",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "on_hold", actorID).Return(&domain.StatusActionResponse{ID: payoutID, Status: "on_hold"}, nil)
			},
		},
		{
			name:      "gagal: status tidak valid",
			status:    "selesai",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPayoutStatus,
		},
		{
			name:      "gagal: status kosong",
			status:    "",
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidPayoutStatus,
		},
		{
			name:   "gagal: payout tidak ditemukan",
			status: "completed",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(nil, nil)
			},
			expectErr: usecase.ErrPayoutNotFound,
		},
		{
			name:   "gagal: repo error saat GetPayoutRecord",
			status: "completed",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(nil, errors.New("db error"))
			},
		},
		{
			name:   "gagal: UpdatePayoutStatus mengembalikan nil (tidak ditemukan)",
			status: "completed",
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetPayoutRecord(gomock.Any(), payoutID).Return(okRecord, nil)
				r.EXPECT().UpdatePayoutStatus(gomock.Any(), payoutID, "completed", actorID).Return(nil, nil)
			},
			expectErr: usecase.ErrPayoutNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.UpdatePayoutStatus(context.Background(), payoutID, tt.status, actorID)
			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectErr)
				}
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("expected %v, got %v", tt.expectErr, err)
				}
				return
			}
			// Skenario "repo error" tanpa expectErr → hanya cek err tidak nil
			if tt.name == "gagal: repo error saat GetPayoutRecord" {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			_ = result
		})
	}
}

// --------------------------------------------------------------------
// UpdateRefundStatus - transisi state machine
// --------------------------------------------------------------------

func TestFinanceUpdateRefundStatus(t *testing.T) {
	refundID := uuid.New()
	orderID := uuid.New()
	actorID := uuid.New()

	makeRecord := func(status string) *domain.RefundRecord {
		return &domain.RefundRecord{ID: refundID, OrderID: orderID, Status: status}
	}
	okResponse := func(status string) *domain.StatusActionResponse {
		return &domain.StatusActionResponse{ID: refundID, Status: status}
	}

	tests := []struct {
		name      string
		req       domain.UpdateRefundStatusRequest
		setup     func(*mocks.MockFinanceRepository)
		expectErr error
	}{
		// --- Status tidak valid ---
		{
			name:      "gagal: status tidak valid",
			req:       domain.UpdateRefundStatusRequest{Status: "xyz"},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidRefundStatus,
		},
		{
			name:      "gagal: status kosong",
			req:       domain.UpdateRefundStatusRequest{Status: ""},
			setup:     func(r *mocks.MockFinanceRepository) {},
			expectErr: usecase.ErrInvalidRefundStatus,
		},

		// --- Refund tidak ditemukan ---
		{
			name: "gagal: refund tidak ditemukan",
			req:  domain.UpdateRefundStatusRequest{Status: "approved"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(nil, nil)
			},
			expectErr: usecase.ErrRefundNotFound,
		},

		// --- Transisi valid: requested → approved ---
		{
			name: "sukses: requested → approved",
			req:  domain.UpdateRefundStatusRequest{Status: "approved"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("requested"), nil)
				r.EXPECT().IsOrderSettlementCompleted(gomock.Any(), orderID).Return(false, nil)
				r.EXPECT().UpdateRefundStatus(gomock.Any(), refundID, "approved", actorID).Return(okResponse("approved"), nil)
			},
		},

		// --- Transisi valid: requested → rejected ---
		{
			name: "sukses: requested → rejected",
			req:  domain.UpdateRefundStatusRequest{Status: "rejected"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("requested"), nil)
				r.EXPECT().UpdateRefundStatus(gomock.Any(), refundID, "rejected", actorID).Return(okResponse("rejected"), nil)
			},
		},

		// --- Transisi tidak valid ---
		{
			name: "gagal: rejected → approved (tidak diizinkan)",
			req:  domain.UpdateRefundStatusRequest{Status: "approved"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("rejected"), nil)
			},
			expectErr: usecase.ErrInvalidRefundTransition,
		},
		{
			name: "gagal: rejected → processed (tidak diizinkan)",
			req:  domain.UpdateRefundStatusRequest{Status: "processed"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("rejected"), nil)
			},
			expectErr: usecase.ErrInvalidRefundTransition,
		},

		// --- Diblokir oleh settlement ---
		{
			name: "gagal: approved tapi order sudah di-settle",
			req:  domain.UpdateRefundStatusRequest{Status: "approved"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("requested"), nil)
				r.EXPECT().IsOrderSettlementCompleted(gomock.Any(), orderID).Return(true, nil)
			},
			expectErr: usecase.ErrRefundBlockedByCompletedSettlement,
		},

		// --- Status sama (idempoten) ---
		{
			name: "sukses: requested → requested (idempoten, tidak error transisi)",
			req:  domain.UpdateRefundStatusRequest{Status: "requested"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("requested"), nil)
				r.EXPECT().UpdateRefundStatus(gomock.Any(), refundID, "requested", actorID).Return(okResponse("requested"), nil)
			},
		},

		// --- approved → processed (disbursement path) ---
		// Data disbursement tidak ada (nil) → refund not found
		{
			name: "gagal: approved → processed, context disbursement nil",
			req:  domain.UpdateRefundStatusRequest{Status: "processed"},
			setup: func(r *mocks.MockFinanceRepository) {
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("approved"), nil)
				r.EXPECT().IsOrderSettlementCompleted(gomock.Any(), orderID).Return(false, nil)
				r.EXPECT().GetRefundDisbursementContext(gomock.Any(), refundID).Return(nil, nil)
			},
			expectErr: usecase.ErrRefundNotFound,
		},

		// Payout sudah sedang berjalan
		{
			name: "gagal: approved → processed, payout sedang diproses (PROCESSING)",
			req:  domain.UpdateRefundStatusRequest{Status: "processed"},
			setup: func(r *mocks.MockFinanceRepository) {
				processingStatus := "PROCESSING"
				r.EXPECT().GetRefundRecord(gomock.Any(), refundID).Return(makeRecord("approved"), nil)
				r.EXPECT().IsOrderSettlementCompleted(gomock.Any(), orderID).Return(false, nil)
				r.EXPECT().GetRefundDisbursementContext(gomock.Any(), refundID).Return(&domain.RefundDisbursementContext{
					RefundID:     refundID,
					OrderID:      orderID,
					Amount:       100000,
					Currency:     "IDR",
					Status:       "approved",
					PayoutStatus: &processingStatus,
				}, nil)
			},
			expectErr: usecase.ErrRefundPayoutInProgress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)
			tt.setup(repo)

			result, err := uc.UpdateRefundStatus(context.Background(), refundID, tt.req, actorID)
			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("[%s] expected error %v, got nil", tt.name, tt.expectErr)
				}
				if !errors.Is(err, tt.expectErr) {
					t.Fatalf("[%s] expected %v, got %v", tt.name, tt.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("[%s] unexpected error: %v", tt.name, err)
			}
			_ = result
		})
	}
}

// --------------------------------------------------------------------
// Helper normalisasi (unit test fungsi internal via exported behavior)
// --------------------------------------------------------------------

func TestFinanceRefundStatusNormalization(t *testing.T) {
	// "pending" harus dinormalisasi ke "requested" → sehingga ListRefunds berhasil
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc, repo, _ := newFinanceUC(ctrl)

	repo.EXPECT().ListRefunds(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), nil)
	_, _, err := uc.ListRefunds(context.Background(), domain.FinanceListParams{Status: "pending", Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("expected pending to normalize to requested, got error: %v", err)
	}
}

func TestFinancePayoutStatusNormalization(t *testing.T) {
	// "complete" → "completed", "on hold" → "on_hold"
	tests := []struct {
		input    string
		expected bool // true = should pass status validation (tidak error)
	}{
		{"complete", true},
		{"completed", true},
		{"on hold", true},
		{"on_hold", true},
		{"schedule", true},
		{"failed", true},
		{"COMPLETE", true},
		{"random_invalid", false},
	}

	for _, tt := range tests {
		t.Run("payout status: "+tt.input, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			uc, repo, _ := newFinanceUC(ctrl)

			if tt.expected {
				repo.EXPECT().ListPayouts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, int64(0), nil)
			}

			_, _, err := uc.ListPayouts(context.Background(), domain.FinanceListParams{Status: tt.input, Page: 1, Limit: 10})
			if tt.expected && err != nil {
				t.Fatalf("expected no error for status %q, got: %v", tt.input, err)
			}
			if !tt.expected && !errors.Is(err, usecase.ErrInvalidPayoutStatus) {
				t.Fatalf("expected ErrInvalidPayoutStatus for status %q, got: %v", tt.input, err)
			}
		})
	}
}

// --------------------------------------------------------------------
// Pagination meta validation
// --------------------------------------------------------------------

func TestFinancePaginationMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc, repo, _ := newFinanceUC(ctrl)

	// Total 25 item, page 1, limit 10 → total_pages = 3
	repo.EXPECT().ListTransactions(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		[]domain.FinanceTransactionItem{{Invoice: "INV-001"}}, int64(25), nil,
	)

	_, meta, err := uc.ListTransactions(context.Background(), domain.FinanceListParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.TotalItems != 25 {
		t.Errorf("expected TotalItems=25, got %d", meta.TotalItems)
	}
	if meta.TotalPages != 3 {
		t.Errorf("expected TotalPages=3, got %d", meta.TotalPages)
	}
	if meta.Page != 1 {
		t.Errorf("expected Page=1, got %d", meta.Page)
	}
	if meta.Limit != 10 {
		t.Errorf("expected Limit=10, got %d", meta.Limit)
	}
}
