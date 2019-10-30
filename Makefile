BIN_DIR  ?= ./bin
PKG_NAME ?= event-forwarder
LDFLAGS  ?= "-s -w"

COVERAGE_PROFILE ?= coverage.out

AWS_PROFILE ?= sandbox
AWS_REGION ?= us-west-2

GOTOOLS := \
golang.org/x/tools/cmd/cover \

default: build

.PHONY: build
build:
	@echo "---> Building"
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(PKG_NAME) -ldflags $(LDFLAGS) ./cmd/main.go

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
setup: $(BIN_DIR)/golangci-lint
	@echo "--> Installing development tools"
	go get -u $(GOTOOLS)

$(BIN_DIR)/golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.18.0

.PHONY: test
test:
	@echo "---> Testing"
	AWS_PROFILE=$(AWS_PROFILE) AWS_REGION=$(AWS_REGION) ENVIRONMENT=test go test ./pkg/... -race -coverprofile $(COVERAGE_PROFILE)

.PHONY: test_ci
test_ci:
	@echo "---> Testing"
	AWS_REGION=$(AWS_REGION) ENVIRONMENT=test go test ./pkg/... -race -coverprofile $(COVERAGE_PROFILE)
