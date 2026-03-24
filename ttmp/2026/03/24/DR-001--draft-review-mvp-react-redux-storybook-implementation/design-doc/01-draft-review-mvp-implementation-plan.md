---
Title: Draft Review MVP Implementation Plan
Ticket: DR-001
Status: active
Topics:
    - react
    - redux
    - storybook
    - rtk-query
    - msw
    - frontend
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - imports/reader-view.jsx
    - imports/review-system.jsx
    - imports/draft-review-screen-spec.md
ExternalSources: []
Summary: "Modular, themeable React/Redux implementation plan for the Draft Review MVP -- a personal draft beta-reading tool with shareable review links."
LastUpdated: 2026-03-24T18:06:15.456203779-04:00
WhatFor: "Planning the architecture and phased implementation of the Draft Review application"
WhenToUse: "When starting implementation work on any phase of the Draft Review MVP"
---

# Draft Review MVP Implementation Plan

## Executive Summary

Build a **personal draft review application** that lets an author (Manuel) upload article drafts, generate shareable review links, and collect structured paragraph-level feedback from readers. The UI follows a retro Mac OS 1 aesthetic already prototyped in two JSX sketches (`reader-view.jsx` and `review-system.jsx`).

**MVP Scope**: Single author, no payment, no team settings. The core loop is:
1. Author creates/uploads an article with sections
2. Author generates a review link
3. Reader opens the link, reads section-by-section, leaves reactions (useful/confused/slow/favorite) with optional notes
4. Author views a dashboard with reader progress, reaction analytics, and draft-killer detection

**Tech stack**: React 18+, Redux Toolkit, RTK Query, Storybook, MSW (Mock Service Worker), CSS custom properties for theming, Vite.

---

## Problem Statement

The two prototype JSX files demonstrate the complete UX but are monolithic, have hardcoded inline styles, embed sample data, and have no state management or API layer. They need to be decomposed into a production-ready modular component library with proper state management, API integration, and test infrastructure.

---

## Source Analysis

### `reader-view.jsx` (869 lines) -- Reader-facing experience

**Components identified:**

| Component | Lines | Purpose | Reusability |
|-----------|-------|---------|-------------|
| `TitleBar` | 79-123 | Mac OS 1 window title bar with stripes, close box | Core chrome -- shared |
| `Btn` | 125-147 | Retro rounded button (primary/small variants) | Core chrome -- shared |
| `ReactionBadge` | 150-187 | Inline display of a submitted reaction | Reader + Author views |
| `Paragraph` | 190-342 | Paragraph with hover-to-react affordance + picker | Reader view only |
| `SectionNav` | 346-371 | Dot navigation showing read/unread sections | Reader toolbar |
| `HatchBar` | 374-398 | Retro hatched progress bar | Shared (reader + dashboard) |
| `DoneDialog` | 401-470 | Modal shown when review is complete | Reader flow |
| `WelcomeSplash` | 473-579 | Landing/intro screen before reading begins | Reader flow |
| `ReaderView` | 582-869 | Main reader app (menu bar + window + navigation) | Reader page |

**State managed locally:**
- `started` (boolean), `currentSection` (string), `readSections` (Set), `reactions` (array), `showDone` (boolean)
- Keyboard navigation (arrow keys)

**Hardcoded data:** `ARTICLE` object, `REACTIONS` array, `P` palette, font constants.

### `review-system.jsx` (800 lines) -- Author-facing dashboard

**Components identified:**

| Component | Lines | Purpose | Reusability |
|-----------|-------|---------|-------------|
| `MacWindow` | 88-157 | Draggable/maximizable window with title bar | Core chrome -- shared |
| `MacButton` | 160-176 | Same as `Btn` but different name | Merge with `Btn` |
| `MenuBar` | 179-244 | Mac OS 1 menu bar with dropdowns | App shell |
| `MenuDropdown` | 246-273 | Dropdown menu items | App shell |
| `ProgressBar` | 276-289 | Same concept as `HatchBar` | Merge with `HatchBar` |
| `Dashboard` | 292-481 | Article selector, stats, readers, reactions chart, draft-killer alert, recent feedback | Author home |
| `ArticleReader` | 485-625 | Split-pane: section sidebar + content + reaction panel with filters | Author review |
| `InviteDialog` | 628-703 | Modal to invite a reader by email | Author flow |
| `AboutDialog` | 706-748 | About modal | App shell |
| `App` | 751-800 | Main app with view routing | App root |

