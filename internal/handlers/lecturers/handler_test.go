package lecturers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Contract-parity test against internal/handlers/testdata/baseline_responses.json
// (teacher_list_limit_offset case). Requires a reachable MySQL "sekolah" database.

func TestGetLecturers_DefaultAndShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/teacher/list", GetLecturers)

	req := httptest.NewRequest(http.MethodGet, "/teacher/list?limit=2&offset=0", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, ok := body["records"]; !ok {
		t.Errorf("expected 'records' field to be present")
	}
	if got := body["records"].(float64); got != 2 {
		t.Errorf("records = %v, want 2", got)
	}

	data, ok := body["data"].([]any)
	if !ok || len(data) == 0 {
		t.Fatalf("expected non-empty data array, got %v", body["data"])
	}
	row := data[0].(map[string]any)
	if _, present := row["kd_guru"]; present {
		t.Errorf("kd_guru must not appear in JSON output")
	}
	if _, present := row["kd_kompetensi"]; present {
		t.Errorf("kd_kompetensi must not appear in JSON output")
	}
	for _, field := range []string{"nip", "nama", "alamat", "telp"} {
		if _, present := row[field]; !present {
			t.Errorf("expected field %q to be present in lecturer row", field)
		}
	}
}
