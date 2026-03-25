---
Title: Draft Review Backend Architecture and Implementation Guide
Ticket: DR-002
Status: active
Topics:
    - backend
    - golang
    - glazed
    - postgresql
    - docker-compose
    - api
    - auth
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go
      Note: Reference Glazed serve command construction and runtime boot sequence
    - Path: ../../../../../../../../../workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go
      Note: Reference Cobra plus Glazed root command wiring
    - Path: ../../../../../../../../../workspaces/2026-03-19/hair-signup/hair-booking/pkg/config/backend.go
      Note: Reference backend settings section pattern for app-local config
    - Path: ../../../../../../../../../workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go
      Note: Reference embedded SQL migration runner
    - Path: ../../../../../../../../../workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/postgres.go
      Note: Reference pgxpool database bootstrap
    - Path: ../../../../../../../../../workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go
      Note: Reference HTTP server dependency wiring and route registration layout
    - Path: ../../../../../../../corporate-headquarters/clay/pkg/sql/config.go
      Note: Reference implementation for Clay-style Glazed DB config decoding and connection logic
    - Path: ../../../../../../../corporate-headquarters/clay/pkg/sql/flags/sql-connection.yaml
      Note: Reference SQL connection flag schema to mirror in Draft Review
    - Path: ../../../../../../../corporate-headquarters/clay/pkg/sql/settings.go
      Note: Reference implementation for Clay-style embedded SQL connection sections
    - Path: frontend/src/App.tsx
      Note: Defines the author and reader route split the Go backend must support
    - Path: frontend/src/api/articleApi.ts
      Note: Defines the current frontend article and reaction client contracts
    - Path: frontend/src/api/baseApi.ts
      Note: Defines the existing /api prefix that the Go backend should preserve
    - Path: frontend/src/api/readerApi.ts
      Note: Defines token-based reader link resolution on the frontend
    - Path: frontend/src/author/Dashboard.tsx
      Note: Shows current browser-side analytics and draft-killer logic
    - Path: frontend/src/mocks/db.ts
      Note: Defines the current mock article
    - Path: frontend/src/mocks/handlers.ts
      Note: Defines the mock endpoints the Go backend should replace incrementally
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: Shows current client-side review progress and reaction behavior
    - Path: imports/draft-review-screen-spec.md
      Note: Defines the broader backend requirements including auth
ExternalSources: []
Summary: Evidence-backed Go backend design for Draft Review using Glazed command verbs, Clay-style database configuration sections, PostgreSQL, docker compose, and a phased implementation plan grounded in the current frontend and the hair-booking runtime patterns.
LastUpdated: 2026-03-24T21:10:00-04:00
WhatFor: Designing and implementing the first Go backend for Draft Review.
WhenToUse: Use this guide when scaffolding the Draft Review Go service, Glazed CLI verbs, PostgreSQL schema, docker compose stack, or frontend-to-backend integration.
---


# Draft Review Backend Architecture and Implementation Guide

## Executive Summary

Draft Review should not use a TypeScript backend. It should use Go, with Glazed as the CLI/verb framework and PostgreSQL as the persistent store. The Glazed layer should own command verbs such as `serve`, `migrate`, and `seed`, and the database configuration surface should follow the Clay pattern rather than inventing an ad hoc flag parser. That recommendation is grounded in two reference codebases you pointed me at:

1. `clay` already defines a reusable SQL connection/configuration layer with Glazed sections and YAML-backed flags in `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go`, `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go`, and `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/flags/sql-connection.yaml`.
2. `hair-booking` already shows a working Glazed + Cobra Go application structure for a web backend in `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go`, `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go`, `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/config/backend.go`, `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/postgres.go`, and `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go`.

The current Draft Review repo is still frontend-only. The React app already defines the product behavior, routes, and mock API surface, but the real server, database, migrations, auth, and analytics layers do not exist. The backend described here turns that mock app into a real Go service with:

- a Glazed-powered CLI entrypoint,
- a PostgreSQL-backed persistence model,
- a Docker Compose local stack,
- incremental replacement of the current MSW endpoints,
- a package layout that a new intern can navigate without guessing.

## Problem Statement And Scope

Draft Review already has clear product behavior in the frontend:

- authors access the author app on `/*`,
- readers access tokenized review links on `/r/:token`,
- the frontend talks to `/api/*`,
- current data is stored only in frontend mocks,
- the screen spec already expands the product into auth, sharing, article editing, versioning, reader sessions, analytics, exports, and management views.

