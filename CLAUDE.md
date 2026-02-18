# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make run                        # Start API server (go run cmd/api/main.go)
make build                      # Build binary to bin/api
make test                       # Run all tests (go test ./... -v)
make tidy                       # go mod tidy

# Run a single test
go test -v -run TestFunctionName ./internal/usecase/...

# Database
make docker-up                  # Start PostgreSQL via Docker Compose
make docker-down                # Stop Docker Compose services
make db-create                  # Create haji_umroh_store database
make db-setup                   # db-create + migrate-up
make migrate-up                 # Apply all pending migrations
make migrate-down               # Rollback last migration
make migrate-create name=<name> # Create new migration pair in migrations/
make seed                       # Seed admin user (runs migrate-up first)
```

Setup from scratch: `cp .env.example .env && make docker-up && make db-setup && make seed && make run`

## Architecture

Clean Architecture with inward dependency flow. Go 1.25+, Gin framework, GORM + PostgreSQL, Viper config.

```
Delivery (HTTP handlers) → Usecase (business logic) → Domain (entities & interfaces)
                                                              ↑
                                    Repository (GORM data access) ──────┘
```

### Layer responsibilities and locations

| Layer | Path | Role |
|---|---|---|
| Domain | `internal/domain/` | Entities, repository/usecase/provider interfaces, DTOs. Zero external deps. |
| Usecase | `internal/usecase/` | Business logic implementations. Depends only on domain interfaces. |
| Repository | `internal/repository/` | GORM-based data access implementing domain repository interfaces. |
| Delivery | `internal/delivery/http/` | Gin handlers (`handler/`), middleware (`middleware/`), routes (`router/`). |
| Infrastructure | `internal/infrastructure/` | External integrations: `database/` (Postgres), `email/` (SMTP), `storage/` (S3). |
| Config | `internal/config/` | Viper-based config loader from `.env` + env vars. |
| App | `internal/app/app.go` | Dependency injection wiring — `Initialize()` builds the entire object graph. |
| Pkg | `pkg/` | Reusable packages: `response/` (JSON envelope), `utils/auth/` (JWT, bcrypt). |

### Key patterns

- **Interface-driven**: Domain defines all interfaces; repository/infrastructure implement them. New implementations swap in without touching business logic.
- **DI wiring**: `internal/app/app.go` `Initialize()` constructs all dependencies in order: config → DB → storage → email → repositories → usecases → handlers → router.
- **Standard response envelope**: All handlers use `pkg/response/` helpers (`response.OK()`, `response.Created()`, `response.BadRequest()`, etc.) returning `{success, message, data, errors, meta}`.
- **Auth middleware**: JWT auth + role-based access via `middleware.AuthMiddleware()` and `middleware.RequireRoles()`. Claims carry UserID, Email, Role, VendorID.
- **Mocks**: Generated with `go.uber.org/mock` in `internal/usecase/mocks/`. Usecase tests mock repository interfaces.

### Adding a new feature

1. Define entities and interfaces in `internal/domain/`.
2. Implement repository in `internal/repository/`.
3. Implement usecase in `internal/usecase/`.
4. Add handler in `internal/delivery/http/handler/`.
5. Register routes in `internal/delivery/http/router/router.go`.
6. Wire everything in `internal/app/app.go`.

### API routes

All routes under `/api/v1`. Route groups: public, user auth (`/users`), vendor auth (`/vendors`), admin auth (`/admin`). Auth-protected groups use `AuthMiddleware` + `RequireRoles` middleware.

### Database migrations

Sequential SQL files in `migrations/` managed by golang-migrate. 9 migrations covering: identity/access, vendors, catalog, cart, orders/shipping, payments/refunds, payouts, double-entry ledger, and indexes.
