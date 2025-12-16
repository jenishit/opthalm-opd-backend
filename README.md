# Opthalmic Management System — Backend

A multi-tenant (SaaS) backend for an ophthalmology clinic management system: patient records, visits, billing/POS, inventory, lab jobs, reports, and optical calculators. Go, hexagonal architecture, Postgres, Redis.

## Quick start

```bash
cp .env.example .env   # then fill in TOKEN_SECRET, POSTGRES_PASSWORD, REDIS_PASSWORD
docker compose up -d --build
```

This brings up Postgres → Redis → the app, in that dependency order (`depends_on: condition: service_healthy`). The app listens on `HTTP_PORT` (default `8082`).

Apply database migrations with [goose](https://github.com/pressly/goose) (not run automatically):

```bash
goose -dir internal/adapter/storage/postgres/migrations postgres \
  "postgres://$DB_USER:$DB_PASSWORD@localhost:5434/$DB_NAME?sslmode=disable" up
```

API docs (Swagger UI) are served at `/swagger/index.html` in any non-`production` environment — useful for exploring the API and for generating a frontend client. Regenerate after changing handler annotations:

```bash
swag init -g cmd/main.go -o docs --parseInternal
```

## Running tests

The integration suite (`internal/adapter/handler/http/*_test.go`) runs against a real Postgres + Redis — it does not mock the database, since the bugs it's meant to catch (SQL/type mismatches, missing constraints, broken tenant isolation) only show up against the real thing.

```bash
docker compose up -d postgres redis
go test ./...
```

Tests skip (not fail) if Postgres/Redis aren't reachable. `internal/testutil/harness.go` builds the exact same dependency graph as `cmd/main.go` and serves it via `httptest`, so what's tested is the real wiring, not a stand-in.

## Architecture

Hexagonal / ports-and-adapters, three layers:

```
cmd/main.go                    composition root — wires every dependency, starts the server

internal/core/                 business logic — zero knowledge of HTTP, SQL, or JWT
  domain/                      entities, sentinel errors, value objects (e.g. Password)
  port/                        interfaces only — *Repository (outbound), *Service (inbound)
  services/                    use cases; depend only on port interfaces

internal/adapter/              infrastructure
  handler/http/                Gin handlers + DTOs + router + middleware
  storage/postgres/            pgx repositories implementing the *Repository ports
  auth/jwt/                    JWT access tokens + opaque refresh tokens
  cache/redis/                 Redis client
  email/{logsender,smtpsender}/  EmailSender implementations (logsender is the dev default)
```

Every repository method that touches a tenant-scoped table takes a `clinicID` and filters on it — `patient.go`'s repository is the template for the pattern used throughout. `role` and the reference catalogs (`medicines`, `diagnosis_catalog`, `history_conditions`) are deliberately global, not tenant-scoped.

## Auth model

- **Access token**: short-lived JWT (`TOKEN_DURATION`, default 15m), stateless — not checked against the database on every request.
- **Refresh token**: opaque random token, SHA-256-hashed at rest in `sessions`, rotated on every use. Presenting an already-rotated-out token revokes the entire session chain for that user (theft/reuse detection).
- **Signup** (`POST /api/signup`): creates a clinic, its first `ROLE_ADMIN` user, and a 14-day trial subscription, atomically.
- **Subscription gating**: every tenant-scoped route requires the caller's clinic to have an active/trialing subscription (402 otherwise). Only `ROLE_SUPERADMIN` (a manually-created platform-operator account) can manage a clinic's subscription — a clinic can never reactivate itself.
- **Password reset / email verification**: `POST /api/auth/password-reset/{request,confirm}` and `/api/auth/email/verify/{resend,confirm}`. Emails are sent via a pluggable `EmailSender` — logged instead of delivered (`internal/adapter/email/logsender`) unless `SMTP_HOST` is set, in which case real SMTP delivery is used instead.

## Operational notes

- **Graceful shutdown**: the server drains in-flight requests (up to 15s) on `SIGINT`/`SIGTERM` before exiting — see `Router.Serve`.
- **Structured logging**: every request logs one `log/slog` line (request ID, method, path, status, latency, and the authenticated user/clinic when present). JSON in `APP_ENV=production`, human-readable text otherwise. Set `LOG_LEVEL` to override (default `info`).
- **Rate limiting**: Redis-backed, applied to `/api/signup`, `/api/auth/login`, and the password-reset endpoints — the endpoints most worth protecting from brute-force/abuse. Fails open if Redis is unreachable.

## Environment variables

See `.env.example` for the full list. Notable ones beyond the obvious DB/HTTP settings:

| Variable | Purpose |
|---|---|
| `TOKEN_SECRET`, `TOKEN_DURATION` | Access token signing key + lifetime |
| `REFRESH_TOKEN_DURATION` | Refresh token lifetime (default `168h`) |
| `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB` | Subscription cache + rate limiter |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | Real email delivery; omit `SMTP_HOST` to use the log-only dev sender |
| `LOG_LEVEL` | `debug` / `info` (default) / `warn` / `error` |

Inside the docker-compose network, services address each other by service name and internal port (`postgres:5432`, `redis:6379`) — not the host-side published ports, which only matter for connecting from outside the compose network (e.g. running `go test` on the host, or `psql` from your terminal).
