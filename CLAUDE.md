# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
make run                              # Start API server (go run cmd/api/main.go)
make build                            # Build binary to bin/api
make test                             # Run all tests (go test -v ./...)
make tidy                             # Clean Go module dependencies

make docker-up                        # Start CockroachDB via Docker Compose
make docker-down                      # Stop Docker Compose services
make db-create                        # Create haji_umroh_store database in CockroachDB

make migrate-up                       # Apply all pending migrations
make migrate-down                     # Rollback last migration
make migrate-create name=<name>       # Create new migration pair in migrations/
```

**Run a single test:**
```bash
go test -v -run TestFunctionName ./internal/usecase/...
```

**First-time setup:** `make docker-up` → `make db-create` → `make migrate-up` → `make run`

## Architecture

Clean Architecture with four layers (dependency flows inward):

```
Delivery (internal/delivery/http/)  →  Usecase (internal/usecase/)  →  Domain (internal/domain/)
                                                                              ↑
                                       Repository (internal/repository/)  ────┘
```

- **Domain** (`internal/domain/`): Core entities and interface contracts. No external dependencies.
- **Usecase** (`internal/usecase/`): Business logic implementing domain interfaces.
- **Repository** (`internal/repository/`): Data access layer implementing domain interfaces (uses GORM).
- **Delivery** (`internal/delivery/http/`): Gin HTTP handlers, middleware, and route registration.
- **Infrastructure** (`internal/infrastructure/database/`): Database connection setup (CockroachDB via GORM Postgres driver).
- **Config** (`internal/config/`): Viper-based configuration loading from `.env` and environment variables.
- **Response** (`pkg/response/`): Shared JSON response envelope used by all handlers.

## Key Patterns

- **Manual constructor injection** — dependencies are wired in `internal/app/app.go` via `Initialize()`, shared by both entry points. New features follow: create domain interface → implement usecase → implement repository → create handler → inject in `app.go` → register routes in router.
- **Dual entry points**: `cmd/api/main.go` (HTTP server with graceful shutdown) and `cmd/lambda/main.go` (AWS Lambda via API Gateway v2 proxy). Both call `app.Initialize()` for identical dependency wiring.
- **Standard response envelope** (`pkg/response/`): all API responses use `response.Success()`, `response.Error()`, `response.Abort()`, etc. with structure `{success, message, data, errors, meta}`.
- **Middleware chain**: Recovery → CORS → Request ID (UUID via `X-Request-ID` header), applied in `internal/delivery/http/router/router.go`.
- **Route grouping**: all routes under `/api/v1` group in `internal/delivery/http/router/router.go`.

## AWS Lambda / SAM Deployment

```bash
make build-lambda                     # Cross-compile bootstrap binary for Lambda
make sam-build                        # Build SAM application
make sam-local                        # Local testing via SAM (uses .env.lambda.json)
make sam-deploy                       # Deploy to AWS
make sam-deploy-guided                # First-time guided deployment
```

Lambda uses `provided.al2023` runtime with API Gateway v2 HTTP API. Configuration is in `template.yaml` (SAM) and `samconfig.toml` (deploy settings, region `ap-southeast-1`). Local Lambda testing uses `host.docker.internal` as DB host to reach the host machine's CockroachDB.

## Tech Stack

- **Go** with **Gin** (HTTP framework), **GORM** (ORM), **Viper** (config)
- **CockroachDB** v25.1 (PostgreSQL-compatible) — port 26257 for SQL, 8090 for admin UI
- **golang-migrate** for database migrations (CockroachDB driver)
- **go.uber.org/mock** for test mocking
