package http

import (
	"fmt"
)

var (
	ErrBadRequest          = &HTTPError{400, "bad request", nil}
	ErrUnauthorized        = &HTTPError{401, "unauthorized", nil}
	ErrForbidden           = &HTTPError{403, "forbidden", nil}
	ErrNotFound            = &HTTPError{404, "not found", nil}
	ErrUnprocessableEntity = &HTTPError{422, "validation error", nil}
	ErrInternal            = &HTTPError{500, "internal error", nil}
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// wrapped is used internally to share the error context.
	wrapped error `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Wrap(target error) *HTTPError {
	return &HTTPError{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", e.Message, target.Error()),
		wrapped: target,
	}
}

func (e *HTTPError) Unwrap() error {
	return e.wrapped
}
