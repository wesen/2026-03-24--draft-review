---
Title: Implementation Diary
Ticket: DR-005
Status: active
Topics:
    - implementation
    - frontend
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: "Added readOnly + onBackToEditor props for preview mode"
    - Path: frontend/src/reader/Paragraph.tsx
      Note: "Added readOnly prop — chips at reduced opacity, non-interactive"
    - Path: frontend/src/reader/SectionView.tsx
      Note: "Threads readOnly through to Paragraph"
    - Path: frontend/src/author/ArticleSettings.tsx
      Note: "Added title/intro/version editing + delete article with confirmation"
    - Path: frontend/src/author/ArticleEditor.tsx
      Note: "Unsaved changes detection, auto-grow textarea, removed Danger Zone"
    - Path: frontend/src/api/articleApi.ts
      Note: "Added deleteArticle mutation (DELETE /api/articles/:id)"
    - Path: frontend/src/author/Dashboard.tsx
      Note: "Renamed Share -> Settings button"
    - Path: frontend/src/author/ArticleEditor.css
      Note: "Dirty indicator style, improved reorder arrows, removed danger zone CSS"
ExternalSources: []
Summary: "Step-by-step diary of DR-005 implementation: preview reactions fix, article management, editor improvements."
LastUpdated: 2026-03-26T16:00:00-04:00
WhatFor: "Document what was done, why, and how for future reference and code review."
WhenToUse: "When reviewing DR-005 changes or understanding implementation decisions."
---

# DR-005 Implementation Diary

## Session: 2026-03-26

### Commit 1: `5bf9ea2` — WI-1, WI-4, WI-5b, WI-5c

**WI-1: Fix preview-as-reader reactions (Bug)**

The core bug: `AuthorApp.tsx` rendered `<ReaderPage>` in preview mode without `onReactionAdd`/`onReactionRemove`, and without a `reviewToken`. Since `ReaderPage` line 150 guards API calls with `if (sessionId)` and preview has no session, reactions went to local `useState` only — silently lost on navigation.

Implemented Option A from the design guide: disable reactions in preview mode.

- Added `readOnly` prop threaded through `ReaderPage` -> `SectionView` -> `Paragraph`
- In `Paragraph`, when `readOnly=true`:
  - Added `dr-para--readonly` CSS class
  - Chips shown at 25% opacity with `pointer-events: none`
  - No hover highlight on paragraphs (border stays transparent)
- Added hatched preview banner at top of `ReaderPage` with "Preview mode -- reactions are disabled"
  - Banner uses `--dr-hatch-pattern` for Mac OS 1 aesthetic consistency
  - Text on solid black background for contrast
  - Includes "Back to Editor" button (wired to `onBackToEditor` prop)
- In `AuthorApp.tsx`, passed `readOnly` and `onBackToEditor={() => setView("edit")}` to preview ReaderPage

**Why Option A over Option C (ephemeral banner)?** Option A eliminates the confusion entirely. With Option C the user still sees reactions appear and disappear, which is confusing. With A, the visual affordance (faded chips) makes it obvious they're not interactive.

**WI-4: Rename "Share" to "Settings"**

One-line change: `Dashboard.tsx` button label `"Share"` -> `"Settings"`. The button calls `onArticleSettings()` which opens `ArticleSettings`, not the invite dialog. "Share" was misleading since the "Share" action (generating/copying invite links) is already handled by "+ Invite Reader".

**WI-5b: Remove Danger Zone label**

Removed `<div className="dr-editor__danger-zone">` wrapper and `<div className="dr-editor__danger-label">Danger Zone</div>` from `ArticleEditor.tsx`. The confirmation dialog ("Delete Section? Are you sure...") is sufficient protection. Also removed the associated CSS (dotted border, uppercase label).

**WI-5c: Improve reorder arrow visibility**

- Font size: 10px -> 13px (TriangleUp/Down glyphs)
- Padding: `1px 2px` -> `2px 4px`
- Added transition on opacity and background
- Hover state: inverts to `background: var(--dr-on-surface); color: var(--dr-surface)` — matches the Mac OS 1 button invert pattern used elsewhere

---

