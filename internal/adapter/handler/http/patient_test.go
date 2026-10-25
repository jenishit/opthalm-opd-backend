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
	var list []struct {
		PatientID uuid.UUID `json:"patient_id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/patient", token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, list, 1)

	// SEARCH
	var searchResults []struct {
		PatientID uuid.UUID `json:"patient_id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/patient/search?query=Alice", token, nil, &searchResults)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, searchResults, 1)
	assert.Equal(t, created.PatientID, searchResults[0].PatientID)

	// UPDATE with address/occupation omitted — must not panic (this used to
	// crash with a nil-pointer dereference before the fix).
	resp = ts.Do(t, http.MethodPatch, "/api/patient/"+created.PatientID.String(), token, map[string]any{
		"full_name": "Alice Updated",
		"dob":       "1990-01-01",
		"gender":    "female",
		"phone":     "9800000001",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.DoData(t, http.MethodGet, "/api/patient/"+created.PatientID.String(), token, nil, &fetched)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Alice Updated", fetched.FullName)
}

func TestPatient_GetNonexistentReturns404(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	resp := ts.Do(t, http.MethodGet, "/api/patient/"+uuid.New().String(), token, nil, nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPatient_SoftDeleteExcludesFromReads(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	patientID := ts.CreatePatient(t, token, "To Delete", "9800000002")

	resp := ts.Do(t, http.MethodPatch, "/api/patient/"+patientID.String()+"/delete", token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodGet, "/api/patient/"+patientID.String(), token, nil, nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "a soft-deleted patient should not be gettable")
}