That behavior is visible in:

- `frontend/src/App.tsx:1-13`
- `frontend/src/api/baseApi.ts:1-7`
- `frontend/src/api/articleApi.ts:4-48`
- `frontend/src/api/readerApi.ts:9-17`
- `frontend/src/mocks/handlers.ts:4-77`
- `frontend/src/mocks/db.ts:5-152`
- `imports/draft-review-screen-spec.md:440-1633`

The backend work in scope for this ticket is:

1. A Go server runtime.
2. A Glazed/Cobra CLI root and backend verbs.
3. PostgreSQL persistence.
4. Docker Compose local development.
5. Author auth and article management foundations.
6. Reader invite links, review sessions, progress, reactions, and summaries.
7. Author analytics and reader-management queries.

This ticket does not require first-phase implementation of:

1. Team billing.
2. Full organization admin.
3. Background workers beyond what is minimally useful.
4. Rich document import pipelines for every format on day one.

## Current-State Analysis

### Current frontend runtime

The frontend currently defines a stable route split:

- `frontend/src/App.tsx:7-10` routes `/r/:token` to the reader app and everything else to the author app.
- `frontend/src/api/baseApi.ts:3-7` hardcodes the `/api` prefix.

The currently implemented mock endpoint surface is:

- `GET /api/articles`
- `GET /api/articles/:id`
- `GET /api/articles/:id/readers`
- `GET /api/articles/:id/reactions`
- `POST /api/articles/:id/reactions`
- `POST /api/articles/:id/invite`
- `GET /api/r/:token`

Those are defined in `frontend/src/mocks/handlers.ts:4-77` and consumed in `frontend/src/api/articleApi.ts:4-48` and `frontend/src/api/readerApi.ts:9-17`.

### Current mock data model

The current repo only persists data in memory and only models three top-level concepts:

1. articles,
2. readers,
3. reactions.

That is visible in `frontend/src/mocks/db.ts:5-152` and typed in:

- `frontend/src/types/article.ts:1-26`
- `frontend/src/types/reader.ts:1-17`
- `frontend/src/types/reaction.ts:1-24`

Observed consequence: the current mock model is enough for the dashboard demo, but it is not enough for durable auth, versioning, review sessions, access-control modes, or analytics history.

### Current frontend business logic that the backend must absorb

The current author dashboard computes average progress and "draft-killer" heuristics in the browser (`frontend/src/author/Dashboard.tsx:28-60`). The current reader page manages section progress and reaction creation entirely in local React state (`frontend/src/reader/ReaderPage.tsx:27-125`). These are useful MVP prototypes, but they reveal exactly which responsibilities belong on the server:

1. durable section progress,
2. durable reaction IDs and timestamps,
3. article-level aggregate analytics,
4. reader-session lifecycle.

### Product requirements visible in the screen spec

The screen spec defines several backend-relevant requirements that are not implemented today:

1. author authentication and session continuity,
2. article manager, article creation, section editing, article settings,
3. share modes: unique invite links, open link, password-protected link,
4. per-article reaction configuration,
5. reader landing with optional reader name,
6. post-read summary with recommendability and notification opt-in,
7. analytics, exports, version history, reader management.

Those appear in `imports/draft-review-screen-spec.md:141-372`, `imports/draft-review-screen-spec.md:440-873`, `imports/draft-review-screen-spec.md:964-1433`, and `imports/draft-review-screen-spec.md:1437-1633`.

## Why Go + Glazed

### Why Go

Go is the right backend language here for pragmatic reasons:

1. The desired runtime is a web service plus CLI verbs plus database migrations. Go handles that cleanly with one compiled binary.
2. The application’s backend shape is mostly HTTP handlers, SQL access, and small service layers. That maps well to idiomatic Go packages.
3. The local developer story becomes simpler: one backend binary, one Postgres container, one frontend dev server if desired.

### Why Glazed

Glazed is a strong fit because you explicitly want verbs and a DB configuration layer. The relevant patterns already exist:

1. `hair-booking` uses a Cobra root command plus a Glazed-built `serve` command in `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go:14-41`.
2. `hair-booking` constructs the `serve` verb as a Glazed command description with sections in `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go:39-99`.
3. `clay` embeds reusable SQL connection sections from YAML in `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:14-48`.
4. `clay` decodes parsed Glazed section values into a `DatabaseConfig` struct in `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go:19-33` and `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go:251-260`.

