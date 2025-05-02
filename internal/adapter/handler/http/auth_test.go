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

