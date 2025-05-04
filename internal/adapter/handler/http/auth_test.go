package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
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
