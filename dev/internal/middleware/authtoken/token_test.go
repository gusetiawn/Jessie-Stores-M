package authtoken

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nicklaros/gopointer"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedRepo := mock.NewMockRepository(ctrl)
	testAuthTokenMiddleware := New(mockedRepo)

	dummyRequestToken := "dummy-request-token"

	t.Run("With repository.GetUserLastToken returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/some/path", nil)
		c.Request.Header.Add("x-token", dummyRequestToken)

		mockedRepo.
			EXPECT().
			GetUserLastToken(c, dummyRequestToken).
			Return(nil, errors.New("dummy-error"))

		testAuthTokenMiddleware.ValidateToken(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "Internal Server Error", resp["error"])
		assert.Equal(t, "Something went wrong or server is under maintenance. Please contact app support", resp["message"])
	})

	t.Run("With repository.GetUserLastToken returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/some/path", nil)
		c.Request.Header.Add("x-token", dummyRequestToken)

		mockedRepo.
			EXPECT().
			GetUserLastToken(c, dummyRequestToken).
			Return(nil, nil)

		testAuthTokenMiddleware.ValidateToken(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Equal(t, "invalid_token", resp["error"])
		assert.Equal(t, "Invalid token, please use the correct url from whatsapp message", resp["message"])
	})

	dummyLastToken := model.UserToken{
		Token: gopointer.StringPointer("dummy-last-token"),
	}

	t.Run("With repository.GetUserLastToken returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/some/path", nil)
		c.Request.Header.Add("x-token", dummyRequestToken)

		mockedRepo.
			EXPECT().
			GetUserLastToken(c, dummyRequestToken).
			Return(&dummyLastToken, nil)

		testAuthTokenMiddleware.ValidateToken(c)

		respBody, _ := io.ReadAll(w.Result().Body)
		resp := make(map[string]interface{})
		json.Unmarshal(respBody, &resp)

		assert.Nil(t, resp["error"])
	})
}
