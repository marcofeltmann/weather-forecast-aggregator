all: build qa

build: 
	go build -o sampleservice ./main.go

qa: lint test

lint:
	golangci-lint run .

test:
	go test ./...
