.PHONY: tidy run test lint

tidy:
	go mod tidy

run:
	SCAVO_ENV=local SCAVO_HTTP_ADDR=:8080 go run ./cmd/scavo-server

test:
	go test ./...

lint:
	golangci-lint run

