.PHONY: run build clean

GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=createp
ROOT=cd /home/crative/dev/go/project-creator-go

run:
	cd src && $(GOBUILD) -o $(BINARY_NAME)
	./src/$(BINARY_NAME)

build:
	$(ROOT) && nix build

clean:
	$(ROOT) &&	rm -f src/$(BINARY_NAME)
	$(ROOT) &&	rm -rf result
