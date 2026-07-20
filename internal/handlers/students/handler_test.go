package students

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Contract-parity tests against internal/handlers/testdata/baseline_responses.json
// (captured from the original go-chi server). Requires a reachable MySQL "sekolah"
// database (see tools.Init) -- run with the local DSN wired up.

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/student/list", LoadList)
	r.GET("/student/detail", LoadByID)
	return r
}

func TestStudentList_Default(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/student/list", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q (no charset suffix)", ct, "application/json")
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body["code"].(float64) != 200 {
		t.Errorf("code = %v, want 200", body["code"])
	}
	data, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("data is not an array: %T", body["data"])
	}
	if len(data) == 0 {
		t.Errorf("expected default list to return rows")
	}

	// Byte-exact parity with baseline_responses.json (student_list_default),
	// including the trailing newline emitted by json.NewEncoder.Encode (the
	// original go-chi handler used json.NewEncoder(w).Encode, NOT gin's c.JSON,
	// which would otherwise drop that trailing byte).
	wantBody := "{\"code\":200,\"info\":\"Data Found\",\"data\":[{\"nisn\":581928,\"jurusan\":\"TI\",\"nama\":\"Ade Reza Fadilah\",\"alamat\":\"Jalan. Gongseng Raya No 6\",\"tgl_lahir\":\"2015-02-02\",\"file_foto\":\"ade.jpg\"},{\"nisn\":582312,\"jurusan\":\"TI\",\"nama\":\"Muhamad mika\",\"alamat\":\"Jl. sepakat 8, Gg Dumay No 13\",\"tgl_lahir\":\"1991-02-05\",\"file_foto\":\"Koala.jpg\"},{\"nisn\":583621,\"jurusan\":\"TI\",\"nama\":\"Urip Elviana Yusuf\",\"alamat\":\"Jln. Protokol Ip Server, Router Cisco\",\"tgl_lahir\":\"1990-01-01\",\"file_foto\":\"ana.jpg\"},{\"nisn\":587173,\"jurusan\":\"TI\",\"nama\":\"Risky Akbar\",\"alamat\":\"Jl. Gongseng Raya No.20 , Cilangkap\",\"tgl_lahir\":\"1997-02-04\",\"file_foto\":\"risky.jpg\"},{\"nisn\":588120,\"jurusan\":\"TI\",\"nama\":\"Adesta Gumay\",\"alamat\":\"Jl. Kupu - Kupu no. 16 Cilangkap\",\"tgl_lahir\":\"1998-01-05\",\"file_foto\":\"eta.jpg\"},{\"nisn\":588332,\"jurusan\":\"TI\",\"nama\":\"Laras Hani Ariyati\",\"alamat\":\"Jl. Pendekar III Taman Mini\",\"tgl_lahir\":\"1997-01-12\",\"file_foto\":\"laras.jpg\"},{\"nisn\":588413,\"jurusan\":\"TI\",\"nama\":\"Rosmawati\",\"alamat\":\"Cibubur Komp. Duta Block D \",\"tgl_lahir\":\"1997-06-13\",\"file_foto\":\"rosma.jpg\"},{\"nisn\":588881,\"jurusan\":\"TI\",\"nama\":\"Firman Teguh Prabowo\",\"alamat\":\"Jln. Gempol Raya no. 16 Cipayung\",\"tgl_lahir\":\"1997-10-21\",\"file_foto\":\"firman.jpg\"},{\"nisn\":598199,\"jurusan\":\"TI\",\"nama\":\"Mas Kukuh Prakoso\",\"alamat\":\"Jl. Kampung Pulo No.13, Jakarta Timur\",\"tgl_lahir\":\"1994-04-09\",\"file_foto\":\"kukuh.jpg\"}]}\n"
	if rec.Body.String() != wantBody {
		t.Errorf("response body not byte-identical to baseline:\ngot:  %q\nwant: %q", rec.Body.String(), wantBody)
	}
}

func TestStudentList_InvalidLimit_Returns500(t *testing.T) {
	// Documented quirk (baseline_responses.json: student_list_invalid_limit):
	// an invalid limit=abc fails gorilla/schema decode BEFORE the strconv.Atoi
	// fallback logic runs, so it returns 500, NOT a silent fallback to limit=10.
	r := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/student/list?limit=abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500, body=%s", rec.Code, rec.Body.String())
	}
}

func TestStudentList_ByNISN(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/student/list?NISN=581928", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	json.Unmarshal(rec.Body.Bytes(), &body)
	data := body["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("expected exactly 1 row for NISN filter, got %d", len(data))
	}
	row := data[0].(map[string]any)
	if row["nisn"].(float64) != 581928 {
		t.Errorf("nisn = %v, want 581928", row["nisn"])
	}
}

func TestStudentDetail_AlwaysReturnsArray(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/student/detail?NISN=581928", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	json.Unmarshal(rec.Body.Bytes(), &body)
	data, ok := body["data"].([]any)
	if !ok {
		t.Fatalf("data must be an array even for a single result, got %T", body["data"])
	}
	if len(data) != 1 {
		t.Fatalf("expected 1 row, got %d", len(data))
	}
}

func TestStudentDetail_IgnoresLimitOffset(t *testing.T) {
	r := newTestRouter()

	req1 := httptest.NewRequest(http.MethodGet, "/student/detail?NISN=581928", nil)
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/student/detail?NISN=581928&limit=1&offset=0", nil)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec1.Body.String() != rec2.Body.String() {
		t.Errorf("expected identical response with/without limit&offset (must be ignored):\nwithout=%s\nwith=%s", rec1.Body.String(), rec2.Body.String())
	}
}

func TestStudentDetail_NonNumericNISN_Returns500(t *testing.T) {
	r := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/student/detail?NISN=abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 (preserved quirk, not sanitized to 400), body=%s", rec.Code, rec.Body.String())
	}
}
