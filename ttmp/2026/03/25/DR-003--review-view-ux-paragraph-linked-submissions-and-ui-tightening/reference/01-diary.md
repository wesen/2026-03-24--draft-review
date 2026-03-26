---
Title: Diary
Ticket: DR-003
Status: active
Topics:
    - frontend
    - react
    - ux
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.tsx
      Note: "Primary target: author review view with disconnected reaction list"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.css
      Note: "Review view styles — needs paragraph hover additions"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx
      Note: "Reference: reader paragraph with hover highlighting"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.css
      Note: "Reference: hover CSS to replicate in review view"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx
      Note: "Dashboard recent feedback — also missing paragraph context"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/app/AuthorApp.tsx
      Note: "Top-level author routing and state"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/ReaderPage.tsx
      Note: "Reader page — useEffect dep array issue found"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/SectionView.tsx
      Note: "Section rendering with paragraph ID generation"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/types/reaction.ts
      Note: "Reaction type — already has paragraphId"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleEditor.tsx
      Note: "Editor — no feedback indicators while editing"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleSettings.tsx
      Note: "Settings — most options are client-only"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleManager.tsx
      Note: "Article list — reviewed for UX consistency"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/WelcomeSplash.tsx
      Note: "Welcome screen — 'click +' copy is inaccurate"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/theme/mac-os-1.css
      Note: "Design tokens — reviewed for consistency"
ExternalSources: []
Summary: "Investigation diary for DR-003: analyzing the review view UX gaps and producing a comprehensive tightening plan."
LastUpdated: 2026-03-25T11:55:04.293911409-04:00
WhatFor: "Track the analysis process, decisions made, and findings for future reference."
WhenToUse: "When continuing work on this ticket or reviewing what was analyzed."
---

# Diary

## Goal

Document the investigation and analysis of the Draft Review application's UX/UI gaps, focusing on paragraph-linked submissions in the review view and a general tightening pass across all views.

## Step 1: Codebase Exploration and Gap Analysis

Explored the entire Draft Review frontend codebase to understand the architecture, data flow, and existing UX patterns before writing any recommendations. The goal was to understand *exactly* how reactions flow from reader to author view and identify every point where context is lost.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr ticket: when reviewing the submissions for a section, I want to see the submissions linked to the paragraph they are related to, and also highlight the hovered section of the document like in the reader view.

Analyze the issue and create a document in the ticket on how to tackle the issue, and in genral do a tightening pass on the UX and the UI of the application, now that it is kind of ready.

Keep a diary as you work (look up skills for docmgr and diary)"

**Assistant interpretation:** Create a docmgr ticket (DR-003), analyze the ArticleReader review view's disconnected reaction display, propose paragraph-level linking with hover highlighting, and audit the full application for UX/UI improvements. Document the process in a diary.

**Inferred user intent:** The application is feature-complete enough to start polishing. The author wants to shift from "make it work" to "make it good" — the review view is the most impactful place to improve because it's where the author spends time making decisions about their draft.

### What I did

- Read all 14 frontend components, their CSS files, and the type definitions
- Traced the reaction data flow end-to-end: reader creates reaction → stored with `paragraphId` → API returns to author → ArticleReader displays flat list
- Compared reader view (Paragraph.tsx) with author review view (ArticleReader.tsx) to identify the specific UX pattern the user wants replicated
- Audited every view (Dashboard, ArticleManager, ArticleEditor, ArticleSettings, ReaderPage, WelcomeSplash, DoneDialog) for consistency issues
- Created docmgr ticket DR-003 with design doc and diary

### Why

- Need complete understanding before writing recommendations — can't propose changes without knowing what exists
- The user said "tightening pass" which implies a comprehensive audit, not just the one feature request

### What worked

