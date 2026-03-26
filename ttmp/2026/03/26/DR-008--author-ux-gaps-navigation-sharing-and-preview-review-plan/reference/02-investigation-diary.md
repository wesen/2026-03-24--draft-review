---
Title: Investigation diary
Ticket: DR-008
Status: active
Topics:
    - diary
    - frontend
    - backend
    - ux
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Chronological diary for the DR-008 investigation and documentation work.
LastUpdated: 2026-03-26T12:35:00-04:00
WhatFor: Preserve what was inspected, what was learned, and how the conclusions in this ticket were reached.
WhenToUse: Use this when reviewing the investigation history or extending the ticket later.
---

# Investigation Diary

## 2026-03-26

### Objective

Create a new docmgr ticket that analyzes the remaining author UX gaps, with special focus on:

- delete article
- proper URL routes and browser navigation
- the broken share/generate-link modal
- tracked links for readers with and without email
- optional preview-mode test reviews

Also create a separate rolling code-review document that captures subpar or fragile areas noticed during the investigation.

### Work Log

1. Created a new ticket workspace with `docmgr create ticket`.
   Result:
   - new ticket `DR-008`
   - ticket path `ttmp/2026/03/26/DR-008--author-ux-gaps-navigation-sharing-and-preview-review-plan`

2. Created the initial document scaffold with `docmgr add`.
   Notes:
   - the first add attempt for the main design doc raced and reported that the new ticket directory could not be found
   - rerunning the command worked

3. Inspected the current ticket landscape to avoid duplicating nearby work.
   Files checked:
   - `ttmp/2026/03/26/DR-005--author-ux-flow-and-functional-fixes/index.md`
   - `ttmp/2026/03/26/DR-006--react-router-migration-for-author-views/index.md`
   - `ttmp/2026/03/26/DR-007--wire-unused-redux-slices-for-authorapp-and-readerpage/index.md`
   Key finding:
   - DR-006 already tracks route migration specifically
   - DR-005 overlaps with author UX bugs
   - DR-007 overlaps with state wiring but does not solve the route and product-model problems on its own

4. Inspected the application shell and author navigation.
   Files checked:
   - `frontend/src/App.tsx`
   - `frontend/src/app/AuthorApp.tsx`
   - `frontend/src/store/uiSlice.ts`
   Key findings:
   - React Router is already mounted at the app edge
   - author navigation still depends on `uiSlice.view`
   - `goBack()` just resets to dashboard instead of using real history

5. Inspected deletion flow from UI to backend.
   Files checked:
   - `frontend/src/author/ArticleSettings.tsx`
   - `frontend/src/api/articleApi.ts`
   - `pkg/server/http.go`
   - `pkg/articles/service.go`
   Key finding:
   - the delete flow is only half-implemented
   - frontend exposes `DELETE /articles/{id}`
   - backend has no handler or service method

6. Inspected sharing and invite behavior.
   Files checked:
   - `frontend/src/author/InviteDialog.tsx`
   - `pkg/reviewlinks/types.go`
   - `pkg/reviewlinks/service.go`
   - `pkg/reviewlinks/postgres.go`
   Key findings:
   - generic share links and email invites are supported
   - the modal says email is optional but the button requires a valid email
   - the backend invite model requires an email address today
   - the current share URL is treated as modal-local state, not article state

7. Inspected review-session identity and reader preview behavior.
   Files checked:
   - `frontend/src/api/readerApi.ts`
   - `frontend/src/app/ReaderApp.tsx`
   - `frontend/src/reader/ReaderPage.tsx`
   - `frontend/src/store/readerSlice.ts`
   - `pkg/reviews/postgres.go`
   Key findings:
   - invited readers are attributed via invite identity
   - generic share readers fall back to session identity once started
   - preview mode intentionally disables persisted reactions

8. Wrote the main design document.
   Focus:
   - explain the current system for an intern
   - connect the UX bugs to the current architecture
   - propose a phased implementation plan instead of isolated patches

9. Wrote the rolling code-review document.
   Focus:
   - capture broader observations that should not be lost
   - identify likely cleanup tickets

10. Updated ticket metadata and task tracking docs.

11. Ran documentation validation.
   Command:
   - `docmgr doctor --ticket DR-008 --stale-after 30`
   Result:
   - initial warning about unregistered topic slugs `routing` and `sharing`
   - normalized the topic metadata to the repository vocabulary
   - reran the doctor successfully

