.PHONY: build build-release clean test install

BINARY_NAME := chroma-control
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

build-release:
	go build -ldflags "-s -w $(LDFLAGS)" -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)

test:
	go test ./...

install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/ 2>/dev/null || cp $(BINARY_NAME) ~/go/bin/ 2>/dev/null || echo "Install failed: GOPATH/bin not found"

.DEFAULT_GOAL := build
