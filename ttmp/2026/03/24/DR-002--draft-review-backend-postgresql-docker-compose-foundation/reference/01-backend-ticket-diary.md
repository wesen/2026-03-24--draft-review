---
Title: Backend Ticket Diary
Ticket: DR-002
Status: active
Topics:
    - backend
    - postgresql
    - docker-compose
    - api
    - auth
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: imports/draft-review-screen-spec.md
      Note: Primary external input that shaped the documented backend requirements
    - Path: ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md
      Note: Primary planning artifact produced during this diary entry
    - Path: ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md
      Note: Ticket overview updated during this diary entry
    - Path: ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md
      Note: Task checklist updated during this diary entry
ExternalSources: []
Summary: Chronological diary for the DR-002 backend planning ticket, including repository investigation, document creation, validation, and delivery.
LastUpdated: 2026-03-24T23:28:00-04:00
WhatFor: Recording how the backend planning ticket was assembled and validated.
WhenToUse: Use this diary when reviewing the documentation work, reconstructing decisions, or continuing backend implementation later.
---


# Backend Ticket Diary

## Goal

Capture the chronological research, documentation, validation, and delivery steps used to produce the backend planning ticket for Draft Review.

## Step 1: Create The Ticket And Gather Evidence

I started by treating this as a documentation-and-planning task rather than jumping into backend scaffolding. The repository already contained a frontend implementation ticket and a non-trivial screen specification, so the first job was to create a fresh ticket workspace and then map the current application shape before writing recommendations.

The main outcome of this step was a grounded evidence base: current routes, mock API contracts, in-memory entity shapes, and the broader screen-spec requirements that the backend must eventually satisfy. That made the later design work much more concrete and reduced the risk of writing a generic "backend plan" that did not actually fit this application.

### Prompt Context

**User prompt (verbatim):**
```text
Create a new ticket to create the backend for this applicaton, using postgresql (use docker compose).

Keep a diary as you work.

Create a detailed analysis / design / implementation guide that is very detailed for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file
  references.
  It should be very clear and detailed. Store in the ticket and the nupload to remarkable.
```

**Assistant interpretation:** Create a new `docmgr` ticket for the Draft Review backend, document the design in detail for a new intern, keep a proper diary, and publish the resulting bundle to reMarkable.

**Inferred user intent:** Establish an implementation-ready backend plan that is detailed enough for handoff and onboarding, while also fitting the repository’s existing ticket/document workflow.

**Commit (code):** N/A — documentation-only work in ticket files

### What I did
- Read the `ticket-research-docmgr-remarkable` and `diary` skill instructions from `/home/manuel/.codex/skills/.../SKILL.md`.
- Inspected the repository layout with `rg --files`.
- Checked current ticket/doc status with `docmgr status --summary-only`.
- Read the existing frontend implementation ticket at `ttmp/2026/03/24/DR-001--draft-review-mvp-react-redux-storybook-implementation/...`.
- Inspected the key frontend runtime files and screen spec:
- `frontend/src/App.tsx`
- `frontend/src/api/baseApi.ts`
- `frontend/src/api/articleApi.ts`
- `frontend/src/api/readerApi.ts`
- `frontend/src/mocks/db.ts`
- `frontend/src/mocks/handlers.ts`
- `frontend/src/author/Dashboard.tsx`
- `frontend/src/author/ArticleReader.tsx`
- `frontend/src/reader/ReaderPage.tsx`
- `frontend/src/types/article.ts`
- `frontend/src/types/reader.ts`
- `frontend/src/types/reaction.ts`
- `imports/draft-review-screen-spec.md`
- Created ticket `DR-002` and added the primary design doc plus the diary document using `docmgr`.

### Why
- The backend design needed to be evidence-based.
- The existing screen spec includes more backend responsibilities than the currently implemented frontend surfaces, so I needed both the current code and the broader product definition.
- Creating the ticket first ensured all later work landed in the correct ticket paths.

### What worked
- The repository structure and `docmgr` workflow were straightforward to inspect.
- The mock API handlers made the current HTTP contract immediately visible.
- The screen spec clearly exposed future backend requirements like auth, sharing modes, versioning, summaries, and analytics.

### What didn't work
- No failures in this step.

### What I learned
- The current frontend is already opinionated about route shape and `/api` prefixing.
- The current mock model is intentionally minimal and does not represent the eventual domain model well enough for persistence-heavy features.
- The screen spec is the source of truth for several backend concepts that do not yet exist in code.

### What was tricky to build
- The tricky part was not the tooling; it was separating observed facts from proposed architecture. The repo does not contain a backend, so any framework recommendation had to be clearly marked as an inference rather than presented as if it were already chosen.

### What warrants a second pair of eyes
- The chosen backend stack recommendation: TypeScript + Fastify + Drizzle is reasonable, but it is still a recommendation, not a repository constraint.
- The boundary between first-phase backend scope and later admin/team/billing work should stay deliberate.

### What should be done in the future
- Keep the same evidence-first approach when implementation begins so endpoint and schema work stays tied to current frontend needs.

### Code review instructions
- Start with the design doc and compare its "current state" sections against the referenced frontend files.
- Verify that every major backend requirement in the design doc can be traced either to an implemented frontend file or to the screen spec.

### Technical details
- Commands run:
```text
docmgr status --summary-only
docmgr ticket create-ticket --ticket DR-002 --title "Draft Review Backend - PostgreSQL Docker Compose Foundation" --topics backend,postgresql,docker-compose,api,auth
docmgr doc add --ticket DR-002 --doc-type design-doc --title "Draft Review Backend Architecture and Implementation Guide"
docmgr doc add --ticket DR-002 --doc-type reference --title "Backend Ticket Diary"
```

## Step 9: Add A Local Keycloak Dev Stack And Prove The OIDC Flow End To End

Once the backend could sign browser sessions and map OIDC claims onto local authors, the remaining gap was operational: there was no repository-local way to boot Keycloak, no imported development realm, and no live proof that the full browser round-trip actually worked. I closed that by adding a standalone compose file, a realm import, helper `make` targets, and a real browser smoke test against Keycloak rather than stopping at unit tests.

This step was important because OIDC code that passes isolated tests can still fail in practice due to redirect URI mismatches, issuer timing, cookies, or local port conflicts. I wanted one verified path inside this repo that another engineer can reproduce without piecing together commands from `hair-booking`.

### Prompt Context

**User prompt (verbatim):**
```text
we want to do auth using keycloak, look at /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking
```

followed by:

```text
ok, continue,
```

and then:

```text
go ahead, don't stop, go do it all
```

**Assistant interpretation:** finish the Keycloak direction properly by adding local bootstrap assets, proving the live OIDC flow, documenting the exact runbook, and then continue with the next backend phases.

**Inferred user intent:** do not stop at static auth code or partial docs; leave behind a runnable local OIDC stack and a trustworthy diary of the exact commands, failures, and fixes.

**Commit (code):** pending during diary capture for this step

### What I did
- Added `docker-compose.local.yml` with:
- one Postgres service for the app,
- one Postgres service for Keycloak,
- one Keycloak service running `start-dev --import-realm`.
- Added `dev/keycloak/realm-import/draft-review-dev-realm.json` containing:
- realm `draft-review-dev`,
- confidential OIDC client `draft-review-web`,
- redirect URIs for `127.0.0.1` and `localhost` on ports `8080` and `8081`,
- test user `author` with password `secret`.
- Added `Makefile` helper targets:
- `local-keycloak-up`
- `local-keycloak-down`
- `local-keycloak-config`
- `run-local-dev`
- `run-local-oidc`
- Fixed the helper targets so `PG_PORT` and `KEYCLOAK_PORT` overrides are propagated into `docker compose`.
- Updated `README.md` with:
- local Keycloak setup details,
- helper target usage,
- port override examples,
- a reproducible live smoke path.
- Booted the local stack on conflict-free ports:
- Postgres `25432`,
- Keycloak `18190`,
- backend `8081`.
- Ran a real browser-auth smoke test through Playwright:
- opened `/auth/login`,
- signed in through Keycloak as `author`,
- landed on `/api/me`,
- confirmed an authenticated OIDC identity,
- created an article with the authenticated browser session,
- logged out through `/auth/logout`,
- confirmed `/api/me` became unauthenticated again.

### Why
- The backend needed a repository-local OIDC harness, not just borrowed assumptions from `hair-booking`.
- The main failure modes for OIDC are operational rather than algorithmic.
- A new intern should be able to run one documented path and see the full login and article-ownership flow themselves.

### What worked
- The imported realm and client were enough for a clean Keycloak login flow.
- The backend’s signed browser session cookie worked correctly after callback.
- The local-author bootstrap logic correctly mapped the OIDC user onto a local `users` row and article ownership.
- The `Makefile` helper flow is now good enough for repeatable local use.

### What didn't work
- The first attempt to start the local compose stack on the default app Postgres port failed because another local service already owned `15432`.
- The first attempt to reuse the default Keycloak port also failed because `18080` was already occupied.
- Early compose revisions with explicit `container_name` values caused stale-container friction during recreation, so I removed the fixed names.
- I briefly ran `migrate up` and `seed dev` concurrently against the same fresh database and hit:
```text
ERROR: duplicate key value violates unique constraint "pg_type_typname_nsp_index"
```
That was self-inflicted; `seed dev` already runs migrations internally.
- One initial backend start in OIDC mode failed because Keycloak’s discovery endpoint was not fully ready yet, so the backend saw a connection reset while booting the authenticator.

### What I learned
- The helper targets need explicit port override plumbing. Relying on raw shell env overrides makes the runbook harder to trust.
- `seed dev` should be treated as "migrate plus seed", not as a pure data insert command.
- The OIDC path was stable once the compose stack was healthy and the redirect URIs matched the actual app port.

### What was tricky to build
- The subtle part was not the realm JSON itself; it was making the local stack resilient to port collisions that are common on this machine.
- Another subtle point was logout behavior: Keycloak shows its own intermediate logout confirmation screen, so the smoke test needed one additional click before checking `/api/me`.

### What warrants a second pair of eyes
- Whether the local realm should also preconfigure alternate redirect URIs for other common dev ports such as Vite-driven author app ports.
- Whether we want a dedicated scripted smoke command in the repo instead of keeping the live browser verification in the diary/runbook only.

### What should be done in the future
- Add a scripted OIDC smoke check once the frontend starts consuming `/api/me` and `/auth/*`.
- Keep the realm import aligned with any future production Keycloak claim mapping changes.

### Code review instructions
- Review `docker-compose.local.yml`, the realm import JSON, and the `Makefile` targets together; they are one operational unit.
- Re-run the documented smoke flow with an alternate port set to confirm the override path still works.
- Confirm that article ownership after OIDC login uses the local `users` row produced by the auth service rather than a synthetic fallback.

### Technical details
- Main commands run:
```text
make local-keycloak-config
make local-keycloak-up PG_PORT=25432 KEYCLOAK_PORT=18190
go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:25432/draft_review?sslmode=disable'
make run-local-oidc PG_PORT=25432 KEYCLOAK_PORT=18190 APP_PORT=8081
```
- Manual browser flow validated:
```text
http://127.0.0.1:8081/auth/login?return_to=%2Fapi%2Fme
username: author
password: secret
```
- Verified outcomes:
```text
GET /api/me -> authenticated: true, authMode: "oidc", email: "author@example.com"
POST /api/articles -> creates an article owned by the authenticated local user row
GET /auth/logout?return_to=%2Fapi%2Fme -> logout confirmation through Keycloak
GET /api/me -> authenticated: false
```