12. Uploaded the ticket bundle to reMarkable.
   Commands:
   - `remarquee upload bundle --dry-run ... --name "DR-008 Author UX Analysis Bundle" --remote-dir "/ai/2026/03/26/DR-008" --toc-depth 2`
   - `remarquee upload bundle --force ... --name "DR-008 Author UX Analysis Bundle" --remote-dir "/ai/2026/03/26/DR-008" --toc-depth 2`
   - `remarquee cloud ls /ai/2026/03/26/DR-008`
   Result:
   - uploaded bundle name `DR-008 Author UX Analysis Bundle`
   - confirmed remote listing in `/ai/2026/03/26/DR-008`

### Notable Conclusions

- The current pain points are connected. The route model, the share/invite model, and the review-session identity model need to be aligned together.
- The backend already has the beginnings of session-based identity for generic links, which makes tracked non-email links feasible without redesigning the entire review system.
- The author shell has outgrown the prototype-style `view` state machine.

### Suggested Next Work

1. Implement article deletion end to end.
2. Fix the share dialog wording and make the reusable link visible immediately.
3. Decide the tracked-link product model for email, named, anonymous, and preview identities.
4. Migrate author navigation to real routes.
5. Add preview-session support if product wants real test reviews rather than visual-only preview.

### Execution Kickoff

After the investigation finished, the next execution plan was expanded into concrete slices:

1. backend delete article support
2. author route migration
3. share modal repair
4. tracked non-email link support
5. persisted preview test reviews
6. final verification and documentation refresh

The implementation diary will continue under these slices so the ticket records not just what was planned, but what was actually changed, tested, and committed.

### Slice 1: Delete Article End to End

Objective:

- repair the broken delete contract by implementing the missing backend pieces behind the existing frontend affordance

Work performed:

1. Added `DeleteArticle` to the article service repository interface in `pkg/articles/service.go`.
   Reason:
   - the service layer had create, update, and create-version operations but no delete operation, even though the frontend already exposed deletion

2. Implemented repository deletion in `pkg/articles/postgres.go`.
   Notes:
   - deletion uses the existing article ownership guard: `where id = $1 and owner_user_id = $2`
   - because the schema already uses `on delete cascade` widely, deleting the article row is the natural root operation

3. Registered and implemented `DELETE /api/articles/{id}` in `pkg/server/http.go`.
   Behavior:
   - `204 No Content` on success
   - `404` when the article is not found for the authenticated author
   - `503` if the article service is not configured

4. Added focused handler coverage in `pkg/server/http_test.go`.
   Test added:
   - `TestHandleDeleteArticleUsesAuthenticatedAuthor`
   What it verifies:
   - the handler resolves the ensured authenticated author identity
   - the delete call targets the requested article id
   - the response code is `204`

5. Tightened the frontend delete trigger in `frontend/src/app/AuthorApp.tsx`.
   Change:
   - switched from awaiting the mutation trigger directly to `unwrap()`
   Why:
   - previously the UI would call `goBack()` even if the delete failed
   - `unwrap()` ensures navigation only continues on actual mutation success

Validation:

- `go test ./pkg/server ./pkg/articles`
  Result:
  - passed

Commit plan:

- stage the delete slice together with the ticket diary/task/changelog updates
- commit before starting the routing refactor so the contract repair is isolated

### Slice 2: Route-Driven Author Navigation

Objective:

- replace the author shell’s `uiSlice.view` navigation with real URL-based navigation so browser back/forward and refresh reflect the active screen

Work performed:

1. Reworked `frontend/src/app/AuthorApp.tsx` to derive the current screen from React Router path matches instead of Redux view strings.
   New route shapes handled in the author shell:
   - `/`
   - `/articles`
   - `/articles/:articleId`
   - `/articles/:articleId/edit`
   - `/articles/:articleId/settings`
   - `/articles/:articleId/share`
   - `/articles/:articleId/preview`

2. Switched author navigation callbacks from Redux actions to `navigate(...)`.
   Examples:
   - selecting an article now navigates to `/articles/:id`
   - creating an article now navigates to `/articles/:id/edit`
   - preview now navigates to `/articles/:id/preview`
   - share now navigates to `/articles/:id/share`

3. Moved section focus from Redux into the URL query string.
   Example:
   - dashboard “review this section” links now encode `?section=...`

4. Kept only the transient preview draft in Redux for now.
   Reason:
   - previewing unsaved edits still needs a temporary in-memory article object
   - this is a smaller and safer bridge than keeping the full author screen identity in Redux

5. Updated `frontend/src/author/Dashboard.tsx` so invite/share actions carry the target article id explicitly.

