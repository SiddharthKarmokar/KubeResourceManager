.PHONY: run build test lint docker-build swagger docs

APP_NAME=k8s-opt
MAIN_API=cmd/api/main.go
MAIN_CLI=cmd/cli/main.go
BIN_DIR=bin

run:
	@echo "Starting API server..."
	@go run $(MAIN_API)

build:
	@echo "Building API and CLI..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME)-api $(MAIN_API)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_CLI)
	@echo "Build complete. Binaries are in $(BIN_DIR)/"

test:
	@echo "Running tests..."
	@go test -v ./... -cover

lint:
	@echo "Running linters..."
	@go vet ./...
	@if command -v golangci-lint >/dev/null; then golangci-lint run; else echo "golangci-lint not installed, skipping"; fi

docker-build:
	@echo "Building docker image..."
	@docker build -t $(APP_NAME):latest .

swagger:
	@echo "Generating swagger docs..."
	@swag init -g cmd/api/main.go -o swagger --parseDependency --parseInternal

docs:
	@echo "Building MkDocs..."
	@pip install -r docs/requirements.txt
	@mkdocs build --strict
