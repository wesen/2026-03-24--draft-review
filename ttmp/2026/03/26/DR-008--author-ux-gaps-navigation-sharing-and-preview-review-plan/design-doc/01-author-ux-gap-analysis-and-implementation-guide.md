---
Title: Author UX gap analysis and implementation guide
Ticket: DR-008
Status: active
Topics:
    - frontend
    - react
    - ui
    - ux
    - backend
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/App.tsx
      Note: React Router is already mounted at the application edge
    - Path: frontend/src/app/AuthorApp.tsx
      Note: Author navigation, modal wiring, and article create/update/delete flows
    - Path: frontend/src/store/uiSlice.ts
      Note: Current view state machine and back behavior
    - Path: frontend/src/author/ArticleSettings.tsx
      Note: Delete UI exists already, but backend delete does not
    - Path: frontend/src/author/InviteDialog.tsx
      Note: Current share-link and reader-link modal behavior
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: Reader preview mode currently disables review persistence
    - Path: frontend/src/api/articleApi.ts
      Note: Frontend assumes DELETE /articles/{id} exists
    - Path: frontend/src/api/readerApi.ts
      Note: Review session and review persistence client contract
    - Path: pkg/server/http.go
      Note: HTTP routing surface and missing DELETE article endpoint
    - Path: pkg/articles/service.go
      Note: Article service interface has no delete method
    - Path: pkg/reviewlinks/types.go
      Note: Invite and resolved-link domain model
    - Path: pkg/reviewlinks/service.go
      Note: Invite validation currently requires an email address
    - Path: pkg/reviewlinks/postgres.go
      Note: Share token vs invite token persistence model
    - Path: pkg/reviews/postgres.go
      Note: Session identity behavior for invited vs generic readers
ExternalSources: []
Summary: Detailed analysis of the remaining author UX gaps around article deletion, URL routing, reader-link generation, reader identity, and preview review behavior, with a phased implementation plan for an intern.
LastUpdated: 2026-03-26T12:35:00-04:00
WhatFor: Plan and explain the remaining product and implementation work for the author-facing UX and reader-link system.
WhenToUse: Use this when implementing the next author UX slice or when onboarding a new engineer or intern to the Draft Review frontend and sharing model.
---

# Author UX Gap Analysis and Implementation Guide

## Executive Summary

Draft Review is close to a coherent author-and-reader workflow, but the current implementation still has several product seams that are visible to users:

- article deletion is exposed in the frontend but not implemented in the backend
- author navigation is driven by Redux view strings instead of URL routes, so browser back/forward is mostly fake
- the share-link modal mixes several concepts and does not cleanly support generic links, tracked reader links, and anonymous-but-identifiable review sessions
- author preview mode renders the reader UI but explicitly disables reaction persistence, which limits testing and demo workflows

The important architectural point is that these are not isolated bugs. They are all connected by one underlying issue: the app has grown from a mock-first single-screen prototype into a real full-stack application, but some frontend state models still assume a local-only app while the backend now has real identities, persisted articles, persisted review sessions, and share tokens.

This document explains the current system, identifies the exact failure points, and proposes a phased implementation plan that keeps the codebase moving toward a cleaner architecture instead of adding more special-case UI behavior.

## Product Goals

The missing product behavior can be expressed more clearly as five goals:

1. An author can permanently delete an article and all associated data in a predictable, safe way.
2. The author area has real URL routes so deep links, browser back/forward, refresh, and page reload all preserve the current view.
3. The share flow allows the author to create:
   - a generic reusable share link
   - a person-specific tracked link tied to an email or name
   - optionally an anonymous tracked link that still creates a stable reviewer identity once opened
4. The share dialog explains these link types clearly and does not label a field as optional while enforcing it as required.
5. Preview mode can optionally create test review sessions so authors can experience and validate the end-to-end reader flow before sending links to real readers.

## System Overview

### Frontend structure

The frontend is a React application with Redux Toolkit and RTK Query.

Important entry points:

- `frontend/src/App.tsx` mounts React Router and already distinguishes author routes from reader routes.
- `frontend/src/app/AuthorApp.tsx` is the main author shell.
- `frontend/src/app/ReaderApp.tsx` is the reader shell for `/r/:token`.
- `frontend/src/store/uiSlice.ts` holds author view state such as `dashboard`, `edit`, and `settings`.
- `frontend/src/store/readerSlice.ts` holds reader session state such as the active section and started status.
- `frontend/src/api/articleApi.ts` and `frontend/src/api/readerApi.ts` are the main frontend API contracts.

### Backend structure

The backend is a Go server with package-level domain services.

Important backend packages:

- `pkg/server/http.go` registers HTTP routes and translates HTTP requests into service calls.
- `pkg/articles/service.go` and `pkg/articles/postgres.go` own author article data.
- `pkg/reviewlinks/service.go` and `pkg/reviewlinks/postgres.go` own share tokens and reader invites.
- `pkg/reviews/postgres.go` owns review session startup, progress, reactions, and summary persistence.

### Current architecture diagram

```text
Author Browser
    |
    |  React + Redux + RTK Query
    v
AuthorApp
    |
    |  /api/articles, /api/articles/{id}/invite, /api/articles/{id}/share-token
    v
Go HTTP Handlers
    |
    +--> articles.Service ------> articles.PostgresRepository
    |
    +--> reviewlinks.Service ---> reviewlinks.PostgresRepository
    |
    +--> reviews.Service -------> reviews.PostgresRepository
    |
    v
PostgreSQL

Reader Browser
    |
    |  /r/{token}, /api/r/{token}/start, /api/reviews/{sessionId}/*
    v
ReaderApp / ReaderPage
```

### The most important conceptual mismatch

React Router already exists at the application edge, but `AuthorApp` still acts like a desktop prototype that manually swaps screens using Redux `view` state. That means the app looks like a routed application, but the author experience does not behave like one.

At the same time, the sharing and review flow already uses real backend sessions and tokens. This means user-visible state is now partly local and partly persisted. The remaining UX bugs largely come from the boundaries between those two worlds.

## Current State by Area

### 1. Article deletion

The frontend already assumes deletion exists.

Evidence:

- `frontend/src/author/ArticleSettings.tsx` renders a delete section and confirmation dialog.
- `frontend/src/app/AuthorApp.tsx` passes `onDelete` into `ArticleSettings`.
- `frontend/src/api/articleApi.ts` defines `deleteArticle` as `DELETE /articles/{id}`.

The backend does not implement deletion.

Evidence:

- `pkg/server/http.go` registers `GET`, `POST`, and `PATCH` article routes, but no `DELETE /api/articles/{id}` route.
- `pkg/articles/service.go` does not define `DeleteArticle` in the repository interface or service layer.

Practical result:

- the UI offers deletion
- the RTK Query client calls a route that does not exist
- the user experiences a broken flow even though the product appears to support it

### 2. Author routing and browser navigation

The app root already mounts React Router:

- `frontend/src/App.tsx` routes `/r/:token` to `ReaderApp`
- all remaining paths are routed to `AuthorApp`

Inside `AuthorApp`, view changes are Redux-only:

- selecting an article dispatches `setView("article")`
- editing dispatches `setView("edit")`
- settings dispatches `setView("settings")`
- preview dispatches `setView("reader-preview")`

The current back behavior is not real navigation:

- `frontend/src/store/uiSlice.ts` implements `goBack()` by resetting everything to dashboard
- there is no history stack
- there is no serialization of `selectedArticleId`, `focusSectionId`, preview state, or modal state into the URL
- browser back/forward does not correspond to author navigation

Practical result:

- refresh loses author context
- a user cannot deep link an article editor or settings page
- browser back can feel broken or surprising
- bugs become harder to diagnose because the current screen is not addressable

### 3. Share modal and reader-link generation

The current modal in `frontend/src/author/InviteDialog.tsx` presents two paths:

- a generic share link
- a reader email flow

However, several details are inconsistent:

