---
Title: Rolling code review and improvement backlog
Ticket: DR-008
Status: active
Topics:
    - code-review
    - frontend
    - backend
    - react
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/app/AuthorApp.tsx
      Note: Large orchestrator component with routing, data loading, and modal control mixed together
    - Path: frontend/src/store/uiSlice.ts
      Note: View state duplicates what should increasingly be route state
    - Path: frontend/src/author/InviteDialog.tsx
      Note: UI copy and validation rules are inconsistent
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: Component mixes orchestration, persistence, navigation, and optimistic reaction state
    - Path: frontend/src/api/articleApi.ts
      Note: Client contract drifted ahead of server implementation
    - Path: pkg/server/http.go
      Note: Router and handlers are growing into a very large single file
    - Path: pkg/articles/service.go
      Note: Missing operations reveal end-to-end contract gaps
    - Path: pkg/reviewlinks/service.go
      Note: Business rules are still too narrow for the intended product
ExternalSources: []
Summary: Ongoing code-review notes collected while investigating DR-008. This is intentionally broader than the user-requested features and should continue to grow as the codebase evolves.
LastUpdated: 2026-03-26T12:35:00-04:00
WhatFor: Capture subpar or fragile areas noticed during investigation so they can become future cleanup tickets instead of disappearing into chat history.
WhenToUse: Use this document when planning refactors, code cleanup work, or intern onboarding.
---

# Rolling Code Review and Improvement Backlog

## Purpose

This document is intentionally broader than the immediate DR-008 feature request. While investigating the missing author UX pieces, several architectural and implementation patterns stood out as likely future sources of bugs, rework, or confusion. The goal here is not to criticize the codebase for being imperfect. The goal is to preserve observations while they are fresh so they can become deliberate follow-up work.

## High-Priority Findings

### 1. Client/server contract drift is already happening

The clearest example is article deletion:

- `frontend/src/api/articleApi.ts` exposes `DELETE /articles/{id}`
- `frontend/src/author/ArticleSettings.tsx` exposes a delete button
- `frontend/src/app/AuthorApp.tsx` calls the mutation
- `pkg/server/http.go` does not register a delete route
- `pkg/articles/service.go` does not offer a delete method

Why this matters:

- this is exactly how “it looks implemented” bugs are born
- it means frontend and backend work are not being verified together often enough

Recommended improvement:

- add contract-level integration tests for all non-trivial author actions
- adopt a habit that any new RTK Query mutation must be verified against a real handler in the same slice of work

### 2. `AuthorApp` is carrying too many responsibilities

`frontend/src/app/AuthorApp.tsx` currently owns:

- authentication gating
- article fetching
- reader and reaction fetching
- top-level menu construction
- view-state navigation
- modal orchestration
- create/update/delete article flow
- share-link generation and invite creation
- preview-mode wiring

Why this matters:

- it is getting harder to reason about
- route migration will become more painful if everything remains centralized
- subtle bugs become more likely because too many concerns share one render tree

Recommended improvement:

- split route pages from shared layout
- move data-specific actions closer to each page or feature module
- treat `AuthorApp` as an author shell, not a page god object

### 3. Route state and Redux UI state are currently overlapping

`frontend/src/store/uiSlice.ts` stores:

- `view`
- `selectedArticleId`
- `focusSectionId`
- `previewArticle`

This is convenient in a local prototype, but brittle in a routed web application.

Why this matters:

- some state belongs in the URL
- some state belongs in RTK Query cache
- some state belongs in component-local draft buffers
- currently these responsibilities are blurred

Recommended improvement:

- define a routing/state ownership document
- aggressively remove any state that simply mirrors the current URL or current query result

### 4. Share-link product terminology and domain model are not aligned

The modal says one thing, the validation rules say another, and the backend model says something narrower still.

Why this matters:

- the more the codebase grows, the more expensive it becomes to retrofit the domain language
- UX confusion often indicates deeper model confusion

Recommended improvement:

