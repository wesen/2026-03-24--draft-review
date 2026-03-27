---
Title: "Local Development Workflow"
Slug: "local-development"
Short: "Run the Go backend, PostgreSQL, Vite frontend, and optional Keycloak stack together."
Topics:
- development
- backend
- frontend
- oidc
Commands:
- serve
- seed
- migrate
Flags:
- frontend-dev-proxy-url
- auto-migrate
- auth-mode
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This page covers the normal end-to-end development loop for Draft Review, how the Go backend and Vite frontend are combined during development, and why the backend origin is the correct browser entrypoint for auth-aware testing.

The working model is two processes plus PostgreSQL. Vite compiles the frontend and serves hot-reload assets, while `draft-review serve` owns the API, auth endpoints, and the browser origin used for OIDC callbacks. When `--frontend-dev-proxy-url http://127.0.0.1:5173` is set, the backend forwards non-API browser routes such as `/` and `/r/<token>` to the Vite dev server. That keeps the browser on the backend origin at `:8080`, which is important because `/auth/*` and the callback flow also live there.

Image uploads add one more local concern: `draft-review serve` now writes uploaded
article media under `--media-root` (default `.draft-review/media`). In local
development the default path is fine. In hosted environments that path must point at
persisted storage such as a mounted Coolify volume.

Use this sequence for the standard local loop:

1. Start PostgreSQL with `docker compose up -d postgres`.
2. Seed the local database with `go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'`.
3. Start the backend with `go run ./cmd/draft-review serve --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --auto-migrate --auth-mode dev --media-root .draft-review/media --frontend-dev-proxy-url http://127.0.0.1:5173`.
4. Start Vite with `cd frontend && npm run dev`.
5. Open `http://127.0.0.1:8080/` in the browser.

For OIDC work, replace `auth-mode=dev` with `auth-mode=oidc` and provide the issuer, client, secret, redirect URL, and session secret. The repository `Makefile` includes `make local-keycloak-up` and `make run-local-oidc` to shorten that path.

The most important operational rule is that the full browser flow should be exercised through the backend origin, not the raw Vite origin. Vite proxies `/api`, which is useful for isolated UI iteration, but the authoritative end-to-end test path is `http://127.0.0.1:8080/` because that covers API requests, auth redirects, and the callback cookie on one origin.

Problem | Cause | Solution
--- | --- | ---
The page loads on `:5173` but login flow breaks | The browser is running on the Vite origin instead of the backend origin | Open `http://127.0.0.1:8080/` for real auth testing
The backend serves JSON but not the frontend shell | `--frontend-dev-proxy-url` was not set or Vite is not running | Start Vite and pass the proxy URL to `serve`
Seed or migrate fails during local boot | PostgreSQL is not running or the DSN points to the wrong port | Start `docker compose up -d postgres` and verify the `--dsn` value
Uploaded images disappear between runs | `--media-root` points at a temporary directory | Use the default local path or a durable mounted path and avoid ephemeral temp dirs

See Also:

- `draft-review help draft-review-overview`
- `draft-review help auth-modes`
- `draft-review help database-workflow`
