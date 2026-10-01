package http_test

import (
	"net/http"
	"testing"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// TestRateLimit_LoginIsLimited exercises the Redis-backed rate limiter on
// /api/auth/login (10 requests/minute/IP): once the limit is exceeded, the
// server must respond 429 with a Retry-After header, regardless of whether
// the credentials themselves are valid.
func TestRateLimit_LoginIsLimited(t *testing.T) {
	ts := testutil.NewTestServer(t)

	body := map[string]string{"email": "nobody@ratelimit.test", "password": "wrong-password"}

	var last *http.Response
	for range 11 {
		last = ts.Do(t, http.MethodPost, "/api/auth/login", "", body, nil)
	}

	require.Equal(t, http.StatusTooManyRequests, last.StatusCode, "11th request within the window should be rate-limited")
	require.NotEmpty(t, last.Header.Get("Retry-After"), "rate-limited response should advertise Retry-After")
}
