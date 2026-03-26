# Changelog

## 2026-03-26

- Initial workspace created

## 2026-03-26 — Implementation

Wired both pre-built Redux slices across 2 commits:

- **2958b2a**: WI-1+2 — Aligned uiSlice View type (5→6 values), added focusSectionId + previewArticle + goBack. Migrated AuthorApp: 4 useState → Redux (22 call sites). Preview reads from previewArticle to preserve unsaved drafts.
- **3f15311**: WI-3+4 — Added sessionId to readerSlice. Migrated ReaderPage: 4 useState → Redux (8 call sites) + unmount cleanup via resetReader().
- WI-5 (MenuBar openMenu) skipped — transient dropdown state, no benefit from Redux.

