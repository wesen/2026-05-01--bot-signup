---
Title: Investigation Diary
Ticket: bot-signup
Status: active
Topics:
    - go
    - react
    - sqlite
    - discord
    - signup
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
Summary: "Chronological diary of the bot-signup platform design and analysis work."
LastUpdated: 2026-05-01
WhatFor: "Track what was done, what was decided, and what to do next."
WhenToUse: "Read this before resuming work on the bot-signup ticket."
---

# Investigation Diary

## 2026-05-01 — Initial design and analysis

### What was done

1. Created ticket `bot-signup` with docmgr
2. Cloned and read the go-go-golems/discord-bot repo (README.md, tutorial, example bots)
3. Read the go-web-frontend-embed skill for the Go+React+go:embed pattern
4. Wrote a comprehensive 16-section design document covering:
   - Executive summary, problem statement, system overview
   - Architecture (client-server with go:embed)
   - Database schema (users + bot_credentials tables in SQLite)
   - Full API reference (auth, profile, admin, health endpoints)
   - Frontend page wireframes for all 10+ pages
   - Authentication with JWT + bcrypt
   - Admin backend design (approval workflow, credential management)
   - Tutorial content from discord-bot
   - Complete project structure with every file
   - 11-phase implementation plan (Day 1-7)
   - Pseudocode for key flows (server startup, route registration, signup handler, RTK Query setup, Storybook stories)
   - Testing strategy (Go table-driven tests, Storybook interaction tests, manual E2E)
   - Risks, alternatives, open questions
   - References

### What was decided

- **RTK Query** for API state management (instead of manual fetch) — gives caching, loading/error states, auto-generated hooks
- **Storybook** with stories created alongside every component — serves as visual testing and living documentation
- **Tailwind** (not themable CSS) for styling
- **SQLite** with WAL mode for the database
- **JWT** with 24-hour expiry for auth
- **Manual Discord ID entry** (not OAuth) for V1 simplicity

### What was tricky

- Balancing the level of detail for an intern audience — needed to explain concepts like SQLite, JWT, and go:embed before using them
- Structuring the RTK Query API slice to cover all endpoints while keeping the example readable

### Next steps

1. Upload the design doc to reMarkable
2. Begin Phase 1 (project scaffolding) when ready to implement

## 2026-05-01 — Phase 1 implementation

### What was done

1. Committed the initial docmgr design ticket as a clean baseline (`3ac2707`).
2. Initialized the Go module as `github.com/go-go-golems/bot-signup`.
3. Added Cobra-based CLI entrypoint at `cmd/bot-signup/main.go` with a `serve` command.
4. Added the initial server package with `GET /api/health`.
5. Added `Makefile` targets for `dev-backend`, `test`, and `build`.
6. Added `.gitignore` and kept `data/.gitkeep` while ignoring runtime database files.

### Commands run

```bash
go mod init github.com/go-go-golems/bot-signup
go get github.com/spf13/cobra@latest
gofmt -w cmd/bot-signup/main.go internal/server/*.go
go mod tidy
go test ./...
```

### What worked

- `go test ./...` passes for the initial scaffold.
- The health route returns a JSON payload from a Go 1.22-style `http.ServeMux` route.

### What was tricky

- Nothing significant in Phase 1; this was a small scaffold.

### Next steps

1. Commit Phase 1 scaffold.
2. Start Phase 2 database layer.

## 2026-05-01 — Phase 2 implementation

### What was done

1. Added `modernc.org/sqlite` as the pure-Go SQLite driver.
2. Added `internal/database` with:
   - `database.go` for opening SQLite, configuring pragmas, and running embedded migrations
   - `models.go` for `User` and `BotCredentials`
   - `users.go` for user CRUD and list-by-status
   - `credentials.go` for Discord bot credential CRUD
   - `migrations/001_initial.sql` for `users`, `bot_credentials`, and `schema_migrations`
3. Wired server startup to open the SQLite database and run migrations before serving HTTP.
4. Added tests for migrations, user CRUD, uniqueness constraints, and credential CRUD/cascade delete.

### Commands run

