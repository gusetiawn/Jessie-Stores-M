package service

import (
	"context"
	"errors"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository"
	mockRepository "git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/nicklaros/gopointer"
	"github.com/stretchr/testify/assert"
)

func TestGetNearestStores(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedRepository := mockRepository.NewMockRepository(ctrl)

	testService := New(mockedRepository)

	dummyReq := constant.GetNearestStoresRequest{Token: &model.UserToken{}}

	t.Run("With addressID is invalid", func(t *testing.T) {
		_, err := testService.GetNearestStores(context.TODO(), dummyReq)

		assert.Error(t, err)
	})

	dummyReq.Token.AddressID = gopointer.Int64Pointer(1)

	t.Run("With merchantID is invalid", func(t *testing.T) {
		_, err := testService.GetNearestStores(context.TODO(), dummyReq)

		assert.Error(t, err)
	})

	dummyReq.Token.MerchantID = gopointer.Int64Pointer(1)

	t.Run("With categoryType is invalid", func(t *testing.T) {
		stores, err := testService.GetNearestStores(context.TODO(), dummyReq)

		assert.Empty(t, stores)
		assert.NoError(t, err)
	})

	dummyReq.CategoryType = constant.ALLOWED_CATEGORY_TYPE
	dummyReq.Token.UserID = 1
	dummyReq.CategoryID = 1
	dummyReq.Name = "dummy-name"

	dummyInput := repository.GetNearestStoresInput{
		UserID:     dummyReq.Token.UserID,
		MerchantID: *dummyReq.Token.MerchantID,
		CategoryID: dummyReq.CategoryID,
		Name:       dummyReq.Name,
	}

	ctx := context.TODO()
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	t.Run("With repository.GetAddress returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(nil, errors.New("dummy-error"))
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(nil, nil)

		_, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Error(t, err)
	})

	t.Run("With repository.GetAddress returns nil address", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(nil, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(nil, nil)

		_, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Error(t, err)
	})

	dummyAddress := model.Address{LocationLat: 1, LocationLong: 1}

	t.Run("With repository.GetNearestStores returns nil stores", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(nil, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	dummyStores := []model.StoreWithCategory{{Store: model.Store{ID: 1}}}

	t.Run("With repository.GetStoreOpenConfig returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1}).
			Return(nil, errors.New("dummy-error"))

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	t.Run("With closed is true", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed": "true",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	t.Run("With closed is false but service is invalid", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "invalid-json",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	t.Run("With closed is false but service startTime is invalid", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"invalid-time\",\"end\":\"21:00\"},\"1\":{\"start\":\"invalid-time\",\"end\":\"21:00\"},\"2\":{\"start\":\"invalid-time\",\"end\":\"23:59\"},\"3\":{\"start\":\"invalid-time\",\"end\":\"23:59\"},\"4\":{\"start\":\"invalid-time\",\"end\":\"22:00\"},\"5\":{\"start\":\"invalid-time\",\"end\":\"23:59\"},\"6\":{\"start\":\"invalid-time\",\"end\":\"21:00\"}}",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	t.Run("With closed is false but service endTime is invalid", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"08:00\",\"end\":\"invalid-time\"},\"1\":{\"start\":\"06:00\",\"end\":\"invalid-time\"},\"2\":{\"start\":\"08:00\",\"end\":\"invalid-time\"},\"3\":{\"start\":\"08:00\",\"end\":\"invalid-time\"},\"4\":{\"start\":\"00:00\",\"end\":\"invalid-time\"},\"5\":{\"start\":\"08:00\",\"end\":\"invalid-time\"},\"6\":{\"start\":\"00:00\",\"end\":\"invalid-time\"}}",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	t.Run("With closed is false and now is during opening hour but store's location is undefined", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Empty(t, res)
		assert.NoError(t, err)
	})

	dummyReq.Limit = 2
	dummyStores[0].LocationLat = gopointer.Float64Pointer(1)
	dummyStores[0].LocationLong = gopointer.Float64Pointer(1)
	dummyStores = append(dummyStores, model.StoreWithCategory{Store: model.Store{
		ID:           2,
		LocationLat:  gopointer.Float64Pointer(0.5),
		LocationLong: gopointer.Float64Pointer(0.5),
	}})

	t.Run("With closed is false and now is during opening hour but store's location is not undefined", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
			"2": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1, 2}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.NotEmpty(t, res)
		assert.NoError(t, err)
	})

	dummyReq.Limit = 1
	dummyStores[0].MaxRadius = gopointer.Int64Pointer(0)
	dummyStores = append(dummyStores, model.StoreWithCategory{Store: model.Store{
		ID:           3,
		LocationLat:  gopointer.Float64Pointer(0.5),
		LocationLong: gopointer.Float64Pointer(0.5),
	}})

	t.Run("With closed is false and now is during opening hour, store's location is not undefined, but max_radius is zero", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
			"2": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
			"3": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1, 2, 3}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Len(t, res, 1)
		assert.NoError(t, err)
	})

	t.Run("With closed is false and now is never in the opening hour", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"00:00\"},\"1\":{\"start\":\"00:00\",\"end\":\"00:00\"},\"2\":{\"start\":\"00:00\",\"end\":\"00:00\"},\"3\":{\"start\":\"00:00\",\"end\":\"00:00\"},\"4\":{\"start\":\"00:00\",\"end\":\"00:00\"},\"5\":{\"start\":\"00:00\",\"end\":\"00:00\"},\"6\":{\"start\":\"00:00\",\"end\":\"00:00\"}}",
			},
			"2": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
			"3": map[string]string{
				"closed":  "false",
				"service": "{\"0\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"1\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"2\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"3\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"4\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"5\":{\"start\":\"00:00\",\"end\":\"23:59\"},\"6\":{\"start\":\"00:00\",\"end\":\"23:59\"}}",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1, 2, 3}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Len(t, res, 1)
		assert.NoError(t, err)
	})

	t.Run("With closed is false and service is undefined", func(t *testing.T) {
		dummyConfig := model.StoreConfig{
			"1": map[string]string{
				"closed": "false",
			},
		}

		mockedRepository.
			EXPECT().
			GetAddress(newCtx, *dummyReq.Token.AddressID).
			Return(&dummyAddress, nil)
		mockedRepository.
			EXPECT().
			GetNearestStores(newCtx, dummyInput).
			Return(dummyStores, nil)
		mockedRepository.
			EXPECT().
			GetStoreOpenConfig(ctx, []int64{1, 2, 3}).
			Return(dummyConfig, nil)

		res, err := testService.GetNearestStores(ctx, dummyReq)

		assert.Len(t, res, 1)
		assert.NoError(t, err)
	})
}
