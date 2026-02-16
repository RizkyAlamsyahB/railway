# Repository Guidelines

## Project Structure & Module Organization
This backend is written in Go and follows Clean Architecture.

- `cmd/api/main.go`: API entrypoint.
- `cmd/seed/main.go`: admin/bootstrap seed runner.
- `internal/domain/`: core entities and interfaces.
- `internal/usecase/`: business logic orchestration.
- `internal/repository/`: PostgreSQL access via GORM.
- `internal/delivery/http/`: Gin handlers, middleware, and routing.
- `internal/infrastructure/`: DB wiring and external services (for example S3).
- `pkg/`: shared helpers such as `pkg/response` and auth utilities.
- `migrations/`: ordered `*.up.sql` and `*.down.sql` schema changes.
- `docs/`: OpenAPI specs and schema references.

Place tests next to implementation files using `*_test.go`.

## Build, Test, and Development Commands
Use `make` targets as the standard workflow:

- `make run`: start the API locally.
- `make build`: compile binary to `bin/api`.
- `make test`: run all tests with verbose output (`go test ./... -v`).
- `make tidy`: synchronize `go.mod` and `go.sum`.
- `make docker-up` / `make docker-down`: start/stop local PostgreSQL stack.
- `make db-setup`: create DB and apply migrations.
- `make migrate-up` / `make migrate-down`: apply or rollback one migration.
- `make migrate-create name=add_orders_table`: scaffold a new migration pair.
- `make seed`: run admin seeding flow.

## Coding Style & Naming Conventions
- Format code with `gofmt` before committing.
- Keep package names lowercase and short (no underscores).
- Use `CamelCase` for exported identifiers and `camelCase` for private ones.
- Prefer dependency-injection constructors like `NewAuthHandler(...)`.
- Use feature-oriented snake case filenames, such as `vendor_usecase.go`.

## Testing Guidelines
- Use Go’s `testing` package; use `go.uber.org/mock` for mocks when needed.
- Prefer table-driven tests for validation and edge cases.
- Run focused tests with commands like:
  `go test -v -run TestVendorCreate ./internal/usecase/...`
- Prioritize coverage on use cases and auth/security-critical paths.

## Commit & Pull Request Guidelines
- Follow Conventional Commits, for example:
  `feat(vendor): add product management flow`
- Keep commits atomic; include related migration changes in the same commit.
- PRs should include a concise summary, affected modules, and test evidence (for example `make test` output).
- Update `docs/` and `README.md` when API contracts or behavior change.
