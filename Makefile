# Makefile for lcss

BINARY := bin/lcss
CMD := ./cmd/lcss

# Version information - can be overridden by environment (e.g. CI or manual invocations)
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

.PHONY: build build-release clean test

build:
	@mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

# Example usage: make build-release GOOS=linux GOARCH=amd64
build-release:
	@mkdir -p dist/$(GOOS)_$(GOARCH)
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
	VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) \
	go build -ldflags "$(LDFLAGS)" -o dist/$(GOOS)_$(GOARCH)/lcss $(CMD)

test:
	go test ./...

clean:
	rm -rf bin dist
