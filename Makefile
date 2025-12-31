.PHONY: help test lint fmt vet tidy clean build run

# --- Project Variables ---
BINARY_NAME := gotraffic
CMD_API_PATH := ./cmd/api/main.go

# --- Build Variables ---
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME := $(shell date +%FT%T%z)

LDFLAGS := -ldflags "-s -w -X main.CommitHash=$(COMMIT_HASH) -X main.BuildTime=$(BUILD_TIME)"
GO_BUILD_FLAGS := -trimpath

# --- Docker Variables ---
DOCKER_TAG := $(COMMIT_HASH)

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1mUsage:\033[0m\n  make \033[36m<target>\033[0m\n"} \
	/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } \
	/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Env

env: ## Create .env file from example
	@cp -n .env.example .env || true

##@ Development

mocks: ## Generate all mocks
	@echo "Generating mocks..."
	@go generate ./...

install-swag:
	@if [ ! -f $(SWAG_BIN) ]; then \
		@echo "Installing swag..."; \
		GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@latest; \
	fi

SWAG_CMD := go run github.com/swaggo/swag/cmd/swag@latest
swagger: install-swag ## Generate Swagger documentation
	@echo "Generating swagger docs..."
	@$(SWAG_CMD) init -g $(CMD_API_PATH) --output docs --parseDependency --parseInternal

##@ Run (Backend)

run-local: ## Run app in local mode
	@go run $(CMD_API_PATH) --env=local --env-path=.env

run-dev: ## Run app in dev mode
	@go run $(CMD_API_PATH) --env=dev --env-path=.env

run-prod: ## Run app in prod mode
	@go run $(CMD_API_PATH) --env=prod --env-path=.env

##@ Run (Frontend)

run-frontend: ## Runs frontend application
	cd frontend && yarn dev --host

run-frontend-docker: ## Runs frontend application in Docker
	docker-compose -f docker-compose.frontend.yml up --build

##@ Builds

build: swagger ## Build binary for current OS
	@echo "Building $(BINARY_NAME)..."
	@CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) $(LDFLAGS) -o build/$(BINARY_NAME) $(CMD_API_PATH)
	@echo "Build complete: build/$(BINARY_NAME)"

build-linux: swagger ## Build binary for Linux (AMD64)
	@echo "Building Linux binary..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) $(LDFLAGS) -o build/$(BINARY_NAME)-linux $(CMD_API_PATH)
	@echo "Build complete: build/$(BINARY_NAME)-linux"

build-windows: swagger ## Build binary for Windows (AMD64)
	@echo "Building Windows binary..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) $(LDFLAGS) -o build/$(BINARY_NAME)-windows.exe $(CMD_API_PATH)
	@echo "Build complete: build/$(BINARY_NAME)-windows.exe"

clean: ## Remove build artifacts
	@rm -rf build/

##@ Tests

test: ## Runs tests
	@go test -v ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

lint-install: ## Install golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

##@ Quality Control

format: ## Format code
	@go fmt ./...

check: format lint test ## Run all checks before commit
