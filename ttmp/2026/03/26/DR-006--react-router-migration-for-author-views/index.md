---
Title: React Router migration for author views
Ticket: DR-006
Status: active
Topics:
    - frontend
    - react
    - routing
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/App.tsx
      Note: BrowserRouter already set up — add nested author routes
    - Path: frontend/src/app/AuthorApp.tsx
      Note: 5 useState calls to replace with Redux + Router
    - Path: frontend/src/reader/ReaderPage.tsx
      Note: 4 useState calls replaceable with readerSlice
    - Path: frontend/src/store/readerSlice.ts
      Note: Pre-built Redux slice — startReading
    - Path: frontend/src/store/uiSlice.ts
      Note: Pre-built Redux slice — setView
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-26T11:08:18.130007328-04:00
WhatFor: ""
WhenToUse: ""
---


# React Router migration for author views

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
- routing

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
