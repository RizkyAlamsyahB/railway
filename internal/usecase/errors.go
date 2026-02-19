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

// --- Vendor errors ---

var (
	ErrEmailAlreadyRegistered   = errors.New("email already registered")
	ErrVendorAlreadyExists      = errors.New("user already has a vendor")
	ErrVendorNotFound           = errors.New("vendor not found")
	ErrDocumentNotFound         = errors.New("document not found for this vendor")
	ErrObjectNotUploaded        = errors.New("object not found in storage")
	ErrInvalidDocumentContent   = errors.New("invalid document content type")
	ErrDocumentSizeOverflow     = errors.New("document file size exceeds supported limit")
	ErrVendorInvalidCredentials = errors.New("invalid email or password")
	ErrVendorAccountBlocked     = errors.New("vendor account is blocked")
	ErrNotVendor                = errors.New("user does not have vendor access")
	ErrNoVendorProfile          = errors.New("no vendor profile found for this user")
	ErrInvalidStatusTransition  = errors.New("invalid status transition")
)

// --- Product errors ---

var (
	ErrVendorNotActive           = errors.New("vendor is not active")
	ErrCategoryNotFound          = errors.New("category not found")
	ErrShippingServiceNotFound   = errors.New("one or more shipping services not found")
	ErrProductNotFound           = errors.New("product not found")
	ErrProductNotOwned           = errors.New("product does not belong to this vendor")
	ErrImageNotFound             = errors.New("image not found for this product")
	ErrImageNotUploaded          = errors.New("image not found in storage")
	ErrInvalidImageContentType   = errors.New("invalid image content type")
	ErrImageSizeOverflow         = errors.New("image file size exceeds supported limit")
	ErrDuplicatePrimaryImage     = errors.New("only one image can be marked as primary")
	ErrNoPrimaryImage            = errors.New("at least one image must be marked as primary")
	ErrPublishedRequiresShipping = errors.New("published product must have at least one shipping service")
	ErrTooManyImages             = errors.New("maximum 10 images per product")
)

// --- Cart errors ---

var (
	ErrCartItemNotFound    = errors.New("cart item not found")
	ErrCartItemNotOwned    = errors.New("cart item does not belong to your cart")
	ErrVariantNotFound     = errors.New("product variant not found")
	ErrVariantNotActive    = errors.New("product variant is not active")
	ErrProductNotAvailable = errors.New("product is not available")
	ErrInsufficientStock   = errors.New("insufficient stock")
)
