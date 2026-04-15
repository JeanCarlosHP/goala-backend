MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_DIR := $(dir $(MAKEFILE_PATH))

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make docker-up    - Start Docker services"
	@echo "  make docker-down  - Stop Docker services"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make sqlc         - Generate code from SQL queries"
	@echo "  make import-openfoodfacts - Import Open Food Facts dump into local catalog"
	@echo "  make clean        - Clean build artifacts"

.PHONY: run
run:
	air -c .air.toml

.PHONY: build
build:
	go build -o bin/api cmd/api/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-unit
test-unit:
	go test -v -short ./internal/services/... ./internal/repositories/...

.PHONY: test-integration
test-integration:
	go test -v -run Integration ./internal/handlers/...

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: docker-build
docker-build:
	docker-compose up --build -d

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/calorie_ai?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/calorie_ai?sslmode=disable" down

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: clean
clean:
	rm -rf bin/
	go clean

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: format
format:
	go fmt ./...

.PHONY: import-openfoodfacts
import-openfoodfacts:
	mkdir -p .gocache
	GOCACHE=$(PROJECT_DIR).gocache \
	go run ./cmd/import-openfoodfacts \
		-file ../en.openfoodfacts.org.products.csv.gz \
		-batch-size $${BATCH_SIZE:-2000} \
		-progress-every $${PROGRESS_EVERY:-10000} \
		$${LIMIT:+-limit $$LIMIT} \
		$${TRUNCATE:+-truncate} \
		$${DRY_RUN:+-dry-run} \
		$${INDEX_MEILI:+-index-meili} \
		$${DATABASE_URL:+-database-url $$DATABASE_URL} \
		$${MEILISEARCH_URL:+-meili-url $$MEILISEARCH_URL} \
		$${MEILISEARCH_API_KEY:+-meili-api-key $$MEILISEARCH_API_KEY} \
		$${MEILISEARCH_FOODS_INDEX:+-meili-index $$MEILISEARCH_FOODS_INDEX} \
		$${MEILISEARCH_TIMEOUT:+-meili-timeout $$MEILISEARCH_TIMEOUT} \
		$${MEILISEARCH_RETRIES:+-meili-retries $$MEILISEARCH_RETRIES} \
		$${MEILISEARCH_RETRY_BACKOFF:+-meili-backoff $$MEILISEARCH_RETRY_BACKOFF}
