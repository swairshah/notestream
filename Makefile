.PHONY: build run clean

build:
	go build -o notestream ./cmd/main.go

run: build
	./notestream

clean:
	rm -f notestream
	rm -rf output
	mkdir output

test:
	go test ./...

deps:
	go mod tidy
	go mod verify
