package usecase

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/utils/auth"
)

// Required document types for vendor registration (must be uploaded for draft→submitted).
var requiredDocTypes = []string{
	"owner_ktp",
	"owner_passport",
	"store_photo",
	"bank_account_proof",
	"business_logo",
	"business_banner",
}

// Optional document types (presigned URLs generated, but not required for submission).
var optionalDocTypes = []string{
	"business_npwp",
}

// allDocTypes combines required + optional for presigned URL generation.
var allDocTypes = append(append([]string{}, requiredDocTypes...), optionalDocTypes...)

type vendorUseCase struct {
	userRepo     domain.UserRepository
	vendorRepo   domain.VendorRepository
	storage      domain.StorageProvider
	xenditPayout domain.XenditPayoutProvider
	jwtSecret    string
	jwtExpiry    int
	jwtIssuer    string
}

// NewVendorUseCase creates a new VendorUseCase.
func NewVendorUseCase(
	userRepo domain.UserRepository,
	vendorRepo domain.VendorRepository,
	storage domain.StorageProvider,
	xenditPayout domain.XenditPayoutProvider,
	jwtSecret string,
	jwtExpiry int,
	jwtIssuer string,
) domain.VendorUseCase {
	return &vendorUseCase{
		userRepo:     userRepo,
		vendorRepo:   vendorRepo,
		storage:      storage,
		xenditPayout: xenditPayout,
		jwtSecret:    jwtSecret,
		jwtExpiry:    jwtExpiry,
		jwtIssuer:    jwtIssuer,
	}
}

func (uc *vendorUseCase) Register(ctx context.Context, req domain.VendorRegisterRequest) (*domain.VendorRegisterResponse, error) {
	// 1. Check if email is already taken. If user exists, check for existing vendor.
	existing, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		existingVendor, err := uc.vendorRepo.FindByOwnerUserID(ctx, existing.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing vendor: %w", err)
		}
		if existingVendor != nil {
			return nil, ErrVendorAlreadyExists
		}
		return nil, ErrEmailAlreadyRegistered
	}

	// 2. Hash password.
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Prepare all entities in memory.
	now := time.Now()
	phone := req.Phone
	userID := uuid.New()
	vendorID := uuid.New()

	user := &domain.User{
		ID:           userID,
		Email:        req.Email,
		FullName:     req.OwnerName,
		Phone:        &phone,
		PasswordHash: hash,
		Status:       domain.UserStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	vendor := &domain.Vendor{
		ID:                    vendorID,
		OwnerUserID:           userID,
		VendorType:            req.StoreType,
		LegalName:             req.LegalName,
		DisplayName:           req.StoreName,
		ResponsiblePersonName: req.ResponsiblePersonName,
		Status:                domain.VendorStatusDraft,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	bankAccount := &domain.VendorBankAccount{
		ID:                 uuid.New(),
		VendorID:           vendorID,
		BankName:           req.BankName,
		AccountNumber:      req.BankAccountNumber,
		AccountHolderName:  req.BankAccountHolderName,
		VerificationStatus: domain.VerificationStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// 4. Generate presigned upload URLs BEFORE any DB writes.
	//    Also prepare VendorDocument records for each document type.
	uploadURLs := make([]domain.PresignedUploadInfo, 0, len(allDocTypes))
	documents := make([]domain.VendorDocument, 0, len(allDocTypes))

	for _, docType := range allDocTypes {
		objectKey := fmt.Sprintf("vendors/%s/documents/%s/%s", vendorID.String(), docType, uuid.New().String())

		// Keep content type unsigned so client can upload with the file's actual MIME type.
		uploadURL, err := uc.storage.GeneratePresignedUploadURL(ctx, objectKey, "", PresignedUploadExpiry)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for %s: %w", docType, err)
		}

		uploadURLs = append(uploadURLs, domain.PresignedUploadInfo{
			DocType:   docType,
			UploadURL: uploadURL,
			ObjectKey: objectKey,
		})

		documents = append(documents, domain.VendorDocument{
			ID:                 uuid.New(),
			VendorID:           vendorID,
			DocType:            docType,
			FileURL:            objectKey,
			VerificationStatus: domain.VerificationStatusPending,
			CreatedAt:          now,
			UpdatedAt:          now,
		})
	}

	// 5. Persist user to DB.
	if err := uc.userRepo.Create(ctx, user, domain.RoleUMKM); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 6. Persist vendor + bank account + documents in one transaction.
	if err := uc.vendorRepo.Create(ctx, vendor, bankAccount, documents); err != nil {
		return nil, fmt.Errorf("failed to create vendor: %w", err)
	}

	// 7. Initialize vendor balance (zero balance).
	if err := uc.vendorRepo.InitBalance(ctx, vendorID); err != nil {
		return nil, fmt.Errorf("failed to initialize vendor balance: %w", err)
	}

	// 7. Generate JWT token.
	token, err := auth.GenerateToken(userID, user.Email, domain.RoleUMKM, &vendorID, uc.jwtSecret, uc.jwtExpiry, uc.jwtIssuer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &domain.VendorRegisterResponse{
		VendorID:   vendorID,
		UserID:     userID,
		Token:      token,
		UploadURLs: uploadURLs,
	}, nil
}

func (uc *vendorUseCase) ConfirmDocuments(ctx context.Context, userID uuid.UUID, req domain.ConfirmDocumentsRequest) (*domain.ConfirmDocumentsResponse, error) {
	// 1. Find vendor by owner_user_id.
	vendor, err := uc.vendorRepo.FindByOwnerUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find vendor: %w", err)
	}
	if vendor == nil {
		return nil, ErrVendorNotFound
	}

	// 2. Load existing documents for this vendor.
	existingDocs, err := uc.vendorRepo.FindDocumentsByVendorID(ctx, vendor.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	// Build lookup map by doc_type.
	docMap := make(map[string]*domain.VendorDocument, len(existingDocs))
	for i := range existingDocs {
		docMap[existingDocs[i].DocType] = &existingDocs[i]
	}

	// 3. Validate each requested document and verify it exists in S3.
	now := time.Now()
	updatedDocs := make([]domain.VendorDocument, 0, len(req.Documents))

	for _, item := range req.Documents {
		doc, exists := docMap[item.DocType]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrDocumentNotFound, item.DocType)
		}

		// Verify object actually exists in S3 via HeadObject.
		info, err := uc.storage.HeadObject(ctx, item.ObjectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to verify object %s: %w", item.ObjectKey, err)
		}
		if info == nil {
			return nil, fmt.Errorf("%w: %s", ErrObjectNotUploaded, item.ObjectKey)
		}
		contentType := normalizeContentType(info.ContentType)
		if !isAllowedDocumentContentType(contentType) {
			return nil, fmt.Errorf("%w: %s (%s)", ErrInvalidDocumentContent, item.DocType, info.ContentType)
		}
		if info.ContentLength > int64(math.MaxInt32) {
			return nil, fmt.Errorf("%w: %s (%d bytes)", ErrDocumentSizeOverflow, item.DocType, info.ContentLength)
		}

		// Prepare update.
		fileSize := int(info.ContentLength)
		doc.FileURL = item.ObjectKey
		doc.MimeType = &contentType
		doc.FileSizeBytes = &fileSize
		doc.UploadedBy = &userID
		doc.UpdatedAt = now

		updatedDocs = append(updatedDocs, *doc)
	}

	// 4. Determine whether all required documents now have uploads.
	for i := range updatedDocs {
		docMap[updatedDocs[i].DocType] = &updatedDocs[i]
	}

	allUploaded := true
	for _, docType := range requiredDocTypes {
		doc, exists := docMap[docType]
		if !exists || doc.UploadedBy == nil {
			allUploaded = false
			break
		}
	}

	// 5. Determine new status.
	newStatus := ""
	if allUploaded && vendor.Status == domain.VendorStatusDraft {
		newStatus = domain.VendorStatusSubmitted
	}

	// 6. Persist all document updates + optional status change in one transaction.
	if err := uc.vendorRepo.ConfirmDocumentsAndUpdateStatus(ctx, vendor.ID, updatedDocs, newStatus); err != nil {
		return nil, fmt.Errorf("failed to confirm documents: %w", err)
	}

	vendorStatus := vendor.Status
	if newStatus != "" {
		vendorStatus = newStatus
	}

	return &domain.ConfirmDocumentsResponse{
		VendorID:       vendor.ID,
		VendorStatus:   vendorStatus,
		DocumentsCount: len(req.Documents),
	}, nil
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
		Token:        token,
		VendorID:     vendor.ID,
		VendorStatus: vendor.Status,
		DisplayName:  vendor.DisplayName,
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
		TotalEarned:      balance.TotalEarned,
		TotalWithdrawn:   balance.TotalWithdrawn,
	}, nil
}

