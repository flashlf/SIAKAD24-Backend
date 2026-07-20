package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"lumen/go-siakad/internal/auth"

	"github.com/gin-gonic/gin"
)

func newCORSTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(NewCORSMiddleware())
	r.GET("/open", func(c *gin.Context) { c.Status(http.StatusOK) })

	protected := r.Group("/protected")
	protected.Use(RequireRole(auth.RoleAdministrator))
	protected.GET("/thing", func(c *gin.Context) { c.Status(http.StatusOK) })

	return r
}

func TestCORS_PermissiveByDefault(t *testing.T) {
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	r := newCORSTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/open", nil)
	req.Header.Set("Origin", "https://anywhere.example.com")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	got := rec.Header().Get("Access-Control-Allow-Origin")
	if got != "https://anywhere.example.com" && got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want permissive match", got)
	}
}

func TestCORS_AllowedOriginRestriction(t *testing.T) {
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://allowed.example.com")
	defer os.Unsetenv("CORS_ALLOWED_ORIGINS")
	r := newCORSTestRouter()

	// Allowed origin
	req1 := httptest.NewRequest(http.MethodGet, "/open", nil)
	req1.Header.Set("Origin", "https://allowed.example.com")
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)
	if got := rec1.Header().Get("Access-Control-Allow-Origin"); got != "https://allowed.example.com" {
		t.Errorf("allowed origin: Access-Control-Allow-Origin = %q, want %q", got, "https://allowed.example.com")
	}

	// Disallowed origin
	req2 := httptest.NewRequest(http.MethodGet, "/open", nil)
	req2.Header.Set("Origin", "https://not-allowed.example.com")
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if got := rec2.Header().Get("Access-Control-Allow-Origin"); got == "https://not-allowed.example.com" {
		t.Errorf("disallowed origin must not get a matching Access-Control-Allow-Origin, got %q", got)
	}
}

func TestCORS_PreflightBypassesAuth(t *testing.T) {
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	r := newCORSTestRouter()

	req := httptest.NewRequest(http.MethodOptions, "/protected/thing", nil)
	req.Header.Set("Origin", "https://anywhere.example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code == http.StatusUnauthorized || rec.Code == http.StatusForbidden {
		t.Fatalf("preflight OPTIONS must not require auth, got status %d", rec.Code)
	}
}
