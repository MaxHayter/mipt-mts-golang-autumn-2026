.PHONY: generate test check run-service run-api

generate:
	mkdir -p gen
	protoc -I api --go_out=gen --go_opt=paths=source_relative --go-grpc_out=gen --go-grpc_opt=paths=source_relative api/incident/v1/incident.proto

test:
	go test -race ./...

# Сначала запустите оба приложения в отдельных терминалах.
check:
	go test -race -tags=integration -count=1 -timeout=90s ./apitest

run-service:
	go run -race ./cmd/incident-service

run-api:
	go run -race ./cmd/incident-api
