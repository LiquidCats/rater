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
	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/cex"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinapi"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coingate"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coingecko"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinmarketcap"
	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"

	_ "go.uber.org/automaxprocs"
)

const app = "rater"

func main() {
	logger := zerolog.New(os.Stdout).With().Caller().Stack().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logger // nolint:reassign

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logger.WithContext(ctx)

	cfg, err := configs.Load(app)
	if err != nil {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("failed to load config")
	}

	zerolog.SetGlobalLevel(cfg.App.LogLevel)

	cache, err := redis.NewCacheRepository(cfg.Redis, app)
	if nil != err {
		logger.Fatal().
			Any("err", eris.ToJSON(err, true)).
			Msg("app: cant connect to cache")
	}

	apiRegistry := api.Registry{}
	apiRegistry.Register(entity.ProviderNameCex, cex.NewRepository(cfg.Cex))
	apiRegistry.Register(entity.ProviderNameCoinApi, coinapi.NewRepository(cfg.CoinApi))
	apiRegistry.Register(entity.ProviderNameCoinGate, coingate.NewRepository(cfg.CoinGate))
	apiRegistry.Register(entity.ProviderNameCoinGecko, coingecko.NewRepository(cfg.CoinGecko))
	apiRegistry.Register(entity.ProviderNameCoinMarketCap, coinmarketcap.NewReposiotry(cfg.CoinMarketCap))

	providerErrRateMetric := prometheus.NewProviderErrRate(app)
	responseTimeMetric := prometheus.NewResponseTime(app)

	rateUsecase := usecase.NewRateUsecase(
		cache,
		apiRegistry,
		usecase.RateUsecaseMetrics{
			ProviderErrRate: providerErrRateMetric,
		},
	)

	rootHandler := routes.NewRootHandler()
	rateHandler := routes.NewRateHandler(rateUsecase, routes.Metrics{ResponseTime: responseTimeMetric})

	baseCurrencyMiddleware := middlware.NewPairValidation(cfg.App.Pairs)

	router := http.NewRouter()

	router.Any("/", rootHandler.Handle)

	v1Router := router.Group("/v1")
	v1Router.GET("/", rootHandler.Handle)
	v1Router.GET(
		"/rate/:pair",
		baseCurrencyMiddleware.Handle,
		rateHandler.Handle,
	)

	runners := []graceful.Runner{
		graceful.Signals,
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
