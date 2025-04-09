package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/LiquidCats/rater/configs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Srv struct {
	http *http.Server
}

func NewServer(cfg configs.AppConfig, router *gin.Engine) *Srv {
	server := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second, // nolint:mnd
		WriteTimeout:   10 * time.Second, // nolint:mnd
		MaxHeaderBytes: 1 << 20,          // nolint:mnd
	}

	return &Srv{
		http: server,
	}
}

func (s *Srv) Start(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("server: starting server")

	if err := s.http.ListenAndServe(); nil != err && !errors.Is(err, http.ErrServerClosed) {
		logger.Error().Err(err).Msg("app: cant start server")
	}
}

func (s *Srv) Stop(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("server: stopping server")

	if err := s.http.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server: server shutdown failed")
	}
}
