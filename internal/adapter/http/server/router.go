package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewRouter(logger *zerolog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info().
			Int("method", param.StatusCode).
			Str("method", param.Method).
			Str("path", param.Path).
			Str("client_ip", param.ClientIP).
			Dur("latency", param.Latency).
			Msg("server: request")

		return ""
	}))

	return router
}
