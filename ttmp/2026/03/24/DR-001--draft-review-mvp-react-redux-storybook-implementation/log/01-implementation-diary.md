---
Title: Implementation Diary
Ticket: DR-001
Status: active
Topics:
    - react
    - redux
    - storybook
    - rtk-query
    - msw
    - frontend
DocType: log
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Step-by-step narrative of implementing the Draft Review MVP"
LastUpdated: 2026-03-24T18:17:51.41778593-04:00
WhatFor: "Recording what changed, why, what worked, what failed, and how to review"
WhenToUse: "When reviewing implementation history or resuming work on the Draft Review MVP"
---

# Implementation Diary

## 2026-03-24 18:20 - Project kickoff & planning

**What:** Created docmgr ticket DR-001, wrote implementation plan design doc, analyzed both JSX prototypes.

**Key findings from source analysis:**
- `reader-view.jsx` (869 lines): 9 components, reader-facing flow (welcome -> read sections -> react -> done)
- `review-system.jsx` (800 lines): 10 components, author dashboard + article reader + invite dialog
- 5 merge opportunities: Btn/MacButton, HatchBar/ProgressBar, P palettes, TitleBar, REACTIONS constants
- Both files use inline styles exclusively -- will migrate to CSS custom properties

**MVP scope decision:** 8 screens from the 27-screen spec. No auth, no payments, no teams. Personal tool for Manuel.

**Next:** Scaffold Vite + React 18 + TypeScript project (task 1.1).

## 2026-03-24 18:45 - Phase 1 complete: Foundation

**What:** Built the entire Phase 1 foundation layer.

**Scaffolding (1.1):** `npm create vite@latest frontend -- --template react-ts`, installed deps (react-redux, @reduxjs/toolkit, react-router-dom, msw).

**Theme (1.5):** Created `mac-os-1.css` with 30+ CSS custom properties extracting the hardcoded `P` palette from both prototypes. All tokens prefixed `--dr-` (draft-review). Created `tokens.ts` with `REACTION_TYPES` constant (merging `REACTIONS` + `REACTION_ICONS` from both files).

**Store (1.3):** Redux Toolkit store with 3 slices:
- `uiSlice` -- view routing, modal/menu state, selected article
- `readerSlice` -- reader session (started, currentSection, readSections)
- RTK Query `baseApi` with Article/Reaction/Reader tag types

**MSW (1.4):** Full mock layer:
- `db.ts` -- seeded with all 3 articles, 7 readers, 12 reactions from the prototypes
- `handlers.ts` -- 6 endpoints: articles CRUD, reactions, reader token resolution, invites
- `browser.ts` -- service worker setup, auto-started in dev mode

**Chrome components (1.6):** Extracted and CSS-ified 6 components:
- `TitleBar` -- stripes + close box (from both files, unified)
- `MacWindow` -- draggable/maximizable container (from review-system.jsx)
- `MacButton` -- merged Btn + MacButton (identical API, one component)
- `MacInput` / `MacTextArea` -- styled form controls (new, needed by multiple screens)
- `MenuBar` + `MenuDropdown` -- app shell (from review-system.jsx)

**Primitive components (1.7):** Extracted 5 components:
- `ProgressBar` -- merged HatchBar + ProgressBar
- `ReactionBadge` -- from reader-view.jsx
- `SectionNav` -- dot navigation from reader-view.jsx
- `ReactionPicker` -- extracted from Paragraph's inline picker (was deeply nested)
- `StatCard` -- extracted from Dashboard stats grid

**Storybook (1.2, 1.8):** Set up Storybook 10 with react-vite framework. Wrote 10 story files covering all chrome and primitive components with multiple variants each. Preview imports the theme CSS.

**Types (bonus):** TypeScript types for Article, Section, Reaction, Reader, plus DTOs.

**What worked well:**
- CSS custom properties migration was straightforward -- every `P.xxx` became `var(--dr-xxx)`
- Merging Btn/MacButton and HatchBar/ProgressBar was trivial since they had identical APIs
- Extracting ReactionPicker from the deeply-nested Paragraph component made it independently testable

**What was tricky:**
- The Vite scaffold template changed significantly -- had to replace the entire App.tsx and index.css
- Storybook init prompted interactively for Playwright -- needed `--yes` flag

**Build:** `npm run build` passes, `tsc --noEmit` clean.

**Next:** Phase 2 -- Reader experience components (WelcomeSplash, Paragraph, SectionView, ReaderPage).

## 2026-03-24 19:15 - Phase 2 complete: Reader Experience

**What:** Built the complete reader-facing flow as modular components.