This means Draft Review does not need to invent either:

- a CLI architecture,
- or a database flag/configuration model.

It should reuse these patterns directly.

## Reference Architecture Patterns

### Pattern 1: Cobra root + Glazed verbs

The `hair-booking` root command is simple and correct:

```text
root cobra command
  -> logging section on root
  -> BuildCobraCommandFromCommand(serveCmd)
  -> root.Execute()
```

Evidence:

- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go:14-41`

Draft Review should use the same shape:

```text
cmd/draft-review/main.go
  rootCmd
    persistent logging init
    register Glazed verbs:
      serve
      migrate
      seed
      doctor
```

### Pattern 2: Command-local sections

`hair-booking` composes its `serve` command from:

1. a default server section,
2. an auth section,
3. a backend section,
4. command settings.

Evidence:

- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go:39-92`

Draft Review should follow the same composition style, but replace the ad hoc `database-url` field with a Clay-style SQL connection section.

### Pattern 3: Clay-style DB configuration sections

Clay already provides the exact pattern you asked for:

1. `SqlConnectionSettings` and `DatabaseConfig` carry Glazed-tagged DB fields.
2. `NewSqlConnectionParameterLayer()` loads the section from embedded YAML.
3. `OpenDatabaseFromSqlConnectionLayer()` decodes the parsed sections into a config and opens the database.

Evidence:

- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:17-48`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:79-108`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go:19-33`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go:164-249`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/flags/sql-connection.yaml:1-56`

Draft Review should either:

1. reuse Clay’s SQL section directly if dependency boundaries make that acceptable, or
2. copy the same pattern into a Draft Review package with only the Postgres-relevant fields preserved.

### Pattern 4: Embedded migrations in Go

`hair-booking` uses embedded SQL migration files with a tiny migration runner:

- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go:14-89`

That approach is appropriate here too. It keeps:

- local bootstrap simple,
- deploy artifacts self-contained,
- migration execution visible and understandable to an intern.

## Proposed Runtime Architecture

### System diagram

```text
React frontend
  |  BrowserRouter, RTK Query
  v
/api/*
  |
  v
draft-review Go server
  |
  +--> auth package
  +--> articles package
  +--> reviewlinks package
  +--> reviews package
  +--> analytics package
  +--> readers package
  |
  v
PostgreSQL
```

### Runtime boundaries

The runtime should separate into four layers:

1. CLI bootstrap and verbs.
2. HTTP server and route wiring.
3. domain services and repositories.
4. PostgreSQL schema and migrations.

The most important rule for maintainability is:

- Glazed commands should configure and launch behavior.
- HTTP handlers should stay thin.
- service packages should own business logic.
- repository packages should own SQL.

## Proposed Repository Layout

```text
cmd/draft-review/main.go
cmd/draft-review/cmds/serve.go
cmd/draft-review/cmds/migrate_up.go
cmd/draft-review/cmds/migrate_status.go
cmd/draft-review/cmds/seed_dev.go

pkg/config/auth.go
pkg/config/backend.go
pkg/config/sql.go

pkg/db/postgres.go
pkg/db/migrations.go
pkg/db/migrations/
  0001_init.sql
  0002_seed_reaction_types.sql

pkg/auth/
pkg/articles/
pkg/reviewlinks/
pkg/reviews/
pkg/readers/
pkg/analytics/
pkg/server/http.go
pkg/server/handlers_auth.go
pkg/server/handlers_articles.go
pkg/server/handlers_reader.go
pkg/server/handlers_analytics.go
pkg/web/

docker-compose.yml
Dockerfile
```

### Intern orientation guide

If a new intern opens this repo, they should be able to answer these questions quickly:

1. How do I run the service locally?
   Look in `cmd/draft-review/main.go`, `cmd/draft-review/cmds/serve.go`, and `docker-compose.yml`.
2. How is database config passed in?
   Look in `pkg/config/sql.go` and the Clay-style embedded SQL section.
3. Where do migrations live?
   Look in `pkg/db/migrations.go` and `pkg/db/migrations/*.sql`.
4. Where do article routes live?
   Look in `pkg/server/handlers_articles.go`.
5. Where is article business logic?
   Look in `pkg/articles/service.go`.
6. Where is the SQL for article queries?
   Look in `pkg/articles/postgres.go`.

That organization matters more than being "clever".

## Glazed Verb Design

