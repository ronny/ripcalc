all: checks ripcalc

BIN_DIR ?= ./bin/
ripcalc:
	go build -o ${BIN_DIR} ./...

test:
	env go test -shuffle=on -race -covermode=atomic -coverprofile=coverage.txt -count=1 -timeout=5m  ./...

lint:
	golangci-lint run --fix

vet:
	go vet ./...

format:
	golangci-lint fmt

tidy:
	go mod tidy

generate:
	go generate -v ./...

checks: tidy format generate lint vet

ci-checks: checks
	git diff --exit-code

.PHONY: all ripcalc test lint vet format tidy generate checks ci-checks