**Components built:**
- `WelcomeSplash` -- pre-reading landing page with article info, author note, reaction guide. Refactored from reader-view.jsx lines 473-579. Changed from inline styles to CSS classes.
- `Paragraph` -- reactable paragraph with hover-to-add (+) button and inline ReactionPicker. Refactored from reader-view.jsx lines 190-342. **Key design change:** extracted the reaction picker into its own component (ReactionPicker from Phase 1), making Paragraph a composition of primitives rather than a monolith.
- `SectionView` -- section header (badge + title + stripe divider) + paragraph list + section reaction summary. New composition component that didn't exist as a standalone in the prototype.
- `ReaderToolbar` -- SectionNav + ProgressBar + percentage display. Thin composition component.
- `DoneDialog` -- review completion modal with reaction stats grid. Refactored from reader-view.jsx lines 401-470.
- `ReaderPage` -- the full composed reader page: WelcomeSplash -> ReaderToolbar + SectionView + navigation + DoneDialog. This replaces the monolithic `ReaderView` from the prototype. Includes keyboard navigation (arrow keys), section read tracking, and local reactions state.

**Architecture decision:** ReaderPage manages its own local state (started, currentSection, readSections, reactions) rather than wiring directly to Redux. It exposes `onReactionAdd` and `onReactionRemove` callbacks so the parent can dispatch to the store/API. This makes it independently testable and usable in Storybook without a Redux provider.

**Stories written:** WelcomeSplash (default + short article), DoneDialog (default + no reactions), ReaderPage (full article + short article). The ReaderPage story imports mock data from `mocks/db.ts`.

**What worked well:**
- Decomposing the monolithic ReaderView into 6 focused components was clean -- each component has a single responsibility
- ReactionPicker (extracted in Phase 1) plugs right into Paragraph without duplication

**Build:** `tsc --noEmit` clean.

**Next:** Phase 3 -- Author dashboard and article management.

## 2026-03-24 19:45 - Phase 3 complete: Author Dashboard & App Wiring

**What:** Built the author-facing components and wired the full App.tsx with MSW data.

**RTK Query API (3.1):** Created `articleApi.ts` with 6 endpoints: getArticles, getArticle, getReaders, getReactions, inviteReader, addReaction. All use the MSW handlers from Phase 1.

**Dashboard (3.2):** Full dashboard refactored from review-system.jsx's `Dashboard` component (lines 292-481):
- Article selector tabs
- Stat cards grid (readers, reactions, sections, avg progress)
- Two-column layout: readers panel with progress bars + reactions-by-section bar chart with hatched fill patterns per reaction type
- Draft-killer alert (auto-detects section with most confused+slow flags)
- Recent feedback list

**ArticleReader:** Split-pane author review view refactored from review-system.jsx's `ArticleReader` (lines 485-625):
- Section sidebar with reaction counts per type + warning icons
- Article text panel with paragraph display
- Bottom reactions panel with type filter tabs (All/Useful/Confused/Slow/Favorite)

**InviteDialog:** Reader invitation modal refactored from review-system.jsx's `InviteDialog` (lines 628-703). Two states: form -> sent confirmation.

**App.tsx wiring:** Connected everything with RTK Query hooks. Dashboard fetches articles/readers/reactions from MSW and passes to child components. View switching: dashboard -> article reader. InviteDialog opens as overlay.

**What was tricky:**
- TypeScript narrowing issue with `bookKiller` variable. Using `let` + `forEach` caused TS to narrow it to `never` inside the JSX conditional. Fixed by refactoring to `reduce()` which returns a properly typed result object.

**Build:** `npm run build` passes, `tsc` clean. Production bundle: ~290KB JS, ~20KB CSS.

**Next:** Phase 4 -- routing, end-to-end flow, polish.

## 2026-03-24 20:00 - Phase 4: Routing & end-to-end integration

**What:** Added react-router and split the app into author and reader entry points.

**Architecture:**
- `App.tsx` is now a pure router: `/r/:token` -> ReaderApp, `/*` -> AuthorApp
- `AuthorApp` (extracted from previous App.tsx) -- the dashboard/article-reader flow
- `ReaderApp` -- new component that resolves a reader token via `readerApi`, shows loading/error states, then renders ReaderPage

**RTK Query (readerApi):** New `getReaderLink` endpoint that hits `/api/r/:token`. The MSW handler resolves the token to a reader + article pair from the mock DB.

**Reader link flow (end-to-end):**
1. Author creates article (mock data for now)
2. Author invites reader -> gets a token (e.g. `tok-sarah-1`)
3. Reader opens `/r/tok-sarah-1`
4. ReaderApp fetches article data via token
5. ReaderPage renders: WelcomeSplash -> read sections -> react -> DoneDialog

