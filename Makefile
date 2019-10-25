BIN_DIR  ?= ./bin
PKG_NAME ?= event-forwarder
LDFLAGS  ?= "-s -w"

COVERAGE_PROFILE ?= coverage.out

AWS_PROFILE ?= sandbox

GOTOOLS := \
golang.org/x/tools/cmd/cover \

default: build

.PHONY: build
build:
	@echo "---> Building"
	go build -o ./bin/$(PKG_NAME) -ldflags $(LDFLAGS) ./cmd

.PHONY: clean
clean:
	@echo "---> Cleaning"
	rm -rf $(BIN_DIR)

.PHONY: enforce
enforce:
	@echo "---> Enforcing coverage"
	./scripts/coverage.sh $(COVERAGE_PROFILE)

.PHONY: html
html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

.PHONY: install
install:
	@echo "---> Installing dependencies"
	go mod download

.PHONY: lint
lint:
	@echo "---> Linting"
	$(BIN_DIR)/golangci-lint run

.PHONY: setup
setup:
	@echo "--> Installing development tools"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.18.0
	go get -u $(GOTOOLS)

.PHONY: test
test:
	@echo "---> Testing"
	AWS_PROFILE=$(AWS_PROFILE) ENVIRONMENT=test go test ./pkg/... -race -coverprofile $(COVERAGE_PROFILE)

.PHONY: test_ci
test_ci:
	@echo "---> Testing"
	ENVIRONMENT=test go test ./pkg/... -race -coverprofile $(COVERAGE_PROFILE)
