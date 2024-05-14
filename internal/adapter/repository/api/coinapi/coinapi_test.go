package coinapi

import (
	"context"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

var response = `{"time":"2023-07-24T11:31:56.0000000Z","asset_id_base":"BTC","asset_id_quote":"USD","rate":29295.929694597355}`

func TestCoinApiRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "test1234", r.Header.Get("X-CoinAPI-Key"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := &Repository{
		url:      ts.URL,
		apiToken: "test1234",
	}

	rate, err := repo.Get(context.Background(), "USD", "BTC")
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, rate.String(), big.NewFloat(29295.929694597355).String())
}
