package handler

import (
	"encoding/json"
	"fmt"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (h *handler) StoreState(c *gin.Context) {
	token := c.Keys["token"].(*model.UserToken)
	logger := logrus.
		WithField("func", "handler.StoreState").
		WithField("token.UserID", token.UserID).
		WithField("token.StoreID", &token.StoreID)

	var a int64
	a = *token.StoreID
	fmt.Println("a")
	fmt.Println(a)

	jsonData, err := io.ReadAll(c.Request.Body)
	logger.Infoln("payload: ", string(jsonData))
	if err != nil {
		logger.WithError(err).Error("failed to read body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "Can't decode request",
		})
		return
	}

	var request constant.StoreStateRequest
	if err = json.Unmarshal(jsonData, &request); err != nil {
		logger.WithError(err).Error("failed to decode request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "Can't decode request",
		})
		return
	}

	store, err := h.service.StoreState(c, token, request)
	if err != nil {
		logger.WithError(err).Error("failed at service.StoreState")
		code, msg := errors.ErrorHandling(err)
		c.JSON(code, gin.H{
			"error":   http.StatusText(code),
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Store updated",
		"data": gin.H{
			"store": store,
		},
	})
}
