package usecase

import (
	"errors"
	"time"
)

// Presigned URL expiry durations.
const (
	PresignedUploadExpiry   = 30 * time.Minute
	PresignedDownloadExpiry = 30 * time.Minute
)

// --- Auth errors ---

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotAdmin           = errors.New("user does not have admin access")
	ErrNotCS              = errors.New("user does not have customer service access")
	ErrAccountInactive    = errors.New("account is not active")
)

// --- User errors ---

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrEmailExists              = errors.New("email already exists")
	ErrInvalidBirthDate         = errors.New("invalid birth_date format, expected YYYY-MM-DD")
	ErrCannotDeleteSelf         = errors.New("cannot delete your own account")
	ErrPhoneAlreadyRegistered   = errors.New("phone number already registered")
	ErrInvalidVerificationToken = errors.New("invalid or expired verification token")
	ErrUserAlreadyActive        = errors.New("user is already active")
	ErrResendTooSoon            = errors.New("please wait before requesting another verification email")
	ErrUserNotPending           = errors.New("user is not in pending status")
	ErrUserInvalidCredentials   = errors.New("invalid email or password")
	ErrUserAccountNotActive     = errors.New("account is not active")
	ErrUserEmailNotVerified     = errors.New("email is not verified")
	ErrUserAccountBlocked       = errors.New("account is blocked")
)

// --- OTP errors ---

var (
	ErrOTPInvalid             = errors.New("invalid otp code")
	ErrOTPExpired             = errors.New("otp code has expired")
	ErrOTPTooManyAttempts     = errors.New("too many invalid otp attempts")
	ErrOTPResendTooSoon       = errors.New("please wait before requesting another otp")
	ErrOTPPurposeInvalid      = errors.New("invalid otp purpose")
	ErrOTPProofInvalid        = errors.New("invalid otp proof token")
	ErrOTPProofExpired        = errors.New("otp proof token has expired")
	ErrOTPSecretNotConfigured = errors.New("otp secret is not configured")
)

// --- Vendor errors ---

var (
	ErrEmailAlreadyRegistered    = errors.New("email already registered")
	ErrVendorAlreadyExists       = errors.New("user already has a vendor")
	ErrVendorNotFound            = errors.New("vendor not found")
	ErrVendorOnboardingNotFound  = errors.New("vendor onboarding not found")
	ErrVendorOnboardingInvalid   = errors.New("invalid vendor onboarding token")
	ErrVendorOnboardingExpired   = errors.New("vendor onboarding token has expired")
	ErrVendorOnboardingStep      = errors.New("invalid vendor onboarding step")
	ErrVendorOnboardingDone      = errors.New("vendor onboarding already completed")
	ErrDocumentNotFound          = errors.New("document not found for this vendor")
	ErrInvalidVendorDocumentType = errors.New("invalid vendor document type")
	ErrObjectNotUploaded         = errors.New("object not found in storage")
	ErrInvalidDocumentContent    = errors.New("invalid document content type")
	ErrInvalidDocumentObjectKey  = errors.New("invalid document object key")
	ErrDocumentSizeOverflow      = errors.New("document file size exceeds supported limit")
	ErrVendorInvalidCredentials  = errors.New("invalid email or password")
	ErrVendorAccountBlocked      = errors.New("vendor account is blocked")
	ErrNotVendor                 = errors.New("user does not have vendor access")
	ErrNoVendorProfile           = errors.New("no vendor profile found for this user")
	ErrVendorNotOwned            = errors.New("vendor does not belong to this user")
	ErrVendorResponsibleEmail    = errors.New("responsible person email must match owner account email")
	ErrInvalidStatusTransition   = errors.New("invalid status transition")
)

// --- Product errors ---

