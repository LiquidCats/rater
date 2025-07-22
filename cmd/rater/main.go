package main

import (
	"context"
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/http"
	"github.com/LiquidCats/rater/internal/adapter/http/middlware"
	"github.com/LiquidCats/rater/internal/adapter/http/routes"
	"github.com/LiquidCats/rater/internal/adapter/metrics/prometheus"
	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
	"github.com/LiquidCats/rater/internal/adapter/scheduler/cron"
	"github.com/LiquidCats/rater/internal/app/bootstrap"
	"github.com/LiquidCats/rater/internal/app/service"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"

	_ "go.uber.org/automaxprocs"
)

func main() { //nolint:funlen
	logger := zerolog.New(os.Stdout).With().Caller().Stack().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger // nolint:reassign

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logger.WithContext(ctx)

	cfg, err := configs.Load()
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("failed to load config")
	}

	zerolog.SetGlobalLevel(cfg.App.LogLevel)

	poolConfig, err := pgxpool.ParseConfig(cfg.DB.ToDSN())
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("parse db config")
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("database config")
	}
	defer pool.Close()
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("connect to database")
	}

	migrationConn, err := pool.Acquire(ctx)
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("acquire pool connection")
	}

	if err = postgres.Migrate(migrationConn.Conn()); err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("migrate")
	}

	cache := redis.New(cfg.Redis)
	db := postgres.New(pool)

	apiRegistry, err := bootstrap.NewRegistry(ctx, cfg, db)
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("failed to create registry")
	}
	providerErrRateMetric := prometheus.NewProviderErrRate()
	responseTimeMetric := prometheus.NewResponseTime()

	cacheService := service.NewCacheService(cache)
	rateService := service.NewRateService(apiRegistry, db, providerErrRateMetric)

	rateUseCase := usecase.NewRateUseCase(usecase.RateUseCaseDeps{
		Cache:   cacheService,
		Service: rateService,
	})

	collectRateTask := cron.NewCollectRateTask(&logger, cfg.App, db, rateUseCase)

	rootHandler := routes.NewRootHandler()
	rateHandler := routes.NewRateHandler(rateUseCase, routes.Metrics{ResponseTime: responseTimeMetric})

	pairValidationMiddleware := middlware.NewPairValidation(db)
	dateValidationMiddleware := middlware.NewDateValidation()

	router := http.NewRouter()

	router.Any("/", rootHandler.Handle)

	v1Router := router.Group("/v1")
	v1Router.GET("/", rootHandler.Handle)
	v1Router.GET(
		"/rate/:pair/*date",
		pairValidationMiddleware.Handle,
		dateValidationMiddleware.Handle,
		rateHandler.Handle,
	)

	runners := []graceful.Runner{
		graceful.Signals,
		graceful.ScheduleRunner(
			collectRateTask,
		),
		graceful.ServerRunner(router, cfg.HTTP),
		graceful.ServerRunner(prometheus.GinHandler(), cfg.Metrics),
	}

	if err = graceful.WaitContext(
		ctx,
		runners...,
	); err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("server terminated")
	}

	logger.Info().Msg("application stopped")
}
