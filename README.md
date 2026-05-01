# Bot Signup

A Go + SQLite + React/Vite application for signing people up to vibe-code Discord bots.

## Current stack

- Go backend using `net/http` / Go 1.22+ `ServeMux` patterns
- SQLite via `modernc.org/sqlite`
- Discord OAuth for signup/login
- Signed HTTP-only session cookies
- React + Vite + TypeScript frontend
- Tailwind CSS v4
- Redux Toolkit / RTK Query
- Storybook
- Dagger-backed frontend build pipeline for embedded assets

## Development

Start the backend:

```bash
make dev-backend
```

Start the frontend:

```bash
make dev-frontend
```

Start Storybook:

```bash
make storybook
```

## Discord OAuth configuration

For real OAuth, configure:

```bash
export DISCORD_CLIENT_ID=...
export DISCORD_CLIENT_SECRET=...
export DISCORD_REDIRECT_URL=http://localhost:8080/auth/discord/callback
export SESSION_SECRET=replace-me-with-a-long-random-secret
```

The Discord Developer Portal redirect URL must exactly match `DISCORD_REDIRECT_URL`.

## Build

The reproducible frontend build pipeline is:

```bash
go run ./cmd/build-web
```

This uses Dagger by default with a cached pnpm store and copies `ui/dist/` to `internal/web/embed/public/` for Go embedding.

If Docker/Dagger is unavailable, force the local pnpm fallback:

```bash
BUILD_WEB_LOCAL=1 go run ./cmd/build-web
```

Build the single binary with embedded UI:

```bash
make build
```

## Validation

```bash
go test ./...
pnpm --dir ui lint
pnpm --dir ui build
pnpm --dir ui build-storybook
```
