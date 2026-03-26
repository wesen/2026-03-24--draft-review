---
Title: Implementation Diary
Ticket: IMPROVE-UX
Status: active
Topics:
    - implementation
    - frontend
    - ux
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/author/Dashboard.tsx
      Note: "Redesigned: tab strip + actions bar → sidebar article list with inline action icons"
    - Path: frontend/src/author/Dashboard.css
      Note: "New layout: sidebar + main split, article row styles, action icon hover reveal"
ExternalSources: []
Summary: "Diary of dashboard redesign from tab-based to sidebar list layout."
LastUpdated: 2026-03-26T19:00:00-04:00
WhatFor: "Document what changed and design decisions for code review."
WhenToUse: "When reviewing the dashboard redesign."
---

# IMPROVE-UX Implementation Diary

## Session: 2026-03-26

### Commit 1: `0f78be1` — Dashboard redesign

**What changed:**

The dashboard previously had:
1. A **tab strip** at the top with article titles as clickable tabs (pill-shaped, active one inverted)
2. A separate **actions bar** below with Edit / Settings / + Invite Reader buttons
3. Main content below: stat bar, readers panel, reactions chart, draft-killer alert, recent feedback

The new layout:
1. A **sidebar** (260px) on the left with:
   - Header: "Articles" label + "All" button
   - Scrollable article list — each row shows:
     - Article title (truncated with ellipsis)
     - Meta line: status badge + reader count + feedback count (e.g. "draft · 3R 12F")
     - Action icons (hover-revealed): ✎ Edit, ⚙ Settings, ✉ Invite, ▶ Review
   - Footer: "+ New Article" button
2. **Main content** on the right (flex: 1, scrollable) with the same panels as before

**Design decisions:**

- **Action icons hidden until hover** — keeps the list clean. On the active (selected) row, icons are always visible since the row is inverted.
- **Icons use Unicode glyphs** matching the Mac OS 1 aesthetic: ✎ (U+270E), ⚙ (U+2699), ✉ (U+2709), ▶ (U+25B6). These render consistently cross-platform.
- **Per-row meta**: shows `{count}R {count}F` (readers / feedback) — gives at-a-glance comparison across articles without selecting each one.
- **Status badge** inline in the meta line with `text-transform: capitalize` and a `1px solid currentColor` border.
- **Selected state**: full-width inverted row (`background: var(--dr-on-surface); color: var(--dr-surface)`) — unmistakable.
- **Action icon hover**: inverts to white background / black text on hover, matching Mac OS 1 button style.
- **Removed**: the actions bar (Edit / Settings / Invite buttons), the tab strip, the `onViewArticles` button (moved to sidebar header as "All").
- **Kept**: stat bar, panels grid, draft-killer alert, recent feedback, "Open Full Review" bottom action — all unchanged, just moved into the main content area.

**What I removed from CSS:**
- `.dr-dashboard__tabs`, `.dr-dashboard__tab`, `.dr-dashboard__tab--active`, `.dr-dashboard__tab-status` — replaced by article list
- `.dr-dashboard__top-bar`, `.dr-dashboard__actions-bar` — replaced by sidebar

**What I added to CSS:**
- `.dr-dashboard--split` — flex container for sidebar + main
- `.dr-dashboard__sidebar`, `__sidebar-header`, `__sidebar-footer` — sidebar frame
- `.dr-dashboard__article-list` — scrollable article container
- `.dr-dashboard__article-row`, `--active` — article items with selection state
- `.dr-dashboard__article-info`, `__article-title`, `__article-meta`, `__article-status` — row content
- `.dr-dashboard__article-actions`, `__action-icon` — hover-revealed action buttons
