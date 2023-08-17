package service

import (
	"context"
	"errors"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	mockRepository "git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStoreState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedRepository := mockRepository.NewMockRepository(ctrl)

	testService := New(mockedRepository)

	dummyReq := constant.StoreStateRequest{
		ReferenceID: "xxx",
		Type:        "items",
	}
	dummyToken := model.UserToken{ID: 1}

	t.Run("With repository.GetStoreByRef returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetStoreByRef(context.TODO(), dummyReq.Type, dummyReq.ReferenceID).
			Return(nil, errors.New("dummy-error"))

		_, err := testService.StoreState(context.TODO(), &dummyToken, dummyReq)

		assert.Error(t, err)
	})

	t.Run("With repository.GetStoreByRef returns nil store", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetStoreByRef(context.TODO(), dummyReq.Type, dummyReq.ReferenceID).
			Return(nil, nil)

		_, err := testService.StoreState(context.TODO(), &dummyToken, dummyReq)

		assert.Error(t, err)
	})

	dummyStore := model.Store{ID: 1}

	t.Run("With repository.UpdateTokenStore returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetStoreByRef(context.TODO(), dummyReq.Type, dummyReq.ReferenceID).
			Return(&dummyStore, nil)
		mockedRepository.
			EXPECT().
			UpdateTokenStore(context.TODO(), dummyToken.ID, dummyStore.ID).
			Return(errors.New("dummy-error"))

		_, err := testService.StoreState(context.TODO(), &dummyToken, dummyReq)

		assert.Error(t, err)
	})

	t.Run("With repository.UpdateTokenStore returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			GetStoreByRef(context.TODO(), dummyReq.Type, dummyReq.ReferenceID).
			Return(&dummyStore, nil)
		mockedRepository.
			EXPECT().
			UpdateTokenStore(context.TODO(), dummyToken.ID, dummyStore.ID).
			Return(nil)

		store, err := testService.StoreState(context.TODO(), &dummyToken, dummyReq)

		assert.NotEmpty(t, store)
		assert.NoError(t, err)
	})
}
