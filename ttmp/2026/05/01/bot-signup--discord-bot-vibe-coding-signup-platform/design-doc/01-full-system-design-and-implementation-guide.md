---
Title: Full System Design and Implementation Guide
Ticket: bot-signup
Status: active
Topics:
    - go
    - react
    - sqlite
    - discord
    - signup
    - authentication
    - admin
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources:
    - https://github.com/go-go-golems/discord-bot
Summary: "Complete design and intern-ready implementation guide for a Discord bot vibe-coding signup platform: Go + SQLite backend, React + Vite + Tailwind frontend, waiting-list workflow, admin approval with Discord bot credentials, embedded tutorial from go-go-golems/discord-bot."
LastUpdated: 2026-05-01T09:53:40.866194437-04:00
WhatFor: "Guide a new intern through every part of building, understanding, and deploying the bot-signup platform"
WhenToUse: "Read this before writing any code. Reference during implementation of each phase."
---

# Full System Design and Implementation Guide

## 1. Executive Summary

This document is the complete blueprint for the **Bot Signup Platform** — a web application where people sign up to "vibe code" their own Discord bot. The system collects their Discord identity, places them on a waiting list, lets them set a password for their account, and then gives an administrator the power to approve their request by filling in the four Discord credentials needed to run a real bot (application ID, bot token, guild ID, public key). Once approved, the user gets a profile page showing their bot status and a rich tutorial drawn from the [go-go-golems/discord-bot](https://github.com/go-go-golems/discord-bot) project.

The platform is a **single binary** in production: a Go web server that embeds a React + Vite + Tailwind frontend using `go:embed`. During development, you run the Go API on `:8080` and the Vite dev server on `:5173` with hot module replacement. The database is SQLite — a single file, no external database server needed.

This guide is written for an **intern joining the project cold**. Every concept is explained before it is used. Every file is named before it is referenced. Every decision has a reason.

---

## 2. Problem Statement and Scope

### The problem

The go-go-golems team runs a Discord bot runtime ([github.com/go-go-golems/discord-bot](https://github.com/go-go-golems/discord-bot)). They want to let people sign up to create and run their own bots — but not everyone should get access immediately. There needs to be a controlled onboarding flow:

1. A visitor arrives at the site and learns what the platform offers.
2. The visitor signs up with their Discord ID and email.
3. The visitor is placed on a **waiting list** and prompted to create a password.
4. An **admin** reviews the waiting list and approves users one at a time.
5. On approval, the admin fills in the four Discord credentials that the user's bot will need.
6. The approved user returns, logs in, sees their profile with bot status, and reads the tutorial to start coding their bot.

### Scope

- **In scope**: signup, waiting list, password setup, admin approval with Discord credentials, user profile, tutorial content, API, database, frontend.
- **Out of scope** (for now): actual bot process management, real-time bot logs, payment/billing, multi-tenant isolation of running bots, OAuth Discord login (we use manual Discord ID entry for V1).

---

## 3. System Overview — What Is This Thing?

If you have never built a web application before, here is the mental model:

```
┌─────────────────────────────────────────────────────────────────────┐
│                        YOUR LAPTOP / SERVER                         │
│                                                                     │
│  ┌──────────────────────────┐    ┌───────────────────────────────┐  │
│  │   React Frontend (SPA)   │    │      Go Backend (API)         │  │
│  │                          │    │                               │  │
│  │  - Signup page           │◄──►│  - POST /api/signup           │  │
│  │  - Login page            │    │  - POST /api/login            │  │
│  │  - Waiting list page     │    │  - GET  /api/profile          │  │
│  │  - Profile page          │    │  - GET  /api/admin/waitlist   │  │
│  │  - Tutorial page         │    │  - POST /api/admin/approve    │  │
│  │  - Admin dashboard       │    │  - ... more routes            │  │
│  │                          │    │                               │  │
│  │  Built with:             │    │          ┌─────────────┐      │  │
│  │   React + Vite           │    │          │  SQLite DB  │      │  │
│  │   Tailwind CSS           │    │          │  (one file) │      │  │
│  │   TypeScript             │    │          └─────────────┘      │  │
│  └──────────────────────────┘    └───────────────────────────────┘  │
│                                                                     │
│  Production: Go binary embeds the React build output                │
│  (go:embed) and serves everything on port 8080.                     │
└─────────────────────────────────────────────────────────────────────┘
```

### The user journey (happy path)

```
Visitor browses site
    │
    ▼
Reads landing page + tutorial
    │
    ▼
Clicks "Sign Up"
    │
    ▼
Fills in: Discord ID, Email, Display Name
    │
    ▼
Creates a Password
    │
    ▼
Account created → status: "waiting"
    │
    ▼
Sees "You are on the waiting list" page
    │
    ▼
... time passes ...
    │
    ▼
Admin approves + fills in Discord credentials
    │
    ▼
User logs in → sees Profile with bot status: "approved"
    │
    ▼
User reads tutorial, starts coding their bot
```

### The admin journey

```
Admin logs in (special credentials)
    │
    ▼
Sees Admin Dashboard with waiting-list table
    │
    ▼
Clicks "Approve" on a user row
    │
    ▼
Form appears: enter Application ID, Bot Token, Guild ID, Public Key
    │
    ▼
Submits → user status becomes "approved"
    │
    ▼
Can also edit/revoke credentials later
```

## 4. Architecture — How the Pieces Fit Together

### 4.1 High-level architecture

The system follows a classic **client-server** pattern with a twist: in production, the client (the React SPA) is compiled into static files and embedded *inside* the Go binary using Go's `//go:embed` directive. This means deployment is a single file.

```
┌──────────────── Production Binary ────────────────┐
│                                                    │
│  Go HTTP Server (net/http, Go 1.22+ ServeMux)      │
│    │                                               │
│    ├── /api/*     → JSON API handlers              │
│    │                 (auth middleware, DB queries)   │
│    │                                               │
│    ├── /assets/*  → Static files (JS, CSS, images)  │
│    │                 served from embedded FS         │
│    │                                               │
│    └── /*         → SPA fallback (index.html)       │
│                    for client-side routing           │
│                                                    │
│  go:embed ← ui/dist/ (Vite build output)           │
│                                                    │
│  SQLite DB file: ./data/bot-signup.db               │
└────────────────────────────────────────────────────┘
```

### 4.2 Why these technologies?

| Technology | Why we chose it | What it gives us |
|---|---|---|
| **Go** | Fast compilation, single-binary output, excellent HTTP standard library | No runtime dependencies, easy deployment |
| **SQLite** | Zero-config, file-based, embedded in the process | No separate database server, easy backups (copy the file) |
| **React** | Component-based UI, huge ecosystem | Reusable UI components, declarative rendering |
| **Vite** | Fast dev server with HMR, optimized production builds | Sub-second page reloads during development |
| **Tailwind CSS** | Utility-first CSS, no separate CSS files to manage | Rapid styling without context-switching |
| **TypeScript** | Type safety catches bugs at compile time | Better IDE support, fewer runtime errors |
| **RTK Query** | Built into Redux Toolkit; auto-generates hooks for API calls | Caching, loading/error states, optimistic updates — no manual fetch boilerplate |
| **Storybook** | Isolated component development environment | Build, test, and document each UI component in isolation before wiring it into pages |
| **go:embed** | Go 1.16+ feature to embed files into the binary | Single binary deployment |

### 4.3 Development vs Production topology

**During development** you run two processes:

```
Terminal 1: Go API server
  $ go run ./cmd/bot-signup serve --dev
  → listens on :8080
  → serves /api/* routes
  → reads SQLite from ./data/bot-signup.db

Terminal 2: Vite dev server
  $ cd ui && pnpm dev
  → listens on :5173
  → proxies /api/* to :8080 (see vite.config.ts)
  → hot-reloads on file change
```

You open `http://localhost:5173` in your browser. The Vite dev server handles the frontend and forwards any `/api/*` request to the Go backend.

**In production** there is only one process:

```
$ go generate ./internal/web/     # builds the React app and copies it into the Go tree
$ go build -tags embed -o bot-signup ./cmd/bot-signup
$ ./bot-signup serve
→ listens on :8080
→ serves API + static files + SPA fallback
```

### 4.4 Request lifecycle

Here is exactly what happens when a user clicks "Sign Up":

```
Browser (React)                    Go Server                      SQLite
     │                                 │                              │
     │  POST /api/signup                │                              │
     │  { discord_id, email, name,      │                              │
     │    password }                    │                              │
     │─────────────────────────────────►│                              │
     │                                 │  1. Validate input           │
     │                                 │  2. Hash password (bcrypt)   │
     │                                 │  3. Check discord_id unique  │
     │                                 │     SELECT FROM users        │
     │                                 │─────────────────────────────►│
     │                                 │◄─────────────────────────────│
     │                                 │  4. Insert user              │
     │                                 │     status='waiting'         │
     │                                 │─────────────────────────────►│
     │                                 │◄─────────────────────────────│
     │                                 │  5. Create JWT session token │
     │                                 │                              │
     │  { token, user }                │                              │
     │◄─────────────────────────────────│                              │
     │                                 │                              │
     │  Browser stores JWT in          │                              │
     │  localStorage, redirects to     │                              │
     │  /waiting-list page             │                              │
```

## 5. Database Schema — SQLite Tables

### 5.1 Why SQLite?

SQLite is not "MySQL lite". It is a full-featured relational database engine that runs inside your process. There is no separate server to install, configure, or secure. The entire database is one file on disk. For a signup platform that will handle hundreds or low thousands of users, SQLite is more than sufficient.

Key properties:
- **ACID transactions** — your writes are safe even if the power goes out
- **WAL mode** — enables concurrent reads while writing (we enable this on startup)
- **Zero configuration** — just open the file and go

### 5.2 Entity-relationship diagram

```
┌──────────────────────────────────────┐
│              users                    │
├──────────────────────────────────────┤
│ id            INTEGER PRIMARY KEY     │
│ discord_id    TEXT    UNIQUE NOT NULL │
│ email         TEXT    UNIQUE NOT NULL │
│ display_name  TEXT    NOT NULL        │
│ password_hash TEXT    NOT NULL        │
│ status        TEXT    DEFAULT 'waiting' │
│               ('waiting','approved',  │
│                'rejected','suspended') │
│ role          TEXT    DEFAULT 'user'  │
│               ('user','admin')        │
│ created_at    DATETIME NOT NULL       │
│ updated_at    DATETIME NOT NULL       │
└──────────────┬───────────────────────┘
               │
               │ 1:1
               ▼
┌──────────────────────────────────────┐
│         bot_credentials               │
├──────────────────────────────────────┤
│ id              INTEGER PRIMARY KEY   │
│ user_id         INTEGER UNIQUE FK     │
│ application_id  TEXT    NOT NULL      │
│ bot_token       TEXT    NOT NULL      │
│ guild_id        TEXT    NOT NULL      │
│ public_key      TEXT    NOT NULL      │
│ approved_by     INTEGER FK (users)    │
│ approved_at     DATETIME              │
│ created_at      DATETIME NOT NULL     │
│ updated_at      DATETIME NOT NULL     │
└──────────────────────────────────────┘
```

### 5.3 SQL migration (the single source of truth)

The file `internal/database/migrations/001_initial.sql` contains:

```sql
-- Users table: everyone who signs up
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_id    TEXT    UNIQUE NOT NULL,
    email         TEXT    UNIQUE NOT NULL,
    display_name  TEXT    NOT NULL,
    password_hash TEXT    NOT NULL,
    status        TEXT    NOT NULL DEFAULT 'waiting'
                  CHECK(status IN ('waiting','approved','rejected','suspended')),
    role          TEXT    NOT NULL DEFAULT 'user'
                  CHECK(role IN ('user','admin')),
    created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);

-- Index for fast lookup by discord_id (used in signup/login)
CREATE INDEX IF NOT EXISTS idx_users_discord_id ON users(discord_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Bot credentials: filled in by admin on approval
CREATE TABLE IF NOT EXISTS bot_credentials (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER UNIQUE NOT NULL
                    REFERENCES users(id) ON DELETE CASCADE,
    application_id  TEXT    NOT NULL,
    bot_token       TEXT    NOT NULL,
    guild_id        TEXT    NOT NULL,
    public_key      TEXT    NOT NULL,
    approved_by     INTEGER REFERENCES users(id),
    approved_at     TEXT,
    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_bot_credentials_user_id ON bot_credentials(user_id);
```

### 5.4 Status values explained

A user's `status` field controls what they can see and do:

| Status | Meaning | What the user sees |
|---|---|---|
| `waiting` | Just signed up, waiting for admin approval | "You're on the waiting list" page. Can log in but cannot see credentials. |
| `approved` | Admin has approved them and assigned Discord credentials | Profile page with their bot credentials, tutorial, and getting-started guide. |
| `rejected` | Admin has rejected their signup | "Your application was not approved" page. |
| `suspended` | Admin has temporarily disabled their account | "Your account has been suspended" page. |

### 5.5 How we interact with the database in Go

We use the standard library `database/sql` package with the `modernc.org/sqlite` driver (a pure-Go SQLite implementation — no CGO needed, which means easy cross-compilation).

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

type DB struct {
    db *sql.DB
}

func Open(path string) (*DB, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, err
    }
    // Enable WAL mode for concurrent read/write performance
    _, _ = db.Exec("PRAGMA journal_mode=WAL")
    // Enable foreign keys
    _, _ = db.Exec("PRAGMA foreign_keys=ON")
    return &DB{db: db}, nil
}
```

Every database operation lives in a Go function on the `*DB` struct. For example:

```go
func (db *DB) CreateUser(ctx context.Context, discordID, email, displayName, passwordHash string) (*User, error) {
    // ...
}

func (db *DB) GetUserByDiscordID(ctx context.Context, discordID string) (*User, error) {
    // ...
}
```

## 6. API Reference — Every Endpoint

All API endpoints live under `/api/`. They accept and return JSON. The Go server uses Go 1.22+ `http.ServeMux` with method+path pattern matching (no third-party router).

### 6.1 Authentication endpoints

#### `POST /api/auth/signup`

Creates a new user account. The user starts with `status='waiting'`.

**Request body:**
```json
{
  "discord_id": "123456789012345678",
  "email": "user@example.com",
  "display_name": "CoolBotDev",
  "password": "my-secure-password"
}
```

**Success response (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "discord_id": "123456789012345678",
    "email": "user@example.com",
    "display_name": "CoolBotDev",
    "status": "waiting",
    "role": "user"
  }
}
```

**Validation rules:**
- `discord_id`: required, must be a numeric string (Discord snowflake ID)
- `email`: required, must be a valid email format
- `display_name`: required, 2–50 characters
- `password`: required, minimum 8 characters

**Error responses:**
- `400 Bad Request` — validation failure, with specific field errors
- `409 Conflict` — discord_id or email already registered

---

#### `POST /api/auth/login`

Authenticates an existing user and returns a JWT token.

**Request body:**
```json
{
  "email": "user@example.com",
  "password": "my-secure-password"
}
```

**Success response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { ... }
}
```

**Error responses:**
- `401 Unauthorized` — invalid email or password

---

#### `POST /api/auth/logout`

Invalidates the current session (client-side: delete the JWT from localStorage).

**Response (200):** `{ "message": "logged out" }`

---

#### `GET /api/auth/me`

Returns the currently authenticated user's profile. Requires a valid JWT in the `Authorization: Bearer <token>` header.

**Success response (200):**
```json
{
  "id": 1,
  "discord_id": "123456789012345678",
  "email": "user@example.com",
  "display_name": "CoolBotDev",
  "status": "approved",
  "role": "user",
  "created_at": "2026-05-01T10:00:00Z"
}
```

**Error responses:**
- `401 Unauthorized` — missing or invalid token

### 6.2 User profile endpoints

#### `GET /api/profile`

Returns the current user's profile **including bot credentials** (if approved). Requires auth.

**Success response (200) for approved user:**
```json
{
  "user": { ... },
  "bot_credentials": {
    "application_id": "987654321098765432",
    "bot_token": "MTIz...ODk.GHx...",
    "guild_id": "111222333444555666",
    "public_key": "abcdef123456...",
    "approved_at": "2026-05-02T14:30:00Z"
  }
}
```

**Success response (200) for waiting user:**
```json
{
  "user": { "status": "waiting", ... },
  "bot_credentials": null,
  "message": "Your account is pending approval."
}
```

---

#### `PUT /api/profile`

Updates the current user's display name or email. Requires auth.

**Request body:**
```json
{
  "display_name": "NewName",
  "email": "new@example.com"
}
```

**Success response (200):** `{ "user": { ... updated ... } }`

---

#### `PUT /api/profile/password`

Changes the user's password. Requires the current password for verification.

**Request body:**
```json
{
  "current_password": "old-password",
  "new_password": "new-secure-password"
}
```

**Success response (200):** `{ "message": "password updated" }`

### 6.3 Admin endpoints

All admin endpoints require a JWT for a user with `role='admin'`. If a non-admin calls these, the server returns `403 Forbidden`.

#### `GET /api/admin/waitlist`

Returns all users with `status='waiting'`, ordered by signup date (oldest first).

**Success response (200):**
```json
{
  "users": [
    {
      "id": 5,
      "discord_id": "123456789012345678",
      "email": "user@example.com",
      "display_name": "CoolBotDev",
      "created_at": "2026-05-01T10:00:00Z"
    },
    ...
  ]
}
```

---

#### `GET /api/admin/users`

Returns all users (for the admin user management view). Supports pagination.

**Query parameters:** `?page=1&per_page=20&status=waiting`

**Success response (200):**
```json
{
  "users": [ ... ],
  "total": 47,
  "page": 1,
  "per_page": 20
}
```

---

#### `POST /api/admin/users/{id}/approve`

Approves a waiting user and assigns their Discord bot credentials.

**Request body:**
```json
{
  "application_id": "987654321098765432",
  "bot_token": "MTIz...ODk.GHx...",
  "guild_id": "111222333444555666",
  "public_key": "abcdef123456..."
}
```

**Validation rules:**
- All four fields are required and must be non-empty strings
- `application_id`: must be numeric
- `bot_token`: must match the general Discord token format
- `guild_id`: must be numeric
- `public_key`: must be a non-empty hex/base64 string

**Success response (200):**
```json
{
  "message": "User approved successfully",
  "user": { "status": "approved", ... },
  "bot_credentials": { ... }
}
```

**Error responses:**
- `404 Not Found` — user does not exist
- `409 Conflict` — user is not in 'waiting' status
- `400 Bad Request` — validation failure

---

#### `POST /api/admin/users/{id}/reject`

Rejects a waiting user.

**Success response (200):** `{ "message": "User rejected" }`

---

#### `POST /api/admin/users/{id}/suspend`

Suspends an approved user.

**Success response (200):** `{ "message": "User suspended" }`

---

#### `PUT /api/admin/users/{id}/credentials`

Updates the bot credentials for an approved user.

**Request body:** same fields as approve.

**Success response (200):** `{ "message": "Credentials updated" }`

---

#### `DELETE /api/admin/users/{id}`

Deletes a user and their credentials (hard delete). Use with caution.

**Success response (200):** `{ "message": "User deleted" }`

### 6.4 Health and info endpoints

#### `GET /api/health`

Returns server health status. No auth required.

```json
{ "status": "ok", "version": "0.1.0" }
```

---

#### `GET /api/stats`

Returns public statistics (no auth required, optional for landing page):

```json
{
  "total_users": 47,
  "approved_users": 32,
  "waiting_users": 12,
  "bots_running": 28
}
```

## 7. Frontend Pages and Components

### 7.1 Page map

The frontend is a **Single Page Application (SPA)**. The browser loads `index.html` once, and React Router handles URL changes without full-page reloads. Here is every page the user can visit:

```
/                     → Landing page (hero + features + CTA to sign up)
/signup               → Signup form (discord_id, email, name, password)
/login                → Login form (email, password)
/waiting-list         → "You're on the waiting list" status page
/profile              → User profile + bot credentials (if approved)
/profile/password     → Change password form
/tutorial             → Discord bot tutorial (from go-go-golems/discord-bot)
/admin                → Admin dashboard (requires admin role)
/admin/waitlist       → Waiting-list management
/admin/users/{id}     → Single-user detail + approve/reject/edit credentials
*                     → 404 Not Found page
```

### 7.2 Landing page (`/`)

This is the first thing a visitor sees. It needs to accomplish three things:

1. **Explain what the platform is** — "Sign up to create your own Discord bot using JavaScript"
2. **Show credibility** — link to the go-go-golems/discord-bot repo, show stats
3. **Get them to sign up** — a clear call-to-action button

Layout (Tailwind classes as mental model):

```
┌───────────────────────────────────────────────────────────┐
│  Navbar: [Logo] [Tutorial] [Login] [Sign Up]              │
├───────────────────────────────────────────────────────────┤
│                                                           │
│  Hero Section (bg-gradient, dark theme)                   │
│  ┌──────────────────────────────────────────────────┐     │
│  │  🤖 Build Your Own Discord Bot                    │     │
│  │                                                    │     │
│  │  Sign up, get approved, and start coding your     │     │
│  │  bot in JavaScript — powered by the go-go-golems   │     │
│  │  discord-bot runtime.                              │     │
│  │                                                    │     │
│  │  [Get Started →]  [Read the Tutorial]              │     │
│  └──────────────────────────────────────────────────┘     │
│                                                           │
│  Features Grid (3 cards)                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐               │
│  │ 🎯 Easy  │  │ 📚 Learn │  │ ⚡ Fast   │               │
│  │ Signup   │  │ Tutorial │  │ Deploy   │               │
│  │ with     │  │ included │  │ instantly│               │
│  │ Discord  │  │ step by  │  │ with     │               │
│  │ ID       │  │ step     │  │ Go power │               │
│  └──────────┘  └──────────┘  └──────────┘               │
│                                                           │
│  Stats Bar                                                │
│  [47 Users]  [32 Bots Running]  [12 on Waitlist]         │
│                                                           │
│  Footer                                                   │
└───────────────────────────────────────────────────────────┘
```

**React component tree:**
- `LandingPage` (page component)
  - `Navbar`
  - `HeroSection`
  - `FeaturesGrid`
    - `FeatureCard` × 3
  - `StatsBar`
  - `Footer`

### 7.3 Signup page (`/signup`)

A clean form that collects Discord ID, email, display name, and password.

```
┌────────────────────────────────────────────┐
│  Create Your Account                       │
│                                            │
│  Discord User ID                           │
│  ┌──────────────────────────────────────┐  │
│  │ 123456789012345678                   │  │
│  └──────────────────────────────────────┘  │
│  💡 How to find your Discord ID:           │
│     Settings → Advanced → Developer Mode   │
│     Right-click your name → Copy User ID   │
│                                            │
│  Email                                     │
│  ┌──────────────────────────────────────┐  │
│  │ you@example.com                      │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  Display Name                              │
│  ┌──────────────────────────────────────┐  │
│  │ CoolBotDev                           │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  Password (min 8 chars)                    │
│  ┌──────────────────────────────────────┐  │
│  │ ••••••••                             │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  [Create Account]                          │
│                                            │
│  Already have an account? [Log in]         │
└────────────────────────────────────────────┘
```

**Key UX details:**
- The Discord ID field has a helper tooltip explaining how to find it
- Password field has a show/hide toggle
- Form validates in real-time (client-side) before submitting
- On success, the user is redirected to `/waiting-list`
- On error (e.g., duplicate Discord ID), the form shows the error inline

**React component tree:**
- `SignupPage`
  - `SignupForm`
    - `FormField` (reusable) × 4
    - `DiscordIDHelper` (tooltip/info box)

### 7.4 Login page (`/login`)

Simple email + password form.

```
┌────────────────────────────────────────────┐
│  Welcome Back                              │
│                                            │
│  Email                                     │
│  ┌──────────────────────────────────────┐  │
│  │ you@example.com                      │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  Password                                  │
│  ┌──────────────────────────────────────┐  │
│  │ ••••••••                             │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  [Log In]                                  │
│                                            │
│  Don't have an account? [Sign up]          │
└────────────────────────────────────────────┘
```

After login, the user is redirected based on their status:
- `waiting` → `/waiting-list`
- `approved` → `/profile`
- `rejected` → `/waiting-list` (shows rejection message)
- `suspended` → `/waiting-list` (shows suspension message)
- `admin` → `/admin`

### 7.5 Waiting-list page (`/waiting-list`)

This page tells the user where they stand.

```
┌────────────────────────────────────────────┐
│  Your Signup Status                        │
│                                            │
│  ┌──────────────────────────────────────┐  │
│  │  🕐 You are on the waiting list.     │  │
│  │                                       │  │
│  │  Your request is being reviewed by    │  │
│  │  our team. We'll notify you by email   │  │
│  │  when you're approved.                 │  │
│  │                                       │  │
│  │  Signed up: May 1, 2026              │  │
│  │  Position: ~#5                        │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  While you wait, check out the tutorial:   │
│  [Read the Bot Tutorial →]                 │
│                                            │
│  [Log out]                                 │
└────────────────────────────────────────────┘
```

### 7.6 Profile page (`/profile`)

This is what approved users see. It shows their account details and **all four Discord credentials** they need to configure their bot.

```
┌────────────────────────────────────────────────────────┐
│  Your Bot Dashboard                                     │
│                                                         │
│  ┌─ Account Info ────────────────────────────────────┐ │
│  │  Discord ID: 123456789012345678                   │ │
│  │  Email:      user@example.com                     │ │
│  │  Status:     ✅ Approved                          │ │
│  │  [Edit Profile]  [Change Password]                │ │
│  └──────────────────────────────────────────────────┘ │
│                                                         │
│  ┌─ Bot Credentials ────────────────────────────────┐ │
│  │  ⚠️  Keep these secret! Never share your token.  │ │
│  │                                                   │ │
│  │  Application ID:  987654321098765432              │ │
│  │  Bot Token:       MTIz...ODk  [Show] [Copy]      │ │
│  │  Guild ID:        111222333444555666              │ │
│  │  Public Key:      abcdef123456... [Copy]          │ │
│  │                                                   │ │
│  │  Approved by admin on May 2, 2026                │ │
│  └──────────────────────────────────────────────────┘ │
│                                                         │
│  ┌─ Next Steps ─────────────────────────────────────┐ │
│  │  1. Install discord-bot                           │ │
│  │     brew install go-go-golems/tap/discord-bot     │ │
│  │                                                   │ │
│  │  2. Set environment variables:                    │ │
│  │     export DISCORD_BOT_TOKEN="your-token"          │ │
│  │     export DISCORD_APPLICATION_ID="your-app-id"    │ │
│  │     export DISCORD_GUILD_ID="your-guild-id"        │ │
│  │                                                   │ │
│  │  3. Run the bot:                                  │ │
│  │     discord-bot bots ping run --sync-on-start     │ │
│  │                                                   │ │
│  │  [Read the Full Tutorial →]                       │ │
│  └──────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────┘
```

**Key UX detail:** The bot token is hidden by default (masked with `••••••••`). There is a "Show" toggle and a "Copy" button.

### 7.7 Tutorial page (`/tutorial`)

This page embeds the tutorial content from the go-go-golems/discord-bot project. It is a long-form, readable guide.

```
┌────────────────────────────────────────────┐
│  [Back to Profile]    Tutorial             │
├────────────────────────────────────────────┤
│                                            │
│  # Building Your Discord Bot               │
│                                            │
│  ## What you need                          │
│  - Your bot credentials (from profile)     │
│  - A Discord server where you are admin    │
│  - Basic JavaScript knowledge              │
│                                            │
│  ## Step 1: Install discord-bot            │
│  ```bash                                   │
│  brew install go-go-golems/tap/discord-bot │
│  ```                                       │
│                                            │
│  ## Step 2: Write your first bot           │
│  ```js                                     │
│  const { defineBot } = require("discord") │
│  // ... (full tutorial content)            │
│  ```                                       │
│                                            │
│  ... (continues with full tutorial)        │
│                                            │
│  [← Back to Profile]  [Back to Top ↑]     │
└────────────────────────────────────────────┘
```

### 7.8 Admin dashboard (`/admin`)

The admin sees an overview of all users and a waiting-list management interface.

```
┌───────────────────────────────────────────────────────────────┐
│  Admin Dashboard                                    [Log out] │
├───────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─ Stats ─────────────────────────────────────────────────┐ │
│  │  Total Users: 47  │  Waiting: 12  │  Approved: 32       │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─ Waiting List ──────────────────────────────────────────┐ │
│  │                                                         │ │
│  │  Discord ID     │ Name       │ Email       │ Joined    │ │
│  │  ─────────────── ──────────── ───────────── ─────────── │ │
│  │  12345...        │ CoolBotDev │ user@ex...  │ May 1     │ │
│  │  98765...        │ BotMaster  │ bot@ex...   │ May 1     │ │
│  │                                                         │ │
│  │  [Approve] [Reject]                                     │ │
│  │                                                         │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                               │
│  [View All Users]                                             │
└───────────────────────────────────────────────────────────────┘
```

### 7.9 Admin user detail / approval (`/admin/users/{id}`)

When the admin clicks "Approve" on a waiting user, they see this form:

```
┌───────────────────────────────────────────────────────────┐
│  Approve User: CoolBotDev                                 │
│                                                           │
│  Discord ID: 123456789012345678                           │
│  Email:      user@example.com                             │
│  Signed up:  May 1, 2026                                  │
│                                                           │
│  ┌─ Discord Bot Credentials ──────────────────────────┐  │
│  │                                                     │  │
│  │  Application ID *                                    │  │
│  │  ┌────────────────────────────────────────────┐     │  │
│  │  │                                            │     │  │
│  │  └────────────────────────────────────────────┘     │  │
│  │  From: Discord Developer Portal → My Application   │  │
│  │                                                     │  │
│  │  Bot Token *                                         │  │
│  │  ┌────────────────────────────────────────────┐     │  │
│  │  │                                            │     │  │
│  │  └────────────────────────────────────────────┘     │  │
│  │  From: Discord Developer Portal → Bot → Token      │  │
│  │                                                     │  │
│  │  Guild ID *                                          │  │
│  │  ┌────────────────────────────────────────────┐     │  │
│  │  │                                            │     │  │
│  │  └────────────────────────────────────────────┘     │  │
│  │  The Discord server (guild) where the bot lives     │  │
│  │                                                     │  │
│  │  Public Key *                                        │  │
│  │  ┌────────────────────────────────────────────┐     │  │
│  │  │                                            │     │  │
│  │  └────────────────────────────────────────────┘     │  │
│  │  From: Discord Developer Portal → General Info      │  │
│  │                                                     │  │
│  └─────────────────────────────────────────────────────┘  │
│                                                           │
│  [Approve User]  [Reject User]  [Cancel]                  │
└───────────────────────────────────────────────────────────┘
```

### 7.10 Shared components

These components are reused across pages:

| Component | Used on | Purpose |
|---|---|---|
| `Navbar` | Every page | Navigation bar with auth-aware links |
| `Footer` | Every page | Links to GitHub, copyright |
| `FormField` | Signup, Login, Profile, Admin | Label + input + error message |
| `StatusBadge` | Profile, Admin | Colored pill showing user status |
| `CredentialCard` | Profile | Displays one credential field with copy/mask |
| `ProtectedRoute` | All auth-required pages | Redirects to /login if no JWT |
| `AdminRoute` | All admin pages | Redirects to / if user is not admin |
| `ErrorBoundary` | App root | Catches React rendering errors gracefully |

### 7.11 Tailwind styling approach

We use Tailwind utility classes directly in JSX. There is no separate CSS file (beyond what Tailwind generates). A few conventions:

- **Color palette**: define a custom palette in `tailwind.config.ts` matching the brand (e.g., a purple/indigo primary, gray neutrals)
- **Dark mode**: default to dark mode for the bot-dev aesthetic
- **Responsive**: mobile-first, with `md:` and `lg:` breakpoints for tablets/desktop
- **Components**: use `@apply` sparingly — prefer inline utilities for clarity

Example button:
```tsx
<button className="bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-lg transition-colors">
  Sign Up
</button>
```

## 8. Authentication and Security

### 8.1 How authentication works

The platform uses **JSON Web Tokens (JWT)** for authentication. Here is the flow:

```
1. User submits email + password to POST /api/auth/login
2. Server verifies password against bcrypt hash in database
3. Server generates a JWT containing: { user_id, role, exp }
4. Server returns the JWT to the browser
5. Browser stores the JWT in localStorage
6. Every subsequent API request includes: Authorization: Bearer <token>
7. Server middleware validates the JWT signature and expiry
8. If valid, the handler receives the user_id and role in the request context
```

### 8.2 Why JWT (and not sessions)?

| Aspect | JWT | Server sessions |
|---|---|---|
| Server state | Stateless — no session store needed | Server must store session data |
| Scaling | Easy — any server instance can validate | Requires shared session store (Redis, DB) |
| Mobile-friendly | Yes — just send the token | Cookies can be tricky on mobile |
| Logout | Client-side delete (or short expiry) | Server destroys session |

For our single-server, SQLite-backed app, JWT is the simpler choice. We use **short-lived tokens** (24 hours) to limit the damage if a token is stolen.

### 8.3 Password storage

Passwords are **never stored in plain text**. We use bcrypt with a cost factor of 12:

```go
import "golang.org/x/crypto/bcrypt"

// Hashing a password on signup
hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

// Verifying a password on login
err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(inputPassword))
```

**Why bcrypt?** It is designed specifically for passwords. It is slow (by design), which makes brute-force attacks expensive. The cost factor of 12 means each hash takes ~250ms — fast enough for logins, slow enough to defeat attackers.

### 8.4 JWT implementation in Go

We use the `github.com/golang-jwt/jwt/v5` package:

```go
import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID int    `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, role string, secret []byte) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "bot-signup",
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
```

