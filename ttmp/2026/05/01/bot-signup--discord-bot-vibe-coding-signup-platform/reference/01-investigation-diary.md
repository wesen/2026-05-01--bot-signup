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
