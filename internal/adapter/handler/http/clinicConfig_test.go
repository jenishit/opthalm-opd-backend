package http_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenish-brainztechs/go-backend/internal/testutil"
)

func TestClinicProfile_ExtendedFieldsRoundTrip(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var created struct {
		ID uuid.UUID `json:"id"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/admin/clinic", token, map[string]any{
		"clinic_name":     "Clearview Eye Clinic",
		"registration_no": "REG-900",
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.Do(t, http.MethodPatch, fmt.Sprintf("/api/admin/clinic/%s", created.ID), token, map[string]any{
		"clinic_name":                 "Clearview Eye Clinic",
		"registration_no":             "REG-900",
		"workspace_slug":              "clearview",
		"default_currency":            "AUD",
		"timezone":                    "Australia/Melbourne",
		"require_consent_before_exam": true,
		"low_stock_notifications":     false,
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var got struct {
		WorkspaceSlug            string `json:"workspace_slug"`
		DefaultCurrency          string `json:"default_currency"`
		Timezone                 string `json:"timezone"`
		RequireConsentBeforeExam bool   `json:"require_consent_before_exam"`
		LowStockNotifications    bool   `json:"low_stock_notifications"`
	}
	resp = ts.DoData(t, http.MethodGet, fmt.Sprintf("/api/admin/clinic/%s", created.ID), token, nil, &got)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "clearview", got.WorkspaceSlug)
	assert.Equal(t, "AUD", got.DefaultCurrency)
	assert.Equal(t, "Australia/Melbourne", got.Timezone)
	assert.True(t, got.RequireConsentBeforeExam)
	assert.False(t, got.LowStockNotifications)
}

func TestClinicLocations_CreateAndList(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var created struct {
		ID     uuid.UUID `json:"id"`
		Status string    `json:"status"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/admin/clinic/locations", token, map[string]any{
		"name":       "Main Branch",
		"exam_rooms": 3,
	}, &created)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.Equal(t, "active", created.Status) // default when omitted

	resp = ts.Do(t, http.MethodPost, "/api/admin/clinic/locations", token, map[string]any{
		"name":       "New Branch",
		"exam_rooms": 1,
		"status":     "opening_soon",
	}, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/locations", token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, list, 2)
	assert.Equal(t, "Main Branch", list[0].Name)
	assert.Equal(t, "opening_soon", list[1].Status)
}

