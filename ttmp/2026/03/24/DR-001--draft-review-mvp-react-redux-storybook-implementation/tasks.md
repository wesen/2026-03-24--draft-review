# Tasks

## Phase 1: Foundation (Storybook-first)

- [x] 1.1 Scaffold Vite + React 18 + TypeScript project
- [x] 1.2 Set up Storybook 10
- [x] 1.3 Set up Redux Toolkit + RTK Query wiring
- [x] 1.4 Set up MSW (browser + handlers with seeded data)
- [x] 1.5 Create CSS custom properties theme file (mac-os-1.css)
- [x] 1.6 Extract chrome components: TitleBar, MacWindow, MacButton, MacInput, MacTextArea, MenuBar
- [x] 1.7 Extract primitive components: ProgressBar, ReactionBadge, SectionNav, ReactionPicker, StatCard
- [x] 1.8 Write Storybook stories for all chrome + primitive components

## Phase 2: Reader Experience

- [x] 2.1 Define TypeScript types (Article, Section, Reaction, Reader) -- done in Phase 1
- [x] 2.2 RTK Query endpoints + MSW handlers for reader flow -- done in Phase 1
- [x] 2.6 Reader Redux slice (currentSection, readSections, started) -- done in Phase 1
- [ ] 2.3 Build reader components: WelcomeSplash, Paragraph, SectionView, ReaderToolbar
- [ ] 2.4 Build DoneDialog (review completion)
- [ ] 2.5 Build ReaderPage (composed page with keyboard nav)
- [ ] 2.7 Storybook stories for reader components

## Phase 3: Author Dashboard & Article Management

- [ ] 3.1 Author API layer (RTK Query + MSW with seeded ARTICLES data)
- [ ] 3.2 Build Dashboard (stats, readers, reactions chart, draft-killer, recent feedback)
- [ ] 3.3 Build ArticleManager (article list with search/sort)
- [ ] 3.4 Build ArticleEditor (split-pane section editor)
- [ ] 3.5 Build ArticleSettings (sharing, link generation, reaction toggles)
- [ ] 3.6 Storybook stories for author components

## Phase 4: Integration & Polish

- [ ] 4.1 Set up routing (react-router)
- [ ] 4.2 Connect reader + author flows end-to-end
- [ ] 4.3 Error & empty states
- [ ] 4.4 Keyboard shortcuts
