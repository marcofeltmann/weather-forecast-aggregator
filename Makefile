all: build qa

build: 
	go build -o sampleservice ./main.go

qa: lint test

lint:
	golangci-lint run .

test:
	go test --count=1 ./...
	# blackbox test with k6 is missing but not really needed