**Hardcoded data:** `ARTICLES` array (3 articles with sections, readers, reactions), `REACTION_ICONS`, `P` palette.

### Overlap & Merge Opportunities

1. **`TitleBar` / `MacWindow` title bar**: `MacWindow` is a superset (adds dragging, maximizing). Standardize on `MacWindow` as the container, extract `TitleBar` as a sub-component.
2. **`Btn` / `MacButton`**: Identical API. Merge into single `MacButton` component.
3. **`HatchBar` / `ProgressBar`**: Same visual. Merge into `ProgressBar` with hatched fill style.
4. **`P` palette**: Both files define nearly identical palettes. Extract to a single theme tokens file.
5. **`REACTIONS` / `REACTION_ICONS`**: Same data, different shapes. Normalize into one `REACTION_TYPES` constant.

---

## Proposed Solution: Modular Component Architecture

### Package Structure

```
packages/
  draft-review/
    src/
      theme/
        tokens.ts              # CSS custom properties + palette constants
        mac-os-1.css           # Default Mac OS 1 theme
        parts.ts               # CSS ::part() definitions for external styling

      chrome/                  # Mac OS 1 window chrome
        MacWindow.tsx           # Draggable/maximizable window container
        TitleBar.tsx            # Title bar with stripes + close box
        MenuBar.tsx             # Top menu bar with dropdowns
        MenuDropdown.tsx        # Dropdown menu
        MacButton.tsx           # Retro rounded button
        MacInput.tsx            # Styled text input
        MacTextArea.tsx         # Styled textarea
        index.ts

      primitives/              # Low-level reusable pieces
        ProgressBar.tsx         # Hatched progress bar (merges HatchBar+ProgressBar)
        ReactionBadge.tsx       # Reaction type icon + text badge
        SectionNav.tsx          # Dot navigation for sections
        ReactionPicker.tsx      # Reaction type selector + comment input
        StatCard.tsx            # Stats display card (icon + number + label)
        index.ts

      reader/                  # Reader-facing screens
        WelcomeSplash.tsx       # Pre-reading landing page
        Paragraph.tsx           # Reactable paragraph with hover affordance
        SectionView.tsx         # Section header + paragraphs
        ReaderToolbar.tsx       # Section nav + progress bar
        DoneDialog.tsx          # Review completion modal
        ReaderPage.tsx          # Composed reader experience
        index.ts

      author/                  # Author-facing screens
        Dashboard.tsx           # Article selector + stats + readers + chart
        ArticleReader.tsx       # Split-pane review with reaction filters
        ArticleManager.tsx      # List of all articles (from spec screen 06)
        ArticleEditor.tsx       # Section management editor (from spec screen 08)
        ArticleSettings.tsx     # Sharing & feedback settings (from spec screen 09)
        InviteDialog.tsx        # Reader invitation modal
        AboutDialog.tsx         # About modal
        index.ts

      store/                   # Redux Toolkit
        store.ts                # configureStore
        articleSlice.ts         # Articles CRUD + current selection
        readerSlice.ts          # Reader session state (for reader view)
        reactionSlice.ts        # Reactions state
        uiSlice.ts              # UI state (current view, modals, menus)

      api/                     # RTK Query
        articleApi.ts           # Articles endpoints
        reactionApi.ts          # Reactions endpoints
        readerApi.ts            # Reader/invite endpoints
        baseApi.ts              # createApi base

      mocks/                   # MSW
        handlers.ts             # Request handlers
        db.ts                   # In-memory mock database (seeded from current ARTICLES data)
        browser.ts              # Browser service worker setup
        server.ts               # Node service worker for tests

      types/                   # Shared TypeScript types
        article.ts              # Article, Section, Version
        reaction.ts             # Reaction, ReactionType
        reader.ts               # Reader, ReaderProgress

      app/                     # App shell
        App.tsx                 # Router + shell
        routes.tsx              # Route definitions
```

### Theming Strategy (CSS Custom Properties + `::part()`)

All visual tokens from the hardcoded `P` palette become CSS custom properties:

