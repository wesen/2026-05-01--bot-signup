.PHONY: dev-backend dev-frontend storybook storybook-build frontend-check test build

dev-backend:
	go run ./cmd/bot-signup serve --addr :8080

dev-frontend:
	pnpm --dir ui dev

storybook:
	pnpm --dir ui storybook

storybook-build:
	pnpm --dir ui build-storybook

frontend-check:
	pnpm --dir ui lint
	pnpm --dir ui build

test:
	go test ./...

build: frontend-check storybook-build
	go build -o bin/bot-signup ./cmd/bot-signup
