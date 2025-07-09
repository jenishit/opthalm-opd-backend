package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestVisit_CreateDefaultsStatusToScheduled(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Visit Patient", "9800000010")

	var created struct {
		ID     uuid.UUID `json:"id"`
		Status string    `json:"status"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/visit", token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"chief_complaint": "eye pain",
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "scheduled", created.Status)
}

func TestVisit_GetByID_VisitDateRoundTrips(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Visit Patient 2", "9800000011")

	var created struct {
		ID uuid.UUID `json:"id"`
	}
	ts.DoData(t, http.MethodPost, "/api/visit", token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"chief_complaint": "blurry vision",
	}, &created)

	var fetched struct {
		ID        uuid.UUID `json:"id"`
		VisitDate string    `json:"visit_date"`
		Status    string    `json:"status"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/visit/"+created.ID.String(), token, nil, &fetched)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, created.ID, fetched.ID)
	assert.NotEmpty(t, fetched.VisitDate, "visit_date must round-trip through the TIMESTAMP column scan")
}

func TestVisit_UpdateStatusDoesNotBlankOtherFields(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Visit Patient 3", "9800000012")
