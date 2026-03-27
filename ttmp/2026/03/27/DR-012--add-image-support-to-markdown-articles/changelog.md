# Changelog

## 2026-03-27

- Initial workspace created
- Inspected the current article model, section persistence, author editor, markdown renderer, review reactions, and reader APIs to determine what already supports markdown images and where the system is still paragraph-centric
- Added a detailed design and implementation guide for intentional image support in markdown articles, including phased rollout recommendations, API sketches, pseudocode, diagrams, testing, and file references
- Added investigation notes summarizing the current state, key constraints, recommended phases, and relevant code anchors

## 2026-03-27

Added a detailed scoping packet for markdown article image support, covering the current paragraph-centric article model, intentional markdown rendering, optional managed uploads, and phased implementation guidance for interns.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/27/DR-012--add-image-support-to-markdown-articles/design-doc/01-markdown-article-image-support-analysis-and-implementation-guide.md — Primary design deliverable
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/27/DR-012--add-image-support-to-markdown-articles/reference/01-investigation-notes.md — Quick-reference architecture notes


## 2026-03-27

Expanded DR-012 into an explicit phase-1 execution plan and added a detailed implementation diary for the markdown-canonical migration.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/27/DR-012--add-image-support-to-markdown-articles/reference/02-phase-1-implementation-diary.md — Chronological implementation record
- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/27/DR-012--add-image-support-to-markdown-articles/tasks.md — Defines the concrete phase-1 execution slices


## 2026-03-27

Completed the backend half of phase 1 by making markdown canonical in the API and persistence layer, dropping the duplicated plaintext section body through a migration, and updating Go fixtures and validation accordingly.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0007_drop_article_section_plaintext.sql — Backfills markdown if needed and drops `body_plaintext`
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go — Article API sections now expose `bodyMarkdown`
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/types.go — Reader API sections now expose `bodyMarkdown`
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go — Article persistence now reads and writes only `body_markdown`
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go — Reader link resolution now loads markdown section bodies
- /home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/service.go — Section validation now normalizes markdown text


## 2026-03-27

Completed the main frontend half of phase 1 by making `bodyMarkdown` the canonical section shape in React, adding shared markdown-block derivation helpers, updating the editor to author markdown directly, and rewiring reader/review/dashboard surfaces to render and excerpt markdown-derived blocks.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/types/article.ts — Section type now uses `bodyMarkdown`
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/lib/markdownBlocks.ts — Shared block derivation, excerpts, and word-count helpers
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx — Markdown-native section editing and image snippet insertion
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx — Reader blocks now derive from markdown
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.tsx — Author review view now resolves reactions against derived markdown blocks
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx — Recent feedback excerpts now come from derived block text
