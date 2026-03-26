# Changelog

## 2026-03-26

- Initial workspace created
- Added a detailed design and implementation guide explaining why Draft Review's
  hosted browser session likely expires too quickly and how to decouple the app
  session lifetime from the raw Keycloak token lifetime
- Split the auth roadmap into two tickets: DR-010 now owns the medium-term opaque
  server-side session implementation, while DR-011 captures the full long-term
  refresh-token renewal architecture
- Added a DR-010 implementation diary and expanded DR-010 into concrete
  implementation tasks covering repository work, session manager replacement,
  callback changes, request auth resolution, logout, tests, and validation
- Implemented the first DR-010 code slice: opaque author-session persistence using
  the existing `author_sessions` table, OIDC callback-driven local user + session
  creation, server-side session lookup for `/api/me` and author routes, and logout
  revocation of the backing session row
- Updated the repo runbook, hosted deployment notes, and DR-010 design guide so they
  now describe the opaque server-side session model and the new `auth-session-ttl`
  setting instead of the old stateless cookie recommendation

## 2026-03-26

Added author-session activity tracking with last_used_at persistence and touch-on-read session updates.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/postgres.go — Loads and updates last_used_at in Postgres
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go — Touches server-side sessions on authenticated reads
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0006_author_session_activity.sql — Adds last_used_at to persisted author sessions


## 2026-03-26

Added sliding opaque-session renewal with configurable threshold and cookie refresh on active authenticated requests.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/config.go — Adds session renewal settings and validation
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/auth/session.go — Renews session expiry and reissues cookies near expiry
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http.go — Passes response writer through auth resolution so renewal can refresh cookies

