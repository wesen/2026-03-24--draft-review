---
Title: Long-term durable auth sessions with server-side refresh-token renewal
Ticket: DR-011
Status: active
Topics:
    - auth
    - backend
    - oidc
    - keycloak
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-26T15:35:30.52951536-04:00
WhatFor: Track the future-state auth/session architecture beyond the medium-term implementation slice.
WhenToUse: Use when reasoning about refresh-token renewal, revocation, and final durable auth design.
---

# Long-term durable auth sessions with server-side refresh-token renewal

## Overview

This ticket captures the final auth target for Draft Review: opaque server-side app
sessions backed by encrypted refresh-token storage, silent renewal, revocation, and
better operator control. It is deliberately separate from DR-010, which is the
medium-term implementation ticket for moving to opaque server-side sessions without
refresh-token renewal yet.

Primary design document:

- [01-long-term-durable-auth-sessions-with-server-side-refresh-token-renewal-guide.md](./design-doc/01-long-term-durable-auth-sessions-with-server-side-refresh-token-renewal-guide.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- auth
- backend
- oidc
- keycloak

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
