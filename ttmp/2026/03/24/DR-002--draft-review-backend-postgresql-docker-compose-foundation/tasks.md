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
- [x] Create `cmd/draft-review/cmds/serve.go`
- [x] Create `cmd/draft-review/cmds/migrate_up.go`
- [x] Create `cmd/draft-review/cmds/migrate_status.go`
- [x] Create `cmd/draft-review/cmds/seed_dev.go`
- [x] Add backend README/run instructions
- [x] Document docker-compose startup, migration, seed, and local `go run` flows
- [x] Document current API surface and temporary implementation limits

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
- [ ] Add local Keycloak bootstrap/run instructions
- [ ] Add authenticated author bootstrap from OIDC claims to local article ownership
- [ ] Protect article mutation endpoints with authenticated author identity

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
- [ ] Add article version creation flow
- [ ] Add article section editing flow
- [ ] Add share-token reset flow
- [ ] Add invite creation flow

### Phase 5: Reader Sessions And Feedback

- [ ] Add `pkg/reviewlinks` package
- [ ] Add `pkg/reviews/types.go`
- [ ] Add `pkg/reviews/postgres.go`
- [ ] Add `pkg/reviews/service.go`
- [ ] Add `GET /api/r/{token}`
- [ ] Add `POST /api/r/{token}/start`
- [ ] Add `POST /api/reviews/{sessionId}/progress`
- [ ] Add `POST /api/reviews/{sessionId}/reactions`
- [ ] Add `POST /api/reviews/{sessionId}/summary`
- [ ] Return article payloads shaped for the reader experience, not the author dashboard
- [ ] Persist paragraph-level progress and reaction authorship details

### Phase 6: Analytics And Reader Management

- [ ] Add `pkg/analytics/postgres.go`
- [ ] Add `pkg/analytics/service.go`
- [ ] Add `GET /api/articles/{id}/analytics`
- [ ] Add `GET /api/articles/{id}/feedback`
- [ ] Add `GET /api/articles/{id}/readers`
- [ ] Add cross-article reader management query surface
- [ ] Add draft-killer heuristic implementation
- [ ] Add export/report generation stub
- [ ] Replace current placeholder reader/reaction endpoints with real query-backed responses

### Phase 7: Frontend Integration And Operations

- [ ] Align frontend RTK Query clients with the real Go backend contracts
- [ ] Add `docker-compose.yml` for backend plus PostgreSQL
- [ ] Add Dockerfile for the Go service
- [ ] Add backend unit tests
- [ ] Add database integration tests
- [ ] Add HTTP handler tests
- [ ] Add local developer runbook
- [ ] Add at least one end-to-end smoke check covering migrate -> seed -> serve -> article API
