package service

import (
	"context"

	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/rotisserie/eris"
)

type CacheService struct {
	cache repository.RateCache
}

func NewCacheService(cache repository.RateCache) *CacheService {
	return &CacheService{
		cache: cache,
	}
}

func (s *CacheService) GetRate(ctx context.Context, pair entity.Pair) (*entity.Rate, error) {
	r, err := s.cache.GetRate(ctx, redis.RateKey{
		From: pair.From.String(),
		To:   pair.To.String(),
	})
	if err != nil {
		return nil, eris.Wrap(err, "get rate from cache")
	}
	if r.Price.IsZero() && r.Provider == "" && r.From == "" && r.To == "" {
		return nil, nil //nolint:nilnil
	}

	from := entity.CurrencyISO(r.From).ToUpper()
	to := entity.CurrencyISO(r.To).ToUpper()

	return &entity.Rate{
		Pair: entity.Pair{
			From:   from,
			To:     to,
			Symbol: entity.NewSymbol(from, to),
		},
		Price:    r.Price,
		Provider: entity.ProviderName(r.Provider),
	}, nil
}

func (s *CacheService) PutRate(ctx context.Context, rate entity.Rate) error {
	err := s.cache.PutRate(ctx, redis.RateKey{
		From: rate.Pair.From.String(),
		To:   rate.Pair.To.String(),
	}, redis.Rate{
		From:     rate.Pair.From.String(),
		To:       rate.Pair.To.String(),
		Price:    rate.Price,
		Provider: rate.Provider.String(),
	})

	if err != nil {
		return eris.Wrap(err, "put rate to cache")
	}

	return nil
}
