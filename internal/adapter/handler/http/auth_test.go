package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
)

func TestAuth_RoleCreate(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		ID       uuid.UUID `json:"ID"`
		RoleName string    `json:"RoleName"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/role/create", token, map[string]string{"role_name": "ROLE_NEWLYCREATED"}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, uuid.Nil, out.ID)
	assert.Equal(t, "ROLE_NEWLYCREATED", out.RoleName)
}

func TestAuth_RoleCreateRequiresAdmin(t *testing.T) {
	ts := testutil.NewTestServer(t)

	resp := ts.Do(t, http.MethodPost, "/api/role/create", "", map[string]string{"role_name": "ROLE_SNEAKY"}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_UserCreate(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		ID uuid.UUID `json:"ID"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/user/create", token, map[string]string{
		"first_name": "New",
		"last_name":  "User",
		"email":      "newuser@test.local",
		"password":   testutil.TestPassword,
		"role_name":  "ROLE_ADMIN",
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, uuid.Nil, out.ID)
}

func TestAuth_LoginSuccess(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)

	assert.NotEmpty(t, token)
	assert.NotEqual(t, uuid.Nil, userID)
}

func TestAuth_LoginReturnsRefreshToken(t *testing.T) {
	ts := testutil.NewTestServer(t)
	access, refresh, userID := ts.LoginWithRefresh(t, "refreshuser1@test.local", "ROLE_ADMIN")

	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)
	assert.NotEqual(t, uuid.Nil, userID)
}

func TestAuth_RefreshIssuesNewPairAndRotatesOldOne(t *testing.T) {
	ts := testutil.NewTestServer(t)
	_, refresh1, _ := ts.LoginWithRefresh(t, "refreshuser2@test.local", "ROLE_ADMIN")

	var out struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": refresh1,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, out.AccessToken)
	require.NotEmpty(t, out.RefreshToken)
	assert.NotEqual(t, refresh1, out.RefreshToken, "refresh should rotate to a new token")

	// The old (now-rotated-out) refresh token must no longer work on its own.
	resp = ts.Do(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": refresh1,
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_RefreshReuseDetectionRevokesWholeChain(t *testing.T) {
	ts := testutil.NewTestServer(t)
	_, refresh1, _ := ts.LoginWithRefresh(t, "refreshuser3@test.local", "ROLE_ADMIN")

	var out struct {
		RefreshToken string `json:"refresh_token"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": refresh1,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	refresh2 := out.RefreshToken

	// Replay the already-rotated-out refresh1 — this should be treated as
	// theft and revoke the whole chain, including the still-fresh refresh2.
	resp = ts.Do(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": refresh1,
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp = ts.Do(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": refresh2,
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "reuse detection should have revoked the whole session chain")
}

func TestAuth_LogoutRevokesSessionSoRefreshFails(t *testing.T) {
	ts := testutil.NewTestServer(t)
	_, refresh, _ := ts.LoginWithRefresh(t, "refreshuser4@test.local", "ROLE_ADMIN")

	resp := ts.Do(t, http.MethodPost, "/api/auth/logout", "", map[string]string{
		"refresh_token": refresh,
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": refresh,
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_RefreshWithGarbageTokenRejected(t *testing.T) {
	ts := testutil.NewTestServer(t)

	resp := ts.Do(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{
		"refresh_token": "not-a-real-refresh-token",
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_LoginWrongPassword(t *testing.T) {
	ts := testutil.NewTestServer(t)
	ts.CreateUser(t, "wrongpw@test.local", "ROLE_ADMIN")

	resp := ts.Do(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email":    "wrongpw@test.local",
		"password": "not-the-real-password",
	}, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_LoginUnknownEmail(t *testing.T) {
	ts := testutil.NewTestServer(t)

	resp := ts.Do(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email":    "nobody@test.local",
		"password": "whatever",
	}, nil)
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_ProtectedRouteWithoutToken(t *testing.T) {
	ts := testutil.NewTestServer(t)

	resp := ts.Do(t, http.MethodGet, "/api/patient", "", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuth_ProtectedRouteWithMalformedToken(t *testing.T) {
	ts := testutil.NewTestServer(t)

	resp := ts.Do(t, http.MethodGet, "/api/patient", "not-a-real-jwt", nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRBAC_NonAdminRejectedFromBillingInventoryReports(t *testing.T) {
	ts := testutil.NewTestServer(t)
	doctorToken, _ := ts.Login(t, "doctor@test.local", "ROLE_DOCTOR")

	for _, path := range []string{"/api/billing/invoice", "/api/inventory/items", "/api/reports/dues"} {
		resp := ts.Do(t, http.MethodGet, path, doctorToken, nil, nil)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "ROLE_DOCTOR should be rejected from %s", path)
	}
}

func TestRBAC_NonAdminAllowedOnCalculators(t *testing.T) {
	ts := testutil.NewTestServer(t)
	doctorToken, _ := ts.Login(t, "doctor2@test.local", "ROLE_DOCTOR")

	resp := ts.Do(t, http.MethodPost, "/api/calculators/spherical-equivalent", doctorToken, map[string]float64{
		"sphere": 1, "cylinder": 1,
	}, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRBAC_AdminAllowedEverywhere(t *testing.T) {
	ts := testutil.NewTestServer(t)
	adminToken, _ := ts.AdminToken(t)

	for _, path := range []string{"/api/billing/invoice", "/api/inventory/items", "/api/reports/dues"} {
		resp := ts.Do(t, http.MethodGet, path, adminToken, nil, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "ROLE_ADMIN should be allowed on %s", path)
	}
}

func TestCORS_PreflightAllowedOrigin(t *testing.T) {
	ts := testutil.NewTestServer(t)

	req, err := http.NewRequest(http.MethodOptions, ts.Server.URL+"/api/patient", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORS_DisallowedOriginGetsNoCORSHeader(t *testing.T) {
	ts := testutil.NewTestServer(t)

	req, err := http.NewRequest(http.MethodOptions, ts.Server.URL+"/api/patient", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://evil.example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
}
