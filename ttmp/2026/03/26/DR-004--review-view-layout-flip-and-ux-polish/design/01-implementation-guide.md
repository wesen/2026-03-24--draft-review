---
Title: Implementation Guide
Ticket: DR-004
Status: active
Topics:
    - frontend
    - react
    - ux
    - ui
DocType: design
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.tsx
      Note: "Primary target — layout flip from top/bottom to side-by-side"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.css
      Note: "All layout CSS changes for the flip"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx
      Note: "Stat card compaction and hover states on tabs/rows"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.css
      Note: "Stat card and hover CSS"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleManager.css
      Note: "Hover states for article cards"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/InviteDialog.css
      Note: "Already has shadow and backdrop — make backdrop darker"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/primitives/Prose.css
      Note: "Already has code block styling — verified working"
ExternalSources: []
Summary: "Step-by-step implementation guide for the review view side-by-side layout, invite dialog backdrop fix, hover states across all interactive elements, and dashboard stat card compaction."
LastUpdated: 2026-03-26T09:09:06.0693423-04:00
WhatFor: "Guide implementation of the four highest-priority UX improvements from the DR-003 design review."
WhenToUse: "When starting implementation. Follow the work items in order."
---

# DR-004 Implementation Guide

## Overview

Four high-impact UX improvements from the DR-003 design review, ordered by impact:

| # | Item | Impact | Effort | Files |
|---|------|--------|--------|-------|
| 1 | Review view side-by-side layout | Very high | Medium | `ArticleReader.tsx`, `ArticleReader.css` |
| 2 | Invite dialog backdrop darkening | Low | Trivial | `InviteDialog.css` |
| 3 | Hover states on tabs, cards, rows | Medium | Low | `Dashboard.css`, `ArticleManager.css`, `ArticleReader.css` |
| 4 | Compact stat cards | Medium | Low | `Dashboard.tsx`, `Dashboard.css` |
| 5 | Increase font sizes to normal web sizes | High | Low | `mac-os-1.css`, `Prose.css`, `Paragraph.css`, `ArticleReader.css` |

---

## Work Item 1: Review View Side-by-Side Layout

### Current state

`ArticleReader.tsx` uses a three-panel layout:

```
┌──────────┬────────────────────────────────┐
│ Sections │  Article text (flex: 1)        │
│ sidebar  │  (scrolls independently)       │
│ (180px)  │                                │
│          ├────────────────────────────────┤
│          │  Reactions panel               │
│          │  (min 150px, max 50%)          │
│          │  filter tabs + grouped list    │
│          └────────────────────────────────┘
└──────────┘
```

The main area (`dr-article-reader__main`) is a vertical flex column with `__text` on top and `__reactions-panel` on the bottom. The text area gets `flex: 1` (most of the space) and the reactions panel gets the remainder.

**Problem**: The author came here to read feedback. The feedback panel is cramped at the bottom while the text (which the author wrote) dominates the view. When hovering a reaction card to highlight its paragraph, the paragraph may be scrolled out of view in the text area above.

### Target state

```
┌──────────┬──────────────────┬──────────────────┐
│ Sections │  Article text    │  Reactions panel  │
│ sidebar  │  (scrolls)       │  (scrolls)        │
│ (180px)  │                  │  filter tabs      │
│          │                  │  grouped list     │
│          │                  │                   │
└──────────┴──────────────────┴──────────────────┘
```

The main area becomes a horizontal flex row: text on the left, reactions on the right. Both scroll independently. The bidirectional hover highlighting works much better because both the paragraph and its reactions are visible simultaneously.

### Implementation

#### CSS changes (`ArticleReader.css`)

The key change is flipping `__main` from `flex-direction: column` to `row`:

```css
/* Main — CHANGE from column to row */
.dr-article-reader__main {
  flex: 1;
  display: flex;
  flex-direction: row;  /* was: column */
}

.dr-article-reader__text {
  flex: 1;
  padding: 20px 28px;
  overflow: auto;
  border-right: var(--dr-border-width) solid var(--dr-border-color);  /* was: border-bottom */
  /* REMOVE: border-bottom */
}

.dr-article-reader__reactions-panel {
  width: 340px;           /* fixed width instead of height constraints */
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  /* REMOVE: min-height, max-height, flex-grow */
}
```