func TestClinicNotifications_DefaultsThenUpsert(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var prefs struct {
		EmailNotifications bool `json:"email_notifications"`
		WeeklyDigest       bool `json:"weekly_digest"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/admin/clinic/notifications", token, nil, &prefs)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, prefs.EmailNotifications) // column default before any save

	resp = ts.DoData(t, http.MethodPut, "/api/admin/clinic/notifications", token, map[string]any{
		"email_notifications":       false,
		"sms_appointment_reminders": true,
		"low_stock_alerts":          true,
		"new_order_notifications":   false,
		"weekly_digest":             true,
	}, &prefs)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, prefs.EmailNotifications)
	assert.True(t, prefs.WeeklyDigest)

	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/notifications", token, nil, &prefs)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, prefs.EmailNotifications)
	assert.True(t, prefs.WeeklyDigest)
}

func TestClinicPaymentConfig_UpsertWithNestedBankDetailsAndQRUpload(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var cfg struct {
		TaxLabel string  `json:"tax_label"`
		TaxRate  float64 `json:"tax_rate"`
		Methods  struct {
			Cash         bool `json:"cash"`
			BankTransfer bool `json:"bank_transfer"`
		} `json:"methods"`
		BankDetail *struct {
			BankName string `json:"bank_name"`
		} `json:"bank_details"`
	}
	resp := ts.DoData(t, http.MethodPut, "/api/admin/clinic/payment-config", token, map[string]any{
		"tax_label": "GST",
		"tax_rate":  13,
		"methods": map[string]any{
			"cash":          true,
			"bank_transfer": true,
			"qr_wallet":     true,
		},
		"bank_details": map[string]any{
			"bank_name":      "Nepal Bank",
			"account_holder": "Clearview Eye Clinic",
			"account_number": "1234567890",
		},
		"qr_wallet": map[string]any{
			"provider":    "eSewa",
			"merchant_id": "MERCH-1",
		},
	}, &cfg)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GST", cfg.TaxLabel)
	assert.Equal(t, float64(13), cfg.TaxRate)
	assert.True(t, cfg.Methods.Cash)
	assert.True(t, cfg.Methods.BankTransfer)
	require.NotNil(t, cfg.BankDetail)
	assert.Equal(t, "Nepal Bank", cfg.BankDetail.BankName)

	// multipart QR code upload
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// CreateFormFile hardcodes Content-Type: application/octet-stream, which
	// the handler would reject — build the part header manually so it's a
	// real image/png upload.
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="qr.png"`)
	header.Set("Content-Type", "image/png")
	part, err := writer.CreatePart(header)
	require.NoError(t, err)
	// A minimal valid PNG header is enough — the handler only checks
	// Content-Type, it doesn't decode the image.
	_, err = part.Write([]byte("\x89PNG\r\n\x1a\n"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, ts.Server.URL+"/api/admin/clinic/payment-config/qr-code", &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	httpResp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer httpResp.Body.Close()
	require.Equal(t, http.StatusOK, httpResp.StatusCode)

	var afterUpload struct {
		QRWallet *struct {
			QRAssetURL string `json:"qr_asset_url"`
		} `json:"qr_wallet"`
	}
	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/payment-config", token, nil, &afterUpload)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, afterUpload.QRWallet)
	assert.Contains(t, afterUpload.QRWallet.QRAssetURL, "/uploads/clinics/")
	assert.Contains(t, afterUpload.QRWallet.QRAssetURL, "qr-code.png")
}

func TestClinicPosConfig_DefaultsThenUpsert(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var cfg struct {
		BarcodeSymbology string `json:"barcode_symbology"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/admin/clinic/pos-config", token, nil, &cfg)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "code128", cfg.BarcodeSymbology)

	resp = ts.DoData(t, http.MethodPut, "/api/admin/clinic/pos-config", token, map[string]any{
		"barcode_symbology":        "ean13",
		"auto_generate_barcode":    true,
		"barcode_scanning_enabled": true,
		"require_pin_for_refunds":  true,
		"till_float_amount":        500,
	}, &cfg)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ean13", cfg.BarcodeSymbology)
}

func TestClinicPrivacyConfig_DefaultsThenUpsertAndAuditLog(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var cfg struct {
		RetentionYears int  `json:"retention_years"`
		LoginAlerts    bool `json:"login_alerts"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/admin/clinic/privacy-config", token, nil, &cfg)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 7, cfg.RetentionYears)

	resp = ts.DoData(t, http.MethodPut, "/api/admin/clinic/privacy-config", token, map[string]any{
		"retention_years":                 10,
		"allow_patient_data_export":       true,
		"allow_patient_deletion_requests": true,
		"require_2fa":                     true,
		"login_alerts":                    true,
	}, &cfg)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 10, cfg.RetentionYears)

	var auditLog []any
	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/audit-log", token, nil, &auditLog)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, auditLog) // nothing writes to it yet
}

func TestClinicIntegrations_ListConnectDisconnect(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var list []struct {
		Provider  string `json:"provider"`
		Connected bool   `json:"connected"`
	}
	resp := ts.DoData(t, http.MethodGet, "/api/admin/clinic/integrations", token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, list, 4) // every known provider, even unconnected
	for _, i := range list {
		assert.False(t, i.Connected)
	}

	var connected struct {
		Provider    string `json:"provider"`
		Connected   bool   `json:"connected"`
		ConnectedAt string `json:"connected_at"`
	}
	resp = ts.DoData(t, http.MethodPost, "/api/admin/clinic/integrations/xero/connect", token, nil, &connected)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "xero", connected.Provider)
	assert.True(t, connected.Connected)
	assert.NotEmpty(t, connected.ConnectedAt)

	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/integrations", token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	for _, i := range list {
		if i.Provider == "xero" {
			assert.True(t, i.Connected)
		}
	}

	resp = ts.Do(t, http.MethodDelete, "/api/admin/clinic/integrations/xero", token, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp = ts.DoData(t, http.MethodGet, "/api/admin/clinic/integrations", token, nil, &list)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	for _, i := range list {
		assert.False(t, i.Connected)
	}

	// Unknown provider is rejected rather than silently accepted.
	resp = ts.Do(t, http.MethodPost, "/api/admin/clinic/integrations/not-a-real-provider/connect", token, nil, nil)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
