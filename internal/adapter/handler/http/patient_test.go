package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestPatient_CreateGetListSearchUpdate(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var created struct {
		PatientID  uuid.UUID `json:"patient_id"`
		FullName   string    `json:"full_name"`
		DOB        string    `json:"dob"`
		Occupation *string   `json:"occupation"`
		Address    *string   `json:"address"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/patient", token, map[string]any{
		"full_name":  "Alice Test",
		"dob":        "1990-01-01",
		"gender":     "female",
		"phone":      "9800000001",
		"address":    "123 Main St",
		"occupation": "Teacher",
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEqual(t, uuid.Nil, created.PatientID)
	assert.Equal(t, "Alice Test", created.FullName)
	assert.Equal(t, "1990-01-01", created.DOB)
	require.NotNil(t, created.Occupation)
	assert.Equal(t, "Teacher", *created.Occupation)

	// GET by ID
	var fetched struct {
		PatientID uuid.UUID `json:"patient_id"`
		FullName  string    `json:"full_name"`
		DOB       string    `json:"dob"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/patient/"+created.PatientID.String(), token, nil, &fetched)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, created.PatientID, fetched.PatientID)
	assert.Equal(t, "1990-01-01", fetched.DOB, "DOB must round-trip through the DATE column scan")

	// LIST