### Commit 2: `22a0b55` — WI-3a, WI-3b, WI-3c, WI-5a

**WI-3a: Article title/intro/version editing in Settings**

`ArticleSettings` previously only saved `{ status }` on save. Added three fields at the top in a new "ARTICLE DETAILS" section:

- Title (text input)
- Author Intro (textarea, 3 rows — the text shown on reader welcome splash)
- Version Label (text input — e.g. "Draft 1", "Revision 2")

All three are `useState` initialized from `article.title`, `article.intro`, `article.version`. `handleSave` now sends `{ title, intro, version, status }`.

These fields already existed on the `Article` type and the `updateArticle` PATCH mutation accepts `Partial<Article>`, so no backend/API changes needed.

**WI-3b: Delete article**

- Added `deleteArticle` mutation to `articleApi.ts`: `DELETE /api/articles/:id`, invalidates `["Article"]` tags
- Added `onDelete` optional prop to `ArticleSettings` — renders a "Delete Article" button under a new "DELETE" section header (only shown when `onDelete` is provided)
- Confirmation dialog with overlay backdrop, same style as editor delete confirmation: warns about permanent removal of sections, readers, and reactions
- Wired in `AuthorApp.tsx`: `onDelete` calls `deleteArticle(selectedArticle.id)` then `goBack()`
- Added MSW mock handler: `http.delete("/api/articles/:id")` — splices article from array, returns 204

**WI-3c: Unsaved changes indicator**

This was the most involved change. The editor works entirely with local state, so changes are invisible to the save system until "Save" is clicked. Users could add sections, edit content, then click Back and lose everything.

Approach:
- `initialJson = useMemo(() => JSON.stringify({ s: initialArticle.sections }), [initialArticle])` — serialized snapshot of initial state
- `isDirty` computed by comparing `JSON.stringify({ s: currentWithDraft.sections })` against `initialJson`
- `currentWithDraft = buildArticleWithCommittedDraft(article)` includes any uncommitted textarea content

UI changes:
- "Unsaved changes" italic label appears in the action bar (left-aligned via `margin-right: auto`)
- Back button: if dirty, shows confirmation dialog with "Discard" and "Save" options
- Delete confirmation dialog also got backdrop click-to-close and proper `&ldquo;...&rdquo;` quotes

**Tricky part:** The dirty check must include the current textarea content even if the user hasn't blurred (which triggers `commitContentDraft`). Using `buildArticleWithCommittedDraft` in the comparison ensures the textarea's uncommitted text is included.

**WI-5a: Auto-grow editor textarea**

- Added `textareaRef = useRef<HTMLTextAreaElement>(null)` and bound to textarea
- Removed `rows={16}` from textarea element
- Added `useEffect` that fires on `contentDraft` changes: sets `height: auto` then `height: max(scrollHeight, 200)px`
- CSS: `resize: none; overflow: hidden; min-height: 200px` (was `min-height: 300px`)

The double scroll problem (scrolling inside a textarea inside a scrolling window) is eliminated — the textarea always shows its full content.

---

## What's Not Done (Deferred)

**WI-2: URL routing** — Medium effort, requires replacing all `useState<View>` with `useNavigate()`/`useParams()`. Not implemented in this session. The infrastructure change is significant and should be its own focused commit.

**WI-6: Empty states cleanup** — Low priority, mostly already addressed by DR-004.

## How to Review

1. Start dev server: `cd frontend && npm run dev`
2. Preview reactions (WI-1): Open any article, click "Edit", click "Preview as Reader" — hatched banner should appear at top, reaction chips should be faded and non-clickable
3. Settings (WI-3a): Click "Settings" on dashboard — should see Title, Intro, Version fields at top, pre-filled from article data
4. Delete article (WI-3b): In Settings, scroll to DELETE section, click "Delete Article" — confirmation dialog should appear
5. Unsaved changes (WI-3c): Edit an article (change section title or content), click Back — should see "Unsaved changes" label and confirmation dialog
6. Auto-grow (WI-5a): Type multiple paragraphs in editor — textarea should grow; no double scrollbar
7. Reorder arrows (WI-5c): Hover the triangle arrows in editor sidebar — should invert on hover, clearly clickable
