---
Title: Implementation Guide
Ticket: DR-005
Status: active
Topics:
    - frontend
    - react
    - ux
    - bugs
DocType: design
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/app/AuthorApp.tsx
      Note: "All author view switching — no URL routing, goBack() always goes to dashboard"
    - Path: frontend/src/author/ArticleEditor.tsx
      Note: "Editor — delete section works locally but unclear if persisted; preview opens reader without API hooks"
    - Path: frontend/src/author/ArticleSettings.tsx
      Note: "Status changes work (PATCH /api/articles/:id), but COMING SOON settings are local-only"
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: "onReactionAdd/onReactionRemove are optional — preview mode passes neither"
    - Path: frontend/src/author/Dashboard.tsx
      Note: "Dashboard duplicates article selection with tab clicks but has no way to edit article title/intro"
    - Path: frontend/src/store/
      Note: "Redux store — uiSlice manages view state without URL routing"
ExternalSources: []
Summary: "Author UX flow improvements: fix broken preview reactions, navigation model, duplicate paths, and missing article management features."
LastUpdated: 2026-03-26T14:30:00-04:00
WhatFor: "Guide implementation of author-facing UX fixes validated by Playwright testing."
WhenToUse: "When implementing author UX improvements. Follow work items in priority order."
---

# DR-005 Implementation Guide

## Overview

Findings from the DR-003 design review validated with Playwright testing against the running application. Covers functional bugs, navigation model issues, and missing admin features.

### Testing Summary

| Area | Status | Notes |
|------|--------|-------|
| Delete Section | **Works** | Local state delete → confirmed via dialog, BUT only saves when user clicks "Save" |
| Status Change | **Works** | PATCH /api/articles/:id fires on Save Settings → persists |
| Preview as Reader reactions | **Broken** | Reactions go to local state only — no API call, silently lost |
| Navigation / URL routing | **Missing** | All author views on `/` — no back button, no bookmarks, no deep links |
| Article title/intro editing | **Missing** | No UI to edit article title, intro, or version outside of Settings |
| Delete article | **Missing** | No way to delete an entire article |

---

## Work Item 1: Fix Preview-as-Reader Reactions (Bug)

### Current state

`AuthorApp.tsx` line 403 renders `<ReaderPage>` for preview mode without `onReactionAdd` or `onReactionRemove` props:

```tsx
{view === "reader-preview" && selectedArticle && (
  <MacWindow ...>
    <ReaderPage
      article={{...selectedArticle}}
      // ← no onReactionAdd, no onReactionRemove
    />
  </MacWindow>
)}
```

`ReaderPage.tsx` defines these as optional (`onReactionAdd?: ...`), so reactions in preview mode go into the component's local `useState` and are silently lost when the user navigates away. No API call is made. The user sees the reaction appear in the UI (badge count, reaction card) but nothing persists.

### Problem

An author testing their article in preview mode submits a reaction, sees it appear, and assumes it's saved. It isn't. This is misleading — the UI provides no indication that reactions are ephemeral in preview mode.

### Options

**Option A — Disable reactions in preview mode** (recommended):
Show the paragraph hover highlighting and chips at reduced opacity but make them non-clickable. Add a banner at the top: "Preview mode — reactions are disabled."

```tsx
<ReaderPage
  article={{...selectedArticle}}
  readOnly={true}  // new prop
/>
```

This is the simplest fix and sets correct expectations. The author is previewing layout, not testing the feedback flow.

**Option B — Wire reactions to the real API**:
Pass `onReactionAdd` that calls the reaction mutation. This requires the preview to have a valid reader identity (currently it doesn't — there's no `readerId` or `inviteCode` in preview mode).

**Option C — Show ephemeral banner**:
Keep local-only reactions but show a yellow/hatched banner: "Preview mode — reactions shown here are not saved."

### Recommendation

Option A. Preview is for checking paragraph layout and formatting, not for testing the feedback system. Disabling reactions avoids confusion entirely.

---

## Work Item 2: Navigation Model — Add URL Routing

### Current state

All author views are rendered on `http://localhost:5174/` with view switching managed by `useState<View>("dashboard")` in `AuthorApp.tsx`. The `goBack()` function always sets view to `"dashboard"` regardless of navigation history.

Available views: `"dashboard" | "articles" | "article" | "edit" | "settings" | "reader-preview"`

### Problems

1. **Browser back button doesn't work** — pressing Back navigates away from the app entirely
2. **Can't bookmark or share a URL** to a specific article's editor or review view
3. **Page refresh always returns to dashboard** — loses the user's current context
4. **`goBack()` always goes to dashboard** — if you go Dashboard → Articles → Edit, back takes you to Dashboard, not Articles

### Implementation

