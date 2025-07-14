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
	profileHandler := httpadapter.NewProfileHandler(profileService)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	signupRepo := repository.NewSignupRepository(db)
	verificationRepo := repository.NewVerificationTokenRepository(db)
	emails := &emailCapture{}
	refreshDuration, err := time.ParseDuration(cfg.Refresh.Duration)
	require.NoError(t, err, "parse refresh duration")
	authService := services.NewAuthService(userRepo, sessionRepo, signupRepo, verificationRepo, emails, tokenService, refreshDuration)
	authHandler := httpadapter.NewAuthHandler(authService)
	userService := services.NewUserService(userRepo, roleService, profileService)
	userHandler := httpadapter.NewUsersHandler(userService)

	clinicRepo := repository.NewClinicRepository(db)
	clinicService := services.NewClinicService(clinicRepo)
	clinicHandler := httpadapter.NewClinicHandler(clinicService)

	patientRepo := repository.NewPatientRepository(db)
	patientService := services.NewPatientService(patientRepo)
	patientHandler := httpadapter.NewPatientHandler(patientService)

	visitRepo := repository.NewVisitsRepository(db)
	visitService := services.NewVisitsService(visitRepo)
	visitHandler := httpadapter.NewVisitHandler(visitService)

	medicineRepo := repository.NewMedicineRepository(db)
	medicineService := services.NewMedicineService(medicineRepo)

	diagnosisRepo := repository.NewDiagnosisCatalogRepository(db)
	diagnosisService := services.NewDiagnosisCatalogService(diagnosisRepo)

	conditionRepo := repository.NewHistoryConditionRepository(db)
	conditionService := services.NewHistoryConditionService(conditionRepo)

	catalogHandler := httpadapter.NewCatalogHandler(medicineService, diagnosisService, conditionService)

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := services.NewInventoryService(inventoryRepo)
	inventoryHandler := httpadapter.NewInventoryHandler(inventoryService)

	vendorRepo := repository.NewVendorRepository(db)
	vendorService := services.NewVendorService(vendorRepo)
	vendorHandler := httpadapter.NewVendorHandler(vendorService)

	stockPurchaseRepo := repository.NewStockPurchaseRepository(db, inventoryRepo)
	stockPurchaseService := services.NewStockPurchaseService(stockPurchaseRepo)
	stockPurchaseHandler := httpadapter.NewStockPurchaseHandler(stockPurchaseService)

	invoiceRepo := repository.NewInvoiceRepository(db, inventoryRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)
	invoiceHandler := httpadapter.NewInvoiceHandler(invoiceService, clinicService)

	labJobRepo := repository.NewLabJobRepository(db)
	labJobService := services.NewLabJobService(labJobRepo)
	labJobHandler := httpadapter.NewLabJobHandler(labJobService)

	reportsRepo := repository.NewReportsRepository(db)
	reportsService := services.NewReportsService(reportsRepo, inventoryService)
	reportsHandler := httpadapter.NewReportsHandler(reportsService)

	calculatorService := services.NewCalculatorService()
	calculatorHandler := httpadapter.NewCalculatorHandler(calculatorService)

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, redisClient)
	subscriptionHandler := httpadapter.NewSubscriptionHandler(subscriptionService)

	router, err := httpadapter.NewRouter(
		cfg, tokenService,
		*roleHandler, *userHandler, *profileHandler, *authHandler, *clinicHandler,
		*patientHandler, *visitHandler, *catalogHandler, *invoiceHandler,
		*inventoryHandler, *vendorHandler, *stockPurchaseHandler, *labJobHandler,
		*reportsHandler, *calculatorHandler,
		*subscriptionHandler, subscriptionService,
		redisClient,
	)
	require.NoError(t, err, "build router")

	srv := httptest.NewServer(router)
	ts := &TestServer{Server: srv, DB: db, Emails: emails}

	ts.Reset(t)
	ts.ClinicID = ts.NewClinic(t, "Primary Test Clinic")
	t.Cleanup(func() {
		ts.Reset(t)
		srv.Close()
		db.Close()
		redisClient.Close()
	})

	return ts
}

