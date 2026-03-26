---
Title: UX UI Design Review
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
    - Path: frontend/src/author/ArticleEditor.tsx
      Note: Editor — where author writes
    - Path: frontend/src/author/ArticleManager.tsx
      Note: All Articles list view
    - Path: frontend/src/author/ArticleReader.tsx
      Note: |-
        Review view — where author reads feedback
        Review view proportion issues documented
    - Path: frontend/src/author/ArticleSettings.tsx
      Note: Settings — sharing and access config
    - Path: frontend/src/author/Dashboard.tsx
      Note: |-
        Dashboard — the main author view
        Dashboard layout issues documented
    - Path: frontend/src/author/InviteDialog.css
      Note: Invite dialog styles
    - Path: frontend/src/chrome/MacButton.css
      Note: |-
        Button styles — the primary interactive element
        Button styles audited
    - Path: frontend/src/chrome/MacWindow.css
      Note: Window chrome — the Mac OS 1 frame
    - Path: frontend/src/chrome/MenuBar.css
      Note: |-
        Menu bar — top system bar
        Menu bar styles — working well
    - Path: frontend/src/reader/Paragraph.css
      Note: |-
        Paragraph styles — reaction chips and hover
        Hover highlighting — gold standard reference
    - Path: frontend/src/reader/Paragraph.tsx
      Note: Paragraph component — hover highlighting reference
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: Reader page — the reading experience
    - Path: frontend/src/reader/WelcomeSplash.tsx
      Note: Welcome screen — first impression for readers
    - Path: frontend/src/theme/mac-os-1.css
      Note: |-
        Design tokens — the foundation of the visual language
        Design tokens — foundation of visual language
ExternalSources: []
Summary: Comprehensive UX/UI design review of the Draft Review application covering visual identity, layout, typography, interaction patterns, information architecture, and actionable improvement recommendations across all views.
LastUpdated: 2026-03-25T17:54:10.318698215-04:00
WhatFor: Guide UX/UI improvements now that the application is feature-complete. Provides a design critique with concrete fixes.
WhenToUse: When planning UX improvements, during design discussions, or when onboarding a designer to understand the current state.
---


# Draft Review — UX/UI Design Review

## Executive Summary

Draft Review is a feedback collection tool for draft articles, built with a distinctive Mac OS System 1 aesthetic. The application has two user flows — **Author** (dashboard, editor, review, settings) and **Reader** (welcome, reading, reacting). The retro design language is compelling and memorable, but the execution has drifted from its original vision in several places, and key interaction patterns are missing or inconsistent. This review covers every view and component, identifies what's working, what's broken, and provides actionable recommendations.

**Overall assessment**: The visual identity is strong and distinctive. The bones are good. What's needed now is a tightening pass — not a redesign, but a disciplined application of the existing design language to every corner of the application. The Mac OS 1 aesthetic should be a *constraint*, not a costume.

---

## 1. Visual Identity Audit

### 1.1 What's Working

The Mac OS 1 theme is genuinely well-executed in several areas:

- **Menu bar** (screenshot 12): Pixel-perfect. The Apple logo glyph, the black-on-white inversion on hover, the dropdown shadow — this feels authentic. The keyboard shortcut display in the dropdown is a nice touch.
- **Scrollbars**: The custom hatched scrollbar thumb is a committed design choice. It sells the retro feel instantly.
- **Desktop background**: The subtle `#aaa` with line pattern is exactly right — not distracting, just textured enough to feel like a CRT.
- **Typography**: The Chicago/Geneva/Monaco stack for UI chrome, Georgia for prose content — this separation is correct. The font sizes (9–20px range) feel appropriate for the retro scale.
- **Border language**: Black 2px borders, 0px radii, hard-pixel shadows (`4px 4px 0 #000`) — this is consistent and deliberate.

### 1.2 What's Drifted

Several areas don't feel like Mac OS 1 anymore:

