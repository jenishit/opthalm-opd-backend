package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
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
