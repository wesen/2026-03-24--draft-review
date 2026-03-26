---
Title: 'Review View UX: Paragraph-Linked Submissions and UI Tightening'
Ticket: DR-003
Status: active
Topics:
    - frontend
    - react
    - ux
    - ui
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/app/AuthorApp.tsx
      Note: Top-level author routing and state management
    - Path: frontend/src/author/ArticleReader.css
      Note: Review view styles needing paragraph hover additions
    - Path: frontend/src/author/ArticleReader.tsx
      Note: Primary target — author review view with disconnected reaction list
    - Path: frontend/src/author/Dashboard.tsx
      Note: Dashboard recent feedback also missing paragraph context
    - Path: frontend/src/reader/Paragraph.css
      Note: Hover CSS patterns to replicate in review view
    - Path: frontend/src/reader/Paragraph.tsx
      Note: Reference implementation for paragraph hover highlighting
    - Path: frontend/src/types/reaction.ts
      Note: Reaction type already has paragraphId field
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-25T11:55:00.091882939-04:00
WhatFor: ""
WhenToUse: ""
---


# Review View UX: Paragraph-Linked Submissions and UI Tightening

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
- ux
- ui

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
