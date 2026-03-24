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
