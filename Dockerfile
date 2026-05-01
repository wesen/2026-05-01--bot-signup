# syntax=docker/dockerfile:1.7

FROM node:22-bookworm-slim AS web
WORKDIR /src/ui

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN corepack enable \
  && corepack prepare pnpm@10.15.1 --activate \
  && pnpm install --frozen-lockfile

COPY ui/ ./
RUN pnpm run build

FROM golang:1.25-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web /src/ui/dist ./internal/web/embed/public
RUN CGO_ENABLED=0 go build -tags embed -trimpath -ldflags="-s -w" -o /out/bot-signup ./cmd/bot-signup

FROM debian:bookworm-slim AS runtime

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* \
  && useradd --system --uid 10001 --home-dir /app --shell /usr/sbin/nologin appuser

WORKDIR /app

COPY --from=build /out/bot-signup /usr/local/bin/bot-signup

RUN mkdir -p /data \
  && chown -R appuser:appuser /app /data \
  && chmod +x /usr/local/bin/bot-signup

USER appuser

EXPOSE 8080
VOLUME ["/data"]

ENTRYPOINT ["bot-signup"]
CMD ["serve"]
