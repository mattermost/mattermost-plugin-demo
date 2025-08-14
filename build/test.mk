# ====================================================================================
# Testing and Quality Assurance
# ====================================================================================

GOLANGCI_LINT_BINARY = ./build/bin/golangci-lint
GOTESTSUM_BINARY = ./build/bin/gotestsum

## Install go tools
install-go-tools:
	@echo "Installing development tools..."
	@pluginctl tools install --bin-dir ./build/bin

## Runs eslint and golangci-lint
.PHONY: check-style
check-style: manifest-check apply webapp/node_modules install-go-tools
	@echo Checking for style guide compliance

ifneq ($(HAS_WEBAPP),)
	cd webapp && npm run lint
	cd webapp && npm run check-types
endif

# It's highly recommended to run go-vet first
# to find potential compile errors that could introduce
# weird reports at golangci-lint step
ifneq ($(HAS_SERVER),)
	@echo Running golangci-lint
	$(GO) vet ./...
	$(GOLANGCI_LINT_BINARY) run ./...
endif

## Runs any lints and unit tests defined for the server and webapp, if they exist.
.PHONY: test
test: apply webapp/node_modules install-go-tools
ifneq ($(HAS_SERVER),)
	$(GOTESTSUM_BINARY) -- -v ./...
endif
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) run test;
endif

## Runs any lints and unit tests defined for the server and webapp, if they exist, optimized
## for a CI environment.
.PHONY: test-ci
test-ci: apply webapp/node_modules install-go-tools
ifneq ($(HAS_SERVER),)
	$(GOTESTSUM_BINARY) --format standard-verbose --junitfile report.xml -- ./...
endif
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) run test;
endif

## Creates a coverage report for the server code.
.PHONY: coverage
coverage: apply webapp/node_modules
ifneq ($(HAS_SERVER),)
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=server/coverage.txt ./server/...
	$(GO) tool cover -html=server/coverage.txt
endif
