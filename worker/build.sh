#!/bin/bash

REPO="github.com/dpolansky/go-poet"

export GOPATH=$HOME
go get -v $REPO

cd "$GOPATH/src/$REPO"
go get -v -t ./...
go test -v ./...
