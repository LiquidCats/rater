package coinapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinapi"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const response = `{"time":"2023-07-24T11:31:56.0000000Z","asset_id_base":"BTC","asset_id_quote":"USD","rate":29295.929694597355}`

func TestCoinApiRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test1234", r.Header.Get("X-CoinAPI-Key")) //nolint:canonicalheader

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := coinapi.NewRepository(configs.CoinApiConfig{
		URL:    ts.URL,
		Secret: "test1234",
	})

	rate, err := repo.GetRate(t.Context(), entity.Pair{
		From: "BTC",
		To:   "USD",
	})
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, "29295.929694597355", rate.String())
}