## Step 10: Replace The Placeholder Reader Backend With Real Persistence

After the Keycloak/OIDC slice was stable, I moved into the next missing chunk of real product behavior: invite links, share-token reset, reader link resolution, review-session persistence, author analytics, and version creation. Up to this point the backend still had placeholder readers and reactions endpoints, even though the schema already had the tables needed to support the actual flow.

I implemented this as one cohesive vertical slice because these features depend on each other. A share token without `/api/r/{token}` is not useful. A reader link without a review session cannot persist reactions. Analytics without readers and reactions is just another placeholder. Splitting those into too many tiny steps would have produced a lot of integration churn with little value.

### Prompt Context

**User prompt (verbatim):**
```text
go ahead, don't stop, go do it all
```

**Assistant interpretation:** continue through the next substantial backend milestones instead of stopping after the auth infrastructure.

**Inferred user intent:** turn the remaining placeholder backend surfaces into real PostgreSQL-backed behavior and keep the ticket current while doing it.

**Commit (code):** pending during diary capture for this step

### What I did
- Added `pkg/db/migrations/0004_review_paragraph_progress.sql` to persist paragraph-level progress separately from section-level progress.
- Added `pkg/reviewlinks` with:
- share-token reset,
- invite creation,
- token resolution for both article share tokens and unique invite tokens,
- reader-facing article payload shaping.
- Added `pkg/reviews` with:
- review-session start,
- progress persistence,
- paragraph progress persistence,
- reaction persistence,
- summary submission.
- Added `pkg/analytics` with:
- real `GET /api/articles/{id}/readers`,
- real `GET /api/articles/{id}/reactions`,
- `GET /api/articles/{id}/analytics`,
- cross-article `GET /api/readers`,
- draft-killer heuristic calculation.
- Added `POST /api/articles/{id}/versions` to clone the current article version, copy sections, and activate the new version.
- Added `POST /api/articles/{id}/export` as an authenticated stub route that reserves the export contract without pretending report generation is done.
- Updated `pkg/server/http.go` to wire all of the new packages and routes.
- Added focused tests in:
- `pkg/reviewlinks/service_test.go`
- `pkg/reviews/service_test.go`
- `pkg/server/http_test.go`
- Ran a live dev-mode smoke test covering:
- share-token reset,
- invite creation,
- `GET /api/r/{token}`,
- `POST /api/r/{token}/start`,
- `POST /api/reviews/{sessionId}/progress`,
- `POST /api/reviews/{sessionId}/reactions`,
- `POST /api/reviews/{sessionId}/summary`,
- `GET /api/articles/{id}/readers`,
- `GET /api/articles/{id}/reactions`,
- `GET /api/articles/{id}/analytics`,
- `GET /api/readers`,
- `POST /api/articles/{id}/versions`.

### Why
- The existing placeholder routes were blocking real frontend-to-backend alignment.
- The schema already expected versioning, review sessions, reactions, and analytics, so the missing piece was repository and HTTP wiring rather than new product invention.
- Version creation belonged in the same implementation window because the article/version tables were already live and it was the most obvious unused author-side capability.

### What worked
- The first full smoke path succeeded in dev auth mode once the server was restarted on the updated code.
- The analytics queries returned sensible live values after one invited reader session and one reaction.
- Version creation cleanly cloned the current section set and advanced the article’s current version pointer.
- The service-layer validation tests were enough to pin the new payload rules without building a much heavier test harness.

### What didn't work
- The first implementation of `/api/r/{token}` leaked internal resolver fields like `ArticleID`, `InviteID`, and `AllowAnonymous` because the resolver struct did not mark those fields `json:"-"`. I caught that immediately during the live curl smoke and fixed the JSON tags before proceeding.
- I did not add automated database integration tests yet. The SQL behavior is currently covered by compile-time checks, unit tests around the service layer, and live smoke runs.

### What I learned
- The live smoke mattered again. The leaked internal resolver fields would not have been visible from pure service tests.
- The existing schema was already close enough that most of the work was about correctly projecting between reader-facing, author-facing, and internal storage shapes.
- A small export stub route is better than leaving the export task invisible; it makes the gap explicit and keeps the route contract available for later implementation.

### What was tricky to build
- The trickiest part was shaping the same underlying data differently for different audiences:
- `GET /api/r/{token}` needs only reader-facing `reader` and `article` fields,
- author analytics needs normalized reader and reaction records,
- internal services still need hidden fields like `ArticleVersionID`, invite identity, and access rules.
- Another subtle point was deciding where to compute display names and avatars. I kept that logic close to `reviewlinks` because invites, reader-directory views, and reader link resolution all need consistent derivation from email addresses.

### What warrants a second pair of eyes
- The current cross-article reader directory route shape (`GET /api/readers`) is sensible, but it is still app-local rather than driven by an existing frontend screen.
- The version-cloning endpoint currently copies the active version’s intro, author note, and sections wholesale. That is the right default for now, but it may need more explicit UI choices later.
- The export stub is intentionally minimal and should not be mistaken for a full report pipeline.

### What should be done in the future
- Wire the frontend reader experience to the real review-session endpoints instead of keeping progress and reactions in local React state.
- Add automated database integration tests around invites, progress, reactions, analytics, and version cloning.
- Replace the export stub with real Markdown/PDF generation once the desired export formats are settled.

### Code review instructions
- Review `pkg/reviewlinks`, `pkg/reviews`, `pkg/analytics`, and `pkg/server/http.go` together; they form one end-to-end slice.
- Re-run the live smoke path in dev auth mode to confirm the repository code still matches the route contracts.
- Pay special attention to the JSON surface of `/api/r/{token}` and the ownership checks on author routes.

### Technical details
- Key commands run:
```text
go test ./cmd/... ./pkg/...
make local-keycloak-up PG_PORT=25432 KEYCLOAK_PORT=18190
go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:25432/draft_review?sslmode=disable'
make run-local-dev PG_PORT=25432 APP_PORT=8082
```
- Representative smoke calls:
```text
POST /api/articles/{id}/share-token
POST /api/articles/{id}/invite
GET /api/r/{token}
POST /api/r/{token}/start
POST /api/reviews/{sessionId}/progress
POST /api/reviews/{sessionId}/reactions
POST /api/reviews/{sessionId}/summary
GET /api/articles/{id}/readers
GET /api/articles/{id}/reactions
GET /api/articles/{id}/analytics
GET /api/readers
POST /api/articles/{id}/versions
```

## Step 11: Align The Frontend With The Real Backend Routes

Once the backend routes were stable, the remaining practical issue was that the frontend still behaved like a mock-only prototype in development. MSW always started in dev mode, the author app did not consult `/api/me`, and the reader page kept all of its progress and feedback state locally. I treated this as a separate slice because it crosses from backend implementation into app integration and deserved its own validation cycle.

The goal here was not to redesign the frontend. It was to make the current UI actually capable of talking to the backend that now exists, while still preserving an explicit way to opt back into MSW for mock-only work.

### Prompt Context

**User prompt (verbatim):**
```text
go ahead, don't stop, go do it all
```

**Assistant interpretation:** continue into the frontend integration slice once the backend endpoints exist, instead of leaving the app hard-wired to the mock layer.

**Inferred user intent:** the repository should behave like one application, not a backend sidecar plus an unrelated prototype frontend.

**Commit (code):** pending during diary capture for this step

### What I did
- Added `frontend/src/api/authApi.ts` with `/api/me`.
- Expanded `frontend/src/api/readerApi.ts` with:
- `startReview`
- `updateReviewProgress`
- `addReviewReaction`
- `submitReviewSummary`
- Expanded `frontend/src/api/articleApi.ts` with:
- `createVersion`
- `exportArticle`
- Changed `frontend/src/main.tsx` so MSW only starts when `VITE_USE_MSW=1` is set explicitly.
- Updated `frontend/src/app/AuthorApp.tsx` to:
- load `/api/me` when not in mock mode,
- block the author UI behind login when unauthenticated,
- add an account/logout menu,
- wire the invite dialog to the real invite mutation.
- Updated `frontend/src/app/ReaderApp.tsx` to pass the active review token through to `ReaderPage`.
- Updated `frontend/src/reader/ReaderPage.tsx` to:
- start a review session when reading begins,
- persist section progress,
- persist new reactions,
- submit a minimal summary on finish.
- Built the frontend with `npm run build`.

### Why
- Leaving MSW as the default meant the repo could still appear "working" while silently bypassing the new backend.
- `/api/me` is now the backend entrypoint for author auth state, so the author app needed to respect it.
- The reader flow was the most important place to move from local-only state to real persistence because the new review-session endpoints now exist specifically for that purpose.

### What worked
- The frontend TypeScript build passed cleanly after the new RTK Query endpoints and auth gate were added.
- The MSW opt-in approach keeps mock mode available without hiding the real backend by default.
- The reader flow could be upgraded incrementally without rewriting the existing UI components.

### What didn't work
- I did not add a richer reader summary UI in this step. The reader now submits a minimal empty summary on completion so the backend session can close cleanly, but the recommendability and notify-on-new-version fields still have no frontend controls.
- I did not touch the unrelated untracked `frontend/README.md` because it is outside the tracked repository state for this task.

### What I learned
- The cleanest migration path was to make MSW explicit rather than trying to support both real and mock behavior invisibly.
- Author auth integration can stay narrow and still be useful: one `/api/me` query plus login/logout redirects gets the current app much closer to the real deployment model.
- Reader persistence did not require a major state-management refactor; the existing component callbacks were already enough to call the new backend endpoints.

### What was tricky to build
- The most delicate part was balancing real-backend alignment with not destroying the existing prototype workflow. Making MSW opt-in was a better compromise than trying to keep the old implicit mock behavior and the new backend behavior active at the same time.
- The second subtle point was React hook ordering in `AuthorApp`: the auth query needed to gate the article queries without turning the component into a conditional-hook bug.

### What warrants a second pair of eyes
- The reader completion flow now submits an empty summary payload. That is operationally correct but product-incomplete.
- The author login gate should be verified once a real frontend dev server is run against Keycloak/OIDC mode end to end.

### What should be done in the future
- Add a real summary form in the reader completion dialog for recommendability and notification opt-in.
- Add a small frontend runbook note for `VITE_USE_MSW=1` versus real-backend mode once the tracked frontend docs are in scope.

### Code review instructions
- Review the API client additions and the app boot changes together; the point is to make the frontend hit the real backend by default.
- Confirm that `VITE_USE_MSW=1` still restores mock mode intentionally.
- Run `npm run build` and then boot the frontend against the Go backend to check the login gate and reader flow manually.

### Technical details
- Main command run:
```text
cd frontend
npm run build
```

## Step 12: Wire The Missing Dev Proxy Path

After the frontend/backend alignment pass, one practical blocker still remained: there was no reliable way to run the real app on one browser origin during development. The backend already exposed a `frontend-dev-proxy-url` setting, but it did not actually proxy any browser routes. That meant auth and API testing still depended on awkward split-origin behavior. I fixed that by implementing the backend reverse proxy for non-API routes and tightening the Vite-side dev behavior around it.

### Prompt Context

**User prompt (verbatim):**
```text
alright wire up the proxy so I can test this for real.
```

**Assistant interpretation:** make the actual development proxy path work, not just document it, so the browser can exercise the real app against the real backend.

**Inferred user intent:** be able to start the dev loop and open one URL that behaves like the real application, especially for auth-sensitive browser testing.

**Commit (code):** pending during diary capture for this step

