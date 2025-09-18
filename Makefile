# Makefile for Kawasan Digital Ingat Dok Backend Project

# Project variables
PROJECT_NAME := itsrama-portfolio-backend
GO_VERSION := 1.24.4
DOCKER_REGISTRY ?= holycann
VERSION ?= $(shell git describe --tags --always --dirty)

# Go parameters
GO := go
GOBUILD := $(GO) build
GOTEST := $(GO) test
GOMOD := $(GO) mod
GORUN := $(GO) run

# Directories
CMD_DIR := ./cmd
BUILD_DIR := ./build
DIST_DIR := ./dist
MIGRATIONS_DIR := ./db/migrations

# Main application entry point
MAIN_APP := $(CMD_DIR)/main.go

# Database migration tool
MIGRATE := migrate

# Linting and formatting
GOLANGCI_LINT := golangci-lint

# Docker parameters
DOCKER_COMPOSE := docker-compose
DOCKERFILE := Dockerfile

# Postman parameters
COLLECTION_DIR := postman
COLLECTION_NAME := Itsrama
COLLECTION_FILES := *_collection.json
COLLECTION_OUTPUT_FILE := itsrama_portfolio_backend.json

# Targets
.PHONY: all clean build test lint run docker-build docker-push deps \
		migrate-up migrate-down migrate-create swagger dev-up dev-down help

# Default target
all: clean deps lint test build

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@$(GO) clean -cache

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GOMOD) tidy
	@$(GOMOD) download

# Run linter
lint:
	@echo "Running linter..."
	@$(GOLANGCI_LINT) run ./...

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Build the application
build:
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -o $(BUILD_DIR)/$(PROJECT_NAME) \
		-ldflags "-X main.Version=$(VERSION) \
				  -X main.BuildTime=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
				  -X main.GitCommit=$(shell git rev-parse HEAD)" \
		$(MAIN_APP)

# Run the application locally
run:
	@echo "Running application..."
	@$(GORUN) $(MAIN_APP)

# Database Migrations
migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)


migrate-up:
	@echo "Running database migrations (up)..."
	@go run scripts/migrate.go --up

migrate-down:
	@echo "Reverting database migrations (down)..."
	@go run scripts/migrate.go --down

# Docker targets
docker-build:
	@echo "Building Docker image..."
	@docker build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		-t $(DOCKER_REGISTRY)/$(PROJECT_NAME):$(VERSION) \
		-f $(DOCKERFILE) .

docker-push: docker-build
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_REGISTRY)/$(PROJECT_NAME):$(VERSION)

# Development environment
dev-up: dev-down
	@echo "Starting development environment..."
	@$(DOCKER_COMPOSE) up -d --build

dev-down:
	@echo "Stopping development environment..."
	@$(DOCKER_COMPOSE) down

# Swagger documentation
swagger-init:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go

postman-merge:
	@npx postman-combine-collections --name $(COLLECTION_NAME) -f '$(COLLECTION_DIR)/$(COLLECTION_FILES)' -o $(COLLECTION_DIR)/$(COLLECTION_OUTPUT_FILE)

# Help target
help:
	@echo "Itsrama Portfolio Backend - Makefile Help"
	@echo ""
	@echo "Project Management:"
	@echo "  all           - Clean, install deps, lint, test, and build the project"
	@echo "  clean         - Remove build artifacts and clean cache"
	@echo "  deps          - Install and download project dependencies"
	@echo ""
	@echo "Development Workflow:"
	@echo "  lint          - Run code quality checks with golangci-lint"
	@echo "  test          - Execute all unit and integration tests"
	@echo "  build         - Compile the application binary"
	@echo "  run           - Start the application locally"
	@echo ""
	@echo "Database Management:"
	@echo "  migrate-create- Interactively create a new database migration"
	@echo "  migrate-up    - Apply pending database migrations"
	@echo "  migrate-down  - Revert last applied database migrations"
	@echo ""
	@echo "Docker & Deployment:"
	@echo "  docker-build  - Build Docker image for the application"
	@echo "  docker-push   - Push Docker image to registry"
	@echo "  dev-up        - Start development environment containers"
	@echo "  dev-down      - Stop development environment containers"
	@echo ""
	@echo "Documentation:"
	@echo "  swagger-init - Generate Swagger API documentation"
	@echo "  postman-merge - Merge Postman API collections"
