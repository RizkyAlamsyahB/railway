package usecase

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

const (
	defaultWithdrawalFeeEstimateFixed = 0.0
	fixedPayoutFeeActual              = 3000.0
	payoutChannelCurrencyIDR          = "IDR"
	payoutChannelCategoryBank         = "BANK"
)

// VendorWithdrawalPolicy defines fee estimation behavior for vendor withdrawals.
type VendorWithdrawalPolicy struct {
	FeeEstimateFixed float64
	MinNetAmount     float64
}

type vendorUseCase struct {
	otpUseCase                 domain.OTPUseCase
	userRepo                   domain.UserRepository
	vendorRepo                 domain.VendorRepository
	onboardingRepo             domain.VendorOnboardingRepository
	storage                    domain.StorageProvider
	xenditPayout               domain.XenditPayoutProvider
	jwtSecret                  string
	jwtExpiry                  int
	jwtIssuer                  string
	withdrawalFeeEstimateFixed float64
	withdrawalMinNetAmount     float64
}

// NewVendorUseCase creates a new VendorUseCase.
func NewVendorUseCase(
	otpUseCase domain.OTPUseCase,
	userRepo domain.UserRepository,
	vendorRepo domain.VendorRepository,
	onboardingRepo domain.VendorOnboardingRepository,
	storage domain.StorageProvider,
	xenditPayout domain.XenditPayoutProvider,
	jwtSecret string,
	jwtExpiry int,
	jwtIssuer string,
	policy ...VendorWithdrawalPolicy,
) domain.VendorUseCase {
	cfg := VendorWithdrawalPolicy{
		FeeEstimateFixed: defaultWithdrawalFeeEstimateFixed,
		MinNetAmount:     MinWithdrawalAmount,
	}
	if len(policy) > 0 {
		cfg = policy[0]
	}
	if cfg.FeeEstimateFixed < 0 {
		cfg.FeeEstimateFixed = defaultWithdrawalFeeEstimateFixed
	}
	if cfg.MinNetAmount < MinWithdrawalAmount {
		cfg.MinNetAmount = MinWithdrawalAmount
	}

	return &vendorUseCase{
		otpUseCase:                 otpUseCase,
		userRepo:                   userRepo,
		vendorRepo:                 vendorRepo,
		onboardingRepo:             onboardingRepo,
		storage:                    storage,
		xenditPayout:               xenditPayout,
		jwtSecret:                  jwtSecret,
		jwtExpiry:                  jwtExpiry,
		jwtIssuer:                  jwtIssuer,
		withdrawalFeeEstimateFixed: cfg.FeeEstimateFixed,
		withdrawalMinNetAmount:     cfg.MinNetAmount,
	}
}

func (uc *vendorUseCase) Login(ctx context.Context, req domain.VendorLoginRequest) (*domain.VendorLoginResponse, error) {
	// 1. Find user by email.
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrVendorInvalidCredentials
	}

	// 2. Verify password.
	if err := auth.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrVendorInvalidCredentials
	}

	// 3. Check user has "umkm" role.
	if user.Role == nil || user.Role.Code != domain.RoleUMKM {
		return nil, ErrNotVendor
	}

	// 4. Find vendor by owner_user_id.
	vendor, err := uc.vendorRepo.FindByOwnerUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrNoVendorProfile
	}

	// 5. Check vendor status is not "blocked".
	if vendor.Status == domain.VendorStatusBlocked {
		return nil, ErrVendorAccountBlocked
	}

	// 6. Generate JWT with vendor_id in claims.
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role.Code, &vendor.ID, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.VendorLoginResponse{
		AccessToken:  token,
		VendorID:     vendor.ID,
		Email:        user.Email,
		ImageURL:     user.ImageURL,
		StoreName:    vendor.DisplayName,
		VendorType:   vendor.VendorType,
		VendorStatus: vendor.Status,
	}, nil
}

func (uc *vendorUseCase) GetMe(ctx context.Context, vendorID uuid.UUID) (*domain.VendorMeResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	user, err := uc.userRepo.FindByID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &domain.VendorMeResponse{
		VendorID:     vendor.ID,
		ImageURL:     user.ImageURL,
		Email:        user.Email,
		VendorType:   vendor.VendorType,
		VendorStatus: vendor.Status,
		StoreName:    vendor.DisplayName,
	}, nil
}

func (uc *vendorUseCase) GetBalance(ctx context.Context, vendorID uuid.UUID) (*domain.VendorBalanceResponse, error) {
	balance, err := uc.vendorRepo.GetBalance(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	if balance == nil {
		return nil, ErrBalanceNotFound
	}

	return &domain.VendorBalanceResponse{
		AvailableBalance: balance.AvailableBalance,
		PendingBalance:   balance.PendingBalance,
		EscrowBalance:    balance.EscrowBalance,
		TotalEarned:      balance.TotalEarned,
		TotalWithdrawn:   balance.TotalWithdrawn,
	}, nil
}

func (uc *vendorUseCase) ListPayoutChannels(ctx context.Context, vendorID uuid.UUID) (*domain.VendorPayoutChannelsResponse, error) {
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}
	if vendor.Status != domain.VendorStatusActive {
		return nil, ErrVendorNotActive
	}

	channels, err := uc.fetchIDRBankPayoutChannels(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]domain.VendorPayoutChannelItem, 0, len(channels))
	seenCodes := make(map[string]struct{}, len(channels))
	for _, channel := range channels {
		channelCode := strings.ToUpper(strings.TrimSpace(channel.ChannelCode))
		if channelCode == "" || !channel.IsActivated {
			continue
		}
		if _, exists := seenCodes[channelCode]; exists {
			continue
		}
		seenCodes[channelCode] = struct{}{}

		items = append(items, domain.VendorPayoutChannelItem{
			ChannelCode:     channelCode,
			ChannelName:     strings.TrimSpace(channel.ChannelName),
			Currency:        strings.TrimSpace(channel.Currency),
			ChannelCategory: strings.TrimSpace(channel.ChannelCategory),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ChannelCode < items[j].ChannelCode
	})

	return &domain.VendorPayoutChannelsResponse{Channels: items}, nil
}

