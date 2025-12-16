# Project Report: Opthalmic Management System

**Updated:** September 30, 2026 (supersedes the July 7, 2026 report — the project has grown substantially since then: multi-tenancy, full SaaS auth, and five business modules were added)

---

## Overview

A Go backend for a multi-tenant ophthalmology clinic management system, following **Clean Architecture / Hexagonal** pattern with 3 layers: `core` (domain/port/services), `adapter` (handler/repository/config/auth/cache/email), and `cmd` (entrypoint).

**Module:** `github.com/jenish-brainztechs/go-backend`
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
