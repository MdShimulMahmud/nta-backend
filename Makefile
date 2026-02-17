.PHONY: help build run test clean lint format swagger docker-build docker-run install

# Variables
APP_NAME=nta-backend
BUILD_DIR=bin
MAIN_PATH=./cmd/server
DOCKER_IMAGE=nta-backend:latest

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@go install github.com/swaggo/swag/cmd/swag@latest

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)/main.go

dev: ## Run in development mode with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	@air

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

coverage: test ## Run tests with coverage
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run ./...

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

swagger: ## Generate swagger documentation
	@echo "Generating swagger docs..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf docs
	@rm -f coverage.txt coverage.html

docker-build: ## Build docker image
	@echo "Building docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run docker container
	@echo "Running docker container..."
	@docker run -p 8080:8080 --privileged -v /sys/class/net:/sys/class/net:ro $(DOCKER_IMAGE)

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

check: lint test ## Run linter and tests

all: clean install swagger build ## Clean, install, generate swagger and build

.DEFAULT_GOAL := help
