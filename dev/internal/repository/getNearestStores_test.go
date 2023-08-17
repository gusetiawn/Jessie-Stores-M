package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nicklaros/gopointer"
	"github.com/stretchr/testify/assert"
)

func newTestRepository(t *testing.T) (Repository, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed opening a stub database connection: %s", err.Error())
	}

	return New(db), db, mock
}

func TestGetNearestStores(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	dummyInput := GetNearestStoresInput{
		UserID:     1,
		MerchantID: 1,
	}

	query, args := buildGetStoresQuery(dummyInput)

	t.Run("With QueryContext returns error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1]).
			WillReturnError(errors.New("dummy-error"))

		_, err := testRepository.GetNearestStores(context.TODO(), dummyInput)

		assert.Error(t, err)
	})

	dummyInput.CategoryID = 1
	query, args = buildGetStoresQuery(dummyInput)

	t.Run("With QueryContext returns empty rows", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2]).
			WillReturnRows(&sqlmock.Rows{})

		result, err := testRepository.GetNearestStores(context.TODO(), dummyInput)

		assert.Empty(t, result)
		assert.NoError(t, err)
	})

	dummyInput.Name = "dummy-name"
	query, args = buildGetStoresQuery(dummyInput)

	dummyNearestStores := []model.StoreWithCategory{
		{
			Store: model.Store{
				ID:           1,
				Name:         "dummy-name",
				Address:      gopointer.StringPointer("dummy-address"),
				City:         gopointer.StringPointer("dummy-city"),
				Country:      gopointer.StringPointer("dummy-country"),
				MaxRadius:    gopointer.Int64Pointer(1),
				Pic:          gopointer.StringPointer("dummy-pic"),
				LocationLat:  gopointer.Float64Pointer(1),
				LocationLong: gopointer.Float64Pointer(1),
				MerchantID:   gopointer.Int64Pointer(1),
				ZipCode:      gopointer.StringPointer("1234"),
				CategoryID:   gopointer.Int64Pointer(1),
			},
			CategoryName: "dummy-category",
		},
	}

	t.Run("With Scan returns error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2], args[3]).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"id", "name", "address", "city", "country", "max_radius", "pic", "location_lat", "location_long", "merchant_id", "zip_code", "category_id", "category_name"}).
					AddRow("invalid-store-id", dummyNearestStores[0].Store.Name, dummyNearestStores[0].Store.Address, dummyNearestStores[0].Store.City, dummyNearestStores[0].Store.Country, dummyNearestStores[0].Store.MaxRadius, dummyNearestStores[0].Store.Pic, dummyNearestStores[0].Store.LocationLat, dummyNearestStores[0].Store.LocationLong, dummyNearestStores[0].Store.MerchantID, dummyNearestStores[0].Store.ZipCode, dummyNearestStores[0].Store.CategoryID, dummyNearestStores[0].CategoryName),
			)

		result, err := testRepository.GetNearestStores(context.TODO(), dummyInput)

		assert.Empty(t, result)
		assert.Error(t, err)
	})

	t.Run("With Scan returns no error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2], args[3]).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"id", "name", "address", "city", "country", "max_radius", "pic", "location_lat", "location_long", "merchant_id", "zip_code", "category_id", "category_name"}).
					AddRow(dummyNearestStores[0].Store.ID, dummyNearestStores[0].Store.Name, dummyNearestStores[0].Store.Address, dummyNearestStores[0].Store.City, dummyNearestStores[0].Store.Country, dummyNearestStores[0].Store.MaxRadius, dummyNearestStores[0].Store.Pic, dummyNearestStores[0].Store.LocationLat, dummyNearestStores[0].Store.LocationLong, dummyNearestStores[0].Store.MerchantID, dummyNearestStores[0].Store.ZipCode, dummyNearestStores[0].Store.CategoryID, dummyNearestStores[0].CategoryName),
			)

		result, err := testRepository.GetNearestStores(context.TODO(), dummyInput)

		assert.NotEmpty(t, result)
		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
