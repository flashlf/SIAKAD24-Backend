package middleware

import (
	"net/http"
	"strings"

	"lumen/go-siakad/api"
	"lumen/go-siakad/internal/auth"

	"github.com/gin-gonic/gin"
)

const claimsContextKey = "auth_claims"

// RequireRole returns a gin.HandlerFunc that only allows requests carrying a
// valid JWT Bearer token whose role claim is one of the given roles. See
// contracts/auth-middleware.md for the exact status/body contract.
func RequireRole(roles ...auth.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if header == "" || !strings.HasPrefix(header, prefix) {
			api.CustomErrorHandler(c.Writer, errUnauthorized, http.StatusUnauthorized)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(header, prefix)
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			api.CustomErrorHandler(c.Writer, errUnauthorized, http.StatusUnauthorized)
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if claims.Role == r {
				allowed = true
				break
			}
		}
		if !allowed {
			api.CustomErrorHandler(c.Writer, errForbidden, http.StatusForbidden)
			c.Abort()
			return
		}

		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

var (
	errUnauthorized = unauthorizedError{}
	errForbidden    = forbiddenError{}
)

type unauthorizedError struct{}

func (unauthorizedError) Error() string { return "Unauthorized" }

type forbiddenError struct{}

func (forbiddenError) Error() string { return "Forbidden" }