- **Blue checkboxes** (screenshot 05, settings view): The default browser `<input type="checkbox">` and `<input type="radio">` are rendering with the system's blue accent color. This is the single most visible break from the aesthetic. Mac OS 1 had black-and-white checkboxes with a chunky X mark, not blue fills. **Fix**: Style custom checkbox/radio components with black-and-white pixel art, or at minimum use `accent-color: black`.
- **Emoji as icons** (screenshot 02, "📭" for empty readers): Emoji are full-color and modern — they clash with the 1-bit palette. The reaction icons (★ ? ◎ ♥) are text glyphs which work well. But 📭 and 📝 are color emoji that break the illusion. **Fix**: Replace emoji with bitmap-style 1-bit icons (could be Unicode box-drawing characters, or a tiny SVG icon set in the style of Susan Kare's originals).
- **The "How to leave feedback" card** (screenshot 09, welcome): The colored reaction icons (orange ★, blue ?) break the monochrome palette when they're in the chrome context. They work when they're small and inline in the reader view, but the welcome splash uses them at display size.
- **Article tabs** (screenshot 01, top of dashboard): The status badges ("Draft", "In Review") use a `border: 2px solid` style that's correct, but the selected tab has an inverted background that makes the status badge inside it hard to read. The tab strip itself feels more like a modern tab bar than a Mac OS 1 element.

### 1.3 Color Palette Integrity

Current palette is almost purely monochrome, which is correct:

| Token | Value | Usage | On-brand? |
|-------|-------|-------|-----------|
| `--dr-bg` | `#aaa` | Desktop | Yes |
| `--dr-surface` | `#fff` | Window fills | Yes |
| `--dr-on-surface` | `#000` | Text, borders | Yes |
| `--dr-muted` | `#999` | De-emphasized | Yes |
| `--dr-subtle` | `#555` | Secondary text | Yes |
| `--dr-hover-bg` | `rgba(0,0,0,0.04)` | Hover state | Yes |
| Blue checkboxes | system default | Form controls | **No** |
| Color emoji | system default | Empty states | **No** |

**Recommendation**: The palette is correct. The problem is that browser default form controls and color emoji break out of it. Every interactive element needs to be custom-styled within the monochrome constraint.

---

## 2. Layout and Spatial Composition

### 2.1 Dashboard (screenshots 01, 02)

The dashboard uses a vertical stack layout: article tabs → action bar → stat cards → readers/reactions split → recent feedback → CTA button. This is functional but has spatial issues:

**Problems**:
- **Stat cards are too wide**: Each card gets 25% of the width, but the content is tiny ("◉ 1" / "READERS"). The cards are mostly empty space. On a maximized window, each card could be 200+ px wide showing a single digit.
- **Readers list and Reactions by Section are jammed together**: The two-column layout at this level works, but there's no breathing room. The reader row (avatar + name + progress bar + percentage) is crammed into ~400px and the progress bar is tiny.
- **"Open Full Review →" button is centered but feels orphaned**: It sits alone below "Recent Feedback" with no clear visual connection to what it does. It's the most important action on the page but has no spatial prominence.
- **Article tabs overflow horizontally**: With 4+ articles, the tab strip wraps. There's no truncation, scrolling, or dropdown fallback.
- **Massive empty space below** (screenshot 02): When "No feedback yet" is shown, the bottom 60% of the window is empty white.

**Recommendations**:
1. **Tighten stat cards**: Make them narrower or use a compact inline layout: `◉ 1 readers · ✦ 0 reactions · § 1 sections · ▸ 0% avg progress` — a single line instead of four boxes.
2. **Give the CTA button more prominence**: Move it to the action bar (next to Edit/Share/Invite) or make it a persistent element. An author's main job is reviewing feedback — this should be the primary action.
3. **Collapse empty states**: When there's no feedback, don't show a large empty box. Show a single-line prompt: "No reactions yet. Share your article to start."
4. **Add article tab overflow handling**: Horizontal scroll with fade, or a dropdown after N tabs.

### 2.2 Review View (screenshots 03, 07)

The review view uses a three-panel layout: left sidebar (sections) + main content (article text + reactions panel). This is the right structure, but the proportions are wrong.

**Problems**:
- **Article text gets ~70% of the vertical space, reactions get a fixed 220px sliver at the bottom**. With only one paragraph of content, 90% of the screen is white space while the reactions panel (the thing the author actually came here for) is tiny.
- **The section sidebar is 180px fixed width** — appropriate for the content, but the section items don't show any feedback density information. The author can't see at a glance which sections need attention.
- **Article paragraphs are plain text** with no interactivity (the core issue from the DR-003 ticket).
- **"No reactions for this section yet"** occupies the entire reactions panel with centered gray text — wasted space.

**Recommendations**:
1. **Flip the proportion**: Use a flexible split where the reactions panel gets more space when there are reactions. Consider a 60/40 or even side-by-side layout (text left, reactions right) instead of top/bottom.
2. **Add reaction density to section sidebar**: Small dots or a bar next to each section name showing how much feedback it has.
3. **Make article paragraphs interactive** (see DR-003 implementation plan).
4. **When empty, collapse the reactions panel** to a single line, giving the article text the full height.

### 2.3 Editor View (screenshot 04)

The editor is a left sidebar + main content area: section list on the left, title + content textarea on the right.

**Problems**:
- **The content textarea is unstyled** — it renders in the default monospace font, which is correct for editing, but the textarea height is fixed and doesn't grow with content. Long articles require scrolling inside a fixed box inside a scrolling window — double scroll.
- **"Paragraph breaks" hint text** is useful but easy to miss. It's a single gray line below the textarea.
- **"Preview as Reader" button is in the bottom right** — positioned as a secondary action when it should feel more prominent. The author needs to toggle between editing and previewing.
- **"Danger Zone" with "Delete Section"** is visually disconnected — it sits at the bottom of the left sidebar with a dashed border. The label "DANGER ZONE" is dramatic but the button is disabled (for the last section, correctly), making the drama feel hollow.
- **Section reorder arrows** (▲▼) are small and not clearly interactive — they could be mistaken for decoration.

**Recommendations**:
1. **Make the textarea grow with content**: Use `min-height` + auto-growth or a simple contenteditable div.
2. **Add a live paragraph preview**: Below or beside the textarea, show the paragraph splits as they'll appear to the reader. This helps the author understand the blank-line → paragraph mapping.
3. **Move "Preview as Reader" to the toolbar** or make it a toggle in the title bar.
4. **Soften "Danger Zone"**: Remove the dramatic label. Just show "Delete Section" with a confirmation dialog.

### 2.4 Settings View (screenshot 05)

Settings is a single-column form with sections: Sharing → Access → Reader visibility → Feedback → Status.

**Problems**:
- **Form controls are unstyled browser defaults**: Blue checkboxes and radio buttons. This is the biggest aesthetic violation in the app.
- **No visual grouping**: The section headers ("SHARING", "FEEDBACK", "STATUS") are styled with a line below, which is good, but the controls within each section run together with no padding between groups.
- **"Generate" and "Copy" buttons** for the share link look correct but sit on the same line as the text input, creating a cramped row.
- **Status radio buttons** are on one line ("Draft · In Review · Complete · Archived") which works but doesn't show what changing the status will do. No confirmation, no explanation.
- **"Back" and "Save Settings" buttons** are at the bottom, spaced far apart with no visual relationship. The "Back" button has a ← arrow but looks the same as "Save Settings" except for the fill.

**Recommendations**:
1. **Custom-style all form controls**: Black-and-white checkboxes with a pixel X mark. Radio buttons with a solid black dot. This is non-negotiable for the aesthetic.
2. **Add cards or bordered groups** around each settings section (Sharing, Access, Visibility, etc.) to create clear visual structure.
3. **Add status change confirmation**: When changing from "In Review" to "Complete", show what that means ("Readers will no longer be able to submit reactions").
4. **Make "Save" a primary button** and "Back" a text link or secondary button.

### 2.5 All Articles (screenshot 06)

The article list is a vertical card layout with action buttons on each card.

**Problems**:
- **Cards are tall**: Each article card takes ~120px of vertical space, but the content is sparse. Title + metadata + buttons. With only 4 articles, the page fills the viewport, but with 10+ it would require excessive scrolling.
- **Inconsistent button sets**: The first article (which has readers) shows "Review | Edit | Share". The others show "Review | Edit | Share | Invite →". The conditional "Invite" button makes the rows visually uneven.
- **No visual hierarchy between articles**: All cards look the same regardless of status. "In Review" articles with active readers should feel more alive than stale "Draft" articles.
- **The progress bar** is only shown on articles that have readers, but it's a thin black line that looks like a border or divider, not a progress indicator.
- **"✎" pen icon before each title** is a nice touch but adds visual noise when every article has it.

**Recommendations**:
1. **Compact the card layout**: Consider a table-like layout for the article list. Title | Status | Sections | Readers | Reactions | Actions — all on one line. This is more information-dense and feels more like a Mac OS 1 list view.
2. **Differentiate by status**: Give "In Review" articles a highlighted border or a subtle indicator that there's active work happening.
3. **Unify the button set**: Always show all four buttons but disable "Invite" when not applicable, or always show it (inviting is always useful).
4. **Style the progress bar**: Add a retro Mac-style progress indicator (chunky black fill with white track, bordered).

### 2.6 Invite Dialog (screenshot 08)

The dialog appears as a centered overlay with two sections: Share Link and Reader Email.

**Problems**:
- **No dialog backdrop dimming**: The dialog sits on top of the dashboard but the background remains fully visible and interactive-looking. It's unclear that this is a modal that blocks the background.
- **The dialog border is `2px solid black` with a white background** — correct for the aesthetic, but there's no shadow. Mac OS 1 dialogs had a distinct drop shadow.
- **"Share Link" and "Reader Email" are two distinct flows** crammed into one dialog. The user has to understand both options and choose. The separation between them is a subtle `<hr>`.
- **"Personal Note" textarea** has a great default text but the label says "Personal Note" without explaining that this goes to the reader.

**Recommendations**:
1. **Add dialog shadow**: Use `var(--dr-shadow-dialog)` which is already defined as `6px 6px 0 #000`.
2. **Add a dimmed backdrop**: `background: rgba(0,0,0,0.3)` behind the dialog.
3. **Separate the two flows more clearly**: Use tabs or a toggle ("Share link" vs "Invite by email") instead of stacking both.
4. **Label the note**: "Personal Note (sent to the reader)" or similar.

### 2.7 Reader Welcome Splash (screenshot 09)

The welcome splash is a centered card with article info, author note, feedback instructions, and a "Begin Reading" CTA.

**What's working**:
- The vertical flow is natural: title → author → version → note → how-to → CTA → stats.
- The "How to leave feedback" card is genuinely helpful for onboarding.
- "Begin Reading →" is prominent and clear.
- Article stats at the bottom ("1 sections · ~1360 words · 7 min read") set expectations.

**Problems**:
- **Too much vertical space**: The page is ~60% empty space above the content. The card floats in the center of a huge white window.
- **Color emoji** (📝) breaks the monochrome theme.
- **"click +" in the instructions** is inaccurate — the reader actually clicks text chips like "★ Useful!".
- **The "Draft 1" version badge** is inverted (white on black) but feels like a status badge. On the welcome screen, the reader doesn't need to know the internal version label.

**Recommendations**:
1. **Reduce top padding**: Anchor the welcome content to the upper third instead of dead center.
2. **Replace 📝 with a 1-bit icon**: Use a text glyph or a small SVG that matches the monochrome palette.
3. **Fix the instructions copy**: "Hover any paragraph and click a reaction chip" instead of "click +".
4. **Consider hiding "Draft 1"**: The reader doesn't care about the version. If it's needed, make it subtle.

### 2.8 Reader View (screenshots 10, 11)

The reader view is the most polished part of the application.

**What's working**:
- **Paragraph hover highlighting** is beautiful: left border + subtle background + chips fading in. This is the gold standard that the review view should replicate.
- **Reaction chips** are compact and well-placed below each paragraph.
- **The section progress bar** at the top is clean and informative.
- **Bottom navigation** (← Previous / Finish Review ✓) is clear.
- **"Use arrow keys to navigate · hover paragraphs to react"** hint at the bottom is helpful and unobtrusive.

**Problems**:
- **Reaction chips are always visible** (at 35% opacity), which is correct for discoverability, but they create visual noise when reading long articles. Every paragraph has four chips below it, making the reading experience feel like an evaluation form rather than a document.
- **Code blocks render without styling**: The markdown renderer shows code in a default pre block. No border, no background — it looks like broken text.
- **Headings within the article** (h2, h3, h4, h5 — visible in screenshot 10) don't have reaction chips but they do break the visual flow. They're rendered as plain headings mixed with paragraphs that have chips, creating an uneven rhythm.
- **The section label** "§ UNTITLED SECTION" at the top uses letterspacing and all-caps, which is correct, but it's very close to the toolbar. Needs more top margin.
- **Long paragraphs** (like the "Available Commands" block in screenshot 10) become wall-of-text. No visual break within a paragraph.

**Recommendations**:
1. **Consider hiding chips until hover**: Instead of showing all four chips at 35% opacity under every paragraph, show a single "+" button or nothing, and expand to the four chips on hover. This reduces visual noise significantly.
2. **Style code blocks**: Add a monochrome code block style — black border, light gray background, monospace font. Consistent with Mac OS 1.
3. **Add spacing between headings and paragraph groups**: Headings should feel like section breaks, not just bold text.
4. **Warn authors about long paragraphs**: In the editor, flag paragraphs over ~200 words. These are hard to react to meaningfully.

---

## 3. Typography

### 3.1 Font Stack Assessment

| Context | Font | Size | Verdict |
|---------|------|------|---------|
| Menu bar | `Chicago_12 / ChicagoFLF` | 12px | Excellent — authentic |
| Window titles | `Geneva / Monaco / monospace` | 11px | Good — readable and retro |
| Dashboard labels | `Geneva / Monaco` | 10–12px | Good |
| Article prose (reader) | `Georgia, serif` | (inherited from Prose.css) | Good — readable and warm |
| Article prose (review) | `Georgia, serif` | Same | Good |
| Editor textarea | Default monospace | (browser default) | Adequate but could be Geneva |
| Status badges | `Geneva / Monaco` | 9px | Borderline too small |

### 3.2 Recommendations

- **Lock in Chicago**: The `--dr-font-title` variable specifies `Chicago_12, ChicagoFLF, Geneva, monospace`. If the Chicago font isn't loading (it may not be bundled), the fallback to Geneva is fine, but the app should include the ChicagoFLF webfont for authenticity.
- **Editor textarea**: Use `font-family: var(--dr-font-body)` instead of the browser default monospace.
- **Status badges at 9px**: Ensure these are legible on non-retina displays. 10px minimum.

---

## 4. Interaction Patterns

### 4.1 Navigation Model

The application uses a **single-window model** with view switching managed by Redux state (`uiSlice.view`). The menu bar persists across views. Navigation between views happens through:

1. Button clicks (Edit, Share, Review, All Articles, Back)
2. Menu bar items (File → New Article)
3. Article tab clicks (dashboard)

**Issue**: There's no URL routing for author views. Everything is `http://localhost:5174/`. This means:
- Browser back button doesn't work
- Can't bookmark a specific view
- Can't share a link to a specific article's review view
- Refreshing the page always returns to the dashboard

**Recommendation**: Add URL routing for author views: `/articles`, `/articles/:id/edit`, `/articles/:id/review`, `/articles/:id/settings`. This is a significant improvement for usability.

### 4.2 Button Patterns

Buttons are well-styled (black border, white fill, inverted on press) but used inconsistently:

| Pattern | Example | Consistency |
|---------|---------|-------------|
| Primary action | "Open Full Review →" (inverted, black bg) | Used in some places |
| Secondary action | "Edit" / "Share" (white bg) | Consistent |
| Back navigation | "← Back" / "← Dashboard" (white bg, left arrow) | Consistent |
| Destructive action | "Delete Section" (same as secondary, disabled) | Needs differentiation |

**Issue**: Primary and secondary buttons are distinguishable (inverted vs not), but there's no tertiary/text-link button. "Back" looks the same weight as "Save" even though one is navigation and the other is a commit action.

**Recommendation**: Add a text-link button style for navigation actions (← Back, ← Dashboard). These should look like text with an underline or just a clickable label, not a full bordered button. Save the bordered button for actual actions.

### 4.3 Empty States

The app has several empty states, all inconsistent:

| View | Empty State | Style |
|------|-------------|-------|
| Dashboard / No readers | "📭 No readers yet. Invite your first beta reader!" | Color emoji + text |
| Dashboard / No feedback | "No feedback yet. Share your article to start collecting reactions." | Text only, bordered box |
| Review / No reactions | "No reactions for this section yet." | Text only, centered, small font |
| All Articles / No readers | "No readers yet — invite your first beta reader!" | Italic orange text |

**Recommendation**: Unify all empty states to a single pattern: a centered label in `var(--dr-subtle)` color with no emoji, no box, and an optional action link. Example: "No readers yet. [Invite a reader]"

### 4.4 Hover and Focus States

The reader view has excellent hover states (paragraph highlight, chip fade-in). The author views have minimal hover states:
- Buttons: `background: #eee` on hover — subtle and correct
- Menu items: inverted (black bg, white text) — authentic
- Article tabs: no visible hover state
- Article cards (All Articles): no visible hover state
- Reader list rows: no visible hover state

**Recommendation**: Add hover states to all clickable elements. Use the same `var(--dr-hover-bg)` token. Article cards and list rows should show a left border or background tint on hover, matching the reader paragraph pattern.

---

## 5. Information Architecture

### 5.1 Author Flow Analysis

The current author flow is: Dashboard → (Edit | Review | Settings | All Articles). Each view is a full-screen replacement. The dashboard tries to be both an overview AND a launcher, which creates tension:

- It shows metrics (readers, reactions, sections, progress) that belong on an overview
- It shows action buttons (Edit, Share, Invite) that belong on a toolbar
- It shows recent feedback that belongs in the review view
- It has a huge CTA to "Open Full Review" which suggests the dashboard itself isn't the review

**Recommendation**: Consider whether the dashboard and review view should be merged. The author's primary workflow is: "I have an article in review. I want to see the feedback." The dashboard currently serves as an intermediate step that provides summary stats before the detailed view. If the review view showed the stats in a sidebar or header, the dashboard could be simplified to just the article list.

### 5.2 Reader Flow Analysis

The reader flow is clean: Welcome → Reading → Done. This is correct and doesn't need restructuring.

One improvement: After the reader finishes and sees the "Done" dialog, there's no way to go back and review their own reactions. Consider a "Review your reactions" post-completion view.

---

## 6. Specific Component Issues

### 6.1 Progress Bar (Readers list, Reader toolbar)

The progress bar in the readers list (screenshot 01, next to "Wesen 0%") is a thin black line in a container. It looks like a border, not a progress indicator. The reader toolbar progress bar (screenshot 10, top) is better — it's a thicker bar with visible fill.

**Fix**: Unify to one progress bar style. Use a 4px-high bar with `var(--dr-on-surface)` fill and `var(--dr-border-light)` track, with a 1px black border.

### 6.2 Stat Cards

The stat cards (◉ 1 READERS / ✦ 0 REACTIONS / etc.) use large text for the number and small caps for the label. The icons (◉ ✦ § ▸) are text glyphs which is correct for the aesthetic, but the cards feel oversized for the information they carry.

**Fix**: Compact to an inline stat bar or reduce card width with tighter padding.

### 6.3 Section Sidebar (Review View)

The section sidebar shows numbered section titles with a highlighted background for the selected section. The highlight is a solid black background with white text — correct for Mac OS 1 list selection.

**Missing**: No reaction count or density indicator per section. The author can't see which sections need attention without clicking through each one.

**Fix**: Add a small count badge or bar chart next to each section title showing total reactions.

### 6.4 Reactions by Section (Dashboard)

The "Reactions by Section" panel shows section name + count + a row of reaction type legends (★ Useful! ? Confused ◎ Slow ♥ Love it). When there are zero reactions, this panel is mostly wasted space showing the legend for data that doesn't exist.

**Fix**: Only show the reaction type legend when there are reactions. When empty, collapse to "No reactions" one-liner.

---

## 7. Prioritized Recommendations

### Tier 1 — Must Fix (Visual Integrity)

1. **Style form controls** (checkboxes, radios) to match the monochrome theme. Use `accent-color: black` as a quick fix, or build custom checkbox/radio components.
2. **Replace color emoji** (📭📝) with 1-bit text glyphs or monochrome SVGs.
3. **Add dialog shadow and backdrop dimming** to the Invite dialog.
4. **Fix "click +" copy** in WelcomeSplash to match the actual UI.

### Tier 2 — Should Fix (Layout and Proportion)

5. **Compact stat cards** on the dashboard to reduce visual waste.
6. **Flip the review view proportion** — give the reactions panel more space, or use a side-by-side layout.
7. **Add hover states** to article tabs, article cards, and reader list rows.
8. **Unify empty states** across all views to a single pattern.
9. **Make the "Open Full Review" CTA more prominent** or move it to the action bar.
10. **Style code blocks** in the reader view with monochrome borders and background.

### Tier 3 — Nice to Have (Interaction Polish)

11. **Add URL routing** for author views (Edit, Review, Settings).
12. **Use text-link style** for back navigation buttons instead of bordered buttons.
13. **Add reaction density indicators** to the review view section sidebar.
14. **Add textarea auto-growth** in the editor.
15. **Add loading and error states** for API calls.
16. **Bundle the ChicagoFLF webfont** for authentic typography.

### Tier 4 — Strategic (Information Architecture)

17. **Consider merging Dashboard and Review view** for a streamlined author workflow.
18. **Add compact table-style article list** as an alternative to card layout.
19. **Add post-completion reaction review** for readers.
20. **Add article tab overflow handling** (horizontal scroll or dropdown).

---

## 8. Design Language Summary

The Mac OS 1 aesthetic is the application's strongest differentiator. Here are the rules that should be enforced consistently:

| Principle | Rule |
|-----------|------|
| Color | Monochrome only. No color except reaction type icons in inline context. |
| Borders | 2px solid black. No rounded corners. No gradients. |
| Shadows | Hard pixel shadows: `Npx Npx 0 #000`. No blur. |
| Typography | Chicago for chrome, Geneva for UI text, Georgia for prose. |
| Icons | Text glyphs or 1-bit pixel art. No color emoji. No SVG icon libraries. |
| Interaction | Inversion on press/active (black bg, white text). Subtle gray on hover. |
| Forms | Custom-styled to match the monochrome palette. No browser defaults. |
| Layout | Sharp edges, no padding-heavy whitespace. Information-dense over airy. |
| Patterns | Use hatching and stripe patterns for visual interest, not color. |

---

## Appendix: Screenshots Captured

| # | View | File |
|---|------|------|
| 01 | Dashboard (default article) | `screenshots/01-loading.png` |
| 02 | Dashboard (article selected) | `screenshots/02-dashboard-article-selected.png` |
| 03 | Review view (single section) | `screenshots/03-review-view.png` |
| 04 | Editor view | `screenshots/04-editor-view.png` |
| 05 | Settings view | `screenshots/05-settings-view.png` |
| 06 | All Articles list | `screenshots/06-all-articles.png` |
| 07 | Review view (multi-section) | `screenshots/07-review-multi-section.png` |
| 08 | Invite dialog | `screenshots/08-invite-dialog.png` |
| 09 | Reader welcome splash | `screenshots/09-reader-welcome.png` |
| 10 | Reader reading view | `screenshots/10-reader-view.png` |
| 11 | Reader paragraph hover | `screenshots/11-reader-hover-paragraph.png` |
| 12 | File menu dropdown | `screenshots/12-file-menu.png` |