var (
	ErrVendorNotActive          = errors.New("vendor is not active")
	ErrVendorNoWarehouseAddress = errors.New("vendor must set a default address before adding products")
	ErrCategoryNotFound         = errors.New("category not found")
	ErrProductNotFound          = errors.New("product not found")
	ErrProductNotOwned          = errors.New("product does not belong to this vendor")
	ErrImageNotFound            = errors.New("image not found for this product")
	ErrImageNotUploaded         = errors.New("image not found in storage")
	ErrInvalidImageContentType  = errors.New("invalid image content type")
	ErrImageSizeOverflow        = errors.New("image file size exceeds supported limit")
	ErrDuplicatePrimaryImage    = errors.New("only one image can be marked as primary")
	ErrNoPrimaryImage           = errors.New("at least one image must be marked as primary")
	ErrTooManyImages            = errors.New("maximum 10 images per product")
)

// --- Xendit errors ---

var (
	ErrXenditAccountCreation = errors.New("failed to create xendit account")
)

// --- Withdrawal errors ---

var (
	ErrInvalidChannelCode          = errors.New("invalid payout channel code")
	ErrXenditPayoutFailed          = errors.New("failed to create payout via Xendit")
	ErrPayoutChannelsUnavailable   = errors.New("payout channels unavailable")
	ErrVendorBankNotFound          = errors.New("vendor bank account not found")
	ErrInsufficientBalance         = errors.New("insufficient available balance")
	ErrBelowMinWithdrawal          = errors.New("amount is below minimum withdrawal of Rp 10.000")
	ErrWithdrawalNetAmountTooSmall = errors.New("net payout amount after fee is below minimum withdrawal of Rp 10.000")
	ErrBalanceNotFound             = errors.New("vendor balance record not found")
)

// MinWithdrawalAmount is the minimum withdrawal amount in IDR.
const MinWithdrawalAmount = 10000

// --- Cart errors ---

var (
	ErrCartItemNotFound    = errors.New("cart item not found")
	ErrCartItemNotOwned    = errors.New("cart item does not belong to your cart")
	ErrVariantNotFound     = errors.New("product variant not found")
	ErrVariantNotActive    = errors.New("product variant is not active")
	ErrProductNotAvailable = errors.New("product is not available")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrNoSelectedCartItems = errors.New("no selected cart items")
)

// --- Checkout errors ---

var (
	ErrCartEmpty                        = errors.New("cart is empty")
	ErrCartHasUnavailableItems          = errors.New("cart contains unavailable items")
	ErrVendorNoXenditAccount            = errors.New("vendor does not have a Xendit account")
	ErrCheckoutStockInsufficient        = errors.New("insufficient stock during checkout")
	ErrInvoiceCreationFailed            = errors.New("failed to create payment invoice")
	ErrCheckoutCompensationFailed       = errors.New("checkout compensation failed")
	ErrOrderNotFound                    = errors.New("order not found")
	ErrOrderNotOwned                    = errors.New("order does not belong to this user")
	ErrInvalidOrderStatus               = errors.New("invalid order status")
	ErrInvoiceNotFound                  = errors.New("payment invoice not found")
	ErrInvalidOrderCompletionTransition = errors.New("order status cannot be completed from current status")
	ErrAddressNoDistrict                = errors.New("selected address has no district; please update your address")
	ErrVendorWarehouseNotFound          = errors.New("vendor warehouse address not found; vendor must set a default address")
	ErrVendorNoCouriersConfigured       = errors.New("vendor has no couriers configured")
	ErrShippingCostFailed               = errors.New("failed to calculate shipping cost")
	ErrInvalidShippingChoice            = errors.New("invalid shipping choice")
	ErrShippingServiceNotFound          = errors.New("selected shipping service not found in available options")
)

// --- Finance errors ---

