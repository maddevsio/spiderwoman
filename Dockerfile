FROM golang:1.8
WORKDIR /go/src/github.com/maddevsio/spiderwoman
COPY . .
ENV DB-PATH root:root@tcp(spiderdb:3306)/spiderwoman?multiStatements=true
COPY config.production.yaml config.yaml
CMD make run-forever
