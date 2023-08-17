package http

import (
	"git-rbi.jatismobile.com/jatis_chatcommerce/mi-storesapi/internal/middleware/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupRouter(r *gin.Engine, handlers *ServerHandlers) {
	prometheus.RegisterPromMetrics()
	r.Use(prometheus.PrometheusMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/ping", handlers.Ping)
	r.Use(handlers.ValidateToken)
	r.POST("/m/store-state", handlers.StoreState)
	r.GET("/m/v4/nearest-store", handlers.GetNearestStore)
	r.POST("/m/select-store", handlers.SelectStore)
}
