objects = main.go callbacks.go actions.go crawl.go data.go init.go grab.go

test:
	go test github.com/maddevsio/spiderwoman
	go test github.com/maddevsio/spiderwoman/api
	go test github.com/maddevsio/spiderwoman/grabber
	go test github.com/maddevsio/spiderwoman/lib

run:
	go run $(objects)

run-forever:
	go run $(objects) forever

run-once:
	go run $(objects) once

run-excel:
	go run $(objects) excel

run-grab:
	go run $(objects) grab

run-forever-log:
	go run $(objects) forever > log 2>&1

build:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o crawler

