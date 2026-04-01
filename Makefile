.PHONY: build test clean install lint

BINARY_NAME=ep
INSTALL_PATH=/usr/local/bin
VERSION ?= 0.1.0
LDFLAGS := -X main.version=$(VERSION)

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)

install: build
	mv $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

lint:
	golangci-lint run ./...

run:
	go run ./cmd $(ARGS)

.PHONY: version
version:
	@echo $(VERSION)
