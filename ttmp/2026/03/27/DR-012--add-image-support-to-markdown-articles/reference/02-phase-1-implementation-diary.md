---
Title: Phase 1 implementation diary
Ticket: DR-012
Status: active
Topics:
    - backend
    - frontend
    - editor
    - markdown
    - images
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/author/ArticleEditor.tsx
      Note: Shows the blank-line-to-paragraph authoring model being replaced
    - Path: frontend/src/author/Dashboard.tsx
      Note: Shows ancillary UI that still indexes into section.paragraphs
    - Path: pkg/articles/types.go
      Note: Defines the section API shape being migrated away from paragraph arrays
    - Path: pkg/reviewlinks/postgres.go
      Note: Shows reader payload assembly that still loads body_plaintext
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-27T08:16:45.073517611-04:00
WhatFor: Keep a detailed chronological record of the DR-012 phase-1 implementation work that removes body_plaintext and makes markdown canonical end-to-end.
WhenToUse: Use while reviewing or continuing the DR-012 phase-1 implementation.
---


# Phase 1 implementation diary

## Goal

Capture the step-by-step implementation history for DR-012 phase 1: drop
`body_plaintext`, make `body_markdown` canonical, remove the paragraph-array API,
and update the frontend/editor/reader pipeline accordingly.

## Context

This diary accompanies the DR-012 design guide. The design guide explains the
architecture and recommended plan; this diary records what was actually changed,
what broke, what was learned, and how to review each step.

## Diary

## Step 1: Phase-1 Execution Plan And Scope Lock

This step converts DR-012 from a research ticket into an execution ticket. The
important scope decision is that phase 1 will be opinionated: `body_markdown`
becomes the only canonical section body, `body_plaintext` is removed from the
database through a migration, and the API does not preserve backwards compatibility
for `paragraphs[]`.

That choice narrows the implementation in a useful way. Instead of trying to keep
two competing representations alive, the codebase can move cleanly toward one model
and derive review blocks from markdown wherever the UI still needs paragraph-like
anchors.

### Prompt Context

**User prompt (verbatim):** "for phase1 just get rid of body_plaintext entirely and jus tuse the markdown, no need to  maintain backwards compat in the api (but db migration). 
then add tasks for phse1 to the ticket and work on them one by one, committing at appropriate intervals and writing a frequent detailed diary."

**Assistant interpretation:** The user wants DR-012 turned into an implementation ticket, with phase 1 defined as a hard switch to canonical markdown and no paragraph-array compatibility in the API.

**Inferred user intent:** Move from analysis into execution with a clean, non-hedged design, and keep the implementation process documented in detail.

### What I did
- Added explicit phase-1 execution tasks to [tasks.md](../tasks.md).
- Added this diary document so each implementation slice can be recorded as it lands.
- Re-checked the frontend and backend usage sites for `section.paragraphs`, `body_markdown`, and `body_plaintext` to understand the full impact radius before editing code.

### Why
- The original ticket only described a recommended queue; it did not yet describe the concrete execution order requested by the user.
- The diary needs to exist before code changes begin so the implementation history is chronological and not reconstructed later.

### What worked
- The current code search made the migration scope clear: the paragraph-array assumption is spread across backend types, section persistence, reader link resolution, author review UI, dashboard feedback excerpts, mocks, stories, and tests.

### What didn't work
- N/A

### What I learned
- The most dangerous assumption is not the renderer. It is the number of small UI surfaces that index into `section.paragraphs[pIdx]` to generate excerpts and reaction anchors.
- Removing `body_plaintext` is straightforward in the database, but removing `paragraphs[]` from the API requires a coordinated frontend refactor rather than a backend-only slice.

### What was tricky to build
- The tricky part at this stage was scoping, not code. The phrase “just use markdown” sounds smaller than it is because review semantics still depend on stable paragraph-style IDs. The implementation must keep those anchors stable while making markdown canonical.

### What warrants a second pair of eyes
- The eventual choice of markdown-to-block splitting rules, because that logic will determine whether historical reaction anchors remain intuitive.

### What should be done in the future
- Record each code slice after it lands, including exact validation commands and any breakages encountered during the migration.

### Code review instructions
- Start with [tasks.md](../tasks.md) to see the ordered phase-1 slices.
- Review the design guide for the target end state before reviewing the code diffs.