func (uc *vendorUseCase) RequestWithdrawal(ctx context.Context, vendorID uuid.UUID, req domain.VendorWithdrawRequest) (*domain.VendorWithdrawResponse, error) {
	// 1. Validate channel code.
	if !domain.ValidPayoutChannelCodes[req.ChannelCode] {
		return nil, ErrInvalidChannelCode
	}

	// 2. Validate minimum withdrawal amount.
	if req.Amount < MinWithdrawalAmount {
		return nil, ErrBelowMinWithdrawal
	}

	// 3. Find vendor and validate it's active + has Xendit account.
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
	description := fmt.Sprintf("Withdrawal for vendor %s", vendor.DisplayName)

	withdrawal := &domain.VendorWithdrawal{
		ID:          withdrawalID,
		VendorID:    vendorID,
		Amount:      req.Amount,
		ChannelCode: req.ChannelCode,
		Status:      domain.WithdrawalStatusPending,
		Description: &description,
		CreatedAt:   now,
		UpdatedAt:   now,
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
		ChannelCode: req.ChannelCode,
		ChannelProperties: domain.XenditPayoutChannelProperties{
			AccountNumber:     bankAccount.AccountNumber,
			AccountHolderName: bankAccount.AccountHolderName,
		},
		Amount:      req.Amount,
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
		withdrawal.FailedReason = &failedReason
		_ = uc.vendorRepo.UpdateWithdrawal(ctx, withdrawal)
		return nil, fmt.Errorf("%w: %v", ErrXenditPayoutFailed, err)
	}

	// 11. Update withdrawal with Xendit response.
	xenditPayoutID := xenditResp.ID
	xenditStatus := xenditResp.Status
	channelCode := xenditResp.ChannelCode

	withdrawal.XenditPayoutID = &xenditPayoutID
	withdrawal.XenditStatus = &xenditStatus
	withdrawal.ChannelCode = channelCode
	withdrawal.Status = domain.WithdrawalStatusProcessing

	if err := uc.vendorRepo.UpdateWithdrawal(ctx, withdrawal); err != nil {
		return nil, fmt.Errorf("failed to update withdrawal: %w", err)
	}

	return &domain.VendorWithdrawResponse{
		WithdrawalID:   withdrawalID,
		XenditPayoutID: xenditPayoutID,
		Status:         withdrawal.Status,
		XenditStatus:   xenditStatus,
		Amount:         req.Amount,
		ChannelCode:    channelCode,
	}, nil
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
	); err != nil {
		return fmt.Errorf("failed to apply payout webhook update: %w", err)
	}

	return nil
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
