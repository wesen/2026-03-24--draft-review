---
Title: "Database Workflow"
Slug: "database-workflow"
Short: "Use migrate and seed commands correctly against the embedded PostgreSQL schema."
Topics:
- database
- migrations
- seed
- postgresql
Commands:
- migrate
- seed
Flags:
- dsn
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Application
---

This page covers the database lifecycle commands, how they fit together during normal development, and why `seed dev` should be treated as a migration-aware bootstrap rather than a pure fixture inserter.

The database workflow is intentionally simple. Migrations are embedded in the Go binary, so `draft-review migrate up` is the canonical way to bring a PostgreSQL database to the expected schema. `draft-review migrate status` reports which embedded migrations have already been applied. `draft-review seed dev` assumes a local development scenario and loads a deterministic dataset after ensuring the schema is current.

Use `migrate up` when you want schema control without loading fixture data. Use `seed dev` when you want a ready-to-test local author, article, and reader state. Because `seed dev` already applies migrations, it should not be run concurrently with a separate `migrate up` against the same fresh database.

Typical command sequence:

```bash
go run ./cmd/draft-review migrate up \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'

go run ./cmd/draft-review migrate status \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'

go run ./cmd/draft-review seed dev \
  --dsn 'postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable'
```

During development, the most common failure is not a bad migration file. It is usually one of these: PostgreSQL is not reachable, the DSN points to the wrong port, or `seed dev` is being treated as if it were safe to race against another migration process.

Problem | Cause | Solution
--- | --- | ---
`migrate up` fails to connect | PostgreSQL is not listening on the DSN host or port | Start the database and verify the DSN explicitly
`seed dev` fails with schema-related database errors | Another process is migrating the same database concurrently | Run `seed dev` by itself on a fresh database, or wait for migrations to complete
`migrate status` shows pending rows after seeding | The wrong database was seeded or the DSN changed between commands | Re-run both commands with the exact same DSN and verify the target database name

See Also:

- `draft-review help local-development`
- `draft-review help auth-modes`
- `draft-review help draft-review-overview`
