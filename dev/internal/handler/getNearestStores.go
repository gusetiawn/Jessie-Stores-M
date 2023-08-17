package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/constant"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/model"
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/util/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *handler) GetNearestStore(c *gin.Context) {
	token := c.Keys["token"].(*model.UserToken)

	name := c.Query("name")
	if name != "" {
		name = strings.ReplaceAll(name, "'", "''")
		name = strings.ToLower(name)
	}

	categoryType := c.Query("category_type")
	if categoryType == "" {
		categoryType = constant.ALLOWED_CATEGORY_TYPE
	}

	logger := logrus.
		WithField("func", "handler.GetNearestStore").
		WithField("name", name).
		WithField("categoryType", categoryType)

	categoryID, err := strconv.Atoi(c.Query("category_id"))
	if err != nil {
		logger.Warn("invalid category_id query, fallback to 0")
	}
	logger.WithField("categoryId", categoryID)

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit == 0 {
		logger.Warn("invalid limit query, fallback to DEFAULT_LIMIT (5)")
		limit = constant.DEFAULT_LIMIT
	}
	logger.WithField("limit", limit)

	req := constant.GetNearestStoresRequest{
		Token:        token,
		Name:         name,
		CategoryType: categoryType,
		CategoryID:   categoryID,
		Limit:        limit,
	}
	res, err := h.service.GetNearestStores(c, req)
	if err != nil {
		logger.WithError(err).Error("failed at service.GetNearestStores")
		code, msg := errors.ErrorHandling(err)
		c.JSON(code, gin.H{
			"error":   http.StatusText(code),
			"message": msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       res,
		"token":      token.Token,
		"state":      token.State,
		"is_expired": isExpired(token.ExpireAt),
	})
}

func isExpired(expiredAt *time.Time) bool {
	if expiredAt != nil {
		return expiredAt.Before(time.Now().UTC())
	}

	return false
}
