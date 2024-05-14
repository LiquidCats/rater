package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"rater/configs"
	"rater/internal/port"
	"time"
)

type Srv struct {
	http   *http.Server
	logger port.Logger
}

func NewServer(cfg configs.Config, router *gin.Engine, log port.Logger) *Srv {
	server := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Srv{
		http:   server,
		logger: log,
	}
}

func (s *Srv) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.logger.Info("server: starting server")
	if err := s.http.ListenAndServe(); nil != err && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("app: cant start server", zap.Error(err))
		cancel()
		return
	}
}

func (s *Srv) Stop(ctx context.Context) {
	s.logger.Info("server: shutting down server")
	if err := s.http.Shutdown(ctx); err != nil {
		log.Fatal("server: server shutdown failed", zap.Error(err))
	}
}