var (
	ErrInvalidMonth                       = errors.New("invalid month format, expected YYYY-MM")
	ErrInvalidDate                        = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrInvalidDateRange                   = errors.New("invalid date range: date_from must be before or equal to date_to")
	ErrInvalidPaymentStatus               = errors.New("invalid payment status")
	ErrInvalidRefundStatus                = errors.New("invalid refund status")
	ErrInvalidPayoutStatus                = errors.New("invalid payout status")
	ErrInvalidRefundTransition            = errors.New("invalid refund status transition")
	ErrRefundBlockedByCompletedSettlement = errors.New("refund cannot be processed because order settlement is already completed")
	ErrRefundDestinationRequired          = errors.New("destination account data is required for disbursement refund")
	ErrRefundPayoutInProgress             = errors.New("refund payout is already in progress")
	ErrRefundDisbursementFailed           = errors.New("failed to initiate refund disbursement")
	ErrRefundStrategyNotSupported         = errors.New("refund strategy for this payment channel is not supported")
	ErrRefundNotFound                     = errors.New("refund not found")
	ErrRefundAlreadyRequested             = errors.New("refund request already exists for this order")
	ErrOrderRefundNotEligible             = errors.New("order is not eligible for refund request")
	ErrRefundReasonNotFound               = errors.New("return reason not found")
	ErrInvalidRefundEvidenceContentType   = errors.New("invalid refund evidence content type; allowed: image/png, image/jpeg, image/webp, video/mp4")
	ErrRefundEvidenceTooLarge             = errors.New("refund evidence exceeds maximum size of 10 MB")
	ErrRefundEvidenceNotUploaded          = errors.New("refund evidence not found in storage; upload it first")
	ErrTooManyRefundEvidenceImages        = errors.New("maximum 5 images for refund evidence")
	ErrTooManyRefundEvidenceVideos        = errors.New("maximum 1 video for refund evidence")
	ErrPayoutNotFound                     = errors.New("payout batch not found")
)

// --- Notification errors ---

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

// --- Review errors ---

var (
	ErrReviewNotAllowed              = errors.New("order is not eligible for review")
	ErrReviewAlreadyExists           = errors.New("review for this purchased item already exists")
	ErrOrderItemNotFound             = errors.New("order item not found")
	ErrInvalidReviewRating           = errors.New("review rating must be between 1 and 5")
	ErrReviewTextRequired            = errors.New("review text is required")
	ErrTooManyReviewImages           = errors.New("maximum 5 images per review")
	ErrReviewImageNotUploaded        = errors.New("review image not found in storage")
	ErrInvalidReviewImageContentType = errors.New("invalid review image content type")
	ErrReviewImageTooLarge           = errors.New("review image size exceeds maximum limit")
)

// --- CS errors ---

var (
	ErrTicketNotFound           = errors.New("ticket not found")
	ErrConversationNotFound     = errors.New("conversation not found")
	ErrConversationNotAllowed   = errors.New("chat between these roles is not allowed")
	ErrConversationUnauthorized = errors.New("you are not a participant of this conversation")
	ErrReplyTemplateNotFound    = errors.New("reply template not found")
	ErrCSUserNotFound           = errors.New("user not found")
	ErrNoCSAvailable            = errors.New("no customer service agent is available right now")

	// ticket attachment
	ErrInvalidAttachmentContentType = errors.New("invalid attachment content type; allowed: image/png, image/jpeg, video/mp4")
	ErrAttachmentTooLarge           = errors.New("attachment exceeds maximum size of 5 MB")
	ErrAttachmentNotUploaded        = errors.New("attachment not found in storage; upload it first")

	// ticket ownership
	ErrTicketAlreadyTaken     = errors.New("ticket is already assigned to another CS")
	ErrTicketNotAssignedToYou = errors.New("only the assigned CS can perform this action")
	ErrTicketClosed           = errors.New("ticket is already closed")

	// ticket subject
	ErrTicketSubjectNotFound = errors.New("ticket subject not found")

	// reply template
	ErrInvalidReplyTemplateCategory = errors.New("invalid reply template category")
	ErrShortcutAlreadyExists        = errors.New("shortcut already exists")

	// chat attachment
	ErrInvalidChatAttachmentContentType = errors.New("invalid chat attachment content type; allowed: image/png, image/jpeg, video/mp4")
	ErrChatAttachmentTooLarge           = errors.New("chat attachment exceeds maximum size of 10 MB")
	ErrChatAttachmentNotUploaded        = errors.New("chat attachment not found in storage; upload it first")
)

