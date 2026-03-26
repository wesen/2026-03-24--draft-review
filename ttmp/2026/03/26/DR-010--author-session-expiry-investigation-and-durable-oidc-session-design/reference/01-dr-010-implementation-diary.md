---
Title: DR-010 implementation diary
Ticket: DR-010
Status: active
Topics:
    - auth
    - backend
    - oidc
    - keycloak
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: pkg/auth/oidc.go
    - Path: pkg/auth/postgres.go
    - Path: pkg/auth/session.go
      Note: Touches author_sessions on each authenticated read
    - Path: pkg/auth/session_test.go
      Note: Verifies session reads update activity timestamps
    - Path: pkg/db/migrations/0006_author_session_activity.sql
      Note: Introduces last_used_at activity metadata for opaque author sessions
    - Path: pkg/server/http.go
ExternalSources: []
Summary: Detailed step-by-step diary for the server-side session implementation work.
LastUpdated: 2026-03-26T15:50:00-04:00
WhatFor: Keep a running implementation diary for the medium-term opaque session migration.
WhenToUse: Use while reviewing the DR-010 implementation history and validation steps.
---


# DR-010 implementation diary

## Goal

Implement the medium-term auth repair: move Draft Review from a stateless signed
browser cookie carrying identity claims to an opaque server-side author session backed
by the existing `author_sessions` table.

## Context

The repo already contains the essential building blocks for this migration:

- OIDC login and logout flow
- local `users` rows keyed by OIDC identity
- an `author_sessions` table in the initial schema

The old implementation still stored the full identity payload in the cookie and tied
session lifetime directly to provider token expiry. DR-010 is the migration away from
that model.

## Quick Reference

### Execution strategy

1. Create the future-state design ticket separately so the medium-term implementation
   does not accumulate long-term refresh-token complexity.
2. Expand DR-010 tasks into concrete implementation slices.
3. Rebuild session management around opaque tokens and Postgres session lookup.
4. Update callback/login flow to create local user + session immediately.
5. Update request auth resolution and logout.
6. Add tests, docs, and validation notes.

## Diary

### 2026-03-26 15:50 EDT

Started the medium-term implementation planning pass.

What I confirmed before editing code:

- the current session cookie logic in `pkg/auth/session.go` is still a signed payload
  cookie containing identity claims
- `pkg/auth/oidc.go` still copies `token.Expiry` into `SessionClaims.ExpiresAt`
- the initial schema already contains `author_sessions`, which gives the migration a
  clean persistence anchor

Why this matters:

- I do not need to invent a new session table for the medium-term fix
- the main code change is architectural rewiring, not schema invention
- the long-term refresh-token design should stay in a separate ticket so DR-010 stays
  focused

### 2026-03-26 16:05 EDT

Implemented the first real code slice for opaque server-side author sessions.

Files changed in this slice:

- `pkg/auth/config.go`
- `pkg/auth/types.go`
- `pkg/auth/postgres.go`
- `pkg/auth/session.go`
- `pkg/auth/oidc.go`
- `pkg/server/http.go`
- `pkg/auth/session_test.go`
- `pkg/server/http_test.go`

What changed:

- added `auth-session-ttl` parsing so the backend has an explicit app-managed session
  duration
- implemented repository methods for `author_sessions` using the already existing
  table in Postgres
- replaced the old signed-claims cookie model with an opaque token cookie whose hash
  is looked up in the session table
- changed OIDC callback handling so it now ensures the local user exists immediately,
  then creates a server-side author session
- changed request auth resolution in the HTTP handler so `/api/me` and author routes
  read identity from the session store instead of from cookie payload claims
- changed logout to revoke the server-side session and clear the browser cookie

Important design note from this slice:

- I kept `SessionClaims` as the request-time identity shape because it lets the rest of
  the handler code stay mostly stable while the storage model underneath changes

Validation run:

```bash
go test ./cmd/... ./pkg/...
```

Result:

- all tests passed

What was intentionally deferred:

- long-term refresh-token storage and renewal remain in DR-011
- richer session metadata like `last_used_at` is not part of this medium-term slice

### 2026-03-26 16:20 EDT

Aligned the documentation with the code that landed in the first implementation
commit.

Files updated in this follow-up slice:

- `README.md`
- `docs/deployments/draft-review-coolify.md`
- `ttmp/.../DR-010/design-doc/01-author-session-expiry-investigation-and-durable-oidc-session-design-guide.md`

Why this mattered:

- the original DR-010 guide still described the pre-implementation recommendation too
  literally and needed to reflect the new opaque-session model
- operators need to know about `DRAFT_REVIEW_AUTH_SESSION_TTL`
- local developers need the updated `serve --auth-session-ttl 12h` example

### 2026-03-26 18:55 EDT

Started the first hardening follow-up after the initial opaque-session deploy.
I deliberately kept this as a narrow activity-tracking slice instead of mixing it
with sliding renewal, because both features touch `pkg/auth/session.go` and would
be harder to review if they landed together.

Files changed in this slice:

- `pkg/db/migrations/0006_author_session_activity.sql`
- `pkg/auth/types.go`
- `pkg/auth/postgres.go`
- `pkg/auth/session.go`
- `pkg/auth/session_test.go`
- `pkg/server/http_test.go`

What changed:

- added a schema migration that introduces `author_sessions.last_used_at`, backfills
  it from `created_at`, and indexes `(user_id, last_used_at)` for future session
  inspection and idle-timeout work
- extended the auth repository to scan and update `last_used_at` whenever an opaque
  author session is loaded and accepted
- added a `TouchAuthorSession(...)` repository method and wired the session manager
  to call it on every authenticated session read
- updated unit tests and handler tests so fake session stores now implement the
  new touch method and assert that a successful session read updates activity state

Why I split it this way:

- `last_used_at` is useful on its own for observability and future idle policies
- sliding renewal needs response-writer-aware cookie refresh behavior, which is a
  separate behavioral change and should get its own commit

Validation run:

```bash
gofmt -w pkg/auth/config.go pkg/auth/postgres.go pkg/auth/session.go pkg/auth/session_test.go pkg/auth/types.go pkg/server/http_test.go
go test ./pkg/auth ./pkg/server
```

Result:

- both package test suites passed after the touch-on-read wiring was added

What was tricky:

- I had already started changing the session manager API for sliding renewal before
  this step. I backed those renewal-specific pieces out temporarily so the first
  follow-up commit would stay focused on activity tracking and not drag response
  cookie semantics into the same diff.

What to watch in review:

- `TouchAuthorSession(...)` now happens on every authenticated request, so reviewers
  should confirm this write pattern is acceptable for the expected author traffic
- the migration backfills existing rows and makes `last_used_at` non-null, so schema
  reviewers should confirm that aligns with any external DB inspection scripts

### Next diary entries

Add entries after each implementation slice, including:

- which files changed
- what command or test was run
- what failed
- what was learned

## Usage Examples

Use this diary together with:

- [DR-010 design guide](../design-doc/01-author-session-expiry-investigation-and-durable-oidc-session-design-guide.md)
- [tasks.md](../tasks.md)
- [changelog.md](../changelog.md)

## Related

- [DR-011 long-term guide](/home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/26/DR-011--long-term-durable-auth-sessions-with-server-side-refresh-token-renewal/design-doc/01-long-term-durable-auth-sessions-with-server-side-refresh-token-renewal-guide.md)
