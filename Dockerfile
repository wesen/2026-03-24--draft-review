FROM node:22-bookworm AS frontend-builder

WORKDIR /src/frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

FROM golang:1.25.8 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /src/frontend/dist ./frontend/dist

RUN DRAFT_REVIEW_SKIP_FRONTEND_BUILD=1 GOWORK=off go generate ./pkg/web && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/draft-review ./cmd/draft-review

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl tzdata && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/draft-review /usr/local/bin/draft-review

EXPOSE 8080

ENTRYPOINT ["draft-review"]
CMD ["serve"]
