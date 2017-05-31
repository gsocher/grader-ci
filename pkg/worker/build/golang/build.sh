#!/bin/bash

REPO=$1
export GOPATH=$HOME

mkdir -p $GOPATH/src
mv $HOME/$REPO $GOPATH/src/$REPO
cd "$GOPATH/src/$REPO"
go get -v -t ./...
go test -v ./...
