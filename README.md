# Draft Review Backend

This repository now includes the first Go backend scaffold for Draft Review.

The backend uses:

- Go with a Glazed/Cobra CLI entrypoint
- PostgreSQL via `pgxpool`
- Embedded SQL migrations
- Docker Compose for a local backend-plus-database workflow

## Current Commands

Start the full local stack with Docker Compose:

```bash
docker compose up --build
```

Run PostgreSQL only, then start the backend locally:

```bash
docker compose up postgres
go run ./cmd/draft-review serve \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' \
  --auto-migrate \
  --auth-mode dev \
  --frontend-dev-proxy-url http://127.0.0.1:5173
```

Start the backend against Keycloak / OIDC:

```bash
go run ./cmd/draft-review serve \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:15432/draft_review?sslmode=disable' \
  --auth-mode oidc \
  --auth-session-secret local-session-secret \
  --oidc-issuer-url http://127.0.0.1:18080/realms/draft-review-dev \
  --oidc-client-id draft-review-web \
  --oidc-client-secret draft-review-web-secret \
  --oidc-redirect-url http://127.0.0.1:8080/auth/callback
```

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

## Current HTTP API

The current backend exposes these routes:

- `GET /healthz`
- `GET /api/info`
- `GET /api/me`
- `GET /api/articles`
- `POST /api/articles`
- `GET /api/articles/{id}`
- `PATCH /api/articles/{id}`
- `GET /api/articles/{id}/readers`
- `GET /api/articles/{id}/reactions`

When `auth-mode=oidc`, the backend also exposes:

- `GET /auth/login`
- `GET /auth/callback`
- `GET /auth/logout`
- `GET /auth/logout/callback`

The article read and write routes are backed by PostgreSQL. The reader and reaction routes are still placeholders and currently return empty arrays.

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
- The `serve` command can auto-run embedded migrations with `--auto-migrate`.
- `auth-mode=dev` gives a synthetic local author identity through `/api/me`.
- `auth-mode=oidc` expects a Keycloak-compatible issuer and signs its own HTTP-only browser session cookie after callback.
- `seed dev` inserts a stable local author plus one sample article and its first version.
- Author article routes now resolve the current browser identity into a local `users` row and scope article access by `owner_user_id`.

## Known Gaps

- The frontend still does not consume `/api/me` or drive the `/auth/*` browser flow.
- Article version creation is still in progress; `PATCH /api/articles/{id}` currently updates the current version in place.
- Reader session, invite, analytics, and feedback APIs are not implemented yet.
- Frontend RTK Query clients still assume the original mock API behavior in some areas.
- `seed dev` already runs migrations internally, so do not run `migrate up` concurrently against the same fresh database.
