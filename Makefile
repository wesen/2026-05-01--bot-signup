.PHONY: dev-backend test build

dev-backend:
	go run ./cmd/bot-signup serve --addr :8080

test:
	go test ./...

build:
	go build -o bin/bot-signup ./cmd/bot-signup
