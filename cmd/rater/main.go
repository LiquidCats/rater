package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"rater/configs"
	"rater/internal/adapter/api/middlware"
	"rater/internal/adapter/api/routes"
	"rater/internal/adapter/logger"
	"rater/internal/adapter/repository/api/cex"
	"rater/internal/adapter/repository/api/coinapi"
	"rater/internal/adapter/repository/api/coingate"
	"rater/internal/adapter/repository/cache/redis"
	httpserver "rater/internal/adapter/server"
	"rater/internal/app/usecase"
	"syscall"
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

	cache, err := redis.NewCacheRepository(cfg.Redis, app)
	if nil != err {
		log.Fatal("app: cant connect to cache", zap.Error(err))
	}

	rateUsecase := usecase.NewRateUsecase(log, cache)

	rateUsecase.SetAdapter(coinapi.NewRepository())
	rateUsecase.SetAdapter(coingate.NewRepository())
	rateUsecase.SetAdapter(cex.NewRepository())

	rootHandler := routes.NewRootHandler()
	rateHandler := routes.NewRateHandler(rateUsecase)

	baseCurrencyMiddleware := middlware.NewCurrencyParamMiddleware("base", cfg.BaseCurrencies)
	quoteCurrencyMiddleware := middlware.NewCurrencyParamMiddleware("quote", cfg.QuoteCurrencies)

	router := httpserver.NewRouter(log)

	router.Any("/", rootHandler.GetRoot)

	v1Router := router.Group("/v1")
	v1Router.GET("/", rootHandler.GetRoot)
	v1Router.GET(
		"/rate/:base/:quote",
		baseCurrencyMiddleware.Handle,
		quoteCurrencyMiddleware.Handle,
		rateHandler.GetRate,
	)

	server := httpserver.NewServer(cfg, router, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		server.Stop(ctx)
	case <-ctx.Done():
		log.Info("app: shutting down server...", zap.Error(ctx.Err()))
	}

	log.Info("Done")
}