### What I did
- Added a real catch-all reverse proxy in `pkg/server/http.go` that forwards non-API browser routes to the configured frontend dev server when `frontend-dev-proxy-url` is set.
- Added an HTTP test in `pkg/server/http_test.go` proving that `/r/...` is proxied to an upstream dev server instead of returning `404`.
- Updated `frontend/vite.config.ts` to:
- listen on `0.0.0.0`,
- proxy `/api` to `VITE_BACKEND_ORIGIN` or `http://127.0.0.1:8080`.
- Updated `frontend/src/app/AuthorApp.tsx` so login/logout always target the backend origin explicitly instead of relying on Vite to proxy `/auth`.
- Updated `README.md` to define the canonical real-app dev loop:
- backend on `8080`,
- Vite on `5173`,
- browser opened on `http://127.0.0.1:8080/`.

### Why
- Full browser auth flows should run on the backend origin because the OIDC redirect and cookie behavior belong there.
- The repo already had the right configuration knob; the missing work was implementation and a clear local workflow.
- A split setup can still be useful, but it should be deliberate: `5173` for UI iteration, `8080` for end-to-end app testing.

### What worked
- The backend origin now returns the Vite HTML shell on `/` and `/r/...` when the dev proxy is enabled.
- `http://127.0.0.1:5173/api/me` now works through Vite’s `/api` proxy.
- `go test ./cmd/... ./pkg/...` still passes.
- `npm run build` still passes.

### What didn't work
- Proxying `/auth` through Vite was the wrong abstraction for this app. It produced a `431 Request Header Fields Too Large` response during direct testing and was removed.
- The correct answer was to treat `8080` as the canonical full-app origin and keep `5173` as a UI-focused dev origin.

### What I learned
- The missing backend proxy was the actual blocker, not the frontend API clients.
- Browser auth is much cleaner when the app is accessed through the backend origin from the start.
- Vite still benefits from an `/api` proxy, but `/auth` should stay anchored on the backend origin for this setup.

### What was tricky to build
- The subtle part was not writing `httputil.NewSingleHostReverseProxy`; it was deciding where auth should live in development. Once I verified that Vite’s `/auth` proxy path was unhealthy, the right split became obvious:
- `8080` for end-to-end testing,
- `5173` for direct frontend iteration with API proxying only.

### What warrants a second pair of eyes
- The backend reverse proxy should be exercised with a real login/logout browser flow in OIDC mode after the user runs the full stack locally.
- If embedded frontend serving is added later, the fallback precedence between embedded assets and the dev proxy should be reviewed carefully.

### What should be done in the future
- Consider adding a single `make dev` target that starts backend plus Vite together.
- If desired, add a small health or debug route that reports whether the frontend dev proxy target is reachable.

### Code review instructions
- Read `pkg/server/http.go`, `pkg/server/http_test.go`, and `frontend/vite.config.ts` together; they define the new development routing model.
- Validate the documented loop by opening `http://127.0.0.1:8080/` with the backend proxy enabled.

### Technical details
- Live checks run:
```text
docker compose up -d postgres
go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --auto-migrate --auth-mode dev --listen-host 127.0.0.1 --listen-port 8080 --frontend-dev-proxy-url http://127.0.0.1:5173
cd frontend && npm run dev -- --host 0.0.0.0
curl -s http://127.0.0.1:8080/ | sed -n '1,5p'
curl -s http://127.0.0.1:8080/r/tok-demo -I
curl -s http://127.0.0.1:5173/api/me
```

## Step 2: Write The Guide, Update Bookkeeping, Validate, And Deliver

After the ticket existed and the evidence was collected, I replaced the generated placeholders with the actual deliverables: the intern-facing backend design guide, a useful ticket overview, and this diary. I then related the key repository files, validated the ticket, and uploaded the final document bundle to reMarkable.

This step converted the investigation into something operational: a ticket another engineer can open, understand, and use to implement the backend in phases.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Finish the ticket artifacts end to end, including bookkeeping, validation, and reMarkable delivery.

**Inferred user intent:** Produce a handoff-grade backend planning packet, not just notes.

**Commit (code):** N/A — documentation-only work in ticket files

### What I did
- Replaced placeholder text in:
- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md`
- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md`
- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md`
- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md`

## Step 3: Rewrite The Backend Guide Around Go, Glazed, And Clay SQL Sections

After the first version of the guide was complete, the direction changed: the backend should use Go, Glazed for command verbs, and the Clay DB configuration pattern instead of a TypeScript/Fastify/Drizzle stack. I treated that as a full rewrite of the technical recommendation rather than a light edit, because the runtime, package layout, CLI shape, and configuration model all changed.

The core work in this step was to inspect the exact reference implementations you pointed me at, extract the concrete patterns that matter, and rewrite the guide so it speaks in terms of real Go files, real Glazed sections, and real Postgres bootstrap code rather than abstract framework preferences.

### Prompt Context

**User prompt (verbatim):**
```text
update the document entirely, we want to use golang and use the glazed framework for verbs and db configuration layer (see ~/code/wesen/corporate-headquarters/clay for glazed section definitions for db configuration). You can look at projects like /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking
```

**Assistant interpretation:** Rewrite the backend planning document so the recommended implementation is Go-based, uses Glazed verbs, and uses Clay-style section definitions for DB configuration, with `hair-booking` as the runtime reference.

**Inferred user intent:** Align the backend design with the existing Go/Glazed patterns already used elsewhere in the user’s codebase so implementation can follow known conventions instead of a generic stack choice.

**Commit (code):** N/A — documentation-only work in ticket files

### What I did
- Read the Glazed command-authoring skill.
- Inspected Clay SQL configuration files:
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/flags/sql-connection.yaml`
- Inspected `hair-booking` runtime files:
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/config/backend.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/postgres.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go`
- Rewrote the primary design doc to recommend:
- Go,
- Glazed/Cobra verbs,
- Clay-style SQL config sections,
- `pgxpool`,
- embedded SQL migrations,
- a `hair-booking`-style runtime/package layout.
- Updated the ticket index and next-step task list so the ticket no longer points toward the old TypeScript backend direction.

### Why
- The prior version of the design doc used the wrong implementation language and framework direction.
- The new direction is stronger because it matches existing code patterns in the user’s ecosystem rather than introducing a separate stack.

### What worked
- Clay provided exactly the database-section reference pattern needed for the rewrite.
- `hair-booking` provided an immediately reusable model for the Glazed root command, serve command, backend config section, Postgres bootstrap, migration runner, and HTTP server package layout.

### What didn't work
- No tooling failed in this step.
- The main cost was that the first version of the backend recommendation had to be discarded and replaced, rather than incrementally edited.

### What I learned
- The most useful part of the rewrite was not just "use Go", but being able to anchor the recommendation in concrete local examples the implementation can copy from.
- The Clay SQL section pattern is the right level of abstraction for this ticket because it cleanly separates DB configuration from verb logic.

### What was tricky to build
- The tricky part was preserving the valid domain-model work from the earlier guide while replacing the runtime and tooling assumptions. The schema and API requirements still largely hold, but the command layout, configuration layer, migration strategy, and package structure all needed to be restated in Go terms.

### What warrants a second pair of eyes
- Whether Draft Review should import Clay’s SQL section helpers directly or copy the pattern into a narrower app-local package.
- Whether author authentication should be fully app-managed in phase 1 or kept minimal until the article/review flows are stable.

### What should be done in the future
- If implementation starts, create follow-up diary steps for:
- Glazed root/verb scaffolding,
- Clay-style SQL section integration,
- initial migration/schema creation,
- first end-to-end frontend-to-Go API cutover.

### Code review instructions
- Compare the current design doc against the Go reference files in `clay` and `hair-booking`.
- Verify that the revised guide now points to Glazed sections, Cobra wiring, `pgxpool`, embedded migrations, and a Go package layout instead of the old TypeScript stack.

### Technical details
- Key inspection commands included:
```text
rg -n "SectionDefinition|section definition|db config|database|postgres|glazed" /home/manuel/code/wesen/corporate-headquarters/clay -S
rg -n "glazed|cobra|postgres|database|sqlx|gorm|ent|pgx|migrate|verb|cmd|section" /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking -S
nl -ba /home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go | sed -n '1,280p'
nl -ba /home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go | sed -n '1,260p'
nl -ba /home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/flags/sql-connection.yaml | sed -n '1,260p'
nl -ba /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go | sed -n '1,220p'
nl -ba /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go | sed -n '1,280p'
nl -ba /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/config/backend.go | sed -n '1,260p'
nl -ba /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/postgres.go | sed -n '1,260p'
nl -ba /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go | sed -n '1,260p'
nl -ba /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go | sed -n '1,280p'
```

## Step 4: Revalidate And Re-upload The Revised Bundle

After rewriting the guide, I reran `docmgr doctor`, fixed the new topic vocabulary entries, and re-uploaded the ticket bundle to reMarkable. The main operational issue in this step was that the first upload attempt skipped because the bundle already existed remotely under the same name. I resolved that by rerunning the upload with `--force` and then verifying the remote listing again.

This step matters because the ticket is not complete if the stored reMarkable copy still points to the old TypeScript-based version of the guide.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Finish the rewrite all the way through validation and overwrite the previously uploaded bundle so the reMarkable copy matches the new Go/Glazed guide.

**Inferred user intent:** Keep the ticket and the delivered artifact in sync after the architecture direction changed.

**Commit (code):** N/A — documentation-only work in ticket files

### What I did
- Updated ticket file relations and changelog for the Go/Glazed rewrite.
- Added `golang` and `glazed` to the local ticket vocabulary.
- Ran `docmgr doctor --ticket DR-002 --stale-after 30` until it passed cleanly again.
- Ran a dry-run bundle upload.
- Attempted a normal bundle upload and observed that it skipped because the file already existed.
- Re-ran the upload with `--force`.
- Verified the remote listing at `/ai/2026/03/24/DR-002`.

### Why
- The first uploaded bundle represented the earlier backend direction.
- After the rewrite, the remote copy needed to be replaced, not merely left alongside the new local docs.

### What worked
- Vocabulary and related-file fixes were enough to get `docmgr doctor` back to a clean state.
- `remarquee upload bundle --force` successfully replaced the remote document bundle.

### What didn't work
- The first non-dry-run upload skipped because the file already existed:
```text
remarquee upload bundle ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/changelog.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md --name "DR-002 Draft Review Backend Planning Bundle" --remote-dir "/ai/2026/03/24/DR-002" --toc-depth 2
SKIP: DR-002 Draft Review Backend Planning Bundle already exists in /ai/2026/03/24/DR-002 (use --force to overwrite)
```

### What I learned
- For ticket revisions that keep the same bundle name, reMarkable delivery needs an explicit overwrite step.
- External file relations are safer on the design doc than on the ticket index because `docmgr` normalized those index paths into broken relatives during validation.

### What was tricky to build
- The main sharp edge was not the content rewrite itself but keeping validation, related-file metadata, and bundle delivery consistent after the rewrite changed the ticket’s topic vocabulary and reference set.

### What warrants a second pair of eyes
- Whether future revisions should use a versioned bundle name instead of forced overwrite, depending on how much historical artifact preservation matters.

### What should be done in the future
- If the design changes materially again, either:
- upload with `--force`, or
- intentionally version the bundle name instead of silently relying on the previous artifact.

### Code review instructions
- Re-run:
```text
docmgr doctor --ticket DR-002 --stale-after 30
remarquee cloud ls /ai/2026/03/24/DR-002 --long --non-interactive
```
- Confirm that the ticket validates and that the remote bundle still exists after the overwrite.

