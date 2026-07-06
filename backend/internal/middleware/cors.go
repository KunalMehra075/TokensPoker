package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS applies an env-driven origin allowlist. An unset allowlist blocks every
// browser request, matching the playbook's gotcha about empty CORS_ORIGINS.
func CORS(allowed []string) gin.HandlerFunc {
	allow := map[string]struct{}{}
	for _, o := range allowed {
		allow[o] = struct{}{}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if _, ok := allow[origin]; ok && origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "600")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