- the share link is local-only state and begins empty every time the modal opens
- the dialog does not receive the article’s current share URL as a prop
- the author must click Generate before seeing any generic link
- the email label says `optional`, but the submit button is disabled unless the email looks valid
- there is no UI for “tracked reader without email”

The backend model is also narrower than the UI language suggests.

In `pkg/reviewlinks/service.go`:

- `CreateInvite` lowercases and trims `Email`
- it validates the email using `mail.ParseAddress`
- therefore an invite currently requires a real email address

In `pkg/reviewlinks/postgres.go`:

- generic share links are stored on the `articles` table as `share_token`
- person-specific invites are stored in `reader_invites`

This means the current domain model supports exactly two issuance modes:

1. one shared article link
2. one invite row per email

There is no third mode for:

- name-only tracked reader links
- anonymous-but-precreated tracked links
- batches of unique reader links without requiring email addresses

### 4. Reader identity and attribution

This part is subtle and worth explaining carefully.

The review system already has a useful identity fallback for generic links.

For invited readers:

- `pkg/reviewlinks/postgres.go` resolves the token to a specific `InviteID`
- `pkg/reviews/postgres.go` reuses that `InviteID` as the reader identity

For generic share links:

- `pkg/reviewlinks/postgres.go` resolves the token to a synthetic reader identity such as `share:<article-id>`
- `pkg/reviews/postgres.go` creates a new review session on start
- later persistence falls back to the `sessionID` when there is no invite

The identity rule is approximately:

```text
if invite exists:
    reviewer identity = invite ID
else:
    reviewer identity = review session ID
```

That means the backend already knows how to distinguish multiple readers coming from a generic share link after they start their own sessions. The missing piece is product modeling, not core persistence.

What is missing today:

- an author-facing concept of a tracked link that is not email-backed
- a UI that explains how generic sessions are attributed
- a first-class backend type for a named tracked invite without email

### 5. Preview review mode

The author preview path lives in `frontend/src/app/AuthorApp.tsx` and renders `ReaderPage` with `readOnly`.

In `frontend/src/reader/ReaderPage.tsx`:

- preview mode shows a banner that says reactions are disabled
- `SectionView` receives `readOnly`
- review-session startup and review persistence are skipped because `reviewToken` is absent and the preview path is intentionally read-only

This is useful for purely visual previewing, but it blocks an important author use case:

- “I want to test the exact review flow and see my fake reactions appear in the system.”

So the current preview mode is really a rendering preview, not a workflow preview.

## Related Existing Tickets

Three adjacent tickets already overlap with parts of this work:

- `DR-005` covers author UX flow and functional gaps
- `DR-006` covers a React Router migration for author views
- `DR-007` covers wiring unused Redux slices

DR-008 should not replace those tickets. Instead, it should serve as the consolidation and implementation guide for the specific user-visible workflow gaps described in this document. The main overlap is:

- DR-006 is directly relevant to URL routing
- DR-005 overlaps with modal and preview UX bugs
- DR-007 is related but lower-level, because Redux wiring alone does not solve browser navigation or sharing semantics

## Root Cause Analysis

### Root cause 1: the frontend API contract got ahead of the backend

The delete flow is the clearest example:

- UI exists
- RTK Query endpoint exists
- backend route and service method do not exist

This suggests the frontend was designed assuming the next backend slice would arrive, but the contract was never completed end to end.

### Root cause 2: the author shell still uses a prototype-style navigation model

The app started as a desktop-like view switcher. That was fast for design iteration, but now it fights:

- browser navigation
- refresh stability
- sharable URLs
- route-aware testing

The current Redux state is useful, but it should no longer be the primary source of navigation truth.

### Root cause 3: the share model is under-specified

The product language wants several reader-link concepts, but the backend only models two of them. The modal tries to compensate with wording, which creates contradictory UX:

- “email optional”
- but email required for tracked invite

This is a sign that the domain model and UI language have diverged.

### Root cause 4: preview mode is not explicit about its purpose

There are really two kinds of preview:

1. render preview: “show me the reader UI”
2. workflow preview: “let me simulate a reader and persist test feedback”

