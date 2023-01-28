.PHONY: build test clean

build:
	go build -v ./...

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	go clean
	rm -f coverage.out
	rm -rf ./code_regions

lint:
	golangci-lint run

.DEFAULT_GOAL := build

generate:
	./aquila generate

update:
	./aquila update-md