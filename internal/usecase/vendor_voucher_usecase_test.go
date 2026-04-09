package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

func setupVendorVoucherUseCase(t *testing.T) (
	*mocks.MockVendorVoucherRepository,
	*mocks.MockProductRepository,
	domain.VendorVoucherUseCase,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	voucherRepo := mocks.NewMockVendorVoucherRepository(ctrl)
	productRepo := mocks.NewMockProductRepository(ctrl)
	uc := NewVendorVoucherUseCase(voucherRepo, productRepo)
	return voucherRepo, productRepo, uc
}

func TestVendorVoucherGetSummary(t *testing.T) {
	vendorID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	errDB := errors.New("db error")

	tests := []struct {
		name    string
		setup   func(vr *mocks.MockVendorVoucherRepository)
		want    *domain.VendorVoucherSummary
		wantErr bool
	}{
		{
			name: "success",
			setup: func(vr *mocks.MockVendorVoucherRepository) {
				vr.EXPECT().GetSummary(gomock.Any(), vendorID).Return(&domain.VendorVoucherSummary{
					ActiveCount:   24,
					UpcomingCount: 8,
					ExpiredCount:  5,
					TotalUsage:    12,
				}, nil)
			},
			want: &domain.VendorVoucherSummary{
				ActiveCount:   24,
				UpcomingCount: 8,
				ExpiredCount:  5,
				TotalUsage:    12,
			},
		},
		{
			name: "error - repository failure",
			setup: func(vr *mocks.MockVendorVoucherRepository) {
				vr.EXPECT().GetSummary(gomock.Any(), vendorID).Return(nil, errDB)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			voucherRepo, _, uc := setupVendorVoucherUseCase(t)
			tc.setup(voucherRepo)

			got, err := uc.GetSummary(context.Background(), vendorID)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ActiveCount != tc.want.ActiveCount ||
				got.UpcomingCount != tc.want.UpcomingCount ||
				got.ExpiredCount != tc.want.ExpiredCount ||
				got.TotalUsage != tc.want.TotalUsage {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}
