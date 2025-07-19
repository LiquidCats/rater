package service_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
	"github.com/LiquidCats/rater/internal/app/bootstrap"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/service"
	"github.com/LiquidCats/rater/internal/app/utils/timeutils"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateService_Current(t *testing.T) {
	ctx := t.Context()

	rateAPI := mocks.NewRateAPI(t)
	rateDB := mocks.NewRateDatabase(t)
	metr := mocks.NewProviderErrRateMetric(t)

	reg := make(bootstrap.Registry, 1)
	reg["testprovider"] = rateAPI

	pair := entity.Pair{From: "BTC", To: "USD", Symbol: "BTC_USD"}
	price := decimal.RequireFromString("65256.31")

	rateAPI.On("GetRate", ctx, pair).Once().Return(price, nil)
	rateDB.On("HasRate", ctx, postgres.HasRateParams{
		Ts: pgtype.Timestamp{
			Time:  timeutils.RoundToNearest(time.Now(), timeutils.FiveMinuteBucket),
			Valid: true,
		},
		Pair: "testprovider",
	}).Once().Return(false, nil)
	rateDB.On("SaveRate", ctx, postgres.SaveRateParams{
		Price: pgtype.Numeric{
			Int:   price.Coefficient(),
			Exp:   price.Exponent(),
			Valid: true,
		},
		Pair:     pair.Symbol.String(),
		Provider: "testprovider",
		Ts: pgtype.Timestamp{
			Time:  timeutils.RoundToNearest(time.Now(), timeutils.FiveMinuteBucket),
			Valid: true,
		},
	}).Once().Return(postgres.Rate{
		ID: 12,
		Price: pgtype.Numeric{
			Int:   price.Coefficient(),
			Exp:   price.Exponent(),
			Valid: true,
		},
		Pair:     pair.Symbol.String(),
		Provider: "testprovider",
		Ts: pgtype.Timestamp{
			Time: timeutils.RoundToNearest(time.Now(), timeutils.FiveMinuteBucket),
		},
		CreatedAt: pgtype.Timestamp{
			Time: time.Now(),
		},
	}, nil)

	svc := service.NewRateService(reg, rateDB, metr)

	rate, provider, err := svc.Current(ctx, pair)
	require.NoError(t, err)

	assert.Equal(t, "testprovider", provider.String())
	assert.Equal(t, pair.From.String(), rate.Pair.From.String())
	assert.Equal(t, pair.To.String(), rate.Pair.To.String())
	assert.Equal(t, pair.Symbol.String(), rate.Pair.Symbol.String())
	assert.Equal(t, price.String(), rate.Price.String())
	assert.Equal(t, "testprovider", rate.Provider.String())
}

func TestRateService_Historical(t *testing.T) {
	ctx := t.Context()
	now := time.Now()
	pair := entity.Pair{From: "BTC", To: "USD", Symbol: "BTC_USD"}
	price := decimal.RequireFromString("65256.31")

	rateAPI := mocks.NewRateAPI(t)
	rateDB := mocks.NewRateDatabase(t)
	metr := mocks.NewProviderErrRateMetric(t)

	rateDB.On("GetRate", ctx, postgres.GetRateParams{
		Ts: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
		Pair: pair.Symbol.String(),
	}).Once().Return(postgres.Rate{
		ID: 12,
		Price: pgtype.Numeric{
			Int:   price.Coefficient(),
			Exp:   price.Exponent(),
			Valid: true,
		},
		Pair:     pair.Symbol.String(),
		Provider: "testprovider",
		Ts: pgtype.Timestamp{
			Time: timeutils.RoundToNearest(time.Now(), timeutils.FiveMinuteBucket),
		},
		CreatedAt: pgtype.Timestamp{
			Time: time.Now(),
		},
	}, nil)

	reg := make(bootstrap.Registry, 1)
	reg["testprovider"] = rateAPI

	svc := service.NewRateService(reg, rateDB, metr)

	rate, err := svc.Historical(ctx, pair, now)
	require.NoError(t, err)
	require.NotNil(t, rate)

	assert.Equal(t, "testprovider", rate.Provider.String())
	assert.Equal(t, pair.From.String(), rate.Pair.From.String())
	assert.Equal(t, pair.To.String(), rate.Pair.To.String())
	assert.Equal(t, pair.Symbol.String(), rate.Pair.Symbol.String())
	assert.Equal(t, price.String(), rate.Price.String())
}
