# Project Report: Reports-Management-System

**Generated:** July 7, 2026

---

## Overview

A Go backend for a lab reports management system following **Clean Architecture / Hexagonal** pattern with 3 layers: `core` (domain/port/services), `adapter` (handler/repository/config/auth), and `cmd` (entrypoint).

**Module:** `github.com/jenish-brainztechs/go-backend`
**Go Version:** 1.26.4
**Framework:** Gin (v1.12.0)
**Database:** PostgreSQL 16 (pgx/v5)

---

## Architecture: Clean Architecture / Hexagonal Pattern

The project follows **Clean Architecture** (also known as Hexagonal Architecture or Ports & Adapters), which enforces a strict dependency rule: **dependencies point inward** — outer layers depend on inner layers, never the reverse.

### Layer Structure

```
cmd/                        ← Entrypoint / Composition Root
  └── main.go               ← Wires all dependencies together

internal/
  └── adapter/              ← Outer layer (infrastructure / IO)
      ├── auth/jwt/         ←   JWT adapter (implements TokenService port)
      ├── config/           ←   Env config loader
      ├── handler/http/     ←   HTTP handlers (Gin controllers)
      │   ├── dto/          ←   Request/Response DTOs
      │   ├── middleware.go ←   CORS + Auth middleware
      │   ├── response.go   ←   Standardized response helpers
      │   └── router.go     ←   Gin router + route definitions
      └── storage/postgres/ ←   PostgreSQL adapter
          ├── db.go         ←   Connection pool + migration runner
          ├── migrations/   ←   SQL migration files
          └── repository/   ←   Repository implementations

  └── core/                 ← Inner layer (business logic, zero external deps)
      ├── domain/           ←   Enterprise business rules (models, errors, VOs)
      │   └── valueobjects/ ←   Value objects (e.g., Password with bcrypt)
      ├── port/             ←   Interface definitions (ports)
      │   ├── *Repository   ←   Outbound ports (driven interfaces)
      │   └── *Service      ←   Inbound ports (driving interfaces)
      └── services/         ←   Application business rules (use cases)
```

### How It Works

1. **Domain Layer** (`core/domain/`) — Pure Go types, no imports from other packages. Contains entity structs (`User`, `Role`, `Profile`, `Patient`, `LabSettings`), value objects (`Password` with bcrypt), sentinel errors, and business rules.

2. **Port Layer** (`core/port/`) — Defines **interfaces only**. `UserRepository` / `ProfileRepository` / etc. are outbound ports (what the app needs from the outside world). `UserService` / `AuthService` / etc. are inbound ports (what the app offers to the outside world).

3. **Service Layer** (`core/services/`) — Implements **application use cases**. Depends only on port interfaces, never on concrete adapters. Dependencies are injected via constructors (e.g., `NewAuthService(repo, tokenService)`).

4. **Adapter Layer** (`internal/adapter/`) — Implements **port interfaces** using concrete technologies:
   - `handler/http/` — Gin-based HTTP handlers that call service ports
   - `auth/jwt/` — JWT implementation of `TokenService` port
   - `storage/postgres/repository/` — pgx-based implementations of repository ports

5. **Composition Root** (`cmd/main.go`) — Creates all concrete implementations and wires them together through constructor injection. No service locator, no global state.

### Data Flow Example (Login)

```
HTTP Request
  → Router (Gin)
    → AuthHandler.Login()
      → domain.Login (request DTO)
        → AuthService.Login()  [core/services]
          → UserRepository.GetUserByEmail()  [port interface]
            → UserRepository (Postgres)  [adapter/storage]
          ← domain.BasicDetails
          → Password.Verify(plaintext)  [domain/valueobjects]
          → JWTToken.CreateAccessToken()  [port interface]
            → JWTToken (adapter/auth/jwt)
          ← token string
        ← domain.LoginResponse
      ← dto.LoginResponse (response DTO)
    ← handleSuccess(ctx, response)
  ← JSON Response
```

### Key Benefits

- **Testability**: Services can be tested with mock repositories/services
- **Swapability**: Swap Postgres for MongoDB by writing a new adapter
- **Isolation**: Business logic has zero knowledge of HTTP, JWT, or SQL
- **Single Responsibility**: Each layer has a clear, non-overlapping concern

---

## Project Structure

```
go-backend/
├── .air.toml                        # Air live-reload config
├── .env                             # Local environment variables
├── .env.example                     # Example env vars template
├── docker-compose.yml               # Postgres 16 service (port 5433)
├── go.mod / go.sum
├── cmd/
│   └── main.go                      # Entrypoint — DI wiring + server start
└── internal/
    ├── adapter/
    │   ├── auth/jwt/jwt.go          # JWT token creation & verification (HS256)
    │   ├── config/config.go         # Env-based config loader (godotenv)
    │   ├── handler/http/
    │   │   ├── auth.go              # Login handler
    │   │   ├── middleware.go        # CORS + Bearer auth middleware
    │   │   ├── profile.go           # Profile CRUD handlers
    │   │   ├── response.go          # Standardized JSON responses + error mapping
    │   │   ├── role.go              # Role creation handler
    │   │   ├── router.go            # Gin router + route definitions
    │   │   ├── users.go             # User creation handler
    │   │   └── dto/                 # Request/Response DTOs
    │   └── storage/postgres/
    │       ├── db.go                # Connection pool + migration runner
    │       ├── migrations/          # 3 SQL migration files
    │       └── repository/          # CRUD implementations
    └── core/
        ├── domain/                  # Domain models, errors, value objects
        ├── port/                    # Repository & Service interfaces
        └── services/                # Business logic implementations
```

