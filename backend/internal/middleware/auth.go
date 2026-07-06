package middleware

import (
	"net/http"
	"strings"

	"freetokenspoker/internal/auth"
	"freetokenspoker/internal/dto"

	"github.com/gin-gonic/gin"
)

// Context keys for values stashed by the auth middleware.
const (
	CtxUserID = "userId"
	CtxEmail  = "email"
	CtxName   = "name"
)

// Auth validates the Bearer token and stashes the identity on the context.
func Auth(jwt *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			dto.Fail(c, http.StatusUnauthorized, dto.CodeUnauthorized, "missing or malformed authorization header")
			return
		}
		claims, err := jwt.Verify(parts[1])
		if err != nil {
			dto.Fail(c, http.StatusUnauthorized, dto.CodeUnauthorized, "invalid or expired token")
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxEmail, claims.Email)
		c.Set(CtxName, claims.Name)
		c.Next()
	}
}
