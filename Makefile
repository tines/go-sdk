default: lint

# Run unit tests
.PHONY: tests
tests:
	go test ./... -v $(TESTARGS) -cover -timeout 120m

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Install dependencies for local development
.PHONY: setup
setup:
	brew install golangci-lint

