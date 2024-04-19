NAME="github.com/goto/stencil"
VERSION=$(shell git describe --always --tags 2>/dev/null)
PROTON_COMMIT := "4acac160b663d97d64cab017dcb08c787eec3e1d"
EXCLUDE_FILES :=./test_utils/testutils.go
.PHONY: all build test clean dist vet proto install ui

all: build

build: ui ## Build the stencil binary
	go build -ldflags "-X config.Version=${VERSION}" ${NAME}

test: ui ## Run the tests
	go test ./... -coverprofile=coverage.out

coverage: ui ## Print code coverage
	go test -race -coverprofile coverage.txt -covermode=atomic ./... & go tool cover -html=coverage.out

coverage-exclude:
	# Run tests and generate coverage profile
	go test -race -coverprofile coverage.txt -covermode=atomic ./...
	# Filter out coverage data for specified files
	@for exclude_file in $(EXCLUDE_FILES); do \
		grep -v "$$exclude_file" coverage.out > coverage.tmp && mv coverage.tmp coverage.out; \
	done
	go tool cover -html=coverage.out

vet: ## Run the go vet tool
	go vet ./...

lint: ## Run golang-ci lint 
	golangci-lint run

proto: ## Generate the protobuf files
	@echo " > generating protobuf stub files"
	@echo " > [info] make sure correct version of dependencies are installed using 'make install'"
	@buf generate https://github.com/goto/proton/archive/${PROTON_COMMIT}.zip#strip_components=1 --template buf.gen.yaml --path gotocompany/stencil
	@echo " > protobuf compilation finished"

clean: ## Clean the build artifacts
	rm -rf stencil dist/ ui/build/

ui:
	@echo " > generating ui build"
	@cd ui && $(MAKE) dep && $(MAKE) dist

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
