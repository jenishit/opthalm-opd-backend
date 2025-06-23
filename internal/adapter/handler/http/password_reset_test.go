package http_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// extractToken pulls the plaintext token out of a captured email body. Every
// AuthService verification email follows "...: <token>\n\n...", so splitting
// on the first ": " and then the first "\n" isolates it.
func extractToken(t *testing.T, body string) string {
	t.Helper()
	_, rest, ok := strings.Cut(body, ": ")
	require.True(t, ok, "expected email body to contain a token after ': ' — got %q", body)
	token, _, _ := strings.Cut(rest, "\n")
	require.NotEmpty(t, token)
	return token
}

func TestPasswordReset_FullFlow(t *testing.T) {
	ts := testutil.NewTestServer(t)
	email := "resetme@test.local"
	access, refresh, _ := ts.LoginWithRefresh(t, email, "ROLE_ADMIN")
	require.NotEmpty(t, access)

	resp := ts.Do(t, http.MethodPost, "/api/auth/password-reset/request", "", map[string]string{"email": email}, nil)
