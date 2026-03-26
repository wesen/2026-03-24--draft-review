# Draft Review Backend

This repository now includes the first Go backend scaffold for Draft Review.

The backend uses:

- Go with a Glazed/Cobra CLI entrypoint
- PostgreSQL via `pgxpool`
- Embedded SQL migrations
- Embedded frontend assets for production and Docker deployment
- Docker Compose for a local backend-plus-database workflow

Hosted deployment guidance now lives in:

- [docs/deployments/draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md)
- [docs/deployments/draft-review-coolify-playbook.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify-playbook.md)

## Current Commands

Start the full local stack with Docker Compose:

```bash
docker compose up --build
```

That path now serves the browser app directly from the Go container at `http://127.0.0.1:8080/`; it does not require Vite.

Run PostgreSQL only, then start the backend locally:

```bash
docker compose up postgres
go run ./cmd/draft-review serve \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' \
  --auto-migrate \
  --auth-mode dev \
  --frontend-dev-proxy-url http://127.0.0.1:5173
```

Run the frontend against the real backend:

```bash
cd frontend
npm run dev
```

If you want the old mock-only frontend behavior instead, opt back into MSW explicitly:

```bash
cd frontend
VITE_USE_MSW=1 npm run dev
```

Canonical real-app dev loop:

1. Start Postgres.
2. Start the backend with `--frontend-dev-proxy-url http://127.0.0.1:5173`.
3. Start Vite in `frontend/`.
4. Open `http://127.0.0.1:8080/`.

Use `http://127.0.0.1:8080/` for full end-to-end testing because that keeps the browser on the backend origin for `/api`, `/auth`, and the OIDC callback. The Vite origin at `http://127.0.0.1:5173/` is still useful for UI iteration and now proxies `/api`, but browser auth should be exercised through `8080`.

Start the backend against Keycloak / OIDC:

```bash
go run ./cmd/draft-review serve \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:15432/draft_review?sslmode=disable' \
  --auth-mode oidc \
  --auth-session-secret local-session-secret \
  --auth-session-ttl 12h \
  --auth-session-sliding-renewal true \
  --auth-session-renew-before 1h \
  --oidc-issuer-url http://127.0.0.1:18080/realms/draft-review-dev \
  --oidc-client-id draft-review-web \
  --oidc-client-secret draft-review-web-secret \
  --oidc-redirect-url http://127.0.0.1:8080/auth/callback
```

In OIDC mode, the app now uses an opaque browser session token backed by the
`author_sessions` table. `--auth-session-ttl` controls the Draft Review session
lifetime separately from the upstream Keycloak token lifetime.
`--auth-session-sliding-renewal` and `--auth-session-renew-before` control whether
active sessions extend themselves before expiry.

Or use the local helper targets:

```bash
make local-keycloak-up
make run-local-oidc
```

If those default local ports are already occupied, override them explicitly:

```bash
make local-keycloak-up PG_PORT=25432 KEYCLOAK_PORT=18190
make run-local-oidc PG_PORT=25432 KEYCLOAK_PORT=18190 APP_PORT=8081
```

Stop the local stack:

```bash
make local-keycloak-down
```

Apply migrations without starting the HTTP server:

```bash
go run ./cmd/draft-review migrate up \
  --dsn postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable
```

Inspect migration status:

```bash
go run ./cmd/draft-review migrate status \
  --dsn postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable
```

Seed the development dataset:

```bash
go run ./cmd/draft-review seed dev \
  --dsn postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable
```

Run tests:

```bash
go test ./cmd/... ./pkg/...
```

Browse the embedded Glazed help pages:

```bash
go run ./cmd/draft-review help
go run ./cmd/draft-review help --topics
go run ./cmd/draft-review help local-development
go run ./cmd/draft-review help auth-modes
```

The root command now carries the shared logging flags and the embedded help system. That means the same logger configuration applies across `serve`, `migrate`, `seed`, and the help command itself.

## Current HTTP API

The current backend exposes these routes:

- `GET /healthz`
- `GET /api/info`
- `GET /api/me`
- `GET /api/debug/session`
- `GET /api/articles`
- `POST /api/articles`
- `GET /api/articles/{id}`
- `PATCH /api/articles/{id}`
- `POST /api/articles/{id}/versions`
- `POST /api/articles/{id}/share-token`
- `POST /api/articles/{id}/invite`
- `GET /api/articles/{id}/readers`
- `GET /api/articles/{id}/reactions`
- `GET /api/articles/{id}/analytics`
- `GET /api/articles/{id}/feedback`
- `POST /api/articles/{id}/export`
- `GET /api/readers`
- `GET /api/r/{token}`
- `POST /api/r/{token}/start`
- `POST /api/reviews/{sessionId}/progress`
- `POST /api/reviews/{sessionId}/reactions`
- `POST /api/reviews/{sessionId}/summary`