The app already uses React Router (it's in `package.json` and used for reader routes in `App.tsx`). Extend it to author routes:

```
/                           → Dashboard (default article selected)
/articles                   → All Articles list
/articles/:id               → Dashboard with specific article selected
/articles/:id/edit          → Editor
/articles/:id/review        → Review view
/articles/:id/settings      → Settings
/articles/:id/preview       → Preview as Reader
```

Key changes:
- Replace `useState<View>` with `useNavigate()` / `useParams()`
- Replace `setView("edit")` calls with `navigate(`/articles/${id}/edit`)`
- Replace `goBack()` with `navigate(-1)` (uses browser history)
- Keep the menu bar's View menu items as navigation links

### Effort

Medium — the component tree stays the same, but the view switching logic moves from state to routing. Each view component needs to read `articleId` from URL params instead of local state.

---

## Work Item 3: Missing Article Management Features

### 3a. Edit article title, intro, author note, and version

Currently the only way to edit an article's title is... there isn't one visible in the UI. The article title shown in dashboard tabs comes from the API but there's no editor for it. The `ArticleSettings` view has status and sharing config but no title/intro/version fields.

**Fix**: Add a "Details" section at the top of ArticleSettings (or a separate "Article Info" view) with:
- Article title (text input)
- Author intro/note (textarea — the text shown on the reader welcome splash)
- Version label (text input — "Draft 1", "Revision 2", etc.)

These fields already exist on the `Article` type and are already passed to `onSave` via PATCH.

### 3b. Delete article

There's no way to delete an article. The "Archived" status exists but doesn't remove the article from the list.

**Fix**: Add a "Delete Article" button in ArticleSettings under the Status section, inside a danger zone with a confirmation dialog. Wire to `DELETE /api/articles/:id`.

### 3c. Editor "Save" only saves to local state until explicit save

The editor works entirely with local `useState`. Changes to section title, content, adding/deleting sections — all happen in local state. Only when the user clicks "Save" does `onSave(article)` fire, which triggers the API mutation.

This is actually correct behavior (draft editing), but the UX provides no indication of unsaved changes. The user can add 5 sections, edit content, then click "← Back" and lose everything.

**Fix**:
- Add a dirty state indicator: "Unsaved changes" label next to the Save button
- Confirm on navigation away: "You have unsaved changes. Save before leaving?"
- Auto-save on a timer (every 30s) or on blur

---

## Work Item 4: Navigation Clarity and Duplicate Paths

### Current duplicate paths

The design review identified confusion about how to reach each view. Validated findings:

| Action | Path 1 | Path 2 | Path 3 |
|--------|--------|--------|--------|
| View reactions | Dashboard → "Open Full Review" | All Articles → "Review" button | Dashboard → click section in "Reactions by Section" |
| Edit article | Dashboard → "Edit" button | All Articles → "Edit" button | Menu bar → Edit (doesn't exist) |
| Settings | Dashboard → "Share" button | All Articles → "Share" button | — |
| Invite reader | Dashboard → "+ Invite Reader" | All Articles → "Invite →" button | Dashboard → "Share" → Generate |

The "Share" button on the dashboard opens ArticleSettings, not the invite dialog. This is confusing — "Share" suggests sharing a link, but it opens the full settings page. The "+ Invite Reader" opens the invite dialog which is specifically about sharing.

### Recommendations

1. **Rename "Share" to "Settings"** on the dashboard action bar — it opens the settings page
2. **Or split**: keep "Share" but make it open the invite dialog (same as "+ Invite Reader"), and add a separate "Settings" button
3. **Remove duplicate "Invite →"** from article cards in the All Articles view — it does the same thing as the dashboard's "+ Invite Reader"
4. **Add "Edit" to the menu bar** under a new "Article" menu: Edit, Review, Settings, Preview

---

## Work Item 5: Editor UX Improvements

### 5a. Textarea doesn't auto-grow

The content textarea has `rows={16}` — fixed height. Long articles require scrolling inside the textarea inside the scrolling window (double scroll).

**Fix**: Use a `ref` to auto-grow the textarea on input:
```tsx
const autoGrow = (el: HTMLTextAreaElement) => {
  el.style.height = "auto";
  el.style.height = el.scrollHeight + "px";
};
```

### 5b. "Danger Zone" label is dramatic

The delete section area has a dashed border with "DANGER ZONE" in uppercase. For deleting a single section (with a confirmation dialog), this is over the top.

**Fix**: Remove the "DANGER ZONE" label. Keep the button with its confirmation dialog — that's sufficient protection.

### 5c. Section reorder arrows are subtle

The ▲▼ arrows are 10px and low contrast. They could be mistaken for decoration.

**Fix**: Increase to 12px, add cursor: pointer, and show a hover state (invert on hover like buttons).

---

## Work Item 6: Empty States Cleanup

The DR-003 review noted inconsistent empty states. Current state after DR-004 fixes:

| View | Empty State | Issue |
|------|-------------|-------|
| Dashboard / No readers | "✉ No readers yet. Invite your first beta reader!" | Text glyph ✉ is good |
| Dashboard / No feedback | "No feedback yet. Share your article to start collecting reactions." | Good, but inside a bordered panel that takes 80px for one line |
| Review / No reactions | "No reactions for this section yet." | Good |
| All Articles / No articles | "✎ No articles yet..." with CTA | Good |

Mostly fixed by DR-003/DR-004. One remaining issue: the "Recent Feedback" panel on the dashboard takes substantial vertical space to show "No feedback yet." When empty, collapse it to a single borderless line or hide it entirely.

---

## Priority Order

| # | Item | Impact | Effort | Type |
|---|------|--------|--------|------|
| 1 | Fix preview reactions (WI-1) | High | Low | Bug fix |
| 2 | Article title/intro editing (WI-3a) | High | Low | Missing feature |
| 3 | Rename "Share" → "Settings" (WI-4) | Medium | Trivial | UX clarity |
| 4 | Unsaved changes indicator (WI-3c) | Medium | Low | UX safety |
| 5 | URL routing (WI-2) | High | Medium | UX infrastructure |
| 6 | Delete article (WI-3b) | Medium | Low | Missing feature |
| 7 | Editor textarea auto-grow (WI-5a) | Low | Low | UX polish |
| 8 | Remove Danger Zone label (WI-5b) | Low | Trivial | UX polish |
