---
Title: "Draft Review CLI Overview"
Slug: "draft-review-overview"
Short: "Understand how the Draft Review CLI fits together for local backend development."
Topics:
- draft-review
- backend
- cli
Commands:
- serve
- migrate
- seed
Flags:
- dsn
- frontend-dev-proxy-url
- auth-mode
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This page covers the shape of the `draft-review` command line application, how its command groups map onto the local backend workflow, and why the root command is responsible for both logging and discoverable help topics.

In practice, the CLI has three jobs. `serve` runs the HTTP backend, `migrate` manages the embedded PostgreSQL schema, and `seed` loads a deterministic development dataset. Those jobs are intentionally separate so local development, smoke testing, and CI can drive the same binary without inventing one-off scripts.

The root command matters because it is where application-wide behavior is configured. The logging flags come from the Glazed logging section added to the root Cobra command, and `PersistentPreRunE` initializes the zerolog-backed logger once for every command invocation. The help system is also attached at the root, which means `draft-review help topics` and `draft-review help <slug>` work consistently across the entire command tree.

The main command groups are:

- `serve`: boot the HTTP API, auth routes, and optional frontend dev proxy.
- `migrate up`: apply embedded SQL migrations to the configured PostgreSQL database.
- `migrate status`: inspect which embedded migrations are already applied.
- `seed dev`: migrate if needed and then load a stable local article, author, and review dataset.

The main configuration domains are:

- SQL connection settings, which define how the CLI reaches PostgreSQL.
- Auth settings, which define `dev` mode or OIDC / Keycloak mode.
- Backend settings, which define local conveniences such as auto-migration and frontend dev proxying.

Typical local usage starts by pointing every command at the same DSN, then picking the auth mode that matches the scenario being tested. Use `auth-mode=dev` for fast backend iteration without Keycloak, and use `auth-mode=oidc` when testing browser login, callback, and logout behavior.

Problem | Cause | Solution
--- | --- | ---
`draft-review help topics` shows only default Cobra help | The embedded help system is not loaded at the root | Make sure `AddDocToHelpSystem(...)` and `help_cmd.SetupCobraRootCommand(...)` are called from `cmd/draft-review/main.go`
Logging flags are missing on subcommands | The logging section was not attached to the root command | Keep `logging.AddLoggingSectionToRootCommand(rootCmd, "draft-review")` on the root before executing
`migrate` or `seed` fails with database configuration errors | The SQL section was not populated | Pass `--dsn ...` or the equivalent SQL connection flags

See Also:

- `draft-review help local-development`
- `draft-review help auth-modes`
- `draft-review help database-workflow`
