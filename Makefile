-include .env
export

ifeq ($(strip $(DB_PASSWORD)),)
MIGRATE_URL ?= postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
else
MIGRATE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
endif

.PHONY: run build test tidy seed db-setup migrate-up migrate-down migrate-create docker-up docker-down db-create

## run: Start the API server
run:
	go run cmd/api/main.go

## build: Build the API binary
build:
	go build -o bin/api cmd/api/main.go

## test: Run all tests
test:
	go test ./... -v

## tidy: Tidy Go modules
tidy:
	go mod tidy

## db-setup: Ensure database exists and all migrations are applied
db-setup: db-create migrate-up

## seed: Seed the admin user from ADMIN_* environment variables
seed: migrate-up
	go run ./cmd/seed/...

## migrate-up: Run all pending database migrations
migrate-up:
	migrate -path migrations -database "$(MIGRATE_URL)" up

## migrate-down: Roll back the last database migration
migrate-down:
	migrate -path migrations -database "$(MIGRATE_URL)" down 1

## migrate-create: Create a new migration file (usage: make migrate-create name=create_users)
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

## docker-up: Start Docker Compose services
docker-up:
	docker compose up -d

## docker-down: Stop Docker Compose services
docker-down:
	docker compose down

## db-create: Create the application database in PostgreSQL
db-create:
	docker compose exec postgres psql -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);" 2>/dev/null || true