### Technical details
- Commands run:
```text
docmgr vocab add --category topics --slug golang --description "Go language backend implementation and package design"
docmgr vocab add --category topics --slug glazed --description "Glazed command framework, sections, and CLI verb architecture"
docmgr doctor --ticket DR-002 --stale-after 30
remarquee upload bundle --dry-run ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/changelog.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md --name "DR-002 Draft Review Backend Planning Bundle" --remote-dir "/ai/2026/03/24/DR-002" --toc-depth 2
remarquee upload bundle --force ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/changelog.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md --name "DR-002 Draft Review Backend Planning Bundle" --remote-dir "/ai/2026/03/24/DR-002" --toc-depth 2
remarquee cloud ls /ai/2026/03/24/DR-002 --long --non-interactive
```
- Related the ticket docs to the repository files that informed the design.
- Updated the changelog and task status.
- Ran `docmgr doctor --ticket DR-002 --stale-after 30`.
- Performed a dry-run bundle upload to reMarkable, then the real upload, then verified the remote listing.

### Why
- The ticket needed to be useful to another engineer immediately.
- `docmgr doctor` and the reMarkable upload are part of the expected completion workflow for this kind of ticket.

### What worked
- The generated ticket structure was easy to fill in once the evidence was gathered.
- The design doc format was a good fit for a long intern-facing guide.
- The bundle upload workflow allowed verification before the final upload.

### What didn't work
- The first verification listing ran before the upload process completed:
```text
remarquee cloud ls /ai/2026/03/24/DR-002 --long --non-interactive
Error: no matches for 'DR-002'
```
- I resolved it by waiting for the upload command to finish and then rerunning the listing successfully.

### What I learned
- The best way to make this kind of ticket useful is to write it as both a system explanation and a file-by-file implementation map.
- Keeping the diary in the same ticket makes future continuation much easier.

### What was tricky to build
- The sharp edge was balancing detail with scope control. The screen spec includes significantly more product surface than the current frontend implements, so the guide needed to acknowledge future features without turning the first backend slice into an unbounded rewrite.

### What warrants a second pair of eyes
- The proposed schema for versions, invites, review sessions, and analytics should be reviewed before implementation begins because those tables will shape a large amount of follow-on code.
- The exact DTO contract for reader-start and reaction endpoints should be reviewed alongside frontend integration plans.

### What should be done in the future
- When backend implementation starts, create a follow-up diary step per implementation phase.
- Keep the design doc updated if the chosen backend framework or DTO shapes change.

### Code review instructions
- Review `index.md` first for ticket framing.
- Review the design doc next, especially the schema, API reference, and implementation phases.
- Review the diary last to reconstruct how the ticket was produced and which commands were run.
- Validate with:
```text
docmgr doctor --ticket DR-002 --stale-after 30
```

### Technical details
- Ticket root:
```text
ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation
```
- Validation and delivery commands:
```text
docmgr doc relate --ticket DR-002 --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/imports/draft-review-screen-spec.md:Primary product specification that defines backend requirements beyond the current frontend implementation" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/mocks/handlers.ts:Current mock HTTP contract that the real backend should replace incrementally" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/mocks/db.ts:Current in-memory data model used to infer the first persistent entities"
docmgr doc relate --doc ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/App.tsx:Defines the author and reader route split that the backend must support" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/baseApi.ts:Defines the current /api prefix used by the frontend" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/articleApi.ts:Defines the current article, reader, invite, and reaction client contracts" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/readerApi.ts:Defines token-based reader link resolution on the frontend" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/mocks/handlers.ts:Defines the mock endpoints and payload shapes that the backend should replace" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/mocks/db.ts:Defines the current mock article, reader, and reaction fields" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx:Shows current browser-side analytics and draft-killer logic" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/ReaderPage.tsx:Shows current client-side review progress and reaction behavior" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/imports/draft-review-screen-spec.md:Defines the broader backend requirements including auth, sharing modes, summaries, analytics, and future admin features"
docmgr doc relate --doc ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md:Ticket overview updated during this diary entry" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md:Task checklist updated during this diary entry" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md:Primary planning artifact produced during this diary entry" --file-note "/home/manuel/code/wesen/2026-03-24--draft-review/imports/draft-review-screen-spec.md:Primary external input that shaped the documented backend requirements"
docmgr doctor --ticket DR-002 --stale-after 30
remarquee upload bundle --dry-run ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/changelog.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md --name "DR-002 Draft Review Backend Planning Bundle" --remote-dir "/ai/2026/03/24/DR-002" --toc-depth 2
remarquee upload bundle ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/tasks.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/changelog.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md --name "DR-002 Draft Review Backend Planning Bundle" --remote-dir "/ai/2026/03/24/DR-002" --toc-depth 2
remarquee cloud ls /ai/2026/03/24/DR-002 --long --non-interactive
```

## Step 5: Scaffold The Go Backend Foundation

I moved the ticket from planning into implementation by landing the first backend foundation slice. The goal in this step was not to build product behavior yet; it was to establish a stable Go runtime with Glazed verbs, Clay-style SQL configuration, PostgreSQL bootstrap, embedded migrations, dev seed support, and a minimal HTTP surface that later tasks can extend.

This step is the architectural spine for the rest of the backend. Once it existed, later work no longer had to make open-ended decisions about module layout, command wiring, migration strategy, or how PostgreSQL configuration is passed into the process.

### Prompt Context

**User prompt (verbatim):**
```text
Add detailed tasks to the ticket, then work on them one by one, comitting at appropriate interval, keeping a diary.
```

**Assistant interpretation:** Expand the ticket into an execution plan and start implementing the first backend milestones with focused commits and diary entries.

**Inferred user intent:** Convert the backend effort from pure documentation into real, incremental implementation while keeping the work reviewable and well-documented.

**Commit (code):** `01c8c9fa2d4b1bbcdb9f1bdbed51bb5a956d0bf9` — `"Scaffold Go backend foundation"`

### What I did
- Expanded `tasks.md` into a phase-by-phase implementation queue.
- Added the Go module and dependencies in `go.mod`.
- Added the Glazed/Cobra runtime:
- `cmd/draft-review/main.go`
- `cmd/draft-review/cmds/serve.go`
- `cmd/draft-review/cmds/migrate_up.go`
- `cmd/draft-review/cmds/migrate_status.go`
- `cmd/draft-review/cmds/seed_dev.go`
- Added Clay-style SQL connection config:
- `pkg/config/sql.go`
- `pkg/config/flags/sql-connection.yaml`
- Added app-local backend config:
- `pkg/config/backend.go`
- Added PostgreSQL bootstrap, migration runner, and seed path:
- `pkg/db/postgres.go`
- `pkg/db/migrations.go`
- `pkg/db/seed.go`
- `pkg/db/migrations/0001_init.sql`
- `pkg/db/migrations/0002_seed_reaction_types.sql`
- Added the first HTTP server package and routes:
- `pkg/server/http.go`
- Added `Dockerfile` and `docker-compose.yml`.
- Updated `.gitignore` for Go runtime output.

### Why
- The backend needed a real executable base before feature work on articles, sessions, reactions, and auth could proceed cleanly.
- The Glazed and Clay-style config decisions are foundational, so they belong in the first implementation slice rather than being retrofitted later.

### What worked
- The `hair-booking` command and server pattern translated cleanly into this repo.
- The Clay-style SQL section design fit Draft Review well with only a narrower Postgres-focused field set.
- The migration runner and seed path worked well as a first operational DB layer.

### What didn't work
- `go mod tidy` upgraded the module’s Go version line because `glazed v1.0.5` requires at least Go `1.25.7`:
```text
go: github.com/go-go-golems/glazed@v1.0.5 requires go >= 1.25.7; switching to go1.25.8
```
- This did not block the work, but it was a real environment constraint that had to be accepted in the module metadata.

### What I learned
- The backend foundation is much easier to review when it lands as a single focused slice with commands, config, DB, server, and container assets together.
- `pgxpool` plus embedded SQL migrations is a clean and low-friction first backend approach for this repo.

### What was tricky to build
- The main sharp edge was keeping the first code slice broad enough to support later work without pulling feature logic into it prematurely. The right line was: runtime spine, not domain behavior.

### What warrants a second pair of eyes
- The initial schema in `pkg/db/migrations/0001_init.sql`, because it sets the structural direction for every later API.
- The narrowed local copy of the Clay SQL configuration pattern, especially if the project later wants to depend on Clay directly.

### What should be done in the future
- Implement the first real article endpoints next.
- Add tests around config parsing and migrations once there is a slightly larger behavior surface.

### Code review instructions
- Start with `go.mod`, `cmd/draft-review/main.go`, and `cmd/draft-review/cmds/serve.go`.
- Then review `pkg/config/sql.go`, `pkg/db/migrations.go`, and `pkg/server/http.go`.
- Validate with:
```text
go test ./cmd/... ./pkg/...
```

### Technical details
- Commands run:
```text
mkdir -p cmd/draft-review/cmds pkg/config/flags pkg/db/migrations pkg/server
go mod tidy
gofmt -w cmd/draft-review/main.go cmd/draft-review/cmds/*.go pkg/config/*.go pkg/db/*.go pkg/server/*.go
go test ./cmd/... ./pkg/...
git add .gitignore go.mod go.sum Dockerfile docker-compose.yml cmd pkg
git commit -m "Scaffold Go backend foundation"
git rev-parse HEAD
```

## Step 6: Add The First Real Article Read Endpoints

With the runtime foundation committed, I moved to the next smallest useful slice: article reads. This step adds the first Postgres-backed feature package and wires it into the server so the backend starts replacing the frontend’s MSW article mocks instead of only exposing health/info routes.

I kept this slice intentionally narrow. It adds `GET /api/articles` and `GET /api/articles/{id}` backed by PostgreSQL, and it leaves readers/reactions as temporary empty-array endpoints so the current frontend route shape can be exercised without 404s while the reader/reaction packages remain pending.

### Prompt Context

**User prompt (verbatim):** (same as Step 5)

**Assistant interpretation:** Continue down the execution queue after the foundation commit and land the next practical backend feature slice.

**Inferred user intent:** Make steady, reviewable progress on the actual backend rather than stopping after the runtime scaffold.

**Commit (code):** `845eeb67ff0d75ffa5586037cf82a900f05469ab` — `"Add initial article read endpoints"`

### What I did
- Added `pkg/articles/types.go`.
- Added `pkg/articles/postgres.go`.
- Added `pkg/articles/service.go`.
- Updated `pkg/server/http.go` to:
- wire the article service when the database is available,
- add `GET /api/articles`,
- add `GET /api/articles/{id}`,
- add temporary `GET /api/articles/{id}/readers`,
- add temporary `GET /api/articles/{id}/reactions`.
- Preserved the current frontend DTO shape by serializing section bodies back into `paragraphs[]`.

### Why
- Article reads are the smallest meaningful slice that starts replacing the current frontend mock API.
- The author app already depends on article list and detail responses, so these endpoints create a useful seam for later integration.

### What worked
- The repository/service split in `pkg/articles` fit cleanly into the foundation from Step 5.
- Keeping the frontend DTO shape stable avoided unnecessary integration churn.
- The slice compiled and tested cleanly without reopening any of the runtime decisions from the previous commit.

### What didn't work
- There were no failures in this step.
- The readers and reactions endpoints are still placeholders returning empty arrays; that is intentional, but they remain follow-up work rather than complete implementations.

### What I learned
- The easiest way to replace the existing mocks incrementally is to preserve the frontend response contracts exactly, even if the database schema stores richer internal forms.
- A small feature package plus thin server wiring is a much better review unit than embedding SQL in handlers.

