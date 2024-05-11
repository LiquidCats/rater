package coingecko

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var response = `{"bitcoin": {"eur": 60715.43364698}}`

func TestCoinGateRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := &Repository{
		url: ts.URL,
	}

	rate, err := repo.Get(context.Background(), "EUR", "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, rate.String(), big.NewFloat(60715.43364698).String())
}
