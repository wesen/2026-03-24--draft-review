---
Title: Long-term durable auth sessions with server-side refresh-token renewal guide
Ticket: DR-011
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
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/postgres.go
    - /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql
ExternalSources: []
Summary: Detailed future-state design for Draft Review auth using opaque app sessions, encrypted refresh tokens, renewal, and revocation.
LastUpdated: 2026-03-26T15:50:00-04:00
WhatFor: Capture the full long-term auth target beyond the current implementation slice.
WhenToUse: Use when planning the final durable auth architecture, not the immediate medium-term repair.
---

# Long-term durable auth sessions with server-side refresh-token renewal guide

## Executive Summary

The ideal long-term auth architecture for Draft Review is a full server-side session
system where the browser only holds an opaque session cookie, while the backend owns
session state, encrypted refresh tokens, token renewal, revocation, and auditability.

This is better than both of the simpler approaches:

- better than the old stateless signed-cookie approach because the browser no longer
  carries the full identity payload and the app can revoke sessions centrally
- better than the medium-term opaque-session-only repair because it allows silent
  renewal against Keycloak without forcing frequent re-login

The final architecture should look like this:

```text
Browser
  -> opaque app session cookie
Draft Review backend
  -> author_sessions table
  -> encrypted refresh token store
  -> local user and claim cache
Keycloak
  -> login
  -> refresh
  -> logout / revocation
```

This ticket is intentionally future-state. It is not the immediate implementation
slice. The immediate implementation work belongs in DR-010.

## Problem Statement

Draft Review needs durable author sessions, but it should not rely forever on either:

- copying provider token expiry into a signed client-side cookie
- or a basic opaque session table with no renewal model

The first model gives poor UX. The second model improves durability but still lacks
some important long-term properties:

- no silent renewal after the initial login
- no clear place to manage provider token rotation
- limited visibility into active devices and revocation
- weak support for stronger security policies later

If Draft Review becomes a real product, operators and authors will eventually need:

- long-lived but controlled sessions
- logout across devices
- revocation after suspected compromise
- a clear story for deploys, provider policy changes, and incident response

## Proposed Solution

The long-term solution has five main parts.

### 1. Opaque browser cookie

The browser should only hold a random session token. It should not hold the full
identity payload or provider tokens.

### 2. Server-side app session record

Each author session should be stored server-side with:

- session ID
- user ID
- hashed session token
- created at
- last used at
- idle expiry
- absolute max expiry
- revoked at
- user agent and IP metadata if desired

### 3. Encrypted refresh-token material

The backend should store the Keycloak refresh token encrypted at rest so it can renew
provider access without forcing the user through login again.

### 4. Renewal policy

The app should use:

- short-lived provider access/id tokens
- backend-managed refresh on demand
- sliding idle expiry for the app session
- absolute max lifetime for the app session

### 5. Revocation and operator tooling

The system should support:

- logout current session
- logout all sessions
- revoke a specific session
- inspect active sessions in an internal operator view or debug endpoint

## System Design

### Data model

The current schema already has an `author_sessions` table in:

- [0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql)

For the long-term design, that table likely needs to grow with fields like:

- `last_used_at`
- `idle_expires_at`
- `absolute_expires_at`
- `refresh_token_ciphertext`
- `refresh_token_key_id`
- `provider_session_subject` if we want explicit provider linkage beyond user row data
- optional `user_agent`
- optional `ip_address`

### Request flow

```text
1. Browser sends opaque session cookie
2. Backend hashes token and loads author session row
3. Backend checks revoked / idle / max expiry
4. Backend refreshes provider token if needed
5. Backend loads local user identity
6. Request proceeds
```

### Login flow

```text
1. Browser completes Keycloak login
2. Backend exchanges auth code
3. Backend verifies id_token
4. Backend ensures local user exists
5. Backend encrypts and stores refresh token
6. Backend creates author session row
7. Backend sends opaque cookie
```

### Logout flow

```text
1. Browser requests logout
2. Backend revokes or marks session revoked
3. Backend clears local cookie
4. Backend optionally hits Keycloak logout endpoint
5. Browser returns unauthenticated
```

## Design Decisions

### Decision 1: Store refresh tokens encrypted, never plaintext

Refresh tokens are high-value credentials. If the database is leaked, plaintext refresh
tokens would be much worse than opaque session token hashes.

### Decision 2: Separate idle timeout and absolute lifetime

This gives good UX without creating effectively permanent sessions.

Example:

- idle timeout: 12 hours
- absolute lifetime: 14 days

### Decision 3: Keep app session policy independent from provider token TTL

The provider token TTL should influence renewal behavior, not directly define the app
session lifetime.

### Decision 4: Keep revocation server-side

This is the main benefit of moving away from a fully stateless cookie model.

## Alternatives Considered

### Alternative A: Long-lived opaque sessions without refresh tokens forever

Pros:

- simpler than full renewal architecture

Cons:

- eventually forces re-login once the app session dies
- no way to silently renew provider state
- weaker fit for a polished product

Verdict:

- acceptable medium-term, not ideal long-term

### Alternative B: Stateless signed JWT sessions forever

Pros:

- no session table

Cons:

- weak revocation story
- browser carries too much auth state
- less flexible for future controls

Verdict:

- not recommended long-term

## Implementation Plan

The final implementation should be broken into phases.

### Phase 1: Medium-term base

Complete DR-010:

- opaque session cookie
- server-side session lookup
- no refresh-token renewal yet

### Phase 2: Schema extension

Add the fields needed for:

- refresh token storage
- idle vs absolute expiry
- session metadata

### Phase 3: Encryption support

Introduce encryption for refresh tokens, ideally with:

- a key-encryption strategy
- key rotation support
- clear separation between encryption config and session logic

### Phase 4: Refresh flow

Teach the auth layer to:

- detect expiring provider credentials
- use the stored refresh token
- update stored token material if Keycloak rotates it
- fail gracefully when renewal is no longer allowed

### Phase 5: Revocation and tooling

Add:

- logout current session
- logout all sessions
- revoke session by ID
- session inspection/debug support

## Open Questions

- Where should encryption keys live in the long-term deployment model?
- Should refresh happen synchronously on request or via a cached renewal step?
- Do we want session/device naming for authors later?
- How much operator-facing session inspection is appropriate for this app?

## References

- [oidc.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/oidc.go)
- [session.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go)
- [postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/postgres.go)
- [0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql)
