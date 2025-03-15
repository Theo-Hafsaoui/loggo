build:
	go build -o loggo main.go

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

test:
	go test -v ./...
