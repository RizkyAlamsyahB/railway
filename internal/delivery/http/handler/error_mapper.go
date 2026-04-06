package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/usecase"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

// errorRule maps a sentinel error to an HTTP status code and optional fixed message.
// When message is empty, err.Error() is used (for wrapped errors with contextual details).
type errorRule struct {
	sentinel error
	status   int
	message  string
}

// errorRules defines the complete mapping from usecase sentinel errors to HTTP responses.
// Order matters: the first matching rule wins.
var errorRules = []errorRule{
	// Auth
	{usecase.ErrInvalidCredentials, http.StatusUnauthorized, "invalid email or password"},
	{usecase.ErrAccountInactive, http.StatusUnauthorized, "account is not active"},
	{usecase.ErrNotAdmin, http.StatusForbidden, "admin access required"},
	{usecase.ErrNotCS, http.StatusForbidden, "customer service access required"},

	// User
	{usecase.ErrUserNotFound, http.StatusNotFound, "user not found"},
	{usecase.ErrEmailExists, http.StatusBadRequest, "email already exists"},
	{usecase.ErrInvalidBirthDate, http.StatusBadRequest, "invalid birth_date format, expected YYYY-MM-DD"},
	{usecase.ErrCannotDeleteSelf, http.StatusBadRequest, "cannot delete your own account"},
	{usecase.ErrEmailAlreadyRegistered, http.StatusConflict, ""},
	{usecase.ErrPhoneAlreadyRegistered, http.StatusConflict, ""},
	{usecase.ErrInvalidVerificationToken, http.StatusBadRequest, ""},
	{usecase.ErrUserAlreadyActive, http.StatusConflict, ""},
	{usecase.ErrResendTooSoon, http.StatusTooManyRequests, ""},
	{usecase.ErrUserNotPending, http.StatusBadRequest, ""},
	{usecase.ErrUserInvalidCredentials, http.StatusUnauthorized, "invalid email or password"},
	{usecase.ErrUserAccountBlocked, http.StatusForbidden, "account is blocked"},
	{usecase.ErrUserAccountNotActive, http.StatusUnauthorized, "account is not active"},
	{usecase.ErrUserEmailNotVerified, http.StatusUnauthorized, "email is not verified"},

	// OTP
	{usecase.ErrOTPInvalid, http.StatusBadRequest, ""},
	{usecase.ErrOTPExpired, http.StatusBadRequest, ""},
	{usecase.ErrOTPTooManyAttempts, http.StatusTooManyRequests, ""},
	{usecase.ErrOTPResendTooSoon, http.StatusTooManyRequests, ""},
	{usecase.ErrOTPPurposeInvalid, http.StatusBadRequest, ""},
	{usecase.ErrOTPProofInvalid, http.StatusUnauthorized, ""},
	{usecase.ErrOTPProofExpired, http.StatusUnauthorized, ""},
	{usecase.ErrOTPSecretNotConfigured, http.StatusInternalServerError, ""},

	// Vendor
	{usecase.ErrVendorAlreadyExists, http.StatusConflict, "user already has a vendor"},
	{usecase.ErrVendorNotFound, http.StatusNotFound, "vendor not found"},
	{usecase.ErrVendorOnboardingNotFound, http.StatusNotFound, "vendor onboarding not found"},
	{usecase.ErrVendorOnboardingInvalid, http.StatusUnauthorized, ""},
	{usecase.ErrVendorOnboardingExpired, http.StatusUnauthorized, ""},
	{usecase.ErrVendorOnboardingStep, http.StatusBadRequest, ""},
	{usecase.ErrVendorOnboardingDone, http.StatusConflict, ""},
	{usecase.ErrDocumentNotFound, http.StatusBadRequest, ""},
	{usecase.ErrInvalidVendorDocumentType, http.StatusBadRequest, ""},
	{usecase.ErrObjectNotUploaded, http.StatusBadRequest, ""},
	{usecase.ErrInvalidDocumentContent, http.StatusBadRequest, ""},
	{usecase.ErrInvalidDocumentObjectKey, http.StatusBadRequest, ""},
	{usecase.ErrDocumentSizeOverflow, http.StatusBadRequest, ""},
	{usecase.ErrVendorInvalidCredentials, http.StatusUnauthorized, "invalid email or password"},
	{usecase.ErrVendorAccountBlocked, http.StatusForbidden, "vendor account is blocked"},
	{usecase.ErrNotVendor, http.StatusForbidden, "vendor access required"},
	{usecase.ErrNoVendorProfile, http.StatusForbidden, "no vendor profile found"},
	{usecase.ErrVendorNotOwned, http.StatusForbidden, "vendor does not belong to this user"},
	{usecase.ErrInvalidStatusTransition, http.StatusBadRequest, ""},

	// Xendit
	{usecase.ErrXenditAccountCreation, http.StatusBadGateway, ""},

	// Withdrawal
	{usecase.ErrInvalidChannelCode, http.StatusBadRequest, ""},
	{usecase.ErrXenditPayoutFailed, http.StatusBadGateway, ""},
	{usecase.ErrPayoutChannelsUnavailable, http.StatusBadGateway, ""},
	{usecase.ErrVendorBankNotFound, http.StatusNotFound, "vendor bank account not found"},
	{usecase.ErrInsufficientBalance, http.StatusBadRequest, "insufficient available balance"},
	{usecase.ErrBelowMinWithdrawal, http.StatusBadRequest, ""},
	{usecase.ErrWithdrawalNetAmountTooSmall, http.StatusBadRequest, ""},
	{usecase.ErrBalanceNotFound, http.StatusNotFound, "vendor balance not found"},

	// Product
	{usecase.ErrVendorNotActive, http.StatusForbidden, "vendor is not active"},
	{usecase.ErrVendorNoWarehouseAddress, http.StatusUnprocessableEntity, ""},
	{usecase.ErrCategoryNotFound, http.StatusBadRequest, "category not found"},
	{usecase.ErrProductNotFound, http.StatusNotFound, "product not found"},
	{usecase.ErrProductNotOwned, http.StatusForbidden, "product does not belong to this vendor"},
	{usecase.ErrImageNotFound, http.StatusBadRequest, ""},
	{usecase.ErrImageNotUploaded, http.StatusBadRequest, ""},
	{usecase.ErrInvalidImageContentType, http.StatusBadRequest, ""},
	{usecase.ErrImageSizeOverflow, http.StatusBadRequest, ""},
	{usecase.ErrDuplicatePrimaryImage, http.StatusBadRequest, "only one image can be marked as primary"},
	{usecase.ErrNoPrimaryImage, http.StatusBadRequest, "at least one image must be marked as primary"},
	{usecase.ErrTooManyImages, http.StatusBadRequest, "maximum 10 images per product"},

	// Cart
	{usecase.ErrCartItemNotFound, http.StatusNotFound, "cart item not found"},
	{usecase.ErrCartItemNotOwned, http.StatusForbidden, "cart item does not belong to your cart"},
	{usecase.ErrVariantNotFound, http.StatusNotFound, "product variant not found"},
	{usecase.ErrVariantNotActive, http.StatusBadRequest, "product variant is not active"},
	{usecase.ErrProductNotAvailable, http.StatusBadRequest, "product is not available"},
	{usecase.ErrInsufficientStock, http.StatusConflict, "insufficient stock for requested quantity"},
	{usecase.ErrNoSelectedCartItems, http.StatusBadRequest, "no selected cart items"},

	// Checkout
	{usecase.ErrCartEmpty, http.StatusBadRequest, "cart is empty"},
	{usecase.ErrCartHasUnavailableItems, http.StatusConflict, "cart contains unavailable or inactive items"},
	{usecase.ErrVendorNoXenditAccount, http.StatusBadGateway, ""},
	{usecase.ErrCheckoutStockInsufficient, http.StatusConflict, ""},
	{usecase.ErrInvoiceCreationFailed, http.StatusBadGateway, ""},
	{usecase.ErrCheckoutCompensationFailed, http.StatusBadGateway, ""},
	{usecase.ErrOrderNotFound, http.StatusNotFound, "order not found"},
	{usecase.ErrOrderNotOwned, http.StatusForbidden, "order does not belong to this user"},
	{usecase.ErrInvalidOrderStatus, http.StatusBadRequest, "invalid order status"},
	{usecase.ErrInvoiceNotFound, http.StatusNotFound, "payment invoice not found"},
	{usecase.ErrInvalidOrderCompletionTransition, http.StatusBadRequest, "order cannot be completed from current status"},
	{usecase.ErrAddressNoDistrict, http.StatusUnprocessableEntity, ""},
	{usecase.ErrVendorWarehouseNotFound, http.StatusUnprocessableEntity, ""},
	{usecase.ErrVendorNoCouriersConfigured, http.StatusUnprocessableEntity, ""},
	{usecase.ErrShippingCostFailed, http.StatusBadGateway, ""},
	{usecase.ErrInvalidShippingChoice, http.StatusBadRequest, ""},
	{usecase.ErrShippingServiceNotFound, http.StatusBadRequest, ""},

	// Finance
	{usecase.ErrInvalidMonth, http.StatusBadRequest, "invalid month format, expected YYYY-MM"},
	{usecase.ErrInvalidDate, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD"},
	{usecase.ErrInvalidDateRange, http.StatusBadRequest, "invalid date range: date_from must be before or equal to date_to"},
	{usecase.ErrInvalidPaymentStatus, http.StatusBadRequest, "invalid payment status"},
	{usecase.ErrInvalidRefundStatus, http.StatusBadRequest, "invalid refund status"},
	{usecase.ErrInvalidPayoutStatus, http.StatusBadRequest, "invalid payout status"},
	{usecase.ErrInvalidRefundTransition, http.StatusBadRequest, "invalid refund status transition"},
	{usecase.ErrRefundBlockedByCompletedSettlement, http.StatusConflict, "refund cannot be processed because order settlement is already completed"},
	{usecase.ErrRefundDestinationRequired, http.StatusBadRequest, "destination account data is required for disbursement refund"},
	{usecase.ErrRefundPayoutInProgress, http.StatusConflict, "refund payout is already in progress"},
	{usecase.ErrRefundDisbursementFailed, http.StatusBadGateway, "failed to initiate refund disbursement"},
	{usecase.ErrRefundStrategyNotSupported, http.StatusBadRequest, "refund strategy for this payment channel is not supported"},
	{usecase.ErrRefundNotFound, http.StatusNotFound, "refund not found"},
	{usecase.ErrRefundAlreadyRequested, http.StatusConflict, "refund request already exists for this order"},
	{usecase.ErrOrderRefundNotEligible, http.StatusUnprocessableEntity, "order is not eligible for refund request"},
	{usecase.ErrRefundReasonNotFound, http.StatusBadRequest, "return reason not found"},
	{usecase.ErrInvalidRefundEvidenceContentType, http.StatusUnprocessableEntity, "invalid refund evidence content type; allowed: image/png, image/jpeg, image/webp, video/mp4"},
	{usecase.ErrRefundEvidenceTooLarge, http.StatusUnprocessableEntity, "refund evidence exceeds maximum size of 10 MB"},
	{usecase.ErrRefundEvidenceNotUploaded, http.StatusUnprocessableEntity, "refund evidence not found in storage; upload it first"},
	{usecase.ErrTooManyRefundEvidenceImages, http.StatusBadRequest, "maximum 5 images for refund evidence"},
	{usecase.ErrTooManyRefundEvidenceVideos, http.StatusBadRequest, "maximum 1 video for refund evidence"},
	{usecase.ErrPayoutNotFound, http.StatusNotFound, "payout batch not found"},

	// Notifications
	{usecase.ErrNotificationNotFound, http.StatusNotFound, "notification not found"},

	// Review
	{usecase.ErrReviewNotAllowed, http.StatusUnprocessableEntity, "order is not eligible for review"},
	{usecase.ErrReviewAlreadyExists, http.StatusConflict, "review for this purchased item already exists"},
	{usecase.ErrOrderItemNotFound, http.StatusNotFound, "order item not found"},
	{usecase.ErrInvalidReviewRating, http.StatusBadRequest, "review rating must be between 1 and 5"},
	{usecase.ErrReviewTextRequired, http.StatusBadRequest, "review text is required"},
	{usecase.ErrTooManyReviewImages, http.StatusBadRequest, "maximum 5 images per review"},
	{usecase.ErrReviewImageNotUploaded, http.StatusUnprocessableEntity, "review image not found in storage"},
	{usecase.ErrInvalidReviewImageContentType, http.StatusUnprocessableEntity, "invalid review image content type"},
	{usecase.ErrReviewImageTooLarge, http.StatusUnprocessableEntity, "review image size exceeds maximum limit"},

	// CS
	{usecase.ErrTicketNotFound, http.StatusNotFound, "ticket not found"},
	{usecase.ErrConversationNotFound, http.StatusNotFound, "conversation not found"},
	{usecase.ErrConversationNotAllowed, http.StatusForbidden, "chat between these roles is not allowed"},
	{usecase.ErrConversationUnauthorized, http.StatusForbidden, "you are not a participant of this conversation"},
	{usecase.ErrReplyTemplateNotFound, http.StatusNotFound, "reply template not found"},
	{usecase.ErrCSUserNotFound, http.StatusNotFound, "user not found"},
	{usecase.ErrNoCSAvailable, http.StatusServiceUnavailable, "no customer service agent is available right now"},

	// ticket attachment
	{usecase.ErrInvalidAttachmentContentType, http.StatusUnprocessableEntity, ""},
	{usecase.ErrAttachmentTooLarge, http.StatusUnprocessableEntity, ""},
	{usecase.ErrAttachmentNotUploaded, http.StatusUnprocessableEntity, ""},

	// ticket ownership
	{usecase.ErrTicketAlreadyTaken, http.StatusConflict, "ticket is already assigned to another CS"},
	{usecase.ErrTicketNotAssignedToYou, http.StatusForbidden, "only the assigned CS can perform this action"},
	{usecase.ErrTicketClosed, http.StatusUnprocessableEntity, "ticket is already closed"},

	// ticket subject
	{usecase.ErrTicketSubjectNotFound, http.StatusBadRequest, "ticket subject not found"},

	// reply template
	{usecase.ErrInvalidReplyTemplateCategory, http.StatusUnprocessableEntity, ""},
	{usecase.ErrShortcutAlreadyExists, http.StatusConflict, ""},

	// chat attachment
	{usecase.ErrInvalidChatAttachmentContentType, http.StatusUnprocessableEntity, ""},
	{usecase.ErrChatAttachmentTooLarge, http.StatusUnprocessableEntity, ""},
	{usecase.ErrChatAttachmentNotUploaded, http.StatusUnprocessableEntity, ""},

	// Banner
	{usecase.ErrBannerNotFound, http.StatusNotFound, "banner not found"},
	{usecase.ErrInvalidBannerContentType, http.StatusUnprocessableEntity, ""},
	{usecase.ErrBannerImageTooLarge, http.StatusUnprocessableEntity, ""},
	{usecase.ErrBannerImageNotUploaded, http.StatusUnprocessableEntity, ""},

	// Master Lookup
	{usecase.ErrAdminCategoryNotFound, http.StatusNotFound, "category not found"},
	{usecase.ErrCategorySlugExists, http.StatusConflict, ""},
	{usecase.ErrReturnReasonNotFound, http.StatusNotFound, "return reason not found"},
	{usecase.ErrAdminContactNotFound, http.StatusNotFound, "admin contact not found"},
	{usecase.ErrFAQNotFound, http.StatusNotFound, "FAQ not found"},
	{usecase.ErrInvalidFAQCategory, http.StatusUnprocessableEntity, ""},

	// Vendor Banner
	{usecase.ErrVendorBannerNotFound, http.StatusNotFound, "vendor banner not found"},
	{usecase.ErrInvalidVendorBannerContentType, http.StatusUnprocessableEntity, ""},
	{usecase.ErrVendorBannerImageTooLarge, http.StatusUnprocessableEntity, ""},
	{usecase.ErrVendorBannerImageNotUploaded, http.StatusUnprocessableEntity, ""},

	// Vendor Voucher
	{usecase.ErrVoucherNotFound, http.StatusNotFound, "voucher not found"},
	{usecase.ErrVoucherCodeExists, http.StatusConflict, "voucher code already exists"},
	{usecase.ErrInvalidVoucherCode, http.StatusBadRequest, "invalid voucher code"},
	{usecase.ErrInvalidVoucherPeriod, http.StatusBadRequest, ""},
	{usecase.ErrInvalidVoucherQuota, http.StatusBadRequest, ""},
	{usecase.ErrVoucherProductsRequired, http.StatusBadRequest, ""},

	// Address
	{usecase.ErrAddressNotFound, http.StatusNotFound, "address not found"},
	{usecase.ErrAddressNotOwned, http.StatusForbidden, "address does not belong to this user"},
	{usecase.ErrCannotDeleteDefaultAddr, http.StatusUnprocessableEntity, ""},
	{usecase.ErrCannotDeleteLastAddr, http.StatusUnprocessableEntity, ""},

	// Courier
	{usecase.ErrCourierNotFound, http.StatusNotFound, "courier not found"},
	{usecase.ErrCourierNotActive, http.StatusBadRequest, "courier is not active"},
	{usecase.ErrVendorCourierNotFound, http.StatusNotFound, "vendor courier selection not found"},

	// Shipping / RajaOngkir
	{usecase.ErrRajaOngkirFailed, http.StatusBadGateway, ""},

	// Vendor Order Management
	{usecase.ErrOrderNotBelongToVendor, http.StatusForbidden, "order does not belong to this vendor"},
	{usecase.ErrInvalidOrderAcceptTransition, http.StatusBadRequest, ""},
	{usecase.ErrInvalidOrderRejectTransition, http.StatusBadRequest, ""},
	{usecase.ErrInvalidOrderShipTransition, http.StatusBadRequest, ""},
	{usecase.ErrTrackingNumberRequired, http.StatusBadRequest, ""},
	{usecase.ErrShipmentNotFound, http.StatusNotFound, "shipment not found for this order"},
	{usecase.ErrTrackingFailed, http.StatusBadGateway, ""},
	{usecase.ErrInvalidOrderCancelTransition, http.StatusBadRequest, ""},

	// User Bank Account
	{usecase.ErrBankAccountNotFound, http.StatusNotFound, "bank account not found"},
	{usecase.ErrBankAccountNotOwned, http.StatusForbidden, "bank account does not belong to this user"},
}

// HandleUsecaseError maps a usecase error to the appropriate HTTP response.
// It checks the error against all known sentinel errors and returns the
// corresponding status code and message. For unrecognized errors, it returns
// 500 Internal Server Error with a generic message.
func HandleUsecaseError(c *gin.Context, err error) {
	for _, rule := range errorRules {
		if errors.Is(err, rule.sentinel) {
			msg := rule.message
			if msg == "" {
				msg = err.Error()
			}
			response.Error(c, rule.status, msg, nil)
			return
		}
	}
	response.InternalServerError(c, "internal server error", nil)
}
