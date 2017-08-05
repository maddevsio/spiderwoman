#!/bin/bash
cd ../
rm crawler
docker run --rm -v "$HOME"/src/go:/go -w /go/src/github.com/maddevsio/spiderwoman --name go golang:1.8 make build
