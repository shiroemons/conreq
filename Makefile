.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	go build -ldflags "-s -w -X main.version=$$(git describe --tags --always) -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/conreq ./cmd/conreq

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -cover ./...

.PHONY: lint
lint: ## Run golangci-lint via Docker
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.3.0 golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix via Docker
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.3.0 golangci-lint run --fix

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/

.PHONY: deps
deps: ## Download dependencies
	go mod download

.PHONY: mod-tidy
mod-tidy: ## Tidy go.mod
	go mod tidy

.PHONY: run
run: ## Run the application (example: make run ARGS="https://httpbin.org/get")
	go run -ldflags "-s -w -X main.version=$$(git describe --tags --always) -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" ./cmd/conreq $(ARGS)

.PHONY: install
install: ## Install conreq to GOPATH/bin
	go install -ldflags "-s -w -X main.version=$$(git describe --tags --always) -X main.commit=$$(git rev-parse --short HEAD) -X main.date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" ./cmd/conreq

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t conreq:latest .

.PHONY: all
all: clean deps fmt lint test build ## Run all tasks

.PHONY: setup-git-hooks
setup-git-hooks: ## Setup git hooks for the project
	@scripts/setup-git-hooks.sh

.PHONY: pre-push
pre-push: ## Run all checks that would be run by pre-push hook
	@echo "📝 Checking code formatting..."
	@if ! go fmt ./... | grep -q .; then \
		echo "✓ Code formatting check passed"; \
	else \
		echo "✗ Code needs formatting"; \
		echo "Run 'go fmt ./...' to fix formatting issues"; \
		exit 1; \
	fi
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✓ go vet passed"
	@echo "🔍 Running golangci-lint..."
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v2.3.0 golangci-lint run
	@echo "✓ golangci-lint passed"
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅ All pre-push checks passed!"