When `auth-mode=oidc`, the backend also exposes:

- `GET /auth/login`
- `GET /auth/callback`
- `GET /auth/logout`
- `GET /auth/logout/callback`

The article, sharing, reader-session, and analytics routes are backed by PostgreSQL. `POST /api/articles/{id}/export` is intentionally still a stub response that reserves the route shape without generating files yet.

## Local Keycloak Setup

The repository includes a standalone Keycloak stack in [docker-compose.local.yml](/home/manuel/code/wesen/2026-03-24--draft-review/docker-compose.local.yml) and a realm import in [draft-review-dev-realm.json](/home/manuel/code/wesen/2026-03-24--draft-review/dev/keycloak/realm-import/draft-review-dev-realm.json).

Local defaults:

- Keycloak admin: `http://127.0.0.1:18080/admin`
- Keycloak bootstrap admin username: `admin`
- Keycloak bootstrap admin password: `admin`
- Realm: `draft-review-dev`
- OIDC client: `draft-review-web`
- OIDC client secret: `draft-review-web-secret`
- Test user username: `author`
- Test user password: `secret`
- Test user email: `author@example.com`

Useful commands:

```bash
make local-keycloak-up
make local-keycloak-config
make local-keycloak-down
make run-local-dev
make run-local-oidc
```

Port overrides are supported on all helper targets:

```bash
make local-keycloak-config PG_PORT=25432 KEYCLOAK_PORT=18190
make local-keycloak-down PG_PORT=25432 KEYCLOAK_PORT=18190
```

The live OIDC smoke path validated during implementation was:

1. `make local-keycloak-up PG_PORT=25432 KEYCLOAK_PORT=18190`
2. `go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:25432/draft_review?sslmode=disable'`
3. `make run-local-oidc PG_PORT=25432 KEYCLOAK_PORT=18190 APP_PORT=8081`
4. Open `http://127.0.0.1:8081/auth/login?return_to=%2Fapi%2Fme`
5. Log in with `author` / `secret`
6. Confirm `/api/me` returns an authenticated OIDC identity
7. Create an article through `/api/articles`
8. Visit `/auth/logout?return_to=%2Fapi%2Fme` and confirm `/api/me` returns `authenticated: false`

## Local Development Notes

- Docker Compose starts PostgreSQL on `127.0.0.1:5432` and the backend on `127.0.0.1:8080`.
- The production container now embeds the built frontend, so `/`, `/assets/*`, and reader SPA routes work without Vite.
- The `serve` command can auto-run embedded migrations with `--auto-migrate`.
- `auth-mode=dev` gives a synthetic local author identity through `/api/me`.
- `auth-mode=oidc` expects a Keycloak-compatible issuer and signs its own HTTP-only browser session cookie after callback.
- OIDC mode now records `author_sessions.last_used_at`, can renew active sessions before expiry, and exposes `GET /api/debug/session` for current-session inspection.
- `seed dev` inserts a stable local author plus one sample article and its first version.
- Author article routes now resolve the current browser identity into a local `users` row and scope article access by `owner_user_id`.
- The frontend now targets the real Go backend by default; set `VITE_USE_MSW=1` only when you intentionally want the legacy mock layer.
- The backend now really proxies browser routes to the frontend dev server when `--frontend-dev-proxy-url` is set, so `/` and `/r/:token` can be tested through the backend origin.
- Vite now listens on `0.0.0.0` and proxies `/api`, `/auth`, and `/healthz` to `VITE_BACKEND_ORIGIN` or `http://127.0.0.1:8080` by default.

## Known Gaps

- The frontend reader flow now starts review sessions and persists progress, reactions, and an empty summary, but it still does not expose the richer summary fields such as recommendability or notify-on-new-version.
- `POST /api/articles/{id}/export` is only a stub acknowledgement, not a real report generator.
- Database integration tests are still missing; the current validation story is unit tests plus live smoke runs.
- `seed dev` already runs migrations internally, so do not run `migrate up` concurrently against the same fresh database.
