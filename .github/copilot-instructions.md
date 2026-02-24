# GitHub Copilot Instructions

## Architecture Overview

Go backend for a Haji/Umroh souvenir e-commerce platform using **Clean Architecture**. Dependencies flow strictly inward:

```
delivery/http  →  usecase  →  domain  ←  repository
     │                                        │
     └──── infrastructure (db, s3, smtp) ─────┘
```

| Layer | Path | Key Rule |
|---|---|---|
| Domain | `internal/domain/` | Zero external imports. Entities, request/response DTOs (with `binding` tags), interface contracts, and **all status constants** (`status.go`). |
| Usecase | `internal/usecase/` | Business logic per feature file. Returns only sentinel errors from `errors.go`. Mappers in `mappers.go`. File type validation in `content_type.go`. |
| Repository | `internal/repository/` | GORM-backed. Uses **private model structs** (not domain entities) with `gorm:"column:..."` tags and `TableName()` for DB mapping. |
| Delivery | `internal/delivery/http/` | Gin handlers, middleware (`Auth`, `RequireRoles`, `CORS`, `Recovery`, `RequestID`), route registration in `router/router.go`. |
| Infrastructure | `internal/infrastructure/` | PostgreSQL (GORM), S3 (AWS SDK v2, supports custom endpoint for MinIO), SMTP email. |

All dependencies are wired in `internal/app/app.go` via `Initialize()`. Adding a new feature requires touching all layers plus updating the wiring there.

## Adding a New Feature (Checklist)

1. **Domain** (`internal/domain/`) — add entity, request/response DTOs with `binding` tags, repository interface, usecase interface, and any new status constants to `status.go`.
2. **Usecase** (`internal/usecase/`) — implement logic; add new sentinel errors to `errors.go`.
3. **Repository** (`internal/repository/`) — define private GORM model struct with `TableName()`, implement domain interface. Never pass domain entities to GORM directly.
4. **Error mapping** (`internal/delivery/http/handler/error_mapper.go`) — add a rule to `errorRules` for each new sentinel error. Empty `message` field → `err.Error()` is used verbatim.
5. **Handler** (`internal/delivery/http/handler/`) — bind request, call usecase, use `HandleUsecaseError(c, err)` for errors, `response.*` helpers for success.
6. **Router** (`internal/delivery/http/router/router.go`) — register routes under the correct role group.
7. **Wire** (`internal/app/app.go`) — instantiate repo → usecase → handler; pass handler to `router.NewRouter(...)`.
8. **Migration** (`migrations/`) — create numbered `.up.sql`/`.down.sql` pair with `make migrate-create name=<name>`.

## Key Patterns

**Error flow:** sentinel in `usecase/errors.go` → returned from usecase → matched by `errors.Is` in `handler/error_mapper.go` → HTTP status + message.

**Response envelope:** All handlers use `pkg/response/` — wraps output in `{success, message, data, errors, meta}`. Use helpers: `response.OK`, `response.Created`, `response.BadRequest`, `response.SuccessWithMeta` (for paginated lists). Use `response.Abort` inside middleware (calls `c.AbortWithStatusJSON`).

**Auth middleware chain:**
```go
group.Use(middleware.Auth(jwtSecret))       // validates Bearer JWT, sets context keys
group.Use(middleware.RequireRoles("umkm"))  // enforces role
```
Context key constants (use these, never hardcode strings):
```go
middleware.ContextKeyUserID   // "auth_user_id"   → uuid.UUID
middleware.ContextKeyEmail    // "auth_email"      → string
middleware.ContextKeyRole     // "auth_role"       → string
middleware.ContextKeyVendorID // "auth_vendor_id"  → uuid.UUID (only set for umkm role)
```
JWT `Claims` struct (`pkg/utils/auth/claims.go`): `UserID`, `Email`, `Role`, `VendorID *uuid.UUID`.

**Status constants** — always use `domain.*` constants, never raw strings:
- User: `UserStatusPending`, `UserStatusActive`, `UserStatusBlocked`
- Vendor: `VendorStatusDraft`, `VendorStatusSubmitted`, `VendorStatusActive`, `VendorStatusRejected`, `VendorStatusBlocked`
- Product: `ProductStatusDraft`, `ProductStatusPublished`
- Cart: `CartStatusActive`
- Roles: `RoleAdmin`, `RoleUMKM`, `RoleCustomer`

**Presigned upload pattern** (vendor documents and product images — two-step flow):
1. Client sends create/register request → usecase generates presigned upload URLs via `storage.GeneratePresignedUploadURL`.
2. Client uploads files directly to S3.
3. Client calls confirm endpoint → usecase calls `storage.HeadObject` to verify upload, then saves metadata.
Expiry constants: `usecase.PresignedUploadExpiry` and `usecase.PresignedDownloadExpiry` (both 30 min).

