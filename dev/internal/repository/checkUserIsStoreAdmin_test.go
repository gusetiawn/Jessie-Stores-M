package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCheckUserIsStoreAdmin(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	t.Run("With Scan returns sql.ErrNoRows", func(t *testing.T) {
		mock.
			ExpectQuery(CHECK_USER_IS_STORE_ADMIN_QUERY).
			WithArgs(int64(1), int64(1)).
			WillReturnError(sql.ErrNoRows)

		isAdmin, err := testRepository.CheckUserIsStoreAdmin(context.TODO(), 1, 1)

		assert.False(t, isAdmin)
		assert.NoError(t, err)
	})

	t.Run("With Scan returns generic error", func(t *testing.T) {
		mock.
			ExpectQuery(CHECK_USER_IS_STORE_ADMIN_QUERY).
			WithArgs(int64(1), int64(1)).
			WillReturnError(errors.New("dummy-error"))

		_, err := testRepository.CheckUserIsStoreAdmin(context.TODO(), 1, 1)

		assert.Error(t, err)
	})

	t.Run("With Scan returns no error", func(t *testing.T) {
		mock.
			ExpectQuery(CHECK_USER_IS_STORE_ADMIN_QUERY).
			WithArgs(int64(1), int64(1)).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"is_admin"}).
					AddRow(true),
			)

		isAdmin, err := testRepository.CheckUserIsStoreAdmin(context.TODO(), 1, 1)

		assert.True(t, isAdmin)
		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