```bash
go get modernc.org/sqlite@latest
gofmt -w internal/database/*.go cmd/bot-signup/main.go internal/server/*.go
go mod tidy
go test ./...
```

### What worked

- `go test ./...` passes.
- Migrations are embedded and recorded in `schema_migrations`.
- SQLite foreign keys are enabled, so deleting a user cascades to credentials.

### What was tricky

- The migration runner creates `schema_migrations` before applying embedded migrations so it can track the first migration cleanly.
- Runtime database files are ignored by `.gitignore`; only `data/.gitkeep` is committed.

### Next steps

1. Commit Phase 2 database layer.
2. Start Phase 3 authentication (bcrypt, JWT, signup/login handlers).

## 2026-05-01 — Phase 3 implementation

### What was done

1. Added `internal/auth` with bcrypt password hashing and JWT generation/parsing.
2. Added auth middleware that validates `Authorization: Bearer <token>` and injects user ID / role into request context.
3. Added `AdminOnly` middleware for later admin endpoints.
4. Added auth routes:
   - `POST /api/auth/signup`
   - `POST /api/auth/login`
   - `POST /api/auth/logout`
   - `GET /api/auth/me`
5. Added validation for Discord ID, email, display name, and password length.
6. Added server tests covering signup → login → authenticated `/me`, and signup validation failures.

### Commands run

```bash
go get golang.org/x/crypto/bcrypt github.com/golang-jwt/jwt/v5@latest
gofmt -w internal/auth/*.go internal/server/*.go cmd/bot-signup/main.go
go mod tidy
go test ./...
```

### What worked

- `go test ./...` passes.
- Signup creates a waiting user and returns a JWT.
- Login verifies the bcrypt password and returns a JWT.
- `/api/auth/me` resolves the current user via the signed token.

### What was tricky

- Kept auth context values typed through unexported context keys to avoid collisions.
- The CLI currently defaults to `dev-insecure-change-me` for local JWT signing; production should set `JWT_SECRET`.

### Next steps

1. Commit Phase 3 auth implementation.
2. Start Phase 4 profile and admin handlers.

## 2026-05-01 — Phase 4 implementation

### What was done

1. Added profile routes:
   - `GET /api/profile`
   - `PUT /api/profile`
   - `PUT /api/profile/password`
2. Added public stats route `GET /api/stats`.
3. Added admin routes:
   - `GET /api/admin/waitlist`
   - `GET /api/admin/users`
   - `POST /api/admin/users/{id}/approve`
   - `POST /api/admin/users/{id}/reject`
   - `POST /api/admin/users/{id}/suspend`
   - `PUT /api/admin/users/{id}/credentials`
   - `DELETE /api/admin/users/{id}`
4. Added database helpers for profile updates, password updates, role updates, all-user listing, stats, and transactional approval.
5. Added tests for profile updates/password changes, admin approval, and non-admin rejection.

### Commands run

```bash
gofmt -w internal/database/*.go internal/server/*.go
go test ./...
```

### What worked

- `go test ./...` passes.
- Admin approval updates the user status to `approved` and inserts credentials in a transaction.
- Non-admin users receive `403 Forbidden` on admin routes.

### What was tricky

- Approval needs to be atomic because it touches both `users` and `bot_credentials`.
- Tests generate JWTs with the role stored in the DB; admin role changes must happen before token generation.

### Next steps

1. Commit Phase 4 profile/admin implementation.
2. Start Phase 5 frontend scaffolding with Vite, Tailwind, Storybook, Redux Toolkit, and RTK Query.

## 2026-05-01 — Requirement pivot: Discord OAuth only and VibeBot Sessions visual reference

### What changed

The product direction changed before starting the frontend: the app should not have passwords at all. Authentication should be Discord OAuth only, using HTTP-only same-site session cookies. The earlier bcrypt/JWT/password signup implementation is now explicitly superseded and must be replaced rather than preserved for backwards compatibility.

The user also supplied a landing-page reference image at `/tmp/pi-clipboard-92d825d5-a5a0-4f6c-be68-3edd25c51e5c.png`. I copied it into the ticket as:

