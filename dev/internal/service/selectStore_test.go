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

func TestSelectStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedRepository := mockRepository.NewMockRepository(ctrl)

	testService := New(mockedRepository)

	dummyReq := constant.SelectStoreRequest{ID: 9}
	dummyToken := model.UserToken{UserID: 1}

	t.Run("With repository.CheckUserIsStoreAdmin returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			CheckUserIsStoreAdmin(context.TODO(), dummyToken.UserID, dummyReq.ID).
			Return(false, errors.New("dummy-error"))

		err := testService.SelectStore(context.TODO(), &dummyToken, dummyReq)

		assert.Error(t, err)
	})

	t.Run("With repository.CheckUserIsStoreAdmin returns isAdmin true", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			CheckUserIsStoreAdmin(context.TODO(), dummyToken.UserID, dummyReq.ID).
			Return(true, nil)

		err := testService.SelectStore(context.TODO(), &dummyToken, dummyReq)

		assert.Error(t, err)
	})

	t.Run("With repository.UpdateTokenStore returns error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			CheckUserIsStoreAdmin(context.TODO(), dummyToken.UserID, dummyReq.ID).
			Return(false, nil)
		mockedRepository.
			EXPECT().
			UpdateTokenStore(context.TODO(), dummyToken.ID, dummyReq.ID).
			Return(errors.New("dummy-error"))

		err := testService.SelectStore(context.TODO(), &dummyToken, dummyReq)

		assert.Error(t, err)
	})

	t.Run("With repository.UpdateTokenStore returns no error", func(t *testing.T) {
		mockedRepository.
			EXPECT().
			CheckUserIsStoreAdmin(context.TODO(), dummyToken.UserID, dummyReq.ID).
			Return(false, nil)
		mockedRepository.
			EXPECT().
			UpdateTokenStore(context.TODO(), dummyToken.ID, dummyReq.ID).
			Return(nil)

		err := testService.SelectStore(context.TODO(), &dummyToken, dummyReq)

		assert.NoError(t, err)
	})
}
