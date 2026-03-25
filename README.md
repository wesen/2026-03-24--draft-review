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
  --dsn postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable \
  --auto-migrate \
  --frontend-dev-proxy-url http://127.0.0.1:5173
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
- `GET /api/articles`
- `GET /api/articles/{id}`
- `GET /api/articles/{id}/readers`
- `GET /api/articles/{id}/reactions`

The article read routes are backed by PostgreSQL. The reader and reaction routes are still placeholders and currently return empty arrays.

## Local Development Notes

- Docker Compose starts PostgreSQL on `127.0.0.1:5432` and the backend on `127.0.0.1:8080`.
- The `serve` command can auto-run embedded migrations with `--auto-migrate`.
- `seed dev` inserts a stable local author plus one sample article and its first version.
- The backend currently uses a temporary development owner for article creation because author authentication has not been implemented yet.

## Known Gaps

- Author auth endpoints are not implemented yet.
- Article write endpoints are still in progress.
- Reader session, invite, analytics, and feedback APIs are not implemented yet.
- Frontend RTK Query clients still assume the original mock API behavior in some areas.
