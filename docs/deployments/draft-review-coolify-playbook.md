# Draft Review Coolify Deployment Playbook

This is the operator runbook for deploying `draft-review` to Coolify and verifying the hosted runtime.

## Known-Good Target Shape

- Public hostname: `https://draft-review.app.scapegoat.dev`
- Keycloak issuer: `https://auth.scapegoat.dev/realms/draft-review`
- Coolify host: `89.167.52.236`
- Coolify app UUID: `p128frx236jo6e1py4ajslhn`
- Coolify project UUID: `n8xkgqpbjj04m4pishy3su5e`

## Preconditions

Before attempting a deploy, make sure these are true:

1. The branch you want Coolify to build is pushed to GitHub.
2. The hosted Keycloak realm and `draft-review-web` client already exist.
3. You have the browser client secret that matches the Terraform apply.
4. The Coolify app is configured to build from the repository root `Dockerfile`.
5. The Coolify app env includes the OIDC, session, and database values listed in [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md).
6. If the app uses Coolify's public Git source instead of a GitHub App integration, the repository is public.

## Hosted Keycloak Apply

The realm and browser client are managed from:

```text
/home/manuel/code/wesen/terraform/keycloak/apps/draft-review/envs/hosted
```

Run from the Terraform repo root:

```bash
direnv exec /home/manuel/code/wesen/terraform bash -lc '
  export TF_VAR_realm_name=draft-review
  export TF_VAR_public_app_url=https://draft-review.app.scapegoat.dev
  export TF_VAR_web_client_secret="$(cat /home/manuel/.cache/draft-review-web-client-secret)"
  terraform -chdir=keycloak/apps/draft-review/envs/hosted validate
  terraform -chdir=keycloak/apps/draft-review/envs/hosted plan -input=false
  terraform -chdir=keycloak/apps/draft-review/envs/hosted apply -input=false -auto-approve
'
```

Check for:

- realm `draft-review` exists
- browser client `draft-review-web` exists
- callback redirect is `https://draft-review.app.scapegoat.dev/auth/callback`
- post-logout redirect is `https://draft-review.app.scapegoat.dev/auth/logout/callback*`

## Local Pre-Deploy Validation

Run these before touching hosted deployment:

```bash
go test ./cmd/... ./pkg/...
cd frontend && npm run build
cd /home/manuel/code/wesen/2026-03-24--draft-review
go generate ./pkg/web
docker build -t draft-review:local .
docker run --rm -p 18081:8080 draft-review:local
curl -fsS http://127.0.0.1:18081/healthz
curl -fsS http://127.0.0.1:18081/api/info
curl -I http://127.0.0.1:18081/
```

Exit criteria:

- tests pass
- frontend build passes
- Docker image builds
- `/` returns HTML instead of `404`

## Coolify Create Or Redeploy

Preferred path: use the Coolify UI.

If UI access is not available and the API is restricted to localhost, run requests from inside the `coolify` container over SSH.

Deploy trigger pattern:

```bash
ssh root@89.167.52.236 \
  'docker exec coolify sh -lc '\''curl -sS \
    -H "Authorization: Bearer <token>" \
    -H "Accept: application/json" \
    "http://127.0.0.1:8080/api/v1/deploy?uuid=p128frx236jo6e1py4ajslhn"'\'''
```

Required runtime env for the existing app:

```env
DRAFT_REVIEW_DSN=postgres://draft_review:<password>@go1o5tbegalwy3kesshq3hcp:5432/draft_review?sslmode=disable
DRAFT_REVIEW_AUTO_MIGRATE=true
DRAFT_REVIEW_AUTH_MODE=oidc
DRAFT_REVIEW_AUTH_SESSION_SECRET=<long-random-secret>
DRAFT_REVIEW_OIDC_ISSUER_URL=https://auth.scapegoat.dev/realms/draft-review
DRAFT_REVIEW_OIDC_CLIENT_ID=draft-review-web
DRAFT_REVIEW_OIDC_CLIENT_SECRET=<same secret used by Terraform>
DRAFT_REVIEW_OIDC_REDIRECT_URL=https://draft-review.app.scapegoat.dev/auth/callback
```

If the API returns `403`, check Coolify's `instance_settings.allowed_ips`. Earlier investigation showed that valid tokens were rejected when requests originated outside the allowed source IP set.

## Hosted Validation

After the redeploy finishes:

```bash
curl -ksSI https://draft-review.app.scapegoat.dev/
curl -ksS https://draft-review.app.scapegoat.dev/healthz
curl -ksS https://draft-review.app.scapegoat.dev/api/info
curl -ksSI https://draft-review.app.scapegoat.dev/auth/login
```

Check for:

- `/` returns `200 OK`
- `/healthz` returns a healthy JSON payload
- `/api/info` reports `authMode: "oidc"` and `issuerUrl: "https://auth.scapegoat.dev/realms/draft-review"`
- `/auth/login` redirects into the `draft-review` realm
- the login button on `/` does not send the browser to `http://127.0.0.1:8080/auth/login`

Then do one real browser pass:

1. Open `https://draft-review.app.scapegoat.dev/`
2. Start login
3. Complete Keycloak auth
4. Confirm `/api/me` shows an authenticated user
5. Start logout from the app and confirm Keycloak accepts the logout callback instead of showing `Invalid redirect uri`
6. Confirm `/api/me` returns `authenticated: false` after logout
7. Open a reader link such as `/r/<token>` and confirm the SPA shell loads

## Sharp Edges

- Coolify builds from GitHub, not from the local checkout.
- Keycloak redirect URI mismatches fail before the backend sees a callback.
- Keycloak post-logout redirect mismatches fail inside Keycloak before Draft Review sees the request.
- A healthy `/healthz` is not enough; `/` must return the embedded shell.
- The browser client secret in Coolify must exactly match the value used in Terraform.
- A production frontend bundle can still contain localhost-only auth behavior even when `/api/info` and `/healthz` look healthy, so always click the real login button in the browser.
