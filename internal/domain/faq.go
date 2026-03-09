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

// AdminFAQResponse is the admin response for a single FAQ item (includes all fields).
type AdminFAQResponse struct {
	ID        int       `json:"id"`
	Category  string    `json:"category"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateFAQRequest is the request DTO to create a new FAQ.
type CreateFAQRequest struct {
	Category string `json:"category" binding:"required,max=60"`
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}

// UpdateFAQRequest is the request DTO to update an existing FAQ.
type UpdateFAQRequest struct {
	Category *string `json:"category" binding:"omitempty,max=60"`
	Question *string `json:"question"`
	Answer   *string `json:"answer"`
}

// ============================================================
// Repository Interface
// ============================================================

// FAQRepository provides data-access for FAQs.
type FAQRepository interface {
	List(ctx context.Context, params FAQListParams) ([]FAQ, error)
	ListCategories(ctx context.Context) ([]string, error)
	FindByID(ctx context.Context, id int) (*FAQ, error)
	Create(ctx context.Context, faq *FAQ) error
	Update(ctx context.Context, faq *FAQ) error
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]FAQ, error)
}

// ============================================================
// UseCase Interface
// ============================================================

// FAQUseCase provides business logic for the help center / FAQ feature.
type FAQUseCase interface {
	ListFAQs(ctx context.Context, params FAQListParams) ([]FAQCategoryGroup, error)
	ListCategories(ctx context.Context) ([]string, error)
}

// AdminFAQUseCase provides admin CRUD for FAQs.
type AdminFAQUseCase interface {
	CreateFAQ(ctx context.Context, req CreateFAQRequest) (*AdminFAQResponse, error)
	ListFAQs(ctx context.Context) ([]AdminFAQResponse, error)
	UpdateFAQ(ctx context.Context, id int, req UpdateFAQRequest) (*AdminFAQResponse, error)
	DeleteFAQ(ctx context.Context, id int) error
}