### Technical details
- Relevant code anchors gathered before implementation:
  - [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go)
  - [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go)
  - [pkg/reviewlinks/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go)
  - [frontend/src/author/ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx)
  - [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx)
  - [frontend/src/author/Dashboard.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx)

## Step 2: Backend Becomes Markdown-Canonical

This step completes the backend half of phase 1. The database gains an explicit
migration that removes the duplicated plaintext section body, the article and
reader payloads stop exposing `paragraphs[]`, and both persistence paths now load
and store only `body_markdown`.

This is the structural point of no return for the phase-1 migration. After this
step, any frontend still expecting `section.paragraphs` is supposed to break, which
is acceptable because the user explicitly requested no backwards compatibility in
the API.

### What I changed
- Added [0007_drop_article_section_plaintext.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0007_drop_article_section_plaintext.sql) to backfill `body_markdown` from `body_plaintext` when needed and then drop `body_plaintext`.
- Updated [types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go) so article section types expose `bodyMarkdown`.
- Updated [types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/types.go) so reader article sections expose `bodyMarkdown`.
- Updated [service.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/service.go) so validation normalizes markdown text instead of collapsing sections into paragraph arrays.
- Updated [postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go) so article creation, version cloning, section replacement, and section listing all operate only on `body_markdown`.
- Updated [postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go) so reader link resolution assembles sections from `body_markdown`.
- Updated [seed.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/seed.go) to seed only the canonical markdown column.
- Updated [http_test.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/server/http_test.go) test fixtures to use the new section shape.

### Why this step matters
- The old dual-column design was artificial. The only reason `body_plaintext` still
  existed was to support the paragraph-array API shape, not because the product
  needed two independent section-body representations.
- Image support depends on markdown being the source of truth. If the backend keeps
  pretending the canonical model is plain paragraphs, the editor and renderer can
  never safely converge.

### Validation
- Ran `gofmt -w pkg/articles/types.go pkg/reviewlinks/types.go pkg/articles/service.go pkg/articles/postgres.go pkg/reviewlinks/postgres.go pkg/db/seed.go pkg/server/http_test.go`
- Ran `go test ./cmd/... ./pkg/...`
- Result: all Go tests passed after the backend migration

### What worked
- The backend change set was smaller than the frontend one because the backend was
  already storing markdown in the database. Most of the work was deleting the
  fallback model rather than inventing new behavior.
- The review subsystem is flexible about paragraph keys, which means later frontend
  block derivation can keep the reaction model without another backend schema change.

### What did not work immediately
- I initially had to re-check the reader-link repository carefully because it had a
  separate section-loading path that still scanned `body_plaintext` and split it into
  `paragraphs[]`.

### What I learned
- The real blast radius is the frontend, not the backend. Once the backend exposes
  `bodyMarkdown`, almost every user-facing surface that renders or excerpts article
  content must be updated together.
- The migration is clean because `body_markdown` already exists in production data;
  phase 1 is mostly about removing the duplicate representation and making that
  removal explicit.

### What warrants review
- The migration SQL in [0007_drop_article_section_plaintext.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0007_drop_article_section_plaintext.sql),
  because it is intentionally destructive after the backfill step.
- The section payload shape in [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go)
  and [pkg/reviewlinks/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/types.go),
  because the frontend must now follow that contract exactly.

### Next step
- Refactor the frontend to use `bodyMarkdown` end-to-end, then add a block-derivation
  layer for review anchors and introduce minimal markdown-image authoring affordances.

## Step 3: Frontend Switches To Markdown Bodies And Derived Review Blocks

This step makes the React app compatible with the new backend contract. The
important architectural decision is that the frontend does not try to reconstruct
the old `paragraphs[]` API model. Instead, it keeps markdown canonical in the
section shape and derives reactable reader blocks locally from `bodyMarkdown`.

That keeps the system honest. The API speaks markdown, the editor edits markdown,
the reader renders markdown, and the review UI derives block anchors from markdown
instead of pretending the backend still stores separate paragraph objects.

### What I changed
- Updated [article.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/types/article.ts) so `Section` uses `bodyMarkdown`.
- Added [markdownBlocks.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/lib/markdownBlocks.ts) with shared helpers for:
  - markdown normalization
  - blank-line block derivation
  - markdown-to-plain-text excerpts
  - block-index lookup from `sectionId-pN` reaction IDs
  - rough markdown word counting
