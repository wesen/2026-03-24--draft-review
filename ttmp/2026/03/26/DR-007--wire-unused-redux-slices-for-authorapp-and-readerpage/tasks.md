# Tasks

## WI-1: Align uiSlice with AuthorApp (blocks WI-2)

- [x] Update uiSlice View type to match AuthorApp — 2958b2a
- [x] Add focusSectionId, previewArticle fields + actions — 2958b2a
- [x] Add goBack compound action — 2958b2a
- [x] Export View type — 2958b2a

## WI-2: Replace AuthorApp useState with Redux (22 call sites)

- [x] Replace view, selectedArticle, focusSection, showInvite with Redux — 2958b2a
- [x] Remove selectArticle useCallback + local View type — 2958b2a
- [x] Preview uses previewArticle from Redux (preserves unsaved draft) — 2958b2a
- [x] shareUrl stays as local useState — 2958b2a

## WI-3: Add sessionId to readerSlice (blocks WI-4)

- [x] Add sessionId field + setSessionId action — 3f15311

## WI-4: Replace ReaderPage useState with Redux (8 call sites)

- [x] Replace started, sessionId, currentSectionId, readSectionIds with Redux — 3f15311
- [x] Add useEffect cleanup: resetReader() on unmount — 3f15311
- [x] reactions + showDone stay as local useState — 3f15311

## WI-5: Wire MenuBar openMenu to uiSlice (optional, low priority)

- [ ] Replace MenuBar useState(openMenu) with uiSlice.setActiveMenu — skipped, not needed
