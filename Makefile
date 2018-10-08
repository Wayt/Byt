# Byt makefile

all: generate run

test:
	go test -v ./...

build:
	go build

generate:
	go get -u github.com/jteeuwen/go-bindata/...
	go generate -v ./...

run:
	go run main.go http.go static.go

.PHONY: all test build generate run
