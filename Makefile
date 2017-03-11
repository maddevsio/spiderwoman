test:
	go test -v ./...

run:
	go run main.go forever

runlog:
	go run main.go forever > log 2>&1

build:
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o crawler

