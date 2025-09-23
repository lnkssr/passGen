# Usage: make            -> build
#        make build
#        make run
#        make test
#        make fmt
#        make vet
#        make tidy
#        make clean
#        make install
#        make cross GOOS=linux GOARCH=amd64
#        make release

PROJECT := $(notdir $(CURDIR))
BINARY  ?= $(PROJECT)
BIN_DIR := bin

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.DEFAULT_GOAL := build

.PHONY: all build run test fmt vet tidy clean install cross release lint

all: build

build: $(BIN_DIR)/$(BINARY)

$(BIN_DIR)/$(BINARY):
	@mkdir -p $(BIN_DIR)
	@echo "go build -> $@ (GOOS=$(GOOS) GOARCH=$(GOARCH))"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '$(LDFLAGS)' -o $@ ./...

run: build
	@echo "running $(BIN_DIR)/$(BINARY)"
	./$(BIN_DIR)/$(BINARY)

test:
	@go test ./...

fmt:
	@go fmt ./...

vet:
	@go vet ./...

tidy:
	@go mod tidy

install:
	@echo "go install (ldflags)"
	go install -ldflags '$(LDFLAGS)' ./...

cross:
	@mkdir -p $(BIN_DIR)
	@echo "cross build -> $(BIN_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY)-$(GOOS)-$(GOARCH) ./...

release: clean
	@$(MAKE) build GOOS=linux GOARCH=amd64 BINARY=$(BINARY)-linux-amd64
	@tar -C $(BIN_DIR) -czf $(BINARY)-$(VERSION)-linux-amd64.tar.gz $(BINARY)-linux-amd64
	@echo "release: $(BINARY)-$(VERSION)-linux-amd64.tar.gz"

clean:
	@rm -rf $(BIN_DIR)
	@echo "cleaned $(BIN_DIR)"

lint:
	@golangci-lint run || true

