include .env
export

MIGRATE_URL ?= cockroachdb://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: run build test tidy migrate-up migrate-down migrate-create docker-up docker-down db-create build-lambda sam-build sam-local sam-deploy sam-deploy-guided build-ApiFunction

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

## db-create: Create the application database in CockroachDB
db-create:
	docker compose exec cockroachdb cockroach sql --insecure -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME);"

## build-lambda: Cross-compile the Lambda binary for Amazon Linux
build-lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap cmd/lambda/main.go

## build-ApiFunction: SAM build target (called by SAM BuildMethod: makefile)
build-ApiFunction:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o $(ARTIFACTS_DIR)/bootstrap cmd/lambda/main.go

## sam-build: Build the SAM application
sam-build:
	sam build

## sam-local: Start SAM local API for testing
sam-local:
	sam local start-api --env-vars .env.lambda.json

## sam-deploy: Deploy to AWS via SAM
sam-deploy:
	sam deploy

## sam-deploy-guided: First-time guided SAM deployment
sam-deploy-guided:
	sam deploy --guided
