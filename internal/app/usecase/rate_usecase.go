package usecase

import (
	"context"
	"time"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/service"
	"github.com/LiquidCats/rater/internal/app/utils/timeutils"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type RateUseCaseDeps struct {
	Cache   service.CacheService
	Service service.RateService
}

type RateUseCase struct {
	cache   service.CacheService
	service service.RateService
}

func NewRateUseCase(deps RateUseCaseDeps) *RateUseCase {
	return &RateUseCase{
		cache:   deps.Cache,
		service: deps.Service,
	}
}

func (e *RateUseCase) CollectRate(ctx context.Context, symbol entity.Symbol) error {
	logger := zerolog.Ctx(ctx).
		With().
		Str("name", "use_case.get_rate").
		Any("symbol", symbol).
		Logger()

	rate, provider, err := e.service.Current(ctx, symbol.ToPair())
	if err != nil {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Any("provider", provider).
			Msg("cant get rate")

		return err
	}

	logger.Debug().
		Any("rate_entry", rate).
		Any("provider", provider).
		Msg("rate saved")

	return nil
}

func (e *RateUseCase) GetRate(ctx context.Context, symbol entity.Symbol, date time.Time) (*entity.Rate, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("name", "use_case.get_rate").
		Any("symbol", symbol).
		Logger()

	pairEntity := symbol.ToPair()
	logger.Debug().Msg("get rate")

	ts := timeutils.RoundToNearest(date, timeutils.FiveMinuteBucket)

	if ts.Compare(timeutils.RoundToNearest(time.Now(), timeutils.FiveMinuteBucket)) < 0 {
		rate, err := e.service.Historical(ctx, pairEntity, ts)
		if err != nil {
			logger.Error().
				Any("err", eris.ToJSON(err, true)).
				Msg("get historical rate")

			return nil, err
		}

		return rate, nil
	}

	rate, err := e.cache.GetRate(ctx, pairEntity)
	if err != nil {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Msg("cant get rate value from cache")
	}

	if rate != nil {
		return rate, nil
	}

	rate, provider, err := e.service.Current(ctx, pairEntity)
	if err != nil {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Any("provider", provider).
			Msg("cant get rate")
	}

	if err = e.cache.PutRate(ctx, *rate); nil != err {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Any("provider", provider).
			Msg("cant put rate value into cache")

		return nil, err
	}

	return rate, nil
}
