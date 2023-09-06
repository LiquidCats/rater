package api

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestCoinGateRepository_Get(t *testing.T) {
	repo := &CoinGateRepository{
		url: "http://mocks:3001/coingate/v2/rates/merchant",
	}

	rate, err := repo.Get(context.Background(), "USD", "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, rate.String(), big.NewFloat(29295.92969).String())
}
