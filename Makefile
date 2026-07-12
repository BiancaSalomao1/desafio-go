# ==========================
# Configurações
# ==========================

ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: \
	run build fmt fmt-check vet test tidy clean check \
	up down restart logs ps docker-build \
	swagger \
	migrate-up migrate-down migrate-version migrate-force \
	test-api test-crud

# ==========================
# Aplicação
# ==========================

run:
	go run ./cmd/api

build:
	go build ./...

fmt:
	gofmt -w .

fmt-check:
	test -z "$$(gofmt -l .)"
	
vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

clean:
	go clean

check:
	$(MAKE) fmt-check
	$(MAKE) vet
	$(MAKE) test
	$(MAKE) build

# ==========================
# Swagger
# ==========================

swagger:
	swag init \
		-g cmd/api/main.go \
		--parseInternal \
		--parseDependency

# ==========================
# Docker
# ==========================

docker-build:
	docker build -t desafio-go .

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

# ==========================
# Scripts de Teste
# ==========================

test-api:
	./scripts/smoke_test.sh

test-crud:
	./scripts/crud_test.sh