ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: test, lint, run, sqlc

GOOSE=goose
DB_HOST=localhost
DB_URL=postgres://${DB_USERNAME}:${DB_PASSWORD}@$(DB_HOST):${DB_PORT}/${DB_NAME}?sslmode=disable

test:
	go test -cover -race ./... -v

lint:
	golangci-lint run

run:
	go run cmd/main.go

sqlc:
	sqlc generate

proto-gen:
	@protoc -I=./internal/proto --go_out=./internal/ --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./internal ./internal/proto/*.proto

proto-clear:
	@rm -r ./internal/genproto

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