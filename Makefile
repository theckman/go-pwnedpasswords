.DEFAULT_GOAL := test

test:
	go test -v -covermode=count -coverprofile=profile.out .

build:
	go build -o pp ./cmd/pp/

.PHONY: test build
