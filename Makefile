# Makefile

run-generator: ## запуск генератора с контекстом
	go run ./generator/simple

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

GOLANGCI_LINT_VERSION = v2.8.0

install-linter: ## установка линтера
	@echo "Installing golangci-lint... "
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $$(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

lint:
	@echo "Linting..."
	golangci-lint run ./...

fix:
	@echo "Fixing..."
	golangci-lint run --fix ./...

.PHONY: \
	run-generator \
	help \
	lint \
	fix \
	install-linter