package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCORSMiddleware builds the CORS middleware from CORS_ALLOWED_ORIGINS
// (comma-separated). An empty/unset value means permissive mode: all
// origins allowed. See contracts/cors-middleware.md.
func NewCORSMiddleware() gin.HandlerFunc {
	raw := os.Getenv("CORS_ALLOWED_ORIGINS")

	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	if raw == "" {
		cfg.AllowAllOrigins = true
	} else {
		origins := make([]string, 0)
		for _, o := range strings.Split(raw, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins = append(origins, o)
			}
		}
		cfg.AllowOrigins = origins
	}

	return cors.New(cfg)
}
