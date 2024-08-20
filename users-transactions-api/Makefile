ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifeq (${ENV},dev)
	DB_HOST=localhost
endif

.PHONY: test, lint, run

GOOSE=goose
DB_URL=postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

test:
	go test -cover -race ./... -v

lint:
	golangci-lint run

run:
	go run cmd/api/main.go

migrate-status:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" status

migrate-reset:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" reset

migrate-up:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	$(GOOSE) -dir ./migrations postgres "$(DB_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	$(GOOSE) -dir ./migrations create $$name sql

clean-migrations:
	@rm ./migrations/*.sql