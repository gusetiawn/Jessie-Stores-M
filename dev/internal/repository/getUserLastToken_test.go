package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserLastToken(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	dummyTokenString := "dummy-token-string"

	t.Run("With Scan returns generic error", func(t *testing.T) {
		mock.
			ExpectQuery(GET_USER_LAST_TOKEN_QUERY).
			WithArgs(dummyTokenString).
			WillReturnError(errors.New("dummy-error"))

		_, err := testRepository.GetUserLastToken(context.TODO(), dummyTokenString)

		assert.Error(t, err)
	})

	t.Run("With Scan returns sql.ErrNoRows", func(t *testing.T) {
		mock.
			ExpectQuery(GET_USER_LAST_TOKEN_QUERY).
			WithArgs(dummyTokenString).
			WillReturnError(sql.ErrNoRows)

		token, err := testRepository.GetUserLastToken(context.TODO(), dummyTokenString)

		assert.Nil(t, token)
		assert.NoError(t, err)
	})

	t.Run("With Scan returns no error", func(t *testing.T) {
		mock.
			ExpectQuery(GET_USER_LAST_TOKEN_QUERY).
			WithArgs(dummyTokenString).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"id", "token", "merchant_id", "user_id", "address_id", "state", "store_id", "expire_at", "created_at", "updated_at", "deleted_at"}).
					AddRow(int64(1), "dummy-token", int64(1), int64(1), int64(1), "dummy-state", int64(1), time.Now(), time.Now(), time.Now(), nil),
			)

		token, err := testRepository.GetUserLastToken(context.TODO(), dummyTokenString)

		assert.NotNil(t, token)
		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
