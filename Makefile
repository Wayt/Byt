# Byt makefile

all: generate run

test:
	go test -v ./...

dockertest:
	docker run --rm -v $(shell pwd):/go/src/app -w /go/src/app golang:latest make test

build:
	go build

generate:
	go get -u github.com/jteeuwen/go-bindata/...
	go generate -v ./...

run:
	go run main.go http.go static.go

.PHONY: all test dockertest build generate run
