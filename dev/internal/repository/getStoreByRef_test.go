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

func TestGetStoreByRef(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	dummyType := "dummy-type"
	dummyRef := "dummy-reference"

	t.Run("With Scan returns sql.ErrNoRows", func(t *testing.T) {
		mock.
			ExpectQuery(SELECT_STORE_BY_REF_ITEM_QUERY).
			WithArgs(dummyRef).
			WillReturnError(sql.ErrNoRows)

		store, err := testRepository.GetStoreByRef(context.TODO(), "items", dummyRef)

		assert.Empty(t, store)
		assert.NoError(t, err)
	})

	t.Run("With Scan returns general errors", func(t *testing.T) {
		mock.
			ExpectQuery(SELECT_STORE_BY_REF_CAT_QUERY).
			WithArgs(dummyRef).
			WillReturnError(errors.New("dummy-error"))

		_, err := testRepository.GetStoreByRef(context.TODO(), "categories", dummyRef)

		assert.Error(t, err)
	})

	t.Run("With Scan from SELECT_STORE_BY_REF_STORE_QUERY returns no error", func(t *testing.T) {
		mock.
			ExpectQuery(SELECT_STORE_BY_REF_STORE_QUERY).
			WithArgs(dummyRef).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"id", "merchant_id", "name", "address", "city", "country", "max_radius", "location_lat", "location_long", "pic", "is_active", "zip_code", "category_id", "created_at", "updated_at", "deleted_at"}).
					AddRow(int64(1), int64(1), "name", "address", "city", "country", int64(1), float64(1), float64(2), "pic", true, "zipcode", int64(1), time.Now(), time.Now(), time.Now()),
			)

		store, err := testRepository.GetStoreByRef(context.TODO(), dummyType, dummyRef)

		assert.NotEmpty(t, store)
		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
