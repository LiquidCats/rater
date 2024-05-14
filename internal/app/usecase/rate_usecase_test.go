package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"rater/internal/adapter/logger"
	"rater/internal/adapter/repository/cache/memory"
	"rater/internal/app/domain/types"
	"testing"
)

type TestAdapter struct {
}

func (a *TestAdapter) Get(_ context.Context, _ types.QuoteCurrency, _ types.BaseCurrency) (*big.Float, error) {
	return big.NewFloat(25000.77733333), nil
}

func (a *TestAdapter) Name() string {
	return "test"
}

func TestExchange_Get(t *testing.T) {
	l := logger.NewNilLogger()
	c := memory.NewCacheRepository()

	exch := NewRateUsecase(l, c)
	exch.SetAdapter(&TestAdapter{})

	rate, err := exch.GetRate(context.Background(), "USD", "BTC")

	require.NoError(t, err)
	require.NotNil(t, rate)

	assert.Equal(t, "BTC", string(rate.Base))
	assert.Equal(t, "USD", string(rate.Quote))
	assert.Equal(t, big.NewFloat(25000.77733333).String(), rate.Price.String())
}
