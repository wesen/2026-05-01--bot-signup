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


## 2026-05-01

Implemented Phase 5/6 frontend scaffold: Vite React TypeScript app, Tailwind v4, Redux Toolkit/RTK Query, Discord OAuth auth provider, VibeBot Sessions landing page matching the stored reference image, Storybook config, component stories, Makefile frontend targets, and passing lint/build/storybook checks.


## 2026-05-01

Implemented Phase 7 user pages: protected route, waiting-list page, profile dashboard with credential cards, auth callback placeholder, tutorial placeholder, RTK Query getProfile endpoint, component/page Storybook stories, and validation with UI lint/build/storybook plus Go tests.


## 2026-05-01

Implemented Phase 8 admin pages: admin route guard, stats cards, waitlist table, approval form, admin dashboard/detail routes, RTK Query admin endpoints, admin Storybook stories, and validation with frontend lint/build/storybook plus Go tests.


## 2026-05-01

Implemented Phase 9 tutorial page: copied discord-bot tutorial markdown into ui/src/content, rendered it with react-markdown/remark-gfm, added markdown styling and TutorialPage Storybook story, and validated lint/build/storybook/tests.


## 2026-05-01

Implemented Phase 10/11 build and delivery polish: Dagger-backed cmd/build-web with pnpm cache, internal/web go:embed SPA serving, Makefile build-web/build targets, GitHub Actions CI, README, packageManager pin, and validated make build with Dagger exporting assets before embedded Go build.

