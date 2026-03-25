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
LastUpdated: 2026-03-24T20:04:39.227367839-04:00
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

## Related

- `ttmp/2026/03/24/DR-002--draft-review-backend-postgresql-docker-compose-foundation/design-doc/01-draft-review-backend-architecture-and-implementation-guide.md`
- `ttmp/2026/03/24/DR-001--draft-review-mvp-react-redux-storybook-implementation/design-doc/01-draft-review-mvp-implementation-plan.md`
