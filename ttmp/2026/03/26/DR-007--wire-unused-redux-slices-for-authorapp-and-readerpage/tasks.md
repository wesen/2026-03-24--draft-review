# Tasks

## WI-1: Align uiSlice with AuthorApp (blocks WI-2)

- [ ] Update uiSlice View type to match AuthorApp: "dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"
- [ ] Add focusSectionId: string | null to UiState + setFocusSection action
- [ ] Add previewArticle: Article | null to UiState + setPreviewArticle action (for unsaved draft preview)
- [ ] Add goBack action that resets view/selectedArticleId/focusSectionId/previewArticle
- [ ] Export View type for consumers

## WI-2: Replace AuthorApp useState with Redux (22 call sites)

- [ ] Import useAppDispatch, useAppSelector and uiSlice actions into AuthorApp
- [ ] Replace useState(view) with useAppSelector(s => s.ui.view) + dispatch(setView(...))
- [ ] Replace useState(selectedArticle) with useAppSelector(s => s.ui.selectedArticleId) + hydrate from articles array
- [ ] Replace useState(focusSection) with useAppSelector(s => s.ui.focusSectionId) + dispatch(setFocusSection(...))
- [ ] Replace useState(showInvite) with useAppSelector(s => s.ui.activeModal) === "invite" + dispatch(openModal/closeModal)
- [ ] Remove selectArticle useCallback — replace with inline dispatch(selectArticleAction(id))
- [ ] Update goBack() to dispatch(goBackAction()) + setShareUrl(undefined)
- [ ] Update handleNewArticle to dispatch selectArticle + setView
- [ ] Update onPreview callback to dispatch setPreviewArticle + setView("reader-preview")
- [ ] Update reader-preview render block to read from previewArticle instead of selectedArticle
- [ ] Remove local View type definition (import from uiSlice)
- [ ] Verify: shareUrl stays as local useState (transient, per-article)

## WI-3: Add sessionId to readerSlice (blocks WI-4)

- [ ] Add sessionId: string | null to ReaderState (initial: null)
- [ ] Add setSessionId action
- [ ] Export setSessionId

## WI-4: Replace ReaderPage useState with Redux (8 call sites)

- [ ] Import useAppDispatch, useAppSelector and readerSlice actions into ReaderPage
- [ ] Replace useState(started) with useAppSelector(s => s.reader.started)
- [ ] Replace useState(sessionId) with useAppSelector(s => s.reader.sessionId)
- [ ] Replace useState(currentSectionId) with useAppSelector(s => s.reader.currentSectionId) ?? fallback
- [ ] Replace useState(readSectionIds) with useAppSelector(s => s.reader.readSectionIds)
- [ ] Update goTo() to dispatch(goToSection(id))
- [ ] Update markRead() to dispatch(markSectionRead(sectionId))
- [ ] Update WelcomeSplash onStart to dispatch startReading + setSessionId
- [ ] Add useEffect cleanup: dispatch(resetReader()) on unmount
- [ ] Verify: reactions and showDone stay as local useState

## WI-5: Wire MenuBar openMenu to uiSlice (optional, low priority)

- [ ] Replace MenuBar useState(openMenu) with useAppSelector(s => s.ui.activeMenu) + dispatch(setActiveMenu)

## Verification

- [ ] TypeScript compiles clean (npx tsc --noEmit)
- [ ] Manual test: all view transitions in AuthorApp
- [ ] Manual test: preview round-trip shows unsaved draft
- [ ] Manual test: reader session lifecycle (start → navigate → finish → cleanup)
