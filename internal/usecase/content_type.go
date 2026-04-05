package usecase

import "strings"

// Allowed MIME types for product images.
var allowedImageContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

// Allowed MIME types for vendor documents (images + PDF).
var allowedDocumentContentTypes = map[string]struct{}{
	"image/jpeg":      {},
	"image/png":       {},
	"image/webp":      {},
	"application/pdf": {},
}

// Allowed MIME types for individual onboarding documents.
var allowedIndividualOnboardingContentTypes = map[string]struct{}{
	"image/jpeg":      {},
	"application/pdf": {},
}

// normalizeContentType strips parameters (e.g. charset) and normalizes
// a Content-Type header value to its lowercase media type.
func normalizeContentType(contentType string) string {
	normalized := strings.TrimSpace(strings.ToLower(contentType))
	if idx := strings.Index(normalized, ";"); idx >= 0 {
		normalized = strings.TrimSpace(normalized[:idx])
	}
	return normalized
}

// isAllowedImageContentType reports whether the given (normalized) content
// type is an acceptable product image MIME type.
func isAllowedImageContentType(contentType string) bool {
	_, ok := allowedImageContentTypes[contentType]
	return ok
}

// isAllowedDocumentContentType reports whether the given (normalized) content
// type is an acceptable vendor document MIME type.
func isAllowedDocumentContentType(contentType string) bool {
	_, ok := allowedDocumentContentTypes[contentType]
	return ok
}

// isAllowedIndividualOnboardingContentType reports whether the given
// (normalized) content type is acceptable for individual onboarding uploads.
func isAllowedIndividualOnboardingContentType(contentType string) bool {
	_, ok := allowedIndividualOnboardingContentTypes[contentType]
	return ok
}

// Allowed MIME types for NIB (corporate legal) documents: PDF only.
var allowedNIBDocumentContentTypes = map[string]struct{}{
	"application/pdf": {},
}

// isAllowedNIBDocumentContentType reports whether the given (normalized)
// content type is acceptable for a corporate NIB document.
func isAllowedNIBDocumentContentType(contentType string) bool {
	_, ok := allowedNIBDocumentContentTypes[contentType]
	return ok
}

// Allowed MIME types for ticket attachments (images + mp4 video, max 5 MB).
var allowedTicketAttachmentContentTypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"video/mp4":  {},
}

const TicketAttachmentMaxBytes int64 = 5 * 1024 * 1024 // 5 MB

// isAllowedTicketAttachmentContentType reports whether the given (normalized)
// content type is acceptable for a ticket attachment.
func isAllowedTicketAttachmentContentType(contentType string) bool {
	_, ok := allowedTicketAttachmentContentTypes[contentType]
	return ok
}

// Allowed MIME types for chat attachments (images + mp4 video, max 10 MB).
var allowedChatAttachmentContentTypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"video/mp4":  {},
}

const ChatAttachmentMaxBytes int64 = 10 * 1024 * 1024 // 10 MB

// isAllowedChatAttachmentContentType reports whether the given (normalized)
// content type is acceptable for a chat attachment.
func isAllowedChatAttachmentContentType(contentType string) bool {
	_, ok := allowedChatAttachmentContentTypes[contentType]
	return ok
}

const ReviewImageMaxBytes int64 = 5 * 1024 * 1024 // 5 MB

// Allowed MIME types for refund evidence uploads (images + one optional mp4 video).
var allowedRefundEvidenceContentTypes = map[string]struct{}{
	"image/png":  {},
	"image/jpeg": {},
	"image/webp": {},
	"video/mp4":  {},
}

const RefundEvidenceMaxBytes int64 = 10 * 1024 * 1024 // 10 MB per file

func isAllowedRefundEvidenceContentType(contentType string) bool {
	_, ok := allowedRefundEvidenceContentTypes[contentType]
	return ok
}

func isRefundEvidenceImageContentType(contentType string) bool {
	switch contentType {
	case "image/png", "image/jpeg", "image/webp":
		return true
	default:
		return false
	}
}

func isRefundEvidenceVideoContentType(contentType string) bool {
	return contentType == "video/mp4"
}
