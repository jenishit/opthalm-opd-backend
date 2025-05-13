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
