package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	logrustest "github.com/sirupsen/logrus/hooks/test"
)

func newAccessLogTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AccessLog())
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.GET("/client-error", func(c *gin.Context) { c.Status(http.StatusBadRequest) })
	r.GET("/server-error", func(c *gin.Context) { c.Status(http.StatusInternalServerError) })
	return r
}

func TestAccessLog_FieldsAndSeverity(t *testing.T) {
	hook := logrustest.NewLocal(log.StandardLogger())

	tests := []struct {
		path     string
		wantCode int
		wantLvl  log.Level
	}{
		{"/ok", http.StatusOK, log.InfoLevel},
		{"/client-error", http.StatusBadRequest, log.WarnLevel},
		{"/server-error", http.StatusInternalServerError, log.ErrorLevel},
	}

	r := newAccessLogTestRouter()

	for _, tt := range tests {
		hook.Reset()
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != tt.wantCode {
			t.Fatalf("%s: status = %d, want %d", tt.path, rec.Code, tt.wantCode)
		}

		entries := hook.AllEntries()
		if len(entries) != 1 {
			t.Fatalf("%s: expected exactly 1 log entry, got %d", tt.path, len(entries))
		}
		entry := entries[0]

		if entry.Level != tt.wantLvl {
			t.Errorf("%s: level = %v, want %v", tt.path, entry.Level, tt.wantLvl)
		}
		for _, field := range []string{"method", "path", "status", "latency_ms"} {
			if _, ok := entry.Data[field]; !ok {
				t.Errorf("%s: missing field %q in log entry", tt.path, field)
			}
		}
		if got := entry.Data["status"]; got != tt.wantCode {
			t.Errorf("%s: status field = %v, want %v", tt.path, got, tt.wantCode)
		}
		if got := entry.Data["method"]; got != http.MethodGet {
			t.Errorf("%s: method field = %v, want %v", tt.path, got, http.MethodGet)
		}
		if got := entry.Data["path"]; got != tt.path {
			t.Errorf("%s: path field = %v, want %v", tt.path, got, tt.path)
		}
	}
}
