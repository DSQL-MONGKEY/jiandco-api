package apperror

import (
	"errors"
	"net/http"
)

type Error struct {
	Status  int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}

func Wrap(status int, message string, err error) *Error {
	return &Error{Status: status, Message: message, Err: err}
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *Error {
	return New(http.StatusForbidden, message)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, message)
}

func Conflict(message string) *Error {
	return New(http.StatusConflict, message)
}

func Internal(err error) *Error {
	return Wrap(http.StatusInternalServerError, "internal server error", err)
}

func Is(err error, target *Error) bool {
	return errors.Is(err, target)
}