```css
/* mac-os-1.css (default theme) */
:root {
  --dr-bg: #e8e8e8;
  --dr-surface: #ffffff;
  --dr-on-surface: #000000;
  --dr-muted: #a0a0a0;
  --dr-subtle: #555555;
  --dr-border: #c0c0c0;

  --dr-font-title: "Chicago_12", "ChicagoFLF", "Geneva", monospace;
  --dr-font-body: "Geneva", "Monaco", monospace;

  --dr-border-width: 2px;
  --dr-border-radius-window: 3px;
  --dr-border-radius-button: 8px;

  --dr-stripe-pattern: repeating-linear-gradient(0deg, #fff 0px, #fff 1px, #000 1px, #000 2px);
  --dr-hatch-pattern: repeating-linear-gradient(45deg, var(--dr-on-surface) 0px, var(--dr-on-surface) 2px, var(--dr-surface) 2px, var(--dr-surface) 4px);
  --dr-desktop-pattern: repeating-conic-gradient(var(--dr-border) 0% 25%, transparent 0% 50%) 0 0 / 4px 4px;
}
```

Components use these variables exclusively. A different theme CSS file can override all visuals without touching component code.

### State Management Design

```
                    ┌─────────────────────────┐
                    │      RTK Query API       │
                    │  (articleApi, readerApi,  │
                    │   reactionApi)            │
                    └──────────┬──────────────┘
                               │ cache
                    ┌──────────┴──────────────┐
                    │     Redux Store          │
                    │                          │
                    │  articleSlice:            │
                    │    articles[]             │
                    │    selectedArticleId      │
                    │                          │
                    │  readerSlice:             │
                    │    currentSection         │
                    │    readSections Set       │
                    │    started boolean        │
                    │                          │
                    │  reactionSlice:           │
                    │    reactions[]             │
                    │    (optimistic updates)   │
                    │                          │
                    │  uiSlice:                 │
                    │    view (dashboard|       │
                    │     article|reader)       │
                    │    modals {}              │
                    │    activeMenu             │
                    └──────────────────────────┘
```

**Reader view** uses `readerSlice` + `reactionApi` mutations (optimistic).
**Author view** uses `articleSlice` + all APIs for fetching dashboard data.

### RTK Query API Shape

```typescript
// articleApi.ts
const articleApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getArticles: build.query<Article[], void>(),
    getArticle: build.query<Article, string>(),
    createArticle: build.mutation<Article, CreateArticleDto>(),
    updateArticle: build.mutation<Article, UpdateArticleDto>(),

    // Reader-facing: fetch article by review token
    getArticleByToken: build.query<ReaderArticle, string>(),
  }),
});

// reactionApi.ts
const reactionApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getReactions: build.query<Reaction[], { articleId: string }>(),
    addReaction: build.mutation<Reaction, AddReactionDto>(),
    removeReaction: build.mutation<void, string>(),
  }),
});

// readerApi.ts
const readerApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getReaders: build.query<Reader[], { articleId: string }>(),
    inviteReader: build.mutation<Reader, InviteReaderDto>(),
    updateProgress: build.mutation<void, { token: string; sectionId: string }>(),
  }),
});
```

### MSW Mock Layer

MSW handlers serve as both the development backend and test fixtures:

```typescript
// handlers.ts
export const handlers = [
  rest.get('/api/articles', (req, res, ctx) => {
    return res(ctx.json(db.articles));
  }),
  rest.get('/api/articles/:id', ...),
  rest.post('/api/articles/:id/reactions', ...),
  rest.get('/api/r/:token', ...),  // reader link resolution
  rest.post('/api/r/:token/progress', ...),
  rest.post('/api/invite', ...),
];
```

The mock DB is seeded from the current `ARTICLES` data in `review-system.jsx`, providing realistic demo data from day one.

---

## MVP Screens (Scoped Down from Full Spec)

From the 27-screen spec, the MVP includes **8 screens**:

| # | Screen | Priority | Source |
|---|--------|----------|--------|
| 1 | **Dashboard** (author home) | P0 | `review-system.jsx` Dashboard |
| 2 | **Article Manager** (list) | P0 | Spec screen 06 |
| 3 | **New Article / Upload** | P0 | Spec screen 07 (paste text + markdown only) |
| 4 | **Article Editor** | P0 | Spec screen 08 |
| 5 | **Article Settings & Sharing** | P0 | Spec screen 09 (link generation only) |
| 6 | **Reader Welcome** | P0 | `reader-view.jsx` WelcomeSplash |
| 7 | **Reader View** | P0 | `reader-view.jsx` ReaderView |
| 8 | **Reader Summary** (post-read) | P1 | `reader-view.jsx` DoneDialog (expanded) |

