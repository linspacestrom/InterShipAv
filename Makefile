.PHONY: run, env, create_db, test

ifneq (,$(wildcard .env))
    include .env
    export
endif

run: format
	@go run cmd/server/main.go || true

env:
	@cp .env.example .env

create_db:
	@echo "Creating database $(POSTGRES_NAME) as user $(POSTGRES_USER)..."
	@sudo -u ${POSTGRES_USER} createdb ${POSTGRES_NAME} || echo "Database may already exist"

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_NAME)
MIGRATIONS_DIR=./migrations

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down

test:
	@go test -count=1 -v ./tests/...

format: fmt vet tidy
	@echo "Code formatted & vetted successfully"

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

tidy:
	@echo "Tidying go modules..."
	@go mod tidy