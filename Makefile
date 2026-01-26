.PHONY: help run test test-watch clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the server (with fresh database)
	@rm -rf api/data/pos.db
	@cd api && go run cmd/server/main.go

test: ## Run Bruno API tests (requires server to be running)
	@cd api/bruno && bru run --env local --tags entities

test-watch: ## Run Bruno tests in watch mode
	@cd api/bruno && bru run --env local --tags entities --watch

clean: ## Remove generated files and databases
	@rm -rf api/data/pos.db
	@echo "Cleaned database files"
