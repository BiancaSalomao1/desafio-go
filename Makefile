# ==========================
# Configurações
# ==========================

ifneq (,$(wildcard .env))
include .env
export
endif

MIGRATE=$(shell go env GOPATH)/bin/migrate

# ==========================
# Aplicação
# ==========================

run:
	go run ./cmd/api

build:
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

# ==========================
# Docker
# ==========================

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up -d

logs:
	docker compose logs -f

ps:
	docker ps

# ==========================
# Migrations
# ==========================

migrate-up:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" up

migrate-down:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" down

migrate-version:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" version

migrate-force:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(DATABASE_URL)" force 1