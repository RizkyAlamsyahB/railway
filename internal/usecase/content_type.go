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
