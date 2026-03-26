---
Title: Author session expiry investigation and durable OIDC session design guide
Ticket: DR-010
Status: active
Topics:
    - auth
    - backend
    - oidc
    - keycloak
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../terraform/keycloak/apps/draft-review/envs/hosted/main.tf
    - Path: ../../../../../../../terraform/keycloak/modules/realm-base/main.tf
    - Path: README.md
      Note: Documents the new renewal settings and debug endpoint for local operators
    - Path: docs/deployments/draft-review-coolify.md
      Note: Documents hosted env vars and validation for session renewal and inspection
    - Path: pkg/auth/config.go
    - Path: pkg/auth/oidc.go
    - Path: pkg/auth/session.go
    - Path: pkg/server/http.go
      Note: Shows the dedicated opaque-session debug route and current session state lookup
ExternalSources: []
Summary: Detailed investigation and implementation guide for the fast logout issue in Draft Review's Keycloak-backed browser auth flow.
LastUpdated: 2026-03-26T13:35:00-04:00
WhatFor: Explain why Draft Review logs authors out quickly today and provide a detailed implementation plan for a more durable browser session model.
WhenToUse: Use when debugging short author sessions, changing Keycloak realm settings, or redesigning app-side auth persistence.
---


# Author session expiry investigation and durable OIDC session design guide

## Executive Summary

Draft Review currently appears to log authors out too quickly because the backend
treats the provider token lifetime as the same thing as the application session
lifetime. During OIDC callback handling, the app verifies the Keycloak login
correctly, but then it copies the upstream token expiry into the local Draft Review
session cookie. Once that copied timestamp passes, the browser cookie is considered
expired and author API routes become unauthenticated.

This is the key architectural issue an intern needs to understand: an identity
provider session and an application session are related, but they are not the same
object. A provider proves who the user is. The application decides how long it wants
to trust the user's local browser session after successful login. Right now Draft
Review collapses those two concerns together.

The recommended medium-term fix is to keep Keycloak as the identity provider but move
Draft Review to an opaque server-side session model. In that model, the browser only
holds a random session token, the backend stores the real session record in the
existing `author_sessions` table, and the app uses its own session TTL instead of the
raw OIDC token expiry.

## Problem Statement

Authors report that the hosted application logs them out after a short amount of
time. This creates a poor writing experience because:

- a writer can be interrupted mid-session
- a tab left open during ordinary work becomes invalid too soon
- the app feels unreliable even when login technically succeeded

The code strongly suggests one primary cause and one secondary operational risk.

Primary cause:

- the app copies Keycloak token expiry into the Draft Review browser session

Secondary risk:

- if the hosted session signing secret changes across deploys or restarts, all
  existing cookies become invalid immediately

The first issue is definitely present in code. The second issue must be checked in the
hosted environment.

### Files the intern must read first

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)
- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)
- [config.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/config.go)
- [http.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go)
- [main.tf](/home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted/main.tf)
- [main.tf](/home/manuel/code/wesen/terraform/keycloak/modules/realm-base/main.tf)
- [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md)

### Why this is a product problem, not just an auth detail

Draft Review is an authoring tool. Authoring tools need durable, predictable user
sessions because the user is actively creating content, stepping away, coming back,
and continuing work. A short-lived session model is much more damaging here than in a
rarely used admin screen.

## Current System Walkthrough

This section explains the existing system in the order a request flows.

### 1. Keycloak realm and browser client

The hosted Keycloak configuration for Draft Review is managed in:

- [main.tf](/home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted/main.tf)

This Terraform creates:

- the `draft-review` realm
- a confidential browser client `draft-review-web`
- callback and logout redirect URIs for the hosted app

What the current Terraform does not configure:

- access token lifespan
- SSO session idle or max lifespan
- client-specific token lifetime overrides

That matters because the app is likely inheriting Keycloak defaults.

### 2. OIDC callback

