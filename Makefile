test:
	go test -v ./...

run:
	go run main.go callbacks.go

run-forever:
	go run main.go callbacks.go forever

run-once:
	go run main.go callbacks.go forever

run-forever-log:
	go run main.go callbacks.go forever > log 2>&1

build:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o crawler

