package http

import (
	"fmt"
	"net/http"
)

var (
	ErrBadRequest   = HTTPError{http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil}
	ErrUnauthorized = HTTPError{http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil}
	ErrNotFound     = HTTPError{http.StatusNotFound, http.StatusText(http.StatusNotFound), nil}
	ErrInternal     = HTTPError{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil}
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	wrapped error  // wrapped is used internally to share the error context.
}

func (e HTTPError) Error() string {
	return e.Message
}

func (e HTTPError) Wrap(target error) HTTPError {
	return HTTPError{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", e.Message, target.Error()),
		wrapped: target,
	}
}

func (e HTTPError) Unwrap() error {
	return e.wrapped
}
