package apperr

import (
	"net/http"
	"testing"
)

func TestHelpersMapToStatusAndCode(t *testing.T) {
	cases := []struct {
		name       string
		err        *Error
		wantStatus int
		wantCode   string
	}{
		{"bad request", BadRequest("bad"), http.StatusBadRequest, "VALIDATION_ERROR"},
		{"unauthorized", Unauthorized("no"), http.StatusUnauthorized, "UNAUTHORIZED"},
		{"forbidden", Forbidden("nope"), http.StatusForbidden, "FORBIDDEN"},
		{"not found", NotFound("gone"), http.StatusNotFound, "NOT_FOUND"},
		{"conflict", Conflict("dupe"), http.StatusConflict, "CONFLICT"},
		{"internal", Internal("boom"), http.StatusInternalServerError, "INTERNAL_ERROR"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.err.Status != c.wantStatus {
				t.Errorf("Status = %d, want %d", c.err.Status, c.wantStatus)
			}
			if c.err.Code != c.wantCode {
				t.Errorf("Code = %q, want %q", c.err.Code, c.wantCode)
			}
		})
	}
}

func TestErrorMessageIsPreserved(t *testing.T) {
	e := BadRequest("something is wrong")
	if e.Error() != "something is wrong" {
		t.Errorf("Error() = %q, want %q", e.Error(), "something is wrong")
	}
}

func TestNewBuildsError(t *testing.T) {
	e := New(http.StatusTeapot, "TEAPOT", "short and stout")
	if e.Status != http.StatusTeapot || e.Code != "TEAPOT" || e.Message != "short and stout" {
		t.Errorf("New did not set fields: %+v", e)
	}
	// It satisfies the standard error interface.
	var _ error = e
}
