package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
	"github.com/LiquidCats/rater/internal/app/bootstrap"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	domain "github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/LiquidCats/rater/internal/app/port/adapter/metrics"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/LiquidCats/rater/internal/app/utils/timeutils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type RateService struct {
	adapters   bootstrap.Registry
	rateDB     repository.RateDatabase
	errMetrics metrics.ProviderErrRateMetric
}

func NewRateService(
	adapters bootstrap.Registry,
	rateDB repository.RateDatabase,
	errMetrics metrics.ProviderErrRateMetric,
) *RateService {
	return &RateService{
		adapters:   adapters,
		rateDB:     rateDB,
		errMetrics: errMetrics,
	}
}

func (s *RateService) Historical(ctx context.Context, pair entity.Pair, ts time.Time) (*entity.Rate, error) {
	logger := zerolog.Ctx(ctx).
		With().
		Str("name", "service.rate.historical").
		Any("pair", pair).
		Logger()

	historical, err := s.rateDB.GetRate(ctx, postgres.GetRateParams{
		Ts: pgtype.Timestamp{
			Time:  ts,
			Valid: true,
		},
		Pair: pair.Symbol.String(),
	})
	if err != nil {
		if eris.Is(err, sql.ErrNoRows) {
			logger.Error().
				Any("err", eris.ToJSON(err, true)).
				Msg("no historical rate found")
			return nil, domain.ErrNoHistoricalRate
		}
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Msg("get historical rate")
		return nil, eris.Wrap(err, "cant get rate from database")
	}

	logger.Debug().
		Any("rate", historical).
		Msg("get historical rate")

	rate := entity.NewRate(historical)

	return &rate, nil
}

func (s *RateService) Current(ctx context.Context, pair entity.Pair) (*entity.Rate, entity.ProviderName, error) {
	var (
		err      error
		price    decimal.Decimal
		provider entity.ProviderName
	)
	start := time.Now()

	logger := zerolog.Ctx(ctx).
		With().
		Str("name", "service.rate.current").
		Any("pair", pair).
		Logger()

	for name, adapter := range s.adapters {
		price, err = adapter.GetRate(ctx, pair)
		provider = name
		if err != nil {
			var providerErr *domain.ProviderRequestFailedError
			if eris.As(err, &providerErr) {
				logger = logger.With().
					Int("provider_err_code", providerErr.StatusCode).
					Str("provider_err_body", providerErr.Body).
					Logger()
				s.errMetrics.Inc(providerErr.StatusCode, name)
			}

			logger.Error().
				Any("err", eris.ToJSON(err, true)).
				Any("provider", name).
				Msg("cant get rate")
			continue
		}

		if !price.IsZero() {
			break
		}
	}

	if price.IsZero() {
		s.errMetrics.Inc(-1, provider)

		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Any("provider", provider).
			Msg("rate not available")
		return nil, provider, domain.ErrRateNotAvailable
	}

	rate := &entity.Rate{
		Pair:     pair,
		Price:    price,
		Provider: provider,
	}

	logger.Debug().
		Any("rate", rate).
		Msg("get historical rate")

	hasRate, err := s.rateDB.HasRate(ctx, postgres.HasRateParams{
		Ts: pgtype.Timestamp{
			Time:  timeutils.RoundToNearest(start, timeutils.FiveMinuteBucket),
			Valid: true,
		},
		Pair: rate.Provider.String(),
	})

	if err == nil && hasRate {
		return rate, provider, nil
	}

	if err != nil {
		logger.Warn().
			Any("err", eris.ToJSON(err, true)).
			Msg("cant check rate in database")
		return rate, provider, nil
	}

	_, err = s.rateDB.SaveRate(ctx, postgres.SaveRateParams{
		Price:    rate.Price,
		Pair:     rate.Pair.Symbol.String(),
		Provider: rate.Provider.String(),
		Ts: pgtype.Timestamp{
			Time:  timeutils.RoundToNearest(start, timeutils.FiveMinuteBucket),
			Valid: true,
		},
	})

	if err != nil {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Any("provider", provider).
			Any("rate", rate).
			Msg("cant save rate to database")
	}

	return rate, provider, nil
}