### 8.5 Auth middleware

Every route that requires authentication passes through middleware that:
1. Reads the `Authorization` header
2. Strips the `Bearer ` prefix
3. Parses and validates the JWT
4. Injects `user_id` and `role` into the request context

```go
func AuthMiddleware(secret []byte, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
            return
        }
        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            return secret, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
            return
        }
        // Inject user info into context
        ctx := context.WithValue(r.Context(), ctxKeyUserID, claims.UserID)
        ctx = context.WithValue(ctx, ctxKeyUserRole, claims.Role)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        role := r.Context().Value(ctxKeyUserRole).(string)
        if role != "admin" {
            http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    }
}
```

### 8.6 Frontend state management with RTK Query

Instead of writing manual `fetch()` calls and managing loading/error state in every component, we use **RTK Query** (part of Redux Toolkit). RTK Query auto-generates React hooks for each API endpoint, giving us caching, loading indicators, error handling, and optimistic updates out of the box.

**Why RTK Query?**

| Manual `fetch()` | RTK Query |
|---|---|
| Write `fetch()` in every component | Define API once, auto-generated hooks |
| Manually track `isLoading`, `error` states | Built-in `isLoading`, `isError`, `data` fields |
| No caching — duplicate requests everywhere | Automatic caching and deduplication |
| Manual refetch logic | `refetch()` on every hook result |
| No optimistic updates | Supports optimistic updates for snappy UX |

