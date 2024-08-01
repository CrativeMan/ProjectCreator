.PHONY: run build clean clean-all

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
	$(ROOT) && rm -rf test
	$(ROOT) && rm -rf asd

clean-all:
	make clean
	$(ROOT) &&	rm -rf result
