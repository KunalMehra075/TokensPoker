package handlers

import (
	"errors"

	"freetokenspoker/internal/apperr"
	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/middleware"

	"github.com/gin-gonic/gin"
)

// identity pulls the authenticated user from the gin context.
type identity struct {
	UserID string
	Email  string
	Name   string
}

func currentUser(c *gin.Context) identity {
	return identity{
		UserID: c.GetString(middleware.CtxUserID),
		Email:  c.GetString(middleware.CtxEmail),
		Name:   c.GetString(middleware.CtxName),
	}
}

// respondErr maps a service error to the consistent JSON envelope.
func respondErr(c *gin.Context, err error) {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		dto.Fail(c, ae.Status, ae.Code, ae.Message)
		return
	}
	dto.Fail(c, 500, dto.CodeInternal, "internal server error")
}
