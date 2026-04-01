.PHONY: build test clean install lint

BINARY_NAME=ep
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) ./cmd

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
