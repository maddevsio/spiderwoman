FROM golang:1.8
WORKDIR /go/src/github.com/maddevsio/spiderwoman
COPY . .
COPY config.yaml config.yaml
CMD make run-forever
