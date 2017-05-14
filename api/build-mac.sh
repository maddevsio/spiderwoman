#!/bin/bash
rm api
docker run --rm -v "$HOME"/src/go:/go -w /go/src/github.com/maddevsio/spiderwoman/api --name go golang:1.8 make build