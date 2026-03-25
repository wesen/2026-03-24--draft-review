# Draft Review Coolify Deployment

This document describes the hosted deployment shape for `draft-review` on Coolify.

## Current Target Shape

- Public app URL: `https://draft-review.app.scapegoat.dev`
- Keycloak issuer: `https://auth.scapegoat.dev/realms/draft-review`
- Keycloak client: `draft-review-web`
- Health check: `https://draft-review.app.scapegoat.dev/healthz`
- Coolify app UUID: `p128frx236jo6e1py4ajslhn`
- Coolify project UUID: `n8xkgqpbjj04m4pishy3su5e`

The app is deployed as a Dockerfile-based Coolify application. The production image now builds the React frontend, embeds it into the Go binary, and serves the SPA shell directly from the backend process. That means Coolify only needs one container and one exposed port.

## Runtime Contract

The container starts with:

```text
draft-review serve
```

Important runtime environment variables:

```env
DRAFT_REVIEW_LISTEN_HOST=0.0.0.0
DRAFT_REVIEW_LISTEN_PORT=8080
DRAFT_REVIEW_AUTO_MIGRATE=true
DRAFT_REVIEW_DSN=postgres://draft_review:<password>@<postgres-host>:5432/draft_review?sslmode=disable
DRAFT_REVIEW_AUTH_MODE=oidc
DRAFT_REVIEW_AUTH_SESSION_SECRET=<long-random-secret>
DRAFT_REVIEW_OIDC_ISSUER_URL=https://auth.scapegoat.dev/realms/draft-review
DRAFT_REVIEW_OIDC_CLIENT_ID=draft-review-web
DRAFT_REVIEW_OIDC_CLIENT_SECRET=<same secret used by Terraform>
DRAFT_REVIEW_OIDC_REDIRECT_URL=https://draft-review.app.scapegoat.dev/auth/callback
```

Current hosted DSN shape:

```env
DRAFT_REVIEW_DSN=postgres://draft_review:<password>@go1o5tbegalwy3kesshq3hcp:5432/draft_review?sslmode=disable
```

## Coolify Application Shape

- Build pack: `Dockerfile`
- Dockerfile path: `Dockerfile`
- Port: `8080`
- Health check path: `/healthz`
- Branch: the branch pushed to GitHub that contains the production frontend embed changes

Coolify should build from the repository root. There is no separate frontend service in production.

## Hosted Keycloak Provisioning

The hosted Keycloak realm and browser client belong in the central Terraform repo:

- `/home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted`

The hosted browser client must include:

- redirect URI `https://draft-review.app.scapegoat.dev/auth/callback`
- post-logout redirect URI `https://draft-review.app.scapegoat.dev/auth/logout/callback*`
- web origin `https://draft-review.app.scapegoat.dev`
- client ID `draft-review-web`

The shared Terraform browser-client module also needs to pass `valid_post_logout_redirect_uris`; otherwise login succeeds but hosted logout fails inside Keycloak with `Invalid redirect uri`.

## Pre-Deploy Validation

Before triggering the first hosted deployment, make sure these are true:

- the `draft-review` realm exists in hosted Keycloak
- the `draft-review-web` client exists in that realm
- the Coolify app env contains the same client secret Terraform applied
- the public hostname already resolves to the Coolify edge
- the branch being deployed has been pushed to GitHub
- if using Coolify's generic public Git source, the GitHub repository must be public or otherwise reachable from Coolify

## Local Container Smoke Check

Build and run the image locally:

```bash
docker build -t draft-review:local .
docker run --rm -p 8080:8080 draft-review:local
```

Then verify:

```bash
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/api/info
curl -I http://127.0.0.1:8080/
```

## Coolify Control Plane Notes

On this host, Coolify API access is effectively localhost-only and filtered through Coolify's API allowlist. The reliable operator patterns are:

- use the Coolify UI when practical
- if API inspection or deploy triggering is needed, SSH to the host and call the API from inside the `coolify` container

Pattern:

```bash
ssh root@89.167.52.236 \
  'docker exec coolify sh -lc '\''curl -sS \
    -H "Authorization: Bearer <token>" \
    -H "Accept: application/json" \
    "http://127.0.0.1:8080/api/v1/deploy?uuid=<app-uuid>"'\'''
```

If a host-side API request returns `403` immediately, inspect `instance_settings.allowed_ips` before assuming the token is invalid.

## Hosted Verification

After deployment, validate in this order:

```bash
curl -ksS https://draft-review.app.scapegoat.dev/healthz
curl -ksS https://draft-review.app.scapegoat.dev/api/info
curl -ksSI https://draft-review.app.scapegoat.dev/
curl -ksSI https://draft-review.app.scapegoat.dev/auth/login
```

Then complete a real browser login and verify:

- `/auth/login` redirects into `/realms/draft-review/`
- `/api/me` returns an authenticated OIDC identity after login
- `/auth/logout?return_to=%2Fapi%2Fme` reaches the Keycloak logout confirmation page and returns to unauthenticated `/api/me`
- `/` returns the embedded frontend shell rather than `404`
- `/r/<token>` loads the reader shell
