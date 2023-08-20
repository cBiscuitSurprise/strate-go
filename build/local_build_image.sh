#!/bin/sh
# build the docker image
# we first need to copy the binary from the symlink path to the local path
#
# Args:
#   1 - docker image tag

BINARY=bazel-bin/strate-go_/strate-go
BINARY_COPY=.build/strate-go

mkdir -p $(dirname $BINARY_COPY)
cp -f $BINARY $BINARY_COPY

docker build -f build/prod.Dockerfile -t cbiscuit87/strate-go:$1 --build-arg="BINARY=$BINARY_COPY" .