func (uc *vendorUseCase) RequestWithdrawal(ctx context.Context, vendorID uuid.UUID, req domain.VendorWithdrawRequest) (*domain.VendorWithdrawResponse, error) {
	channelCode := strings.ToUpper(strings.TrimSpace(req.ChannelCode))
	if channelCode == "" {
		return nil, ErrInvalidChannelCode
	}

	// 1. Validate minimum withdrawal amount.
	if req.Amount < MinWithdrawalAmount {
		return nil, ErrBelowMinWithdrawal
	}

	estimatedFee := roundCurrencyAmount(uc.withdrawalFeeEstimateFixed)
	estimatedNetAmount := roundCurrencyAmount(req.Amount - estimatedFee)
	if estimatedNetAmount < uc.withdrawalMinNetAmount {
		return nil, ErrWithdrawalNetAmountTooSmall
	}

	// 2. Find vendor and validate it's active + has Xendit account.
	vendor, err := uc.vendorRepo.FindByID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}
	if vendor.Status != domain.VendorStatusActive {
		return nil, ErrVendorNotActive
	}
	if vendor.XenditAccountID == nil || *vendor.XenditAccountID == "" {
		return nil, ErrVendorNoXenditAccount
	}

	// 3. Load dynamic payout channel list from Xendit and validate channel code.
	payoutChannels, err := uc.fetchIDRBankPayoutChannels(ctx)
	if err != nil {
		return nil, err
	}
	if !isPayoutChannelAllowed(channelCode, payoutChannels) {
		return nil, ErrInvalidChannelCode
	}

	// 4. Check available balance.
	balance, err := uc.vendorRepo.GetBalance(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	if balance == nil {
		return nil, ErrBalanceNotFound
	}
	if balance.AvailableBalance < req.Amount {
		return nil, ErrInsufficientBalance
	}

	// 5. Find vendor bank account.
	bankAccount, err := uc.vendorRepo.FindBankAccountByVendorID(ctx, vendorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bank account: %w", err)
	}
	if bankAccount == nil {
		return nil, ErrVendorBankNotFound
	}

	// 6. Find vendor owner user for receipt notification email.
	owner, err := uc.userRepo.FindByID(ctx, vendor.OwnerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor owner: %w", err)
	}
	if owner == nil {
		return nil, ErrUserNotFound
	}

	// 7. Debit balance (available → pending).
	if err := uc.vendorRepo.DebitBalance(ctx, vendorID, req.Amount); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInsufficientBalance, err)
	}

	// 8. Create withdrawal record.
	now := time.Now()
	withdrawalID := uuid.New()
	description := fmt.Sprintf("Withdrawal for vendor %s", vendorDisplayNameOrFallback(vendor))

	withdrawal := &domain.VendorWithdrawal{
		ID:            withdrawalID,
		VendorID:      vendorID,
		Amount:        req.Amount,
		ChannelCode:   channelCode,
		Status:        domain.WithdrawalStatusPending,
		FeeEstimated:  estimatedFee,
		AmountNet:     estimatedNetAmount,
		TotalDeducted: req.Amount,
		FeeStatus:     domain.WithdrawalFeeStatusPending,
		Description:   &description,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := uc.vendorRepo.CreateWithdrawal(ctx, withdrawal); err != nil {
		// Rollback balance debit on failure.
		_ = uc.vendorRepo.FailWithdrawal(ctx, vendorID, req.Amount)
		return nil, fmt.Errorf("failed to create withdrawal record: %w", err)
	}

	// 9. Build Xendit payout request.
	referenceID := withdrawalID.String()
	xenditReq := domain.XenditPayoutRequest{
		ReferenceID: referenceID,
		ChannelCode: channelCode,
		ChannelProperties: domain.XenditPayoutChannelProperties{
			AccountNumber:     bankAccount.AccountNumber,
			AccountHolderName: bankAccount.AccountHolderName,
		},
		Amount:      estimatedNetAmount,
		Description: description,
		Currency:    "IDR",
		ReceiptNotification: &domain.XenditPayoutReceiptNotification{
			EmailTo: []string{owner.Email},
		},
	}

	// 10. Call Xendit Payout API.
	xenditResp, err := uc.xenditPayout.CreatePayout(ctx, *vendor.XenditAccountID, referenceID, xenditReq)
	if err != nil {
		// Rollback: return funds and mark withdrawal as failed.
		_ = uc.vendorRepo.FailWithdrawal(ctx, vendorID, req.Amount)
		failedReason := err.Error()
		withdrawal.Status = domain.WithdrawalStatusFailed
		withdrawal.FeeStatus = domain.WithdrawalFeeStatusFailed
		withdrawal.FailedReason = &failedReason
		_ = uc.vendorRepo.UpdateWithdrawal(ctx, withdrawal)
		return nil, fmt.Errorf("%w: %v", ErrXenditPayoutFailed, err)
	}

	// 11. Update withdrawal with Xendit response.
	xenditPayoutID := xenditResp.ID
	xenditStatus := xenditResp.Status
	channelCode = xenditResp.ChannelCode

	withdrawal.XenditPayoutID = &xenditPayoutID
	withdrawal.XenditStatus = &xenditStatus
	withdrawal.ChannelCode = channelCode
	withdrawal.Status = domain.WithdrawalStatusProcessing

	if err := uc.vendorRepo.UpdateWithdrawal(ctx, withdrawal); err != nil {
		return nil, fmt.Errorf("failed to update withdrawal: %w", err)
	}

	return &domain.VendorWithdrawResponse{
		WithdrawalID:       withdrawalID,
		XenditPayoutID:     xenditPayoutID,
		Status:             withdrawal.Status,
		XenditStatus:       xenditStatus,
		Amount:             req.Amount,
		RequestedAmount:    req.Amount,
		EstimatedFee:       estimatedFee,
		EstimatedNetAmount: estimatedNetAmount,
		ChannelCode:        channelCode,
	}, nil
}

