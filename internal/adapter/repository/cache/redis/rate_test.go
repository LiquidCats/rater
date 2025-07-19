package redis_test

import (
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetRate(t *testing.T) {
	srv := miniredis.RunT(t)
	defer srv.Close()

	srv.FlushAll()

	err := srv.Set("rater:rate:from:btc:to:usd", `{"from":"BTC","to":"USD","price":"112000.34","provider":"testprovider"}`)
	require.NoError(t, err)

	cfg := configs.RedisConfig{
		Address:  srv.Addr(),
		Protocol: 2,
	}

	repo := redis.New(cfg)

	ctx := t.Context()

	rate, err := repo.GetRate(ctx, redis.RateKey{
		From: "BTC",
		To:   "USD",
	})

	require.NoError(t, err)
	require.NotNil(t, rate)
	require.Equal(t, "112000.34", rate.Price.String())
	require.Equal(t, "BTC", rate.From)
	require.Equal(t, "USD", rate.To)
	require.Equal(t, "testprovider", rate.Provider)
}

func TestRepository_PutRate(t *testing.T) {
	srv := miniredis.RunT(t)
	defer srv.Close()

	srv.FlushAll()

	err := srv.Set("rater:rate:from:btc:to:usd", `{"from":"BTC","to":"USD","price":"112000.34","provider":"testprovider"}`)
	require.NoError(t, err)

	cfg := configs.RedisConfig{
		Address:  srv.Addr(),
		Protocol: 2,
	}

	repo := redis.New(cfg)

	ctx := t.Context()

	err = repo.PutRate(ctx, redis.RateKey{
		From: "BTC",
		To:   "USD",
	}, redis.Rate{
		From:     "BTC",
		To:       "USD",
		Price:    decimal.RequireFromString("112011.34"),
		Provider: "testprovider2",
	})
	require.NoError(t, err)

	res, err := srv.Get("rater:rate:from:btc:to:usd")
	require.NoError(t, err)
	require.Equal(t, `{"from":"BTC","to":"USD","price":"112011.34","provider":"testprovider2"}`, res)
}