**Explicitly deferred:**
- Auth screens (01-04) -- use simple local auth or none for MVP
- Invite modal (10) -- simplified inline in settings
- Analytics deep dive (11) -- dashboard stats suffice
- Feedback deep dive (12) -- dashboard recent feedback suffices
- Feedback export (13) -- later
- Version history (14) -- later
- Reader management (15) -- later
- Reader history (19) -- later
- Profile/Account (20-22) -- not needed for single user
- Team/Admin (23-24) -- not needed
- Onboarding (26) -- not needed for single user

---

## Phased Implementation Plan

### Phase 1: Foundation (Storybook-first)

**Goal:** Extract all shared chrome and primitive components, set up theming, Storybook.

1. **Project scaffold**
   - Vite + React 18 + TypeScript
   - Storybook 8 setup
   - Redux Toolkit + RTK Query wiring
   - MSW browser + node setup
   - CSS custom properties theme file

2. **Chrome components** (from both prototypes)
   - `TitleBar` -- extracted from both, unified
   - `MacWindow` -- draggable/maximizable, uses TitleBar
   - `MacButton` -- merged from Btn + MacButton
   - `MacInput` / `MacTextArea` -- styled form controls
   - `MenuBar` + `MenuDropdown` -- from review-system.jsx
   - Storybook stories for each with theme knobs

3. **Primitive components**
   - `ProgressBar` -- merged HatchBar + ProgressBar
   - `ReactionBadge` -- from reader-view.jsx
   - `SectionNav` -- from reader-view.jsx
   - `ReactionPicker` -- extracted from Paragraph's inline picker
   - `StatCard` -- extracted from Dashboard stats grid
   - Storybook stories for each

### Phase 2: Reader Experience

**Goal:** Full reader flow working against MSW mocks.

1. **Types + API layer**
   - TypeScript types for Article, Section, Reaction, Reader
   - RTK Query endpoints (getArticleByToken, addReaction, updateProgress)
   - MSW handlers for reader endpoints

2. **Reader components**
   - `WelcomeSplash` -- refactored from reader-view.jsx, uses theme tokens
   - `Paragraph` -- refactored, uses ReactionPicker, connected to reactionApi
   - `SectionView` -- section header + paragraph list
   - `ReaderToolbar` -- SectionNav + ProgressBar, connected to readerSlice
   - `DoneDialog` -- refactored, reads stats from store
   - `ReaderPage` -- composed page, keyboard nav, routing by token

3. **Reader Redux slice**
   - `readerSlice`: currentSection, readSections, started
   - Selectors: progress%, currentSectionData, sectionReactions

### Phase 3: Author Dashboard & Article Management

**Goal:** Author can see articles, view feedback, manage content.

1. **Author API layer**
   - RTK Query: getArticles, getArticle, createArticle, updateArticle
   - MSW handlers seeded with ARTICLES data
   - Reactions + readers endpoints

2. **Dashboard**
   - Article selector tabs
   - Stat cards (readers, reactions, sections, avg progress)
   - Readers panel with progress bars
   - Reactions-by-section bar chart
   - Draft-killer alert
   - Recent feedback list

3. **Article Manager**
   - Scrollable article card list
   - Search + sort
   - Status badges
   - Actions: edit, share, archive

4. **Article Editor**
   - Split-pane: sortable section list + content editor
   - Add/delete/reorder sections
   - Preview as reader

5. **Article Settings**
   - Review link generation + copy
   - Access mode selection
   - Reaction type toggles
   - Status selector

### Phase 4: Integration & Polish

**Goal:** Connect reader + author flows end-to-end, routing, error states.

1. **Routing**
   - `/` -- Dashboard (author)
   - `/articles` -- Article Manager
   - `/articles/:id/edit` -- Article Editor
   - `/articles/:id/settings` -- Article Settings
   - `/articles/:id/review` -- Author review view
   - `/r/:token` -- Reader entry point

2. **End-to-end flow**
   - Author creates article -> generates link -> reader opens link -> leaves reactions -> author sees them on dashboard

