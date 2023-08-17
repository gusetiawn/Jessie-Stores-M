package handler

import (
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Ping(c *gin.Context)
	StoreState(c *gin.Context)
	GetNearestStore(c *gin.Context)
	SelectStore(c *gin.Context)
}

type handler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handler{
		service: service,
	}
}
