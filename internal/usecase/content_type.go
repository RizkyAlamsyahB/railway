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