### What was tricky to build
- The main balancing act was choosing a slice that was real enough to be useful but still small enough for a clean commit. Article reads plus placeholder readers/reactions hit that balance better than attempting the entire author surface at once.

### What warrants a second pair of eyes
- The query and serialization flow in `pkg/articles/postgres.go`, especially the section/paragraph conversion logic and the decision to return full sections on the list endpoint.

### What should be done in the future
- Replace the placeholder readers and reactions handlers with real packages and queries.
- Add article creation and update endpoints next.

### Code review instructions
- Start with `pkg/articles/postgres.go`, then `pkg/articles/service.go`, then `pkg/server/http.go`.
- Validate with:
```text
go test ./cmd/... ./pkg/...
```

### Technical details
- Commands run:
```text
mkdir -p pkg/articles
gofmt -w pkg/articles/*.go pkg/server/http.go
go test ./cmd/... ./pkg/...
git add pkg/articles pkg/server/http.go
git commit -m "Add initial article read endpoints"
git rev-parse HEAD
```

## Step 7: Expand The Execution Queue And Add A Backend Runbook

Before taking the next API write slice, I tightened the ticket task list and added a repository-level backend README. The goal of this step was to make the implementation queue explicit enough that each later commit can map cleanly to one or two checklist items rather than a vague phase label.

The README also turns the scaffold into something another engineer can actually run without reconstructing command lines from source. That matters now that the repository has a real Go CLI, Docker Compose setup, and migration/seed commands.

### Prompt Context

**User prompt (verbatim):**
```text
Add detailed tasks to the ticket, then work on them one by one, comitting at appropriate interval, keeping a diary.
```

**Assistant interpretation:** Expand the ticket checklist into more concrete implementation items, then start executing them in small reviewable slices while continuing the diary.

**Inferred user intent:** Keep implementation disciplined and traceable. The ticket should act as a live execution plan, not just a high-level design artifact.

**Commit (code):** pending in this step while the runbook/checklist slice is being assembled

### What I did
- Expanded `tasks.md` with more detailed acceptance items around:
- backend runbook coverage,
- article write validation and persistence behavior,
- reader/session expectations,
- placeholder replacement work,
- future smoke-test coverage.
- Added a top-level `README.md` describing:
- Docker Compose startup,
- local `go run` flows,
- migration and seed commands,
- the current HTTP API surface,
- known temporary implementation limits.
- Marked the runbook-related checklist items complete in the ticket after the README existed.

### Why
- The next implementation work is easier to review when each task has a concrete definition of done.
- A backend without a runbook is expensive for the next engineer to pick up, even if the code compiles.

### What worked
- The existing command surface was already stable enough to document directly.
- Breaking the queue down further made the next API slice more obvious: article writes are the next contract-driven gap.

### What didn't work
- No technical failures in this step.
- The queue is still intentionally broad in later phases like analytics because those packages do not exist yet; those tasks will need another pass once the auth/review flows are real.

### What I learned
- The ticket is already shifting from research artifact to implementation tracker, and that means the task list needs to describe executable slices rather than architecture topics.

### What was tricky to build
- The subtle part was documenting current behavior honestly without implying unfinished endpoints are production-ready. The README needs to help someone run the backend today while still making the temporary gaps obvious.

### What warrants a second pair of eyes
- The README wording around "current API" and "known gaps" should stay aligned with the actual implementation as more endpoints land.

### What should be done in the future
- Keep the README current as each API slice lands.
- Continue converting the remaining broad phase tasks into smaller completion units before or during implementation.

### Code review instructions
- Review `README.md` for command accuracy and alignment with the current CLI.
- Review `tasks.md` to make sure the added checklist items match the actual intended implementation order.

### Technical details
- Commands run:
```text
rg --files -g 'README*'
sed -n '1,240p' docker-compose.yml
sed -n '1,240p' Dockerfile
```

## Step 8: Add Initial Article Write Endpoints And Verify Them Live

With the runbook/checklist slice committed, I moved to the next contract-driven backend gap: article writes. The frontend already has create and update mutations, so this step makes those routes real and keeps the backend moving toward replacing the MSW article mock instead of only reading seeded data.

I kept the scope intentionally narrow. The backend now supports `POST /api/articles` and `PATCH /api/articles/{id}` with request decoding, basic validation, transactional PostgreSQL writes, and section replacement logic. I explicitly did not add version snapshots yet; updates still modify the current article version in place until the dedicated versioning task is implemented.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Continue executing the ticket one slice at a time and commit after a coherent unit of backend behavior lands.

**Inferred user intent:** Replace the frontend article mock surface with real backend behavior in small reviewable steps, while keeping the ticket docs and diary current.

**Commit (code):** pending in this step while the article write slice is being finalized

### What I did
- Extended `pkg/articles` with:
- create/update DTOs,
- validation helpers,
- repository methods for article creation and mutation,
- transactional section replacement logic.
- Updated `pkg/server/http.go` to add:
- `POST /api/articles`,
- `PATCH /api/articles/{id}`,
- JSON request decoding with unknown-field rejection,
- validation and not-found error mapping.
- Updated `README.md` and `tasks.md` so the documented API surface and checklist match the implementation.
- Ran `gofmt` and `go test ./cmd/... ./pkg/...`.
- Performed a live smoke test against local Docker Compose PostgreSQL:
- created an article,
- patched title/status/intro/sections,
- re-read the article,
- confirmed invalid status returns HTTP 400.

### Why
- These routes are already present in the frontend RTK Query layer, so implementing them gives the backend immediate integration value.
- This slice is small enough to review clearly while still exercising real DB mutations.

### What worked
- The existing article repository/service split extended cleanly to write operations.
- Replacing all sections on patch is a simple way to support reorder/create/delete behavior before full versioning exists.
- The live smoke test confirmed the route behavior against a real PostgreSQL instance, not just compilation.

### What didn't work
- My first migration smoke-test command failed because `zsh` treated `?sslmode=disable` as a glob when the DSN was not quoted.
- I also briefly got a misleading empty list response because I ran the first create and list curl calls in parallel, so the read could race the insert.

### What I learned
- The fastest path to frontend alignment is to keep the DTO shape stable and add backend validation under the existing contract, instead of trying to redesign the payloads mid-implementation.
- Even small shell details matter in this repo because the default shell is `zsh`; DSNs with `?` need quoting in docs and smoke tests.

### What was tricky to build
- The main tradeoff was how much "edit article" behavior to support before versioning exists. I chose full current-version replacement for sections because it satisfies the immediate editor contract without pretending version history is already implemented.

### What warrants a second pair of eyes
- The section replacement approach in `pkg/articles/postgres.go`, especially the decision to regenerate IDs for client-side temporary section identifiers.
- The decision to use a temporary development owner record until auth is implemented.

### What should be done in the future
- Add real article version creation instead of mutating the current version in place.
- Add share-token reset and invite creation after the reader/review flows are scaffolded.
- Add automated handler or integration tests for the article write paths instead of relying only on the manual smoke test.

### Code review instructions
- Start with `pkg/articles/service.go`, then `pkg/articles/postgres.go`, then `pkg/server/http.go`.
- Re-run:
```text
go test ./cmd/... ./pkg/...
docker compose up -d postgres
go run ./cmd/draft-review migrate up --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --listen-host 127.0.0.1 --listen-port 8080 --auto-migrate
```
- Then exercise:
```text
curl -sS -X POST http://127.0.0.1:8080/api/articles -H 'Content-Type: application/json' -d '{"title":"API Smoke Article","author":"Integration Check","intro":"Created during smoke testing."}'
curl -sS -X PATCH http://127.0.0.1:8080/api/articles/<id> -H 'Content-Type: application/json' -d '{"title":"API Smoke Article Revised","status":"in_review","intro":"Updated during smoke testing.","sections":[{"id":"s-new-1","title":"Rewritten Opening","paragraphs":["First paragraph.","Second paragraph."]},{"id":"s-new-2","title":"Second Section","paragraphs":["More detail here."]}]}'
curl -sS http://127.0.0.1:8080/api/articles/<id>
curl -sS -o /tmp/draft-review-invalid.json -w '%{http_code}' -X PATCH http://127.0.0.1:8080/api/articles/<id> -H 'Content-Type: application/json' -d '{"status":"broken"}'
docker compose down
```

### Technical details
- Commands run:
```text
gofmt -w pkg/articles/*.go pkg/server/http.go
go test ./cmd/... ./pkg/...
docker compose up -d postgres
go run ./cmd/draft-review migrate up --dsn postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable
go run ./cmd/draft-review migrate up --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --listen-host 127.0.0.1 --listen-port 8080 --auto-migrate
curl -sS -X POST http://127.0.0.1:8080/api/articles -H 'Content-Type: application/json' -d '{"title":"API Smoke Article","author":"Integration Check","intro":"Created during smoke testing."}'
curl -sS http://127.0.0.1:8080/api/articles
curl -sS -X PATCH http://127.0.0.1:8080/api/articles/7e6d7677-a626-48ad-8de4-5a2e78ca5fb4 -H 'Content-Type: application/json' -d '{"title":"API Smoke Article Revised","status":"in_review","intro":"Updated during smoke testing.","sections":[{"id":"s-new-1","title":"Rewritten Opening","paragraphs":["First paragraph.","Second paragraph."]},{"id":"s-new-2","title":"Second Section","paragraphs":["More detail here."]}]}'
curl -sS http://127.0.0.1:8080/api/articles/7e6d7677-a626-48ad-8de4-5a2e78ca5fb4
curl -sS -o /tmp/draft-review-invalid.json -w '%{http_code}' -X PATCH http://127.0.0.1:8080/api/articles/7e6d7677-a626-48ad-8de4-5a2e78ca5fb4 -H 'Content-Type: application/json' -d '{"status":"broken"}'
docker compose down
```

## Step 9: Add Auth Domain Types And Core Session Helpers

After the article write slice, I took one smaller infrastructure task from the auth phase: define the auth-facing types and add the core password/token/session-cookie helpers that future auth handlers will rely on. This is a good stopping point for a small commit because it lands reusable primitives without mixing them into incomplete HTTP endpoint code.

I also added focused unit tests for the helper behavior so this slice improves the test surface instead of only adding passive type definitions.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Keep moving through the ticket queue in coherent slices and commit when a clean boundary is reached.

**Inferred user intent:** Build implementation scaffolding that reduces risk for the next endpoint slices rather than leaving critical auth behavior to be improvised later.

**Commit (code):** pending in this step while the helper slice is being finalized

### What I did
- Added `pkg/auth/types.go` with:
- `User`,
- `Session`,
- `PasswordResetToken`,
- `EmailVerificationToken`.
- Added `pkg/auth/helpers.go` with:
- bcrypt password hashing and comparison,
- opaque token generation,
- token hashing,
- session cookie issuance,
- session cookie expiration helpers.
- Added `pkg/auth/helpers_test.go` covering:
- password hash round-trip,
- token generation and hashing,
- session cookie creation and revocation behavior.
- Marked the first three auth checklist items complete in `tasks.md`.

### Why
- Signup/login/logout handlers will need these primitives immediately.
- Pulling them into a dedicated auth package now keeps later HTTP/database code simpler and more consistent.

### What worked
- The helper slice stayed nicely isolated from the rest of the codebase.
- Adding tests here was cheap and immediately useful because the auth helpers are pure functions.

### What didn't work
- No failures in this step.

### What I learned
- The auth phase will be easier to stage if I keep separating pure helper logic from database persistence and HTTP handlers.

