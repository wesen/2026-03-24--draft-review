---
Title: Implementation Diary
Ticket: DR-004
Status: active
Topics:
    - frontend
    - react
    - ux
    - ui
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Step-by-step diary of DR-004 implementation: layout flip, backdrop fix, hover states, stat compaction, and font size bump."
LastUpdated: 2026-03-26T13:30:00-04:00
WhatFor: "Record what changed, why, what worked, and how to verify."
WhenToUse: "When reviewing DR-004 changes or debugging regressions."
---

# DR-004 Implementation Diary

## 2026-03-26

### Session: Implementing all 5 work items

**Context**: DR-004 ticket was created with an implementation guide covering 5 work items from the DR-003 design review. All items are independent CSS/layout changes with no backend work.

---

### WI-1: Review View Side-by-Side Layout Flip

**What changed**:
- `ArticleReader.css`: Flipped `.dr-article-reader__main` from `flex-direction: column` to `row`
- Text panel border changed from `border-bottom` to `border-right`
- Reactions panel changed from height-constrained (min 150px, max 50%) to width-constrained (340px fixed)
- Added `.dr-article-reader__reactions-header` — mirrors the sidebar "SECTIONS" header for visual symmetry
- Filter tabs compacted: changed from `{icon} {label} ({count})` to `{icon} {count}` with `title` attribute for full label
- Added `flex-wrap: wrap` to filter tabs for narrow panel
- Paragraph excerpts compacted: smaller font, `white-space: nowrap`, `text-overflow: ellipsis`

**Why**: The author comes to the review view to read feedback, but the old column layout gave most of the screen to the article text (which the author wrote) and cramped the reactions at the bottom. Side-by-side makes both visible simultaneously, and bidirectional hover highlighting works much better because the paragraph and its reactions are on screen at the same time.

**What worked well**: The CSS-only approach — the JSX structure barely changed (only the filter tab render and the new header). The `scrollIntoView` for bidirectional hover still works correctly in the row layout.

**Verification**: Loaded the review view in Playwright. Text scrolls on the left, reactions panel on the right. Filter tabs fit in the 340px panel with wrapping. No reactions in the test data so couldn't verify bidirectional hover, but the code path is unchanged.

---

### WI-2: Invite Dialog Backdrop and Border-Radius

**What changed**:
- `InviteDialog.css`: Backdrop opacity from `0.15` to `0.35`
- `InviteDialog.css`: `border-radius` from `4px` to `0`

**Why**: The 15% backdrop was nearly invisible — the dialog looked like it was floating on the page without dimming the background. 35% matches standard Mac OS dialog dimming. The 4px border-radius violated the Mac OS 1 aesthetic (everything should be 0px corners).

**Verification**: Opened the invite dialog via Playwright. The backdrop is clearly visible now — the background dims noticeably. Sharp corners look correct.

---

### WI-3: Hover States

**What changed**: Added `transition: background 0.1s` and `:hover { background: var(--dr-hover-bg) }` to:
- `Dashboard.css`: `.dr-dashboard__tab`, `.dr-dashboard__reader-row`, `.dr-dashboard__section-row`, `.dr-dashboard__feedback-item`
- `ArticleManager.css`: `.dr-manager__card`
- `ArticleReader.css`: `.dr-article-reader__section-item` (with override for `--active` state to keep the inverted style)

**Why**: Interactive elements without hover feedback feel dead. Every clickable row/card/tab should respond to the cursor. Using the existing `--dr-hover-bg` token (4% black overlay) keeps it consistent with the paragraph hover in both reader and review views.

**What was tricky**: The active section sidebar item inverts (white on black). Without the `--active:hover` override, hovering it would add a light overlay on top of the black background, looking wrong. The override ensures the active state stays fully inverted on hover.

---

### WI-4: Compact Stat Cards

**What changed**:
- `Dashboard.tsx`: Removed `StatCard` import and the grid of four `<StatCard>` components. Replaced with a single `div.dr-dashboard__stat-bar` containing four `<span>` elements with icon + value + label inline.
- `Dashboard.css`: Replaced `.dr-dashboard__stats` grid with `.dr-dashboard__stat-bar` flex layout — `flex-wrap: wrap`, `gap: 16px`, monospace font, dotted bottom border.

**Why**: Four stat cards took ~80px of vertical space for single-digit metrics. The inline bar takes ~30px and frees vertical space for the more useful readers panel, reactions chart, and feedback list below.

