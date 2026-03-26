# Changelog

## 2026-03-25

- Initial workspace created


## 2026-03-25

Created ticket, design doc, and diary. Completed full codebase analysis of all 14+ frontend components. Identified core gap: ArticleReader renders reactions as flat list without paragraph context despite Reaction.paragraphId being available. Produced 4-phase implementation plan covering paragraph-linked submissions, review view polish, cross-view consistency, and settings wiring. Documented 20+ UX issues across all views.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.tsx — Primary analysis target


## 2026-03-25

Completed comprehensive UX/UI design review. Captured 12 screenshots of every view via Playwright. Identified 20 prioritized issues across 4 tiers: visual integrity breaks (blue checkboxes, color emoji), layout proportion problems (review view split, stat card sizing), interaction gaps (missing hover states, no URL routing), and strategic architecture improvements (dashboard/review merge). Core design language (Mac OS 1 monochrome, 2px borders, hard shadows) is solid — issue is inconsistent application.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/ttmp/2026/03/25/DR-003--review-view-ux-paragraph-linked-submissions-and-ui-tightening/design/02-ux-ui-design-review.md — Full design review document


## 2026-03-25

Implemented all 15 tasks across 4 phases and 4 commits. Phase 1-2 (2153022, e699eca): paragraph-linked submissions with hover highlighting, grouped reactions, timestamps, badges. Phase 3 (01ac26a): Dashboard paragraph context, useEffect fix, WelcomeSplash copy, loading states. Phase 4 (430db43): monochrome form controls, text glyphs, coming-soon badges. All tasks checked off.

### Related Files

- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.css — Review paragraph hover CSS and grouped reactions styles (commit 2153022)
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/author/ArticleReader.tsx — Core feature — paragraph-linked submissions (commit 2153022)
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/index.css — Global accent-color: black (commit 430db43)
- /home/manuel/code/wesen/2026-03-24--draft-review/frontend/src/reader/ReaderPage.tsx — useEffect dep array fix (commit 01ac26a)

