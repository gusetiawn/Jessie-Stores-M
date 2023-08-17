package errors

import (
	"net/http"
	"strings"
)

func ErrorHandling(err error) (errorCode int, errorMessage string) {
	switch {
	case strings.Contains(err.Error(), "invalid addressID from user"):
		return http.StatusBadRequest, "user not choose address yet"
	case strings.Contains(err.Error(), "invalid merchantID from user token"):
		return http.StatusBadRequest, "Bad request"
	case strings.Contains(err.Error(), "store not found"):
		return http.StatusNotFound, "Store not found"
	case strings.Contains(err.Error(), "address not found"):
		return http.StatusBadRequest, "user not choose address yet"
	case strings.Contains(err.Error(), "you are the admin of the store"):
		return http.StatusBadRequest, "You are the admin of the store"
	default:
		return http.StatusInternalServerError, "Something went wrong or server is under maintenance. Please contact app support"
	}
}
