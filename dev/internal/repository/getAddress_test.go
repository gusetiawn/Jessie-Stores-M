package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAddress(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	dummyAddress := model.Address{LocationLat: 1, LocationLong: 1}

	t.Run("With Scan returns sql.ErrNoRows", func(t *testing.T) {
		mock.
			ExpectQuery(GET_ADDRESS_QUERY).
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		address, err := testRepository.GetAddress(context.TODO(), 1)

		assert.Empty(t, address)
		assert.NoError(t, err)
	})

	t.Run("With Scan returns generic error", func(t *testing.T) {
		mock.
			ExpectQuery(GET_ADDRESS_QUERY).
			WithArgs(1).
			WillReturnError(errors.New("dummy-error"))

		_, err := testRepository.GetAddress(context.TODO(), 1)

		assert.Error(t, err)
	})

	t.Run("With Scan returns no error", func(t *testing.T) {
		mock.
			ExpectQuery(GET_ADDRESS_QUERY).
			WithArgs(1).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"location_lat", "location_long"}).
					AddRow(dummyAddress.LocationLat, dummyAddress.LocationLong),
			)

		address, err := testRepository.GetAddress(context.TODO(), 1)

		assert.NotEmpty(t, address)
		assert.NoError(t, err)
	})
}
