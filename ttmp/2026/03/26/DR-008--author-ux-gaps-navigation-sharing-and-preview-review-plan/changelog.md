# Changelog

## 2026-03-26

- Initial workspace created
- Added a detailed design and implementation guide covering deletion, routing, sharing, reader identity, and preview review behavior
- Added a rolling code-review backlog for broader cleanup opportunities discovered during investigation
- Added an investigation diary and explicit task list for the ticket
- Validated the ticket with `docmgr doctor --ticket DR-008 --stale-after 30`
- Uploaded the DR-008 bundle to reMarkable at `/ai/2026/03/26/DR-008`
- Implemented Slice 1: backend article deletion route, service, repository support, focused HTTP coverage, and safer frontend delete mutation handling
- Implemented Slice 2: route-driven author navigation for dashboard, review, edit, settings, preview, and share flows, plus a small Storybook typing fix needed to restore frontend build verification
