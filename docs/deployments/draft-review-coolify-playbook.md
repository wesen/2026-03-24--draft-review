# Draft Review Coolify Deployment Playbook

This is the operator runbook for deploying `draft-review` to Coolify and verifying the hosted runtime.

## Known-Good Target Shape

- Public hostname: `https://draft-review.app.scapegoat.dev`
- Keycloak issuer: `https://auth.scapegoat.dev/realms/draft-review`
- Coolify host: `89.167.52.236`

## Preconditions

Before attempting a deploy, make sure these are true:

1. The branch you want Coolify to build is pushed to GitHub.
2. The hosted Keycloak realm and `draft-review-web` client already exist.
3. You have the browser client secret that matches the Terraform apply.
4. The Coolify app is configured to build from the repository root `Dockerfile`.
5. The Coolify app env includes the OIDC, session, and database values listed in [draft-review-coolify.md](/home/manuel/code/wesen/2026-03-24--draft-review/docs/deployments/draft-review-coolify.md).

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

## Coolify Redeploy

Use the Coolify UI to redeploy the existing application. The established operator guidance in the other repos is still the right default: prefer the UI over the API because the API token path has historically been unreliable for project and application operations.

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

Then do one real browser pass:

1. Open `https://draft-review.app.scapegoat.dev/`
2. Start login
3. Complete Keycloak auth
4. Confirm `/api/me` shows an authenticated user
5. Open a reader link such as `/r/<token>` and confirm the SPA shell loads

## Sharp Edges

- Coolify builds from GitHub, not from the local checkout.
- Keycloak redirect URI mismatches fail before the backend sees a callback.
- A healthy `/healthz` is not enough; `/` must return the embedded shell.
- The browser client secret in Coolify must exactly match the value used in Terraform.
