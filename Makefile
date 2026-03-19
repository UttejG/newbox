.PHONY: build run test test-race test-integration test-update-golden coverage lint clean fmt

build:
	go build -o bin/newbox ./cmd/newbox

run:
	go run ./cmd/newbox

test:
	go test ./...

test-race:
	go test -race ./...

test-integration:
	go test -tags=integration ./...

test-update-golden:
	go test ./... -update

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	go vet ./...

clean:
	rm -rf bin/

fmt:
	gofmt -w .
