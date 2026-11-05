package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
)

// decodeEnvelope unwraps the {success, message, data} response envelope
// from a raw *http.Response (used for multipart requests, which bypass
// TestServer.DoData's JSON-only request encoding) and decodes the `data`
// field into out.
func decodeEnvelope(t *testing.T, resp *http.Response, out any) error {
	t.Helper()
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	if err := json.Unmarshal(body, &env); err != nil {
		return err
	}
	return json.Unmarshal(env.Data, out)
}

// multipartUpload builds a multipart/form-data request with one file part
// and the given extra form fields, mirroring the pattern already used for
// the QR-code upload test in clinicConfig_test.go.
func multipartUpload(t *testing.T, url, token, fileField, filename, contentType string, fileBytes []byte, fields map[string]string) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fileField, filename))
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	require.NoError(t, err)
	_, err = part.Write(fileBytes)
	require.NoError(t, err)

	for k, v := range fields {
		require.NoError(t, writer.WriteField(k, v))
	}
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

const minimalPNG = "\x89PNG\r\n\x1a\n"

func TestAttachment_UploadListDelete(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Attachment Patient", "9800000040")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var uploaded struct {
		ID        uuid.UUID `json:"id"`
		VisitID   uuid.UUID `json:"visit_id"`
		PatientID uuid.UUID `json:"patient_id"`
		Category  string    `json:"category"`
		Eye       string    `json:"eye"`
		URL       string    `json:"url"`
	}
	resp := multipartUpload(t,
		ts.Server.URL+"/api/visit/"+visitID.String()+"/attachments",
		token, "file", "fundus.png", "image/png", []byte(minimalPNG),
		map[string]string{"category": "fundus_photo", "eye": "OD"},
	)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, decodeEnvelope(t, resp, &uploaded))
	assert.Equal(t, visitID, uploaded.VisitID)
	assert.Equal(t, patientID, uploaded.PatientID)
	assert.Equal(t, "fundus_photo", uploaded.Category)
	assert.Equal(t, "OD", uploaded.Eye)
	assert.NotEmpty(t, uploaded.URL)

	// List
	var listed []struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/visit/"+visitID.String()+"/attachments", token, nil, &listed)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, listed, 1)
	assert.Equal(t, uploaded.ID, listed[0].ID)

	// Delete
	resp = ts.Do(t, http.MethodDelete, "/api/visit/"+visitID.String()+"/attachments/"+uploaded.ID.String(), token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.DoData(t, http.MethodGet, "/api/visit/"+visitID.String()+"/attachments", token, nil, &listed)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, listed)
}

func TestAttachment_RejectsUnsupportedContentType(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Attachment Patient 2", "9800000041")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	resp := multipartUpload(t,
		ts.Server.URL+"/api/visit/"+visitID.String()+"/attachments",
		token, "file", "scan.txt", "text/plain", []byte("not an image"),
		map[string]string{"category": "fundus_photo"},
	)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestConsent_RecordAndList(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, userID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, token, "Consent Patient", "9800000042")
	visitID := createTestVisit(t, ts, token, userID, patientID)

	var recorded struct {
		ID               uuid.UUID `json:"id"`
		VisitID          uuid.UUID `json:"visit_id"`
		ConsentType      string    `json:"consent_type"`
		SignedByName     string    `json:"signed_by_name"`
		SignedByRelation string    `json:"signed_by_relation"`
		SignatureURL     string    `json:"signature_url"`
	}
	resp := multipartUpload(t,
		ts.Server.URL+"/api/visit/"+visitID.String()+"/consents",
		token, "signature", "sig.png", "image/png", []byte(minimalPNG),
		map[string]string{"consent_type": "surgery", "signed_by_name": "Jane Doe"},
	)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, decodeEnvelope(t, resp, &recorded))
	assert.Equal(t, visitID, recorded.VisitID)
	assert.Equal(t, "surgery", recorded.ConsentType)
	assert.Equal(t, "Jane Doe", recorded.SignedByName)
	// signed_by_relation defaults to "self" when omitted.
	assert.Equal(t, "self", recorded.SignedByRelation)
	assert.NotEmpty(t, recorded.SignatureURL)

	var listed []struct {
		ID uuid.UUID `json:"id"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/visit/"+visitID.String()+"/consents", token, nil, &listed)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, listed, 1)
	assert.Equal(t, recorded.ID, listed[0].ID)
}

// TestClinical_NonClinicalRoleForbidden asserts the clinical sub-resources
// (chart, attachments, consents) reject a role that isn't admin/doctor/
// optometrist, while the plain visit record stays reachable to it.
func TestClinical_NonClinicalRoleForbidden(t *testing.T) {
	ts := testutil.NewTestServer(t)
	adminToken, adminID := ts.AdminToken(t)
	patientID := ts.CreatePatient(t, adminToken, "Gated Patient", "9800000043")
	visitID := createTestVisit(t, ts, adminToken, adminID, patientID)

	billingToken, _ := ts.Login(t, "billing@test.local", "ROLE_BILLING")

	resp := ts.Do(t, http.MethodGet, "/api/visit/"+visitID.String(), billingToken, nil, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "front-desk/billing roles can still read the visit record itself")

	// roleMiddleware rejects a role mismatch with domain.ErrUnauthorized,
	// which maps to 401 (not 403) throughout this app — same as the
	// existing billing/inventory/lab-job role gates.
	resp = ts.Do(t, http.MethodGet, "/api/visit/"+visitID.String()+"/clinical", billingToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "billing role must not read the clinical chart")

	resp = ts.Do(t, http.MethodGet, "/api/visit/"+visitID.String()+"/attachments", billingToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "billing role must not read clinical attachments")

	resp = ts.Do(t, http.MethodGet, "/api/visit/"+visitID.String()+"/consents", billingToken, nil, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "billing role must not read consents")

	// A doctor, on the other hand, is allowed through the gate.
	doctorToken, _ := ts.Login(t, "doctor@test.local", "ROLE_DOCTOR")
	resp = ts.Do(t, http.MethodGet, "/api/visit/"+visitID.String()+"/clinical", doctorToken, nil, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
