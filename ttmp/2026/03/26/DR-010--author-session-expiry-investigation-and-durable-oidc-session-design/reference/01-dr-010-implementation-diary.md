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
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/postgres.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go
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
