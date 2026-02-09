# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

| Command | Description |
|---------|-------------|
| `make run` | Start API server on port 8080 |
| `make build` | Build binary to `bin/api` |
| `make test` | Run all tests (`go test -v ./...`) |
| `make tidy` | Clean Go module dependencies |
| `make docker-up` / `make docker-down` | Start/stop CockroachDB |
| `make db-setup` | Create database and run all migrations |
| `make migrate-up` / `make migrate-down` | Apply/rollback migrations |
| `make migrate-create name=<name>` | Create new migration pair |
| `make seed` | Seed admin user from ADMIN_* env vars |
| `make sam-build` | Build SAM application for Lambda |
| `make sam-local` | Start SAM local API |
| `make sam-deploy` | Deploy to AWS |

Run a single test: `go test -v -run TestFunctionName ./path/to/package`

## Architecture

Clean Architecture with four layers — dependency flows inward:

```
Delivery (HTTP) → Usecase → Domain ← Repository
```

- **Domain** (`internal/domain/`): Core entities and interface contracts. Zero external dependencies.
- **Usecase** (`internal/usecase/`): Business logic implementing domain interfaces.
- **Repository** (`internal/repository/`): GORM-based data access implementing domain interfaces.
- **Delivery** (`internal/delivery/http/`): Gin handlers, middleware, route registration.
- **Infrastructure** (`internal/infrastructure/database/`): Database connection setup.
- **Config** (`internal/config/`): Viper-based config loading from `.env` and environment variables.

Wiring happens in `internal/app/app.go` via `Initialize()` which constructs all dependencies and returns an `App` holding Config, DB, and Router.

### Entry Points

- `cmd/api/main.go` — HTTP server with graceful shutdown
- `cmd/lambda/main.go` — AWS Lambda handler using `aws-lambda-go-api-proxy`
- `cmd/seed/main.go` — Admin user seeding utility

### Key Patterns

**Constructor injection everywhere:**
```go
func NewHealthHandler(uc domain.HealthUseCase) *HealthHandler
func NewHealthUseCase() domain.HealthUseCase
```

**Response envelope** (`pkg/response/`): All API responses use `response.OK()`, `response.Created()`, `response.Error()`, `response.BadRequest()`, etc. wrapping data in a consistent `{success, message, data, errors, meta}` JSON structure.

**Middleware chain** (applied in order): Recovery → CORS → RequestID → Auth (optional) → RequireRoles (optional).

**Authentication**: JWT (HS256) with Bearer tokens. Claims include UserID (UUID), Email, Roles. Passwords hashed with bcrypt.

## Database

- **Engine**: CockroachDB (PostgreSQL-compatible), single-node insecure mode for local dev
- **ORM**: GORM with `gorm.io/driver/postgres`
- **Migrations**: SQL files in `migrations/` managed by `golang-migrate` CLI. Migration URL uses `cockroachdb://` scheme.
- **Schema**: 9 migration sets covering identity/access, vendors, catalog, cart, orders/shipping, payments/refunds, payouts, ledger (double-entry bookkeeping), and indexes.

## Configuration

Viper loads from environment variables (highest priority), then `.env` file, then defaults. Config struct in `internal/config/config.go` with sections: App, Database, JWT, Admin. See `.env.example` for all variables.
