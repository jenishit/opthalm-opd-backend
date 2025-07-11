// Package testutil provides an integration-test harness that wires up the
// full application (exactly as cmd/main.go does) against a real Postgres
// database and exposes it via httptest, plus small helpers for the
// login/create-patient/etc. boilerplate every test needs.
//
// It deliberately does NOT use mocked repositories: the bugs this suite is
// meant to catch (pgx type-scanning mismatches, missing NOT NULL columns,
// broken SQL fragments) only show up against a real database.
package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	auth "github.com/jenish-brainztechs/go-backend/internal/adapter/auth/jwt"
	redisadapter "github.com/jenish-brainztechs/go-backend/internal/adapter/cache/redis"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/config"
	httpadapter "github.com/jenish-brainztechs/go-backend/internal/adapter/handler/http"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres"
	"github.com/jenish-brainztechs/go-backend/internal/adapter/storage/postgres/repository"
	"github.com/jenish-brainztechs/go-backend/internal/core/services"
)

// appTables lists every table this test suite may write to, in an order
// safe for TRUNCATE ... CASCADE (CASCADE makes strict ordering unnecessary,
// but listing children-of before parents-of keeps intent readable).
var appTables = []string{
	"lab_job_status_history", "lab_jobs",
	"stock_purchase_items", "stock_purchases",
	"inventory_stock_movements", "inventory_items", "vendors",
	"payments", "invoice_items", "invoices",
	"follow_ups", "refraction_readings", "investigations", "examination_findings", "visit_symptoms",
	"visits",
	"medicines", "diagnosis_catalog", "history_conditions",
	"patients",
	"subscriptions",
	"clinic_settings",
	"verification_tokens",
	"sessions",
	"profile", "users", "role",
}

// TestServer bundles a running httptest server (the full app) with direct
// DB access for setup/teardown and assertions the API doesn't expose.
type TestServer struct {
	*httptest.Server
	DB *postgres.DB
	// ClinicID is the primary test clinic, created once per TestServer.
	// Most tests only ever need this one tenant; cross-tenant tests use
	// NewClinic to spin up an additional, isolated one.
	ClinicID uuid.UUID
	// Emails captures every email AuthService would have sent (password
	// reset / email verification codes), in place of real delivery, so
	// tests can pull the plaintext token straight out of the body.
	Emails *emailCapture
}

// SentEmail is one email captured by emailCapture during a test.
type SentEmail struct {
	To, Subject, Body string
}

// emailCapture is a port.EmailSender that records instead of delivering,
// used as AuthService's email sender for every TestServer.
type emailCapture struct {
	mu     sync.Mutex
	emails []SentEmail
}

func (e *emailCapture) Send(_ context.Context, to, subject, body string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.emails = append(e.emails, SentEmail{To: to, Subject: subject, Body: body})
	return nil
}

// Last returns the most recently captured email, failing the test if none
// has been sent yet.
func (e *emailCapture) Last(t *testing.T) SentEmail {
	t.Helper()
	e.mu.Lock()
	defer e.mu.Unlock()
	require.NotEmpty(t, e.emails, "expected an email to have been sent")
	return e.emails[len(e.emails)-1]
}

// NewTestServer builds the full DI graph identically to cmd/main.go and
// serves it via httptest, against the docker-compose Postgres. It skips the
// test (rather than failing it) if that database isn't reachable, so
// `go test ./...` stays safe to run without Postgres up.
func NewTestServer(t *testing.T) *TestServer {
	t.Helper()

	setTestEnvDefaults()

	cfg, err := config.New()
	require.NoError(t, err, "load config")

	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		t.Skipf("test database not reachable (%s:%s): %v — start it with `docker compose up -d`", cfg.DB.Host, cfg.DB.Port, err)
	}

	redisClient, err := redisadapter.New(ctx, cfg.Redis)
	if err != nil {
		t.Skipf("test redis not reachable (%s, db %d): %v — start it with `docker compose up -d`", cfg.Redis.Addr, cfg.Redis.DB, err)
	}
	// Unlike Postgres (truncated fresh per TestServer via Reset), Redis
	// persists across every NewTestServer call in the same `go test` run.
	// Without this, rate-limit counters (and any stale subscription cache
	// entries) accumulate across unrelated tests — e.g. dozens of tests
	// calling /api/auth/login would eventually trip the login rate limit
	// even though each test only logs in once or twice.
	require.NoError(t, redisClient.FlushDB(ctx).Err(), "flush test redis")

	tokenService, err := auth.New(cfg.Token)
	require.NoError(t, err, "init token service")

	roleRepo := repository.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepo)
	roleHandler := httpadapter.NewRoleHandler(roleService)

	profileRepo := repository.NewProfileRepository(db)
	profileService := services.NewProfileService(profileRepo)
