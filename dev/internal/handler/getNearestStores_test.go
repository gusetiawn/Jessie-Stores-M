package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nicklaros/gopointer"
	"github.com/stretchr/testify/assert"
)

func TestGetNearestStores(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService)

	dummyToken := model.UserToken{
		Token:    gopointer.StringPointer("dummy-token"),
		State:    gopointer.StringPointer("dummy-state"),
		ExpireAt: &time.Time{},
	}
	dummyReq := constant.GetNearestStoresRequest{
		Token:        &dummyToken,
		Name:         "dummy-name",
		CategoryType: constant.ALLOWED_CATEGORY_TYPE,
		Limit:        constant.DEFAULT_LIMIT,
	}

	t.Run("With service.GetNearestStores returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodGet, "/m/v4/nearest-store?name=dummy-name", nil)

		mockedService.
			EXPECT().
			GetNearestStores(c, dummyReq).
			Return(nil, errors.New("dummy-error"))

		testHandler.GetNearestStore(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Internal Server Error", resp["error"])
		assert.Equal(t, "Something went wrong or server is under maintenance. Please contact app support", resp["message"])
	})

	t.Run("With service.GetNearestStores returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodGet, "/m/v4/nearest-store?name=dummy-name", nil)

		mockedService.
			EXPECT().
			GetNearestStores(c, dummyReq).
			Return([]constant.NearestStore{}, nil)

		testHandler.GetNearestStore(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.NotNil(t, resp["data"])
		assert.Equal(t, *dummyToken.Token, resp["token"].(string))
		assert.Equal(t, *dummyToken.State, resp["state"].(string))
		assert.Equal(t, isExpired(dummyToken.ExpireAt), resp["is_expired"])
	})
}
