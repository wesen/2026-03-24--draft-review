---
Title: Review View UX Analysis and Implementation Plan
Ticket: DR-003
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
      Note: "Primary target: author's review view, currently missing paragraph-reaction links"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.css
      Note: "Review view styles, needs paragraph hover/highlight additions"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.tsx
      Note: "Reader paragraph component with hover highlighting — reference implementation"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/Paragraph.css
      Note: "Paragraph hover styles to replicate in review view"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/Dashboard.tsx
      Note: "Dashboard recent feedback — also missing paragraph context"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/app/AuthorApp.tsx
      Note: "Top-level author routing and state management"
    - Path: /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/types/reaction.ts
      Note: "Reaction type definition — already has paragraphId field"
ExternalSources: []
Summary: "Analysis of the ArticleReader review view gaps and comprehensive UX/UI tightening plan for the Draft Review application."
LastUpdated: 2026-03-25T11:55:04.038182587-04:00
WhatFor: "Guide implementation of paragraph-linked submissions in the review view and prioritize UX improvements across the application."
WhenToUse: "When implementing the review view overhaul or picking up UX improvements."
---

# Review View UX Analysis and Implementation Plan

## 1. Problem Statement

When an author reviews submissions for a section in `ArticleReader`, reactions are displayed as a flat list filtered by section. There is **no link between a reaction card and the paragraph it was left on**. The author cannot tell which paragraph a reaction refers to without reading the comment text and guessing. Additionally, the article text in the review view has **no hover highlighting** — unlike the reader view where paragraphs light up with a left border and background tint on hover.

This makes the review experience disconnected: the author sees text above and reactions below, but the two don't talk to each other.

## 2. Current State Analysis

### 2.1 ArticleReader.tsx — What exists

The review view (lines 34–163) has a three-panel layout:

1. **Left sidebar** (180px): Section list with reaction counts per type
2. **Top main**: Article text rendered via `<Prose>` — plain, no interactivity
3. **Bottom main**: Reactions panel (fixed 220px height) with filter tabs

**Critical gap**: The `Reaction` object already carries a `paragraphId` (e.g., `s1-p0`), but ArticleReader never uses it. Reactions are filtered by `sectionId` only (line 27–29), then rendered as a flat list (lines 142–157).

The article paragraphs are rendered as plain `<Prose>` elements (lines 93–97) with no data attributes, no hover states, and no way to identify which paragraph is which.

### 2.2 Reader View — What the user expects

In the reader view (`Paragraph.tsx` + `Paragraph.css`), paragraphs have:
- **Hover highlight**: `border-left: 3px solid` + `background: var(--dr-hover-bg)` on hover
- **Commented state**: Paragraphs with reactions always show the left border
- **Reaction count badge**: Small badge showing how many reactions are on this paragraph
- **Smooth transitions**: 0.1s on background and border-color

The user explicitly asked for "highlight the hovered section of the document like in the reader view" — meaning this same visual language should carry over to the review view.

### 2.3 Data Model — Already supports paragraph linking

```typescript
// Reaction type already has paragraphId
interface Reaction {
  id: string;
  articleId: string;
  sectionId: string;
  paragraphId: string;  // e.g., "s1-p0" — THIS IS THE KEY
  readerId: string;
  readerName: string;
  type: ReactionType;
  text: string;
  createdAt: string;
}
```

No backend changes are needed. The `paragraphId` is already stored and returned in the API response.

## 3. Implementation Plan — Paragraph-Linked Submissions

### 3.1 Make paragraphs interactive in ArticleReader

Replace the plain `<Prose>` render (ArticleReader.tsx lines 93–97) with paragraph containers that:

1. Have a `data-paragraph-id` attribute matching the `{sectionId}-p{index}` convention
2. Support hover highlighting (reuse the `.dr-para` visual language)
3. Show a reaction count badge on hover (or always, if reactions exist)
4. Support click-to-filter: clicking a paragraph filters the reactions panel to that paragraph only

**Approach**: Create a new `ReviewParagraph` component (or reuse parts of the reader `Paragraph` component) that accepts a `paragraphId`, `reactionCount`, and `isHighlighted` prop. The component renders the prose text with hover styles but WITHOUT the reaction input chips (those belong in the reader view only).

