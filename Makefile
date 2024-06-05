include .env

.PHONY: run

.DEFAULT_GOAL=run

LOCAL_BIN=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0


#by default is checked using pre-commit
lint:
	$(LOCAL_BIN)/golangci-lint run --fast --config .golangci.pipeline.yaml

run:
	go run cmd/app/main.go

migration-status:
	$(LOCAL_BIN)/migrate -source file://migrations -database $(DB_DSN) version

migration-up:
	$(LOCAL_BIN)/migrate -source file://migrations -database $(DB_DSN) up 1

migration-down:
	$(LOCAL_BIN)/migrate -source file://migrations -database $(DB_DSN) down 1