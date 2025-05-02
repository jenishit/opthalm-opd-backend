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
