ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: test, lint

test:
	go test -cover -race ./... -v

lint:
	golangci-lint run

proto-gen:
	@protoc -I=./internal/proto --go_out=./internal/ --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./internal ./internal/proto/*.proto

proto-clear:
	@rm -r ./internal/genproto