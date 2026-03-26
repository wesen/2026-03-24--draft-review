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
- [ ] Slice 2: Migrate author navigation to URL routes
  - [ ] Define author route map for dashboard, articles, review, edit, settings, preview
  - [ ] Replace `uiSlice.view`-driven screen selection with React Router state in `AuthorApp`
  - [ ] Keep only truly local UI state in Redux
  - [ ] Make browser back/forward and reload preserve the current article screen
  - [ ] Add route-aware frontend tests or smoke coverage
- [ ] Slice 3: Repair and simplify the share link modal
  - [ ] Show the reusable article link immediately when the modal opens
  - [ ] Stop labeling email as optional unless that is actually true for the selected link type
  - [ ] Add clearer sections for reusable link vs tracked reader link
  - [ ] Add success and error surfaces for link generation and invite creation
- [ ] Slice 4: Extend tracked reader links beyond email-only invites
  - [ ] Add a DB migration for invite identity metadata
  - [ ] Extend `reviewlinks.InviteInput` and validation rules
  - [ ] Support at least `email`, `named`, and `anonymous` tracked link modes
  - [ ] Update analytics and resolved reader display naming rules
  - [ ] Update frontend types and invite form submission payloads
- [ ] Slice 5: Enable persisted preview test reviews
  - [ ] Add an explicit preview/test mode for reader sessions
  - [ ] Allow preview mode to start a synthetic tracked session and post reactions
  - [ ] Distinguish preview/test data in author-facing reader and reaction displays
  - [ ] Verify preview reviews do not corrupt normal production reader identity
- [ ] Slice 6: Final verification and documentation refresh
  - [ ] Run backend tests
  - [ ] Run frontend build and relevant smoke checks
  - [ ] Update DR-008 diary, changelog, and implementation notes after each slice
  - [ ] Commit each major slice with a focused message
