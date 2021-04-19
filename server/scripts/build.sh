#!/bin/bash
NAME="github.com/odpf/stencil/server"
VERSION=$(git describe --always --tags 2>/dev/null)

SYS="linux"
ARCH="amd64"
BUILD_DIR="dist"

build() {
    EXECUTABLE_NAME=$1
    LD_FLAGS=$2
    # create a folder named via the combination of os and arch
    TARGET="./$BUILD_DIR/${SYS}-${ARCH}"
    mkdir -p $TARGET

    # place the executable within that folder
    executable="${TARGET}/$EXECUTABLE_NAME"
    echo $executable
    GOOS=$SYS GOARCH=$ARCH go build -ldflags "$LD_FLAGS" -o $executable $NAME
}

build stencil "-X main.Version=${VERSION}"
