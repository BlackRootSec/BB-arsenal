BINARY  := bb-arsenal
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test vet fmt check install clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

# Same checks CI runs: formatting, vet, tests.
check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed on:"; gofmt -l .; exit 1)
	go vet ./...
	go test ./...

install:
	go install -ldflags "$(LDFLAGS)" .

clean:
	rm -f $(BINARY)
	rm -rf dist
