.PHONY: up

include .env

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

server:
	go run ./cmd/server/main.go
db-up:
	docker compose up -d
db-stop:
	docker compose stop
migration:
	@bash scripts/create_migration.sh $(name)
migrate-up:
	go run cmd/migrations/main.go up
migrate-down:
	go run cmd/migrations/main.go down
migrate-rollback:
	go run cmd/migrations/main.go rollback
build-server:
	go mod download; \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o server cmd/server/main.go
test:
	go test -race ./...

