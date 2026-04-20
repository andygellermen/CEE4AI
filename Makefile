.PHONY: run migrate seed tidy test

run:
	go run ./cmd/api

migrate:
	go run ./cmd/migrate

seed:
	go run ./cmd/seed

tidy:
	go mod tidy

test:
	go test ./...