```text
ttmp/2026/05/01/bot-signup--discord-bot-vibe-coding-signup-platform/sources/01-vibebot-sessions-ui-reference.png
```

### UI reference summary

The target UI is a clean SaaS-style "VibeBot Sessions" landing page:

- white/off-white background with subtle purple/blue gradient depth;
- top nav with robot logo, "VibeBot Sessions", About/FAQ links, purple Sign Up button;
- two-column hero: left value prop, right white signup/reservation card;
- purple Discord-blurple accent color;
- badge text `VIBE + CODE + DISCORD`;
- headline: "Build a Discord Bot. Vibe. Code. Deploy.";
- signup card title: "Sign Up for a Session";
- primary CTA should become "Continue with Discord" or equivalent Discord OAuth CTA;
- three feature cards under "What you get".

### Documentation updates

Updated the implementation guide to:

1. Link the image in frontmatter `RelatedFiles`.
2. Make Discord OAuth the only auth path.
3. Remove password signup/login/change-password from the architecture and page plan.
4. Replace JWT/localStorage guidance with HTTP-only session cookie guidance.
5. Add Pyxis-derived Discord OAuth operational notes: exact redirect URL matching, bot guild install if role/member lookup is needed, and Server Members Intent caveat.
6. Update tasks to introduce Phase 3R and Phase 4R because the already committed password/JWT backend must be refactored.

### Commands run

```bash
cp /tmp/pi-clipboard-92d825d5-a5a0-4f6c-be68-3edd25c51e5c.png \
  ttmp/2026/05/01/bot-signup--discord-bot-vibe-coding-signup-platform/sources/01-vibebot-sessions-ui-reference.png
rg -n "password|Password|POST /api/auth/signup|POST /api/auth/login|JWT|bcrypt|ChangePassword|LoginPage|SignupPage|useLoginMutation|useSignupMutation|localStorage|password_hash|auth_provider|profile/password" \
  ttmp/2026/05/01/bot-signup--discord-bot-vibe-coding-signup-platform/design-doc/01-full-system-design-and-implementation-guide.md
```

### What was tricky

The backend already has a committed bcrypt/JWT implementation. Since there is no backwards-compatibility requirement, the clean path is not to layer OAuth on top of password auth, but to replace the auth package, schema, routes, and frontend plan with Discord OAuth/session semantics.

### Next steps

1. Commit the requirement-pivot documentation and stored image.
2. Implement Phase 3R: replace password/JWT auth with Discord OAuth and HTTP-only sessions.

## 2026-05-01 — Phase 3R/4R implementation

### What was done

1. Removed password/JWT auth files and tests.
2. Added `golang.org/x/oauth2` and implemented `internal/auth/discord_oauth.go`:
   - Discord authorize URL generation,
   - authorization-code exchange,
   - `/users/@me` fetch,
   - Discord avatar URL helper.
3. Added `internal/auth/sessions.go`:
   - signed HTTP-only session cookie,
   - signed short-lived OAuth state cookie,
   - return-to sanitization,
   - logout cookie clearing.
4. Changed `users` schema/model to Discord OAuth fields:
   - `discord_id`, `email`, `display_name`, `avatar_url`, `last_login_at`, `status`, `role`, timestamps.
   - Removed `password_hash` and auth-provider/password concepts.
5. Replaced `CreateUser` with `UpsertDiscordUser`.
6. Replaced auth routes with:
   - `GET /auth/discord/login`,
   - `GET /auth/discord/callback`,
   - `POST /api/auth/logout`,
   - `GET /api/auth/me`.
7. Replaced JWT middleware with `SessionMiddleware`.
8. Updated profile/admin handlers and tests to use session cookies and OAuth-created users.
9. Updated implementation guide snippets that still referenced the old password/JWT flow.

### Commands run

