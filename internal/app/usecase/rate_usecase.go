package usecase

import (
	"context"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	domain "github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type RateUsecase struct {
	cache    repository.RateCache
	adapters map[entity.ProviderName]repository.RateAPI
}

func NewRateUsecase(cache repository.RateCache, providers api.Registry) *RateUsecase {
	return &RateUsecase{
		cache:    cache,
		adapters: providers,
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
		logger.Error().Stack().Err(err).Msg("cant get rate value from cache")
	}

	if rate != nil {
		return rate, nil
	}

	for name, adapter := range e.adapters {
		price, err = adapter.GetRate(ctx, pair)
		provider = name

		if err != nil {
			logger.Error().
				Err(err).
				Stack().
				Any("provider", provider).
				Msg("cant get rate")
			continue
		}

		if price.IsZero() {
			logger.Error().Any("provider", provider).Msg("provider gave zero value")
			continue
		} else {
			break
		}

	}

	if price.IsZero() {
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
