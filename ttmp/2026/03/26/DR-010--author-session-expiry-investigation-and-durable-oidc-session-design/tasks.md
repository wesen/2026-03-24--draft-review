# Tasks

## TODO

- [x] Create a focused ticket for the short-lived hosted author session issue
- [x] Inspect the Draft Review auth flow, cookie handling, and hosted Keycloak config
- [x] Write a detailed intern-facing design and implementation guide
- [x] Validate the ticket bundle with `docmgr doctor`
- [x] Upload the ticket bundle to reMarkable

## Implementation Tasks

- [x] Create a separate long-term design ticket for refresh-token renewal so DR-010 can stay focused on opaque server-side sessions
- [x] Add a DR-010 implementation diary document
- [x] Add auth configuration fields for app-managed session TTL and any session-manager options needed for opaque sessions
- [x] Extend the auth repository with author-session persistence methods using the existing `author_sessions` table
- [x] Replace the stateless signed-claims cookie manager with an opaque-token session manager that hashes tokens and loads users from Postgres
- [x] Update the OIDC callback to ensure the authenticated local user exists before creating the browser session
- [x] Update request auth resolution in `pkg/server/http.go` so `/api/me` and author routes derive identity from the server-side session lookup
- [x] Update logout to revoke or invalidate the server-side session and clear the browser cookie
- [x] Add and update tests for callback login, cookie/session reading, `/api/me`, and logout
- [x] Verify that existing dev auth mode behavior still works
- [x] Update DR-010 docs, diary, and changelog as each slice lands
- [x] Run `go test ./cmd/... ./pkg/...` and `docmgr doctor --ticket DR-010 --stale-after 30`

## Follow-Up Session Hardening Tasks

- [x] Add `last_used_at` tracking to `author_sessions` and update it on authenticated requests
- [ ] Add sliding session renewal with explicit config for enablement and renewal threshold
- [ ] Add a current-session inspection endpoint for debugging the opaque auth session state
- [ ] Extend auth and handler tests to cover the three follow-up features
- [ ] Update DR-010 docs and diary with the deployed follow-up behavior
