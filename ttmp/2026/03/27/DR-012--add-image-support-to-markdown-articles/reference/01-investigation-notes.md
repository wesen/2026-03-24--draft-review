---
Title: Investigation notes
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
    - Path: frontend/package.json
      Note: Shows existing markdown rendering dependencies and absence of upload/storage dependencies
    - Path: frontend/src/primitives/Prose.css
      Note: Confirms image styling already exists in the markdown renderer
    - Path: frontend/src/reader/Paragraph.tsx
      Note: Shows the current heuristics that would likely make image-only blocks non-reactable
    - Path: frontend/src/reader/WelcomeSplash.tsx
      Note: Shows paragraph-count-based reading metrics that image support should revisit
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-27T08:01:01.434074605-04:00
WhatFor: Provide a quick reference of the current article/image-related architecture and the recommended rollout for image support.
WhenToUse: Use when you need a fast summary of where image support touches the codebase without reading the full design guide.
---


# Investigation notes

## Goal

Capture the current article, markdown, rendering, and review constraints that matter
for adding images to Draft Review articles.

## Context

This document is the short companion to the main design guide. It focuses on
copy/paste-friendly findings, code anchors, and recommended rollout phases.

## Quick Reference

### Main finding

Draft Review already has partial image support accidentally because:

- section bodies are stored in `body_markdown`,
- the reader uses `react-markdown`,
- image CSS already exists.

But it does not have intentional product support because:

- the API exposes `paragraphs[]`, not canonical markdown,
- the editor is built around blank-line paragraph splitting,
- reactions and analytics are keyed per paragraph,
- there is no upload/storage path.

### Current architecture by layer

| Layer | Current shape | Why it matters |
| --- | --- | --- |
| DB schema | `article_sections.body_markdown`, `body_plaintext` | The schema is richer than the API |
| Backend types | `Section{ title, paragraphs[] }` | Canonical markdown is hidden |
| Editor | textarea split on blank lines | Authoring is prose-first |
| Reader render | `react-markdown` + `remark-gfm` | Images can already render |
| Review anchors | `paragraphId`, `paragraph_key` | Image support collides with paragraph-only semantics |

### Key code anchors

- Backend section types:
  - [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go#L5)
- Backend section persistence:
  - [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L521)
- Section schema:
  - [pkg/db/migrations/0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql#L77)
- Reader section load path:
  - [pkg/reviewlinks/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go#L323)
- Author editor:
  - [frontend/src/author/ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx#L35)
- Markdown render primitive:
  - [frontend/src/primitives/Prose.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.tsx#L14)
- Image styling:
  - [frontend/src/primitives/Prose.css](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.css#L136)
- Reaction anchor generation:
  - [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx#L41)
- Reaction persistence:
  - [pkg/reviews/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/postgres.go#L367)

### Recommended rollout

#### Phase 1

- expose `bodyMarkdown`
- treat markdown as canonical
- support URL-based markdown images
- keep reaction compatibility

#### Phase 2

- add `article_assets`
- add upload API
- add persistent storage adapter

#### Phase 3

- rename/generalize paragraph anchors into review block anchors

### Suggested API sketches

Transitional section:

```json
{
  "id": "section-id",
  "title": "Introduction",
  "bodyMarkdown": "Text\n\n![Alt](https://example.com/image.png)",
  "paragraphs": ["Text", "![Alt](https://example.com/image.png)"]
}
```

Managed upload response:

```json
{
  "asset": {
    "id": "asset-id",
    "url": "https://.../media/articles/a1/assets/asset-id/file.png"
  },
  "markdown": "![Alt text](https://.../media/articles/a1/assets/asset-id/file.png)"
}
```

## Usage Examples

Use this document when:

- deciding whether image support is “already there” or only partial,
- planning the first implementation slice,
- explaining to another engineer why uploads and rendering are separate problems,
- finding the exact code paths that still assume paragraphs.

## Related

- [Markdown article image support analysis and implementation guide](../design-doc/01-markdown-article-image-support-analysis-and-implementation-guide.md)
