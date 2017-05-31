#!/bin/bash
# Script that builds necessary docker images

# build docker images
IMAGE_DIR_PATH=$GOPATH/src/github.com/dpolansky/grader-ci/pkg/worker/build

# go
docker build -t build-golang $IMAGE_DIR_PATH/golang

# java8-maven
docker pull maven
docker build -t build-java8-maven $IMAGE_DIR_PATH/java8-maven

# python3
docker build -t build-python3 $IMAGE_DIR_PATH/python3