**Error states:** ReaderApp shows a Mac OS 1 styled error window for invalid/expired tokens.

**Build:** `npm run build` passes. Final bundle: 332KB JS (105KB gzipped), 20KB CSS (3.6KB gzipped).

**How to validate:**
```bash
cd frontend
npm run dev
# Author dashboard: http://localhost:5173/
# Reader view: http://localhost:5173/r/tok-sarah-1
# Storybook: npm run storybook
```

**Summary of what was built across all 4 phases:**

| Layer | Files | Components |
|-------|-------|------------|
| Theme | 2 | CSS custom properties, tokens.ts |
| Chrome | 12 | TitleBar, MacWindow, MacButton, MacInput, MacTextArea, MenuBar |
| Primitives | 10 | ProgressBar, ReactionBadge, SectionNav, ReactionPicker, StatCard |
| Reader | 14 | WelcomeSplash, Paragraph, SectionView, ReaderToolbar, DoneDialog, ReaderPage |
| Author | 8 | Dashboard, ArticleReader, InviteDialog |
| Store | 4 | store, uiSlice, readerSlice, hooks |
| API | 3 | baseApi, articleApi, readerApi |
| Mocks | 3 | db (seeded data), handlers (6 endpoints), browser |
| App | 3 | App (router), AuthorApp, ReaderApp |
| Types | 4 | article, reaction, reader, index |
| Stories | 13 | All chrome, primitives, reader, and author components |

**Total: ~76 source files, 3500+ lines of component code, decomposed from 2 monolithic JSX prototypes (1669 lines).**

## 2026-03-24 20:30 - Article management trio: Editor, Settings, Manager

**What:** Built the three deferred Phase 3 screens that complete the author content workflow.

**API layer extensions:**
- MSW: added `PATCH /api/articles/:id` (update), `POST /api/articles` (create), `POST /api/articles/:id/share-token` (generate link)
- Mock DB: `updateArticle()`, `createArticle()`, `generateShareToken()` functions
- RTK Query: `createArticle`, `updateArticle`, `generateShareToken` mutations

**ArticleEditor (3.4):**
- Split-pane layout matching spec screen 08: section sidebar (200px) + content editor
- Sidebar: numbered section list with up/down reorder buttons, "Add Section", "Delete Section" (with confirmation dialog)
- Editor pane: title input (Chicago font, large) + content textarea (auto-splits on double newlines into paragraphs)
- Actions: "Preview as Reader" + "Save"
- Content model: `paragraphs.join("\n\n")` for editing, split back on save -- matching the prototype's approach where paragraph breaks define reactable units
- **Design note:** Chose up/down arrow buttons instead of full drag-and-drop for MVP simplicity. Drag-and-drop can be added later with a library like dnd-kit.

**ArticleSettings (3.5):**
- Three sections matching spec screen 09: SHARING, FEEDBACK, STATUS
- Sharing: review link display with Copy/Generate/Reset buttons, access mode radio (unique links / open link / password), reader visibility checkboxes
- Feedback: reaction type toggles (all 4 types), require-note checkbox, allow-anonymous checkbox
- Status: radio group (Draft / In Review / Complete / Archived)
- Share link generation calls `generateShareToken` mutation and displays the URL
- **Scope note:** Custom reactions and password field deferred. The settings are local state for now (only status persists via updateArticle).

**ArticleManager (3.3):**
- Toolbar: "New Article" button, search input, sort dropdown (Recent / Most Reactions / Most Readers / Status)
- Article cards: title + status badge, stats (sections/reactions/readers), version, progress bar, draft-killer alert, action buttons (Review / Edit / Share / Invite)
- Empty state: centered icon + message + "New Article" CTA
- Archived section: collapsible footer with simplified items
- Sorting and filtering via `useMemo`

**AuthorApp wiring:**
- Extended view type: `"dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"`
- View menu now has "Dashboard" and "Articles" actions
- Full navigation flow: Dashboard -> Articles -> Edit/Settings/Review -> Back
- Editor save calls `updateArticle` mutation then returns to dashboard
- Settings save calls `updateArticle` for status, generates share tokens

**Stories:** 6 new story files (2 per component: default + edge case).

**Build:** `npm run build` passes. Bundle: 345KB JS, 27KB CSS.

**How to validate:**
- Dashboard: click any article tab -> "Open Full Review"
- View menu -> "Articles" to see the ArticleManager
- Click "Edit" on any card -> ArticleEditor with sections sidebar
- Click "Share" on any card -> ArticleSettings with link generation
- In editor: add/delete/reorder sections, edit content, preview as reader
