APP := webhix
MAIN := ./cmd/webhix
BUILD_DIR := bin
COVER_DIR := coverage
TOOLS_DIR := .tools/bin
GOLANGCI_LINT_VERSION ?= v2.12.2

GO ?= go
GOLANGCI_LINT ?= $(TOOLS_DIR)/golangci-lint

GOFLAGS ?=
PKGS := ./...
LDFLAGS ?= -s -w

.PHONY: deps tidy tidy-check fmt fmt-check vet lint lint-install test test-race cover check ci build build-prod run clean

# Download Go module dependencies.
deps:
	$(GO) mod download

# Clean up go.mod and go.sum.
tidy:
	$(GO) mod tidy

# Verify that go.mod and go.sum are already tidy.
tidy-check:
	$(GO) mod tidy
	git diff --exit-code -- go.mod go.sum

# Format all Go packages.
fmt:
	$(GO) fmt $(PKGS)

# Verify that Go files are already formatted.
fmt-check:
	$(GO) fmt $(PKGS)
	git diff --exit-code -- '*.go'

# Run Go's built-in static analysis.
vet:
	$(GO) vet $(PKGS)

# Run golangci-lint with the project config.
lint:
	$(GOLANGCI_LINT) run $(PKGS)

# Install the pinned golangci-lint version.
lint-install:
	mkdir -p $(TOOLS_DIR)
	GOBIN=$(CURDIR)/$(TOOLS_DIR) $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# Run unit tests.
test:
	$(GO) test $(GOFLAGS) $(PKGS)

# Run tests with the race detector enabled.
test-race:
	$(GO) test -race $(GOFLAGS) $(PKGS)

# Generate an HTML coverage report.
cover:
	mkdir -p $(COVER_DIR)
	$(GO) test $(GOFLAGS) -coverprofile=$(COVER_DIR)/coverage.out $(PKGS)
	$(GO) tool cover -html=$(COVER_DIR)/coverage.out -o $(COVER_DIR)/coverage.html

# Run the standard local quality gate.
check: tidy fmt vet lint test

# Run the full CI-style pipeline.
ci: deps tidy-check fmt-check vet lint-install lint test build-prod

# Build a local development binary.
build:
	mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP) $(MAIN)

# Build a smaller production binary.
build-prod:
	mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -trimpath -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP) $(MAIN)

# Run the application from source.
run:
	$(GO) run $(MAIN)

# Remove build and coverage artifacts.
clean:
	$(GO) clean
	rm -rf $(BUILD_DIR) $(COVER_DIR) .tools
