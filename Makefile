.PHONY: dev-backend dev-frontend storybook storybook-build frontend-check test build-web build docker-build docker-smoke

IMAGE_REPOSITORY ?= ghcr.io/wesen/2026-05-01--bot-signup
IMAGE_TAG ?= local
IMAGE ?= $(IMAGE_REPOSITORY):$(IMAGE_TAG)

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

build-web:
	go run ./cmd/build-web

build: frontend-check storybook-build build-web
	go build -tags embed -o bin/bot-signup ./cmd/bot-signup

docker-build:
	docker build -t $(IMAGE) .

docker-smoke: docker-build
	docker run --rm $(IMAGE) --help
	docker run --rm $(IMAGE) serve --help