// --- Banner errors ---

var (
	ErrBannerNotFound           = errors.New("banner not found")
	ErrInvalidBannerContentType = errors.New("invalid banner image content type; allowed: image/jpeg, image/png, image/webp")
	ErrBannerImageTooLarge      = errors.New("banner image exceeds maximum size of 5 MB")
	ErrBannerImageNotUploaded   = errors.New("banner image not found in storage; upload it first")
)

// --- Master Lookup errors ---

var (
	ErrAdminCategoryNotFound = errors.New("category not found")
	ErrCategorySlugExists    = errors.New("category with this name already exists")
	ErrReturnReasonNotFound  = errors.New("return reason not found")
	ErrAdminContactNotFound  = errors.New("admin contact not found")
	ErrFAQNotFound           = errors.New("FAQ not found")
	ErrInvalidFAQCategory    = errors.New("invalid FAQ category")
)

// --- Vendor Banner errors ---

var (
	ErrVendorBannerNotFound           = errors.New("vendor banner not found")
	ErrInvalidVendorBannerContentType = errors.New("invalid vendor banner image content type; allowed: image/jpeg, image/png, image/webp")
	ErrVendorBannerImageTooLarge      = errors.New("vendor banner image exceeds maximum size of 5 MB")
	ErrVendorBannerImageNotUploaded   = errors.New("vendor banner image not found in storage; upload it first")
)

// --- Vendor Voucher errors ---

var (
	ErrVoucherNotFound         = errors.New("voucher not found")
	ErrVoucherCodeExists       = errors.New("voucher code already exists")
	ErrInvalidVoucherCode      = errors.New("invalid voucher code")
	ErrInvalidVoucherPeriod    = errors.New("invalid voucher period: ends_at must be after or equal to starts_at")
	ErrInvalidVoucherQuota     = errors.New("invalid voucher quota")
	ErrVoucherProductsRequired = errors.New("voucher must be assigned to at least one product")
)

// --- Address errors ---

var (
	ErrAddressNotFound         = errors.New("address not found")
	ErrAddressNotOwned         = errors.New("address does not belong to this user")
	ErrCannotDeleteDefaultAddr = errors.New("cannot delete default address; set another address as default first")
	ErrCannotDeleteLastAddr    = errors.New("cannot delete your last address; at least one address is required")
)

// --- Courier errors ---

var (
	ErrCourierNotFound       = errors.New("courier not found")
	ErrCourierNotActive      = errors.New("courier is not active")
	ErrVendorCourierNotFound = errors.New("vendor courier selection not found")
)

// --- Shipping / RajaOngkir errors ---

var (
	ErrRajaOngkirFailed = errors.New("failed to fetch data from shipping provider")
)

// --- Vendor Order Management errors ---

var (
	ErrOrderNotBelongToVendor       = errors.New("order does not belong to this vendor")
	ErrInvalidOrderAcceptTransition = errors.New("order can only be accepted when status is paid")
	ErrInvalidOrderRejectTransition = errors.New("order can only be rejected when status is paid or processing")
	ErrInvalidOrderShipTransition   = errors.New("order can only be shipped when status is processing")
	ErrTrackingNumberRequired       = errors.New("tracking number is required")
	ErrShipmentNotFound             = errors.New("shipment not found for this order")
	ErrTrackingFailed               = errors.New("failed to track waybill")
	ErrInvalidOrderCancelTransition = errors.New("order can only be canceled when status is pending_payment, paid, or processing")
)

// --- User Bank Account errors ---

var (
	ErrBankAccountNotFound = errors.New("bank account not found")
	ErrBankAccountNotOwned = errors.New("bank account does not belong to this user")
)
