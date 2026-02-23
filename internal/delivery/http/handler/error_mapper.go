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

	// Vendor
	{usecase.ErrVendorAlreadyExists, http.StatusConflict, "user already has a vendor"},
	{usecase.ErrVendorNotFound, http.StatusNotFound, "vendor not found"},
	{usecase.ErrDocumentNotFound, http.StatusBadRequest, ""},
	{usecase.ErrObjectNotUploaded, http.StatusBadRequest, ""},
	{usecase.ErrInvalidDocumentContent, http.StatusBadRequest, ""},
	{usecase.ErrDocumentSizeOverflow, http.StatusBadRequest, ""},
	{usecase.ErrVendorInvalidCredentials, http.StatusUnauthorized, "invalid email or password"},
	{usecase.ErrVendorAccountBlocked, http.StatusForbidden, "vendor account is blocked"},
	{usecase.ErrNotVendor, http.StatusForbidden, "vendor access required"},
	{usecase.ErrNoVendorProfile, http.StatusForbidden, "no vendor profile found"},
	{usecase.ErrInvalidStatusTransition, http.StatusBadRequest, ""},

	// Xendit
	{usecase.ErrXenditAccountCreation, http.StatusBadGateway, ""},

	// Product
	{usecase.ErrVendorNotActive, http.StatusForbidden, "vendor is not active"},
	{usecase.ErrCategoryNotFound, http.StatusBadRequest, "category not found"},
	{usecase.ErrShippingServiceNotFound, http.StatusBadRequest, "one or more shipping services not found"},
	{usecase.ErrProductNotFound, http.StatusNotFound, "product not found"},
	{usecase.ErrProductNotOwned, http.StatusForbidden, "product does not belong to this vendor"},
	{usecase.ErrImageNotFound, http.StatusBadRequest, ""},
	{usecase.ErrImageNotUploaded, http.StatusBadRequest, ""},
	{usecase.ErrInvalidImageContentType, http.StatusBadRequest, ""},
	{usecase.ErrImageSizeOverflow, http.StatusBadRequest, ""},
	{usecase.ErrDuplicatePrimaryImage, http.StatusBadRequest, "only one image can be marked as primary"},
	{usecase.ErrNoPrimaryImage, http.StatusBadRequest, "at least one image must be marked as primary"},
	{usecase.ErrPublishedRequiresShipping, http.StatusBadRequest, "published product must have at least one shipping service"},
	{usecase.ErrTooManyImages, http.StatusBadRequest, "maximum 10 images per product"},

	// Cart
	{usecase.ErrCartItemNotFound, http.StatusNotFound, "cart item not found"},
	{usecase.ErrCartItemNotOwned, http.StatusForbidden, "cart item does not belong to your cart"},
	{usecase.ErrVariantNotFound, http.StatusNotFound, "product variant not found"},
	{usecase.ErrVariantNotActive, http.StatusBadRequest, "product variant is not active"},
	{usecase.ErrProductNotAvailable, http.StatusBadRequest, "product is not available"},
	{usecase.ErrInsufficientStock, http.StatusConflict, "insufficient stock for requested quantity"},

	// Finance
	{usecase.ErrInvalidMonth, http.StatusBadRequest, "invalid month format, expected YYYY-MM"},
	{usecase.ErrInvalidPaymentStatus, http.StatusBadRequest, "invalid payment status"},
	{usecase.ErrInvalidRefundStatus, http.StatusBadRequest, "invalid refund status"},
	{usecase.ErrInvalidPayoutStatus, http.StatusBadRequest, "invalid payout status"},
	{usecase.ErrInvalidRefundTransition, http.StatusBadRequest, "invalid refund status transition"},
	{usecase.ErrRefundNotFound, http.StatusNotFound, "refund not found"},
	{usecase.ErrPayoutNotFound, http.StatusNotFound, "payout batch not found"},

	// CS
	{usecase.ErrTicketNotFound, http.StatusNotFound, "ticket not found"},
	{usecase.ErrConversationNotFound, http.StatusNotFound, "conversation not found"},
	{usecase.ErrConversationNotAllowed, http.StatusForbidden, "chat between these roles is not allowed"},
	{usecase.ErrConversationUnauthorized, http.StatusForbidden, "you are not a participant of this conversation"},
	{usecase.ErrReplyTemplateNotFound, http.StatusNotFound, "reply template not found"},
	{usecase.ErrCSUserNotFound, http.StatusNotFound, "user not found"},

	// ticket attachment
	{usecase.ErrInvalidAttachmentContentType, http.StatusUnprocessableEntity, ""},
	{usecase.ErrAttachmentTooLarge, http.StatusUnprocessableEntity, ""},
	{usecase.ErrAttachmentNotUploaded, http.StatusUnprocessableEntity, ""},

	// ticket ownership
	{usecase.ErrTicketAlreadyTaken, http.StatusConflict, "ticket is already assigned to another CS"},
	{usecase.ErrTicketNotAssignedToYou, http.StatusForbidden, "only the assigned CS can perform this action"},
	{usecase.ErrTicketClosed, http.StatusUnprocessableEntity, "ticket is already resolved or closed"},

	// ticket subject
	{usecase.ErrTicketSubjectNotFound, http.StatusBadRequest, "ticket subject not found"},

	// reply template
	{usecase.ErrInvalidReplyTemplateCategory, http.StatusUnprocessableEntity, ""},
	{usecase.ErrShortcutAlreadyExists, http.StatusConflict, ""},
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