The current implementation only supports the first kind, but the product need increasingly points to the second.

## Recommended Product Model

### Reader-link types

I recommend explicitly defining three link types:

1. `shared_link`
   - one reusable article-level token
   - good for Slack, group chats, or “anyone can review this”
   - each opening creates its own review session
   - attribution is session-based unless the reader self-identifies

2. `reader_invite`
   - one dedicated token intended for a known recipient
   - may have:
     - email
     - display name
     - both
     - neither, if the author just wants a single-use labeled slot
   - attribution is invite-based

3. `preview_invite`
   - clearly marked non-production invite for author testing
   - persists fake review data but is isolated in UI and analytics

This naming matters because it gives the product and the codebase the same mental model.

### Recommended data model adjustment

Instead of forcing `reader_invites.email` to represent all tracked invite identity, allow tracked invites to use:

- `display_name` nullable
- `email` nullable
- `identity_mode` enum-like string

Suggested modes:

- `email`
- `named`
- `anonymous`
- `preview`

Pseudocode:

```text
reader_invites
    id
    article_id
    invite_token
    identity_mode
    display_name nullable
    email nullable
    invite_note nullable
    is_preview boolean
```

This preserves today’s email invites while making the model honest about future capabilities.

## Recommended Routing Model

### Target author routes

I recommend moving the author shell to explicit nested routes such as:

- `/`
  - dashboard
- `/articles`
  - article manager list
- `/articles/:articleId`
  - author reader view
- `/articles/:articleId/edit`
  - editor
- `/articles/:articleId/settings`
  - settings
- `/articles/:articleId/preview`
  - workflow preview
- `/articles/:articleId/share`
  - share dialog route, if desired

Benefits:

- refresh works
- browser back/forward works
- author views become addressable
- tests can navigate directly to a target screen
- selected article state no longer depends on ephemeral Redux state

### Navigation ownership

After the route migration:

- the URL should be the source of truth for current author screen and selected article
- Redux can still hold local UI state such as active modal, optimistic draft content, and preview session data
- `uiSlice.view` should eventually disappear or become derived compatibility state

### Route diagram

```text
/
├── dashboard
├── articles
│   ├── :articleId
│   │   ├── review
│   │   ├── edit
│   │   ├── settings
│   │   ├── preview
│   │   └── share
└── r/:token
```

## Detailed Implementation Plan

### Phase 1: finish the broken contract surfaces

This phase should land first because it fixes genuine broken behavior.

Tasks:

- implement backend article deletion
- make the share modal internally consistent
- expose the current article share link when the modal opens

#### 1A. Backend delete article

Code changes:

- add `DeleteArticle(ctx, ownerUserID, id string) error` to `pkg/articles/service.go`
- implement repository delete in `pkg/articles/postgres.go`
- register `DELETE /api/articles/{id}` in `pkg/server/http.go`
- add handler tests in `pkg/server/http_test.go`

Important design choice:

- prefer hard delete only if foreign keys already cascade correctly
- otherwise either:
  - add `ON DELETE CASCADE` relationships
  - or implement explicit transactional delete order

Pseudocode:

```text
handler DELETE /api/articles/{id}:
    author = requireCurrentAuthor()
    articleService.DeleteArticle(ctx, author, id)
    if not found:
        return 404
    return 204
```

Validation checklist:

- deleting own article returns `204 No Content`
- deleting missing article returns `404`
- deleting another author’s article returns `404` or `403`, depending on policy
- article list is invalidated in RTK Query and dashboard updates immediately

#### 1B. Fix share modal UX

Code changes:

- change `InviteDialog` props to accept:
  - `shareUrl`
  - `onRefreshShareLink`
  - `onCreateTrackedInvite`
- remove the contradictory `optional` label if email is still required
- or better, add explicit link-type options

Recommended immediate product copy:

- `Reusable article link`
- `Tracked reader link`

Suggested modal flow:

```text
[Reusable article link]
existing share URL shown immediately
[Copy]
[Refresh]

[Tracked reader link]
Reader name (optional)
Reader email (optional)
Personal note (optional)
[Create tracked link]
```

