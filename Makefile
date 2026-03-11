BINARY_NAME=entire-local
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=${VERSION}"

RED=\033[0;31m
GREEN=\033[0;32m
NC=\033[0m

DEMOS=menu status explain doctor rewind resume reset clean

.PHONY: all build clean test install fmt vet lint check deps demo $(addprefix demo-,$(DEMOS)) help

help:
	@echo "entire-local Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build    Build the binary"
	@echo "  test     Run tests"
	@echo "  lint     Run golangci-lint"
	@echo "  check    Run fmt, vet, lint"
	@echo "  install  Install to ~/go/bin"
	@echo "  clean    Remove build artifacts"
	@echo "  deps     Download dependencies"
	@echo "  demo          Record all demo GIFs"
	@echo "  demo-<name>   Record one GIF (menu status explain doctor rewind resume reset clean)"
	@echo "  fmt      Format code"
	@echo "  vet      Run go vet"

all: build

build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${BINARY_NAME} .
	@echo "${GREEN}✓${NC} Build complete: ${BINARY_NAME}"

clean:
	@go clean
	@rm -f ${BINARY_NAME}
	@rm -rf dist/
	@rm -f coverage.txt coverage.html

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "${GREEN}✓${NC} Tests complete"

install: build
	@echo "Installing ${BINARY_NAME}..."
	@go install ${LDFLAGS} .
	@echo "${GREEN}✓${NC} Installed to ~/go/bin"

fmt:
	@go fmt ./...

vet:
	@go vet ./...

lint:
	@golangci-lint run

check: fmt vet lint
	@echo "${GREEN}✓${NC} All checks passed"

demo: build $(addprefix demo-,$(DEMOS))
	@echo "${GREEN}✓${NC} All demos recorded in demo/"

$(addprefix demo-,$(DEMOS)): demo-%: build
	@echo "Recording $*..."
	@vhs demo/$*.tape
	@echo "${GREEN}✓${NC} demo/$*.gif"

deps:
	@go mod download
	@go mod tidy

.DEFAULT_GOAL := help