**How RTK Query is structured in this project:**

```
ui/src/store/
  ├── store.ts              # Redux store setup (configureStore)
  ├── api.ts                # RTK Query API slice definition (all endpoints)
  └── hooks.ts              # Re-exports all auto-generated hooks for convenience
```

**The API slice (`ui/src/store/api.ts`):**

```tsx
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api',
    prepareHeaders: (headers) => {
      const token = localStorage.getItem('token')
      if (token) {
        headers.set('Authorization', `Bearer ${token}`)
      }
      return headers
    },
  }),
  tagTypes: ['User', 'Profile', 'Waitlist', 'Credentials'],
  endpoints: (builder) => ({
    // Auth
    signup: builder.mutation<SignupResponse, SignupRequest>({
      query: (body) => ({ url: '/auth/signup', method: 'POST', body }),
    }),
    login: builder.mutation<LoginResponse, LoginRequest>({
      query: (body) => ({ url: '/auth/login', method: 'POST', body }),
    }),
    getMe: builder.query<User, void>({
      query: () => '/auth/me',
      providesTags: ['User'],
    }),

    // Profile
    getProfile: builder.query<ProfileResponse, void>({
      query: () => '/profile',
      providesTags: ['Profile', 'Credentials'],
    }),
    updateProfile: builder.mutation<User, Partial<User>>({
      query: (body) => ({ url: '/profile', method: 'PUT', body }),
      invalidatesTags: ['Profile'],
    }),
    changePassword: builder.mutation<void, ChangePasswordRequest>({
      query: (body) => ({ url: '/profile/password', method: 'PUT', body }),
    }),

    // Admin
    getWaitlist: builder.query<{ users: User[] }, void>({
      query: () => '/admin/waitlist',
      providesTags: ['Waitlist'],
    }),
    getAdminUsers: builder.query<PaginatedUsers, AdminUsersParams>({
      query: (params) => ({ url: '/admin/users', params }),
      providesTags: ['Waitlist'],
    }),
    approveUser: builder.mutation<ApproveResponse, ApproveRequest>({  
      query: ({ id, ...body }) => ({ url: `/admin/users/${id}/approve`, method: 'POST', body }),
      invalidatesTags: ['Waitlist', 'Profile'],
    }),
    rejectUser: builder.mutation<void, number>({
      query: (id) => ({ url: `/admin/users/${id}/reject`, method: 'POST' }),
      invalidatesTags: ['Waitlist'],
    }),
    suspendUser: builder.mutation<void, number>({
      query: (id) => ({ url: `/admin/users/${id}/suspend`, method: 'POST' }),
      invalidatesTags: ['Waitlist', 'Profile'],
    }),
    updateCredentials: builder.mutation<void, { id: number } & CredentialFields>({
      query: ({ id, ...body }) => ({ url: `/admin/users/${id}/credentials`, method: 'PUT', body }),
      invalidatesTags: ['Credentials'],
    }),
    deleteUser: builder.mutation<void, number>({
      query: (id) => ({ url: `/admin/users/${id}`, method: 'DELETE' }),
      invalidatesTags: ['Waitlist'],
    }),

    // Public
    getStats: builder.query<Stats, void>({
      query: () => '/stats',
    }),
    getHealth: builder.query<Health, void>({
      query: () => '/health',
    }),
  }),
})

// Auto-generated hooks: useSignupMutation, useLoginMutation, useGetProfileQuery, etc.
export const {
  useSignupMutation,
  useLoginMutation,
  useGetMeQuery,
  useGetProfileQuery,
  useUpdateProfileMutation,
  useChangePasswordMutation,
  useGetWaitlistQuery,
  useGetAdminUsersQuery,
  useApproveUserMutation,
  useRejectUserMutation,
  useSuspendUserMutation,
  useUpdateCredentialsMutation,
  useDeleteUserMutation,
  useGetStatsQuery,
} = apiSlice
```

