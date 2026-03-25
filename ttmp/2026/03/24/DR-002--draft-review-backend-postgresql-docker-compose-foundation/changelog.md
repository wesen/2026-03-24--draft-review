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