func (uc *vendorUseCase) fetchIDRBankPayoutChannels(ctx context.Context) ([]domain.XenditPayoutChannel, error) {
	channels, err := uc.xenditPayout.ListPayoutChannels(ctx, domain.XenditListPayoutChannelsParams{
		Currency:        payoutChannelCurrencyIDR,
		ChannelCategory: payoutChannelCategoryBank,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPayoutChannelsUnavailable, err)
	}
	return channels, nil
}

func isPayoutChannelAllowed(channelCode string, channels []domain.XenditPayoutChannel) bool {
	for _, channel := range channels {
		if !channel.IsActivated {
			continue
		}
		if strings.EqualFold(channel.ChannelCode, channelCode) {
			return true
		}
	}
	return false
}

func (uc *vendorUseCase) HandlePayoutWebhook(ctx context.Context, payload domain.XenditPayoutWebhookPayload) error {
	referenceID := strings.TrimSpace(payload.Data.ReferenceID)
	if referenceID == "" {
		return fmt.Errorf("missing payout reference_id")
	}

	withdrawalID, err := uuid.Parse(referenceID)
	if err != nil {
		return fmt.Errorf("invalid payout reference_id: %w", err)
	}

	withdrawal, err := uc.vendorRepo.FindWithdrawalByID(ctx, withdrawalID)
	if err != nil {
		return fmt.Errorf("failed to find withdrawal by id: %w", err)
	}
	if withdrawal == nil {
		// Unknown reference_id should not force retries from Xendit.
		return nil
	}

	xenditStatus := normalizeXenditPayoutWebhookStatus(payload.Event, payload.Data.Status)

	var xenditPayoutID *string
	if id := strings.TrimSpace(payload.Data.ID); id != "" {
		xenditPayoutID = &id
	}

	failedReason := buildPayoutFailureReason(payload)
	newStatus, balanceAction := mapPayoutWebhookTransition(withdrawal.Status, payload.Event, payload.Data.Status)
	if shouldIgnorePayoutWebhookForFinalState(withdrawal.Status, payload.Event, payload.Data.Status) {
		return nil
	}

	if balanceAction == domain.WithdrawalBalanceActionComplete {
		completion, err := buildWithdrawalCompletionDataFromFixedFee(withdrawal, fixedPayoutFeeActual)
		if err != nil {
			return fmt.Errorf("invalid payout completion data: %w", err)
		}

		if err := uc.vendorRepo.ApplyWithdrawalWebhookUpdate(
			ctx,
			withdrawal.ID,
			newStatus,
			xenditStatus,
			xenditPayoutID,
			failedReason,
			balanceAction,
			completion,
		); err != nil {
			return fmt.Errorf("failed to apply payout webhook completion update: %w", err)
		}
		return nil
	}

	// If there is no state/balance mutation and nothing new from Xendit, skip write.
	if newStatus == "" && balanceAction == domain.WithdrawalBalanceActionNone && xenditStatus == "" && xenditPayoutID == nil && failedReason == nil {
		return nil
	}

	if err := uc.vendorRepo.ApplyWithdrawalWebhookUpdate(
		ctx,
		withdrawal.ID,
		newStatus,
		xenditStatus,
		xenditPayoutID,
		failedReason,
		balanceAction,
		nil,
	); err != nil {
		return fmt.Errorf("failed to apply payout webhook update: %w", err)
	}

	return nil
}

func buildWithdrawalCompletionDataFromFixedFee(withdrawal *domain.VendorWithdrawal, fixedFeeActual float64) (*domain.WithdrawalCompletionData, error) {
	if withdrawal == nil {
		return nil, fmt.Errorf("missing withdrawal")
	}

	feeActual := roundCurrencyAmount(fixedFeeActual)
	if feeActual < 0 {
		return nil, fmt.Errorf("negative fee actual")
	}

	netAmount := roundCurrencyAmount(withdrawal.Amount - feeActual)
	if netAmount <= 0 {
		return nil, fmt.Errorf("invalid net amount")
	}

	totalDeducted := roundCurrencyAmount(withdrawal.Amount)
	if totalDeducted <= 0 {
		return nil, fmt.Errorf("invalid total deducted")
	}

	return &domain.WithdrawalCompletionData{
		FeeActual:     feeActual,
		NetAmount:     netAmount,
		TotalDeducted: totalDeducted,
	}, nil
}

