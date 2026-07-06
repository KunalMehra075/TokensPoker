package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger emits one structured JSON line per request.
func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("http_request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"durationMs", time.Since(start).Milliseconds(),
			"clientIp", c.ClientIP(),
		)
	}
}

// Recovery converts panics into a structured 500 instead of crashing.
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Error("panic_recovered", "error", recovered, "path", c.Request.URL.Path)
		c.AbortWithStatusJSON(500, gin.H{
			"success": false, "message": "internal server error", "errorCode": "INTERNAL_ERROR",
		})
	})
}