### What was tricky to build
- The main design choice was deciding how much auth behavior to add before the repository layer exists. I stopped at the boundary where the next slice clearly becomes signup/login/logout plus token persistence.

### What warrants a second pair of eyes
- The cookie defaults in `pkg/auth/helpers.go`, especially the `SameSite`, TTL, and `Secure` flag expectations for local development versus deployed environments.

### What should be done in the future
- Add auth repository methods over `users`, `author_sessions`, `password_reset_tokens`, and `email_verification_tokens`.
- Wire the helpers into real signup/login/logout endpoints.

### Code review instructions
- Review `pkg/auth/helpers.go` and `pkg/auth/helpers_test.go` together.
- Validate with:
```text
go test ./cmd/... ./pkg/...
```

### Technical details
- Commands run:
```text
gofmt -w pkg/auth/*.go
go test ./cmd/... ./pkg/...
```

## Step 10: Pivot Author Auth To Keycloak / OIDC

After the initial auth helper slice, the direction changed: Draft Review should not build app-managed signup/login/password-reset flows first. Instead, it should follow the `hair-booking` model and use Keycloak-backed browser OIDC, with the backend owning signed session cookies and exposing the authenticated author through `/api/me`.

I treated this as a real implementation pivot, not just a task-list note. That meant replacing the local password-helper dead-end, adding Glazed auth settings, adding session and OIDC packages, wiring `/auth/*` routes and `/api/me`, updating the ticket language, and validating the dev-mode runtime behavior.

### Prompt Context

**User prompt (verbatim):**
```text
we want to do auth using keycloak, look at /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking
```

**Assistant interpretation:** Switch the auth plan and implementation toward the same Keycloak / OIDC pattern used in `hair-booking`.

**Inferred user intent:** Avoid investing further in local password-auth flows and align Draft Review with an existing, already-proven Keycloak integration pattern from the user’s own Go codebase.

**Commit (code):** pending in this step while the Keycloak auth slice is being finalized

### What I did
- Inspected the `hair-booking` reference files:
- `pkg/auth/config.go`,
- `pkg/auth/session.go`,
- `pkg/auth/oidc.go`,
- `pkg/server/http.go`,
- `pkg/server/handlers_me.go`,
- related tests and serve wiring.
- Replaced the earlier local-password helper direction by:
- deleting `pkg/auth/helpers.go`,
- deleting `pkg/auth/helpers_test.go`,
- adding `pkg/auth/config.go`,
- adding `pkg/auth/session.go`,
- adding `pkg/auth/oidc.go`,
- adding `pkg/auth/session_test.go`,
- adding `pkg/auth/oidc_test.go`.
- Updated `cmd/draft-review/cmds/serve.go` to:
- register a Glazed auth section,
- load auth settings,
- pass auth settings into the server,
- document OIDC examples.
- Updated `cmd/draft-review/main.go` so the serve command help includes the auth section.
- Updated `pkg/server/http.go` to:
- initialize auth settings,
- create the session manager and OIDC authenticator,
- expose `/api/me`,
- expose `/auth/login`, `/auth/callback`, `/auth/logout`, `/auth/logout/callback` in OIDC mode,
- include auth metadata in `/api/info`,
- support `auth-mode=dev` for local work.
- Added `pkg/server/http_test.go` to cover `/api/me` in dev and unauthenticated OIDC modes.
- Updated `README.md`, `tasks.md`, `changelog.md`, and the design doc so the ticket now points at Keycloak / OIDC rather than local signup/reset flows.

### Why
- `hair-booking` already provides a working Glazed + OIDC + signed-cookie model in the same ecosystem.
- Reusing that pattern is materially better than inventing a different auth architecture for Draft Review.

### What worked
- The `hair-booking` structure copied over cleanly with only project-specific env names and help text changes.
- The dev-mode smoke test verified the new runtime surface:
- `/api/info` now reports auth configuration,
- `/api/me` returns a synthetic local author in `auth-mode=dev`.
- The auth package and server tests passed after wiring.

### What didn't work
- The first `go test ./cmd/... ./pkg/...` failed because the OIDC dependencies were not yet in `go.mod`:
```text
pkg/auth/oidc.go:15:2: no required module provides package github.com/go-jose/go-jose/v3
pkg/auth/oidc.go:18:2: no required module provides package golang.org/x/oauth2
```
- I fixed that with `go get github.com/go-jose/go-jose/v3 github.com/go-jose/go-jose/v3/jwt golang.org/x/oauth2` followed by `go mod tidy`.
- My first auth smoke attempt without a DSN failed because the current SQL config layer still insists on a configured connection string:
```text
Error: failed to compute connection string: dsn is empty and host/database/user are not fully configured
```
- I also hit a short startup race right after `docker compose up -d postgres`:
```text
failed to ping postgres: dial tcp 127.0.0.1:5432: connect: connection refused
```
- Retrying the migration command after Postgres finished booting resolved that.

### What I learned
- The current backend already has a clean seam for browser auth: Glazed section on `serve`, auth initialization in `NewHTTPServer`, session decoding in handler helpers, and `/api/me` as the frontend probe.
- The remaining hard part is no longer Keycloak wiring. It is binding authenticated OIDC claims to local author ownership in the article layer.

### What was tricky to build
- The subtle part was not the OIDC code itself; it was making sure the implementation pivot was reflected everywhere. Without updating the ticket/docs, the repo would have compiled with Keycloak auth while the implementation guide still told an intern to build password resets and email verification first.

### What warrants a second pair of eyes
- The `auth-mode=dev` versus `auth-mode=oidc` behavior split in `pkg/server/http.go`.
- The decision to leave article ownership on the temporary development owner for one more slice until OIDC identity is threaded into article mutations.
- The current CLI/database behavior where `serve` still expects explicit SQL connection settings even for auth-only local checks.

### What should be done in the future
- Add authenticated author bootstrap from OIDC `issuer` + `subject` claims into a local author/user record.
- Protect article mutation routes with the authenticated author identity.
- Add local Keycloak startup instructions and, ideally, a small local Keycloak compose/bootstrap flow similar to `hair-booking`.

### Code review instructions
- Start with `pkg/auth/config.go`, `pkg/auth/session.go`, and `pkg/auth/oidc.go`.
- Then review `cmd/draft-review/cmds/serve.go` and `pkg/server/http.go`.
- Validate with:
```text
go test ./cmd/... ./pkg/...
docker compose up -d postgres
go run ./cmd/draft-review migrate up --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --auth-mode dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --listen-host 127.0.0.1 --listen-port 8080
curl -sS http://127.0.0.1:8080/api/info
curl -sS http://127.0.0.1:8080/api/me
docker compose down
```

### Technical details
- Commands run:
```text
rg -n "keycloak|oidc|openid|oauth|session|auth" /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking -S
sed -n '1,260p' /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/auth/config.go
sed -n '1,340p' /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/auth/session.go
sed -n '1,360p' /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/auth/oidc.go
gofmt -w cmd/draft-review/main.go cmd/draft-review/cmds/serve.go pkg/auth/*.go pkg/server/http.go
go test ./cmd/... ./pkg/...
go get github.com/go-jose/go-jose/v3 github.com/go-jose/go-jose/v3/jwt golang.org/x/oauth2
go mod tidy
go test ./cmd/... ./pkg/...
docker compose up -d postgres
go run ./cmd/draft-review migrate up --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --auth-mode dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --listen-host 127.0.0.1 --listen-port 8080
curl -sS http://127.0.0.1:8080/api/info
curl -sS http://127.0.0.1:8080/api/me
curl -sS -o /tmp/draft-review-auth-login.txt -w '%{http_code}' http://127.0.0.1:8080/auth/login
docker compose down
```

## Step 11: Bind Authenticated Authors To Local Ownership And Scope Article Routes

Once the Keycloak/OIDC slice existed, the next real gap was ownership. The backend could tell who the current browser user was, but article reads and writes still did not use that identity. They either returned everything or used the old hardcoded owner path. This step closes that gap by adding local author bootstrap and making author article routes owner-scoped.

The implementation uses the same high-level pattern as `hair-booking`'s authenticated client bootstrap: stable external identity fields in the local table, a service that ensures a local row from `issuer` + `subject`, and handlers that resolve the current browser identity before calling the domain service.

### Prompt Context

**User prompt (verbatim):** (same as Step 10)

**Assistant interpretation:** Continue from the Keycloak auth pivot by making article ownership actually depend on the authenticated author.

**Inferred user intent:** Finish the first useful auth integration step, not just the browser-session plumbing.

**Commit (code):** pending in this step while the ownership slice is being finalized

### What I did
- Added `pkg/db/migrations/0003_user_auth_identity.sql` to:
- add `auth_subject` and `auth_issuer` to `users`,
- add a unique auth-identity index,
- backfill the seeded development user to `dev/local-author`.
- Added `pkg/auth/service.go` and `pkg/auth/postgres.go` to:
- find local users by auth identity,
- create local users for new identities,
- update existing users from current OIDC claims.
- Added `pkg/auth/service_test.go`.
- Changed `pkg/articles/service.go` and `pkg/articles/postgres.go` so author routes take an owner user ID and filter by `owner_user_id` for list/detail/update/create.
- Updated `pkg/server/http.go` to:
- construct an auth service,
- resolve the current browser identity into a local author,
- require an authenticated author for `/api/articles` list/detail/create/update.
- Updated `pkg/server/http_test.go` with a test that proves article list uses the ensured local author ID.
- Updated `pkg/db/seed.go` so `seed dev` upserts the development author with auth identity fields instead of leaving it detached from the local auth model.
- Updated the README, tasks, and design doc to reflect that article routes are now identity-scoped.

### Why
- Without local author bootstrap, Keycloak auth only affected `/api/me`; it did not actually protect or personalize backend data access.
- Owner-scoped article queries are the minimum requirement for a real author backend.

### What worked
- The `hair-booking`-style ensure-user pattern fit Draft Review cleanly.
- The live smoke test verified the important end-to-end behavior:
- dev-mode `/api/me` resolves the local author,
- seeded articles are visible to that author,
- newly created articles are also owned by that author.

### What didn't work
- I caught a subtle lifecycle issue before it became a user-facing bug: `migrate` normally runs before `seed`, so the earlier auth-identity backfill migration alone was not enough for fresh databases. A brand-new DB would still get a seeded user row without auth identity fields unless `seed dev` itself also set them.
- I fixed that by changing `pkg/db/seed.go` to upsert the development user with `auth_subject='local-author'` and `auth_issuer='dev'`.

### What I learned
- The correct unit of ownership here is not the browser session itself. It is the local `users` row derived from the auth claims. Once that exists, the rest of the backend can stay database-native and not depend directly on cookie/session internals.

### What was tricky to build
- The tricky part was making the seeded dev data and the new auth bootstrap agree on the same author. If the dev claims and seeded user diverge even slightly, the backend appears "correct" but the author sees an empty article list. That kind of mismatch is easy to miss unless you validate a real migrate -> seed -> serve flow.

### What warrants a second pair of eyes
- The `CreateAuthenticatedUser` upsert behavior in `pkg/auth/postgres.go`, especially the email-based conflict path and the assumptions it makes about OIDC-provided email stability.
- The decision to scope even article list/detail routes to the authenticated owner, which is correct for author endpoints but changes the earlier unauthenticated development behavior.

### What should be done in the future
- Run the same ownership flow in real `auth-mode=oidc` against a live Keycloak realm.
- Consider whether owner-scoped article routes should move behind a dedicated author sub-router or middleware once more endpoints exist.

