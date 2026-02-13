# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go backend using Clean Architecture.

- `cmd/api/main.go`: API entrypoint.
- `cmd/seed/main.go`: seed runner for admin bootstrap.
- `internal/domain/`: core entities and interfaces.
- `internal/usecase/`: business logic implementations.
- `internal/repository/`: GORM/PostgreSQL data access.
- `internal/delivery/http/`: Gin handlers, middleware, and router.
- `internal/infrastructure/`: database and external services (e.g., S3).
- `pkg/`: shared utilities (`pkg/response`, `pkg/utils/auth`).
- `migrations/`: sequential `*.up.sql` / `*.down.sql` files.
- `docs/`: OpenAPI specs and schema diagrams.

## Build, Test, and Development Commands
Use Make targets as the standard workflow:

- `make run`: start API locally (`go run cmd/api/main.go`).
- `make build`: build binary to `bin/api`.
- `make test`: run all tests (`go test ./... -v`).
- `make tidy`: sync module dependencies.
- `make docker-up` / `make docker-down`: start/stop PostgreSQL stack.
- `make db-setup`: create DB + run migrations.
- `make migrate-up` / `make migrate-down`: apply/rollback one migration.
- `make migrate-create name=add_orders_table`: create migration files.
- `make seed`: run admin seeding flow.

## Coding Style & Naming Conventions
- Follow standard Go formatting: run `gofmt` (or editor auto-format) before commit.
- Keep package names short, lowercase, no underscores.
- Exported identifiers use `CamelCase`; unexported use `camelCase`.
- Use constructor style `NewXxx(...)` (e.g., handlers/use cases) for dependency injection.
- File naming pattern is feature-oriented snake case, e.g. `vendor_usecase.go`, `auth_handler.go`.

## Testing Guidelines
- Use Go’s built-in `testing` package; mocking uses `go.uber.org/mock` when needed.
- Place tests in `*_test.go` files beside the code under test.
- Prefer table-driven tests for validation and edge-case coverage.
- Run targeted tests with `go test -v -run TestName ./path/to/package/...`.
- No enforced coverage threshold currently; prioritize use case and auth/security-critical paths.

## Commit & Pull Request Guidelines
- Follow Conventional Commit style used in history, e.g.:
  - `feat(vendor): add product management flow`
  - `refactor: relocate admin login endpoint`
  - `chore: dockerize backend services`
- Keep commits scoped and atomic; include migration changes in the same commit when behavior depends on schema updates.
- PRs should include: concise summary, affected modules, test evidence (`make test` output), and API/migration notes.
- If endpoints or contracts change, update relevant docs in `docs/` and `README.md`.
