PROG_NAME := "./bin/sys_monitoring"

build:
	go build -v -o $(PROG_NAME) ./cmd

run: build
	$(PROG_NAME) -config ./configs/config.toml

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	$(shell go env GOPATH)/bin/golangci-lint run ./...

generate:
	protoc api/server.proto --go_out=./internal/server/pb --go-grpc_out=./internal/server/pb