- The `Reaction` type already has `paragraphId` — no backend changes needed. The paragraph ID convention (`{sectionId}-p{index}`) is consistent between SectionView.tsx (where it's generated) and the stored reactions.
- The reader view's hover CSS (Paragraph.css) is clean and self-contained — easy to adapt for the review view.
- The explore agent was effective for getting a full picture of 14+ components quickly.

### What didn't work

- N/A (pure analysis step)

### What I learned

1. **The gap is purely in the view layer**: `ArticleReader.tsx` lines 93–97 render paragraphs as plain `<Prose>` without any `paragraphId` data attribute or hover state. Lines 142–157 render reactions without any paragraph context. The data (`reaction.paragraphId`) is already there — it's just not used.

2. **Paragraph IDs are generated by convention, not stored in the Section type**: `Section.paragraphs` is just `string[]`. The paragraph ID is constructed at render time as `${section.id}-p${index}`. This means paragraph IDs can shift if paragraphs are reordered or inserted in the editor. This is a latent bug — if an author edits between reviews, old reactions may point to wrong paragraphs. Not in scope for this ticket but worth noting.

3. **The 220px fixed height on the reactions panel** (`ArticleReader.css` line 100) is a hard constraint that makes the panel feel cramped with many reactions. This should become flexible.

4. **Dashboard's "Recent Feedback" (Dashboard.tsx lines 238–273)** already resolves section title from `sectionId` but doesn't attempt paragraph resolution — it could show a snippet using the same pattern.

5. **ReaderPage.tsx useEffect for keyboard nav (lines 169–183)** has no dependency array, meaning the event listener is re-added on every render. This is a performance issue and potential source of stale closures.

6. **WelcomeSplash.tsx line 44** says "click + to react" but the actual UI shows text labels (★ Useful, ? Confused, etc.), not a "+" button. Minor copy issue.

### What was tricky to build

N/A — this was an analysis step, not implementation.

### What warrants a second pair of eyes

- The paragraph ID stability issue (paragraph IDs shift when content is edited). This is architectural and might need a stable ID scheme in the future.
- Whether the panel layout (text above, reactions below) is the right long-term choice vs. inline annotations. I recommended keeping the panel approach for now because it fits the Mac OS 1 aesthetic.

### What should be done in the future

- Implement Phase 1 (paragraph-linked submissions) as defined in the design doc
- Address the paragraph ID stability issue before this goes to production with real users
- Consider a "feedback mode" overlay in the editor that shows reaction heat

### Code review instructions

- Start with the design doc: `ttmp/2026/03/25/DR-003--review-view-ux-paragraph-linked-submissions-and-ui-tightening/design/01-review-view-ux-analysis-and-implementation-plan.md`
- Review Section 3 (Implementation Plan) for the core feature approach
- Review Section 4 (UX/UI Tightening Audit) for the complete list of findings
- Review Section 5 (Implementation Order) for the suggested phasing

### Technical details

**Key file locations and what to change:**

| File | Current State | Needed Change |
|------|--------------|---------------|
| `ArticleReader.tsx:93-97` | Plain `<Prose>` render | Replace with interactive paragraph containers |
| `ArticleReader.tsx:27-32` | Filter by sectionId only | Add paragraphId filter + group-by |
| `ArticleReader.tsx:142-157` | Flat reaction card list | Group by paragraph with excerpts |
| `ArticleReader.css:99-104` | Fixed 220px reactions panel | Flex with min/max |
| `ArticleReader.css` (new) | No paragraph hover styles | Add `.dr-review-para` hover rules |

**Paragraph ID convention:**
```
Format: {sectionId}-p{paragraphIndex}
Example: s1-p0 (first paragraph of section "s1")
Source: SectionView.tsx line 43: paragraphId={`${section.id}-p${pi}`}
```

**Reaction filtering pattern:**
```typescript
// Current (ArticleReader.tsx:27-32)
const sectionReactions = reactions.filter(r => r.sectionId === selectedSection);
const filtered = reactionFilter
  ? sectionReactions.filter(r => r.type === reactionFilter)
  : sectionReactions;

// Proposed: add paragraph grouping
const grouped = Map.groupBy(filtered, r => r.paragraphId);
```

## Step 2: Visual UX/UI Design Review

Navigated through every view of the live application on port 5174 using Playwright, capturing 12 screenshots. Examined every component's CSS, the design token system, and the interaction patterns. Wrote a comprehensive 8-section design review covering visual identity, layout, typography, interaction patterns, information architecture, and prioritized recommendations.

### Prompt Context

**User prompt (verbatim):** "now go over the design of the full thing (I't running on :5174 I think) because our UX frontend design has been lost a bit . Load your frontend design skill and make a full pass on the UX / UI of the system, and write a separate design review / report and upload it to remarkable"

**Assistant interpretation:** Open the running app, visually audit every view, apply frontend design expertise to critique the Mac OS 1 aesthetic execution, produce a long-form design review document, and upload to reMarkable.

**Inferred user intent:** The application is at a maturity point where the design needs tightening. The user recognizes the visual language has drifted and wants a systematic audit to guide polish work.

### What I did

- Navigated through all 8+ views via Playwright (Dashboard, Review, Editor, Settings, All Articles, Invite Dialog, Reader Welcome, Reader Reading)
- Captured 12 screenshots including hover states and menu interactions
- Read all 24 CSS files and the design token system
- Compared every component against the Mac OS 1 design language
- Wrote an 8-section design review with 20 prioritized recommendations across 4 tiers

### Why

The user explicitly said the "UX frontend design has been lost a bit" — they need someone to systematically evaluate every pixel and tell them what's on-brand and what's drifted. This isn't about adding features, it's about enforcing the design language they chose.

### What worked

- Playwright screenshots were effective for capturing the exact current state including hover interactions
- The CSS variable system (`mac-os-1.css`) is well-organized and made it easy to understand what the intended design language is
- The reader view paragraph hover highlighting is genuinely well-executed and served as the benchmark for what "good" looks like in this app

### What didn't work

- Playwright browser crashed once when navigating to the reader URL (Chrome data dir conflict). Resolved by clearing `~/.cache/ms-playwright/mcp-chrome-*`.

### What I learned

1. **The Mac OS 1 aesthetic is strongest where custom components exist** (menu bar, scrollbars, buttons) and weakest where browser defaults leak through (form checkboxes, radio buttons, emoji). The fix isn't more styling — it's replacing browser defaults with custom components.

2. **The review view's top/bottom split (text above, reactions below at 220px) is fundamentally wrong for the use case**. The author comes to this view to read feedback, but the feedback panel gets the smallest share of the screen. This is the biggest layout-level issue.

3. **The dashboard tries to do too much**: it's an overview, a launcher, a metrics dashboard, and a recent-activity feed all at once. Consider whether it should be merged with the review view.

4. **The design token system in `mac-os-1.css` is well-structured** but only has ~35 tokens. It's missing tokens for: selected list item background, empty state text color, code block background, dialog backdrop color. Adding these would make the theme more complete.

5. **The "§ UNTITLED SECTION" header uses all-caps letterspacing** which is a nice touch — it feels like a MacWrite section label. This pattern should be used more consistently across section headers.

### What was tricky to build

N/A — this was a review/analysis step, not implementation.

### What warrants a second pair of eyes

- Whether the dashboard and review view should be merged (a significant IA change)
- Whether the reaction chips in the reader view should be hidden until hover (reduced noise vs discoverability tradeoff)
- Whether URL routing for author views is worth the implementation cost (significant change to navigation model)

### What should be done in the future

- Implement Tier 1 fixes (form controls, emoji, dialog shadow, copy fix) — these are small and high-impact
- Build a custom checkbox/radio component that matches the monochrome palette
- Consider a compact article list view (table-style) as an alternative to cards
- Bundle the ChicagoFLF webfont for typography authenticity

### Code review instructions

- Read the full design review: `design/02-ux-ui-design-review.md`
- Section 1 (Visual Identity Audit) covers what's on-brand vs drifted
- Section 2 (Layout and Spatial Composition) has view-by-view analysis
- Section 7 (Prioritized Recommendations) is the actionable list
- Section 8 (Design Language Summary) is a one-page reference for enforcing the aesthetic

### Technical details

**Screenshots captured**: 12 views at `screenshots/01-loading.png` through `screenshots/12-file-menu.png`

**Key design drift points identified**:
- `ArticleSettings.tsx`: Uses `<input type="checkbox">` and `<input type="radio">` — browser blue defaults
- `Dashboard.tsx` lines 99, 233: Uses color emoji 📭 and 📝
- `WelcomeSplash.tsx` line 44: "click +" copy doesn't match actual chip-based UI
- `InviteDialog.css`: Missing `box-shadow: var(--dr-shadow-dialog)` and backdrop dimming

**Design token gaps**: Missing `--dr-selected-bg`, `--dr-empty-state-color`, `--dr-code-bg`, `--dr-backdrop-color`

## Step 3: Phase 1 & 2 — Paragraph-Linked Submissions and Review Polish

Implemented the core feature: replacing plain `<Prose>` paragraphs in ArticleReader with interactive `ReviewParagraph` elements that support hover highlighting, reaction count badges, paragraph-grouped reactions with text excerpts, bidirectional hover highlighting, click-to-filter, and relative timestamps.

### Prompt Context

**User prompt (verbatim):** "Work on the tasks in DR-003 one by one, committing at appropriate intevals, keeping a detailed diary"

**Assistant interpretation:** Implement all 15 docmgr tasks across 4 phases, committing at each phase boundary, and updating the diary with implementation details.

**Inferred user intent:** Move from analysis/planning to actual implementation. The design documents are done; now execute the plan.

**Commit (code):** 2153022 — "Add paragraph-linked submissions to review view"
**Commit (code):** e699eca — "Add relative timestamps to review reaction cards"

### What I did

- Added `hoveredParagraph` and `filteredParagraph` state to ArticleReader
- Replaced flat `<Prose>` render with `<div data-paragraph-id>` wrappers with hover CSS matching the reader view's `.dr-para` pattern
- Grouped reactions by `paragraphId` using `Map` with paragraph text excerpts
- Added bidirectional hover: hovering a paragraph highlights its reaction cards, hovering a card highlights the paragraph and scrolls it into view
- Added click-to-filter by paragraph with a "Clear" bar
- Made the reactions panel flexible (min 150px, grows to 50%) instead of fixed 220px
- Added reaction count badges on paragraphs
- Added `timeAgo()` helper and timestamps on reaction cards

### Why

These six sub-tasks are tightly coupled — they all modify the same component and share state. Implementing them together prevents merge conflicts and ensures the feature works end-to-end.

### What worked

- The reader view's `.dr-para` CSS pattern translated directly to `.dr-review-para` with minimal changes. The hover transition (0.1s background + border-color) feels right.
- `Map.groupBy` (using `new Map` + for-loop for broader compat) cleanly groups reactions by paragraph.
- The `excerpt()` function strips markdown syntax and truncates to 80 chars, giving clean paragraph labels.

### What didn't work

- N/A — implementation was clean.

### What I learned

- The paragraph text lookup requires rebuilding a `Map<paragraphId, text>` from `section.paragraphs` array on each render. This is O(n) per render but n is small (typically <20 paragraphs per section).
- `scrollIntoView({ behavior: "smooth", block: "nearest" })` works well for the bidirectional hover — it only scrolls if the paragraph is off-screen.

### What was tricky to build

The `filteredParagraph` state interacts with `reactionFilter` (type filter). Both can be active simultaneously: "show only 'confused' reactions on paragraph 3". The filter chain is: sectionReactions → type filter → paragraph filter → grouped display. The "All" tab needed to be updated to only highlight when *both* filters are clear.

### What warrants a second pair of eyes

- The reaction panel `max-height: 50%` may not be ideal for all viewport sizes. On very small screens, both the text and panel could be cramped.
- The `excerpt()` function strips markdown with a simple regex — complex markdown (tables, nested lists) may produce odd excerpts.

### What should be done in the future

N/A — Phase 1 & 2 are complete.

### Code review instructions

- Start at `ArticleReader.tsx` — the entire feature is in this file and its CSS
- Key new state: `hoveredParagraph`, `filteredParagraph` (lines 42-45)
- Paragraph render: lines 152-180 (the `dr-review-para` wrapper)
- Grouped reactions: lines 246-287 (the `groupedByParagraph` rendering)
- CSS: `ArticleReader.css` — search for `.dr-review-para` block

### Technical details

Commits: `2153022`, `e699eca`

---

## Step 4: Phase 3 — Cross-View Consistency

Fixed four cross-cutting issues: Dashboard paragraph context, ReaderPage keyboard nav useEffect bug, WelcomeSplash copy, and loading/error states for API queries.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Commit (code):** 01ac26a — "Cross-view consistency: paragraph context, keyboard fix, loading states"

### What I did

1. **Dashboard Recent Feedback**: Added paragraph excerpt resolution — each feedback item now shows `¶ "The key insight..."` below the section title, using the same paragraph ID → text lookup pattern from Step 3
2. **ReaderPage useEffect**: Added `[started, isFirst, isLast, goNext, goPrev]` dependency array. Wrapped `goNext`/`goPrev` in `useCallback` to prevent the dep array from changing on every render
3. **WelcomeSplash**: Changed "click +" to "click a reaction chip" to match the actual UI
4. **AuthorApp**: Destructured `isLoading`/`error` from article query, added loading screen and error screen with retry button, added "Loading..." prefix to menu bar status when readers/reactions are loading

### Why

These fixes address the "cross-view consistency" tier of the design review — issues that span multiple components and affect the overall feel of the application.

### What worked

- The paragraph excerpt in Dashboard uses the same `r.paragraphId.match(/-p(\d+)$/)` pattern to resolve the paragraph index, then looks up the text in `section.paragraphs[idx]`.
- The useEffect fix is clean — adding the dep array plus useCallback on the nav functions is the standard React pattern.

### What didn't work

- N/A

### What I learned

The `goNext` function calls `markRead(currentSectionId)` which means it closes over `currentSectionId`. This creates a dependency chain: `goNext` → `currentSectionId` → re-creates `goNext` → re-subscribes keyboard listener. With `useCallback`, this is explicit and correct. Without it, the keyboard listener was being re-added on *every* render.

### What was tricky to build

The loading/error states in AuthorApp needed to work with the existing `authReady` gate. The sequence is: auth check → articles loading → articles error → data loading (readers/reactions). The `dataLoading` flag combines reader and reaction loading into a single "Loading..." prefix in the menu bar status, which is subtle but informative.

### What warrants a second pair of eyes

- The `articlesError` type from RTK Query is `FetchBaseQueryError | SerializedError | undefined` — the error UI just shows a generic message rather than trying to parse the error type.

### What should be done in the future

N/A — Phase 3 is complete.

### Code review instructions

- `Dashboard.tsx` lines 254-280 — paragraph excerpt resolution in Recent Feedback
- `ReaderPage.tsx` lines 99-106, 169-183 — useCallback wrappers and useEffect deps
- `WelcomeSplash.tsx` line 44 — copy fix
- `AuthorApp.tsx` lines 49-67, 205-250 — loading/error states

### Technical details

Commit: `01ac26a`

---

## Step 5: Phase 4 — Visual Integrity and Settings

Fixed the three biggest visual integrity issues from the design review: blue form controls, color emoji, and unimplemented settings with no indication.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Commit (code):** 430db43 — "Visual integrity: monochrome form controls, text glyphs, coming-soon badges"

### What I did

1. **Form controls**: Added `accent-color: black` globally in `index.css` for all checkboxes and radio buttons. Also added it specifically in `ArticleSettings.css` for redundancy.
2. **Color emoji → text glyphs**: Replaced 📝 (`\uD83D\uDCDD`) with ✎ (`\u270E`) and 📭 (`\uD83D\uDCED`) with ✉ (`\u2709`) in Dashboard and WelcomeSplash.
3. **Coming-soon badges**: Added a `.dr-settings__coming-soon` CSS class (tiny bordered uppercase label) and applied it to Access, Reader visibility, and Enabled reactions settings sections — the three groups whose state is local-only and not persisted to the backend.

### Why

The design review identified these as Tier 1 (must fix) issues. Blue checkboxes are the single most visible break from the Mac OS 1 aesthetic. Color emoji break the 1-bit monochrome palette. Unlabeled decorative settings confuse users who expect them to persist.

### What worked

- `accent-color: black` is a one-line fix that makes every checkbox and radio button monochrome. It's supported in all modern browsers.
- The "COMING SOON" badge uses `font-size: 8px; letter-spacing: 0.5px; text-transform: uppercase` which feels like a tiny Mac OS 1 system label.

### What didn't work

- N/A

### What I learned

The `accent-color` CSS property controls the color of form controls like checkboxes, radio buttons, range sliders, and progress bars. Setting it to `black` gives a clean monochrome look without needing custom checkbox components.

### What was tricky to build

Nothing particularly tricky in this step — it was straightforward CSS and copy changes.

### What warrants a second pair of eyes

- The ✎ (`\u270E`) glyph may render as a color emoji on some platforms (iOS, for example). If so, it would need to be wrapped in a text-only span or replaced with an SVG.

### What should be done in the future

- Build proper custom checkbox/radio components with pixel-art X marks and dots for full Mac OS 1 authenticity (accent-color is a good start but not pixel-perfect)
- Wire the "coming soon" settings to the backend when the API supports them
- Consider adding InviteDialog shadow and backdrop dimming (identified in design review but not in the task list)

### Code review instructions

- `index.css` lines 7-10 — global accent-color rule
- `ArticleSettings.css` — accent-color + `.dr-settings__coming-soon` class
- `ArticleSettings.tsx` — three "coming soon" badge placements
- `Dashboard.tsx` lines 39, 148 — emoji replacements
- `WelcomeSplash.tsx` line 21 — emoji replacement

### Technical details

Commit: `430db43`

---

## Summary

All 15 tasks in DR-003 completed across 4 commits:

| Phase | Commit | Tasks | Description |
|-------|--------|-------|-------------|
| 1 & 2 | `2153022`, `e699eca` | 2–10 | Paragraph-linked submissions, hover highlighting, grouped reactions, timestamps, badges |
| 3 | `01ac26a` | 11–14 | Dashboard paragraph context, useEffect fix, WelcomeSplash copy, loading states |
| 4 | `430db43` | 15 | Monochrome form controls, text glyphs, coming-soon badges |

## Related

- Design doc: `design/01-review-view-ux-analysis-and-implementation-plan.md`
- Design review: `design/02-ux-ui-design-review.md`
