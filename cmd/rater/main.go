package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"rater/configs"
	"rater/internal/adapter/api/middlware"
	"rater/internal/adapter/api/routes"
	"rater/internal/adapter/logger"
	"rater/internal/adapter/repository/api"
	"rater/internal/adapter/repository/redis"
	"rater/internal/app/usecase"
	"syscall"
	"time"
)

const app = "rater"

func main() {
	cfg := configs.Load()
	//
	log, err := logger.NewZapLogger(app)
	if nil != err {
		panic(fmt.Sprintf("app: cant configure logger - %s", err))
	}
	defer log.Sync()

	cache := redis.NewCacheRepository(cfg.Redis, app)

	rateUsecase := usecase.NewRateUsecase(log, cache)

	rateUsecase.SetAdapter(api.NewCoinApiRepository(cfg.CoinApiUrl, cfg.CoinApiSecret))
	rateUsecase.SetAdapter(api.NewCoinGateRepository(cfg.CoinGateUrl))

	rootHandler := routes.NewRootHandler()
	rateHandler := routes.NewRateHandler(rateUsecase)

	baseCurrencyMiddleware := middlware.NewCurrencyParamMiddleware("base", cfg.BaseCurrencies)
	quoteCurrencyMiddleware := middlware.NewCurrencyParamMiddleware("quote", cfg.QuoteCurrencies)

	router := gin.Default()

	router.Any("/", rootHandler.GetRoot)

	v1Router := router.Group("/v1")
	v1Router.GET("/", rootHandler.GetRoot)
	v1Router.GET(
		"/rate/:base/:quote",
		baseCurrencyMiddleware.Handle,
		quoteCurrencyMiddleware.Handle,
		rateHandler.GetRate,
	)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func(log *logger.ZapLogger, server *http.Server) {
		if err := server.ListenAndServe(); nil != err && !errors.Is(err, http.ErrServerClosed) {
			log.Error("cant start server", zap.Error(err))
		}
	}(log, server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create a context with a timeout for our server's graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Info("server shutdown failed", zap.Error(err))
	}

	log.Info("Done")
}
