package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"rater/internal/port"
)

func NewRouter(log port.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Info("server: request",
			zap.Int("method", param.StatusCode),
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.String("client_ip", param.ClientIP),
			zap.Duration("latency", param.Latency),
		)

		return ""
	}))

	return router
}
