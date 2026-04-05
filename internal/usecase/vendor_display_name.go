package usecase

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

func vendorDisplayNameOrEmpty(vendor *domain.Vendor) string {
	if vendor == nil || vendor.DisplayName == nil {
		return ""
	}
	return strings.TrimSpace(*vendor.DisplayName)
}

func vendorDisplayNameOrFallback(vendor *domain.Vendor, fallbacks ...string) string {
	if name := vendorDisplayNameOrEmpty(vendor); name != "" {
		return name
	}
	for _, fallback := range fallbacks {
		fallback = strings.TrimSpace(fallback)
		if fallback != "" {
			return fallback
		}
	}
	if vendor != nil && vendor.ID != uuid.Nil {
		return fmt.Sprintf("vendor-%s", vendor.ID.String()[:8])
	}
	return "vendor"
}