- explicitly define link types and identity modes in one short domain document
- update frontend copy, API payloads, database schema, and analytics language to match

### 5. `pkg/server/http.go` is on track to become a bottleneck file

This file already owns:

- route registration
- auth route wiring
- many different author and reader handlers
- frontend serving and proxy behavior
- request decoding and response shaping across multiple feature areas

Why this matters:

- large handler files are hard to test and easy to accidentally couple
- every new feature increases merge-conflict pressure in the same file

Recommended improvement:

- split handler registration and feature handlers by area:
  - article handlers
  - review-link handlers
  - review-session handlers
  - auth/info handlers

## Medium-Priority Findings

### 6. Preview mode is currently a separate behavior branch instead of a first-class workflow

`ReaderPage` receives a `readOnly` prop that disables reactions and avoids the normal persisted flow.

Why this matters:

- it risks creating duplicated review logic if preview grows
- it is already a product pressure point

Recommended improvement:

- choose whether preview is:
  - visual-only
  - or real test-review mode
- then model it explicitly rather than hanging more logic off `readOnly`

### 7. Reader state may be too thin for future resume/reload behavior

`frontend/src/store/readerSlice.ts` only stores:

- started
- sessionId
- currentSectionId
- readSectionIds

This is fine today, but limited if you later want:

- resume after reload
- local optimistic summaries
- preview/test-session markers
- link identity details

Recommended improvement:

- define a future reader-session state shape before adding more booleans ad hoc

### 8. Invite dialog depends on ephemeral local state too heavily

`frontend/src/author/InviteDialog.tsx` stores the currently generated share URL locally and starts blank every time.

Why this matters:

- a stable reusable link should feel like article state, not throwaway modal state
- reopening the modal should not feel like the app forgot a key part of the article

Recommended improvement:

- either fetch/share the current share link from article data
- or keep article share-link state in RTK Query cache and pass it down explicitly

### 9. Missing explicit error surfaces in author flows

Many async flows call mutations inline and do not present robust user feedback for failures.

Examples:

- create article
- update article
- delete article
- generate share link
- create invite

Recommended improvement:

- standardize optimistic state, loading state, and failure state for author mutations
- avoid silent or console-only failures for important author actions

## Lower-Priority Findings

### 10. Some feature flags are implicit rather than modeled

There are hints of future features such as:

- enabled reactions
- require note
- allow anonymous

but parts of the settings UI are labeled “coming soon” while the backend already has related fields.

Why this matters:

- partial settings create false confidence
- UI and backend may drift again if these flags are not treated as real features with end-to-end tests

### 11. The current component naming is sometimes screen-first and sometimes role-first

Examples:

- `AuthorApp`
- `ArticleManager`
- `ArticleReader`
- `ReaderPage`
- `WelcomeSplash`

This is not a bug, but as the app grows, inconsistent naming makes the codebase harder to scan.

Recommended improvement:

- choose one naming scheme for route pages, feature widgets, and layout shells

### 12. The codebase needs a short architecture note for newcomers

A new intern can understand the system, but only by jumping between:

- React Router
- Redux slices
- RTK Query APIs
- Go handlers
- Go services
- Postgres repositories

Recommended improvement:

- keep one lightweight living architecture note in the main repo docs, not only in ticket workspaces

## Suggested Follow-Up Tickets

- Split `AuthorApp` into route pages and layout shell
- Split `pkg/server/http.go` into feature-specific handler files
- Add end-to-end API contract smoke tests for all author mutations
- Define and implement reader-link identity modes
- Convert preview mode into first-class persisted test-review workflow
- Create a state ownership guide for route state vs Redux state vs query cache vs local draft state

## Review Heuristic for Future Work

When reviewing future Draft Review changes, ask four questions:

1. Is this state truly local, or should it be in the URL?
2. Is this API contract implemented end to end, or only in the client?
3. Is this user-facing term backed by an explicit domain model?
4. Is this preview/test behavior a real workflow or just a mock branch?

If a change cannot answer those questions cleanly, it probably needs more design work before more code is added.
