#!/bin/bash

mkdir -p build/;

go get -d .
go install -i .
go mod download
go mod vendor
# GOPATH="$(pwd)/vendor:$(pwd)"
# GOBIN="$(pwd)/bin"

env go111module=off gomobile bind -v -trimpath -target android -o ./build/ipfs.aar -classpath go.ipfs;