```tsx
// Sketch — ReviewParagraph in ArticleReader
<div
  className={`dr-review-para ${isHighlighted ? "dr-review-para--active" : ""} ${reactionCount > 0 ? "dr-review-para--has-reactions" : ""}`}
  onMouseEnter={() => setHoveredParagraph(paragraphId)}
  onMouseLeave={() => setHoveredParagraph(null)}
  onClick={() => setFilteredParagraph(paragraphId)}
  data-paragraph-id={paragraphId}
>
  <Prose>{text}</Prose>
  {reactionCount > 0 && (
    <span className="dr-review-para__badge">{reactionCount}</span>
  )}
</div>
```

### 3.2 Link reaction cards to paragraphs

In the reactions panel, add two features:

1. **Paragraph snippet**: Show a truncated excerpt (first ~60 chars) of the paragraph text above or beside each reaction card, so the author knows what the reader was reacting to.

2. **Bidirectional highlighting**:
   - When hovering a **paragraph** in the text: highlight matching reaction cards in the panel (add a class or border)
   - When hovering a **reaction card**: highlight the corresponding paragraph in the text (scroll into view if needed)

**State additions to ArticleReader**:
```tsx
const [hoveredParagraph, setHoveredParagraph] = useState<string | null>(null);
const [filteredParagraph, setFilteredParagraph] = useState<string | null>(null);
```

**Filtering logic update** (currently lines 30–32):
```tsx
// Add paragraph filter on top of existing type filter
let filtered = reactionFilter
  ? sectionReactions.filter((r) => r.type === reactionFilter)
  : sectionReactions;
if (filteredParagraph) {
  filtered = filtered.filter((r) => r.paragraphId === filteredParagraph);
}
```

### 3.3 Group reactions by paragraph

Instead of a flat reaction list, group reactions by `paragraphId`:

```
┌─────────────────────────────────────┐
│ ¶ "The key insight is that..."      │  ← paragraph excerpt
│  ★ Alice: "I can use this"          │
│  ? Bob: "What does 'key' mean?"     │
├─────────────────────────────────────┤
│ ¶ "In contrast to prior work..."    │
│  ◯ Carol: "This drags"              │
└─────────────────────────────────────┘
```

This gives immediate spatial context. The paragraph excerpt acts as a visual anchor.

### 3.4 CSS additions to ArticleReader.css

```css
/* Review paragraph — hover highlight like reader view */
.dr-review-para {
  position: relative;
  border-left: 3px solid transparent;
  padding-left: 8px;
  margin-left: -11px;
  padding-bottom: 4px;
  margin-bottom: 4px;
  transition: background 0.1s, border-color 0.1s;
  cursor: pointer;
}

.dr-review-para:hover,
.dr-review-para--active {
  border-left-color: var(--dr-on-surface);
  background: var(--dr-hover-bg);
}

.dr-review-para--has-reactions {
  border-left-color: var(--dr-on-surface);
  background: rgba(0, 0, 0, 0.03);
}

/* Reaction card highlight when its paragraph is hovered */
.dr-article-reader__reaction-card--highlighted {
  background: var(--dr-hover-bg);
  border-color: var(--dr-on-surface);
}

/* Paragraph excerpt in reaction group */
.dr-article-reader__para-excerpt {
  font-size: var(--dr-font-size-sm);
  color: var(--dr-subtle);
  font-style: italic;
  padding: 4px 10px;
  border-left: 2px solid var(--dr-border-light);
  margin-bottom: 4px;
}
```

## 4. UX/UI Tightening Audit

### 4.1 ArticleReader (Review View) — HIGH PRIORITY

| Issue | Severity | Fix |
|-------|----------|-----|
| Reactions panel fixed at 220px — too small for many reactions, too large when empty | Medium | Make panel resizable or use a flex split with min/max heights |
| No paragraph context on reaction cards | High | Add paragraph excerpt (see §3) |
| No visual feedback when hovering article text | High | Add hover highlighting (see §3) |
| No "clear filter" affordance when filtered by paragraph | Medium | Add a "Showing reactions for paragraph X — Clear" bar |
| Section sidebar doesn't show paragraph-level density | Low | Add small heatmap dots or counts per paragraph |
| Reaction card doesn't show timestamp | Low | Add relative time (e.g., "2h ago") |

### 4.2 Dashboard — MEDIUM PRIORITY

| Issue | Severity | Fix |
|-------|----------|-----|
| "Recent Feedback" shows section title but no paragraph context | Medium | Add truncated paragraph excerpt: `"on ¶ 'The key insight...'" ` |
| Draft-killer alert identifies section but not specific paragraph | Low | Add "worst paragraph" detail when data supports it |
| Article tabs don't show article status visually (just tiny text badge) | Low | Use distinct background or icon for in_review vs draft |
| No empty-state illustration — just emoji | Low | Consider replacing emoji with a retro Mac icon |

