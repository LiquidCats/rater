package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GinHandler() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()

	mux.Any("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	return mux
}
