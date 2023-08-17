package authtoken

import (
	"net/http"

	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthTokenMiddleware struct {
	repository repository.Repository
}

func New(repository repository.Repository) *AuthTokenMiddleware {
	return &AuthTokenMiddleware{repository: repository}
}

func (at *AuthTokenMiddleware) ValidateToken(c *gin.Context) {
	reqToken := c.GetHeader("x-token")

	logger := logrus.
		WithField("func", "middleware.microsite_validatetoken").
		WithField("token", reqToken)

	lastToken, err := at.repository.GetUserLastToken(c, reqToken)
	if err != nil {
		logger.WithError(err).Error("failed to get token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "Something went wrong or server is under maintenance. Please contact app support",
		})
		c.Abort()
		return
	}

	if lastToken == nil {
		logger.WithError(err).Error("token does not exist")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_token",
			"message": "Invalid token, please use the correct url from whatsapp message",
		})
		c.Abort()
		return
	}

	if lastToken.Token != nil && *lastToken.Token != reqToken {
		logger.WithField("latest_token", lastToken.Token).Info("token that is used is old, replacing with the latest token")
	}

	c.Set("token", lastToken)
}
