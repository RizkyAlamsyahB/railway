# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go backend for a Haji/Umroh souvenir e-commerce marketplace. Built with Clean Architecture, Gin, GORM, and PostgreSQL.

## Development Commands

```bash
make run                          # Start API server (go run cmd/api/main.go)
make build                        # Build binary to bin/api
make test                         # Run all tests (go test ./... -v)
make tidy                         # Clean Go module dependencies

# Run a single test
go test -v -run TestFunctionName ./path/to/package/...

# Database
make docker-up                    # Start PostgreSQL via Docker Compose
make docker-down                  # Stop Docker Compose services
make db-setup                     # Create database + apply all migrations
make db-create                    # Create database in PostgreSQL
make migrate-up                   # Apply all pending migrations
make migrate-down                 # Rollback last migration
make migrate-create name=xxx      # Create new migration pair in migrations/

# Seeding
make seed                         # Seed admin user from ADMIN_* env vars
```

## Architecture

Clean Architecture with strict inward dependency flow. All wiring happens in `internal/app/app.go` via constructor injection.

```
Delivery (HTTP handlers)  →  Usecase (business logic)  →  Domain (entities + interfaces)
                                                                     ↑
                              Repository (GORM data access)  ────────┘
```

### Layer conventions

| Layer | Path | Pattern |
|---|---|---|
| Domain | `internal/domain/` | Entities, DTOs, and Go interface contracts. Zero external imports. |
| Usecase | `internal/usecase/` | Implements domain interfaces. Receives repos/providers via constructor. |
| Repository | `internal/repository/` | GORM queries with transaction support for multi-entity ops. |
| Delivery | `internal/delivery/http/handler/` | Gin handlers: bind input → call usecase → return JSON envelope. |
| Middleware | `internal/delivery/http/middleware/` | Recovery, CORS, RequestID, JWT Auth, RequireRoles. |
| Router | `internal/delivery/http/router/` | Gin route groups with middleware composition. |
| Infrastructure | `internal/infrastructure/` | Database connection (GORM/PostgreSQL), S3 storage provider. |
| Config | `internal/config/` | Viper-based loading from `.env` and env vars. |
| Response | `pkg/response/` | Shared JSON envelope: `{success, message, data, errors, meta}`. |
| Auth utils | `pkg/utils/auth/` | JWT generation/validation (HS256), bcrypt password hashing. |

### Adding a new feature (typical flow)

1. Define entity structs and interface in `internal/domain/`
2. Implement repository in `internal/repository/`
3. Implement usecase in `internal/usecase/`
4. Create handler in `internal/delivery/http/handler/`
5. Register routes in `internal/delivery/http/router/router.go`
6. Wire dependencies in `internal/app/app.go`

### Authentication & authorization

- JWT (HS256) with Bearer token in Authorization header.
- Claims contain: `user_id`, `email`, `roles` ([]string), `vendor_id` (optional for umkm role).
- Middleware stores claims in Gin context under keys: `ContextKeyClaims`, `ContextKeyUserID`, `ContextKeyEmail`, `ContextKeyRoles`, `ContextKeyVendorID`.
- Role-based access via `middleware.RequireRoles("admin")` etc.
- Roles: `admin`, `umkm`, `customer`, `cs`, `finance`.

### Key domain models

- **User**: status (pending/active/blocked), M:M with roles via user_roles table.
- **Vendor**: types (umrah_souvenir_store, hajj_souvenir_store, general_souvenir_store), status workflow: draft → submitted → active/rejected/blocked.
- **VendorDocument**: 6 required + 1 optional document types. Required: owner_ktp, owner_passport, store_photo, bank_account_proof, business_logo, business_banner. Optional: business_npwp. S3-backed file storage with presigned URLs.
- **VendorBankAccount**: 1:1 with vendor, has verification status.

### Response pattern

All handlers use `pkg/response/` helpers (`response.OK()`, `response.Created()`, `response.BadRequest()`, etc.) which return a standard JSON envelope.

## Database

- PostgreSQL 17 via Docker Compose (`docker-compose.yml`).
- Migrations managed with `golang-migrate` CLI, stored in `migrations/` as sequential `.up.sql`/`.down.sql` pairs.
- GORM models use string UUIDs and `time.Time` for timestamps.

## Configuration

Environment variables loaded from `.env` via Viper. See `.env.example` for all available variables. Key groups: `APP_*`, `DB_*`, `JWT_*`, `ADMIN_*`, `STORAGE_*`.

## Testing

- Standard Go `testing` package with `go.uber.org/mock` for interface mocking.
- Existing tests in `pkg/utils/auth/` cover JWT and password utilities.
