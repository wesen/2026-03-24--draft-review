---
Title: Markdown article image support analysis and implementation guide
Ticket: DR-012
Status: active
Topics:
    - backend
    - frontend
    - editor
    - markdown
    - images
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/author/ArticleEditor.tsx
      Note: Shows the blank-line paragraph editor model
    - Path: frontend/src/primitives/Prose.tsx
      Note: Shows the existing react-markdown rendering path
    - Path: frontend/src/reader/SectionView.tsx
      Note: Shows paragraph-card rendering and paragraphId generation
    - Path: pkg/articles/postgres.go
      Note: Shows how section markdown is currently flattened into paragraph arrays and duplicated into plaintext columns
    - Path: pkg/articles/types.go
      Note: Defines the paragraph-centric section API contract that image support must evolve
    - Path: pkg/db/migrations/0001_init.sql
      Note: Defines the article_sections and reactions schema constraints relevant to image support
    - Path: pkg/reviewlinks/postgres.go
      Note: Shows how reader article payloads still load body_plaintext and split it into paragraphs
    - Path: pkg/reviews/postgres.go
      Note: Shows the paragraph-key-based reaction persistence model
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-27T08:01:01.431030569-04:00
WhatFor: Explain how Draft Review articles work today and provide a detailed plan for adding intentional image support without breaking editing, rendering, sharing, or reactions.
WhenToUse: Use when planning or implementing image support for article sections, especially if the engineer is new to the codebase.
---


# Markdown article image support analysis and implementation guide

## Executive Summary

Draft Review already stores article section bodies in database columns named
`body_markdown` and renders reader-facing content with `react-markdown`. That means
markdown image syntax such as `![Alt text](https://example.com/image.png)` is not a
foreign concept to the codebase. However, the current product surface is still built
around a much narrower abstraction: each section is exposed through the API as
`paragraphs: string[]`, the author editor splits textarea input on blank lines, the
reader maps each array entry to a reactable paragraph card, and analytics/reactions
use paragraph-style keys such as `sectionId-p0`.

Because of that mismatch, adding images is not just a renderer tweak. The work
touches at least five layers:

- persistence and API shape,
- author editing UX,
- reader rendering,
- review anchor identity,
- optional asset upload and storage.

The recommended plan is phased:

1. Make markdown the explicit source of truth for each section.
2. Productize URL-based markdown images first.
3. Add managed uploads only if the product truly needs local file attachments.
4. Keep the review anchor system compatible at first, then clean it up later if richer block types become common.

## Problem Statement

The product request is simple on the surface: authors should be able to include
images inside markdown articles. The difficulty is that the current system does not
treat section bodies as a first-class markdown document in its public contracts.

Current constraints:

- The backend API returns `Section{ title, paragraphs[] }`, not canonical markdown.
  See [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go#L5).
- The database stores both `body_markdown` and `body_plaintext`, but the current
  code writes the same joined paragraph string into both columns. See
  [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L521) and
  [pkg/db/migrations/0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql#L77).
- The author editor converts textarea input into `paragraphs[]` by splitting on
  blank lines. See
  [frontend/src/author/ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx#L35).
- The reader renders each paragraph through `react-markdown`, so markdown syntax
  already works at the paragraph level. See
  [frontend/src/primitives/Prose.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.tsx#L14).
- Reactions and progress are keyed per rendered paragraph via `paragraphId`,
  persisted as `paragraph_key`, and used throughout reader analytics. See
  [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx#L41),
  [pkg/db/migrations/0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql#L147),
  and
  [pkg/reviews/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/postgres.go#L367).

So the real problem is:

`How do we add intentional image support without violating the current editor,
reader, and reaction assumptions?`

This ticket scopes that problem and recommends an implementation order that preserves
today’s product behavior while making room for richer content.

## Current-State Architecture

### 1. Persistence model

The section persistence layer is already suggestive of richer content.

Database facts:

- `article_sections` stores both `body_markdown` and `body_plaintext`.
- The table also stores `estimated_read_seconds`.
- Reactions refer to a section row plus a `paragraph_key`.

Relevant schema:

- [pkg/db/migrations/0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql#L77)

Important observation:

Although the database schema looks markdown-aware, the repository currently loads
sections from `body_plaintext` only, not `body_markdown`. See
[pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L489)
and
[pkg/reviewlinks/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go#L323).

That means the application is not actually using the “markdown vs plaintext”
distinction right now.

### 2. Backend API shape

The backend article model is paragraph-centric.

Relevant types:

- [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go#L5)

Key points:

- `Section` contains `Paragraphs []string`
- `SectionInput` also contains `Paragraphs []string`
- `UpdateArticleInput` accepts `Sections *[]SectionInput`

In practice this means the client cannot send or receive an explicit markdown body
for a section. The API contract itself encourages a “paragraph array” mental model.

### 3. Backend section serialization

When the backend saves a section update, it joins paragraphs with blank lines and
writes the same string into `body_markdown` and `body_plaintext`.

Relevant code:

- [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L536)

Pseudocode of the current save path:

```text
for each section input:
    body = join(section.paragraphs, "\n\n")
    insert article_sections(
        body_markdown = body,
        body_plaintext = body,
        estimated_read_seconds = estimateReadSeconds(body)
    )
```

This is serviceable for plain prose, but it means:

- plaintext is not really plaintext,
- markdown is not exposed directly,
- image-aware metadata is impossible,
- estimated read time counts image markdown text as words.

### 4. Author editor

The author editor is a textarea-per-section editor.

Relevant file:

- [frontend/src/author/ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx#L13)

Behavior:

- each section has a title,
- each section body is shown in one textarea,
- the textarea draft is split into paragraphs on blank lines,
- save writes `paragraphs[]` back into the article object.

Pseudocode:

```text
textarea text
    -> split on blank lines
    -> trim each paragraph
    -> store as section.paragraphs[]
```

This is important because image support must either:

- work inside this existing textarea model, or
- replace it with a more explicit markdown editing experience.

### 5. Reader rendering

The reader is closer to supporting images than the editor is.

Relevant files:

- [frontend/src/primitives/Prose.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.tsx#L14)
- [frontend/src/primitives/Prose.css](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.css#L136)
- [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx#L41)
- [frontend/src/reader/Paragraph.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx#L17)

Facts:

- `Prose` uses `react-markdown` plus `remark-gfm`.
- CSS already defines `.dr-prose img { max-width: 100%; border: 2px solid ... }`.
- Each paragraph string is rendered independently through `Prose`.

This means remote-image markdown syntax is already technically renderable if it
survives the paragraph splitting intact.

### 6. Review anchor model

The reaction system is still named and shaped around paragraphs.

Evidence:

- [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx#L41)
- [pkg/reviews/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/types.go#L56)
- [pkg/reviews/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/postgres.go#L365)

Flow:

```text
section.paragraphs[pi]
    -> paragraphId = section.id + "-p" + pi
    -> reaction POST carries paragraphId
    -> backend stores paragraph_key = paragraphId
```

That works for prose, but it raises product questions for images:

- Should an image itself be reactable?
- Should only the caption be reactable?
- Should an image-only block be non-reactable?
- If markdown block splitting changes, how stable are historical anchors?

### 7. Reader metrics and heuristics

Several reader behaviors assume text paragraphs:

- `WelcomeSplash` estimates word count from paragraph count. See
  [frontend/src/reader/WelcomeSplash.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/WelcomeSplash.tsx#L12).
- `Paragraph.tsx` suppresses chips for headings, code blocks, and short fragments.
  An image-only markdown block will likely fall into “short fragment”. See
  [frontend/src/reader/Paragraph.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx#L29).

This is not fatal, but it means image support is not only about rendering. It also
touches UX rules and analytics assumptions.

## What Already Works Today

This section is important because it prevents over-engineering.

Today, the following already works accidentally:

- a section paragraph can contain markdown,
- the reader will render that markdown,
- `react-markdown` plus `remark-gfm` can render standard markdown images,
- reader CSS already gives images `max-width: 100%`.

So if an author somehow enters:

```markdown
![System overview](https://example.com/system.png)
```

and that string reaches a paragraph intact, the reader can render it.

However, that is not the same as “the app supports article images” because all of
the following remain unresolved:

- there is no image insertion UX,
- there is no upload path,
- the API still hides canonical markdown,
- the review model still assumes paragraphs,
- metrics and heuristics still assume prose.

## Proposed Solution

### Recommendation

Use a phased approach with one primary design decision:

`Make section markdown canonical before adding managed uploads.`

That gives the system a clean foundation and avoids building asset upload flows on
top of the wrong abstraction.

### Phase 1: Explicit markdown-image support without uploads

Goal:

- support image markdown intentionally,
- do not require storage infrastructure yet,
- minimize blast radius.

Recommended changes:

- add `bodyMarkdown` to the section API payload,
- keep `paragraphs[]` temporarily for compatibility if needed,
- update the editor to treat the section body as markdown, not as a prose-only
  paragraph list,
- allow authors to paste image markdown syntax directly,
- update reader rendering logic so image-containing blocks are handled intentionally.

#### Why this is the right first slice

- lowest infrastructure cost,
- aligns with what the renderer already supports,
- clarifies whether the product even needs managed uploads,
- avoids premature object storage work.

### Phase 2: Managed uploads

Goal:

- let authors choose local image files,
- upload them through the backend,
- return a markdown snippet or resolved URL for the editor.

Recommended shape:

- add an `article_assets` table,
- add a storage abstraction,
- use local disk for development,
- use persistent object storage or a mounted persistent volume for hosted.

Do not use the container filesystem directly in production without persistence.
Coolify deployments replace containers; image files written inside the runtime
container would be lost on redeploy unless backed by a mounted volume.

### Phase 3: Block/anchor cleanup

Goal:

- rename and generalize paragraph-only anchor semantics,
- support stable anchors for text, image blocks, and future rich content.

This phase is optional for the first image rollout, but it is the correct long-term
direction if the content model becomes richer.

## Design Decisions

### Decision 1: Canonical section body should be markdown, not `paragraphs[]`

Rationale:

- the database already stores markdown,
- the renderer already consumes markdown,
- `paragraphs[]` is a derived presentation detail, not the canonical document.

Consequence:

- the API and frontend types should expose markdown explicitly.

### Decision 2: Start with URL-based images before managed uploads

Rationale:

- much smaller first slice,
- validates product need,
- avoids blocking on storage infrastructure,
- makes the authoring and reader changes testable immediately.

### Decision 3: Preserve reaction compatibility in the first rollout

Rationale:

- reactions already work and are already in production,
- renaming `paragraph_key` across analytics, reviews, and frontend types increases
  risk,
- first rollout should preserve existing semantics whenever possible.

Practical compromise:

- internally continue using the existing paragraph/block key field,
- but document that it now refers to a rendered review block, not strictly prose.

### Decision 4: Separate image rendering from upload management

Rationale:

- rendering and uploads are distinct problems,
- remote-image markdown support is valuable even without uploads,
- upload design depends on hosting/storage choices that do not exist yet in this repo.

## Proposed Data Model

### Current model

```text
Article
  -> sections[]
       -> title
       -> paragraphs[]
```

### Recommended phase-1 model

```text
Article
  -> sections[]
       -> title
       -> bodyMarkdown
       -> blocks[]          # optional derived view for reader/editor convenience
```

Where:

- `bodyMarkdown` is canonical,
- `blocks[]` is derived for rendering/reactions,
- legacy `paragraphs[]` can temporarily remain for compatibility during migration.

### Optional phase-2 asset model

```text
article_assets
  id
  article_id
  article_version_id?      # optional if assets are version-scoped
  storage_key
  public_url
  content_type
  width
  height
  alt_text
  created_at
```

## API Reference and Sketches

### Current article update API

Current shape:

```json
{
  "title": "Article Title",
  "sections": [
    {
      "id": "uuid-or-temp-id",
      "title": "Introduction",
      "paragraphs": ["One paragraph", "Another paragraph"]
    }
  ]
}
```

Evidence:

- [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go#L11)
- [frontend/src/api/articleApi.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/api/articleApi.ts#L18)

### Recommended phase-1 section payload

Non-breaking transitional shape:

```json
{
  "id": "section-uuid",
  "title": "Introduction",
  "bodyMarkdown": "Intro text\n\n![Architecture](https://example.com/arch.png)",
  "paragraphs": [
    "Intro text",
    "![Architecture](https://example.com/arch.png)"
  ]
}
```

The important rule is:

- `bodyMarkdown` is canonical,
- `paragraphs` is derived and eventually removable.

### Recommended phase-2 upload API

Example:

`POST /api/articles/{id}/assets/images`

Request:

- multipart form upload, or
- direct-to-storage presign flow if object storage is chosen later.

Example response:

```json
{
  "asset": {
    "id": "asset-uuid",
    "url": "https://draft-review.app.scapegoat.dev/media/articles/a1/assets/asset-uuid/system.png",
    "contentType": "image/png",
    "width": 1600,
    "height": 900,
    "altText": ""
  },
  "markdown": "![System diagram](https://draft-review.app.scapegoat.dev/media/articles/a1/assets/asset-uuid/system.png)"
}
```

## Detailed Implementation Plan

This section is written as a sequence an intern could follow.

### Step 1: Expose canonical markdown in backend types

Update:

- [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go#L5)
- [pkg/reviewlinks/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/types.go#L50)

Recommended change:

- add `BodyMarkdown string` to `Section`,
- add `BodyMarkdown string` to `SectionInput`,
- decide whether to keep `Paragraphs []string` temporarily.

Pseudocode:

```text
type Section struct {
    ID
    Title
    BodyMarkdown
    Paragraphs   # optional transitional field
}
```

### Step 2: Read from `body_markdown`, not `body_plaintext`

Update:

- [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L489)
- [pkg/reviewlinks/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go#L323)

Current problem:

- both article owner views and reader views load `body_plaintext`,
- that throws away the semantic distinction the schema already has.

Recommended behavior:

- load `body_markdown`,
- derive legacy paragraph/block arrays from that markdown if needed.

### Step 3: Centralize markdown-to-review-block derivation

Create one helper, probably first in Go, later optionally mirrored in TypeScript.

Goal:

- split canonical markdown into reviewable blocks,
- preserve image-only blocks,
- preserve headings/code blocks/list blocks,
- generate stable block keys like `b0`, `b1`, `b2`.

Pseudocode:

```text
func deriveReviewBlocks(bodyMarkdown string) []ReviewBlock {
    rawBlocks = splitOnBlankLines(bodyMarkdown)
    for i, raw in rawBlocks:
        kind = classify(raw)   # paragraph, image, heading, code, list, quote, etc.
        blocks.append({
            key: "b" + i,
            kind: kind,
            markdown: raw,
            reactable: isReactable(kind, raw),
        })
    return blocks
}
```

Important note:

For the first implementation, a blank-line block splitter is probably enough.
Do not jump straight to a full markdown AST parser unless product requirements force
it. The current system is already block-oriented by blank lines.

### Step 4: Update the frontend section types

Update:

- [frontend/src/types/article.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/types/article.ts#L1)

Recommended phase-1 shape:

```ts
export interface Section {
  id: string;
  title: string;
  bodyMarkdown: string;
  paragraphs?: string[]; // transitional
}
```

### Step 5: Update the author editor to edit markdown explicitly

Update:

- [frontend/src/author/ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx#L35)
- [frontend/src/author/ArticleEditor.css](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.css#L95)

Recommendations:

- rename the internal `contentDraft` mental model to `markdownDraft`,
- stop presenting blank-line splitting as the primary user concept,
- add helper text for image markdown syntax,
- optionally add small buttons:
  - `Insert Image Markdown`
  - `Insert Image Link`

Simple insert example:

```text
![Describe the image](https://example.com/image.png)
```

If upload support exists later, the upload flow can simply insert the returned
markdown snippet at the cursor position.

### Step 6: Update reader rendering from paragraph mapping to review-block mapping

Update:

- [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx#L41)
- [frontend/src/reader/Paragraph.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx#L29)

Recommended direction:

- rename the internal concept from “paragraph” to “review block” in the reader layer,
- still reuse the current UI card if the block is prose-like,
- render image-only blocks with intentional spacing and optional caption styling.

Key product decision:

- If image blocks are not reactable, the UI should hide chips clearly and
  intentionally.
- If they are reactable, the chips should attach to the image block itself, not a
  synthetic caption unless that is a deliberate product rule.

### Step 7: Revisit reaction heuristics

Current heuristics in [Paragraph.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx#L31)
disable reactions on short fragments, headings, and code blocks.

That means a markdown image block will probably become:

- rendered correctly,
- but not reactable.

Recommended first-pass rule:

- image-only blocks should be non-reactable unless product explicitly wants image
  feedback,
- prose paragraphs remain reactable,
- mixed text+image blocks should be discouraged or normalized into separate blocks.

### Step 8: Fix metrics that assume text paragraphs

Update:

- [frontend/src/reader/WelcomeSplash.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/WelcomeSplash.tsx#L12)
- [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L582)

Current issues:

- word count is estimated as `paragraphCount * 40`,
- backend read time counts markdown syntax characters as words.

Recommended behavior:

- calculate metrics from visible text only,
- exclude image-only blocks from word count,
- optionally assign a small fixed read cost to image-only blocks if UX wants that.

### Step 9: Add managed uploads only after phase 1 is stable

There is no existing upload/storage pipeline in this repo. Search evidence:

- [frontend/package.json](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/package.json#L1)
- repo-wide search found no image upload or storage service code.

Recommended backend shape:

```text
POST /api/articles/{id}/assets/images
  validate auth + ownership
  validate content type and size
  persist bytes via storage adapter
  create article_assets row
  return public URL + markdown snippet
```

Recommended storage abstraction:

```text
type ImageStorage interface {
    SaveArticleImage(ctx, articleID, filename, bytes, contentType) -> SavedAsset
    DeleteArticleImage(ctx, storageKey)
}
```

Development:

- local filesystem under a configured path

Hosted:

- persistent volume or object storage

### Step 10: Testing

Backend tests:

- section serialization reads `body_markdown`
- markdown-to-block derivation preserves image-only blocks
- article update accepts `bodyMarkdown`
- upload endpoint validates ownership and file type

Frontend tests:

- editor preserves image markdown syntax during save/load
- reader renders image markdown correctly
- reaction chips behave correctly on prose vs image-only blocks
- preview mode handles sections containing images

Manual smoke test script:

```text
1. Create article
2. Add section body with prose + image markdown
3. Save
4. Re-open article and confirm markdown round-trips unchanged
5. Open reader preview
6. Confirm image renders
7. Confirm reaction behavior matches product rule
8. Generate share link and test reader flow
```

## Design Decisions

See the dedicated section above. The key ones are:

- make markdown canonical,
- phase URL-based support before uploads,
- preserve reaction compatibility first,
- separate rendering from storage concerns.

## Alternatives Considered

### Alternative A: Do nothing except rely on current accidental support

Pros:

- almost no code change

Cons:

- no intentional authoring UX,
- no canonical markdown API,
- no upload story,
- fragile reaction semantics,
- hard to explain to new contributors.

Verdict:

- not recommended

### Alternative B: Full rich-content AST immediately

Pros:

- strongest long-term model
- clean support for many block kinds

Cons:

- too large for the first image feature
- invasive API and frontend rewrite
- unnecessary if the product only needs markdown images soon

Verdict:

- useful later, too large as the first move

### Alternative C: Add uploads first without changing the markdown model

Pros:

- flashy user-facing win quickly

Cons:

- upload API would be built on top of the wrong canonical section abstraction,
- markdown/editor round-tripping remains unclear,
- reactions/metrics still have hidden assumptions.

Verdict:

- wrong ordering

## Implementation Plan

### Suggested delivery order

1. Backend: expose canonical section markdown.
2. Frontend types: consume canonical markdown.
3. Editor: edit markdown and insert image syntax.
4. Reader: derive blocks intentionally and render image blocks cleanly.
5. Reactions/metrics: update heuristics and copy.
6. Optional uploads: add asset table + storage + upload route.

### Sequence diagram: phase-1 save and render

```text
Author textarea
    -> bodyMarkdown
        -> PATCH /api/articles/{id}
            -> store body_markdown
            -> derive compatibility blocks/paragraphs if needed
Reader fetch
    -> section.bodyMarkdown
        -> derive review blocks
            -> render prose/image blocks
            -> attach block ids for reactions
```

### Sequence diagram: phase-2 upload

```text
Author picks image file
    -> frontend POST /api/articles/{id}/assets/images
        -> backend validates auth + ownership
        -> storage adapter saves bytes
        -> backend records article_assets row
        -> response returns public URL + markdown snippet
    -> editor inserts ![alt](url) at cursor
    -> author saves article markdown
```

## Risks and Edge Cases

### Risk 1: Historical reactions tied to paragraph numbering

If block splitting changes too aggressively, old `paragraph_key` values may stop
matching the UI.

Mitigation:

- keep current blank-line block ordering in the first rollout,
- do not rewrite historical reaction keys in phase 1.

### Risk 2: Broken images in shared reader links

If the system allows arbitrary external URLs, images can disappear or become slow.

Mitigation:

- accept that tradeoff in phase 1,
- add managed uploads in phase 2 if product wants stronger control.

### Risk 3: Hosting/storage mismatch

If uploads are added without persistence, images will disappear on redeploy.

Mitigation:

- require persistent storage before shipping managed uploads to hosted.

### Risk 4: Read-time metrics become nonsensical

Image-heavy articles will produce poor estimates if paragraph count or markdown
syntax is used as a text proxy.

Mitigation:

- rework the metrics functions during the same rollout.

## Open Questions

- Are external image URLs sufficient for the first release, or is local upload a
  hard requirement?
- Should image blocks be reactable?
- If image blocks are reactable, should the reaction attach to the image itself or
  to a surrounding/caption block?
- Are uploaded assets article-scoped or article-version-scoped?
- Should intro text also support markdown images, or only section bodies?

## References

Core code references:

- [pkg/articles/types.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/types.go#L5)
- [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L183)
- [pkg/articles/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/articles/postgres.go#L489)
- [pkg/db/migrations/0001_init.sql](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/db/migrations/0001_init.sql#L77)
- [pkg/reviewlinks/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviewlinks/postgres.go#L323)
- [pkg/reviews/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/postgres.go#L164)
- [pkg/reviews/postgres.go](/home/manuel/code/wesen/2026-03-24--draft-review/pkg/reviews/postgres.go#L288)
- [frontend/src/author/ArticleEditor.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx#L35)
- [frontend/src/types/article.ts](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/types/article.ts#L1)
- [frontend/src/primitives/Prose.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.tsx#L14)
- [frontend/src/primitives/Prose.css](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.css#L136)
- [frontend/src/reader/SectionView.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx#L41)
- [frontend/src/reader/Paragraph.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx#L29)
- [frontend/src/reader/WelcomeSplash.tsx](/home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/WelcomeSplash.tsx#L12)

## Open Questions

<!-- List any unresolved questions or concerns -->

## References

<!-- Link to related documents, RFCs, or external resources -->
