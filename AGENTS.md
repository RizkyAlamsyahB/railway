# Repository Guidelines

## Project Structure & Module Organization
This service follows clean architecture with dependencies pointing inward:

- `cmd/api` and `cmd/seed`: application entrypoints.
- `internal/domain`: entities, DTOs, and repository/usecase interfaces.
- `internal/usecase`: business logic and sentinel errors in `errors.go`.
- `internal/repository`: GORM-backed persistence implementations.
- `internal/delivery/http`: Gin handlers, middleware, websockets, and routing under `/api/v1`.
- `internal/infrastructure`: PostgreSQL, SMTP, S3, payment, and shipping integrations.
- `pkg/response` and `pkg/utils`: shared helpers.
- `migrations/`: sequential SQL migrations (`000001_*.up.sql`, `*.down.sql`).
- `docs/`: OpenAPI specs, design notes, and Postman collections.

Keep tests beside the code they verify, for example `internal/usecase/user_usecase_test.go`.

## Build, Test, and Development Commands
- `make run`: start the API with `go run cmd/api/main.go`.
- `make build`: compile the binary to `bin/api`.
- `make test`: run the full test suite with verbose output.
- `go test ./internal/usecase -run TestUserRegister -v`: run one targeted test.
- `make tidy`: sync `go.mod` and `go.sum`.
- `make docker-up` / `make docker-down`: start or stop PostgreSQL locally.
- `make db-setup`, `make migrate-up`, `make migrate-down`: manage schema changes.
- `make seed`: seed the admin user from `ADMIN_*` environment variables.

## Coding Style & Naming Conventions
Use idiomatic Go and run `gofmt` on edited files. Package names stay lowercase; files use descriptive snake_case such as `vendor_order_usecase.go`. Prefer constructor names like `NewXxx`. Preserve layer boundaries: `delivery -> usecase -> domain <- repository`.

When adding a new use case error, define it in `internal/usecase/errors.go` and map it in `internal/delivery/http/handler/error_mapper.go`.

## Testing Guidelines
Tests use Go’s `testing` package with `gomock` for mocks in `internal/usecase/mocks/`. Prefer table-driven tests with clear `name`, setup, expected result, and expected error branches. New behavior should include both success and failure-path coverage.

## Commit & Pull Request Guidelines
Use Conventional Commits, for example `feat(user): add email verification` or `docs(openapi): normalize vendor spec tags`. PRs should include a concise summary, linked task or issue, test commands run, and updated OpenAPI files in `docs/openapi-*.yaml` when API behavior changes.

## Security & Configuration Tips
Load configuration from `.env` with environment overrides via Viper. Never commit secrets. Validate `DB_*`, `JWT_*`, `SMTP_*`, `STORAGE_*`, and `ADMIN_*` values before running locally, and apply migrations before starting the server.

This project already includes field-level sensitive data encryption in `pkg/utils/sensitivedata/field_cipher.go`, initialized from config in `internal/app/app.go`. Use that path for sensitive string fields stored in the database instead of adding ad hoc crypto. Current examples include vendor bank account numbers in `internal/repository/vendor_repository.go` and onboarding `NIK`/`NIB` values in `internal/repository/vendor_onboarding_repository.go`. Preserve decrypt support for legacy plaintext values and add tests when extending encrypted fields.