**Note**: Left the `StatCard` component in `primitives/` — it's a clean, generic component that could be useful elsewhere.

---

### WI-5: Font Size Bump

**What changed**:

**Design tokens** (`mac-os-1.css`):
| Token | Old | New |
|-------|-----|-----|
| `--dr-font-size-xs` | 9px | 11px |
| `--dr-font-size-sm` | 10px | 12px |
| `--dr-font-size-base` | 11px | 14px |
| `--dr-font-size-md` | 12px | 15px |
| `--dr-font-size-lg` | 14px | 17px |
| `--dr-font-size-xl` | 18px | 21px |
| `--dr-font-size-2xl` | 20px | 26px |
| `--dr-titlebar-height` | 22px | 28px |
| `--dr-menubar-height` | 20px | 26px |
| `--dr-close-box-size` | 12px | 14px |

**Hardcoded sizes bumped** across 14 CSS files:
- `Prose.css`: body 13→17px, h1 20→26px, h2 16→22px, h3 14→18px, code/pre/table 11→14px
- `Paragraph.css`: chips 9→11px, badges 8→10px
- `ArticleReader.css`: heading 16→20px, paragraph 13→17px, max-width 540→620px, badge 8→10px
- `SectionView.css`: max-width 540→640px, title 16→20px
- `ReaderPage.css`: max-width 540→640px (nav + hint)
- `WelcomeSplash.css`: max-width 420→500px
- All other 8px badges/status indicators → 10px
- All 13px icon sizes → 16px

**Why**: The original 9-20px range was scaled for pixel-art aesthetic but was too small for actual reading. User compared to Substack (~18px body) and requested "normal sizes". The 17px body is close to Substack while keeping a slightly tighter feel for the retro aesthetic.

**What was tricky**: Max-width constraints. With 17px text, the old 540px max-width on paragraphs/sections was too narrow — text felt cramped. Widened reader content to 640px and review paragraphs to 620px. The welcome splash also needed widening (420→500px).

**Verification**: Took screenshots of dashboard, review view, settings, reader welcome, and reader content view. All fonts are proportional and readable. Menu bar and title bar accommodate the larger text without overflow. Reaction chips in the reader view are legible. The overall feel is much closer to a modern web reading experience while retaining the Mac OS 1 aesthetic.

---

### Commit Strategy

Two commits:
1. `0266233` — WI-1 through WI-4: layout flip, backdrop fix, hover states, stat compaction (6 files, +87 -31)
2. `6afc7c6` — WI-5: font size bump (15 files, +38 -38)

Split WI-5 into its own commit because it touches 15 files and is a distinct concern (visual sizing) from the layout/interaction changes in WI-1–4.

---

### How to Review

1. **Dashboard**: Check stat bar is compact, hover states work on tabs/rows/cards
2. **Review view**: Verify side-by-side layout — text left, reactions right, both scroll independently. "REACTIONS" header visible. Filter tabs show icon+count only.
3. **Invite dialog**: Click "+ Invite Reader" — backdrop should dim the background noticeably. Dialog corners should be sharp.
4. **Reader view**: Navigate to a reader link. Body text should be ~17px. Content area should be wider (~640px). Reaction chips readable.
5. **Article list**: Check hover on article cards
6. **Settings**: Check font sizes on labels, checkboxes, COMING SOON badges

---

### Post-Implementation Fixes

**Broken unicode**: The menu bar logo used U+F8FF (Apple Private Use Area — the Apple logo). This only renders on macOS with Apple fonts; on Linux it shows as a □ box. Replaced with ⌘ (U+2318, command key symbol) which renders on all platforms and is still recognizable as a Mac icon.

**Tight margins**: After bumping font sizes, the padding/margins that were proportional at 9-13px became too cramped at 14-17px. Fixed by:
- Bumping spacing tokens: `--dr-space-sm` 6→8px, `--dr-space-md` 10→12px, `--dr-space-lg` 16→18px, `--dr-space-xl` 24→28px, `--dr-space-2xl` 32→40px
- Increasing padding on: menu bar items (6→10px), title bar (4→8px padding, 4→6px gap), dashboard (16→24px outer padding), panel headers (4→6px), panel body (8→10px), reader rows, feedback items, sidebar headers/items, filter tabs, reaction list, article manager toolbar/list/cards
- Reader section view padding bumped 20→24px top/bottom

Commit: `24c5c8c`
