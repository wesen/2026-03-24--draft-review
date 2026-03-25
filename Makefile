.PHONY: local-keycloak-up local-keycloak-down local-keycloak-config run-local-dev run-local-oidc

APP_PORT ?= 8080
PG_PORT ?= 15432
KEYCLOAK_PORT ?= 18080
SESSION_SECRET ?= local-session-secret
FRONTEND_DEV_PROXY_URL ?=

local-keycloak-up:
	DRAFT_REVIEW_PG_PORT=$(PG_PORT) \
	DRAFT_REVIEW_KEYCLOAK_PORT=$(KEYCLOAK_PORT) \
	docker compose -f docker-compose.local.yml up -d

local-keycloak-down:
	DRAFT_REVIEW_PG_PORT=$(PG_PORT) \
	DRAFT_REVIEW_KEYCLOAK_PORT=$(KEYCLOAK_PORT) \
	docker compose -f docker-compose.local.yml down

local-keycloak-config:
	DRAFT_REVIEW_PG_PORT=$(PG_PORT) \
	DRAFT_REVIEW_KEYCLOAK_PORT=$(KEYCLOAK_PORT) \
	docker compose -f docker-compose.local.yml config

run-local-dev:
	DRAFT_REVIEW_FRONTEND_DEV_PROXY_URL=$(FRONTEND_DEV_PROXY_URL) \
	GOWORK=off go run ./cmd/draft-review serve \
		--auth-mode dev \
		--dsn 'postgres://draft_review:draft_review@127.0.0.1:$(PG_PORT)/draft_review?sslmode=disable' \
		--listen-host 0.0.0.0 \
		--listen-port $(APP_PORT)

run-local-oidc:
	DRAFT_REVIEW_FRONTEND_DEV_PROXY_URL=$(FRONTEND_DEV_PROXY_URL) \
	DRAFT_REVIEW_AUTH_MODE=oidc \
	DRAFT_REVIEW_AUTH_SESSION_SECRET=$(SESSION_SECRET) \
	DRAFT_REVIEW_OIDC_ISSUER_URL=http://127.0.0.1:$(KEYCLOAK_PORT)/realms/draft-review-dev \
	DRAFT_REVIEW_OIDC_CLIENT_ID=draft-review-web \
	DRAFT_REVIEW_OIDC_CLIENT_SECRET=draft-review-web-secret \
	DRAFT_REVIEW_OIDC_REDIRECT_URL=http://127.0.0.1:$(APP_PORT)/auth/callback \
	GOWORK=off go run ./cmd/draft-review serve \
		--dsn 'postgres://draft_review:draft_review@127.0.0.1:$(PG_PORT)/draft_review?sslmode=disable' \
		--listen-host 0.0.0.0 \
		--listen-port $(APP_PORT)
