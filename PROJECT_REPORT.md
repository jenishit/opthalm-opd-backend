# Project Report: Opthalmic Management System

**Updated:** September 30, 2026 (supersedes the July 7, 2026 report — the project has grown substantially since then: multi-tenancy, full SaaS auth, and five business modules were added)

---

## Overview

A Go backend for a multi-tenant ophthalmology clinic management system, following **Clean Architecture / Hexagonal** pattern with 3 layers: `core` (domain/port/services), `adapter` (handler/repository/config/auth/cache/email), and `cmd` (entrypoint).

**Module:** `github.com/jenishit/opthalm-opd-backend`
**Go Version:** 1.26.4
**Framework:** Gin (v1.12.0)
**Database:** PostgreSQL 16 (pgx/v5) · **Cache:** Redis 7

See `README.md` for the architecture diagram, quick-start, and operational notes — this document tracks what's implemented and what's still missing.

---

## What Has Been Implemented

### 1. Multi-tenancy
Shared database, `clinic_id` column on every tenant-scoped table (patients, visits, billing, inventory, lab jobs). `role` and the reference catalogs (medicines, diagnoses, history conditions) are deliberately global. Cross-tenant isolation is covered by `tenancy_test.go`.

### 2. Authentication & subscription gating
- Short-lived access JWT (15m default) + rotating, revocable opaque refresh tokens (`sessions` table), with reuse-detection (a replayed, already-rotated-out refresh token revokes the whole session chain).
- `POST /api/signup` — creates a clinic + first `ROLE_ADMIN` user + 14-day trial subscription, atomically.
- Subscription gating middleware — 402 on any tenant-scoped route if the clinic's subscription isn't active/trialing. Only `ROLE_SUPERADMIN` (manually provisioned) can manage a clinic's subscription via `/api/platform/subscriptions/:clinicId`.
- Password reset and email verification flows (`verification_tokens` table, single-use, expiring), delivered via a pluggable `EmailSender` (log-only in dev, real SMTP if `SMTP_HOST` is set).

### 3. Business modules
| Module | Endpoints (base path) |
|---|---|
| Patients | `/api/patient` — CRUD, search |
| Visits | `/api/visit` — create/get/update, list by patient |
| Catalog | `/api/admin/catalog/{medicines,diagnoses,conditions}` — global reference data |
| Billing/POS | `/api/billing/invoice` — create, list, search, status, payments, PDF/QR/WhatsApp-link |
| Inventory | `/api/inventory/{items,vendors,purchases}` — CRUD, stock movements, barcode image, low-stock |
| Lab jobs | `/api/lab-jobs` — create, list, status updates |
| Reports | `/api/reports/*` — sales (daily/monthly/range), dues, inventory valuation/low-stock, visits summary; CSV export via `?format=csv` |
| Calculators | `/api/calculators/*` — transposition, spherical equivalent, near-add, vertex distance, telescope FOV (stateless optical math) |
| Clinic settings | `/api/admin/clinic` — per-tenant branding/contact info used on invoices |

### 4. Cross-cutting concerns
- **CORS** — configurable allowed origins.
- **Structured JSON responses** — `{success, message, data}` / `{success, messages}`, with a sentinel-error → HTTP-status map covering every domain error.
- **Structured logging** — one `log/slog` line per request (request ID, latency, authenticated user/clinic), JSON in production.
- **Rate limiting** — Redis-backed, on signup/login/password-reset.
- **Graceful shutdown** — drains in-flight requests on `SIGINT`/`SIGTERM`.
- **Swagger/OpenAPI** — generated from handler annotations, served at `/swagger/index.html` outside production.
- **Panic recovery** — `gin.CustomRecovery` returns the standard JSON error envelope instead of crashing the process.

### 5. Testing
Real-database integration suite (no mocked repositories) in `internal/adapter/handler/http/*_test.go`, built on a shared harness (`internal/testutil/harness.go`) that wires the exact same DI graph as `cmd/main.go`. Covers every module above plus auth/session lifecycle, tenant isolation, subscription gating, and rate limiting.

### 6. Infrastructure
- Multi-stage `Dockerfile` (`golang:1.26-alpine` → `alpine:3.20`), `docker-compose.yml` with healthchecked Postgres + Redis and `depends_on: condition: service_healthy` ordering.
- 14 goose-style SQL migrations (`internal/adapter/storage/postgres/migrations/`), applied via the `goose` CLI (not run automatically by the app).

---

## What Is Missing / Deferred

| Item | Status |
|---|---|
| AI/OCR for prescription scanning | Deliberately deferred (explicit early scoping decision) |
| Desktop/hardware integrations (barcode scanners, card readers) | Out of scope — backend-only by design |
| Frontend | Not part of this codebase |
| CI pipeline | Not set up — tests are run manually against docker-compose |
| Payment gateway integration | Subscription billing is an internal/manual flag (`subscriptions` table), not wired to a real payment processor |
| Audit log | No append-only record of who changed what; sentinel errors + structured logs are the only trail today |

---

## Known Rough Edges

| Item | Details |
|---|---|
| Dead config structs | `config.Session`, `config.Cache`, and the legacy `config.Redis` struct are loaded but unused — superseded by `config.RedisConfig`. Left in place rather than risk breaking something that reads them; safe to remove in a follow-up pass. |
| Lambda leftovers | `config.go` still carries `IsLambdaRuntime()` and AWS secret-ARN resolution logic from an earlier serverless-adapted template. Unused in the current docker-compose deployment model. |
| `DB.Migrate()` unused | A golang-migrate-based `Migrate()` method exists on `postgres.DB` but nothing calls it — migrations are applied via the `goose` CLI instead, which is what all 14 existing migration files are formatted for (`-- +goose Up/Down` markers, not golang-migrate's `.up.sql`/`.down.sql` pairs). Worth removing the dead method or wiring it up properly, not leaving both half-present. |