```bash
go get golang.org/x/oauth2@latest
rm internal/auth/jwt.go internal/auth/password.go internal/auth/auth_test.go internal/server/auth_handlers_test.go
gofmt -w cmd/bot-signup/main.go internal/auth/*.go internal/database/*.go internal/server/*.go
go mod tidy
go test ./...
rg -n "CreateUser|HashPassword|CheckPassword|jwtSecret|AuthMiddleware|bcrypt|password_hash|useLoginMutation|localStorage|LoginPage|SignupPage" \
  ttmp/2026/05/01/bot-signup--discord-bot-vibe-coding-signup-platform/design-doc/01-full-system-design-and-implementation-guide.md
```

### What worked

- `go test ./...` passes.
- OAuth login test verifies redirect + HTTP-only OAuth state cookie.
- OAuth callback test verifies user creation, HTTP-only session cookie, and `/api/auth/me`.
- Profile/admin tests now authenticate with session cookies.

### What was tricky

- Because no backwards compatibility is required, the schema migration was rewritten instead of adding incremental compatibility columns. This is clean for a new app, but any local dev DB from the old auth implementation should be deleted.
- OAuth state uses a signed short-lived cookie. That avoids adding an `oauth_states` table while still preventing CSRF and open redirects.

### Next steps

1. Commit Phase 3R/4R.
2. Start Phase 5 frontend scaffolding and build the VibeBot Sessions landing UI to match the stored reference image.

## 2026-05-01 — Phase 5/6 frontend scaffold and VibeBot landing page

### What was done

1. Created `ui/` with Vite + React + TypeScript.
2. Added Tailwind v4 through `@tailwindcss/vite`.
3. Added React Router, Redux Toolkit, React Redux, RTK Query, react-markdown, remark-gfm, and lucide-react.
4. Added Storybook 10 with the Vite framework and docs/a11y addons.
5. Added RTK Query `apiSlice` with same-origin cookie credentials and initial endpoints:
   - `getMe`,
   - `logout`,
   - `getStats`.
6. Added Redux store and AuthProvider for Discord OAuth navigation + session-cookie auth.
7. Built the VibeBot Sessions landing page to match the reference image:
   - off-white/purple SaaS background,
   - logo/nav with About/FAQ/Sign Up,
   - two-column hero,
   - right-side signup card,
   - "Continue with Discord" CTA,
   - three "What you get" feature cards.
8. Added Storybook stories for:
   - `DiscordOAuthButton`,
   - `SessionSignupCard`,
   - `FeatureCard`.
9. Added Makefile targets for frontend dev, Storybook, Storybook build, and frontend checks.

### Commands run

```bash
pnpm create vite ui --template react-ts
pnpm --dir ui install
pnpm --dir ui add react-router-dom react-markdown remark-gfm @reduxjs/toolkit react-redux lucide-react
pnpm --dir ui add -D tailwindcss @tailwindcss/vite storybook @storybook/react-vite @storybook/addon-docs @storybook/addon-a11y
pnpm --dir ui lint
pnpm --dir ui build
pnpm --dir ui build-storybook
rm -rf ui/dist ui/storybook-static
go test ./...
```

### What worked

- `pnpm --dir ui lint` passes.
- `pnpm --dir ui build` passes.
- `pnpm --dir ui build-storybook` passes.
- `go test ./...` still passes.

### What was tricky

- Storybook 10 does not have a matching `@storybook/test` version; I removed `@storybook/test` and kept static Storybook stories for now.
- Storybook build output and Vite dist output must be deleted before commit and are already ignored by `.gitignore`.
- React Fast Refresh required moving `AuthContext` exports into a separate `context.ts` file.

### Next steps

1. Commit Phase 5/6 frontend scaffold.
2. Start Phase 7 user pages: waiting list, profile, credential cards, route protection.

## 2026-05-01 — Phase 7 user pages

### What was done

1. Started Storybook in tmux session `bot-signup-storybook` so progress is visible at `http://localhost:6006/`.
2. Added shared UI components:
   - `StatusBadge`,
   - `CredentialCard`,
   - `ProtectedRoute`.
3. Added Storybook stories for:
   - `StatusBadge`,
   - `CredentialCard`,
   - `WaitingListPage`,
   - `ProfilePage`.
4. Extended RTK Query with `getProfile` and profile/credential response types.
5. Added user-facing pages:
   - `/waiting-list`,
   - `/profile`,
   - `/auth/callback`,
   - placeholder `/tutorial`.
