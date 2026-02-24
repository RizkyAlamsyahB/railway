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
- **`internal/infrastructure/`** — External service implementations: PostgreSQL (GORM), S3 storage (presigned URLs), SMTP email.
- **`internal/app/app.go`** — Wires all dependencies via constructor injection in `Initialize()`.

### Key patterns

**Error handling flow:** Define sentinel error in `usecase/errors.go` → return it from usecase → add mapping rule in `handler/error_mapper.go` → `HandleUsecaseError()` translates to HTTP response.

**API response envelope:** All endpoints use `pkg/response/` which wraps responses in `{success, message, data, errors, meta}`.

**Auth middleware chain:** `middleware.Auth(jwtSecret)` validates JWT and populates context keys (`auth_user_id`, `auth_email`, `auth_role`, `auth_vendor_id`), then `middleware.RequireRoles(...)` checks role.

**Routing:** All routes under `/api/v1`. Role groups: public, `customer`, `umkm` (vendor), `admin`.

**Mocks:** Generated with `go.uber.org/mock` (mockgen) into `internal/usecase/mocks/`. Tests use table-driven patterns with `gomock.Controller` and setup helper functions.

## Configuration

Viper-based, loads from `.env` file with environment variable overrides. Config struct in `internal/config/config.go`. Key prefixes: `APP_`, `DB_`, `JWT_`, `STORAGE_`, `SMTP_`, `ADMIN_`.

## Migrations

SQL files in `migrations/` using golang-migrate. Numbered sequentially (`000001_`, `000002_`, ...). Each migration has `.up.sql` and `.down.sql`.
