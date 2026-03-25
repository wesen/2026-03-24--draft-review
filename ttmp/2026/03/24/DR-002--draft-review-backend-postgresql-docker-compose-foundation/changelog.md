# Changelog

## 2026-03-24

- Initial workspace created


## 2026-03-24

Created the backend planning ticket, wrote the intern-facing architecture guide, recorded the diary, related source files, and validated the ticket.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/imports/draft-review-screen-spec.md — Source specification used to derive backend requirements
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md — Primary backend design artifact for the ticket
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md — Chronological diary for the ticket work


## 2026-03-24

Uploaded the final ticket bundle to reMarkable and verified the remote listing at /ai/2026/03/24/DR-002.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md — Included in the uploaded bundle
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md — Included in the uploaded bundle
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md — Included in the uploaded bundle


## 2026-03-24

Rewrote the backend guide to target Go, Glazed command verbs, Clay-style SQL configuration, pgxpool, and hair-booking-style server wiring instead of the earlier TypeScript stack recommendation.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md — Primary backend guide rewritten around the Go and Glazed direction
- /home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go — Source pattern for the Glazed DB section design
- /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go — Source pattern for the Go server verb architecture


## 2026-03-24

Revalidated the rewritten Go/Glazed ticket and force-overwrote the existing reMarkable bundle so the remote copy matches the revised guide.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md — Revised guide delivered to reMarkable
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md — Diary updated with the rewrite validation and overwrite details


## 2026-03-24

Implemented the Go backend foundation with Glazed verbs, Clay-style SQL config, embedded migrations, seed support, HTTP server bootstrap, Docker assets, and commit 01c8c9fa2d4b1bbcdb9f1bdbed51bb5a956d0bf9.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/cmd/draft-review/cmds/serve.go — Serve verb with config
- /home/manuel/code/wesen/2026-03-24--draft-review/cmd/draft-review/main.go — Root Cobra plus Glazed command wiring
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/config/sql.go — Clay-style SQL connection section and DSN builder
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations.go — Embedded SQL migration runner
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go — Initial health and info routes


## 2026-03-24

Added the first PostgreSQL-backed article read endpoints with commit 845eeb67ff0d75ffa5586037cf82a900f05469ab.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go — Postgres article queries and section serialization
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/service.go — Article service wrapper
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go — Article list/detail routes and placeholder readers/reactions routes


## 2026-03-24

Expanded the execution checklist and added a repository-level backend runbook covering Docker Compose, local CLI flows, current API routes, and temporary implementation limits.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/README.md — Backend quickstart and current API notes
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md — More detailed execution queue
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md — Diary entry for the runbook/checklist slice


## 2026-03-24

Added initial PostgreSQL-backed article write endpoints, including `POST /api/articles`, `PATCH /api/articles/{id}`, request validation, transactional section replacement, and a live Docker/Postgres smoke test.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/service.go — Article write validation and normalization
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go — Transactional create/update persistence and section replacement
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go — HTTP handlers for article create/update plus JSON decoding
- /home/manuel/code/wesen/2026-03-24--draft-review/README.md — Updated API surface and current limitation notes


## 2026-03-24

Added the initial auth package with domain types, password hashing helpers, opaque token/session-cookie helpers, and unit tests.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/types.go — Auth-facing DTO and persistence model shapes
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/helpers.go — Password, token, and session cookie helpers
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/helpers_test.go — Focused tests for auth helper behavior


## 2026-03-24

Pivoted Draft Review auth to the Keycloak / OIDC model used by `hair-booking`, adding Glazed auth settings, signed browser sessions, `/auth/*` routes, `/api/me`, auth-aware `/api/info`, and focused auth/server tests.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/config.go — Glazed auth section and OIDC settings
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go — Signed browser session cookie manager
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go — Keycloak / OIDC login and callback flow
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go — Auth route registration and `/api/me`
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http_test.go — `/api/me` server coverage


## 2026-03-24

Bound authenticated authors to local `users` rows, added auth-identity columns plus migration/backfill, scoped author article routes by `owner_user_id`, and validated the flow with a live migrate/seed/serve smoke test.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0003_user_auth_identity.sql — Adds local auth identity mapping to users
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/service.go — Ensures a local author row from auth claims
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/postgres.go — Local user bootstrap and update persistence
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go — Owner-scoped article queries and writes
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/seed.go — Development author upsert with auth identity fields


## 2026-03-24

Added a standalone local Keycloak development stack, helper `make` targets, README/runbook coverage for port overrides, and a live browser-verified OIDC smoke test against the imported `draft-review-dev` realm.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/docker-compose.local.yml — Standalone local Postgres plus Keycloak stack for OIDC testing
- /home/manuel/code/wesen/2026-03-24--draft-review/dev/keycloak/realm-import/draft-review-dev-realm.json — Imported Keycloak realm, client, and test user
- /home/manuel/code/wesen/2026-03-24--draft-review/Makefile — Helper targets for bringing the local Keycloak stack up and running the backend in dev or OIDC mode
- /home/manuel/code/wesen/2026-03-24--draft-review/README.md — Local Keycloak bootstrap, port override, and smoke-test instructions
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md — Diary entry recording the failures, fixes, and final successful smoke flow


## 2026-03-24

Implemented the reader-sharing, review-session, analytics, and versioning slices: real share-token and invite persistence, token-based reader article resolution, review progress and reaction storage, summary submission, author analytics/readers/reactions queries, a cross-article reader directory, paragraph-progress persistence, and an authenticated export stub.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0004_review_paragraph_progress.sql — Adds paragraph-level progress persistence for review sessions
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go — Share-token reset, invite creation, and reader-link resolution
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/postgres.go — Review session creation, progress tracking, reactions, and summary persistence
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/analytics/postgres.go — Real readers, reactions, analytics, and reader-directory queries
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go — Article version cloning and activation
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go — Route wiring for sharing, reader sessions, analytics, and export stub
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http_test.go — Contract-focused HTTP coverage for the new routes
- /home/manuel/code/wesen/2026-03-24--draft-review/README.md — Updated route inventory and known-gap notes


## 2026-03-24

Aligned the frontend with the real Go backend by adding auth and reader-session RTK Query clients, making MSW opt-in instead of the default dev path, gating the author app on `/api/me`, wiring invite sending, and persisting reader start/progress/reactions/summary through the new backend endpoints.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/authApi.ts — `/api/me` client contract for browser auth state
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/readerApi.ts — RTK Query client for reader session start/progress/reactions/summary
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/articleApi.ts — Added version/export alignment with the real backend routes
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/app/AuthorApp.tsx — Auth-aware author boot and invite wiring
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/app/ReaderApp.tsx — Passes the review token into the real reader-session flow
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/ReaderPage.tsx — Starts review sessions and persists progress/reactions/summary
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/main.tsx — MSW is now opt-in via `VITE_USE_MSW=1`
