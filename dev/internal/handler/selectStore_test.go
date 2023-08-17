package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nicklaros/gopointer"
	"github.com/stretchr/testify/assert"
)

func TestSelectStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService)

	dummyToken := model.UserToken{
		Token: gopointer.StringPointer("dummy-token"),
		State: gopointer.StringPointer("dummy-state"),
	}

	t.Run("With failure to read request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/select-store", &brokenRequestBody{})

		testHandler.SelectStore(c)

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
		c.Request = httptest.NewRequest(http.MethodPost, "/m/select-store", bytes.NewBuffer([]byte("INVALID_JSON")))

		testHandler.SelectStore(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Bad Request", resp["error"])
		assert.Equal(t, "Can't decode request", resp["message"])
	})

	dummyReq := constant.SelectStoreRequest{ID: 1}
	dummyReqByte, _ := json.Marshal(dummyReq)

	t.Run("With service.SelectStore returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/select-store", bytes.NewBuffer(dummyReqByte))

		mockedService.
			EXPECT().
			SelectStore(c, &dummyToken, dummyReq).
			Return(errors.New("dummy-error"))

		testHandler.SelectStore(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Internal Server Error", resp["error"])
		assert.Equal(t, "Something went wrong or server is under maintenance. Please contact app support", resp["message"])
	})

	t.Run("With service.SelectStore returns no error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("token", &dummyToken)
		c.Request = httptest.NewRequest(http.MethodPost, "/m/select-store", bytes.NewBuffer(dummyReqByte))

		mockedService.
			EXPECT().
			SelectStore(c, &dummyToken, dummyReq).
			Return(nil)

		testHandler.SelectStore(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, *dummyToken.Token, resp["token"].(string))
		assert.Equal(t, *dummyToken.State, resp["state"].(string))
		assert.Equal(t, isExpired(dummyToken.ExpireAt), resp["is_expired"])
	})
}

type brokenRequestBody struct{}

func (ths *brokenRequestBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("dummy-error")
}

func (ths *brokenRequestBody) Close() error {
	return errors.New("dummy-error")
}
