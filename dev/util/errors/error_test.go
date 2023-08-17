package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandling(t *testing.T) {
	t.Run("invalid addressID from user", func(t *testing.T) {
		code, msg := ErrorHandling(errors.New("invalid addressID from user"))

		assert.Equal(t, code, http.StatusBadRequest)
		assert.Equal(t, msg, "user not choose address yet")
	})

	t.Run("invalid merchantID from user token", func(t *testing.T) {
		code, msg := ErrorHandling(errors.New("invalid merchantID from user token"))

		assert.Equal(t, code, http.StatusBadRequest)
		assert.Equal(t, msg, "Bad request")
	})

	t.Run("store not found", func(t *testing.T) {
		code, msg := ErrorHandling(errors.New("store not found"))

		assert.Equal(t, code, http.StatusNotFound)
		assert.Equal(t, msg, "Store not found")
	})

	t.Run("address not found", func(t *testing.T) {
		code, msg := ErrorHandling(errors.New("address not found"))

		assert.Equal(t, code, http.StatusBadRequest)
		assert.Equal(t, msg, "user not choose address yet")
	})

	t.Run("you are the admin of the store", func(t *testing.T) {
		code, msg := ErrorHandling(errors.New("you are the admin of the store"))

		assert.Equal(t, code, http.StatusBadRequest)
		assert.Equal(t, msg, "You are the admin of the store")
	})

	t.Run("general error", func(t *testing.T) {
		code, msg := ErrorHandling(errors.New("dummy-error"))

		assert.Equal(t, code, http.StatusInternalServerError)
		assert.Equal(t, msg, "Something went wrong or server is under maintenance. Please contact app support")
	})
}
