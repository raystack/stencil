NAME="github.com/odpf/stencil"
VERSION=$(shell git describe --always --tags 2>/dev/null)
PROTON_COMMIT := "d08a4916075438ed9c3a4510989bda2273e0a4e4"

.PHONY: all build test clean dist vet proto install

all: build

build: ## Build the stencil binary
	go build -ldflags "-X config.Version=${VERSION}" ${NAME}

test: ## Run the tests
	go test ./... -coverprofile=coverage.out

coverage: ## Print code coverage
	go test -race -coverprofile coverage.txt -covermode=atomic ./... & go tool cover -html=coverage.out

vet: ## Run the go vet tool
	go vet ./...

proto: ## Generate the protobuf files
	@echo " > generating protobuf from odpf/proton"
	@echo " > [info] make sure correct version of dependencies are installed using 'make install'"
	@buf generate https://github.com/odpf/proton/archive/${PROTON_COMMIT}.zip#strip_components=1 --template buf.gen.yaml --path odpf/stencil
	@echo " > protobuf compilation finished"

clean: ## Clean the build artifacts
	rm -rf stencil dist/

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'