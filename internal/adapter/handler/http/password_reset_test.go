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
	require.Equal(t, http.StatusOK, resp.StatusCode)

	token := extractToken(t, ts.Emails.Last(t).Body)

	resp = ts.Do(t, http.MethodPost, "/api/auth/password-reset/confirm", "", map[string]string{
		"token": token, "new_password": "BrandNewPassword456!",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Old password must no longer work.
	resp = ts.Do(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email": email, "password": testutil.TestPassword,
	}, nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// New password does.
	resp = ts.Do(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email": email, "password": "BrandNewPassword456!",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// The refresh token issued before the reset must have been revoked along
	// with every other session for this user.
	resp = ts.Do(t, http.MethodPost, "/api/auth/refresh", "", map[string]string{"refresh_token": refresh}, nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// The reset token itself is single-use.
	resp = ts.Do(t, http.MethodPost, "/api/auth/password-reset/confirm", "", map[string]string{
		"token": token, "new_password": "AnotherPassword789!",
	}, nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
