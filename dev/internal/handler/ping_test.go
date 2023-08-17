package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedService := mock.NewMockService(ctrl)
	testHandler := New(mockedService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ping", nil)

	testHandler.Ping(c)

	respBody, _ := io.ReadAll(w.Result().Body)
	resp := make(map[string]interface{})
	json.Unmarshal(respBody, &resp)

	assert.Equal(t, "ok", resp["message"])
}
