NAME="github.com/odpf/stencil"
VERSION=$(shell git describe --always --tags 2>/dev/null)

.PHONY: all build clean

all: build

build:
	go build -ldflags "-X main.Version=${VERSION}" ${NAME}

clean:
	rm -rf stencil dist/

test:
	go test -count 1 -cover -race -timeout 1m ./...

dist:
	@bash ./scripts/build.sh
