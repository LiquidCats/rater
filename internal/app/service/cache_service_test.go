package service_test

import (
	"testing"

	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/service"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCacheService_GetRate(t *testing.T) {
	ctx := t.Context()
	repo := mocks.NewRateCache(t)
	repo.On("GetRate", ctx, redis.RateKey{
		From: "BTC",
		To:   "USD",
	}).Once().Return(redis.Rate{
		From:     "BTC",
		To:       "USD",
		Price:    decimal.RequireFromString("65123.45"),
		Provider: "testprovider",
	}, nil)

	svc := service.NewCacheService(repo)

	rate, err := svc.GetRate(ctx, entity.Symbol("BTC_USD").ToPair())
	require.NoError(t, err)
	require.NotNil(t, rate)
	require.Equal(t, "BTC", string(rate.Pair.From))
	require.Equal(t, "USD", string(rate.Pair.To))
	require.Equal(t, "65123.45", rate.Price.String())
	require.Equal(t, "testprovider", string(rate.Provider))
}

func TestCacheService_PutRate(t *testing.T) {
	ctx := t.Context()
	repo := mocks.NewRateCache(t)

	repo.On("PutRate", ctx, redis.RateKey{
		From: "BTC",
		To:   "USD",
	}, redis.Rate{
		From:     "BTC",
		To:       "USD",
		Price:    decimal.RequireFromString("65123.45"),
		Provider: "testprovider",
	}).Once().Return(nil)

	svc := service.NewCacheService(repo)

	err := svc.PutRate(ctx, entity.Rate{
		Pair: entity.Pair{
			From:   "BTC",
			To:     "USD",
			Symbol: "BTC_USD",
		},
		Price:    decimal.RequireFromString("65123.45"),
		Provider: "testprovider",
	})
	require.NoError(t, err)
}
