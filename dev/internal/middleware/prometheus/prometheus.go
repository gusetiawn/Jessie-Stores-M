package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Number of requests.",
	}, []string{"method", "path"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func RegisterPromMetrics() {
	_ = prometheus.Register(totalRequest)
	_ = prometheus.Register(httpDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() != "/metrics" {
			totalRequest.WithLabelValues(c.Request.Method, c.Request.URL.String()).Inc()
			timer := prometheus.NewTimer(httpDuration.WithLabelValues(c.Request.URL.String()))
			timer.ObserveDuration()
		}
		c.Next()
	}
}