The reactions panel gets a fixed width (340px) which gives it enough room for reaction cards while leaving the majority of the horizontal space for article text. The `max-width: 540px` constraint on `__paragraph` still applies, so text doesn't stretch to fill.

#### Reaction filter tabs

The filter tabs currently run horizontally across the full panel width. In the side-by-side layout they'll be narrower. Two approaches:

**Option A — Keep horizontal tabs, allow wrapping**:
```css
.dr-article-reader__filter-tabs {
  flex-wrap: wrap;
}
```

**Option B — Compact tabs**: Shorten labels to just icons + counts:
```
★ 3 | ? 1 | ◎ 0 | ♥ 2
```
instead of:
```
★ Useful! (3) | ? Confused (1) | ◎ Slow (0) | ♥ Love it (2)
```

**Recommendation**: Option B for the compact side panel. The icons are recognizable enough without labels. Keep the full labels only in the "All (N)" tab. Show a `title` attribute with the full label for accessibility.

#### Paragraph excerpt in grouped reactions

The paragraph excerpt currently appears as an italic label above each group. In the side-by-side layout, the excerpt becomes less important because the reader can see the actual paragraph in the left panel (especially with bidirectional hover highlighting). Consider making the excerpt more compact:

```css
.dr-article-reader__para-excerpt {
  font-size: 9px;       /* was: var(--dr-font-size-sm) = 10px */
  padding: 2px 8px;     /* was: 4px 10px */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
```

#### TSX changes (`ArticleReader.tsx`)

Minimal JSX changes needed. The component structure stays the same — only the CSS layout changes. Two optional improvements:

1. **Compact filter tabs** — change the render of filter tabs to show only icon + count:

```tsx
{REACTION_TYPES.map((rt) => {
  const c = sectionReactions.filter((r) => r.type === rt.type).length;
  return (
    <div
      key={rt.type}
      title={`${rt.label} (${c})`}
      onClick={() => setReactionFilter(reactionFilter === rt.type ? null : rt.type)}
      className={`dr-article-reader__filter-tab ${
        reactionFilter === rt.type ? "dr-article-reader__filter-tab--active" : ""
      }`}
    >
      {rt.icon} {c}
    </div>
  );
})}
```

2. **Reactions panel header**: Add a small "Reactions" header above the filter tabs (matching the sidebar "SECTIONS" header) to make the panel feel like a proper column:

```tsx
<div className="dr-article-reader__reactions-panel">
  <div className="dr-article-reader__reactions-header">Reactions</div>
  {/* filter tabs, list, etc. */}
</div>
```

```css
.dr-article-reader__reactions-header {
  padding: 8px 10px;
  border-bottom: 1px solid var(--dr-on-surface);
  font-family: var(--dr-font-title);
  font-size: var(--dr-font-size-sm);
  text-transform: uppercase;
  letter-spacing: 1px;
  color: var(--dr-subtle);
}
```

This mirrors the sidebar's `__sidebar-header` style exactly, giving visual symmetry.

### Verification

After the layout flip:
- Both panels should scroll independently (test with long article + many reactions)
- Bidirectional hover should work: hovering a paragraph highlights reaction cards visible in the right panel; hovering a reaction card highlights the paragraph (scrolling it into view if needed via `scrollIntoView`)
- Filter tabs should be readable at 340px width
- Paragraph filter bar should fit within 340px
- The paragraph excerpt should truncate with ellipsis if too long

---

## Work Item 2: Invite Dialog Backdrop

### Current state

`InviteDialog.css` already has both shadow and backdrop:

```css
/* line 8 */  background: rgba(0, 0, 0, 0.15);
/* line 17 */ box-shadow: var(--dr-shadow-dialog); /* 6px 6px 0 #000 */
```

The DR-003 design review incorrectly stated these were missing. However, the 15% opacity backdrop is very subtle — easy to miss visually.

### Implementation

One-line change — darken the backdrop:

```css
.dr-invite-overlay {
  background: rgba(0, 0, 0, 0.35);  /* was: 0.15 */
}
```

This matches the standard Mac OS dialog dimming where the background visually recedes.

### Also fix: `border-radius: 4px`

