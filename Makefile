GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=createp

all: build

build: 
	cd src && $(GOBUILD) -o $(BINARY_NAME)

test: 
	$(GOTEST) -v ./...

clean: 
	cd src && $(GOCLEAN)
	cd src && rm -f $(BINARY_NAME)
	cd src && rm -f $(BINARY_UNIX)

run:
	cd src && $(GOBUILD) -o $(BINARY_NAME)
	./src/$(BINARY_NAME)

mod:
	$(GOMOD) tidy

.PHONY: all build test clean run mod