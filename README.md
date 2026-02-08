# Haji Umroh Store BE

Go backend for a Haji/Umroh souvenir e-commerce platform. Built with Clean Architecture, dual deployment support (HTTP server and AWS Lambda), and CockroachDB as the database.

## Tech Stack

- **Go** (1.25+) with **Gin** HTTP framework
- **GORM** ORM with **CockroachDB** (PostgreSQL-compatible)
- **Viper** for configuration management
- **golang-migrate** for database migrations
- **AWS Lambda** via SAM (Serverless Application Model)
- **go.uber.org/mock** for test mocking

## Architecture

Clean Architecture with four layers (dependency flows inward):

```
Delivery (internal/delivery/http/)  -->  Usecase (internal/usecase/)  -->  Domain (internal/domain/)
                                                                                  ^
                                         Repository (internal/repository/)  ------+
```

| Layer | Path | Responsibility |
|---|---|---|
| Domain | `internal/domain/` | Core entities and interface contracts. No external dependencies. |
| Usecase | `internal/usecase/` | Business logic implementing domain interfaces. |
| Repository | `internal/repository/` | Data access layer (GORM-based). |
| Delivery | `internal/delivery/http/` | Gin HTTP handlers, middleware, and route registration. |
| Infrastructure | `internal/infrastructure/database/` | Database connection setup (CockroachDB via GORM Postgres driver). |
| Config | `internal/config/` | Viper-based configuration loading from `.env` and environment variables. |
| Response | `pkg/response/` | Shared JSON response envelope used by all handlers. |

Dependencies are wired in `internal/app/app.go` via the `Initialize()` function, shared by both HTTP server and Lambda entry points.

## Database Design

![Database Design](docs/database-design.svg)

The schema covers identity & access, vendors, product catalog, cart, orders & shipping, payments & refunds, payouts, and double-entry bookkeeping ledger. Migrations are in the `migrations/` directory.

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) (for Lambda deployment only)

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/media-inovasi-strategis/haji-umroh-store-be.git
   cd haji-umroh-store-be
   ```

2. Copy the environment file and adjust as needed:
   ```bash
   cp .env.example .env
   ```

3. Start the database, create the schema, and run the server:
   ```bash
   make docker-up
   make db-create
   make migrate-up
   make run
   ```

The API server will be available at `http://localhost:8080`.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `APP_NAME` | `haji-umroh-store-be` | Application name |
| `APP_PORT` | `8080` | HTTP server port |
| `APP_ENV` | `development` | Environment (development/production) |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `26257` | Database port |
| `DB_USER` | `root` | Database user |
| `DB_PASSWORD` | _(empty)_ | Database password |
| `DB_NAME` | `haji_umroh_store` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `DB_MAX_IDLE_CONNS` | `10` | Max idle connections |
| `DB_MAX_OPEN_CONNS` | `100` | Max open connections |

## Development Commands

| Command | Description |
|---|---|
| `make run` | Start API server (`go run cmd/api/main.go`) |
| `make build` | Build binary to `bin/api` |
| `make test` | Run all tests (`go test -v ./...`) |
| `make tidy` | Clean Go module dependencies |
| `make docker-up` | Start CockroachDB via Docker Compose |
| `make docker-down` | Stop Docker Compose services |
| `make db-create` | Create `haji_umroh_store` database in CockroachDB |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-create name=<name>` | Create new migration pair in `migrations/` |

Run a single test:
```bash
go test -v -run TestFunctionName ./internal/usecase/...
```

## API Endpoints

All routes are grouped under `/api/v1`. Responses use a standard JSON envelope: `{success, message, data, errors, meta}`.

| Method | Path | Description |
|---|---|---|
| GET | `/api/v1/health` | Health check |

**Middleware chain** (applied in order): Recovery, CORS, Request ID (`X-Request-ID`).

## Project Structure

```
.
├── cmd/
│   ├── api/main.go              # HTTP server entry point
│   └── lambda/main.go           # AWS Lambda entry point
├── internal/
│   ├── app/app.go               # Dependency wiring
│   ├── config/                  # Configuration loading
│   ├── delivery/http/
│   │   ├── handler/             # HTTP handlers
│   │   ├── middleware/          # CORS, Recovery, Request ID
│   │   └── router/             # Route registration
│   ├── domain/                  # Entities and interface contracts
│   ├── infrastructure/database/ # Database connection setup
│   ├── repository/              # Data access implementations
│   └── usecase/                 # Business logic
├── pkg/response/                # Shared JSON response helpers
├── migrations/                  # SQL migration files
├── docs/                        # Documentation and diagrams
├── docker-compose.yml           # CockroachDB setup
├── template.yaml                # AWS SAM template
├── samconfig.toml               # SAM deployment config
├── Makefile                     # Development commands
└── .env.example                 # Environment variable template
```

## AWS Lambda Deployment

The application can be deployed as an AWS Lambda function using SAM.

```bash
make build-lambda       # Cross-compile bootstrap binary for Lambda
make sam-build          # Build SAM application
make sam-deploy-guided  # First-time guided deployment
make sam-deploy         # Subsequent deployments
```

For local Lambda testing:
```bash
make sam-local          # Uses .env.lambda.json for configuration
```

The Lambda function uses the `provided.al2023` runtime with an API Gateway v2 HTTP API. Configuration is defined in `template.yaml` and `samconfig.toml` (region: `ap-southeast-1`).
