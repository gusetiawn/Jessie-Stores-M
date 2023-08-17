package repository

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetStoreOpenConfig(t *testing.T) {
	testRepository, db, mock := newTestRepository(t)
	defer db.Close()

	dummyStoreIDs := []int64{1, 2}
	query, args := buildGetStoreOpenConfigQuery(dummyStoreIDs)

	t.Run("With storeIDs is empty", func(t *testing.T) {
		config, err := testRepository.GetStoreOpenConfig(context.TODO(), nil)

		assert.Empty(t, config)
		assert.NoError(t, err)
	})

	t.Run("With Scan returns error", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2], args[3]).
			WillReturnError(errors.New("dummmy-error"))

		_, err := testRepository.GetStoreOpenConfig(context.TODO(), dummyStoreIDs)

		assert.Error(t, err)
	})

	t.Run("With returned result is nil", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2], args[3]).
			WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(nil))

		config, err := testRepository.GetStoreOpenConfig(context.TODO(), dummyStoreIDs)

		assert.Empty(t, config)
		assert.NoError(t, err)
	})

	t.Run("With returned result is invalid JSON", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2], args[3]).
			WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow("INVALID_JSON"))

		_, err := testRepository.GetStoreOpenConfig(context.TODO(), dummyStoreIDs)

		assert.Error(t, err)
	})

	t.Run("With returned result is valid JSON", func(t *testing.T) {
		mock.
			ExpectQuery(query).
			WithArgs(args[0], args[1], args[2], args[3]).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"result"}).
					AddRow(`{ "1" : { "close" : "false", "service" : "{\"0\":{\"start\":\"08:00\",\"end\":\"21:00\"},\"1\":{\"start\":\"06:00\",\"end\":\"21:00\"},\"2\":{\"start\":\"08:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"08:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"22:00\"},\"5\":{\"start\":\"08:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"21:00\"}}" } }`),
			)

		config, err := testRepository.GetStoreOpenConfig(context.TODO(), dummyStoreIDs)

		openingHours := make(map[string]struct {
			Start string `json:"start"`
			End   string `json:"end"`
		})
		json.Unmarshal([]byte(config["1"]["service"]), &openingHours)

		assert.Equal(t, "08:00", openingHours["0"].Start)
		assert.NoError(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
