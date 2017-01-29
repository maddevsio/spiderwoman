test:
	go test -v ./...

run:
	go run main.go forever

runlog:
	go run main.go forever > log 2>&1