3. **Error & empty states**
   - Mac OS 1 styled error dialogs
   - Empty states for no articles, no readers, no reactions

4. **Keyboard shortcuts**
   - Arrow keys for reader navigation
   - Cmd+N for new article
   - Global shortcut modal

---

## Design Decisions

### 1. CSS Custom Properties over CSS-in-JS
**Decision:** Use CSS custom properties for theming instead of the current inline styles.
**Rationale:** The prototypes use inline `style={{}}` everywhere, which prevents theming. CSS custom properties allow theme overrides via a single CSS file swap, support `::part()` for web component-style external styling, and are more performant than runtime style computation.

### 2. Redux Toolkit over local state
**Decision:** Move from per-component `useState` to Redux slices.
**Rationale:** The reader view has interconnected state (current section, read sections, reactions) that benefits from centralized management. The author dashboard aggregates data across articles. RTK Query provides cache management and optimistic updates for reactions.

### 3. MSW as development backend
**Decision:** Use MSW instead of building a real backend first.
**Rationale:** This is a personal tool MVP. MSW lets us build the full UI with realistic data, test all flows, and defer backend decisions. The mock DB can be persisted to localStorage for demo purposes. When a real backend exists, only the MSW layer is removed.

### 4. Merge duplicated components
**Decision:** Merge `Btn`/`MacButton`, `HatchBar`/`ProgressBar`, and the two `P` palettes.
**Rationale:** Both prototypes independently created the same components. The merged versions become the single source of truth in the component library.

### 5. MVP scoping: no auth, no teams, no payments
**Decision:** Skip screens 01-04 (auth), 20-24 (account/team/admin).
**Rationale:** This is a personal tool. Auth can be added later. The reader-link model (readers don't need accounts) is the core innovation and works without auth infrastructure.

---

## Alternatives Considered

### Next.js instead of Vite SPA
Rejected because the app is a single-page experience with no SEO requirements. Reader links resolve client-side against MSW/API. Next.js adds routing complexity we don't need.

### Tailwind CSS instead of CSS custom properties
Rejected because the Mac OS 1 aesthetic requires precise pixel values, custom patterns (stripe/hatch gradients), and specific font stacks that don't map well to utility classes. CSS custom properties with targeted component CSS gives more control.

### Zustand instead of Redux Toolkit
Rejected because RTK Query provides the API cache layer we need, and Redux DevTools are invaluable for debugging the interconnected state (reader progress, reactions, article selection). The overhead of Redux Toolkit is minimal with the slice pattern.

---

## Storybook Strategy

Every component gets stories in three categories:

1. **Default** -- component with default theme and typical props
2. **Variants** -- all visual states (primary/secondary, small/large, empty/populated, error)
3. **Themed** -- same component with an alternative theme CSS applied

Stories also serve as visual regression test targets.

**Story organization:**
```
Chrome/
  TitleBar
  MacWindow
  MacButton
  MacInput
  MenuBar
Primitives/
  ProgressBar
  ReactionBadge
  SectionNav
  ReactionPicker
  StatCard
Reader/
  WelcomeSplash
  Paragraph
  SectionView
  ReaderToolbar
  DoneDialog
  ReaderPage (full page story with MSW)
Author/
  Dashboard
  ArticleManager
  ArticleEditor
  ArticleSettings
  InviteDialog
```

---

## Open Questions

1. **Persistence for MVP**: Should MSW mock data persist to localStorage, or is ephemeral (refreshes reset) acceptable for initial development?
2. **Reader identity**: Readers don't auth, but do we track them by token only, or also collect a name on the welcome screen?
3. **Real-time updates**: Should the author dashboard poll for new reactions, or is manual refresh sufficient for MVP?
4. **Markdown rendering**: Should the reader view render markdown in paragraphs, or treat content as plain text (as in the prototypes)?
5. **Mobile responsiveness**: The Mac OS 1 aesthetic is desktop-centric. Do we need a mobile reader view for MVP?

---

## References

- `imports/reader-view.jsx` -- Reader-facing prototype (869 lines)
- `imports/review-system.jsx` -- Author-facing prototype (800 lines)
- `imports/draft-review-screen-spec.md` -- Full 27-screen specification with ASCII mockups and YAML DSL
- Beta reading methodology: "Help This Book" model (drop-off curves, draft-killer detection)
