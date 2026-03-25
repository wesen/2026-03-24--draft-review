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
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' \
  --auth-mode oidc \
  --auth-session-secret local-session-secret \
  --oidc-issuer-url http://127.0.0.1:18080/realms/draft-review-dev \
  --oidc-client-id draft-review-web \
  --oidc-client-secret draft-review-web-secret \
  --oidc-redirect-url http://127.0.0.1:8080/auth/callback
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

## Local Development Notes

- Docker Compose starts PostgreSQL on `127.0.0.1:5432` and the backend on `127.0.0.1:8080`.
- The `serve` command can auto-run embedded migrations with `--auto-migrate`.
- `auth-mode=dev` gives a synthetic local author identity through `/api/me`.
- `auth-mode=oidc` expects a Keycloak-compatible issuer and signs its own HTTP-only browser session cookie after callback.
- `seed dev` inserts a stable local author plus one sample article and its first version.
- Article creation still uses a temporary development owner until article ownership is bound to authenticated OIDC claims.

## Known Gaps

- Article mutation endpoints are not yet bound to the authenticated OIDC author identity.
- Local Keycloak bootstrap instructions have not been added yet.
- Article version creation is still in progress; `PATCH /api/articles/{id}` currently updates the current version in place.
- Reader session, invite, analytics, and feedback APIs are not implemented yet.
- Frontend RTK Query clients still assume the original mock API behavior in some areas.