Intern note:

The important change is not visual polish. It is making the modal state correspond to real backend capabilities. If the backend still requires an email in Phase 1, then the UI must say that plainly. Do not keep “optional” text if the submit button rejects empty email.

### Phase 2: introduce tracked non-email invites

This phase addresses the bigger product need around attribution.

#### 2A. Expand the review-link domain model

Backend changes:

- extend `reviewlinks.InviteInput`
- update `reviewlinks.Service.CreateInvite`
- update `reviewlinks.PostgresRepository.CreateInvite`
- add a DB migration for nullable email plus nullable display name, or introduce a new invite identity table

Suggested API request:

```json
{
  "email": "reader@example.com",
  "displayName": "Friend from writing group",
  "note": "Tear this apart",
  "identityMode": "email"
}
```

Allowed request shapes:

- email only
- display name only
- email plus display name
- neither only if product intentionally supports unlabeled single-use links

Validation logic:

```text
if identityMode == email:
    require valid email
if identityMode == named:
    require displayName
if identityMode == anonymous:
    require neither email nor displayName
if identityMode == preview:
    mark invite as preview/test data
```

#### 2B. Update the share modal to expose link types explicitly

Recommended UI sections:

- `Reusable link`
- `Tracked link for a person`
- `Tracked anonymous link`

This avoids overloaded wording and makes it clear why one path creates a reader record while the other does not.

### Phase 3: route migration for the author shell

This is a larger refactor, but it is necessary.

#### 3A. Route design

Update `frontend/src/App.tsx` so the author area has nested routes under a small `AuthorLayout`.

Possible structure:

- `AuthorApp` becomes route orchestration
- `DashboardPage`, `ArticleListPage`, `ArticleEditPage`, and `ArticleSettingsPage` become route elements

#### 3B. Remove route-shaped state from Redux

Current route-shaped state:

- `view`
- `selectedArticleId`
- possibly `focusSectionId`

Recommended future state split:

- URL owns `articleId` and screen identity
- Redux owns only transient UI concerns

Migration pseudocode:

```text
before:
    dispatch(selectArticle(id))
    dispatch(setView("edit"))

after:
    navigate(`/articles/${id}/edit`)
```

#### 3C. Back/forward behavior

After route migration:

- the browser handles back/forward by default
- the custom `goBack()` reducer can be removed or reduced to local editor behavior only

### Phase 4: preview workflow mode

There are two viable designs.

#### Option A: preview uses a special preview invite

Flow:

- author opens preview
- app creates or reuses a preview invite/session for that article/version
- `ReaderPage` runs the real review flow
- preview-tagged reactions are visibly separated from real reader feedback

Benefits:

- uses the real system
- lowest conceptual drift
- best for QA and demos

Costs:

- needs a preview/test-data concept in analytics and reader lists

#### Option B: preview writes fake local reviews only

Flow:

- preview mode persists to frontend-only draft state
- no backend writes occur

Benefits:

- safer
- no database changes

Costs:

- duplicates review logic
- not useful for end-to-end testing
- creates two parallel reader systems

Recommendation:

Choose Option A. The system already has real session logic. Extending it with preview-tagged sessions is cleaner than building a second fake path.

Implementation note after execution:

- For the current pass, the codebase intentionally chose a narrower version of preview reviews: preview reactions are now interactive but local-only.
- Reason: editor preview can render unsaved draft content, so backend persistence would risk attaching feedback to the wrong saved article/version.
- Revisit full persisted preview sessions only after the preview flow is explicitly tied to a saved snapshot or saved draft checkpoint.

## API Reference Proposal

### Existing endpoints

- `GET /api/articles`
- `POST /api/articles`
- `GET /api/articles/{id}`
- `PATCH /api/articles/{id}`
- `POST /api/articles/{id}/share-token`
- `POST /api/articles/{id}/invite`
- `GET /api/r/{token}`
- `POST /api/r/{token}/start`
- `POST /api/reviews/{sessionId}/progress`
- `POST /api/reviews/{sessionId}/reactions`
- `POST /api/reviews/{sessionId}/summary`