`InviteDialog.css` line 16 has `border-radius: 4px` which violates the Mac OS 1 aesthetic (should be `0`). Fix:

```css
.dr-invite-dialog {
  border-radius: 0;  /* was: 4px */
}
```

---

## Work Item 3: Hover States

### Current state

Interactive elements with **no hover feedback**:
- Dashboard article tabs (`.dr-dashboard__tab`)
- Dashboard reader rows (`.dr-dashboard__reader-row`)
- Dashboard section rows (`.dr-dashboard__section-row`)
- Dashboard feedback items (`.dr-dashboard__feedback-item`)
- Article Manager cards (`.dr-manager__card`)
- Review view section sidebar items (`.dr-article-reader__section-item`) — already have active state (inverted) but no hover state for non-selected items

### Implementation

Add consistent hover states using `var(--dr-hover-bg)` throughout. Each gets a `background` transition on hover:

#### Dashboard (`Dashboard.css`)

```css
/* Article tabs — hover for non-active */
.dr-dashboard__tab:hover {
  background: var(--dr-hover-bg);
}

/* Reader rows */
.dr-dashboard__reader-row:hover {
  background: var(--dr-hover-bg);
}

/* Section rows */
.dr-dashboard__section-row:hover {
  background: var(--dr-hover-bg);
}

/* Feedback items */
.dr-dashboard__feedback-item:hover {
  background: var(--dr-hover-bg);
}
```

#### Article Manager (`ArticleManager.css`)

```css
/* Article cards */
.dr-manager__card:hover {
  background: var(--dr-hover-bg);
}
```

#### Review sidebar (`ArticleReader.css`)

The section sidebar items already invert when active. Add a subtle hover for non-active items:

```css
.dr-article-reader__section-item:hover {
  background: var(--dr-hover-bg);
}

/* Don't override the active inverted style */
.dr-article-reader__section-item--active:hover {
  background: var(--dr-on-surface);
}
```

#### Transition

Add a shared transition to make the hover feel smooth:

```css
.dr-dashboard__tab,
.dr-dashboard__reader-row,
.dr-dashboard__section-row,
.dr-dashboard__feedback-item,
.dr-manager__card,
.dr-article-reader__section-item {
  transition: background 0.1s;
}
```

Or add `transition: background 0.1s` to each element's base rule individually.

---

## Work Item 4: Compact Stat Cards

### Current state

`Dashboard.tsx` renders four `<StatCard>` components in a grid (`grid-template-columns: repeat(auto-fit, minmax(130px, 1fr))`). Each card takes substantial vertical space (icon + value + label) for a single metric:

```
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│  ◉ 1     │ │  ✦ 0     │ │  § 1     │ │  ▸ 0%    │
│ READERS  │ │REACTIONS │ │ SECTIONS │ │AVG PROG  │
└──────────┘ └──────────┘ └──────────┘ └──────────┘
```

On a maximized window, each card is ~200px wide for a single digit.

### Target state

Replace the four-card grid with a single compact stat bar:

```
◉ 1 readers · ✦ 0 reactions · § 1 sections · ▸ 0% avg progress
```

This frees vertical space for the more useful content below (readers panel, reactions chart, feedback).

### Implementation

#### TSX changes (`Dashboard.tsx`)

Replace the stat cards grid with an inline stat bar:

```tsx
{/* Replace the stats grid */}
<div className="dr-dashboard__stat-bar">
  <span>{"\u25C9"} {articleReaders.length} readers</span>
  <span>{"\u2726"} {totalReactions} reactions</span>
  <span>{"\u00A7"} {selected.sections.length} sections</span>
  <span>
    {"\u25B8"} {articleReaders.length ? `${avgProgress}%` : "\u2014"} avg
    progress
  </span>
</div>
```

This replaces the entire `<div className="dr-dashboard__stats">` block including the four `<StatCard>` components.

#### CSS (`Dashboard.css`)

Replace the stats grid CSS:

```css
/* Replace .dr-dashboard__stats grid */
.dr-dashboard__stat-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 8px 0;
  margin-bottom: 12px;
  font-family: monospace;
  font-size: var(--dr-font-size-base);
  color: var(--dr-subtle);
  border-bottom: 1px dotted var(--dr-border-light);
}

.dr-dashboard__stat-bar span {
  white-space: nowrap;
}
```

