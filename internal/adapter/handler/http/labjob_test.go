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
