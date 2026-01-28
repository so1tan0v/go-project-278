# Variables
BINARY_NAME=link-service
BUILD_DIR=./build
MIGRATIONS_DIR=db/migrations
DB_TYPE=postgres

# Default DATABASE_URL if not set in environment
DATABASE_URL ?= "postgres://user:password@localhost:5432/link_service?sslmode=disable"

.PHONY: build clean run test test-cov run-dev generate-config sqlc migrate-up migrate-down migrate-status migrate-create help

## Build targets
build:
	go build -o $(BUILD_DIR)/app ./main.go

clean:
	rm -rf $(BUILD_DIR)
	rm -rf tmp

## Execution targets
run:
	$(BUILD_DIR)/app

run-dev:
	air -c .air.toml

generate-config:
	cp .env.example .env

## Test targets
test:
	go test ./...

test-cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## Tools & Generation
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

## Migrations (using goose)
migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATIONS_DIR) $(DB_TYPE) $(DATABASE_URL) up

migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATIONS_DIR) $(DB_TYPE) $(DATABASE_URL) down

## Docker targets
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -d --rm -p 8080:8080 --env-file ./.env $(BINARY_NAME)

## Help
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  clean            - Remove build artifacts"
	@echo "  run              - Run the compiled binary"
	@echo "  run-dev          - Run with hot-reload using air"
	@echo "  test             - Run all tests"
	@echo "  test-cov         - Run tests with coverage report"
	@echo "  sqlc             - Generate Go code from SQL queries"
	@echo "  migrate-up       - Apply all up migrations"
	@echo "  migrate-down     - Roll back the last migration"
	@echo "  docker-build     - Build docker image"
	@echo "  docker-run       - Run docker container"
	@echo "  generate-config  - Generate configuration file"
