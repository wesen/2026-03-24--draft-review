# Changelog

## 2026-03-26

- Initial workspace created
- Added a detailed design and implementation guide explaining why Draft Review's
  hosted browser session likely expires too quickly and how to decouple the app
  session lifetime from the raw Keycloak token lifetime
- Split the auth roadmap into two tickets: DR-010 now owns the medium-term opaque
  server-side session implementation, while DR-011 captures the full long-term
  refresh-token renewal architecture
- Added a DR-010 implementation diary and expanded DR-010 into concrete
  implementation tasks covering repository work, session manager replacement,
  callback changes, request auth resolution, logout, tests, and validation