### Recommended verbs

Draft Review should have at least these verbs:

1. `draft-review serve`
2. `draft-review migrate up`
3. `draft-review migrate status`
4. `draft-review seed dev`

Future useful verbs:

1. `draft-review export article-feedback`
2. `draft-review doctor`
3. `draft-review create-dev-invite`

### Root command shape

Recommended root behavior:

```go
rootCmd := &cobra.Command{
    Use:     "draft-review",
    Short:   "Draft Review backend and tooling",
    Version: version,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        return logging.InitLoggerFromCobra(cmd)
    },
}

if err := logging.AddLoggingSectionToRootCommand(rootCmd, "draft-review"); err != nil {
    return err
}
```

This is the same root pattern used in `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go:14-41`.

### `serve` command section composition

Recommended sections:

1. default server settings,
2. auth settings,
3. Clay-style SQL connection section,
4. optional DBT section only if you really need it,
5. backend settings,
6. command settings.

Recommended `serve` pseudocode:

```text
func (c *ServeCommand) Run(ctx, parsedValues):
    decode default server settings
    decode auth settings
    decode sql connection settings
    decode backend settings

    db = open postgres connection from sql connection section
    if backendSettings.autoMigrate:
        db.Migrate(ctx)

    server = server.NewHTTPServer(
        authSettings,
        db,
        backendSettings,
    )

    start server
    wait for shutdown
```

### `migrate up` command

Purpose:

- open Postgres from the same SQL section,
- apply embedded migrations,
- exit cleanly without starting the HTTP server.

This is valuable because it decouples schema operations from runtime startup, even if `serve --auto-migrate` still exists for local convenience.

## Database Configuration Layer

### Why use the Clay pattern

You specifically asked for the Glazed DB configuration layer to follow Clay section definitions. That is the right call because Clay already solved:

1. Glazed tags for DB config fields.
2. YAML-defined flag surface.
3. decoding parsed section values into a Go config object.
4. DSN normalization and connection-opening behavior.

Evidence:

- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go:19-33`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go:150-249`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:37-48`

### Recommended Draft Review adaptation

Do not expose every multi-database feature Clay has if Draft Review only needs Postgres. Keep the Clay pattern, but narrow the domain surface to what the app actually needs.

Recommended fields:

```text
host
port
database
user
password
schema
db-type
dsn
driver
ssl-disable
```

Postgres defaults should differ from Clay’s more generic defaults:

- port default: `5432`
- db-type default: `pgx`
- driver default: `pgx`

### Recommended file layout

```text
pkg/config/sql.go
pkg/config/flags/sql-connection.yaml
```

### Example section YAML

```yaml
slug: sql-connection
name: SQL connection flags
Description: |
  PostgreSQL connection settings for the Draft Review application database.
flags:
  - name: host
    type: string
    help: Database host
    default: ""
  - name: port
    type: int
    help: Database port
    default: 5432
  - name: database
    type: string
    help: Database name
    default: ""
  - name: user
    type: string
    help: Database user
    default: ""
  - name: password
    type: string
    help: Database password
    default: ""
  - name: db-type
    type: string
    help: Database type
    default: pgx
  - name: dsn
    type: string
    help: Full PostgreSQL DSN
    default: ""
  - name: ssl-disable
    type: bool
    help: Disable SSL/TLS for local Postgres
    default: true
```

### Recommended decoding flow

```go
sqlSection, err := draftconfig.NewSqlConnectionParameterLayer()
// add to command sections