### 4.3 ArticleEditor — MEDIUM PRIORITY

| Issue | Severity | Fix |
|-------|----------|-----|
| No indication of existing feedback while editing | Medium | Show reaction heat alongside paragraphs (e.g., colored gutter marks) |
| Content textarea is free-form with no paragraph boundary preview | Low | Could show live paragraph split preview below textarea |
| "Preview as Reader" doesn't show existing reactions in preview | Low | Pass reactions to ReaderPage in preview mode |
| No autosave or dirty indicator | Medium | Add "Unsaved changes" indicator + autosave timer |

### 4.4 ArticleSettings — LOW PRIORITY

| Issue | Severity | Fix |
|-------|----------|-----|
| Most settings are client-only (access mode, reader visibility checkboxes) | Medium | Wire to backend or clearly mark as "coming soon" |
| No confirmation when resetting share link | Medium | Add "This will invalidate existing links" warning |
| Copy button gives no feedback on failure | Low | Handle clipboard API rejection |

### 4.5 Reader View — LOW PRIORITY (already polished)

| Issue | Severity | Fix |
|-------|----------|-----|
| WelcomeSplash howto says "click +" but chips show text labels, not "+" | Low | Update copy to "click a reaction chip" |
| Keyboard navigation hook re-subscribes on every render (no dep array on useEffect) | Medium | Add `[isFirst, isLast, currentSectionId]` to dep array |
| Word count estimate uses `paragraphs.length * 40` — very rough | Low | Count actual words in paragraph text |
| No way to navigate back to a previous section's reactions | Low | Consider letting reader view their own reaction history |

### 4.6 General / Cross-Cutting

| Issue | Severity | Fix |
|-------|----------|-----|
| No loading states for API calls (articles, reactions, readers) | Medium | Add skeleton screens or loading spinners in Mac chrome style |
| No error UI when API calls fail | Medium | Add error alert with retry button |
| Inline styles in several components where CSS classes exist | Low | Move to CSS classes for consistency |
| No keyboard shortcuts documented or implemented beyond arrow navigation | Low | Add ⌘K command palette or shortcut overlay |
| `goNext` / `goPrev` in ReaderPage are not memoized, causing hook churn | Low | Wrap in useCallback |
| Dashboard `onSelectArticle` can navigate to review view for a different article than the one whose readers/reactions are loaded | Medium | Ensure selectedArticle state and query params stay in sync |

## 5. Suggested Implementation Order

### Phase 1 — Core Feature (paragraph-linked submissions)
1. Add `hoveredParagraph` and `filteredParagraph` state to ArticleReader
2. Replace `<Prose>` paragraph render with interactive `ReviewParagraph` elements
3. Add paragraph hover CSS to ArticleReader.css
4. Group reaction list by paragraphId with paragraph excerpts
5. Wire bidirectional hover highlighting (paragraph ↔ reaction card)
6. Add click-to-filter and "clear filter" for paragraph filtering

### Phase 2 — Review View Polish
7. Make reactions panel height flexible (min 150px, flex grow)
8. Add relative timestamps to reaction cards
9. Add reaction count badges on paragraphs in the article text
10. Improve section sidebar with paragraph-level reaction density

### Phase 3 — Cross-View Consistency
11. Add paragraph context to Dashboard "Recent Feedback" items
12. Add feedback heat indicators to ArticleEditor paragraphs
13. Fix ReaderPage useEffect dependency array
14. Add loading/error states for API calls
15. Update WelcomeSplash copy ("click a reaction chip" not "click +")

### Phase 4 — Settings & Polish
16. Wire ArticleSettings to backend (or mark as "coming soon")
17. Add share link reset confirmation dialog
18. Add autosave/dirty indicator to ArticleEditor
19. Clean up inline styles → CSS classes

## 6. Key Architectural Decision: Inline vs. Panel

The current ArticleReader uses a **text-above / reactions-below** split. An alternative is **inline reactions** (like Google Docs comment threads attached to specific text). The inline approach would be more natural for paragraph-level feedback but is significantly more complex to implement (comment positioning, overflow handling, overlapping comments).

**Recommendation**: Keep the panel layout but add strong bidirectional linking. The panel approach works well for the Mac OS 1 aesthetic (feels like a mail/inspector pattern). The key improvement is making the panel and the text talk to each other through hover and click interactions. Inline annotations can be a future enhancement.
