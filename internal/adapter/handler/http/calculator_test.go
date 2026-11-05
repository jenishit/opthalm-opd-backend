package http_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jenishit/opthalm-opd-backend/internal/testutil"
)

func TestCalculator_Transposition(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		Sphere   float64 `json:"sphere"`
		Cylinder float64 `json:"cylinder"`
		Axis     int     `json:"axis"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/calculators/transposition", token, map[string]any{
		"sphere": 2.0, "cylinder": -1.0, "axis": 90,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1.0, out.Sphere)
	assert.Equal(t, 1.0, out.Cylinder)
	assert.Equal(t, 0, out.Axis)
}

func TestCalculator_SphericalEquivalent(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		Result float64 `json:"result"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/calculators/spherical-equivalent", token, map[string]any{
		"sphere": 2.0, "cylinder": -1.0,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1.5, out.Result)
}

func TestCalculator_NearAdd(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		Result float64 `json:"result"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/calculators/near-add", token, map[string]any{
		"distance_sphere": -2.0, "add_power": 2.5,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0.5, out.Result)
}

func TestCalculator_VertexDistance(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		Result float64 `json:"result"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/calculators/vertex-distance", token, map[string]any{
		"power": 10.0, "from_distance_mm": 12.0, "to_distance_mm": 0.0,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// F' = F / (1 - d*F), d = (0-12)/1000 = -0.012 -> 10 / 1.12
	assert.InDelta(t, 8.9286, out.Result, 0.001)
}

func TestCalculator_TelescopeFOV(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		Result float64 `json:"result"`
	}
	resp := ts.DoData(t, http.MethodPost, "/api/calculators/telescope-fov", token, map[string]any{
		"true_fov_degrees": 5.0, "magnification": 8.0,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 40.0, out.Result)
}

func TestCalculator_IOLPower(t *testing.T) {
	ts := testutil.NewTestServer(t)
	token, _ := ts.AdminToken(t)

	var out struct {
		Result float64 `json:"result"`
	}

	// Average axial length (22-24.5mm): no A-constant correction.
	// P = 118 - 0.9*44 - 2.5*23.5 = 118 - 39.6 - 58.75 = 19.65
	resp := ts.DoData(t, http.MethodPost, "/api/calculators/iol-power", token, map[string]any{
		"axial_length_mm": 23.5, "avg_k_diopters": 44.0, "a_constant": 118.0,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.InDelta(t, 19.65, out.Result, 0.001)

	// Short eye (<20mm) gets a +3 A-constant correction.
	// P = 121 - 0.9*44 - 2.5*19.5 = 121 - 39.6 - 48.75 = 32.65
	resp = ts.DoData(t, http.MethodPost, "/api/calculators/iol-power", token, map[string]any{
		"axial_length_mm": 19.5, "avg_k_diopters": 44.0, "a_constant": 118.0,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.InDelta(t, 32.65, out.Result, 0.001)

	// A non-zero target refraction shifts the emmetropic power.
	resp = ts.DoData(t, http.MethodPost, "/api/calculators/iol-power", token, map[string]any{
		"axial_length_mm": 23.5, "avg_k_diopters": 44.0, "a_constant": 118.0, "target_refraction_d": -2.0,
	}, &out)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	assert.InDelta(t, 21.65, out.Result, 0.001)
}
