.PHONY: run migrate seed tidy test docker-build smoke-test-live

BASE_URL ?= https://cpe.geller.men

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

smoke-test-live:
	BASE_URL=$(BASE_URL) ./infra/ansible/smoke-test-live-mvp.sh
