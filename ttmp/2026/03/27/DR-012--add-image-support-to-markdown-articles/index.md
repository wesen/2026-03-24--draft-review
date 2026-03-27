---
Title: Add image support to markdown articles
Ticket: DR-012
Status: active
Topics:
    - backend
    - frontend
    - editor
    - markdown
    - images
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-27T08:00:48.302650473-04:00
WhatFor: Scope the work required to add intentional image support to Draft Review markdown articles, including editor, API, rendering, review semantics, and asset handling.
WhenToUse: Use when planning or implementing article image support, deciding between URL-only images and managed uploads, or onboarding a new engineer onto the article pipeline.
---

# Add image support to markdown articles

## Overview

This ticket scopes the work needed to support images inside Draft Review articles.
The current system already stores section bodies as markdown and renders them through
`react-markdown`, so image syntax is partially possible today. However, that support
is accidental rather than productized: the public API exposes `paragraphs[]` instead
of canonical markdown, the editor is built around blank-line paragraph splitting,
reactions and analytics are keyed as paragraph-level anchors, and there is no image
upload or asset lifecycle.

The primary deliverable is an intern-facing design and implementation guide that
maps the current system, explains the architectural constraints, and recommends a
phased rollout:

1. explicit support for markdown image syntax in section bodies,
2. optional managed image uploads and asset records,
3. later cleanup of the paragraph-only anchor model if richer content blocks are needed.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- backend
- frontend
- editor
- markdown
- images

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
