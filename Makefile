.PHONY: run

.DEFAULT_GOAL=run

LOCAL_BIN=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3


#by default is checked using pre-commit
lint:
	$(LOCAL_BIN)/golangci-lint run --fast --config .golangci.pipeline.yaml

run:
	go run cmd/app/main.go
