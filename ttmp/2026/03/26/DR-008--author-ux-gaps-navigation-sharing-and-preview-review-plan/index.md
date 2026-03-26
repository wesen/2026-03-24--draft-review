---
Title: Author UX gaps, navigation, sharing, and preview review plan
Ticket: DR-008
Status: active
Topics:
    - frontend
    - react
    - ui
    - ux
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Analysis ticket for the remaining author UX gaps around delete article, real URL routing, share-link generation, reader identity modes, and preview-test reviews.
LastUpdated: 2026-03-26T12:35:00-04:00
WhatFor: Plan the next user-visible author workflow fixes and capture adjacent code-quality concerns discovered during investigation.
WhenToUse: Use this ticket when implementing the next author UX slice or onboarding a developer to the remaining navigation and sharing work.
---

# Author UX gaps, navigation, sharing, and preview review plan

## Overview

This ticket documents the remaining high-impact author UX gaps in Draft Review and proposes a structured implementation plan. The current focus areas are article deletion, proper URL-based author navigation, the broken share/generate-link modal, richer tracked reader-link creation, and optional persisted preview reviews for testing.

A second companion document in this ticket captures broader code-quality findings that surfaced during the investigation so they can turn into future cleanup work instead of being lost.

## Key Links

- **Main Design Doc**: `design-doc/01-author-ux-gap-analysis-and-implementation-guide.md`
- **Rolling Code Review**: `reference/01-rolling-code-review-and-improvement-backlog.md`
- **Investigation Diary**: `reference/02-investigation-diary.md`
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- frontend
- ux
- routing
- sharing

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
