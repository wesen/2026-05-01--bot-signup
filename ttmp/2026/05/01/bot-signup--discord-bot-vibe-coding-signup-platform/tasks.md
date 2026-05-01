# Tasks

## TODO

- [x] Phase 1: Project scaffolding (Go server + health endpoint)
- [x] Phase 2: Database layer (SQLite + migrations + CRUD tests)
- [x] Phase 3R: Replace password/JWT auth with Discord OAuth + HTTP-only sessions
- [x] Phase 4R: Adjust profile/admin handlers to use session middleware and OAuth-created users
- [x] Phase 5: Frontend scaffolding (Vite + React + Tailwind + Storybook + RTK Query)
- [x] Phase 6: OAuth landing/signup card (VibeBot Sessions visual reference + Storybook stories)
- [x] Phase 7: User pages (profile + waiting list + credential display + stories)
- [ ] Phase 8: Admin pages (dashboard + approval form + stories)
- [ ] Phase 9: Tutorial page (embedded markdown from discord-bot)
- [ ] Phase 10: Frontend embedding (go:embed + single-binary build)
- [ ] Phase 11: Polish and deploy (error handling, CI, Storybook build)

## Done

- [x] Create ticket and initial workspace
- [x] Write comprehensive design and implementation guide (16 sections)
- [x] Upload design doc to reMarkable
- [x] Store VibeBot Sessions UI reference image in ticket sources
- [x] Update implementation guide to make Discord OAuth the only auth path (no passwords)

## Superseded

- [x] Phase 3 old implementation: bcrypt + JWT + signup/login endpoints (committed, now must be replaced by Phase 3R)
- [x] Phase 4 old implementation: profile/admin handlers using JWT middleware (committed, now must be adapted by Phase 4R)
