#!/bin/bash

export PATH="$PATH:$GOPATH/bin"

go get -u github.com/jteeuwen/go-bindata/...

cd ./front
./build.sh

cd ..
./build.sh