// Reset truncates every app table so each top-level test starts from a
// clean slate, and restarts the invoice-number sequence so INV-000001 is
// predictable across test runs.
func (ts *TestServer) Reset(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	_, err := ts.DB.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(appTables, ", ")))
	require.NoError(t, err, "truncate app tables")

	_, err = ts.DB.Exec(ctx, "ALTER SEQUENCE IF EXISTS invoice_no_seq RESTART WITH 1")
	require.NoError(t, err, "reset invoice_no_seq")
}

// --- Request helpers -------------------------------------------------

// Do performs a request against the test server, optionally authenticated,
// and JSON-decodes the response body into out (skipped if out is nil).
func (ts *TestServer) Do(t *testing.T, method, path, token string, body any, out any) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, ts.Server.URL+path, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	if out != nil {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		if len(b) > 0 {
			require.NoError(t, json.Unmarshal(b, out), "decode response body: %s", string(b))
		}
	}

	return resp
}

type envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
}

// DoData is like Do but unwraps the {success, message, data} envelope and
// decodes just the `data` field into out.
func (ts *TestServer) DoData(t *testing.T, method, path, token string, body any, out any) *http.Response {
	t.Helper()

	var env envelope
	resp := ts.Do(t, method, path, token, body, &env)
	if out != nil && len(env.Data) > 0 {
		require.NoError(t, json.Unmarshal(env.Data, out), "decode data field")
	}
	return resp
}

// --- Fixture helpers ---------------------------------------------------

const TestPassword = "Password123!"

// EnsureRole creates the role if it doesn't already exist and returns its
// ID. It goes straight to the DB (rather than the public /role/create
// endpoint, which has no upsert semantics and would error on a second call
// with the same name) so it's safe to call repeatedly for the same role
// within one test.
func (ts *TestServer) EnsureRole(t *testing.T, roleName string) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	err := ts.DB.QueryRow(context.Background(),
		`INSERT INTO role (role_name) VALUES ($1)
		 ON CONFLICT (role_name) DO UPDATE SET role_name = EXCLUDED.role_name
		 RETURNING id`,
		roleName,
	).Scan(&id)
	require.NoError(t, err, "ensure role %q", roleName)
	return id
}

// NewClinic creates an additional, isolated tenant (clinic_settings row +
// trialing subscription), for cross-tenant isolation tests. The primary
// clinic (ts.ClinicID) is already created by NewTestServer; most tests never
// need to call this directly.
func (ts *TestServer) NewClinic(t *testing.T, name string) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var clinicID uuid.UUID
	err := ts.DB.QueryRow(ctx,
		`INSERT INTO clinic_settings (clinic_name, registration_no) VALUES ($1, $2) RETURNING id`,
		name, "REG-"+uuid.New().String()[:8],
	).Scan(&clinicID)
	require.NoError(t, err, "create clinic")

	_, err = ts.DB.Exec(ctx,
		`INSERT INTO subscriptions (clinic_id, status, current_period_end) VALUES ($1, 'trialing', now() + interval '14 days')`,
		clinicID,
	)
	require.NoError(t, err, "create subscription")

	return clinicID
}

// SetSubscriptionStatus overwrites a clinic's subscription status/expiry
// directly via SQL, for tests that need to simulate a cancelled or expired
// subscription without going through the platform-operator API.
func (ts *TestServer) SetSubscriptionStatus(t *testing.T, clinicID uuid.UUID, status string, periodEnd time.Time) {
	t.Helper()
	_, err := ts.DB.Exec(context.Background(),
		`UPDATE subscriptions SET status = $1, current_period_end = $2, updated_at = NOW() WHERE clinic_id = $3`,
		status, periodEnd, clinicID,
	)
	require.NoError(t, err, "set subscription status")
}

