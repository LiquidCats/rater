package coinapi_test

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinapi"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/stretchr/testify/require"
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

	repo := coinapi.NewRepository(configs.CoinApiConfig{
		URL:    ts.URL,
		Secret: "test1234",
	})

	rate, err := repo.GetRate(context.Background(), entity.Pair{
		From: "BTC",
		To:   "USD",
	})
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, rate.String(), big.NewFloat(29295.929694597355).String())
}
