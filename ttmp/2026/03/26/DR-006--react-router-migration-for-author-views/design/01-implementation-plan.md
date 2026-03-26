---
Title: Implementation Plan
Ticket: DR-006
Status: active
Topics:
    - frontend
    - react
    - routing
    - redux
DocType: design
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/App.tsx
      Note: "Top-level BrowserRouter + Routes — already routes /r/:token to ReaderApp, /* to AuthorApp"
    - Path: frontend/src/app/AuthorApp.tsx
      Note: "All author navigation via useState<View> — 5 state vars to replace with router + Redux"
    - Path: frontend/src/store/uiSlice.ts
      Note: "Pre-built Redux slice with setView, selectArticle, openModal, closeModal — completely unused"
    - Path: frontend/src/store/readerSlice.ts
      Note: "Pre-built Redux slice with startReading, goToSection, markSectionRead — completely unused"
    - Path: frontend/src/store/hooks.ts
      Note: "useAppDispatch, useAppSelector typed hooks — ready to use"
    - Path: frontend/src/app/ReaderApp.tsx
      Note: "Already uses useParams() for /r/:token route"
ExternalSources: []
Summary: "Replace useState-based view switching in AuthorApp with React Router routes and wire unused Redux slices for shared state."
LastUpdated: 2026-03-26T17:00:00-04:00
WhatFor: "Guide the migration from local state navigation to URL-based routing and Redux state management."
WhenToUse: "When implementing DR-006. Follow phases in order."
---

# DR-006 Implementation Plan

## Current State

### AuthorApp navigation: entirely useState

`AuthorApp.tsx` manages all view switching with local state:

```tsx
const [view, setView] = useState<View>("dashboard");          // line 38
const [selectedArticle, setSelectedArticle] = useState(null);  // line 39
const [focusSection, setFocusSection] = useState(undefined);   // line 40
const [showInvite, setShowInvite] = useState(false);           // line 41
const [shareUrl, setShareUrl] = useState(undefined);           // line 42
```

View type: `"dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"`

`goBack()` always resets to `"dashboard"` (line 92-96).

### Existing Redux infrastructure: fully wired, completely unused

**uiSlice** — already in store, typed hooks ready:
- `setView(View)` — matches AuthorApp.view (type mismatch: uiSlice has different View enum)
- `selectArticle(string)` — matches AuthorApp.selectedArticle (stores ID, not full object)
- `openModal(string)` / `closeModal()` — matches AuthorApp.showInvite
- `setActiveMenu(string | null)` — could enhance MenuBar

**readerSlice** — already in store:
- `startReading(sectionId)` — matches ReaderPage.started + currentSectionId
- `goToSection(sectionId)` — matches ReaderPage.currentSectionId
- `markSectionRead(sectionId)` — matches ReaderPage.readSectionIds
- `resetReader()` — cleanup on exit

### Existing React Router: BrowserRouter in App.tsx

`App.tsx` already has:
```tsx
<BrowserRouter>
  <Routes>
    <Route path="/r/:token" element={<ReaderApp />} />
    <Route path="/*" element={<AuthorApp />} />
  </Routes>
</BrowserRouter>
```

`ReaderApp` already uses `useParams()`. AuthorApp catches `/*` but ignores the URL.

---

## Target Route Structure

```
/                           → Dashboard (first article selected)
/articles                   → All Articles list
/articles/:id               → Dashboard with specific article selected
/articles/:id/edit          → Editor
/articles/:id/review        → Review view (author's reaction review)
/articles/:id/settings      → Settings
/articles/:id/preview       → Preview as Reader
```

---

## Phase 1: Wire uiSlice for AuthorApp (prerequisite)

Replace AuthorApp's 5 useState calls with Redux state. This decouples view management from the component and makes it available to any child without prop drilling.

### Changes

**1a. Align uiSlice View type with AuthorApp**

Current uiSlice View: `"dashboard" | "articles" | "article-edit" | "article-settings" | "article-review"`
Current AuthorApp View: `"dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"`

Reconcile to: `"dashboard" | "articles" | "review" | "edit" | "settings" | "reader-preview"`

Add missing state to uiSlice:
```ts
interface UiState {
  view: View;
  selectedArticleId: string | null;
  focusSectionId: string | null;    // NEW
  activeModal: string | null;
  activeMenu: string | null;
}
```

**1b. Replace useState in AuthorApp**

```tsx
// BEFORE
const [view, setView] = useState<View>("dashboard");
const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
const [focusSection, setFocusSection] = useState<string | undefined>();
const [showInvite, setShowInvite] = useState(false);

// AFTER
const dispatch = useAppDispatch();
const view = useAppSelector(state => state.ui.view);
const selectedArticleId = useAppSelector(state => state.ui.selectedArticleId);
const focusSection = useAppSelector(state => state.ui.focusSectionId);
const activeModal = useAppSelector(state => state.ui.activeModal);

// Hydrate full Article from RTK Query cache
const selectedArticle = articles.find(a => a.id === selectedArticleId) || null;
const showInvite = activeModal === "invite";
```

Replace all `setView("x")` with `dispatch(setView("x"))`, etc.

**1c. Keep shareUrl as local state** — it's transient per-article, no cross-component use.

### Why phase 1 before routing

The routing migration is simpler if view state is already in Redux. Route changes can dispatch Redux actions, and components can read state from Redux instead of getting it through AuthorApp props.