**Using RTK Query hooks in a component:**

```tsx
// Example: ProfilePage.tsx
function ProfilePage() {
  // RTK Query gives us data, isLoading, error, and refetch automatically
  const { data, isLoading, error } = useGetProfileQuery()
  const [updateProfile] = useUpdateProfileMutation()

  if (isLoading) return <LoadingSpinner />
  if (error) return <div>Failed to load profile</div>

  return (
    <div>
      <h1>Welcome, {data.user.display_name}</h1>
      {data.bot_credentials && <CredentialCard credential={data.bot_credentials} />}
    </div>
  )
}
```

**Auth context still handles token storage and redirect logic**, but delegates all API calls to RTK Query:

```tsx
function AuthProvider({ children }) {
  const [loginMutation] = useLoginMutation()
  const [signupMutation] = useSignupMutation()

  const login = async (email: string, password: string) => {
    const res = await loginMutation({ email, password }).unwrap()
    localStorage.setItem('token', res.token)
    setUser(res.user)
  }

  const signup = async (data: SignupData) => {
    const res = await signupMutation(data).unwrap()
    localStorage.setItem('token', res.token)
    setUser(res.user)
  }

  // ...
}
```

### 8.7 Security checklist

- [x] Passwords hashed with bcrypt (cost 12)
- [x] JWT tokens expire after 24 hours
- [x] JWT secret loaded from environment variable (never hard-coded)
- [x] HTTPS enforced in production
- [x] Bot tokens stored in database (never logged or exposed in error messages)
- [x] Admin routes protected with role check middleware
- [x] Input validation on all endpoints
- [x] CORS restricted to same-origin in production
- [x] SQL injection prevented via parameterized queries
- [x] Rate limiting on login endpoint (5 attempts per minute per IP)

## 9. Admin Backend — Detailed Design

### 9.1 How does someone become an admin?

The first admin is created manually — either via a CLI command or by inserting directly into the database:

```bash
# Option A: CLI command
./bot-signup admin create --email admin@example.com --password secret123

# Option B: Direct SQL (for development)
sqlite3 data/bot-signup.db "UPDATE users SET role='admin' WHERE email='admin@example.com'"
```

Only users with `role='admin'` can access `/api/admin/*` endpoints.

### 9.2 What the admin workflow looks like

```
Admin opens /admin
    │
    ▼
Sees stats bar (total, waiting, approved)
    │
    ▼
Sees waiting-list table
    │
    ├── Clicks "Approve" on user #5
    │       │
    │       ▼
    │   Approval form opens (modal or inline)
    │       │
    │       ▼
    │   Fills in 4 Discord credential fields:
    │     • Application ID  (from Discord Developer Portal)
    │     • Bot Token       (from Discord Developer Portal → Bot)
    │     • Guild ID        (right-click server → Copy Server ID)
    │     • Public Key      (from Discord Developer Portal → General)
    │       │
    │       ▼
    │   Clicks "Approve User"
    │       │
    │       ▼
    │   Server: validates fields → inserts bot_credentials row
    │          → updates user status to 'approved'
    │       │
    │       ▼
    │   Admin sees success toast → user disappears from waiting list
    │
    ├── Clicks "Reject" on user #6
    │       │
    │       ▼
    │   Confirmation dialog: "Reject this user?"
    │       │
    │       ▼
    │   User status → 'rejected'
    │
    └── Clicks "View All Users" → paginated user list
            │
            ▼
        Can filter by status, search by name/email
        Can edit credentials, suspend, or delete users
```

### 9.3 Discord credentials — what they are and where they come from

The admin needs to create a Discord Application for each user's bot. Here is what each field means:

| Field | What it is | Where to find it |
|---|---|---|
| **Application ID** | A unique snowflake ID for the Discord application | Discord Developer Portal → General Information → APPLICATION ID |
| **Bot Token** | A secret string that authenticates the bot with Discord's API | Discord Developer Portal → Bot → Token → "Reset Token" (shown once) |
| **Guild ID** | The ID of the Discord server where the bot will operate | Discord client: Server Settings → Widget → Server ID, or right-click server icon → Copy Server ID (requires Developer Mode) |
| **Public Key** | Used to verify incoming HTTP interactions from Discord | Discord Developer Portal → General Information → PUBLIC KEY |

### 9.4 Admin approval handler pseudocode

```go
func (s *Server) handleApproveUser(w http.ResponseWriter, r *http.Request) {
    // 1. Extract user_id from URL path
    userID, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        respondError(w, 400, "invalid user id")
        return
    }

    // 2. Parse request body
    var req ApproveRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, 400, "invalid request body")
        return
    }

    // 3. Validate all four credential fields
    if req.ApplicationID == "" || req.BotToken == "" ||
       req.GuildID == "" || req.PublicKey == "" {
        respondError(w, 400, "all credential fields are required")
        return
    }

    // 4. Check user exists and is in 'waiting' status
    user, err := s.db.GetUserByID(r.Context(), userID)
    if err != nil {
        respondError(w, 404, "user not found")
        return
    }
    if user.Status != "waiting" {
        respondError(w, 409, "user is not in waiting status")
        return
    }

    // 5. Get admin user_id from context (set by auth middleware)
    adminID := r.Context().Value(ctxKeyUserID).(int)

    // 6. Start a transaction (update user + insert credentials atomically)
    tx, err := s.db.BeginTx(r.Context())
    if err != nil {
        respondError(w, 500, "database error")
        return
    }
    defer tx.Rollback()

    // 7. Update user status
    err = tx.UpdateUserStatus(r.Context(), userID, "approved")
    if err != nil {
        respondError(w, 500, "failed to update user status")
        return
    }

    // 8. Insert bot credentials
    err = tx.InsertBotCredentials(r.Context(), &BotCredentials{
        UserID:         userID,
        ApplicationID:  req.ApplicationID,
        BotToken:       req.BotToken,
        GuildID:        req.GuildID,
        PublicKey:      req.PublicKey,
        ApprovedBy:     adminID,
        ApprovedAt:     time.Now(),
    })
    if err != nil {
        respondError(w, 500, "failed to insert credentials")
        return
    }

    // 9. Commit transaction
    if err := tx.Commit(); err != nil {
        respondError(w, 500, "failed to commit transaction")
        return
    }

    // 10. Return success
    respondJSON(w, 200, map[string]interface{}{
        "message": "User approved successfully",
    })
}
```

