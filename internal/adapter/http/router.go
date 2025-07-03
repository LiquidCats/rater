package http

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(logger.SetLogger(logger.WithUTC(true)))

	return router
}
