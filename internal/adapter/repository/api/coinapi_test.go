package api

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestCoinApiRepository_Get(t *testing.T) {
	repo := &CoinApiRepository{
		url:      "http://mocks:3001/coinapi/v1/exchangerate",
		apiToken: "test1234",
	}

	rate, err := repo.Get(context.Background(), "USD", "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, rate.String(), big.NewFloat(29295.929694597355).String())
}