### Code review instructions
- Start with `pkg/db/migrations/0003_user_auth_identity.sql`.
- Then review `pkg/auth/service.go`, `pkg/auth/postgres.go`, `pkg/articles/service.go`, `pkg/articles/postgres.go`, and `pkg/server/http.go`.
- Validate with:
```text
go test ./cmd/... ./pkg/...
docker compose up -d postgres
go run ./cmd/draft-review migrate up --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --auth-mode dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --listen-host 127.0.0.1 --listen-port 8080
curl -sS http://127.0.0.1:8080/api/me
curl -sS http://127.0.0.1:8080/api/articles
curl -sS -X POST http://127.0.0.1:8080/api/articles -H 'Content-Type: application/json' -d '{"title":"Owned By Dev Author","author":"Development Author","intro":"Ownership smoke test."}'
docker compose down
```

### Technical details
- Commands run:
```text
gofmt -w pkg/auth/*.go pkg/articles/*.go pkg/server/http.go pkg/server/http_test.go
go test ./cmd/... ./pkg/...
docker compose up -d postgres
go run ./cmd/draft-review migrate up --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review seed dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
go run ./cmd/draft-review serve --auth-mode dev --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable' --listen-host 127.0.0.1 --listen-port 8080
curl -sS http://127.0.0.1:8080/api/me
curl -sS http://127.0.0.1:8080/api/articles
curl -sS -X POST http://127.0.0.1:8080/api/articles -H 'Content-Type: application/json' -d '{"title":"Owned By Dev Author","author":"Development Author","intro":"Ownership smoke test."}'
docker compose down
```

## Step 13: Upgrade The CLI Root To Real Glazed Help Pages

The backend CLI had already adopted Glazed-style flags and root logging, but the command surface was still using plain Cobra help. That meant the user-facing command documentation lagged behind the actual local workflow, and the root command was missing the help-topic discovery experience that the Glazed skills expect. I fixed that by wiring the embedded help system into the root, keeping the existing persistent logger path, and then adding the first set of actual help pages for this application.

This step also flushed out a concrete dependency gap: the repository had never imported `github.com/go-go-golems/glazed/pkg/help`, so the first real `go run ./cmd/draft-review help ...` invocation failed until the new help-related transitive modules were recorded in `go.sum`. That was useful to catch now, because otherwise the help integration would have looked correct in code review but still failed in a fresh environment.

### Prompt Context

**User prompt (verbatim):**
```text
wire up the persistent logger and setup the help entries for the glazed command properly (see skill)
```

**Assistant interpretation:** use the Glazed command and help-authoring patterns, make sure the root-level logger setup is the canonical persistent path, and replace the plain Cobra help surface with an embedded Glazed help system plus real topic pages.

**Inferred user intent:** leave the CLI in the same operational shape as other Glazed-based tools in this ecosystem so a developer can discover the app workflow directly from `draft-review help`.

**Commit (code):** pending during diary capture for this step

### What I did
- Read the `glazed-help-page-authoring` and `glazed-command-authoring` skill instructions again before editing.
- Verified that the root command already had:
- `logging.AddLoggingSectionToRootCommand(rootCmd, "draft-review")`
- `PersistentPreRunE: logging.InitLoggerFromCobra(...)`
- Added `cmd/draft-review/doc/doc.go` with an embedded filesystem and `AddDocToHelpSystem(...)`.
- Added four Glazed help pages:
- `cmd/draft-review/doc/01-overview.md`
- `cmd/draft-review/doc/02-local-development.md`
- `cmd/draft-review/doc/03-auth-modes.md`
- `cmd/draft-review/doc/04-database-workflow.md`
- Updated `cmd/draft-review/main.go` to:
- create `help.NewHelpSystem()`,
- load the embedded docs,
- call `help_cmd.SetupCobraRootCommand(helpSystem, rootCmd)`,
- add slightly fuller `Long` descriptions on the root and group commands.
- Ran `gofmt` on the Go files touched in the CLI root/doc package.
- Ran `go mod tidy` after the first help invocation exposed missing `go.sum` entries for Glazed help dependencies such as `glamour`, `bubbletea`, `frontmatter`, and `go-sqlite3`.
- Validated the resulting CLI with:
- `go run ./cmd/draft-review help`
- `go run ./cmd/draft-review help --topics`
- `go run ./cmd/draft-review help local-development`
- `go test ./cmd/... ./pkg/...`

### Why
- The root command was in an incomplete state: logging was already rooted correctly, but documentation discovery still behaved like a plain Cobra app.
- Draft Review now has enough local workflow complexity that a real help index is materially useful, especially for a new engineer who needs to understand `serve`, `migrate`, `seed`, the frontend dev proxy, and the Keycloak modes.
- The transitive dependency gap needed to be resolved in versioned module metadata, not left implicit in a developer cache.

### What worked
- The existing persistent logger setup was already correct, so no invasive logging refactor was necessary.
- `help_cmd.SetupCobraRootCommand(...)` attached cleanly once the embedded docs package existed.
- The root `help` page now lists the new top-level topics and renders them correctly.
- `go test ./cmd/... ./pkg/...` still passed after the root help upgrade.

### What didn't work
- My first validation command was `go run ./cmd/draft-review help topics`, following the wording in the skill guidance. In this Glazed version that is not a subcommand; it is a flag-based view, so the correct invocation is:
```text
go run ./cmd/draft-review help --topics
```
- The first help invocation also failed because the repository was missing the help system's transitive `go.sum` entries. `go mod tidy` resolved that cleanly.

### What I learned
- The logging half of the user request was already satisfied by the root command; the missing part was that the help surface had not caught up with the Glazed conventions.
- This Glazed release exposes topic listing via `help --topics` and the default root help page, not a dedicated `help topics` subcommand.

### What was tricky to build
- The tricky part was not the embedded filesystem or the frontmatter itself. It was validating the actual behavior of the imported Glazed help package rather than assuming the skill wording matched this exact released version.

### What warrants a second pair of eyes
- The wording and discoverability of the first four help pages. They are functional and accurate, but as the command surface grows it may be worth adding more command-specific examples.
- Whether we want to expose a custom `help topics` alias in the app for consistency with the wording used in internal skills, even though the upstream package already works via `help` and `help --topics`.

### What should be done in the future
- Add command-specific examples as more operational workflows stabilize, especially for Keycloak/OIDC and frontend-integrated end-to-end testing.
- If Draft Review grows more command groups, move the explicit Cobra group construction out of `main.go` into per-group `root.go` files for cleaner composition.

### Code review instructions
- Start with `cmd/draft-review/main.go` and verify the root initialization order:
- logging section first,
- help system load second,
- command registration after that.
- Then read the docs in `cmd/draft-review/doc/` and check that the slugs are unique, top-level sections are intentional, and the pages map to real commands and flags in the repo.
- Validate with:
```text
go run ./cmd/draft-review help
go run ./cmd/draft-review help --topics
go run ./cmd/draft-review help local-development
go test ./cmd/... ./pkg/...
```

### Technical details
- Commands run:
```text
gofmt -w cmd/draft-review/main.go cmd/draft-review/doc/doc.go
go run ./cmd/draft-review help topics
go mod tidy
go run ./cmd/draft-review help
go run ./cmd/draft-review help --topics
go run ./cmd/draft-review help local-development
go test ./cmd/... ./pkg/...
```
- Observed error before the dependency fix:
```text
missing go.sum entry for module providing package github.com/adrg/frontmatter
missing go.sum entry for module providing package github.com/charmbracelet/glamour
missing go.sum entry for module providing package github.com/mattn/go-sqlite3
```

## Step 14: Make Reader Invites Immediately Copyable

The backend already returned a unique per-reader token from the invite API, but the frontend threw that information away and replaced it with a generic "Invitation Sent!" screen. That was enough to prove the mutation worked, but not enough to support the actual workflow you asked for: create an invite and immediately copy a real reader link into email, Slack, or wherever else you talk to beta readers.

I kept this intentionally small. The backend did not need a schema or API change because `POST /api/articles/{id}/invite` already returns the created reader with its `token`. The missing piece was in the invite dialog: await the mutation instead of marking success optimistically, build the real review URL from the returned token, and expose a copy button in the success state.

### Prompt Context

**User prompt (verbatim):**
```text
do the small slice
```

**Assistant interpretation:** implement the smallest practical version of “invite this reader and give me a link I can copy/paste”.

**Inferred user intent:** keep the existing per-reader invite model, but surface it in the UI immediately so invite creation turns into a usable communication step rather than just a database write.

**Commit (code):** pending during diary capture for this step

### What I did
- Updated `frontend/src/author/InviteDialog.tsx` so the dialog:
- awaits the async invite callback,
- only enters the success state after the mutation resolves,
- shows backend errors inline,
- renders the generated per-reader reader URL,
- exposes copy actions for that URL.
- Updated `frontend/src/author/InviteDialog.css` to style:
- the new copyable link row,
- the sent-state action layout,
- the inline error message.
- Updated `frontend/src/app/AuthorApp.tsx` so the real app invite path:
- unwraps the invite mutation result,
- builds a stable backend-origin URL using `backendOrigin` plus `/r/${reader.token}`,
- returns that URL into the dialog.
- Updated `frontend/src/author/InviteDialog.stories.tsx` so Storybook still satisfies the new async contract and demonstrates the success state with a realistic link.
- Validated with `npm run build`.

### Why
- The backend already had the token. Not surfacing it in the UI was wasted value.
- The old optimistic success state could claim "sent" even if the invite mutation failed.
- Copy-paste is the shortest route to real usage and does not require adding email-delivery infrastructure.

### What worked
- No backend change was required because the API contract already returned the reader token.
- The existing `backendOrigin` setting in `AuthorApp` made it straightforward to generate a backend-owned reader URL rather than accidentally using the Vite origin.
- The build passed once the dialog contract and Storybook stub were aligned.

### What didn't work
- The first build failed because the new `onInvite` contract was stricter than the old fire-and-forget callback:
```text
Type '(email: string, note: string) => Promise<{ email: string; inviteUrl: string; } | undefined>' is not assignable...
```
- I fixed that by making the real app throw if no active article is selected and by updating the Storybook story to return a mock async result instead of `void`.

### What I learned
- The frontend had all the data it needed already. The gap was purely in how the success path was modeled in the dialog.
- Using the backend origin explicitly matters here for the same reason it mattered in auth: the shareable link should point at the actual app origin that owns `/r/:token`.

### What was tricky to build
- The main subtlety was keeping the invite contract strict enough to be useful without forcing a larger redesign of the dialog API. Returning an explicit invite result turned out to be the cleanest minimal shape.

### What warrants a second pair of eyes
- Whether we also want a second convenience action like "Copy email + link" in addition to the plain link copy.
- Whether the readers list should surface the same full invite URL directly so authors can re-copy it later without reopening the dialog.

### What should be done in the future
- Consider adding a “copy invite message” template that includes the link and the personal note.
- Consider surfacing the full invite URL in the reader-management table so the dialog is not the only place where it can be retrieved.

### Code review instructions
- Review `frontend/src/author/InviteDialog.tsx` first to check the async success and error states.
- Then review `frontend/src/app/AuthorApp.tsx` to confirm the URL is built from the backend origin and the returned token.
- Validate with:
```text
cd frontend
npm run build
```

### Technical details
- Commands run:
```text
npm run build
```
- Final behavior:
```text
create invite -> wait for mutation -> show /r/<token> URL -> copy to clipboard
```

## Context

This diary belongs to the backend planning ticket for Draft Review. The app is currently a React frontend using MSW and in-memory mock data; this ticket defines the first real backend built on PostgreSQL and local Docker Compose.

