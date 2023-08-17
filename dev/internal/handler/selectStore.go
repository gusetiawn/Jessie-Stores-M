package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *handler) SelectStore(c *gin.Context) {
	token := c.Keys["token"].(*model.UserToken)
	logger := logrus.
		WithField("func", "handler.SelectStore").
		WithField("token.UserID", token.UserID)

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.WithError(err).Error("failed to read request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "Can't decode request",
		})
		return
	}

	var request constant.SelectStoreRequest
	if err = json.Unmarshal(jsonData, &request); err != nil {
		logger.WithError(err).Error("failed to decode request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "Can't decode request",
		})
		return
	}

	err = h.service.SelectStore(c, token, request)
	if err != nil {
		logger.WithError(err).Error("failed at service.SelectStore")
		code, msg := errors.ErrorHandling(err)
		c.JSON(code, gin.H{
			"error":   http.StatusText(code),
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token.Token,
		"state":      token.State,
		"is_expired": isExpired(token.ExpireAt),
	})
}
