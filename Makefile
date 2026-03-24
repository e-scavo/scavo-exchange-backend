.PHONY: tidy test link run up down migrate migrate-down

tidy:
	go mod tidy

test:
	go test ./...

lint:
	golangci-lint run

run:
	go run ./cmd/scavo-server

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	SCAVO_POSTGRES_URL=$$SCAVO_POSTGRES_URL ./scripts/migrate.sh up

migrate-down:
	SCAVO_POSTGRES_URL=$$SCAVO_POSTGRES_URL ./scripts/migrate.sh down