## Quick Reference

Primary documents:

- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/index.md`
- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md`
- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/reference/01-backend-ticket-diary.md`

## Usage Examples

Use this diary when:

- reviewing how the ticket was assembled,
- continuing backend implementation in a later session,
- checking which repository files shaped the design,
- confirming which validation and delivery commands were used.

## Step 12: Package The Frontend For Real Deployment And Document The Hosted Contract

The deployment request forced me to re-check the actual runtime path rather than trusting the existing local-development docs. That was the right call: the backend could proxy Vite in development, but the production image still only contained the Go API binary, which meant a Coolify rollout would have shipped a healthy `/healthz` and a broken `/`. I fixed that by moving Draft Review onto the same production pattern used successfully in `hair-booking`: build the frontend, copy the generated assets into a Go web package, embed them, and serve the SPA shell directly from the backend when no dev proxy is configured.

This step also exposed two deployment-only bugs that normal `go test` did not catch. First, the Dockerfile pinned Go `1.25.5` while the module required `1.25.7+`. Second, the new Glazed help integration used an in-memory SQLite store that crashes a `CGO_ENABLED=0` binary on startup. I corrected both so the container actually boots in a hosted shape, then verified the image against a real Postgres instance and the embedded shell paths.

### Prompt Context

**User prompt (verbatim):**
```text
setup keycloak realm and stuff in ~/code/wesen/terraform, and then deploy to coolify (look at /home/manuel/workspaces/2026-03-19/hair-signup/hair-booking for example. Let me know which documentation you find about coolify / keycloak setup
```

followed by:

```text
look at ticket and diaries and other paybooks potentially in these repos
```

**Assistant interpretation:** inspect the existing deployment and identity playbooks first, then make Draft Review actually deployable in the same hosted shape before attempting the Terraform and Coolify rollout.

**Inferred user intent:** do the real hosted setup rather than stopping at local dev instructions, and reuse the established `hair-booking` and infra-repo operator patterns wherever possible.

**Commit (code):** af42fbe799d081e4cba790440872d3fbc352b1d0 — "Embed frontend for production deployment"

### What I did
- Read the existing Draft Review ticket checklist and diary plus the shared infra and `hair-booking` deployment playbooks.
- Confirmed the pre-existing production gap:
- `Dockerfile` only built the Go binary.
- `pkg/server/http.go` returned `404` for browser routes unless `frontend-dev-proxy-url` was set.
- Added `pkg/web` with:
- embedded production assets,
- `go generate` support,
- a generator that builds `frontend/` and copies `frontend/dist` into `pkg/web/public`.
- Updated `pkg/server/http.go` so non-API/non-auth browser routes:
- proxy to Vite in dev when configured,
- otherwise serve embedded assets directly,
- fall back to `index.html` for SPA routes like `/r/<token>`.
- Added HTTP tests for embedded `/`, `/assets/app.js`, and SPA fallback behavior.
- Updated `frontend/vite.config.ts` to proxy `/auth` and `/healthz` in addition to `/api`.
- Reworked the Docker build into a multi-stage frontend-plus-Go image.
- Added `curl` and `tzdata` to the runtime image for hosted health-check/debug parity.
- Fixed the Docker Go toolchain version mismatch by moving from `golang:1.25.5` to `golang:1.25.8`.
- Fixed the Glazed help/runtime crash by switching the production image build from `CGO_ENABLED=0` to `CGO_ENABLED=1`.
- Added SQL env fallbacks in `pkg/config/sql.go` so `DRAFT_REVIEW_DSN` works in container-style deployments.
- Added hosted docs:
- `docs/deployments/draft-review-coolify.md`
- `docs/deployments/draft-review-coolify-playbook.md`
- Validated with:
- `npm run build`
- `go generate ./pkg/web`
- `go test ./cmd/... ./pkg/...`
- `docker build -t draft-review:local .`
- live container smoke with Postgres:
- `curl http://127.0.0.1:18081/healthz`
- `curl -I http://127.0.0.1:18081/`
- `curl http://127.0.0.1:18081/api/info`

### Why
- Coolify deployment only makes sense if the container serves the browser shell itself.
- The existing repo docs had drifted ahead of the actual production packaging.
- Reusing the `hair-booking` embed pattern reduced deployment risk and kept the operator model consistent across repos.

### What worked
- The `hair-booking` frontend-embed pattern transferred cleanly.
- The embedded shell tests passed immediately once the server stopped hard-failing to `404`.
- The final Docker image built successfully and served `/`, `/healthz`, and `/api/info` against a live Postgres database.
- The DSN env fallback made the container runnable with a Coolify-style environment contract instead of requiring CLI flag overrides.

### What didn't work
- The first Docker build failed because the builder image was too old for the module:
```text
go: go.mod requires go >= 1.25.7 (running go 1.25.5; GOTOOLCHAIN=local)
```
- The first container smoke failed before serving HTTP because the Glazed help system tried to initialize SQLite inside a `CGO_ENABLED=0` binary:
```text
failed to create tables: failed to create sections table: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
```
- The first container run using env-only DB config still failed because the SQL section was not actually reading `DRAFT_REVIEW_DSN`:
```text
Error: failed to compute connection string: dsn is empty and host/database/user are not fully configured
```
- The first host-bridge run on Linux failed because `host.docker.internal` was not resolvable until I added:
```text
--add-host=host.docker.internal:host-gateway
```

### What I learned
- The current Glazed embedded-help setup is not compatible with a `CGO_ENABLED=0` production binary.
- The repo had a subtle contract mismatch: docs and local dev workflow implied a deployable full-stack app, but only the API had actually been packaged.
- For this repo, env-driven DSN support is important because the deployment story is container-first rather than shell-wrapper-first.

### What was tricky to build
- The hard part was not embedding the frontend assets themselves. The sharp edges were the deployment-only interactions between Glazed help, CGO, Docker toolchain versions, and container-style configuration injection. Each one looked unrelated at first, but together they would have made the first Coolify rollout fail even after the frontend embed was correct.

### What warrants a second pair of eyes
- Whether the Glazed help system should later move to a non-SQLite backend so the image can return to `CGO_ENABLED=0`.
- Whether the generated frontend assets should remain committed in `pkg/web/public` or be treated as build artifacts only.
- Whether we want a dedicated `make docker-smoke` target so the Postgres-backed container check becomes routine.

### What should be done in the future
- Finalize the hosted Keycloak realm settings and secret flow in `~/code/wesen/terraform`.
- Deploy the app to Coolify and verify the real hosted OIDC browser loop.
- Consider code-splitting the frontend bundle; Vite currently warns about a chunk just over 500 kB after minification.

### Code review instructions
- Start with `pkg/server/http.go` and `pkg/web/generate_build.go`.
- Then review `Dockerfile` and `pkg/config/sql.go` together because they define the real deployment contract.
- Confirm the new operator docs reflect the actual image/runtime shape:
- `docs/deployments/draft-review-coolify.md`
- `docs/deployments/draft-review-coolify-playbook.md`
- Re-run the key validations:
```text
npm run build
go generate ./pkg/web
go test ./cmd/... ./pkg/...
docker build -t draft-review:local .
docker compose up -d postgres
docker run --rm --add-host=host.docker.internal:host-gateway -p 18081:8080 \
  -e 'DRAFT_REVIEW_DSN=postgres://draft_review:draft_review@host.docker.internal:5432/draft_review?sslmode=disable' \
  -e DRAFT_REVIEW_AUTO_MIGRATE=true \
  -e DRAFT_REVIEW_AUTH_MODE=dev \
  draft-review:local
curl -fsS http://127.0.0.1:18081/healthz
curl -fsSI http://127.0.0.1:18081/
curl -fsS http://127.0.0.1:18081/api/info
```

### Technical details
- Local playbooks consulted:
- `/home/manuel/code/wesen/terraform/docs/shared-keycloak-platform-playbook.md`
- `/home/manuel/code/wesen/terraform/ttmp/2026/03/19/TF-001-SHARED-KEYCLOAK-INFRA--centralize-shared-keycloak-terraform-and-coolify-deployment-guidance/playbook/01-shared-keycloak-terraform-and-coolify-deployment-playbook.md`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/docs/deployments/hair-booking-coolify.md`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/docs/deployments/hair-booking-coolify-playbook.md`
- Official docs consulted:
- Coolify environment-variable docs explaining build vs runtime variable injection
- Keycloak hostname docs for production reverse-proxy exposure
- Keycloak OIDC client docs covering valid redirect URIs and web origins

## Related

- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md`
- `ttmp/2026/03/24/DR-001--draft-review-mvp-react-redux-storybook-implementation/design-doc/01-draft-review-mvp-implementation-plan.md`

## 2026-03-25 - Hosted New-Article Follow-Up

### Goal
- Reproduce the user report that `New Article` on the hosted app did nothing and emitted no HTTP request, verify the live deployment state, and either prove a cache issue or ship the missing fix.

### What I changed
- Confirmed the public HTML was already serving the newer frontend bundle:
```text
curl -sS https://draft-review.app.scapegoat.dev/ | sed -n '1,120p'
```
- Confirmed the live JS bundle already contained the `Untitled Article` create path and `/api/articles` mutation wiring:
```text
curl -sS https://draft-review.app.scapegoat.dev/assets/index-muJ7WUsy.js | rg -n 'Untitled Article|/api/articles'
```
- Used a scripted browser check against the hosted app to log in through Keycloak, then trigger `File -> New Article…` and capture requests/responses.
- Verified the hosted click emitted:
```text
POST https://draft-review.app.scapegoat.dev/api/articles
201 POST https://draft-review.app.scapegoat.dev/api/articles
```
- Found the remaining real issue: the empty dashboard state showed only “No articles yet” with no visible create CTA, so first-run users had to know about the menu item.
- Added a primary `New Article` button to the empty dashboard in `frontend/src/author/Dashboard.tsx`.
- Wired the dashboard empty-state CTA to the shared `handleNewArticle` mutation handler in `frontend/src/app/AuthorApp.tsx`.
- Rebuilt the frontend and embedded assets with:
```text
npm run build
go generate ./pkg/web
go test ./cmd/... ./pkg/...
```

### What worked
- The hosted app was already on the corrected bundle; the prior deploy was not stale.
- The menu path worked exactly as intended once exercised with a real browser session.
- The server returned `201` and immediately loaded the new article editor plus reader/reaction queries for the created article.
- The dashboard CTA patch was small and reused the same creation handler, so there was no duplicate article-creation logic.

### What didn't work
- My first browser automation attempt failed because the Keycloak login page exposes both the password input and the “show password” toggle under a loose `/password/i` label match:
```text
locator.fill: Error: strict mode violation: getByLabel(/password/i) resolved to 2 elements
```
- A direct search for a dashboard `New Article` button after login found nothing because there actually was no button in the empty state.

### What I learned
- The user report was directionally right, but the live root cause was split in two: the hidden menu path worked, while the obvious first-run path did not exist.
- For this app, first-run empty states need to expose the primary workflow directly instead of relying on desktop-style menu discovery.

### What was tricky
- Verifying the production issue required separating “the latest deploy is live” from “the visible UX still fails the user goal.” The JS bundle and menu handler were correct, but the empty dashboard still produced a broken first impression.

### Review notes
- Re-run the hosted browser flow after deploy:
```text
1. Log into https://draft-review.app.scapegoat.dev/
2. From the empty dashboard, click New Article ->
3. Confirm POST /api/articles appears in the network tab
4. Confirm the editor opens on Untitled Article
```
