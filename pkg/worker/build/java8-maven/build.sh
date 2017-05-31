#!/bin/bash

REPO=$1
cd "$HOME/$REPO"
mvn clean install
