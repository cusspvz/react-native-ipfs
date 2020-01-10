#!/bin/bash


go get -d github.com/cusspvz/rn-ipfs/go-ipfs
go install -i github.com/cusspvz/rn-ipfs/go-ipfs
go mod download

gomobile bind -v -trimpath -target android -o ./ipfs.aar -classpath io.ipfs.rn github.com/cusspvz/rn-ipfs/go-ipfs;

