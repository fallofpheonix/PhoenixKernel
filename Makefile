.PHONY: all build test clean bpf

all: build test

build:
	go build ./...

test:
	go test -v -race ./...

bpf:
	go generate ./...

clean:
	go clean
	rm -f probes/*.o