func normalizeXenditPayoutWebhookStatus(event, status string) string {
	status = strings.TrimSpace(strings.ToUpper(status))
	if status != "" {
		return status
	}
	return strings.TrimSpace(strings.ToLower(event))
}

func buildPayoutFailureReason(payload domain.XenditPayoutWebhookPayload) *string {
	code := strings.TrimSpace(payload.Data.FailureCode)
	if code == "" {
		return nil
	}
	reason := fmt.Sprintf("xendit failure code: %s", code)
	return &reason
}

func mapPayoutWebhookTransition(currentStatus, event, xenditStatus string) (string, string) {
	event = strings.TrimSpace(strings.ToLower(event))
	xenditStatus = strings.TrimSpace(strings.ToUpper(xenditStatus))

	isSucceeded := event == domain.XenditPayoutWebhookEventSucceeded || xenditStatus == domain.XenditPayoutStatusSucceeded
	isFailed := event == domain.XenditPayoutWebhookEventFailed || xenditStatus == domain.XenditPayoutStatusFailed || xenditStatus == domain.XenditPayoutStatusCancelled
	isReversed := event == domain.XenditPayoutWebhookEventReversed || xenditStatus == domain.XenditPayoutStatusReversed

	switch {
	case isSucceeded:
		switch currentStatus {
		case domain.WithdrawalStatusPending, domain.WithdrawalStatusProcessing:
			return domain.WithdrawalStatusCompleted, domain.WithdrawalBalanceActionComplete
		default:
			return "", domain.WithdrawalBalanceActionNone
		}
	case isReversed:
		switch currentStatus {
		case domain.WithdrawalStatusCompleted:
			return domain.WithdrawalStatusFailed, domain.WithdrawalBalanceActionReverseCompleted
		case domain.WithdrawalStatusPending, domain.WithdrawalStatusProcessing:
			return domain.WithdrawalStatusFailed, domain.WithdrawalBalanceActionFail
		default:
			return "", domain.WithdrawalBalanceActionNone
		}
	case isFailed:
		switch currentStatus {
		case domain.WithdrawalStatusPending, domain.WithdrawalStatusProcessing:
			return domain.WithdrawalStatusFailed, domain.WithdrawalBalanceActionFail
		default:
			return "", domain.WithdrawalBalanceActionNone
		}
	default:
		return "", domain.WithdrawalBalanceActionNone
	}
}

func shouldIgnorePayoutWebhookForFinalState(currentStatus, event, xenditStatus string) bool {
	event = strings.TrimSpace(strings.ToLower(event))
	xenditStatus = strings.TrimSpace(strings.ToUpper(xenditStatus))

	isSucceeded := event == domain.XenditPayoutWebhookEventSucceeded || xenditStatus == domain.XenditPayoutStatusSucceeded
	isFailed := event == domain.XenditPayoutWebhookEventFailed || xenditStatus == domain.XenditPayoutStatusFailed || xenditStatus == domain.XenditPayoutStatusCancelled
	isReversed := event == domain.XenditPayoutWebhookEventReversed || xenditStatus == domain.XenditPayoutStatusReversed

	switch currentStatus {
	case domain.WithdrawalStatusCompleted:
		// Only REVERSED is allowed to move a completed withdrawal.
		return isFailed && !isReversed
	case domain.WithdrawalStatusFailed:
		// A failed withdrawal should not be resurrected by late/duplicate success callbacks.
		return isSucceeded
	default:
		return false
	}
}

func roundCurrencyAmount(amount float64) float64 {
	return math.Round(amount*100) / 100
}
