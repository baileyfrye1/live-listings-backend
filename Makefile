# Simple Makefile for a Go project
include .env

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o bin/main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Run database migration
up:
	@goose -dir=database/migrations postgres postgresql://$(BLUEPRINT_DB_USERNAME):$(BLUEPRINT_DB_PASSWORD)@$(BLUEPRINT_DB_HOST):$(BLUEPRINT_DB_PORT)/$(BLUEPRINT_DB_DATABASE) up

down:
	@goose -dir=database/migrations postgres postgresql://$(BLUEPRINT_DB_USERNAME):$(BLUEPRINT_DB_PASSWORD)@$(BLUEPRINT_DB_HOST):$(BLUEPRINT_DB_PORT)/$(BLUEPRINT_DB_DATABASE) down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
dev:

	@if command -v air > /dev/null; then \
            air --build.cmd "go build -o bin/main cmd/api/main.go" --build.bin "./bin/main"; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest
