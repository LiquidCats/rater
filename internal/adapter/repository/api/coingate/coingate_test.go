package coingate

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

var response = `29295.929694597355`

func TestCoinGateRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := &Repository{
		url: ts.URL,
	}

	rate, err := repo.Get(context.Background(), "USD", "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	println(rate.String())

	require.Equal(t, rate.String(), big.NewFloat(29295.92969).String())
}