#### Cleanup

The `StatCard` component (`primitives/StatCard.tsx`) and its CSS can be left in place — it may be useful elsewhere later. Just remove the import from `Dashboard.tsx` if it's the only consumer.

---

## Work Item 5: Increase Font Sizes

### Current state

The design token font sizes in `mac-os-1.css` are very small — scaled for a 1:1 pixel-art aesthetic but uncomfortable for actual reading:

```css
--dr-font-size-xs: 9px;
--dr-font-size-sm: 10px;
--dr-font-size-base: 11px;
--dr-font-size-md: 12px;
--dr-font-size-lg: 14px;
--dr-font-size-xl: 18px;
--dr-font-size-2xl: 20px;
```

The Prose component uses `font-size: 13px` (hardcoded in `Prose.css`), and many components use the `--dr-font-size-base` (11px) and `--dr-font-size-md` (12px) tokens. These are too small for comfortable reading, especially the reader view where users spend extended time.

### Implementation

**Reference**: Substack article pages use ~18px body text, ~28px h2, ~1.7 line-height. Our current 13px body is nearly half that. We don't need to match Substack exactly (our retro aesthetic can be a bit tighter), but we should land in the 15–17px range for prose body text and scale everything proportionally.

Bump all token values to normal web sizes:

```css
--dr-font-size-xs: 11px;    /* was: 9px */
--dr-font-size-sm: 12px;    /* was: 10px */
--dr-font-size-base: 14px;  /* was: 11px */
--dr-font-size-md: 15px;    /* was: 12px */
--dr-font-size-lg: 17px;    /* was: 14px */
--dr-font-size-xl: 21px;    /* was: 18px */
--dr-font-size-2xl: 26px;   /* was: 20px */
```

Also update the hardcoded sizes in `Prose.css` — this is what readers spend the most time looking at:

```css
.dr-prose {
  font-size: 17px;   /* was: 13px — approaching Substack's 18px */
  line-height: 1.75;
}

.dr-prose h1 { font-size: 26px; }  /* was: 20px */
.dr-prose h2 { font-size: 22px; }  /* was: 16px */
.dr-prose h3 { font-size: 18px; }  /* was: 14px */

.dr-prose code { font-size: 14px; }  /* was: 11px */
.dr-prose pre { font-size: 14px; }   /* was: 11px */
.dr-prose table { font-size: 14px; } /* was: 11px */
```

And in `ArticleReader.css`, the paragraph font size:

```css
.dr-article-reader__paragraph {
  font-size: 17px;  /* was: 13px — match Prose base */
}
```

### What to watch for

- The menu bar height (`--dr-menubar-height: 20px`) and title bar height (`--dr-titlebar-height: 22px`) will need bumping — likely to 26px and 28px respectively — since the font inside them grows
- The reaction chips in `Paragraph.css` use `font-size: 9px` (hardcoded) — bump to 11px
- The badge in `Paragraph.css` uses `font-size: 8px` — bump to 10px
- The heading in `ArticleReader.css` uses `font-size: 16px` — bump to 20px
- The section badge uses `font-size: 9px` via `--dr-font-size-xs` — will auto-update
- The close box size (`--dr-close-box-size: 12px`) may need bumping to 14–16px
- Stat cards, status badges, and other small labels should be checked visually after the bump
- The `max-width: 540px` on paragraphs may need widening to ~620px to accommodate the larger text comfortably

### Approach

Since all component sizes derive from the CSS variable tokens, changing them in `mac-os-1.css` propagates everywhere automatically. The hardcoded sizes in `Prose.css`, `Paragraph.css`, and `ArticleReader.css` need individual updates.

---

## Dependency Graph

```
Work Item 1 (layout flip)     — standalone, no dependencies
Work Item 2 (backdrop)         — standalone, trivial
Work Item 3 (hover states)     — standalone, touches 3 CSS files
Work Item 4 (stat compaction)  — standalone, Dashboard only
```

All four are independent and can be done in any order. Recommended commit strategy:

1. Commit Work Item 1 alone (it's the biggest and most impactful)
2. Commit Work Items 2 + 3 together (both are CSS-only polish)
3. Commit Work Item 4 alone (Dashboard layout change)
