package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTokenStore(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	dummyTokenID := int64(1)
	dummyStoreID := int64(1)

	t.Run("With ExecContext update token store query returns error", func(t *testing.T) {
		mock.
			ExpectExec(UPDATE_STORE_ID_TOKEN_QUERY).
			WithArgs(dummyStoreID, dummyTokenID).
			WillReturnError(errors.New("dummy-error"))

		err := testRepository.UpdateTokenStore(context.TODO(), dummyTokenID, dummyStoreID)

		assert.Error(t, err)
	})

	t.Run("With ExecContext update token store query returns no error", func(t *testing.T) {
		mock.
			ExpectExec(UPDATE_STORE_ID_TOKEN_QUERY).
			WithArgs(dummyStoreID, dummyTokenID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := testRepository.UpdateTokenStore(context.TODO(), dummyTokenID, dummyStoreID)

		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
