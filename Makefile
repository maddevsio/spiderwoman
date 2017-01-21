test:
	go test -v ./...

run:
	go run main.go > log 2>&1