---

## Phase 2: Add React Router routes to AuthorApp

### Changes

**2a. Define nested routes in App.tsx**

```tsx
<Route path="/*" element={<AuthorApp />}>
  <Route index element={<DashboardView />} />
  <Route path="articles" element={<ArticlesView />} />
  <Route path="articles/:id" element={<DashboardView />} />
  <Route path="articles/:id/edit" element={<EditorView />} />
  <Route path="articles/:id/review" element={<ReviewView />} />
  <Route path="articles/:id/settings" element={<SettingsView />} />
  <Route path="articles/:id/preview" element={<PreviewView />} />
</Route>
```

OR keep AuthorApp as a switch on `useLocation()` — simpler if we don't want route-level components.

**2b. Replace dispatch(setView(...)) with navigate(...)**

```tsx
// BEFORE
dispatch(setView("edit"));
dispatch(selectArticle(id));

// AFTER
navigate(`/articles/${id}/edit`);
```

The route change drives the Redux state via a sync effect:

```tsx
const { id } = useParams();
const location = useLocation();

useEffect(() => {
  // Derive view from URL path
  const path = location.pathname;
  if (path.endsWith("/edit")) dispatch(setView("edit"));
  else if (path.endsWith("/review")) dispatch(setView("review"));
  // ...etc

  if (id) dispatch(selectArticle(id));
}, [location, id, dispatch]);
```

**2c. Replace goBack() with navigate(-1)**

The current `goBack()` always goes to dashboard. With React Router, `navigate(-1)` uses browser history, so Dashboard → Articles → Edit → Back goes to Articles (correct), not Dashboard (current wrong behavior).

**2d. Update menu bar items**

Menu "View" items become navigation links:
```tsx
{ label: "Dashboard", action: () => navigate("/") },
{ label: "Articles", action: () => navigate("/articles") },
```

---

## Phase 3: Wire readerSlice for ReaderPage

### Current ReaderPage useState (all replaceable)

```tsx
const [started, setStarted] = useState(false);                    // → readerSlice.started
const [sessionId, setSessionId] = useState<string | null>(null);  // → readerSlice (new field)
const [currentSectionId, setCurrentSectionId] = useState(...);    // → readerSlice.currentSectionId
const [readSectionIds, setReadSectionIds] = useState<string[]>([]); // → readerSlice.readSectionIds
```

### Changes

**3a. Add sessionId to readerSlice**

```ts
interface ReaderState {
  started: boolean;
  sessionId: string | null;      // NEW
  currentSectionId: string | null;
  readSectionIds: string[];
}
```

Add action: `setSessionId(string)`

**3b. Replace useState in ReaderPage**

```tsx
const dispatch = useAppDispatch();
const { started, sessionId, currentSectionId, readSectionIds } =
  useAppSelector(state => state.reader);
```

Replace:
- `setStarted(true)` → `dispatch(startReading(sectionId))`
- `setCurrentSectionId(id)` → `dispatch(goToSection(id))`
- `setReadSectionIds(...)` → `dispatch(markSectionRead(id))`

**3c. Keep reactions as local state** — optimistic updates work better locally; reactions are also ephemeral in preview mode.

**3d. Reset on unmount** — dispatch `resetReader()` when leaving ReaderPage.

---

## Phase 4: Cleanup and edge cases

- Remove the `View` type from AuthorApp (use from uiSlice or derive from route)
- Remove `goBack()` function (replaced by `navigate(-1)`)
- Update `handleSelectArticle` to navigate instead of dispatching multiple state changes
- Test: browser back button, page refresh, deep linking, bookmark
- Test: preview mode → back to editor preserves article state
- Handle 404 for invalid article IDs in URL

---

## useState/useEffect audit: what stays local

These are all fine as local state and should NOT move to Redux:

| Component | State | Why local |
|-----------|-------|-----------|
| ArticleEditor | article, activeSectionId, contentDraft, showDeleteConfirm, showUnsavedConfirm | Editor working draft — component-scoped, never shared |
| ArticleSettings | title, intro, version, status, accessMode, all toggles, copied, showDeleteConfirm | Form state — collected and submitted together |
| Dashboard | selectedId | Tab selection within dashboard — no cross-component use |
| ArticleManager | search, sortKey, expanded | Filter/sort state — component-scoped |
| ArticleReader | selectedSection, reactionFilter, hoveredParagraph, filteredParagraph | View filters — component-scoped |
| InviteDialog | email, note, inviteResult, isSending, isGeneratingShareLink, etc. | Modal form state — dies when dialog closes |
| Paragraph | activeType, comment | Reaction chip interaction — per-paragraph |
| MenuBar | openMenu | Dropdown state — transient |
| ReactionPicker | chosenType, comment | Selection state — per-picker |
| ReaderPage | reactions, showDone | Reactions: optimistic local; showDone: dialog toggle |

---

## Risk assessment

| Risk | Mitigation |
|------|-----------|
| URL changes break existing links | Only internal — app is not public yet |
| AuthorApp refactor is large | Phase 1 (Redux) is independent of Phase 2 (Router) — can ship separately |
| Race between route change and data loading | RTK Query handles loading states already; show loading spinner per route |
| Preview mode needs article data in memory | Editor sets selectedArticle before navigating to preview — data stays in Redux |
