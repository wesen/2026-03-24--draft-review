# Tasks

## Investigation Deliverables

- [x] Inspect the backend article schema, persistence model, and article section API shape
- [x] Inspect the frontend editor, reader renderer, and reaction model for markdown/image assumptions
- [x] Write a detailed intern-facing analysis / design / implementation guide
- [x] Write a compact investigation reference document with file-backed findings
- [x] Validate the ticket bundle with `docmgr doctor`
- [x] Upload the ticket bundle to reMarkable

## Recommended Implementation Queue

## Phase 1 Execution

- [x] Add a DR-012 phase-1 implementation diary and record each implementation slice chronologically
- [x] Add a database migration that backfills canonical markdown if needed and drops `article_sections.body_plaintext`
- [x] Refactor backend article and reader section types to expose `bodyMarkdown` instead of `paragraphs[]`
- [x] Refactor backend article persistence and reader link resolution to read and write only `body_markdown`
- [x] Refactor backend validation logic so section bodies validate markdown text instead of normalized paragraph arrays
- [x] Refactor frontend article and reader types to use `bodyMarkdown`
- [x] Add frontend markdown block helpers that derive reader/review blocks from `bodyMarkdown`
- [x] Refactor the author editor to edit canonical markdown bodies and add minimal image-markdown insertion affordances
- [x] Refactor reader-facing components and author review components to render derived markdown blocks instead of `paragraphs[]`
- [x] Update dashboard, welcome, and other ancillary UI that currently assumes paragraph arrays
- [x] Update stories, mocks, and tests for the new canonical markdown section shape
- [x] Run `go test ./cmd/... ./pkg/...`, `npm run build`, and `docmgr doctor --ticket DR-012 --stale-after 30`
- [x] Update DR-012 docs, changelog, and diary after each completed slice

## Phase 2 Execution

- [x] Extend the DR-012 implementation diary with explicit phase-2 upload work entries
- [x] Add a database migration for `article_assets`
- [x] Add a backend asset package with asset types, repository methods, a local-disk storage adapter, and service-level validation for file type, size, filenames, and markdown snippet generation
- [x] Extend backend configuration with media storage settings suitable for local dev and hosted persistent volumes
- [x] Add `POST /api/articles/{id}/assets` for authenticated author uploads
- [x] Add `GET /media/article-assets/{assetId}/{filename}` for serving persisted uploaded assets
- [x] Add backend tests for upload validation and media serving
- [x] Add frontend upload API wiring for article assets
- [x] Refactor the author editor to upload a chosen local file and insert the returned markdown snippet at the cursor
- [x] Add basic editor UI feedback for upload in-flight and upload errors
- [x] Update README / local-dev docs with media-root configuration and persistence notes
- [x] Run `go test ./cmd/... ./pkg/...`, `npm run build`, and `docmgr doctor --ticket DR-012 --stale-after 30`
- [x] Update DR-012 docs, changelog, and diary after each completed slice

## Later Phases

- [ ] Decide whether image blocks should be reactable, non-reactable, or caption-reactable and align analytics with that choice
- [ ] Add version-scoped assets and richer asset metadata if article versioning needs hard media snapshots
- [ ] Add cloud/object-storage backends if hosted deployment outgrows mounted-volume storage
