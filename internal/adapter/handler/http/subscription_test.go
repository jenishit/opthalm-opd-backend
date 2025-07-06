package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestSubscription_CancelledClinicGets402OnTenantRoutesButCanStillLogin(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	ts.SetSubscriptionStatus(t, ts.ClinicID, "cancelled", time.Now().Add(24*time.Hour))

	resp := ts.Do(t, http.MethodGet, "/api/patient", token, nil, nil)
	assert.Equal(t, http.StatusPaymentRequired, resp.StatusCode)

	// Login must still work — the client needs to be able to authenticate to
	// even find out the subscription is the problem.
	var login struct {
		AccessToken string `json:"access_token"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/auth/login", "", map[string]string{
		"email":    "admin@test.local",
		"password": testutil.TestPassword,
	}, &login)
