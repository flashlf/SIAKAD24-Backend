package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"lumen/go-siakad/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	os.Setenv("JWT_SECRET", "test-secret-for-auth-middleware")
}

func newAuthTestRouter(roles ...auth.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", RequireRole(roles...), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func signToken(t *testing.T, role auth.Role, expiresAt time.Time) string {
	t.Helper()
	claims := auth.Claims{
		Subject: "test-user",
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func TestRequireRole_NoAuthorizationHeader(t *testing.T) {
	r := newAuthTestRouter(auth.RoleAdministrator)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestRequireRole_MalformedHeader(t *testing.T) {
	r := newAuthTestRouter(auth.RoleAdministrator)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "NotBearer sometoken")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestRequireRole_InvalidSignature(t *testing.T) {
	r := newAuthTestRouter(auth.RoleAdministrator)
	badToken := signTokenWithSecret(t, auth.RoleAdministrator, time.Now().Add(time.Hour), "wrong-secret")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+badToken)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestRequireRole_ExpiredToken(t *testing.T) {
	r := newAuthTestRouter(auth.RoleAdministrator)
	token := signToken(t, auth.RoleAdministrator, time.Now().Add(-time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestRequireRole_RoleNotAllowed(t *testing.T) {
	r := newAuthTestRouter(auth.RoleAdministrator)
	token := signToken(t, auth.RoleStaff, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}

func TestRequireRole_ValidTokenAllowed(t *testing.T) {
	r := newAuthTestRouter(auth.RoleAdministrator, auth.RoleStaff)
	token := signToken(t, auth.RoleStaff, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body=%s", rec.Code, rec.Body.String())
	}
}

func signTokenWithSecret(t *testing.T, role auth.Role, expiresAt time.Time, secret string) string {
	t.Helper()
	claims := auth.Claims{
		Subject: "test-user",
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}
