package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
)

func TestClinic_InsertGetAllGetByIDUpdate(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var created struct {
		ID         uuid.UUID `json:"id"`
		ClinicName string    `json:"clinic_name"`
	}
	// Deliberately omit optional fields (tagline/address/phone/email/report_footer)
	// — this used to panic with a nil-pointer dereference in the response mapper.
	resp := ts.DoData(t, http.MethodPost, "/api/admin/clinic", token, map[string]any{
		"clinic_name":     "Test Eye Clinic",
		"registration_no": "REG-001",
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "Test Eye Clinic", created.ClinicName)

	var all []struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic", token, nil, &all)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// GetAllClinics is a deliberately global, unscoped listing (a
	// superadmin-style operation) — the harness's own primary test clinic is
	// already in there alongside the one just created.
	assert.GreaterOrEqual(t, len(all), 2)

	var got struct {
		ClinicName string `json:"clinic_name"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/"+created.ID.String(), token, nil, &got)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Test Eye Clinic", got.ClinicName)

	resp = ts.Do(t, http.MethodPatch, "/api/admin/clinic/"+created.ID.String(), token, map[string]any{
		"clinic_name":     "Renamed Clinic",
		"registration_no": "REG-001",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/"+created.ID.String(), token, nil, &got)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Renamed Clinic", got.ClinicName)
}

func TestClinic_NonAdminRejected(t *testing.T) {
	ts := testutil.NewTestServer(t)
	doctorToken, _ := ts.Login(t, "clinicuser@test.local", "ROLE_DOCTOR")

	resp := ts.Do(t, http.MethodGet, "/api/admin/clinic", doctorToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
