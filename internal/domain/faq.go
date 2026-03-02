package domain

import (
	"context"
	"time"
)

// ============================================================
// Entity
// ============================================================

// FAQ represents a frequently-asked question in the help center.
type FAQ struct {
	ID        int       `json:"id"`
	Category  string    `json:"category"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ValidFAQCategories is the allowed set for FAQ categories.
var ValidFAQCategories = []string{
	"Akun",
	"Produk",
	"Pemesanan",
	"Pembayaran",
	"Pengiriman",
	"Pengembalian & Refund",
	"Promo & Voucher",
	"Chat & Customer Service",
	"Keamanan & Privasi",
	"Vendor/UMKM",
	"Teknis/Aplikasi",
	"Lainnya",
}

// ============================================================
// DTOs
// ============================================================

// FAQListParams holds filters for listing FAQs.
type FAQListParams struct {
	Category string
	Search   string
}

// FAQResponse is the public response for a single FAQ item.
type FAQResponse struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// FAQCategoryGroup groups FAQs by category for structured response.
type FAQCategoryGroup struct {
	Category string        `json:"category"`
	Items    []FAQResponse `json:"items"`
}

// ============================================================
// Repository Interface
// ============================================================

// FAQRepository provides data-access for FAQs.
type FAQRepository interface {
	List(ctx context.Context, params FAQListParams) ([]FAQ, error)
	ListCategories(ctx context.Context) ([]string, error)
}

// ============================================================
// UseCase Interface
// ============================================================

// FAQUseCase provides business logic for the help center / FAQ feature.
type FAQUseCase interface {
	ListFAQs(ctx context.Context, params FAQListParams) ([]FAQCategoryGroup, error)
	ListCategories(ctx context.Context) ([]string, error)
}
