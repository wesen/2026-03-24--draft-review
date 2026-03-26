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
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/config.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go
    - /home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted/main.tf
    - /home/manuel/code/wesen/terraform/keycloak/modules/realm-base/main.tf
    - /home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md
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

The recommended fix is to keep Keycloak as the identity provider but stop using the
raw OIDC token expiry as the application cookie expiry. Instead, Draft Review should
define its own session TTL, write cookies against that policy, keep explicit logout
behavior intact, and verify that the hosted session signing secret remains stable
across deploys.

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

Conceptually, the current code does this:

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

### 3. Session cookie manager

Cookie writing and reading are handled in:

- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)

The important behavior is:

- `WriteSession` takes `claims.ExpiresAt`
- it computes `MaxAge` from that timestamp
- it writes `Expires` and `MaxAge` on the cookie
- `ReadSession` rejects the session once that time is reached

So once the copied expiry arrives, the browser session is over from Draft Review's
perspective, regardless of whether the broader Keycloak login might still be alive.

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

If `DRAFT_REVIEW_AUTH_SESSION_SECRET` changes, all existing session cookies become
invalid immediately. That would cause logout across a deploy or restart even if the
TTL problem were fixed.

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

### Recommended flow

```text
browser -> /auth/login
        -> Keycloak login page
        -> /auth/callback
        -> verify token
        -> create Draft Review app session using app TTL
        -> write cookie using app TTL
```

Diagram:

```text
Keycloak login success
        |
        v
Draft Review creates app session policy
        |
        +--> browser cookie expiry based on app TTL
        |
        +--> optional future revalidation / refresh flow
```

The important change is that Draft Review becomes responsible for its own browser
session duration.

## Proposed Solution

The recommended near-term design is to give Draft Review an application-defined
session policy.

### Proposed behavior

- keep Keycloak for login and logout
- verify OIDC tokens exactly as today
- stop copying `token.Expiry` into the app cookie
- introduce an app session TTL, for example `12h`
- optionally add sliding renewal later
- verify the hosted session secret is stable across restarts

### Why this solves the user-facing issue

If Keycloak currently issues short-lived access tokens, the current code makes the
app feel like it has short-lived logins. Once the app defines its own session TTL, a
successful login results in a stable author session suitable for an editing workflow.

### Minimal viable settings

Recommended first settings:

- `auth-session-ttl = 12h`
- optional later: `auth-session-sliding = true`

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

### Decision 4: Keep `/api/me` as the main debugging surface

Rationale:

- it already exposes session metadata
- it gives an immediate way to confirm effective expiry
- it is simple for operators and developers to use

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
- optional `auth-session-sliding`

Recommended environment shape:

- `DRAFT_REVIEW_AUTH_SESSION_TTL=12h`
- `DRAFT_REVIEW_AUTH_SESSION_SLIDING=true`

The intern should parse these once and choose defaults that are easy to reason about.

### Step 3: Change OIDC callback session creation

Update:

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)

Current conceptual code:

```go
expiry := token.Expiry
sessionClaims.ExpiresAt = expiry.UTC()
```

Recommended conceptual code:

```go
issuedAt := now()
sessionClaims.IssuedAt = issuedAt
sessionClaims.ExpiresAt = issuedAt.Add(settings.SessionTTL)
```

Important note:

Do not remove token verification. The app must still verify the Keycloak callback.
Only the local cookie lifetime should change.

### Step 4: Optionally add sliding renewal

If the product wants inactivity-based sessions:

- when an authenticated request arrives
- if the cookie is nearing expiry
- rewrite it with a fresh `ExpiresAt`

Pseudocode:

```text
claims = readSession()
if slidingSessionsEnabled and claims.expiresSoon():
    claims.ExpiresAt = now + sessionTTL
    writeSession(claims)
```

This should not rewrite the cookie on every single request unless that behavior is
deliberately chosen.

### Step 5: Keep logout intact

Retain the existing logout shape in:

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)
- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)

Requirements:

- local cookie is cleared
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

- callback uses app TTL instead of provider token expiry
- expired cookies are rejected correctly
- sliding renewal extends nearing-expiry sessions if enabled
- logout still clears the cookie

### Step 8: Update deployment documentation

Update:

- [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md)

Document:

- the new app TTL setting
- the difference between app session TTL and Keycloak token lifetime
- how to verify the active session deadline through `/api/me`

## API and Runtime Reference

### Runtime values that matter

- `DRAFT_REVIEW_AUTH_MODE`
- `DRAFT_REVIEW_AUTH_SESSION_SECRET`
- `DRAFT_REVIEW_AUTH_SESSION_TTL`
- `DRAFT_REVIEW_OIDC_ISSUER_URL`
- `DRAFT_REVIEW_OIDC_CLIENT_ID`
- `DRAFT_REVIEW_OIDC_CLIENT_SECRET`
- `DRAFT_REVIEW_OIDC_REDIRECT_URL`

### Routes involved

- `GET /auth/login`
- `GET /auth/callback`
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
- Should the first fix be absolute TTL only, or absolute TTL plus sliding renewal?
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