The login callback is handled in:

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)

Conceptually, the old code path did this:

```text
receive auth code
exchange auth code for tokens
verify id_token
extract claims
read token.Expiry
write Draft Review session claims with ExpiresAt = token.Expiry
```

That final step is the important one. The callback path currently assumes:

```text
provider token expiry == application session expiry
```

That assumption is the likely source of the quick logout behavior.

### 3. Session manager

Cookie writing and reading are handled in:

- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)

The current medium-term implementation now does this instead:

- `WriteSession` generates an opaque random token
- it stores only a hash of that token in `author_sessions`
- it writes the raw token into the browser cookie
- `ReadSession` hashes the incoming cookie, loads the session row, joins the linked
  user, and reconstructs request-time identity claims from Postgres

That is the core architectural change in DR-010. The browser no longer carries the
full identity payload.

### 4. API routes that expose the symptom

The main routes affected are wired in:

- [http.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go)

Relevant browser/API endpoints:

- `GET /auth/login`
- `GET /auth/callback`
- `GET /auth/logout`
- `GET /auth/logout/callback`
- `GET /api/me`
- author routes such as `GET /api/articles` and `POST /api/articles`

Once the app cookie expires, `/api/me` reports the user as unauthenticated and the
author UI stops behaving like an active editing session.

### 5. Runtime configuration

The hosted runtime contract is documented in:

- [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md)

Important environment variables:

- `DRAFT_REVIEW_AUTH_MODE`
- `DRAFT_REVIEW_AUTH_SESSION_SECRET`
- `DRAFT_REVIEW_OIDC_ISSUER_URL`
- `DRAFT_REVIEW_OIDC_CLIENT_ID`
- `DRAFT_REVIEW_OIDC_CLIENT_SECRET`
- `DRAFT_REVIEW_OIDC_REDIRECT_URL`

`DRAFT_REVIEW_AUTH_SESSION_SECRET` is still part of the session-token hashing path in
the medium-term implementation. If it changes, existing server-side session cookies
will no longer match their stored hashes. That would still cause logout across a
deploy or restart.

## Current and Recommended Flows

### Current flow

```text
browser -> /auth/login
        -> Keycloak login page
        -> /auth/callback
        -> verify token
        -> copy token expiry into app session
        -> write cookie with same expiry
```

Diagram:

```text
+---------+        +----------------+        +-----------+
| Browser | -----> | Draft Review   | -----> | Keycloak  |
|         |        | /auth/login    |        |           |
+---------+        +----------------+        +-----------+
     ^                                            |
     |                                            v
     |        +----------------------------------------+
     +--------| /auth/callback copies token expiry     |
              | into app session cookie                |
              +----------------------------------------+
```

### Recommended medium-term flow

```text
browser -> /auth/login
        -> Keycloak login page
        -> /auth/callback
        -> verify token
        -> ensure local user exists
        -> create Draft Review server-side session using app TTL
        -> write opaque cookie token
```

Diagram:

```text
Keycloak login success
        |
        v
Draft Review creates app session row
        |
        +--> browser receives opaque token
        +--> backend keeps real session state in Postgres
```

The important change is that Draft Review becomes responsible for its own browser
session duration and session lookup, while still delegating login to Keycloak.

## Proposed Solution

The recommended medium-term design is an opaque server-side session model with an
application-defined TTL.

### Proposed behavior

- keep Keycloak for login and logout
- verify OIDC tokens exactly as today
- stop copying `token.Expiry` into the cookie payload because there is no cookie
  payload anymore
- create or update the local user during callback
- create a row in `author_sessions`
- use an app session TTL, for example `12h`
- record `author_sessions.last_used_at` on authenticated requests
- optionally renew active sessions before expiry when they are close to timing out
- expose a current-session inspection endpoint for debugging
- verify the hosted session secret is stable across restarts

### Why this solves the user-facing issue

