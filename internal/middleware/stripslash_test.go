package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStripSlashes(t *testing.T) {
	tests := []struct {
		name         string
		requestPath  string
		expectedPath string
	}{
		{"trailing slash trimmed", "/student/list/", "/student/list"},
		{"no trailing slash unchanged", "/student/list", "/student/list"},
		{"root path untouched", "/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotStatus = http.StatusOK

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
			})

			handler := StripSlashes(next)

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if gotPath != tt.expectedPath {
				t.Errorf("path = %q, want %q", gotPath, tt.expectedPath)
			}
			if rec.Code != gotStatus {
				t.Errorf("status = %d, want %d", rec.Code, gotStatus)
			}
			if rec.Header().Get("Location") != "" {
				t.Errorf("unexpected redirect Location header: %q", rec.Header().Get("Location"))
			}
		})
	}
}
