package usecase

import (
	"context"
	"math/big"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	domain "github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/rs/zerolog"
)

type RateUsecase struct {
	cache    repository.RateCache
	adapters map[entity.ProviderName]repository.RateApi
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

		price    big.Float
		provider entity.ProviderName
	)

	logger := zerolog.Ctx(ctx).
		With().
		Str("name", "user.get_rate").
		Stack().
		Logger()

	rate, err := e.cache.GetRate(ctx, pair)
	if err != nil {
		logger.Error().Err(err).Msg("cant get rate value from cache")
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
				Any("pair", pair).
				Any("provider", provider).
				Msg("usecase: cant get rate")
			continue
		}

		if price.Cmp(big.NewFloat(0)) == 0 {
			continue
		}
	}

	if price.Cmp(big.NewFloat(0)) == 0 {
		return nil, domain.ErrRateNotAvailable
	}

	rate = &entity.Rate{
		Pair:     pair,
		Price:    price,
		Provider: provider,
	}

	if err := e.cache.PutRate(ctx, *rate, 5*time.Minute); nil != err {
		logger.Error().
			Err(err).
			Any("pair", pair).
			Any("provider", provider).
			Msg("usecase: cant put rate value into cache")
	}

	return rate, nil
}