db, err := draftconfig.OpenDatabaseFromSqlConnectionLayer(
    ctx,
    parsedValues,
    draftconfig.SqlConnectionSlug,
)
```

If you copy Clay helpers directly, keep the same decode logic shape from `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:88-108`.

## Backend Package Design

### `pkg/server`

Purpose: HTTP routing and dependency wiring.

Use the `hair-booking` pattern from `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go:102-240`:

- `NewHTTPServer(...)` wires dependencies,
- `NewHandler(...)` registers routes,
- handlers delegate to services.

Draft Review should do the same with a smaller route set first.

### `pkg/articles`

Purpose:

- article CRUD,
- version creation,
- section editing,
- article settings,
- serializer logic for the frontend DTO shape.

Recommended files:

```text
pkg/articles/service.go
pkg/articles/postgres.go
pkg/articles/types.go
pkg/articles/serialize.go
```

### `pkg/reviewlinks`

Purpose:

- resolve `/r/:token`,
- enforce `invite_link`, `link`, and `password` modes,
- rotate share tokens,
- manage invite tokens.

### `pkg/reviews`

Purpose:

- create review sessions,
- persist section progress,
- persist reactions,
- persist reader summary.

### `pkg/analytics`

Purpose:

- drop-off queries,
- reaction heatmap queries,
- draft-killer heuristics,
- version comparison queries.

### `pkg/auth`

Purpose:

- Keycloak / OIDC configuration,
- signed browser session cookies,
- `/auth/login`, `/auth/callback`, `/auth/logout`, and `/auth/logout/callback`,
- `/api/me` identity projection for the frontend,
- a `dev` fallback mode for local work without Keycloak.

`hair-booking` already demonstrates this pattern, so Draft Review should follow the same split instead of building local password-reset or email-verification flows first.

## PostgreSQL Schema Design

### Core tables

The relational model from the earlier version of this document still holds. The language/runtime changed, but the domain constraints did not.

#### `users`

Purpose: local author records and article ownership projections derived from OIDC identities.

#### `author_sessions`

Purpose: optional future durable session or revocation records if self-contained cookies become insufficient.

#### `password_reset_tokens`

Purpose: not required for the Keycloak-first path. Keep only if Draft Review later adds app-managed credentials.

#### `email_verification_tokens`

Purpose: not required for the Keycloak-first path. Keep only if Draft Review later adds app-managed credentials.

#### `articles`

Purpose: stable article identity, ownership, and sharing/settings state.

Key fields:

- `id uuid primary key`
- `owner_user_id uuid not null`
- `title text not null`
- `author_display_name text not null`
- `status text not null`
- `access_mode text not null`
- `share_token text unique null`
- `share_password_hash text null`
- `show_author_note boolean not null`
- `reader_can_see_reactions boolean not null`
- `reader_can_see_names boolean not null`
- `require_note boolean not null`
- `allow_anonymous boolean not null`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

#### `article_versions`

Purpose: immutable snapshots for version comparisons and reader continuity.

#### `article_sections`

Purpose: ordered sections per article version.

Recommendation: store section body in one text column and derive paragraph arrays for the current frontend response shape.

#### `article_reaction_types`

Purpose: per-article reaction configuration, including custom types later.

#### `reader_invites`

Purpose: invitations sent by the author.

#### `review_sessions`

Purpose: actual reading sessions created after a link is opened.

#### `review_section_progress`

Purpose: section reach/completion tracking per review session.

#### `reactions`

Purpose: paragraph-level feedback.

#### `review_summaries`

Purpose: final overall thoughts, recommendability, and notify-new-version opt-in.

### Entity relationship diagram

```text
users
  └──< articles
         ├──< article_versions
         │      └──< article_sections
         ├──< article_reaction_types
         ├──< reader_invites
         └──< review_sessions
                ├──< review_section_progress
                ├──< reactions
                └─── review_summaries
```

### Why this schema still fits the Go architecture

This schema remains correct even after switching to Go because it reflects domain relationships rather than framework choices. What changes in Go is implementation style:

- migrations live as embedded SQL files,
- repositories use `pgxpool.Pool`,
- services operate on explicit Go structs,
- handlers serialize domain objects into frontend DTOs.

## API Reference

Keep `/api` as the prefix because the current frontend requires it (`frontend/src/api/baseApi.ts:3-7`).

### Author/auth endpoints

1. `GET /api/me`
2. `GET /auth/login`
3. `GET /auth/callback`
4. `GET /auth/logout`
5. `GET /auth/logout/callback`

### Author/article endpoints

1. `GET /api/articles`
2. `POST /api/articles`
3. `GET /api/articles/{id}`
4. `PATCH /api/articles/{id}`
5. `GET /api/articles/{id}/readers`
6. `GET /api/articles/{id}/reactions`
7. `POST /api/articles/{id}/invites`
8. `POST /api/articles/{id}/share-token/reset`
9. `GET /api/articles/{id}/analytics`
10. `GET /api/articles/{id}/feedback`

### Reader/public endpoints

1. `GET /api/r/{token}`
2. `POST /api/r/{token}/start`
3. `POST /api/reviews/{sessionId}/progress`
4. `POST /api/reviews/{sessionId}/reactions`
5. `POST /api/reviews/{sessionId}/summary`

### Error envelope

Use the same stable JSON error-envelope style used in `hair-booking`:

```json
{
  "error": {
    "code": "ARTICLE_NOT_FOUND",
    "message": "The requested article does not exist."
  }
}
```

Reference style:

- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go:93-100`

