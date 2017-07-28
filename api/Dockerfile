FROM golang:1.8
WORKDIR /go/src/github.com/maddevsio/spiderwoman/
COPY . .
ENV DB-PATH spiderwoman
COPY config.production.yaml config.yaml
CMD go run main.go
