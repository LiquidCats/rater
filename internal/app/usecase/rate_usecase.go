package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/big"
	"rater/internal/app/domain/entity"
	"rater/internal/port"
	"time"
)

type RateUsecase struct {
	logger   port.Logger
	cache    port.CacheRepository
	adapters map[string]port.RateRepository
}

func NewRateUsecase(logger port.Logger, cache port.CacheRepository) *RateUsecase {
	return &RateUsecase{
		logger:   logger,
		cache:    cache,
		adapters: make(map[string]port.RateRepository),
	}
}

func (e *RateUsecase) GetRate(ctx context.Context, quote, base string) (*entity.Rate, error) {
	var price *big.Float
	var provider string

	key := fmt.Sprintf("rate:base:%s:quote:%s", base, quote)
	if e.cache.Has(ctx, key) {
		priceString, err := e.cache.Get(ctx, key)
		if nil != err {
			e.logger.Error("usecase: cant get rate value from cache", zap.Error(err))
		}

		if "" != priceString {
			f, _, err := big.ParseFloat(priceString, 10, 0, big.ToNearestEven)
			if nil != err {
				e.logger.Error("usecase: incorrect rate value from cache", zap.Error(err))
			}
			price = f
		}
	}

	if nil != price {
		return &entity.Rate{
			Base:  base,
			Quote: quote,
			Price: price,
		}, nil
	}

	for name, adapter := range e.adapters {
		p, err := adapter.Get(ctx, quote, base)
		if err != nil {
			e.logger.Error(
				"usecase: cant get rate",
				zap.Error(err),
				zap.String("quote", quote),
				zap.String("base", base),
				zap.String("provider", name),
			)
			continue
		}

		if p != nil {
			price = p
			provider = name
			break
		}
	}

	if nil == price {
		return nil, errors.New("usecase: exchange rate is not available right now")
	}

	if err := e.cache.Set(ctx, key, price.String(), 5*time.Minute); nil != err {
		e.logger.Error(
			"usecase: cant put rate value into cache",
			zap.Error(err),
			zap.String("quote", quote),
			zap.String("base", base),
			zap.String("provider", provider),
		)
	}

	return &entity.Rate{
		Base:     base,
		Quote:    quote,
		Price:    price,
		Provider: provider,
	}, nil
}

func (e *RateUsecase) SetAdapter(adapter port.NamedRateRepository) {
	_, ok := e.adapters[adapter.Name()]
	if !ok {
		e.adapters[adapter.Name()] = adapter.(port.RateRepository)
	}
}
