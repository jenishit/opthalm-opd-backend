package http_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

// seedMedicine inserts a medicine row directly (there's no create-via-API
// for catalog entries — they're seeded, not user-created) and returns its ID.
func seedMedicine(t *testing.T, ts *testutil.TestServer, name string) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	err := ts.DB.QueryRow(context.Background(),
		`INSERT INTO medicines (medicine_name) VALUES ($1) RETURNING id`, name,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestCatalog_MedicinesListSearchGetUpdateDelete(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)
	medID := seedMedicine(t, ts, "Timolol")

	var list []struct {
		ID uuid.UUID `json:"id"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/admin/catalog/medicines", token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, list, 1)

	var search []struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/catalog/medicines/search?query=Timo", token, nil, &search)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, search, 1)

	var got struct {
		ID           uuid.UUID `json:"id"`
		MedicineName string    `json:"medicine_name"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/catalog/medicines/"+medID.String(), token, nil, &got)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Timolol", got.MedicineName)

	resp = ts.Do(t, http.MethodPatch, "/api/admin/catalog/medicines/"+medID.String(), token, map[string]any{
		"medicine_name": "Timolol Maleate",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.DoData(t, http.MethodGet, "/api/admin/catalog/medicines/"+medID.String(), token, nil, &got)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Timolol Maleate", got.MedicineName)

	resp = ts.Do(t, http.MethodPatch, "/api/admin/catalog/medicines/"+medID.String()+"/delete", token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodGet, "/api/admin/catalog/medicines/"+medID.String(), token, nil, nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCatalog_GetNonexistentMedicineReturns404(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	resp := ts.Do(t, http.MethodGet, "/api/admin/catalog/medicines/"+uuid.New().String(), token, nil, nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCatalog_NonAdminRejected(t *testing.T) {
	ts := testutil.NewTestServer(t)
	doctorToken, _ := ts.Login(t, "cataloguser@test.local", "ROLE_DOCTOR")

	resp := ts.Do(t, http.MethodGet, "/api/admin/catalog/medicines", doctorToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
