# Repository Guidelines

## Project Structure & Module Organization
This project is a Go backend using Clean Architecture.

- `cmd/api/main.go`: API entry point
- `cmd/seed/main.go`: seed script for admin data
- `internal/domain`: core entities and interfaces
- `internal/usecase`: business logic
- `internal/repository`: GORM-based data access
- `internal/delivery/http`: handlers, middleware, router
- `internal/infrastructure`: database, storage, email providers
- `pkg/response`, `pkg/utils`: shared utilities
- `migrations`: SQL schema changes (`*.up.sql` / `*.down.sql`)
- `docs`: OpenAPI specs and DB design artifacts
- `scripts/manual-tests`: manual HTML test pages

## Build, Test, and Development Commands
- `make run`: run the API locally (`go run cmd/api/main.go`)
- `make build`: compile binary to `bin/api`
- `make test`: run all unit tests (`go test ./... -v`)
- `make tidy`: clean module dependencies
- `make docker-up` / `make docker-down`: start/stop PostgreSQL via Docker Compose
- `make db-create`: create app database
- `make migrate-up` / `make migrate-down`: apply or rollback migrations
- `make migrate-create name=add_orders_table`: create a new migration pair
- `make seed`: run migrations and seed admin user

## Coding Style & Naming Conventions
Use standard Go formatting and idioms.

- Format with `gofmt` (or editor auto-format) before commit.
- Keep package names short, lowercase, and domain-oriented.
- Use `CamelCase` for exported identifiers, `camelCase` for internal symbols.
- Prefer constructor-style wiring in `internal/app/app.go` and keep dependency flow inward.
- Migration names should be descriptive and snake_case (example: `add_payment_refund`).

## Testing Guidelines
- Frameworks: Go `testing` + `go.uber.org/mock/gomock`.
- Place tests beside source files using `_test.go`.
- Prefer table-driven tests for usecase and utility logic.
- Reuse generated mocks in `internal/usecase/mocks`.
- Run focused tests with commands like:
  - `go test -v ./internal/usecase -run TestVendorUseCase`

## Commit & Pull Request Guidelines
Commit history follows Conventional Commits with scopes, e.g.:
- `feat(auth): implement user login`
- `test(product): add ProductUseCase unit tests`
- `docs(users): add OpenAPI spec`

For PRs:
- Keep changes scoped and include a clear description.
- Link related issue/ticket when available.
- Note migration and env changes explicitly.
- Include test evidence (`make test` output summary) and sample request/response for API changes.

## Security & Configuration Tips
- Copy `.env.example` to `.env` for local setup.
- Never commit secrets or credentials.
- Validate DB and migration settings before running `migrate` commands.
