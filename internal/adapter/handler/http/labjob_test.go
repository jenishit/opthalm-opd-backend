package http_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestLabJob_CreateAndStatusPipeline(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Lab Patient", "9800000030")

	var created struct {
		ID     uuid.UUID `json:"id"`
		Status string    `json:"status"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/lab-jobs", token, map[string]any{
		"patient_id": patientID,
		"job_type":   "lens_fitting",
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "in_fitting", created.Status)

	resp = ts.Do(t, http.MethodPatch, "/api/lab-jobs/"+created.ID.String()+"/status", token, map[string]any{
		"status": "ready_to_deliver", "notes": "grinding done",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodPatch, "/api/lab-jobs/"+created.ID.String()+"/status", token, map[string]any{
		"status": "delivered",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched struct {
		Status      string  `json:"status"`
		DeliveredAt *string `json:"delivered_at"`
		History     []struct {
			Status string `json:"status"`
		} `json:"history"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/lab-jobs/"+created.ID.String(), token, nil, &fetched)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "delivered", fetched.Status)
	require.NotNil(t, fetched.DeliveredAt)
	require.Len(t, fetched.History, 3)
	assert.Equal(t, "in_fitting", fetched.History[0].Status)
	assert.Equal(t, "ready_to_deliver", fetched.History[1].Status)
	assert.Equal(t, "delivered", fetched.History[2].Status)
}
