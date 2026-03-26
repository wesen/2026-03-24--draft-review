---
Title: Wire unused Redux slices for AuthorApp and ReaderPage
Ticket: DR-007
Status: active
Topics:
    - frontend
    - react
    - redux
    - state-management
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/app/AuthorApp.tsx
      Note: '4 useState to replace: view'
    - Path: frontend/src/chrome/MenuBar.tsx
      Note: 'Optional: openMenu useState could use uiSlice.setActiveMenu'
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: '4 useState to replace: started'
    - Path: frontend/src/store/hooks.ts
      Note: useAppDispatch + useAppSelector typed hooks — ready to use
    - Path: frontend/src/store/readerSlice.ts
      Note: Target — startReading
    - Path: frontend/src/store/store.ts
      Note: Store configured with ui + reader reducers — no changes needed
    - Path: frontend/src/store/uiSlice.ts
      Note: Target — setView
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-26T11:27:24.44710007-04:00
WhatFor: ""
WhenToUse: ""
---


# Wire unused Redux slices for AuthorApp and ReaderPage

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- frontend
- react
- redux
- state-management

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
