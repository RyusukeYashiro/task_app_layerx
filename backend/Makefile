.SILENT:

# コード生成
gen-oapi:
	oapi-codegen -generate types        -package api -o internal/presentation/http/echo/types.gen.go  api/openapi.yaml
	oapi-codegen -generate echo-server  -package api -o internal/presentation/http/echo/server.gen.go api/openapi.yaml
	@echo "openapi code generated successfully"

# マイグレーション
migrate-up:
	docker compose exec -T db sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE' < migrations/0001_init.up.sql
	@echo "migration up completed"

migrate-down:
	docker compose exec -T db sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE' < migrations/0001_init.down.sql
	@echo "migration down completed"

migrate-status:
	docker compose exec db sh -c 'mysql -u$$MYSQL_USER -p$$MYSQL_PASSWORD $$MYSQL_DATABASE -e "SHOW TABLES;"'

# サーバー起動
run:
	@if [ -f .env ]; then export $$(cat .env | grep -v '^#' | xargs) && go run ./cmd/api; else go run ./cmd/api; fi

# Lint
lint:
	golangci-lint run

# テスト
test:
	go test ./... -race -cover

# Docker
build:
	docker compose build

rebuild:
	docker compose build --no-cache

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f app

ps:
	docker compose ps

bash:
	docker compose exec app sh

# クリーンアップ
clean:
	rm -rf tmp dist
	go clean -cache -testcache -modcache

# ヘルプ
help:
	@echo "Available commands:"
	@echo "  gen-oapi      - Generate OpenAPI code"
	@echo "  run           - Run application locally"
	@echo "  lint          - Run golangci-lint"
	@echo "  test          - Run tests with race detector"
	@echo "  build         - Build docker images"
	@echo "  up            - Start containers"
	@echo "  down          - Stop containers"
	@echo "  logs          - Show app logs"
	@echo "  ps            - Show running containers"
	@echo "  bash          - Open shell in app container"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  migrate-status - Show current database tables"
	@echo "  clean         - Clean build artifacts and caches"

.PHONY: gen-oapi run lint test build up down logs ps bash migrate-up migrate-down migrate-status clean help
