# Changelog

## 2026-05-01

- Initial workspace created


## 2026-05-01

Created comprehensive 16-section design document covering full system architecture, database schema, API reference, frontend pages, RTK Query state management, Storybook component stories, authentication, admin backend, tutorial content, project structure, 11-phase implementation plan, pseudocode, testing strategy, risks, and references.


## 2026-05-01

Implemented Phase 1-3 backend foundation: Go server scaffold, SQLite database layer with migrations and CRUD tests, bcrypt/JWT auth, signup/login/logout/me endpoints, auth middleware, and diary/task updates.


## 2026-05-01

Implemented Phase 4 profile/admin backend: profile endpoints, password change, stats, admin waitlist/users/approve/reject/suspend/update-credentials/delete routes, transactional approval, and route tests.


## 2026-05-01

Stored VibeBot Sessions UI reference image in ticket sources and updated the implementation guide: Discord OAuth is now the only auth path, password/JWT/localStorage flows are removed from the design, session cookies are recommended, and tasks now include Phase 3R/4R refactors.


## 2026-05-01

Implemented Phase 3R/4R: replaced password/JWT auth with Discord OAuth, signed HTTP-only session cookies, OAuth state cookies, Discord-user upsert schema/model changes, session middleware, and updated profile/admin tests to use session auth.