### 9.5 What happens when an admin edits credentials?

The admin can update the four Discord fields at any time (e.g., if the user needs a new bot token). This is a simple UPDATE on the `bot_credentials` row:

```sql
UPDATE bot_credentials
SET application_id = ?, bot_token = ?, guild_id = ?, public_key = ?, updated_at = datetime('now')
WHERE user_id = ?
```

### 9.6 What happens when an admin suspends a user?

Suspension sets the user's status back to `suspended`. The user can still log in but sees a "suspended" message. The bot credentials are NOT deleted (so the admin can re-approve later):

```sql
UPDATE users SET status = 'suspended', updated_at = datetime('now') WHERE id = ?
```

## 10. Tutorial Content — The Discord Bot Guide

### 10.1 Where the tutorial comes from

The tutorial is based on the official documentation from [github.com/go-go-golems/discord-bot](https://github.com/go-go-golems/discord-bot). Specifically:

- The `README.md` — covers installation, quick start, Go API, and JS bot authoring
- `pkg/doc/tutorials/building-and-running-discord-js-bots.md` — the full step-by-step tutorial
- `examples/discord-bots/ping/index.js` — the richest example bot

### 10.2 What the tutorial page shows

The tutorial page at `/tutorial` is a **rendered Markdown document** embedded in the React frontend. It covers:

1. **What is discord-bot?** — A Go runtime that hosts JavaScript Discord bots. You write bot logic in JS; the Go process handles the Discord gateway, authentication, and reconnection.

2. **Prerequisites** — What you need before starting:
   - Your Discord bot credentials (from your profile page)
   - A Discord server where you have admin permissions
   - Basic JavaScript knowledge

3. **Installation** — How to install the `discord-bot` binary:
   ```bash
   brew install go-go-golems/tap/discord-bot
   ```

4. **Your first bot** — A minimal bot with a `/ping` command:
   ```js
   const { defineBot } = require("discord")

   module.exports = defineBot(({ command, event, configure }) => {
     configure({
       name: "my-first-bot",
       description: "My first Discord bot",
     })

     command("hello", {
       description: "Say hello!",
     }, async () => {
       return { content: "Hello from my first bot! 🎉" }
     })

     event("ready", async (ctx) => {
       ctx.log.info("Bot is ready!", { user: ctx.me.username })
     })
   })
   ```

5. **Running your bot** — Using the credentials from your profile:
   ```bash
   export DISCORD_BOT_TOKEN="your-token-here"
   export DISCORD_APPLICATION_ID="your-app-id-here"
   export DISCORD_GUILD_ID="your-guild-id-here"

   discord-bot bots my-first-bot run \
     --bot-repository ./my-bots \
     --sync-on-start
   ```

6. **Adding commands** — How to add slash commands with options, autocomplete, and deferred responses.

7. **Adding buttons and modals** — Interactive components for richer bot experiences.

8. **Using the store** — The built-in key-value store for bot state (`ctx.store.get/set`).

9. **Discord API operations** — Using `ctx.discord.channels.*`, `ctx.discord.messages.*`, etc.

10. **Next steps** — Links to the full API reference and more example bots.

### 10.3 How the tutorial is embedded in the frontend

The tutorial content is stored as a Markdown file in the frontend source tree:

```
ui/src/content/tutorial.md
```

At build time, Vite imports this as a raw string. A React component renders it using a Markdown renderer (e.g., `react-markdown` with syntax highlighting):

```tsx
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import tutorialContent from '../content/tutorial.md?raw'

export function TutorialPage() {
  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>
        {tutorialContent}
      </ReactMarkdown>
    </div>
  )
}
```

### 10.4 JavaScript bot API quick reference

For the tutorial page sidebar, we include a quick reference of the `require("discord")` API:

| Function | Purpose |
|---|---|
| `defineBot(factory)` | Entry point — defines the entire bot |
| `command(name, spec, handler)` | Register a slash command |
| `event(name, handler)` | Register a Discord event handler |
| `component(customId, handler)` | Register a button/select menu handler |
| `modal(customId, handler)` | Register a modal submit handler |
| `autocomplete(cmdName, handler)` | Register an autocomplete handler |
| `configure(spec)` | Set bot metadata and runtime config fields |

**Context (`ctx`) methods available in handlers:**

| Method | Purpose |
|---|---|
| `ctx.reply(payload)` | Respond to an interaction |
| `ctx.defer(opts)` | Acknowledge with a "thinking" state |
| `ctx.edit(payload)` | Edit a deferred response |
| `ctx.followUp(payload)` | Send a follow-up message |
| `ctx.showModal(spec)` | Open a modal dialog |
| `ctx.log.info/warn/error/debug(msg, data)` | Structured logging |
| `ctx.store.get/set/delete/keys()` | Key-value store operations |
| `ctx.discord.channels/messages/members/guilds/roles/threads.*` | Discord API operations |
| `ctx.config.*` | Runtime config values |
| `ctx.me` | Current bot user object |
| `ctx.args.*` | Parsed command arguments |
| `ctx.values.*` | Modal submitted values |

## 11. Project Structure — Every File and What It Does

### 11.1 Complete directory tree

```
bot-signup/
├── cmd/
│   └── bot-signup/
│       └── main.go                    # CLI entrypoint (serve, admin, migrate commands)
│
├── internal/
│   ├── database/
│   │   ├── database.go               # DB struct, Open(), migration runner
│   │   ├── users.go                  # User CRUD operations
│   │   ├── credentials.go            # Bot credentials CRUD operations
│   │   └── migrations/
│   │       └── 001_initial.sql        # Schema migration
│   │
│   ├── server/
│   │   ├── server.go                 # Server struct, New(), route registration
│   │   ├── auth_handlers.go          # POST /api/auth/* handlers
│   │   ├── profile_handlers.go       # GET/PUT /api/profile handlers
│   │   ├── admin_handlers.go         # GET/POST /api/admin/* handlers
│   │   ├── middleware.go             # AuthMiddleware, AdminOnly, CORS, logging
│   │   └── helpers.go                # respondJSON, respondError utilities
│   │
│   ├── auth/
│   │   ├── jwt.go                    # JWT generation and parsing
│   │   └── password.go               # bcrypt hashing and verification
│   │
│   └── web/
│       ├── embed.go                  # //go:build embed — embeds the frontend FS
│       ├── embed_none.go             # //go:build !embed — reads from disk
│       ├── spa.go                    # SPA handler (serve files, fallback to index.html)
│       ├── generate.go               # //go:generate directive
│       └── generate_build.go         # Builds React + copies to internal/web/embed/
│
├── ui/
│   ├── index.html                    # Vite HTML entrypoint
│   ├── vite.config.ts                # Vite config (dev proxy to :8080)
│   ├── tailwind.config.ts            # Tailwind configuration
│   ├── postcss.config.js             # PostCSS (needed by Tailwind)
│   ├── tsconfig.json                 # TypeScript configuration
│   ├── package.json                  # Frontend dependencies
│   ├── pnpm-lock.yaml                # Locked dependencies
│   │
│   ├── .storybook/
│   │   ├── main.ts                   # Storybook config (Vite builder, Tailwind)
│   │   ├── preview.ts                # Global decorators (Tailwind CSS import)
│   │   └── README.md                 # How to run Storybook
│   │
│   └── src/
│       ├── main.tsx                  # React entrypoint
│       ├── App.tsx                   # Router setup, Redux Provider + auth provider
│       ├── vite-env.d.ts             # Vite type declarations
│       │
│       ├── store/
│       │   ├── store.ts              # Redux store setup (configureStore + api middleware)
│       │   ├── api.ts                # RTK Query API slice (all endpoints + auto hooks)
│       │   └── authSlice.ts           # Auth state slice (user, token, login/logout actions)
│       │
│       ├── auth/
│       │   ├── AuthContext.tsx        # React context: login/logout/signup using RTK Query hooks
│       │   ├── useAuth.ts            # Hook: user, token, login, logout, signup
│       │   └── ProtectedRoute.tsx     # Redirects to /login if not authenticated
│       │
│       ├── pages/
│       │   ├── LandingPage.tsx        # / — hero, features, stats
│       │   ├── SignupPage.tsx         # /signup — signup form
│       │   ├── LoginPage.tsx          # /login — login form
│       │   ├── WaitingListPage.tsx    # /waiting-list — status display
│       │   ├── ProfilePage.tsx        # /profile — user info + credentials
│       │   ├── ChangePasswordPage.tsx  # /profile/password
│       │   ├── TutorialPage.tsx       # /tutorial — rendered markdown
│       │   ├── NotFoundPage.tsx       # /* — 404 page
│       │   └── admin/
│       │       ├── AdminDashboard.tsx  # /admin — overview + waiting list
│       │       ├── AdminWaitlist.tsx   # /admin/waitlist
│       │       ├── AdminUserDetail.tsx # /admin/users/{id} — approve/reject/edit
│       │       └── AdminRoute.tsx      # Redirects non-admins to /
│       │
│       ├── components/
│       │   ├── Navbar.tsx             # Top navigation bar
│       │   ├── Footer.tsx             # Bottom footer
│       │   ├── FormField.tsx          # Reusable form input with label + error
│       │   ├── StatusBadge.tsx        # Colored status pill
│       │   ├── CredentialCard.tsx     # Single credential display with copy/mask
│       │   ├── ErrorBoundary.tsx      # Catches rendering errors
│       │   └── LoadingSpinner.tsx     # Loading indicator
│       │
│       ├── components/**/*.stories.tsx  # Storybook stories (one per component)
│       │   ├── Navbar.stories.tsx
│       │   ├── Footer.stories.tsx
│       │   ├── FormField.stories.tsx
│       │   ├── StatusBadge.stories.tsx
│       │   ├── CredentialCard.stories.tsx
│       │   ├── LoadingSpinner.stories.tsx
│       │   ├── SignupForm.stories.tsx
│       │   ├── LoginForm.stories.tsx
│       │   ├── WaitingListStatus.stories.tsx
│       │   ├── ProfileCard.stories.tsx
│       │   ├── AdminUserTable.stories.tsx
│       │   └── ApprovalForm.stories.tsx
│       │
│       └── content/
│           └── tutorial.md            # Tutorial content (imported as raw string)
│
├── data/                              # Git-ignored directory for SQLite + uploads
│   └── .gitkeep
│
├── Makefile                            # Developer convenience targets
├── go.mod                              # Go module definition
├── go.sum                              # Go dependency checksums
├── .gitignore                          # Ignore data/, ui/dist/, ui/node_modules/
├── .github/
│   └── workflows/
│       └── ci.yml                     # GitHub Actions: test, lint, build
│
└── README.md                          # Project README with setup instructions
```

### 11.2 Key files explained (read these first)

If you are new to the project, read files in this order:

| Order | File | Why read it |
|---|---|---|
| 1 | `cmd/bot-signup/main.go` | The entrypoint. Shows how the CLI is wired (Cobra commands). |
| 2 | `internal/server/server.go` | The HTTP server. Shows every route registration. |
| 3 | `internal/database/database.go` | Database initialization and migration. |
| 4 | `internal/database/users.go` | User CRUD — the core data layer. |
| 5 | `internal/server/auth_handlers.go` | Signup and login logic — the main user flow. |
| 6 | `ui/src/App.tsx` | Frontend routing and auth context setup. |
| 7 | `ui/src/pages/SignupPage.tsx` | The signup form — typical React page pattern. |

## 12. Implementation Phases — What to Build and In What Order

### Phase 1: Project scaffolding (Day 1)

**Goal**: A Go binary that starts an HTTP server and returns "hello world" on `/api/health`.

**Steps:**
1. Initialize Go module: `go mod init github.com/go-go-golems/bot-signup`
2. Create `cmd/bot-signup/main.go` with a `serve` command using Cobra
3. Create `internal/server/server.go` with a basic `ServeMux` and `/api/health` handler
4. Create `Makefile` with `dev-backend` target
5. Verify: `make dev-backend` starts, `curl localhost:8080/api/health` returns `{"status":"ok"}`

**Files created:**
```
cmd/bot-signup/main.go
internal/server/server.go
internal/server/helpers.go
Makefile
go.mod
```

---

### Phase 2: Database layer (Day 1–2)

**Goal**: SQLite database initializes on startup, migration runs, and we can CRUD users.

**Steps:**
1. Add `modernc.org/sqlite` dependency
2. Create `internal/database/database.go` with `Open()` that:
   - Creates the `data/` directory if needed
   - Opens SQLite with WAL mode and foreign keys enabled
   - Runs migrations from embedded SQL files
3. Create `internal/database/migrations/001_initial.sql` with the schema from Section 5.3
4. Create `internal/database/users.go` with functions:
   - `CreateUser(ctx, discordID, email, displayName, passwordHash) (*User, error)`
   - `GetUserByID(ctx, id) (*User, error)`
   - `GetUserByEmail(ctx, email) (*User, error)`
   - `GetUserByDiscordID(ctx, discordID) (*User, error)`
   - `UpdateUserStatus(ctx, id, status) error`
   - `ListUsersByStatus(ctx, status, page, perPage) ([]*User, int, error)`
   - `DeleteUser(ctx, id) error`
5. Create `internal/database/credentials.go` with functions:
   - `InsertBotCredentials(ctx, creds) error`
   - `GetCredentialsByUserID(ctx, userID) (*BotCredentials, error)`
   - `UpdateBotCredentials(ctx, creds) error`
6. Write table-driven tests for each function using an in-memory SQLite database

**Files created:**
```
internal/database/database.go
internal/database/users.go
internal/database/credentials.go
internal/database/migrations/001_initial.sql
internal/database/users_test.go
internal/database/credentials_test.go
```

---

### Phase 3: Authentication (Day 2–3)

**Goal**: Signup and login endpoints work, returning JWTs.

**Steps:**
1. Add `golang.org/x/crypto` (bcrypt) and `github.com/golang-jwt/jwt/v5` dependencies
2. Create `internal/auth/password.go` with `HashPassword()` and `CheckPassword()`
3. Create `internal/auth/jwt.go` with `GenerateToken()` and `ParseToken()`
4. Create `internal/server/middleware.go` with `AuthMiddleware()` and `AdminOnly()`
5. Create `internal/server/auth_handlers.go` with:
   - `POST /api/auth/signup` — validates input, hashes password, creates user, returns JWT
   - `POST /api/auth/login` — checks password, returns JWT
   - `POST /api/auth/logout` — returns success (client deletes token)
   - `GET /api/auth/me` — returns current user (requires auth)
6. Test with `curl`:
   ```bash
   # Signup
   curl -X POST http://localhost:8080/api/auth/signup \
     -H 'Content-Type: application/json' \
     -d '{"discord_id":"123","email":"test@test.com","display_name":"Test","password":"password123"}'

   # Login
   curl -X POST http://localhost:8080/api/auth/login \
     -H 'Content-Type: application/json' \
     -d '{"email":"test@test.com","password":"password123"}'

   # Get profile (use token from login response)
   curl http://localhost:8080/api/auth/me \
     -H 'Authorization: Bearer <token>'
   ```

**Files created:**
```
internal/auth/password.go
internal/auth/jwt.go
internal/server/middleware.go
internal/server/auth_handlers.go
```

---

### Phase 4: Profile and admin handlers (Day 3)

**Goal**: All remaining API endpoints work.

**Steps:**
1. Create `internal/server/profile_handlers.go`:
   - `GET /api/profile` — returns user + credentials (if approved)
   - `PUT /api/profile` — update display name / email
   - `PUT /api/profile/password` — change password
2. Create `internal/server/admin_handlers.go`:
   - `GET /api/admin/waitlist` — list waiting users
   - `GET /api/admin/users` — list all users (paginated)
   - `POST /api/admin/users/{id}/approve` — approve with credentials
   - `POST /api/admin/users/{id}/reject` — reject
   - `POST /api/admin/users/{id}/suspend` — suspend
   - `PUT /api/admin/users/{id}/credentials` — update credentials
   - `DELETE /api/admin/users/{id}` — delete user
3. Wire all routes in `server.go`
4. Add `GET /api/stats` public endpoint
5. Test every endpoint with `curl` or a simple test script

**Files created:**
```
internal/server/profile_handlers.go
internal/server/admin_handlers.go
```

---

### Phase 5: Frontend scaffolding (Day 4)

**Goal**: Vite + React + Tailwind + Storybook + RTK Query app renders a landing page.

**Steps:**
1. Create `ui/` directory with Vite + React + TypeScript template:
   ```bash
   pnpm create vite ui --template react-ts
   ```
2. Install dependencies:
   ```bash
   cd ui
   pnpm add react-router-dom react-markdown remark-gfm @reduxjs/toolkit react-redux
   pnpm add -D tailwindcss @tailwindcss/vite
   pnpm add -D @storybook/react-vite @storybook/builder-vite storybook
   ```
3. Configure Tailwind in `vite.config.ts` and add `@import "tailwindcss"` to main CSS
4. Set up dev proxy in `vite.config.ts`:
   ```ts
   export default defineConfig({
     plugins: [react(), tailwindcss()],
     server: {
       proxy: {
         '/api': 'http://localhost:8080',
       },
     },
   })
   ```
5. Initialize Storybook:
   ```bash
   pnpm storybook init
   ```
   This creates `ui/.storybook/main.ts` and `ui/.storybook/preview.ts`. Verify Tailwind classes work in stories by importing the Tailwind CSS in `preview.ts`:
   ```ts
   // .storybook/preview.ts
   import '../src/index.css'
   export const parameters = { ... }
   ```
6. Create `ui/src/store/store.ts` — Redux store with RTK Query middleware:
   ```ts
   import { configureStore } from '@reduxjs/toolkit'
   import { apiSlice } from './api'
   import authReducer from './authSlice'

   export const store = configureStore({
     reducer: {
       [apiSlice.reducerPath]: apiSlice.reducer,
       auth: authReducer,
     },
     middleware: (getDefault) =>
       getDefault().concat(apiSlice.middleware),
   })
   ```
7. Create `ui/src/store/api.ts` — RTK Query API slice with all endpoints (see Section 8.6)
8. Create `ui/src/store/authSlice.ts` — minimal auth slice for user/token state:
   ```ts
   import { createSlice } from '@reduxjs/toolkit'
   const authSlice = createSlice({
     name: 'auth',
     initialState: { user: null, token: localStorage.getItem('token') },
     reducers: {
       setCredentials: (state, action) => {
         state.user = action.payload.user
         state.token = action.payload.token
         localStorage.setItem('token', action.payload.token)
       },
       logout: (state) => {
         state.user = null
         state.token = null
         localStorage.removeItem('token')
       },
     },
   })
   ```
9. Wrap `<App />` with Redux `<Provider>` in `main.tsx`:
   ```tsx
   import { Provider } from 'react-redux'
   import { store } from './store/store'
   createRoot(document.getElementById('root')!).render(
     <Provider store={store}>
       <App />
     </Provider>
   )
   ```
10. Create `Navbar`, `Footer`, and `LandingPage` components
11. **Create first Storybook story** alongside each component:
    ```tsx
    // ui/src/components/Navbar.stories.tsx
    import type { Meta, StoryObj } from '@storybook/react-vite'
    import { Navbar } from './Navbar'

    const meta: Meta<typeof Navbar> = {
      title: 'Components/Navbar',
      component: Navbar,
      args: { isLoggedIn: false },
      argTypes: { isLoggedIn: { control: 'boolean' } },
    }
    export default meta
    type Story = StoryObj<typeof Navbar>

    export const LoggedOut: Story = { args: { isLoggedIn: false } }
    export const LoggedIn: Story = { args: { isLoggedIn: true, displayName: 'CoolBotDev' } }
    ```
12. Verify: `cd ui && pnpm dev` shows landing page, `pnpm storybook` shows stories at `http://localhost:6006`
13. Add `make dev-frontend` and `make storybook` targets to Makefile

**Files created:**
```
ui/* (entire Vite project)
ui/.storybook/main.ts
ui/.storybook/preview.ts
ui/src/store/store.ts
ui/src/store/api.ts
ui/src/store/authSlice.ts
ui/src/App.tsx
ui/src/main.tsx (with Provider wrapper)
ui/src/components/Navbar.tsx
ui/src/components/Navbar.stories.tsx
ui/src/components/Footer.tsx
ui/src/components/Footer.stories.tsx
ui/src/pages/LandingPage.tsx
```

---

### Phase 6: Auth pages (Day 4–5)

**Goal**: Signup, login, and auth context work end-to-end using RTK Query mutations.

**Steps:**
1. Create `ui/src/auth/AuthContext.tsx` — uses `useLoginMutation` and `useSignupMutation` from RTK Query, dispatches `setCredentials`/`logout` to the auth slice
2. Create `ui/src/auth/ProtectedRoute.tsx` — route guard component
3. Create `FormField.tsx` reusable component:
   ```tsx
   // ui/src/components/FormField.tsx
   interface FormFieldProps {
     label: string
     type?: string
     value: string
     onChange: (v: string) => void
     error?: string
     hint?: string
   }
   ```
4. **Create `FormField.stories.tsx`** — stories for default, with-error, with-hint, password type:
   ```tsx
   export const Default: Story = { args: { label: 'Email', value: '' } }
   export const WithError: Story = { args: { label: 'Email', value: 'bad', error: 'Invalid email' } }
   export const Password: Story = { args: { label: 'Password', type: 'password' } }
   ```
5. Create `SignupPage.tsx` using `useSignupMutation()`:
   ```tsx
   const [signup, { isLoading, error }] = useSignupMutation()
   const handleSubmit = async (e) => {
     const res = await signup(formData).unwrap()
     dispatch(setCredentials(res))
     navigate('/waiting-list')
   }
   ```
6. **Create `SignupForm.stories.tsx`** — stories for empty, partially filled, submitting state, error state
7. Create `LoginPage.tsx` using `useLoginMutation()`
8. **Create `LoginForm.stories.tsx`** — stories for empty and error states
9. Verify full flow: signup → redirect to waiting list → logout → login → see waiting list

**Storybook stories created this phase:**
```
ui/src/components/FormField.stories.tsx
ui/src/components/SignupForm.stories.tsx
ui/src/components/LoginForm.stories.tsx
```

---

### Phase 7: User pages (Day 5)

**Goal**: Profile, waiting list, change password pages work with RTK Query hooks.

**Steps:**
1. Create `StatusBadge.tsx` — colored status pill:
   ```tsx
   const colors = {
     waiting: 'bg-yellow-100 text-yellow-800',
     approved: 'bg-green-100 text-green-800',
     rejected: 'bg-red-100 text-red-800',
     suspended: 'bg-gray-100 text-gray-800',
   }
   ```
2. **Create `StatusBadge.stories.tsx`** — one story per status value:
   ```tsx
   export const Waiting: Story = { args: { status: 'waiting' } }
   export const Approved: Story = { args: { status: 'approved' } }
   export const Rejected: Story = { args: { status: 'rejected' } }
   export const Suspended: Story = { args: { status: 'suspended' } }
   ```
3. Create `CredentialCard.tsx` — credential display with mask/copy:
   ```tsx
   interface CredentialCardProps {
     label: string
     value: string
     maskable?: boolean
   }
   ```
4. **Create `CredentialCard.stories.tsx`** — stories for normal, masked, revealed:
   ```tsx
   export const Normal: Story = { args: { label: 'Guild ID', value: '111222333' } }
   export const Masked: Story = { args: { label: 'Bot Token', value: 'MTIz...', maskable: true } }
   ```
5. Create `WaitingListPage.tsx` using `useGetMeQuery()` to check status:
   ```tsx
   const { data: user } = useGetMeQuery()
   // Show different content based on user.status
   ```
6. **Create `WaitingListStatus.stories.tsx`** — stories for waiting, rejected, suspended states
7. Create `ProfilePage.tsx` using `useGetProfileQuery()`:
   ```tsx
   const { data, isLoading } = useGetProfileQuery()
   ```
8. **Create `ProfileCard.stories.tsx`** — stories for approved (with credentials) and waiting user
9. Create `ChangePasswordPage.tsx` using `useChangePasswordMutation()`
10. Verify: all three pages work end-to-end

**Storybook stories created this phase:**
```
ui/src/components/StatusBadge.stories.tsx
ui/src/components/CredentialCard.stories.tsx
ui/src/components/WaitingListStatus.stories.tsx
ui/src/components/ProfileCard.stories.tsx
```

---

### Phase 8: Admin pages (Day 5–6)

**Goal**: Admin can view waiting list, approve users, and manage all users via RTK Query hooks.

**Steps:**
1. Create `AdminRoute.tsx` — route guard for admin role
2. Create `AdminDashboard.tsx` using `useGetWaitlistQuery()` and `useGetStatsQuery()`:
   ```tsx
   const { data: waitlist } = useGetWaitlistQuery()
   const { data: stats } = useGetStatsQuery()
   // Stats bar + waiting-list table
   ```
3. **Create `AdminUserTable.stories.tsx`** — stories for empty list, single user, multiple users:
   ```tsx
   export const Empty: Story = { args: { users: [] } }
   export const WithUsers: Story = {
     args: {
       users: [
         { id: 1, discord_id: '123', display_name: 'CoolBotDev', email: 'user@ex.com', created_at: '...' },
         { id: 2, discord_id: '456', display_name: 'BotMaster', email: 'bot@ex.com', created_at: '...' },
       ]
     }
   }
   ```
4. Create `AdminUserDetail.tsx` using `useApproveUserMutation()` and `useUpdateCredentialsMutation()`:
   ```tsx
   const [approve, { isLoading: approving }] = useApproveUserMutation()
   const handleApprove = async () => {
     await approve({ id: userId, ...credentials }).unwrap()
     navigate('/admin')
   }
   ```
5. **Create `ApprovalForm.stories.tsx`** — stories for empty form, partially filled, validation errors:
   ```tsx
   export const Empty: Story = { args: { onSubmit: fn() } }
   export const WithValidationErrors: Story = {
     args: { onSubmit: fn(), errors: { application_id: 'Required' } }
   }
   ```
6. Create `AdminWaitlist.tsx` — dedicated waiting-list management
7. Verify full admin flow end-to-end

**Storybook stories created this phase:**
```
ui/src/components/AdminUserTable.stories.tsx
ui/src/components/ApprovalForm.stories.tsx
```

---

### Phase 9: Tutorial page (Day 6)

**Goal**: The tutorial renders from embedded Markdown.

**Steps:**
1. Write `ui/src/content/tutorial.md` based on the discord-bot tutorial
2. Create `TutorialPage.tsx` with react-markdown rendering
3. Add code syntax highlighting (e.g., `rehype-highlight`)
4. Style with Tailwind for readability

---

### Phase 10: Frontend embedding (Day 6–7)

**Goal**: `go generate` builds the React app and embeds it in the Go binary.

**Steps:**
1. Create `internal/web/embed.go` (`//go:build embed`)
2. Create `internal/web/embed_none.go` (`//go:build !embed`)
3. Create `internal/web/spa.go` — SPA handler
4. Create `internal/web/generate.go` and `generate_build.go`
5. Wire SPA handler into server.go (last route, catches everything else)
6. Test: `go generate ./internal/web/ && go build -tags embed -o bot-signup ./cmd/bot-signup`
7. Run `./bot-signup serve` and verify the frontend loads

---

### Phase 11: Polish and deploy (Day 7)

**Goal**: Error handling, loading states, edge cases, CI, Storybook build.

**Steps:**
1. Add `ErrorBoundary` component
2. Add loading spinners on all async operations (RTK Query `isLoading` states)
3. Add form validation feedback (client-side)
4. Add `ErrorBoundary` for unexpected crashes
5. Create `.github/workflows/ci.yml` for automated testing + Storybook build
6. Create admin CLI command: `./bot-signup admin create --email ... --password ...`
7. Add `make storybook-build` target for static Storybook output (deploy to GitHub Pages or similar)
8. Final end-to-end walkthrough
9. Write `README.md` with setup instructions
10. **Verify all Storybook stories render correctly**: `cd ui && pnpm storybook` → check all components

**Makefile targets to have:**
```makefile
dev-backend:     # go run ./cmd/bot-signup serve
dev-frontend:   # cd ui && pnpm dev
storybook:       # cd ui && pnpm storybook
storybook-build: # cd ui && pnpm build-storybook
build:           # go generate + go build -tags embed
```

## 13. Pseudocode for Key Flows

### 13.1 Server startup

```go
func main() {
    // Parse CLI flags
    cmd :=cobra.Command{Use: "bot-signup"}
    serveCmd := &cobra.Command{
        Use: "serve",
        Run: func(cmd *cobra.Command, args []string) {
            // 1. Load config from env vars
            jwtSecret := os.Getenv("JWT_SECRET")     // required
            dbPath := os.Getenv("DB_PATH")           // default: ./data/bot-signup.db
            port := os.Getenv("PORT")                // default: 8080

            // 2. Open database + run migrations
            db, err := database.Open(dbPath)
            if err != nil { log.Fatal(err) }
            defer db.Close()

            // 3. Create server
            srv := server.New(db, []byte(jwtSecret))

            // 4. Create ServeMux and register routes
            mux := http.NewServeMux()
            srv.RegisterRoutes(mux) // registers /api/* and SPA

            // 5. Start server
            addr := fmt.Sprintf(":%s", port)
            log.Printf("Server starting on %s", addr)
            log.Fatal(http.ListenAndServe(addr, mux))
        },
    }
    cmd.AddCommand(serveCmd)
    cmd.Execute()
}
```

### 13.2 Route registration

```go
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
    // ── Public endpoints ──────────────────────────────
    mux.HandleFunc("GET /api/health", s.handleHealth)
    mux.HandleFunc("GET /api/stats", s.handleStats)

    // ── Auth endpoints ────────────────────────────────
    mux.HandleFunc("POST /api/auth/signup", s.handleSignup)
    mux.HandleFunc("POST /api/auth/login", s.handleLogin)
    mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
    mux.HandleFunc("GET /api/auth/me", AuthMiddleware(s.secret, s.handleMe))

    // ── Profile endpoints (authenticated) ─────────────
    mux.HandleFunc("GET /api/profile", AuthMiddleware(s.secret, s.handleGetProfile))
    mux.HandleFunc("PUT /api/profile", AuthMiddleware(s.secret, s.handleUpdateProfile))
    mux.HandleFunc("PUT /api/profile/password", AuthMiddleware(s.secret, s.handleChangePassword))

    // ── Admin endpoints (authenticated + admin role) ──
    mux.HandleFunc("GET /api/admin/waitlist", AuthMiddleware(s.secret, AdminOnly(s.handleWaitlist)))
    mux.HandleFunc("GET /api/admin/users", AuthMiddleware(s.secret, AdminOnly(s.handleListUsers)))
    mux.HandleFunc("POST /api/admin/users/{id}/approve", AuthMiddleware(s.secret, AdminOnly(s.handleApproveUser)))
    mux.HandleFunc("POST /api/admin/users/{id}/reject", AuthMiddleware(s.secret, AdminOnly(s.handleRejectUser)))
    mux.HandleFunc("POST /api/admin/users/{id}/suspend", AuthMiddleware(s.secret, AdminOnly(s.handleSuspendUser)))
    mux.HandleFunc("PUT /api/admin/users/{id}/credentials", AuthMiddleware(s.secret, AdminOnly(s.handleUpdateCredentials)))
    mux.HandleFunc("DELETE /api/admin/users/{id}", AuthMiddleware(s.secret, AdminOnly(s.handleDeleteUser)))

    // ── SPA fallback (MUST be last) ──────────────────
    RegisterSPA(mux, publicFS, SPAOptions{APIPrefix: "/api"})
}
```

### 13.3 Signup handler

```go
func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
    // 1. Decode request
    var req SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, 400, "invalid JSON")
        return
    }

    // 2. Validate fields
    errors := validateSignup(req)
    if len(errors) > 0 {
        respondJSON(w, 400, map[string]interface{}{"errors": errors})
        return
    }

    // 3. Check for existing user (discord_id + email must be unique)
    existing, _ := s.db.GetUserByDiscordID(r.Context(), req.DiscordID)
    if existing != nil {
        respondError(w, 409, "Discord ID already registered")
        return
    }
    existing, _ = s.db.GetUserByEmail(r.Context(), req.Email)
    if existing != nil {
        respondError(w, 409, "Email already registered")
        return
    }

    // 4. Hash password
    hash, err := auth.HashPassword(req.Password)
    if err != nil {
        respondError(w, 500, "internal error")
        return
    }

    // 5. Create user in database (status='waiting')
    user, err := s.db.CreateUser(r.Context(), req.DiscordID, req.Email, req.DisplayName, hash)
    if err != nil {
        respondError(w, 500, "failed to create user")
        return
    }

    // 6. Generate JWT
    token, err := auth.GenerateToken(user.ID, user.Role, s.secret)
    if err != nil {
        respondError(w, 500, "failed to generate token")
        return
    }

    // 7. Return token + user
    respondJSON(w, 201, SignupResponse{Token: token, User: user})
}
```

### 13.4 Frontend auth context (using RTK Query + Redux)

```tsx
// store/authSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface AuthState {
  user: User | null
  token: string | null
}

const initialState: AuthState = {
  user: null,
  token: localStorage.getItem('token'),
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials: (state, action: PayloadAction<{ user: User; token: string }>) => {
      state.user = action.payload.user
      state.token = action.payload.token
      localStorage.setItem('token', action.payload.token)
    },
    logout: (state) => {
      state.user = null
      state.token = null
      localStorage.removeItem('token')
    },
  },
})

export const { setCredentials, logout } = authSlice.actions
export default authSlice.reducer
```

```tsx
// auth/AuthContext.tsx
import { useDispatch, useSelector } from 'react-redux'
import { useLoginMutation, useSignupMutation, useGetMeQuery } from '../store/api'
import { setCredentials, logout, selectAuth } from '../store/authSlice'

export function AuthProvider({ children }: { children: ReactNode }) {
  const dispatch = useDispatch()
  const { token } = useSelector(selectAuth)
  const [loginMutation] = useLoginMutation()
  const [signupMutation] = useSignupMutation()

  // RTK Query auto-fetches and caches /auth/me when token exists
  const { data: me, error: meError } = useGetMeQuery(undefined, {
    skip: !token, // don't fetch if no token
  })

  // If token is invalid, auto-logout
  useEffect(() => {
    if (meError) dispatch(logout())
  }, [meError])

  const login = async (email: string, password: string) => {
    const res = await loginMutation({ email, password }).unwrap()
    dispatch(setCredentials(res))
  }

  const signup = async (data: SignupData) => {
    const res = await signupMutation(data).unwrap()
    dispatch(setCredentials(res))
  }

  const handleLogout = () => {
    dispatch(logout())
  }

  return (
    <AuthContext.Provider value={{ user: me ?? null, login, signup, logout: handleLogout }}>
      {children}
    </AuthContext.Provider>
  )
}
```

### 13.5 React Router setup (with Redux Provider)

```tsx
// main.tsx
import { createRoot } from 'react-dom/client'
import { Provider } from 'react-redux'
import { store } from './store/store'
import App from './App'

createRoot(document.getElementById('root')!).render(
  <Provider store={store}>
    <App />
  </Provider>
)
```

```tsx
// App.tsx
import { BrowserRouter, Routes, Route } from 'react-router-dom'

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Navbar />
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/tutorial" element={<TutorialPage />} />
          <Route path="/waiting-list" element={
            <ProtectedRoute><WaitingListPage /></ProtectedRoute>
          } />
          <Route path="/profile" element={
            <ProtectedRoute><ProfilePage /></ProtectedRoute>
          } />
          <Route path="/profile/password" element={
            <ProtectedRoute><ChangePasswordPage /></ProtectedRoute>
          } />
          <Route path="/admin" element={
            <AdminRoute><AdminDashboard /></AdminRoute>
          } />
          <Route path="/admin/waitlist" element={
            <AdminRoute><AdminWaitlist /></AdminRoute>
          } />
          <Route path="/admin/users/:id" element={
            <AdminRoute><AdminUserDetail /></AdminRoute>
          } />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
        <Footer />
      </AuthProvider>
    </BrowserRouter>
  )
}
```

### 13.6 Storybook story pattern

Every component gets a `.stories.tsx` file **next to it** in the same directory. The pattern is always the same:

```tsx
// components/CredentialCard.stories.tsx
import type { Meta, StoryObj } from '@storybook/react-vite'
import { CredentialCard } from './CredentialCard'

const meta: Meta<typeof CredentialCard> = {
  title: 'Components/CredentialCard',
  component: CredentialCard,
  tags: ['autodocs'],  // generates docs tab automatically
  argTypes: {
    maskable: { control: 'boolean' },
  },
}
export default meta
type Story = StoryObj<typeof CredentialCard>

// Default state: a non-sensitive value
export const Normal: Story = {
  args: {
    label: 'Application ID',
    value: '987654321098765432',
  },
}

// Masked state: a sensitive value hidden behind dots
export const Masked: Story = {
  args: {
    label: 'Bot Token',
    value: 'MTIz...ODk.GHx...',
    maskable: true,
  },
}

// Revealed state: user has clicked "Show"
export const Revealed: Story = {
  args: {
    label: 'Bot Token',
    value: 'MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNA',
    maskable: true,
  },
}
```

**Storybook conventions for this project:**

| Convention | Why |
|---|---|
| One `.stories.tsx` per component | Easy to find, easy to review in PRs |
| `tags: ['autodocs']` on every story | Auto-generates documentation tab |
| Name stories by visual state (`Empty`, `WithError`, `Loading`) | Makes it easy to visually verify all states |
| Use `fn()` from `@storybook/test` for callbacks | Avoids Storybook warnings, shows callback in Actions tab |
| Pass mock data as `args` | Stories work as living documentation |

**Running Storybook:**
```bash
make storybook     # or: cd ui && pnpm storybook
# Opens http://localhost:6006
```

## 14. Testing Strategy

### 14.1 Backend testing (Go)

Every database function has a table-driven test using an in-memory SQLite database:

```go
func TestCreateUser(t *testing.T) {
    db := openTestDB(t) // in-memory SQLite
    defer db.Close()

    tests := []struct {
        name        string
        discordID   string
        email       string
        expectError bool
    }{
        {"valid user", "123", "test@test.com", false},
        {"duplicate discord_id", "123", "other@test.com", true},
        {"duplicate email", "456", "test@test.com", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := db.CreateUser(context.Background(), tt.discordID, tt.email, "Test", "hash")
            if tt.expectError && err == nil { t.Error("expected error") }
            if !tt.expectError && err != nil { t.Error(err) }
        })
    }
}
```

HTTP handler tests use `httptest.NewRecorder`:

```go
func TestHandleSignup(t *testing.T) {
    srv := newTestServer(t)
    body := `{"discord_id":"123","email":"t@t.com","display_name":"T","password":"password123"}`
    req := httptest.NewRequest("POST", "/api/auth/signup", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    srv.handleSignup(w, req)
    if w.Code != 201 { t.Errorf("expected 201, got %d", w.Code) }
}
```

### 14.2 Frontend testing (Storybook + interaction tests)

Storybook serves as our primary component test environment. For each component:

1. **Visual stories** — verify the component looks correct in each state
2. **Interaction stories** (using `@storybook/test`) — verify click/type behavior:

```tsx
import { expect, fn, userEvent, within } from '@storybook/test'

export const SubmitForm: Story = {
  args: { onSubmit: fn() },
  play: async ({ args, canvas }) => {
    const canvas = within(canvas)
    await userEvent.type(canvas.getByLabelText('Email'), 'test@test.com')
    await userEvent.type(canvas.getByLabelText('Password'), 'password123')
    await userEvent.click(canvas.getByRole('button', { name: /log in/i }))
    await expect(args.onSubmit).toHaveBeenCalledWith({
      email: 'test@test.com',
      password: 'password123',
    })
  },
}
```

3. **Build-time smoke test** — CI runs `pnpm build-storybook` to verify no stories are broken

### 14.3 End-to-end testing

Manual end-to-end walkthrough script (run before every release):

1. Start backend: `make dev-backend`
2. Start frontend: `make dev-frontend`
3. Open `http://localhost:5173`
4. Sign up as a new user
5. Verify waiting-list page shows correct status
6. Open a new browser, log in as admin
7. Verify admin dashboard shows the new user
8. Approve the user with test credentials
9. Switch back to user browser, refresh profile page
10. Verify credentials appear with mask/copy
11. Visit tutorial page, verify content renders

## 15. Risks, Alternatives, and Open Questions

### 15.1 Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Bot tokens leaked through XSS | Users' Discord bots compromised | Content Security Policy, sanitize all user input, never render tokens in URLs |
| SQLite database corruption | All data lost | WAL mode, regular backups (cron job copies the `.db` file) |
| Admin accidentally deletes wrong user | User data lost irrecoverably | Soft-delete option (add `deleted_at` column), confirmation dialogs |
| No rate limiting on signup | Spam signups overwhelm the waiting list | Add rate limiting middleware (5 signups per IP per hour) |
| JWT secret compromised | Anyone can forge tokens | Load from env var, rotate periodically, short expiry |

### 15.2 Alternatives considered

| Decision | Alternative | Why we didn't choose it |
|---|---|---|
| SQLite | PostgreSQL | Overkill for this scale; requires separate server |
| JWT | Server sessions | More complex for a single-binary app |
| RTK Query | React Query (TanStack Query) | Both are excellent; RTK Query integrates with Redux, which we use for auth state |
| Storybook | No component isolation | Storybook catches visual bugs early and serves as living documentation |
| Manual Discord ID entry | Discord OAuth | OAuth requires a registered Discord application; manual entry is simpler for V1 |
| Tailwind CSS | CSS Modules / styled-components | Tailwind is faster for prototyping and consistent across the team |

### 15.3 Open questions

1. **Email notifications**: Should the system send an email when a user is approved? (Requires SMTP setup or a transactional email service like Resend/SendGrid.)
2. **Discord OAuth**: Should V2 add a "Sign in with Discord" button that auto-fills the Discord ID?
3. **Bot health monitoring**: Should the platform check if a user's bot is actually running and show status?
4. **Multi-bot support**: Should users be able to run multiple bots (multiple credential sets)?
5. **Waitlist ordering**: Should the admin see users in strict FIFO order, or should they be able to prioritize?

## 16. References

### External resources

- [go-go-golems/discord-bot](https://github.com/go-go-golems/discord-bot) — The Discord bot runtime this platform is built around
- [discord-bot tutorial](https://github.com/go-go-golems/discord-bot/blob/main/pkg/doc/tutorials/building-and-running-discord-js-bots.md) — Full bot authoring tutorial
- [Go 1.22+ ServeMux](https://pkg.go.dev/net/http#ServeMux) — New handler syntax with `{...}` pattern matching
- [RTK Query documentation](https://redux-toolkit.js.org/rtk-query/overview) — Redux Toolkit Query overview and API
- [Storybook for React](https://storybook.js.org/docs/get-started/frameworks/react) — Component development environment
- [Tailwind CSS documentation](https://tailwindcss.com/docs) — Utility-first CSS framework
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — Pure-Go SQLite driver (no CGO)
- [bcrypt via x/crypto](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — Password hashing
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) — JWT library for Go

### Key files from the discord-bot repo

| File | Why it matters |
|---|---|
| `README.md` | Installation, quick start, architecture |
| `pkg/doc/tutorials/building-and-running-discord-js-bots.md` | Full tutorial content for our /tutorial page |
| `examples/discord-bots/ping/index.js` | Richest example bot — used in tutorial |
| `examples/discord-bots/knowledge-base/index.js` | Runtime config example |
| `internal/jsdiscord/host.go` | How the JS runtime works (background knowledge) |
