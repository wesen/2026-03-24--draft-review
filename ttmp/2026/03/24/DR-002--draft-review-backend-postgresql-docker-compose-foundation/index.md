---
Title: Draft Review Backend - PostgreSQL Docker Compose Foundation
Ticket: DR-002
Status: active
Topics:
    - backend
    - golang
    - glazed
    - postgresql
    - docker-compose
    - api
    - auth
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: frontend/src/mocks/db.ts
      Note: Current in-memory data model used to infer the first persistent entities
    - Path: frontend/src/mocks/handlers.ts
      Note: Current mock HTTP contract that the real backend should replace incrementally
    - Path: imports/draft-review-screen-spec.md
      Note: |-
        Primary product specification that defines backend requirements beyond the current frontend implementation
        Primary product specification that defines the backend domain and flows
ExternalSources: []
Summary: Ticket workspace and detailed implementation guide for building the first Draft Review backend in Go with Glazed command verbs, Clay-style DB configuration, PostgreSQL, docker compose, durable APIs, and phased delivery guidance.
LastUpdated: 2026-03-24T20:04:39.168007321-04:00
WhatFor: Planning the first real Go backend for Draft Review, including Glazed CLI verbs, Clay-style PostgreSQL configuration, local Docker Compose development, authentication, article/review APIs, and a phased implementation path.
WhenToUse: Use this ticket when scoping, scaffolding, or implementing the Draft Review Go backend and when onboarding engineers who need to understand how the frontend maps to backend responsibilities.
---




# Draft Review Backend - PostgreSQL Docker Compose Foundation

## Overview

This ticket captures the backend plan for Draft Review. The current repository contains a React/Vite frontend with MSW-backed mock APIs and an extensive screen specification, but no real server, database, authentication flow, or persistence layer. The goal of this ticket is to turn that frontend contract into an implementation-ready Go backend design centered on Glazed command verbs, Clay-style SQL configuration, PostgreSQL, and a local `docker compose` workflow.

The primary deliverables in this ticket are:

- a detailed architecture and implementation guide for an intern,
- a chronological diary of the research and documentation work,
- ticket bookkeeping that relates the design back to the frontend files and screen specification that drove the recommendations.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- backend
- golang
- glazed
- postgresql
- docker-compose
- api
- auth

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Deliverables

- `design-doc/01-draft-review-backend-architecture-and-implementation-guide.md`
- `reference/01-backend-ticket-diary.md`

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
