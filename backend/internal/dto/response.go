package dto

import "github.com/gin-gonic/gin"

// APIError is the consistent error envelope from the architecture doc.
type APIError struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

// Common error codes used across handlers.
const (
	CodeValidation   = "VALIDATION_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
)

// Fail writes a structured JSON error and aborts the request.
func Fail(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, APIError{
		Success:   false,
		Message:   message,
		ErrorCode: code,
	})
}

// OK writes a success envelope wrapping arbitrary data.
func OK(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"success": true, "data": data})
}
