# Tasks

## Phase 1: Wire uiSlice for AuthorApp

- [ ] P1-1: Align uiSlice View type with AuthorApp views (add "review", "reader-preview"; add focusSectionId field)
- [ ] P1-2: Replace AuthorApp useState(view/selectedArticle/focusSection/showInvite) with useAppSelector + dispatch
- [ ] P1-3: Update all setView/setSelectedArticle/setShowInvite call sites to dispatch Redux actions
- [ ] P1-4: Hydrate selectedArticle from RTK Query cache using selectedArticleId from Redux

## Phase 2: React Router routes

- [ ] P2-1: Define nested author routes in App.tsx (/articles, /articles/:id, /articles/:id/edit, etc.)
- [ ] P2-2: Add URL-to-Redux sync effect in AuthorApp (derive view + selectedArticleId from route params)
- [ ] P2-3: Replace all navigate-by-dispatch with React Router navigate() calls
- [ ] P2-4: Replace goBack() with navigate(-1) for proper browser history
- [ ] P2-5: Update menu bar View items to use navigate()
- [ ] P2-6: Handle 404 / invalid article ID in URL params

## Phase 3: Wire readerSlice for ReaderPage

- [ ] P3-1: Add sessionId field + setSessionId action to readerSlice
- [ ] P3-2: Replace ReaderPage useState(started/sessionId/currentSectionId/readSectionIds) with Redux
- [ ] P3-3: Dispatch resetReader() on ReaderPage unmount

## Phase 4: Cleanup

- [ ] P4-1: Remove View type and goBack() from AuthorApp (dead code after routing migration)
- [ ] P4-2: Test browser back button, page refresh, deep linking, bookmarks
- [ ] P4-3: Test preview mode round-trip (editor → preview → back preserves article state)
