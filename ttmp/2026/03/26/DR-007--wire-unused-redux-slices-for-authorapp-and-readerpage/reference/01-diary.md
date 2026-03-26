---
Title: Implementation Diary
Ticket: DR-007
Status: active
Topics:
    - implementation
    - redux
    - frontend
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/store/uiSlice.ts
      Note: "Expanded: aligned View type, added focusSectionId, previewArticle, goBack"
    - Path: frontend/src/store/readerSlice.ts
      Note: "Expanded: added sessionId field + setSessionId action"
    - Path: frontend/src/app/AuthorApp.tsx
      Note: "Migrated: 4 useState → Redux (22 call sites)"
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: "Migrated: 4 useState → Redux (8 call sites) + unmount cleanup"
ExternalSources: []
Summary: "Diary of wiring pre-built but unused uiSlice and readerSlice Redux state."
LastUpdated: 2026-03-26T18:00:00-04:00
WhatFor: "Document implementation decisions and tricky parts for code review."
WhenToUse: "When reviewing DR-007 changes."
---

# DR-007 Implementation Diary

## Session: 2026-03-26

### Commit 1: `2958b2a` — WI-1 + WI-2 (uiSlice + AuthorApp)

**WI-1: Align uiSlice**

The existing uiSlice View type didn't match AuthorApp's:
- uiSlice had: `"dashboard" | "articles" | "article-edit" | "article-settings" | "article-review"` (5 values)
- AuthorApp had: `"dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"` (6 values)

Replaced with AuthorApp's values since those are what the app actually uses. Exported the `View` type so AuthorApp no longer defines its own.

Added three new pieces of state:
- `focusSectionId: string | null` — replaces `useState<string | undefined>` in AuthorApp. Changed from undefined to null for Redux serialization consistency.
- `previewArticle: Article | null` — needed because the editor passes its *unsaved draft* to preview. If we only stored the article ID and hydrated from RTK Query cache, the preview would show the last-saved version, not the draft. This is the most important design decision in the migration.
- `goBack` action — compound reducer that resets view, selectedArticleId, focusSectionId, and previewArticle in one dispatch.

**WI-2: AuthorApp migration**

22 call sites changed. Key patterns:

1. **`setView("x")` → `dispatch(setView("x"))`** — straightforward, 10 sites
2. **`selectArticle(id)` callback removed** — was `useCallback` that found article in array and called `setSelectedArticle(a)`. Now just `dispatch(selectArticleAction(id))` inline. The full Article object is derived: `const selectedArticle = articles.find(a => a.id === selectedArticleId) || null`
3. **`goBack()` simplified** — was 3 lines setting view + article + shareUrl. Now dispatches `goBackAction()` (which resets 4 Redux fields) + `setShareUrl(undefined)` for the local state.
4. **Preview mode** — `onPreview` callback now dispatches `setPreviewArticle(a)` + `setView("reader-preview")`. The render block reads from `previewArticle` instead of `selectedArticle`. This preserves the unsaved draft.
5. **Invite dialog** — `setShowInvite(true)` → `dispatch(openModal("invite"))`, close → `dispatch(closeModal())`, check → `activeModal === "invite"`.

**What stayed local:** `shareUrl` — it's transient per-article, generated on demand, never read by child components. Not worth putting in Redux.

**Tricky part:** The `focusSection` prop passed to ArticleReader was `string | undefined`. Redux stores it as `string | null`. Added `?? undefined` at the call site: `focusSection={focusSection ?? undefined}`.

---

### Commit 2: `3f15311` — WI-3 + WI-4 (readerSlice + ReaderPage)

**WI-3: Add sessionId to readerSlice**

Added `sessionId: string | null` to ReaderState and `setSessionId` action. The existing `resetReader` already returns `initialState` which now includes `sessionId: null`, so no changes needed there.

**WI-4: ReaderPage migration**

8 call sites changed. The migration was simpler than AuthorApp because:
- No naming conflicts (readerSlice action names match the local function names)
- No computed state derivation (unlike selectedArticle hydration)
- `markSectionRead` reducer already handles dedup, so the local callback that checked `prev.includes(sectionId)` was simplified

Key changes:
1. **`goTo()` callback** — now dispatches `goToSection(id)` + scrolls. The `useCallback` stays (scroll ref side effect).
2. **`markRead()` callback** — dispatches `markSectionRead(sectionId)` instead of local setState. Dedup logic removed (handled by reducer).
3. **WelcomeSplash `onStart`** — dispatches `setSessionIdAction(result.session.id)` + `startReading(sectionId)`. The `startReading` reducer sets both `started: true` and `currentSectionId` in one action.
4. **Unmount cleanup** — added `useEffect(() => { return () => dispatch(resetReader()); }, [dispatch])`. This ensures reader state doesn't leak between articles.

**What stayed local:**
- `reactions` — optimistic updates with API side-effect. Putting in Redux would require an optimistic update middleware pattern that adds complexity for no benefit here.
- `showDone` — dialog toggle, component-scoped.

**Observation:** `persistProgress` reads `readSectionIds` and `sessionId` via the closure over selector values. This works because the selector re-renders the component with new values, so the closure always captures current state. No stale closure issue.

---

## What's left

- WI-5 (MenuBar openMenu → uiSlice.setActiveMenu) — **skipped**, low priority, purely optional. MenuBar dropdown state is transient and component-scoped.

## Summary

| File | Before | After | Net change |
|------|--------|-------|------------|
| uiSlice.ts | 5 actions, 4 state fields | 8 actions, 6 state fields | +3 actions, +2 fields |
| readerSlice.ts | 4 actions, 3 state fields | 5 actions, 4 state fields | +1 action, +1 field |
| AuthorApp.tsx | 5 useState, 0 Redux | 1 useState, 4 from Redux | -4 useState |
| ReaderPage.tsx | 6 useState, 0 Redux | 2 useState, 4 from Redux | -4 useState |
