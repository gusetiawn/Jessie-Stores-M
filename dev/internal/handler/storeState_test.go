package handler

import (
	"bytes"
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

func TestStoreState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService)

	dummyToken := model.UserToken{
		Token:    gopointer.StringPointer("dummy-token"),
		State:    gopointer.StringPointer("dummy-state"),
		ExpireAt: &time.Time{},
	}

	t.Run("With failure to read request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/store-state", &brokenRequestBody{})

		testHandler.StoreState(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Bad Request", resp["error"])
		assert.Equal(t, "Can't decode request", resp["message"])
	})

	t.Run("With failure to decode request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/store-state", bytes.NewBuffer([]byte("INVALID_JSON")))

		testHandler.StoreState(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Bad Request", resp["error"])
		assert.Equal(t, "Can't decode request", resp["message"])
	})

	dummyReq := constant.StoreStateRequest{ReferenceID: "dummy-reference-id"}
	dummyReqByte, _ := json.Marshal(dummyReq)

	t.Run("With service.StoreState returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/store-state", bytes.NewBuffer(dummyReqByte))

		mockedService.
			EXPECT().
			StoreState(c, &dummyToken, dummyReq).
			Return(nil, errors.New("dummy-error"))

		testHandler.StoreState(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Internal Server Error", resp["error"])
		assert.Equal(t, "Something went wrong or server is under maintenance. Please contact app support", resp["message"])
	})

	t.Run("With service.StoreState returns no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/store-state", bytes.NewBuffer(dummyReqByte))

		mockedService.
			EXPECT().
			StoreState(c, &dummyToken, dummyReq).
			Return(&model.Store{ID: 1}, nil)

		testHandler.StoreState(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "success", resp["status"])
		assert.Equal(t, "Store updated", resp["message"])
		assert.NotEmpty(t, resp["data"])
	})
}