// SuperadminToken creates a ROLE_SUPERADMIN user directly (there is no
// self-serve way to become one — it's the platform operator's own account)
// scoped to the primary test clinic (superadmin routes don't use ClinicID,
// but a user row still needs one to satisfy the NOT NULL constraint).
func (ts *TestServer) SuperadminToken(t *testing.T) (token string, userID uuid.UUID) {
	t.Helper()
	return ts.Login(t, "superadmin@platform.local", "ROLE_SUPERADMIN")
}

var (
	testPasswordHashOnce sync.Once
	testPasswordHash     string
)

// hashedTestPassword bcrypt-hashes TestPassword once and reuses it — bcrypt
// is deliberately slow, and re-hashing it per CreateUser call across dozens
// of tests adds up.
func hashedTestPassword(t *testing.T) string {
	t.Helper()
	testPasswordHashOnce.Do(func() {
		h, err := bcrypt.GenerateFromPassword([]byte(TestPassword), bcrypt.DefaultCost)
		require.NoError(t, err, "hash test password")
		testPasswordHash = string(h)
	})
	return testPasswordHash
}

// CreateUser creates a user (with the given role, creating the role first if
// needed) in the primary test clinic and returns the user's ID. This goes
// straight to the DB rather than through the now admin-gated /api/user/create
// endpoint — creating fixture users is arrange-phase setup, not the thing
// under test, and routing it through HTTP would need an admin token that
// itself needs a user to exist first.
func (ts *TestServer) CreateUser(t *testing.T, email, roleName string) uuid.UUID {
	t.Helper()
	return ts.CreateUserInClinic(t, ts.ClinicID, email, roleName)
}

// CreateUserInClinic is CreateUser scoped to an arbitrary clinic, for
// cross-tenant tests.
func (ts *TestServer) CreateUserInClinic(t *testing.T, clinicID uuid.UUID, email, roleName string) uuid.UUID {
	t.Helper()
	roleID := ts.EnsureRole(t, roleName)
	ctx := context.Background()

	var userID uuid.UUID
	err := ts.DB.QueryRow(ctx,
		`INSERT INTO users (role_id, clinic_id, email, password) VALUES ($1, $2, $3, $4) RETURNING id`,
		roleID, clinicID, email, hashedTestPassword(t),
	).Scan(&userID)
	require.NoError(t, err, "create user")

	_, err = ts.DB.Exec(ctx,
		`INSERT INTO profile (user_id, first_name, last_name, phone) VALUES ($1, 'Test', 'User', '')`,
		userID,
	)
	require.NoError(t, err, "create profile")

	return userID
}

// Login creates a user (with the given role) in the primary test clinic,
// logs in, and returns the access token plus user ID.
func (ts *TestServer) Login(t *testing.T, email, roleName string) (token string, userID uuid.UUID) {
	t.Helper()
	return ts.LoginInClinic(t, ts.ClinicID, email, roleName)
}

// LoginInClinic is Login scoped to an arbitrary clinic, for cross-tenant tests.
func (ts *TestServer) LoginInClinic(t *testing.T, clinicID uuid.UUID, email, roleName string) (token string, userID uuid.UUID) {
	t.Helper()
	ts.CreateUserInClinic(t, clinicID, email, roleName)

	var out struct {
		AccessToken string    `json:"access_token"`
		UserID      uuid.UUID `json:"user_id"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email":    email,
		"password": TestPassword,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return out.AccessToken, out.UserID
}

// AdminToken is the common case: an admin user in the primary test clinic, logged in.
func (ts *TestServer) AdminToken(t *testing.T) (token string, userID uuid.UUID) {
	t.Helper()
	return ts.Login(t, "admin@test.local", "ROLE_ADMIN")
}

// SecondClinicAdminToken spins up a second, isolated clinic with its own
