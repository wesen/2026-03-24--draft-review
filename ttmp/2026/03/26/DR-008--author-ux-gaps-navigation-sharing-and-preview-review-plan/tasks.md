# Tasks

## Completed Investigation Work

- [x] Review DR-005, DR-006, and DR-007 to avoid duplicating adjacent UX and routing work
- [x] Inspect the author navigation flow in `frontend/src/App.tsx`, `frontend/src/app/AuthorApp.tsx`, and `frontend/src/store/uiSlice.ts`
- [x] Inspect the delete article flow across `ArticleSettings`, `articleApi`, `pkg/server/http.go`, and `pkg/articles/service.go`
- [x] Inspect the share modal and reader-link persistence flow across `InviteDialog`, `reviewlinks`, and `reviews`
- [x] Inspect preview-mode behavior in `ReaderPage`, `ReaderApp`, and the review session API
- [x] Write the main intern-facing analysis and implementation guide
- [x] Write a separate rolling code-review and improvement backlog document
- [x] Update the investigation diary, ticket metadata, and changelog
- [x] Validate the ticket with `docmgr doctor`
- [x] Upload the ticket bundle to reMarkable

## Proposed Implementation Sequence

- [x] Slice 1: Add backend article deletion support
  - [x] Add `DeleteArticle` to `pkg/articles/service.go`
  - [x] Implement repository delete in `pkg/articles/postgres.go`
  - [x] Register `DELETE /api/articles/{id}` in `pkg/server/http.go`
  - [x] Add focused HTTP coverage for delete behavior
  - [x] Verify the frontend delete trigger now unwraps success instead of navigating away on failure
- [x] Slice 2: Migrate author navigation to URL routes
  - [x] Define author route map for dashboard, articles, review, edit, settings, preview, and share
  - [x] Replace `uiSlice.view`-driven screen selection with React Router state in `AuthorApp`
  - [x] Keep preview drafts as transient Redux state while route ownership moves to the URL
  - [x] Make browser back/forward and reload preserve the current article screen
  - [x] Add route-aware frontend smoke coverage via production build validation
- [x] Slice 3: Repair and simplify the share link modal
  - [x] Show the reusable article link immediately when the modal opens
  - [x] Stop labeling email as optional unless that is actually true for the selected link type
  - [x] Add clearer sections for reusable link vs tracked reader link
  - [x] Add success and error surfaces for link generation and invite creation
- [x] Slice 4: Extend tracked reader links beyond email-only invites
  - [x] Add a DB migration for invite identity metadata
  - [x] Extend `reviewlinks.InviteInput` and validation rules
  - [x] Support `email`, `named`, and `anonymous` tracked link modes, with `preview` reserved for test workflows
  - [x] Update analytics and resolved reader display naming rules
  - [x] Update frontend types, article share-link state, and invite form submission payloads
- [x] Slice 5: Enable preview test reviews
  - [x] Decide that preview reactions remain local-only while previewing unsaved editor drafts
  - [x] Add an explicit preview mode banner separate from read-only mode
  - [x] Allow preview mode to submit local reactions instead of disabling interaction entirely
  - [x] Avoid backend persistence so preview feedback cannot corrupt normal production reader identity
- [x] Slice 6: Final verification and documentation refresh
  - [x] Run backend tests
  - [x] Run frontend build and relevant smoke checks
  - [x] Update DR-008 diary, changelog, and implementation notes after each slice
  - [x] Commit each major slice with a focused message
  - [x] Re-run `docmgr doctor`
  - [x] Refresh the reMarkable bundle
