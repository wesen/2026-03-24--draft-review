# Tasks

## TODO

- [ ] Inspect the current reader flow in `frontend/src/reader/ReaderPage.tsx`, `frontend/src/app/ReaderApp.tsx`, `frontend/src/api/readerApi.ts`, and `pkg/reviews/postgres.go`
- [ ] Confirm which review data is already persisted immediately and which data still depends on manual completion
- [ ] Change the reader UX so review progress/comments persist without requiring a manual `Finish Review` click
- [ ] Decide whether the last-section transition should auto-complete the review session, or whether completion should happen on first persisted feedback / idle exit / page leave
- [ ] Update reader UI copy so the flow no longer implies a required manual completion action
- [ ] Add backend and frontend validation coverage for the new completion behavior
- [ ] Update the ticket diary and changelog with implementation notes
