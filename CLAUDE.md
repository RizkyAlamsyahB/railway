# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make run              # Start API server (go run cmd/api/main.go)
make build            # Build binary to bin/api
make test             # Run all tests (go test ./... -v)
make tidy             # go mod tidy

# Run a single test function
go test ./internal/usecase/ -run TestUserRegister -v

# Database
make docker-up        # Start PostgreSQL via Docker Compose
make db-setup         # Create database + apply all migrations
make migrate-up       # Apply pending migrations
make migrate-down     # Rollback last migration
make migrate-create name=<name>  # Create new migration pair
make seed             # Seed admin user from ADMIN_* env vars

# Mock generation (no go:generate directives; run manually)
mockgen -destination=internal/usecase/mocks/mock_<name>.go -package=mocks \
  github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain <InterfaceName>
```

## Architecture

Clean architecture with four layers. Dependencies point inward only.

```
delivery/http  →  usecase  →  domain  ←  repository
     │                                        │
     └──── infrastructure (db, s3, smtp) ─────┘
```

### Layer responsibilities

- **`internal/domain/`** — Entities, DTOs (request/response structs with `binding` tags), and interfaces (repository, usecase, StorageProvider, EmailProvider). Zero external dependencies.
- **`internal/usecase/`** — Business logic. Each feature has its own file (e.g., `user_usecase.go`, `product_usecase.go`). Returns sentinel errors defined in `usecase/errors.go`.
- **`internal/repository/`** — GORM-based data access implementing domain repository interfaces. Uses internal GORM models (not domain entities) for DB mapping.
- **`internal/delivery/http/`** — Gin handlers, middleware, router. `handler/error_mapper.go` maps usecase sentinel errors to HTTP status codes.
- **`internal/infrastructure/`** — External service implementations: PostgreSQL (GORM), S3 storage (presigned URLs), SMTP/IMAP email, Xendit payment, RajaOngkir shipping.
- **`pkg/utils/sensitivedata/`** — AES-256-GCM field cipher for encrypting sensitive DB columns.
- **`internal/app/app.go`** — Wires all dependencies via constructor injection in `Initialize()`.

### Key patterns

**Error handling flow:** Define sentinel error in `usecase/errors.go` → return it from usecase → add mapping rule in `handler/error_mapper.go` → `HandleUsecaseError()` translates to HTTP response.

**API response envelope:** All endpoints use `pkg/response/` which wraps responses in `{success, message, data, errors, meta}`.

**Auth middleware chain:** `middleware.Auth(jwtSecret)` validates JWT and populates context keys (`auth_user_id`, `auth_email`, `auth_role`, `auth_vendor_id`), then `middleware.RequireRoles(...)` checks role.

**Routing:** All routes under `/api/v1`. Role groups: public, `customer`, `umkm` (vendor), `admin`.

**Mocks:** Generated with `go.uber.org/mock` (mockgen) into `internal/usecase/mocks/`. Tests use table-driven patterns with `gomock.Controller` and setup helper functions.

## Sensitive Data Encryption

Field-level AES-256-GCM encryption for sensitive vendor data, implemented transparently in the repository layer.

**Implementation:** `pkg/utils/sensitivedata/field_cipher.go`
- Ciphertext format: `enc:v1:{base64(nonce+ciphertext)}` with a random 12-byte nonce per value
- Plaintext values without the `enc:v1:` prefix are passed through unchanged (backward compatibility)

**Encrypted fields:**
- `vendor_bank_accounts.account_number` — via `vendor_repository.go`
- `vendor_onboardings.nik` and `vendor_onboardings.nib` — via `vendor_onboarding_repository.go`

**Flow:** `toXxxModel()` encrypts before DB write; `toDomainXxx()` decrypts after DB read. Use cases operate on plaintext domain objects — encryption is invisible above the repository layer.

**Configuration:** `SENSITIVE_DATA_ENCRYPTION_KEY` must be exactly 32 bytes (AES-256). Set in `.env`. The app fails fast at startup if this key is missing or wrong length (see `app.go` → `newSensitiveDataCipher()`).

**Adding a new encrypted field:** encrypt in `toXxxModel()`, decrypt in `toDomainXxx()`, and change the column type to `TEXT` in a migration if the original type is too short for the ciphertext.

## Configuration

Viper-based, loads from `.env` file with environment variable overrides. Config struct in `internal/config/config.go`. Key prefixes: `APP_`, `DB_`, `JWT_`, `STORAGE_`, `SMTP_`, `ADMIN_`, `SENSITIVE_DATA_`.

## Migrations

SQL files in `migrations/` using golang-migrate. Numbered sequentially (`000001_`, `000002_`, ...). Each migration has `.up.sql` and `.down.sql`.

## Coding Style & Naming Conventions

Use standard Go formatting (`gofmt`) and idiomatic Go naming.
- Package names: lowercase, no underscores.
- Constructors: `NewXxx(...)`.
- Files: descriptive snake_case (e.g., `user_usecase.go`, `user_handler.go`).
- Keep dependency direction inward (`delivery -> usecase -> domain`).
- Add use case sentinel errors in `internal/usecase/errors.go`, then map them in `internal/delivery/http/handler/error_mapper.go`.

## Commit & Pull Request Guidelines

Follow Conventional Commit style: `feat(scope): ...`, `fix(scope): ...`, `test(scope): ...`, `docs(scope): ...`, `refactor(scope): ...`.
- Keep commits scoped to one logical change.
- Mention migration files and config impacts in commit/PR descriptions.
- PRs should include: concise summary, linked issue/task, test commands run, and OpenAPI updates in `docs/openapi-*.yaml` when API behavior changes.

## Security & Configuration Tips

- Copy `.env.example` to `.env`; never commit secrets.
- Validate DB, JWT, SMTP, storage, and `SENSITIVE_DATA_ENCRYPTION_KEY` variables before running locally.
- `SENSITIVE_DATA_ENCRYPTION_KEY` must be exactly 32 bytes; generate one with `openssl rand -hex 16` (produces 32 hex chars = 32 bytes).
- Run migrations before starting the API to avoid schema drift.