If Keycloak currently issues short-lived access tokens, the old code made the app
feel like it had short-lived logins. Once the app creates its own server-side session
row with an app-defined TTL, a successful login results in a stable author session
suitable for an editing workflow.

### Minimal viable settings

Recommended first settings:

- `auth-session-ttl = 12h`
- `auth-session-sliding-renewal = true`
- `auth-session-renew-before = 1h`

That is a pragmatic middle ground:

- long enough for normal authoring work
- short enough to avoid effectively permanent login

## Design Decisions

### Decision 1: Separate provider token lifetime from app session lifetime

Rationale:

- provider tokens are identity infrastructure
- app sessions are product behavior
- coupling them tightly creates bad UX and hides the real policy boundary

### Decision 2: Prefer app-defined TTL before implementing refresh-token renewal

Rationale:

- smaller first slice
- easier for an intern to implement safely
- fixes the likely symptom without adding storage and rotation complexity

### Decision 3: Keep explicit logout flowing through Keycloak

Rationale:

- local app cookie still needs to be cleared
- provider logout still matters for the broader browser session
- this keeps sign-out behavior coherent

### Decision 4: Add a dedicated session inspection surface instead of overloading `/api/me`

Rationale:

- `/api/me` should stay product-facing and stable for the frontend
- session debugging needs more detail than the normal user-info payload should expose
- a dedicated route can show renewal state, `last_used_at`, and local session metadata

## Alternatives Considered

### Alternative A: Only increase Keycloak token/session lifespans

Pros:

- little or no app code change

Cons:

- still leaves the conceptual design flaw in place
- changes provider policy to fix app UX
- keeps future debugging muddy

Verdict:

- not recommended as the primary fix

### Alternative B: Full refresh-token support immediately

Pros:

- closer to a conventional OIDC renewal design
- preserves short-lived access tokens while extending usability

Cons:

- more code and more operational complexity
- refresh token storage and rotation must be handled carefully
- larger test matrix

Verdict:

- valid later, too heavy as the first repair slice

### Alternative C: Server-side session store in Postgres or Redis

Pros:

- easier centralized invalidation
- easier inspection of active sessions

Cons:

- extra stateful infrastructure and cleanup logic
- not required for the current scale of the app

Verdict:

- unnecessary for the present bug unless revocation needs become stronger later

## Detailed Implementation Plan

This section is written as a task sequence for an intern.

### Step 1: Confirm the effective current expiry

After logging in, call:

- `GET /api/me`

Record:

- current timestamp
- `sessionExpiresAt`

If the expiration is only a few minutes in the future, the current theory is
confirmed.

Pseudocode:

```text
login()
me = GET /api/me
print(now, me.sessionExpiresAt)
assert me.sessionExpiresAt - now is suspiciously short
```

### Step 2: Add new auth settings

Extend:

- [config.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/config.go)

Recommended settings:

- `auth-session-ttl`
- `auth-session-sliding-renewal`
- `auth-session-renew-before`

Recommended environment shape:

- `DRAFT_REVIEW_AUTH_SESSION_TTL=12h`
- `DRAFT_REVIEW_AUTH_SESSION_SLIDING_RENEWAL=true`
- `DRAFT_REVIEW_AUTH_SESSION_RENEW_BEFORE=1h`

The intern should parse these once and choose defaults that are easy to reason about.

### Step 3: Change OIDC callback session creation and create a server-side session

Update:

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)

Old conceptual code:

```go
expiry := token.Expiry
sessionClaims.ExpiresAt = expiry.UTC()
```

Medium-term conceptual code:

```go
user := authService.EnsureAuthenticatedUser(identity)
sessionManager.WriteSession(ctx, w, r, user)
```

Important note:

Do not remove token verification. The app must still verify the Keycloak callback.
Only the local session creation model should change.

### Step 4: Read identity from the session table on each request

The handler path should:

Pseudocode:

```text
rawToken = cookie value
tokenHash = sessionManager.hash(rawToken)
resolved = repository.findSessionByTokenHash(tokenHash)
claims = buildSessionClaims(resolved.user, resolved.session)
```

This gives the backend full control over revocation and expiry checks.

### Step 5: Keep logout intact

Retain the existing logout shape in:

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)
- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)

Requirements:

- local cookie is cleared
- backing session row is revoked
- provider logout redirect still succeeds
- post-logout callback remains valid

### Step 6: Audit hosted secret stability

In the hosted Coolify environment, confirm that:

- `DRAFT_REVIEW_AUTH_SESSION_SECRET` is persistent
- deploys do not rotate it accidentally

If this value changes, users will still be logged out on deploy even after the TTL
fix.

### Step 7: Add tests

Relevant files:

- [session_test.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session_test.go)
- [oidc_test.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc_test.go)

Test cases to add:

- server-side session creation writes the expected cookie and database record
- unknown or revoked session tokens are rejected correctly
- logout still clears the cookie
- `/api/me` resolves identity from the server-side session lookup
- reading a valid session updates `last_used_at`
- active sessions close to expiry renew their cookie and DB expiry
- `GET /api/debug/session` returns local session metadata for the current browser session

### Step 8: Update deployment documentation

Update:

- [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md)

Document:

- the new app TTL setting
- the sliding renewal settings
- the difference between app session TTL and Keycloak token lifetime
- how to verify the active session deadline through `/api/me`
- how to inspect the underlying opaque session through `GET /api/debug/session`

## API and Runtime Reference

### Runtime values that matter

- `DRAFT_REVIEW_AUTH_MODE`
- `DRAFT_REVIEW_AUTH_SESSION_SECRET`
- `DRAFT_REVIEW_AUTH_SESSION_TTL`
- `DRAFT_REVIEW_AUTH_SESSION_SLIDING_RENEWAL`
- `DRAFT_REVIEW_AUTH_SESSION_RENEW_BEFORE`
- `DRAFT_REVIEW_OIDC_ISSUER_URL`
- `DRAFT_REVIEW_OIDC_CLIENT_ID`
- `DRAFT_REVIEW_OIDC_CLIENT_SECRET`
- `DRAFT_REVIEW_OIDC_REDIRECT_URL`

### Routes involved

- `GET /auth/login`
- `GET /auth/callback`
- `GET /api/debug/session`
- `GET /auth/logout`
- `GET /auth/logout/callback`
- `GET /api/me`

## Risks and Edge Cases

### Risk 1: App session TTL set too long

If the cookie lifetime is too long, a stolen browser session remains valid longer than
desired. The fix should be durable, not effectively permanent.

### Risk 2: Sliding renewal implemented too broadly

If every polling request renews the cookie, sessions may remain active forever without
meaningful user action. Renewal should be threshold-based or limited to user-driven
traffic.

### Risk 3: App session and provider session diverge

After decoupling, it becomes possible that:

- Draft Review local session is still valid
- Keycloak access token is already expired

That is acceptable for a simple app-side session design as long as Keycloak remains
the login authority and Draft Review keeps explicit logout support. If stronger
synchronization is needed later, that is the point to add refresh or introspection.

### Risk 4: Deploy-time logout from secret rotation

This risk remains even if the TTL design is fixed. Operators must keep the signing
secret stable.

## Open Questions

- What exact app TTL is acceptable for authoring UX and security?
- Should the medium-term implementation stop at absolute TTL, or later add sliding renewal before refresh-token work?
- Should `/api/me` expose both app session expiry and original provider token expiry
  for easier debugging?
- Is the hosted Coolify session secret already confirmed stable across deploys?

## References

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)
- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)
- [config.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/config.go)
- [http.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go)
- [main.tf](/home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted/main.tf)
- [main.tf](/home/manuel/code/wesen/terraform/keycloak/modules/realm-base/main.tf)
- [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md)
