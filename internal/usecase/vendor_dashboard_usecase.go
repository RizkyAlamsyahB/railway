package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type vendorDashboardUseCase struct {
	vendorOrderRepo  domain.VendorOrderRepository
	vendorRepo       domain.VendorRepository
	notificationRepo domain.NotificationRepository
}

// NewVendorDashboardUseCase creates a new VendorDashboardUseCase.
func NewVendorDashboardUseCase(
	vendorOrderRepo domain.VendorOrderRepository,
	vendorRepo domain.VendorRepository,
	notificationRepo domain.NotificationRepository,
) domain.VendorDashboardUseCase {
	return &vendorDashboardUseCase{
		vendorOrderRepo:  vendorOrderRepo,
		vendorRepo:       vendorRepo,
		notificationRepo: notificationRepo,
	}
}

func (uc *vendorDashboardUseCase) GetDashboard(ctx context.Context, vendorID uuid.UUID, month, year int) (*domain.VendorDashboardResponse, error) {
	// 1. Verify vendor exists.
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	// 2. Calculate period (start/end of month).
	loc := time.FixedZone("WIB", 7*3600)
	periodStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	periodEnd := periodStart.AddDate(0, 1, 0)

	// 3. Get order stats for the period.
	total, successful, err := uc.vendorOrderRepo.DashboardOrderStats(ctx, vendorID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	successPct := float64(0)
	if total > 0 {
		successPct = math.Round(float64(successful)/float64(total)*10000) / 100
	}

	// 4. Get revenue (available_balance) and pending settlement from vendor balance.
	revenue := float64(0)
	pendingSettlement := float64(0)
	balance, err := uc.vendorRepo.GetBalance(ctx, vendorID)
	if err == nil && balance != nil {
		revenue = balance.AvailableBalance
		pendingSettlement = balance.EscrowBalance + balance.PendingBalance
	}

	// 5. Get today's transactions (max 5).
	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayTx, err := uc.vendorOrderRepo.DashboardTodayTransactions(ctx, vendorID, todayStart, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get today transactions: %w", err)
	}

	// 6. Get recent notifications for the vendor owner.
	var notifications []domain.NotificationItem
	notifs, _, err := uc.notificationRepo.List(ctx, vendor.OwnerUserID, domain.NotificationListParams{Page: 1, Limit: 5})
	if err == nil {
		notifications = notifs
	}
	if notifications == nil {
		notifications = []domain.NotificationItem{}
	}

	// 7. Payment flow distribution (last 30 days).
	flowEnd := now
	flowStart := now.AddDate(0, 0, -30)
	paymentFlow, err := uc.vendorOrderRepo.DashboardPaymentFlow(ctx, vendorID, flowStart, flowEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment flow: %w", err)
	}
	if paymentFlow == nil {
		paymentFlow = []domain.VendorPaymentFlowItem{}
	}

	return &domain.VendorDashboardResponse{
		Month: fmt.Sprintf("%04d-%02d", year, month),
		Year:  year,
		Summary: domain.VendorDashboardSummary{
			TotalTransactions:            total,
			SuccessfulTransactionPercent: successPct,
			Revenue:                      revenue,
			PendingSettlement:            pendingSettlement,
		},
		TodayTransactions:   todayTx,
		RecentNotifications: notifications,
		PaymentFlow:         paymentFlow,
	}, nil
}
