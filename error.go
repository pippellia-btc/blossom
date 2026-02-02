package blossom

import (
	"fmt"
	"net/http"
)

// Error represent an HTTP error with the specified code and reason.
// If the reason is not empty, it is written in the "X-Reason" header as per BUD-01.
type Error struct {
	Code   int
	Reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, reason: %s", e.Code, e.Reason)
}

func (e Error) String() string {
	return e.Error()
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	err, ok := target.(Error)
	if !ok {
		return false
	}
	return e.Code == err.Code && e.Reason == err.Reason
}

// WriteError writes the error to the http response. If the reason is non-empty,
// it writes it to the "X-Reason" header as per BUD-01.
func WriteError(w http.ResponseWriter, e *Error) {
	if e == nil {
		return
	}
	if e.Reason != "" {
		w.Header().Set("X-Reason", e.Reason)
	}
	http.Error(w, "", e.Code)
}

// ErrBadRequest returns a 400 Bad Request error with the given reason.
func ErrBadRequest(reason string) *Error {
	return &Error{Code: http.StatusBadRequest, Reason: reason}
}

// ErrUnauthorized returns a 401 Unauthorized error with the given reason.
func ErrUnauthorized(reason string) *Error {
	return &Error{Code: http.StatusUnauthorized, Reason: reason}
}

// ErrPaymentRequired returns a 402 Payment Required error with the given reason.
func ErrPaymentRequired(reason string) *Error {
	return &Error{Code: http.StatusPaymentRequired, Reason: reason}
}

// ErrForbidden returns a 403 Forbidden error with the given reason.
func ErrForbidden(reason string) *Error {
	return &Error{Code: http.StatusForbidden, Reason: reason}
}

// ErrNotFound returns a 404 Not Found error with the given reason.
func ErrNotFound(reason string) *Error {
	return &Error{Code: http.StatusNotFound, Reason: reason}
}

// ErrNotAllowed returns a 405 Method Not Allowed error with the given reason.
func ErrNotAllowed(reason string) *Error {
	return &Error{Code: http.StatusMethodNotAllowed, Reason: reason}
}

// ErrTooLarge returns a 413 Payload Too Large error with the given reason.
func ErrTooLarge(reason string) *Error {
	return &Error{Code: http.StatusRequestEntityTooLarge, Reason: reason}
}

// ErrUnsupportedMedia returns a 415 Unsupported Media Type error with the given reason.
func ErrUnsupportedMedia(reason string) *Error {
	return &Error{Code: http.StatusUnsupportedMediaType, Reason: reason}
}

// ErrTooMany returns a 429 Too Many Requests error with the given reason.
func ErrTooMany(reason string) *Error {
	return &Error{Code: http.StatusTooManyRequests, Reason: reason}
}

// ErrInternal returns a 500 Internal Server Error with the given reason.
func ErrInternal(reason string) *Error {
	return &Error{Code: http.StatusInternalServerError, Reason: reason}
}

// ErrNotImplemented returns a 501 Not Implemented error with the given reason.
func ErrNotImplemented(reason string) *Error {
	return &Error{Code: http.StatusNotImplemented, Reason: reason}
}

// ErrUnavailable returns a 503 Service Unavailable error with the given reason.
func ErrUnavailable(reason string) *Error {
	return &Error{Code: http.StatusServiceUnavailable, Reason: reason}
}
