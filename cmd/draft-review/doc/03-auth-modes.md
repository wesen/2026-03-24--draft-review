---
Title: "Auth Modes"
Slug: "auth-modes"
Short: "Choose between dev auth and Keycloak-backed OIDC for the backend."
Topics:
- auth
- oidc
- keycloak
- backend
Commands:
- serve
Flags:
- auth-mode
- oidc-issuer-url
- oidc-client-id
- oidc-client-secret
- oidc-redirect-url
- auth-session-secret
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This page covers the two backend auth modes, how each one changes runtime behavior, and why both modes exist in the repository instead of forcing every developer task through Keycloak.

`auth-mode=dev` is the fast path. It synthesizes an authenticated local author session so article and analytics routes can be exercised without booting Keycloak. This mode is appropriate when the task is database behavior, route shaping, frontend integration, or smoke testing that does not depend on a real identity provider.

`auth-mode=oidc` is the real browser-auth path. In this mode the backend redirects the browser to Keycloak, validates the callback, signs its own session cookie, and resolves the authenticated claims into a local Draft Review `users` row. This mode is required when validating `/auth/login`, `/auth/callback`, `/auth/logout`, author ownership binding, or any issue that depends on real OIDC claims.

The main required OIDC fields are:

- `--oidc-issuer-url`: the Keycloak realm issuer URL.
- `--oidc-client-id`: the configured OIDC client.
- `--oidc-client-secret`: the shared secret for confidential client auth.
- `--oidc-redirect-url`: the browser callback URL served by the backend.
- `--auth-session-secret`: the secret used to sign the backend session cookie.

The repository ships a local Keycloak realm import and helper `make` targets so the development OIDC loop does not have to be assembled by hand. Use those defaults unless you are explicitly testing non-default ports.

Problem | Cause | Solution
--- | --- | ---
`/api/me` shows unauthenticated in OIDC mode after login | Redirect URL, client, or issuer settings do not match the running Keycloak realm | Recheck the `make local-keycloak-up` defaults and align the OIDC flags
Author writes succeed in dev mode but fail in OIDC mode | The backend did not map the OIDC user into a local author row | Check `/api/me` first, then verify the auth service can upsert the local user
Logout appears to work locally but a fresh browser load is still authenticated | The backend cookie or Keycloak session is still active | Use `/auth/logout` on the backend origin and complete the provider logout flow

See Also:

- `draft-review help local-development`
- `draft-review help database-workflow`
- `draft-review help draft-review-overview`
