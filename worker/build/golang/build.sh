#!/bin/bash

REPO=$1

export GOPATH=$HOME
go get -v $REPO

cd "$GOPATH/src/$REPO"
go get -v -t ./...
go test -v ./...
