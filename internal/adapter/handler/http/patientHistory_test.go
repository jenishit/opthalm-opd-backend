package http_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

// createTestCondition inserts a row directly into the history_conditions
// catalog (there's no POST endpoint for it — conditions are seeded, not
// created through the API) and returns its ID.
func createTestCondition(t *testing.T, ts *testutil.TestServer, name string) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	err := ts.DB.QueryRow(context.Background(),
		`INSERT INTO history_conditions (name) VALUES ($1) RETURNING id`, name,
	).Scan(&id)
	require.NoError(t, err, "create test condition")
	return id
}

func TestPatientHistory_CreateAppendsEntry(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "History Patient", "9800000030")
	diabetesID := createTestCondition(t, ts, "Diabetes")
	hypertensionID := createTestCondition(t, ts, "Hypertension")

	systemic := "Metformin 500mg"
	notes := "Diagnosed 2024"

	var created struct {
		ID             uuid.UUID   `json:"id"`
		PatientID      uuid.UUID   `json:"patient_id"`
		UsesSpectacles bool        `json:"uses_spectacles"`
		ConditionIDs   []uuid.UUID `json:"condition_ids"`
		CreatedBy      uuid.UUID   `json:"created_by"`
	}
	resp := ts.DoData(t, http.MethodPost, fmt.Sprintf("/api/patient/%s/history", patientID), token, map[string]any{
		"uses_spectacles":      true,
		"systemic_medications": systemic,
		"notes":                notes,
		"condition_ids":        []uuid.UUID{diabetesID, hypertensionID},
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, patientID, created.PatientID)
	assert.True(t, created.UsesSpectacles)
	assert.NotEqual(t, uuid.Nil, created.CreatedBy)
	assert.ElementsMatch(t, []uuid.UUID{diabetesID, hypertensionID}, created.ConditionIDs)
}

func TestPatientHistory_ListReturnsNewestFirstAndNeverEdits(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "History Patient 2", "9800000031")
	conditionID := createTestCondition(t, ts, "Glaucoma")

	resp := ts.Do(t, http.MethodPost, fmt.Sprintf("/api/patient/%s/history", patientID), token, map[string]any{
		"uses_spectacles": false,
		"notes":           "first entry",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodPost, fmt.Sprintf("/api/patient/%s/history", patientID), token, map[string]any{
		"uses_spectacles": true,
		"notes":           "second entry",
		"condition_ids":   []uuid.UUID{conditionID},
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list []struct {
		Notes        string      `json:"notes"`
		ConditionIDs []uuid.UUID `json:"condition_ids"`
	}
	resp = ts.DoData(t, http.MethodGet, fmt.Sprintf("/api/patient/%s/history", patientID), token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, list, 2)
	// Newest first, and the older entry is untouched — a log, never edited.
	assert.Equal(t, "second entry", list[0].Notes)
	assert.Equal(t, []uuid.UUID{conditionID}, list[0].ConditionIDs)
	assert.Equal(t, "first entry", list[1].Notes)
	assert.Empty(t, list[1].ConditionIDs)
}

func TestPatientHistory_ListEmptyForNewPatient(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "History Patient 3", "9800000032")

	var list []any
	resp := ts.DoData(t, http.MethodGet, fmt.Sprintf("/api/patient/%s/history", patientID), token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, list)
}
