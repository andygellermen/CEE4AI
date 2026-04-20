.PHONY: run migrate seed tidy test docker-build

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

docker-build:
	docker build -t cee4ai:local .
