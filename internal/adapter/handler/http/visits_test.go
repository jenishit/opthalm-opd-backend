package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
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

	var created struct {
		ID uuid.UUID `json:"id"`
	}
	ts.DoData(t, http.MethodPost, "/api/visit", token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"chief_complaint": "original complaint",
	}, &created)

	var before struct {
		VisitDate string `json:"visit_date"`
	}
	ts.DoData(t, http.MethodGet, "/api/visit/"+created.ID.String(), token, nil, &before)

	// Update status only, still resending required fields per the DTO contract.
	resp := ts.Do(t, http.MethodPatch, "/api/visit/"+created.ID.String(), token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"status":          "completed",
		"chief_complaint": "resolved",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var after struct {
		Status         string `json:"status"`
		VisitDate      string `json:"visit_date"`
		CheifComplaint string `json:"chief_complaint"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/visit/"+created.ID.String(), token, nil, &after)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "completed", after.Status)
	assert.Equal(t, "resolved", after.CheifComplaint)
	assert.Equal(t, before.VisitDate, after.VisitDate, "visit_date must not be blanked out by an update that omits it")
}

func TestVisit_ListByPatient(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Visit Patient 4", "9800000013")

	ts.Do(t, http.MethodPost, "/api/visit", token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"chief_complaint": "first visit",
	}, nil)
	ts.Do(t, http.MethodPost, "/api/visit", token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"chief_complaint": "second visit",
	}, nil)

	var out struct {
		Visits []struct {
			ID uuid.UUID `json:"id"`
		} `json:"visits"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/visit/patient/"+patientID.String(), token, nil, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, out.Visits, 2)
}
