package util

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

type DomainError struct {
	Message    string
	HttpStatus int
	Code       int
	WrappedErr error
}

func newDomainError(code int, httpStatus int, defaultMessage string) func(err error, message string) DomainError {
	f := func(err error, message string) DomainError {
		var stackInfo, callerInfo string
		if _, file, line, ok := runtime.Caller(1); ok {
			stackInfo = fmt.Sprintf(" source=\"%s:%d\"", file, line)
		}
		if _, file, line, ok := runtime.Caller(2); ok {
			callerInfo = fmt.Sprintf(" caller=\"%s:%d\"", file, line)
		}
		if message == "" {
			message = defaultMessage
		}
		return DomainError{
			Message:    message,
			HttpStatus: httpStatus,
			Code:       code,
			WrappedErr: fmt.Errorf("domain error%s%s code=%d message=\"%s\": %w", stackInfo, callerInfo, code, message, err),
		}
	}
	return f
}

func (e DomainError) Error() string {
	return e.WrappedErr.Error()
}

func (e DomainError) Unwrap() error {
	return errors.Unwrap(e.WrappedErr)
}

var (
	ErrBadRequest          = newDomainError(40000, http.StatusBadRequest, "bad request")
	ErrUnauthorized        = newDomainError(40100, http.StatusUnauthorized, "unauthorized")
	ErrForbidden           = newDomainError(40300, http.StatusForbidden, "forbidden")
	ErrNotFound            = newDomainError(40400, http.StatusNotFound, "not found")
	ErrInternalServerError = newDomainError(50000, http.StatusInternalServerError, "internal server error")
	ErrNotImplemented      = newDomainError(50100, http.StatusNotImplemented, "not implemented")
)