## Docker Compose Design

### Recommended local stack

```yaml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: draft_review
      POSTGRES_USER: draft_review
      POSTGRES_PASSWORD: draft_review
    ports:
      - "5432:5432"
    volumes:
      - draft_review_postgres:/var/lib/postgresql/data

  backend:
    build:
      context: .
      dockerfile: Dockerfile
    command:
      - draft-review
      - serve
      - --dsn
      - postgres://draft_review:draft_review@postgres:5432/draft_review?sslmode=disable
      - --auto-migrate
      - --frontend-dev-proxy-url
      - http://host.docker.internal:5173
    depends_on:
      - postgres
    ports:
      - "8080:8080"

volumes:
  draft_review_postgres:
```

### Why this compose shape

It mirrors the practical runtime from `hair-booking`:

- one Go backend process,
- one Postgres instance,
- optional frontend dev proxy in local development.

The relevant server-side proxy pattern exists in `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go:272-280` and surrounding code.

## Detailed File-Level Implementation Guide

### `cmd/draft-review/main.go`

Responsibilities:

1. create root Cobra command,
2. add Glazed logging section,
3. register verbs.

Start from the structure in `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go:14-41`.

### `cmd/draft-review/cmds/serve.go`

Responsibilities:

1. define the `serve` command description,
2. compose the Glazed sections,
3. decode section values,
4. open the database,
5. run migrations if configured,
6. start the HTTP server.

Start from `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go:39-181`.

### `pkg/config/sql.go`

Responsibilities:

1. define or import the Clay-style SQL connection section,
2. expose helper functions to decode config and open a DB pool.

Start from:

- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:17-48`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go:88-108`

### `pkg/config/backend.go`

Responsibilities:

1. app-local non-SQL settings such as uploads, public base URL, auto-migrate, frontend proxy URL.

Start from `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/config/backend.go:13-116`.

### `pkg/db/postgres.go`

Responsibilities:

1. open `pgxpool.Pool`,
2. ping the DB,
3. wrap pool access,
4. call migrations.

Start from `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/postgres.go:11-55`.

### `pkg/db/migrations.go`

Responsibilities:

1. embed SQL files,
2. maintain `schema_migrations`,
3. apply unapplied files transactionally.

Start from `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go:14-89`.

### `pkg/server/http.go`

Responsibilities:

1. build the route mux,
2. initialize dependencies,
3. keep route registration readable,
4. keep handlers thin.

Start from `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go:102-240`.

## Key Flow Pseudocode

### Author creates an article

```text
POST /api/articles
  -> authenticate author
  -> validate payload
  -> insert articles row
  -> insert article_versions row
  -> insert article_sections rows
  -> insert default article_reaction_types rows
  -> return article DTO
```

### Reader starts reading from a token

```text
GET /api/r/{token}
  -> resolve invite token or shared token
  -> load current article version
  -> return landing payload

POST /api/r/{token}/start
  -> resolve token
  -> enforce access mode
  -> create review_sessions row
  -> return session id + reader article payload
```

### Reader submits a reaction

```text
POST /api/reviews/{sessionId}/reactions
  -> load review session
  -> validate reaction type is enabled
  -> validate note if article requires note
  -> insert reactions row
  -> update last_active_at
  -> maybe upsert section progress
  -> return persisted reaction DTO
```

### Analytics query

```text
GET /api/articles/{id}/analytics
  -> verify article ownership
  -> query per-section reach counts
  -> query per-section reaction counts by type
  -> compute draft-killer heuristic
  -> query version comparison metrics
  -> return analytics DTO
```

## Implementation Phases

### Phase 0: bootstrap Go runtime

Tasks:

1. add `go.mod` and command root,
2. add `serve` verb,
3. add Clay-style SQL config section,
4. add `docker-compose.yml`,
5. add `/healthz`.

### Phase 1: replace the current mock reader/article paths

Tasks:

1. implement `GET /api/articles`,
2. implement `GET /api/articles/{id}`,
3. implement `GET /api/articles/{id}/readers`,
4. implement `GET /api/articles/{id}/reactions`,
5. implement `POST /api/articles/{id}/invites`,
6. implement `GET /api/r/{token}`,
7. implement `POST /api/r/{token}/start`,
8. implement `POST /api/reviews/{sessionId}/reactions`.

