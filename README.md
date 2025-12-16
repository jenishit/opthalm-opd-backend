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
