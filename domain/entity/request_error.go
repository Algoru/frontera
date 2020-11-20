package entity

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidAuthCredentials = errors.New("email and/or password are incorrect")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidAuthHeader      = errors.New("invalid authorization header")
)

func GetStatusCodeForError(err error) int {
	statusCode := http.StatusInternalServerError

	switch err {
	case ErrInvalidAuthHeader:
		fallthrough
	case ErrInvalidAuthCredentials:
		statusCode = http.StatusUnauthorized
	case ErrUserNotFound:
		statusCode = http.StatusBadRequest
	}

	return statusCode
}
