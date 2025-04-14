package coinmarketcap_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinmarketcap"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const response = `{"data":{"symbol":"BTC","id":"1","name":"Bitcoin","amount":50,"last_updated":"2018-06-06T08:04:36.000Z","quote":{"GBP":{"price":284656.08465608465,"last_updated":"2018-06-06T06:00:00.000Z"},"LTC":{"price":3128.7279766396537,"last_updated":"2018-06-06T08:04:02.000Z"},"USD":{"price":381442,"last_updated":"2018-06-06T08:06:51.968Z"}}},"status":{"timestamp":"2024-05-08T04:21:23.526Z","error_code":0,"error_message":"","elapsed":10,"credit_count":1,"notice":""}}`

func TestCoinmarketcapRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := coinmarketcap.NewReposiotry(configs.CoinMarketCapConfig{
		URL:    ts.URL,
		Secret: "test",
	})

	rate, err := repo.GetRate(t.Context(), entity.Pair{
		From: "BTC",
		To:   "USD",
	})
	require.NoError(t, err)
	require.NotNil(t, rate)

	assert.Equal(t, "381442", rate.String())
}
