# Makefile for API Gateway Service
# Aligned with go-infrastructure architectural standards

# Variables
APP_NAME := api-gateway
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := 1.25.5

# Docker variables
DOCKER_REGISTRY ?= ghcr.io
DOCKER_IMAGE := $(DOCKER_REGISTRY)/vhvplatform/$(APP_NAME)
DOCKER_TAG ?= $(VERSION)

# Go build variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOFMT := $(GOCMD) fmt

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all
all: clean fmt vet test build

.PHONY: help
help: ## Display this help screen
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(GOBIN)
	$(GOBUILD) $(LDFLAGS) -o $(GOBIN)/$(APP_NAME) ./cmd/main.go
	@echo "Build complete: $(GOBIN)/$(APP_NAME)"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(GOBIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(GOBIN)/$(APP_NAME)-linux ./cmd/main.go
	@echo "Build complete: $(GOBIN)/$(APP_NAME)-linux"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(GOBIN)
	@rm -f coverage.txt coverage.html
	@echo "Clean complete"

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	@$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./... 2>&1 | tee test-output.log | grep -v "no such tool" || true
	@if grep -q "FAIL" test-output.log; then \
		echo "❌ Tests failed"; \
		rm -f test-output.log; \
		exit 1; \
	fi
	@rm -f test-output.log
	@echo "Tests complete"

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-coverage-check
test-coverage-check: test ## Check test coverage meets minimum threshold (80%)
	@echo "Checking coverage threshold..."
	@coverage=$$(go tool cover -func=coverage.txt | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$coverage" ]; then \
		echo "❌ Unable to calculate coverage"; \
		exit 1; \
	fi; \
	if [ $$(echo "$$coverage" | awk '{print ($$1 < 80)}') -eq 1 ]; then \
		echo "❌ Coverage $$coverage% is below 80%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$coverage% meets threshold"; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Format complete"

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run --timeout=5m ./...
	@echo "Lint complete"

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "Dependencies downloaded"

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "Dependencies updated"

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	$(GOMOD) verify
	@echo "Dependencies verified"

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

.PHONY: docker-build-no-cache
docker-build-no-cache: ## Build Docker image without cache
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG) without cache..."
	docker build --no-cache -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	@echo "Pushing Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Docker image pushed"

.PHONY: docker-run
docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	@if [ -f .env.local ]; then \
		docker run -p 8080:8080 --env-file .env.local $(DOCKER_IMAGE):$(DOCKER_TAG); \
	elif [ -f .env ]; then \
		docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG); \
	else \
		echo "Warning: No .env or .env.local file found. Running without environment variables."; \
		docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG); \
	fi

.PHONY: run
run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run ./cmd/main.go

.PHONY: install
install: build ## Install the application
	@echo "Installing $(APP_NAME) to $(GOPATH)/bin..."
	@cp $(GOBIN)/$(APP_NAME) $(GOPATH)/bin/
	@echo "Install complete"

.PHONY: security-scan
security-scan: ## Run security scan with gosec
	@echo "Running security scan..."
	@if ! command -v gosec &> /dev/null; then \
		echo "gosec not found. Installing..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...
	@echo "Security scan complete"

.PHONY: validate
validate: fmt vet lint test-coverage-check ## Run all validation checks
	@echo "✅ All validation checks passed"

.PHONY: ci
ci: validate docker-build ## Run CI pipeline locally
	@echo "✅ CI pipeline complete"

.PHONY: version
version: ## Display version information
	@echo "Application: $(APP_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"
