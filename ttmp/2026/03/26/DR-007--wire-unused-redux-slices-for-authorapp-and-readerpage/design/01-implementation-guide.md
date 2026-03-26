---
Title: Implementation Guide
Ticket: DR-007
Status: active
Topics:
    - frontend
    - react
    - redux
    - state-management
DocType: design
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/store/uiSlice.ts
      Note: "Target slice — setView, selectArticle, openModal, closeModal, setActiveMenu all unused"
    - Path: frontend/src/store/readerSlice.ts
      Note: "Target slice — startReading, goToSection, markSectionRead, resetReader all unused"
    - Path: frontend/src/store/hooks.ts
      Note: "useAppDispatch, useAppSelector typed hooks — ready to use"
    - Path: frontend/src/store/store.ts
      Note: "Store already configured with ui + reader reducers"
    - Path: frontend/src/app/AuthorApp.tsx
      Note: "4 useState to move: view, selectedArticle, focusSection, showInvite"
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: "4 useState to move: started, sessionId, currentSectionId, readSectionIds"
    - Path: frontend/src/chrome/MenuBar.tsx
      Note: "1 useState (openMenu) — optional migration to uiSlice.setActiveMenu"
ExternalSources: []
Summary: "Wire the pre-built but completely unused uiSlice and readerSlice Redux state to replace useState in AuthorApp and ReaderPage."
LastUpdated: 2026-03-26T17:30:00-04:00
WhatFor: "Guide the migration of local state to Redux for app-level and reader session state."
WhenToUse: "When implementing DR-007. Work items are ordered by dependency."
---

# DR-007 Implementation Guide

## Overview

The app has two fully-built Redux slices — `uiSlice` and `readerSlice` — wired into the store with typed hooks ready. Neither is used anywhere. Meanwhile, `AuthorApp.tsx` and `ReaderPage.tsx` manage the same state with `useState`, causing:

- No state persistence across component remounts
- Prop drilling of view/selectedArticle through the component tree
- goBack() always resetting to dashboard instead of remembering context
- Reader progress lost if component unmounts

This ticket wires the existing slices to replace local state where it represents shared/app-level concerns.

---

## What exists today (unused)

### uiSlice (`frontend/src/store/uiSlice.ts`)

```ts
interface UiState {
  view: View;                    // "dashboard" | "articles" | "article-edit" | "article-settings" | "article-review"
  selectedArticleId: string | null;
  activeModal: string | null;
  activeMenu: string | null;
}
```

Actions: `setView`, `selectArticle`, `openModal`, `closeModal`, `setActiveMenu`

### readerSlice (`frontend/src/store/readerSlice.ts`)

```ts
interface ReaderState {
  started: boolean;
  currentSectionId: string | null;
  readSectionIds: string[];
}
```

Actions: `startReading`, `goToSection`, `markSectionRead`, `resetReader`

### Typed hooks (`frontend/src/store/hooks.ts`)

```ts
export const useAppDispatch = useDispatch.withTypes<AppDispatch>();
export const useAppSelector = useSelector.withTypes<RootState>();
```

---

## Work Item 1: Align uiSlice View type with AuthorApp

### Problem

uiSlice View type: `"dashboard" | "articles" | "article-edit" | "article-settings" | "article-review"`
AuthorApp View type: `"dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"`

These don't match. AuthorApp has 6 values, uiSlice has 5 and uses different names.

### Fix

Update `uiSlice.ts` View type to match AuthorApp's actual values:

```ts
type View =
  | "dashboard"
  | "articles"
  | "article"        // was "article-review" — author's reaction review view
  | "edit"            // was "article-edit"
  | "settings"        // was "article-settings"
  | "reader-preview"; // NEW — preview as reader

// Also export the type for consumers
export type { View };
```

Add `focusSectionId` to state (needed for AuthorApp line 40):

```ts
interface UiState {
  view: View;
  selectedArticleId: string | null;
  focusSectionId: string | null;      // NEW
  activeModal: string | null;
  activeMenu: string | null;
}
```

Add new action:

```ts
setFocusSection(state, action: PayloadAction<string | null>) {
  state.focusSectionId = action.payload;
},
```

Add a compound `goBack` action that resets view + selection:

```ts
goBack(state) {
  state.view = "dashboard";
  state.selectedArticleId = null;
  state.focusSectionId = null;
},
```

### Export updates

```ts
export const { setView, selectArticle, setFocusSection, goBack, openModal, closeModal, setActiveMenu } =
  uiSlice.actions;
```

---

## Work Item 2: Replace AuthorApp useState with Redux

### Current state (lines 38-42)

