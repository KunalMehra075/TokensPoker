// Package apperr provides typed application errors that carry an HTTP status
// and a stable error code, keeping handlers thin.
package apperr

import "net/http"

// Error is a domain error with HTTP mapping.
type Error struct {
	Status  int
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }

// New builds an Error.
func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

// Helpers for common cases.
func BadRequest(msg string) *Error   { return New(http.StatusBadRequest, "VALIDATION_ERROR", msg) }
func Unauthorized(msg string) *Error { return New(http.StatusUnauthorized, "UNAUTHORIZED", msg) }
func Forbidden(msg string) *Error    { return New(http.StatusForbidden, "FORBIDDEN", msg) }
func NotFound(msg string) *Error     { return New(http.StatusNotFound, "NOT_FOUND", msg) }
func Conflict(msg string) *Error     { return New(http.StatusConflict, "CONFLICT", msg) }
func Internal(msg string) *Error     { return New(http.StatusInternalServerError, "INTERNAL_ERROR", msg) }
