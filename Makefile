.PHONY: up

include .env

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

server:
	go run ./cmd/server/main.go
frontend:
	pnpm -C ./app install && pnpm -C ./app run dev
db-up:
	docker compose up -d
db-stop:
	docker compose stop
migrate:
	go run ./cmd/migrate/main.go
create-migration:
	@bash scripts/create_migration.sh $(name)
migrate-up:
	go run ./cmd/migrate/main.go up
migrate-down:
	go run ./cmd/migrate/main.go down
migrate-rollback:
	go run ./cmd/migrate/main.go rollback
test:
	go test -race ./...
	## build-server: build the server
build-server:
	export GO111MODULE="on"; \
		go mod download; \
		go mod vendor; \
		CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=amd64 \
		go build -o bin/server cmd/server/main.go
build-migrate:
	export GO111MODULE="on"; \
		go mod download; \
		go mod vendor; \
		CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=amd64 \
		go build -o bin/migrate cmd/migrate/main.go