6. Wired React Router routes through `ProtectedRoute`.
7. Kept output artifacts (`ui/dist`, `ui/storybook-static`) out of git.

### Commands run

```bash
tmux new-session -d -s bot-signup-storybook 'cd /home/manuel/code/wesen/2026-05-01--bot-signup/ui && pnpm exec storybook dev -p 6006 --host 0.0.0.0'
pnpm --dir ui lint
pnpm --dir ui build
pnpm --dir ui build-storybook
rm -rf ui/dist ui/storybook-static
go test ./...
```

### What worked

- Storybook is running in tmux and rebuilds as files change.
- `pnpm --dir ui lint` passes.
- `pnpm --dir ui build` passes.
- `pnpm --dir ui build-storybook` passes.
- `go test ./...` passes.

### What was tricky

- `TutorialPage` originally used Tailwind typography `prose` classes, but the phase only needs a placeholder; I removed the dependency on typography-specific styling for now.
- `ProfilePage` story uses RTK Query cache seeding via `apiSlice.util.upsertQueryData` so it can render without a backend.

### Next steps

1. Commit Phase 7.
2. Start Phase 8 admin pages.

## 2026-05-01 — Phase 8 admin pages

### What was done

1. Extended RTK Query with admin endpoints:
   - `getWaitlist`,
   - `getAdminUsers`,
   - `approveUser`,
   - `rejectUser`,
   - `suspendUser`,
   - `deleteUser`.
2. Added admin guard component `AdminRoute`.
3. Added admin UI components:
   - `AdminStats`,
   - `AdminUserTable`,
   - `ApprovalForm`.
4. Added admin pages:
   - `/admin` dashboard with stats + waitlist,
   - `/admin/users/:id` approval detail form.
5. Added stories for admin components/pages:
   - `AdminUserTable.stories.tsx`,
   - `ApprovalForm.stories.tsx`,
   - `AdminDashboard.stories.tsx`.
6. Wired admin routes into `App.tsx`.

### Commands run

```bash
pnpm --dir ui lint
pnpm --dir ui build
pnpm --dir ui build-storybook
rm -rf ui/dist ui/storybook-static
go test ./...
```

### What worked

- UI lint, build, and Storybook build all pass.
- Go tests continue to pass.
- Admin dashboard can display seeded Storybook waitlist data and stats.

### What was tricky

- The backend currently has no `GET /api/admin/users/{id}` endpoint, so `AdminUserDetail` loads the admin user list and finds the target by ID. This works for the current small-scale admin UI but should become a dedicated endpoint if the list grows.

### Next steps

1. Commit Phase 8.
2. Start Phase 9 tutorial markdown content.

## 2026-05-01 — Phase 9 tutorial page

### What was done

1. Copied the discord-bot tutorial from `/tmp/discord-bot/pkg/doc/tutorials/building-and-running-discord-js-bots.md` to `ui/src/content/tutorial.md`.
2. Updated `TutorialPage` to render the markdown with `react-markdown` and `remark-gfm`.
3. Added local markdown styling in `ui/src/index.css` for headings, code blocks, inline code, tables, lists, links, and blockquotes.
4. Added `TutorialPage.stories.tsx`.

### Commands run

```bash
cp /tmp/discord-bot/pkg/doc/tutorials/building-and-running-discord-js-bots.md ui/src/content/tutorial.md
pnpm --dir ui lint
pnpm --dir ui build
pnpm --dir ui build-storybook
rm -rf ui/dist ui/storybook-static
go test ./...
```

### What worked

- The full tutorial markdown imports via Vite `?raw` and renders correctly.
- Frontend lint/build and Storybook build pass.
- Go tests continue to pass.

### What was tricky

- Importing the full tutorial increases the app bundle size above the default Vite warning threshold. This is acceptable for now, but if the tutorial grows much more we should lazy-load `/tutorial` or fetch the markdown separately.

### Next steps

1. Commit Phase 9.
2. Implement Phase 10/11 with Dagger build-web, go:embed SPA serving, and CI/polish targets.