### Proposed additions or changes

- `DELETE /api/articles/{id}`
- broaden `POST /api/articles/{id}/invite` request body
- optionally add `POST /api/articles/{id}/preview-invite`
- optionally add `GET /api/articles/{id}/share-link` if the frontend should fetch rather than infer the current reusable link

Example revised tracked-link request:

```json
{
  "identityMode": "named",
  "displayName": "Workshop Reader 1",
  "email": "",
  "note": "Focus on the ending"
}
```

Example response:

```json
{
  "id": "8f3a...",
  "name": "Workshop Reader 1",
  "email": "",
  "token": "invite-...",
  "inviteUrl": "/r/invite-...",
  "identityMode": "named",
  "isPreview": false
}
```

## Testing Strategy

### Backend tests

- article delete service tests
- article delete HTTP tests
- invite validation tests for new identity modes
- review-session tests for preview invites and generic share sessions

### Frontend tests

- modal rendering with an existing share URL
- tracked invite form validation for each identity mode
- author route navigation and back button behavior
- preview-review flow with persisted test session

### Manual smoke test script

```text
1. Login as author
2. Create article
3. Navigate directly to /articles/{id}/edit
4. Reload page
5. Confirm the editor still loads the same article
6. Open share dialog
7. Copy reusable link
8. Create tracked link without email if supported
9. Open tracked link in separate browser profile
10. Leave reactions
11. Return to author dashboard and confirm attribution
12. Delete article
13. Confirm it disappears from list and related links stop resolving
```

## File-by-File Work Plan

### Frontend

- `frontend/src/App.tsx`
  - add nested author routes
- `frontend/src/app/AuthorApp.tsx`
  - replace route-like Redux navigation with `navigate()`
  - decouple modal state from hidden article selection assumptions
- `frontend/src/store/uiSlice.ts`
  - remove or shrink `view` ownership
- `frontend/src/author/InviteDialog.tsx`
  - redesign dialog inputs and prop contract
- `frontend/src/author/ArticleSettings.tsx`
  - keep delete UI, add better loading/error behavior after delete
- `frontend/src/reader/ReaderPage.tsx`
  - support persisted preview mode via a preview session flag
- `frontend/src/api/articleApi.ts`
  - align delete and invite contracts with backend

### Backend

- `pkg/server/http.go`
  - add delete route
  - update invite request parsing
- `pkg/articles/service.go`
  - add delete operation
- `pkg/articles/postgres.go`
  - implement transactional delete
- `pkg/reviewlinks/types.go`
  - expand invite identity model
- `pkg/reviewlinks/service.go`
  - validate new tracked-link modes
- `pkg/reviewlinks/postgres.go`
  - persist new invite identity data
- `pkg/reviews/postgres.go`
  - support preview/test session semantics cleanly

## Open Questions

- Should delete be hard delete, soft delete, or archive-first?
- Should a “tracked anonymous link” be single-use or multi-use?
- Should preview/test reviews appear in normal analytics by default?
- Should a reusable share link always be stable, or should regeneration revoke the previous one?
- Is “reader identity” meant to represent a person, an invite slot, or a device/browser session?

## Recommended Execution Order

If an intern were implementing this ticket, I would advise this order:

1. Finish delete end to end.
2. Fix the share modal wording and load the current reusable link immediately.
3. Decide the tracked-link product model with product/author input.
4. Extend backend invite identity support.
5. Update the modal to support new tracked-link types.
6. Migrate author navigation to real routes.
7. Add preview-session support.
8. Add tests and run full manual smoke coverage.

This order keeps broken existing flows ahead of architecture refactors, while still moving toward the cleaner route and identity model.

## Final Recommendation

Do not treat the current issues as isolated UI polish tasks. The correct next step is to align three layers together:

- route model
- share/invite domain model
- review-session identity model

If those are aligned, the remaining UX work becomes straightforward. If they are not aligned, the app will keep accumulating “almost works” screens where the copy, the client contract, and the backend behavior all imply slightly different rules.
