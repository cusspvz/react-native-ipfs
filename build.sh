#!/bin/bash

go get golang.org/x/mobile/bind

go get -d .
go install -i .
go mod download
go mod vendor

gomobile bind -v -trimpath -target android -o ./android/libs/ipfs.aar -classpath io.ipfs github.com/cusspvz/rn-ipfs;

