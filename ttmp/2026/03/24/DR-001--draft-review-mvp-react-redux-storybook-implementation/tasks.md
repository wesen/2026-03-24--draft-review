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
- [x] 2.3 Build reader components: WelcomeSplash, Paragraph, SectionView, ReaderToolbar
- [x] 2.4 Build DoneDialog (review completion)
- [x] 2.5 Build ReaderPage (composed page with keyboard nav)
- [x] 2.7 Storybook stories for reader components

## Phase 3: Author Dashboard & Article Management

- [x] 3.1 Author API layer (RTK Query endpoints for articles, readers, reactions)
- [x] 3.2 Build Dashboard (stats, readers, reactions chart, draft-killer, recent feedback)
- [x] 3.6 Storybook stories for author components (Dashboard, ArticleReader, InviteDialog)
- [x] 3.x Build ArticleReader (author review view with section sidebar + reaction filters)
- [x] 3.x Build InviteDialog (reader invitation modal)
- [x] 3.x Wire App.tsx with MSW data (dashboard -> article reader flow)
- [x] 3.3 Build ArticleManager (article list with search/sort/cards/archived section)
- [x] 3.4 Build ArticleEditor (split-pane section editor with reorder/add/delete)
- [x] 3.5 Build ArticleSettings (sharing link gen, access mode, reaction toggles, status)
- [x] 3.x Extended MSW handlers (create/update article, generate share token)
- [x] 3.x Extended RTK Query (createArticle, updateArticle, generateShareToken)
- [x] 3.x Wired ArticleManager/Editor/Settings into AuthorApp with full navigation
- [x] 3.x Storybook stories for ArticleEditor, ArticleSettings, ArticleManager

## Phase 4: Integration & Polish

- [x] 4.1 Set up routing (react-router) — `/r/:token` for readers, `/*` for author
- [x] 4.2 Connect reader + author flows end-to-end — AuthorApp + ReaderApp
- [x] 4.2b Reader token resolution via readerApi + error/loading states
- [ ] 4.3 Error & empty states (further polish)
- [ ] 4.4 Keyboard shortcuts (further polish)