**Pagination:** usecases accept `domain.*ListParams` (Page, Limit, filter fields), return `([]items, *domain.PaginationMeta, error)`. Default page=1, limit=10, max=100.

**Cart:** lazily created on first access via `getOrCreateCart`. One active cart per user (DB `UNIQUE` on `carts.user_id`). `UpsertItem` increments qty if the variant already exists in cart.

**Vendor document types** — required (must be uploaded for vendor to transition to `submitted`): `owner_ktp`, `owner_passport`, `store_photo`, `bank_account_proof`, `business_logo`, `business_banner`. Optional: `business_npwp`.

**File type validation** (`usecase/content_type.go`): images accept `image/jpeg`, `image/png`, `image/webp`; vendor documents additionally accept `application/pdf`.

**Product slug** — auto-generated from name; slug collisions resolved by appending a count suffix. New products always start with `halal_ai_status = "pending"`.

**Admin vendor status machine** — transitions enforced in `admin_vendor_usecase.go`, returns `ErrInvalidStatusTransition` on illegal move:
- `Approve`: `submitted` or `rejected` → `active`
- `Reject`: `submitted` → `rejected`
- `Block`: `active` or `submitted` → `blocked`
- `Unblock`: `blocked` → `active`

**Usecase mappers:** `internal/usecase/mappers.go` converts domain entities to response DTOs (e.g., `toUserResponse`). Add all new mappers here.

## Tests

Tests live beside code in `*_test.go` files inside the `usecase` package (white-box). Pattern:

```go
// 1. Setup helper returns mocks + usecase under test
func setupUserUseCase(t *testing.T) (*mocks.MockUserRepository, ..., domain.UserUseCase) {
    ctrl := gomock.NewController(t)
    // ... create mocks, build usecase
}

// 2. Table-driven test cases with per-case mock setup
testCases := []struct {
    name       string
    setupMocks func(...)
    wantErr    error
}{ ... }
```

**Mock generation** (run manually when interfaces change — no `go:generate` directives exist):
```bash
mockgen -destination=internal/usecase/mocks/mock_<name>.go -package=mocks \
  github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain <InterfaceName>
```

Run a single test:
```bash
go test ./internal/usecase/ -run TestUserRegister -v
```

## Development Commands

```bash
make run              # Start API (go run cmd/api/main.go)
make test             # Run all tests (go test ./... -v)
make docker-up        # Start PostgreSQL via Docker Compose
make db-setup         # Create DB + apply all migrations
make migrate-create name=add_orders  # Create numbered SQL migration pair
make seed             # Seed admin user from ADMIN_* env vars
```

## Configuration

Viper-based, loaded from `.env`. Key prefixes: `APP_`, `DB_`, `JWT_`, `STORAGE_`, `SMTP_`, `ADMIN_`. Config struct at `internal/config/config.go`. S3 supports `STORAGE_S3_ENDPOINT` + `STORAGE_S3_FORCE_PATH_STYLE` for local MinIO. Never commit `.env`.

## Database Schema & Migrations

SQL files in `migrations/` numbered sequentially. Schema coverage:
- `000001` — `users`, `roles`, `email_verification_tokens`, `addresses`
- `000002` — `vendors`, `vendor_documents`, `vendor_bank_accounts`
- `000003` — `categories`, `products`, `product_variants`, `product_images`, `shipping_services`, `product_shipping_services`
- `000004` — `carts`, `cart_items`
- `000005` — `orders`, `order_items`, `order_status_history`, `shipments`
- `000006` — `payment_invoices`, `payment_events`, `refunds`
- `000007` — `payout_batches`, `payout_items`
- `000008` — `ledger_accounts`, `ledger_journals`, `ledger_lines` (double-entry bookkeeping)

**Migrations 5–8 (orders, payments, payouts, ledger) have DB schema but no Go handlers yet** — planned features. Do not add handlers without confirming domain design first.

Note: `users.email` unique index is case-insensitive (`lower(email)`). All monetary columns use `NUMERIC(18,2)`. All PKs are UUID.

## Routes Summary (`/api/v1`)

| Path prefix | Auth | Role |
|---|---|---|
| `GET /health`, `GET /categories`, `GET /shipping-services` | None | — |
| `POST /users/register`, `GET /users/verify-email`, `POST /users/resend-verification`, `POST /users/login` | None | — |
| `GET /users/me`, `GET\|POST\|PATCH\|DELETE /users/cart*` | JWT | `customer` |
| `POST /vendors/register`, `POST /vendors/login` | None | — |
| `POST /vendors/documents/confirm`, `POST /vendors/products*` | JWT | `umkm` |
| `POST /admin/login` | None | — |
| All other `/admin/*` | JWT | `admin` |

## OpenAPI Docs

Specs split by domain in `docs/`: `openapi-admin.yaml`, `openapi-users.yaml`, `openapi-vendor.yaml`, `openapi-catalog.yaml`, `openapi-cart.yaml`. Update the relevant file when changing API behavior.
