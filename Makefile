# Makefile

simple-generator: ## Генератор возвращает канал, останавливается контекстом
	go run ./generator/simple

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: \
	simple-generator \
	help