**Total: 46 files** (36 Go source, 3 SQL migrations, 7 config/support)

---

## What Has Been Implemented

### 1. Database & Storage
- **3 migrations** — `role`, `users`/`profile`, and `lab_settings` tables
- **PostgreSQL connection pool** via pgx/v5
- **SQL query builder** (Masterminds/squirrel) for parameterized, injection-safe queries
- **Repositories:** Role, User, Profile, Settings — all implement port interfaces

### 2. Authentication
- **JWT token service** — HS256, configurable duration, custom claims (user_id, role_name, session_id)
- **Login flow** — email lookup → bcrypt verify → JWT generation → response
- **Bearer token middleware** — extracts, verifies, and injects payload into Gin context
- **Password value object** — bcrypt hashing (cost 14), validation, marshaling-safe

### 3. API Endpoints
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/login` | No | User login |
| POST | `/api/role/create` | No | Create new role |
| POST | `/api/user/create` | No | Register user |
| GET | `/api/profile/getme` | Yes | Get own profile |
| GET | `/api/profile/profile-details` | Yes | Get all profiles |
| PATCH | `/api/profile/update-profile/:id` | Yes | Update profile by user ID |

### 4. Business Logic (Services)
- **AuthService** — login with credential verification
- **UserService** — user creation with role lookup & profile creation in one transaction-like flow
- **RoleService** — role creation, name ↔ ID lookups
- **ProfileService** — CRUD for profiles with JOIN queries
- **SettingsService** — upsert lab settings (incomplete — see below)

### 5. Cross-Cutting Concerns
- **CORS middleware** — allows configurable origins
- **Standardized JSON responses** — `{success, message, data}` / `{success, messages}`
- **Sentinel error mapping** — 30+ domain errors mapped to appropriate HTTP status codes
- **Air live-reload** configured for development

---

## What Is Missing / Incomplete

### 🔴 Not Implemented (defined but unused)

| Feature | Evidence |
|---------|----------|
| **Settings handler & routes** | `SettingsRepository` & `SettingsService` exist but are never instantiated in `main.go`. No settings endpoints in `router.go`. |
| **Settings service GetSettings** | `port/settings.go` defines `GetSettings` but `services/settings.go` only implements `UpsertSettings`. |
| **Patient module** | `domain/patients.go` defines a `Patient` struct. No repository, service, handler, routes, or migration exist. Completely unused. |
| **Session / Redis** | `Session`, `Redis` config structs defined and loaded. `SessionState`, `RefreshTokenPayload` domain types exist. No session repository, Redis connection, or session middleware implemented. |
| **Refresh tokens** | `RefreshTokenPayload` domain type + `Refresh.Duration` config exist. No refresh token creation, storage, verification, or endpoint. |
| **Logout** | No token invalidation / blacklist mechanism. |

### 🟡 Partially Implemented / Issues

| Issue | Details |
|-------|---------|
| **Migrations never called** | `DB.Migrate()` exists in `db.go` but is **not invoked** in `main.go`. Must be run manually. |
| **Docker port mismatch** | `docker-compose.yml` maps port 5433:5433, but `.env` uses port 5432. |
| **Profile update authorization** | `PATCH /api/profile/update-profile/:id` accepts any `:id` param but has no role-based authorization check. |
| **~Validation error handling unused~** | **FIXED** — All handlers now use `validationError()`, `handleError()`, and `handleSuccess()` helpers instead of raw `ctx.JSON`. |
| **BasicDetails unused fields** | `Username` and `FullName` fields in `domain.BasicDetails` are never populated by the repository query. |
| **Profile email duality** | `profile` table has its own `email` column, but the profile repository JOINs with `users.email`. Potential confusion. |
| **Error message typo** | Several handlers return `"Something went wrong internally while creating error"` (should be "creating the resource"). |
| **Zero tests** | No unit or integration tests exist anywhere in the project. |
| **Lambda leftovers** | Config contains `IsLambdaRuntime()` and secret ARN resolution logic, likely from a cloned/serverless-adapted project. |
| **Project name inconsistency** | `.env.example` references "inventory-management" but app name is "Reports-Management-System". |

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/gin-gonic/gin` | v1.12.0 | HTTP router/framework |
| `github.com/jackc/pgx/v5` | v5.10.0 | PostgreSQL driver |
| `github.com/Masterminds/squirrel` | v1.5.4 | SQL query builder |
| `github.com/golang-jwt/jwt/v4` | v4.5.2 | JWT tokens |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | DB migrations |
| `github.com/joho/godotenv` | v1.5.1 | .env loader |
| `golang.org/x/crypto` | v0.48.0 | bcrypt |
| `github.com/google/uuid` | v1.6.0 | UUID generation |

---

## Recommendations

1. **Call `DB.Migrate()` in `main.go`** before starting the server.
2. **Wire up Settings module** — instantiate repo/service/handler and add routes.
3. **Add tests** — at minimum unit tests for services, integration tests for repositories.
4. **Fix Docker port** to match `.env` (5432) or update `.env` to use 5433.
5. **Decide on Patients** — either implement the full module (repo, service, handler, migration) or remove the domain type.
6. **Add authorization middleware** for role-based access control.
7. **Fix error messages** and enable proper validation error handling.
8. **Clean up legacy code** — remove Lambda/secret ARN logic if not needed.
