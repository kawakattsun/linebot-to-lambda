BUILD_FILES =  $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{printf "%s\n" .}}{{end}}' ./...)
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null)
REVISION ?= $(shell git rev-parse --short HEAD 2>/dev/null)

GO_LDFLAGS := -X main.version=$(VERSION)
GO_LDFLAGS += -X main.revision=$(REVISION)

all: build

build: bin/linebot2lambda

bin/linebot-to-lambda: $(BUILD_FILES)
	@go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" .

build-lambda: deploy/lambda/linebot-to-lambda/main

deploy/lambda/linebot-to-lambda/main: $(BUILD_FILES)
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" .

lint:
	@echo "golint running..."
	@golint ./...
	@echo "go vet running..."
	@go vet ./...
