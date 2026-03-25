# Tasks

## Done

- [x] Create ticket workspace `DR-002`
- [x] Analyze the current frontend runtime, mock API layer, and screen specification
- [x] Write the backend architecture and implementation guide
- [x] Record the investigation diary
- [x] Relate key repository files to the ticket docs
- [x] Validate the ticket with `docmgr doctor`
- [x] Upload the document bundle to reMarkable

## Implementation Queue

### Phase 0: Runtime Scaffold

- [x] Create `go.mod` and baseline Go workspace metadata
- [x] Create `cmd/draft-review/main.go` root command with Glazed logging setup
- [x] Wire embedded Glazed help-system integration at the CLI root
- [x] Create `cmd/draft-review/cmds/serve.go`
- [x] Create `cmd/draft-review/cmds/migrate_up.go`
- [x] Create `cmd/draft-review/cmds/migrate_status.go`
- [x] Create `cmd/draft-review/cmds/seed_dev.go`
- [x] Add backend README/run instructions
- [x] Document docker-compose startup, migration, seed, and local `go run` flows
- [x] Document current API surface and temporary implementation limits
- [x] Add Glazed help entries for overview, local development, auth, and database workflow

### Phase 1: Configuration And Database Foundation

- [x] Add `pkg/config/sql.go` with Clay-style SQL connection section helpers
- [x] Add embedded SQL connection flag schema under `pkg/config/flags/sql-connection.yaml`
- [x] Add `pkg/config/backend.go` for app-local settings like auto-migrate and frontend proxy URL
- [x] Add `pkg/db/postgres.go` with `pgxpool` bootstrap
- [x] Add `pkg/db/migrations.go` with embedded migration runner
- [x] Add `pkg/db/migrations/0001_init.sql`
- [x] Add `pkg/db/migrations/0002_seed_reaction_types.sql`
- [x] Add dev seed/reset workflow

### Phase 2: HTTP Server Foundation

- [x] Add `pkg/server/http.go`
- [x] Add `GET /healthz`
- [x] Add `GET /api/info`
- [x] Add dependency wiring for database-backed services
- [x] Add graceful shutdown handling in `serve`

### Phase 3: Author Auth

- [x] Add `pkg/auth/types.go`
- [x] Add `pkg/auth/config.go` with a Glazed auth section
- [x] Add signed browser session helpers
- [x] Add OIDC / Keycloak login flow helpers
- [x] Add `GET /api/me`
- [x] Add `/auth/login`
- [x] Add `/auth/callback`
- [x] Add `/auth/logout`
- [x] Add `/auth/logout/callback`
- [x] Add `auth-mode=dev` fallback for local work without Keycloak
- [x] Add local Keycloak bootstrap/run instructions
- [x] Add authenticated author bootstrap from OIDC claims to local article ownership
- [x] Protect article mutation endpoints with authenticated author identity
- [x] Scope author article list/detail endpoints to the authenticated owner
- [x] Add manual OIDC callback smoke coverage against a live Keycloak realm

### Phase 4: Articles And Sharing

- [x] Add `pkg/articles/types.go`
- [x] Add `pkg/articles/postgres.go`
- [x] Add `pkg/articles/service.go`
- [x] Add `GET /api/articles`
- [x] Add JSON request validation for article writes
- [x] Add `POST /api/articles`
- [x] Add `GET /api/articles/{id}`
- [x] Add `PATCH /api/articles/{id}`
- [x] Persist section reorder/create/delete operations against `article_sections`
- [x] Support article metadata updates (`title`, `author`, `intro`, `status`)
- [x] Seed default reaction types for newly created articles
- [x] Add article version creation flow
- [x] Add article section editing flow
- [x] Add share-token reset flow
- [x] Add invite creation flow

### Phase 5: Reader Sessions And Feedback

- [x] Add `pkg/reviewlinks` package
- [x] Add `pkg/reviews/types.go`
- [x] Add `pkg/reviews/postgres.go`
- [x] Add `pkg/reviews/service.go`
- [x] Add `GET /api/r/{token}`
- [x] Add `POST /api/r/{token}/start`
- [x] Add `POST /api/reviews/{sessionId}/progress`
- [x] Add `POST /api/reviews/{sessionId}/reactions`
- [x] Add `POST /api/reviews/{sessionId}/summary`
- [x] Return article payloads shaped for the reader experience, not the author dashboard
- [x] Persist paragraph-level progress and reaction authorship details

### Phase 6: Analytics And Reader Management

- [x] Add `pkg/analytics/postgres.go`
- [x] Add `pkg/analytics/service.go`
- [x] Add `GET /api/articles/{id}/analytics`
- [x] Add `GET /api/articles/{id}/feedback`
- [x] Add `GET /api/articles/{id}/readers`
- [x] Add cross-article reader management query surface
- [x] Add draft-killer heuristic implementation
- [x] Add export/report generation stub
- [x] Replace current placeholder reader/reaction endpoints with real query-backed responses

### Phase 7: Frontend Integration And Operations

- [x] Align frontend RTK Query clients with the real Go backend contracts
- [x] Add `docker-compose.yml` for backend plus PostgreSQL
- [x] Add Dockerfile for the Go service
- [x] Add backend unit tests
- [ ] Add database integration tests
- [x] Add HTTP handler tests
- [x] Add local developer runbook
- [x] Add at least one end-to-end smoke check covering migrate -> seed -> serve -> article API

### Phase 8: Hosted Deployment Readiness

- [x] Embed the production frontend into the Go binary for non-dev serving
- [x] Add SPA fallback and static asset serving for `/`, `/assets/*`, and reader routes
- [x] Add Docker multi-stage build for frontend plus backend packaging
- [x] Add Docker runtime dependencies needed for health checks
- [x] Add deployment docs for Coolify and hosted OIDC configuration
- [x] Add SQL config env fallback for container-style DSN injection
- [x] Finalize hosted Keycloak realm settings in `~/code/wesen/terraform`
- [x] Apply the hosted `draft-review` Terraform environment
- [x] Configure or update the Coolify application for `draft-review`
- [x] Verify the live hosted browser login flow against `auth.scapegoat.dev`
- [x] Verify the live hosted browser logout flow against `auth.scapegoat.dev`
- [x] Fix the hosted frontend auth-origin fallback so login/logout stay on the current app origin
- [x] Verify the live hosted new-article flow and add a visible empty-dashboard create action
- [x] Set HTML shell cache headers so hosted deploys do not strand browsers on stale frontend bundles
- [x] Add a generic share-link path alongside the per-email invite flow
- [x] Fix article editor textarea so `Enter` inserts line breaks while drafting