6. Added a route-aware “article not found” fallback window in the author shell.

7. During verification, the frontend build failed on an unrelated Storybook type issue in `frontend/src/chrome/MacWindow.stories.tsx`.
   Problem:
   - the story passed obsolete `w` and `h` args that no longer exist on `MacWindowProps`
   Fix:
   - removed the stale args so the frontend typecheck/build pipeline is green again

Validation:

- `cd frontend && npm run build`
  Result:
  - passed
  - one existing Vite chunk-size warning remains, but it is only a warning

Commit plan:

- commit the route migration separately from the upcoming share/invite domain changes
- keep the Storybook typing fix in the same commit because it was necessary to validate the frontend slice

### Slices 3 and 4: Share Modal Repair and Invite Identity Modes

Objective:

- fix the broken share dialog and make the backend honest about the different types of tracked reader links the product now wants to support

Work performed:

1. Expanded the database invite model.
   Files:
   - `pkg/db/migrations/0001_init.sql`
   - `pkg/db/migrations/0005_reader_invite_identity_modes.sql`
   Changes:
   - `reader_invites.email` is now nullable
   - added `display_name`
   - added `identity_mode`
   - added `is_preview`

2. Added invite identity modes to the backend domain model in `pkg/reviewlinks/types.go`.
   Modes introduced:
   - `email`
   - `named`
   - `anonymous`
   - `preview`

3. Replaced the old “email required always” validation with mode-aware validation in `pkg/reviewlinks/service.go`.
   Behavior now:
   - `email` mode requires a valid email
   - `named` mode requires a display name
   - `anonymous` mode requires neither email nor display name
   - `preview` mode reserves a tracked test-reader slot and defaults the display name when omitted

4. Updated invite persistence and token resolution in `pkg/reviewlinks/postgres.go`.
   Notes:
   - invite rows now store display name, identity mode, and preview flag
   - resolved reader identity uses a helper that can derive a stable display name from display name, email, or mode fallback

5. Updated article loading to expose persistent reusable share links.
   Files:
   - `pkg/articles/types.go`
   - `pkg/articles/postgres.go`
   Result:
   - author article payloads now include `shareUrl` when a share token already exists

6. Updated analytics behavior for the richer invite model in `pkg/analytics/postgres.go`.
   Notes:
   - article reader lists now carry `identityMode` and `isPreview`
   - display names come from display name, email, or identity-mode fallback
   - the reader directory now ignores invite rows with no email, which is correct because it is a contact directory, not a universal reader ledger

7. Updated session name resolution in `pkg/reviews/postgres.go`.
   Important nuance:
   - anonymous tracked invite links now resolve to anonymous session names
   - generic article share links were intentionally left as generic session-based links rather than being collapsed into the anonymous invite mode

8. Rewrote the frontend invite/share dialog.
   Files:
   - `frontend/src/author/InviteDialog.tsx`
   - `frontend/src/author/InviteDialog.css`
   - `frontend/src/author/InviteDialog.stories.tsx`
   Changes:
   - the reusable article link is shown immediately
   - when there is no existing reusable link yet, the dialog auto-generates one on open
   - tracked reader links now have explicit mode buttons for `Email`, `Named`, and `Anonymous`
   - the form fields and submit button now match the selected mode instead of pretending email is optional at all times

9. Updated frontend types and API contracts.
   Files:
   - `frontend/src/types/article.ts`
   - `frontend/src/types/reader.ts`
   - `frontend/src/types/index.ts`
   - `frontend/src/api/articleApi.ts`
   - `frontend/src/app/AuthorApp.tsx`
   Changes:
   - articles now expose `shareUrl`
   - invite mutation payloads now include `identityMode`, optional `displayName`, and optional `email`
   - share-token mutation now invalidates article cache so the reusable link refreshes cleanly

10. Added and updated validation coverage.
    Files:
    - `pkg/reviewlinks/service_test.go`
    - `pkg/server/http_test.go`
    Tests covered:
    - named invite without email succeeds
    - anonymous invite without email succeeds
    - anonymous invite with email fails
    - display-name fallback helper behavior

Validation:

- `go test ./pkg/reviewlinks ./pkg/reviews ./pkg/analytics ./pkg/server ./pkg/articles`
  Result:
  - passed
- `cd frontend && npm run build`
  Result:
  - passed

Commit plan:

- commit the share/invite model changes separately from the preview-review work because the preview decision still has one remaining product tradeoff around persisted vs local-only behavior for unsaved editor drafts
