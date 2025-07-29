.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	go build -o bin/conreq ./cmd

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
	go run ./cmd $(ARGS)

.PHONY: install
install: ## Install conreq to GOPATH/bin
	go install ./cmd

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t conreq:latest .

.PHONY: all
all: clean deps fmt lint test build ## Run all tasks