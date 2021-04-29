BUILD_FILES =  $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{printf "%s\n" .}}{{end}}' ./...)

all: build

build: bin/linebot2lambda

bin/linebot-to-lambda: $(BUILD_FILES)
	@go build -trimpath -o "$@" .

lint:
	@echo "golint running..."
	@golint ./...
	@echo "go vet running..."
	@go vet ./...
