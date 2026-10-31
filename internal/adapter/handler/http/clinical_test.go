package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
)

func createTestVisit(t *testing.T, ts *testutil.TestServer, token string, userID, patientID uuid.UUID) uuid.UUID {
	t.Helper()
	var created struct {
		ID uuid.UUID `json:"id"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/visit", token, map[string]any{
		"patient_id":      patientID,
		"examine_by":      userID,
		"chief_complaint": "eye pain",
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return created.ID
}

func TestClinical_SymptomsUpsertRoundTrips(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient", "9800000020")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var symptoms struct {
		VisitID      uuid.UUID `json:"visit_id"`
		Redness      string    `json:"redness"`
		Swelling     string    `json:"swelling"`
		Headache     string    `json:"headache"`
		DoubleVision string    `json:"double_vision"`
		RecentFever  bool      `json:"recent_fever"`
	}
	resp := ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/symptoms", visitID), token, map[string]any{
		"redness":       "OD",
		"headache":      "OU",
		"double_vision": "OS",
		"recent_fever":  true,
	}, &symptoms)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, visitID, symptoms.VisitID)
	assert.Equal(t, "OD", symptoms.Redness)
	assert.Equal(t, "OU", symptoms.Headache)
	assert.Equal(t, "OS", symptoms.DoubleVision)
	assert.True(t, symptoms.RecentFever)
	// Fields not sent default to "none", matching the column defaults.
	assert.Equal(t, "none", symptoms.Swelling)

	// Re-PUT replaces the row rather than erroring.
	resp = ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/symptoms", visitID), token, map[string]any{
		"redness": "OU",
	}, &symptoms)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "OU", symptoms.Redness)
	assert.False(t, symptoms.RecentFever)
}

func TestClinical_ExaminationUpsertRoundTrips(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient 2", "9800000021")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var exam struct {
		VisitID uuid.UUID `json:"visit_id"`
		IOPOD   float64   `json:"iop_od"`
		BP      string    `json:"bp"`
	}
	resp := ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/examination", visitID), token, map[string]any{
		"iop_od": 16.5,
		"iop_os": 15.0,
		"bp":     "120/80",
	}, &exam)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, visitID, exam.VisitID)
	assert.Equal(t, 16.5, exam.IOPOD)
	assert.Equal(t, "120/80", exam.BP)
}

func TestClinical_ReadingsAppendAndUpsertByEyeStage(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient 3", "9800000022")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var reading struct {
		ID    uuid.UUID `json:"id"`
		Eye   string    `json:"eye"`
		Stage string    `json:"stage"`
	}
	resp := ts.DoData(t, http.MethodPost, fmt.Sprintf("/api/visit/%s/readings", visitID), token, map[string]any{
		"eye":    "OD",
		"stage":  "unaided",
		"sphere": -1.25,
		"va":     "6/9",
	}, &reading)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	firstID := reading.ID

	resp = ts.DoData(t, http.MethodPost, fmt.Sprintf("/api/visit/%s/readings", visitID), token, map[string]any{
		"eye":   "OS",
		"stage": "unaided",
		"va":    "6/6",
	}, &reading)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Re-submitting the same eye+stage upserts instead of duplicating.
	resp = ts.DoData(t, http.MethodPost, fmt.Sprintf("/api/visit/%s/readings", visitID), token, map[string]any{
		"eye":    "OD",
		"stage":  "unaided",
		"sphere": -1.50,
		"va":     "6/6",
	}, &reading)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, firstID, reading.ID)

	var list []struct {
		Eye string `json:"eye"`
		VA  string `json:"va"`
	}
	resp = ts.DoData(t, http.MethodGet, fmt.Sprintf("/api/visit/%s/readings", visitID), token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, list, 2)
}

func TestClinical_DiagnosesReplaceSwapsFullSet(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient 4", "9800000023")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	diagID1 := uuid.New()
	diagID2 := uuid.New()

	var diagnoses []struct {
		DiagnosisID uuid.UUID `json:"diagnosis_id"`
		Laterality  string    `json:"laterality"`
	}
	resp := ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/diagnoses", visitID), token, map[string]any{
		"diagnoses": []map[string]any{
			{"diagnosis_id": diagID1, "laterality": "OD", "diagnosis_name": "Conjunctivitis"},
			{"diagnosis_id": diagID2, "laterality": "OU", "diagnosis_name": "Dry eye"},
		},
	}, &diagnoses)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, diagnoses, 2)

	// Replacing with a smaller set drops the ones not resent.
	resp = ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/diagnoses", visitID), token, map[string]any{
		"diagnoses": []map[string]any{
			{"diagnosis_id": diagID1, "laterality": "OU"},
		},
	}, &diagnoses)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, diagnoses, 1)
	assert.Equal(t, diagID1, diagnoses[0].DiagnosisID)
	assert.Equal(t, "OU", diagnoses[0].Laterality)
}

func TestClinical_TreatmentsReplaceSwapsFullSet(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient 5", "9800000024")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var treatments []struct {
		Type       string `json:"type"`
		AdviceText string `json:"advice_text"`
	}
	resp := ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/treatments", visitID), token, map[string]any{
		"treatments": []map[string]any{
			{"type": "advice", "advice_text": "Warm compress"},
			{"type": "procedure", "procedure_name": "Lid scrub"},
		},
	}, &treatments)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, treatments, 2)
}

func TestClinical_FollowUpUpsertRoundTrips(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient 6", "9800000025")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var followUp struct {
		VisitID      uuid.UUID `json:"visit_id"`
		FollowUpDate string    `json:"follow_up_date"`
		Reason       string    `json:"reason"`
	}
	resp := ts.DoData(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/follow-up", visitID), token, map[string]any{
		"follow_up_date": "2026-11-15",
		"reason":         "Review IOP",
	}, &followUp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "2026-11-15", followUp.FollowUpDate)
	assert.Equal(t, "Review IOP", followUp.Reason)
}

func TestClinical_GetFullChartAggregatesEverything(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Clinical Patient 7", "9800000026")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	require.Equal(t, http.StatusOK, ts.Do(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/symptoms", visitID), token, map[string]any{
		"redness": "OD",
	}, nil).StatusCode)

	require.Equal(t, http.StatusOK, ts.Do(t, http.MethodPost, fmt.Sprintf("/api/visit/%s/readings", visitID), token, map[string]any{
		"eye": "OD", "stage": "unaided",
	}, nil).StatusCode)

	require.Equal(t, http.StatusOK, ts.Do(t, http.MethodPut, fmt.Sprintf("/api/visit/%s/follow-up", visitID), token, map[string]any{
		"follow_up_date": "2026-12-01",
	}, nil).StatusCode)

	var chart struct {
		Symptoms struct {
			Redness string `json:"redness"`
		} `json:"symptoms"`
		Examination any `json:"examination"`
		Readings    []struct {
			Eye string `json:"eye"`
		} `json:"readings"`
		Diagnoses  []any `json:"diagnoses"`
		Treatments []any `json:"treatments"`
		FollowUp   struct {
			FollowUpDate string `json:"follow_up_date"`
		} `json:"follow_up"`
	}
	resp := ts.DoData(t, http.MethodGet, fmt.Sprintf("/api/visit/%s/clinical", visitID), token, nil, &chart)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "OD", chart.Symptoms.Redness)
	assert.Nil(t, chart.Examination)
	require.Len(t, chart.Readings, 1)
	assert.Equal(t, "OD", chart.Readings[0].Eye)
	assert.Empty(t, chart.Diagnoses)
	assert.Empty(t, chart.Treatments)
	assert.Equal(t, "2026-12-01", chart.FollowUp.FollowUpDate)
}
