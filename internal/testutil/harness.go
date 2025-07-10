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
