---
Title: Immediate review persistence and automatic review completion
Ticket: DR-009
Status: active
Topics:
    - frontend
    - react
    - ux
    - backend
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Track the reader-flow change where feedback should be saved immediately and review completion should no longer rely on a manual Finish Review action.
LastUpdated: 2026-03-26T13:06:53.016261384-04:00
WhatFor: Plan and track the change to immediate review persistence and automatic completion behavior.
WhenToUse: Use this ticket when implementing or reviewing the reader-flow persistence and completion UX.
---

# Immediate review persistence and automatic review completion

## Overview

This ticket tracks the reader-flow change where feedback should be persisted immediately and the review should no longer depend on a manual `Finish Review` action that is easy to miss.

The main goal is to make the reader experience more forgiving: if someone leaves feedback and closes the page, the author should still see that work without depending on the reader to press an extra completion control.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- frontend
- react
- ux
- backend

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
