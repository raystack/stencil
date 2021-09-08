NAME="github.com/odpf/stencil"
VERSION=$(shell git describe --always --tags 2>/dev/null)
COVERFILE="/tmp/stencil.coverprofile"
DIST_PATH=out

.PHONY: all build clean

all: build

build:
	go build -ldflags "-X main.Version=${VERSION}" ${NAME}

clean:
	rm -rf stencil dist/

lint:
	go vet ./...

test:
	go test ./... -coverprofile=coverage.out

test-coverage: test
	go tool cover -html=coverage.out

dist:
	@bash ./scripts/build.sh