```tsx
const [view, setView] = useState<View>("dashboard");
const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
const [focusSection, setFocusSection] = useState<string | undefined>();
const [showInvite, setShowInvite] = useState(false);
const [shareUrl, setShareUrl] = useState<string | undefined>();  // KEEP LOCAL
```

### Replacement

```tsx
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  setView,
  selectArticle as selectArticleAction,
  setFocusSection,
  goBack as goBackAction,
  openModal,
  closeModal,
} from "../store/uiSlice";

// ...inside component:
const dispatch = useAppDispatch();
const view = useAppSelector((s) => s.ui.view);
const selectedArticleId = useAppSelector((s) => s.ui.selectedArticleId);
const focusSection = useAppSelector((s) => s.ui.focusSectionId);
const activeModal = useAppSelector((s) => s.ui.activeModal);

// Hydrate full Article object from RTK Query cache
const selectedArticle = articles.find((a) => a.id === selectedArticleId) || null;
const showInvite = activeModal === "invite";

// Keep shareUrl local — it's transient per-article
const [shareUrl, setShareUrl] = useState<string | undefined>();
```

### Call site migration (22 sites in AuthorApp)

Each `setView("x")` becomes `dispatch(setView("x"))`, etc. Full list:

| Line | Current | Replacement |
|------|---------|-------------|
| 80 | `setSelectedArticle(a)` | `dispatch(selectArticleAction(a.id))` |
| 88 | `setFocusSection(sectionId)` | `dispatch(setFocusSection(sectionId ?? null))` |
| 89 | `setView("article")` | `dispatch(setView("article"))` |
| 95-97 | `goBack()` body | replace with `dispatch(goBackAction())` + `setShareUrl(undefined)` |
| 118 | `setSelectedArticle(article)` | `dispatch(selectArticleAction(article.id))` |
| 119 | `setFocusSection(undefined)` | `dispatch(setFocusSection(null))` |
| 120 | `setView("edit")` | `dispatch(setView("edit"))` |
| 144 | `setView("dashboard")` | `dispatch(setView("dashboard"))` |
| 145 | `setView("articles")` | `dispatch(setView("articles"))` |
| 289-290 | `selectArticle(id); setView("edit")` | `dispatch(selectArticleAction(id)); dispatch(setView("edit"))` |
| 293-295 | `selectArticle(id); setShareUrl(...); setView("settings")` | `dispatch(selectArticleAction(id)); setShareUrl(undefined); dispatch(setView("settings"))` |
| 298 | `setView("articles")` | `dispatch(setView("articles"))` |
| 299 | `setShowInvite(true)` | `dispatch(openModal("invite"))` |
| 311-313 | `selectArticle(id); setView("edit")` | `dispatch(selectArticleAction(id)); dispatch(setView("edit"))` |
| 316-318 | `selectArticle(id); setShareUrl(...); setView("settings")` | same pattern |
| 322 | `setShowInvite(true)` | `dispatch(openModal("invite"))` |
| 365-366 | `setSelectedArticle(a); setView("reader-preview")` | `dispatch(selectArticleAction(a.id)); dispatch(setView("reader-preview"))` |
| 419 | `setView("edit")` | `dispatch(setView("edit"))` |
| 426 | `setShowInvite(false)` | `dispatch(closeModal())` |

### selectArticle helper

The current `selectArticle` callback (lines 77-83) finds the article in the array and sets it. With Redux storing only the ID, this simplifies to just dispatching:

```tsx
// BEFORE
const selectArticle = useCallback(
  (id: string) => {
    const a = articles.find((x) => x.id === id);
    if (a) setSelectedArticle(a);
  },
  [articles]
);

// AFTER — no callback needed, just dispatch inline
// dispatch(selectArticleAction(id))
```

The full `Article` object is derived: `const selectedArticle = articles.find(a => a.id === selectedArticleId) || null;`

### handleNewArticle

Lines 112-121 create an article, then select + navigate. After migration:

```tsx
const handleNewArticle = useCallback(async () => {
  const article = await createArticle({
    title: "Untitled Article",
    author: me?.displayName || me?.preferredUsername || me?.email || "You",
    intro: "",
  }).unwrap();
  dispatch(selectArticleAction(article.id));
  dispatch(setFocusSection(null));
  dispatch(setView("edit"));
}, [createArticle, dispatch, me?.displayName, me?.email, me?.preferredUsername]);
```

### goBack

Lines 94-98 become:

```tsx
const goBack = () => {
  dispatch(goBackAction());
  setShareUrl(undefined);
};
```

### Preview mode caveat

