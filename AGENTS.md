# Repository Guidelines

## Project Structure & Module Organization
This service is a Go backend using Clean Architecture.
- `cmd/api/main.go`: HTTP entrypoint.
- `cmd/seed/main.go`: admin seeding command.
- `internal/domain/`: core entities and interfaces.
- `internal/usecase/`: business logic and use case orchestration.
- `internal/repository/`: GORM-backed persistence implementations.
- `internal/delivery/http/`: Gin handlers, middleware, and router wiring.
- `internal/infrastructure/`: external integrations (PostgreSQL, SMTP, S3).
- `pkg/`: shared utilities (`pkg/response`, `pkg/utils/auth`).
- `migrations/`: SQL migrations (`.up.sql` / `.down.sql` pairs).
- `docs/`: OpenAPI specs and design docs.

## Build, Test, and Development Commands
- `make run`: run API locally (`go run cmd/api/main.go`).
- `make build`: build binary to `bin/api`.
- `make test`: run all tests (`go test ./... -v`).
- `make tidy`: tidy Go modules.
- `make docker-up` / `make docker-down`: start/stop local PostgreSQL.
- `make db-setup`: create DB and apply migrations.
- `make migrate-up` / `make migrate-down`: apply or rollback one migration.
- `make migrate-create name=add_orders`: create a numbered migration pair.
- `make seed`: seed admin account from `ADMIN_*` environment variables.

## Coding Style & Naming Conventions
Use standard Go formatting (`gofmt`) and idiomatic Go naming.
- Package names: lowercase, no underscores.
- Constructors: `NewXxx(...)`.
- Files: descriptive snake_case (for example `user_usecase.go`, `user_handler.go`).
- Keep dependency direction inward (`delivery -> usecase -> domain`).
- Add use case sentinel errors in `internal/usecase/errors.go`, then map them in `internal/delivery/http/handler/error_mapper.go`.

## Testing Guidelines
Testing uses Go’s `testing` package with `gomock` for dependency mocking.
- Place tests beside code in `*_test.go`.
- Prefer table-driven tests (`name`, inputs, mock setup, expected error/result).
- Run targeted tests with `go test ./internal/usecase -run TestUserRegister -v`.
- Generate/update mocks under `internal/usecase/mocks/` with `mockgen` when interfaces change.
- No hard coverage gate is enforced; new features should include success and failure-path tests.

## Commit & Pull Request Guidelines
Follow Conventional Commit style seen in history: `feat(scope): ...`, `fix(scope): ...`, `test(scope): ...`, `docs(scope): ...`, `refactor(scope): ...`.
- Keep commits scoped to one logical change.
- Mention migration files and config impacts in commit/PR descriptions.
- PRs should include: concise summary, linked issue/task, test commands run, and OpenAPI updates in `docs/openapi-*.yaml` when API behavior changes.

## Security & Configuration Tips
- Copy `.env.example` to `.env`; never commit secrets.
- Validate DB, JWT, SMTP, and storage variables before running locally.
- Run migrations before starting the API to avoid schema drift.
