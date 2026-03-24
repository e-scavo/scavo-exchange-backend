.PHONY: run up down migrate migrate-down migrate-status test test-unit test-integration smoke-login

run:
	go run ./cmd/scavo-server

up:
	docker compose up -d

down:
	docker compose down

migrate:
	SCAVO_POSTGRES_URL=$$SCAVO_POSTGRES_URL ./scripts/migrate.sh up

migrate-down:
	SCAVO_POSTGRES_URL=$$SCAVO_POSTGRES_URL ./scripts/migrate.sh down

migrate-status:
	SCAVO_POSTGRES_URL=$$SCAVO_POSTGRES_URL ./scripts/migrate.sh status

test:
	go test ./...

test-unit:
	go test ./internal/...

test-integration:
	SCAVO_TEST_POSTGRES_URL=$$SCAVO_TEST_POSTGRES_URL go test ./internal/modules/user -run TestPostgresRepository_UpsertDevUser -v

smoke-login:
	./scripts/smoke_login.sh