- Updated [ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx) so section editing is markdown-native instead of paragraph-native.
- Added simple image-insertion affordances in [ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx) and corresponding layout styles in [ArticleEditor.css](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.css).
- Updated [SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx) so reader reactions attach to derived markdown blocks.
- Updated [ArticleReader.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.tsx) so author review mode renders and excerpts markdown-derived blocks.
- Updated [Dashboard.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx) so recent feedback excerpts resolve against derived block text instead of `section.paragraphs`.
- Updated [WelcomeSplash.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/WelcomeSplash.tsx) to estimate reading length from markdown text rather than paragraph count.
- Updated fixtures and stories in:
  - [db.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/mocks/db.ts)
  - [ReaderPage.stories.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/ReaderPage.stories.tsx)
  - [WelcomeSplash.stories.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/WelcomeSplash.stories.tsx)
  - [SectionNav.stories.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/SectionNav.stories.tsx)

### Why this step matters
- This is the step that actually makes image support real for phase 1. Once the
  editor persists markdown directly and the reader renders markdown-derived blocks,
  ordinary markdown image syntax can flow end-to-end without a second article-body
  representation.
- The shared helper file avoids a subtle class of bugs where the dashboard, reader,
  and author review surfaces would disagree about what block `sectionId-p2` refers
  to.

### Validation
- Ran `npm run build` inside `frontend/`
- Result: TypeScript compilation and Vite production build both succeeded

### What worked
- The existing [Prose.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.tsx)
  component meant the rendering side was already close; the real work was changing
  the data model and block derivation.
- The existing reaction ID format (`sectionId-pN`) survived the migration cleanly
  because the frontend can continue generating those IDs from derived markdown
  blocks.

### What was tricky
- The dashboard and author review surfaces were the easy places to miss. They do
  not look like “editor” code, but they both depended on direct indexing into
  `section.paragraphs`.
- Block derivation has to be shared. If one screen split on different blank-line
  rules than another, the reaction anchors would drift immediately.

### Product note
- The new editor affordance is intentionally minimal. It inserts markdown image
  snippets; it does not yet handle local uploads, asset picking, captions as a
  first-class content model, or image-specific reaction policy.

### Next step
- Run the full phase-1 validation pass (`go test`, `npm run build`, `docmgr doctor`)
  and then update the ticket docs and commit the frontend slice.

## Step 4: Full Validation And Ticket Closure For Phase 1

This final step is the integration checkpoint. The implementation is only credible
 if the backend tests still pass after the API break, the frontend production build
 succeeds against the new section shape, and the ticket itself validates cleanly
 under `docmgr doctor`.

### Validation run
- Ran `go test ./cmd/... ./pkg/...`
- Ran `npm run build` in [frontend/](/home/manuel/code/wesen/2026-03-24--draft-review/frontend)
- Ran `docmgr doctor --ticket DR-012 --stale-after 30`

### Results
- Go tests: passed
- Frontend TypeScript + Vite build: passed
- `docmgr doctor`: passed

### What I verified conceptually
- The backend API and the frontend now agree that a section body is `bodyMarkdown`.
- The editor, reader, and author review mode all derive block anchors from the same
  helper layer.
- The phase-1 ticket tasks now describe a completed migration rather than a plan.

### Remaining limits after phase 1
- Images are currently markdown URLs, not uploaded assets.
- Image blocks inherit the generic reader-block reaction behavior; the product still
  needs a later decision on whether image-only blocks should remain reactable,
  become non-reactable, or use caption-specific reactions.
- The editor affordance inserts markdown snippets, but there is not yet an upload
  picker, media library, or drag-and-drop image workflow.

### Review notes for the next engineer
- Start with [markdownBlocks.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/lib/markdownBlocks.ts)
  if you need to understand reaction anchors.
- Review [0007_drop_article_section_plaintext.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0007_drop_article_section_plaintext.sql)
  before deploying this slice to a real database.
- Treat phase 1 as “markdown URL images are now structurally possible,” not as a
  complete digital asset system.

## Related

- [Markdown article image support analysis and implementation guide](../design-doc/01-markdown-article-image-support-analysis-and-implementation-guide.md)
- [Investigation notes](./01-investigation-notes.md)
