---
Title: Author session expiry investigation and durable OIDC session design
Ticket: DR-010
Status: active
Topics:
    - auth
    - backend
    - oidc
    - keycloak
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted/main.tf
      Note: Hosted realm and browser client configuration that currently does not set token/session lifespans
    - Path: docs/deployments/draft-review-coolify.md
      Note: Hosted runtime contract and session secret expectations for Coolify deployment
    - Path: pkg/auth/config.go
      Note: Current auth settings surface where app-defined session TTL should be added
    - Path: pkg/auth/oidc.go
      Note: OIDC callback currently copies provider token expiry into app session claims
    - Path: pkg/auth/session.go
      Note: Session cookie expiry and validation logic that likely causes short browser sessions
ExternalSources: []
Summary: ""
LastUpdated: 2026-03-26T13:14:42.689599482-04:00
WhatFor: Track the investigation and implementation planning for Draft Review's short-lived hosted author sessions.
WhenToUse: Use when reviewing the auth expiry problem, planning the repair slice, or onboarding someone to the current OIDC session design.
---


# Author session expiry investigation and durable OIDC session design

## Overview

This ticket captures the design investigation for the hosted "logged out too
quickly" issue in Draft Review. The current code path verifies Keycloak login
correctly but appears to copy provider token expiry directly into the app's own
session cookie, which likely makes the browser session expire much sooner than an
authoring tool should.

The main output is the intern-facing design and implementation guide:

- [01-author-session-expiry-investigation-and-durable-oidc-session-design-guide.md](./design-doc/01-author-session-expiry-investigation-and-durable-oidc-session-design-guide.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

This ticket is documentation-only so far. No code changes are included in this
workspace.

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
