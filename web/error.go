package web

import (
	"errors"
	"net/http"
)

var (
	ErrLoginFailed            = Error{errors.New("Login failed"), http.StatusUnauthorized}
	ErrNeedLogin              = Error{errors.New("Please login"), http.StatusForbidden}
	ErrMalformedAuthorization = Error{errors.New("Malformed Authorization Header"), http.StatusBadRequest}
	ErrSigningMethod          = Error{errors.New("Unexpected signing method"), http.StatusBadRequest}
	ErrMethodNotAllowed       = Error{errors.New("Method not allowed"), http.StatusMethodNotAllowed}
)

type Error struct {
	Err  error
	Code int
}

func (e Error) Error() string {
	return e.Err.Error()
}
