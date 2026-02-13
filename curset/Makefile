BINARY_NAME := curset
VERSION ?= 0.1.0
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build install clean tidy

build: tidy
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

install: build
	sudo cp bin/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed to /usr/local/bin/$(BINARY_NAME)"

clean:
	rm -rf bin/

tidy:
	go mod tidy