### Phase 2: auth and settings

Tasks:

1. add Keycloak / OIDC browser auth,
2. expose `/api/me`,
3. map authenticated OIDC users onto local author ownership,
4. article settings,
5. access mode enforcement,
6. reaction configuration.

### Phase 3: versions, editor support, summaries

Tasks:

1. section editing endpoints,
2. version creation,
3. review summary submission,
4. notification opt-in persistence.

### Phase 4: analytics and management

Tasks:

1. drop-off queries,
2. heatmap queries,
3. draft-killer insights,
4. reader management,
5. feedback export.

## Testing Strategy

### Unit tests

Focus on:

1. token resolution,
2. access-mode enforcement,
3. draft-killer insight logic,
4. summary validation,
5. SQL config decoding.

### Integration tests

Run against Postgres and cover:

1. migrations,
2. article creation,
3. invite creation,
4. review-session creation,
5. reaction submission,
6. analytics retrieval.

### HTTP tests

Use the style already present in `hair-booking`:

- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http_test.go`

That is the right place to verify route envelopes and auth behavior.

## Design Decisions

### Decision 1: use Go + net/http style handlers

Reason: `hair-booking` already shows a clean, maintainable pattern for this exact class of application.

### Decision 2: use Glazed for command verbs

Reason: you explicitly asked for it, and the repo references give a stronger existing pattern than inventing plain Cobra-only commands.

### Decision 3: use Clay-style SQL sections

Reason: this is the cleanest existing configuration model for DB flags and DSN handling in your codebase.

### Decision 4: use `pgxpool` plus embedded SQL migrations

Reason: it is simple, production-ready, and already demonstrated in `hair-booking`.

### Decision 5: keep the frontend `/api` contract stable

Reason: the React app already expects it, so incremental backend adoption is easier.

## Alternatives Considered

### Alternative: plain Cobra flags without Glazed sections

Rejected because the command surface would drift and the DB configuration layer would be weaker than the Clay pattern you explicitly requested.

### Alternative: reuse Clay’s entire SQL package unchanged

Possible, but only if dependency weight is acceptable. The safer plan may be to copy the section-definition pattern and narrow it to Draft Review’s Postgres-only needs.

### Alternative: use an ORM-heavy Go stack

Rejected for the first recommendation. This application’s strongest complexity is workflow and analytics, not rich object graphs. `pgx` plus explicit repository SQL is easier for an intern to understand and debug.

## Risks And Open Questions

### Risks

1. Importing Word/Google Docs/Markdown can balloon scope if treated as part of phase 1.
2. Password-protected links and invite links must be modeled carefully to avoid authorization bugs.
3. Analytics queries can become expensive if written carelessly.

### Open questions

1. Should Draft Review own full author auth itself, or should external identity eventually be introduced?
2. Do open-link articles create one review session per browser, per entered reader name, or per explicit "start" event only?
3. Should custom reaction types ship in the first implementation or wait until after the default four are stable?

## Suggested First Build Order For An Intern

1. Create `cmd/draft-review/main.go`.
2. Add `serve` using the `hair-booking` command structure.
3. Add `pkg/config/sql.go` using the Clay-style SQL section pattern.
4. Add `pkg/db/postgres.go` and `pkg/db/migrations.go`.
5. Add `docker-compose.yml`.
6. Add `GET /healthz`.
7. Add the first article/list/read routes.
8. Add reader link resolution and reaction submission.
9. Only then add analytics and auth expansions.

That order keeps the feedback loop short and keeps the project implementation-first instead of architecture-first.

## References

- `frontend/src/App.tsx`
- `frontend/src/api/baseApi.ts`
- `frontend/src/api/articleApi.ts`
- `frontend/src/api/readerApi.ts`
- `frontend/src/mocks/db.ts`
- `frontend/src/mocks/handlers.ts`
- `frontend/src/author/Dashboard.tsx`
- `frontend/src/reader/ReaderPage.tsx`
- `imports/draft-review-screen-spec.md`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/config.go`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/settings.go`
- `/home/manuel/code/wesen/corporate-headquarters/clay/pkg/sql/flags/sql-connection.yaml`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/main.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/cmd/hair-booking/cmds/serve.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/config/backend.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/postgres.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/db/migrations.go`
- `/home/manuel/workspaces/2026-03-19/hair-signup/hair-booking/pkg/server/http.go`
