package usecase

import (
	"context"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	domain "github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/LiquidCats/rater/internal/app/port/adapter/metrics"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type Metrics struct {
	providerErrRate metrics.ProviderErrRateMetric
}

type RateUsecase struct {
	cache    repository.RateCache
	adapters map[entity.ProviderName]repository.RateAPI
	metrics  Metrics
}

func NewRateUsecase(
	cache repository.RateCache,
	providers api.Registry,
	providerErrRateMetric metrics.ProviderErrRateMetric,
) *RateUsecase {
	return &RateUsecase{
		cache:    cache,
		adapters: providers,
		metrics: Metrics{
			providerErrRate: providerErrRateMetric,
		},
	}
}

func (e *RateUsecase) GetRate(ctx context.Context, pair entity.Pair) (*entity.Rate, error) {
	var (
		rate *entity.Rate

		price    decimal.Decimal
		provider entity.ProviderName
	)

	logger := zerolog.Ctx(ctx).
		With().
		Str("name", "use_case.get_rate").
		Any("pair", pair).
		Stack().
		Logger()

	logger.Debug().Msg("get rate")

	rate, err := e.cache.GetRate(ctx, pair)
	if err != nil {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Msg("cant get rate value from cache")
	}

	if rate != nil {
		return rate, nil
	}

	for name, adapter := range e.adapters {
		price, err = adapter.GetRate(ctx, pair)
		provider = name
		if err != nil {
			var providerErr *domain.ErrProviderRequestFailed
			if eris.As(err, &providerErr) {
				logger = logger.With().
					Int("provider_err_code", providerErr.StatusCode).
					Str("provider_err_body", providerErr.Body).
					Logger()
				e.metrics.providerErrRate.Inc(providerErr.StatusCode, name)
			}

			logger.Error().
				Any("err", eris.ToJSON(err, true)).
				Stack().
				Any("provider", name).
				Msg("cant get rate")
			continue
		}

		if !price.IsZero() {
			break
		}
	}

	if price.IsZero() {
		e.metrics.providerErrRate.Inc(-1, provider)

		logger.Error().
			Err(err).
			Stack().
			Any("provider", provider).
			Msg("rate not available")
		return nil, domain.ErrRateNotAvailable
	}

	rate = &entity.Rate{
		Pair:     pair,
		Price:    price,
		Provider: provider,
	}

	if err = e.cache.PutRate(ctx, *rate, 5*time.Second); nil != err { // nolint:mnd
		logger.Error().
			Err(err).
			Stack().
			Any("provider", provider).
			Msg("cant put rate value into cache")

		return nil, err
	}

	return rate, nil
}