Line 365: `setSelectedArticle(a)` passes the *in-memory edited article* (from ArticleEditor's local state) for preview. With Redux storing only the ID, the preview would show the *last-saved* version from RTK Query cache, not the unsaved draft.

**Fix:** For preview, keep the in-memory article in a separate piece of state. Options:

**Option A** (recommended): Add `previewArticle: Article | null` to uiSlice. The preview view reads from `previewArticle` instead of hydrating from cache.

```ts
// uiSlice addition:
previewArticle: null as Article | null,

// New action:
setPreviewArticle(state, action: PayloadAction<Article | null>) {
  state.previewArticle = action.payload;
},
```

Then the preview render block reads from `previewArticle`:
```tsx
{view === "reader-preview" && previewArticle && (
  <ReaderPage article={previewArticle} readOnly ... />
)}
```

And the editor's preview button dispatches:
```tsx
onPreview={(a) => {
  dispatch(setPreviewArticle(a));
  dispatch(setView("reader-preview"));
}}
```

**Option B**: Keep `selectedArticle` as a local override that takes precedence when non-null. This preserves current behavior but mixes local + Redux state.

---

## Work Item 3: Add sessionId to readerSlice

### Current state

`readerSlice` has `started`, `currentSectionId`, `readSectionIds` — but no `sessionId`. `ReaderPage` line 42 manages sessionId locally.

### Fix

```ts
interface ReaderState {
  started: boolean;
  sessionId: string | null;       // NEW
  currentSectionId: string | null;
  readSectionIds: string[];
}

const initialState: ReaderState = {
  started: false,
  sessionId: null,               // NEW
  currentSectionId: null,
  readSectionIds: [],
};
```

Add action:

```ts
setSessionId(state, action: PayloadAction<string | null>) {
  state.sessionId = action.payload;
},
```

Update `resetReader` to also clear sessionId (already handled by returning `initialState`).

Export:

```ts
export const { startReading, setSessionId, goToSection, markSectionRead, resetReader } =
  readerSlice.actions;
```

---

## Work Item 4: Replace ReaderPage useState with Redux

### Current state (lines 41-46)

```tsx
const [started, setStarted] = useState(false);
const [sessionId, setSessionId] = useState<string | null>(null);
const [currentSectionId, setCurrentSectionId] = useState(article.sections[0]?.id);
const [readSectionIds, setReadSectionIds] = useState<string[]>([]);
const [reactions, setReactions] = useState<Reaction[]>(initialReactions);   // KEEP LOCAL
const [showDone, setShowDone] = useState(false);                            // KEEP LOCAL
```

### Replacement

```tsx
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  startReading,
  setSessionId as setSessionIdAction,
  goToSection,
  markSectionRead,
  resetReader,
} from "../store/readerSlice";

// ...inside component:
const dispatch = useAppDispatch();
const started = useAppSelector((s) => s.reader.started);
const sessionId = useAppSelector((s) => s.reader.sessionId);
const currentSectionId = useAppSelector((s) => s.reader.currentSectionId) ?? article.sections[0]?.id;
const readSectionIds = useAppSelector((s) => s.reader.readSectionIds);

// Keep local:
const [reactions, setReactions] = useState<Reaction[]>(initialReactions);
const [showDone, setShowDone] = useState(false);
```

### Call site migration (8 sites)

| Line | Current | Replacement |
|------|---------|-------------|
| 63 | `setCurrentSectionId(id)` | `dispatch(goToSection(id))` |
| 95-98 | `setReadSectionIds(prev => ...)` | `dispatch(markSectionRead(sectionId))` |
| 205 | `setSessionId(result.session.id)` | `dispatch(setSessionIdAction(result.session.id))` |
| 210 | `setStarted(true)` | (handled by `startReading` below) |
| 211 | `setCurrentSectionId(article.sections[0]?.id)` | (handled by `startReading` below) |
| 210-211 combined | `setStarted(true); setCurrentSectionId(...)` | `dispatch(startReading(article.sections[0]?.id))` |

### WelcomeSplash onStart callback

Lines 201-213 become:

```tsx
onStart={async () => {
  if (reviewToken && !sessionId) {
    try {
      const result = await startReview({ token: reviewToken }).unwrap();
      dispatch(setSessionIdAction(result.session.id));
    } catch (error) {
      console.error("Failed to start review session", error);
    }
  }
  dispatch(startReading(article.sections[0]?.id));
}}
```

### Reset on unmount

Add cleanup effect to dispatch `resetReader()` when the component unmounts:

```tsx
useEffect(() => {
  return () => {
    dispatch(resetReader());
  };
}, [dispatch]);
```

This ensures the reader state doesn't leak between different articles/sessions.

### markRead simplification

The current `markRead` callback (lines 93-103) manually deduplicates in local state. The Redux `markSectionRead` action already handles deduplication (`if (!state.readSectionIds.includes(...))`), so:

```tsx
// BEFORE
const markRead = useCallback(
  (sectionId: string) => {
    setReadSectionIds((prev) =>
      prev.includes(sectionId) ? prev : [...prev, sectionId]
    );
    void persistProgress(sectionId);
  },
  [persistProgress]
);

// AFTER
const markRead = useCallback(
  (sectionId: string) => {
    dispatch(markSectionRead(sectionId));
    void persistProgress(sectionId);
  },
  [dispatch, persistProgress]
);
```

### persistProgress adjustment

`persistProgress` (lines 67-91) reads `readSectionIds` and `sessionId` from local state. After migration, it reads from Redux via the selector values already destructured at the top of the component. No change needed to the function body — the closure will capture the selector-derived values.

---

## Work Item 5: Wire MenuBar openMenu to uiSlice (optional)

### Current

`MenuBar.tsx` line 43: `const [openMenu, setOpenMenu] = useState<string | null>(null);`

`uiSlice` already has `activeMenu: string | null` and `setActiveMenu(string | null)`.

### Benefit

Low. MenuBar's dropdown state is transient and fully component-scoped. The only benefit is if other components need to programmatically open/close menus (e.g., keyboard shortcuts).

### If migrating

```tsx
// MenuBar.tsx
const dispatch = useAppDispatch();
const openMenu = useAppSelector((s) => s.ui.activeMenu);
// Replace setOpenMenu(x) with dispatch(setActiveMenu(x))
```

### Recommendation

**Skip this.** Keep local. Only migrate if a future feature needs cross-component menu control.

---

## What stays local (no migration)

| Component | State vars | Reason |
|-----------|-----------|--------|
| ArticleEditor | `article`, `activeSectionId`, `contentDraft`, `showDeleteConfirm`, `showUnsavedConfirm` | Editor draft — component-scoped working state. Only persisted on explicit Save. |
| ArticleSettings | `title`, `intro`, `version`, `status`, `accessMode`, 6 toggle states, `copied`, `showDeleteConfirm` | Form state — collected and submitted together. Dies on unmount. |
| Dashboard | `selectedId` | Tab selection within dashboard. Not shared with siblings. |
| ArticleManager | `search`, `sortKey`, `expanded` | Filter/sort within article list view. |
| ArticleReader | `selectedSection`, `reactionFilter`, `hoveredParagraph`, `filteredParagraph` | View filters within review component. |
| InviteDialog | `email`, `note`, `inviteResult`, `isSending`, `isGeneratingShareLink`, `currentShareUrl`, `shareCopied`, `inviteCopied`, `error`, `shareError` | Modal form state — dies when dialog closes. |
| Paragraph | `activeType`, `comment` | Reaction chip interaction per paragraph. |
| ReactionPicker | `chosenType`, `comment` | Selection state per picker. |
| ReaderPage | `reactions`, `showDone` | Reactions: optimistic local updates with API side-effect. showDone: dialog toggle. |
| MenuBar | `openMenu` | Transient dropdown state (optional migration in WI-5). |

---

## Priority order

| # | Item | Impact | Effort | Notes |
|---|------|--------|--------|-------|
| 1 | WI-1: Align uiSlice View type | Blocking | Low | Must fix before WI-2 |
| 2 | WI-2: Replace AuthorApp useState | High | Medium | 22 call sites; preview article caveat |
| 3 | WI-3: Add sessionId to readerSlice | Blocking | Trivial | Must add before WI-4 |
| 4 | WI-4: Replace ReaderPage useState | High | Low | 8 call sites + unmount cleanup |
| 5 | WI-5: Wire MenuBar openMenu | Low | Trivial | Optional — skip unless needed |

---

## Testing checklist

After migration, verify:

1. **Dashboard loads** — view defaults to "dashboard", first article auto-selected
2. **Article navigation** — clicking tabs, Edit, Settings, Review all switch views correctly
3. **Back button** — goBack returns to dashboard and clears selection
4. **New article** — creates article, selects it, navigates to editor
5. **Preview round-trip** — Editor → Preview shows unsaved draft (not cached version) → Back returns to editor
6. **Invite dialog** — opens/closes via Redux modal state
7. **Reader session** — WelcomeSplash → reading → section navigation → progress tracking → Finish all work
8. **Reader cleanup** — navigating away from reader resets Redux state (no stale section/progress)
9. **Page refresh** — Redux state resets to defaults (this is expected; URL routing in DR-006 will preserve state across refreshes)
